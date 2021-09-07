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
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) updateHighestCommittedRound(round *big.Int) {
	if c.smr.HighCommitRound().Cmp(round) < 0 {
		c.smr.SetHighCommitRound(round)
	}
}

func (c *core) handleTimeout(src hotstuff.Validator, data *hotstuff.Message) error {
	var evt *TimeoutEvent
	logger := c.newLogger()
	msgTyp := MsgTypeTimeout

	if err := data.Decode(&evt); err != nil {
		logger.Trace("Failed to decode", "msg", msgTyp, "err", err)
		return err
	}
	if err := c.checkView(evt.View); err != nil {
		logger.Trace("Failed to check view", "msg", msgTyp, "err", err)
		return err
	}
	if err := c.signer.VerifyHash(c.valset, evt.Digest, data.CommittedSeal); err != nil {
		logger.Trace("Failed to verify hash", "msg", msgTyp, "err", err)
		return err
	}

	round := evt.View.Round
	if err := c.messages.AddTimeout(round.Uint64(), data); err != nil {
		logger.Trace("Failed to add timeout", "msg", msgTyp, "err", err)
		return err
	}

	if c.isSelf(src.Address()) {
		c.increaseLastVoteRound(round)
	}

	size := c.messages.TimeoutSize(round.Uint64())
	logger.Trace("Accept Timeout", "msg", msgTyp, "from", src.Address(), "size", size)
	if size != c.Q() {
		return nil
	}

	tc := c.aggregateTC(evt, size)
	logger.Trace("Aggregate TC", "msg", msgTyp, "hash", tc.Hash)

	if err := c.advanceRoundByTC(tc, true); err != nil {
		logger.Trace("Failed to advance round by TC", "msg", msgTyp, "err", err)
		return nil
	}

	return nil
}

// advanceRoundByQC
// 使用qc或者tc驱动paceMaker进入下一轮，有个前提就是qc.round >= curRound.
// 一般而言只有leader才能收到qc，
func (c *core) advanceRoundByQC(qc *hotstuff.QuorumCert) error {
	if qc.Round().Cmp(c.smr.Round()) < 0 || qc.Height().Cmp(c.smr.Height()) < 0 {
		return errOldMessage
	}

	// catch up view
	var (
		height, round *big.Int
	)
	if bigEq(c.smr.Height(), qc.Height()) {
		height = bigAdd1(c.smr.Height())
	} else {
		height = bigAdd1(qc.Height())
	}
	if bigEq(c.smr.Round(), qc.Round()) {
		round = bigAdd1(c.smr.Round())
	} else {
		round = bigAdd1(qc.Round())
	}
	c.smr.SetHeight(height)
	c.smr.SetRound(round)

	c.valset.CalcProposerByIndex(c.smr.RoundU64())
	c.newRoundChangeTimer()
	c.logger.Trace("AdvanceQC", "view", c.currentView(), "hash", qc.Hash)

	return c.handleNewRound()
}

func (c *core) advanceRoundByTC(tc *TimeoutCert, broadcast bool) error {
	// broadcast to next leader first, we will use `curRound` again in broadcasting.
	if !c.IsProposer() && broadcast {
		c.encodeAndBroadcast(MsgTypeTC, tc)
	}

	// catch up view
	var (
		round *big.Int
	)
	if bigEq(c.smr.Round(), tc.Round()) {
		round = bigAdd1(c.smr.Round())
	} else {
		round = bigAdd1(tc.Round())
	}
	c.smr.SetRound(round)

	c.valset.CalcProposerByIndex(c.smr.RoundU64())
	c.newRoundChangeTimer()
	c.logger.Trace("AdvanceTC", "round", c.smr.Round())

	return c.handleNewRound()
}

func (c *core) stopTimer() {
	if c.timer != nil {
		c.timer.Stop()
	}
}

// todo: add in config
const standardTmoDuration = time.Second * 4

func (c *core) newRoundChangeTimer() {
	c.stopTimer()

	curRd := c.smr.RoundU64()
	hgRd := c.smr.HighCommitRoundU64()
	index := curRd - hgRd
	timeout := standardTmoDuration
	if index > 0 {
		timeout += time.Duration(math.Pow(2, float64(index))) * time.Second
	}

	c.timer = time.AfterFunc(timeout, func() {
		evt := c.generateTimeoutEvent()
		c.encodeAndBroadcast(MsgTypeTimeout, evt)
	})
}
