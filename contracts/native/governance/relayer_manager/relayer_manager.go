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

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/contract"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "relayer manager"

const (
	//function name
	MethodContractName           = "name"
	MethodRegisterRelayer        = "registerRelayer"
	MethodApproveRegisterRelayer = "approveRegisterRelayer"
	MethodRemoveRelayer          = "removeRelayer"
	MethodApproveRemoveRelayer   = "approveRemoveRelayer"

	//key prefix
	RELAYER        = "relayer"
	RELAYER_APPLY  = "relayerApply"
	RELAYER_REMOVE = "relayerRemove"
	APPLY_ID       = "applyID"
	REMOVE_ID      = "removeID"
)

var (
	this     = native.NativeContractAddrMap[native.NativeRelayerManager]
	gasTable = map[string]uint64{
		MethodContractName:           0,
		MethodRegisterRelayer:        0,
		MethodApproveRegisterRelayer: 100000,
		MethodRemoveRelayer:          0,
		MethodApproveRemoveRelayer:   0,
	}

	ABI *abi.ABI
)

func InitRelayerManager() {
	ABI = GetABI()
	native.Contracts[this] = RegisterRelayerManagerContract
}

func RegisterRelayerManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodContractName, Name)
	s.Register(MethodRegisterRelayer, RegisterRelayer)
	s.Register(MethodApproveRegisterRelayer, ApproveRegisterRelayer)
	s.Register(MethodRemoveRelayer, RemoveRelayer)
	s.Register(MethodApproveRemoveRelayer, ApproveRemoveRelayer)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(ABI, MethodContractName, contractName)
}

func RegisterRelayer(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &RelayerListParam{}
	if err := utils.UnpackMethod(ABI, MethodRegisterRelayer, params, ctx.Payload); err != nil {
		return nil, err
	}

	//check witness
	err := contract.ValidateOwner(native, params.Address)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, checkWitness error: %v", err)
	}

	if err := putRelayerApply(native, params); err != nil {
		return nil, fmt.Errorf("RegisterRelayer, putRelayer error: %v", err)
	}

	return utils.PackOutputs(ABI, MethodRegisterRelayer, true)
}

func ApproveRegisterRelayer(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &ApproveRelayerParam{}
	if err := utils.UnpackMethod(ABI, MethodApproveRegisterRelayer, params, ctx.Payload); err != nil {
		return nil, err
	}

	//check witness
	err := contract.ValidateOwner(native, params.Address)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, checkWitness error: %v", err)
	}

	relayerListParam, err := getRelayerApply(native, params.ID)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterRelayer, getRelayerApply error: %v", err)
	}

	//check consensus signs
	ok, err := node_manager.CheckConsensusSigns(native, MethodApproveRegisterRelayer, utils.GetUint64Bytes(params.ID), params.Address)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterRelayer, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.PackOutputs(ABI, MethodApproveRegisterRelayer, true)
	}

	for _, address := range relayerListParam.AddressList {
		err = putRelayer(native, address)
		if err != nil {
			return nil, fmt.Errorf("ApproveRegisterRelayer, putRelayer error: %v", err)
		}
	}

	native.GetCacheDB().Delete(utils.ConcatKey(utils.RelayerManagerContractAddress, []byte(RELAYER_APPLY), utils.GetUint64Bytes(params.ID)))

	err = native.AddNotify(ABI, []string{EventApproveRegisterRelayer}, params.ID)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterRelayer, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodApproveRegisterRelayer, true)
}

func RemoveRelayer(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &RelayerListParam{}
	if err := utils.UnpackMethod(ABI, MethodRemoveRelayer, params, ctx.Payload); err != nil {
		return nil, err
	}

	//check witness
	err := contract.ValidateOwner(native, params.Address)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, checkWitness error: %v", err)
	}

	err = putRelayerRemove(native, params)
	if err != nil {
		return nil, fmt.Errorf("RemoveRelayer, putRelayer error: %v", err)
	}

	return utils.PackOutputs(ABI, MethodRemoveRelayer, true)
}

func ApproveRemoveRelayer(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &ApproveRelayerParam{}
	if err := utils.UnpackMethod(ABI, MethodApproveRemoveRelayer, params, ctx.Payload); err != nil {
		return nil, err
	}

	//check witness
	err := contract.ValidateOwner(native, params.Address)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, checkWitness error: %v", err)
	}

	relayerListParam, err := getRelayerRemove(native, params.ID)
	if err != nil {
		return nil, fmt.Errorf("ApproveRemoveRelayer, getRelayerRemove error: %v", err)
	}

	//check consensus signs
	ok, err := node_manager.CheckConsensusSigns(native, MethodApproveRemoveRelayer, utils.GetUint64Bytes(params.ID), params.Address)
	if err != nil {
		return nil, fmt.Errorf("ApproveRemoveRelayer, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.PackOutputs(ABI, MethodApproveRemoveRelayer, true)
	}

	for _, address := range relayerListParam.AddressList {
		native.GetCacheDB().Delete(utils.ConcatKey(utils.RelayerManagerContractAddress, []byte(RELAYER), address[:]))
	}
	err = native.AddNotify(ABI, []string{EventApproveRemoveRelayer}, params.ID)
	if err != nil {
		return nil, fmt.Errorf("ApproveRemoveRelayer, AddNotify error: %v", err)
	}

	return utils.PackOutputs(ABI, MethodApproveRemoveRelayer, true)
}
