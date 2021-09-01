/*
 * Copyright (C) 2021 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */
package eth

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	scom "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/eth/rlp"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/zhiqiangxu/go-ethereum/common"
	"github.com/zhiqiangxu/go-ethereum/consensus"
	"github.com/zhiqiangxu/go-ethereum/crypto"
	"github.com/zhiqiangxu/go-ethereum/params"
	"golang.org/x/crypto/sha3"
)

var (
	BIG_1             = big.NewInt(1)
	BIG_2             = big.NewInt(2)
	BIG_9             = big.NewInt(9)
	BIG_MINUS_99      = big.NewInt(-99)
	BLOCK_DIFF_FACTOR = big.NewInt(2048)
	DIFF_PERIOD       = big.NewInt(100000)
	BOMB_DELAY        = big.NewInt(8999999)
)

type ETHHandler struct {
}

func NewETHHandler() *ETHHandler {
	return &ETHHandler{}
}

func (this *ETHHandler) SyncGenesisHeader(native *native.NativeContract) error {
	ctx := native.ContractRef().CurrentContext()
	params := &scom.SyncGenesisHeaderParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodSyncGenesisHeader, params, ctx.Payload); err != nil {
		return fmt.Errorf("SyncGenesisHeader, contract params deserialize error: %v", err)
	}

	// Get current epoch operator
	ok, err := node_manager.CheckConsensusSigns(native, scom.MethodSyncGenesisHeader, ctx.Payload, native.ContractRef().MsgSender())
	if err != nil {
		return fmt.Errorf("SyncGenesisHeader, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return nil
	}

	var header Header
	err = json.Unmarshal(params.GenesisHeader, &header)
	if err != nil {
		return fmt.Errorf("SyncGenesisHeader, json.Unmarshal header err: %v", err)
	}

	headerStore, err := native.GetCacheDB().Get(utils.ConcatKey(utils.HeaderSyncContractAddress, []byte(scom.GENESIS_HEADER), utils.GetUint64Bytes(params.ChainID)))
	if err != nil {
		return fmt.Errorf("ETHHandler GetHeaderByHeight, get blockHashStore error: %v", err)
	}
	if headerStore != nil {
		return fmt.Errorf("ETHHandler GetHeaderByHeight, genesis header had been initialized")
	}

	//block header storage
	err = putGenesisBlockHeader(native, header, params.ChainID)
	if err != nil {
		return fmt.Errorf("ETHHandler SyncGenesisHeader, put blockHeader error: %v", err)
	}

	return nil
}

func (this *ETHHandler) SyncBlockHeader(native *native.NativeContract) error {
	headerParams := &scom.SyncBlockHeaderParam{}
	{
		ctx := native.ContractRef().CurrentContext()
		if err := utils.UnpackMethod(scom.ABI, scom.MethodSyncBlockHeader, headerParams, ctx.Payload); err != nil {
			return err
		}
	}
	caches := NewCaches(3, native)
	for _, v := range headerParams.Headers {
		var header Header
		err := json.Unmarshal(v, &header)
		if err != nil {
			return fmt.Errorf("SyncBlockHeader, deserialize header err: %v", err)
		}
		headerHash := header.Hash()
		exist, err := IsHeaderExist(native, headerHash.Bytes(), headerParams.ChainID)
		if err != nil {
			return fmt.Errorf("SyncBlockHeader, check header exist err: %v", err)
		}
		if exist == true {
			log.Warnf("SyncBlockHeader, header has exist. Header: %s", string(v))
			continue
		}
		// get pre header
		parentHeader, parentDifficultySum, err := GetHeaderByHash(native, header.ParentHash.Bytes(), headerParams.ChainID)
		if err != nil {
			return fmt.Errorf("SyncBlockHeader, get the parent block failed. Error:%s, header: %s", err, string(v))
		}
		parentHeaderHash := parentHeader.Hash()
		/**
		this code source refer to https://github.com/ethereum/go-ethereum/blob/master/consensus/ethash/consensus.go
		verify header need to verify:
		1. parent hash
		2. extra size
		3. current time
		*/
		//verify whether parent hash validity
		if !bytes.Equal(parentHeaderHash.Bytes(), header.ParentHash.Bytes()) {
			return fmt.Errorf("SyncBlockHeader, parent header is not right. Header: %s", string(v))
		}
		//verify whether extra size validity
		if uint64(len(header.Extra)) > params.MaximumExtraDataSize {
			return fmt.Errorf("SyncBlockHeader, SyncBlockHeader extra-data too long: %d > %d, header: %s", len(header.Extra), params.MaximumExtraDataSize, string(v))
		}
		//verify current time validity
		if header.Time > uint64(time.Now().Add(allowedFutureBlockTime).Unix()) {
			return fmt.Errorf("SyncBlockHeader,  verify header time error:%s, checktime: %d, header: %s", consensus.ErrFutureBlock, time.Now().Add(allowedFutureBlockTime).Unix(), string(v))
		}
		//verify whether current header time and prevent header time validity
		if header.Time <= parentHeader.Time {
			return fmt.Errorf("SyncBlockHeader, verify header time fail. Header: %s", string(v))
		}
		// Verify that the gas limit is <= 2^63-1
		cap := uint64(0x7fffffffffffffff)
		if header.GasLimit > cap {
			return fmt.Errorf("SyncBlockHeader, invalid gasLimit: have %v, max %v, header: %s", header.GasLimit, cap, string(v))
		}
		// Verify that the gasUsed is <= gasLimit
		if header.GasUsed > header.GasLimit {
			return fmt.Errorf("SyncBlockHeader, invalid gasUsed: have %d, gasLimit %d, header: %s", header.GasUsed, header.GasLimit, string(v))
		}
		if isLondon(&header) {
			err = VerifyEip1559Header(parentHeader, &header)
		} else {
			err = VerifyGaslimit(parentHeader.GasLimit, header.GasLimit)
		}
		if err != nil {
			return fmt.Errorf("SyncBlockHeader, err:%v", err)
		}

		//verify difficulty
		var expected *big.Int
		if isLondon(&header) {
			expected = makeDifficultyCalculator(big.NewInt(9700000))(header.Time, parentHeader)
		} else {
			return fmt.Errorf("SyncBlockHeader, header before london fork is no longer supported")
		}
		if expected.Cmp(header.Difficulty) != 0 {
			return fmt.Errorf("SyncBlockHeader, invalid difficulty: have %v, want %v, header: %s", header.Difficulty, expected, string(v))
		}
		// verfify header
		err = this.verifyHeader(&header, caches)
		if err != nil {
			return fmt.Errorf("SyncBlockHeader, verify header error: %v, header: %s", err, string(v))
		}
		//block header storage
		hederDifficultySum := new(big.Int).Add(header.Difficulty, parentDifficultySum)
		err = putBlockHeader(native, header, hederDifficultySum, headerParams.ChainID)
		if err != nil {
			return fmt.Errorf("SyncGenesisHeader, put blockHeader error: %v, header: %s", err, string(v))
		}
		// get current header of main
		currentHeader, currentDifficultySum, err := GetCurrentHeader(native, headerParams.ChainID)
		if err != nil {
			return fmt.Errorf("SyncBlockHeader, get the current block failed. error:%s", err)
		}
		if bytes.Equal(currentHeader.Hash().Bytes(), header.ParentHash.Bytes()) {
			appendHeader2Main(native, header.Number.Uint64(), headerHash, headerParams.ChainID)
		} else {
			//
			if hederDifficultySum.Cmp(currentDifficultySum) > 0 {
				RestructChain(native, currentHeader, &header, headerParams.ChainID)
			}
		}
	}
	caches.deleteCaches()
	return nil
}

