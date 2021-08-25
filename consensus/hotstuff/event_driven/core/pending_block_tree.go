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
	"github.com/ethereum/go-ethereum/common"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type TreeNode struct {
	block    *types.Block
	round    *big.Int
	clildren []string
}

func (n *TreeNode) HasChild(childHash string) bool {
	if n == nil {
		return false
	}
	if n.clildren == nil || len(n.clildren) == 0 {
		return false
	}
	for _, v := range n.clildren {
		if v == childHash {
			return true
		}
	}
	return false
}

type PendingBlockTree struct {
	root      *TreeNode
	nodeHashs map[string]*TreeNode
	capacity  int
}

// todo:
func NewPendingBlockTree() *PendingBlockTree {
	tree := &PendingBlockTree{}
	return tree
}

func (tr *PendingBlockTree) Add(block *types.Block, round *big.Int) error {
	blockHash := block.Hash().Hex()
	parentHash := block.ParentHash().Hex()
	parentNode, ok := tr.nodeHashs[parentHash]
	if !ok {
		return fmt.Errorf("tree node %s parent node %s not exist", blockHash, parentHash)
	}

	if parentNode.clildren == nil || len(parentNode.clildren) == 0 {
		parentNode.clildren = []string{blockHash}
	} else if !parentNode.HasChild(blockHash) {
		parentNode.clildren = append(parentNode.clildren, blockHash)
	} else {
		// todo: child already exist
		return nil
	}

	node := &TreeNode{
		block: block,
		round: new(big.Int).Set(round),
	}

	tr.nodeHashs[blockHash] = node
	return nil
}

// Branch retrieve the clean branch which has an 3-chain
func (tr *PendingBlockTree) Branch(block *types.Block) []*types.Block {
	return nil
}

func (tr *PendingBlockTree) GetBlockByHash(hash common.Hash) *types.Block {
	node, ok := tr.nodeHashs[hash.Hex()]
	if !ok {
		return nil
	}
	return node.block
}

// Pure delete the forked branches in this tree and reset an new root which always is an locked block.
func (tr *PendingBlockTree) Pure(block *types.Block, round *big.Int) (*types.Block, error) {
	return nil, nil
}
