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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/message_set"
	"github.com/ethereum/go-ethereum/core/types"
)

// todo: if actually need it
type BlockTree struct {
	tree *PendingBlockTree
	highQC *hotstuff.QuorumCert		// the highest qc
	pendingVotes *message_set.MessageSet
}

// Insert insert new block into pending block tree, calculate and return the highestQC
func (tr *BlockTree) Insert(block *types.Block) *hotstuff.QuorumCert {
	return nil
}

// ProcessVote caching participants votes and drive `paceMaker` into next round if the
// vote message number arrived the quorum size.
func (tr *BlockTree) ProcessVote(vote *Vote) {

}

// ProcessCommit commit the block into ledger and pure the `pendingBlockTree`
func (tr *BlockTree) ProcessCommit(hash common.Hash) {

}