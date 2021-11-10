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
	"os"
	"testing"

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
	db := rawdb.NewMemoryDatabase()
	testStateDB, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	ref := generateContractRef()
	testEmptyCtx = native.NewNativeContract(testStateDB, ref)

	InitABI()
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

func TestStoreProxy(t *testing.T) {
	resetTestContext()
	s := testEmptyCtx

	var testcases = []struct {
		ChainID uint64
		Proxy   common.Address
	}{
		{
			ChainID: 12,
			Proxy:   common.HexToAddress("0xcbc84f846c4afabd5a8adcef92b40c1c4448f31a"),
		},
		{
			ChainID: 0,
			Proxy:   common.HexToAddress("0x846c4afabd5a8adcef92b40c1c4448f31a"),
		},
		{
			ChainID: 12,
			Proxy:   common.EmptyAddress,
		},
		{
			ChainID: 0,
			Proxy:   common.EmptyAddress,
		},
	}

	for _, v := range testcases {
		storeProxy(s, v.ChainID, v.Proxy[:])
		blob, err := getProxy(s, v.ChainID)
		assert.NoError(t, err)
		got := common.BytesToAddress(blob)
		assert.Equal(t, v.Proxy, got)
	}

	// test store nil
	storeProxy(s, 12, nil)
	blob, err := getProxy(s, 12)
	assert.NoError(t, err)
	assert.Nil(t, blob)
}

func TestStoreAsset(t *testing.T) {
	resetTestContext()
	s := testEmptyCtx

	var testcases = []struct {
		FromAsset     common.Address
		TargetChainID uint64
		ToAssetHash   []byte
	}{
		{
			FromAsset:     common.HexToAddress("0xcbc84f846c4afabd5a8adcef92b40c1c4448f31a"),
			TargetChainID: 12,
			ToAssetHash:   []byte{'1', 'a', '3'},
		},
		{
			FromAsset:     common.EmptyAddress,
			TargetChainID: 0,
			ToAssetHash:   nil,
		},
	}

	for _, v := range testcases {
		storeAsset(s, v.FromAsset, v.TargetChainID, v.ToAssetHash)
		got, err := getAsset(s, v.FromAsset, v.TargetChainID)
		assert.NoError(t, err)
		assert.Equal(t, v.ToAssetHash, got)
	}
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
