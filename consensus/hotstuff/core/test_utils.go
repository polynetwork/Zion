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

package core

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/signer"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	elog "github.com/ethereum/go-ethereum/log"
)

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

func newTestValidatorSet(n int) (hotstuff.ValidatorSet, []*ecdsa.PrivateKey) {
	mkeys := make(map[common.Address]*ecdsa.PrivateKey)
	addrs := make([]common.Address, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(privateKey.PublicKey)
		addrs[i] = addr
		mkeys[addr] = privateKey
	}
	vset := validator.NewSet(addrs, hotstuff.RoundRobin)
	keys := make(Keys, n)
	for i, addr := range vset.AddressList() {
		keys[i] = mkeys[addr]
	}
	return vset, keys
}

func makeBlock(number int) *types.Block {
	header := &types.Header{
		Difficulty: big.NewInt(0),
		Number:     big.NewInt(int64(number)),
		GasLimit:   0,
		GasUsed:    0,
		Time:       0,
	}
	block := &types.Block{}
	return block.WithSeal(header)
}

func newTestProposal() hotstuff.Proposal {
	return makeBlock(1)
}

func makeBlockWithParentHash(number int, parentHash common.Hash) *types.Block {
	header := &types.Header{
		Difficulty: big.NewInt(0),
		Number:     big.NewInt(int64(number)),
		GasLimit:   0,
		GasUsed:    0,
		Time:       0,
	}
	if parentHash != common.EmptyHash {
		header.ParentHash = parentHash
	}
	block := &types.Block{}
	return block.WithSeal(header)
}

// ==============================================
//
// define the mock backend
//
// ==============================================

var testLogger = elog.New()

type testSystemBackend struct {
	id  int
	sys *testSystem

	engine *core
	peers  hotstuff.ValidatorSet
	events *event.TypeMux

	committedMsgs []testCommittedMsgs
	sentMsgs      [][]byte // store the message when Send is called by core

	address common.Address
	db      ethdb.Database
}

type testCommittedMsgs struct {
	commitProposal hotstuff.Proposal
	committedSeals [][]byte
}

func (ts *testSystemBackend) Address() common.Address {
	return ts.address
}

// Peers returns all connected peers
func (ts *testSystemBackend) Validators(hash common.Hash, mining bool) hotstuff.ValidatorSet {
	return ts.peers
}

func (ts *testSystemBackend) EventMux() *event.TypeMux {
	return ts.events
}

func (ts *testSystemBackend) Broadcast(valSet hotstuff.ValidatorSet, message []byte) error {
	//return nil
	testLogger.Info("enqueuing a message...", "address", ts.Address())
	ts.sentMsgs = append(ts.sentMsgs, message)
	ts.sys.queuedMessage <- hotstuff.MessageEvent{
		//Code:    code,
		Payload: message,
	}
	return nil
}

func (ts *testSystemBackend) Gossip(valSet hotstuff.ValidatorSet, message []byte) error {
	return nil
	testLogger.Warn("not sign any data")
	return nil
}

func (ts *testSystemBackend) Unicast(valSet hotstuff.ValidatorSet, message []byte) error {
	return nil
	testLogger.Info("enqueuing a message...", "address", ts.Address())
	ts.sentMsgs = append(ts.sentMsgs, message)
	ts.sys.queuedMessage <- hotstuff.MessageEvent{
		//Code:    code,
		Payload: message,
	}
	return nil
}

func (ts *testSystemBackend) PreCommit(proposal hotstuff.Proposal, seals [][]byte) (hotstuff.Proposal, error) {
	// todo:
	return nil, nil
}

func (ts *testSystemBackend) Commit(proposal hotstuff.Proposal) error {
	testLogger.Info("commit message", "address", ts.Address())
	ts.committedMsgs = append(ts.committedMsgs, testCommittedMsgs{
		commitProposal: proposal,
		//committedSeals: seals,
	})

	// fake new head events
	go ts.events.Post(hotstuff.FinalCommittedEvent{})
	return nil
}

func (ts *testSystemBackend) Verify(proposal hotstuff.Proposal) (time.Duration, error) {
	return 0, nil
}

