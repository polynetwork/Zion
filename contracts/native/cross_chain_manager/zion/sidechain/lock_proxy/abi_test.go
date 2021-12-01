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

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	zutils "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/utils"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
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
	ref.PushContext(&native.Context{
		Caller:          this,
		ContractAddress: this,
		Payload:         nil,
	})
}

func TestABIMethodContractNameOutput(t *testing.T) {
	expect := &MethodContractNameOutput{Name: contractName}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodContractNameOutput)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)
}

func TestABIMethodBurnInput(t *testing.T) {
	expect := &MethodBurnInput{
		ToChainId: 3,
		ToAddress: common.HexToAddress("0x3"),
		Amount:    big.NewInt(145),
	}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodBurnInput)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)
}

func TestABIMethodMintInput(t *testing.T) {
	expect := &MethodMintInput{
		ArgsBs:           []byte{'a'},
		FromContractAddr: []byte{'x'},
		FromChainId:      12,
	}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodMintInput)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)

	// test mint packed
	id := crypto.Keccak256(utils.EncodePacked([]byte("mint"), []byte("(bytes,bytes,uint64)")))[:4]
	args := abi.Arguments{
		{Type: zutils.BytesTy, Name: "_argsBs"},
		{Type: zutils.BytesTy, Name: "_fromContractAddr"},
		{Type: zutils.Uint64Ty, Name: "_fromChainId"},
	}
	data, err := args.Pack(expect.ArgsBs, expect.FromContractAddr, expect.FromChainId)
	assert.NoError(t, err)

	packed := utils.EncodePacked(id, data)
	assert.Equal(t, payload, packed)

	// test name packed
	id = crypto.Keccak256(utils.EncodePacked([]byte("name"), []byte("()")))[:4]
	nameInput := new(MethodContractNameInput)
	namePayload, err := nameInput.Encode()
	assert.NoError(t, err)
	assert.Equal(t, namePayload, id)
}

func TestEmitBurn(t *testing.T) {
	resetTestContext()
	s := testEmptyCtx

	fromAsset := common.HexToAddress("0x2")
	fromAddr := common.HexToAddress("0x3")
	toChainID := uint64(3)
	toAsset := common.HexToAddress("0x5").Bytes()
	toAddr := common.HexToAddress("0x6").Bytes()
	amount := big.NewInt(13)

	err := emitBurnEvent(s, fromAsset, fromAddr, toChainID, toAsset, toAddr, amount)
	assert.NoError(t, err)
}
