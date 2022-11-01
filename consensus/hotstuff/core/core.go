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
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

type core struct {
	config *hotstuff.Config
	logger log.Logger

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

	validateFn func([]byte, []byte) (common.Address, error)
	isRunning  bool
}

// New creates an HotStuff consensus core
func New(backend hotstuff.Backend, config *hotstuff.Config, signer hotstuff.Signer) *core {
	c := &core{
		config:            config,
		logger:            log.New("address", backend.Address()),
		backend:           backend,
		signer:            signer,
		backlogs:          newBackLog(),
		pendingRequests:   prque.New(nil),
		pendingRequestsMu: new(sync.Mutex),
	}
	c.validateFn = c.checkValidatorSignature

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
	// * last proposal is equal to current state height -1, and the round is 0, it denote that last proposal has't chained.
	// * last proposal is equal to current state height -1, and the round lower than current state round, it should not happen.
	// * last proposal is equal to current state height -1, and the round greater or equal to current state round, it denote change view.
	if c.current == nil {
		logger.Trace("Start for the initial round")
	} else if lastProposal.NumberU64() >= c.current.HeightU64() {
		logger.Trace("Catch up latest proposal", "number", lastProposal.NumberU64(), "hash", lastProposal.Hash())
	} else if lastProposal.NumberU64() < c.current.HeightU64()-1 {
		logger.Warn("New height should be larger than current height", "new_height", lastProposal.NumberU64)
		return
	} else if round.Sign() == 0 {
		logger.Debug("Latest proposal not chained", "chained", lastProposal.NumberU64(), "current", c.current.HeightU64())
		return
	} else if round.Cmp(c.current.Round()) < 0 {
		logger.Warn("New round should not be smaller than current round", "height", lastProposal.NumberU64(), "new_round", round, "old_round", c.current.Round())
		return
	} else {
		changeView = true
	}

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

	c.valSet = c.backend.Validators(common.EmptyHash, true)
	c.valSet.CalcProposer(lastProposer, newView.Round.Uint64())
	c.updateRoundState(newView, changeView, lastProposal, c.valSet)

	logger.Debug("New round", "state", c.currentState(), "newView", newView, "new_proposer", c.valSet.GetProposer(), "valSet", c.valSet.List(), "size", c.valSet.Size(), "IsProposer", c.IsProposer())

	// set init state and process backlogs
	c.setCurrentState(StateAcceptRequest)
	// stop last timer and regenerate new timer
	c.newRoundChangeTimer()
	c.sendNewView(newView)
}

// check point and return true if the engine is stopped, return false if the validators not changed
func (c *core) checkPoint(view *View) bool {
	if c.backend.CheckPoint(view.Height.Uint64()) {
		c.backend.ReStart()
	}
	if !c.isRunning {
		return true
	}
	return false
}

func (c *core) updateRoundState(newView *View, changeView bool, lastProposal hotstuff.Proposal, valset hotstuff.ValidatorSet) {
	prepareQC := proposal2QC(lastProposal, common.Big0)

	if !changeView {
		c.current = newRoundState(newView, c.valSet, prepareQC)
		return
	}

	// reuse pending request if round changed
	lastPendingRequest := c.current.PendingRequest()
	c.current = newRoundState(newView, c.valSet, prepareQC)
	c.current.SetPendingRequest(lastPendingRequest)
}

func (c *core) setCurrentState(s State) {
	c.current.SetState(s)
	if s == StateAcceptRequest || s == StateHighQC {
		c.processPendingRequests()
	}
	c.processBacklog()
}

func (c *core) checkValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	return c.signer.CheckSignature(c.valSet, data, sig)
}
