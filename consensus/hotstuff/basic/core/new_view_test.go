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
	"github.com/ethereum/go-ethereum/consensus/hotstuff/basic/validator"
	"github.com/stretchr/testify/assert"
)

func newTestQC(c *core, h, r uint64) *hotstuff.QuorumCert {
	view := makeView(h, r)
	block := makeBlock(int64(h))
	N := c.valSet.Size()
	coinbase := c.valSet.GetByIndex(h % uint64(N))
	return &hotstuff.QuorumCert{
		View:     view,
		Hash:     block.Hash(),
		Proposer: coinbase.Address(),
	}
}

func newTestNewViewMsg(c *core, index int, h, r uint64, prepareQC *hotstuff.QuorumCert) *message {
	curView := makeView(h, r)
	val := c.valSet.GetByIndex(uint64(index))
	newViewMsg := &MsgNewView{
		View:      curView,
		PrepareQC: prepareQC,
	}
	payload, err := Encode(newViewMsg)
	if err != nil {
		panic(err)
	}
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

	addQC := func(index int, h, r uint64) {
		prepareQC := newTestQC(c, h-1, r)
		msg := newTestNewViewMsg(c, index, h, r, prepareQC)
		assert.NoError(t, c.current.AddNewViews(msg))
	}

	maxHeight := uint64(10)
	addQC(0, H, R)
	addQC(1, H, R)
	addQC(2, H, R)
	addQC(3, maxHeight, 0)

	highQC := c.getHighQC()
	assert.Equal(t, maxHeight-1, highQC.View.Height.Uint64())
}

func TestHandleNewView(t *testing.T) {
	N := uint64(4)
	F := uint64(1)
	H := uint64(1)
	R := uint64(0)

	sys := NewTestSystemWithBackend(N, F, H, R)
	msgList := make([]*message, N)
	for index, node := range sys.backends {
		c := node.core()
		prepareQC := newTestQC(c, H-1, R)
		c.current.SetPrepareQC(prepareQC)
		msg := newTestNewViewMsg(c, index, H, R, prepareQC)
		msgList[index] = msg
	}

	leader := sys.getLeader()
	for _, msg := range msgList {
		val := validator.New(msg.Address)
		assert.NoError(t, leader.handleNewView(msg, val))
	}

	highQC := leader.getHighQC()
	t.Log(highQC.View.Height.Uint64())
}

func TestHandleNewViewFailed(t *testing.T) {
	N := uint64(4)
	F := uint64(1)
	H := uint64(5)
	R := uint64(2)

	sys := NewTestSystemWithBackend(N, F, H, R)
	leader := sys.getLeader()
	repos := sys.getRepos()

	testMessage := func(data *MsgNewView, index int) *message {
		src := repos[index].address
		payload, _ := Encode(data)
		return &message{
			Code:    MsgTypeNewView,
			Msg:     payload,
			Address: src,
		}
	}
	testValidator := func(index int) hotstuff.Validator {
		return validator.New(repos[index].address)
	}

	testcases := []struct {
		Msg       *message
		Src       hotstuff.Validator
		ExpectErr error
	}{
		{
			Msg:       &message{Msg: []byte("123456")},
			Src:       testValidator(0),
			ExpectErr: errFailedDecodeNewView,
		},
		{
			Msg: testMessage(&MsgNewView{
				View:      makeView(H-1, R),
				PrepareQC: newTestQC(repos[0], H-1, R),
			}, 0),
			Src:       testValidator(0),
			ExpectErr: errOldMessage,
		},
		{
			Msg: testMessage(&MsgNewView{
				View:      makeView(H+1, R),
				PrepareQC: newTestQC(repos[0], H, R),
			}, 0),
			Src:       testValidator(0),
			ExpectErr: errFutureMessage,
		},
	}

	for _, v := range testcases {
		assert.Equal(t, v.ExpectErr, leader.handleNewView(v.Msg, v.Src))
	}
}
