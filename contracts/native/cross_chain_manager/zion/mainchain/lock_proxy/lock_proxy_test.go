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
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	nu "github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/stretchr/testify/assert"
)

func TestLockAndUnlock(t *testing.T) {
	targetChainID := uint64(12)
	fromAsset := common.EmptyAddress
	sender := common.HexToAddress("0x4")
	receiver := common.HexToAddress("0x5")
	minBalance, _ := new(big.Int).SetString("1000000000000000000", 10)
	amount := minBalance
	cnt := 10

	testStateDB.SetBalance(this, new(big.Int).Mul(minBalance, big.NewInt(int64(cnt))))
	testStateDB.SetBalance(sender, new(big.Int).Mul(minBalance, big.NewInt(int64(cnt))))
	assert.NoError(t, testSetSideChain(targetChainID))

	for i := 0; i < cnt; i++ {
		_, _, err := testLock(sender, receiver, targetChainID, amount)
		assert.NoError(t, err)

		total, err := testGetLockAmount(targetChainID)
		assert.NoError(t, err)
		t.Logf("current lock amount %v", total)

		txArgs, err := utils.EncodeTxArgs(fromAsset.Bytes(), sender.Bytes(), amount)
		assert.NoError(t, err)
		
		txParams := &scom.MakeTxParam{
			CrossChainID:        []byte{'1', 'a'},
			FromContractAddress: this[:],
			ToChainID:           native.ZionMainChainID,
			ToContractAddress:   this.Bytes(),
			Method:              "unlock",
			Args:                txArgs,
		}
		_, err = testUnlock(receiver, targetChainID, txParams, amount)
		assert.NoError(t, err)

		total, err = testGetLockAmount(targetChainID)
		assert.NoError(t, err)
		t.Logf("current lock amount %v", total)
	}
}

func testLock(sender, toAddress common.Address, toChainID uint64, amount *big.Int) (*native.NativeContract, []byte, error) {
	input := &MethodLockInput{
		ToChainId:     toChainID,
		ToAddress:     toAddress,
		Amount:        amount,
	}
	payload, err := input.Encode()
	if err != nil {
		return nil, nil, err
	}

	ctx := generateTestSenderTx(sender, sender, payload)
	ctx.ContractRef().SetValue(amount)
	ctx.ContractRef().SetTo(this)
	if ret, err := Lock(ctx); err != nil {
		return nil, nil, err
	} else {
		return ctx, ret, nil
	}
}

func testUnlock(sender common.Address, srcChainID uint64, makeTxParams *scom.MakeTxParam, amount *big.Int) (*native.NativeContract, error) {
	entrance := nu.CrossChainManagerContractAddress
	ctx := generateTestSenderTx(sender, entrance,nil)
	ctx.ContractRef().SetValue(amount)
	ctx.ContractRef().SetTo(entrance)
	if err := Unlock(ctx, srcChainID, makeTxParams); err != nil {
		return nil, err
	} else {
		return ctx, nil
	}
}

func testGetLockAmount(chainID uint64) (*big.Int, error) {
	input := &MethodGetSideChainLockAmountInput{ChainId: chainID}
	payload, err := input.Encode()
	if err != nil {
		return nil, err
	}
	ctx := generateTestCallCtx(payload)

	enc, err := GetSideChainLockAmount(ctx)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(enc), nil
}

func generateTestSenderTx(sender, caller common.Address, payload []byte) *native.NativeContract {
	txHash := nm.GenerateTestHash(rand.Int())
	ref := native.NewContractRef(testStateDB, sender, caller, big.NewInt(testBlockNum), txHash, testSupplyGas, nil)
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

func testSetSideChain(chainID uint64) error {
	sc := new(side_chain_manager.SideChain)
	sc.Router = nu.ZION_ROUTER
	sc.ChainId = chainID
	s := generateTestCallCtx(nil)
	return side_chain_manager.PutSideChain(s, sc)
}
