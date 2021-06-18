package core

import (
	"github.com/ethereum/go-ethereum/consensus"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/prque"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

type core struct {
	config  *hotstuff.Config
	address common.Address
	logger  log.Logger

	state         State
	current       *roundState
	changeViewSet *changeViewSet
	backend       hotstuff.Backend
	valSet        hotstuff.ValidatorSet
	chain         consensus.ChainReader

	events            *event.TypeMuxSubscription
	finalCommittedSub *event.TypeMuxSubscription
	timeoutSub        *event.TypeMuxSubscription
	handlerWg         *sync.WaitGroup

	futureProposalTimer *time.Timer

	pendingRequest   *prque.Prque
	pendingRequestMu *sync.Mutex

	backlogs   map[common.Address]*prque.Prque // todo back log module, processing future messages
	backlogsMu *sync.Mutex

	consensusTimestamp time.Time
	roundChangeTimer   *time.Timer
	futurePrepareTimer *time.Timer

	waitingForRoundChange bool
	validateFn            func([]byte, []byte) (common.Address, error) // == c.checkValidatorSignature
}

// New creates an HotStuff consensus core
func New(backend hotstuff.Backend, config *hotstuff.Config) CoreEngine {
	c := &core{
		config:             config,
		state:              StateAcceptRequest,
		address:            backend.Address(),
		logger:             log.New("address", backend.Address()),
		handlerWg:          new(sync.WaitGroup),
		backend:            backend,
		backlogs:           make(map[common.Address]*prque.Prque),
		backlogsMu:         new(sync.Mutex),
		pendingRequest:     prque.New(),
		pendingRequestMu:   new(sync.Mutex),
		consensusTimestamp: time.Time{},
	}
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
	return c.current != nil && c.current.pendingRequest != nil && c.current.pendingRequest.Proposal.Hash() == blockHash
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

	var newView *hotstuff.View
	if roundChange {
		newView = &hotstuff.View{
			Height: new(big.Int).Set(c.current.Height()),
			Round:    new(big.Int).Set(round),
		}
	} else {
		newView = &hotstuff.View{
			Height: new(big.Int).Add(lastProposal.Number(), common.Big1),
			Round:    new(big.Int),
		}
		c.valSet = c.backend.Validators(lastProposal)
	}

	// Update logger
	logger = logger.New("old_proposer", c.valSet.GetProposer())
	// Clear invalid ROUND CHANGE messages
	c.changeViewSet = newChangeViewSet(c.valSet)
	// New snapshot for new round
	c.updateRoundState(newView, c.valSet, roundChange)
	// Calculate new proposer
	c.valSet.CalcProposer(lastProposer, newView.Round.Uint64())
	c.waitingForRoundChange = false
	c.setState(StateAcceptRequest)
	if roundChange && c.IsProposer() && c.current != nil {
		// If it is locked, propose the old proposal
		// If we have pending request, propose pending request
		if c.current.IsHashLocked() {
			r := &hotstuff.Request{
				Proposal: c.current.Proposal(), //c.current.Proposal would be the locked proposal by previous proposer, see updateRoundState
			}
			c.sendPrepare(r)
		} else if c.current.pendingRequest != nil {
			c.sendPrepare(c.current.pendingRequest)
		}
	}
	c.newRoundChangeTimer()

	logger.Debug("New round", "last block number", lastProposal.Number().Uint64(), "new_round", newView.Round, "new_heigth", newView.Height, "new_proposer", c.valSet.GetProposer(), "valSet", c.valSet.List(), "size", c.valSet.Size(), "IsProposer", c.IsProposer())
}

func (c *core) currentView() *hotstuff.View {
	return &hotstuff.View{
		Height: new(big.Int).Set(c.current.Height()),
		Round:  new(big.Int).Set(c.current.Round()),
	}
}

func (c *core) catchUpRound(view *hotstuff.View) {
	logger := c.logger.New("old_round", c.current.Round(), "old_height", c.current.Height(), "old_proposer", c.valSet.GetProposer())

	c.waitingForRoundChange = true

	// Need to keep block locked for round catching up
	c.updateRoundState(view, c.valSet, true)
	c.changeViewSet.Clear(view.Round)
	c.newRoundChangeTimer()

	logger.Trace("Catch up round", "new_round", view.Round, "new_height", view.Height, "new_proposer", c.valSet)
}

func (c *core) updateRoundState(view *hotstuff.View, valSet hotstuff.ValidatorSet, changeView bool) {
	// Lock only if both roundChange is true and it is locked
	if changeView && c.current != nil {
		if c.current.IsHashLocked() {
			c.current = newRoundState(view, valSet, c.current.GetLockedHash(), c.current.prepare, c.current.pendingRequest, c.backend.HasBadProposal)
		} else {
			c.current = newRoundState(view, valSet, common.Hash{}, nil, c.current.pendingRequest, c.backend.HasBadProposal)
		}
	} else {
		c.current = newRoundState(view, valSet, common.Hash{}, nil, nil, c.backend.HasBadProposal)
	}
}

func (c *core) handleFinalCommitted() error {
	logger := c.logger.New("state", c.state)
	logger.Trace("Received a final committed proposal")
	c.startNewRound(common.Big0)
	return nil
}

func (c *core) setState(state State) {
	if state != c.state {
		c.state = state
	}
	if state == StateAcceptRequest {
		c.processPendingRequest()
	}
	//c.processBacklog()
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