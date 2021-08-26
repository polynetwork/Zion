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
	"github.com/ethereum/go-ethereum/core/types"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/prque"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (e *EventDrivenEngine) handleRequest(req *hotstuff.Request) error {
	logger := e.newLogger()

	if err := e.requests.checkRequest(e.currentView(), req); err != nil {
		if err == errFutureMessage {
			e.requests.StoreRequest(req)
			return nil
		} else {
			logger.Warn("receive request", "err", err)
			return err
		}
	} else {
		e.requests.StoreRequest(req)
	}

	logger.Trace("handleRequest", "height", req.Proposal.Number(), "proposal", req.Proposal.Hash())
	return nil
}

func (e *EventDrivenEngine) generateProposalMessage() (*MsgProposal, error) {
	req := e.requests.GetRequest(e.currentView())
	proposal, ok := req.Proposal.(*types.Block)
	if !ok {
		return nil, errProposalConvert
	}
	_, justifyQC, err := extraHeader(req.Parent)
	if err != nil {
		return nil, errExtraHeader
	}

	msg := &MsgProposal{
		Epoch:     e.epoch,
		View:      e.currentView(),
		Proposal:  proposal,
		JustifyQC: justifyQC,
	}
	return msg, nil
}

type requestSet struct {
	mtx *sync.RWMutex

	pendingRequest *prque.Prque
}

func newRequestSet() *requestSet {
	return &requestSet{
		mtx:            new(sync.RWMutex),
		pendingRequest: prque.New(nil),
	}
}

func (s *requestSet) checkRequest(view *hotstuff.View, req *hotstuff.Request) error {
	if req == nil || req.Proposal == nil {
		return errInvalidMessage
	}

	// todo(fuk): how to process future block, store or throw?
	if c := view.Height.Cmp(req.Proposal.Number()); c < 0 {
		return errFutureMessage
	} else if c > 0 {
		return errOldMessage
	} else {
		return nil
	}
}

func (s *requestSet) StoreRequest(req *hotstuff.Request) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	priority := -req.Proposal.Number().Int64()
	s.pendingRequest.Push(req, priority)
}

func (s *requestSet) GetRequest(view *hotstuff.View) *hotstuff.Request {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	maxRetry := 20
retry:
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
	if maxRetry -= 1; maxRetry > 0 {
		time.Sleep(500 * time.Millisecond)
		goto retry
	}
	return nil
}

func (s *requestSet) Size() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.pendingRequest.Size()
}
