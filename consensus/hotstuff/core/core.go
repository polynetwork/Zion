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
	"bytes"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

type core struct {
	config *hotstuff.Config
	logger log.Logger

	current *roundState
	backend hotstuff.Backend
	signer  hotstuff.Signer

	lastEpochValSet     hotstuff.ValidatorSet
	curEpochStartHeight uint64
	valSet              hotstuff.ValidatorSet
	requests            *requestSet
	backlogs            *backlog

	events            *event.TypeMuxSubscription
	timeoutSub        *event.TypeMuxSubscription
	finalCommittedSub *event.TypeMuxSubscription

	roundChangeTimer *time.Timer

	validateFn func([]byte, []byte) (common.Address, error)
	isRunning  bool
}

// New creates an HotStuff consensus core
func New(backend hotstuff.Backend, config *hotstuff.Config, signer hotstuff.Signer) hotstuff.CoreEngine {
	c := &core{
		config:  config,
		logger:  log.New("address", backend.Address()),
		backend: backend,
	}
	c.requests = newRequestSet()
	c.backlogs = newBackLog()
	c.validateFn = c.checkValidatorSignature
	c.signer = signer

	// todo(fuk): delete after test
	rand.Seed(time.Now().UnixNano())
	return c
}

func (c *core) Address() common.Address {
	return c.signer.Address()
}

func (c *core) IsProposer() bool {
	return c.valSet.IsProposer(c.backend.Address())
}

func (c *core) IsCurrentProposal(blockHash common.Hash) bool {
	if c.current == nil {
		return false
	}
	if proposal := c.current.Proposal(); proposal != nil && proposal.Hash() == blockHash {
		return true
	}
	if req := c.current.PendingRequest(); req != nil && req.Proposal != nil && req.Proposal.Hash() == blockHash {
		return true
	}
	return false
}

func (c *core) PrepareExtra(header *types.Header, valSet hotstuff.ValidatorSet) ([]byte, error) {
	var (
		buf  bytes.Buffer
		vals = valSet.AddressList()
	)

	// compensate the lack bytes if header.Extra is not enough IstanbulExtraVanity bytes.
	if len(header.Extra) < types.HotstuffExtraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, types.HotstuffExtraVanity-len(header.Extra))...)
	}
	buf.Write(header.Extra[:types.HotstuffExtraVanity])

	ist := &types.HotstuffExtra{
		Validators:    vals,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return nil, err
	}

	return append(buf.Bytes(), payload...), nil
}

func (c *core) GetHeader(hash common.Hash, number uint64) *types.Header {
	return nil
}

func (c *core) SubscribeRequest(ch chan<- consensus.AskRequest) event.Subscription {
	return nil
}

const maxRetry uint64 = 10

func (c *core) startNewRound(round *big.Int) {
	logger := c.logger.New()

	if !c.isRunning {
		logger.Trace("Start engine first")
		return
	}

	changeView := false
	catchUpRetryCnt := maxRetry
	retryPeriod := time.Duration(c.config.RequestTimeout/maxRetry) * time.Millisecond

catchup:
	lastProposal, lastProposer := c.backend.LastProposal()
	if c.current == nil {
		logger.Trace("Start to the initial round")
	} else if lastProposal == nil {
		logger.Warn("Last proposal should not be nil")
		return
	} else if lastProposal.Number().Cmp(c.current.Height()) >= 0 {
		logger.Trace("Catch up latest proposal", "number", lastProposal.Number().Uint64(), "hash", lastProposal.Hash())
	} else if lastProposal.Number().Cmp(big.NewInt(c.current.Height().Int64()-1)) == 0 {
		if round.Cmp(common.Big0) == 0 {
			// chain reader sync last proposal
			if catchUpRetryCnt -= 1; catchUpRetryCnt <= 0 {
				logger.Warn("Sync last proposal failed", "height", c.current.Height())
				return
			} else {
				time.Sleep(retryPeriod)
				goto catchup
			}
		} else if round.Cmp(c.current.Round()) < 0 {
			logger.Warn("New round should not be smaller than current round", "height", lastProposal.Number().Int64(), "new_round", round, "old_round", c.current.Round())
			return
		}
		changeView = true
	} else {
		logger.Warn("New height should be larger than current height", "new_height", lastProposal.Number().Int64())
		return
	}

	newView := &hotstuff.View{
		Height: new(big.Int).Add(lastProposal.Number(), common.Big1),
		Round:  common.Big0,
	}
	if changeView {
		newView.Height = new(big.Int).Set(c.current.Height())
		newView.Round = new(big.Int).Set(round)
	}

	var (
		lastProposalLocked bool
		lastLockedProposal hotstuff.Proposal
		lastPendingRequest *hotstuff.Request
	)
	if c.current != nil {
		lastProposalLocked, lastLockedProposal = c.current.LastLockedProposal()
		lastPendingRequest = c.current.PendingRequest()
	}

	// calculate new proposal and init round state
	c.valSet.CalcProposer(lastProposer, newView.Round.Uint64())
	prepareQC := proposal2QC(lastProposal, common.Big0)
	c.current = newRoundState(newView, c.valSet, prepareQC)
	if changeView && lastProposalLocked && lastLockedProposal != nil {
		c.current.SetProposal(lastLockedProposal)
		c.current.LockProposal()
	}
	if changeView && lastPendingRequest != nil {
		c.current.SetPendingRequest(lastPendingRequest)
	}

	logger.Debug("New round", "state", c.currentState(), "newView", newView, "new_proposer", c.valSet.GetProposer(), "valSet", c.valSet.List(), "size", c.valSet.Size(), "IsProposer", c.IsProposer())

	// process pending request
	c.setCurrentState(StateAcceptRequest)
	c.sendNewView(newView)

	// stop last timer and regenerate new timer
	c.newRoundChangeTimer()
}

func (c *core) currentView() *hotstuff.View {
	return &hotstuff.View{
		Height: new(big.Int).Set(c.current.Height()),
		Round:  new(big.Int).Set(c.current.Round()),
	}
}

func (c *core) currentState() State {
	return c.current.State()
}

func (c *core) setCurrentState(s State) {
	c.current.SetState(s)
	c.processBacklog()
}

func (c *core) currentProposer() hotstuff.Validator {
	return c.valSet.GetProposer()
}

func (c *core) Q() int {
	return c.valSet.Q()
}

func (c *core) stopTimer() {
	if c.roundChangeTimer != nil {
		c.roundChangeTimer.Stop()
	}
}

func (c *core) newRoundChangeTimer() {
	c.stopTimer()

	// set timeout based on the round number
	timeout := time.Duration(c.config.RequestTimeout) * time.Millisecond
	round := c.current.Round().Uint64()
	if round > 0 {
		timeout += time.Duration(math.Pow(2, float64(round))) * time.Second
	}
	c.roundChangeTimer = time.AfterFunc(timeout, func() {
		c.sendEvent(timeoutEvent{})
	})
}
