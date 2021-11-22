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
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/utils"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/stretchr/testify/assert"
)

func TestLockAndUnlock(t *testing.T) {
	targetChainID := uint64(12)
	fromAsset := common.EmptyAddress
	targetCaller := common.HexToAddress("0x2").Bytes()
	sender := common.HexToAddress("0x4")
	receiver := common.HexToAddress("0x5")
	amount := minBalance

	testStateDB.SetBalance(this, new(big.Int).Mul(minBalance, big.NewInt(1)))
	testStateDB.SetBalance(sender, new(big.Int).Mul(minBalance, big.NewInt(2)))

	_, _, err := testLock(sender, fromAsset, targetChainID, receiver.Bytes(), amount)
	assert.NoError(t, err)

	txArgs := utils.EncodeTxArgs(fromAsset.Bytes(), sender.Bytes(), amount)
	txParams := &scom.MakeTxParam{
		CrossChainID:        []byte{'1', 'a'},
		FromContractAddress: targetCaller,
		ToChainID:           native.ZionMainChainID,
		ToContractAddress:   this.Bytes(),
		Method:              "unlock",
		Args:                txArgs,
	}
	entranParams := &scom.EntranceParam{SourceChainID: targetChainID}
	_, err = testUnlock(receiver, entranParams, txParams, amount)
	assert.NoError(t, err)
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
	ctx.ContractRef().SetValue(amount)
	if ret, err := Lock(ctx); err != nil {
		return nil, nil, err
	} else {
		return ctx, ret, nil
	}
}

func testUnlock(sender common.Address, entranParams *scom.EntranceParam, makeTxParams *scom.MakeTxParam, amount *big.Int) (*native.NativeContract, error) {
	ctx := generateTestSenderTx(sender, nil)
	ctx.ContractRef().SetValue(amount)
	if err := Unlock(ctx, entranParams, makeTxParams); err != nil {
		return nil, err
	} else {
		return ctx, nil
	}
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
