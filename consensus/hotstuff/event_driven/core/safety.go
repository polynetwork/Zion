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
	"github.com/ethereum/go-ethereum/core/types"
)

// IncreaseLastVoteRound commit not to vote in rounds lower than target
func (e *EventDrivenEngine) IncreaseLastVoteRound(rd *big.Int) {
	if e.lastVoteRound.Cmp(rd) < 0 {
		e.lastVoteRound = rd
	}
}

// UpdateLockQC update the latest quorum certificate after voteRule judgement succeed.
func (e *EventDrivenEngine) UpdateLockQCRound(round *big.Int) {
	e.lockQCRound = round
}

// VoteRule validator should check vote in consensus round:
// first, the proposal should be exist in the `PendingBlockTree`
// second, the proposal round should be greater than `lastVoteRound`
// third, the proposal's justify qc round should NOT be smaller than `lockQCRound`
// we should ensure that only one vote in different round with first two items,
// and the last item used to make sure that there were `2F + 1` votes have been locked last in 3-chain round,
// and the proposal of that round should be current proposal's grand pa or justifyQC's parent.
func (e *EventDrivenEngine) VoteRule(proposalRound, proposalJustifyQCRound *big.Int) bool {
	if proposalRound == nil || proposalJustifyQCRound == nil {
		e.logger.Error("[safety voteRule]", "some params invalid", "nil")
		return false
	}

	if proposalRound.Cmp(e.lastVoteRound) <= 0 {
		e.logger.Error("[safety voteRule]", "this round already voted, proposalRound ", proposalRound, "lastVoteRound", e.lastVoteRound)
		return false
	}
	if proposalJustifyQCRound.Cmp(e.lockQCRound) < 0 {
		e.logger.Error("[safety voteRule]", "proposal parent qc round should be greater, justifyQCRound ", proposalJustifyQCRound, "lockQCRound", e.lockQCRound)
		return false
	}

	return true
}

func (e *EventDrivenEngine) MakeVote(proposal *types.Block) (*Vote, error) {
	justifyQC, proposalRound, err := extraProposal(proposal)
	if err != nil {
		return nil, err
	}
	justifyQCRound := justifyQC.View.Round
	justifyQCHeight := justifyQC.View.Height

	justifyQCBlock := e.blkTree.GetBlockAndCheckHeight(justifyQC.Hash, justifyQCHeight)
	if justifyQCBlock == nil {
		return nil, fmt.Errorf("justifyQC block (hash, height)not exist, (%v, %v)", justifyQC.Hash, justifyQCHeight)
	}

	if proposalRound.Cmp(e.lastVoteRound) <= 0 {
		//sr.logger.Error("[safety voteRule]", "this round already voted, proposalRound ", proposalRound, "lastVoteRound", sr.lastVoteRound)
		return nil, fmt.Errorf("proposalRound <= lastVoteRound, (%v, %v)", proposalRound, e.lastVoteRound)
	}
	if justifyQCRound.Cmp(e.lockQCRound) < 0 {
		//sr.logger.Error("[safety voteRule]", "proposal parent qc round should be greater, justifyQCRound ", proposalJustifyQCRound, "lockQCRound", sr.lockQCRound)
		return nil, fmt.Errorf("justifyQCRound < lockQCRound, (%v, %v)", justifyQCRound, e.lockQCRound)
	}

	// qc.parent.round + 1 == qc.round
	qcParentHash := justifyQCBlock.ParentHash()
	qcParentHeight := new(big.Int).Sub(justifyQCHeight, common.Big1)
	qcParentBlock := e.blkTree.GetBlockAndCheckHeight(qcParentHash, qcParentHeight)
	_, qcParentRound, err := extraProposal(qcParentBlock)
	if err != nil {
		return nil, err
	}
	if qcParentBlock == nil {
		return nil, fmt.Errorf("justifyQC parent (hash, height) not exist, (%v, %v)", qcParentHash, qcParentHeight)
	}

	qcGrandHash := qcParentBlock.ParentHash()
	qcGrandHeight := new(big.Int).Sub(qcParentHeight, common.Big1)
	qcGrandBlock := e.blkTree.GetBlockAndCheckHeight(qcGrandHash, qcGrandHeight)
	if qcGrandBlock == nil {
		return nil, fmt.Errorf("justifyQC grand-pa (hash, height) not exist, (%v, %v)", qcGrandHash, qcGrandHeight)
	}

	vote := &Vote{
		Epoch:       e.epoch,
		Hash:        proposal.Hash(),
		Round:       proposalRound,
		ParentHash:  justifyQC.Hash,
		ParentRound: justifyQCRound,
		GrandHash:   qcParentHash,
		GrandRound:  qcParentRound,
	}
	return vote, nil
}

// CommitRule validator should find out the parent block and grand pa block of the committed block by parent hash,
// and their height should be decreased by one.
/*
// find the committed id in case a qc is formed in the vote round
if (qc.parent round + 1 == qc.round) ∧ (qc.round + 1 == vote round) then
return qc.parent id
else
return nil
*/
// todo: useless
func (e *EventDrivenEngine) CommitRule(proposalJustifyQCRound, proposalJustifyQCParentRound *big.Int) bool {
	return false
}
