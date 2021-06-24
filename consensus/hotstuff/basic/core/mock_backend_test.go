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

package core

import (
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/basic/validator"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	elog "github.com/ethereum/go-ethereum/log"
)

var (
	needBroadCast = false
	testLogger    = elog.New()
)

type mockBackend struct {
	id  uint64
	sys *testSystem

	engine CoreEngine
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

func (m *mockBackend) core() *core {
	return m.engine.(*core)
}

// ==============================================
//
// define the functions that needs to be provided for Istanbul.

func (m *mockBackend) Address() common.Address {
	return m.address
}

// Peers returns all connected peers
func (m *mockBackend) Validators(proposal hotstuff.Proposal) hotstuff.ValidatorSet {
	return m.peers
}

func (m *mockBackend) EventMux() *event.TypeMux {
	return m.events
}

func (m *mockBackend) Send(message []byte, target common.Address) error {
	testLogger.Info("enqueuing a message...", "address", m.Address())
	m.sentMsgs = append(m.sentMsgs, message)
	m.sys.queuedMessage <- hotstuff.MessageEvent{
		Payload: message,
	}
	return nil
}

func (m *mockBackend) Broadcast(valSet hotstuff.ValidatorSet, message []byte) error {
	if !needBroadCast {
		return nil
	}
	testLogger.Info("enqueuing a message...", "address", m.Address())
	m.sentMsgs = append(m.sentMsgs, message)
	m.sys.queuedMessage <- hotstuff.MessageEvent{
		Payload: message,
	}
	return nil
}

func (m *mockBackend) Gossip(valSet hotstuff.ValidatorSet, message []byte) error {
	testLogger.Warn("not sign any data")
	return nil
}

func (m *mockBackend) Unicast(valSet hotstuff.ValidatorSet, payload []byte) error {
	return nil
}

func (m *mockBackend) PreCommit(view *hotstuff.View, proposal hotstuff.Proposal, seals [][]byte) (hotstuff.Proposal, *hotstuff.QuorumCert, error) {
	qc := &hotstuff.QuorumCert{
		View: view,
		Hash: proposal.Hash(),
	}
	return proposal, qc, nil
}

func (m *mockBackend) Commit(proposal hotstuff.Proposal) error {
	return nil
}

func (m *mockBackend) Verify(proposal hotstuff.Proposal) (time.Duration, error) {
	return 0, nil
}

func (m *mockBackend) VerifyUnsealedProposal(proposal hotstuff.Proposal) (time.Duration, error) {
	return 0, nil
}

func (s *mockBackend) VerifyQuorumCert(qc *hotstuff.QuorumCert) error {
	return nil
}

func (m *mockBackend) Sign(data []byte) ([]byte, error) {
	testLogger.Info("returning current backend address so that CheckValidatorSignature returns the same value")
	return m.address.Bytes(), nil
}

// SignTx signs transaction data with backend's private key
func (m *mockBackend) SignTx(tx *types.Transaction, signer types.Signer) (*types.Transaction, error) {
	return nil, nil
}

func (m *mockBackend) CheckSignature([]byte, common.Address, []byte) error {
	return nil
}

// todo: delete after test
func (m *mockBackend) CheckValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	return common.BytesToAddress(sig), nil
}

func (m *mockBackend) Hash(b interface{}) common.Hash {
	return common.HexToHash("Test")
}

func (m *mockBackend) NewRequest(request hotstuff.Proposal) {
	go m.events.Post(hotstuff.RequestEvent{
		Proposal: request,
	})
}

func (m *mockBackend) HasBadProposal(hash common.Hash) bool {
	return false
}

func (m *mockBackend) LastProposal() (hotstuff.Proposal, common.Address) {
	l := len(m.committedMsgs)
	if l > 0 {
		return m.committedMsgs[l-1].commitProposal, common.Address{}
	}
	return makeBlock(0), common.Address{}
}

func (m *mockBackend) CurrentProposer() (*big.Int, common.Address) {
	return nil, common.Address{}
}

// Only block height 5 will return true
func (m *mockBackend) HasProposal(hash common.Hash, number *big.Int) bool {
	return number.Cmp(big.NewInt(5)) == 0
}

func (m *mockBackend) GetProposer(number uint64) common.Address {
	return common.Address{}
}

func (m *mockBackend) ParentValidators(proposal hotstuff.Proposal) hotstuff.ValidatorSet {
	return m.peers
}

func (m *mockBackend) Close() error {
	return nil
}

// ==============================================
//
// define the struct that need to be provided for integration tests.

type testSystem struct {
	backends []*mockBackend

	queuedMessage chan hotstuff.MessageEvent
	quit          chan struct{}
}

func (s *testSystem) getLeader() *core {
	for _, v := range s.backends {
		if v.engine.IsProposer() {
			return v.core()
		}
	}
	return nil
}

func (s *testSystem) getRepos() []*core {
	list := make([]*core, 0)
	for _, v := range s.backends {
		if !v.engine.IsProposer() {
			list = append(list, v.core())
		}
	}
	return list
}

func newTestSystem(n uint64) *testSystem {
	testLogger.SetHandler(elog.StdoutHandler)
	return &testSystem{
		backends: make([]*mockBackend, n),

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

func newTestValidatorSet(n int) hotstuff.ValidatorSet {
	return validator.NewSet(generateValidators(n), hotstuff.RoundRobin)
}

// FIXME: int64 is needed for N and F
func NewTestSystemWithBackend(n, f, h, r uint64) *testSystem {
	testLogger.SetHandler(elog.StdoutHandler)

	addrs := generateValidators(int(n))
	sys := newTestSystem(n)
	config := hotstuff.DefaultConfig

	for i := uint64(0); i < n; i++ {
		vset := validator.NewSet(addrs, hotstuff.RoundRobin)
		backend := sys.NewBackend(i)
		backend.peers = vset
		backend.address = vset.GetByIndex(i).Address()

		core := New(backend, config, vset).(*core)
		core.current = newRoundState(&hotstuff.View{
			Height: new(big.Int).SetUint64(h),
			Round:  new(big.Int).SetUint64(r),
		}, vset, nil)
		core.valSet = vset
		core.logger = testLogger
		core.backend = backend
		core.validateFn = backend.CheckValidatorSignature

		backend.engine = core

		core.subscribeEvents()
		defer core.unsubscribeEvents()
	}

	//backend := &testSystemBackend{
	//	events: new(event.TypeMux),
	//	peers:  vset,
	//}
	//c := &core{
	//	logger:     log.New("backend", "test", "id", 0),
	//	backlogs:   make(map[common.Address]*prque.Prque),
	//	backlogsMu: new(sync.Mutex),
	//	valSet:     vset,
	//	backend:    backend,
	//	state:      State(msg.Code),
	//	current: newRoundState(&istanbul.View{
	//		Sequence: big.NewInt(1),
	//		Round:    big.NewInt(0),
	//	}, newTestValidatorSet(4), common.Hash{}, nil, nil, nil),
	//}
	//c.subscribeEvents()
	//defer c.unsubscribeEvents()
	//
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
			b.engine.Start() // start hotstuff core
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

func (t *testSystem) NewBackend(id uint64) *mockBackend {
	// assume always success
	ethDB := rawdb.NewMemoryDatabase()
	backend := &mockBackend{
		id:     id,
		sys:    t,
		events: new(event.TypeMux),
		db:     ethDB,
	}

	t.backends[id] = backend
	return backend
}

// ==============================================
//
// helper functions.

func getPublicKeyAddress(privateKey *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(privateKey.PublicKey)
}
