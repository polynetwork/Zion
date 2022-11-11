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
//	"testing"
//
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
//	"github.com/stretchr/testify/assert"
//)
//
//func newTestNewViewMsg(sender common.Address, h, r int, prepareQC *QuorumCert) *Message {
//	view := makeView(h, r)
//	newViewMsg := &MsgNewView{
//		View:      view,
//		PrepareQC: prepareQC,
//	}
//	payload, err := Encode(newViewMsg)
//	if err != nil {
//		panic(err)
//	}
//	msg := &Message{
//		View:    view,
//		Code:    MsgTypeNewView,
//		Msg:     payload,
//		Address: sender,
//	}
//	return msg
//}
//
//// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestMaxView
//func TestMaxView(t *testing.T) {
//	N, H, R := 4, 1, 0
//
//	sys := NewTestSystemWithBackend(N, H, R)
//	c := sys.backends[0].engine
//
//	addQC := func(index int, h, r int) {
//		prepareQC := newTestQCWithoutExtra(c, h-1, r)
//		sender := c.valSet.GetByIndex(uint64(index))
//		msg := newTestNewViewMsg(sender.Address(), h, r, prepareQC)
//		assert.NoError(t, c.current.AddNewViews(msg))
//	}
//
//	maxHeight := 10
//	addQC(0, H, R)
//	addQC(1, H, R)
//	addQC(2, H, R)
//	addQC(3, maxHeight, 0)
//
//	highQC, err := c.getHighQC()
//	assert.NoError(t, err)
//	assert.Equal(t, maxHeight-1, int(highQC.HeightU64()))
//}
//
//// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestHandleNewView
//func TestHandleNewView(t *testing.T) {
//	N, H, R := 4, 1, 0
//
//	sys := NewTestSystemWithBackend(N, H, R)
//	msgList := make([]*Message, N)
//	for index, node := range sys.backends {
//		c := node.engine
//		prepareQC := newTestQCWithoutExtra(c, H-1, R)
//		c.current.SetPrepareQC(prepareQC)
//		sender := c.valSet.GetByIndex(uint64(index))
//		msg := newTestNewViewMsg(sender.Address(), H, R, prepareQC)
//		assert.NoError(t, c.current.AddNewViews(msg))
//		msgList[index] = msg
//	}
//
//	leader := sys.getLeader()
//	for _, msg := range msgList {
//		val := validator.New(msg.Address)
//		assert.NoError(t, leader.handleNewView(msg, val))
//	}
//
//	highQC, err := leader.getHighQC()
//	assert.NoError(t, err)
//	t.Log(highQC.HeightU64())
//}
//
//// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestHandleNewViewFailed
//func TestHandleNewViewFailed(t *testing.T) {
//	N, H, R := 4, 5, 2
//
//	sys := NewTestSystemWithBackend(N, H, R)
//	leader := sys.getLeader()
//	repo := sys.getRepos()[0]
//	val := validator.New(repo.Address())
//
//	// decode failed
//	err := leader.handleNewView(&Message{Msg: []byte("123456")}, val)
//	assert.Equal(t, errFailedDecodeNewView, err)
//
//	// too old message
//	qc := newTestQCWithExtra(t, sys, H-2)
//	msg := newTestNewViewMsg(repo.Address(), H-1, 0, qc)
//	err = leader.handleNewView(msg, val)
//	assert.Equal(t, errOldMessage, err)
//
//	qc = newTestQCWithExtra(t, sys, H-1)
//	msg = newTestNewViewMsg(repo.Address(), H, R-1, qc)
//	err = leader.handleNewView(msg, val)
//	assert.Equal(t, errOldMessage, err)
//
//	// future message
//	qc = newTestQCWithExtra(t, sys, H-1)
//	msg = newTestNewViewMsg(repo.Address(), H, R+1, qc)
//	err = leader.handleNewView(msg, val)
//	assert.Equal(t, errFutureMessage, err)
//
//	qc = newTestQCWithExtra(t, sys, H)
//	msg = newTestNewViewMsg(repo.Address(), H+1, 0, qc)
//	err = leader.handleNewView(msg, val)
//	assert.Equal(t, errFutureMessage, err)
//
//	// error leader
//	qc = newTestQCWithExtra(t, sys, H-1)
//	msg = newTestNewViewMsg(repo.Address(), H, R, qc)
//	err = repo.handleNewView(msg, val)
//	assert.Equal(t, errNotToProposer, err)
//}
