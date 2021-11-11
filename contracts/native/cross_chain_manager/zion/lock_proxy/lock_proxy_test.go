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
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/stretchr/testify/assert"
)

func TestBindProxy(t *testing.T) {
	epoch := prepare(t)

	toChainID := uint64(12)
	targetProxy := common.HexToAddress("0x123a234d3")
	for index, v := range epoch.Peers.List {
		_, _, err := testCallBindProxy(v.Address, toChainID, targetProxy[:])
		if err != nil {
			t.Logf("proposer %d bind message: %s", index, err.Error())
		}
	}

	_, blob, err := testCallGetProxy(toChainID)
	assert.NoError(t, err)
	assert.Equal(t, targetProxy[:], blob)
}

func TestBindAsset(t *testing.T) {
	epoch := prepare(t)
	testStateDB.SetBalance(this, new(big.Int).Mul(minBalance, big.NewInt(2)))

	fromAsset := common.EmptyAddress
	toChainID := uint64(12)
	toAsset := []byte{'1', 'a', '3', 'd', '9'}
	for index, v := range epoch.Peers.List {
		_, _, err := testCallBindAsset(v.Address, fromAsset, toChainID, toAsset)
		if err != nil {
			t.Logf("proposer %d bind message: %s", index, err.Error())
		}
	}

	_, blob, err := testCallGetAsset(fromAsset, toChainID)
	assert.NoError(t, err)
	assert.Equal(t, toAsset, blob)
}

func TestLock(t *testing.T) {
	a := big.NewInt(3)
	data := a.Bytes()
	t.Log(len(data))
}

func TestUnlock(t *testing.T) {

}

func testCallBindProxy(sender common.Address, toChainID uint64, targetProxy []byte) (*native.NativeContract, []byte, error) {
	input := &MethodBindProxyInput{
		ToChainId:       toChainID,
		TargetProxyHash: targetProxy,
	}
	payload, err := input.Encode()
	if err != nil {
		return nil, nil, err
	}

	ctx := generateTestSenderTx(sender, payload)
	if ret, err := BindProxy(ctx); err != nil {
		return nil, nil, err
	} else {
		return ctx, ret, nil
	}
}

func testCallGetProxy(toChainID uint64) (*native.NativeContract, []byte, error) {
	input := &MethodGetProxyInput{ToChainId: toChainID}
	payload, err := input.Encode()
	if err != nil {
		return nil, nil, err
	}

	ctx := generateTestCallCtx(payload)
	if ret, err := GetProxy(ctx); err != nil {
		return nil, nil, err
	} else {
		return ctx, ret, nil
	}
}

func testCallBindAsset(sender, fromAsset common.Address, toChainID uint64, toAsset []byte) (*native.NativeContract, []byte, error) {
	input := &MethodBindAssetHashInput{
		FromAssetHash: fromAsset,
		ToChainId:     toChainID,
		ToAssetHash:   toAsset,
	}
	payload, err := input.Encode()
	if err != nil {
		return nil, nil, err
	}

	ctx := generateTestSenderTx(sender, payload)
	if ret, err := BindAsset(ctx); err != nil {
		return nil, nil, err
	} else {
		return ctx, ret, nil
	}
}

func testCallGetAsset(fromAsset common.Address, toChainID uint64) (*native.NativeContract, []byte, error) {
	input := &MethodGetAssetInput{
		FromAssetHash: fromAsset,
		ToChainId:     toChainID,
	}
	payload, err := input.Encode()
	if err != nil {
		return nil, nil, err
	}

	ctx := generateTestCallCtx(payload)
	if ret, err := GetAsset(ctx); err != nil {
		return nil, nil, err
	} else {
		return ctx, ret, nil
	}
}

func testLock(sender, fromAsset common.Address, toChainID uint64, toAddress []byte, amount *big.Int) (*native.NativeContract, []byte, error) {
	input := &MethodLockInput{
		FromAssetHash: fromAsset,
		ToChainId:     toChainID,
		ToAddress:     toAddress,
		Amount:        amount,
	}
	payload, err := input.Encode()
	if err != nil {
		return nil, nil, err
	}

	ctx := generateTestSenderTx(sender, payload)
	if ret, err := Lock(ctx); err != nil {
		return nil, nil, err
	} else {
		return ctx, ret, nil
	}
}

func testUnlock(sender common.Address, entranParams *scom.EntranceParam, makeTxParams *scom.MakeTxParam) (*native.NativeContract, error) {
	ctx := generateTestSenderTx(sender, nil)
	if err := Unlock(ctx, entranParams, makeTxParams); err != nil {
		return nil, err
	} else {
		return ctx, nil
	}
}

func prepare(t *testing.T) *nm.EpochInfo {
	resetTestContext()
	s := testEmptyCtx
	nm.InitABI()

	epoch := nm.GenerateTestEpochInfo(1, uint64(testBlockNum-1), 4)
	if err := nm.StoreTestEpoch(s, epoch); err != nil {
		t.Fatal(err)
	}
	return epoch
}

func generateTestSenderTx(sender common.Address, payload []byte) *native.NativeContract {
	txHash := nm.GenerateTestHash(rand.Int())
	ref := native.NewContractRef(testStateDB, sender, sender, big.NewInt(testBlockNum), txHash, testSupplyGas, nil)
	ref.PushContext(&native.Context{
		Caller:          sender,
		ContractAddress: this,
		Payload:         payload,
	})
	return native.NewNativeContract(testStateDB, ref)
}

func generateTestCallCtx(payload []byte) *native.NativeContract {
	caller := common.EmptyAddress
	txHash := nm.GenerateTestHash(rand.Int())
	ref := native.NewContractRef(testStateDB, caller, caller, big.NewInt(testBlockNum), txHash, testSupplyGas, nil)
	ref.PushContext(&native.Context{
		Caller:          caller,
		ContractAddress: this,
		Payload:         payload,
	})
	return native.NewNativeContract(testStateDB, ref)
}
