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

package zion

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/assert"
)

var testGenesisHeaderJsonStr = `{"parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","miner":"0x0000000000000000000000000000000000000000","stateRoot":"0xae517c2a24eccd0a44c939fb83a5d4123cfd945b8135cc79b4ab50d787c4bdda","transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","difficulty":"0x1","number":"0x0","gasLimit":"0xffffffff","gasUsed":"0x0","timestamp":"0x0","extraData":"0x0000000000000000000000000000000000000000000000000000000000000000f89bf85494258af48e28e4a6846e931ddff8e1cdf8579821e5946a708455c8777630aac9d1e7702d13f7a865b27c948c09d936a1b408d6e0afaa537ba4e06c4504a0ae94ad3bf5ed640cc72f37bd21d64a65c3c756e9c88cb8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c080","mixHash":"0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365","nonce":"0x4510809143055965","hash":"0x7d15dbe08441fc2b891af7a168e0df94d4363994652afa5528650bc8376844d9"}`

func TestStoreGenesis(t *testing.T) {
	resetTestContext()
	s := testEmptyCtx

	chainID := params.DevnetMainChainID
	expectHeader := testGenesisHeader(t)

	assert.False(t, isGenesisStored(s, chainID))
	assert.NoError(t, storeGenesis(s, chainID, expectHeader))
	assert.True(t, isGenesisStored(s, chainID))

	header, err := getGenesisHeader(s, chainID)
	assert.NoError(t, err)
	assert.Equal(t, expectHeader, header)
}

func TestStoreEpoch(t *testing.T) {
	var testCases = []struct {
		ChainID    uint64
		Height     uint64
		Validators []common.Address
		Err        error
	}{
		{
			ChainID: 212,
			Height:  199,
			Validators: []common.Address{
				utils.Big2Address(big.NewInt(12)),
				utils.Big2Address(big.NewInt(13)),
				utils.Big2Address(big.NewInt(14)),
			},
			Err: nil,
		},
		{
			ChainID: 0, // allow chainID to be zero
			Height:  199,
			Validators: []common.Address{
				utils.Big2Address(big.NewInt(12)),
				utils.Big2Address(big.NewInt(13)),
				utils.Big2Address(big.NewInt(14)),
			},
			Err: nil,
		},
		{
			ChainID: 12,
			Height:  0, // allow height to be nil
			Validators: []common.Address{
				utils.Big2Address(big.NewInt(12)),
				utils.Big2Address(big.NewInt(13)),
				utils.Big2Address(big.NewInt(14)),
			},
			Err: nil,
		},
		{
			ChainID:    12,
			Height:     245,
			Validators: []common.Address{}, // allow validators to be empty
			Err:        nil,
		},
	}

	for _, tc := range testCases {
		resetTestContext()
		s := testEmptyCtx

		assert.NoError(t, storeEpoch(s, tc.ChainID, tc.Height, tc.Validators))

		gotHeight, gotValidators, err := getEpoch(s, tc.ChainID)
		assert.NoError(t, err)

		assert.Equal(t, tc.Height, gotHeight)
		assert.Equal(t, tc.Validators, gotValidators)
	}
}

func testGenesisHeader(t *testing.T) *types.Header {
	header := new(types.Header)
	if err := header.UnmarshalJSON([]byte(testGenesisHeaderJsonStr)); err != nil {
		t.Fatal(err)
	}
	return header
}
