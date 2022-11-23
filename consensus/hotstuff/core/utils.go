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

func (c *core) HeightU64() uint64 {
	if c.current == nil {
		return 0
	} else {
		return c.current.HeightU64()
	}
}

func (c *core) RoundU64() uint64 {
	if c.current == nil {
		return 0
	} else {
		return c.current.RoundU64()
	}
}

func (c *core) checkMsgSource(src common.Address) error {
	if !c.valSet.IsProposer(src) {
		return errNotFromProposer
	}
	return nil
}

func (c *core) checkMsgDest() error {
	if !c.IsProposer() {
		return errNotToProposer
	}
	return nil
}

// checkVote vote should equal to current node hash
func (c *core) checkVote(data *Message, vote common.Hash) error {
	if vote == common.EmptyHash {
		return fmt.Errorf("external vote is empty hash")
	}
	if c.current.Vote() == common.EmptyHash {
		return fmt.Errorf("current vote is nil")
	}
	if !reflect.DeepEqual(c.current.Vote(), vote) {
		return fmt.Errorf("expect %s, got %s", c.current.Vote().Hex(), vote.Hex())
	}
	return nil
}

// checkBlock check the extend relationship between remote block and latest chained block.
// ensure that the remote block equals to locked block if it exist.
func (c *core) checkBlock(block *types.Block) error {
	lastChainedBlock := c.current.LastChainedBlock()
	if lastChainedBlock.NumberU64()+1 != block.NumberU64() {
		return fmt.Errorf("expect block number %v, got %v", lastChainedBlock.NumberU64()+1, block.NumberU64())
	}
	if lastChainedBlock.Hash() != block.ParentHash() {
		return fmt.Errorf("expect parent hash %v, got %v", lastChainedBlock.Hash(), block.ParentHash())
	}

	if lockedBlock := c.current.LockedBlock(); lockedBlock != nil {
		if block.NumberU64() != lockedBlock.NumberU64() {
			return fmt.Errorf("expect locked block number %v, got %v", lockedBlock.NumberU64(), block.NumberU64())
		}
		if block.Hash() != lockedBlock.Hash() {
			return fmt.Errorf("expect locked block hash %v, got %v", lockedBlock.Hash(), block.Hash())
		}
	}

	return nil
}

// checkNode repo compare remote node with local node
func (c *core) checkNode(node *Node, compare bool) error {
	if node == nil || node.Parent == common.EmptyHash ||
		node.Block == nil || node.Block.Header() == nil {
		return errInvalidNode
	}

	if !compare {
		return nil
	}

	local := c.current.Node()
	if local == nil {
		return fmt.Errorf("current node is nil")
	}
	if local.Hash() != node.Hash() {
		return fmt.Errorf("expect node %v but got %v", local.Hash(), node.Hash())
	}
	if local.Block.Hash() != node.Block.Hash() {
		return fmt.Errorf("expect block %v but got %v", local.Block.Hash(), node.Block.Hash())
	}
	return nil
}

