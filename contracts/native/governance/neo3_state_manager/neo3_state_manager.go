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
package neo3_state_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/contract"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "neo3 state manager"

const (
	MethodContractName                  = "name"
	MethodGetCurrentStateValidator      = "getCurrentStateValidator"
	MethodRegisterStateValidator        = "registerStateValidator"
	MethodApproveRegisterStateValidator = "approveRegisterStateValidator"
	MethodRemoveStateValidator          = "removeStateValidator"
	MethodApproveRemoveStateValidator   = "approveRemoveStateValidator"

	//key prefix
	STATE_VALIDATOR           = "stateValidator"
	STATE_VALIDATOR_APPLY     = "stateValidatorApply"
	STATE_VALIDATOR_REMOVE    = "stateValidatorRemove"
	STATE_VALIDATOR_APPLY_ID  = "stateValidatorApplyID"
	STATE_VALIDATOR_REMOVE_ID = "stateValidatorRemoveID"
)

var (
	this     = native.NativeContractAddrMap[native.NativeNeo3StateManager]
	gasTable = map[string]uint64{
		MethodContractName:                  0,
		MethodGetCurrentStateValidator:      0,
		MethodRegisterStateValidator:        100000,
		MethodApproveRegisterStateValidator: 0,
		MethodRemoveStateValidator:          0,
		MethodApproveRemoveStateValidator:   0,
	}

	ABI *abi.ABI
)

func InitNeo3StateManager() {
	ABI = GetABI()
	native.Contracts[this] = RegisterNeo3StateManagerContract
}

func RegisterNeo3StateManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodContractName, Name)
	s.Register(MethodGetCurrentStateValidator, GetCurrentStateValidator)
	s.Register(MethodRegisterStateValidator, RegisterStateValidator)
	s.Register(MethodApproveRegisterStateValidator, ApproveRegisterStateValidator)
	s.Register(MethodRemoveStateValidator, RemoveStateValidator)
	s.Register(MethodApproveRemoveStateValidator, ApproveRemoveStateValidator)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(ABI, MethodContractName, contractName)
}

func GetCurrentStateValidator(native *native.NativeContract) ([]byte, error) {
	data, err := getStateValidators(native)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentStateValidator, getStateValidators error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetCurrentStateValidator, data)
}

func RegisterStateValidator(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	param := &StateValidatorListParam{}
	if err := utils.UnpackMethod(ABI, MethodRegisterStateValidator, param, ctx.Payload); err != nil {
		return nil, err
	}

	//check witness
	err := contract.ValidateOwner(native, param.Address)
	if err != nil {
		return nil, fmt.Errorf("RegisterStateValidator, ValidateOwner error: %v", err)
	}

	if err := putStateValidatorApply(native, param); err != nil {
		return nil, fmt.Errorf("RegisterStateValidator, putStateValidatorApply error: %v", err)
	}

	return utils.PackOutputs(ABI, MethodRegisterStateValidator, true)
}

func ApproveRegisterStateValidator(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	param := &ApproveStateValidatorParam{}
	if err := utils.UnpackMethod(ABI, MethodApproveRegisterStateValidator, param, ctx.Payload); err != nil {
		return nil, err
	}

	//check witness
	err := contract.ValidateOwner(native, param.Address)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterStateValidator, ValidateOwner error: %v", err)
	}

	stateValidatorListParam, err := getStateValidatorApply(native, param.ID)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterStateValidator, getStateValidatorApply error: %v", err)
	}

	//check consensus signs
	ok, err := node_manager.CheckConsensusSigns(native, MethodApproveRegisterStateValidator, utils.GetUint64Bytes(param.ID), param.Address)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterStateValidator, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.PackOutputs(ABI, MethodApproveRegisterStateValidator, true)
	}

	err = putStateValidators(native, stateValidatorListParam.StateValidators)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterStateValidator, putStateValidators error: %v", err)
	}
	native.GetCacheDB().Delete(utils.ConcatKey(utils.Neo3StateManagerContractAddress, []byte(STATE_VALIDATOR_APPLY), utils.GetUint64Bytes(param.ID)))

	err=native.AddNotify(ABI, []string{EventApproveRegisterStateValidator}, param.ID)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterStateValidator, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodApproveRegisterStateValidator, true)
}

func RemoveStateValidator(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	param := &StateValidatorListParam{}
	if err := utils.UnpackMethod(ABI, MethodRemoveStateValidator, param, ctx.Payload); err != nil {
		return nil, err
	}

	//check witness
	err := contract.ValidateOwner(native, param.Address)
	if err != nil {
		return nil, fmt.Errorf("RemoveStateValidator, ValidateOwner error: %v", err)
	}

	err = putStateValidatorRemove(native, param)
	if err != nil {
		return nil, fmt.Errorf("RemoveStateValidator, putStateValidatorRemove error: %v", err)
	}

	return utils.PackOutputs(ABI, MethodRemoveStateValidator, true)
}

func ApproveRemoveStateValidator(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	param := &ApproveStateValidatorParam{}
	if err := utils.UnpackMethod(ABI, MethodApproveRemoveStateValidator, param, ctx.Payload); err != nil {
		return nil, err
	}

	//check witness
	err := contract.ValidateOwner(native, param.Address)
	if err != nil {
		return nil, fmt.Errorf("ApproveRemoveStateValidator, ValidateOwner error: %v", err)
	}
	//get sv list
	svListParam, err := getStateValidatorRemove(native, param.ID)
	if err != nil {
		return nil, fmt.Errorf("ApproveRemoveStateValidator, getStateValidatorRemove error: %v", err)
	}
	//check consensus signs
	ok, err := node_manager.CheckConsensusSigns(native, MethodApproveRemoveStateValidator, utils.GetUint64Bytes(param.ID), param.Address)
	if err != nil {
		return nil, fmt.Errorf("ApproveRemoveStateValidator, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.PackOutputs(ABI, MethodApproveRemoveStateValidator, true)
	}
	//remove sv
	err = removeStateValidators(native, svListParam.StateValidators)
	if err != nil {
		return nil, fmt.Errorf("ApproveRemoveStateValidator, removeStateValidators error: %v", err)
	}
	//delete removed sv
	native.GetCacheDB().Delete(utils.ConcatKey(utils.Neo3StateManagerContractAddress, []byte(STATE_VALIDATOR_REMOVE), utils.GetUint64Bytes(param.ID)))
	err = native.AddNotify(ABI, []string{EventApproveRemoveStateValidator}, param.ID)
	if err != nil {
		return nil, fmt.Errorf("ApproveRemoveStateValidator, AddNofity error: %v", err)
	}

	return utils.PackOutputs(ABI, MethodApproveRemoveStateValidator, true)
}
