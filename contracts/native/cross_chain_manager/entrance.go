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
package cross_chain_manager

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "cross chain manager"

const (
	MethodContractName        = "name"
	MethodImportOuterTransfer = "importOuterTransfer"
	MethodMultiSign           = "MultiSign"
	MethodBlackChain          = "BlackChain"
	MethodWhiteChain          = "WhiteChain"
)

var (
	this     = native.NativeContractAddrMap[native.NativeCrossChain]
	gasTable = map[string]uint64{
		MethodContractName:        0,
		MethodImportOuterTransfer: 0,
		MethodMultiSign:           100000,
		MethodBlackChain:          0,
		MethodWhiteChain:          0,
	}

	ABI *abi.ABI
)

func InitCrossChainManager() {
	ABI = GetABI()
	native.Contracts[this] = RegisterCrossChainManagerContract
}

func RegisterCrossChainManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodContractName, Name)
	s.Register(MethodImportOuterTransfer, ImportOuterTransfer)
	s.Register(MethodMultiSign, MultiSign)
	s.Register(MethodBlackChain, BlackChain)
	s.Register(MethodWhiteChain, WhiteChain)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(ABI, MethodContractName, contractName)
}

func ImportOuterTransfer(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &EntranceParam{}
	if err := utils.UnpackMethod(ABI, MethodImportOuterTransfer, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodImportOuterTransfer, true)
}

func MultiSign(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &MultiSignParam{}
	if err := utils.UnpackMethod(ABI, MethodMultiSign, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodMultiSign, true)
}

func BlackChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &BlackChainParam{}
	if err := utils.UnpackMethod(ABI, MethodBlackChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodBlackChain, true)
}

func WhiteChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &BlackChainParam{}
	if err := utils.UnpackMethod(ABI, MethodWhiteChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodWhiteChain, true)
}
