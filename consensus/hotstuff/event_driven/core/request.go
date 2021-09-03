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
	logger := e.newLogger()

	height := copyNum(e.curHeight)
	proposal := e.blkPool.GetHighProposal()
	e.feed.Send(consensus.AskRequest{
		Number: height,
		Parent: proposal,
	})

	logger.Trace("Ask Request", "num", height, "parent hash", proposal.Hash())
	return nil
}

func (e *core) handleRequest(req *hotstuff.Request) error {
	logger := e.newLogger()

	msgTyp := MsgTypeRequest

	if req == nil || req.Proposal == nil {
		logger.Trace("Receive Request", "msg", msgTyp, "request", "is nil")
		return nil
	}
	proposal, ok := req.Proposal.(*types.Block)
	if !ok {
		logger.Trace("Receive Request", "msg", msgTyp, "convert proposal", "type invalid")
		return nil
	}
	if proposal.Number().Cmp(e.curHeight) != 0 {
		logger.Trace("Receive Request", "msg", msgTyp, "expect height", e.curHeight, "got", proposal.Number())
		return nil
	}

	e.blkPool.UpdateHighProposal(proposal)
	logger.Trace("Received request", "msg", msgTyp, "num", req.Proposal.Number(), "hash", req.Proposal.Hash())

	return e.sendProposal()
}
