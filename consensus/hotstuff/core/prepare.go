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

	code := MsgTypePrepare
	request := c.current.PendingRequest()
	if request.Proposal.NumberU64() != c.current.HeightU64() {
		logger.Trace("Failed to send prepare", "msg", code, "err", "request height invalid")
		return
	}
	//expectAddr := request.Proposal.Coinbase()
	//if got := c.Address(); expectAddr != c.Address() {
	//	logger.Trace("Failed to send prepare", "msg", code, "coinbase expect", expectAddr, "got", got)
	//	return
	//}
	if c.current.HighQC().hash != request.Proposal.ParentHash() {
		logger.Trace("Failed to send prepare", "msg", code, "err", "request parent hash invalid")
		return
	}

	prepare := &MsgPrepare{
		Proposal: request.Proposal,
		HighQC:   c.current.HighQC(),
	}
	payload, err := Encode(prepare)
	if err != nil {
		logger.Trace("Failed to encode", "msg", code, "err", err)
		return
	}

	// consensus spent time always less than a block period, waiting for `delay` time to catch up the system time.
	delay := time.Unix(int64(prepare.Proposal.Time()), 0).Sub(time.Now())
	time.Sleep(delay)
	logger.Trace("delay to broadcast proposal", "msg", code, "time", delay.Milliseconds())

	c.broadcast(code, payload)
	logger.Trace("sendPrepare", "msg", code, "proposal", prepare.Proposal.Hash())
}

func (c *core) handlePrepare(data *Message) error {
	logger := c.newLogger()

	var (
		msg    *MsgPrepare
		code = MsgTypePrepare
		src    = data.address
	)
	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to decode", "msg", code, "src", src, "err", err)
		return errFailedDecodePrepare
	}
	if err := c.checkView(code, data.View); err != nil {
		logger.Trace("Failed to check view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkMsgFromProposer(src); err != nil {
		logger.Trace("Failed to check proposer", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkProposalView(msg.Proposal, data.View); err != nil {
		logger.Trace("Failed to check proposal and msg view", "msg", code, "src", src, "err", err)
		return err
	}

	if _, err := c.backend.VerifyUnsealedProposal(msg.Proposal); err != nil {
		logger.Trace("Failed to verify unsealed proposal", "msg", code, "src", src, "err", err)
		return errVerifyUnsealedProposal
	}
	if err := c.extend(msg.Proposal, msg.HighQC); err != nil {
		logger.Trace("Failed to check extend", "msg", code, "src", src, "err", err)
		return errExtend
	}
	if err := c.safeNode(msg.Proposal, msg.HighQC); err != nil {
		logger.Trace("Failed to check safeNode", "msg", code, "src", src, "err", err)
		return errSafeNode
	}
	if err := c.preExecuteBlock(msg.Proposal); err != nil {
		logger.Trace("Failed to pre-execute block", "msg", code, "src", src, "err", err)
		return err
	}

	logger.Trace("handlePrepare", "msg", code, "src", src, "hash", msg.Proposal.Hash())

	// leader accept proposal
	if c.IsProposer() && c.currentState() < StatePrepared {
		c.current.SetProposal(msg.Proposal)
		c.sendPrepareVote()
	}

	// repo accept proposal and high qc
	// todo(fuk): 是否真的需要stateHighQC
	if !c.IsProposer() && c.currentState() < StateHighQC {
		c.current.SetHighQC(msg.HighQC)
		c.current.SetProposal(msg.Proposal)
		c.setCurrentState(StateHighQC)
		logger.Trace("acceptHighQC", "msg", code, "highQC", msg.HighQC.hash)

		c.sendPrepareVote()
	}

	return nil
}

func (c *core) sendPrepareVote() {
	logger := c.newLogger()

	code := MsgTypePrepareVote
	vote := c.current.Vote()
	if vote == nil {
		logger.Trace("Failed to send vote", "msg", code, "err", "current vote is nil")
		return
	}
	payload, err := Encode(vote)
	if err != nil {
		logger.Trace("Failed to encode", "msg", code, "err", err)
		return
	}
	c.broadcast(code, payload)
	logger.Trace("sendPrepareVote", "msg", code, "hash", vote.Digest)
}

func (c *core) extend(proposal hotstuff.Proposal, highQC *QuorumCert) error {
	block, ok := proposal.(*types.Block)
	if !ok {
		return fmt.Errorf("invalid proposal: hash %s", proposal.Hash())
	}
	if err := c.verifyVoteQC(highQC.hash, highQC); err != nil {
		return err
	}
	if highQC.hash != block.ParentHash() {
		return fmt.Errorf("block %v (parent %v) not extend hiqhQC %v", block.Hash(), block.ParentHash(), highQC.hash)
	}
	return nil
}

// proposal extend lockedQC `OR` hiqhQC.view > lockedQC.view
func (c *core) safeNode(proposal hotstuff.Proposal, highQC *QuorumCert) error {
	//logger := c.newLogger()

	if proposal.Number().Uint64() == 1 && c.current.lockedQC == nil {
		return nil
	}
	//safety := false
	//liveness := false

	if highQC.view.Cmp(c.current.lockedQC.view) > 0 ||
		proposal.ParentHash() == c.current.lockedQC.hash {
		return nil
	} else {
		return errSafeNode
	}

	//if c.current.PreCommittedQC() == nil {
	//	logger.Trace("safeNodeChecking", "lockQC", "is nil")
	//	return errSafeNode
	//}
	//if err := c.extend(proposal, c.current.PreCommittedQC()); err == nil {
	//	safety = true
	//} else {
	//	logger.Trace("safeNodeChecking", "extend err", err)
	//}
	//if highQC.view.Cmp(c.current.PreCommittedQC().view) > 0 {
	//	liveness = true
	//}
	//if safety || liveness {
	//	return nil
	//} else {
	//	return errSafeNode
	//}
}
