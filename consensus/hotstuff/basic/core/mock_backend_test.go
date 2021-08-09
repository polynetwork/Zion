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
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
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
	signer hotstuff.Signer
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

func (m *mockBackend) NewRequest(request hotstuff.Proposal) {
	go m.events.Post(hotstuff.RequestEvent{
		Proposal: request,
	})
}

// ==============================================
//
// define the functions that needs to be provided for Istanbul.

func (m *mockBackend) Address() common.Address {
	return m.address
}

// Peers returns all connected peers
func (m *mockBackend) Validators() hotstuff.ValidatorSet {
	return m.peers
}

func (m *mockBackend) EventMux() *event.TypeMux {
	return m.events
}

func (m *mockBackend) Broadcast(valSet hotstuff.ValidatorSet, payload []byte) error {
	if !needBroadCast {
		return nil
	}
	testLogger.Info("leader broadcast", m.Address().Hex())
	m.sentMsgs = append(m.sentMsgs, payload)
	m.sys.broadcastQueuedMessage <- innerEvent{
		Event:   hotstuff.MessageEvent{Payload: payload},
		Address: m.address,
		View:    m.core().currentView(),
	}
	return nil
}

func (m *mockBackend) Gossip(valSet hotstuff.ValidatorSet, message []byte) error {
	testLogger.Warn("not sign any data")
	return nil
}

func (m *mockBackend) Unicast(valSet hotstuff.ValidatorSet, payload []byte) error {
	m.sentMsgs = append(m.sentMsgs, payload)
	m.sys.unicastQueuedMessage <- innerEvent{
		Event:   hotstuff.MessageEvent{Payload: payload},
		Address: m.address,
		View:    m.core().currentView(),
	}
	return nil
}

func (m *mockBackend) PreCommit(proposal hotstuff.Proposal, seals [][]byte) (hotstuff.Proposal, error) {
	//qc := &hotstuff.QuorumCert{
	//	View: view,
	//	Hash: proposal.Hash(),
	//}
	//qc.Proposer = proposal.(*types.Block).Header().Coinbase
	//qc.Extra = seals[0]
	return proposal, nil
}

func (m *mockBackend) Commit(proposal hotstuff.Proposal) error {
	testLogger.Info("commit message", "address", m.Address())
	m.committedMsgs = append(m.committedMsgs, testCommittedMsgs{
		commitProposal: proposal,
	})
	// fake new head events
	// go self.events.Post(istanbul.FinalCommittedEvent{})
	return nil
}

func (m *mockBackend) Verify(proposal hotstuff.Proposal) (time.Duration, error) {
	return 0, nil
}

func (m *mockBackend) VerifyUnsealedProposal(proposal hotstuff.Proposal) (time.Duration, error) {
	return 0, nil
}

func (m *mockBackend) HasBadProposal(hash common.Hash) bool {
	return false
}

func (m *mockBackend) LastProposal() (hotstuff.Proposal, common.Address) {
	l := len(m.committedMsgs)
	if l == 0 {
		return nil, EmptyAddress
	} else {
		proposal := m.committedMsgs[l-1].commitProposal
		block := proposal.(*types.Block)
		return proposal, block.Coinbase()
	}
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
// define the mock singer

type mockSinger struct {
	address common.Address
}

func (s *mockSinger) Address() common.Address {
	return s.address
}

func (m *mockSinger) Sign(data []byte) ([]byte, error) {
	return m.address.Bytes(), nil
}

// todo
func (m *mockSinger) SigHash(header *types.Header) (hash common.Hash) {
	return header.Hash()
}

func (m *mockSinger) SignVote(p hotstuff.Proposal) ([]byte, error) {
	return nil, nil
}

func (s *mockSinger) Recover(h *types.Header) (common.Address, error) {
	return h.Coinbase, nil
}

func (s *mockSinger) PrepareExtra(header *types.Header, valSet hotstuff.ValidatorSet) ([]byte, error) {
	return nil, nil
}

func (s *mockSinger) SealBeforeCommit(h *types.Header) error {
	return nil
}

func (s *mockSinger) SealAfterCommit(h *types.Header, committedSeals [][]byte) error {
	return nil
}

func (m *mockSinger) VerifyHeader(header *types.Header, valSet hotstuff.ValidatorSet, seal bool) error {
	return nil
}

func (s *mockSinger) VerifyQC(qc *hotstuff.QuorumCert, valSet hotstuff.ValidatorSet) error {
	return nil
}

func (s *mockSinger) CheckQCParticipant(qc *hotstuff.QuorumCert, signer common.Address) error {
	return nil
}

func (m *mockSinger) CheckSignature(valSet hotstuff.ValidatorSet, data []byte, signature []byte) (common.Address, error) {
	return common.BytesToAddress(signature), nil
}

func (s *mockSinger) WrapCommittedSeal(hash common.Hash) []byte {
	return hash.Bytes()
}

// ==============================================
//
// define the struct that need to be provided for integration tests.

type innerEvent struct {
	Event   hotstuff.MessageEvent
	Address common.Address
	View    *hotstuff.View
}

type testSystem struct {
	backends []*mockBackend

	broadcastQueuedMessage chan innerEvent //hotstuff.MessageEvent
	unicastQueuedMessage   chan innerEvent //hotstuff.MessageEvent
	quit                   chan struct{}
}

func (s *testSystem) getLeader() *core {
	for _, v := range s.backends {
		if v.engine.IsProposer() {
			return v.core()
		}
	}
	return nil
}

func (s *testSystem) getLeaderByRound(lastProposer common.Address, round *big.Int) *core {
	valset := s.backends[0].peers.Copy()
	valset.CalcProposer(lastProposer, round.Uint64())
	proposer := valset.GetProposer().Address()
	for _, v := range s.backends {
		core := v.core()
		if core.Address() == proposer {
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

		broadcastQueuedMessage: make(chan innerEvent, 128),
		unicastQueuedMessage:   make(chan innerEvent, 128),
		quit:                   make(chan struct{}),
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
	config := hotstuff.DefaultBasicConfig

	for i := uint64(0); i < n; i++ {
		vset := validator.NewSet(addrs, hotstuff.RoundRobin)
		backend := sys.NewBackend(i)
		backend.peers = vset
		backend.address = vset.GetByIndex(i).Address()

		signer := &mockSinger{address: backend.address}
		backend.signer = signer

		core := New(backend, config, signer, vset).(*core)
		core.current = newRoundState(&hotstuff.View{
			Height: new(big.Int).SetUint64(h),
			Round:  new(big.Int).SetUint64(r),
		}, vset, nil)
		core.valSet = vset
		core.logger = testLogger
		core.backend = backend
		core.signer = signer
		core.validateFn = func(data []byte, sig []byte) (common.Address, error) {
			return signer.CheckSignature(vset, data, sig)
		}

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
		case m := <-t.broadcastQueuedMessage:
			for _, backend := range t.backends {
				go backend.EventMux().Post(m.Event)
			}
			testLogger.Info("broadcast", "leader", m.Address.Hex(), "height", m.View.Height, "round", m.View.Round)
		case m := <-t.unicastQueuedMessage:
			leader := t.getLeader()
			go leader.sendEvent(m.Event)
			testLogger.Info("unicast", "Address", m.Address.Hex(), "leader", leader.Address().Hex(), "height", m.View.Height, "round", m.View.Round)
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
		signer: nil,
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
