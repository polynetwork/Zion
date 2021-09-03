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

func (e *core) sendProposal() error {
	logger := e.newLogger()

	proposal := e.blkPool.GetHighProposal()
	justifyQC := e.blkPool.GetHighQC()
	view := e.currentView()
	msg := &MsgProposal{
		Epoch:     e.epoch,
		View:      view,
		Proposal:  proposal,
		JustifyQC: justifyQC,
	}

	e.encodeAndBroadcast(MsgTypeProposal, msg)
	logger.Trace("Send proposal", "hash", msg.Proposal.Hash(), "justifyQC hash", msg.JustifyQC.Hash)
	return nil
}

// handleProposal validate proposal info and vote to the next leader if the proposal is valid
func (e *core) handleProposal(src hotstuff.Validator, data *hotstuff.Message) error {
	logger := e.newLogger()

	var (
		msg    *MsgProposal
		msgTyp = MsgTypeProposal
	)
	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to decode", "type", msgTyp, "from", src.Address(), "err", err)
		return errFailedDecodePrepare
	}

	view := msg.View
	proposal := msg.Proposal
	justifyQC := msg.JustifyQC

	if err := e.checkEpoch(msg.Epoch, proposal.Number()); err != nil {
		logger.Trace("Failed to check epoch", "msg", msgTyp, "from", src.Address(), "err", err)
		return err
	}

	logger.Trace("Accept proposal", "msg", msgTyp, "proposer", src.Address(), "hash", proposal.Hash(), "height", proposal.Number())

	// try to advance into new round, it will update proposer and current view
	_ = e.processQC(justifyQC)

	// validate proposal and justify qc
	if err := e.checkView(view); err != nil {
		logger.Trace("Failed to check view", "msg", msgTyp, "from", src.Address(), "err", err)
		return err
	}
	if err := e.validateProposal(proposal); err != nil {
		logger.Trace("Failed to validate proposal", "msg", msgTyp, "from", src.Address(), "err", err)
		return err
	}
	if err := e.checkProposer(proposal.Coinbase()); err != nil {
		logger.Trace("Failed to check proposer", "msg", msgTyp, "from", src.Address(), "err", err)
		return err
	}
	if err := e.checkJustifyQC(proposal, justifyQC); err != nil {
		logger.Trace("Failed to check justify", "msg", msgTyp, "from", src.Address(), "err", err)
		return err
	}
	if err := e.signer.VerifyQC(justifyQC, e.valset); err != nil {
		logger.Trace("Failed to verify justifyQC", "msg", msgTyp, "from", src.Address(), "err", err)
		return err
	}

	// add proposal and update highQC as next justifyQC
	if err := e.blkPool.AddBlock(proposal, view.Round); err != nil {
		logger.Trace("Failed to insert block into block pool", "msg", msgTyp, "from", src.Address(), "err", err)
		return err
	}
	e.blkPool.UpdateHighProposal(proposal)
	e.blkPool.UpdateHighQC(justifyQC)

	return e.sendVote()
}
