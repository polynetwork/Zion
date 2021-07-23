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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func TestMessageSetWithNewView(t *testing.T) {
	valSet := newTestValidatorSet(4)

	ms := newMessageSet(valSet)

	view := &hotstuff.View{
		Round:  new(big.Int),
		Height: new(big.Int),
	}
	pp := &MsgNewView{
		View: view,
	}
	payload, err := Encode(pp)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	msg := &message{
		Code:    MsgTypeNewView,
		Msg:     payload,
		Address: valSet.GetProposer().Address(),
	}

	if err = ms.Add(msg); err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	if err = ms.Add(msg); err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	if ms.Size() != 1 {
		t.Errorf("the size of message set mismatch: have %v, want 1", ms.Size())
	}
}

func TestMessageSetVote(t *testing.T) {
	valSet := newTestValidatorSet(4)

	ms := newMessageSet(valSet)

	view := &hotstuff.View{
		Round:  new(big.Int),
		Height: new(big.Int),
	}

	sub := &Vote{
		View:   view,
		Digest: common.HexToHash("1234567890"),
	}

	payload, err := Encode(sub)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	msg := &message{
		Code:    MsgTypePrepareVote,
		Msg:     payload,
		Address: valSet.GetProposer().Address(),
	}

	if err := ms.Add(msg); err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	if err := ms.Add(msg); err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	if ms.Size() != 1 {
		t.Errorf("the size of message set mismatch: have %v, want 1", ms.Size())
	}
}
