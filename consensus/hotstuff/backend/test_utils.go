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
	snr "github.com/ethereum/go-ethereum/consensus/hotstuff/signer"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	elog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	testLogger = elog.New()
)

func getGenesisAndKeys(n int) (*core.Genesis, []*ecdsa.PrivateKey, hotstuff.ValidatorSet) {
	// Setup validators
	var nodeKeys = make([]*ecdsa.PrivateKey, n)
	var addrs = make([]common.Address, n)
	for i := 0; i < n; i++ {
		nodeKeys[i], _ = crypto.GenerateKey()
		addrs[i] = crypto.PubkeyToAddress(nodeKeys[i].PublicKey)
	}

	// generate genesis block
	genesis := core.TestGenesisBlock()
	genesis.Config = params.TestChainConfig
	// force enable Istanbul engine
	genesis.Config.HotStuff = &params.HotStuffConfig{}
	genesis.Config.Ethash = nil
	genesis.Difficulty = defaultDifficulty
	genesis.Nonce = emptyNonce.Uint64()
	genesis.Mixhash = types.HotstuffDigest

	appendValidators(genesis, addrs)
	valset := makeValSet(addrs)
	return genesis, nodeKeys, valset
}

func appendValidators(genesis *core.Genesis, addrs []common.Address) {

	if len(genesis.ExtraData) < types.HotstuffExtraVanity {
		genesis.ExtraData = append(genesis.ExtraData, bytes.Repeat([]byte{0x00}, types.HotstuffExtraVanity)...)
	}
	genesis.ExtraData = genesis.ExtraData[:types.HotstuffExtraVanity]

	ist := &types.HotstuffExtra{
		Validators:    addrs,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	istPayload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		panic("failed to encode istanbul extra")
	}
	genesis.ExtraData = append(genesis.ExtraData, istPayload...)
}

// in this test, we can set n to 1, and it means we can process Istanbul and commit a
// block by one node. Otherwise, if n is larger than 1, we have to generate
// other fake events to process Istanbul.
func singleNodeChain() (*core.BlockChain, *backend) {
	testLogger.SetHandler(elog.StdoutHandler)

	genesis, nodeKeys, _ := getGenesisAndKeys(1)
	memDB := rawdb.NewMemoryDatabase()
	config := hotstuff.DefaultBasicConfig
	// Use the first key as private key
	b, _ := New(config, nodeKeys[0], memDB).(*backend)
	genesis.MustCommit(memDB)

	txLookUpLimit := uint64(100)
	cacheConfig := &core.CacheConfig{
		TrieCleanLimit: 256,
		TrieDirtyLimit: 256,
		TrieTimeLimit:  5 * time.Minute,
	}
	blockchain, err := core.NewBlockChain(memDB, cacheConfig, genesis.Config, b, vm.Config{}, nil, &txLookUpLimit)
	if err != nil {
		panic(err)
	}

	b.Start(blockchain, nil)

	return blockchain, b
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
		Proposal: block,
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

var emptySigner = &snr.SignerImpl{}

func makeValSet(validators []common.Address) hotstuff.ValidatorSet {
	return validator.NewSet(validators, hotstuff.RoundRobin)
}
