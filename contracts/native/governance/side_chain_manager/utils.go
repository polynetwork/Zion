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

package side_chain_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
)

func GetSideChainApply(native *native.NativeContract, chanid uint64) (*SideChain, error) {
	contract := utils.SideChainManagerContractAddress
	chainidByte := utils.GetUint64Bytes(chanid)

	sideChainStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(SIDE_CHAIN_APPLY),
		chainidByte))
	if err != nil {
		return nil, fmt.Errorf("getRegisterSideChain,get registerSideChainRequestStore error: %v", err)
	}
	sideChain := new(SideChain)
	if sideChainStore != nil {
		if err := rlp.DecodeBytes(sideChainStore, sideChain); err != nil {
			return nil, fmt.Errorf("getRegisterSideChain, deserialize sideChain error: %v", err)
		}
		return sideChain, nil
	} else {
		return nil, nil
	}
}

func putSideChainApply(native *native.NativeContract, sideChain *SideChain) error {
	contract := utils.SideChainManagerContractAddress
	chainidByte := utils.GetUint64Bytes(sideChain.ChainID)

	blob, err := rlp.EncodeToBytes(sideChain)
	if err != nil {
		return fmt.Errorf("putRegisterSideChain, sideChain.Serialization error: %v", err)
	}

	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(SIDE_CHAIN_APPLY), chainidByte), blob)
	return nil
}

func GetSideChainObject(native *native.NativeContract, chainID uint64) (*SideChain, error) {
	contract := utils.SideChainManagerContractAddress
	chainIDByte := utils.GetUint64Bytes(chainID)

	sideChainStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(SIDE_CHAIN),
		chainIDByte))
	if err != nil {
		return nil, fmt.Errorf("getSideChain,get registerSideChainRequestStore error: %v", err)
	}
	sideChain := new(SideChain)
	if sideChainStore != nil {
		if err := rlp.DecodeBytes(sideChainStore, sideChain); err != nil {
			return nil, fmt.Errorf("getSideChain, deserialize sideChain error: %v", err)
		}
		return sideChain, nil
	}
	return nil, nil

}

func PutSideChain(native *native.NativeContract, sideChain *SideChain) error {
	contract := utils.SideChainManagerContractAddress
	chainidByte := utils.GetUint64Bytes(sideChain.ChainID)

	blob, err := rlp.EncodeToBytes(sideChain)
	if err != nil {
		return fmt.Errorf("putSideChain, sideChain.Serialization error: %v", err)
	}

	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(SIDE_CHAIN), chainidByte), blob)
	return nil
}

func getUpdateSideChain(native *native.NativeContract, chanid uint64) (*SideChain, error) {
	contract := utils.SideChainManagerContractAddress
	chainidByte := utils.GetUint64Bytes(chanid)

	sideChainStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(UPDATE_SIDE_CHAIN_REQUEST),
		chainidByte))
	if err != nil {
		return nil, fmt.Errorf("getUpdateSideChain,get registerSideChainRequestStore error: %v", err)
	}
	sideChain := new(SideChain)
	if sideChainStore != nil {
		if err := rlp.DecodeBytes(sideChainStore, sideChain); err != nil {
			return nil, fmt.Errorf("getUpdateSideChain, deserialize sideChain error: %v", err)
		}
		return sideChain, nil
	} else {
		return nil, nil
	}
}

func putUpdateSideChain(native *native.NativeContract, sideChain *SideChain) error {
	contract := utils.SideChainManagerContractAddress
	chainidByte := utils.GetUint64Bytes(sideChain.ChainID)

	blob, err := rlp.EncodeToBytes(sideChain)
	if err != nil {
		return fmt.Errorf("putUpdateSideChain, sideChain.Serialization error: %v", err)
	}

	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(UPDATE_SIDE_CHAIN_REQUEST), chainidByte), blob)
	return nil
}

func getQuitSideChain(native *native.NativeContract, chainid uint64) error {
	contract := utils.SideChainManagerContractAddress
	chainidByte := utils.GetUint64Bytes(chainid)

	chainIDStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(QUIT_SIDE_CHAIN_REQUEST),
		chainidByte))
	if err != nil {
		return fmt.Errorf("getQuitSideChain, get registerSideChainRequestStore error: %v", err)
	}
	if chainIDStore != nil {
		return nil
	}
	return fmt.Errorf("getQuitSideChain, no record")
}

func putQuitSideChain(native *native.NativeContract, chainid uint64) error {
	contract := utils.SideChainManagerContractAddress
	chainidByte := utils.GetUint64Bytes(chainid)

	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(QUIT_SIDE_CHAIN_REQUEST), chainidByte), chainidByte)
	return nil
}

func PutFee(native *native.NativeContract, chainId uint64, fee *Fee) error {
	contract := utils.SideChainManagerContractAddress
	chainIdBytes := utils.GetUint64Bytes(chainId)
	blob, err := rlp.EncodeToBytes(fee)
	if err != nil {
		return fmt.Errorf("PutFee, rlp.EncodeToBytes fee error: %v", err)
	}
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(FEE), chainIdBytes), blob)
	return nil
}

