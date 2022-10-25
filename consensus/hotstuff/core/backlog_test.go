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
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/consensus/hotstuff/signer"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestStoreBacklog
func TestStoreBacklog(t *testing.T) {
	vals, keys := newTestValidatorSet(2)
	signer := signer.NewSigner(keys[0])
	sender := vals.GetByIndex(1)

	c := &core{
		logger: log.New("backend", "test", "id", 0),
		valSet: vals,
		current: newRoundState(&View{
			Round:  big.NewInt(1),
			Height: big.NewInt(0),
		}, vals, nil),
		signer:   signer,
		backlogs: newBackLog(),
	}

	// push new view msg
	view := &View{
		Round:  big.NewInt(10),
		Height: big.NewInt(10),
	}
	newView := &MsgNewView{View: view}
	newViewPayload, _ := Encode(newView)
	m := &Message{
		View:    view,
		Address: sender.Address(),
		Code:    MsgTypeNewView,
		Msg:     newViewPayload,
	}

	c.storeBacklog(m, sender)
	msg, _ := c.backlogs.Pop(sender.Address())
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}

	// push prepare msg
	prepare := &MsgPrepare{View: view}
	subject, _ := Encode(prepare)
	m = &Message{
		View:    view,
		Address: sender.Address(),
		Code:    MsgTypePrepare,
		Msg:     subject,
	}
	c.storeBacklog(m, sender)
	msg, _ = c.backlogs.Pop(sender.Address())
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}

	// push pre-commit msg
	m = &Message{
		View:    view,
		Address: sender.Address(),
		Code:    MsgTypePreCommit,
		Msg:     subject,
	}
	c.storeBacklog(m, sender)
	msg, _ = c.backlogs.Pop(sender.Address())
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}

	// push commit msg
	m = &Message{
		View:    view,
		Address: sender.Address(),
		Code:    MsgTypeCommit,
		Msg:     subject,
	}
	c.storeBacklog(m, sender)
	msg, _ = c.backlogs.Pop(sender.Address())
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestProcessFutureBacklog
func TestProcessFutureBacklog(t *testing.T) {
	backend := &testSystemBackend{
		events: new(event.TypeMux),
	}
	vals, pks := newTestValidatorSet(2)
	signer := signer.NewSigner(pks[0])
	sender := vals.GetByIndex(1)

	c := &core{
		logger:   log.New("backend", "test", "id", 0),
		valSet:   vals,
		signer:   signer,
		backlogs: newBackLog(),
		backend:  backend,
		current: newRoundState(&View{
			Height: big.NewInt(1),
			Round:  big.NewInt(0),
		}, vals, nil),
	}
	c.subscribeEvents()
	defer c.unsubscribeEvents()

	v := &View{
		Round:  big.NewInt(10),
		Height: big.NewInt(10),
	}
	// push a future msg
	subject := &MsgPreCommit{
		View: v,
	}
	subjectPayload, _ := Encode(subject)
	m := &Message{
		View:    v,
		Address: sender.Address(),
		Code:    MsgTypeCommit,
		Msg:     subjectPayload,
	}
	c.storeBacklog(m, sender)
	c.processBacklog()

	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	select {
	case e, ok := <-c.events.Chan():
		if !ok {
			return
		}
		t.Errorf("unexpected events comes: %v", e)
	case <-timeout.C:
		// success
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestProcessBacklog
func TestProcessBacklog(t *testing.T) {
	view := &View{
		Height: big.NewInt(1),
		Round:  big.NewInt(0),
	}
	prepare := &MsgPrepare{
		View:     view,
		Proposal: makeBlock(1),
	}
	prepreparePayload, _ := Encode(prepare)

	subject := &Vote{
		View:   view,
		Digest: common.HexToHash("0x1234567890"),
	}
	subjectPayload, _ := Encode(subject)

	msgs := []*Message{
		{
			View: view,
			Code: MsgTypePrepare,
			Msg:  prepreparePayload,
		},
		{
			View: view,
			Code: MsgTypePreCommit,
			Msg:  subjectPayload,
		},
		{
			View: view,
			Code: MsgTypeCommit,
			Msg:  subjectPayload,
		},
		{
			View: view,
			Code: MsgTypeDecide,
			Msg:  subjectPayload,
		},
	}
	for i := 0; i < len(msgs); i++ {
		testProcessBacklog(t, msgs[i])
	}
}

func testProcessBacklog(t *testing.T, msg *Message) {
	vals, keys := newTestValidatorSet(2)
	signer := signer.NewSigner(keys[0])
	sender := vals.GetByIndex(1)

	msg.Address = sender.Address()
	backend := &testSystemBackend{
		events: new(event.TypeMux),
		peers:  vals,
	}
	c := &core{
		logger:   log.New("backend", "test", "id", 0),
		backlogs: newBackLog(),
		valSet:   vals,
		signer:   signer,
		backend:  backend,
		current: newRoundState(&View{
			Height: big.NewInt(1),
			Round:  big.NewInt(0),
		}, vals, nil),
	}
	c.subscribeEvents()
	defer c.unsubscribeEvents()

	c.storeBacklog(msg, sender)
	c.processBacklog()

	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	select {
	case ev := <-c.events.Chan():
		e, ok := ev.Data.(backlogEvent)
		if !ok {
			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
		}
		if e.msg.Code != msg.Code {
			t.Errorf("message code mismatch: have %v, want %v", e.msg.Code, msg.Code)
		}
		// success
	case <-timeout.C:
		t.Error("unexpected timeout occurs")
	}
}
