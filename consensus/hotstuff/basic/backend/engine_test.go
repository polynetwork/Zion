// Copyright 2017 The go-ethereum Authors
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
	"bytes"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

// TestSealStopChannel stop consensus before result committed
func TestSealStopChannel(t *testing.T) {
	chain, engine := singleNodeChain()
	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	stop := make(chan struct{}, 1)
	eventSub := engine.EventMux().Subscribe(hotstuff.RequestEvent{})
	eventLoop := func() {
		ev := <-eventSub.Chan()
		_, ok := ev.Data.(hotstuff.RequestEvent)
		if !ok {
			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
		}
		stop <- struct{}{}
		eventSub.Unsubscribe()
	}
	go eventLoop()
	resultCh := make(chan *types.Block, 10)
	go func() {
		err := engine.Seal(chain, block, resultCh, stop)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
	}()

	finalBlock := <-resultCh
	assert.Nil(t, finalBlock)
}

func TestSealPreCommitOtherHash(t *testing.T) {
	chain, engine := singleNodeChain()
	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	otherBlock := makeBlockWithoutSeal(chain, engine, block)
	expectedCommittedSeal := append([]byte{1, 2, 3}, bytes.Repeat([]byte{0x00}, types.HotstuffExtraSeal-3)...)

	eventSub := engine.EventMux().Subscribe(hotstuff.RequestEvent{})
	blockOutputChannel := make(chan *types.Block)
	stopChannel := make(chan struct{})

	go func() {
		ev := <-eventSub.Chan()
		if _, ok := ev.Data.(hotstuff.RequestEvent); !ok {
			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
		}
		if _, err := engine.PreCommit(otherBlock, [][]byte{expectedCommittedSeal}); err != nil {
			t.Error(err.Error())
		}
		eventSub.Unsubscribe()
	}()

	go func() {
		if err := engine.Seal(chain, block, blockOutputChannel, stopChannel); err != nil {
			t.Error(err.Error())
		}
	}()

	select {
	case <-blockOutputChannel:
		t.Error("Wrong block found!")
	default:
		//no block found, stop the sealing
		close(stopChannel)
	}

	output := <-blockOutputChannel
	if output != nil {
		t.Error("Block not nil!")
	}
}

func TestSealCommitted(t *testing.T) {
	chain, engine := singleNodeChain()
	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header := block.Header()

	err := engine.signer.SealBeforeCommit(header)
	assert.NoError(t, err, "fillExtraBeforeCommit err", err)
	expectedBlock := block.WithSeal(header)

	resultCh := make(chan *types.Block, 10)
	go func() {
		err := engine.Seal(chain, block, resultCh, make(chan struct{}))

		if err != nil {
			t.Errorf("error mismatch: have %v, want %v", err, block)
		}
	}()

	finalBlock := <-resultCh
	if finalBlock.Hash() != expectedBlock.Hash() {
		t.Errorf("hash mismatch: have %v, want %v", finalBlock.Hash(), expectedBlock.Hash())
	}
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/basic/backend -run TestInsertChain
func TestInsertChain(t *testing.T) {
	chain, engine := singleNodeChain()
	expectBlock := makeBlock(t, chain, engine, chain.Genesis())
	chain.InsertChain(types.Blocks{expectBlock})
	block := chain.GetBlockByNumber(1)
	assert.Equal(t, expectBlock, block)
}

// todo: mock p2p