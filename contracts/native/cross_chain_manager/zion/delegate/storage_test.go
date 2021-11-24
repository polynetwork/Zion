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

package delegate

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

	os.Exit(m.Run())
}

func generateContractRef() *native.ContractRef {
	return native.NewContractRef(testStateDB, testCaller, testCaller, big.NewInt(testBlockNum), testTxHash, testSupplyGas, nil)
}

func resetTestContext() {
	db := rawdb.NewMemoryDatabase()
	testStateDB, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	ref := generateContractRef()
	ref.PushContext(&native.Context{
		Caller:          testCaller,
		ContractAddress: this,
		Payload:         nil,
	})
	testEmptyCtx = native.NewNativeContract(testStateDB, ref)
}

func TestStoreAllowance(t *testing.T) {
	resetTestContext()
	s := testEmptyCtx
	owner := common.HexToAddress("0x13")
	spender := this
	expect := big.NewInt(899)

	setAllowance(s, owner, spender, expect)
	got := getAllowance(s, owner, spender)
	assert.Equal(t, expect, got)
}
