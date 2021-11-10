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
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/stretchr/testify/assert"
)

func TestBindProxy(t *testing.T) {
	resetTestContext()
	s := testEmptyCtx
	nm.InitABI()

	epoch := nm.GenerateTestEpochInfo(1, uint64(testBlockNum-1), 4)
	assert.NoError(t, nm.StoreTestEpoch(s, epoch))

	toChainID := uint64(12)
	targetProxy := common.HexToAddress("0x123a234d3")
	for index, v := range epoch.Peers.List {
		_, _, err := testCallBindProxy(v.Address, toChainID, targetProxy[:])
		if err != nil {
			t.Logf("proposer %d bind message: %s", index, err.Error())
		}
	}

	blob, err := getProxy(s, toChainID)
	assert.NoError(t, err)
	assert.Equal(t, targetProxy[:], blob)
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

	txHash := nm.GenerateTestHash(120)
	ref := native.NewContractRef(testStateDB, sender, sender, big.NewInt(testBlockNum), txHash, testSupplyGas, nil)
	ref.PushContext(&native.Context{
		Caller:          sender,
		ContractAddress: this,
		Payload:         payload,
	})
	ctx := native.NewNativeContract(testStateDB, ref)
	if ret, err := BindProxy(ctx); err != nil {
		return nil, nil, err
	} else {
		return ctx, ret, nil
	}
}
