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

package info_sync

import (
	"fmt"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	iscommon "github.com/ethereum/go-ethereum/contracts/native/info_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/info_sync/consensus_vote"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
)

const contractName = "cross chain info sync"

var (
	this = native.NativeContractAddrMap[native.NativeSyncCrossChainInfo]
)

func InitInfoSync() {
	native.Contracts[this] = RegisterInfoSyncContract
	iscommon.ABI = iscommon.GetABI()
}

func RegisterInfoSyncContract(s *native.NativeContract) {
	s.Prepare(iscommon.ABI, iscommon.GasTable)

	s.Register(iscommon.MethodContractName, Name)
	s.Register(iscommon.MethodSyncRootInfo, SyncRootInfo)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(iscommon.ABI, iscommon.MethodContractName, contractName)
}

func SyncRootInfo(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &iscommon.SyncRootInfoParam{}
	if err := utils.UnpackMethod(iscommon.ABI, iscommon.MethodSyncRootInfo, params, ctx.Payload); err != nil {
		return nil, err
	}

	chainID := params.ChainID

	//check if chainid exist
	sideChain, err := side_chain_manager.GetSideChain(s, chainID)
	if err != nil {
		return nil, fmt.Errorf("SyncRootInfo, side_chain_manager.GetSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("SyncRootInfo, side chain is not registered")
	}

	//sync root infos
	for _, v := range params.RootInfos {
		var rootInfo *iscommon.RootInfo
		err := rlp.DecodeBytes(v, &rootInfo)
		if err != nil {
			return nil, fmt.Errorf("SyncRootInfo, decode root info error")
		}
		//use chain id, info key and value as unique id
		unique := &consensus_vote.RootInfoUnique{
			ChainID: params.ChainID,
			Height:  rootInfo.Height,
			Info:    rootInfo.Info,
		}
		blob, err := rlp.EncodeToBytes(unique)
		if err != nil {
			return nil, err
		}

		ok, err := consensus_vote.CheckConsensusSigns(s, blob)
		if err != nil {
			return nil, fmt.Errorf("SyncRootInfo, CheckConsensusSigns error: %v", err)
		}
		if ok {
			err := iscommon.PutRootInfo(s, chainID, rootInfo.Height, rootInfo.Info)
			if err != nil {
				return nil, fmt.Errorf("SyncRootInfo, PutCrossChainInfo error: %v", err)
			}
		}
	}

	return utils.PackOutputs(iscommon.ABI, iscommon.MethodSyncRootInfo, true)
}