func (this *ETHHandler) SyncCrossChainMsg(native *native.NativeContract) error {
	return nil
}

func (this *ETHHandler) verifyHeader(header *Header, caches *Caches) error {
	// try to verfify header
	number := header.Number.Uint64()
	size := datasetSize(number)
	headerHash := HashHeader(header).Bytes()
	nonce := header.Nonce.Uint64()
	// get seed and seed head
	seed := make([]byte, 40)
	copy(seed, headerHash)
	binary.LittleEndian.PutUint64(seed[32:], nonce)
	seed = crypto.Keccak512(seed)
	// get mix
	mix := make([]uint32, mixBytes/4)
	for i := 0; i < len(mix); i++ {
		mix[i] = binary.LittleEndian.Uint32(seed[i%16*4:])
	}
	// get cache
	cache := caches.getCache(number)
	if len(cache) <= 1 {
		return fmt.Errorf("cache of proof-of-work is not generated!")
	}
	// get new mix with DAG data
	rows := uint32(size / mixBytes)
	temp := make([]uint32, len(mix))
	seedHead := binary.LittleEndian.Uint32(seed)
	for i := 0; i < loopAccesses; i++ {
		parent := fnv(uint32(i)^seedHead, mix[i%len(mix)]) % rows
		for j := uint32(0); j < mixBytes/hashBytes; j++ {
			xx := lookup(cache, 2*parent+j)
			copy(temp[j*hashWords:], xx)
		}
		fnvHash(mix, temp)
	}
	// get new mix by compress
	for i := 0; i < len(mix); i += 4 {
		mix[i/4] = fnv(fnv(fnv(mix[i], mix[i+1]), mix[i+2]), mix[i+3])
	}
	mix = mix[:len(mix)/4]
	// get digest by compressed mix
	digest := make([]byte, common.HashLength)
	for i, val := range mix {
		binary.LittleEndian.PutUint32(digest[i*4:], val)
	}
	// get header result hash
	result := crypto.Keccak256(append(seed, digest...))
	// Verify the calculated digest against the ones provided in the header
	if !bytes.Equal(header.MixDigest[:], digest) {
		return fmt.Errorf("invalid mix digest!")
	}
	// compare result hash with target hash
	target := new(big.Int).Div(two256, header.Difficulty)
	if new(big.Int).SetBytes(result).Cmp(target) > 0 {
		return fmt.Errorf("invalid proof-of-work!")
	}
	return nil
}

func HashHeader(header *Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	enc := []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra,
	}
	if header.BaseFee != nil {
		enc = append(enc, header.BaseFee)
	}
	rlp.Encode(hasher, enc)
	hasher.Sum(hash[:0])
	return hash
}

type hasher func(dest []byte, data []byte)

func makeHasher(h hash.Hash) hasher {
	// sha3.state supports Read to get the sum, use it to avoid the overhead of Sum.
	// Read alters the state but we reset the hash before every operation.
	type readerHash interface {
		hash.Hash
		Read([]byte) (int, error)
	}
	rh, ok := h.(readerHash)
	if !ok {
		panic("can't find Read method on hash")
	}
	outputLen := rh.Size()
	return func(dest []byte, data []byte) {
		rh.Reset()
		rh.Write(data)
		rh.Read(dest[:outputLen])
	}
}

func seedHash(block uint64) []byte {
	seed := make([]byte, 32)
	if block < epochLength {
		return seed
	}
	keccak256 := makeHasher(sha3.NewLegacyKeccak256())
	for i := 0; i < int(block/epochLength); i++ {
		keccak256(seed, seed)
	}
	return seed
}
