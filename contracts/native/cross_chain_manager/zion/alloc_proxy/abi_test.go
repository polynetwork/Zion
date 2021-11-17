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

package alloc_proxy

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

func TestABIMethodInitGenesisHeaderInput(t *testing.T) {
	expect := &MethodInitGenesisHeaderInput{
		Header: []byte{'1', 'a'},
		Proof:  []byte{'d', '3'},
		Extra:  []byte{'8', '6', 'a'},
		Epoch:  []byte{'d', '3'},
	}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodInitGenesisHeaderInput)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)
}

func TestABIMethodChangeEpochInput(t *testing.T) {
	expect := &MethodChangeEpochInput{
		Header: []byte{'1', 'a'},
		Proof:  []byte{'d', '3'},
		Extra:  []byte{'8', '6', 'a'},
		Epoch:  []byte{'d', '3'},
	}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodChangeEpochInput)
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
	expect := &MethodVerifyHeaderAndMintInput{
		Header:     []byte{'a'},
		RawCrossTx: []byte{},
		Proof:      []byte{'a'},
		Extra:      []byte{},
	}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodVerifyHeaderAndMintInput)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)
}
