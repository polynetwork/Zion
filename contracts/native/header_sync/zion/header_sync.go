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

package zion

import (
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/eth/types"
)

type ZionHandler struct {
}

func NewZionHandler() *ZionHandler {
	return &ZionHandler{}
}

func (this *ZionHandler) SyncGenesisHeader(native *native.NativeContract) error {
	//ctx := native.ContractRef().CurrentContext()
	//params := &scom.SyncGenesisHeaderParam{}
	//if err := utils.UnpackMethod(scom.ABI, scom.MethodSyncGenesisHeader, params, ctx.Payload); err != nil {
	//	return fmt.Errorf("SyncGenesisHeader, contract params deserialize error: %v", err)
	//}
	//
	//// Get current epoch operator
	//ok, err := node_manager.CheckConsensusSigns(native, scom.MethodSyncGenesisHeader, ctx.Payload, native.ContractRef().MsgSender())
	//if err != nil {
	//	return fmt.Errorf("SyncGenesisHeader, CheckConsensusSigns error: %v", err)
	//}
	//if !ok {
	//	return nil
	//}
	//
	//var header Header
	//err = json.Unmarshal(params.GenesisHeader, &header)
	//if err != nil {
	//	return fmt.Errorf("SyncGenesisHeader, json.Unmarshal header err: %v", err)
	//}
	//
	//headerStore, err := native.GetCacheDB().Get(utils.ConcatKey(utils.HeaderSyncContractAddress, []byte(scom.GENESIS_HEADER), utils.GetUint64Bytes(params.ChainID)))
	//if err != nil {
	//	return fmt.Errorf("ETHHandler GetHeaderByHeight, get blockHashStore error: %v", err)
	//}
	//if headerStore != nil {
	//	return fmt.Errorf("ETHHandler GetHeaderByHeight, genesis header had been initialized")
	//}
	//
	////block header storage
	//err = putGenesisBlockHeader(native, header, params.ChainID)
	//if err != nil {
	//	return fmt.Errorf("ETHHandler SyncGenesisHeader, put blockHeader error: %v", err)
	//}

	return nil
}

