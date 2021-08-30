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
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/prque"
	"github.com/ethereum/go-ethereum/core/types"
)

//Node save one block and its children
type Node struct {
	block    *types.Block
	round    uint64
	children []*Node // the blockHash with children's block
}

//GetBlock get block
func (n *Node) GetBlock() *types.Block {
	return n.block
}

//GetChildren get children
func (n *Node) GetChildren() []*Node {
	return n.children
}

func generateNode(block *types.Block, round uint64) *Node {
	return &Node{
		block:    block,
		round:    round,
		children: make([]*Node, 0),
	}
}

//BlockTree maintains a block tree of parent and children links and this struct is not thread safety.
type BlockTree struct {
	nodes         map[common.Hash]*Node // store block and its' children blockHash
	roundTable    map[uint64][]*Node    // store nodes in the same round
	root          *Node                 // the latest block is committed to the chain
	prunedBlocks  []common.Hash         // caches the block hash that will be deleted
	maxPrunedSize int                   // the maximum number of cached blocks that will be deleted
}

//NewBlockTree init a block tree with root, rootQC and maxPrunedSize
func NewBlockTree(rootBlock *types.Block, rootRound uint64, maxPrunedSize int) (*BlockTree, error) {
	rootNode := generateNode(rootBlock, rootRound)
	blockTree := &BlockTree{
		nodes:         make(map[common.Hash]*Node, 10),
		root:          rootNode,
		prunedBlocks:  make([]common.Hash, 0, maxPrunedSize),
		maxPrunedSize: maxPrunedSize,
		roundTable:    make(map[uint64][]*Node),
	}
	blockTree.nodes[rootBlock.Hash()] = rootNode
	blockTree.roundTable[rootRound] = append(blockTree.roundTable[rootRound], rootNode)
	if err := blockTree.Add(rootBlock, rootRound); err != nil {
		return nil, err
	}
	return blockTree, nil
}

//Add insert block to tree
func (bt *BlockTree) Add(block *types.Block, round uint64) error {
	if block == nil || block.Hash() == common.EmptyHash {
		return fmt.Errorf("block is invalid")
	}
	if round > 0 && block.ParentHash() == common.EmptyHash {
		return fmt.Errorf("block parent hash is empty")
	}

	var (
		ok     bool
		hash   = block.Hash()
		parent *Node
	)

	if _, ok = bt.nodes[hash]; ok {
		return nil
	}

	if round > 0 {
		parentHash := block.ParentHash()
		if parent, ok = bt.nodes[parentHash]; !ok {
			return fmt.Errorf("block's parent not exist")
		}
	}

	node := generateNode(block, round)
	bt.nodes[hash] = node
	bt.roundTable[round] = append(bt.roundTable[round], node)
	if parent != nil {
		parent.children = append(parent.children, node)
	}

	return nil
}

//GetRootBlock get root block from tree
func (bt *BlockTree) GetRootBlock() *types.Block {
	return bt.root.GetBlock()
}

//GetBlockByID get block by block hash
func (bt *BlockTree) GetBlockByHash(hash common.Hash) *types.Block {
	if node, ok := bt.nodes[hash]; ok {
		return node.GetBlock()
	}
	return nil
}

// GetBlocksByRound
func (bt *BlockTree) GetBlocksByRound(round uint64) []*types.Block {
	list := make([]*types.Block, 0)
	for _, v := range bt.roundTable[round] {
		list = append(list, v.block)
	}
	return list
}

//Branch get branch from root to input block
func (bt *BlockTree) Branch(hash common.Hash) []*types.Block {
	cur, ok := bt.nodes[hash]
	if !ok {
		return nil
	}
	var branch []*types.Block

	for cur.round > bt.root.round {
		curBlock := cur.GetBlock()
		parentHash := curBlock.ParentHash()
		branch = append(branch, curBlock)
		if cur = bt.nodes[parentHash]; cur == nil {
			break
		}
	}

	if cur == nil || cur.block.Hash() != bt.root.block.Hash() {
		return nil
	}
	for i, j := 0, len(branch)-1; i < j; i, j = i+1, j-1 {
		branch[i], branch[j] = branch[j], branch[i]
	}
	return branch
}

//Prune prune block and update root
func (bt *BlockTree) Prune(newRootHash common.Hash) []common.Hash {
	toPruned := bt.findBlockToPrune(newRootHash)
	if toPruned == nil {
		return nil
	}

	newRootNode, ok := bt.nodes[newRootHash]
	if !ok {
		return nil
	}
	bt.root = newRootNode
	bt.prunedBlocks = append(bt.prunedBlocks, toPruned[0:]...)

	var pruned []common.Hash
	if len(bt.prunedBlocks) > bt.maxPrunedSize {
		num := len(bt.prunedBlocks) - bt.maxPrunedSize
		for i := 0; i < num; i++ {
			bt.cleanBlock(bt.prunedBlocks[i])
			pruned = append(pruned, bt.prunedBlocks[i])
		}
		bt.prunedBlocks = bt.prunedBlocks[num:]
	}
	return pruned
}

const defaultPriority = 1

//findBlockToPrune get blocks to prune by the newRootID
func (bt *BlockTree) findBlockToPrune(hash common.Hash) []common.Hash {
	if hash == common.EmptyHash || hash == bt.root.block.Hash() {
		return nil
	}

	var (
		toPruned      = make([]common.Hash, 0)
		toPrunedQueue = prque.New(nil)
	)
	toPrunedQueue.Push(bt.root, defaultPriority)
	for !toPrunedQueue.Empty() {
		data, _ := toPrunedQueue.Pop()
		curID := data.(*Node)
		curNode := bt.nodes[curID.block.Hash()]
		for _, child := range curNode.GetChildren() {
			if child.block.Hash() == hash {
				continue
			}
			toPrunedQueue.Push(child, defaultPriority)
		}
		toPruned = append(toPruned, curID.block.Hash())
	}
	return toPruned
}

//cleanBlock remove block from tree
func (bt *BlockTree) cleanBlock(hash common.Hash) {
	blk := bt.nodes[hash]
	if blk != nil {
		delete(bt.nodes, hash)
		delete(bt.roundTable, blk.round)
	}
}

func (bt *BlockTree) Details() string {
	content := bytes.NewBufferString(fmt.Sprintf("block tree size: %d\n", len(bt.nodes)))
	for _, nodes := range bt.roundTable {
		for _, node := range nodes {
			content.WriteString(fmt.Sprintf("hash: %v, height: %d, round:%d\n", node.block.Hash(), node.block.NumberU64(), node.round))
		}
	}
	return content.String()
}
