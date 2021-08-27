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

package backup
//
//import (
//	"fmt"
//	"math/big"
//
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/core/types"
//)
//
//type TreeNode struct {
//	block    *types.Block
//	round    *big.Int
//	children []common.Hash
//}
//
//func (n *TreeNode) HasChild(childHash common.Hash) bool {
//	if n == nil {
//		return false
//	}
//	if n.children == nil || len(n.children) == 0 {
//		return false
//	}
//	for _, v := range n.children {
//		if v == childHash {
//			return true
//		}
//	}
//	return false
//}
//
//type PendingBlockTree struct {
//	root     *TreeNode
//	nodes    map[common.Hash]*TreeNode
//	capacity int
//}
//
//func NewPendingBlockTree(root *TreeNode, capacity int) *PendingBlockTree {
//	tree := &PendingBlockTree{
//		root:     root,
//		nodes:    make(map[common.Hash]*TreeNode),
//		capacity: capacity,
//	}
//	return tree
//}
//
//func (tr *PendingBlockTree) Add(block *types.Block, round *big.Int) error {
//	blockHash := block.Hash()
//	parentHash := block.ParentHash()
//	parentNode, ok := tr.nodes[parentHash]
//	if !ok {
//		return fmt.Errorf("tree node %s parent node %s not exist", blockHash, parentHash)
//	}
//
//	if parentNode.children == nil || len(parentNode.children) == 0 {
//		parentNode.children = []common.Hash{blockHash}
//	} else if !parentNode.HasChild(blockHash) {
//		parentNode.children = append(parentNode.children, blockHash)
//	} else {
//		return nil
//		// return fmt.Errorf("tree node %v already exist", blockHash)
//	}
//
//	node := &TreeNode{
//		block: block,
//		round: new(big.Int).Set(round),
//	}
//
//	tr.nodes[blockHash] = node
//	return nil
//}
//
//func (tr *PendingBlockTree) GetBlockByHash(hash common.Hash) *types.Block {
//	node, ok := tr.nodes[hash]
//	if !ok {
//		return nil
//	}
//	return node.block
//}
//
//// Pure delete the forked branches in this tree and reset an new root which always is an locked block.
//func (tr *PendingBlockTree) Pure(hash common.Hash, round *big.Int) {
//	node, ok := tr.nodes[hash]
//	if !ok {
//		return
//	}
//
//
//}
