package core

import (
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

type core struct {
	config *hotstuff.Config
	//address common.Address
	logger log.Logger

	current *roundState
	backend hotstuff.Backend
	signer  hotstuff.Signer

	valSet   hotstuff.ValidatorSet
	requests *requestSet

	events            *event.TypeMuxSubscription
	timeoutSub        *event.TypeMuxSubscription
	finalCommittedSub *event.TypeMuxSubscription

	roundChangeTimer *time.Timer

	validateFn func([]byte, []byte) (common.Address, error) // == c.checkValidatorSignature
}

// New creates an HotStuff consensus core
func New(backend hotstuff.Backend, config *hotstuff.Config, valSet hotstuff.ValidatorSet) CoreEngine {
	c := &core{
		config:  config,
		logger:  log.New("address", backend.Address()),
		backend: backend,
	}
	c.requests = newRequestSet()
	c.validateFn = c.checkValidatorSignature
	c.valSet = valSet
	return c
}

func (c *core) Address() common.Address {
	return c.signer.Address()
}

func (c *core) IsProposer() bool {
	return c.valSet.IsProposer(c.backend.Address())
}

func (c *core) IsCurrentProposal(blockHash common.Hash) bool {
	if c.current != nil && c.current.Proposal() != nil && c.current.Proposal().Hash() == blockHash {
		return true
	}
	if c.current != nil && c.current.PendingRequest() != nil && c.current.PendingRequest().Proposal.Hash() == blockHash {
		return true
	}
	return false
}

func (c *core) startNewRound(round *big.Int) {
	var logger log.Logger
	if c.current == nil {
		logger = c.logger.New("old_round", -1, "old_height", 0)
	} else {
		logger = c.logger.New("old_round", c.current.Round(), "old_height", c.current.Height())
	}

	changeView := false
	// Try to get last proposal
	lastProposal, lastProposer := c.backend.LastProposal()
	if c.current == nil {
		logger.Trace("Start to the initial round")
	} else if lastProposal.Number().Cmp(big.NewInt(c.current.Height().Int64()-1)) == 0 {
		if round.Cmp(common.Big0) == 0 {
			// same height and round, don't need to start new round
			return
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

	c.valSet.CalcProposer(lastProposer, newView.Round.Uint64())
	if c.current == nil {
		prepareQC := Proposal2QC(lastProposal)
		c.current = newRoundState(newView, c.valSet, prepareQC)
	} else {
		c.current = c.current.Spawn(newView, c.valSet)
	}
	if c.current.PendingRequest() != nil {
		c.sendNewView(newView)
	}
	c.newRoundChangeTimer()

	logger.Debug("New round", "last block number", lastProposal.Number().Uint64(), "new_round", newView.Round, "new_height", newView.Height, "new_proposer", c.valSet.GetProposer(), "valSet", c.valSet.List(), "size", c.valSet.Size(), "IsProposer", c.IsProposer())
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

func (c *core) currentProposer() hotstuff.Validator {
	return c.valSet.GetProposer()
}

// todo: 检查是否收到了新的proposal
func (c *core) catchUpRound(view *hotstuff.View) {
	logger := c.logger.New("old_round", c.current.Round(), "old_height", c.current.Height(), "old_proposer", c.valSet.GetProposer())

	c.newRoundChangeTimer()

	logger.Trace("Catch up round", "new_round", view.Round, "new_height", view.Height, "new_proposer", c.valSet)
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
