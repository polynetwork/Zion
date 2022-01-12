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
package relayer_manager

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func init() {
	InitRelayerManager()
	node_manager.InitNodeManager()
	db := rawdb.NewMemoryDatabase()
	sdb, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)

	putPeerMapPoolAndView(sdb)
}

func putPeerMapPoolAndView(db *state.StateDB) {
	height := uint64(120)
	epoch := node_manager.GenerateTestEpochInfo(1, height, 4)
	peer := epoch.Peers.List[0]
	rawPubKey, _ := hexutil.Decode(peer.PubKey)
	pubkey, _ := crypto.DecompressPubkey(rawPubKey)
	acct = pubkey
	caller := peer.Address

	txhash := common.HexToHash("0x123")
	ref := native.NewContractRef(db, caller, caller, new(big.Int).SetUint64(height), txhash, 0, nil)
	s := native.NewNativeContract(db, ref)
	node_manager.StoreTestEpoch(s, epoch)
}

var (
	sdb  *state.StateDB
	acct *ecdsa.PublicKey
)

func TestRegisterRelayer(t *testing.T) {
	{
		params := new(RelayerListParam)
		params.AddressList = []common.Address{{1, 2, 4, 6}, {1, 4, 5, 7}, {1, 3, 5, 7, 9}}

		input, err := utils.PackMethodWithStruct(ABI, MethodRegisterRelayer, params)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		contractRef := native.NewContractRef(sdb, common.Address{}, common.Address{}, blockNumber, common.Hash{}, gasTable[MethodRegisterRelayer]+extra, nil)
		ret, leftOverGas, err := contractRef.NativeCall(common.Address{}, utils.RelayerManagerContractAddress, input)

		assert.Nil(t, err)

		result, err := utils.PackOutputs(ABI, MethodRegisterRelayer, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, extra)

		contract := native.NewNativeContract(sdb, contractRef)
		relayerListParam, err := getRelayerApply(contract, 0)
		assert.Nil(t, err)
		assert.Equal(t, params, relayerListParam)
	}

	// none consensus acct should not be able to approve register relayer
	{
		caller := crypto.PubkeyToAddress(*acct)
		param := new(ApproveRelayerParam)
		param.ID = 0
		param.Address = caller

		input, err := utils.PackMethodWithStruct(ABI, MethodApproveRegisterRelayer, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, gasTable[MethodApproveRegisterRelayer]+extra, nil)
		ret, leftOverGas, err := contractRef.NativeCall(caller, utils.RelayerManagerContractAddress, input)

		assert.Nil(t, err)

		result, err := utils.PackOutputs(ABI, MethodApproveRegisterRelayer, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, extra)

		contract := native.NewNativeContract(sdb, contractRef)
		ok, err := node_manager.CheckConsensusSigns(contract, MethodApproveRegisterRelayer, utils.GetUint64Bytes(0), caller)
		assert.Nil(t, err)
		assert.Equal(t, true, ok)
	}

}

func TestRemoveRelayer(t *testing.T) {
	{
		params := new(RelayerListParam)
		params.AddressList = []common.Address{{1, 2, 4, 6}, {1, 4, 5, 7}}

		input, err := utils.PackMethodWithStruct(ABI, MethodRemoveRelayer, params)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		contractRef := native.NewContractRef(sdb, common.Address{}, common.Address{}, blockNumber, common.Hash{}, gasTable[MethodRemoveRelayer]+extra, nil)
		ret, leftOverGas, err := contractRef.NativeCall(common.Address{}, utils.RelayerManagerContractAddress, input)

		assert.Nil(t, err)

		result, err := utils.PackOutputs(ABI, MethodRemoveRelayer, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, extra)

		contract := native.NewNativeContract(sdb, contractRef)
		relayerListParam, err := getRelayerRemove(contract, 0)
		assert.Nil(t, err)
		assert.Equal(t, params, relayerListParam)
	}

	{
		caller := crypto.PubkeyToAddress(*acct)
		param := new(ApproveRelayerParam)
		param.ID = 0
		param.Address = caller

		input, err := utils.PackMethodWithStruct(ABI, MethodApproveRemoveRelayer, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, gasTable[MethodApproveRemoveRelayer]+extra, nil)
		ret, leftOverGas, err := contractRef.NativeCall(caller, utils.RelayerManagerContractAddress, input)

		assert.Nil(t, err)

		result, err := utils.PackOutputs(ABI, MethodApproveRemoveRelayer, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, extra)

		contract := native.NewNativeContract(sdb, contractRef)
		ok, err := node_manager.CheckConsensusSigns(contract, MethodApproveRemoveRelayer, utils.GetUint64Bytes(0), caller)
		assert.Nil(t, err)
		assert.Equal(t, true, ok)
	}
}
