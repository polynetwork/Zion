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
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/stretchr/testify/assert"
)

func TestABIShowJonString(t *testing.T) {
	t.Log(abijson)
}

func TestABIMethodContractName(t *testing.T) {
	enc, err := utils.PackOutputs(ABI, MethodContractName, contractName)
	assert.NoError(t, err)
	params := new(MethodContractNameOutput)
	assert.NoError(t, utils.UnpackOutputs(ABI, MethodContractName, params, enc))
	assert.Equal(t, contractName, params.Name)
}

func TestABIMethodProposeInput(t *testing.T) {
	expectEpochID := uint64(0)
	expectPeers := generateTestPeers(10)
	expect := new(MethodProposeInput)

	enc, err := expect.Encode(expectEpochID, expectPeers)
	assert.NoError(t, err)

	got := new(MethodProposeInput)
	gotEpochID, gotPeers, err := got.Decode(enc)
	assert.NoError(t, err)

	assert.Equal(t, expectEpochID, gotEpochID)
	assert.Equal(t, expectPeers, gotPeers)
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
		expect := new(MethodProposeOutput)
		enc, err := expect.Encode(testdata.Result)
		assert.NoError(t, err)

		got := new(MethodProposeOutput)
		assert.NoError(t, got.Decode(enc))

		assert.Equal(t, expect, got)
	}
}

func TestABIMethodVoteInput(t *testing.T) {
	var cases = []struct {
		EpochID uint64
		Hash    common.Hash
	}{
		{
			EpochID: 0,
			Hash:    generateTestHash(1),
		},
		{
			EpochID: 1,
			Hash:    generateTestHash(11),
		},
	}

	for _, testdata := range cases {
		expect := new(MethodVoteInput)
		enc, err := expect.Encode(testdata.EpochID, testdata.Hash)
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
		expect := new(MethodVoteOutput)
		enc, err := expect.Encode(testdata.Result)
		assert.NoError(t, err)

		got := new(MethodVoteOutput)
		assert.NoError(t, got.Decode(enc))

		assert.Equal(t, expect, got)
	}
}

func TestABIMethodEpochOutput(t *testing.T) {
	origin := new(MethodEpochOutput)
	expect := generateTestEpochInfo(1, 12, 15)
	enc, err := origin.Encode(expect)
	assert.NoError(t, err)

	dst := new(MethodEpochOutput)
	got, err := dst.Decode(enc)
	assert.NoError(t, err)

	assert.Equal(t, expect, got)
}

func TestABIMethodNextEpochOutput(t *testing.T) {
	origin := new(MethodNextEpochOutput)
	expect := generateTestEpochInfo(1, 12, 15)
	enc, err := origin.Encode(expect)
	assert.NoError(t, err)

	dst := new(MethodNextEpochOutput)
	got, err := dst.Decode(enc)
	assert.NoError(t, err)

	assert.Equal(t, expect, got)
}
