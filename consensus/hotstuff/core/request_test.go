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
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/prque"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/event"
)

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestCheckRequestMsg
func TestCheckRequestMsg(t *testing.T) {
	c, _ := singerTestCore(t, 1, 1, 0)
	c.pendingRequests = prque.New(nil)
	c.pendingRequestsMu = new(sync.Mutex)

	// invalid request
	err := c.checkRequestMsg(nil)
	if err != errInvalidMessage {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidMessage)
	}
	r := &Request{Proposal: nil}
	err = c.checkRequestMsg(r)
	if err != errInvalidMessage {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidMessage)
	}

	// old request
	r = &Request{
		Proposal: makeBlock(0),
	}
	err = c.checkRequestMsg(r)
	if err != errOldMessage {
		t.Errorf("error mismatch: have %v, want %v", err, errOldMessage)
	}

	// future request
	r = &Request{
		Proposal: makeBlock(2),
	}
	err = c.checkRequestMsg(r)
	if err != errFutureMessage {
		t.Errorf("error mismatch: have %v, want %v", err, errFutureMessage)
	}

	// current request
	r = &Request{
		Proposal: makeBlock(1),
	}
	err = c.checkRequestMsg(r)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run  TestStoreRequestMsg
func TestStoreRequestMsg(t *testing.T) {
	c, _ := singerTestCore(t, 4, 0, 0)
	c.backend = &testSystemBackend{events: new(event.TypeMux)}
	c.pendingRequests = prque.New(nil)
	c.pendingRequestsMu = new(sync.Mutex)
	c.subscribeEvents()
	defer c.unsubscribeEvents()

	requests := []*Request{
		{
			Proposal: makeBlock(1),
		},
		{
			Proposal: makeBlock(2),
		},
		{
			Proposal: makeBlock(3),
		},
	}

	c.storeRequestMsg(requests[1])
	c.storeRequestMsg(requests[0])
	c.storeRequestMsg(requests[2])
	if c.pendingRequests.Size() != len(requests) {
		t.Errorf("the size of pending requests mismatch: have %v, want %v", c.pendingRequests.Size(), len(requests))
	}

	c.current.height = big.NewInt(3)
	c.processPendingRequests()

	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	select {
	case ev := <-c.events.Chan():
		e, ok := ev.Data.(hotstuff.RequestEvent)
		if !ok {
			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
		}
		if e.Proposal.Number().Cmp(requests[2].Proposal.Number()) != 0 {
			t.Errorf("the number of proposal mismatch: have %v, want %v", e.Proposal.Number(), requests[2].Proposal.Number())
		}
	case <-timeout.C:
		t.Error("unexpected timeout occurs")
	}
}
