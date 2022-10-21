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
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *core) handleCommitVote(data *Message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		vote   *Vote
		msgTyp = MsgTypeCommitVote
	)
	if err := data.Decode(&vote); err != nil {
		logger.Trace("Failed to decode", "msg", msgTyp, "src", src.Address(), "err", err)
		return errFailedDecodeCommitVote
	}
	if err := c.checkView(msgTyp, vote.View); err != nil {
		logger.Trace("Failed to check view", "msg", msgTyp, "src", src.Address(), "err", err)
		return err
	}
	if err := c.checkVote(vote); err != nil {
		logger.Trace("Failed to check vote", "msg", msgTyp, "src", src.Address(), "err", err)
		return err
	}
	if c.current.PreCommittedQC() == nil || vote.Digest != c.current.PreCommittedQC().Hash() {
		logger.Trace("Failed to check hash", "msg", msgTyp, "src", src.Address(), "got", vote.Digest)
		return errInvalidDigest
	}
	if err := c.checkMsgToProposer(); err != nil {
		logger.Trace("Failed to check proposer", "msg", msgTyp, "src", src.Address(), "err", err)
		return err
	}

	if err := c.current.AddCommitVote(data); err != nil {
		logger.Trace("Failed to add vote", "msg", msgTyp, "src", src.Address(), "err", err)
		return errAddPreCommitVote
	}

	logger.Trace("handleCommitVote", "msg", msgTyp, "src", src.Address(), "hash", vote.Digest)

	if size := c.current.CommitVoteSize(); size >= c.Q() && c.currentState() < StateCommitted {
		c.current.SetState(StateCommitted)
		c.current.SetCommittedQC(c.current.PreCommittedQC())
		logger.Trace("acceptCommit", "msg", msgTyp, "msgSize", size)

		c.sendDecide()
	}

	return nil
}

func (c *core) sendDecide() {
	logger := c.newLogger()

	msgTyp := MsgTypeDecide
	sub := c.current.CommittedQC()
	payload, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&Message{Code: msgTyp, Msg: payload})
	logger.Trace("sendDecide", "msg", msgTyp, "proposal", sub.Hash())
}

func (c *core) handleDecide(data *Message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		msg    *QuorumCert
		msgTyp = MsgTypeDecide
	)
	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to decode", "msg", msgTyp, "src", src.Address(), "err", err)
		return errFailedDecodeCommit
	}
	if err := c.checkView(msgTyp, msg.view); err != nil {
		logger.Trace("Failed to check view", "msg", msgTyp, "src", src.Address(), "err", err)
		return err
	}
	if err := c.checkMsgFromProposer(src); err != nil {
		logger.Trace("Failed to check proposer", "msg", msgTyp, "src", src.Address(), "err", err)
		return err
	}
	if err := c.checkPreCommittedQC(msg); err != nil {
		logger.Trace("Failed to check prepareQC", "msg", msgTyp, "src", src.Address(), "err", err)
		return err
	}
	if err := c.signer.VerifyQC(msg, c.valSet); err != nil {
		logger.Trace("Failed to check verify qc", "msg", msgTyp, "src", src.Address(), "err", err)
		return err
	}

	logger.Trace("handleDecide", "msg", msgTyp, "address", src.Address(), "proposal", msg.Hash())

	if c.IsProposer() && c.currentState() == StateCommitted {
		if err := c.backend.Commit(c.current.Proposal()); err != nil {
			logger.Trace("Failed to commit proposal", "err", err)
			return err
		}
	}

	if !c.IsProposer() && c.currentState() >= StatePreCommitted && c.currentState() < StateCommitted {
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
