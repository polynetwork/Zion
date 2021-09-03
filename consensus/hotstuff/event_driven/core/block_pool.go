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
	tree  *BlockTree
	qcMap map[common.Hash]*hotstuff.QuorumCert // caches the quorum certificates
}

func NewBlockPool(tr *BlockTree) *BlockPool {
	return &BlockPool{
		tree:  tr,
		qcMap: make(map[common.Hash]*hotstuff.QuorumCert),
	}
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

func (tr *BlockPool) GetQCByHash(hash common.Hash) *hotstuff.QuorumCert {
	return tr.qcMap[hash]
}

// AddBlock insert new block into pending block tree, calculate and return the highestQC
// allow to store the sealed and unsealed block with same hash.
func (tr *BlockPool) AddBlock(block *types.Block, round *big.Int) error {
	return tr.tree.Add(block, round.Uint64())
}

func (tr *BlockPool) AddQC(qc *hotstuff.QuorumCert) {
	if _, ok := tr.qcMap[qc.Hash]; !ok {
		tr.qcMap[qc.Hash] = qc
	}
}

// GetCommitBlock get highQC's grand-parent block which should be committed at current round
func (tr *BlockPool) GetCommitBlock(highQC, lockQC common.Hash) *types.Block {
	block := tr.GetBlockByHash(highQC)
	parent := tr.GetBlockAndCheckHeight(block.ParentHash(), bigSub1(block.Number()))
	if parent == nil {
		return nil
	}
	grand := tr.GetBlockAndCheckHeight(parent.ParentHash(), bigSub1(parent.Number()))
	if grand == nil {
		return nil
	}
	if grand.Hash() != lockQC {
		return nil
	}
	return tr.GetBlockByHash(lockQC)
}

// Pure delete useless blocks
func (tr *BlockPool) Pure(committedBlock common.Hash) {
	tr.tree.Prune(committedBlock)
}
