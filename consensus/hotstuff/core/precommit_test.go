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
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/stretchr/testify/assert"
)

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestHandlePrepareVote
func TestHandlePrepareVote(t *testing.T) {
	N, H, R := 4, 5, 1

	newVoteMsg := func(hash common.Hash, sender common.Address, h, r int) *Message {
		return &Message{
			Code:    MsgTypePrepareVote,
			Msg:     hash.Bytes(),
			View:    makeView(h, r),
			address: sender,
		}
	}

	type testcase struct {
		Sys       *testSystem
		Votes     map[hotstuff.Validator]*Message
		ExpectErr error
	}

	parentNode := common.HexToHash("0x123")
	testcases := []*testcase{
		func() *testcase {
			sys := NewTestSystemWithBackend(N, H, R)
			proposal := makeBlock(H)
			node := NewNode(parentNode, proposal)
			votes := make(map[hotstuff.Validator]*Message)
			for _, v := range sys.backends {
				core := v.engine
				core.current.SetNode(node)
				msg := newVoteMsg(node.Hash(), core.Address(), H, R)
				msg.address = core.Address()
				val := validator.New(msg.address)
				msg.PayloadNoSig()
				sig, _ := v.engine.signer.SignHash(msg.hash)
				msg.Signature = sig
				votes[val] = msg
			}
			sys.getLeader().current.state = StateHighQC
			sys.Run(false)

			return &testcase{
				Sys:       sys,
				Votes:     votes,
				ExpectErr: nil,
			}
		}(),
	}

	for _, v := range testcases {
		leader := v.Sys.getLeader()
		for _, vote := range v.Votes {
			assert.Equal(t, v.ExpectErr, leader.handlePrepareVote(vote))
		}
		if v.ExpectErr == nil {
			assert.Equal(t, StatePrepared, leader.current.State())
			assert.Equal(t, N, leader.current.PrepareVoteSize())
		}
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestHandlePreCommit
func TestHandlePreCommit(t *testing.T) {
	N, H, R := 4, 5, 1

	parentNode := common.HexToHash("0x123")
	newPreCommitMsg := func(backend *testSystemBackend) (*Subject, *Message) {
		coreView := backend.engine.currentView()
		h := int(coreView.Height.Uint64())
		r := int(coreView.Round.Uint64())
		proposal := makeBlock(h)
		node := NewNode(parentNode, proposal)
		qc := newTestQCWithExtra(t, backend.sys, node.Hash(), MsgTypePrepareVote, h, r)
		sub := NewSubject(node, qc)
		payload, _ := Encode(qc)
		leader := backend.sys.getLeader()
		msg := &Message{
			address: leader.Address(),
			Code:    MsgTypePreCommit,
			Msg:     payload,
			View:    coreView,
		}
		msg.PayloadNoSig()
		sig, _ := leader.signer.SignHash(msg.hash)
		msg.Signature = sig
		return sub, msg
	}

	type testcase struct {
		Sys       *testSystem
		Msg       *Message
		Leader    hotstuff.Validator
		ExpectErr error
	}
	testcases := []*testcase{
		func() *testcase {
			sys := NewTestSystemWithBackend(N, H, R)
			sys.Run(false)
			leader := sys.getLeader()
			val := validator.New(leader.Address())
			msgs := make([]*Message, N)
			for i, backend := range sys.backends {
				core := backend.engine
				subject, msg := newPreCommitMsg(backend)
				core.current.SetNode(subject.Node)
				msgs[i] = msg
			}
			return &testcase{
				Sys:       sys,
				Msg:       msgs[0],
				Leader:    val,
				ExpectErr: nil,
			}
		}(),
	}

	for _, c := range testcases {
		for _, backend := range c.Sys.backends {
			core := backend.engine
			assert.Equal(t, c.ExpectErr, core.handlePreCommit(c.Msg))
		}
	}
}