func (ts *testSystemBackend) VerifyUnsealedProposal(proposal hotstuff.Proposal) (time.Duration, error) {
	return 0, nil
}

func (ts *testSystemBackend) ValidateBlock(block *types.Block) error { return nil }
func (ts *testSystemBackend) HasBadProposal(hash common.Hash) bool   { return false }

func (ts *testSystemBackend) LastProposal() (hotstuff.Proposal, common.Address) {
	l := len(ts.committedMsgs)
	if l > 0 {
		return ts.committedMsgs[l-1].commitProposal, common.Address{}
	}
	return makeBlock(0), common.Address{}
}

// Only block height 5 will return true
func (ts *testSystemBackend) HasPropsal(hash common.Hash, number *big.Int) bool {
	return number.Cmp(big.NewInt(5)) == 0
}

func (ts *testSystemBackend) Close() error                  { return nil }
func (ts *testSystemBackend) ReStart()                      {}
func (ts *testSystemBackend) CheckPoint(height uint64) bool { return false }

// ==============================================
//
// define the struct that need to be provided for integration tests.

type testSystem struct {
	backends []*testSystemBackend

	queuedMessage chan hotstuff.MessageEvent
	quit          chan struct{}
}

func newTestSystem(n int) *testSystem {
	testLogger.SetHandler(elog.StdoutHandler)
	return &testSystem{
		backends:      make([]*testSystemBackend, n),
		queuedMessage: make(chan hotstuff.MessageEvent),
		quit:          make(chan struct{}),
	}
}

func generateValidators(n int) []common.Address {
	vals := make([]common.Address, 0)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		vals = append(vals, crypto.PubkeyToAddress(privateKey.PublicKey))
	}
	return vals
}

func NewTestSystemWithBackend(n, h, r int) *testSystem {
	testLogger.SetHandler(elog.StdoutHandler)

	vset, keys := newTestValidatorSet(n)
	sys := newTestSystem(n)
	config := hotstuff.DefaultBasicConfig

	for i := 0; i < n; i++ {
		backend := sys.NewBackend(i)
		backend.peers = vset
		backend.address = vset.GetByIndex(uint64(i)).Address()

		core := New(backend, config, signer.NewSigner(keys[i]))
		core.current = newRoundState(makeView(h, r), vset)
		core.valSet = vset
		core.logger = testLogger
		core.validateFn = nil

		backend.engine = core
	}

	return sys
}

// listen will consume messages from queue and deliver a message to core
func (t *testSystem) listen() {
	for {
		select {
		case <-t.quit:
			return
		case queuedMessage := <-t.queuedMessage:
			testLogger.Info("consuming a queue message...")
			for _, backend := range t.backends {
				go backend.EventMux().Post(queuedMessage)
			}
		}
	}
}

// Run will start system components based on given flag, and returns a closer
// function that caller can control lifecycle
//
// Given a true for core if you want to initialize core engine.
func (t *testSystem) Run(core bool) func() {
	for _, b := range t.backends {
		if core {
			b.engine.Start(nil) // start Istanbul core
		}
	}

	go t.listen()
	closer := func() { t.stop(core) }
	return closer
}

func (t *testSystem) stop(core bool) {
	close(t.quit)

	for _, b := range t.backends {
		if core {
			b.engine.Stop()
		}
	}
}

func (t *testSystem) NewBackend(id int) *testSystemBackend {
	// assume always success
	ethDB := rawdb.NewMemoryDatabase()
	backend := &testSystemBackend{
		id:     id,
		sys:    t,
		events: new(event.TypeMux),
		db:     ethDB,
	}

	t.backends[id] = backend
	return backend
}

func (s *testSystem) getLeader() *core {
	for _, v := range s.backends {
		if v.engine.IsProposer() {
			return v.engine
		}
	}
	return nil
}

func (s *testSystem) getRepos() []*core {
	list := make([]*core, 0)
	for _, v := range s.backends {
		if !v.engine.IsProposer() {
			list = append(list, v.engine)
		}
	}
	return list
}

// ==============================================
//
// mock signer
//
// ==============================================

type testSigner struct {
	pk      *ecdsa.PrivateKey
	address common.Address
}

