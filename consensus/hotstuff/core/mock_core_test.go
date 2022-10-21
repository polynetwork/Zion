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
//	"time"
//
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/consensus/hotstuff"
//)
//
//// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/basic/core -run TestNewRound
//func TestNewRound(t *testing.T) {
//	N := uint64(1)
//	F := uint64(1)
//	H := uint64(1)
//	R := uint64(0)
//
//	needBroadCast = true
//	sys := NewTestSystemWithBackend(N, F, H, R)
//
//	// prepare genesis qc
//	lastLeader := sys.getLeaderByRound(EmptyAddress, common.Big0)
//	genesisBlock, _ := newProposalAndQC(lastLeader, 0, 0)
//	for _, v := range sys.backends {
//		v.committedMsgs = append(v.committedMsgs, testCommittedMsgs{
//			commitProposal: genesisBlock,
//		})
//		v.core().current = nil
//	}
//
//	close := sys.Run(true)
//	defer close()
//
//	block := makeBlockWithParentHash(1, genesisBlock.Hash())
//	for _, backend := range sys.backends {
//		go backend.EventMux().Post(hotstuff.RequestEvent{
//			Proposal: block,
//		})
//	}
//
//	<-time.After(2 * time.Second)
//
//	for _, v := range sys.backends {
//		if len(v.committedMsgs) > 1 {
//			block := v.committedMsgs[1].commitProposal
//			t.Logf("proposer %s committed block %d hash %s", v.address.Hex(), block.Number().Uint64(), block.Hash().Hex())
//		}
//	}
//}
