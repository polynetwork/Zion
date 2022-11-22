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
)

func (c *core) handlePreCommitVote(data *Message) error {
	var (
		logger = c.newLogger()
		code   = MsgTypePreCommitVote
		src    = data.address
		vote   = common.BytesToHash(data.Msg)
	)

	if err := c.checkView(data.View); err != nil {
		logger.Trace("Failed to check view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkVote(data, vote); err != nil {
		logger.Trace("Failed to check vote", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkMsgDest(); err != nil {
		logger.Trace("Failed to check proposal", "msg", code, "src", src, "err", err)
		return err
	}

	if err := c.current.AddPreCommitVote(data); err != nil {
		logger.Trace("Failed to add vote", "msg", code, "src", src, "err", err)
		return errAddPreCommitVote
	}

	logger.Trace("handlePreCommitVote", "msg", code, "src", src, "hash", vote)

	if size := c.current.PreCommitVoteSize(); size >= c.Q() && c.currentState() < StateLocked {
		lockQC, err := c.messages2qc(code)
		if err != nil {
			logger.Trace("Failed to assemble lockQC", "msg", code, "err", err)
			return err
		}
		if err := c.acceptLockQC(lockQC); err != nil {
			logger.Trace("Failed to accept lockQC", "msg", code, "err", err)
			return err
		}

		logger.Trace("acceptPreCommitted", "msg", code, "msgSize", size)
		c.sendCommit(lockQC)
	}
	return nil
}

func (c *core) sendCommit(lockQC *QuorumCert) {
	logger := c.newLogger()

	code := MsgTypeCommit
	payload, err := Encode(lockQC)
	if err != nil {
		logger.Error("Failed to encode", "msg", code, "err", err)
		return
	}
	c.broadcast(code, payload)
	logger.Trace("sendCommit", "msg", code, "node", lockQC.node)
}

func (c *core) handleCommit(data *Message) error {
	logger := c.newLogger()

	var (
		lockQC  *QuorumCert
		code = MsgTypeCommit
		src  = data.address
	)
	if err := data.Decode(&lockQC); err != nil {
		logger.Trace("Failed to decode", "msg", code, "src", src, "err", err)
		return errFailedDecodeCommit
	}
	if err := c.checkView(lockQC.view); err != nil {
		logger.Trace("Failed to check view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkMsgSource(src); err != nil {
		logger.Trace("Failed to check proposer", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.verifyQC(data, lockQC); err != nil {
		logger.Trace("Failed to check verify qc", "msg", code, "src", src, "err", err)
		return err
	}

	logger.Trace("handleCommit", "msg", code, "src", src, "lockQC", lockQC.node)

	// accept lockQC
	if c.IsProposer() && c.currentState() < StateCommitted {
		c.sendVote(MsgTypeCommitVote)
	}
	if !c.IsProposer() && c.currentState() < StateLocked {
		if err := c.acceptLockQC(lockQC); err != nil {
			logger.Trace("Failed to accept lockQC", "msg", code, "err", err)
			return err
		}
		logger.Trace("acceptLockQC", "msg", code, "lockQC", lockQC.node)

		c.sendVote(MsgTypeCommitVote)
	}

	return nil
}

func (c *core) acceptLockQC(qc *QuorumCert) error {
	if err := c.current.Lock(qc); err != nil {
		return err
	}
	c.current.SetState(StateLocked)
	return nil
}
