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
	"time"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *core) sendPrepare() {
	logger := c.newLogger()

	if !c.IsProposer() {
		return
	}
	if c.currentState() != StateHighQC {
		return
	}
	if c.current.PendingRequest() == nil || c.current.PendingRequest().Proposal == nil {
		return
	}

	request := c.current.PendingRequest()
	expectHeight := c.current.HeightU64()
	if got := request.Proposal.NumberU64(); expectHeight != got {
		logger.Trace("Failed to send prepare", "height expect", expectHeight, "got", got)
		return
	}
	expectAddr := request.Proposal.Coinbase()
	if got := c.Address(); expectAddr != c.Address() {
		logger.Trace("Failed to send prepare", "coinbase expect", expectAddr, "got", got)
		return
	}
	expectParentHash := c.current.HighQC().Hash()
	if got := request.Proposal.ParentHash(); expectParentHash != got {
		logger.Trace("Failed to send prepare", "expect parent hash", expectParentHash, "got", got)
		return
	}

	msgTyp := MsgTypePrepare
	prepare := &MsgPrepare{
		View:     c.currentView(),
		Proposal: request.Proposal,
		HighQC:   c.current.HighQC(),
	}
	payload, err := Encode(prepare)
	if err != nil {
		logger.Trace("Failed to encode", "msg", msgTyp, "err", err)
		return
	}

	// consensus spent time always less than a block period, waiting for `delay` time to catch up the system time.
	delay := time.Unix(int64(prepare.Proposal.Time()), 0).Sub(time.Now())
	time.Sleep(delay)
	logger.Trace("delay to broadcast proposal", "time", delay.Milliseconds())

	c.broadcast(&Message{Code: msgTyp, Msg: payload})
	logger.Trace("sendPrepare", "prepare view", prepare.View, "proposal", prepare.Proposal.Hash())
}

func (c *core) handlePrepare(data *Message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		msg    *MsgPrepare
		msgTyp = MsgTypePrepare
	)
	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to decode", "type", msgTyp, "err", err)
		return errFailedDecodePrepare
	}
	if err := c.checkView(msgTyp, msg.View); err != nil {
		logger.Trace("Failed to check view", "msg", msgTyp, "err", err)
		return err
	}
	if err := c.checkMsgFromProposer(src); err != nil {
		logger.Trace("Failed to check proposer", "msg", msgTyp, "err", err)
		return err
	}
	if err := c.checkProposalView(msg.Proposal, msg.View); err != nil {
		logger.Trace("Failed to check proposal and msg view", "msg", msgTyp, "err", err)
		return err
	}

	if _, err := c.backend.VerifyUnsealedProposal(msg.Proposal); err != nil {
		logger.Trace("Failed to verify unsealed proposal", "msg", msgTyp, "err", err)
		return errVerifyUnsealedProposal
	}
	if err := c.extend(msg.Proposal, msg.HighQC); err != nil {
		logger.Trace("Failed to check extend", "msg", msgTyp, "err", err)
		return errExtend
	}
	if err := c.safeNode(msg.Proposal, msg.HighQC); err != nil {
		logger.Trace("Failed to check safeNode", "msg", msgTyp, "err", err)
		return errSafeNode
	}
	// todo: delete after test
	//if err := c.checkLockedProposal(msg.Proposal); err != nil {
	//	logger.Trace("Failed to check locked proposal", "msg", msgTyp, "err", err)
	//	return err
	//}
	if err := c.preExecuteBlock(msg.Proposal); err != nil {
		logger.Trace("Failed to pre-execute block", "msg", msgTyp, "err", err)
		return err
	}

	logger.Trace("handlePrepare", "msg", msgTyp, "src", src.Address(), "hash", msg.Proposal.Hash())

	// leader accept proposal
	if c.IsProposer() && c.currentState() < StatePrepared {
		c.current.SetProposal(msg.Proposal)
		c.sendPrepareVote()
	}

	// repo accept proposal and high qc
	if !c.IsProposer() && c.currentState() < StateHighQC {
		c.current.SetHighQC(msg.HighQC)
		c.current.SetProposal(msg.Proposal)
		c.setCurrentState(StateHighQC)
		logger.Trace("acceptHighQC", "msg", msgTyp, "src", src.Address(), "highQC", msg.HighQC.Hash)

		c.sendPrepareVote()
	}

	return nil
}

func (c *core) sendPrepareVote() {
	logger := c.newLogger()

	msgTyp := MsgTypePrepareVote
	vote := c.current.Vote()
	if vote == nil {
		logger.Trace("Failed to send vote", "msg", msgTyp, "err", "current vote is nil")
		return
	}
	payload, err := Encode(vote)
	if err != nil {
		logger.Trace("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&Message{Code: msgTyp, Msg: payload})
	logger.Trace("sendPrepareVote", "vote view", vote.View, "vote", vote.Digest)
}

func (c *core) extend(proposal hotstuff.Proposal, highQC *QuorumCert) error {
	block, ok := proposal.(*types.Block)
	if !ok {
		return fmt.Errorf("invalid proposal: hash %s", proposal.Hash())
	}
	if err := c.verifyCrossEpochQC(highQC); err != nil {
		return err
	}
	if highQC.Hash() != block.ParentHash() {
		return fmt.Errorf("block %v (parent %v) not extend hiqhQC %v", block.Hash(), block.ParentHash(), highQC.Hash)
	}
	return nil
}

// proposal extend lockedQC `OR` hiqhQC.view > lockedQC.view
func (c *core) safeNode(proposal hotstuff.Proposal, highQC *QuorumCert) error {
	logger := c.newLogger()

	if proposal.Number().Uint64() == 1 {
		return nil
	}
	safety := false
	liveness := false
	if c.current.PreCommittedQC() == nil {
		logger.Trace("safeNodeChecking", "lockQC", "is nil")
		return errSafeNode
	}
	if err := c.extend(proposal, c.current.PreCommittedQC()); err == nil {
		safety = true
	} else {
		logger.Trace("safeNodeChecking", "extend err", err)
	}
	if highQC.view.Cmp(c.current.PreCommittedQC().view) > 0 {
		liveness = true
	}
	if safety || liveness {
		return nil
	} else {
		return errSafeNode
	}
}
