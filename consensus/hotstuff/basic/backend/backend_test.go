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

import (
	"bytes"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	b := newBackend()
	data := []byte("Here is a string....")
	sig, err := b.Sign(data)
	assert.NoError(t, err, "error mismatch: have %v, want nil", err)

	//Check signature recover
	hashData := crypto.Keccak256(data)
	pubkey, _ := crypto.Ecrecover(hashData, sig)
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])
	assert.Equal(t, signer, getAddress(), "address mismatch: have %v, want %s", signer.Hex(), getAddress().Hex())
}

func TestCheckSignature(t *testing.T) {
	key, _ := generatePrivateKey()
	data := []byte("Here is a string....")
	hashData := crypto.Keccak256([]byte(data))
	sig, _ := crypto.Sign(hashData, key)
	b := newBackend()
	a := getAddress()
	assert.NoError(t, b.CheckSignature(data, a, sig), "error mismatch: want nil")

	a = getInvalidAddress()
	assert.Error(t, b.CheckSignature(data, a, sig), "error mismatch: want%v", errInvalidSignature)
}

func TestCheckValidatorSignature(t *testing.T) {
	vset, keys := newTestValidatorSet(5)

	// 1. Positive test: sign with validator's key should succeed
	data := []byte("dummy data")
	hashData := crypto.Keccak256([]byte(data))
	for i, k := range keys {
		// Sign
		sig, err := crypto.Sign(hashData, k)
		assert.NoError(t, err, "error mismatch: have %v, want nil", err)

		// CheckValidatorSignature should succeed
		addr, err := hotstuff.CheckValidatorSignature(vset, data, sig)
		assert.NoError(t, err, "error mismatch: have %v, want nil", err)

		val := vset.GetByIndex(uint64(i))
		assert.Equal(t, addr, val.Address(), "validator address mismatch: have %v, want %v", addr, val.Address())
	}

	// 2. Negative test: sign with any key other than validator's key should return error
	key, err := crypto.GenerateKey()
	assert.NoError(t, err, "error mismatch: have %v, want nil", err)

	// Sign
	sig, err := crypto.Sign(hashData, key)
	assert.NoError(t, err, "error mismatch: have %v, want nil", err)

	// CheckValidatorSignature should return ErrUnauthorizedAddress
	addr, err := hotstuff.CheckValidatorSignature(vset, data, sig)
	assert.Equal(t, err, hotstuff.ErrUnauthorizedAddress, "error mismatch: have %v, want %v", err, hotstuff.ErrUnauthorizedAddress)

	emptyAddr := common.Address{}
	assert.Equal(t, emptyAddr, common.Address{}, "address mismatch: have %v, want %v", addr, emptyAddr)
}

func TestCommit(t *testing.T) {
	backend := newBackend()

	commitCh := make(chan *types.Block)
	// Case: it's a proposer, so the backend.commit will receive channel result from backend.Commit function
	testCases := []struct {
		expectedErr       error
		expectedSignature [][]byte
		expectedBlock     func() *types.Block
	}{
		{
			// normal case
			nil,
			[][]byte{append([]byte{1}, bytes.Repeat([]byte{0x00}, types.HotstuffExtraSeal-1)...)},
			func() *types.Block {
				chain, engine := newBlockChain(1)
				block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
				expectedBlock, _ := engine.UpdateBlock(block)
				return expectedBlock
			},
		},
		{
			// invalid signature
			errInvalidCommittedSeals,
			nil,
			func() *types.Block {
				chain, engine := newBlockChain(1)
				block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
				expectedBlock, _ := engine.UpdateBlock(block)
				return expectedBlock
			},
		},
	}

	for _, test := range testCases {
		expBlock := test.expectedBlock()
		go func() {
			result := <-backend.commitCh
			commitCh <- result
		}()

		backend.proposedBlockHash = expBlock.Hash()
		view := &hotstuff.View{
			Round:  new(big.Int).SetUint64(0),
			Height: expBlock.Number(),
		}
		// todo: modify this function to be TestPreCommit
		if _, _, err := backend.PreCommit(view, expBlock, test.expectedSignature); err != nil {
			if err != test.expectedErr {
				t.Errorf("error mismatch: have %v, want %v", err, test.expectedErr)
			}
		}

		if test.expectedErr == nil {
			// to avoid race condition is occurred by goroutine
			select {
			case result := <-commitCh:
				if result.Hash() != expBlock.Hash() {
					t.Errorf("hash mismatch: have %v, want %v", result.Hash(), expBlock.Hash())
				}
			case <-time.After(10 * time.Second):
				t.Fatal("timeout")
			}
		}
	}
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/basic/backend -run TestGetProposer
func TestGetProposer(t *testing.T) {
	chain, engine := newBlockChain(4)
	block := makeBlock(chain, engine, chain.Genesis())
	chain.InsertChain(types.Blocks{block})
	expected := engine.GetProposer(1)
	actual := engine.Address()
	assert.Equal(t, expected, actual, "proposer mismatch: have %v, want %v", actual.Hex(), expected.Hex())
}
