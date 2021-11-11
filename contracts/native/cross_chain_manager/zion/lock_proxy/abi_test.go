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
	"github.com/stretchr/testify/assert"
)

func TestABIMethodContractNameOutput(t *testing.T) {
	expect := &MethodContractNameOutput{Name: contractName}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodContractNameOutput)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)
}

func TestABIMethodBindProxyInput(t *testing.T) {
	expect := &MethodBindProxyInput{
		ToChainId:       12,
		TargetProxyHash: common.HexToHash("0x2347123de324").Bytes(),
	}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodBindProxyInput)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)
}

func TestABIMethodGetProxyInput(t *testing.T) {
	expect := &MethodGetProxyInput{ToChainId: 33}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodGetProxyInput)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)
}

func TestABIMethodBindAssetHashInput(t *testing.T) {
	expect := &MethodBindAssetHashInput{
		FromAssetHash: common.HexToAddress("0x12e345d"),
		ToChainId:     39,
		ToAssetHash:   []byte{'1', 'a', 'c', '2'},
	}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodBindAssetHashInput)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)
}

func TestABIMethodGetAssetInput(t *testing.T) {
	expect := &MethodGetAssetInput{
		FromAssetHash: common.HexToAddress("0x234abc234"),
		ToChainId:     79,
	}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodGetAssetInput)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)
}

func TestABIMethodLockInput(t *testing.T) {
	expect := &MethodLockInput{
		FromAssetHash: common.HexToAddress("0x123456"),
		ToChainId:     13,
		ToAddress:     common.HexToAddress("0x335").Bytes(),
		Amount:        big.NewInt(123648),
	}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodLockInput)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)
}
