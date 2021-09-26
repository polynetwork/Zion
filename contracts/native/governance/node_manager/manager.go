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
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

var (
	gasTable = map[string]uint64{
		MethodContractName: 0,
		MethodPropose:      100000000,
		MethodVote:         100000,
		MethodEpoch:        0,
	}
)

const (
	MinEpochValidPeriod     uint64 = 100
	DefaultEpochValidPeriod uint64 = 86400
	MaxEpochValidPeriod     uint64 = 86400 * 10

	MinProposalPeersLen int = 4   // F = 1, n >= 3f + 1
	MaxProposalPeersLen int = 100 // F = 33
)

func InitNodeManager() {
	ABI = GetABI()
	native.Contracts[this] = RegisterNodeManagerContract
}

func RegisterNodeManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodContractName, Name)
	s.Register(MethodPropose, Propose)
	s.Register(MethodVote, Vote)
	s.Register(MethodEpoch, Epoch)
	s.Register(MethodProof, EpochProof)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return new(MethodContractNameOutput).Encode()
}

// todo(fuk): hash sum all validator addresses and block height after voted success
func Propose(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	height := s.ContractRef().BlockHeight().Uint64()

	// check authority
	curEpoch, err := getCurEpoch(s.GetCacheDB())
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("propose, get current epoch info err: %v", err)
	}
	if err := checkAuthority(s.ContractRef().TxOrigin(), curEpoch); err != nil {
		return utils.ByteFailed, fmt.Errorf("propose, propose authority err: %v", err)
	}

	// decode input
	input := new(MethodProposeInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("propose, propose input decode err: %v", err)
	}
	peers := input.Peers
	startHeight := input.StartHeight

	// check peers, try to match all peer's public key and address
	if peers == nil || peers.List == nil {
		return utils.ByteFailed, fmt.Errorf("propose, proposal peers invalid")
	}
	if len(peers.List) < MinProposalPeersLen || len(peers.List) > MaxProposalPeersLen {
		return utils.ByteFailed, fmt.Errorf("propose, peers length should be in range of [%d, %d]",
			MinProposalPeersLen, MaxProposalPeersLen)
	}
	for _, peer := range peers.List {
		if err := checkPeer(peer); err != nil {
			return utils.ByteFailed, fmt.Errorf("propose, check proposal peers err: %v", err)
		}
	}

	// check peers, number for proposal's peers should be at least 2/3 of old members
	curMembers := curEpoch.Members()
	if curMembers == nil {
		return utils.ByteFailed, fmt.Errorf("propose, get current epoch members err: %v", err)
	}
	oldMemberSize := 0
	for _, peer := range peers.List {
		if _, ok := curMembers[peer.Address]; ok {
			oldMemberSize += 1
		}
	}
	if 3*oldMemberSize < 2*len(curMembers) {
		return utils.ByteFailed, fmt.Errorf("propose, proposal peers should be at least 2/3 old members")
	}

	// proposal start height should be in range of [height + minEpochValidPeriod, height + maxEpochValidPeriod]
	if startHeight > 0 {
		latestStartHeight := height + MinEpochValidPeriod
		farawayStartHeight := height + MaxEpochValidPeriod
		if startHeight < latestStartHeight || startHeight > farawayStartHeight {
			return utils.ByteFailed, fmt.Errorf("propose, proposal start height should be in range of [%d,  %d]",
				latestStartHeight, farawayStartHeight)
		}
	} else {
		startHeight = height + DefaultEpochValidPeriod
	}

	// sort and store proposal
	epoch := &EpochInfo{
		ID:          curEpoch.ID + 1,
		Peers:       peers,
		StartHeight: startHeight,
	}
	sort.Sort(epoch.Peers)
	if findProposal(s, epoch.Hash()) {
		return utils.ByteFailed, fmt.Errorf("propose, proposal %s already exist", epoch.Hash().Hex())
	}

	if err := storeEpoch(s, epoch); err != nil {
		return utils.ByteFailed, fmt.Errorf("propose, store epoch failed, err: %v", err)
	}
	if err := storeProposal(s, epoch.Hash()); err != nil {
		return utils.ByteFailed, fmt.Errorf("propose, store proposal failed, err: %v", err)
	}
	// vote to self proposal
	if err := storeVote(s, epoch.Hash(), s.ContractRef().TxOrigin()); err != nil {
		return utils.ByteFailed, fmt.Errorf("propose, proposer vote to self proposal err: %v", err)
	}

	// emit event log
	if err := emitEventProposed(s, epoch); err != nil {
		return utils.ByteFailed, fmt.Errorf("propose, emit proposed event log err: %v", err)
	}

	return utils.ByteSuccess, nil
}

