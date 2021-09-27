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
	"github.com/ethereum/go-ethereum/log"
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
	MinProposalPeersLen     int    = 4   // F = 1, n >= 3f + 1
	MaxProposalPeersLen     int    = 100 // F = 33
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

func Propose(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	height := s.ContractRef().BlockHeight().Uint64()
	proposer := s.ContractRef().TxOrigin()
	caller := ctx.Caller

	// check authority
	curEpoch, err := getCurEpoch(s.GetCacheDB())
	if err != nil {
		log.Trace("propose", "get current epoch failed", err)
		return utils.ByteFailed, ErrEpochNotExist
	}
	if err := checkAuthority(proposer, caller, curEpoch); err != nil {
		log.Trace("propose", "check authority failed", err, "tx origin", proposer.Hex())
		return utils.ByteFailed, ErrInvalidAuthority
	}

	// decode input
	input := new(MethodProposeInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("propose", "decode input failed", err)
		return utils.ByteFailed, ErrInvalidInput
	}

	peers := input.Peers
	startHeight := input.StartHeight
	// check peers, try to match all peer's public key and address
	if peers == nil || peers.List == nil || len(peers.List) == 0 {
		log.Trace("propose", "check peers", "peer list is nil")
		return utils.ByteFailed, ErrInvalidPeers
	}
	if len(peers.List) < MinProposalPeersLen || len(peers.List) > MaxProposalPeersLen {
		log.Trace("propose", "check peers number",
			fmt.Errorf("propose, peers length should be in range of [%d, %d]",
				MinProposalPeersLen, MaxProposalPeersLen))
		return utils.ByteFailed, ErrProposalPeersOutOfRange
	}
	for _, peer := range peers.List {
		if err := checkPeer(peer); err != nil {
			log.Trace("propose", "check peer public key", "public key not match address")
			return utils.ByteFailed, ErrInvalidPubKey
		}
	}

	// check peers, number for proposal's peers should be at least 2/3 of old members
	curMembers := curEpoch.Members()
	if curMembers == nil {
		log.Trace("propose", "check current epoch members", "members is nil")
		return utils.ByteFailed, ErrInvalidEpoch
	}
	oldMemberSize := 0
	for _, peer := range peers.List {
		if _, ok := curMembers[peer.Address]; ok {
			oldMemberSize += 1
		}
	}
	if 3*oldMemberSize < 2*len(curMembers) {
		log.Trace("propose", "check old members", "proposal peers should be at least 2/3 old members")
		return utils.ByteFailed, ErrOldParticipantsNumber
	}

	// proposal start height should be in range of [height + minEpochValidPeriod, height + maxEpochValidPeriod]
	if startHeight > 0 {
		latestStartHeight := height + MinEpochValidPeriod
		farawayStartHeight := height + MaxEpochValidPeriod
		if startHeight < latestStartHeight || startHeight > farawayStartHeight {
			log.Trace("propose", "check start height", fmt.Errorf("propose, proposal start height should be in range of [%d,  %d]",
				latestStartHeight, farawayStartHeight))
			return utils.ByteFailed, ErrProposalStartHeight
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
		log.Trace("propose", "check proposal hash, dump proposal", epoch.Hash().Hex())
		return utils.ByteFailed, ErrDuplicateProposal
	}

	if err := storeEpoch(s, epoch); err != nil {
		log.Trace("propose", "store epoch failed", err)
		return utils.ByteFailed, ErrStorage
	}
	if err := storeProposal(s, epoch.Hash()); err != nil {
		log.Trace("propose", "store proposal hash failed", err)
		return utils.ByteFailed, ErrStorage
	}
	// vote to self proposal
	if err := storeVote(s, epoch.Hash(), s.ContractRef().TxOrigin()); err != nil {
		log.Trace("propose", "store vote failed", err)
		return utils.ByteFailed, ErrStorage
	}

	// emit event log
	if err := emitEventProposed(s, epoch); err != nil {
		log.Trace("propose", "emit event log failed", err)
		return utils.ByteFailed, ErrEmitLog
	}

	return utils.ByteSuccess, nil
}

func Vote(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	voter := s.ContractRef().TxOrigin()
	caller := ctx.Caller

	// check authority
	curEpoch, err := getCurEpoch(s.GetCacheDB())
	if err != nil {
		log.Trace("vote", "get current epoch failed", err)
		return utils.ByteFailed, ErrEpochNotExist
	}
	if err := checkAuthority(voter, caller, curEpoch); err != nil {
		log.Trace("vote", "check authority failed", err, "voter", voter.Hex())
		return utils.ByteFailed, ErrInvalidAuthority
	}

	// decode and check epoch info
	input := new(MethodVoteInput)
	if err := input.Decode(ctx.Payload); err != nil {
		log.Trace("vote", "decode input failed", err)
		return utils.ByteFailed, ErrInvalidInput
	}
	if expectEpochID := curEpoch.ID + 1; input.EpochID != expectEpochID {
		log.Trace("vote", "check epoch ID failed, expect", expectEpochID, "got", curEpoch.ID)
		return utils.ByteFailed, ErrInvalidEpoch
	}
	proposal := input.Hash
	if !findProposal(s, proposal) {
		log.Trace("vote", "find proposal failed", proposal.Hex())
		return utils.ByteFailed, ErrProposalNotExist
	}
	epoch, err := getEpoch(s, proposal)
	if err != nil {
		log.Trace("vote", "get epoch failed", proposal.Hex())
		return utils.ByteFailed, ErrEpochNotExist
	}
	if input.EpochID != epoch.ID {
		log.Trace("vote", "check epoch id failed, expect", input.EpochID, "got", epoch.ID)
		return utils.ByteFailed, ErrInvalidEpoch
	}
	if proposal != epoch.Hash() {
		log.Trace("vote", "check epoch hash failed, expect", proposal.Hex(), "got", epoch.Hash().Hex())
		return utils.ByteFailed, ErrInvalidEpoch
	}

	// forbid duplicate vote
	if findVote(s, proposal, voter) {
		log.Trace("vote", "check vote", "duplicate vote", "proposal", proposal.Hex(), "vote", voter.Hex())
		return utils.ByteFailed, ErrDuplicateVote
	}

	// already reach quorum size
	sizeBeforeVote := getVoteSize(s, proposal)
	if sizeBeforeVote > curEpoch.QuorumSize() {
		log.Trace("vote", "check size", "already reach quorum size", "num", sizeBeforeVote, "quorum size", curEpoch.QuorumSize())
		return utils.ByteSuccess, nil
	}

	// store vote
	if err := storeVote(s, proposal, voter); err != nil {
		log.Trace("vote", "store vote failed", err)
		return utils.ByteFailed, ErrStorage
	}
	sizeAfterVote := getVoteSize(s, proposal)
	groupSize := len(curEpoch.Members())
	if err := emitEventVoted(s, input.EpochID, proposal, sizeAfterVote, groupSize); err != nil {
		log.Trace("vote", "emit voted log failed", err)
		return utils.ByteFailed, ErrEmitLog
	}

	// check point:
	// 1. store current epoch
	// 2. store current epoch proof
	// 3. emit event log
	if sizeAfterVote == curEpoch.QuorumSize() {
		storeCurrentEpochHash(s, epoch.Hash())
		storeEpochProof(s, epoch.ID, epoch.Hash())
		if err := emitEpochChange(s, curEpoch, epoch); err != nil {
			log.Trace("vote", "emit epoch change log failed", err)
			return utils.ByteFailed, ErrEmitLog
		}
	}

	return utils.ByteSuccess, nil
}

func Epoch(s *native.NativeContract) ([]byte, error) {
	hash, err := getCurrentEpochHash(s)
	if err != nil {
		log.Trace("epoch", "get current proposal failed", err)
		return utils.ByteFailed, ErrProposalNotExist
	}
	epoch, err := getEpoch(s, hash)
	if err != nil {
		log.Trace("epoch", "get current epoch failed", err)
		return utils.ByteFailed, ErrEpochNotExist
	}

	output := &MethodEpochOutput{Epoch: epoch}
	return output.Encode()
}

func EpochProof(s *native.NativeContract) ([]byte, error) {
	hash, err := getCurrentEpochHash(s)
	if err != nil {
		log.Trace("epoch proof", "get current proposal failed", err)
		return utils.ByteFailed, ErrProposalNotExist
	}
	epoch, err := getEpoch(s, hash)
	if err != nil {
		log.Trace("epoch proof", "get current epoch failed", err)
		return utils.ByteFailed, ErrEpochNotExist
	}
	proof, err := getEpochProof(s, epoch.ID)
	if err != nil {
		log.Trace("epoch proof", "get current epoch proof failed", err)
		return utils.ByteFailed, ErrEpochProofNotExist
	}
	output := &MethodProofOutput{Hash: proof}
	return output.Encode()
}

func CheckConsensusSigns(s *native.NativeContract, method string, input []byte, signer common.Address) (bool, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	// get epoch info
	epochHash, err := getCurrentEpochHash(s)
	if err != nil {
		log.Trace("checkConsensusSign", "get current proposal failed", err)
		return false, ErrProposalNotExist
	}
	epoch, err := getEpoch(s, epochHash)
	if err != nil {
		log.Trace("checkConsensusSign", "get current epoch failed", err)
		return false, ErrEpochNotExist
	}

	// check authority
	if err := checkAuthority(signer, caller, epoch); err != nil {
		log.Trace("checkConsensusSign", "check authority failed", err)
		return false, ErrInvalidAuthority
	}

	// get or set consensus sign info
	sign := &ConsensusSign{Method: method, Input: input}
	if exist, err := getSign(s, sign.Hash()); err != nil {
		if err.Error() == "EOF" {
			if err := storeSign(s, sign); err != nil {
				log.Trace("checkConsensusSign", "store sign failed", err, "hash", sign.Hash().Hex())
				return false, ErrStorage
			}
		} else {
			log.Trace("checkConsensusSign", "get sign failed", err, "hash", sign.Hash().Hex())
			return false, ErrConsensusSignNotExist
		}
	} else if exist.Hash() != sign.Hash() {
		log.Trace("checkConsensusSign", "check sign hash failed, expect", exist.Hash().Hex(), "got", sign.Hash().Hex())
		return false, ErrInvalidSign
	}

	// check duplicate signature
	if findSigner(s, sign.Hash(), signer) {
		log.Trace("checkConsensusSign", "signer already exist", signer.Hex(), "hash", sign.Hash().Hex())
		return false, ErrDuplicateSigner
	}

	// do not store redundancy sign
	sizeBeforeSign := getSignerSize(s, sign.Hash())
	if sizeBeforeSign >= epoch.QuorumSize() {
		return true, nil
	}

	// store signer address and emit event log
	if err := storeSigner(s, sign.Hash(), signer); err != nil {
		log.Trace("checkConsensusSign", "store signer failed", err, "hash", sign.Hash().Hex())
		return false, ErrStorage
	}
	sizeAfterSign := getSignerSize(s, sign.Hash())
	if err := emitConsensusSign(s, sign, signer, sizeAfterSign); err != nil {
		log.Trace("checkConsensusSign", "emit consensus sign log failed", err, "hash", sign.Hash().Hex())
		return false, ErrEmitLog
	}

	return sizeAfterSign >= epoch.QuorumSize(), nil
}
