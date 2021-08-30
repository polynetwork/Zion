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
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
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
	chain   consensus.ChainReader
	backend hotstuff.Backend

	state  State
	addr   common.Address
	signer hotstuff.Signer
	valset hotstuff.ValidatorSet

	epoch            uint64
	epochHeightStart *big.Int // [epochHeightStart, epochHeightEnd] is an closed interval
	epochHeightEnd   *big.Int
	curRound         *big.Int // 从genesis block 0开始
	curHeight        *big.Int // 从genesis block 0开始

	requests *requestSet
	messages *MessagePool
	blkPool  *BlockPool

	// pace maker
	highestCommitRound *big.Int    // used to calculate timeout duration
	timer              *time.Timer // drive consensus round

	// safety
	//lockQCRound   *big.Int
	lockQC        *hotstuff.QuorumCert
	lastVoteRound *big.Int

	events *event.TypeMuxSubscription
	//timeoutSub        *event.TypeMuxSubscription
	finalCommittedSub *event.TypeMuxSubscription

	validateFn func([]byte, []byte) (common.Address, error)
}

// todo: chain
func New(
	backend hotstuff.Backend,
	c *hotstuff.Config,
	db ethdb.Database,
	signer hotstuff.Signer,
	valset hotstuff.ValidatorSet,
) *EventDrivenEngine {

	e := &EventDrivenEngine{
		config:  c,
		logger:  log.New("address", backend.Address()),
		backend: backend,
		db:      db,
		//chain:   chain,
	}
	e.addr = backend.Address()
	e.requests = newRequestSet()
	e.valset = valset
	e.signer = signer
	e.requests = newRequestSet()
	e.messages = NewMessagePool(valset)
	e.validateFn = e.checkValidatorSignature

	if err := e.initialize(); err != nil {
		e.logger.Error("initialize event-driven hotstuff protocol", "err", err)
		return nil
	}
	return e
}

// handleNewRound proposer at this round get an new proposal and broadcast to all validators.
func (e *EventDrivenEngine) handleNewRound() error {
	logger := e.newLogger()
	e.state = StateAcceptRequest

	if !e.IsProposer() {
		return nil
	}

	req := e.requests.GetRequest(e.currentView())
	if req == nil {
		return nil
	}
	proposal, ok := req.Proposal.(*types.Block)
	if !ok {
		return errProposalConvert
	}

	justifyQC := e.blkPool.GetHighQC()
	view := e.currentView()
	msg := &MsgProposal{
		Epoch:     e.epoch,
		View:      view,
		Proposal:  proposal,
		JustifyQC: justifyQC,
	}
	e.state = StateAcceptProposal

	logger.Debug("----------New round", "state", e.currentState(), "newRound", e.curRound, "new_proposer", e.valset.GetProposer(), "valSet", e.valset.List(), "size", e.valset.Size(), "IsProposer", e.IsProposer())

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
	proposal := msg.Proposal
	hash := proposal.Hash()
	proposer := proposal.Coinbase()
	header := proposal.Header()
	justifyQC := msg.JustifyQC

	if err := e.checkEpoch(msg.Epoch, proposal.Number()); err != nil {
		return err
	}
	if err := e.checkJustifyQC(proposal, justifyQC); err != nil {
		return err
	}
	if err := e.signer.VerifyQC(justifyQC, e.valset); err != nil {
		return err
	}
	if err := e.signer.VerifyHeader(header, e.valset, false); err != nil {
		return err
	}

	// try to advance into new round, it will update proposer and current view
	if err := e.processQC(justifyQC); err != nil {
		logger.Error("Failed to process qc", "err", err)
	}

	if err := e.checkProposer(proposer); err != nil {
		return err
	}
	if err := e.checkView(view); err != nil {
		return err
	}

	e.blkPool.UpdateHighQC(justifyQC)
	if err := e.blkPool.Insert(proposal, view.Round); err != nil {
		return err
	}

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

	if err := e.checkVote(vote); err != nil {
		return err
	}
	if err := e.checkEpoch(vote.Epoch, vote.View.Height); err != nil {
		return err
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

	qc, err := e.aggregateQC(vote, size)
	if err != nil {
		return err
	}

	e.blkPool.UpdateHighQC(qc)
	highQC := e.blkPool.GetHighQC()

	if err := e.advanceRoundByQC(highQC, false); err != nil {
		return err
	}
	e.state = StateVoted

	return nil
}

func (e *EventDrivenEngine) handleQC(src hotstuff.Validator, data *hotstuff.Message) error {
	var (
		qc *hotstuff.QuorumCert
	)
	if err := data.Decode(&qc); err != nil {
		return err
	}

	if err := e.signer.VerifyQC(qc, e.valset); err != nil {
		return err
	}

	return e.processQC(qc)
}

func (e *EventDrivenEngine) handleTC(src hotstuff.Validator, data *hotstuff.Message) error {
	var (
		tc *TimeoutCert
	)
	if err := data.Decode(&tc); err != nil {
		return err
	}

	if err := e.signer.VerifyCommittedSeal(e.valset, tc.Hash, tc.Seals); err != nil {
		return err
	}

	if err := e.advanceRoundByTC(tc, false); err != nil {
		return err
	}

	return nil
}

// try to advance into new round, it will update proposer and current view
// commit the proposal
func (e *EventDrivenEngine) processQC(qc *hotstuff.QuorumCert) error {
	// try to advance consensus into next round
	e.advanceRoundByQC(qc, false)

	// commit qc grand (proposal's great-grand parent block)
	lastLockQC := e.getLockQC()
	if committedBlock := e.blkPool.GetCommitBlock(lastLockQC.Hash); committedBlock != nil {
		if existProposal := e.backend.GetProposal(committedBlock.Hash()); existProposal == nil {
			// todo: 如果节点此时宕机怎么办？还是说允许所有的节点一起提交区块
			if e.isSelf(committedBlock.Coinbase()) {
				e.backend.Commit(committedBlock)
			}
		}
		e.blkPool.Pure(committedBlock.Hash())
	}

	return e.updateLockQC(qc)
}
