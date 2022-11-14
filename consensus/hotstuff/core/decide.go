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

func (c *core) handleCommitVote(data *Message) error {
	logger := c.newLogger()

	var (
		vote = common.BytesToHash(data.Msg)
		code = MsgTypeCommitVote
		src  = data.address
	)
	if err := c.checkView(code, data.View); err != nil {
		logger.Trace("Failed to check view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkVote(data, vote); err != nil {
		logger.Trace("Failed to check vote", "msg", code, "src", src, "err", err)
		return err
	}
	if c.current.PreCommittedQC() == nil || vote != c.current.PreCommittedQC().hash {
		logger.Trace("Failed to check hash", "msg", code, "src", src, "got", vote)
		return errInvalidDigest
	}
	if err := c.checkMsgToProposer(); err != nil {
		logger.Trace("Failed to check proposer", "msg", code, "src", src, "err", err)
		return err
	}

	// check committed seal
	if addr, err := c.validateFn(vote, data.CommittedSeal, true); err != nil {
		logger.Trace("Failed to check vote", "msg", code, "src", src, "err", err, "expect", src, "got", addr)
		return err
	}

	if err := c.current.AddCommitVote(data); err != nil {
		logger.Trace("Failed to add vote", "msg", code, "src", src, "err", err)
		return errAddPreCommitVote
	}

	logger.Trace("handleCommitVote", "msg", code, "src", src, "hash", vote)

	if size := c.current.CommitVoteSize(); size >= c.Q() && c.currentState() < StateCommitted {

		seals := c.current.GetCommittedSeals(size)
		newProposal, err := c.backend.PreCommit(c.current.Proposal(), seals)
		if err != nil {
			logger.Trace("Failed to assemble committed proposal", "msg", code, "err", err)
			return err
		}
		commitQC, err := c.messages2qc(c.proposer(), newProposal.Hash(), c.current.CommitVotes())
		if err != nil {
			logger.Trace("Failed to assemble commitQC", "msg", code, "err", err)
			return err
		}
		c.current.SetProposal(newProposal)
		c.current.SetState(StateCommitted)
		c.current.SetCommittedQC(commitQC)
		logger.Trace("acceptCommit", "msg", code, "msgSize", size)

		c.sendDecide()
	}

	return nil
}

func (c *core) sendDecide() {
	logger := c.newLogger()

	code := MsgTypeDecide
	msg := &Subject{
		Proposal: c.current.Proposal(),
		QC:       c.current.CommittedQC(),
	}
	payload, err := Encode(msg)
	if err != nil {
		logger.Trace("Failed to encode", "msg", code, "err", err)
		return
	}
	c.broadcast(code, payload)

	logger.Trace("sendDecide", "msg", code, "proposal", msg.Proposal.Hash())
}

func (c *core) handleDecide(data *Message) error {
	logger := c.newLogger()

	var (
		msg  *Subject
		code = MsgTypeDecide
		src  = data.address
	)
	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to decode", "msg", code, "src", src, "err", err)
		return errFailedDecodeCommit
	}
	commitQC := msg.QC

	if err := c.checkView(code, data.View); err != nil {
		logger.Trace("Failed to check view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkMsgFromProposer(src); err != nil {
		logger.Trace("Failed to check proposer", "msg", code, "src", src, "err", err)
		return err
	}
	//if err := c.checkPreCommittedQC(msg.CommitQC); err != nil {
	//	logger.Trace("Failed to check prepareQC", "msg", code, "src", src, "err", err)
	//	return err
	//}
	if err := c.verifyVoteQC(commitQC.hash, commitQC); err != nil {
		logger.Trace("Failed to check verify qc", "msg", code, "src", src, "err", err)
		return err
	}

	logger.Trace("handleDecide", "msg", code, "src", src, "proposal", commitQC.hash)

	if c.IsProposer() && c.currentState() == StateCommitted {
		if err := c.backend.Commit(c.current.Proposal()); err != nil {
			logger.Trace("Failed to commit proposal", "err", err)
			return err
		}
	}

	if !c.IsProposer() && c.currentState() >= StateLocked && c.currentState() < StateCommitted {
		c.current.SetState(StateCommitted)
		c.current.SetCommittedQC(c.current.PreCommittedQC())
		if err := c.backend.Commit(c.current.Proposal()); err != nil {
			logger.Trace("Failed to commit proposal", "err", err)
			return err
		}
	}

	c.startNewRound(common.Big0)
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
