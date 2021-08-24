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
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}

func (e *EventDrivenEngine) checkValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	return e.signer.CheckSignature(e.valset, data, sig)
}

// todo
func (e *EventDrivenEngine) newLogger() log.Logger {
	logger := e.logger.New("state")
	return logger
}

func (e *EventDrivenEngine) address() common.Address {
	return e.addr
}

func (e *EventDrivenEngine) isProposer() bool {
	if e.valset.IsProposer(e.address()) {
		return true
	}
	return false
}

func (e *EventDrivenEngine) currentView() *hotstuff.View {
	return &hotstuff.View{
		Round:  e.paceMaker.CurrentRound(),
		Height: e.paceMaker.CurrentHeight(),
	}
}

func (e *EventDrivenEngine) finalizeMessage(msg *hotstuff.Message) ([]byte, error) {
	var err error

	// Add sender address
	msg.Address = e.address()
	msg.View = e.currentView()

	// Add proof of consensus
	// todo: sign proposal into committed seal
	//proposal := c.current.Proposal()
	//if msg.Code == MsgTypePrepareVote && proposal != nil {
	//	seal, err := c.signer.SignVote(proposal)
	//	if err != nil {
	//		return nil, err
	//	}
	//	msg.CommittedSeal = seal
	//}

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

//func (c *EventDrivenEngine) getMessageSeals(n int) [][]byte {
//	seals := make([][]byte, n)
//	for i, data := range c.current.PrepareVotes() {
//		if i < n {
//			seals[i] = data.CommittedSeal
//		}
//	}
//	return seals
//}

func (e *EventDrivenEngine) encodeAndBroadcast(msgTyp MsgType, val interface{}) error {
	payload, err := Encode(val)
	if err != nil {
		return err
	}

	msg := &hotstuff.Message{
		Code: msgTyp,
		Msg: payload,
	}

	return e.broadcast(msg)
}

func (e *EventDrivenEngine) broadcast(msg *hotstuff.Message) error {
	logger := e.newLogger()

	payload, err := e.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize Message", "msg", msg, "err", err)
		return err
	}

	switch msg.Code {
	case MsgTypeVote:
		if err := e.backend.Unicast(e.valset, payload); err != nil {
			logger.Error("Failed to unicast Message", "msg", msg, "err", err)
			return err
		}
	case MsgTypeNewView, MsgTypeTimeout:
		if err := e.backend.Broadcast(e.valset, payload); err != nil {
			logger.Error("Failed to broadcast Message", "msg", msg, "err", err)
			return err
		}
	default:
		logger.Error("invalid msg type", "msg", msg)
		return errInvalidMessage
	}
	return nil
}

// todo: extra block into justifyQC and round
func extraProposal(proposal hotstuff.Proposal) (*hotstuff.QuorumCert, *big.Int, error) {
	block := proposal.(*types.Block)
	h := block.Header()
	qc := new(hotstuff.QuorumCert)
	qc.View = &hotstuff.View{
		Height: block.Number(),
		Round:  big.NewInt(0),
	}
	qc.Hash = h.Hash()
	qc.Proposer = h.Coinbase
	qc.Extra = h.Extra
	return qc, big.NewInt(0), nil
}
