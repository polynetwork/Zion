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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

// increaseLastVoteRound commit not to vote in rounds lower than target
func (e *EventDrivenEngine) increaseLastVoteRound(rd *big.Int) {
	if e.lastVoteRound.Cmp(rd) < 0 {
		e.lastVoteRound = rd
	}
}

// UpdateLockQC update the latest quorum certificate after voteRule judgement succeed.
func (e *EventDrivenEngine) updateLockQCRound(round *big.Int) {
	e.lockQCRound = round
}

func (e *EventDrivenEngine) makeVote(hash common.Hash, proposer common.Address, view *hotstuff.View, justifyQC *hotstuff.QuorumCert) (*Vote, error) {
	justifyQCRound := justifyQC.View.Round
	justifyQCHeight := justifyQC.View.Height

	justifyQCBlock := e.blkTree.GetBlockAndCheckHeight(justifyQC.Hash, justifyQCHeight)
	if justifyQCBlock == nil {
		return nil, fmt.Errorf("justifyQC block (hash, height)not exist, (%v, %v)", justifyQC.Hash, justifyQCHeight)
	}

	if view.Round.Cmp(e.lastVoteRound) <= 0 {
		return nil, fmt.Errorf("proposalRound <= lastVoteRound, (%v, %v)", view.Round, e.lastVoteRound)
	}
	if justifyQCRound.Cmp(e.lockQCRound) < 0 {
		return nil, fmt.Errorf("justifyQCRound < lockQCRound, (%v, %v)", justifyQCRound, e.lockQCRound)
	}

	vote := &Vote{
		Epoch:      e.epoch,
		Hash:       hash,
		Proposer:   proposer,
		View:       view,
		ParentHash: justifyQC.Hash,
		ParentView: justifyQC.View,
	}

	if e.curHeight.Cmp(e.chain2Height()) >= 0 {
		qcParentHash := justifyQCBlock.ParentHash()
		qcParentHeight := new(big.Int).Sub(justifyQCHeight, common.Big1)
		qcParentBlock := e.blkTree.GetBlockAndCheckHeight(qcParentHash, qcParentHeight)
		if qcParentBlock == nil {
			return nil, fmt.Errorf("justifyQC parent (hash, height) not exist, (%v, %v)", qcParentHash, qcParentHeight)
		}
		qcParentSalt, _, err := extraProposal(qcParentBlock)
		if err != nil {
			return nil, err
		}
		qcParentView := &hotstuff.View{
			Round:  qcParentSalt.Round,
			Height: qcParentHeight,
		}
		vote.GrandHash = qcParentHash
		vote.GrandView = qcParentView
	}

	if e.curHeight.Cmp(e.chain3Height()) >= 0 {
		qcGrandHash := vote.GrandHash
		qcGrandHeight := new(big.Int).Sub(vote.GrandView.Height, common.Big1)
		qcGrandBlock := e.blkTree.GetBlockAndCheckHeight(qcGrandHash, qcGrandHeight)
		if qcGrandBlock == nil {
			return nil, fmt.Errorf("justifyQC grand-pa (hash, height) not exist, (%v, %v)", qcGrandHash, qcGrandHeight)
		}
		qcGrandSalt, _, err := extraProposal(qcGrandBlock)
		if err != nil {
			return nil, err
		}
		qcGrandView := &hotstuff.View{
			Round:  qcGrandSalt.Round,
			Height: qcGrandHeight,
		}
		vote.GreatGrandHash = qcGrandHash
		vote.GreatGrandView = qcGrandView
	}

	return vote, nil
}

func (e *EventDrivenEngine) validateVote(vote *Vote) error {
	// validate current proposal
	if vote.View == nil || vote.Hash == utils.EmptyHash {
		return errInvalidVote
	}
	if err := e.validateSingleChain(vote.Hash, vote.View, false, utils.EmptyHash); err != nil {
		return err
	}

	// validate parent block
	if vote.ParentHash == utils.EmptyHash || vote.ParentView == nil {
		return errInvalidVote
	}
	if new(big.Int).Add(vote.ParentView.Height, common.Big1).Cmp(vote.View.Height) != 0 {
		return errInvalidVote
	}
	if err := e.validateSingleChain(vote.ParentHash, vote.ParentView, true, vote.Hash); err != nil {
		return err
	}

	// validate grand block
	if e.curHeight.Cmp(e.chain2Height()) >= 0 {
		if vote.GrandHash == utils.EmptyHash || vote.GrandView == nil {
			return errInvalidVote
		}
		if new(big.Int).Add(vote.GrandView.Height, common.Big1).Cmp(vote.ParentView.Height) != 0 {
			return errInvalidVote
		}
		if err := e.validateSingleChain(vote.GrandHash, vote.GrandView, true, vote.ParentHash); err != nil {
			return err
		}
	}

	// validate great-grand block
	if e.curHeight.Cmp(e.chain3Height()) >= 0 {
		if vote.GreatGrandHash == utils.EmptyHash || vote.GreatGrandView == nil {
			return errInvalidVote
		}
		if new(big.Int).Add(vote.GreatGrandView.Height, common.Big1).Cmp(vote.GrandView.Height) != 0 {
			return errInvalidVote
		}
		if err := e.validateSingleChain(vote.GreatGrandHash, vote.GreatGrandView, true, vote.GrandHash); err != nil {
			return err
		}
	}
	return nil
}

func (e *EventDrivenEngine) validateSingleChain(hash common.Hash, view *hotstuff.View, checkParent bool, childParent common.Hash) error {
	block := e.blkTree.GetBlockAndCheckHeight(hash, view.Height)
	if block == nil {
		return fmt.Errorf("proposal %v not exist", hash)
	}
	salt, _, err := extraProposal(block)
	if err != nil {
		return err
	}
	if salt.Round.Cmp(view.Round) != 0 {
		return fmt.Errorf("vote proposal round expect %v, got %v", salt.Round, view.Round)
	}
	if checkParent && childParent != block.Hash() {
		return fmt.Errorf("vote proposal parent hash expect %v, got %v", childParent, block.Hash())
	}
	return nil
}
