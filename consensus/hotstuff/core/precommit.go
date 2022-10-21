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
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) handlePrepareVote(data *Message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		vote   *Vote
		msgTyp = MsgTypePrepareVote
	)
	if err := data.Decode(&vote); err != nil {
		logger.Trace("Failed to decode", "msg", msgTyp, "err", err)
		return errFailedDecodePrepareVote
	}
	if err := c.checkView(msgTyp, vote.View); err != nil {
		logger.Trace("Failed to check view", "msg", msgTyp, "err", err)
		return err
	}
	if err := c.checkVote(vote); err != nil {
		logger.Trace("Failed to check vote", "msg", msgTyp, "err", err)
		return err
	}
	if err := c.checkMsgToProposer(); err != nil {
		logger.Trace("Failed to check proposer", "msg", msgTyp, "err", err)
		return err
	}
	if err := c.current.AddPrepareVote(data); err != nil {
		logger.Trace("Failed to add vote", "msg", msgTyp, "err", err)
		return errAddPrepareVote
	}

	// check committed seal
	if addr, err := c.validateFn(vote.Digest[:], data.CommittedSeal); err != nil {
		logger.Trace("Failed to check vote", "msg", msgTyp, "err", err)
	} else if addr != src.Address() {
		logger.Trace("Failed to check vote", "msg", msgTyp, "expect", src.Address().Hex(), "got", addr.Hex())
	}

	logger.Trace("handlePrepareVote", "msg", msgTyp, "src", src.Address(), "hash", vote.Digest)

	if size := c.current.PrepareVoteSize(); size >= c.Q() && c.currentState() < StatePrepared {
		seals := c.getMessageSeals(size)
		newProposal, err := c.backend.PreCommit(c.current.Proposal(), seals)
		if err != nil {
			logger.Trace("Failed to assemble committed seal", "err", err)
			return err
		}

		prepareQC := proposal2QC(newProposal, c.current.Round())
		c.acceptPrepare(prepareQC, newProposal)
		logger.Trace("acceptPrepare", "msg", msgTyp, "src", src.Address(), "hash", newProposal.Hash(), "msgSize", size)

		c.sendPreCommit()
	}

	return nil
}

func (c *core) sendPreCommit() {
	logger := c.newLogger()

	msgTyp := MsgTypePreCommit
	msg := &MsgPreCommit{
		View:      c.currentView(),
		Proposal:  c.current.Proposal(),
		PrepareQC: c.current.PrepareQC(),
	}
	payload, err := Encode(msg)
	if err != nil {
		logger.Trace("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&Message{Code: msgTyp, Msg: payload})
	logger.Trace("sendPreCommit", "msg view", msg.View, "proposal", msg.Proposal.Hash())
}

func (c *core) handlePreCommit(data *Message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		msg    *MsgPreCommit
		msgTyp = MsgTypePreCommit
	)
	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to check decode", "msg", msgTyp, "err", err)
		return errFailedDecodePreCommit
	}
	if err := c.checkView(MsgTypePreCommit, msg.View); err != nil {
		logger.Trace("Failed to check view", "msg", msgTyp, "err", err)
		return err
	}
	if err := c.checkProposalView(msg.Proposal, msg.View); err != nil {
		logger.Trace("Failed to check proposal and msg view", "msg", msgTyp, "err", err)
		return err
	}
	if err := c.checkMsgFromProposer(src); err != nil {
		logger.Trace("Failed to check proposer", "msg", msgTyp, "err", err)
		return err
	}
	if msg.Proposal.Hash() != msg.PrepareQC.Hash() {
		logger.Trace("Failed to check msg", "msg", msgTyp, "expect prepareQC hash", msg.Proposal.Hash().Hex(), "got", msg.PrepareQC.hash.Hex())
		return errInvalidProposal
	}
	if _, err := c.backend.Verify(msg.Proposal); err != nil {
		logger.Trace("Failed to check verify proposal", "msg", msgTyp, "err", err)
		return err
	}
	if err := c.signer.VerifyQC(msg.PrepareQC, c.valSet); err != nil {
		logger.Trace("Failed to verify prepareQC", "msg", msgTyp, "err", err)
		return err
	}

	logger.Trace("handlePreCommit", "msg", msgTyp, "src", src.Address(), "hash", msg.Proposal.Hash())

	if c.IsProposer() && c.currentState() < StatePreCommitted {
		c.sendPreCommitVote()
	}
	if !c.IsProposer() && c.currentState() < StatePrepared {
		c.acceptPrepare(msg.PrepareQC, msg.Proposal)
		logger.Trace("acceptPrepare", "msg", msgTyp, "src", src.Address(), "prepareQC", msg.PrepareQC.Hash)

		c.sendPreCommitVote()
	}

	return nil
}

func (c *core) acceptPrepare(prepareQC *QuorumCert, proposal hotstuff.Proposal) {
	c.current.SetPrepareQC(prepareQC)
	c.current.SetProposal(proposal)
	c.current.SetState(StatePrepared)
}

func (c *core) sendPreCommitVote() {
	logger := c.newLogger()

	msgTyp := MsgTypePreCommitVote
	vote := c.current.Vote()
	if vote == nil {
		logger.Trace("Failed to send vote", "msg", msgTyp, "err", "current vote is nil")
		return
	}
	payload, err := Encode(vote)
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&Message{Code: msgTyp, Msg: payload})
	logger.Trace("sendPreCommitVote", "vote view", vote.View, "vote", vote.Digest)
}
