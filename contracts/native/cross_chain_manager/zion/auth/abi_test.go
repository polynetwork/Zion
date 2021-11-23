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

package auth

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestABIMethodApproveInput(t *testing.T) {
	expect := &MethodApproveInput{
		Spender: common.HexToAddress("0x12"),
		Amount:  big.NewInt(3),
	}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodApproveInput)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)
}

func TestABIMethodAllowanceInput(t *testing.T) {
	expect := &MethodAllowanceInput{
		Owner:   common.HexToAddress("0x12"),
		Spender: this,
	}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodAllowanceInput)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)
}

func TestEmitApproval(t *testing.T) {
	resetTestContext()
	s := testEmptyCtx
	assert.NoError(t, emitApprovedEvent(s, common.HexToAddress("0x13"), this, big.NewInt(13)))
}