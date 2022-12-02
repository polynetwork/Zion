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

package backend

import (
	"math/big"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/consensus/hotstuff/signer"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	tu "github.com/ethereum/go-ethereum/consensus/hotstuff/testutils"
	"github.com/ethereum/go-ethereum/contracts/native/boot"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	boot.InitNativeContracts()
	os.Exit(m.Run())
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestGenesisBlock
func TestGenesisBlock(t *testing.T) {
	memDB := rawdb.NewMemoryDatabase()
	g, _, err := tu.GenesisAndKeys(4)
	if err != nil {
		t.Error(err)
	}
	block := tu.GenesisBlock(g, memDB)
	raw, err := block.Header().MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	t.Logf("genesis block, header %s, mixDigest %v, number %d, hash %v, size %v",
		string(raw), block.MixDigest(), block.NumberU64(), block.Hash(), block.Size())
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestPrepare
func TestPrepare(t *testing.T) {
	chain, engine := singleNodeChain()
	defer engine.Stop()

	header := makeHeader(chain.Genesis())
	assert.NoError(t, engine.Prepare(chain, header))

	header.ParentHash = common.HexToHash("0x1234567890")
	assert.Error(t, engine.Prepare(chain, header), consensus.ErrUnknownAncestor)
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestSealStopChannel
// TestSealStopChannel stop seal and result channel should be empty
func TestSealStopChannel(t *testing.T) {
	chain, engine := singleNodeChain()
	defer engine.Stop()

	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	stop := make(chan struct{}, 1)
	eventSub := engine.EventMux().Subscribe(hotstuff.RequestEvent{})
	blockSub := engine.SubscribeBlock(make(chan consensus.ExecutedBlock))
	eventLoop := func() {
		ev := <-eventSub.Chan()
		if _, ok := ev.Data.(hotstuff.RequestEvent); !ok {
			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
		}
		stop <- struct{}{}
		eventSub.Unsubscribe()
		blockSub.Unsubscribe()
	}
	resultCh := make(chan *types.Block, 10)
	go func() {
		if err := engine.Seal(chain, block, resultCh, stop); err != nil {
			t.Errorf("error mismatch: have error %v, want nil", err)
		}
	}()
	go eventLoop()

	finalBlock := <-resultCh
	if finalBlock != nil {
		t.Errorf("block mismatch: have final block %v, want nil", finalBlock)
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestSealOtherHash
// TestSealCommittedOtherHash result channel should be empty if engine commit another block before seal
func TestSealOtherHash(t *testing.T) {
	chain, engine := singleNodeChain()
	defer engine.Stop()
	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	otherBlock := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	otherBlock.Header().GasUsed = 10

	blockCh := make(chan consensus.ExecutedBlock)
	blockSub := engine.SubscribeBlock(blockCh)
	eventSub := engine.EventMux().Subscribe(hotstuff.RequestEvent{})
	blockOutputChannel := make(chan *types.Block)
	stopChannel := make(chan struct{})

	go func() {
		ev := <-eventSub.Chan()
		if _, ok := ev.Data.(hotstuff.RequestEvent); !ok {
			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
		}
		if err := engine.Commit(&consensus.ExecutedBlock{Block: otherBlock}); err != nil {
			t.Error(err.Error())
		}
		eventSub.Unsubscribe()
		blockSub.Unsubscribe()
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

func updateTestBlock(block *types.Block, addr common.Address) *types.Block {
	header := block.Header()
	header.Coinbase = addr
	return block.WithSeal(header)
}

func updateTestBlockWithoutExtra(block *types.Block, addr common.Address) *types.Block {
	header := block.Header()
	header.Extra = []byte{}
	return block.WithSeal(header)
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestSealCommitted
// TestSealCommitted block hash WONT change after seal.
func TestSealCommitted(t *testing.T) {
	chain, engine := singleNodeChain()
	defer engine.Stop()

	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	expectedBlock := updateTestBlock(block, engine.Address())

	resultCh := make(chan *types.Block, 10)
	go func() {
		if err := engine.Seal(chain, block, resultCh, make(chan struct{})); err != nil {
			t.Errorf("error mismatch: have %v, want %v", err, expectedBlock)
		}
	}()

	finalBlock := <-resultCh
	if finalBlock.Hash() != expectedBlock.Hash() {
		t.Errorf("hash mismatch: have %v, want %v", finalBlock.Hash(), expectedBlock.Hash())
	}
}

// go test -count=1 -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestVerifyHeader
func TestVerifyHeader(t *testing.T) {
	chain, engine := singleNodeChain()
	defer engine.Stop()

	// errEmptyCommittedSeals case
	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header := engine.chain.GetHeader(block.ParentHash(), block.NumberU64()-1)
	block = updateTestBlock(block, engine.Address())
	if err := engine.VerifyHeader(chain, block.Header(), false); err != signer.ErrInvalidSignature {
		t.Errorf("error mismatch: have %v, want %v", err, "invalid signature length")
	}

	// short extra data
	header = block.Header()
	header.Extra = []byte{}
	if err := engine.VerifyHeader(chain, header, false); err != types.ErrInvalidHotstuffHeaderExtra {
		t.Errorf("error mismatch: have %v, want %v", err, types.ErrInvalidHotstuffHeaderExtra)
	}
	// incorrect extra format
	header.Extra = []byte("0000000000000000000000000000000012300000000000000000000000000000000000000000000000000000000000000000")
	if err := engine.VerifyHeader(chain, header, false); err != types.ErrInvalidHotstuffHeaderExtra {
		t.Errorf("error mismatch: have %v, want %v", err, types.ErrInvalidHotstuffHeaderExtra)
	}

	// non zero MixDigest
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.MixDigest = common.HexToHash("0x123456789")
	if err := engine.VerifyHeader(chain, header, false); err != errInvalidMixDigest {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidMixDigest)
	}

	// invalid uncles hash
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.UncleHash = common.HexToHash("0x123456789")
	if err := engine.VerifyHeader(chain, header, false); err != errInvalidUncleHash {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidUncleHash)
	}

	// invalid difficulty
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.Difficulty = big.NewInt(2)
	if err := engine.VerifyHeader(chain, header, false); err != errInvalidDifficulty {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidDifficulty)
	}

	// invalid timestamp
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.Time = chain.Genesis().Time() - 1

	if err := engine.VerifyHeader(chain, header, false); err != errInvalidTimestamp {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidTimestamp)
	}

	// future block
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.Time = uint64(time.Now().Unix() + 1)
	if err := engine.VerifyHeader(chain, header, false); err != consensus.ErrFutureBlock {
		t.Errorf("error mismatch: have %v, want %v", err, consensus.ErrFutureBlock)
	}
}
