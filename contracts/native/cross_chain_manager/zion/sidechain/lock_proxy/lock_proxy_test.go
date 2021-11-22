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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
)

func testBurn(toChainID uint64, sender, toAddress common.Address, amount *big.Int) (*native.NativeContract, []byte, error) {
	input := &MethodBurnInput{
		ToChainId: toChainID,
		ToAddress: toAddress,
		Amount:    amount,
	}
	payload, err := input.Encode()
	if err != nil {
		return nil, nil, err
	}

	ctx := generateTestSenderTx(sender, payload)
	if ret, err := Burn(ctx); err != nil {
		return nil, nil, err
	} else {
		return ctx, ret, nil
	}
}

func testUnlock(sender common.Address, header, crossTx, proof, extra []byte) (*native.NativeContract, []byte, error) {
	input := &MethodVerifyHeaderAndExecuteTxInput{
		Header:     header,
		RawCrossTx: crossTx,
		Proof:      proof,
		Extra:      extra,
	}
	payload, err := input.Encode()
	if err != nil {
		return nil, nil, err
	}

	ctx := generateTestSenderTx(sender, payload)
	if ret, err := Mint(ctx); err != nil {
		return nil, nil, err
	} else {
		return ctx, ret, nil
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
