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
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
)

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/event-driven/core -run TestNewRound
func TestNewRound(t *testing.T) {
	N := uint64(1)
	F := uint64(1)
	H := uint64(1)
	R := uint64(0)

	needBroadCast = true
	sys := NewTestSystemWithBackend(N, F, H, R)

	close := sys.Run(true)
	defer close()

	genesisBlock, _ := sys.backends[0].LastProposal()
	vs := sys.backends[0].Validators()
	blocks := makeContinueBlocks(vs, genesisBlock.(*types.Block), 5)[1:]
	for _, backend := range sys.backends {
		go func() {
			for _, block := range blocks {
				_ = backend.EventMux().Post(hotstuff.RequestEvent{
					Proposal: block,
				})
			}
		}()
	}

	<-time.After(1 * time.Second)

	for _, v := range sys.backends {
		if len(v.committedMsgs) > 1 {
			block := v.committedMsgs[1].commitProposal
			t.Logf("proposer %s committed block %d hash %s", v.address.Hex(), block.Number().Uint64(), block.Hash().Hex())
		}
	}
}
