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

package signature_manager

import (
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/signature_manager_abi"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	GetABI()
	nm.InitNodeManager()
	InitSignatureManager()
	os.Exit(m.Run())
}

func TestAddSignature(t *testing.T) {

	name := MethodAddSignature

	// 1.unpack error, subject should be bytes, use Uint256 instead
	{
		var (
			Uint256, _ = abi.NewType("uint256", "", nil)
			Bytes, _   = abi.NewType("bytes", "", nil)
			Address, _ = abi.NewType("address", "", nil)
		)

		input := abi.Arguments{
			{Name: "addr", Type: Address, Indexed: false},
			{Name: "sideChainID", Type: Uint256, Indexed: false},
			{Name: "subject", Type: Uint256, Indexed: false},
			{Name: "sig", Type: Bytes, Indexed: false},
		}

		var (
			addr         = common.HexToAddress("0x123")
			sideChainID  = big.NewInt(2)
			errorSubject = big.NewInt(3)
			sig          = []byte{'1'}
		)

		args, err := input.Pack(addr, sideChainID, errorSubject, sig)
		assert.NoError(t, err)
		methodID := ABI.Methods[name].ID
		payload := append(methodID, args...)
		supplyGas := gasTable[name]

		_, err = native.TestNativeCall(t, this, name, payload, supplyGas)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "abi", "marshal")
	}

	// 2.supply gas not enough,
	{
		payload, err := utils.PackMethod(ABI, name, common.HexToAddress("0x12"), big.NewInt(2), []byte{'a'}, []byte{'1'})
		assert.NoError(t, err)
		supplyGas := gasTable[name] - 1

		_, err = native.TestNativeCall(t, this, name, payload, supplyGas)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "gas", "supply", "left")
	}

	// 3.witness checking failed
	{
		var (
			sender    = common.HexToAddress("0x123")
			addr      = common.HexToAddress("0x12")
			supplyGas = gasTable[name]
		)
		payload, err := utils.PackMethod(ABI, name, addr, big.NewInt(2), []byte{'a'}, []byte{'1'})
		assert.NoError(t, err)

		_, err = native.TestNativeCall(t, this, name, payload, sender, supplyGas)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "checkWitness", "authentication")
	}

	// 4.signature invalid???
	{
		peers, _ := native.GenerateTestPeers(4)

		var (
			sender    = peers[0]
			chainID   = big.NewInt(2)
			supplyGas = gasTable[name]
			subject   = []byte{'a'}
			errSig    = []byte{'b'}
		)

		payload, err := utils.PackMethod(ABI, name, sender, chainID, subject, errSig)
		assert.NoError(t, err)

		_, err = native.TestNativeCall(t, this, name, payload, sender, supplyGas, func(state *state.StateDB) {
			nm.StoreGenesisEpoch(state, peers, peers)
		})
		t.Error(err)
	}
}
