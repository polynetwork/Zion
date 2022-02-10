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
package bsc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	scom "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/eth/types"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

// Handler ...
type Handler struct {
}

// NewHandler ...
func NewHandler() *Handler {
	return &Handler{}
}

// GenesisHeader ...
type GenesisHeader struct {
	Header         types.Header
	PrevValidators []HeightAndValidators
}

// HeightAndValidators ...
type HeightAndValidators struct {
	Height     *big.Int
	Validators []common.Address
	Hash       *common.Hash
}

// HeaderWithDifficultySum ...
type HeaderWithDifficultySum struct {
	Header          *types.Header `json:"header"`
	DifficultySum   *big.Int      `json:"difficultySum"`
	EpochParentHash *common.Hash  `json:"epochParentHash"`
}

// ExtraInfo ...
type ExtraInfo struct {
	ChainID *big.Int // for bsc
}

// Context ...
type Context struct {
	ExtraInfo ExtraInfo
	ChainID   uint64
}

// HeaderWithChainID ...
type HeaderWithChainID struct {
	Header  *HeaderWithDifficultySum
	ChainID uint64
}

var (
	inMemoryHeaders = 400
	inMemoryGenesis = 40
	extraVanity     = 32                       // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal       = crypto.SignatureLength   // Fixed number of extra-data suffix bytes reserved for signer seal
	uncleHash       = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
	diffInTurn      = big.NewInt(2)            // Block difficulty for in-turn signatures
	diffNoTurn      = big.NewInt(1)            // Block difficulty for out-of-turn signatures

	GasLimitBoundDivisor uint64 = 256 // The bound divisor of the gas limit, used in update calculations.
)

// SyncGenesisHeader ...
func (h *Handler) SyncGenesisHeader(native *native.NativeContract) (err error) {
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

	// can only store once
	stored, err := isGenesisStored(native, params)
	if err != nil {
		return fmt.Errorf("bsc Handler SyncGenesisHeader, isGenesisStored error: %v", err)
	}
	if stored {
		return fmt.Errorf("bsc Handler SyncGenesisHeader, genesis had been initialized")
	}

	var genesis GenesisHeader
	err = json.Unmarshal(params.GenesisHeader, &genesis)
	if err != nil {
		return fmt.Errorf("bsc Handler SyncGenesisHeader, deserialize GenesisHeader err: %v", err)
	}

	signersBytes := len(genesis.Header.Extra) - extraVanity - extraSeal
	if signersBytes == 0 || signersBytes%common.AddressLength != 0 {
		return fmt.Errorf("invalid signer list, signersBytes:%d", signersBytes)
	}

	if len(genesis.PrevValidators) != 1 {
		return fmt.Errorf("invalid PrevValidators")
	}
	if genesis.Header.Number.Cmp(genesis.PrevValidators[0].Height) <= 0 {
		return fmt.Errorf("invalid height orders")
	}
	validators, err := ParseValidators(genesis.Header.Extra[extraVanity : extraVanity+signersBytes])
	if err != nil {
		return
	}
	genesis.PrevValidators = append([]HeightAndValidators{
		{Height: genesis.Header.Number, Validators: validators},
	}, genesis.PrevValidators...)

	err = storeGenesis(native, params, &genesis)
	if err != nil {
		return fmt.Errorf("bsc Handler SyncGenesisHeader, storeGenesis error: %v", err)
	}

	return
}

