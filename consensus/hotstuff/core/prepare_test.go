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
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func newTestPrepareMsg(t *testing.T, s *testSystem, sender common.Address, parentBlock common.Hash, h, r int) (*Subject, *Message) {
	view := makeView(h, r)
	highQC := newTestQCWithExtra(t, s, common.HexToHash("0x123"), MsgTypePrepareVote, h-1, r)
	proposal := makeBlockWithParentHash(h, parentBlock)
	node := NewNode(highQC.node, proposal)
	prepare := NewSubject(node, highQC)
	payload, err := Encode(prepare)
	if err != nil {
		t.Error(err)
	}

	msg := &Message{
		Code:    MsgTypePrepare,
		View:    view,
		Msg:     payload,
		address: sender,
	}
	return prepare, msg
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestHandlePrepare
func TestHandlePrepare(t *testing.T) {
	N, H, R := 4, 5, 1

	sys := NewTestSystemWithBackend(N, H, R)
	sys.Run(false)
	leader := sys.getLeader()
	parent := makeBlock(H - 1)
	prepare, msg := newTestPrepareMsg(t, sys, leader.Address(), parent.Hash(), H, R)
	for _, backend := range sys.backends {
		core := backend.engine
		core.current.lastChainedBlock = parent
		core.current.prepareQC = prepare.QC
		err := core.handlePrepare(msg)
		assert.NoError(t, err)
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestPrepareFailed
func TestPrepareFailed(t *testing.T) {
	N, H, R := 4, 5, 1

	sys := NewTestSystemWithBackend(N, H, R)
	leader := sys.getLeader()
	parent := makeBlock(H - 1)
	prepare, msg := newTestPrepareMsg(t, sys, leader.Address(), parent.Hash(), H, R)

	// too old message
	{
		for _, backend := range sys.backends {
			core := backend.engine
			core.current.lastChainedBlock = parent
			core.current.height = big.NewInt(int64(H + 1))
			core.current.prepareQC = prepare.QC
			err := core.handlePrepare(msg)
			assert.Equal(t, errOldMessage, err)
		}
	}
	{
		for _, backend := range sys.backends {
			core := backend.engine
			core.current.lastChainedBlock = parent
			core.current.height = big.NewInt(int64(H))
			core.current.round = big.NewInt(int64(R + 1))
			core.current.prepareQC = prepare.QC
			err := core.handlePrepare(msg)
			assert.Equal(t, errOldMessage, err)
		}
	}

	// future message
	{
		for _, backend := range sys.backends {
			core := backend.engine
			core.current.lastChainedBlock = parent
			core.current.height = big.NewInt(int64(H - 1))
			core.current.prepareQC = prepare.QC
			err := core.handlePrepare(msg)
			assert.Equal(t, errFutureMessage, err)
		}
	}
	{
		for _, backend := range sys.backends {
			core := backend.engine
			core.current.lastChainedBlock = parent
			core.current.height = big.NewInt(int64(H))
			core.current.round = big.NewInt(int64(R - 1))
			core.current.prepareQC = prepare.QC
			err := core.handlePrepare(msg)
			assert.Equal(t, errFutureMessage, err)
		}
	}
}
