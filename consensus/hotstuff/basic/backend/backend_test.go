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

package backend
//
//import (
//	"bytes"
//	"math/big"
//	"testing"
//	"time"
//
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/consensus/hotstuff"
//	"github.com/ethereum/go-ethereum/core/types"
//	"github.com/ethereum/go-ethereum/crypto"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestCommit(t *testing.T) {
//	backend := newBackend()
//
//	commitCh := make(chan *types.Block)
//	// Case: it's a proposer, so the backend.commit will receive channel result from backend.Commit function
//	testCases := []struct {
//		expectedErr       error
//		expectedSignature [][]byte
//		expectedBlock     func() *types.Block
//	}{
//		{
//			// normal case
//			nil,
//			[][]byte{append([]byte{1}, bytes.Repeat([]byte{0x00}, types.HotstuffExtraSeal-1)...)},
//			func() *types.Block {
//				chain, engine := newBlockChain(1)
//				block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
//				expectedBlock, _ := engine.UpdateBlock(block)
//				return expectedBlock
//			},
//		},
//		{
//			// invalid signature
//			errInvalidCommittedSeals,
//			nil,
//			func() *types.Block {
//				chain, engine := newBlockChain(1)
//				block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
//				expectedBlock, _ := engine.UpdateBlock(block)
//				return expectedBlock
//			},
//		},
//	}
//
//	for _, test := range testCases {
//		expBlock := test.expectedBlock()
//		go func() {
//			result := <-backend.commitCh
//			commitCh <- result
//		}()
//
//		backend.proposedBlockHash = expBlock.Hash()
//		// todo: modify this function to be TestPreCommit
//		if _, err := backend.PreCommit(expBlock, test.expectedSignature); err != nil {
//			if err != test.expectedErr {
//				t.Errorf("error mismatch: have %v, want %v", err, test.expectedErr)
//			}
//		}
//
//		if test.expectedErr == nil {
//			// to avoid race condition is occurred by goroutine
//			select {
//			case result := <-commitCh:
//				if result.Hash() != expBlock.Hash() {
//					t.Errorf("hash mismatch: have %v, want %v", result.Hash(), expBlock.Hash())
//				}
//			case <-time.After(10 * time.Second):
//				t.Fatal("timeout")
//			}
//		}
//	}
//}
//
//// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/basic/backend -run TestGetProposer
//func TestGetProposer(t *testing.T) {
//	chain, engine := newBlockChain(4)
//	block := makeBlock(chain, engine, chain.Genesis())
//	chain.InsertChain(types.Blocks{block})
//	expected := engine.GetProposer(1)
//	actual := engine.Address()
//	assert.Equal(t, expected, actual, "proposer mismatch: have %v, want %v", actual.Hex(), expected.Hex())
//}
