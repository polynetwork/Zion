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
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common/prque"
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

//func TestStoreRequestMsg(t *testing.T) {
//	backend := &mockBackend{
//		events: new(event.TypeMux),
//	}
//	height := big.NewInt(1)
//	reqNum := 3
//	c := &core{
//		logger:  log.New("backend", "test", "id", 0),
//		backend: backend,
//		current: newRoundState(&view{
//			Height: height,
//			Round:  big.NewInt(0),
//		}, newTestValidatorSet(4), nil),
//		requests: newRequestSet(),
//	}
//	c.subscribeEvents()
//	defer c.unsubscribeEvents()
//
//	go c.handleEvents()
//
//	requests := make([]hotstuff.Request, reqNum)
//	for i := 1; i <= reqNum; i++ {
//		requests[i-1] = hotstuff.Request{
//			Proposal: makeBlock(int64(i)),
//		}
//	}
//
//	for _, req := range requests {
//		c.sendEvent(hotstuff.RequestEvent{
//			Proposal: req.Proposal,
//		})
//	}
//
//	<-time.After(1 * time.Second)
//	t.Log("request size after sendEvent", c.requests.Size())
//
//	req := c.requests.GetRequest(c.currentView())
//	assert.NotNil(t, req)
//	assert.Equal(t, req.Proposal.Number().Uint64(), c.current.Height().Uint64())
//	assert.Equal(t, c.requests.Size(), int(reqNum-1))
//}
