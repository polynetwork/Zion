// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package backend

import (
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/assert"
)

func TestHotstuffMessage(t *testing.T) {
	_, backend := newBlockChain(1)

	// generate one msg
	data := []byte("data1")
	hash := hotstuff.RLPHash(data)
	msg := makeMsg(hotstuffMsg, data)
	addr := common.HexToAddress("address")

	// 1. this message should not be in cache
	// for peers
	_, ok := backend.recentMessages.Get(addr)
	assert.False(t, ok, "the cache of messages for this peer should be nil")

	// for self
	_, ok = backend.knownMessages.Get(hash)
	assert.False(t, ok, "the cache of messages should be nil")

	// 2. this message should be in cache after we handle it
	_, err := backend.HandleMsg(addr, msg)
	assert.NoError(t, err, "handle message failed:", err)

	// for peers
	if ms, ok := backend.recentMessages.Get(addr); ms == nil || !ok {
		t.Fatalf("the cache of messages for this peer cannot be nil")
	} else if m, ok := ms.(*lru.ARCCache); !ok {
		t.Fatalf("the cache of messages for this peer cannot be casted")
	} else if _, ok := m.Get(hash); !ok {
		t.Fatalf("the cache of messages for this peer cannot be found")
	}

	// for self
	_, ok = backend.knownMessages.Get(hash)
	assert.True(t, ok, "the cache of messages cannot be found")
}

func TestHandleNewBlockMessage_whenTypical(t *testing.T) {
	_, backend := newBlockChain(1)
	arbitraryAddress := common.HexToAddress("arbitrary")
	arbitraryBlock, arbitraryP2PMessage := buildArbitraryP2PNewBlockMessage(t, false)
	postAndWait(backend, arbitraryBlock, t)

	handled, err := backend.HandleMsg(arbitraryAddress, arbitraryP2PMessage)
	assert.NoError(t, err, "expected message being handled successfully but got", err)
	assert.False(t, handled, "expected message not being handled")

	_, err = ioutil.ReadAll(arbitraryP2PMessage.Payload)
	assert.NoError(t, err, "expected p2p message payload is restored")
}

func TestHandleNewBlockMessage_whenNotAProposedBlock(t *testing.T) {
	_, backend := newBlockChain(1)
	arbitraryAddress := common.HexToAddress("arbitrary")
	_, arbitraryP2PMessage := buildArbitraryP2PNewBlockMessage(t, false)
	postAndWait(backend, types.NewBlock(&types.Header{
		Number:    big.NewInt(1),
		Root:      common.HexToHash("someroot"),
		GasLimit:  1,
		MixDigest: types.HotstuffDigest,
	}, nil, nil, nil, nil), t)

	handled, err := backend.HandleMsg(arbitraryAddress, arbitraryP2PMessage)
	assert.NoError(t, err, "expected message being handled successfully but got", err)
	assert.False(t, handled, "expected message not being handled")

	_, err = ioutil.ReadAll(arbitraryP2PMessage.Payload)
	assert.NoError(t, err, "expected p2p message payload is restored")
}

func TestHandleNewBlockMessage_whenFailToDecode(t *testing.T) {
	_, backend := newBlockChain(1)
	arbitraryAddress := common.HexToAddress("arbitrary")
	_, arbitraryP2PMessage := buildArbitraryP2PNewBlockMessage(t, true)
	postAndWait(backend, types.NewBlock(&types.Header{
		Number:    big.NewInt(1),
		GasLimit:  1,
		MixDigest: types.HotstuffDigest,
	}, nil, nil, nil, nil), t)

	handled, err := backend.HandleMsg(arbitraryAddress, arbitraryP2PMessage)
	assert.NoError(t, err, "expected message being handled successfully but got", err)
	assert.False(t, handled, "expected message not being handled")

	_, err = ioutil.ReadAll(arbitraryP2PMessage.Payload)
	assert.NoError(t, err, "expected p2p message payload is restored")
}
