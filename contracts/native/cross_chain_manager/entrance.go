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
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native/contract"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/btc"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"

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

	BLACKED_CHAIN = "BlackedChain"
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
)

func InitCrossChainManager() {
	native.Contracts[this] = RegisterCrossChainManagerContract
}

func RegisterCrossChainManagerContract(s *native.NativeContract) {
	s.Prepare(scom.ABI, gasTable)

	s.Register(MethodContractName, Name)
	s.Register(MethodImportOuterTransfer, ImportOuterTransfer)
	s.Register(MethodMultiSign, MultiSign)
	s.Register(MethodBlackChain, BlackChain)
	s.Register(MethodWhiteChain, WhiteChain)
}

func GetChainHandler(router uint64) (scom.ChainHandler, error) {
	return nil, nil
}

func Name(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(scom.ABI, MethodContractName, contractName)
}

func ImportOuterTransfer(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &scom.EntranceParam{}
	if err := utils.UnpackMethod(scom.ABI, MethodImportOuterTransfer, params, ctx.Payload); err != nil {
		return nil, err
	}

	chainID := params.SourceChainID
	blacked, err := CheckIfChainBlacked(native, chainID)
	if err != nil {
		return nil, fmt.Errorf("ImportExTransfer, CheckIfChainBlacked error: %v", err)
	}
	if blacked {
		return nil, fmt.Errorf("ImportExTransfer, source chain is blacked")
	}

	//check if chainid exist
	sideChain, err := side_chain_manager.GetSideChain(native, chainID)
	if err != nil {
		return nil, fmt.Errorf("ImportExTransfer, side_chain_manager.GetSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("ImportExTransfer, side chain %d is not registered", chainID)
	}

	handler, err := GetChainHandler(sideChain.Router)
	if err != nil {
		return nil, err
	}
	//1. verify tx
	txParam, err := handler.MakeDepositProposal(native)
	if err != nil {
		return nil, err
	}

	//2. make target chain tx
	targetid := txParam.ToChainID
	blacked, err = CheckIfChainBlacked(native, targetid)
	if err != nil {
		return nil, fmt.Errorf("ImportExTransfer, CheckIfChainBlacked error: %v", err)
	}
	if blacked {
		return nil, fmt.Errorf("ImportExTransfer, target chain is blacked")
	}

	//check if chainid exist
	sideChain, err = side_chain_manager.GetSideChain(native, targetid)
	if err != nil {
		return nil, fmt.Errorf("ImportExTransfer, side_chain_manager.GetSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("ImportExTransfer, side chain %d is not registered", targetid)
	}
	if sideChain.Router == utils.BTC_ROUTER {
		err := btc.NewBTCHandler().MakeTransaction(native, txParam, chainID)
		if err != nil {
			return nil, err
		}
		return utils.PackOutputs(scom.ABI, MethodImportOuterTransfer, true)
	}

	return utils.PackOutputs(scom.ABI, MethodImportOuterTransfer, true)
}

func MultiSign(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &scom.MultiSignParam{}
	if err := utils.UnpackMethod(scom.ABI, MethodMultiSign, params, ctx.Payload); err != nil {
		return nil, err
	}

	handler := btc.NewBTCHandler()

	//1. multi sign
	err := handler.MultiSign(native, params)
	if err != nil {
		return nil, fmt.Errorf("MultiSign fail:%v", err)
	}

	return utils.PackOutputs(scom.ABI, MethodMultiSign, true)
}

func BlackChain(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &scom.BlackChainParam{}
	if err := utils.UnpackMethod(scom.ABI, MethodBlackChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	// Get current epoch operator
	operatorAddress, err := node_manager.GetCurConOperator(native)
	if err != nil {
		return nil, fmt.Errorf("BlackChain, get current consensus operator address error: %v", err)
	}

	//check witness
	err = contract.ValidateOwner(native, operatorAddress)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, checkWitness error: %v", err)
	}

	PutBlackChain(native, params.ChainID)
	return utils.PackOutputs(scom.ABI, MethodBlackChain, true)
}

func WhiteChain(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &scom.BlackChainParam{}
	if err := utils.UnpackMethod(scom.ABI, MethodWhiteChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	// Get current epoch operator
	operatorAddress, err := node_manager.GetCurConOperator(native)
	if err != nil {
		return nil, fmt.Errorf("BlackChain, get current consensus operator address error: %v", err)
	}

	//check witness
	err = contract.ValidateOwner(native, operatorAddress)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, checkWitness error: %v", err)
	}

	RemoveBlackChain(native, params.ChainID)
	return utils.PackOutputs(scom.ABI, MethodWhiteChain, true)
}
