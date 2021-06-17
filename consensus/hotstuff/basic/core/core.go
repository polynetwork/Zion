package core

import (
	"github.com/ethereum/go-ethereum/consensus"
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
	state   State
	logger  log.Logger

	current *roundState

	chain consensus.ChainReader

	backend             hotstuff.Backend
	events              *event.TypeMuxSubscription
	finalCommittedSub   *event.TypeMuxSubscription
	timeoutSub          *event.TypeMuxSubscription
	futureProposalTimer *time.Timer

	pendingRQ   *prque.Prque
	pendingRQMu *sync.Mutex

	valSet                hotstuff.ValidatorSet
	waitingForRoundChange bool
	validateFn            func([]byte, []byte) (common.Address, error) // == c.checkValidatorSignature
}

func (c *core) Address() common.Address {
	return c.backend.Address()
}

func (c *core) Start() error {
	return nil
}

func (c *core) Stop() error {
	return nil
}

func (c *core) IsProposer() bool {
	return c.valSet.IsProposer(c.backend.Address())
}

func (c *core) IsCurrentProposal(blockHash common.Hash) bool {
	return false
}

func (c *core) startNewRound(round *big.Int) {

}

func (c *core) currentView() *hotstuff.View {
	return &hotstuff.View{
		Height: new(big.Int).Set(c.current.Height()),
		Round:  new(big.Int).Set(c.current.Round()),
	}
}

func (c *core) setState(st State) {
	if st != c.state {
		c.state = st
	}
}

func (c *core) finalizeMessage(msg *message) ([]byte, error) {
	var err error

	// Add sender address
	msg.Address = c.Address()

	// Add proof of consensus
	msg.CommittedSeal = []byte{}
	// Assign the CommittedSeal if it's a COMMIT message and proposal is not nil
	if msg.Code == MsgTypeCommit && c.current.Proposal() != nil {
		seal := PrepareCommittedSeal(c.current.Proposal().Hash())
		msg.CommittedSeal, err = c.backend.Sign(seal)
		if err != nil {
			return nil, err
		}
	}

	// Sign message
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}
	msg.Signature, err = c.backend.Sign(data)
	if err != nil {
		return nil, err
	}

	// Convert to payload
	payload, err := msg.Payload()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *core) broadcast(msg *message) {
	logger := c.logger.New("state", c.state)

	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	if msg.Code == MsgTypeRoundChange { // todo: judge current proposal is not nil
		_, lastProposer := c.backend.LastProposal()
		proposedNewSet := c.valSet.Copy()
		newRound := new(big.Int).Add(c.currentView().Round, common.Big1)
		proposedNewSet.CalcProposer(lastProposer, newRound.Uint64())
		if !proposedNewSet.IsProposer(c.Address()) {
			if err = c.backend.Unicast(proposedNewSet, payload); err != nil {
				logger.Error("Failed to unicast message", "msg", msg, "err", err)
				return
			}
		} else {
			logger.Trace("Local is the next proposer", "msg", msg)
			return
		}
	} else if msg.Code == MsgTypePrepareVote || msg.Code == MsgTypePreCommitVote || msg.Code == MsgTypeCommitVote { // todo: judge current proposal is not nil
		if err := c.backend.Unicast(c.valSet, payload); err != nil {
			logger.Error("Failed to unicast message", "msg", msg, "err", err)
			return
		}
	} else {
		if err := c.backend.Broadcast(c.valSet, payload); err != nil {
			logger.Error("Failed to broadcast message", "msg", msg, "err", err)
			return
		}
	}
}

//func (c *core) broadcast(msg *message) {
//	logger := c.logger.New("state", c.state)
//
//	payload, err := c.finalizeMessage(msg)
//	if err != nil {
//		logger.Error("Failed to finalize message", "msg", msg, "err", err)
//		return
//	}
//
//	// Broadcast payload
//	if err = c.backend.Broadcast(c.valSet, payload); err != nil {
//		logger.Error("Failed to broadcast message", "msg", msg, "err", err)
//		return
//	}
//}

func (c *core) catchUpRound(view *hotstuff.View) {

}

func (c *core) checkValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	return hotstuff.CheckValidatorSignature(c.valSet, data, sig)
}
