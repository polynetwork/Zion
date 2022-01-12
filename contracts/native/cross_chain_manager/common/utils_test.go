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

package common

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestUint256ToBytes(t *testing.T) {
	var testcases = []struct {
		Num    *big.Int
		Expect *big.Int
	}{
		{
			Num:    nil,
			Expect: common.Big0,
		},
		{
			Num:    common.Big0,
			Expect: common.Big0,
		},
		{
			Num:    big.NewInt(32),
			Expect: big.NewInt(32),
		},
	}

	for _, v := range testcases {
		blob := Uint256ToBytes(v.Num)
		expectHex := hexutil.Encode(blob)

		gotHex := hexutil.Encode(blob)
		assert.Equal(t, expectHex, gotHex)

		gotNum := BytesToUint256(blob)
		assert.Equal(t, v.Expect, gotNum)
	}
}
