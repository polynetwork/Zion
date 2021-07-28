package core

import (
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

func (c *core) checkMsgFromProposer(src hotstuff.Validator) error {
	if !c.valSet.IsProposer(src.Address()) {
		return errNotFromProposer
	}
	return nil
}

func (c *core) checkMsgToProposer() error {
	if !c.IsProposer() {
		return errNotToProposer
	}
	return nil
}

func (c *core) checkPrepareQC(qc *hotstuff.QuorumCert) error {
	if !reflect.DeepEqual(c.current.PrepareQC(), qc) {
		return errInconsistentPrepareQC
	}
	return nil
}

func (c *core) checkPreCommittedQC(qc *hotstuff.QuorumCert) error {
	if !reflect.DeepEqual(c.current.PreCommittedQC(), qc) {
		return errInconsistentPrepareQC
	}
	return nil
}

func (c *core) checkVote(vote *Vote) error {
	if !reflect.DeepEqual(c.current.Vote(), vote) {
		return errInconsistentVote
	}
	return nil
}

func (c *core) checkLockedProposal(msg hotstuff.Proposal) error {
	isLocked, proposal := c.current.LastLockedProposal()
	if !isLocked {
		return nil
	}
	if proposal == nil {
		return errLockedProposal
	}
	if !reflect.DeepEqual(proposal, msg) {
		return errLockedProposal
	}
	return nil
}

// checkView checks the message state, msg view should not be nil. if the view is ahead of current view
// we name the message to be future message, and if the view is behind of current view, we name it as old
// message. `old message` and `invalid message` will be dropped . and we use the storage of `backlog` to
// cache the future message, it only allow the message height not bigger than `current height + 1` to ensure
// that the `backlog` memory won't be too large, it won't interrupt the consensus process, because that the
// `core` instance will sync block until the current height to the correct value.
//
// if the view is equal the current view, compare the message type and round state, with the right round state sequence,
// message ahead of certain state is `old message`, and message behind certain state is `future message`.
// message type and round state table as follow:
func (c *core) checkView(msgCode MsgType, view *hotstuff.View) error {
	if view == nil || view.Height == nil || view.Round == nil {
		return errInvalidMessage
	}

	// validators not in the same view
	if hdiff, rdiff := view.Sub(c.currentView()); hdiff < 0 {
		return errOldMessage
	} else if hdiff > 1 {
		return errFarAwayFutureMessage
	} else if hdiff == 1 {
		return errFutureMessage
	} else if rdiff < 0 {
		return errOldMessage
	} else if rdiff == 0 {
		return nil
	} else {
		return errFutureMessage
	}
}

func (c *core) finalizeMessage(msg *message) ([]byte, error) {
	var err error

	// Add sender address
	msg.Address = c.Address()
	msg.View = c.currentView()

	// Add proof of consensus
	proposal := c.current.Proposal()
	if msg.Code == MsgTypePrepareVote && proposal != nil {
		seal, err := c.signer.SignVote(proposal)
		if err != nil {
			return nil, err
		}
		msg.CommittedSeal = seal
	}

	// Sign message
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}
	msg.Signature, err = c.signer.Sign(data)
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

func (c *core) getMessageSeals(n int) [][]byte {
	seals := make([][]byte, n)
	for i, data := range c.current.PrepareVotes() {
		if i < n {
			seals[i] = data.CommittedSeal
		}
	}
	return seals
}

func (c *core) broadcast(msg *message) {
	logger := c.logger.New("state", c.currentState())

	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	// todo: judge current proposal is not nil
	switch msg.Code {
	case MsgTypeNewView, MsgTypePrepareVote, MsgTypePreCommitVote, MsgTypeCommitVote:
		if err := c.backend.Unicast(c.valSet, payload); err != nil {
			logger.Error("Failed to unicast message", "msg", msg, "err", err)
		}
	case MsgTypePrepare, MsgTypePreCommit, MsgTypeCommit, MsgTypeDecide:
		if err := c.backend.Broadcast(c.valSet, payload); err != nil {
			logger.Error("Failed to broadcast message", "msg", msg, "err", err)
		}
	default:
		logger.Error("invalid msg type", "msg", msg)
	}
}

func (c *core) checkValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	return c.signer.CheckSignature(c.valSet, data, sig)
}

func (c *core) newLogger() log.Logger {
	logger := c.logger.New("state", c.currentState(), "view", c.currentView())
	return logger
}

func proposal2QC(proposal hotstuff.Proposal, round *big.Int) *hotstuff.QuorumCert {
	block := proposal.(*types.Block)
	h := block.Header()
	qc := new(hotstuff.QuorumCert)
	qc.View = &hotstuff.View{
		Height: block.Number(),
		Round:  round,
	}
	qc.Hash = h.Hash()
	qc.Proposer = h.Coinbase
	qc.Extra = h.Extra
	return qc
}
