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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

// Start implements core.Engine.Start
func (c *core) Start(chain consensus.ChainReader) {
	c.isRunning = true
	c.current = nil

	c.subscribeEvents()
	go c.handleEvents()

	// Start a new round from last sequence + 1
	c.startNewRound(common.Big0)
}

// Stop implements core.Engine.Stop
func (c *core) Stop() {
	c.stopTimer()
	c.unsubscribeEvents()
	c.isRunning = false
}

// Address implement core.Engine.Address
func (c *core) Address() common.Address {
	return c.signer.Address()
}

// IsProposer implement core.Engine.IsProposer
func (c *core) IsProposer() bool {
	return c.valSet.IsProposer(c.backend.Address())
}

func (c *core) IsCurrentProposal(blockHash common.Hash) bool {
	if c.current == nil {
		return false
	}
	if node := c.current.Node(); node != nil && node.Block != nil && node.Block.Hash() == blockHash {
		return true
	}
	if req := c.current.PendingRequest(); req != nil && req.block != nil && req.block.Hash() == blockHash {
		return true
	}
	return false
}

// ----------------------------------------------------------------------------

// Subscribe both internal and external events
func (c *core) subscribeEvents() {
	c.events = c.backend.EventMux().Subscribe(
		// external events
		hotstuff.RequestEvent{},
		// internal events
		hotstuff.MessageEvent{},
		backlogEvent{},
	)
	c.timeoutSub = c.backend.EventMux().Subscribe(
		timeoutEvent{},
	)
	c.finalCommittedSub = c.backend.EventMux().Subscribe(
		hotstuff.FinalCommittedEvent{},
	)
}

// Unsubscribe all events
func (c *core) unsubscribeEvents() {
	c.events.Unsubscribe()
	c.timeoutSub.Unsubscribe()
	c.finalCommittedSub.Unsubscribe()
}

func (c *core) handleEvents() {
	logger := c.logger.New("handleEvents")

	for {
		select {
		case event, ok := <-c.events.Chan():
			if !ok {
				logger.Error("Failed to receive msg Event", "err", "subscribe event chan out empty")
				return
			}
			// A real Event arrived, process interesting content
			switch ev := event.Data.(type) {
			case hotstuff.RequestEvent:
				c.handleRequest(&Request{block: ev.Block})

			case hotstuff.MessageEvent:
				c.handleMsg(ev.Src, ev.Payload)

			case backlogEvent:
				c.handleCheckedMsg(ev.msg)
			}

		case _, ok := <-c.timeoutSub.Chan():
			if !ok {
				logger.Error("Failed to receive timeout Event")
				return
			}
			c.handleTimeoutMsg()

		case evt, ok := <-c.finalCommittedSub.Chan():
			if !ok {
				logger.Error("Failed to receive finalCommitted Event")
				return
			}
			switch ev := evt.Data.(type) {
			case hotstuff.FinalCommittedEvent:
				c.handleFinalCommitted(ev.Header)
			}
		}
	}
}

// sendEvent sends events to mux
func (c *core) sendEvent(ev interface{}) {
	c.backend.EventMux().Post(ev)
}

func (c *core) handleMsg(val common.Address, payload []byte) error {
	logger := c.logger.New()

	// Decode Message and check its signature
	msg := new(Message)
	if err := msg.FromPayload(val, payload, c.validateFn); err != nil {
		logger.Error("Failed to decode Message from payload", "err", err)
		return errFailedDecodeMessage
	}

	// Only accept message if the src is consensus participant
	index, src := c.valSet.GetByAddress(val)
	if index < 0 || src == nil {
		logger.Error("Invalid address in Message", "msg", msg)
		return errInvalidSigner
	}

	// handle checked Message
	return c.handleCheckedMsg(msg)
}

func (c *core) handleCheckedMsg(msg *Message) (err error) {
	if c.current == nil {
		c.logger.Error("engine state not prepared...")
		return
	}

	switch msg.Code {
	case MsgTypeNewView:
		err = c.handleNewView(msg)
	case MsgTypePrepare:
		err = c.handlePrepare(msg)
	case MsgTypePrepareVote:
		err = c.handlePrepareVote(msg)
	case MsgTypePreCommit:
		err = c.handlePreCommit(msg)
	case MsgTypePreCommitVote:
		err = c.handlePreCommitVote(msg)
	case MsgTypeCommit:
		err = c.handleCommit(msg)
	case MsgTypeCommitVote:
		err = c.handleCommitVote(msg)
	case MsgTypeDecide:
		err = c.handleDecide(msg)
	default:
		err = errInvalidMessage
		c.logger.Error("msg type invalid", "unknown type", msg.Code)
	}

	if err == errFutureMessage {
		c.storeBacklog(msg)
	}
	return
}

func (c *core) handleTimeoutMsg() {
	c.logger.Trace("handleTimeout", "state", c.currentState(), "view", c.currentView())
	round := new(big.Int).Add(c.current.Round(), common.Big1)
	c.startNewRound(round)
}

func (c *core) broadcast(code MsgType, payload []byte) {
	logger := c.logger.New("state", c.currentState())

	// forbid unConsensus nodes send message to leader
	if index, _ := c.valSet.GetByAddress(c.Address()); index < 0 {
		return
	}

	msg := NewCleanMessage(c.currentView(), code, payload)
	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize Message", "msg", msg, "err", err)
		return
	}

	// leader set source for message and add it in set directly if msg.code is kind of vote.
	// the self voting happened before qc assembling to ensure the field of qc.seal WONT miss.
	switch msg.Code {
	case MsgTypeNewView, MsgTypePrepareVote, MsgTypePreCommitVote, MsgTypeCommitVote:
		if err = c.backend.Unicast(c.valSet, payload); err != nil {
			logger.Error("Failed to unicast Message", "msg", msg, "err", err)
		}

	case MsgTypePrepare, MsgTypePreCommit, MsgTypeCommit, MsgTypeDecide:
		if err = c.backend.Broadcast(c.valSet, payload); err != nil {
			logger.Error("Failed to broadcast Message", "msg", msg, "err", err)
		}
	default:
		logger.Error("invalid msg type", "msg", msg)
	}
}

func (c *core) finalizeMessage(msg *Message) ([]byte, error) {
	var (
		seal, sig []byte
		err       error
	)

	// Add proof of consensus
	node := c.current.Node()
	if msg.Code == MsgTypeCommitVote && node != nil && node.Block != nil {
		if seal, err = c.signer.SignHash(node.Block.Hash(), true); err != nil {
			return nil, err
		}
		msg.CommittedSeal = seal
	}

	// Sign Message
	if _, err = msg.PayloadNoSig(); err != nil {
		return nil, err
	}
	if sig, err = c.signer.SignHash(msg.hash, false); err != nil {
		return nil, err
	} else {
		msg.Signature = sig
	}

	// Convert to payload
	return msg.Payload()
}
