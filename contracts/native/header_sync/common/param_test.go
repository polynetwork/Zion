/*
 * Copyright (C) 2021 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package common

import (
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

func TestSyncGenesisHeaderParam(t *testing.T) {
	expect := &SyncGenesisHeaderParam{
		ChainID:       123,
		GenesisHeader: []byte{1, 2, 3},
	}

	blob, err := rlp.EncodeToBytes(expect)
	assert.NoError(t, err)

	got := new(SyncGenesisHeaderParam)
	assert.NoError(t, rlp.DecodeBytes(blob, got))

	assert.Equal(t, expect, got)
}

func TestSyncBlockHeaderParam(t *testing.T) {
	expect := &SyncBlockHeaderParam{
		ChainID: 123,
		Headers: [][]byte{{1, 2, 3}},
	}

	blob, err := rlp.EncodeToBytes(expect)
	assert.NoError(t, err)

	got := new(SyncBlockHeaderParam)
	assert.NoError(t, rlp.DecodeBytes(blob, got))

	assert.Equal(t, expect, got)
}
