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
	"math/big"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
)

// SafetyRules contains 3 group variables which used to judge that if the new proposal can be voted or committed.
// lockQC denotes an 2-chain header's qc, and the proposal which related by it can be committed in the next round.
// lastVote recorded the latest vote and round which used to ensure that there are only one vote in single round.
// lastCommittedQC represent an 3-chain header which have already been committed into pendingBlockTree/chain.
type SafetyRules struct {
	lockQC      *hotstuff.QuorumCert
	lockQCRound *big.Int

	lastVoteMsg   *Vote
	lastVoteRound *big.Int

	lastCommittedQC    *hotstuff.QuorumCert
	lastCommittedRound *big.Int
}

func NewSafetyRules() *SafetyRules {
	return nil
}

// UpdateLockQC update the latest quorum certificate after voteRule judgement succeed.
func (sr *SafetyRules) UpdateLockQC(qc *hotstuff.QuorumCert, round *big.Int) {
	sr.lockQC = qc
	sr.lockQCRound = round
}

// VoteRule validator should check vote in consensus round:
// first, the proposal should be exist in the `PendingBlockTree`
// second, the proposal round should be greater than `lastVoteRound`
// third, the proposal's justify qc round should NOT be smaller than `lockQCRound`
// we should ensure that only one vote in different round with first two items,
// and the last item used to make sure that there were `2F + 1` votes have been locked last in 3-chain round,
// and the proposal of that round should be current proposal's grand pa or justifyQC's parent.
func (sr *SafetyRules) VoteRule(proposalRound *big.Int, justifyQC *hotstuff.QuorumCert) bool {
	return false
}

// CommitRule validator should find out the parent block and grand pa block of the committed block by parent hash,
// and their height should be decreased by one.
func (sr *SafetyRules) CommitRule(block *types.Block, propoosalRound *big.Int, justifyQC *hotstuff.QuorumCert) bool {
	return false
}
