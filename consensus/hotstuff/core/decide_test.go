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

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestHandleCommitVote
func TestHandleCommitVote(t *testing.T) {
	N, H, R := 4, 5, 1

	newVoteMsg := func(hash common.Hash, sender common.Address, h, r int) *Message {
		return &Message{
			Code:    MsgTypeCommitVote,
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
				core.current.node.node = node
				core.current.node.temp = node
				core.current.Lock(&QuorumCert{node: node.Hash()})
				msg := newVoteMsg(node.Hash(), core.Address(), H, R)
				msg.PayloadNoSig()
				sig, _ := v.engine.signer.SignHash(msg.hash)
				msg.Signature = sig
				votes[validator.New(msg.address)] = msg
			}
			sys.getLeader().current.state = StateLocked
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
			assert.Equal(t, v.ExpectErr, leader.handleCommitVote(vote))
		}
		if v.ExpectErr == nil {
			assert.Equal(t, StateCommitted, leader.current.State())
			assert.Equal(t, N, leader.current.CommitVoteSize())
		}
	}
}
