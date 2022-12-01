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

	"github.com/ethereum/go-ethereum/consensus"

	"github.com/ethereum/go-ethereum/core/types"
)

// sendPrepare leader send message of prepare(view, node, highQC)
func (c *core) sendPrepare() {

	// filter incorrect proposer and state
	if !c.IsProposer() || c.currentState() != StateHighQC {
		return
	}

	var (
		block  *types.Block
		code   = MsgTypePrepare
		highQC = c.current.HighQC()
		logger = c.newLogger()
	)

	// fetch block with locked node or miner pending request
	if lockedBlock := c.current.LockedBlock(); lockedBlock != nil {
		if lockedBlock.NumberU64() != c.HeightU64() {
			logger.Trace("Locked block height invalid", "msg", code, "expect", c.HeightU64(), "got", lockedBlock.NumberU64())
			return
		}
		block = lockedBlock
		logger.Trace("Reuse lock block", "msg", code, "hash", block.Hash(), "number", block.NumberU64())
	} else {
		request := c.current.PendingRequest()
		if request == nil || request.block == nil || request.block.NumberU64() != c.HeightU64() {
			logger.Trace("Pending request invalid", "msg", code)
			return
		} else {
			block = c.current.PendingRequest().block
			logger.Trace("Use pending request", "msg", code, "hash", block.Hash(), "number", block.NumberU64())
		}
	}

	// consensus spent time always less than a block period, waiting for `delay` time to catch up the system time.
	// todo(fuk): waiting in `startNewRound`
	if block.Time() > uint64(time.Now().Unix()) {
		delay := time.Unix(int64(block.Time()), 0).Sub(time.Now())
		time.Sleep(delay)
		logger.Trace("delay to broadcast proposal", "msg", code, "time", delay.Milliseconds())
	}

	// assemble message as formula: MSG(view, node, prepareQC)
	parent := highQC.node
	node := NewNode(parent, block)
	prepare := NewSubject(node, highQC)
	payload, err := Encode(prepare)
	if err != nil {
		logger.Trace("Failed to encode", "msg", code, "err", err)
		return
	}

	// store the node before `handlePrepare` to prevent the replica from receiving the message and voting earlier
	// than the leader, and finally causing `handlePrepareVote` to fail.
	if err := c.current.SetNode(node); err != nil {
		logger.Trace("Failed to set node", "msg", code, "err", err)
		return
	}

	c.broadcast(code, payload)
	logger.Trace("sendPrepare", "msg", code, "node", node.Hash(), "block", block.Hash())
}

// handlePrepare implement description as follow:
// ```
//  repo wait for message m : matchingMsg(m, prepare, curView) from leader(curView)
//	if m.node extends from m.justify.node ∧
//	safeNode(m.node, m.justify) then
//	send voteMsg(prepare, m.node, ⊥) to leader(curView)
// ```
func (c *core) handlePrepare(data *Message) error {
	var (
		logger = c.newLogger()
		code   = data.Code
		src    = data.address
		msg    *Subject
	)

	// check message
	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to decode", "msg", code, "src", src, "err", err)
		return errFailedDecodePrepare
	}
	if err := c.checkView(data.View); err != nil {
		logger.Trace("Failed to check view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkMsgSource(src); err != nil {
		logger.Trace("Failed to check proposer", "msg", code, "src", src, "err", err)
		return err
	}

	// local node is nil before `handlePrepare`, only check fields here.
	node := msg.Node
	if err := c.checkNode(node, false); err != nil {
		logger.Trace("Failed to check node", "msg", code, "src", src, "err", err)
		return err
	}

	// ensure remote block is legal.
	block := node.Block
	if err := c.checkBlock(block); err != nil {
		logger.Trace("Failed to check block", "msg", code, "src", src, "err", err)
		return err
	}
	if duration, err := c.backend.Verify(block, false); err != nil {
		logger.Trace("Failed to verify unsealed proposal", "msg", code, "src", src, "err", err, "duration", duration)
		return errVerifyUnsealedProposal
	}
	if err := c.executeBlock(block); err != nil {
		logger.Trace("Failed to execute block", "msg", code, "src", src, "err", err)
		return err
	}

	// safety and liveness rules judgement.
	highQC := msg.QC
	if err := c.verifyQC(data, highQC); err != nil {
		logger.Trace("Failed to verify highQC", "msg", code, "src", src, "err", err, "highQC", highQC)
		return err
	}
	if err := c.extend(node, highQC); err != nil {
		logger.Trace("Failed to check extend", "msg", code, "src", src, "err", err)
		return errExtend
	}
	if err := c.safeNode(node, highQC); err != nil {
		logger.Trace("Failed to check safeNode", "msg", code, "src", src, "err", err)
		return errSafeNode
	}

	logger.Trace("handlePrepare", "msg", code, "src", src, "node", node.Hash(), "block", block.Hash())

	// accept msg info, DONT persist node before accept `prepareQC`
	if c.IsProposer() && c.currentState() == StateHighQC {
		c.sendVote(MsgTypePrepareVote, node.Hash())
	}
	if !c.IsProposer() && c.currentState() < StateHighQC {
		if err := c.current.SetNode(node); err != nil {
			logger.Trace("Failed to set node", "msg", code, "err", err)
			return err
		}
		c.setCurrentState(StateHighQC)
		logger.Trace("acceptHighQC", "msg", code, "highQC", highQC.node, "node", node.Hash())
		c.sendVote(MsgTypePrepareVote, node.Hash())
	}

	return nil
}

// proposer do not need execute block again after miner.worker commitNewWork.
func (c *core) executeBlock(block *types.Block) error {
	if c.IsProposer() {
		c.current.executed = &consensus.ExecutedBlock{Block: block}
		return nil
	}

	executed, err := c.backend.ExecuteBlock(block)
	if err != nil {
		return err
	}
	c.current.executed = executed
	return nil
}

// remote node's parent should equals to highQC's node
func (c *core) extend(node *Node, highQC *QuorumCert) error {
	if highQC == nil || highQC.view == nil {
		return errInvalidQC
	}
	if highQC.node != node.Parent {
		return fmt.Errorf("expect parent %v, got %v", highQC.node, node.Parent)
	}
	return nil
}

// proposal extend lockQC `OR` highQC.view > lockQC.view
func (c *core) safeNode(node *Node, highQC *QuorumCert) error {
	if highQC == nil || highQC.view == nil {
		return errInvalidQC
	}

	// skip epoch start block
	lockQC := c.current.LockQC()
	if lockQC == nil {
		c.logger.Warn("LockQC be nil should only happen at `startUp`")
		return nil
	}

	if highQC.view.Cmp(lockQC.view) > 0 || node.Parent == lockQC.node {
		return nil
	}

	return errSafeNode
}
