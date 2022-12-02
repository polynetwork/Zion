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
//	"math/big"
//	"reflect"
//	"testing"
//	"time"
//
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/consensus/hotstuff"
//	"github.com/ethereum/go-ethereum/core/types"
//	elog "github.com/ethereum/go-ethereum/log"
//)
//
//func makeBlock(number int64) *types.Block {
//	return makeBlockWithParentHash(number, EmptyHash)
//}
//
//func makeBlockWithParentHash(number int64, parentHash common.Hash) *types.Block {
//	header := &types.Header{
//		Difficulty: big.NewInt(0),
//		Number:     big.NewInt(number),
//		GasLimit:   0,
//		GasUsed:    0,
//		Time:       0,
//	}
//	if parentHash != EmptyHash {
//		header.ParentHash = parentHash
//	}
//	block := &types.Block{}
//	return block.WithSeal(header)
//}
//
//func makeAddress(i int) common.Address {
//	num := new(big.Int).SetUint64(uint64(i))
//	return common.BytesToAddress(num.Bytes())
//}
//
//func makeHash(i int) common.Hash {
//	num := new(big.Int).SetUint64(uint64(i))
//	return common.BytesToHash(num.Bytes())
//}
//
//func makeView(h, r uint64) *View {
//	return &View{
//		Height: new(big.Int).SetUint64(h),
//		Round:  new(big.Int).SetUint64(r),
//	}
//}
//
//func TestNewRequest(t *testing.T) {
//	testLogger.SetHandler(elog.StdoutHandler)
//
//	N := uint64(4)
//	F := uint64(1)
//	H := uint64(1)
//	R := uint64(0)
//
//	sys := NewTestSystemWithBackend(N, F, H, R)
//
//	close := sys.Run(true)
//	defer close()
//
//	request1 := makeBlock(1)
//	sys.backends[0].NewRequest(request1)
//
//	<-time.After(1 * time.Second)
//
//	request2 := makeBlock(2)
//	sys.backends[0].NewRequest(request2)
//
//	<-time.After(1 * time.Second)
//
//	for _, backend := range sys.backends {
//		if len(backend.committedMsgs) != 2 {
//			t.Errorf("the number of executed requests mismatch: have %v, want 2", len(backend.committedMsgs))
//		}
//		if !reflect.DeepEqual(request1.Number(), backend.committedMsgs[0].commitProposal.Number()) {
//			t.Errorf("the number of requests mismatch: have %v, want %v", request1.Number(), backend.committedMsgs[0].commitProposal.Number())
//		}
//		if !reflect.DeepEqual(request2.Number(), backend.committedMsgs[1].commitProposal.Number()) {
//			t.Errorf("the number of requests mismatch: have %v, want %v", request2.Number(), backend.committedMsgs[1].commitProposal.Number())
//		}
//	}
//}
//
//func TestQuorumSize(t *testing.T) {
//	N, H, R := 4, 1, 0
//
//	sys := NewTestSystemWithBackend(N, H, R)c
//	backend := sys.bakends[0]
//	c := backend.engine.(*core)
//
//	valSet := c.valSet
//	for i := 1; i <= 2; i++ {
//		valSet.AddValidator(makeAddress(i))
//		if 2*c.Q() <= (valSet.Size()+valSet.F()) || 2*c.Q() > (valSet.Size()+valSet.F()+2) {
//			t.Errorf("quorumSize constraint failed, expected value (2*QuorumSize > Size+F && 2*QuorumSize <= Size+F+2) to be:%v, got: %v, for size: %v", true, false, valSet.Size())
//		}
//	}
//}
