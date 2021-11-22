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

package lock_proxy

import (
	"math/big"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/stretchr/testify/assert"
)

var (
	testStateDB  *state.StateDB
	testEmptyCtx *native.NativeContract

	testSupplyGas uint64 = 100000000000000000
	testBlockNum         = int64(12)
	testTxHash           = common.EmptyHash
	testCaller           = common.EmptyAddress
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())

	db := rawdb.NewMemoryDatabase()
	testStateDB, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	ref := generateContractRef()
	testEmptyCtx = native.NewNativeContract(testStateDB, ref)

	InitLockProxy()

	os.Exit(m.Run())
}

func generateContractRef() *native.ContractRef {
	return native.NewContractRef(testStateDB, testCaller, testCaller, big.NewInt(testBlockNum), testTxHash, testSupplyGas, nil)
}
func resetTestContext() {
	db := rawdb.NewMemoryDatabase()
	testStateDB, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	ref := generateContractRef()
	testEmptyCtx = native.NewNativeContract(testStateDB, ref)
}

func TestStoreTxIndex(t *testing.T) {
	resetTestContext()
	s := testEmptyCtx

	var testcases = []struct {
		Index  *big.Int
		Expect *big.Int
	}{
		{
			Index:  big.NewInt(0),
			Expect: common.Big0,
		},
		{
			Index:  nil,
			Expect: common.Big0,
		},
		{
			Index:  common.Big0,
			Expect: common.Big0,
		},
		{
			Index:  big.NewInt(1),
			Expect: common.Big1,
		},
	}

	for _, v := range testcases {
		storeTxIndex(s, v.Index)
		got := getTxIndex(s)
		assert.Equal(t, v.Expect.Uint64(), got.Uint64())
	}
}

func TestStoreAmount(t *testing.T) {
	resetTestContext()
	s := testEmptyCtx
	chainID := uint64(12)

	var testcases = []struct {
		Amount  uint64
		Add bool
		Expect uint64
	}{
		{
			Amount: 4,
			Add: true,
			Expect: 4,
		},
		{
			Amount: 3,
			Add: true,
			Expect: 7,
		},
		{
			Amount: 2,
			Add: true,
			Expect: 9,
		},
		{
			Amount: 1,
			Add: true,
			Expect: 10,
		},
		{
			Amount: 10,
			Add: false,
			Expect: 0,
		},
	}

	for _, v := range testcases {
		if v.Add {
			addTotalAmount(s, chainID, new(big.Int).SetUint64(v.Amount))
		} else {
			subTotalAmount(s, chainID, new(big.Int).SetUint64(v.Amount))
		}
		data := getTotalAmount(s, chainID)
		assert.Equal(t, v.Expect, data.Uint64())
	}
}