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
	config  *hotstuff.Config
	address common.Address
	logger  log.Logger

	current  *roundState
	backend  hotstuff.Backend
	valSet   hotstuff.ValidatorSet
	requests *requestSet

	events            *event.TypeMuxSubscription
	timeoutSub        *event.TypeMuxSubscription
	finalCommittedSub *event.TypeMuxSubscription

	futureProposalTimer *time.Timer
	consensusTimestamp  time.Time
	roundChangeTimer    *time.Timer
	futurePrepareTimer  *time.Timer

	validateFn func([]byte, []byte) (common.Address, error) // == c.checkValidatorSignature
}

// New creates an HotStuff consensus core
func New(backend hotstuff.Backend, config *hotstuff.Config, valSet hotstuff.ValidatorSet) CoreEngine {
	c := &core{
		config:             config,
		address:            backend.Address(),
		valSet:             valSet,
		logger:             log.New("address", backend.Address()),
		backend:            backend,
		consensusTimestamp: time.Time{},
	}
	c.requests = newRequestSet()
	c.validateFn = c.checkValidatorSignature
	return c
}

func (c *core) Address() common.Address {
	return c.backend.Address()
}

func (c *core) IsProposer() bool {
	return c.valSet.IsProposer(c.backend.Address())
}

func (c *core) IsCurrentProposal(blockHash common.Hash) bool {
	return c.current != nil && c.current.Proposal() != nil && c.current.Proposal().Hash() == blockHash
}

func (c *core) startNewRound(round *big.Int) {
	var logger log.Logger
	if c.current == nil {
		logger = c.logger.New("old_round", -1, "old_height", 0)
	} else {
		logger = c.logger.New("old_round", c.current.Round(), "old_height", c.current.Height())
	}

	roundChange := false
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
		roundChange = true
	} else {
		logger.Warn("New height should be larger than current height", "new_height", lastProposal.Number().Int64())
		return
	}

	newView := &hotstuff.View{
		Height: new(big.Int).Add(lastProposal.Number(), common.Big1),
		Round:  new(big.Int),
	}
	if roundChange {
		newView.Height = new(big.Int).Set(c.current.Height())
		newView.Round = new(big.Int).Set(round)
	}
	prepareQC := &QuorumCert{
		View:     newView,
		Proposal: lastProposal,
	}
	c.valSet.CalcProposer(lastProposer, newView.Round.Uint64())
	if c.current == nil {
		c.current = newRoundState(newView, c.valSet, prepareQC)
	} else {
		c.current = c.current.Spawn(newView)
	}
	c.current = c.current.Spawn(newView)
	c.sendNewView(newView)
	c.newRoundChangeTimer()

	logger.Debug("New round", "last block number", lastProposal.Number().Uint64(), "new_round", newView.Round, "new_heigth", newView.Height, "new_proposer", c.valSet.GetProposer(), "valSet", c.valSet.List(), "size", c.valSet.Size(), "IsProposer", c.IsProposer())
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

// todo: 检查是否收到了新的proposal
func (c *core) catchUpRound(view *hotstuff.View) {
	logger := c.logger.New("old_round", c.current.Round(), "old_height", c.current.Height(), "old_proposer", c.valSet.GetProposer())

	c.newRoundChangeTimer()

	logger.Trace("Catch up round", "new_round", view.Round, "new_height", view.Height, "new_proposer", c.valSet)
}

func (c *core) stopFuturePrepareTimer() {
	if c.futurePrepareTimer != nil {
		c.futurePrepareTimer.Stop()
	}
}

func (c *core) stopTimer() {
	c.stopFuturePrepareTimer()
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

//func (c *core) setState(state State) {
//	c.current.SetState(state)
//	//if state == StateAcceptRequest {
//	//	c.GetRequest()
//	//}
//	//c.processBacklog()
//}
