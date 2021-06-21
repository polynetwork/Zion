package core

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"math/big"
)

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

	if msg.Code == MsgTypeChangeView { // todo: judge current proposal is not nil
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

func (c *core) checkValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	return hotstuff.CheckValidatorSignature(c.valSet, data, sig)
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(MsgTypeCommit)})
	return buf.Bytes()
}