// SyncBlockHeader ...
func (h *Handler) SyncBlockHeader(native *native.NativeContract) error {
	headerParams := &scom.SyncBlockHeaderParam{}
	{
		ctx := native.ContractRef().CurrentContext()
		if err := utils.UnpackMethod(scom.ABI, scom.MethodSyncBlockHeader, headerParams, ctx.Payload); err != nil {
			return err
		}
	}

	side, err := side_chain_manager.GetSideChain(native, headerParams.ChainID)
	if err != nil {
		return fmt.Errorf("bsc Handler SyncBlockHeader, GetSideChain error: %v", err)
	}
	var extraInfo ExtraInfo
	err = json.Unmarshal(side.ExtraInfo, &extraInfo)
	if err != nil {
		return fmt.Errorf("bsc Handler SyncBlockHeader, ExtraInfo Unmarshal error: %v", err)
	}

	ctx := &Context{ExtraInfo: extraInfo, ChainID: headerParams.ChainID}

	for _, v := range headerParams.Headers {
		var header types.Header
		err := json.Unmarshal(v, &header)
		if err != nil {
			return fmt.Errorf("bsc Handler SyncBlockHeader, deserialize header err: %v", err)
		}
		headerHash := header.Hash()

		exist, err := isHeaderExist(native, headerHash, ctx)
		if err != nil {
			return fmt.Errorf("bsc Handler SyncBlockHeader, isHeaderExist headerHash err: %v", err)
		}
		if exist {
			log.Warnf("bsc Handler SyncBlockHeader, header has exist. Header: %s", string(v))
			continue
		}

		parentExist, err := isHeaderExist(native, header.ParentHash, ctx)
		if err != nil {
			return fmt.Errorf("bsc Handler SyncBlockHeader, isHeaderExist ParentHash err: %v", err)
		}
		if !parentExist {
			return fmt.Errorf("bsc Handler SyncBlockHeader, parent header not exist. Header: %s", string(v))
		}

		signer, err := verifySignature(native, &header, ctx)
		if err != nil {
			return fmt.Errorf("bsc Handler SyncBlockHeader, verifySignature err: %v", err)
		}

		// get prev epochs, also checking recent limit
		phv, pphv, lastSeenHeight, err := getPrevHeightAndValidators(native, &header, ctx)
		if err != nil {
			return fmt.Errorf("bsc Handler SyncBlockHeader, getPrevHeightAndValidators err: %v", err)
		}

		var (
			inTurnHV *HeightAndValidators
		)

		diffWithLastEpoch := big.NewInt(0).Sub(header.Number, phv.Height).Int64()
		if diffWithLastEpoch <= int64(len(pphv.Validators)/2) {
			// pphv is in effect
			inTurnHV = pphv

			if len(header.Extra) > extraVanity+extraSeal {
				return fmt.Errorf("bsc Handler SyncBlockHeader: can not change epoch continuously")
			}
		} else {
			// phv is in effect
			inTurnHV = phv
		}

		if lastSeenHeight > 0 {
			limit := int64(len(inTurnHV.Validators) / 2)
			if header.Number.Int64() <= lastSeenHeight+limit {
				return fmt.Errorf("bsc Handler SyncBlockHeader, RecentlySigned, lastSeenHeight:%d currentHeight:%d #V:%d", lastSeenHeight, header.Number.Int64(), len(inTurnHV.Validators))
			}
		}

		indexInTurn := int(header.Number.Uint64()) % len(inTurnHV.Validators)
		if indexInTurn < 0 {
			return fmt.Errorf("indexInTurn is negative:%d inTurnHV.Height:%d header.Number:%d", indexInTurn, inTurnHV.Height.Int64(), header.Number.Int64())
		}
		valid := false
		for idx, v := range inTurnHV.Validators {
			if v == signer {
				valid = true
				if indexInTurn == idx {
					if header.Difficulty.Cmp(diffInTurn) != 0 {
						return fmt.Errorf("invalid difficulty, got %v expect %v index:%v", header.Difficulty.Int64(), diffInTurn.Int64(), int(indexInTurn)%len(inTurnHV.Validators))
					}
				} else {
					if header.Difficulty.Cmp(diffNoTurn) != 0 {
						return fmt.Errorf("invalid difficulty, got %v expect %v index:%v", header.Difficulty.Int64(), diffNoTurn.Int64(), int(indexInTurn)%len(inTurnHV.Validators))
					}
				}
			}
		}
		if !valid {
			return fmt.Errorf("bsc Handler SyncBlockHeader, invalid signer")
		}

		err = addHeader(native, &header, phv, ctx)
		if err != nil {
			return fmt.Errorf("bsc Handler SyncBlockHeader, addHeader err: %v", err)
		}

		scom.NotifyPutHeader(native, headerParams.ChainID, header.Number.Uint64(), header.Hash().Hex())
	}
	return nil
}

