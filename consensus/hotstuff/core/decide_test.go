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
//
//import (
//	"math/big"
//	"testing"
//
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/consensus/hotstuff"
//	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestHandleCommitVote(t *testing.T) {
//	N := uint64(4)
//	F := uint64(1)
//	H := uint64(5)
//	R := uint64(1)
//
//	newVote := func(c *core, hash common.Hash) *Vote {
//		view := c.currentView()
//		return &Vote{
//			View:   view,
//			Digest: hash,
//		}
//	}
//	newVoteMsg := func(vote *Vote) *hotstuff.Message {
//		payload, _ := Encode(vote)
//		return &hotstuff.Message{
//			Code: MsgTypeCommitVote,
//			Msg:  payload,
//		}
//	}
//
//	type testcase struct {
//		Sys       *testSystem
//		Votes     map[hotstuff.Validator]*hotstuff.Message
//		ExpectErr error
//	}
//
//	testcases := []*testcase{
//
//		// normal case
//		func() *testcase {
//			sys := NewTestSystemWithBackend(N, F, H, R)
//			proposal := makeBlock(int64(H))
//			votes := make(map[hotstuff.Validator]*hotstuff.Message)
//			for _, v := range sys.backends {
//				core := v.core()
//				core.current.SetProposal(proposal)
//				core.current.SetPreCommittedQC(&QuorumCert{Hash: proposal.Hash()})
//
//				vote := newVote(core, proposal.Hash())
//				msg := newVoteMsg(vote)
//				msg.Address = core.Address()
//				val := validator.New(msg.Address)
//
//				votes[val] = msg
//			}
//			return &testcase{
//				Sys:       sys,
//				Votes:     votes,
//				ExpectErr: nil,
//			}
//		}(),
//
//		// errOldMessage
//		func() *testcase {
//			sys := NewTestSystemWithBackend(N, F, H, R)
//			proposal := makeBlock(int64(H))
//			votes := make(map[hotstuff.Validator]*hotstuff.Message)
//			for _, v := range sys.backends {
//				core := v.core()
//				core.current.SetProposal(proposal)
//				core.current.SetPreCommittedQC(&QuorumCert{Hash: proposal.Hash()})
//
//				vote := newVote(core, proposal.Hash())
//				vote.View.Height = new(big.Int).SetUint64(H - 1)
//
//				msg := newVoteMsg(vote)
//				msg.Address = core.Address()
//				val := validator.New(msg.Address)
//
//				votes[val] = msg
//			}
//			return &testcase{
//				Sys:       sys,
//				Votes:     votes,
//				ExpectErr: errOldMessage,
//			}
//		}(),
//
//		// errFutureMessage
//		func() *testcase {
//			sys := NewTestSystemWithBackend(N, F, H, R)
//			proposal := makeBlock(int64(H))
//			votes := make(map[hotstuff.Validator]*hotstuff.Message)
//			for _, v := range sys.backends {
//				core := v.core()
//				core.current.SetProposal(proposal)
//				core.current.SetPreCommittedQC(&QuorumCert{Hash: proposal.Hash()})
//
//				vote := newVote(core, proposal.Hash())
//				vote.View.Round = new(big.Int).SetUint64(R + 1)
//
//				msg := newVoteMsg(vote)
//				msg.Address = core.Address()
//				val := validator.New(msg.Address)
//
//				votes[val] = msg
//			}
//			return &testcase{
//				Sys:       sys,
//				Votes:     votes,
//				ExpectErr: errFutureMessage,
//			}
//		}(),
//
//		// errInconsistentVote
//		func() *testcase {
//			sys := NewTestSystemWithBackend(N, F, H, R)
//			proposal := makeBlock(int64(H))
//			votes := make(map[hotstuff.Validator]*hotstuff.Message)
//			for _, v := range sys.backends {
//				core := v.core()
//				core.current.SetProposal(proposal)
//				core.current.SetPreCommittedQC(&QuorumCert{Hash: proposal.Hash()})
//
//				vote := newVote(core, proposal.Hash())
//				vote.Digest = common.HexToHash("0x1234")
//				msg := newVoteMsg(vote)
//				msg.Address = core.Address()
//				val := validator.New(msg.Address)
//
//				votes[val] = msg
//			}
//			return &testcase{
//				Sys:       sys,
//				Votes:     votes,
//				ExpectErr: errInconsistentVote,
//			}
//		}(),
//
//		// errInvalidDigest
//		func() *testcase {
//			sys := NewTestSystemWithBackend(N, F, H, R)
//			proposal := makeBlock(int64(H))
//			votes := make(map[hotstuff.Validator]*hotstuff.Message)
//			for _, v := range sys.backends {
//				core := v.core()
//				core.current.SetProposal(proposal)
//				core.current.SetPreCommittedQC(&QuorumCert{Hash: common.HexToHash("0x124")})
//
//				vote := newVote(core, proposal.Hash())
//				msg := newVoteMsg(vote)
//				msg.Address = core.Address()
//				val := validator.New(msg.Address)
//
//				votes[val] = msg
//			}
//			return &testcase{
//				Sys:       sys,
//				Votes:     votes,
//				ExpectErr: errInvalidDigest,
//			}
//		}(),
//	}
//
//	for _, v := range testcases {
//		leader := v.Sys.getLeader()
//		for src, vote := range v.Votes {
//			assert.Equal(t, v.ExpectErr, leader.handleCommitVote(vote, src))
//		}
//		if v.ExpectErr == nil {
//			assert.Equal(t, StateCommitted, leader.current.State())
//			assert.Equal(t, int(N), leader.current.CommitVoteSize())
//		}
//	}
//}
