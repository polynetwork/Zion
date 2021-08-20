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
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "relayer manager"

const (
	MethodContractName           = "name"
	MethodRegisterRelayer        = "registerRelayer"
	MethodApproveRegisterRelayer = "approveRegisterRelayer"
	MethodRemoveRelayer          = "removeRelayer"
	MethodApproveRemoveRelayer   = "approveRemoveRelayer"
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

func RegisterRelayer(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &RelayerListParam{}
	if err := utils.UnpackMethod(ABI, MethodRegisterRelayer, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodRegisterRelayer, true)
}

func ApproveRegisterRelayer(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ApproveRelayerParam{}
	if err := utils.UnpackMethod(ABI, MethodApproveRegisterRelayer, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodApproveRegisterRelayer, true)
}

func RemoveRelayer(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &RelayerListParam{}
	if err := utils.UnpackMethod(ABI, MethodRemoveRelayer, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodRemoveRelayer, true)
}

func ApproveRemoveRelayer(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ApproveRelayerParam{}
	if err := utils.UnpackMethod(ABI, MethodApproveRemoveRelayer, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodApproveRemoveRelayer, true)
}
