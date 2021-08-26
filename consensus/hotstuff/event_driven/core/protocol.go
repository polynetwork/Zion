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
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

// EventDrivenEngine implement event-driven hotstuff protocol, it obtains:
// 1.validator set which represent consensus participants
type EventDrivenEngine struct {
	config  *hotstuff.Config
	logger  log.Logger
	db      ethdb.Database
	backend hotstuff.Backend

	addr   common.Address
	signer hotstuff.Signer
	valset hotstuff.ValidatorSet

	epoch            uint64
	epochHeightStart *big.Int
	epochHeightEnd   *big.Int
	curRound         *big.Int // 从genesis block 0开始
	curHeight        *big.Int // 从genesis block 0开始

	requests *requestSet
	messages *MessagePool
	blkTree  *BlockTree

	// pace maker
	highestCommitRound *big.Int
	timer              *time.Timer

	// safety
	lockQCRound   *big.Int
	lastVoteRound *big.Int

	events *event.TypeMuxSubscription
	//timeoutSub        *event.TypeMuxSubscription
	finalCommittedSub *event.TypeMuxSubscription

	validateFn func([]byte, []byte) (common.Address, error)
}

func NewEventDrivenEngine(valset hotstuff.ValidatorSet) *EventDrivenEngine {
	return nil
}

// handleNewRound proposer at this round get an new proposal and broadcast to all validators.
func (e *EventDrivenEngine) handleNewRound() error {
	if !e.IsProposer() {
		return nil
	}
	msg, err := e.generateProposalMessage()
	if err != nil {
		return err
	}
	return e.encodeAndBroadcast(MsgTypeProposal, msg)
}

// handleProposal validate proposal info and vote to the next leader if the proposal is valid
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

	view := msg.View
	epoch := msg.Epoch
	proposal := msg.Proposal
	hash := proposal.Hash()
	proposer := proposal.Coinbase()
	header := proposal.Header()
	justifyQC := msg.JustifyQC

	if epoch != e.epoch {
		return errInvalidEpoch
	}
	if err := e.signer.VerifyQC(justifyQC, e.valset); err != nil {
		return err
	}
	if err := e.signer.VerifyHeader(header, e.valset, false); err != nil {
		return err
	}

	// try to advance into new round, it will update proposer and current view
	if err := e.advanceRound(justifyQC, false); err != nil {
		logger.Trace("Failed to advance new round", "err", err)
	} else {
		e.updateLockQCRound(justifyQC.View.Round)
		e.blkTree.ProcessCommit(justifyQC.Hash)
	}

	if err := e.checkProposer(proposer); err != nil {
		return err
	}
	if err := e.checkView(view); err != nil {
		return err
	}

	e.blkTree.Insert(proposal)

	vote, err := e.makeVote(hash, proposer, view, justifyQC)
	if err != nil {
		return err
	}

	e.increaseLastVoteRound(view.Round)

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

	if err := e.validateVote(vote); err != nil {
		return err
	}
	if err := e.messages.AddVote(vote.Hash, data); err != nil {
		return err
	}

	size := e.messages.VoteSize(vote.Hash)
	if size != e.Q() {
		return nil
	}

	qc, err := e.aggregate(vote, size)
	if err != nil {
		return err
	}

	// paceMaker send qc to next leader
	if err := e.advanceRound(qc, false); err != nil {
		return err
	}
	e.blkTree.UpdateHighQC(qc)

	return nil
}

func (e *EventDrivenEngine) handleCertificate(src hotstuff.Validator, data *hotstuff.Message) error {
	var (
		certEvt *CertificateEvent
	)
	if err := data.Decode(&certEvt); err != nil {
		return err
	}

	qc := certEvt.Cert
	if err := e.signer.VerifyQC(qc, e.valset); err != nil {
		return err
	}

	if err := e.advanceRound(qc, false); err != nil {
		return err
	}

	e.updateLockQCRound(qc.View.Round)

	// try to commit locked block and pure the `pendingBlockTree`
	e.blkTree.ProcessCommit(qc.Hash)
	return nil
}
