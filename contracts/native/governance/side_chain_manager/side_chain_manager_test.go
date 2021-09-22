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
package side_chain_manager

import (
	"encoding/hex"
	"math/big"
	"testing"

	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	polycomm "github.com/polynetwork/poly/common"
	cstates "github.com/polynetwork/poly/core/states"
	"github.com/stretchr/testify/assert"
)

func init() {
	InitSideChainManager()
	node_manager.InitNodeManager()
	db := rawdb.NewMemoryDatabase()
	sdb, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)

	cacheDB := (*state.CacheDB)(sdb)
	putPeerMapPoolAndView(cacheDB)
}

func putPeerMapPoolAndView(db *state.CacheDB) {
	key, _ := crypto.GenerateKey()
	acct = &key.PublicKey

	peerPoolMap := new(node_manager.PeerPoolMap)
	peerPoolMap.PeerPoolMap = make(map[string]*node_manager.PeerPoolItem)
	pkStr := hex.EncodeToString(crypto.FromECDSAPub(acct))
	peerPoolMap.PeerPoolMap[pkStr] = &node_manager.PeerPoolItem{
		Index:      uint32(0),
		PeerPubkey: pkStr,
		Address:    crypto.PubkeyToAddress(*acct),
		Status:     node_manager.ConsensusStatus,
	}
	view := uint32(0)
	viewBytes := utils.GetUint32Bytes(view)
	sink := polycomm.NewZeroCopySink(nil)
	peerPoolMap.Serialization(sink)
	db.Put(utils.ConcatKey(utils.NodeManagerContractAddress, []byte(node_manager.PEER_POOL), viewBytes), cstates.GenRawStorageItem(sink.Bytes()))

	sink.Reset()

	govView := node_manager.GovernanceView{
		View: view,
	}
	govView.Serialization(sink)
	db.Put(utils.ConcatKey(utils.NodeManagerContractAddress, []byte(node_manager.GOVERNANCE_VIEW)), cstates.GenRawStorageItem(sink.Bytes()))
}

var (
	sdb  *state.StateDB
	acct *ecdsa.PublicKey
)

func TestRegisterSideChainManager(t *testing.T) {
	param := new(RegisterSideChainParam)
	param.BlocksToWait = 4
	param.ChainId = 8
	param.Name = "mychain"
	param.Router = 3

	input, err := utils.PackMethodWithStruct(ABI, MethodRegisterSideChain, param)
	assert.Nil(t, err)

	blockNumber := big.NewInt(1)
	extra := uint64(10)
	contractRef := native.NewContractRef(sdb, common.Address{}, common.Address{}, blockNumber, common.Hash{}, gasTable[MethodRegisterSideChain]+extra, nil)
	ret, leftOverGas, err := contractRef.NativeCall(common.Address{}, utils.SideChainManagerContractAddress, input)

	assert.Nil(t, err)

	result, err := utils.PackOutputs(ABI, MethodRegisterSideChain, true)
	assert.Nil(t, err)
	assert.Equal(t, ret, result)
	assert.Equal(t, leftOverGas, extra)

	contract := native.NewNativeContract(sdb, contractRef)
	sideChain, err := GetSideChainApply(contract, 8)
	assert.Equal(t, sideChain.Name, "mychain")
	assert.Nil(t, err)

	_, _, err = contractRef.NativeCall(common.Address{}, utils.SideChainManagerContractAddress, input)
	assert.NotNil(t, err)
}

func TestApproveRegisterSideChain(t *testing.T) {

	TestRegisterSideChainManager(t)

	caller := crypto.PubkeyToAddress(*acct)
	param := new(ChainidParam)
	param.Chainid = 8
	param.Address = caller

	input, err := utils.PackMethodWithStruct(ABI, MethodApproveRegisterSideChain, param)
	assert.Nil(t, err)

	blockNumber := big.NewInt(1)
	extra := uint64(10)
	contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, gasTable[MethodApproveRegisterSideChain]+extra, nil)
	ret, leftOverGas, err := contractRef.NativeCall(caller, utils.SideChainManagerContractAddress, input)

	assert.Nil(t, err)

	result, err := utils.PackOutputs(ABI, MethodApproveRegisterSideChain, true)
	assert.Nil(t, err)
	assert.Equal(t, ret, result)
	assert.Equal(t, leftOverGas, extra)
}

func TestUpdateSideChain(t *testing.T) {
	TestApproveRegisterSideChain(t)

	param := new(RegisterSideChainParam)
	param.Address = common.Address{}
	param.BlocksToWait = 10
	param.ChainId = 8
	param.Name = "own"
	param.Router = 3

	input, err := utils.PackMethodWithStruct(ABI, MethodUpdateSideChain, param)
	assert.Nil(t, err)

	blockNumber := big.NewInt(1)
	extra := uint64(10)
	contractRef := native.NewContractRef(sdb, common.Address{}, common.Address{}, blockNumber, common.Hash{}, gasTable[MethodUpdateSideChain]+extra, nil)
	ret, leftOverGas, err := contractRef.NativeCall(common.Address{}, utils.SideChainManagerContractAddress, input)

	assert.Nil(t, err)

	result, err := utils.PackOutputs(ABI, MethodUpdateSideChain, true)
	assert.Nil(t, err)
	assert.Equal(t, ret, result)
	assert.Equal(t, leftOverGas, extra)
}

func TestApproveUpdateSideChain(t *testing.T) {
	TestUpdateSideChain(t)

	caller := crypto.PubkeyToAddress(*acct)

	param := new(ChainidParam)
	param.Chainid = 8
	param.Address = caller

	input, err := utils.PackMethodWithStruct(ABI, MethodApproveUpdateSideChain, param)
	assert.Nil(t, err)

	blockNumber := big.NewInt(1)
	extra := uint64(10)
	contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, gasTable[MethodApproveUpdateSideChain]+extra, nil)
	ret, leftOverGas, err := contractRef.NativeCall(caller, utils.SideChainManagerContractAddress, input)

	assert.Nil(t, err)

	result, err := utils.PackOutputs(ABI, MethodApproveUpdateSideChain, true)
	assert.Nil(t, err)
	assert.Equal(t, ret, result)
	assert.Equal(t, leftOverGas, extra)

	contract := native.NewNativeContract(sdb, contractRef)
	sideChain, err := GetSideChain(contract, 8)
	assert.Equal(t, sideChain.Name, "own")
	assert.Nil(t, err)
}
