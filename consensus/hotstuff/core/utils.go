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
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

func (c *core) proposer() common.Address {
	return c.valSet.GetProposer().Address()
}

func (c *core) checkMsgFromProposer(src common.Address) error {
	if !c.valSet.IsProposer(src) {
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

func (c *core) checkPrepareQC(qc *QuorumCert) error {
	if qc == nil {
		return fmt.Errorf("external prepare qc is nil")
	}

	localQC := c.current.PrepareQC()
	if localQC == nil {
		return fmt.Errorf("current prepare qc is nil")
	}

	if localQC.view.Cmp(qc.view) != 0 {
		return fmt.Errorf("view unsame, expect %v, got %v", localQC.view, qc.view)
	}
	if localQC.proposer != qc.proposer {
		return fmt.Errorf("proposer unsame, expect %v, got %v", localQC.Proposer(), qc.Proposer())
	}
	if localQC.hash != qc.hash {
		return fmt.Errorf("expect %v, got %v", localQC.hash, qc.hash)
	}
	return nil
}

func (c *core) checkPreCommittedQC(qc *QuorumCert) error {
	if qc == nil {
		return fmt.Errorf("external pre-committed qc is nil")
	}
	if c.current.PreCommittedQC() == nil {
		return fmt.Errorf("current pre-committed qc is nil")
	}
	if !reflect.DeepEqual(c.current.PreCommittedQC(), qc) {
		return fmt.Errorf("expect %s, got %s", c.current.PreCommittedQC().String(), qc.String())
	}
	return nil
}

func (c *core) checkVote(data *Message, vote common.Hash) error {
	if vote == common.EmptyHash {
		return fmt.Errorf("external vote is empty hash")
	}
	if c.current.Vote() == common.EmptyHash {
		return fmt.Errorf("current vote is nil")
	}
	if !reflect.DeepEqual(c.current.Vote(), vote) {
		return fmt.Errorf("expect %s, got %s", c.current.Vote().String(), vote.String())
	}
	// todo:
	//if hash, err := c.current.SelfVoteHash(data.View, data.Code); err != nil {
	//	return fmt.Errorf("get self vote hash failed, err: %v", err)
	//} else if hash != data.hash {
	//	return fmt.Errorf("expect vote hash %v, got %v", hash, data.hash)
	//}
	return nil
}

// todo(fuk): 这里应该考虑将qc分成两种，一种用于投票，一种用于区块，两种qc公用同一套接口，然后在signer那边可以共同验证
//func (c *core) votes2qc(hash common.Hash, seals [][]byte) {
//	QuorumCert{
//		view:          nil,
//		hash:          common.Hash{},
//		proposer:      common.Address{},
//		seal:          nil,
//		committedSeal: nil,
//	}
//}

func (c *core) checkProposal(hash common.Hash) error {
	if c.current == nil || c.current.Proposal() == nil {
		return fmt.Errorf("current proposal is nil")
	}
	if expect := c.current.Proposal().Hash(); hash != expect {
		return fmt.Errorf("hash expect %s got %s", expect.Hex(), hash.Hex())
	}
	return nil
}

func (c *core) checkProposalView(proposal hotstuff.Proposal, view *View) error {
	if proposal == nil || view == nil || view.Height == nil {
		return fmt.Errorf("proposal or view is invalid")
	} else if proposal.NumberU64() != view.Height.Uint64() {
		return fmt.Errorf("proposal height %v != view %v", proposal.NumberU64(), view.Height)
	} else {
		return nil
	}
}

// verifyCrossEpochQC verify quorum certificate with current validator set or
// last epoch's val set if current height equals to epoch start height
func (c *core) verifyCrossEpochQC(qc *QuorumCert) error {
	valset := c.backend.Validators(qc.hash, false)
	if err := c.signer.VerifyQC(qc, valset); err != nil {
		return err
	}
	return nil
}

// checkView checks the Message state, remote msg view should not be nil(local view WONT be nil).
// if the view is ahead of current view we name the Message to be future Message, and if the view
// is behind of current view, we name it as old Message. `old Message` and `invalid Message` will
// be dropped. and we use the storage of `backlog` to cache the future Message, it only allow the
// Message height not bigger than `current height + 1` to ensure that the `backlog` memory won't be
// too large, it won't interrupt the consensus process, because that the `core` instance will sync
// block until the current height to the correct value.
//
// if the view is equal the current view, compare the Message type and round state, with the right
// round state sequence, Message ahead of certain state is `old Message`, and Message behind certain
// state is `future Message`. Message type and round state table as follow:
func (c *core) checkView(msgCode MsgType, view *View) error { //todo
	if view == nil || view.Height == nil || view.Round == nil {
		return errInvalidMessage
	}

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

func (c *core) finalizeMessage(msg *Message) ([]byte, error) {
	var (
		seal, sig []byte
		err       error
	)

	// Add proof of consensus
	proposal := c.current.Proposal()
	if msg.Code == MsgTypeCommitVote && proposal != nil {
		if seal, err = c.signer.SignHash(proposal.Hash(), true); err != nil {
			return nil, err
		}
		msg.CommittedSeal = seal
	}

	// Sign Message
	if _, err = msg.PayloadNoSig(); err != nil {
		return nil, err
	}
	if sig, err = c.signer.SignHash(msg.hash, false); err != nil {
		return nil, err
	} else {
		msg.Signature = sig
	}

	// Convert to payload
	return msg.Payload()
}

func (c *core) broadcast(code MsgType, payload []byte) {
	logger := c.logger.New("state", c.currentState())

	// todo(fuk): forbid at start at new round
	// forbid normal nodes send message to leader
	if index, _ := c.valSet.GetByAddress(c.Address()); index < 0 {
		return
	}

	msg := NewCleanMessage(c.currentView(), code, payload)
	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize Message", "msg", msg, "err", err)
		return
	}

	// todo(fuk): 这里需要注意的是，prepareQC和preCommitQC以及commitQC所需要的3轮投票中要提前将
	// 自己的投票
	switch msg.Code {
	case MsgTypeNewView, MsgTypePrepareVote, MsgTypePreCommitVote, MsgTypeCommitVote:
		if err := c.backend.Unicast(c.valSet, payload); err != nil {
			logger.Error("Failed to unicast Message", "msg", msg, "err", err)
		}

	case MsgTypePrepare, MsgTypePreCommit, MsgTypeCommit, MsgTypeDecide:
		if err := c.backend.Broadcast(c.valSet, payload); err != nil {
			logger.Error("Failed to broadcast Message", "msg", msg, "err", err)
		}
	default:
		logger.Error("invalid msg type", "msg", msg)
	}
}

func (c *core) preExecuteBlock(proposal hotstuff.Proposal) error {
	if c.IsProposer() {
		return nil
	}

	block, ok := proposal.(*types.Block)
	if !ok {
		return errInvalidProposal
	}
	return c.backend.ValidateBlock(block)
}

func (c *core) newLogger() log.Logger {
	logger := c.logger.New("state", c.currentState(), "view", c.currentView())
	return logger
}

func (c *core) Q() int {
	return c.valSet.Q()
}

func proposal2QC(proposal hotstuff.Proposal) (*QuorumCert, error) {
	block := proposal.(*types.Block)
	h := block.Header()
	qc := EmptyQC()

	extra, err := types.ExtractHotstuffExtra(h)
	if err != nil {
		return nil, err
	}

	if h.Hash() == common.EmptyHash || extra.Seal == nil || extra.CommittedSeal == nil {
		return nil, errInvalidProposal
	}

	qc.view.Height = new(big.Int).SetUint64(h.Number.Uint64())
	qc.view.Round = big.NewInt(0)
	copy(qc.hash[:], h.Hash().Bytes())
	copy(qc.seal, extra.Seal)
	copy(qc.committedSeal, extra.CommittedSeal)
	fmt.Println("----qc.hash", qc.hash.Hex())
	return qc, nil
}

// assemble messages to quorum cert.
func (c *core) messages2qc(proposer common.Address, hash common.Hash, msgs []*Message) (*QuorumCert, error) {
	if len(msgs) == 0 {
		return nil, fmt.Errorf("assemble qc: not enough message")
	}

	qc := &QuorumCert{
		view:          msgs[0].View,
		code:          msgs[0].Code,
		hash:          hash,
		proposer:      proposer,
		committedSeal: make([][]byte, len(msgs)),
	}

	for i, msg := range msgs {
		if msg.address == proposer {
			qc.seal = msg.Signature
		}
		qc.committedSeal[i] = msg.Signature
	}

	// proposer self vote should be add in message set first.
	if qc.seal == nil {
		return nil, fmt.Errorf("assemble qc: proposer self vote not exist")
	}

	return qc, nil
}

func (c *core) verifyVoteQC(digest common.Hash, qc *QuorumCert) error {
	// check qc view
	if qc == nil || qc.view == nil {
		return errInvalidQC
	}

	// verify genesis qc
	if qc.HeightU64() == 0 {
		return nil
	}

	// check qc hash and sigs
	if qc.seal == nil || qc.committedSeal == nil ||
		qc.hash == common.EmptyHash || qc.proposer == common.EmptyAddress {
		return errInvalidQC
	}

	// qc code should be vote
	if !checkQCCode(qc.code) {
		return errInvalidQC
	}

	// resturct msg payload and compare msg.hash with qc.hash
	msg := NewCleanMessage(qc.view, qc.code, qc.hash.Bytes())
	if _, err := msg.PayloadNoSig(); err != nil {
		return err
	}
	if sealHash := qc.SealHash(); msg.hash != sealHash {
		return fmt.Errorf("expect qc hash %v, got %v", msg.hash, sealHash)
	}

	// verify seal and committed seals
	return c.signer.VerifyQC(qc, c.valSet)
}

// todo:
func (c *core) verifyProposalQC() error {
	return nil
}

// qc comes from vote
func checkQCCode(code MsgType) bool {
	if code == MsgTypePrepareVote || code == MsgTypePreCommitVote || code == MsgTypeCommitVote {
		return true
	}
	return false
}
