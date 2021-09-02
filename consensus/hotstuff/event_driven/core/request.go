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
	"github.com/ethereum/go-ethereum/common/prque"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"time"
)

func (e *EventDrivenEngine) handleRequest(req *hotstuff.Request) error {
	logger := e.newLogger()
	//
	//if err := e.requests.checkRequest(e.currentView(), req); err == nil || err == errFutureMessage {
	//	if e.curProposal == nil {
	//		e.requests.StoreRequest(req)
	//	} else if e.curProposal.Number().Cmp(req.Proposal.Number()) < 0 {
	//		e.requests.StoreRequest(req)
	//	}
	//}
	//logger.Trace("Receive Request", "store request, number", req.Proposal.Number(), "hash", req.Proposal.Hash())
	//
	//if e.state != StateAcceptRequest {
	//	return nil
	//}
	//
	//return e.handleNewRound()
	if req == nil || req.Proposal == nil {
		logger.Trace("Receive Request", "request", "is nil")
		return nil
	}
	proposal, ok := req.Proposal.(*types.Block)
	if !ok {
		logger.Trace("Receive Request", "convert proposal", "type invalid")
		return nil
	}

	if proposal.Number().Cmp(e.curHeight) != 0 {
		logger.Trace("Receive Request", "expect height", e.curHeight, "got", proposal.Number())
	}

	e.curRequest = proposal
	return nil
}

func (e *EventDrivenEngine) getRequest() *types.Block {
	maxRetry := 20
	for maxRetry > 0 && (e.curRequest == nil || e.curRequest.Number().Cmp(e.curHeight) != 0) {
		e.askReq()
		maxRetry -= 1
		time.Sleep(500 * time.Millisecond)
	}
	return e.curRequest
}

type requestSet struct {
	pendingRequest *prque.Prque
}

func newRequestSet() *requestSet {
	return &requestSet{
		pendingRequest: prque.New(nil),
	}
}

func (s *requestSet) checkRequest(view *hotstuff.View, req *hotstuff.Request) error {
	if req == nil || req.Proposal == nil {
		return errInvalidMessage
	}

	if c := view.Height.Cmp(req.Proposal.Number()); c < 0 {
		return errFutureMessage
	} else if c > 0 {
		return errOldMessage
	} else {
		return nil
	}
}

func (s *requestSet) StoreRequest(req *hotstuff.Request) {
	priority := -req.Proposal.Number().Int64()
	s.pendingRequest.Push(req, priority)
}

func (s *requestSet) GetRequest(view *hotstuff.View) *hotstuff.Request {
	for !s.pendingRequest.Empty() {
		m, prior := s.pendingRequest.Pop()
		req, ok := m.(*hotstuff.Request)
		if !ok {
			continue
		}

		// push back if it's future message
		if err := s.checkRequest(view, req); err != nil {
			if err == errFutureMessage {
				s.pendingRequest.Push(m, prior)
				// todo: 是否为continue
				break
			}
			continue
		}
		return req
	}

	return nil
}

func (s *requestSet) Size() int {
	return s.pendingRequest.Size()
}
