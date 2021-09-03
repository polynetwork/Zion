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

func (e *core) sendRequest() error {
	logger := e.newLogger("msg", MsgTypeSendRequest)

	height := e.smr.Height()

	// use existing pending request
	request := e.smr.Request()
	if request != nil && bigEq(height, request.Number()) {
		logger.Trace("Got pending request", "num", request.Number(), "parent hash", request.ParentHash())
		return e.sendProposal()
	}

	// ask miner new proposal, use latest high proposal as parent block
	parent := e.smr.Proposal()
	if parent == nil {
		logger.Trace("Failed to get parent block", "num", height, "err", "parent is nil")
		return nil
	}
	if expect, eq := bigSub1Eq(height, parent.Number()); !eq {
		logger.Trace("Invalid parent block", "expect height", expect, "got", parent.Number())
		return nil
	}

	e.feed.Send(consensus.AskRequest{
		Number: height,
		Parent: parent,
	})

	logger.Trace("Send Request", "num", height, "parent hash", parent.Hash())
	return nil
}

func (e *core) handleRequest(req *hotstuff.Request) error {
	logger := e.newLogger("msg", MsgTypeRequest)

	if req == nil || req.Proposal == nil {
		logger.Trace("Invalid request", "err", "is nil")
		return nil
	}
	proposal, ok := req.Proposal.(*types.Block)
	if !ok {
		logger.Trace("Failed to convert proposal", "err", "type invalid")
		return nil
	}
	if !bigEq(proposal.Number(), e.smr.Height()) {
		logger.Trace("Invalid proposal", "expect height", e.smr.HeightU64(), "got", proposal.Number())
		return nil
	}
	if parent := e.smr.Proposal(); parent == nil || proposal.ParentHash() != parent.Hash() {
		logger.Trace("Invalid parent", "err", "parent is nil or hash not equal")
		return nil
	}

	e.smr.SetRequest(proposal)
	logger.Trace("Received request", "num", req.Proposal.Number(), "hash", req.Proposal.Hash())

	return e.sendProposal()
}
