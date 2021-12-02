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

func TestABIMethodLockInput(t *testing.T) {
	expect := &MethodLockInput{
		ToChainId:     13,
		ToAddress:     common.HexToAddress("0x335"),
		Amount:        big.NewInt(123648),
	}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodLockInput)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)
}

func TestABIMethodGetSideChainLockAmountInput(t *testing.T) {
	expect := &MethodGetSideChainLockAmountInput{ChainId: 12}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodGetSideChainLockAmountInput)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)
}
