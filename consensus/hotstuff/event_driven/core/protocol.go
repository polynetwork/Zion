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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

// EventDrivenEngine implement event-driven hotstuff protocol, it obtains:
// 1.validator set which represent consensus participants
type EventDrivenEngine struct {
	config *hotstuff.Config
	logger log.Logger

	addr   common.Address
	signer hotstuff.Signer
	valset hotstuff.ValidatorSet

	requests  *requestSet
	paceMaker *PaceMaker
	blkTree   *BlockTree
	safety    *SafetyRules

	backend hotstuff.Backend

	events            *event.TypeMuxSubscription
	timeoutSub        *event.TypeMuxSubscription
	finalCommittedSub *event.TypeMuxSubscription

	validateFn func([]byte, []byte) (common.Address, error)
}

func NewEventDrivenEngine(valset hotstuff.ValidatorSet) *EventDrivenEngine {
	return nil
}

// ProcessNewRoundEvent proposer at this round get an new proposal and broadcast to all validators.
func (e *EventDrivenEngine) ProcessNewRoundEvent() error {
	view := e.currentView()
	e.valset.CalcProposerByIndex(view.Round.Uint64())

	if !e.isProposer() {
		return nil
	}

	proposal := e.getCurrentPendingRequest()
	msg := &MsgNewView{
		View:     view,
		Proposal: proposal,
	}
	return e.encodeAndBroadcast(MsgTypeNewView, msg)
}

// ProcessProposal validate proposal info and vote to the next leader if the proposal is valid
func (e *EventDrivenEngine) ProcessProposal(proposal *types.Block) error {
	justifyQC, proposalRound, err := extraProposal(proposal)
	if err != nil {
		return err
	}

	if err := e.ProcessCertificates(justifyQC); err != nil {
		return err
	}

	currentRound := e.paceMaker.CurrentRound()
	if currentRound.Cmp(proposalRound) != 0 {
		// todo: modify err type
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

	vote, err := e.safety.VoteRule(proposal, proposalRound, justifyQC)
	if err != nil {
		return err
	}

	// todo: vote to next proposal
	return e.encodeAndBroadcast(MsgTypeVote, vote)
}

// ProcessVoteMsg validate vote message and try to assemble qc
func (e *EventDrivenEngine) ProcessVoteMsg(vote *Vote) error {
	e.blkTree.ProcessVote(vote)
	return nil
}

// ProcessCertificates validate and handle QC/TC
func (e *EventDrivenEngine) ProcessCertificates(qc *hotstuff.QuorumCert) error {
	if err := e.paceMaker.AdvanceRound(qc); err != nil {
		return err
	}

	e.safety.UpdateLockQC(qc, qc.View.Round)

	// try to commit locked block and pure the `pendingBlockTree`
	e.blkTree.ProcessCommit(qc.Hash)
	return nil
}
