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

package validator

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	testAddress  = "70524d664ffe731100208a0154e556f9bb679ae6"
	testAddress2 = "b37866a925bccd69cfa98d43b510f1d23d78a851"
)

func TestValidatorSet(t *testing.T) {
	testNewValidatorSet(t)
	testNormalValSet(t)
	testCalcProposerByIndex(t)
	testEmptyValSet(t)
	testStickyProposer(t)
	testAddAndRemoveValidator(t)
}

func testNewValidatorSet(t *testing.T) {
	var validators []hotstuff.Validator
	const ValCnt = 100

	// Create 100 validators with random addresses
	b := []byte{}
	for i := 0; i < ValCnt; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		val := New(addr)
		validators = append(validators, val)
		b = append(b, val.Address().Bytes()...)
	}

	// Create ValidatorSet
	valSet := NewSet(ExtractValidators(b), hotstuff.RoundRobin)
	if valSet == nil {
		t.Errorf("the validator byte array cannot be parsed")
		t.FailNow()
	}

	// Check validators sorting: should be in ascending order
	for i := 0; i < ValCnt-1; i++ {
		val := valSet.GetByIndex(uint64(i))
		nextVal := valSet.GetByIndex(uint64(i + 1))
		if strings.Compare(val.String(), nextVal.String()) >= 0 {
			t.Errorf("validator set is not sorted in ascending order")
		}
	}
}

func testCalcProposerByIndex(t *testing.T) {
	const ValCnt = 100

	var addrs []common.Address
	for i := 0; i < ValCnt; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		addrs = append(addrs, addr)
	}
	valset := newDefaultSet(addrs, hotstuff.RoundRobin)
	expected := valset.validators

	// Check validators sorting: should be in ascending order
	for i := 0; i < ValCnt-1; i++ {
		valset.CalcProposerByIndex(uint64(i))
		assert.Equal(t, expected[i], valset.GetProposer())
	}
}

func testNormalValSet(t *testing.T) {
	b1 := common.Hex2Bytes(testAddress)
	b2 := common.Hex2Bytes(testAddress2)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	val1 := New(addr1)
	val2 := New(addr2)

	valSet := newDefaultSet([]common.Address{addr1, addr2}, hotstuff.RoundRobin)
	if valSet == nil {
		t.Errorf("the format of validator set is invalid")
		t.FailNow()
	}

	// check size
	if size := valSet.Size(); size != 2 {
		t.Errorf("the size of validator set is wrong: have %v, want 2", size)
	}
	// test get by index
	if val := valSet.GetByIndex(uint64(0)); !reflect.DeepEqual(val, val1) {
		t.Errorf("validator mismatch: have %v, want %v", val, val1)
	}
	// test get by invalid index
	if val := valSet.GetByIndex(uint64(2)); val != nil {
		t.Errorf("validator mismatch: have %v, want nil", val)
	}
	// test get by address
	if _, val := valSet.GetByAddress(addr2); !reflect.DeepEqual(val, val2) {
		t.Errorf("validator mismatch: have %v, want %v", val, val2)
	}
	// test get by invalid address
	invalidAddr := common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
	if _, val := valSet.GetByAddress(invalidAddr); val != nil {
		t.Errorf("validator mismatch: have %v, want nil", val)
	}
	// test get Proposer
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("Proposer mismatch: have %v, want %v", val, val1)
	}
	// test calculate Proposer
	lastProposer := addr1
	valSet.CalcProposer(lastProposer, uint64(0))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("Proposer mismatch: have %v, want %v", val, val2)
	}
	valSet.CalcProposer(lastProposer, uint64(3))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("Proposer mismatch: have %v, want %v", val, val1)
	}
	// test empty last Proposer
	lastProposer = common.Address{}
	valSet.CalcProposer(lastProposer, uint64(3))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("Proposer mismatch: have %v, want %v", val, val2)
	}
}

func testEmptyValSet(t *testing.T) {
	valSet := NewSet(ExtractValidators([]byte{}), hotstuff.RoundRobin)
	if valSet == nil {
		t.Errorf("validator set should not be nil")
	}
}

func testAddAndRemoveValidator(t *testing.T) {
	valSet := NewSet(ExtractValidators([]byte{}), hotstuff.RoundRobin)
	if !valSet.AddValidator(common.HexToAddress("0x2")) {
		t.Error("the validator should be added")
	}
	if valSet.AddValidator(common.HexToAddress("0x2")) {
		t.Error("the existing validator should not be added")
	}
	valSet.AddValidator(common.HexToAddress("0x1"))
	valSet.AddValidator(common.HexToAddress("0x0"))
	if len(valSet.List()) != 3 {
		t.Error("the size of validator set should be 3")
	}

	for i, v := range valSet.List() {
		expected := common.HexToAddress(fmt.Sprintf("0x%d", i))
		if v.Address() != expected {
			t.Errorf("the order of validators is wrong: have %v, want %v", v.Address().Hex(), expected.Hex())
		}
	}

	if !valSet.RemoveValidator(common.HexToAddress("0x2")) {
		t.Error("the validator should be removed")
	}
	if valSet.RemoveValidator(common.HexToAddress("0x2")) {
		t.Error("the non-existing validator should not be removed")
	}
	if len(valSet.List()) != 2 {
		t.Error("the size of validator set should be 2")
	}
	valSet.RemoveValidator(common.HexToAddress("0x1"))
	if len(valSet.List()) != 1 {
		t.Error("the size of validator set should be 1")
	}
	valSet.RemoveValidator(common.HexToAddress("0x0"))
	if len(valSet.List()) != 0 {
		t.Error("the size of validator set should be 0")
	}
}

func testStickyProposer(t *testing.T) {
	b1 := common.Hex2Bytes(testAddress)
	b2 := common.Hex2Bytes(testAddress2)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	val1 := New(addr1)
	val2 := New(addr2)

	valSet := newDefaultSet([]common.Address{addr1, addr2}, hotstuff.Sticky)

	// test get Proposer
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}
	// test calculate proposer
	lastproposer := addr1
	valSet.CalcProposer(lastproposer, uint64(0))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}

	valSet.CalcProposer(lastproposer, uint64(1))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}
	// test empty last proposer
	lastproposer = common.Address{}
	valSet.CalcProposer(lastproposer, uint64(3))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}
}

func TestFAndQ(t *testing.T) {
	vs := newDefaultSet([]common.Address{
		common.HexToAddress("0x1"),
		common.HexToAddress("0x2"),
		common.HexToAddress("0x3"),
		common.HexToAddress("0x4"),
		common.HexToAddress("0x5"),
		common.HexToAddress("0x6"),
		common.HexToAddress("0x7"),
	}, hotstuff.RoundRobin)

	faultySize := vs.F()
	quorumSize := vs.Q()
	assert.Equal(t, 2, faultySize)
	assert.Equal(t, 5, quorumSize)
	t.Logf("faulty size %d, quorum size %d", faultySize, quorumSize)
}
