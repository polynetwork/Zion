package core

import (
	"bytes"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/crypto"
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

func (c *core) checkLockedQC(qc *hotstuff.QuorumCert) error {
	if !reflect.DeepEqual(c.current.LockedQC(), qc) {
		return errInconsistentPrepareQC
	}
	return nil
}

func (c *core) checkVote(vote *hotstuff.Vote) error {
	if !reflect.DeepEqual(c.current.Vote(), vote) {
		return errInconsistentVote
	}
	return nil
}

// checkView checks the message state
// return errInvalidMessage if the message is invalid
// return errFutureMessage if the message view is larger than current view
// return errOldMessage if the message view is smaller than current view
func (c *core) checkView(msgCode MsgType, view *hotstuff.View) error {
	if view == nil || view.Height == nil || view.Round == nil {
		return errInvalidMessage
	}

	if msgCode == MsgTypeNewView {
		if view.Height.Cmp(c.currentView().Height) > 0 {
			return errFutureMessage
		} else if view.Cmp(c.currentView()) < 0 {
			return errOldMessage
		}
		return nil
	}

	if view.Cmp(c.currentView()) > 0 {
		return errFutureMessage
	}

	if view.Cmp(c.currentView()) < 0 {
		return errOldMessage
	}
	return nil
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
	logger := c.logger.New("state", c.currentState())

	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	if msg.Code == MsgTypeNewView { // todo: judge current proposal is not nil
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

func (c *core) Q() int {
	return c.valSet.Q()
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(MsgTypeCommit)})
	return buf.Bytes()
}

// GetSignatureAddress gets the signer address from the signature
func GetSignatureAddress(data []byte, sig []byte) (common.Address, error) {
	// 1. Keccak data
	hashData := crypto.Keccak256(data)
	// 2. Recover public key
	pubkey, err := crypto.SigToPub(hashData, sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}

func Proposal2QC(view *hotstuff.View, proposal hotstuff.Proposal) *hotstuff.QuorumCert {
	block := proposal.(*types.Block)
	h := block.Header()
	qc := new(hotstuff.QuorumCert)
	qc.View = view
	qc.Hash = h.Hash()
	qc.Proposer = h.Coinbase
	qc.Extra = h.Extra

	return qc
}
