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

func (c *core) handlePreCommitVote(data *Message) error {
	logger := c.newLogger()

	var (
		vote   *Vote
		code = MsgTypePreCommitVote
		src    = data.address
	)
	if err := data.Decode(&vote); err != nil {
		logger.Trace("Failed to decode", "type", code, "src", src, "err", err)
		return errFailedDecodePreCommitVote
	}
	if err := c.checkView(code, data.View); err != nil {
		logger.Trace("Failed to check view", "type", code, "src", src, "err", err)
		return err
	}
	if err := c.checkVote(data, vote); err != nil {
		logger.Trace("Failed to check vote", "type", code, "src", src, "err", err)
		return err
	}
	if err := c.checkProposal(vote.Digest); err != nil {
		logger.Trace("Failed to check hash", "type", code, "src", src, "expect vote", vote.Digest)
		return errInvalidDigest
	}
	if err := c.checkMsgToProposer(); err != nil {
		logger.Trace("Failed to check proposal", "type", code, "src", src, "err", err)
		return err
	}

	if err := c.current.AddPreCommitVote(data); err != nil {
		logger.Trace("Failed to add vote", "type", code, "src", src, "err", err)
		return errAddPreCommitVote
	}

	logger.Trace("handlePreCommitVote", "msg", code, "src", src, "hash", vote.Digest)

	if size := c.current.PreCommitVoteSize(); size >= c.Q() && c.currentState() < StatePreCommitted {
		if preCommitQC, err := c.messages2qc(c.proposer(), vote.Digest, c.current.PreCommitVotes()); err != nil {
			logger.Trace("Failed to assemble preCommitQC", "type", code, "err", err)
			return err
		} else {
			c.lockQCAndProposal(preCommitQC)
		}

		logger.Trace("acceptPreCommitted", "msg", code, "msgSize", size)
		c.sendCommit()
	}
	return nil
}

func (c *core) sendCommit() {
	logger := c.newLogger()

	code := MsgTypeCommit
	sub := c.current.PreCommittedQC()
	payload, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "msg", code, "err", err)
		return
	}
	c.broadcast(code, payload)
	logger.Trace("sendCommit", "msg", code, "proposal", sub.hash)
}

func (c *core) handleCommit(data *Message) error {
	logger := c.newLogger()

	var (
		msg    *QuorumCert
		code = MsgTypeCommit
		src    = data.address
	)
	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to decode", "msg", code, "src", src, "err", err)
		return errFailedDecodeCommit
	}
	if err := c.checkView(MsgTypeCommit, msg.view); err != nil {
		logger.Trace("Failed to check view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkMsgFromProposer(src); err != nil {
		logger.Trace("Failed to check proposer", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkPrepareQC(msg); err != nil {
		logger.Trace("Failed to check prepareQC", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.verifyVoteQC(msg.hash, msg); err != nil {
		logger.Trace("Failed to check verify qc", "msg", code, "src", src, "err", err)
		return err
	}

	logger.Trace("handleCommit", "msg", code, "src", src, "proposal", msg.hash)

	if c.IsProposer() && c.currentState() < StateCommitted {
		c.sendCommitVote()
	}
	if !c.IsProposer() && c.currentState() < StatePreCommitted {
		c.lockQCAndProposal(msg)
		logger.Trace("acceptPreCommitted", "msg", code, "lockQC", msg.hash)
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

	code := MsgTypeCommitVote
	vote := c.current.Vote()
	if vote == nil {
		logger.Error("Failed to send vote", "msg", code, "err", "current vote is nil")
		return
	}
	payload, err := Encode(vote)
	if err != nil {
		logger.Error("Failed to encode", "msg", code, "err", err)
		return
	}
	c.broadcast(code, payload)
	logger.Trace("sendCommitVote", "msg", code, "hash", vote.Digest)
}
