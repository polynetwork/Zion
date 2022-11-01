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
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func newTestPrepareMsg(t *testing.T, s *testSystem, sender common.Address, h, r int) (*MsgPrepare, *Message) {
	view := makeView(h, r)
	highQC := newTestQCWithExtra(t, s, h-1)
	proposal := makeBlockWithParentHash(h, highQC.hash)
	prepare := &MsgPrepare{
		View:     view,
		Proposal: proposal,
		HighQC:   highQC,
	}
	payload, err := Encode(prepare)
	if err != nil {
		t.Error(err)
	}

	msg := &Message{
		Code:    MsgTypePrepare,
		View:    view,
		Msg:     payload,
		Address: sender,
	}
	return prepare, msg
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestHandlePrepare
func TestHandlePrepare(t *testing.T) {
	N, H, R := 4, 5, 1

	sys := NewTestSystemWithBackend(N, H, R)
	leader := sys.getLeader()
	val := leader.valSet.GetProposer()
	prepare, msg := newTestPrepareMsg(t, sys, leader.Address(), H, R)
	for _, backend := range sys.backends {
		core := backend.engine
		core.current.SetPreCommittedQC(prepare.HighQC)
		err := core.handlePrepare(msg, val)
		assert.NoError(t, err)
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestPrepareFailed
func TestPrepareFailed(t *testing.T) {
	N, H, R := 4, 5, 1

	sys := NewTestSystemWithBackend(N, H, R)
	leader := sys.getLeader()
	val := leader.valSet.GetProposer()
	prepare, msg := newTestPrepareMsg(t, sys, leader.Address(), H, R)

	// too old message
	{
		for _, backend := range sys.backends {
			core := backend.engine
			core.current.height = big.NewInt(int64(H + 1))
			core.current.SetPreCommittedQC(prepare.HighQC)
			err := core.handlePrepare(msg, val)
			assert.Equal(t, errOldMessage, err)
		}

		for _, backend := range sys.backends {
			core := backend.engine
			core.current.height = big.NewInt(int64(H))
			core.current.round = big.NewInt(int64(R + 1))
			core.current.SetPreCommittedQC(prepare.HighQC)
			err := core.handlePrepare(msg, val)
			assert.Equal(t, errOldMessage, err)
		}
	}

	// future message
	{
		for _, backend := range sys.backends {
			core := backend.engine
			core.current.height = big.NewInt(int64(H - 1))
			core.current.SetPreCommittedQC(prepare.HighQC)
			err := core.handlePrepare(msg, val)
			assert.Equal(t, errFutureMessage, err)
		}

		for _, backend := range sys.backends {
			core := backend.engine
			core.current.height = big.NewInt(int64(H))
			core.current.round = big.NewInt(int64(R - 1))
			core.current.SetPreCommittedQC(prepare.HighQC)
			err := core.handlePrepare(msg, val)
			assert.Equal(t, errFutureMessage, err)
		}
	}
}

//type testcase struct {
//	Sys       *testSystem
//	Msg       *Message
//	Leader    hotstuff.Validator
//	ExpectErr error
//}
//
//}
//	// errMsgOld
//	func() *testcase {
//		sys := NewTestSystemWithBackend(N, H, R)
//		var data *MsgPrepare
//		for _, backend := range sys.backends {
//			core := backend.engine
//			data = newPrepareMsg(core)
//			core.current.height = new(big.Int).SetUint64(uint64(H + 1))
//		}
//		msg := newP2PMsg(data)
//		leader := validator.New(sys.getLeader().Address())
//		return &testcase{
//			Sys:       sys,
//			Msg:       msg,
//			Leader:    leader,
//			ExpectErr: errOldMessage,
//		}
//	}(),
//
//	// errFutureMsg
//		func() *testcase {
//			sys := NewTestSystemWithBackend(N, H, R)
//			var data *MsgPrepare
//			for _, backend := range sys.backends {
//				core := backend.engine
//				data = newPrepareMsg(core)
//				core.current.round = new(big.Int).SetUint64(uint64(R - 1))
//			}
//			msg := newP2PMsg(data)
//			leader := validator.New(sys.getLeader().Address())
//			return &testcase{
//				Sys:       sys,
//				Msg:       msg,
//				Leader:    leader,
//				ExpectErr: errFutureMessage,
//			}
//		}(),
//
//	// errNotFromProposer
//		func() *testcase {
//			sys := NewTestSystemWithBackend(N, H, R)
//			var data *MsgPrepare
//			for _, backend := range sys.backends {
//				core := backend.engine
//				data = newPrepareMsg(core)
//			}
//			msg := newP2PMsg(data)
//			wrongLeader := validator.New(sys.getRepos()[0].Address())
//			return &testcase{
//				Sys:       sys,
//				Msg:       msg,
//				Leader:    wrongLeader,
//				ExpectErr: errNotFromProposer,
//			}
//		}(),
//
//	// errExtend
//		func() *testcase {
//			sys := NewTestSystemWithBackend(N, H, R)
//			leader := sys.getLeader()
//			val := validator.New(leader.Address())
//			var data *MsgPrepare
//			for _, backend := range sys.backends {
//				core := backend.engine
//				data = newPrepareMsg(core)
//				core.current.SetPreCommittedQC(data.HighQC)
//			}
//			// msg.proposal.parentHash not equal to the field of `lockedQC.Hash`
//			data.HighQC.hash = common.HexToHash("0x124")
//			msg := newP2PMsg(data)
//			return &testcase{
//				Sys:       sys,
//				Msg:       msg,
//				Leader:    val,
//				ExpectErr: errExtend,
//			}
//		}(),
//
//	// errSafeNode
//		func() *testcase {
//			sys := NewTestSystemWithBackend(N, H, R)
//			leader := sys.getLeader()
//			val := validator.New(leader.Address())
//			var data *MsgPrepare
//			for _, backend := range sys.backends {
//				core := backend.engine
//				data = newPrepareMsg(core)
//				// safety is false, and liveness false:
//				// msg.proposal is not extend lockedQC
//				// msg.highQC.view is smaller than lockedQC.view
//				// or just set lockedQC is nil
//				//core.current.SetPreCommittedQC(data.HighQC)
//			}
//			msg := newP2PMsg(data)
//			return &testcase{
//				Sys:       sys,
//				Msg:       msg,
//				Leader:    val,
//				ExpectErr: errSafeNode,
//			}
//		}(),
//}