func Vote(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	// check authority
	voter := s.ContractRef().TxOrigin()
	curEpoch, err := getCurEpoch(s.GetCacheDB())
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("vote, get current epoch info err: %v", err)
	}
	if err := checkAuthority(voter, curEpoch); err != nil {
		return utils.ByteFailed, fmt.Errorf("vote, propose authority err: %v", err)
	}

	// decode and check epoch info
	input := new(MethodVoteInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("vote, decode input failed, err: %v", err)
	}
	if expectEpochID := curEpoch.ID + 1; input.EpochID != expectEpochID {
		return utils.ByteFailed, fmt.Errorf("vote, epoch id expect %d, got %d", expectEpochID, curEpoch.ID)
	}
	proposal := input.Hash
	if !findProposal(s, proposal) {
		return utils.ByteFailed, fmt.Errorf("vote, can not find proposal %s", proposal.Hex())
	}
	epoch, err := getEpoch(s, proposal)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("vote, can not find proposal %s epoch info", proposal.Hex())
	}
	if input.EpochID != epoch.ID {
		return utils.ByteFailed, fmt.Errorf("vote, vote epoch id %d not match proposal epoch id %d", input.EpochID, epoch.ID)
	}
	if proposal != epoch.Hash() {
		return utils.ByteFailed, fmt.Errorf("vote, vote proposal hash %s not match epoch proposal %s", proposal.Hex(), epoch.Hash().Hex())
	}

	// forbid duplicate vote
	if findVote(s, proposal, voter) {
		return utils.ByteFailed, fmt.Errorf("vote, validator %s already voted to proposal %s", voter.Hex(), proposal.Hex())
	}

	// already reach quorum size
	sizeBeforeVote := getVoteSize(s, proposal)
	if sizeBeforeVote > curEpoch.QuorumSize() {
		return utils.ByteSuccess, nil
	}

	// store vote
	if err := storeVote(s, proposal, voter); err != nil {
		return utils.ByteFailed, fmt.Errorf("vote, store proposal %s vote failed, err: %v", proposal.Hex(), err)
	}
	sizeAfterVote := getVoteSize(s, proposal)
	groupSize := len(curEpoch.Members())
	if err := emitEventVoted(s, input.EpochID, proposal, sizeAfterVote, groupSize); err != nil {
		return utils.ByteFailed, fmt.Errorf("vote, emit event voted failed, err: %v", err)
	}

	// check point:
	// 1. store current epoch
	// 2. store current epoch proof
	// 3. emit event log
	if sizeAfterVote == curEpoch.QuorumSize() {
		storeCurrentEpochHash(s, epoch.Hash())
		storeEpochProof(s, epoch.ID, epoch.Hash())
		if err := emitEpochChange(s, curEpoch, epoch); err != nil {
			return utils.ByteFailed, fmt.Errorf("vote, emit epoch changed failed, err: %v", err)
		}
	}

	return utils.ByteSuccess, nil
}

func Epoch(s *native.NativeContract) ([]byte, error) {
	hash, err := getCurrentEpochHash(s)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("epoch, get current epoch hash err: %v", err)
	}
	epoch, err := getEpoch(s, hash)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("epoch, get current epoch info err: %v", err)
	}

	output := &MethodEpochOutput{Epoch: epoch}
	return output.Encode()
}

func EpochProof(s *native.NativeContract) ([]byte, error) {
	hash, err := getCurrentEpochHash(s)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("proof, get current epoch hash err: %v", err)
	}
	epoch, err := getEpoch(s, hash)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("proof, get current epoch info err: %v", err)
	}
	proof, err := getEpochProof(s, epoch.ID)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("proof, get current epoch proof err: %v", err)
	}
	output := &MethodProofOutput{Hash: proof}
	return output.Encode()
}

func CheckConsensusSigns(s *native.NativeContract, method string, input []byte, address common.Address) (bool, error) {
	epochHash, err := getCurrentEpochHash(s)
	if err != nil {
		return false, fmt.Errorf("checkConsensusSign, get current epoch hash failed, err: %v", err)
	}
	epoch, err := getEpoch(s, epochHash)
	if err != nil {
		return false, fmt.Errorf("checkConsensusSign, get current epoch info failed, err: %v", err)
	}
	if err := checkAuthority(address, epoch); err != nil {
		return false, fmt.Errorf("checkConsensusSign, check authority failed, err: %v", err)
	}

	sign := &ConsensusSign{Method: method, Input: input}
	if exist, err := getSign(s, sign.Hash()); err != nil {
		return false, fmt.Errorf("checkConsensusSign, get sign failed, err: %v", err)
	} else if exist == nil {
		if err := storeSign(s, sign); err != nil {
			return false, fmt.Errorf("checkConsensusSign, store sign failed, err: %v", err)
		}
	} else if exist.Hash() != sign.Hash() {
		return false, fmt.Errorf("checkConsensusSign, expect sign hash %s got %s", exist.Hash().Hex(), sign.Hash().Hex())
	}

	if findSigner(s, sign.Hash(), address) {
		return false, fmt.Errorf("checkConsensusSign, signer %s already exist", address.Hex())
	}
	if err := storeSigner(s, sign.Hash(), address); err != nil {
		return false, fmt.Errorf("checkConsensusSign, store signers failed, hash %s, err: %v", sign.Hash().Hex(), err)
	}
	if size := getSignerSize(s, sign.Hash()); size >= epoch.QuorumSize() {
		delSign(s, sign.Hash())
		//todo: emit
	} else {
		// todo: emit event log
	}

	return true, nil
}