func (ts *testSigner) Address() common.Address                         { return ts.address }
func (ts *testSigner) Sign(data []byte) ([]byte, error)                { return common.EmptyHash.Bytes(), nil }
func (ts *testSigner) SigHash(header *types.Header) (hash common.Hash) { return common.EmptyHash }
func (ts *testSigner) SignHash(hash common.Hash) ([]byte, error)       { return common.EmptyHash.Bytes(), nil }
func (ts *testSigner) SignTx(tx *types.Transaction, signer types.Signer) (*types.Transaction, error) {
	return tx, nil
}
func (ts *testSigner) Recover(h *types.Header) (common.Address, *types.HotstuffExtra, error) {
	return h.Coinbase, nil, nil
}
func (ts *testSigner) SealBeforeCommit(h *types.Header) error                         { return nil }
func (ts *testSigner) SealAfterCommit(h *types.Header, committedSeals [][]byte) error { return nil }
func (ts *testSigner) VerifyHeader(header *types.Header, valSet hotstuff.ValidatorSet, seal bool) (*types.HotstuffExtra, error) {
	return nil, nil
}
func (ts *testSigner) VerifyQC(qc hotstuff.QC, valSet hotstuff.ValidatorSet) error    { return nil }
func (ts *testSigner) CheckQCParticipant(qc hotstuff.QC, signer common.Address) error { return nil }
func (ts *testSigner) CheckSignature(valSet hotstuff.ValidatorSet, data []byte, signature []byte) (common.Address, error) {
	return common.EmptyAddress, nil
}
func (ts *testSigner) VerifyHash(valSet hotstuff.ValidatorSet, hash common.Hash, sig []byte) error {
	return nil
}
func (ts *testSigner) VerifyCommittedSeal(valSet hotstuff.ValidatorSet, hash common.Hash, committedSeals [][]byte) error {
	return nil
}

// ==============================================
//
// helper functions.

func getPublicKeyAddress(privateKey *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(privateKey.PublicKey)
}

func singerTestCore(t *testing.T, n int, height, round int64) (*core, hotstuff.ValidatorSet) {
	if n < 1 {
		t.Error("invalid participants")
	}

	vals, keys := newTestValidatorSet(n)

	c := &core{
		logger: log.New("backend", "test", "id", 0),
		valSet: vals,
		current: newRoundState(&View{
			Height: big.NewInt(height),
			Round:  big.NewInt(round),
		}, vals),
		signer:   signer.NewSigner(keys[0]),
		backlogs: newBackLog(),
	}

	return c, vals
}

func singerAddress() common.Address {
	num := rand.Int()
	return common.HexToAddress(fmt.Sprintf("0x%d", num))
}

func makeView(h, r int) *View {
	return &View{
		Height: big.NewInt(int64(h)),
		Round:  big.NewInt(int64(r)),
	}
}

func newTestQCWithoutExtra(c *core, h, r int) *QuorumCert {
	view := makeView(h, r)
	block := makeBlock(h)
	N := c.valSet.Size()
	coinbase := c.valSet.GetByIndex(uint64(h % N))
	return &QuorumCert{
		view:     view,
		hash:     block.Hash(),
		proposer: coinbase.Address(),
	}
}

func newTestQCWithExtra(t *testing.T, s *testSystem, h int) *QuorumCert {
	view := makeView(h, 0)
	block := makeBlock(h)
	hash := block.Hash()
	vset := s.backends[0].engine.valSet
	N := vset.Size()
	coinbase := vset.GetByIndex(uint64(h % N))

	leader := s.getLeader()
	seal, _ := leader.signer.SignHash(hash, true)
	committedSeal := make([][]byte, N-1)
	for i, v := range s.getRepos() {
		sig, err := v.signer.SignHash(hash, true)
		if err != nil {
			t.Errorf("sign block hash failed, err: %v", err)
		}
		committedSeal[i] = sig
	}
	//extra, err := types.GenerateExtraWithSignature(0, 1000000000, vset.AddressList(), seal, committedSeal)
	//if err != nil {
	//	t.Errorf("generate extra with signatures failed, err: %v", err)
	//}
	return &QuorumCert{
		view:          view,
		hash:          hash,
		proposer:      coinbase.Address(),
		seal:          seal,
		committedSeal: committedSeal,
	}
}
