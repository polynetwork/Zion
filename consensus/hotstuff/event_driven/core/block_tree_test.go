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
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

var (
	testRoot = newBlockTreeTestBlock(common.EmptyHash, 0)
)

func newBlockTreeTestBlock(parentHash common.Hash, height int64) *types.Block {
	header := &types.Header{
		ParentHash: parentHash,
		Number:     big.NewInt(height),
		Time:       uint64(time.Now().UnixNano()), // `time` and `gasLimit` used to make hash different
		GasLimit:   rand.Uint64(),
	}
	return types.NewBlock(header, nil, nil, nil, nil)
}

func newTestBlockTree(pureSize int) *BlockTree {
	tr, _ := NewBlockTree(testRoot, 0, pureSize)
	return tr
}

func TestBlockTree(t *testing.T) {
	tr := newTestBlockTree(0)

	block1 := newBlockTreeTestBlock(testRoot.Hash(), 1)
	block2 := newBlockTreeTestBlock(block1.Hash(), 2)
	block3 := newBlockTreeTestBlock(block2.Hash(), 3)

	// test GetBlockByHash
	_ = tr.Add(block1, 1)
	_ = tr.Add(block2, 2)
	_ = tr.Add(block3, 3)

	assert.NotNil(t, tr.GetBlockByHash(block1.Hash()))
	assert.NotNil(t, tr.GetBlockByHash(block2.Hash()))
	assert.NotNil(t, tr.GetBlockByHash(block3.Hash()))

	t.Log(tr.Details())

	// test branch
	branch3 := tr.Branch(block3.Hash())
	assert.Equal(t, 3, len(branch3))
	for i, node := range branch3 {
		v := node.GetBlock()
		t.Logf("branch3: block%d hash %v, height %d", i, v.Hash(), v.NumberU64())
	}

	branch2 := tr.Branch(block2.Hash())
	assert.Equal(t, 2, len(branch2))
	for i, node := range branch2 {
		v := node.GetBlock()
		t.Logf("branch2: block%d hash %v, height %d", i, v.Hash(), v.NumberU64())
	}

	branch1 := tr.Branch(block1.Hash())
	assert.Equal(t, 1, len(branch1))
	for i, node := range branch1 {
		v := node.GetBlock()
		t.Logf("branch1: block%d hash %v, height %d", i, v.Hash(), v.NumberU64())
	}

	// test pure
	assert.Equal(t, 1, len(tr.Prune(block1.Hash())))
	assert.Equal(t, 1, len(tr.Prune(block2.Hash())))
	assert.Equal(t, 1, len(tr.Prune(block3.Hash())))
}

func TestBlockTreeBranch(t *testing.T) {
	t.Log("fork from block (round: 4, height: 3)")
	tr := newTestBlockTree(0)

	block1 := newBlockTreeTestBlock(testRoot.Hash(), 1)
	block2 := newBlockTreeTestBlock(block1.Hash(), 2)
	block30 := newBlockTreeTestBlock(block2.Hash(), 3)
	block31 := newBlockTreeTestBlock(block2.Hash(), 3)
	block32 := newBlockTreeTestBlock(block2.Hash(), 3)
	block33 := newBlockTreeTestBlock(block2.Hash(), 3)

	// test GetBlockByHash
	_ = tr.Add(block1, 1)
	_ = tr.Add(block2, 2)
	_ = tr.Add(block30, 3)
	_ = tr.Add(block31, 4)
	_ = tr.Add(block32, 5)
	_ = tr.Add(block33, 6)

	t.Log("before pure", tr.Details())

	br := tr.Branch(block31.Hash())
	t.Log("branch detail", br.Details())

	tr.Prune(block31.Hash())
	t.Log("after pure", tr.Details())
}

func TestBlockTreeFork(t *testing.T) {
	{
		t.Log("fork from block3")
		tr := newTestBlockTree(0)

		block1 := newBlockTreeTestBlock(testRoot.Hash(), 1)
		block2 := newBlockTreeTestBlock(block1.Hash(), 2)
		block30 := newBlockTreeTestBlock(block2.Hash(), 3)
		block31 := newBlockTreeTestBlock(block2.Hash(), 3)
		block32 := newBlockTreeTestBlock(block2.Hash(), 3)
		block33 := newBlockTreeTestBlock(block2.Hash(), 3)

		// test GetBlockByHash
		_ = tr.Add(block1, 1)
		_ = tr.Add(block2, 2)
		_ = tr.Add(block30, 3)
		_ = tr.Add(block31, 4)
		_ = tr.Add(block32, 5)
		_ = tr.Add(block33, 6)

		assert.Equal(t, 6, len(tr.Prune(block31.Hash())))
		assert.Equal(t, tr.root.block.Hash(), block31.Hash())
		t.Log(tr.Details())
		t.Log("======================================")
	}

	{
		t.Log("fork from block1")
		tr := newTestBlockTree(0)

		block1 := newBlockTreeTestBlock(testRoot.Hash(), 1)
		block20 := newBlockTreeTestBlock(block1.Hash(), 2)
		block21 := newBlockTreeTestBlock(block1.Hash(), 2)
		block30 := newBlockTreeTestBlock(block21.Hash(), 3)
		block32 := newBlockTreeTestBlock(block21.Hash(), 3)
		block33 := newBlockTreeTestBlock(block21.Hash(), 3)

		// test GetBlockByHash
		_ = tr.Add(block1, 1)
		_ = tr.Add(block20, 2)
		_ = tr.Add(block21, 3)
		_ = tr.Add(block30, 4)
		_ = tr.Add(block32, 5)
		_ = tr.Add(block33, 6)

		assert.Equal(t, 6, len(tr.Prune(block33.Hash())))
		assert.Equal(t, tr.root.block.Hash(), block33.Hash())
		t.Log(tr.Details())
		t.Log("======================================")
	}
}