func (this *ZionHandler) SyncBlockHeader(native *native.NativeContract) error {
	//headerParams := &scom.SyncBlockHeaderParam{}
	//{
	//	ctx := native.ContractRef().CurrentContext()
	//	if err := utils.UnpackMethod(scom.ABI, scom.MethodSyncBlockHeader, headerParams, ctx.Payload); err != nil {
	//		return err
	//	}
	//}
	//caches := NewCaches(3, native)
	//for _, v := range headerParams.Headers {
	//	var header Header
	//	err := json.Unmarshal(v, &header)
	//	if err != nil {
	//		return fmt.Errorf("SyncBlockHeader, deserialize header err: %v", err)
	//	}
	//	headerHash := header.Hash()
	//	exist, err := IsHeaderExist(native, headerHash.Bytes(), headerParams.ChainID)
	//	if err != nil {
	//		return fmt.Errorf("SyncBlockHeader, check header exist err: %v", err)
	//	}
	//	if exist == true {
	//		log.Warnf("SyncBlockHeader, header has exist. Header: %s", string(v))
	//		continue
	//	}
	//	// get pre header
	//	parentHeader, parentDifficultySum, err := GetHeaderByHash(native, header.ParentHash.Bytes(), headerParams.ChainID)
	//	if err != nil {
	//		return fmt.Errorf("SyncBlockHeader, get the parent block failed. Error:%s, header: %s", err, string(v))
	//	}
	//	parentHeaderHash := parentHeader.Hash()
	//	/**
	//	this code source refer to https://github.com/ethereum/go-ethereum/blob/master/consensus/ethash/consensus.go
	//	verify header need to verify:
	//	1. parent hash
	//	2. extra size
	//	3. current time
	//	*/
	//	//verify whether parent hash validity
	//	if !bytes.Equal(parentHeaderHash.Bytes(), header.ParentHash.Bytes()) {
	//		return fmt.Errorf("SyncBlockHeader, parent header is not right. Header: %s", string(v))
	//	}
	//	//verify whether extra size validity
	//	if uint64(len(header.Extra)) > params.MaximumExtraDataSize {
	//		return fmt.Errorf("SyncBlockHeader, SyncBlockHeader extra-data too long: %d > %d, header: %s", len(header.Extra), params.MaximumExtraDataSize, string(v))
	//	}
	//	//verify current time validity
	//	if header.Time > uint64(time.Now().Add(allowedFutureBlockTime).Unix()) {
	//		return fmt.Errorf("SyncBlockHeader,  verify header time error:%s, checktime: %d, header: %s", consensus.ErrFutureBlock, time.Now().Add(allowedFutureBlockTime).Unix(), string(v))
	//	}
	//	//verify whether current header time and prevent header time validity
	//	if header.Time <= parentHeader.Time {
	//		return fmt.Errorf("SyncBlockHeader, verify header time fail. Header: %s", string(v))
	//	}
	//	// Verify that the gas limit is <= 2^63-1
	//	cap := uint64(0x7fffffffffffffff)
	//	if header.GasLimit > cap {
	//		return fmt.Errorf("SyncBlockHeader, invalid gasLimit: have %v, max %v, header: %s", header.GasLimit, cap, string(v))
	//	}
	//	// Verify that the gasUsed is <= gasLimit
	//	if header.GasUsed > header.GasLimit {
	//		return fmt.Errorf("SyncBlockHeader, invalid gasUsed: have %d, gasLimit %d, header: %s", header.GasUsed, header.GasLimit, string(v))
	//	}
	//	if isLondon(&header) {
	//		err = VerifyEip1559Header(parentHeader, &header)
	//	} else {
	//		err = VerifyGaslimit(parentHeader.GasLimit, header.GasLimit)
	//	}
	//	if err != nil {
	//		return fmt.Errorf("SyncBlockHeader, err:%v", err)
	//	}
	//
	//	//verify difficulty
	//	var expected *big.Int
	//	if isLondon(&header) {
	//		expected = makeDifficultyCalculator(big.NewInt(9700000))(header.Time, parentHeader)
	//	} else {
	//		return fmt.Errorf("SyncBlockHeader, header before london fork is no longer supported")
	//	}
	//	if expected.Cmp(header.Difficulty) != 0 {
	//		return fmt.Errorf("SyncBlockHeader, invalid difficulty: have %v, want %v, header: %s", header.Difficulty, expected, string(v))
	//	}
	//	// verfify header
	//	err = this.verifyHeader(&header, caches)
	//	if err != nil {
	//		return fmt.Errorf("SyncBlockHeader, verify header error: %v, header: %s", err, string(v))
	//	}
	//	//block header storage
	//	hederDifficultySum := new(big.Int).Add(header.Difficulty, parentDifficultySum)
	//	err = putBlockHeader(native, header, hederDifficultySum, headerParams.ChainID)
	//	if err != nil {
	//		return fmt.Errorf("SyncGenesisHeader, put blockHeader error: %v, header: %s", err, string(v))
	//	}
	//	// get current header of main
	//	currentHeader, currentDifficultySum, err := GetCurrentHeader(native, headerParams.ChainID)
	//	if err != nil {
	//		return fmt.Errorf("SyncBlockHeader, get the current block failed. error:%s", err)
	//	}
	//	if bytes.Equal(currentHeader.Hash().Bytes(), header.ParentHash.Bytes()) {
	//		appendHeader2Main(native, header.Number.Uint64(), headerHash, headerParams.ChainID)
	//	} else {
	//		//
	//		if hederDifficultySum.Cmp(currentDifficultySum) > 0 {
	//			RestructChain(native, currentHeader, &header, headerParams.ChainID)
	//		}
	//	}
	//}
	//caches.deleteCaches()
	return nil
}

func (this *ZionHandler) SyncCrossChainMsg(native *native.NativeContract) error {
	return nil
}

func (this *ZionHandler) verifyHeader(header *types.Header) error {
	// todo
	return nil
}