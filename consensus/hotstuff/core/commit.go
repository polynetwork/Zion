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
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) handlePreCommitVote(data *Message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		vote   *Vote
		msgTyp = MsgTypePreCommitVote
	)
	if err := data.Decode(&vote); err != nil {
		logger.Trace("Failed to decode", "type", msgTyp, "src", src.Address(), "err", err)
		return errFailedDecodePreCommitVote
	}
	if err := c.checkView(msgTyp, vote.View); err != nil {
		logger.Trace("Failed to check view", "type", msgTyp, "src", src.Address(), "err", err)
		return err
	}
	if err := c.checkVote(vote); err != nil {
		logger.Trace("Failed to check vote", "type", msgTyp, "src", src.Address(), "err", err)
		return err
	}
	if err := c.checkProposal(vote.Digest); err != nil {
		logger.Trace("Failed to check hash", "type", msgTyp, "src", src.Address(), "expect vote", vote.Digest)
		return errInvalidDigest
	}
	if err := c.checkMsgToProposer(); err != nil {
		logger.Trace("Failed to check proposal", "type", msgTyp, "src", src.Address(), "err", err)
		return err
	}

	if err := c.current.AddPreCommitVote(data); err != nil {
		logger.Trace("Failed to add vote", "type", msgTyp, "src", src.Address(), "err", err)
		return errAddPreCommitVote
	}

	logger.Trace("handlePreCommitVote", "msg", msgTyp, "src", src.Address(), "hash", vote.Digest)

	if size := c.current.PreCommitVoteSize(); size >= c.Q() && c.currentState() < StatePreCommitted {
		c.lockQCAndProposal(c.current.PrepareQC())
		logger.Trace("acceptPreCommitted", "msg", msgTyp, "msgSize", size)
		c.sendCommit()
	}
	return nil
}

func (c *core) sendCommit() {
	logger := c.newLogger()

	msgTyp := MsgTypeCommit
	sub := c.current.PreCommittedQC()
	payload, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&Message{Code: msgTyp, Msg: payload})
	logger.Trace("sendCommit", "msg", msgTyp, "proposal", sub.hash)
}

func (c *core) handleCommit(data *Message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		msg    *QuorumCert
		msgTyp = MsgTypeCommit
	)
	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to decode", "msg", msgTyp, "src", src.Address(), "err", err)
		return errFailedDecodeCommit
	}
	if err := c.checkView(MsgTypeCommit, msg.view); err != nil {
		logger.Trace("Failed to check view", "msg", msgTyp, "src", src.Address(), "err", err)
		return err
	}
	if err := c.checkMsgFromProposer(src); err != nil {
		logger.Trace("Failed to check proposer", "msg", msgTyp, "src", src.Address(), "err", err)
		return err
	}
	if err := c.checkPrepareQC(msg); err != nil {
		logger.Trace("Failed to check prepareQC", "msg", msgTyp, "src", src.Address(), "err", err)
		return err
	}
	if err := c.signer.VerifyQC(msg, c.valSet); err != nil {
		logger.Trace("Failed to check verify qc", "msg", msgTyp, "src", src.Address(), "err", err)
		return err
	}

	logger.Trace("handleCommit", "msg", msgTyp, "src", src.Address(), "proposal", msg.hash)

	if c.IsProposer() && c.currentState() < StateCommitted {
		c.sendCommitVote()
	}
	if !c.IsProposer() && c.currentState() < StatePreCommitted {
		c.lockQCAndProposal(msg)
		logger.Trace("acceptPreCommitted", "msg", msgTyp, "lockQC", msg.hash)
		c.sendCommitVote()
	}
	return nil
}

func (c *core) lockQCAndProposal(qc *QuorumCert) {
	c.current.SetPreCommittedQC(qc)
	c.current.SetState(StatePreCommitted)
	c.current.LockProposal()
}

func (c *core) sendCommitVote() {
	logger := c.newLogger()

	msgTyp := MsgTypeCommitVote
	vote := c.current.Vote()
	if vote == nil {
		logger.Error("Failed to send vote", "msg", msgTyp, "err", "current vote is nil")
		return
	}
	payload, err := Encode(vote)
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&Message{Code: msgTyp, Msg: payload})
	logger.Trace("sendCommitVote", "msg", msgTyp, "hash", vote.Digest)
}
