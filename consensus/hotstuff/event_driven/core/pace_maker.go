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

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (e *core) updateHighestCommittedRound(round *big.Int) {
	if e.smr.HighCommitRound().Cmp(round) < 0 {
		e.smr.SetHighCommitRound(round)
	}
}

func (e *core) handleTimeout(src hotstuff.Validator, data *hotstuff.Message) error {
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
func (e *core) advanceRoundByQC(qc *hotstuff.QuorumCert) error {
	if qc == nil || qc.View == nil || qc.Round().Cmp(e.smr.Round()) < 0 || qc.Height().Cmp(e.smr.Height()) < 0 {
		return fmt.Errorf("qc invalid")
	}

	// catch up view
	var (
		height, round *big.Int
	)
	if bigEq(e.smr.Height(), qc.Height()) {
		height = bigAdd1(e.smr.Height())
	} else {
		height = bigAdd1(qc.Height())
	}
	if bigEq(e.smr.Round(), qc.Round()) {
		round = bigAdd1(e.smr.Round())
	} else {
		round = bigAdd1(qc.Round())
	}
	e.smr.SetHeight(height)
	e.smr.SetRound(round)

	e.valset.CalcProposerByIndex(e.smr.RoundU64())
	e.newRoundChangeTimer()
	e.logger.Trace("AdvanceQC", "view", e.currentView(), "hash", qc.Hash)

	return e.handleNewRound()
}

func (e *core) advanceRoundByTC(tc *TimeoutCert, broadcast bool) error {
	if tc == nil || tc.View == nil || tc.Round().Cmp(e.smr.Round()) < 0 || tc.Height().Cmp(e.smr.Height()) < 0 {
		return fmt.Errorf("tc invalid")
	}

	// broadcast to next leader first, we will use `curRound` again in broadcasting.
	if !e.IsProposer() && broadcast {
		e.encodeAndBroadcast(MsgTypeTC, tc)
	}

	// catch up view
	var (
		round *big.Int
	)
	if bigEq(e.smr.Round(), tc.Round()) {
		round = bigAdd1(e.smr.Round())
	} else {
		round = bigAdd1(tc.Round())
	}
	e.smr.SetRound(round)

	e.valset.CalcProposerByIndex(e.smr.RoundU64())
	e.newRoundChangeTimer()
	e.logger.Trace("AdvanceTC", "round", e.smr.Round())

	return e.handleNewRound()
}

func (e *core) stopTimer() {
	if e.timer != nil {
		e.timer.Stop()
	}
}

// todo: add in config
const standardTmoDuration = time.Second * 4

func (e *core) newRoundChangeTimer() {
	e.stopTimer()

	curRd := e.smr.RoundU64()
	hgRd := e.smr.HighCommitRoundU64()
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
