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

	"github.com/stretchr/testify/assert"
)

// notice: we need only 3 test case:
// 1. `newView` send quorumCert, e.g: sendPreCommit, sendCommit
// 2. `prepare` send msgNewProposal
// 3. `prepareVote` send vote, e.g: sendPrepareVote, sendPreCommitVote, sendCommitVote
func TestHandleMsg(t *testing.T) {
	N := uint64(4)
	F := uint64(1)
	H := uint64(5)
	R := uint64(0)

	sys := NewTestSystemWithBackend(N, F, H, R)

	closer := sys.Run(true)
	defer closer()

	v0 := sys.backends[0]
	r0 := v0.core()
	_, val := v0.Validators().GetByAddress(v0.Address())

	// decode new view
	{
		block := makeBlock(1)
		payload, _ := Encode(&MsgNewView{
			View:      makeView(H, R),
			PrepareQC: newTestQC(r0, H-1, R),
		})
		// with a matched payload. msg prepare vote should match with *hotstuff.MsgPrepareVote in normal case.
		msg := &message{
			Code:    MsgTypeNewView,
			Msg:     payload,
			Address: v0.Address(),
		}
		r0.current.SetProposal(block)
		assert.NoError(t, r0.handleCheckedMsg(msg, val))
	}

	// decode prepare failed
	{
		payload, _ := Encode(&MsgPrepare{
			View:     makeView(H, R),
			Proposal: makeBlock(int64(H - 1)),
			HighQC:   newTestQC(r0, H, R),
		})
		msg := &message{
			Code:    MsgTypePrepare,
			Msg:     payload,
			Address: v0.Address(),
		}
		enc, err := r0.finalizeMessage(msg)
		assert.NoError(t, err)
		assert.Equal(t, errExtend, r0.handleMsg(enc))
	}

	// decode prepareVote failed
	{
		block := makeBlock(int64(H))
		payload, _ := Encode(&Vote{
			View:   makeView(H, R),
			Digest: block.Hash(),
		})
		// with a matched payload. msg prepare vote should match with *hotstuff.MsgPrepareVote in normal case.
		msg := &message{
			Code:    MsgTypePrepareVote,
			Msg:     payload,
			Address: v0.Address(),
		}
		assert.Equal(t, errInconsistentVote, r0.handleCheckedMsg(msg, val))
	}
}