func GetFeeObj(native *native.NativeContract, chainID uint64) (*Fee, error) {
	chainIDBytes := utils.GetUint64Bytes(chainID)
	key := utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(FEE), chainIDBytes)
	store, err := native.GetCacheDB().Get(key)
	if err != nil {
		return nil, fmt.Errorf("GetFee, get fee info store error: %v", err)
	}
	fee := &Fee{
		Fee: new(big.Int),
	}
	if store != nil {
		if err := rlp.DecodeBytes(store, fee); err != nil {
			return nil, fmt.Errorf("GetFee, deserialize fee error: %v", err)
		}
	}
	return fee, nil
}

func PutFeeInfo(native *native.NativeContract, chainId, view uint64, feeInfo *FeeInfo) error {
	chainIdBytes := utils.GetUint64Bytes(chainId)
	viewBytes := utils.GetUint64Bytes(view)
	key := utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(FEE_INFO), chainIdBytes, viewBytes)
	blob, err := rlp.EncodeToBytes(feeInfo)
	if err != nil {
		return fmt.Errorf("PutFeeInfo, rlp.EncodeToBytes fee info error: %v", err)
	}
	native.GetCacheDB().Put(key, blob)
	return nil
}

func GetFeeInfo(native *native.NativeContract, chainID, view uint64) (*FeeInfo, error) {
	chainIDBytes := utils.GetUint64Bytes(chainID)
	viewBytes := utils.GetUint64Bytes(view)
	key := utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(FEE_INFO), chainIDBytes, viewBytes)
	store, err := native.GetCacheDB().Get(key)
	if err != nil {
		return nil, fmt.Errorf("GetFeeInfo, get fee info store error: %v", err)
	}
	feeInfo := &FeeInfo{
		FeeInfo: make(map[common.Address]*big.Int),
	}
	if store != nil {
		if err := rlp.DecodeBytes(store, feeInfo); err != nil {
			return nil, fmt.Errorf("GetFeeInfo, deserialize fee info error: %v", err)
		}
	}
	return feeInfo, nil
}

func GetRippleExtraInfo(native *native.NativeContract, chainId uint64) (*RippleExtraInfo, error) {
	sideChainInfo, err := GetSideChainObject(native, chainId)
	if err != nil {
		return nil, fmt.Errorf("GetRippleExtraInfo, GetSideChainObject error: %v", err)
	}
	if sideChainInfo == nil {
		return nil, fmt.Errorf("GetRippleExtraInfo, side chain info is nil")
	}
	rippleExtraInfo := &RippleExtraInfo{
		Pks: make([][]byte, 0),
	}
	if err := rlp.DecodeBytes(sideChainInfo.ExtraInfo, rippleExtraInfo); err != nil {
		return nil, fmt.Errorf("GetRippleExtraInfo, deserialize info error: %v", err)
	}
	return rippleExtraInfo, nil
}

func PutRippleExtraInfo(native *native.NativeContract, chainId uint64, rippleExtraInfo *RippleExtraInfo) error {
	blob, err := rlp.EncodeToBytes(rippleExtraInfo)
	if err != nil {
		return fmt.Errorf("PutRippleExtraInfo, rlp.EncodeToBytes info error: %v", err)
	}
	sideChainInfo, err := GetSideChainObject(native, chainId)
	if err != nil {
		return fmt.Errorf("PutRippleExtraInfo, GetSideChainObject error: %v", err)
	}
	sideChainInfo.ExtraInfo = blob
	err = PutSideChain(native, sideChainInfo)
	if err != nil {
		return fmt.Errorf("PutRippleExtraInfo, PutSideChain error: %v", err)
	}
	return nil
}

func PutAssetBind(native *native.NativeContract, chainId uint64, assetBind *AssetBind) error {
	chainIDBytes := utils.GetUint64Bytes(chainId)
	key := utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(ASSET_BIND), chainIDBytes)
	blob, err := rlp.EncodeToBytes(assetBind)
	if err != nil {
		return fmt.Errorf("PutAssetBind, rlp.EncodeToBytes asset bind error: %v", err)
	}
	native.GetCacheDB().Put(key, blob)
	return nil
}

func GetAssetBind(native *native.NativeContract, chainId uint64) (*AssetBind, error) {
	chainIDBytes := utils.GetUint64Bytes(chainId)
	key := utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(ASSET_BIND), chainIDBytes)
	store, err := native.GetCacheDB().Get(key)
	if err != nil {
		return nil, fmt.Errorf("GetAssetBind, get asset map store error: %v", err)
	}
	assetBind := &AssetBind{
		AssetMap:     make(map[uint64][]byte),
		LockProxyMap: make(map[uint64][]byte),
	}
	if store != nil {
		if err := rlp.DecodeBytes(store, assetBind); err != nil {
			return nil, fmt.Errorf("GetAssetBind, deserialize info error: %v", err)
		}
	}
	return assetBind, nil
}
