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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockPool struct {
	// todo: qc cache
	tree   *BlockTree
	highQC *hotstuff.QuorumCert // the highest qc, 从genesis 0开始
}

func (tr *BlockPool) GetHighQC() *hotstuff.QuorumCert {
	return tr.highQC
}

func (tr *BlockPool) GetBlockAndCheckHeight(hash common.Hash, height *big.Int) *types.Block {
	parentBlock := tr.GetBlockByHash(hash)
	if parentBlock == nil {
		return nil
	}
	if parentBlock.Number().Cmp(height) != 0 {
		return nil
	}
	return parentBlock
}

func (tr *BlockPool) GetBlockByHash(hash common.Hash) *types.Block {
	return tr.tree.GetBlockByHash(hash)
}

// Insert insert new block into pending block tree, calculate and return the highestQC
func (tr *BlockPool) Insert(block *types.Block, round *big.Int) error {
	return tr.tree.Add(block, round.Uint64())
}

func (tr *BlockPool) UpdateHighQC(qc *hotstuff.QuorumCert) {
	if qc == nil || qc.View == nil {
		return
	}
	if tr.highQC == nil || tr.highQC.View == nil {
		tr.highQC = qc
	} else if tr.highQC.View.Round.Cmp(qc.View.Round) < 0 {
		tr.highQC = qc
	}
}

// GetCommitBlock commit the block into ledger and pure the `pendingBlockTree`
func (tr *BlockPool) GetCommitBlock(qc common.Hash) *types.Block {
	qcBlock := tr.tree.GetBlockByHash(qc)
	if qcBlock == nil {
		return nil
	}
	qcParentHash := qcBlock.ParentHash()
	return tr.tree.GetBlockByHash(qcParentHash)
}

// Pure delete useless blocks
func (tr *BlockPool) Pure(committedBlock common.Hash) {
	tr.tree.Prune(committedBlock)
}