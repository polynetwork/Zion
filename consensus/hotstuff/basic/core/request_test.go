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
	"time"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
)

func TestCheckRequestMsg(t *testing.T) {
	c := &core{
		current: newRoundState(&hotstuff.View{
			Height: big.NewInt(1),
			Round:  big.NewInt(0),
		}, newTestValidatorSet(4), nil),
	}

	// invalid request
	err := c.requests.checkRequest(nil, nil)
	if err != errInvalidMessage {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidMessage)
	}
	r := &hotstuff.Request{
		Proposal: nil,
	}
	err = c.requests.checkRequest(nil, r)
	if err != errInvalidMessage {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidMessage)
	}

	// old request
	r = &hotstuff.Request{
		Proposal: makeBlock(0),
	}
	err = c.requests.checkRequest(c.currentView(), r)
	if err != errOldMessage {
		t.Errorf("error mismatch: have %v, want %v", err, errOldMessage)
	}

	// future request
	r = &hotstuff.Request{
		Proposal: makeBlock(2),
	}
	err = c.requests.checkRequest(c.currentView(), r)
	if err != errFutureMessage {
		t.Errorf("error mismatch: have %v, want %v", err, errFutureMessage)
	}

	// current request
	r = &hotstuff.Request{
		Proposal: makeBlock(1),
	}
	err = c.requests.checkRequest(c.currentView(), r)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
}

func TestStoreRequestMsg(t *testing.T) {
	backend := &mockBackend{
		events: new(event.TypeMux),
	}
	height := big.NewInt(1)
	reqNum := 3
	c := &core{
		logger:  log.New("backend", "test", "id", 0),
		backend: backend,
		current: newRoundState(&hotstuff.View{
			Height: height,
			Round:  big.NewInt(0),
		}, newTestValidatorSet(4), nil),
		requests: newRequestSet(),
	}
	c.subscribeEvents()
	defer c.unsubscribeEvents()

	go c.handleEvents()

	requests := make([]hotstuff.Request, reqNum)
	for i := 1; i <= reqNum; i++ {
		requests[i-1] = hotstuff.Request{
			Proposal: makeBlock(int64(i)),
		}
	}

	for _, req := range requests {
		c.sendEvent(hotstuff.RequestEvent{
			Proposal: req.Proposal,
		})
	}

	<-time.After(1 * time.Second)
	t.Log("request size after sendEvent", c.requests.Size())

	req := c.requests.GetRequest(c.currentView())
	assert.NotNil(t, req)
	assert.Equal(t, req.Proposal.Number().Uint64(), c.current.Height().Uint64())
	assert.Equal(t, c.requests.Size(), int(reqNum-1))
}