// SyncCrossChainMsg ...
func (h *Handler) SyncCrossChainMsg(native *native.NativeContract) error {
	return nil
}

func isHeaderExist(native *native.NativeContract, headerHash common.Hash, ctx *Context) (bool, error) {
	headerStore, err := scom.GetHeaderIndex(native, ctx.ChainID, headerHash.Bytes())
	if err != nil {
		return false, fmt.Errorf("bsc Handler isHeaderExist error: %v", err)
	}

	return headerStore != nil, nil
}

func verifySignature(native *native.NativeContract, header *types.Header, ctx *Context) (signer common.Address, err error) {
	return verifyHeader(native, header, ctx)
}

func verifyHeader(native *native.NativeContract, header *types.Header, ctx *Context) (signer common.Address, err error) {

	// Don't waste time checking blocks from the future
	if header.Time > uint64(time.Now().Unix()) {
		err = errors.New("block in the future")
		return
	}

	// Check that the extra-data contains both the vanity and signature
	if len(header.Extra) < extraVanity {
		err = errors.New("extra-data 32 byte vanity prefix missing")
		return
	}
	if len(header.Extra) < extraVanity+extraSeal {
		err = errors.New("extra-data 65 byte signature suffix missing")
		return
	}

	// Ensure that the extra-data contains a signer list on checkpoint, but none otherwise
	signersBytes := len(header.Extra) - extraVanity - extraSeal

	if signersBytes%common.AddressLength != 0 {
		err = errors.New("invalid signer list")
		return
	}

	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != (common.Hash{}) {
		err = errors.New("non-zero mix digest")
		return
	}

	// Ensure that the block doesn't contain any uncles which are meaningless in PoA
	if header.UncleHash != uncleHash {
		err = errors.New("non empty uncle hash")
		return
	}

	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if header.Difficulty == nil || (header.Difficulty.Cmp(diffInTurn) != 0 && header.Difficulty.Cmp(diffNoTurn) != 0) {
		err = errors.New("invalid difficulty")
		return
	}

	// All basic checks passed, verify cascading fields
	return verifyCascadingFields(native, header, ctx)
}

func verifyCascadingFields(native *native.NativeContract, header *types.Header, ctx *Context) (signer common.Address, err error) {

	number := header.Number.Uint64()

	parent, err := getHeader(native, header.ParentHash, ctx.ChainID)
	if err != nil {
		return
	}

	if parent.Header.Number.Uint64() != number-1 {
		err = errors.New("unknown ancestor")
		return
	}

	// Verify that the gas limit is <= 2^63-1
	capacity := uint64(0x7fffffffffffffff)
	if header.GasLimit > capacity {
		err = fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, capacity)
		return
	}
	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed > header.GasLimit {
		err = fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit)
		return
	}

	// Verify that the gas limit remains within allowed bounds
	diff := int64(parent.Header.GasLimit) - int64(header.GasLimit)
	if diff < 0 {
		diff *= -1
	}
	limit := parent.Header.GasLimit / GasLimitBoundDivisor

	if uint64(diff) >= limit || header.GasLimit < params.MinGasLimit {
		err = fmt.Errorf("invalid gas limit: have %d, want %d += %d", header.GasLimit, parent.Header.GasLimit, limit)
		return
	}

	return verifySeal(native, header, ctx)
}

// for test
var mockSigner common.Address

func verifySeal(native *native.NativeContract, header *types.Header, ctx *Context) (signer common.Address, err error) {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		err = errors.New("unknown block")
		return
	}

	if mockSigner != (common.Address{}) {
		return mockSigner, nil
	}
	// Resolve the authorization key and check against validators
	signer, err = ecrecover(header, ctx.ExtraInfo.ChainID)
	if err != nil {
		return
	}

	if signer != header.Coinbase {
		err = errors.New("coinbase do not match with signature")
		return
	}

	return
}

