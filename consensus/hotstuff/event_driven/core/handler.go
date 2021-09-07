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
	"fmt"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

// Subscribe both internal and external events
func (c *core) subscribeEvents() {
	c.events = c.backend.EventMux().Subscribe(
		hotstuff.RequestEvent{},
		hotstuff.MessageEvent{},
		backlogEvent{},
	)
	c.finalCommittedSub = c.backend.EventMux().Subscribe(
		hotstuff.FinalCommittedEvent{},
	)
}

// Unsubscribe all events
func (c *core) unsubscribeEvents() {
	c.events.Unsubscribe()
	//e.timeoutSub.Unsubscribe()
	c.finalCommittedSub.Unsubscribe()
}

func (c *core) handleEvents() {
	logger := c.logger.New("handleEvents")

	for {
		select {
		case evt, ok := <-c.events.Chan():
			if !ok {
				logger.Error("Failed to receive msg Event")
				return
			}
			// A real Event arrived, process interesting content
			switch ev := evt.Data.(type) {
			case hotstuff.RequestEvent:
				c.handleRequest(&hotstuff.Request{Proposal: ev.Proposal})

			case hotstuff.MessageEvent:
				c.handleMsg(ev.Payload)

			case backlogEvent:
				c.handleCheckedMsg(ev.src, ev.msg)
			}

		case _, ok := <-c.finalCommittedSub.Chan():
			if !ok {
				logger.Error("Failed to receive finalCommitted Event")
				return
			}
			//e.handleFinalCommitted()
		}
	}
}

// sendEvent sends events to mux
func (c *core) sendEvent(ev interface{}) {
	c.backend.EventMux().Post(ev)
}

func (c *core) handleMsg(payload []byte) error {
	logger := c.logger.New()

	// Decode Message and check its signature
	msg := new(hotstuff.Message)
	if err := msg.FromPayload(payload, c.validateFn); err != nil {
		logger.Error("Failed to decode Message from payload", "err", err)
		return err
	}

	// Only accept Message if the address is valid
	_, src := c.valset.GetByAddress(msg.Address)
	if src == nil {
		logger.Error("Invalid address in Message", "msg", msg)
		return errInvalidSigner
	}

	// handle checked Message
	if err := c.handleCheckedMsg(src, msg); err != nil {
		return err
	}
	return nil
}

func (c *core) handleCheckedMsg(src hotstuff.Validator, msg *hotstuff.Message) (err error) {
	switch msg.Code {
	case MsgTypeProposal:
		err = c.handleProposal(src, msg)
	case MsgTypeVote:
		err = c.handleVote(src, msg)
	case MsgTypeTC:
		err = c.handleTC(src, msg)
	case MsgTypeTimeout:
		err = c.handleTimeout(src, msg)
	default:
		err = errInvalidMessage
	}

	if err == errFutureMessage {
		c.storeBacklog(msg, src)
	}
	return
}

func (c *core) finalizeMessage(msg *hotstuff.Message, val interface{}) ([]byte, error) {
	var err error

	// Add sender address
	msg.Address = c.address
	msg.View = c.currentView()

	if msg.Code == MsgTypeVote {
		vote, ok := val.(*Vote)
		if !ok {
			return nil, fmt.Errorf("msg is not vote")
		}
		seal, err := c.signer.SignHash(vote.Hash)
		if err != nil {
			return nil, err
		}
		msg.CommittedSeal = seal
	}

	if msg.Code == MsgTypeTimeout {
		tm, ok := val.(*TimeoutEvent)
		if !ok {
			return nil, fmt.Errorf("msg is not timeout")
		}
		tm.Digest = tm.Hash()
		seal, err := c.signer.SignHash(tm.Digest)
		if err != nil {
			return nil, err
		}
		msg.CommittedSeal = seal
	}

	// Sign Message
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}
	msg.Signature, err = c.signer.Sign(data)
	if err != nil {
		return nil, err
	}

	// Convert to payload
	payload, err := msg.Payload()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *core) encodeAndBroadcast(msgTyp MsgType, val interface{}) {
	payload, err := Encode(val)
	if err != nil {
		c.logger.Trace("Failed to encode broadcast msg payload", "msg", msgTyp, "instance", val, "err", err)
		return
	}
	if len(payload) == 0 {
		c.logger.Trace("Failed to encode broadcast msg payload", "msg", msgTyp, "instance", val, "err", "payload is nil")
		return
	}

	msg := &hotstuff.Message{
		Code: msgTyp,
		Msg:  payload,
	}
	if err := c.broadcast(msg, val); err != nil {
		c.logger.Trace("Failed to broadcast to peers", "msg", msg, "instance", val, "err", err)
	}
}

func (c *core) broadcast(msg *hotstuff.Message, val interface{}) error {
	payload, err := c.finalizeMessage(msg, val)
	if err != nil {
		return err
	}

	switch msg.Code {
	case MsgTypeProposal, MsgTypeTimeout:
		err = c.backend.Broadcast(c.valset, payload)
	case MsgTypeVote, MsgTypeTC:
		err = c.backend.Unicast(c.nextValSet(), payload)
	default:
		err = errInvalidMessage
	}
	return err
}
