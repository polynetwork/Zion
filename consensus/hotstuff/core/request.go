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
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *core) sendRequest() {
	logger := c.newLogger()

	if c.hasValidPendingRequest() {
		c.sendPrepare()
		return
	}

	proposal, _ := c.backend.LastProposal()
	if proposal == nil {
		logger.Trace("sendRequest", "err", "last proposal is nil")
		return
	}
	parent, ok := proposal.(*types.Block)
	if !ok {
		logger.Trace("sendRequest", "err", "convert proposal to block failed")
		return
	}

	c.backend.AskMiningProposalWithParent(parent)
}

func (c *core) handleRequest(req *hotstuff.Request) error {
	logger := c.newLogger()

	if c.currentState() != StateHighQC {
		logger.Trace("handleRequest state invalid")
		return nil
	}
	if c.current.highQC == nil {
		logger.Trace("handleRequest current highQC is nil")
		return nil
	}

	if !c.hasValidPendingRequest() {
		c.current.SetPendingRequest(req)
	} else {
		logger.Trace("handleRequest", "err", "already has valid pending request")
		return errRequestAlreadyExist
	}

	c.sendPrepare()
	logger.Trace("handleRequest", "height", req.Proposal.Number(), "proposal", req.Proposal.Hash())
	return nil
}

func (c *core) createNewProposal() error {
	if c.current.IsProposalLocked() {
		if c.current.Proposal() == nil {
			return errLockProposalNotExist
		} else {
			return nil
		}
	}

	if cur := c.current.Proposal(); cur != nil && cur.Number().Cmp(c.current.Height()) == 0 {
		return nil
	}

	req := c.current.PendingRequest()
	if req == nil || req.Proposal == nil {
		return errNoRequest
	}

	pending := req.Proposal
	if pending == nil || pending.Number().Cmp(c.current.Height()) != 0 {
		return errInvalidProposal
	}

	c.current.SetProposal(pending)
	return nil
}
