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

func (c *core) handlePrepareVote(data *Message) error {

	var (
		logger = c.newLogger()
		vote   = common.BytesToHash(data.Msg)
		code   = MsgTypePrepareVote
		src    = data.address
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
		logger.Trace("Failed to check proposer", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.current.AddPrepareVote(data); err != nil {
		logger.Trace("Failed to add vote", "msg", code, "src", src, "err", err)
		return errAddPrepareVote
	}

	logger.Trace("handlePrepareVote", "msg", code, "src", src, "vote", vote)

	if size := c.current.PrepareVoteSize(); size >= c.Q() && c.currentState() == StateHighQC {
		prepareQC, err := c.messages2qc(code)
		if err != nil {
			logger.Trace("Failed to assemble prepareQC", "msg", code, "err", err)
			return errInvalidQC
		}
		if err := c.acceptPrepareQC(prepareQC); err != nil {
			logger.Trace("Failed to accept prepareQC", "msg", code, "err", err)
			return err
		}
		logger.Trace("acceptPrepareQC", "msg", code, "prepareQC", prepareQC.node)

		c.sendPreCommit(prepareQC)
	}

	return nil
}

func (c *core) sendPreCommit(prepareQC *QuorumCert) {
	logger := c.newLogger()

	code := MsgTypePreCommit
	payload, err := Encode(prepareQC)
	if err != nil {
		logger.Trace("Failed to encode", "msg", code, "err", err)
		return
	}
	c.broadcast(code, payload)
	logger.Trace("sendPreCommit", "msg", code, "node", prepareQC.node)
}

func (c *core) handlePreCommit(data *Message) error {
	logger := c.newLogger()

	var (
		code      = MsgTypePreCommit
		src       = data.address
		prepareQC *QuorumCert
	)

	// check parameters
	if err := data.Decode(&prepareQC); err != nil {
		logger.Trace("Failed to check decode", "msg", code, "src", src, "err", err)
		return errFailedDecodePreCommit
	}
	if err := c.checkView(data.View); err != nil {
		logger.Trace("Failed to check view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkMsgSource(src); err != nil {
		logger.Trace("Failed to check proposer", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.verifyQC(data, prepareQC); err != nil {
		logger.Trace("Failed to verify prepareQC", "msg", code, "src", src, "err", err)
		return err
	}

	logger.Trace("handlePreCommit", "msg", code, "src", src, "prepareQC", prepareQC.node)

	// accept msg info and state
	if c.IsProposer() && c.currentState() == StatePrepared {
		c.sendVote(MsgTypePreCommitVote)
	}
	if !c.IsProposer() && c.currentState() == StateHighQC {
		if err := c.acceptPrepareQC(prepareQC); err != nil {
			logger.Trace("Failed to accept prepareQC", "msg", code, "err", err)
			return err
		}
		logger.Trace("acceptPrepareQC", "msg", code, "prepareQC", prepareQC.node)
		c.sendVote(MsgTypePreCommitVote)
	}

	return nil
}

func (c *core) acceptPrepareQC(prepareQC *QuorumCert) error {
	if err := c.current.SetNode(c.current.Node()); err != nil {
		return err
	}
	if err := c.current.SetPrepareQC(prepareQC); err != nil {
		return err
	}
	c.current.SetState(StatePrepared)
	return nil
}