// ecrecover extracts the Ethereum account address from a signed header.
func ecrecover(header *types.Header, chainID *big.Int) (common.Address, error) {
	// Retrieve the signature from the header extra-data
	if len(header.Extra) < extraSeal {
		return common.Address{}, errors.New("extra-data 65 byte signature suffix missing")
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]

	// Recover the public key and the Ethereum address
	pubkey, err := crypto.Ecrecover(SealHash(header, chainID).Bytes(), signature)
	if err != nil {
		return common.Address{}, err
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	return signer, nil
}

// SealHash returns the hash of a block prior to it being sealed.
func SealHash(header *types.Header, chainID *big.Int) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	encodeSigHeader(hasher, header, chainID)
	hasher.Sum(hash[:0])
	return hash
}

func encodeSigHeader(w io.Writer, header *types.Header, chainID *big.Int) {
	err := rlp.Encode(w, []interface{}{
		chainID,
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
		header.Extra[:len(header.Extra)-65], // this will panic if extra is too short, should check before calling encodeSigHeader
		header.MixDigest,
		header.Nonce,
	})
	if err != nil {
		panic("can't encode: " + err.Error())
	}
}

func isGenesisStored(native *native.NativeContract, params *scom.SyncGenesisHeaderParam) (stored bool, err error) {
	genesis, err := getGenesis(native, params.ChainID)
	if err != nil {
		return
	}

	stored = genesis != nil
	return
}

func getGenesis(native *native.NativeContract, chainID uint64) (genesisHeader *GenesisHeader, err error) {

	genesisBytes, err := scom.GetGenesisHeader(native, chainID)
	if err != nil {
		err = fmt.Errorf("getGenesis, GetCacheDB err:%v", err)
		return
	}

	if genesisBytes == nil {
		return
	}

	{
		genesisHeader = &GenesisHeader{}
		err = json.Unmarshal(genesisBytes, &genesisHeader)
		if err != nil {
			err = fmt.Errorf("getGenesis, json.Unmarshal err:%v", err)
			return
		}
	}

	return
}

func storeGenesis(native *native.NativeContract, params *scom.SyncGenesisHeaderParam, genesisHeader *GenesisHeader) error {
	genesisBytes, err := json.Marshal(genesisHeader)
	if err != nil {
		return err
	}

	scom.SetGenesisHeader(native, params.ChainID, genesisBytes)
	headerWithSum := &HeaderWithDifficultySum{Header: &genesisHeader.Header, DifficultySum: genesisHeader.Header.Difficulty}

	if err := putHeaderWithSum(native, params.ChainID, headerWithSum); err != nil {
		return err
	}

	putCanonicalHeight(native, params.ChainID, genesisHeader.Header.Number.Uint64())
	putCanonicalHash(native, params.ChainID, genesisHeader.Header.Number.Uint64(), genesisHeader.Header.Hash())

	scom.NotifyPutHeader(native, params.ChainID, genesisHeader.Header.Number.Uint64(), genesisHeader.Header.Hash().Hex())
	return nil
}

// ParseValidators ...
func ParseValidators(validatorsBytes []byte) ([]common.Address, error) {
	if len(validatorsBytes)%common.AddressLength != 0 {
		return nil, errors.New("invalid validators bytes")
	}
	n := len(validatorsBytes) / common.AddressLength
	result := make([]common.Address, n)
	for i := 0; i < n; i++ {
		address := make([]byte, common.AddressLength)
		copy(address, validatorsBytes[i*common.AddressLength:(i+1)*common.AddressLength])
		result[i] = common.BytesToAddress(address)
	}
	return result, nil
}

func putHeaderWithSum(native *native.NativeContract, chainID uint64, headerWithSum *HeaderWithDifficultySum) (err error) {
	headerBytes, err := json.Marshal(headerWithSum)
	if err != nil {
		return
	}

	scom.SetHeaderIndex(native, chainID, headerWithSum.Header.Hash().Bytes(), headerBytes)
	return
}

func putCanonicalHeight(native *native.NativeContract, chainID uint64, height uint64) {
	scom.SetCurrentHeight(native, chainID, utils.GetUint64Bytes(height))
}

func putCanonicalHash(native *native.NativeContract, chainID uint64, height uint64, hash common.Hash) {
	scom.SetMainChain(native, chainID, height, hash.Bytes())
}

