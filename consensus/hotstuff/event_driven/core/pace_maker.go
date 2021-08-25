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
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

// todo input message decode
//func (e *EventDrivenEngine) handleLocalTimeout(evt *TimeoutEvent) error {
//	round := evt.Round
//	if err := e.messages.AddTimeout(round.Uint64(), data); err != nil {
//		return err
//	}
//
//	e.IncreaseLastVoteRound(round)
//
//	return e.broadcast(data)
//}

func (e *EventDrivenEngine) handleTimeout(src hotstuff.Validator, data *hotstuff.Message) error {
	var (
		evt *TimeoutEvent
	)
	if err := data.Decode(&evt); err != nil {
		return err
	}
	round := evt.Round

	if err := e.messages.AddTimeout(round.Uint64(), data); err != nil {
		return err
	}

	if e.isSelf(src.Address()) {
		e.IncreaseLastVoteRound(round)
	}

	if e.messages.TimeoutSize(round.Uint64()) == e.Q() {
		tc := &hotstuff.QuorumCert{}
		return e.advanceRound(tc)
	}

	return nil
}

func (e *EventDrivenEngine) advanceRound(qc *hotstuff.QuorumCert) error {
	qcRound := qc.View.Round
	if qcRound.Cmp(e.curRound) < 0 {
		return fmt.Errorf("qcRound < currentRound, (%v, %v)", qcRound, e.curRound)
	}

	// current round increase
	e.curRound = new(big.Int).Add(qcRound, common.Big1)

	// recalculate proposer
	e.valset.CalcProposerByIndex(e.curRound.Uint64())

	if !e.isProposer() {
		// send qc to next leader
		payload, err := Encode(qc)
		if err != nil {
			return err
		}
		if err := e.broadcast(&hotstuff.Message{
			Code: MsgTypeQC,
			Msg:  payload,
		}); err != nil {
			return err
		}
	}

	e.newRoundChangeTimer()

	return e.handleNewRound()
}

func (e *EventDrivenEngine) stopTimer() {
	if e.timer != nil {
		e.timer.Stop()
	}
}

// todo: add in config
const standardTmoDuration = time.Second * 4

func (e *EventDrivenEngine) newRoundChangeTimer() {
	e.stopTimer()

	index := e.curRound.Uint64() - 1
	if e.highestCommitRound.Uint64() != 0 {
		if e.curRound.Uint64()-e.highestCommitRound.Uint64() < 3 {
			index = 0
		} else {
			index = e.curRound.Uint64() - e.highestCommitRound.Uint64() - 3
		}
	}

	tmoDuration := time.Duration(index) * standardTmoDuration
	e.timer = time.AfterFunc(tmoDuration, func() {
		evt := &TimeoutEvent{
			Epoch: e.epoch,
			Round: e.curRound,
		}
		if payload, err := Encode(evt); err == nil {
			e.broadcast(&hotstuff.Message{
				Code: MsgTypeTimeout,
				Msg:  payload,
			})
		}
	})
}
