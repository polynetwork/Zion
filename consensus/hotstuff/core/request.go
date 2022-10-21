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

func (c *core) handleRequest(request *Request) error {
	logger := c.newLogger()
	if err := c.checkRequestMsg(request); err != nil {
		if err == errInvalidMessage {
			logger.Warn("invalid request")
		} else if err == errFutureMessage {
			c.storeRequestMsg(request)
		} else {
			logger.Warn("unexpected request", "err", err, "number", request.Proposal.Number(), "hash", request.Proposal.Hash())
		}
		return err
	}
	logger.Trace("handleRequest", "number", request.Proposal.Number(), "hash", request.Proposal.Hash())

	switch c.currentState() {
	case StateAcceptRequest:
		// store request and prepare to use it after highQC
		c.storeRequestMsg(request)

	case StateHighQC:
		// consensus step is blocked for proposal is not ready
		if c.current.PendingRequest() == nil {
			c.current.SetPendingRequest(request)
			c.sendPrepare()
		}

	default:
		// store request for `changeView` if node is not the proposer at current round.
		if c.current.PendingRequest() == nil {
			c.current.SetPendingRequest(request)
		}
	}

	return nil
}

// check request state
// return errInvalidMessage if the message is invalid
// return errFutureMessage if the sequence of proposal is larger than current sequence
// return errOldMessage if the sequence of proposal is smaller than current sequence
func (c *core) checkRequestMsg(request *Request) error {
	if request == nil || request.Proposal == nil {
		return errInvalidMessage
	}

	if c := c.current.Height().Cmp(request.Proposal.Number()); c > 0 {
		return errOldMessage
	} else if c < 0 {
		return errFutureMessage
	} else {
		return nil
	}
}

func (c *core) storeRequestMsg(request *Request) {
	logger := c.newLogger()

	logger.Trace("Store future request", "number", request.Proposal.Number(), "hash", request.Proposal.Hash())

	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	c.pendingRequests.Push(request, -request.Proposal.Number().Int64())
}

// todo(fuk): pop too old blocks
func (c *core) processPendingRequests() {
	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	if c.pendingRequests.Empty() {
		return
	}

	for !(c.pendingRequests.Empty()) {
		m, prio := c.pendingRequests.Pop()
		r, ok := m.(*Request)
		if !ok {
			c.logger.Warn("Malformed request, skip", "msg", m)
			continue
		}
		// Push back if it's a future message
		if err := c.checkRequestMsg(r); err != nil {
			if err == errFutureMessage {
				c.logger.Trace("Stop processing request", "number", r.Proposal.Number(), "hash", r.Proposal.Hash())
				c.pendingRequests.Push(m, prio)
				break
			}
			c.logger.Trace("Skip the pending request", "number", r.Proposal.Number(), "hash", r.Proposal.Hash(), "err", err)
			continue
		} else {
			c.logger.Trace("Post pending request", "number", r.Proposal.Number(), "hash", r.Proposal.Hash())
			go c.sendEvent(hotstuff.RequestEvent{
				Proposal: r.Proposal,
			})
		}
	}
}
