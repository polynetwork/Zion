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

package core

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"time"
)

// EventDrivenEngine implement event-driven hotstuff protocol, it obtains:
// 1.validator set which represent consensus participants
type EventDrivenEngine struct {
	config *hotstuff.Config
	logger log.Logger
	epoch uint64

	addr   common.Address
	signer hotstuff.Signer
	valset hotstuff.ValidatorSet

	curRound,
	curHeight *big.Int

	requests  *requestSet
	messages  *MessagePool
	blkTree   *BlockTree
	safety    *SafetyRules
	//paceMaker *PaceMaker

	backend hotstuff.Backend

	events            *event.TypeMuxSubscription
	timeoutSub        *event.TypeMuxSubscription
	finalCommittedSub *event.TypeMuxSubscription

	// pace maker
	highestCommitRound *big.Int
	timer *time.Timer

	validateFn func([]byte, []byte) (common.Address, error)
}

func NewEventDrivenEngine(valset hotstuff.ValidatorSet) *EventDrivenEngine {
	return nil
}

// handleNewRound proposer at this round get an new proposal and broadcast to all validators.
func (e *EventDrivenEngine) handleNewRound() error {
	view := e.currentView()
	e.valset.CalcProposerByIndex(view.Round.Uint64())

	if !e.isProposer() {
		return nil
	}

	// todo: add high qc as justifyQC into pending request
	proposal := e.getCurrentPendingRequest()
	msg := &MsgProposal{
		View:     view,
		Proposal: proposal,
	}
	return e.encodeAndBroadcast(MsgTypeProposal, msg)
}

// handleProposal validate proposal info and vote to the next leader if the proposal is valid
// todo: modify err type
func (e *EventDrivenEngine) handleProposal(src hotstuff.Validator, data *hotstuff.Message) error {
	logger := e.newLogger()

	var (
		msg    *MsgProposal
		msgTyp = MsgTypeProposal
	)

	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to decode", "type", msgTyp, "err", err)
		return errFailedDecodePrepare
	}
	proposal, ok := msg.Proposal.(*types.Block)
	if !ok {
		logger.Trace("Failed to decode", "convert err", "not block")
		return errProposalConvert
	}

	justifyQC, proposalRound, err := extraProposal(proposal)
	if err != nil {
		return err
	}

	// allow the validator get into the new round before vote
	if err := e.ProcessCertificates(justifyQC); err != nil {
		return err
	}

	currentRound := e.curRound
	if currentRound.Cmp(proposalRound) != 0 {
		return errInvalidMessage
	}

	if !e.valset.IsProposer(proposal.Coinbase()) {
		return errNotFromProposer
	}

	if err := e.signer.VerifyQC(justifyQC, e.valset); err != nil {
		return err
	}
	if err := e.signer.VerifyHeader(proposal.Header(), e.valset, false); err != nil {
		return err
	}

	e.blkTree.Insert(proposal)

	//// check parent block existing
	//parentHash := justifyQC.Hash
	//parentRound := justifyQC.View.Round
	//if err := e.checkBlockExist(parentHash, parentRound); err != nil {
	//	return err
	//}
	//
	//// proposal round should be increase by 1
	//if new(big.Int).Sub(proposalRound, parentRound).Cmp(common.Big1) != 0 {
	//	return fmt.Errorf("proposal round != parent round + 1, proposalRound %v, parentRound %v", proposalRound, parentRound)
	//}
	//if !e.safety.VoteRule(proposalRound, justifyQC.View.Round) {
	//	return fmt.Errorf("voteRule failed")
	//}
	//
	//// todo: vote是否应该包含commitInfo用来证明整个3阶段都是有效的
	//vote := &Vote{
	//	Epoch:       e.paceMaker.CurrentEpoch(),
	//	Hash:        proposal.Hash(),
	//	Round:       proposalRound,
	//	ParentHash:  justifyQC.Hash,
	//	ParentRound: justifyQC.View.Round,
	//}
	vote, err := e.safety.MakeVote(proposal)
	if err != nil {
		return err
	}

	e.safety.IncreaseLastVoteRound(proposalRound)

	// todo: vote to next proposal
	return e.encodeAndBroadcast(MsgTypeVote, vote)
}

// handleVote validate vote message and try to assemble qc
func (e *EventDrivenEngine) handleVote(src hotstuff.Validator, data *hotstuff.Message) error {
	var (
		vote    *Vote
		msgType = MsgTypeVote
	)

	logger := e.newLogger()
	if err := data.Decode(&vote); err != nil {
		logger.Trace("Failed to decode", "type", msgType, "err", err)
		return errFailedDecodeNewView
	}

	// todo: first two blocks
	if vote.Hash == utils.EmptyHash || vote.ParentHash == utils.EmptyHash || vote.Round == nil || vote.ParentRound == nil {
		return fmt.Errorf("invalid vote")
	}

	if err := e.checkBlockExist(vote.Hash, vote.Round); err != nil {
		return err
	}
	if err := e.checkBlockExist(vote.ParentHash, vote.ParentRound); err != nil {
		return err
	}

	if err := e.messages.AddVote(vote.Hash, data); err != nil {
		return err
	}

	if e.messages.VoteSize(vote.Hash) < e.Q() {
		return nil
	}

	// todo: format highQC and set block tree high qc
	// todo(fuk): instance qc and broadcast to all validators
	view := e.currentView()
	qc := &hotstuff.QuorumCert{
		View:     view,
		Hash:     vote.Hash,
		Proposer: common.Address{},
		Extra:    nil,
	}

	// paceMaker send qc to next leader
	if err := e.advanceRound(qc); err != nil {
		return err
	}
	e.blkTree.UpdateHighQC(qc)

	return nil
}

//func (e *EventDrivenEngine) handleLocalTimeout() error {
//	return nil
//}
//
//func (e *EventDrivenEngine) handleRemoteTimeout() error {
//	return nil
//}
//
//func (e *EventDrivenEngine) advanceRound() error {
//	return nil
//}

// todo: add this function into handleProposal
// ProcessCertificates validate and handle QC/TC
func (e *EventDrivenEngine) ProcessCertificates(qc *hotstuff.QuorumCert) error {
	if err := e.advanceRound(qc); err != nil {
		return err
	}

	e.safety.UpdateLockQCRound(qc.View.Round)

	// try to commit locked block and pure the `pendingBlockTree`
	e.blkTree.ProcessCommit(qc.Hash)
	return nil
}
