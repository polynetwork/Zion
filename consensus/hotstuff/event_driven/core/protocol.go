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

	state   State
	started bool
	addr    common.Address
	signer  hotstuff.Signer
	valset  hotstuff.ValidatorSet

	epoch            uint64
	epochHeightStart *big.Int // [epochHeightStart, epochHeightEnd] is an closed interval
	epochHeightEnd   *big.Int
	curRound         *big.Int // 从genesis block 0开始
	curHeight        *big.Int // 从genesis block 0开始
	curRequest 		 *types.Block

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

	feed event.Feed // request feed
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
) hotstuff.CoreEngine {

	addr := signer.Address()
	engine := &EventDrivenEngine{
		config:  c,
		db:      db,
		backend: backend,
		//chain:   chain,
		logger: log.New("address", addr),
	}
	engine.addr = addr
	engine.valset = valset
	engine.signer = signer
	engine.started = false
	engine.requests = newRequestSet()
	engine.messages = NewMessagePool(valset)
	engine.validateFn = engine.checkValidatorSignature

	return engine
}

// handleNewRound proposer at this round get an new proposal and broadcast to all validators.
func (e *EventDrivenEngine) handleNewRound() error {
	if !e.started {
		return nil
	}

	logger := e.newLogger()

	e.state = StateAcceptRequest
	msgTyp := MsgTypeNewRound

	logger.Trace("New round", "msg", msgTyp, "state", e.currentState(), "new_proposer", e.valset.GetProposer(), "valSet", e.valset.List(), "size", e.valset.Size(), "IsProposer", e.IsProposer())

	if !e.IsProposer() {
		return nil
	}

	// get or preparing request, although validator may be not the proposer in this round, it can fetch proposal from
	// miner and used it in the correct round.
	proposal := e.getRequest()
	if proposal == nil {
		logger.Error("Failed to get request", "msg", msgTyp, "err", "no request")
		return errNoRequest
	}

	justifyQC, _ := e.blkPool.GetHighQC()
	view := e.currentView()
	msg := &MsgProposal{
		Epoch:     e.epoch,
		View:      view,
		Proposal:  proposal,
		JustifyQC: justifyQC,
	}

	e.state = StateAcceptProposal
	e.encodeAndBroadcast(MsgTypeProposal, msg)

	logger.Trace("Generate proposal", "msg", msgTyp)
	return nil
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
	justifyQC := msg.JustifyQC

	if err := e.checkEpoch(msg.Epoch, proposal.Number()); err != nil {
		logger.Trace("Failed to check epoch", "msg", msgTyp, "err", err)
		return err
	}
	if err := e.validateProposal(proposal); err != nil {
		logger.Trace("Failed to validate proposal", "msg", msgTyp, "err", err)
		return err
	}
	if err := e.checkJustifyQC(proposal, justifyQC); err != nil {
		logger.Trace("Failed to check justify", "msg", msgTyp, "err", err)
		return err
	}
	if err := e.signer.VerifyQC(justifyQC, e.valset); err != nil {
		logger.Trace("Failed to verify justifyQC", "msg", msgTyp, "err", err)
		return err
	}
	if err := e.checkView(view); err != nil {
		logger.Trace("Failed to check view", "msg", msgTyp, "err", err)
		return err
	}

	// try to advance into new round, it will update proposer and current view
	_ = e.processQC(justifyQC)

	if err := e.checkProposer(proposal.Coinbase()); err != nil {
		logger.Trace("Failed to check proposer", "msg", msgTyp, "err", err)
		return err
	}

	if err := e.blkPool.AddBlock(proposal, view.Round); err != nil {
		logger.Trace("Failed to insert block into block pool", "msg", msgTyp, "err", err)
		return err
	}
	e.blkPool.UpdateHighQC(justifyQC)

	logger.Trace("Accept proposal", "msg", msgTyp, "proposer", src.Address(), "hash", proposal.Hash(), "height", proposal.Number())

	vote, err := e.makeVote(proposal.Hash(), proposal.Coinbase(), view, justifyQC)
	if err != nil {
		logger.Trace("Failed to make vote", "msg", msgTyp, "err", err)
		return err
	}

	e.increaseLastVoteRound(view.Round)
	e.encodeAndBroadcast(MsgTypeVote, vote)
	logger.Trace("Send Vote", "msg", msgTyp, "to", e.nextProposer(), "hash", vote.Hash)
	return nil
}

