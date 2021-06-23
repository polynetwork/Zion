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
	"testing"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/stretchr/testify/assert"
)

func newViewMsg(c *core, index, h, r int64) *message {
	view := makeView(h, r)
	lastProposalHeight := h - 1
	block := makeBlock(lastProposalHeight)
	N := c.valSet.Size()
	coinbase := c.valSet.GetByIndex(uint64(lastProposalHeight) % uint64(N))
	val := c.valSet.GetByIndex(uint64(index))
	qc := &hotstuff.QuorumCert{
		View:     view,
		Hash:     block.Hash(),
		Proposer: coinbase.Address(),
	}
	payload, _ := Encode(qc)
	msg := &message{
		Code:    MsgTypeNewView,
		Msg:     payload,
		Address: val.Address(),
	}
	return msg
}

func TestMaxView(t *testing.T) {
	N := uint64(4)
	F := uint64(1)
	H := uint64(1)
	R := uint64(0)

	sys := NewTestSystemWithBackend(N, F, H, R)
	c := sys.backends[0].core()

	addQC := func(index, h, r int64) {
		msg := newViewMsg(c, index, h, r)
		assert.NoError(t, c.current.AddNewViews(msg))
	}

	maxHeight := int64(10)
	addQC(0, int64(H), int64(R))
	addQC(1, int64(H), int64(R))
	addQC(2, int64(H), int64(R))
	addQC(3, maxHeight, 0)

	highQC := c.getHighQC()
	assert.Equal(t, uint64(maxHeight), highQC.View.Height.Uint64())
}

//
//func newTestNewView(v *hotstuff.View) *hotstuff.QuorumCert {
//	block := makeBlock(1)
//	return &hotstuff.QuorumCert{
//		View:     v,
//		Hash: block.Hash(),
//		Proposer:
//	}
//}
//
//func TestHandlePreprepare(t *testing.T) {
//	N := uint64(4) // replica 0 is the proposer, it will send messages to others
//	F := uint64(1) // F does not affect tests
//	H := uint64(1)
//	R := uint64(0)
//
//	testCases := []struct {
//		system          *testSystem
//		expectedRequest *hotstuff.QuorumCert
//		expectedErr     error
//		existingBlock   bool
//	}{
//		{
//			// normal case
//			func() *testSystem {
//				sys := NewTestSystemWithBackend(N, F, H, R)
//
//				for _, backend := range sys.backends {
//					c := backend.engine.(*core)
//					c.valSet = backend.peers
//				}
//				return sys
//			}(),
//			newTestNewView(),
//			nil,
//			false,
//		},
//		//
//	}
//
//OUTER:
//	for _, test := range testCases {
//		test.system.Run(false)
//
//		v0 := test.system.backends[0]
//		r0 := v0.engine.(*core)
//
//		curView := r0.currentView()
//
//		preprepare := &hotstuff.QuorumCert{
//			View:     curView,
//			Proposal: test.expectedRequest,
//		}
//
//		for i, v := range test.system.backends {
//			// i == 0 is primary backend, it is responsible for send PRE-PREPARE messages to others.
//			if i == 0 {
//				continue
//			}
//
//			c := v.engine.(*core)
//
//			m, _ := Encode(preprepare)
//			_, val := r0.valSet.GetByAddress(v0.Address())
//			// run each backends and verify handlePreprepare function.
//			if err := c.handlePreprepare(&message{
//				Code:    msgPreprepare,
//				Msg:     m,
//				Address: v0.Address(),
//			}, val); err != nil {
//				if err != test.expectedErr {
//					t.Errorf("error mismatch: have %v, want %v", err, test.expectedErr)
//				}
//				continue OUTER
//			}
//
//			if c.state != StatePreprepared {
//				t.Errorf("state mismatch: have %v, want %v", c.state, StatePreprepared)
//			}
//
//			if !test.existingBlock && !reflect.DeepEqual(c.current.Subject().View, curView) {
//				t.Errorf("view mismatch: have %v, want %v", c.current.Subject().View, curView)
//			}
//
//			// verify prepare messages
//			decodedMsg := new(message)
//			err := decodedMsg.FromPayload(v.sentMsgs[0], nil)
//			if err != nil {
//				t.Errorf("error mismatch: have %v, want nil", err)
//			}
//
//			expectedCode := MsgTypeNewView
//			if test.existingBlock {
//				expectedCode = MsgTypeCommit
//			}
//			if decodedMsg.Code != expectedCode {
//				t.Errorf("message code mismatch: have %v, want %v", decodedMsg.Code, expectedCode)
//			}
//
//			var subject *hotstuff.Vote
//			err = decodedMsg.Decode(&subject)
//			if err != nil {
//				t.Errorf("error mismatch: have %v, want nil", err)
//			}
//			if !test.existingBlock && !reflect.DeepEqual(subject, c.current.Subject()) {
//				t.Errorf("subject mismatch: have %v, want %v", subject, c.current.Subject())
//			}
//		}
//	}
//}

