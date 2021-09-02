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
	"math"
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
	var evt *TimeoutEvent
	logger := e.newLogger()
	msgTyp := MsgTypeTimeout

	if err := data.Decode(&evt); err != nil {
		logger.Trace("Failed to decode", "msg", msgTyp, "err", err)
		return err
	}
	if err := e.checkView(evt.View); err != nil {
		logger.Trace("Failed to check view", "msg", msgTyp, "err", err)
		return err
	}
	if err := e.signer.VerifyHash(e.valset, evt.Digest, data.CommittedSeal); err != nil {
		logger.Trace("Failed to verify hash", "msg", msgTyp, "err", err)
		return err
	}

	round := evt.View.Round
	if err := e.messages.AddTimeout(round.Uint64(), data); err != nil {
		logger.Trace("Failed to add timeout", "msg", msgTyp, "err", err)
		return err
	}

	if e.isSelf(src.Address()) {
		e.increaseLastVoteRound(round)
	}

	size := e.messages.TimeoutSize(round.Uint64())
	logger.Trace("Accept Timeout", "msg", msgTyp, "from", src.Address(), "size", size)
	if size != e.Q() {
		return nil
	}

	tc := e.aggregateTC(evt, size)
	logger.Trace("Aggregate TC", "msg", msgTyp, "hash", tc.Hash)

	if err := e.advanceRoundByTC(tc, true); err != nil {
		logger.Trace("Failed to advance round by TC", "msg", msgTyp, "err", err)
		return nil
	}

	return nil
}

// advanceRoundByQC
// 使用qc或者tc驱动paceMaker进入下一轮，有个前提就是qc.round >= curRound.
// 一般而言只有leader才能收到qc，
func (e *EventDrivenEngine) advanceRoundByQC(qc *hotstuff.QuorumCert, broadcast bool) error {
	if qc == nil || qc.View == nil || qc.View.Round.Cmp(e.curRound) < 0 {
		return fmt.Errorf("qcRound invalid")
	}

	// broadcast to next leader first, we will use `curRound` again in broadcasting.
	if !e.IsProposer() && broadcast && qc.View.Height.Cmp(e.epochHeightStart) > 0 {
		e.encodeAndBroadcast(MsgTypeQC, qc)
	}

	e.curHeight = new(big.Int).Add(e.curHeight, common.Big1)
	return e.advance(qc.View.Round, true)
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
		e.encodeAndBroadcast(MsgTypeTC, tc)
	}

	return e.advance(tcRound, false)
}

func (e *EventDrivenEngine) advance(round *big.Int, isQC bool) error {
	// current round increase
	e.curRound = new(big.Int).Add(round, common.Big1)

	// recalculate proposer
	e.valset.CalcProposerByIndex(e.curRound.Uint64())

	// reset timer
	e.newRoundChangeTimer()

	if isQC {
		e.logger.Trace("AdvanceQC", "view", e.currentView())
	} else {
		e.logger.Trace("AdvanceTC", "view", e.currentView())
	}

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
	index := curRd - hgRd
	timeout := standardTmoDuration
	if index > 0 {
		timeout += time.Duration(math.Pow(2, float64(index))) * time.Second
	}

	e.timer = time.AfterFunc(timeout, func() {
		evt := e.generateTimeoutEvent()
		e.encodeAndBroadcast(MsgTypeTimeout, evt)
	})
}
