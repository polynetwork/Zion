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
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/stretchr/testify/assert"
)

// notice: we need only 3 test case:
// 1. `newView` send quorumCert, e.g: sendPreCommit, sendCommit
// 2. `prepare` send msgNewProposal
// 3. `prepareVote` send vote, e.g: sendPrepareVote, sendPreCommitVote, sendCommitVote
func TestHandleMsg(t *testing.T) {
	N := uint64(4)
	F := uint64(1)
	H := uint64(1)
	R := uint64(0)

	sys := NewTestSystemWithBackend(N, F, H, R)

	closer := sys.Run(true)
	defer closer()

	v0 := sys.backends[0]
	r0 := v0.engine.(*core)
	_, val := v0.Validators(nil).GetByAddress(v0.Address())

	// decode new view
	{
		block := makeBlock(1)
		payload, _ := Encode(&hotstuff.QuorumCert{
			View: &hotstuff.View{
				Height: big.NewInt(0),
				Round:  big.NewInt(0),
			},
			Hash:          block.Hash(),
			Proposer:      val.Address(),
			Seal:          []byte("123"),
			CommittedSeal: [][]byte{[]byte("12"), []byte("23"), []byte("34")},
		})
		// with a matched payload. msg prepare vote should match with *hotstuff.MsgPrepareVote in normal case.
		msg := &message{
			Code:    MsgTypeNewView,
			Msg:     payload,
			Address: v0.Address(),
		}
		r0.current.SetProposal(block)
		assert.Equal(t, errFailedDecodeNewView, r0.handleCheckedMsg(msg, val))
	}

	// decode prepare failed
	{
		qc := makeBlock(4)
		lastProposer := v0.Validators(nil).GetByIndex(1)
		payload, _ := Encode(&MsgNewView{
			View: &hotstuff.View{
				Height: big.NewInt(5),
				Round:  big.NewInt(0),
			},
			PrepareQC: &hotstuff.QuorumCert{
				View: &hotstuff.View{
					Round:  big.NewInt(0),
					Height: big.NewInt(0),
				},
				Hash:          qc.Hash(),
				Proposer:      lastProposer.Address(),
				Seal:          []byte("1234"),
				CommittedSeal: [][]byte{[]byte("12"), []byte("23"), []byte("34")},
			},
		})
		msg := &message{
			Code:    MsgTypePrepare,
			Msg:     payload,
			Address: v0.Address(),
		}
		enc, err := r0.finalizeMessage(msg)
		assert.NoError(t, err)
		assert.Equal(t, errFailedDecodePrepare, r0.handleMsg(enc))
	}

	// decode prepareVote failed
	{
		block := makeBlock(1)
		payload, _ := Encode(&hotstuff.Vote{
			View: &hotstuff.View{
				Height: big.NewInt(0),
				Round:  big.NewInt(0),
			},
			Digest: block.Hash(),
		})
		// with a matched payload. msg prepare vote should match with *hotstuff.MsgPrepareVote in normal case.
		msg := &message{
			Code:    MsgTypePrepareVote,
			Msg:     payload,
			Address: v0.Address(),
		}
		r0.current.SetProposal(block)
		assert.Equal(t, errFailedDecodePrepareVote, r0.handleCheckedMsg(msg, val))
	}
}
