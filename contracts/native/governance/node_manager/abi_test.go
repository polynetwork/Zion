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

package node_manager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/stretchr/testify/assert"
)

func TestABIShowJonString(t *testing.T) {
	t.Log(INodeManagerABI)
	for name, v := range ABI.Methods {
		t.Logf("method %s, id %s", name, hexutil.Encode(v.ID))
	}
}

func TestABIMethodContractName(t *testing.T) {
	enc, err := utils.PackOutputs(ABI, MethodName, contractName)
	assert.NoError(t, err)
	params := new(MethodContractNameOutput)
	assert.NoError(t, utils.UnpackOutputs(ABI, MethodName, params, enc))
	assert.Equal(t, contractName, params.Name)
}

func TestABIMethodProposeInput(t *testing.T) {
	expectEpochID := uint64(0)
	expectPeers := generateTestPeers(10)
	expect := &MethodProposeInput{StartHeight: expectEpochID, Peers: expectPeers}

	enc, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodProposeInput)
	assert.NoError(t, got.Decode(enc))

	assert.Equal(t, expect, got)
}

func TestABIMethodProposeOutput(t *testing.T) {
	var cases = []struct {
		Result bool
	}{
		{
			Result: true,
		},
		{
			Result: false,
		},
	}

	for _, testdata := range cases {
		expect := &MethodProposeOutput{Success: testdata.Result}
		enc, err := expect.Encode()
		assert.NoError(t, err)

		got := new(MethodProposeOutput)
		assert.NoError(t, got.Decode(enc))

		assert.Equal(t, expect, got)
	}
}

func TestABIMethodVoteInput(t *testing.T) {
	var cases = []struct {
		EpochID   uint64
		EpochHash common.Hash
	}{
		//{
		//	EpochID:   0,
		//	EpochHash: generateTestHash(1),
		//},
		{
			EpochID:   1,
			EpochHash: generateTestHash(11),
		},
	}

	for _, testdata := range cases {
		expect := &MethodVoteInput{EpochID: testdata.EpochID, EpochHash: testdata.EpochHash}
		enc, err := expect.Encode()
		assert.NoError(t, err)

		got := new(MethodVoteInput)
		assert.NoError(t, got.Decode(enc))

		assert.Equal(t, expect, got)
	}
}

func TestABIMethodVoteOutput(t *testing.T) {
	var cases = []struct {
		Result bool
	}{
		{
			Result: true,
		},
		{
			Result: false,
		},
	}

	for _, testdata := range cases {
		expect := &MethodVoteOutput{Success: testdata.Result}
		enc, err := expect.Encode()
		assert.NoError(t, err)

		got := new(MethodVoteOutput)
		assert.NoError(t, got.Decode(enc))

		assert.Equal(t, expect, got)
	}
}

func TestABIMethodEpochOutput(t *testing.T) {
	expect := new(MethodEpochOutput)
	expect.Epoch = generateTestEpochInfo(1, 12, 15)
	enc, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodEpochOutput)
	assert.NoError(t, got.Decode(enc))

	assert.Equal(t, expect, got)
}

func TestABIMethodGetEpochByID(t *testing.T) {
	expect := new(MethodGetEpochByIDInput)
	expect.EpochID = uint64(56)
	enc, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodGetEpochByIDInput)
	assert.NoError(t, got.Decode(enc))

	assert.Equal(t, expect, got)
}

func TestABIMethodProofInput(t *testing.T) {
	expect := new(MethodProofInput)
	expect.EpochID = uint64(9932)
	enc, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodProofInput)
	assert.NoError(t, got.Decode(enc))

	assert.Equal(t, expect, got)
}

func TestABIMethodProofOutput(t *testing.T) {
	proof := generateTestHash(138345729384)
	expect := &MethodProofOutput{Hash: proof}

	enc, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodProofOutput)
	assert.NoError(t, got.Decode(enc))

	assert.Equal(t, expect, got)
}

func TestABIMethodGetEpochListJsonInput(t *testing.T) {
	expect := new(MethodGetEpochListJsonInput)
	expect.EpochID = uint64(9932)
	enc, err := expect.Encode()
	assert.NoError(t, err)

	got := new(MethodGetEpochListJsonInput)
	assert.NoError(t, got.Decode(enc))

	assert.Equal(t, expect, got)
}

func TestABIMethodGetJsonOutput(t *testing.T) {
	expect := new(MethodGetJsonOutput)
	epochID := uint64(2)
	peers := generateTestPeers(7)
	epoch := &EpochInfo{ID: epochID, Proposer: peers.List[0].Address, Peers: peers, StartHeight: 60}
	expect.Result = epoch.Json()
	enc, err := expect.Encode(MethodGetEpochListJson)
	assert.NoError(t, err)

	got := new(MethodGetJsonOutput)
	assert.NoError(t, got.Decode(enc, MethodGetEpochListJson))

	assert.Equal(t, expect, got)
}
