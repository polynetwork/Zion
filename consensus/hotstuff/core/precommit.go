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
	"github.com/ethereum/go-ethereum/core/types"
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

	if size := c.current.PrepareVoteSize(); size >= c.Q() && c.currentState() < StatePrepared {
		prepareQC, err := c.messages2qc(code)
		if err != nil {
			logger.Trace("Failed to assemble prepareQC", "msg", code, "err", err)
			return errInvalidQC
		}
		// save node before repo
		if err := c.acceptPrepare(prepareQC, c.current.Node()); err != nil {
			logger.Trace("Failed to accept prepareQC", "msg", code, "err", err)
			return err
		}
		c.sendPreCommit(c.current.Node(), prepareQC)
	}

	return nil
}

func (c *core) sendPreCommit(node *Node, prepareQC *QuorumCert) {
	logger := c.newLogger()

	code := MsgTypePreCommit
	msg := NewSubject(node, prepareQC)
	payload, err := Encode(msg)
	if err != nil {
		logger.Trace("Failed to encode", "msg", code, "err", err)
		return
	}
	c.broadcast(code, payload)
	logger.Trace("sendPreCommit", "msg", code, "node", msg.Node.Hash())
}

func (c *core) handlePreCommit(data *Message) error {
	logger := c.newLogger()

	var (
		msg       *Subject
		code      = MsgTypePreCommit
		src       = data.address
		node      *Node
		prepareQC *QuorumCert
		block     *types.Block
	)

	// check parameters
	if err := data.Decode(&msg); err != nil {
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
	if err := c.checkSubject(msg); err != nil {
		logger.Trace("Failed to check subject", "msg", code, "src", src, "err", err)
		return err
	} else {
		node = msg.Node
		prepareQC = msg.QC
		block = node.Block
	}

	if err := c.checkNode(node); err != nil {
		logger.Trace("Failed to check node", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkBlock(block); err != nil {
		logger.Trace("Failed to check block", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.verifyQC(data, prepareQC); err != nil {
		logger.Trace("Failed to verify prepareQC", "msg", code, "src", src, "err", err)
		return err
	}
	if _, err := c.backend.Verify(block, false); err != nil {
		logger.Trace("Failed to check verify proposal", "msg", code, "src", src, "err", err)
		return err
	}

	logger.Trace("handlePreCommit", "msg", code, "src", src, "hash", node.Hash())

	// accept msg info and state
	if c.IsProposer() && c.currentState() < StateLocked {
		c.sendVote(MsgTypePreCommitVote)
	}
	if !c.IsProposer() && c.currentState() < StatePrepared {
		if err := c.acceptPrepare(prepareQC, node); err != nil {
			logger.Trace("Failed to accept prepareQC", "msg", code, "err", err)
			return err
		}
		logger.Trace("acceptPrepare", "msg", code, "prepareQC", prepareQC.node)
		c.sendVote(MsgTypePreCommitVote)
	}

	return nil
}

func (c *core) acceptPrepare(prepareQC *QuorumCert, node *Node) error {
	if err := c.current.SetNode(node); err != nil {
		return err
	}
	if err := c.current.SetPrepareQC(prepareQC); err != nil {
		return err
	}
	c.current.SetState(StatePrepared)
	return nil
}
