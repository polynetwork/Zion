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

	"github.com/ethereum/go-ethereum/core/types"
)

func (c *core) sendPrepare() {
	logger := c.newLogger()

	if !c.IsProposer() {
		return
	}

	var (
		code   = MsgTypePrepare
		highQC = c.current.HighQC()
	)

	if c.currentState() != StateHighQC {
		return
	}
	if c.current.PendingRequest() == nil || c.current.PendingRequest().block == nil {
		return
	}

	// todo(fuk): if block is locked, use lock block instead
	request := c.current.PendingRequest()
	if request.block.NumberU64() != c.current.HeightU64() {
		logger.Trace("Failed to send prepare", "msg", code, "err", "request height invalid")
		return
	}
	// todo(fuk): request must extend lastProposal
	node := NewNode(c.current.highQC.node, request.block)
	prepare := &Subject{
		Node: node,
		QC:   highQC,
	}
	payload, err := Encode(prepare)
	if err != nil {
		logger.Trace("Failed to encode", "msg", code, "err", err)
		return
	}

	// consensus spent time always less than a block period, waiting for `delay` time to catch up the system time.
	delay := time.Unix(int64(prepare.Node.Block.Time()), 0).Sub(time.Now())
	time.Sleep(delay)
	logger.Trace("delay to broadcast proposal", "msg", code, "time", delay.Milliseconds())

	c.broadcast(code, payload)
	logger.Trace("sendPrepare", "msg", code, "node", prepare.Node.Hash(), "block", request.block.Hash())
}

func (c *core) handlePrepare(data *Message) error {
	logger := c.newLogger()

	var (
		msg    *Subject
		code   = MsgTypePrepare
		src    = data.address
		node   *Node
		highQC *QuorumCert
		block  *types.Block
	)

	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to decode", "msg", code, "src", src, "err", err)
		return errFailedDecodePrepare
	}
	if err := c.checkSubject(msg); err != nil {
		logger.Trace("Failed to check subject", "msg", code, "src", src, "err", err)
		return err
	} else {
		node = msg.Node
		highQC = msg.QC
		block = node.Block
	}

	if err := c.checkView(data.View); err != nil {
		logger.Trace("Failed to check view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkMsgFromProposer(src); err != nil {
		logger.Trace("Failed to check proposer", "msg", code, "src", src, "err", err)
		return err
	}

	//if err := c.checkProposalView(msg.Proposal, data.View); err != nil {
	//	logger.Trace("Failed to check proposal and msg view", "msg", code, "src", src, "err", err)
	//	return err
	//}
	if err := c.verifyQC(data, highQC); err != nil {
		logger.Trace("Failed to verify highQC", "msg", code, "src", src, "err", err, "highQC", highQC)
		return err
	}
	//// locked proposal hash should be equal to msg proposal hash
	//if c.current.IsProposalLocked() {
	//	if expect := c.current.Proposal().Hash(); expect != msg.Proposal.Hash() {
	//		logger.Trace("Failed to check lock proposal", "msg", code, "src", src, "expect hash", expect, "got", msg.Proposal.Hash())
	//		return errLockProposal
	//	}
	//}

	// todo(fuk): 先验证一下区块头，但是不预执行区块，等到后续安全性及活性得到确认的情况下再预执行区块
	if _, err := c.backend.Verify(block, false); err != nil {
		logger.Trace("Failed to verify unsealed proposal", "msg", code, "src", src, "err", err)
		return errVerifyUnsealedProposal
	}

	if err := c.extend(node, highQC); err != nil {
		logger.Trace("Failed to check extend", "msg", code, "src", src, "err", err)
		return errExtend
	}
	if err := c.safeNode(node, highQC); err != nil {
		logger.Trace("Failed to check safeNode", "msg", code, "src", src, "err", err)
		return errSafeNode
	}

	if err := c.preExecuteBlock(block); err != nil {
		logger.Trace("Failed to pre-execute block", "msg", code, "src", src, "err", err)
		return err
	}

	logger.Trace("handlePrepare", "msg", code, "src", src, "node", node.Hash(), "block", block.Hash())

	// leader accept proposal
	if c.IsProposer() && c.currentState() < StatePrepared {
		if err := c.current.SetNode(node); err != nil {
			logger.Trace("Failed to set proposal", "msg", code, "err", err)
			return err
		}
		c.sendPrepareVote()
	}

	// repo accept proposal and high qc
	// todo(fuk): 是否真的需要stateHighQC
	if !c.IsProposer() && c.currentState() < StateHighQC {
		if err := c.current.SetNode(node); err != nil {
			logger.Trace("Failed to set proposal", "msg", code, "err", err)
			return err
		}
		c.setCurrentState(StateHighQC)
		logger.Trace("acceptHighQC", "msg", code, "highQC", highQC.node)

		c.sendPrepareVote()
	}

	return nil
}

func (c *core) sendPrepareVote() {
	logger := c.newLogger()

	code := MsgTypePrepareVote
	vote := c.current.Vote()
	c.broadcast(code, vote.Bytes())
	logger.Trace("sendPrepareVote", "msg", code, "hash", vote)
}

func (c *core) extend(node *Node, highQC *QuorumCert) error {
	if highQC == nil || highQC.view == nil {
		return errInvalidQC
	}
	// if msgPrepare.proposal is locked, the proposal hash should be equal to qc.hash
	if highQC.node != node.Parent {
		return fmt.Errorf("expect parent %v, got %v", highQC.node, node.Parent)
	}
	return nil
}

// proposal extend lockQC `OR` hiqhQC.view > lockQC.view
func (c *core) safeNode(node *Node, highQC *QuorumCert) error {
	if highQC == nil || highQC.view == nil {
		return errSafeNode
	}

	// skip genesis block
	lockQC := c.current.lockQC
	if lockQC == nil {
		if node.Block.NumberU64() == 1 {
			return nil
		} else {
			return errSafeNode
		}
	}

	if highQC.view.Cmp(lockQC.view) > 0 || node.Parent == lockQC.node {
		return nil
	}

	return errSafeNode
}
