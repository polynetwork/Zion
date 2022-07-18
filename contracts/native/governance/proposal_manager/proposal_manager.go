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
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/proposal_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

const (
	PROPOSE_EVENT              = "Propose"
	VOTE_ACTIVE_PROPOSAL_EVENT = "VoteActiveProposal"
)

var (
	gasTable = map[string]uint64{
		MethodPropose:            0,
		MethodSetActiveProposal:  0,
		MethodVoteActiveProposal: 0,
		MethodGetActiveProposal:  0,
		MethodGetProposalList:    0,
	}
)

func InitProposalManager() {
	InitABI()
	native.Contracts[this] = RegisterProposalManagerContract
}

func RegisterProposalManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodPropose, Propose)
	s.Register(MethodSetActiveProposal, SetActiveProposal)
	s.Register(MethodVoteActiveProposal, VoteActiveProposal)
	s.Register(MethodGetActiveProposal, GetActiveProposal)
	s.Register(MethodGetProposalList, GetProposalList)
}

func Propose(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	height := s.ContractRef().BlockHeight()

	params := &ProposeParam{}
	if err := utils.UnpackMethod(ABI, MethodPropose, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("Propose, unpack params error: %v", err)
	}

	globalConfig, err := node_manager.GetGlobalConfigImpl(s)
	if err != nil {
		return nil, fmt.Errorf("Propose, GetGlobalConfigImpl error: %v", err)
	}

	proposalID, err := getProposalID(s)
	if err != nil {
		return nil, fmt.Errorf("Propose, getProposalID error: %v", err)
	}
	proposalList, err := getProposalList(s)
	if err != nil {
		return nil, fmt.Errorf("Propose, getProposalList error: %v", err)
	}
	proposal := &Proposal{
		ID:      proposalID,
		Address: ctx.Caller,
		PType:   params.PType,
		Content: params.Content,
		Stake:   params.Stake,
	}
	if len(proposalList.ProposalList) == 0 {
		proposal.Status = Active
		proposal.EndHeight = new(big.Int).Add(height, globalConfig.BlockPerEpoch)
	}
	proposalList.ProposalList = append(proposalList.ProposalList, proposal)
	err = setProposalList(s, proposalList)
	if err != nil {
		return nil, fmt.Errorf("Propose, setProposalList error: %v", err)
	}
	setProposalID(s, new(big.Int).Add(proposalID, common.Big1))

	err = s.AddNotify(ABI, []string{PROPOSE_EVENT}, ctx.Caller.Hex(), params.PType, params.Stake.String(), hex.EncodeToString(params.Content))
	if err != nil {
		return nil, fmt.Errorf("Propose, AddNotify error: %v", err)
	}

	return utils.PackOutputs(ABI, MethodPropose, true)
}

func SetActiveProposal(s *native.NativeContract) ([]byte, error) {
	err := setActiveProposal(s)
	if err != nil {
		return nil, fmt.Errorf("SetActiveProposal, setActiveProposal error: %v", err)
	}

	return utils.PackOutputs(ABI, MethodSetActiveProposal, true)
}

func VoteActiveProposal(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	params := &VoteActiveProposalParam{}
	if err := utils.UnpackMethod(ABI, MethodVoteActiveProposal, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("VoteActiveProposal, unpack params error: %v", err)
	}

	activeProposal, flag, err := getActiveProposal(s)
	if err != nil {
		return nil, fmt.Errorf("VoteActiveProposal, getActiveProposal error: %v", err)
	}
	if activeProposal.ID != params.ID {
		return nil, fmt.Errorf("VoteActiveProposal, ID is not active")
	}

	// need to update proposal list
	if flag {
		err = setActiveProposal(s)
		if err != nil {
			return nil, fmt.Errorf("VoteActiveProposal, setActiveProposal error: %v", err)
		}
	}

	success, err := node_manager.CheckConsensusSigns(s, MethodVoteActiveProposal, ctx.Payload, caller, node_manager.Proposer)
	if err != nil {
		return nil, fmt.Errorf("VoteActiveProposal, node_manager.CheckConsensusSigns error: %v", err)
	}
	if success {
		err = setActiveProposal(s)
		if err != nil {
			return nil, fmt.Errorf("VoteActiveProposal, success, setActiveProposal error: %v", err)
		}

		if activeProposal.PType == UpdateGlobalConfig {
			config := new(node_manager.GlobalConfig)
			err := rlp.DecodeBytes(activeProposal.Content, config)
			if err != nil {
				return nil, fmt.Errorf("VoteActiveProposal, deserialize global config error: %v", err)
			}

			globalConfig, err := node_manager.GetGlobalConfigImpl(s)
			if err != nil {
				return nil, fmt.Errorf("VoteActiveProposal, node_manager.GetGlobalConfigImpl error: %v", err)
			}
			if config.ConsensusValidatorNum >= node_manager.GenesisConsensusValidatorNum {
				globalConfig.ConsensusValidatorNum = config.ConsensusValidatorNum
			}
			if config.VoterValidatorNum >= node_manager.GenesisVoterValidatorNum {
				globalConfig.VoterValidatorNum = config.VoterValidatorNum
			}
			if config.BlockPerEpoch != nil {
				globalConfig.BlockPerEpoch = config.BlockPerEpoch
			}
			if config.MaxDescLength != 0 {
				globalConfig.MaxDescLength = config.MaxDescLength
			}
			if config.MaxCommissionChange != nil {
				globalConfig.MaxCommissionChange = config.MaxCommissionChange
			}
			if config.MinInitialStake != nil {
				globalConfig.MinInitialStake = config.MinInitialStake
			}
			err = node_manager.SetGlobalConfig(s, globalConfig)
			if err != nil {
				return nil, fmt.Errorf("VoteActiveProposal, node_manager.SetGlobalConfig error: %v", err)
			}
		}

		err = s.AddNotify(ABI, []string{VOTE_ACTIVE_PROPOSAL_EVENT}, activeProposal.ID.String(), activeProposal.PType)
		if err != nil {
			return nil, fmt.Errorf("VoteActiveProposal, AddNotify error: %v", err)
		}
	}
	return utils.PackOutputs(ABI, MethodVoteActiveProposal, true)
}

func GetActiveProposal(s *native.NativeContract) ([]byte, error) {
	activeProposal, _, err := getActiveProposal(s)
	if err != nil {
		return nil, fmt.Errorf("GetActiveProposal, getActiveProposal error: %v", err)
	}

	enc, err := rlp.EncodeToBytes(activeProposal)
	if err != nil {
		return nil, fmt.Errorf("GetActiveProposal, serialize active proposal error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetActiveProposal, enc)
}

func GetProposalList(s *native.NativeContract) ([]byte, error) {
	proposalList, err := getProposalList(s)
	if err != nil {
		return nil, fmt.Errorf("GetProposalList, getProposalList error: %v", err)
	}

	enc, err := rlp.EncodeToBytes(proposalList)
	if err != nil {
		return nil, fmt.Errorf("GetProposalList, serialize active proposal error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetProposalList, enc)
}
