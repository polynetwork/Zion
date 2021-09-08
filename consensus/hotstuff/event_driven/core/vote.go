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

import "github.com/ethereum/go-ethereum/consensus/hotstuff"

func (c *core) sendVote() error {
	logger := c.newSenderLogger("MSG_SEND_VOTE")

	view := c.currentView()
	justifyQC := c.smr.HighQC()
	proposal := c.smr.Proposal()

	// make vote and send it to next proposer
	vote, err := c.makeVote(proposal.Hash(), proposal.Coinbase(), view, justifyQC)
	if err != nil {
		logger.Trace("[Send Vote], failed to make vote", "err", err)
		return err
	}

	logger.Trace("[Send Vote]", "to", c.nextProposer(), "vote", vote)

	c.increaseLastVoteRound(view.Round)
	c.encodeAndBroadcast(MsgTypeVote, vote)

	return nil
}

// handleVote validate vote message and try to assemble qc
func (c *core) handleVote(src hotstuff.Validator, data *hotstuff.Message) error {
	logger := c.newMsgLogger(MsgTypeVote)

	var vote *Vote
	if err := data.Decode(&vote); err != nil {
		logger.Trace("[Handle Vote], failed to decode", "from", src.Address(), "err", err)
		return errFailedDecodeNewView
	}

	logger.Trace("[Handle Vote], accept Vote", "from", src.Address(), "hash", vote.Hash, "vote view", vote.View)

	if err := c.checkVote(vote); err != nil {
		logger.Trace("[Handle Vote], failed to check vote", "from", src.Address(), "err", err)
		return err
	}
	if err := c.checkEpoch(vote.Epoch, vote.View.Height); err != nil {
		logger.Trace("[Handle Vote], failed to check epoch", "from", src.Address(), "err", err)
		return err
	}
	if err := c.validateVote(vote); err != nil {
		logger.Trace("[Handle Vote], failed to validate vote", "from", src.Address(), "err", err)
		return err
	}
	if err := c.messages.AddVote(vote.Hash, data); err != nil {
		logger.Trace("[Handle Vote], failed to add vote", "from", src.Address(), "err", err)
		return err
	}

	size := c.messages.VoteSize(vote.Hash)
	if size != c.Q() {
		return nil
	}

	highQC, sealedBlock, err := c.aggregateQC(vote, size)
	if err != nil {
		logger.Trace("[Handle Vote], failed to aggregate qc", "err", err)
		return err
	}
	logger.Trace("[Handle Vote], aggregate QC", "qc hash", highQC.Hash, "qc view", highQC.View)

	if err := c.blkPool.AddBlock(sealedBlock, vote.View.Round); err != nil {
		logger.Trace("[Handle Vote], failed to insert block into block pool", "err", err)
		return err
	}

	if err := c.updateHighQCAndProposal(highQC, sealedBlock); err != nil {
		logger.Trace("[Handle Vote], failed to update high qc and proposal", "err", err)
	}

	if err := c.advanceRoundByQC(highQC); err != nil {
		logger.Trace("[Handle Vote], failed to advance round", "err", err)
		return err
	}

	return nil
}
