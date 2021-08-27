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

func (e *EventDrivenEngine) updateHighestCommittedRound(round *big.Int) {
	if e.highestCommitRound.Cmp(round) < 0 {
		e.highestCommitRound = round
	}
}

func (e *EventDrivenEngine) handleTimeout(src hotstuff.Validator, data *hotstuff.Message) error {
	var (
		evt *TimeoutEvent
	)
	if err := data.Decode(&evt); err != nil {
		return err
	}
	if err := e.signer.VerifyHash(e.valset, evt.Digest, data.CommittedSeal); err != nil {
		return err
	}

	round := evt.View.Round
	if err := e.messages.AddTimeout(round.Uint64(), data); err != nil {
		return err
	}

	if e.isSelf(src.Address()) {
		e.increaseLastVoteRound(round)
	}

	if size := e.messages.TimeoutSize(round.Uint64()); size == e.Q() {
		tc := e.aggregateTC(evt, size)
		return e.advanceRoundByTC(tc, true)
	}

	return nil
}

// advanceRoundByQC
// 使用qc或者tc驱动paceMaker进入下一轮，有个前提就是qc.round >= curRound.
// 一般而言只有leader才能收到qc，
func (e *EventDrivenEngine) advanceRoundByQC(qc *hotstuff.QuorumCert, broadcast bool) error {
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
		}, qc)
	}

	e.curHeight = new(big.Int).Add(e.curHeight, common.Big1)
	return e.advance(qcRound)
}

// advanceRoundByQC
// 使用qc或者tc驱动paceMaker进入下一轮，有个前提就是qc.round >= curRound.
// 一般而言只有leader才能收到qc，
func (e *EventDrivenEngine) advanceRoundByTC(tc *TimeoutCert, broadcast bool) error {
	tcRound := tc.View.Round
	if tcRound.Cmp(e.curRound) < 0 {
		return fmt.Errorf("tcRound < currentRound, (%v, %v)", tcRound, e.curRound)
	}

	// broadcast to next leader first, we will use `curRound` again in broadcasting.
	if !e.IsProposer() && broadcast {
		payload, err := Encode(tc)
		if err != nil {
			return err
		}
		_ = e.broadcast(&hotstuff.Message{
			Code: MsgTypeTC,
			Msg:  payload,
		}, tc)
	}

	return e.advance(tcRound)
}

func (e *EventDrivenEngine) advance(round *big.Int) error {
	// current round increase
	e.curRound = new(big.Int).Add(round, common.Big1)

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

	curRd := e.curRound.Uint64()
	hgRd := e.highestCommitRound.Uint64()
	index := e.curRound.Uint64() - 1

	if hgRd != 0 {
		if curRd-hgRd < 3 {
			index = 0
		} else {
			index = curRd - hgRd - 3
		}
	}

	tmoDuration := time.Duration(index) * standardTmoDuration
	e.timer = time.AfterFunc(tmoDuration, func() {
		evt := e.generateTimeoutEvent()
		payload, err := Encode(evt)
		if err != nil {
			e.logger.Error("failed to encode timeout event, err %v", err)
			return
		}
		_ = e.broadcast(&hotstuff.Message{
			Code: MsgTypeTimeout,
			Msg:  payload,
		}, evt)
	})
}
