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

package okex

import (
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

func TestCosmosEpochSwitchInfo(t *testing.T) {
	expect := &CosmosEpochSwitchInfo{
		Height:             12,
		BlockHash:          []byte{'1', '3', '5'},
		NextValidatorsHash: []byte{'2', '4', '6'},
		ChainID:            "hahaha",
	}

	blob, err := rlp.EncodeToBytes(expect)
	assert.NoError(t, err)

	got := new(CosmosEpochSwitchInfo)
	assert.NoError(t, rlp.DecodeBytes(blob, got))

	assert.Equal(t, expect, got)
}