func addHeader(native *native.NativeContract, header *types.Header, phv *HeightAndValidators, ctx *Context) (err error) {
	parentHeader, err := getHeader(native, header.ParentHash, ctx.ChainID)
	if err != nil {
		return
	}

	cheight, err := GetCanonicalHeight(native, ctx.ChainID)
	if err != nil {
		return
	}
	cheader, err := GetCanonicalHeader(native, ctx.ChainID, cheight)
	if err != nil {
		return
	}
	if cheader == nil {
		err = fmt.Errorf("getCanonicalHeader returns nil")
		return
	}

	localTd := cheader.DifficultySum
	externTd := new(big.Int).Add(header.Difficulty, parentHeader.DifficultySum)

	headerWithSum := &HeaderWithDifficultySum{Header: header, DifficultySum: externTd, EpochParentHash: phv.Hash}
	err = putHeaderWithSum(native, ctx.ChainID, headerWithSum)
	if err != nil {
		return
	}

	if externTd.Cmp(localTd) > 0 {
		// Delete any canonical number assignments above the new head
		var headerWithSum *HeaderWithDifficultySum
		for i := header.Number.Uint64() + 1; ; i++ {
			headerWithSum, err = GetCanonicalHeader(native, ctx.ChainID, i)
			if err != nil {
				return
			}
			if headerWithSum == nil {
				break
			}

			deleteCanonicalHash(native, ctx.ChainID, i)
		}

		// Overwrite any stale canonical number assignments
		var (
			hash       common.Hash
			headHeader *HeaderWithDifficultySum
		)
		cheight := header.Number.Uint64() - 1
		headHash := header.ParentHash

		for {
			hash, err = getCanonicalHash(native, ctx.ChainID, cheight)
			if err != nil {
				return
			}
			if hash == headHash {
				break
			}

			putCanonicalHash(native, ctx.ChainID, cheight, headHash)
			headHeader, err = getHeader(native, headHash, ctx.ChainID)
			if err != nil {
				return
			}
			headHash = headHeader.Header.ParentHash
			cheight--
		}

		// Extend the canonical chain with the new header
		putCanonicalHash(native, ctx.ChainID, header.Number.Uint64(), header.Hash())
		putCanonicalHeight(native, ctx.ChainID, header.Number.Uint64())
	}

	return nil
}

