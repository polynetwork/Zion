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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
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
	// validateFn used to check block hash signature, it is an closure function which can be nil in unit test.
	if c.validateFn != nil {
		if addr, err := c.validateFn(lockedBlock.Hash(), data.CommittedSeal); err != nil {
			logger.Trace("Failed to check vote", "msg", code, "src", src, "err", err, "expect", src, "got", addr)
			return err
		}
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

		c.sendDecide(sealedBlock.Hash(), commitQC, seals)
	}

	return nil
}

func (c *core) sendDecide(block common.Hash, commitQC *QuorumCert, committedSeals [][]byte) {
	logger := c.newLogger()

	code := MsgTypeDecide
	msg := &Diploma{
		BlockHash:      block,
		CommitQC:       commitQC,
		CommittedSeals: committedSeals,
	}
	payload, err := Encode(msg)
	if err != nil {
		logger.Trace("Failed to encode", "msg", code, "err", err)
		return
	}
	c.broadcast(code, payload)

	logger.Trace("sendDecide", "msg", code, "node", commitQC.node)
}

// handleDecide repo receive MsgDecide and try to commit the final block.
func (c *core) handleDecide(data *Message) error {
	var (
		logger = c.newLogger()
		code   = data.Code
		src    = data.address
		msg    *Diploma
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

	// ensure commitQC is legal
	commitQC := msg.CommitQC
	if err := c.verifyQC(data, commitQC); err != nil {
		logger.Trace("Failed to verify commitQC", "msg", code, "src", src, "err", err)
		return err
	}

	// ensure the block hash is the correct one
	blockHash := msg.BlockHash
	lockedBlock := c.current.LockedBlock()
	if lockedBlock == nil {
		logger.Trace("Locked block is nil", "msg", code, "src", src)
		return errInvalidBlock
	} else if blockHash == common.EmptyHash || lockedBlock.Hash() != blockHash {
		logger.Trace("Failed to check block hash", "msg", code, "src", src, "expect block", lockedBlock.Hash(), "got", blockHash)
		return errInvalidBlock
	}

	if curNode := c.current.Node(); curNode == nil || curNode.Block == nil {
		logger.Trace("Current node is nil")
		return errInvalidNode
	} else if curNode.Hash() != commitQC.node {
		logger.Trace("Failed to check commitQC", "expect node", curNode.Hash(), "got", commitQC.node)
		return errInvalidQC
	} else if curNode.Block.Hash() != blockHash {
		logger.Trace("Failed to check node", "expect node block hash", curNode.Block.Hash(), "got", blockHash)
		return errInvalidBlock
	}

	if err := c.signer.VerifyCommittedSeal(c.valSet, msg.BlockHash, msg.CommittedSeals); err != nil {
		logger.Trace("Failed to verify committed seals", "msg", code, "src", src, "err", err)
		return errInvalidQC
	}
	logger.Trace("handleDecide", "msg", code, "src", src, "node", commitQC.node)

	// accept commitQC and commit block to miner
	if c.IsProposer() && c.currentState() == StateCommitted {
		if err := c.commit(c.current.LockedBlock()); err != nil {
			logger.Trace("Failed to commit proposal", "msg", code, "err", err)
			return err
		}
	}
	if !c.IsProposer() && c.currentState() == StateLocked {
		sealedBlock, err := c.backend.SealBlock(lockedBlock, msg.CommittedSeals)
		if err != nil {
			logger.Trace("Failed to assemble committed proposal", "msg", code, "err", err)
			return err
		}
		if err := c.acceptCommitQC(sealedBlock, commitQC); err != nil {
			logger.Trace("Failed to accept commitQC", "msg", code, "err", err)
			return err
		}
		if err := c.commit(sealedBlock); err != nil {
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

func (c *core) commit(sealedBlock *types.Block) error {
	if lockedBlock := c.current.LockedBlock(); lockedBlock == nil {
		return fmt.Errorf("locked block is nil")
	} else if lockedBlock.Hash() != sealedBlock.Hash() {
		return fmt.Errorf("expect locked block %v, got %v", lockedBlock.Hash(), sealedBlock.Hash())
	}

	if c.current.executed == nil || c.current.executed.Block == nil || c.current.executed.Block.Hash() != sealedBlock.Hash() {
		if c.IsProposer() {
			c.current.executed = &consensus.ExecutedBlock{Block: sealedBlock}
		} else {
			executed, err := c.backend.ExecuteBlock(sealedBlock)
			if err != nil {
				return fmt.Errorf("failed to execute block %v, err: %v", sealedBlock.Hash(), err)
			}
			executed.Block = sealedBlock
			c.current.executed = executed
		}
	}

	return c.backend.Commit(c.current.executed)
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
