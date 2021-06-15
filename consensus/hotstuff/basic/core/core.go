package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"time"
)

type core struct {
	config  *hotstuff.Config
	address common.Address
	state   State
	logger  log.Logger

	backend             hotstuff.Backend
	events              *event.TypeMuxSubscription
	finalCommittedSub   *event.TypeMuxSubscription
	timeoutSub          *event.TypeMuxSubscription
	futureProposalTimer *time.Timer

	valSet                hotstuff.ValidatorSet
	waitingForRoundChange bool
	validateFn            func([]byte, []byte) (common.Address, error)
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
	return false
}

func (c *core) IsCurrentProposal(blockHash common.Hash) bool {
	return false
}

func (c *core) finalizeMessage(msg *message) ([]byte, error) {
	var err error

	// Add sender address
	msg.Address = c.Address()

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

func (c *core) broadcast(msg *message, round *big.Int) {
	logger := c.logger.New("state", c.state)

	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	if msg.Code == MsgTypeNewView {  // todo: judge current proposal is not nil
		_, lastProposer := c.backend.LastProposal()
		proposedNewSet := c.valSet.Copy()
		proposedNewSet.CalcProposer(lastProposer, round.Uint64())
		if !proposedNewSet.IsProposer(c.Address()) {
			if err = c.backend.Unicast(proposedNewSet, payload); err != nil {
				logger.Error("Failed to unicast message", "msg", msg, "err", err)
				return
			}
		} else {
			logger.Trace("Local is the next proposer", "msg", msg)
			return
		}
	} else if msg.Code == MsgTypePreCommitVote || msg.Code == MsgTypePrepareVote { // todo: judge current proposal is not nil
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