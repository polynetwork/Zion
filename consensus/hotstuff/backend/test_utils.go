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
	"bytes"
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	tu "github.com/ethereum/go-ethereum/consensus/hotstuff/testutils"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	elog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	testLogger = elog.New()
)

// in this test, we can set n to 1, and it means we can process Istanbul and commit a
// block by one node. Otherwise, if n is larger than 1, we have to generate
// other fake events to process Istanbul.
func singleNodeChain() (*core.BlockChain, *backend) {
	testLogger.SetHandler(elog.StdoutHandler)

	genesis, nodeKeys, _ := tu.GenesisAndKeys(1)
	memDB := rawdb.NewMemoryDatabase()
	config := hotstuff.DefaultBasicConfig
	// Use the first key as private key
	backend := New(config, nodeKeys[0], memDB, true)
	genesis.MustCommit(memDB)

	txLookUpLimit := uint64(100)
	cacheConfig := &core.CacheConfig{
		TrieCleanLimit: 256,
		TrieDirtyLimit: 256,
		TrieTimeLimit:  5 * time.Minute,
	}
	blockchain, err := core.NewBlockChain(memDB, cacheConfig, genesis.Config, backend, vm.Config{}, nil, &txLookUpLimit)
	if err != nil {
		panic(err)
	}

	if err := backend.Start(blockchain, nil); err != nil {
		panic(err)
	}

	return blockchain, backend
}

func makeHeader(parent *types.Block) *types.Header {
	const (
		GasFloor, GasCeil = uint64(300000000), uint64(300000000)
	)

	blockNumber := parent.Number().Add(parent.Number(), common.Big1)
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     blockNumber,
		GasLimit:   core.CalcGasLimit(parent, GasFloor, GasCeil),
		GasUsed:    0,
		Time:       parent.Time() + hotstuff.DefaultBasicConfig.BlockPeriod,
		Difficulty: defaultDifficulty,
		Extra:      []byte{},
	}
	return header
}

func makeBlock(chain *core.BlockChain, engine *backend, parent *types.Block) *types.Block {
	block := makeBlockWithoutSeal(chain, engine, parent)
	stopCh := make(chan struct{})
	resultCh := make(chan *types.Block, 10)
	go engine.Seal(chain, block, resultCh, stopCh)
	blk := <-resultCh
	return blk
}

func makeBlockWithoutSeal(chain *core.BlockChain, engine *backend, parent *types.Block) *types.Block {
	header := makeHeader(parent)
	if err := engine.Prepare(chain, header); err != nil {
		panic(err)
	}
	state, err := chain.StateAt(parent.Root())
	if err != nil {
		panic(err)
	}
	if err := engine.FillHeader(state, header); err != nil {
		panic(err)
	}
	block, _, err := engine.FinalizeAndAssemble(chain, header, state, nil, nil, nil)
	if err != nil {
		panic(err)
	}
	return block
}

type Keys []*ecdsa.PrivateKey

func (slice Keys) Len() int {
	return len(slice)
}

func (slice Keys) Less(i, j int) bool {
	return strings.Compare(crypto.PubkeyToAddress(slice[i].PublicKey).String(), crypto.PubkeyToAddress(slice[j].PublicKey).String()) < 0
}

func (slice Keys) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func makeMsg(msgcode uint64, data interface{}) p2p.Msg {
	size, r, _ := rlp.EncodeToReader(data)
	return p2p.Msg{Code: msgcode, Size: uint32(size), Payload: r}
}

func postAndWait(backend *backend, block *types.Block, t *testing.T) {
	eventSub := backend.EventMux().Subscribe(hotstuff.RequestEvent{})
	defer eventSub.Unsubscribe()
	stop := make(chan struct{}, 1)
	eventLoop := func() {
		<-eventSub.Chan()
		stop <- struct{}{}
	}
	go eventLoop()
	if err := backend.EventMux().Post(hotstuff.RequestEvent{
		Block: block,
	}); err != nil {
		t.Fatalf("%s", err)
	}
	<-stop
}

func buildArbitraryP2PNewBlockMessage(t *testing.T, invalidMsg bool) (*types.Block, p2p.Msg) {
	arbitraryBlock := types.NewBlock(&types.Header{
		Number:    big.NewInt(1),
		GasLimit:  0,
		MixDigest: types.HotstuffDigest,
	}, nil, nil, nil, nil)
	request := []interface{}{&arbitraryBlock, big.NewInt(1)}
	if invalidMsg {
		request = []interface{}{"invalid msg"}
	}
	size, r, err := rlp.EncodeToReader(request)
	if err != nil {
		t.Fatalf("can't encode due to %s", err)
	}
	payload, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("can't read payload due to %s", err)
	}
	arbitraryP2PMessage := p2p.Msg{Code: NewBlockMsg, Size: uint32(size), Payload: bytes.NewReader(payload)}
	return arbitraryBlock, arbitraryP2PMessage
}
