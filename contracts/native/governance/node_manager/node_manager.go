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
package node_manager

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "node manager"

const (
	//status
	CandidateStatus Status = iota
	ConsensusStatus
	QuitingStatus
	BlackStatus

	//function name
	MethodContractName        = "name"
	MethodInitConfig          = "initConfig"
	MethodRegisterCandidate   = "registerCandidate"
	MethodUnRegisterCandidate = "unRegisterCandidate"
	MethodApproveCandidate    = "approveCandidate"
	MethodBlackNode           = "blackNode"
	MethodWhiteNode           = "whiteNode"
	MethodQuitNode            = "quitNode"
	MethodUpdateConfig        = "updateConfig"
	MethodCommitDpos          = "commitDpos"

	//key prefix
	GOVERNANCE_VIEW = "governanceView"
	VBFT_CONFIG     = "vbftConfig"
	CANDIDITE_INDEX = "candidateIndex"
	PEER_APPLY      = "peerApply"
	PEER_POOL       = "peerPool"
	PEER_INDEX      = "peerIndex"
	BLACK_LIST      = "blackList"
	CONSENSUS_SIGNS = "consensusSigns"

	//const
	MIN_PEER_NUM = 4
)

var (
	this     = native.NativeContractAddrMap[native.NativeNodeManager]
	gasTable = map[string]uint64{
		MethodContractName:        0,
		MethodInitConfig:          0,
		MethodRegisterCandidate:   100000,
		MethodUnRegisterCandidate: 0,
		MethodApproveCandidate:    0,
		MethodBlackNode:           0,
		MethodWhiteNode:           0,
		MethodQuitNode:            0,
		MethodUpdateConfig:        0,
		MethodCommitDpos:          0,
	}

	ABI *abi.ABI
)

func InitNodeManager() {
	ABI = GetABI()
	native.Contracts[this] = RegisterNodeManagerContract
}

func RegisterNodeManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodContractName, Name)
	s.Register(MethodInitConfig, InitConfig)
	s.Register(MethodRegisterCandidate, RegisterCandidate)
	s.Register(MethodUnRegisterCandidate, UnRegisterCandidate)
	s.Register(MethodQuitNode, QuitNode)
	s.Register(MethodApproveCandidate, ApproveCandidate)
	s.Register(MethodBlackNode, BlackNode)
	s.Register(MethodWhiteNode, WhiteNode)
	s.Register(MethodUpdateConfig, UpdateConfig)
	s.Register(MethodCommitDpos, CommitDpos)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(ABI, MethodContractName, contractName)
}

func InitConfig(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &VBFTConfig{}
	if err := utils.UnpackMethod(ABI, MethodInitConfig, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodInitConfig, true)
}

func RegisterCandidate(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &RegisterPeerParam{}
	if err := utils.UnpackMethod(ABI, MethodRegisterCandidate, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodRegisterCandidate, true)
}

func UnRegisterCandidate(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &PeerParam{}
	if err := utils.UnpackMethod(ABI, MethodUnRegisterCandidate, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodUnRegisterCandidate, true)
}

func QuitNode(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &PeerParam{}
	if err := utils.UnpackMethod(ABI, MethodQuitNode, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodQuitNode, true)
}

func ApproveCandidate(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &PeerParam{}
	if err := utils.UnpackMethod(ABI, MethodApproveCandidate, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodApproveCandidate, true)
}

func BlackNode(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &PeerListParam{}
	if err := utils.UnpackMethod(ABI, MethodBlackNode, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodBlackNode, true)
}

func WhiteNode(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &PeerParam{}
	if err := utils.UnpackMethod(ABI, MethodWhiteNode, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodWhiteNode, true)
}

func UpdateConfig(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &UpdateConfigParam{}
	if err := utils.UnpackMethod(ABI, MethodUpdateConfig, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodUpdateConfig, true)
}

func CommitDpos(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(ABI, MethodCommitDpos, true)
}
