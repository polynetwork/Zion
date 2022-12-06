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
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/prque"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

type core struct {
	db     ethdb.Database
	logger log.Logger
	config *hotstuff.Config

	current  *roundState
	backend  hotstuff.Backend
	signer   hotstuff.Signer
	valSet   hotstuff.ValidatorSet
	backlogs *backlog

	events            *event.TypeMuxSubscription
	timeoutSub        *event.TypeMuxSubscription
	finalCommittedSub *event.TypeMuxSubscription

	roundChangeTimer *time.Timer

	pendingRequests   *prque.Prque
	pendingRequestsMu *sync.Mutex

	lastVals hotstuff.ValidatorSet // validator set for last epoch
	point    uint64                // epoch start height, header's extra contains valset

	validateFn   func(common.Hash, []byte) (common.Address, error)
	checkPointFn func(uint64) (uint64, bool)
	isRunning    bool
}

// New creates an HotStuff consensus core
func New(backend hotstuff.Backend, config *hotstuff.Config, signer hotstuff.Signer, db ethdb.Database, checkPointFn func(uint64) (uint64, bool)) *core {
	c := &core{
		db:                db,
		config:            config,
		logger:            log.New("address", backend.Address()),
		backend:           backend,
		signer:            signer,
		backlogs:          newBackLog(),
		pendingRequests:   prque.New(nil),
		pendingRequestsMu: new(sync.Mutex),
	}
	c.validateFn = c.checkValidatorSignature
	c.checkPointFn = checkPointFn
	return c
}

func (c *core) startNewRound(round *big.Int) {
	logger := c.logger.New()

	if !c.isRunning {
		logger.Trace("Start engine first")
		return
	}

	var (
		changeView                 = false
		lastProposal, lastProposer = c.backend.LastProposal()
	)

	// check last chained block
	if lastProposal == nil {
		logger.Warn("Last proposal should not be nil")
		return
	}

	// compare the chained block height and current state height, there are 6 conditions:
	// * current state is nil, it denote that the engine is initialed just now.
	// * last proposal is greater than current state height, it denote that last proposal has chained.
	// * last proposal is lower that current state height - 1, it should not happen.
	// * last proposal is equal to current state height -1, and the round is 0, it denote that last proposal consensus finished but not chained.
	// * last proposal is equal to current state height -1, and the round lower than current state round, it should not happen.
	// * last proposal is equal to current state height -1, and the round greater or equal to current state round, it denote change view.
	if c.current == nil {
		logger.Trace("Start for the initial round")
	} else if lastProposal.NumberU64() >= c.HeightU64() {
		logger.Trace("Catch up latest proposal", "number", lastProposal.NumberU64(), "hash", lastProposal.Hash())
	} else if lastProposal.NumberU64() < c.HeightU64()-1 {
		logger.Warn("New height should be larger than current height", "new_height", lastProposal.NumberU64)
		return
	} else if round.Sign() == 0 {
		logger.Trace("Latest proposal not chained", "chained", lastProposal.NumberU64(), "current", c.HeightU64())
		return
	} else if round.Uint64() < c.RoundU64() {
		logger.Warn("New round should not be smaller than current round", "height", lastProposal.NumberU64(), "new_round", round, "old_round", c.RoundU64())
		return
	} else {
		changeView = true
	}

	// set view and check point before round start
	newView := &View{
		Height: new(big.Int).Add(lastProposal.Number(), common.Big1),
		Round:  new(big.Int),
	}
	if changeView {
		newView.Height = new(big.Int).Set(c.current.Height())
		newView.Round = new(big.Int).Set(round)
	} else if c.checkPoint(newView) {
		logger.Trace("Stop engine after check point.")
		return
	}

	// calculate validator set
	c.valSet = c.backend.Validators(newView.HeightU64(), true)
	c.valSet.CalcProposer(lastProposer, newView.Round.Uint64())

	// update smr and try to unlock at the round0
	if err := c.updateRoundState(lastProposal, newView); err != nil {
		logger.Error("Update round state failed", "state", c.currentState(), "newView", newView, "err", err)
		return
	}
	if !changeView {
		if err := c.current.Unlock(); err != nil {
			logger.Error("Unlock node failed", "newView", newView, "err", err)
			return
		}
	}

	logger.Debug("New round", "state", c.currentState(), "newView", newView, "new_proposer", c.valSet.GetProposer(), "valSet", c.valSet.List(), "size", c.valSet.Size(), "IsProposer", c.IsProposer())

	// stop last timer and regenerate new timer
	c.newRoundChangeTimer()

	// reset status and process pending request & backlogs
	c.setCurrentState(StateAcceptRequest)
	// start new round from message of `newView`
	c.sendNewView()
}

// check point and return true if the engine is stopped, return false if the validators not changed
func (c *core) checkPoint(view *View) bool {
	if c.checkPointFn == nil {
		return false
	}

	if epochStart, ok := c.checkPointFn(view.HeightU64()); ok {
		c.point = epochStart
		c.lastVals = c.valSet.Copy()
		c.logger.Trace("CheckPoint done", "view", view, "point", c.point)
		c.backend.ReStart()
	}
	if !c.isRunning {
		return true
	}
	return false
}

func (c *core) updateRoundState(lastProposal *types.Block, newView *View) error {
	if c.current == nil {
		c.current = newRoundState(c.db, c.logger.New(), c.valSet, lastProposal, newView)
		c.current.reload(newView)
	} else {
		c.current = c.current.update(c.valSet, lastProposal, newView)
	}

	if !c.isEpochStartQC(c.currentView(), nil) {
		return nil
	}

	prepareQC := c.current.PrepareQC()
	if prepareQC != nil && prepareQC.node == lastProposal.Hash() {
		c.logger.Trace("EpochStartPrepareQC already exist!", "newView", newView, "last block height", lastProposal.NumberU64(), "last block hash", lastProposal.Hash(), "qc.node", prepareQC.node, "qc.view", prepareQC.view, "qc.proposer", prepareQC.proposer)
		return nil
	}

	qc, err := epochStartQC(lastProposal)
	if err != nil {
		return err
	}
	if err := c.current.SetPrepareQC(qc); err != nil {
		return err
	}
	// clear old `lockQC` and `commitQC`
	c.current.lockQC = nil
	c.current.committedQC = nil
	c.logger.Trace("EpochStartPrepareQC settled!", "newView", newView, "last block height", lastProposal.NumberU64(), "last block hash", lastProposal.Hash(), "qc.node", qc.node, "qc.view", qc.view, "qc.proposer", qc.proposer)
	return nil
}

// setCurrentState handle backlog message after round state settled.
func (c *core) setCurrentState(s State) {
	c.current.SetState(s)
	if s == StateAcceptRequest || s == StateHighQC {
		c.processPendingRequests()
	}
	c.processBacklog()
}

func (c *core) checkValidatorSignature(hash common.Hash, sig []byte) (common.Address, error) {
	return c.signer.CheckSignature(c.valSet, hash, sig)
}
