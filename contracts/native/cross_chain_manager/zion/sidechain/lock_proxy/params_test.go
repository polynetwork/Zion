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

func TestCrossTxEncode(t *testing.T) {
	expect := &CrossTx{
		ToChainId:   2,
		FromAddress: common.HexToAddress("0x12"),
		ToAddress:   common.HexToAddress("0x33"),
		Amount:      big.NewInt(34),
		Index:       12,
	}

	payload, err := EncodeCrossTx(expect)
	assert.NoError(t, err)

	got, err := DecodeCrossTx(payload)
	assert.NoError(t, err)

	assert.Equal(t, expect, got)
}
