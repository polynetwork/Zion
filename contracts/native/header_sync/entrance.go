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
package header_sync

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	hscommon "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "header sync"

const (
	MethodContractName      = "name"
	MethodSyncGenesisHeader = "syncGenesisHeader"
	MethodSyncBlockHeader   = "syncBlockHeader"
	MethodSyncCrossChainMsg = "syncCrossChainMsg"
)

var (
	this     = native.NativeContractAddrMap[native.NativeSyncHeader]
	gasTable = map[string]uint64{
		MethodContractName:      0,
		MethodSyncGenesisHeader: 0,
		MethodSyncBlockHeader:   100000,
		MethodSyncCrossChainMsg: 0,
	}

	ABI *abi.ABI
)

func InitHeaderSync() {
	ABI = GetABI()
	native.Contracts[this] = RegisterHeaderSyncContract
}

func RegisterHeaderSyncContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodContractName, Name)
	s.Register(MethodSyncGenesisHeader, SyncGenesisHeader)
	s.Register(MethodSyncBlockHeader, SyncBlockHeader)
	s.Register(MethodSyncCrossChainMsg, SyncCrossChainMsg)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(ABI, MethodContractName, contractName)
}

func SyncGenesisHeader(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &hscommon.SyncGenesisHeaderParam{}
	if err := utils.UnpackMethod(ABI, MethodSyncGenesisHeader, params, ctx.Payload); err != nil {
		return nil, err
	}
	chainID := params.ChainID

	//check if chainid exist
	sideChain, err := side_chain_manager.GetSideChain(native, chainID)
	if err != nil {
		return nil, fmt.Errorf("SyncGenesisHeader, side_chain_manager.GetSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("SyncGenesisHeader, side chain is not registered")
	}

	handler, err := GetChainHandler(sideChain.Router)
	if err != nil {
		return nil, err
	}

	err = handler.SyncGenesisHeader(native)
	if err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodSyncGenesisHeader, true)
}

func SyncBlockHeader(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &hscommon.SyncBlockHeaderParam{}
	if err := utils.UnpackMethod(ABI, MethodSyncBlockHeader, params, ctx.Payload); err != nil {
		return nil, err
	}

	chainID := params.ChainID
	//check if chainid exist
	sideChain, err := side_chain_manager.GetSideChain(native, chainID)
	if err != nil {
		return nil, fmt.Errorf("SyncGenesisHeader, side_chain_manager.GetSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("SyncGenesisHeader, side chain is not registered")
	}

	handler, err := GetChainHandler(sideChain.Router)
	if err != nil {
		return nil, err
	}

	err = handler.SyncBlockHeader(native)
	if err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodSyncBlockHeader, true)
}

func SyncCrossChainMsg(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &hscommon.SyncCrossChainMsgParam{}
	if err := utils.UnpackMethod(ABI, MethodSyncCrossChainMsg, params, ctx.Payload); err != nil {
		return nil, err
	}

	chainID := params.ChainID
	//check if chainid exist
	sideChain, err := side_chain_manager.GetSideChain(native, chainID)
	if err != nil {
		return nil, fmt.Errorf("SyncGenesisHeader, side_chain_manager.GetSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("SyncGenesisHeader, side chain is not registered")
	}

	handler, err := GetChainHandler(sideChain.Router)
	if err != nil {
		return nil, err
	}

	err = handler.SyncCrossChainMsg(native)
	if err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodSyncCrossChainMsg, true)
}

func GetChainHandler(router uint64) (hscommon.HeaderSyncHandler, error) {
	return nil, nil
}
