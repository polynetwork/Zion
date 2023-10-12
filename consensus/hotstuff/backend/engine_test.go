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
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/signer"
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

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestSealOtherHash
// TestSealCommittedOtherHash result channel should be empty if engine commit another block before seal
func TestSealOtherHash(t *testing.T) {
	chain, engine := singleNodeChain()
	defer engine.Stop()
	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())

	blockCh := make(chan consensus.ExecutedBlock)
	blockSub := engine.SubscribeBlock(blockCh)
	requestCh := make(chan hotstuff.RequestEvent)
	eventSub := engine.SubscribeEvent(requestCh)
	blockOutputChannel := make(chan *types.Block)
	stopChannel := make(chan struct{})

	go func() {
		if err := engine.Seal(chain, block, blockOutputChannel, stopChannel); err != nil {
			t.Error(err.Error())
		}
		if err := engine.Commit(&consensus.ExecutedBlock{Block: block}); err != nil {
			t.Error(err.Error())
		}
	}()

	evt := <- requestCh
	assert.Equal(t, evt.Block.SealHash(), block.SealHash())

	data := <-blockCh
	assert.Equal(t, data.Block.SealHash(), block.SealHash())

	eventSub.Unsubscribe()
	blockSub.Unsubscribe()
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

	go func() {
		if err := engine.Seal(chain, block, make(chan *types.Block), make(chan struct{})); err != nil {
			t.Errorf("error mismatch: have %v, want %v", err, expectedBlock)
		}
	}()

	requestCh := make(chan hotstuff.RequestEvent)
	eventSub := engine.SubscribeEvent(requestCh)
	ev := <- requestCh
	eventSub.Unsubscribe()
	finalBlock := ev.Block

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
