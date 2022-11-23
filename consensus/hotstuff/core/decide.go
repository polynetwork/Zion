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

// handleCommitVote implement description as follow:
// ```
// leader wait for (n n f) votes: V ← {v | matchingMsg(v, commit, curView)}
//	commitQC ← QC(V )
//	broadcast Msg(decide, ⊥, commitQC )
// ```
func (c *core) handleCommitVote(data *Message) error {
	var (
		logger = c.newLogger()
		vote   = common.BytesToHash(data.Msg)
		code   = data.Code
		src    = data.address
	)

	// check message
	if err := c.checkView(data.View); err != nil {
		logger.Trace("Failed to check view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkMsgDest(); err != nil {
		logger.Trace("Failed to check proposer", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkVote(data, vote); err != nil {
		logger.Trace("Failed to check vote", "msg", code, "src", src, "err", err)
		return err
	}

	// check locked block's committed seals
	lockedBlock := c.current.LockedBlock()
	if lockedBlock == nil {
		logger.Trace("Failed to get lockBlock", "msg", code, "src", src, "err", "block is nil")
		return errInvalidNode
	}
	if addr, err := c.validateFn(lockedBlock.Hash(), data.CommittedSeal, true); err != nil {
		logger.Trace("Failed to check vote", "msg", code, "src", src, "err", err, "expect", src, "got", addr)
		return err
	}

	// queue vote into messageSet to ensure that at least 2/3 validator vote at the same step.
	if err := c.current.AddCommitVote(data); err != nil {
		logger.Trace("Failed to add vote", "msg", code, "src", src, "err", err)
		return errAddPreCommitVote
	}

	logger.Trace("handleCommitVote", "msg", code, "src", src, "hash", vote)

	// assemble committed signatures to reorg the locked block, and create `commitQC` at the same time.
	if size := c.current.CommitVoteSize(); size >= c.Q() && c.currentState() == StateLocked {
		seals := c.current.GetCommittedSeals(size)
		sealedBlock, err := c.backend.SealBlock(lockedBlock, seals)
		if err != nil {
			logger.Trace("Failed to assemble committed proposal", "msg", code, "err", err)
			return err
		}
		commitQC, err := c.messages2qc(code)
		if err != nil {
			logger.Trace("Failed to assemble commitQC", "msg", code, "err", err)
			return err
		}
		if err := c.acceptCommitQC(sealedBlock, commitQC); err != nil {
			logger.Trace("Failed to accept commitQC")
		}
		logger.Trace("acceptCommit", "msg", code, "msgSize", size)

		c.sendDecide(commitQC)
	}

	return nil
}

func (c *core) sendDecide(commitQC *QuorumCert) {
	logger := c.newLogger()

	code := MsgTypeDecide
	msg := &Subject{
		Node: c.current.Node(),
		QC:   commitQC,
	}
	payload, err := Encode(msg)
	if err != nil {
		logger.Trace("Failed to encode", "msg", code, "err", err)
		return
	}
	c.broadcast(code, payload)

	logger.Trace("sendDecide", "msg", code, "node", msg.Node.Hash())
}

// handleDecide repo receive MsgDecide and try to commit the final block.
func (c *core) handleDecide(data *Message) error {
	var (
		logger = c.newLogger()
		code   = data.Code
		src    = data.address
		msg    *Subject
	)

	// check message
	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to decode", "msg", code, "src", src, "err", err)
		return errFailedDecodeCommit
	}
	if err := c.checkView(data.View); err != nil {
		logger.Trace("Failed to check view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkMsgSource(src); err != nil {
		logger.Trace("Failed to check proposer", "msg", code, "src", src, "err", err)
		return err
	}

	// ensure that remote node equals to locked node.
	if err := c.checkNode(msg.Node, true); err != nil {
		logger.Trace("Failed to check node", "msg", code, "src", src, "err", err)
		return err
	}

	// ensure commitQC is legal
	commitQC := msg.QC
	if err := c.verifyQC(data, commitQC); err != nil {
		logger.Trace("Failed to verify qc", "msg", code, "src", src, "err", err)
		return err
	}

	// validate block safety rule, and verify block with committed seals.
	sealedBlock := msg.Node.Block
	if err := c.checkBlock(sealedBlock); err != nil {
		logger.Trace("Failed to check block", "msg", code, "src", src, "err", err)
		return err
	}
	if _, err := c.backend.Verify(sealedBlock, true); err != nil {
		logger.Trace("Failed to verify block")
	}
	logger.Trace("handleDecide", "msg", code, "src", src, "node", commitQC.node)

	// accept commitQC and commit block to miner
	if c.IsProposer() && c.currentState() == StateCommitted {
		if err := c.backend.Commit(sealedBlock); err != nil {
			logger.Trace("Failed to commit proposal", "msg", code, "err", err)
			return err
		}
	}
	if !c.IsProposer() && c.currentState() == StateLocked {
		if err := c.acceptCommitQC(sealedBlock, commitQC); err != nil {
			logger.Trace("Failed to accept commitQC", "msg", code, "err", err)
			return err
		}
		if err := c.backend.Commit(sealedBlock); err != nil {
			logger.Trace("Failed to commit proposal", "err", err)
			return err
		}
	}

	c.startNewRound(common.Big0)
	return nil
}

func (c *core) acceptCommitQC(sealedBlock *types.Block, commitQC *QuorumCert) error {
	if err := c.current.SetSealedBlock(sealedBlock); err != nil {
		return err
	}
	if err := c.current.SetCommittedQC(commitQC); err != nil {
		return err
	}
	c.current.SetState(StateCommitted)
	return nil
}

// handleFinalCommitted start new round if consensus engine accept notify signal from miner.worker.
// signals should be related with sync header or body. in fact, we DONT need this function to start an new round,
// because that the function `startNewRound` will sync header to preparing new consensus round args.
// we just kept it here for backup.
func (c *core) handleFinalCommitted(header *types.Header) error {
	logger := c.newLogger()
	if height := header.Number.Uint64(); height >= c.current.HeightU64() {
		logger.Trace("handleFinalCommitted", "height", height)
		c.startNewRound(common.Big0)
	}
	return nil
}
