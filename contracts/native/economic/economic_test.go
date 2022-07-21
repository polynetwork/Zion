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

package economic

import (
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/economic_abi"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	InitABI()
	nm.InitNodeManager()
	InitEconomic()
	os.Exit(m.Run())
}

// TestTotalSupply use command as follow to test each cases, and the result contains coverage and output. and use
// the flag of -count=1 to avoid the affect of test cache.
// cmd:
// go test -v -count=1 -cover github.com/ethereum/go-ethereum/contracts/native/economic -run TestTotalSupply
//
func TestTotalSupply(t *testing.T) {

	// genesis supply should be 100,000,000 and total supply has no upper limit.
	testcases := []struct {
		height  int
		expect  *big.Int
		testABI bool
	}{
		{0, big.NewInt(100000000), false},
		{40, big.NewInt(100000040), false},
		{200000000, big.NewInt(300000000), true},
	}

	for _, tc := range testcases {
		var supply *big.Int

		payload, _ := new(MethodTotalSupplyInput).Encode()
		raw, err := native.TestNativeCall(t, this, payload, tc.height, common.EmptyAddress)
		assert.NoError(t, err)

		if tc.testABI {
			output, err := ABI.Unpack(MethodTotalSupply, raw)
			assert.NoError(t, err)
			supply = *abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
		} else {
			assert.NoError(t, utils.UnpackOutputs(ABI, MethodTotalSupply, &supply, raw))
		}

		got := new(big.Int).Div(supply, params.ZNT1)
		assert.Equal(t, tc.expect, got)
	}
}

//func TestReward(t *testing.T) {
//	testcases := []struct {
//		pool common.Address
//		rate int
//		expectPoolAmount *big.Int
//		expectStakeAmount *big.Int
//		errNil bool
//	}{
//		{common.EmptyAddress, 0, common.Big0, params.ZNT1, true},
//		{common.HexToAddress("0x123"), 2000, new(big.Int).SetUint64(1e18), },
//		{common.HexToAddress("0x123"), 10000},
//		{common.HexToAddress("0x123"), 10001},
//	}
//
//
//}
