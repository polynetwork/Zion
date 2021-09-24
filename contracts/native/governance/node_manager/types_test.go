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

package node_manager

import (
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

func TestPeerInfo(t *testing.T) {
	expect := generateTestPeer()

	enc, err := rlp.EncodeToBytes(expect)
	assert.NoError(t, err)

	var got *PeerInfo
	assert.NoError(t, rlp.DecodeBytes(enc, &got))

	assert.Equal(t, expect, got)
	t.Logf("peer info length %d", len(enc))
}

func TestPeers(t *testing.T) {
	expect := generateTestPeers(10)

	enc, err := rlp.EncodeToBytes(&expect)
	assert.NoError(t, err)

	var got *Peers
	assert.NoError(t, rlp.DecodeBytes(enc, &got))
	assert.Equal(t, expect, got)
}

func TestEpochInfo(t *testing.T) {
	expect := generateTestEpochInfo(1, 32, 10)

	enc, err := rlp.EncodeToBytes(expect)
	assert.NoError(t, err)

	var got *EpochInfo
	assert.NoError(t, rlp.DecodeBytes(enc, &got))
	assert.Equal(t, expect, got)

	// test load hash
	assert.Empty(t, expect.hash)
	expectHash := expect.Hash()
	assert.NotEmpty(t, expect.hash)
	assert.Equal(t, expectHash, expect.Hash())
}

func TestHashList(t *testing.T) {
	expect := generateTestHashList(12)

	enc, err := rlp.EncodeToBytes(expect)
	assert.NoError(t, err)

	var got *HashList
	assert.NoError(t, rlp.DecodeBytes(enc, &got))

	assert.Equal(t, expect, got)
}

func TestAddressList(t *testing.T) {
	expect := generateTestAddressList(13)

	enc, err := rlp.EncodeToBytes(expect)
	assert.NoError(t, err)

	var got *AddressList
	assert.NoError(t, rlp.DecodeBytes(enc, &got))

	assert.Equal(t, expect, got)
}
