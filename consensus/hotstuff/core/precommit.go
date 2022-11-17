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
)

func (c *core) handlePrepareVote(data *Message) error {
	logger := c.newLogger()

	var (
		vote = common.BytesToHash(data.Msg)
		code = MsgTypePrepareVote
		src  = data.address
	)
	if err := c.checkView(data.View); err != nil {
		logger.Trace("Failed to check view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkVote(data, vote); err != nil {
		logger.Trace("Failed to check vote", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkMsgToProposer(); err != nil {
		logger.Trace("Failed to check proposer", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.current.AddPrepareVote(data); err != nil {
		logger.Trace("Failed to add vote", "msg", code, "src", src, "err", err)
		return errAddPrepareVote
	}

	logger.Trace("handlePrepareVote", "msg", code, "src", src, "hash", vote)

	if size := c.current.PrepareVoteSize(); size >= c.Q() && c.currentState() < StatePrepared {
		prepareQC, err := c.messages2qc(c.proposer(), c.current.Proposal().Hash(), c.current.PrepareVotes())
		if err != nil {
			logger.Trace("Failed to assemble prepareQC", "msg", code, "err", err)
			return errInvalidQC
		}
		if err := c.acceptPrepare(prepareQC, c.current.Proposal()); err != nil {
			logger.Trace("Failed to accept prepareQC", "msg", code, "err", err)
			return err
		}
		c.sendPreCommit()
	}

	return nil
}

func (c *core) sendPreCommit() {
	logger := c.newLogger()

	code := MsgTypePreCommit
	msg := &Subject{
		Proposal: c.current.Proposal(),
		QC:       c.current.PrepareQC(),
	}
	payload, err := Encode(msg)
	if err != nil {
		logger.Trace("Failed to encode", "msg", code, "err", err)
		return
	}
	c.broadcast(code, payload)
	logger.Trace("sendPreCommit", "msg", code, "proposal", msg.Proposal.Hash())
}

func (c *core) handlePreCommit(data *Message) error {
	logger := c.newLogger()

	var (
		msg  *Subject
		code = MsgTypePreCommit
		src  = data.address
	)
	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to check decode", "msg", code, "src", src, "err", err)
		return errFailedDecodePreCommit
	}
	if err := c.checkView(data.View); err != nil {
		logger.Trace("Failed to check view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkProposalView(msg.Proposal, data.View); err != nil {
		logger.Trace("Failed to check proposal and msg view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkMsgFromProposer(src); err != nil {
		logger.Trace("Failed to check proposer", "msg", code, "src", src, "err", err)
		return err
	}
	if msg.Proposal.Hash() != msg.QC.hash {
		logger.Trace("Failed to check msg", "msg", code, "src", src, "expect prepareQC hash", msg.Proposal.Hash(), "got", msg.QC.hash)
		return errInvalidProposal
	}
	if _, err := c.backend.Verify(msg.Proposal); err != nil {
		logger.Trace("Failed to check verify proposal", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.verifyQC(data, msg.QC); err != nil {
		logger.Trace("Failed to verify prepareQC", "msg", code, "src", src, "err", err)
		return err
	}

	logger.Trace("handlePreCommit", "msg", code, "src", src, "hash", msg.Proposal.Hash())

	if c.IsProposer() && c.currentState() < StateLocked {
		c.sendPreCommitVote()
	}
	if !c.IsProposer() && c.currentState() < StatePrepared {
		if err := c.acceptPrepare(msg.QC, msg.Proposal); err != nil {
			logger.Trace("Failed to accept prepareQC", "msg", code, "err", err)
			return err
		}
		logger.Trace("acceptPrepare", "msg", code, "prepareQC", msg.QC.hash)
		c.sendPreCommitVote()
	}

	return nil
}

func (c *core) acceptPrepare(prepareQC *QuorumCert, proposal hotstuff.Proposal) error {
	if err := c.current.SetPrepareQC(prepareQC); err != nil {
		return err
	}
	if err := c.current.SetProposal(proposal); err != nil {
		return err
	}
	c.current.SetState(StatePrepared)
	return nil
}

func (c *core) sendPreCommitVote() {
	logger := c.newLogger()

	code := MsgTypePreCommitVote
	vote := c.current.Vote()
	c.broadcast(code, vote.Bytes())
	logger.Trace("sendPreCommitVote", "msg", code, "hash", vote)
}