func getPrevHeightAndValidators(native *native.NativeContract, header *types.Header, ctx *Context) (phv, pphv *HeightAndValidators, lastSeenHeight int64, err error) {

	genesis, err := getGenesis(native, ctx.ChainID)
	if err != nil {
		err = fmt.Errorf("bsc Handler getGenesis error: %v", err)
		return
	}

	if genesis == nil {
		err = fmt.Errorf("bsc Handler genesis not set")
		return
	}

	genesisHeaderHash := genesis.Header.Hash()
	if header.Hash() == genesisHeaderHash {
		err = fmt.Errorf("genesis header should not be synced again")
		return
	}

	lastSeenHeight = -1
	targetCoinbase := header.Coinbase
	if header.ParentHash == genesisHeaderHash {
		if genesis.Header.Coinbase == targetCoinbase {
			lastSeenHeight = genesis.Header.Number.Int64()
		}

		phv = &genesis.PrevValidators[0]
		phv.Hash = &genesisHeaderHash
		pphv = &genesis.PrevValidators[1]
		return
	}

	prevHeaderWithSum, err := getHeader(native, header.ParentHash, ctx.ChainID)
	if err != nil {
		err = fmt.Errorf("bsc Handler getHeader error: %v", err)
		return
	}

	if prevHeaderWithSum.Header.Coinbase == targetCoinbase {
		lastSeenHeight = prevHeaderWithSum.Header.Number.Int64()
	} else {
		nextRecentParentHash := prevHeaderWithSum.Header.ParentHash
		defer func() {
			if err == nil {
				maxV := len(phv.Validators)
				if maxV < len(pphv.Validators) {
					maxV = len(pphv.Validators)
				}
				maxLimit := maxV / 2
				for i := 0; i < maxLimit-1; i++ {
					prevHeaderWithSum, err := getHeader(native, nextRecentParentHash, ctx.ChainID)
					if err != nil {
						err = fmt.Errorf("bsc Handler getHeader error: %v", err)
						return
					}
					if prevHeaderWithSum.Header.Coinbase == targetCoinbase {
						lastSeenHeight = prevHeaderWithSum.Header.Number.Int64()
						return
					}

					if nextRecentParentHash == genesisHeaderHash {
						return
					}
					nextRecentParentHash = prevHeaderWithSum.Header.ParentHash
				}
			}
		}()
	}

	var (
		validators     []common.Address
		nextParentHash common.Hash
	)

	currentPV := &phv

	for {

		if len(prevHeaderWithSum.Header.Extra) > extraVanity+extraSeal {
			validators, err = ParseValidators(prevHeaderWithSum.Header.Extra[extraVanity : len(prevHeaderWithSum.Header.Extra)-extraSeal])
			if err != nil {
				err = fmt.Errorf("bsc Handler ParseValidators error: %v", err)
				return
			}
			*currentPV = &HeightAndValidators{
				Height:     prevHeaderWithSum.Header.Number,
				Validators: validators,
			}
			switch *currentPV {
			case phv:
				hash := prevHeaderWithSum.Header.Hash()
				phv.Hash = &hash
				currentPV = &pphv
			case pphv:
				return
			default:
				err = fmt.Errorf("bug in bsc Handler")
				return
			}
		}

		nextParentHash = prevHeaderWithSum.Header.ParentHash
		if prevHeaderWithSum.EpochParentHash != nil {
			nextParentHash = *prevHeaderWithSum.EpochParentHash
		}

		if nextParentHash == genesisHeaderHash {
			switch *currentPV {
			case phv:
				phv = &genesis.PrevValidators[0]
				phv.Hash = &genesisHeaderHash
				pphv = &genesis.PrevValidators[1]
			case pphv:
				pphv = &genesis.PrevValidators[0]
			default:
				err = fmt.Errorf("bug in bsc Handler")
				return
			}
			return
		}

		prevHeaderWithSum, err = getHeader(native, nextParentHash, ctx.ChainID)
		if err != nil {
			err = fmt.Errorf("bsc Handler getHeader error: %v", err)
			return
		}

	}
}

func getHeader(native *native.NativeContract, hash common.Hash, chainID uint64) (headerWithSum *HeaderWithDifficultySum, err error) {

	headerStore, err := scom.GetHeaderIndex(native, chainID, hash.Bytes())
	if err != nil {
		return nil, fmt.Errorf("bsc Handler getHeader error: %v", err)
	}
	if headerStore == nil {
		return nil, fmt.Errorf("bsc Handler getHeader, can not find any header records")
	}

	headerWithSum = &HeaderWithDifficultySum{}
	if err := json.Unmarshal(headerStore, &headerWithSum); err != nil {
		return nil, fmt.Errorf("bsc Handler getHeader, deserialize header error: %v", err)
	}

	return
}

// GetCanonicalHeight ...
func GetCanonicalHeight(native *native.NativeContract, chainID uint64) (height uint64, err error) {
	heightStore, err := scom.GetCurrentHeight(native, chainID)
	if err != nil {
		err = fmt.Errorf("bsc Handler GetCanonicalHeight err:%v", err)
		return
	}

	height = utils.GetBytesUint64(heightStore)
	return
}

// GetCanonicalHeader ...
func GetCanonicalHeader(native *native.NativeContract, chainID uint64, height uint64) (headerWithSum *HeaderWithDifficultySum, err error) {
	hash, err := getCanonicalHash(native, chainID, height)
	if err != nil {
		return
	}

	if hash == (common.Hash{}) {
		return
	}

	headerWithSum, err = getHeader(native, hash, chainID)
	return
}

func getCanonicalHash(native *native.NativeContract, chainID uint64, height uint64) (hash common.Hash, err error) {
	hashBytesStore, err := scom.GetMainChain(native, chainID, height)
	if err != nil {
		return
	}
	if hashBytesStore == nil {
		return
	}

	hash = common.BytesToHash(hashBytesStore)
	return
}

func deleteCanonicalHash(native *native.NativeContract, chainID uint64, height uint64) {
	scom.DelMainChain(native, chainID, height)
}
