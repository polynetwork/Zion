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
package relayer_manager

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	polycomm "github.com/polynetwork/poly/common"
	cstates "github.com/polynetwork/poly/core/states"
)

func putRelayer(native *native.NativeContract, relayer common.Address) error {
	contract := utils.RelayerManagerContractAddress
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(RELAYER), relayer[:]), cstates.GenRawStorageItem(relayer[:]))
	return nil
}

func putRelayerApply(native *native.NativeContract, relayerListParam *RelayerListParam) error {
	contract := utils.RelayerManagerContractAddress
	applyID, err := getApplyID(native)
	if err != nil {
		return fmt.Errorf("putRelayerApply, getApplyID error: %v", err)
	}
	newApplyID := applyID + 1
	err = putApplyID(native, newApplyID)
	if err != nil {
		return fmt.Errorf("putRelayerApply, putApplyID error: %v", err)
	}
	sink := polycomm.NewZeroCopySink(nil)
	relayerListParam.Serialization(sink)
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(RELAYER_APPLY), utils.GetUint64Bytes(applyID)),
		cstates.GenRawStorageItem(sink.Bytes()))

	err = native.AddNotify(ABI, []string{EventRegisterRelayer}, applyID)
	if err != nil {
		return fmt.Errorf("putRelayerApply, AddNotify error: %v", err)
	}
	return nil
}

func getApplyID(native *native.NativeContract) (uint64, error) {
	contract := utils.RelayerManagerContractAddress
	applyIDStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(APPLY_ID)))
	if err != nil {
		return 0, fmt.Errorf("getApplyID, get applyIDStore error: %v", err)
	}
	var applyID uint64 = 0
	if applyIDStore != nil {
		applyIDBytes, err := cstates.GetValueFromRawStorageItem(applyIDStore)
		if err != nil {
			return 0, fmt.Errorf("getApplyID, deserialize from raw storage item err:%v", err)
		}
		applyID = utils.GetBytesUint64(applyIDBytes)
	}
	return applyID, nil
}

func putApplyID(native *native.NativeContract, applyID uint64) error {
	contract := utils.RelayerManagerContractAddress
	applyIDByte := utils.GetUint64Bytes(applyID)

	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(APPLY_ID)), cstates.GenRawStorageItem(applyIDByte))
	return nil
}

func getRelayerApply(native *native.NativeContract, applyID uint64) (*RelayerListParam, error) {
	contract := utils.RelayerManagerContractAddress
	relayerListParamStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(RELAYER_APPLY), utils.GetUint64Bytes(applyID)))
	if err != nil {
		return nil, fmt.Errorf("getRelayerApply, get relayerListParamStore error: %v", err)
	}
	if relayerListParamStore == nil {
		return nil, fmt.Errorf("getRelayerApply, can't find any record")
	}
	relayerListParam := new(RelayerListParam)
	relayerListParamBytes, err := cstates.GetValueFromRawStorageItem(relayerListParamStore)
	if err != nil {
		return nil, fmt.Errorf("getRelayerApply, deserialize from raw storage item err:%v", err)
	}
	err = relayerListParam.Deserialization(polycomm.NewZeroCopySource(relayerListParamBytes))
	if err != nil {
		return nil, fmt.Errorf("getRelayerApply, relayerListParam.Deserialization fail:%v", err)
	}
	return relayerListParam, nil
}

func putRelayerRemove(native *native.NativeContract, relayerListParam *RelayerListParam) error {
	contract := utils.RelayerManagerContractAddress
	removeID, err := getRemoveID(native)
	if err != nil {
		return fmt.Errorf("putRelayerRemove, getRemoveID error: %v", err)
	}
	newRemoveID := removeID + 1
	err = putRemoveID(native, newRemoveID)
	if err != nil {
		return fmt.Errorf("putRelayerRemove, putRemoveID error: %v", err)
	}
	sink := polycomm.NewZeroCopySink(nil)
	relayerListParam.Serialization(sink)
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(RELAYER_REMOVE), utils.GetUint64Bytes(removeID)),
		cstates.GenRawStorageItem(sink.Bytes()))

	err = native.AddNotify(ABI, []string{EventRemoveRelayer}, removeID)
	if err != nil {
		return fmt.Errorf("putRelayerRemove, AddNotify error: %v", err)
	}
	return nil
}

func getRemoveID(native *native.NativeContract) (uint64, error) {
	contract := utils.RelayerManagerContractAddress
	removeIDStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(REMOVE_ID)))
	if err != nil {
		return 0, fmt.Errorf("getRemoveID, get removeIDStore error: %v", err)
	}
	var removeID uint64 = 0
	if removeIDStore != nil {
		removeIDBytes, err := cstates.GetValueFromRawStorageItem(removeIDStore)
		if err != nil {
			return 0, fmt.Errorf("getRemoveID, deserialize from raw storage item err:%v", err)
		}
		removeID = utils.GetBytesUint64(removeIDBytes)
	}
	return removeID, nil
}

func putRemoveID(native *native.NativeContract, removeID uint64) error {
	contract := utils.RelayerManagerContractAddress
	removeIDByte := utils.GetUint64Bytes(removeID)

	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(REMOVE_ID)), cstates.GenRawStorageItem(removeIDByte))
	return nil
}

func getRelayerRemove(native *native.NativeContract, removeID uint64) (*RelayerListParam, error) {
	contract := utils.RelayerManagerContractAddress
	relayerListParamStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(RELAYER_REMOVE), utils.GetUint64Bytes(removeID)))
	if err != nil {
		return nil, fmt.Errorf("getRelayerRemove, get relayerListParamStore error: %v", err)
	}
	if relayerListParamStore == nil {
		return nil, fmt.Errorf("getRelayerRemove, can't find any record")
	}
	relayerListParam := new(RelayerListParam)
	relayerListParamBytes, err := cstates.GetValueFromRawStorageItem(relayerListParamStore)
	if err != nil {
		return nil, fmt.Errorf("getRelayerRemove, deserialize from raw storage item err:%v", err)
	}
	err = relayerListParam.Deserialization(polycomm.NewZeroCopySource(relayerListParamBytes))
	if err != nil {
		return nil, fmt.Errorf("getRelayerRemove, Deserialization fail:%v", err)
	}
	return relayerListParam, nil
}
