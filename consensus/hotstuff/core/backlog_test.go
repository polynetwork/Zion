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
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestStoreBacklog
func TestStoreBacklog(t *testing.T) {
	c, vals := singerTestCore(t, 2, 0, 1)
	sender := vals.GetByIndex(1)

	// push new view msg
	view := makeView(10, 10)
	newViewPayload := []byte{'a', 'b', 'c'}
	m := &Message{
		address: sender.Address(),
		View:    view,
		Code:    MsgTypeNewView,
		Msg:     newViewPayload,
	}
	c.storeBacklog(m)
	msg, _ := c.backlogs.Pop(sender.Address())
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}

	// push prepare msg
	preparePayload := []byte{'b', '1', 'c', 'x'}
	m = &Message{
		View:    view,
		address: sender.Address(),
		Code:    MsgTypePrepare,
		Msg:     preparePayload,
	}
	c.storeBacklog(m)
	msg, _ = c.backlogs.Pop(sender.Address())
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}

	// push pre-commit msg
	preCommitPayload := []byte{'3', '5', '2', 'd'}
	m = &Message{
		address: sender.Address(),
		View:    view,
		Code:    MsgTypePreCommit,
		Msg:     preCommitPayload,
	}
	c.storeBacklog(m)
	msg, _ = c.backlogs.Pop(sender.Address())
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}

	// push commit msg
	commitPayload := []byte{'c', 'm', 't', 'p', 'l', 'd'}
	m = &Message{
		address: sender.Address(),
		View:    view,
		Code:    MsgTypeCommit,
		Msg:     commitPayload,
	}
	c.storeBacklog(m)
	msg, _ = c.backlogs.Pop(sender.Address())
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestProcessFutureBacklog
func TestProcessFutureBacklog(t *testing.T) {
	c, vals := singerTestCore(t, 2, 0, 1)
	c.backend = &testSystemBackend{events: new(event.TypeMux)}
	sender := vals.GetByIndex(1)
	c.subscribeEvents()
	defer c.unsubscribeEvents()

	// push a future msg
	subject := []byte{'c', 'm', 't', 'p', 'l', 'd'}
	subjectPayload, _ := Encode(subject)
	m := &Message{
		address: sender.Address(),
		View:    makeView(10, 10),
		Code:    MsgTypeCommit,
		Msg:     subjectPayload,
	}
	c.storeBacklog(m)
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
		t.Log("timeout")
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestProcessBacklog
func TestProcessBacklog(t *testing.T) {
	view := makeView(1, 0)
	prepreparePayload := []byte{'p', 'p', 'd'}
	vote := common.HexToHash("0x1234567890")
	subjectPayload := vote.Bytes()
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
	c, vals := singerTestCore(t, 2, 1, 0)
	c.backend = &testSystemBackend{events: new(event.TypeMux)}
	c.subscribeEvents()
	defer c.unsubscribeEvents()

	sender := vals.GetByIndex(1)
	msg.address = sender.Address()

	c.storeBacklog(msg)
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
