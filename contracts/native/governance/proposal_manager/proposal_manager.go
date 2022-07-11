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

package proposal_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/contracts/native"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/proposal_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"math/big"
)

const (
	UPDATE_NODE_MANAGER_GLOBAL_CONFIG_EVENT = "UpdateNodeManagerGlobalConfig"
)

var (
	gasTable = map[string]uint64{
		MethodUpdateNodeManagerGlobalConfig: 0,
	}
)

func InitProposalManager() {
	InitABI()
	native.Contracts[this] = RegisterProposalManagerContract
}

func RegisterProposalManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodUpdateNodeManagerGlobalConfig, UpdateNodeManagerGlobalConfig)
}

func UpdateNodeManagerGlobalConfig(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller
	height := s.ContractRef().BlockHeight()
	var ok bool

	params := &UpdateNodeManagerGlobalConfigParam{}
	if err := utils.UnpackMethod(ABI, MethodUpdateNodeManagerGlobalConfig, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("UpdateNodeManagerGlobalConfig, unpack params error: %v", err)
	}
	if height.Cmp(new(big.Int).SetUint64(params.ExpireHeight)) <= 0 {
		return nil, fmt.Errorf("UpdateNodeManagerGlobalConfig, proposal is expired")
	}

	success, err := node_manager.CheckConsensusSigns(s, MethodUpdateNodeManagerGlobalConfig, ctx.Payload, caller)
	if err != nil {
		return nil, fmt.Errorf("UpdateNodeManagerGlobalConfig, node_manager.CheckConsensusSigns error: %v", err)
	}
	if success {
		globalConfig, err := node_manager.GetGlobalConfigImpl(s)
		if err != nil {
			return nil, fmt.Errorf("UpdateNodeManagerGlobalConfig, node_manager.GetGlobalConfig error: %v", err)
		}
		if params.ConsensusValidatorNum >= node_manager.GenesisConsensusValidatorNum {
			globalConfig.ConsensusValidatorNum = params.ConsensusValidatorNum
		}
		if params.VoterValidatorNum >= node_manager.GenesisVoterValidatorNum {
			globalConfig.VoterValidatorNum = params.VoterValidatorNum
		}
		if params.BlockPerEpoch != 0 {
			globalConfig.BlockPerEpoch = new(big.Int).SetUint64(params.BlockPerEpoch)
		}
		if params.MaxDescLength != 0 {
			globalConfig.MaxDescLength = params.MaxDescLength
		}
		if params.MaxCommissionChange != "" {
			globalConfig.MaxCommissionChange, ok = new(big.Int).SetString(params.MaxCommissionChange, 10)
			if !ok {
				return nil, fmt.Errorf("UpdateNodeManagerGlobalConfig, MaxCommissionChange param error")
			}
		}
		if params.MinInitialStake != "" {
			globalConfig.MinInitialStake, ok = new(big.Int).SetString(params.MinInitialStake, 10)
			if !ok {
				return nil, fmt.Errorf("UpdateNodeManagerGlobalConfig, MinInitialStake param error")
			}
		}
		err = node_manager.SetGlobalConfig(s, globalConfig)
		if err != nil {
			return nil, fmt.Errorf("UpdateNodeManagerGlobalConfig, node_manager.SetGlobalConfig error: %v", err)
		}
		err = s.AddNotify(ABI, []string{UPDATE_NODE_MANAGER_GLOBAL_CONFIG_EVENT})
		if err != nil {
			return nil, fmt.Errorf("UpdateNodeManagerGlobalConfig, AddNotify error: %v", err)
		}
	}
	return utils.PackOutputs(ABI, MethodUpdateNodeManagerGlobalConfig, true)
}
