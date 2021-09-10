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
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
)

var testEpoch uint64 = 0

func makeAddress(i int) common.Address {
	num := new(big.Int).SetUint64(uint64(i))
	return common.BytesToAddress(num.Bytes())
}

func makeHash(i int) common.Hash {
	num := new(big.Int).SetUint64(uint64(i))
	return common.BytesToHash(num.Bytes())
}

func makeView(h, r uint64) *hotstuff.View {
	return &hotstuff.View{
		Height: new(big.Int).SetUint64(h),
		Round:  new(big.Int).SetUint64(r),
	}
}

func makeGenesisBlock(valset hotstuff.ValidatorSet) *types.Block {
	vs := valset.Copy()
	vs.CalcProposerByIndex(0)
	proposer := vs.GetProposer().Address()
	header := &types.Header{
		Difficulty: big.NewInt(0),
		Number:     big.NewInt(0),
		GasLimit:   0,
		GasUsed:    0,
		Time:       uint64(time.Now().Unix()),
		Coinbase:   proposer,
	}
	header.Extra, _ = generateExtra(header, valset, testEpoch, big.NewInt(0))
	block := &types.Block{}
	return block.WithSeal(header)
}

func makeBlockWithParentHash(valset hotstuff.ValidatorSet, view *hotstuff.View, parentHash common.Hash) *types.Block {
	vs := valset.Copy()
	vs.CalcProposerByIndex(view.Round.Uint64())
	proposer := vs.GetProposer().Address()
	header := &types.Header{
		Difficulty: big.NewInt(0),
		Number:     view.Height,
		GasLimit:   0,
		GasUsed:    0,
		Time:       uint64(time.Now().Unix()),
		Coinbase:   proposer,
		ParentHash: parentHash,
	}
	header.Extra, _ = generateExtra(header, valset, testEpoch, view.Round)
	block := &types.Block{}
	return block.WithSeal(header)
}

func makeContinueBlocks(valset hotstuff.ValidatorSet, initBlock *types.Block, n int) []*types.Block {
	list := make([]*types.Block, n)
	list[0] = initBlock
	vs := valset.Copy()
	if n <= 1 {
		return list
	}
	salt, _, _ := extraProposal(initBlock)
	view := &hotstuff.View{
		Round:  salt.Round,
		Height: initBlock.Number(),
	}

	for i := 1; i < n; i++ {
		view = &hotstuff.View{
			Round:  new(big.Int).Add(view.Round, common.Big1),
			Height: new(big.Int).Add(view.Height, common.Big1),
		}
		parentHash := list[i-1].Hash()
		list[i] = makeBlockWithParentHash(vs, view, parentHash)
	}
	return list
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	os.Exit(m.Run())
}
