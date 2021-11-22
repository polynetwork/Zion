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

func TestBalanceFor(t *testing.T) {
	resetTestContext()
	s := testEmptyCtx

	owner := common.HexToAddress("0x123")
	expect := big.NewInt(1234234234)
	testStateDB.SetBalance(owner, expect)

	got := getBalanceFor(s, owner)
	assert.Equal(t, expect, got)
}

func TestTransferFromContract(t *testing.T) {
	resetTestContext()
	s := testEmptyCtx

	lockProxyBalance := big.NewInt(10000)
	testStateDB.SetBalance(this, lockProxyBalance)

	expectTransferAmount := big.NewInt(5000)
	expectResAmount := new(big.Int).Sub(lockProxyBalance, expectTransferAmount)
	to := common.HexToAddress("0x12342345")

	assert.NoError(t, transferFromContract(s, to, expectTransferAmount))

	assert.Equal(t, expectResAmount.Uint64(), testStateDB.GetBalance(this).Uint64())
	assert.Equal(t, expectTransferAmount.Uint64(), testStateDB.GetBalance(to).Uint64())
}
