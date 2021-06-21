package core

import (
	"bytes"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) decodeAndCheckProposal(data *message, msgTyp MsgType, msg *MsgNewProposal) error {
	if err := data.Decode(msg); err != nil {
		return err
	}

	if data.Code != msgTyp {
		return errMsgTypeInvalid
	}

	if msgTyp == MsgTypeNewView && !c.IsProposer() {
		return errNotToProposer
	}

	if msgTyp == MsgTypeCommit && !reflect.DeepEqual(c.current.PrepareQC(), msg) {
		return errInconsistentQC
	}
	return c.checkMessage(data.Code, msg.View)
}

func (c *core) decodeAndCheckMessage(data *message, msgTyp MsgType, msg *QuorumCert) error {
	if err := data.Decode(msg); err != nil {
		return err
	}

	if data.Code != msgTyp {
		return errMsgTypeInvalid
	}

	if msgTyp == MsgTypeNewView && !c.IsProposer() {
		return errNotToProposer
	}

	if msgTyp == MsgTypeCommit && !reflect.DeepEqual(c.current.PrepareQC(), msg) {
		return errInconsistentQC
	}
	return c.checkMessage(data.Code, msg.View)
}

func (c *core) decodeAndCheckVote(data *message, msgTyp MsgType, vote *hotstuff.Subject) error {
	if err := data.Decode(vote); err != nil {
		return err
	}

	if !c.IsProposer() {
		return errNotToProposer
	}

	if data.Code != msgTyp {
		return errMsgTypeInvalid
	}

	if err := c.checkMessage(data.Code, vote.View); err != nil {
		return err
	}

	sub := c.current.Subject()
	if !reflect.DeepEqual(sub, vote) {
		return errInconsistentSubject
	}

	return nil
}

// checkMessage checks the message state
// return errInvalidMessage if the message is invalid
// return errFutureMessage if the message view is larger than current view
// return errOldMessage if the message view is smaller than current view
func (c *core) checkMessage(msgCode MsgType, view *hotstuff.View) error {
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

	//if c.waitingForRoundChange {
	//	return errFutureMessage
	//}

	// todo:
	//
	//// StateAcceptRequest only accepts msgPreprepare
	//// other messages are future messages
	//if c.state == StateAcceptRequest {
	//	if msgCode > MsgTypePrepare {
	//		return errFutureMessage
	//	}
	//	return nil
	//}
	//
	//// For states(StatePreprepared, StatePrepared, StateCommitted),
	//// can accept all message types if processing with same view
	//
	//if msgCode == MsgTypeNewView && c.currentState() >= StatePrepared {
	//	return errFutureMessage
	//}
	//if msgCode == MsgTypePrepare && c.currentState() >= StatePrepared {
	//	return errFutureMessage
	//}
	//if msgCode == MsgTypePreCommit && c.currentState() >= StateLocked {
	//	return errFutureMessage
	//}
	//
	//if msgCode == MsgTypePrepareVote && c.currentState() >= StatePrepared {
	//	return errOldVote
	//}
	//if msgCode == MsgTypePreCommitVote && c.currentState() >= StateLocked {
	//	return errOldVote
	//}
	//if msgCode == MsgTypeCommitVote && c.currentState() >= StateCommitted {
	//	return errOldVote
	//}
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