// handleVote validate vote message and try to assemble qc
func (e *EventDrivenEngine) handleVote(src hotstuff.Validator, data *hotstuff.Message) error {
	var (
		vote   *Vote
		msgTyp = MsgTypeVote
	)

	logger := e.newLogger()
	if err := data.Decode(&vote); err != nil {
		logger.Trace("Failed to decode", "msg", msgTyp, "err", err)
		return errFailedDecodeNewView
	}

	if err := e.checkVote(vote); err != nil {
		logger.Trace("Failed to check vote", "msg", msgTyp, "err", err)
		return err
	}
	if err := e.checkEpoch(vote.Epoch, vote.View.Height); err != nil {
		logger.Trace("Failed to check epoch", "msg", msgTyp, "err", err)
		return err
	}
	if err := e.validateVote(vote); err != nil {
		logger.Trace("Failed to validate vote", "msg", msgTyp, "err", err)
		return err
	}
	if err := e.messages.AddVote(vote.Hash, data); err != nil {
		logger.Trace("Failed to add vote", "msg", msgTyp, "err", err)
		return err
	}

	logger.Trace("Accept Vote", "msg", msgTyp, "from", src.Address(), "hash", vote.Hash)

	size := e.messages.VoteSize(vote.Hash)
	if size != e.Q() {
		return nil
	}

	qc, proposal, err := e.aggregateQC(vote, size)
	if err != nil {
		logger.Trace("Failed to aggregate qc", "msg", msgTyp, "err", err)
		return err
	}
	if err := e.blkPool.AddBlock(proposal, vote.View.Round); err != nil {
		logger.Trace("Failed to insert block into block pool", "msg", msgTyp, "err", err)
		return err
	}
	e.blkPool.UpdateHighQC(qc)
	highQC, _ := e.blkPool.GetHighQC()

	if err := e.advanceRoundByQC(highQC, false); err != nil {
		logger.Trace("Failed to advance round", "msg", msgTyp, "err", err)
		return err
	}

	e.state = StateVoted

	logger.Trace("Aggregate QC", "msg", msgTyp, "qc", qc.Hash, "view", qc.View)

	return nil
}

func (e *EventDrivenEngine) handleQC(src hotstuff.Validator, data *hotstuff.Message) error {
	logger := e.newLogger()

	var (
		qc     *hotstuff.QuorumCert
		msgTyp = MsgTypeQC
	)
	if err := data.Decode(&qc); err != nil {
		logger.Trace("Failed to decode", "msg", msgTyp, "err", err)
		return err
	}

	if err := e.signer.VerifyQC(qc, e.valset); err != nil {
		logger.Trace("Failed to verify qc", "msg", msgTyp, "err", err)
		return err
	}

	if err := e.processQC(qc); err != nil {
		logger.Trace("Failed to process qc", "msg", msgTyp, "err", err)
		return err
	}

	logger.Trace("Accept QC", "msg", msgTyp, "src", src.Address(), "qc", qc.Hash, "view", qc.View)
	return nil
}

func (e *EventDrivenEngine) handleTC(src hotstuff.Validator, data *hotstuff.Message) error {
	logger := e.newLogger()

	var (
		tc     *TimeoutCert
		msgTyp = MsgTypeTC
	)
	if err := data.Decode(&tc); err != nil {
		logger.Trace("Failed to decode", "msg", msgTyp, "err", err)
		return err
	}

	if err := e.signer.VerifyCommittedSeal(e.valset, tc.Hash, tc.Seals); err != nil {
		logger.Trace("Failed to verify committed seal", "msg", msgTyp, "err", err)
		return err
	}

	if err := e.advanceRoundByTC(tc, false); err != nil {
		logger.Trace("Failed to advance by tc", "msg", msgTyp, "err", err)
		return err
	}

	logger.Trace("Accept TC", "msg", msgTyp, "src", src.Address(), "tc", tc.Hash, "view", tc.View)
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
