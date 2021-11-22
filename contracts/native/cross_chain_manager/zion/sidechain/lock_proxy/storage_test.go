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

	InitAllocProxy()

	os.Exit(m.Run())
}

func TestStoreCrossTxIndex(t *testing.T) {
	resetTestContext()
	s := testEmptyCtx

	lastCrossTxIndex := getCrossTxIndex(s)
	nextCrossTxIndex := lastCrossTxIndex + 1
	storeCrossTxIndex(s, nextCrossTxIndex)
	curCrossTxIndex := getCrossTxIndex(s)

	assert.Equal(t, nextCrossTxIndex, curCrossTxIndex)
}

func TestStoreCrossTxContent(t *testing.T) {
	resetTestContext()
	s := testEmptyCtx

	expect := &CrossTx{
		ToChainId:   12,
		FromAddress: common.HexToAddress("123"),
		ToAddress:   common.HexToAddress("aba3"),
		Amount:      big.NewInt(120),
		Index:       1,
	}

	assert.NoError(t, storeCrossTxContent(s, expect))
	got, err := getCrossTxContent(s, expect.Hash())
	assert.NoError(t, err)

	assert.Equal(t, expect, got)
}

func TestStoreCrossTxProof(t *testing.T) {
	resetTestContext()
	s := testEmptyCtx

	tx := &CrossTx{
		ToChainId:   12,
		FromAddress: common.HexToAddress("123"),
		ToAddress:   common.HexToAddress("aba3"),
		Amount:      big.NewInt(120),
		Index:       1,
	}
	expect := tx.Proof()
	storeCrossTxProof(s, tx.Index, expect)
	got, err := getCrossTxProof(s, tx.Index)
	assert.NoError(t, err)

	assert.Equal(t, expect, got)
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