// checkView checks the Message sequence remote msg view should not be nil(local view WONT be nil).
// if the view is ahead of current view we name the Message to be future Message, and if the view is behind
// of current view, we name it as old Message. `old Message` and `invalid Message` will be dropped. and we use t
// he storage of `backlog` to cache the future Message, it only allow the Message height not bigger than
// `current height + 1` to ensure that the `backlog` memory won't be too large, it won't interrupt the consensus
// process, because that the `core` instance will sync block until the current height to the correct value.
//
//
// todo(fuk):if the view is equal the current view, compare the Message type and round state, with the right
// round state sequence, Message ahead of certain state is `old Message`, and Message behind certain
// state is `future Message`. Message type and round state table as follow:
func (c *core) checkView(view *View) error {
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

func (c *core) preExecuteBlock(proposal hotstuff.Proposal) error {
	if c.IsProposer() {
		return nil
	}

	block, ok := proposal.(*types.Block)
	if !ok {
		return errInvalidNode
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

var (
	genesisView = &View{
		Round:  big.NewInt(0),
		Height: big.NewInt(0),
	}
	genesisNodeHash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000012345")
)

func genesisQC(lastBlock *types.Block) (*QuorumCert, error) {
	qc := &QuorumCert{
		view:          genesisView,
		code:          MsgTypePrepareVote,
		node:          genesisNodeHash,
		proposer:      common.Address{},
		seal:          make([]byte, 0),
		committedSeal: make([][]byte, 0),
	}

	h := lastBlock.Header()
	extra, err := types.ExtractHotstuffExtra(h)
	if err != nil {
		return nil, err
	}
	if extra.Seal == nil || extra.CommittedSeal == nil {
		return nil, errInvalidNode
	}

	copy(qc.seal, extra.Seal)
	copy(qc.committedSeal, extra.CommittedSeal)
	return qc, nil
}

// assemble messages to quorum cert.
func (c *core) messages2qc(code MsgType) (*QuorumCert, error) {
	var msgs []*Message
	switch code {
	case MsgTypePrepareVote:
		msgs = c.current.PrepareVotes()
	case MsgTypePreCommitVote:
		msgs = c.current.PreCommitVotes()
	case MsgTypeCommitVote:
		msgs = c.current.CommitVotes()
	default:
		return nil, errInvalidCode
	}
	if len(msgs) == 0 {
		return nil, fmt.Errorf("assemble qc: not enough message")
	}

	var (
		proposer = c.proposer()
		view     = c.currentView()
		node     = c.current.Vote()
		sealHash = msgs[0].hash
	)

	if node == common.EmptyHash {
		return nil, fmt.Errorf("current vote is empty")
	}
	if sealHash == common.EmptyHash {
		return nil, fmt.Errorf("message hash is empty")
	}

	qc := &QuorumCert{
		view:          view,
		code:          code,
		node:          node,
		proposer:      proposer,
		committedSeal: make([][]byte, len(msgs)),
	}

	for i, msg := range msgs {
		if msg.hash != sealHash {
			return nil, fmt.Errorf("vote seal hash expect %v got %v", sealHash, msg.hash)
		}
		if msg.View.Cmp(view) != 0 {
			return nil, fmt.Errorf("vote view expect %v, got %v", view, msg.View)
		}
		if msg.Signature == nil {
			return nil, fmt.Errorf("vote signature nil")
		}
		if msg.address == proposer {
			qc.seal = msg.Signature
		}
		qc.committedSeal[i] = msg.Signature
	}

	// proposer self vote should be add in message set first.
	if qc.seal == nil {
		if sig, err := c.signer.SignHash(sealHash, false); err != nil {
			return nil, err
		} else {
			qc.seal = sig
		}
	}

	return qc, nil
}

// verifyQC check and validate qc.
func (c *core) verifyQC(data *Message, qc *QuorumCert) error {
	if qc == nil || qc.view == nil {
		return errInvalidQC
	}

	// skip genesis block
	if qc.HeightU64() == 0 {
		return nil
	}

	// qc fields checking
	if qc.node == common.EmptyHash || qc.seal == nil || qc.committedSeal == nil ||
		(qc.proposer == common.EmptyAddress && qc.HeightU64() > 1) {
		return errInvalidQC
	}

	// matching code
	if (data.Code == MsgTypeNewView && qc.code != MsgTypePrepareVote) ||
		(data.Code == MsgTypePrepare && qc.code != MsgTypePrepareVote) ||
		(data.Code == MsgTypePreCommit && qc.code != MsgTypePrepareVote) ||
		(data.Code == MsgTypeCommit && qc.code != MsgTypePreCommitVote) ||
		(data.Code == MsgTypeDecide && qc.code != MsgTypeCommitVote) {
		return fmt.Errorf("qc.code %s not matching message code", qc.code.String())
	}

	// matching view and compare proposer and local node
	var valset hotstuff.ValidatorSet
	if data.Code == MsgTypePreCommit || data.Code == MsgTypeCommit || data.Code == MsgTypeDecide {
		if qc.view.Cmp(data.View) != 0 {
			return fmt.Errorf("qc.view %v not matching message view", qc.view)
		}
		if qc.proposer != c.proposer() {
			return fmt.Errorf("expect proposer %v, got %v", c.proposer(), qc.proposer)
		}
		if node := c.current.Node(); node == nil {
			return fmt.Errorf("current node is nil")
		} else if node.Hash() != qc.node {
			return fmt.Errorf("expect node %v, got %v", node.Hash(), qc.node)
		}
		valset = c.valSet.Copy()
	} else {
		if hdiff, rdiff := data.View.Sub(qc.view); hdiff < 0 || (hdiff == 0 && rdiff < 0) {
			return fmt.Errorf("view is invalid")
		}
		valset = c.backend.Validators(qc.HeightU64(), false)
	}

	// resturct msg payload and compare msg.hash with qc.hash
	msg := NewCleanMessage(qc.view, qc.code, qc.node.Bytes())
	if _, err := msg.PayloadNoSig(); err != nil {
		return fmt.Errorf("payload no sig")
	}
	if sealHash := qc.SealHash(); msg.hash != sealHash {
		return fmt.Errorf("expect qc hash %v, got %v", msg.hash, sealHash)
	}

	// find the correct validator set and verify seal & committed seals
	return c.signer.VerifyQC(qc, valset)
}

// sendVote repo send kinds of vote to leader, use `current.node` after repo `prepared`.
func (c *core) sendVote(code MsgType, votes ...common.Hash) {
	logger := c.newLogger()

	var vote common.Hash
	if len(votes) == 0 {
		vote = c.current.Vote()
	} else {
		vote = votes[0]
	}
	c.broadcast(code, vote.Bytes())
	prefix := fmt.Sprintf("send%s", code.String())
	logger.Trace(prefix, "msg", code, "hash", vote)
}
