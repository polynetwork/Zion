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
	"github.com/ethereum/go-ethereum/event"
)

type EventSender struct {
	eventMtx *event.TypeMux
}

func NewEventSender(backend hotstuff.Backend) *EventSender {
	return &EventSender{eventMtx: backend.EventMux()}
}

func (s *EventSender) sendEvent(val interface{}) {
	s.eventMtx.Post(val)
}

// ----------------------------------------------------------------------------

// Subscribe both internal and external events
func (e *EventDrivenEngine) subscribeEvents() {
	e.events = e.backend.EventMux().Subscribe(
		// external events
		hotstuff.RequestEvent{},
		hotstuff.MessageEvent{},
		backlogEvent{},
	)
	//e.timeoutSub = e.backend.EventMux().Subscribe(
	//	TimeoutEvent{},
	//)
	e.finalCommittedSub = e.backend.EventMux().Subscribe(
		hotstuff.FinalCommittedEvent{},
	)
}

// Unsubscribe all events
func (e *EventDrivenEngine) unsubscribeEvents() {
	e.events.Unsubscribe()
	//e.timeoutSub.Unsubscribe()
	e.finalCommittedSub.Unsubscribe()
}

func (e *EventDrivenEngine) handleEvents() {
	logger := e.logger.New("handleEvents")

	for {
		select {
		case evt, ok := <-e.events.Chan():
			if !ok {
				logger.Error("Failed to receive msg Event")
				return
			}
			// A real Event arrived, process interesting content
			switch ev := evt.Data.(type) {
			case hotstuff.RequestEvent:
				e.handleRequest(&hotstuff.Request{Proposal: ev.Proposal, Parent: ev.Parent})

			case hotstuff.MessageEvent:
				e.handleMsg(ev.Payload)

			case backlogEvent:
				e.handleCheckedMsg(ev.src, ev.msg)
			}

		case _, ok := <-e.finalCommittedSub.Chan():
			if !ok {
				logger.Error("Failed to receive finalCommitted Event")
				return
			}
			//e.handleFinalCommitted()
		}
	}
}

// sendEvent sends events to mux
func (e *EventDrivenEngine) sendEvent(ev interface{}) {
	e.backend.EventMux().Post(ev)
}

func (e *EventDrivenEngine) handleMsg(payload []byte) error {
	logger := e.logger.New()

	// Decode Message and check its signature
	msg := new(hotstuff.Message)
	if err := msg.FromPayload(payload, e.validateFn); err != nil {
		logger.Error("Failed to decode Message from payload", "err", err)
		return err
	}

	// Only accept Message if the address is valid
	_, src := e.valset.GetByAddress(msg.Address)
	if src == nil {
		logger.Error("Invalid address in Message", "msg", msg)
		return errInvalidSigner
	}

	// handle checked Message
	if err := e.handleCheckedMsg(src, msg); err != nil {
		return err
	}
	return nil
}

func (e *EventDrivenEngine) handleCheckedMsg(src hotstuff.Validator, msg *hotstuff.Message) (err error) {
	switch msg.Code {
	case MsgTypeProposal:
		err = e.handleProposal(src, msg)
	case MsgTypeVote:
		err = e.handleVote(src, msg)
	case MsgTypeQC:
		err = e.handleQuorumCertificate(src, msg)
	case MsgTypeTC:
		err = e.handleTimeoutCertificate(src, msg)
	case MsgTypeTimeout:
		err = e.handleTimeout(src, msg)
	default:
		err = errInvalidMessage
		e.logger.Error("msg type invalid", "unknown type", msg.Code)
	}

	if err == errFutureMessage {
		//e.storeBacklog(msg, src)
	}
	return
}

func (e *EventDrivenEngine) finalizeMessage(msg *hotstuff.Message, val interface{}) ([]byte, error) {
	var err error

	// Add sender address
	msg.Address = e.address()
	msg.View = e.currentView()

	if msg.Code == MsgTypeVote {
		vote, ok := val.(*Vote)
		if !ok {
			return nil, fmt.Errorf("msg is not vote")
		}
		seal, err := e.signer.SignHash(vote.Hash)
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
		digest := tm.Hash()
		tm.Digest = digest
		seal, err := e.signer.SignHash(tm.Digest)
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
	msg.Signature, err = e.signer.Sign(data)
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

func (e *EventDrivenEngine) encodeAndBroadcast(msgTyp MsgType, val interface{}) error {
	payload, err := Encode(val)
	if err != nil {
		return err
	}

	msg := &hotstuff.Message{
		Code: msgTyp,
		Msg:  payload,
	}

	return e.broadcast(msg, val)
}

func (e *EventDrivenEngine) broadcast(msg *hotstuff.Message, val interface{}) error {
	logger := e.newLogger()

	payload, err := e.finalizeMessage(msg, val)
	if err != nil {
		logger.Error("Failed to finalize Message", "msg", msg, "err", err)
		return err
	}

	switch msg.Code {
	case MsgTypeProposal, MsgTypeTimeout:
		if err := e.backend.Broadcast(e.valset, payload); err != nil {
			logger.Error("Failed to broadcast Message", "msg", msg, "err", err)
			return err
		}
	case MsgTypeVote, MsgTypeQC, MsgTypeTC:
		// vote to next round leader
		nextRoundVals := e.valset.Copy()
		nextRoundVals.CalcProposerByIndex(e.curRound.Uint64() + 1)
		if err := e.backend.Unicast(nextRoundVals, payload); err != nil {
			logger.Error("Failed to unicast Message", "msg", msg, "err", err)
			return err
		}
	default:
		logger.Error("invalid msg type", "msg", msg)
		return errInvalidMessage
	}
	return nil
}
