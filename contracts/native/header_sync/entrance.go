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

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/bsc"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/btc"
	hscommon "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/eth"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/heco"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/msc"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/polygon"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/quorum"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "header sync"

var (
	this = native.NativeContractAddrMap[native.NativeSyncHeader]
)

func InitHeaderSync() {
	native.Contracts[this] = RegisterHeaderSyncContract
	hscommon.ABI = hscommon.GetABI()
}

func RegisterHeaderSyncContract(s *native.NativeContract) {
	s.Prepare(hscommon.ABI, hscommon.GasTable)

	s.Register(hscommon.MethodContractName, Name)
	s.Register(hscommon.MethodSyncGenesisHeader, SyncGenesisHeader)
	s.Register(hscommon.MethodSyncBlockHeader, SyncBlockHeader)
	s.Register(hscommon.MethodSyncCrossChainMsg, SyncCrossChainMsg)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(hscommon.ABI, hscommon.MethodContractName, contractName)
}

func SyncGenesisHeader(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &hscommon.SyncGenesisHeaderParam{}
	if err := utils.UnpackMethod(hscommon.ABI, hscommon.MethodSyncGenesisHeader, params, ctx.Payload); err != nil {
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

	return utils.PackOutputs(hscommon.ABI, hscommon.MethodSyncGenesisHeader, true)
}

func SyncBlockHeader(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &hscommon.SyncBlockHeaderParam{}
	if err := utils.UnpackMethod(hscommon.ABI, hscommon.MethodSyncBlockHeader, params, ctx.Payload); err != nil {
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

	return utils.PackOutputs(hscommon.ABI, hscommon.MethodSyncBlockHeader, true)
}

func SyncCrossChainMsg(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &hscommon.SyncCrossChainMsgParam{}
	if err := utils.UnpackMethod(hscommon.ABI, hscommon.MethodSyncCrossChainMsg, params, ctx.Payload); err != nil {
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

	return utils.PackOutputs(hscommon.ABI, hscommon.MethodSyncCrossChainMsg, true)
}

func GetChainHandler(router uint64) (hscommon.HeaderSyncHandler, error) {
	switch router {
	case utils.BTC_ROUTER:
		return btc.NewBTCHandler(), nil
	case utils.BSC_ROUTER:
		return bsc.NewHandler(), nil
	case utils.ETH_ROUTER:
		return eth.NewETHHandler(), nil
	case utils.HECO_ROUTER:
		return heco.NewHecoHandler(), nil
	case utils.MSC_ROUTER:
		return msc.NewHandler(), nil
	case utils.QUORUM_ROUTER:
		return quorum.NewQuorumHandler(), nil
	case utils.POLYGON_HEIMDALL_ROUTER:
		return polygon.NewHeimdallHandler(), nil
	case utils.POLYGON_BOR_ROUTER:
		return polygon.NewBorHandler(), nil
	default:
		return nil, fmt.Errorf("not a supported router:%d", router)
	}
}
