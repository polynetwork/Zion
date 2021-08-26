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
		e.increaseLastVoteRound(round)
	}

	if e.messages.TimeoutSize(round.Uint64()) == e.Q() {
		tc := &hotstuff.QuorumCert{}
		return e.advanceRound(tc, true)
	}

	return nil
}

// advanceRound
// 使用qc或者tc驱动paceMaker进入下一轮，有个前提就是qc.round >= curRound.
// 一般而言只有leader才能收到qc，
func (e *EventDrivenEngine) advanceRound(qc *hotstuff.QuorumCert, broadcast bool) error {
	qcRound := qc.View.Round
	if qcRound.Cmp(e.curRound) < 0 {
		return fmt.Errorf("qcRound < currentRound, (%v, %v)", qcRound, e.curRound)
	}

	// broadcast to next leader first, we will use `curRound` again in broadcasting.
	if !e.IsProposer() && broadcast {
		payload, err := Encode(qc)
		if err != nil {
			return err
		}
		_ = e.broadcast(&hotstuff.Message{
			Code: MsgTypeQC,
			Msg:  payload,
		})
	}

	// current round increase
	e.curRound = new(big.Int).Add(qcRound, common.Big1)

	// recalculate proposer
	e.valset.CalcProposerByIndex(e.curRound.Uint64())

	// reset timer
	e.newRoundChangeTimer()

	// get into new consensus round
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
