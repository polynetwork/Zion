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

func (c *core) sendProposal() error {
	logger := c.newSenderLogger("MSG_SEND_PROPOSAL")

	// proposal and justify qc already checked in request procedure
	proposal := c.smr.Request()
	justifyQC := c.smr.HighQC()
	view := c.currentView()
	msg := &MsgProposal{
		Epoch:     c.smr.Epoch(),
		View:      view,
		Proposal:  proposal,
		JustifyQC: justifyQC,
	}

	c.encodeAndBroadcast(MsgTypeProposal, msg)
	logger.Trace("[Send proposal]", "hash", msg.Proposal.Hash(), "justifyQC hash", msg.JustifyQC.Hash)
	return nil
}

// handleProposal validate proposal info and vote to the next leader if the proposal is valid
func (c *core) handleProposal(src hotstuff.Validator, data *hotstuff.Message) error {
	logger := c.newMsgLogger(MsgTypeProposal)

	var msg *MsgProposal
	if err := data.Decode(&msg); err != nil {
		logger.Trace("[Handle Proposal], failed to decode", "from", src.Address(), "err", err)
		return errFailedDecodePrepare
	}

	view := msg.View
	unsealedBlock := msg.Proposal
	justifyQC := msg.JustifyQC

	if unsealedBlock == nil || justifyQC == nil || view == nil {
		logger.Trace("[Handle Proposal], invalid unsealedBlock msg", "err", "unsealedBlock/justifyQC/view is nil")
		return nil
	}
	if err := c.checkEpoch(msg.Epoch, unsealedBlock.Number()); err != nil {
		logger.Trace("[Handle Proposal], failed to check epoch", "from", src.Address(), "err", err)
		return err
	}
	proposalView, err := c.checkProposalView(unsealedBlock, view)
	if err != nil {
		logger.Trace("[Handle Proposal], failed to check proposal view", "from", src.Address(), "err", err)
		return err
	}
	if err := c.checkJustifyQC(unsealedBlock, justifyQC); err != nil {
		logger.Trace("[Handle Proposal], failed to check justify", "from", src.Address(), "err", err)
		return err
	}
	if err := c.signer.VerifyHeader(unsealedBlock.Header(), c.valset, false); err != nil {
		logger.Trace("[Handle Proposal], failed to validate unsealedBlock header", "from", src.Address(), "err", err)
		return err
	}
	if err := c.signer.VerifyQC(justifyQC, c.valset); err != nil {
		logger.Trace("[Handle Proposal], failed to verify justifyQC", "from", src.Address(), "err", err)
		return err
	}

	logger.Trace("[Handle Proposal], accept unsealedBlock", "proposer", src.Address(), "hash", unsealedBlock.Hash(), "proposal view", proposalView)

	// try to advance into new round, it will update proposer and current view, and reset lockQC as this justify qc.
	// unsealedBlock's great-grand parent will be committed if 3-chain can be generated.
	// logger should be reset for new view.
	_ = c.advanceRoundByQC(justifyQC)
	c.commit3Chain()
	c.updateLockQC(justifyQC)
	logger = c.newMsgLogger(MsgTypeProposal)

	// validate unsealedBlock and justify qc
	if err := c.checkView(data.Code, view); err != nil {
		logger.Trace("[Handle Proposal], failed to check view", "from", src.Address(), "err", err)
		return err
	}
	if err := c.checkProposer(unsealedBlock.Coinbase()); err != nil {
		logger.Trace("[Handle Proposal], failed to check proposer", "from", src.Address(), "err", err)
		return err
	}

	// add unsealedBlock and update highQC as next justifyQC
	if err := c.blkPool.AddBlock(unsealedBlock, view.Round); err != nil {
		logger.Trace("[Handle Proposal], failed to insert block into block pool", "from", src.Address(), "err", err)
		return err
	}
	if err := c.updateHighQCAndProposal(justifyQC, unsealedBlock); err != nil {
		logger.Trace("[Handle Proposal], failed to update high qc and proposal", "err", err)
	} else {
		logger.Trace("[Handle Proposal], update highQC", "highQC hash", justifyQC.Hash, "highQC view", justifyQC.View, "proposal hash", unsealedBlock.Hash(), "proposal view", view)
	}

	c.setCurrentState(StateProposed)

	return c.sendVote()
}