//{
//	// future message
//	func() *testSystem {
//		sys := NewTestSystemWithBackend(N, F)
//
//		for i, backend := range sys.backends {
//			c := backend.engine.(*core)
//			c.valSet = backend.peers
//			if i != 0 {
//				c.state = StateAcceptRequest
//				// hack: force set subject that future message can be simulated
//				c.current = newTestRoundState(
//					&istanbul.View{
//						Round:    big.NewInt(0),
//						Sequence: big.NewInt(0),
//					},
//					c.valSet,
//				)
//
//			} else {
//				c.current.SetSequence(big.NewInt(10))
//			}
//		}
//		return sys
//	}(),
//	makeBlock(1),
//	errFutureMessage,
//	false,
//},
//{
//	// non-proposer
//	func() *testSystem {
//		sys := NewTestSystemWithBackend(N, F)
//
//		// force remove replica 0, let replica 1 be the proposer
//		sys.backends = sys.backends[1:]
//
//		for i, backend := range sys.backends {
//			c := backend.engine.(*core)
//			c.valSet = backend.peers
//			if i != 0 {
//				// replica 0 is the proposer
//				c.state = StatePreprepared
//			}
//		}
//		return sys
//	}(),
//	makeBlock(1),
//	errNotFromProposer,
//	false,
//},
//{
//	// errOldMessage
//	func() *testSystem {
//		sys := NewTestSystemWithBackend(N, F)
//
//		for i, backend := range sys.backends {
//			c := backend.engine.(*core)
//			c.valSet = backend.peers
//			if i != 0 {
//				c.state = StatePreprepared
//				c.current.SetSequence(big.NewInt(10))
//				c.current.SetRound(big.NewInt(10))
//			}
//		}
//		return sys
//	}(),
//	makeBlock(1),
//	errOldMessage,
//	false,
//},

//func TestHandlePreprepareWithLock(t *testing.T) {
//	N := uint64(4) // replica 0 is the proposer, it will send messages to others
//	F := uint64(1) // F does not affect tests
//	proposal := newTestProposal()
//	mismatchProposal := makeBlock(10)
//	newSystem := func() *testSystem {
//		sys := NewTestSystemWithBackend(N, F)
//
//		for i, backend := range sys.backends {
//			c := backend.engine.(*core)
//			c.valSet = backend.peers
//			if i != 0 {
//				c.state = StateAcceptRequest
//			}
//			c.roundChangeSet = newRoundChangeSet(c.valSet)
//		}
//		return sys
//	}
//
//	testCases := []struct {
//		system       *testSystem
//		proposal     istanbul.Proposal
//		lockProposal istanbul.Proposal
//	}{
//		{
//			newSystem(),
//			proposal,
//			proposal,
//		},
//		{
//			newSystem(),
//			proposal,
//			mismatchProposal,
//		},
//	}
//
//	for _, test := range testCases {
//		test.system.Run(false)
//		v0 := test.system.backends[0]
//		r0 := v0.engine.(*core)
//		curView := r0.currentView()
//		preprepare := &istanbul.Preprepare{
//			View:     curView,
//			Proposal: test.proposal,
//		}
//		lockPreprepare := &istanbul.Preprepare{
//			View:     curView,
//			Proposal: test.lockProposal,
//		}
//
//		for i, v := range test.system.backends {
//			// i == 0 is primary backend, it is responsible for send PRE-PREPARE messages to others.
//			if i == 0 {
//				continue
//			}
//
//			c := v.engine.(*core)
//			c.current.SetPreprepare(lockPreprepare)
//			c.current.LockHash()
//			m, _ := Encode(preprepare)
//			_, val := r0.valSet.GetByAddress(v0.Address())
//			if err := c.handlePreprepare(&message{
//				Code:    msgPreprepare,
//				Msg:     m,
//				Address: v0.Address(),
//			}, val); err != nil {
//				t.Errorf("error mismatch: have %v, want nil", err)
//			}
//			if test.proposal == test.lockProposal {
//				if c.state != StatePrepared {
//					t.Errorf("state mismatch: have %v, want %v", c.state, StatePreprepared)
//				}
//				if !reflect.DeepEqual(curView, c.currentView()) {
//					t.Errorf("view mismatch: have %v, want %v", c.currentView(), curView)
//				}
//			} else {
//				// Should stay at StateAcceptRequest
//				if c.state != StateAcceptRequest {
//					t.Errorf("state mismatch: have %v, want %v", c.state, StateAcceptRequest)
//				}
//				// Should have triggered a round change
//				expectedView := &istanbul.View{
//					Sequence: curView.Sequence,
//					Round:    big.NewInt(1),
//				}
//				if !reflect.DeepEqual(expectedView, c.currentView()) {
//					t.Errorf("view mismatch: have %v, want %v", c.currentView(), expectedView)
//				}
//			}
//		}
//	}
//}
