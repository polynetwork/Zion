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
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *core) sendRequest() error {
	logger := c.newSenderLogger("MSG_SEND_REQUEST")

	height := c.smr.Height()

	// use existing pending request
	request := c.smr.Request()
	if request != nil && bigEq(height, request.Number()) {
		logger.Trace("Got pending request", "num", request.Number(), "parent hash", request.ParentHash())
		return c.sendProposal()
	}

	// ask miner new proposal, use latest high proposal as parent block
	parent := c.smr.Proposal()
	if parent == nil {
		logger.Trace("Failed to get parent block", "num", height, "err", "parent is nil")
		return nil
	}
	if expect, eq := bigSub1Eq(height, parent.Number()); !eq {
		if !bigEq(height, parent.Number()) {
			logger.Trace("Invalid parent block", "expect height", expect, "got", parent.Number())
			return nil
		}
		parent = c.blkPool.GetBlockByHash(parent.ParentHash())
		if parent == nil {
			logger.Trace("Failed to get parent block", "err", "parent is nil")
			return nil
		}
	}

	c.feed.Send(consensus.AskRequest{
		Number: height,
		Parent: parent,
	})

	logger.Trace("Send Request", "num", height, "parent hash", parent.Hash())
	return nil
}

func (c *core) handleRequest(req *hotstuff.Request) error {
	logger := c.newSenderLogger("MSG_RECV_REQUEST")

	if req == nil || req.Proposal == nil {
		logger.Trace("Invalid request", "err", "is nil")
		return nil
	}
	proposal, ok := req.Proposal.(*types.Block)
	if !ok {
		logger.Trace("Failed to convert proposal", "err", "type invalid")
		return nil
	}
	if !bigEq(proposal.Number(), c.smr.Height()) {
		logger.Trace("Invalid proposal", "expect height", c.smr.HeightU64(), "got", proposal.Number())
		return nil
	}
	//parent := c.smr.Proposal()
	//if parent == nil {
	//	logger.Trace("Invalid parent", "err", "parent is nil or hash not equal")
	//	return nil
	//}
	//if proposal.ParentHash() != parent.Hash() {
	//	logger.Trace("Invalid parent", "expect parent hash")
	//}
	c.smr.SetRequest(proposal)
	logger.Trace("Received request", "num", req.Proposal.Number(), "hash", req.Proposal.Hash())

	return c.sendProposal()
}
