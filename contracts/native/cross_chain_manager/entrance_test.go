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

package cross_chain_manager

import (
	"crypto/ecdsa"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/cross_chain_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/side_chain_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)


func init() {
	side_chain_manager.InitSideChainManager()
	node_manager.InitNodeManager()
	db := rawdb.NewMemoryDatabase()
	sdb, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	InitCrossChainManager()
}

var (
	sdb  *state.StateDB
	signers []common.Address
	keys []*ecdsa.PrivateKey
)

func init() {
	node_manager.InitNodeManager()
	db := rawdb.NewMemoryDatabase()
	sdb, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	signers, keys = native.GenerateTestPeers(2)
	node_manager.StoreCommunityInfo(sdb, big.NewInt(2000), common.EmptyAddress)
	node_manager.StoreGenesisEpoch(sdb, signers, signers)
	node_manager.StoreGenesisGlobalConfig(sdb)

	param := new(side_chain_manager.RegisterSideChainParam)
	param.BlocksToWait = 4
	param.ChainID = 8
	param.Name = "mychain"

	param1 := new(side_chain_manager.RegisterSideChainParam)
	param1.ChainID = 9
	param1.Name = strings.Repeat("1", 100)
	param1.ExtraInfo = make([]byte, 1000000)
	param1.CCMCAddress = make([]byte, 1000)


	param2 := new(side_chain_manager.RegisterSideChainParam)
	param2.ChainID = 10
	param2.Name = strings.Repeat("1", 100)
	param2.ExtraInfo = make([]byte, 1000000)
	param2.CCMCAddress = make([]byte, 1000)

	param3 := new(side_chain_manager.RegisterSideChainParam)
	param3.ChainID = 11
	param3.Name = strings.Repeat("1", 100)
	param3.ExtraInfo = make([]byte, 1000000)
	param3.CCMCAddress = make([]byte, 1000)


	for _, param := range []*side_chain_manager.RegisterSideChainParam{param, param1, param2, param3} {
		input, err := utils.PackMethodWithStruct(side_chain_manager.ABI, side_chain_manager_abi.MethodRegisterSideChain, param)
		if err != nil { panic(err) }
		caller := signers[0]
		contractRef := native.NewContractRef(sdb, caller, caller, big.NewInt(1), common.Hash{}, 10000000, nil)
		_, _, err = contractRef.NativeCall(caller, utils.SideChainManagerContractAddress, input)
		if err != nil { panic(err) }
		p := new(side_chain_manager.ChainIDParam)
		p.ChainID = param.ChainID

		input, err = utils.PackMethodWithStruct(side_chain_manager.ABI, side_chain_manager_abi.MethodApproveRegisterSideChain, p)
		if err != nil { panic(err) }
		contractRef = native.NewContractRef(sdb, caller, caller, big.NewInt(1), common.Hash{}, 10000000, nil)
		_, _, err = contractRef.NativeCall(caller, utils.SideChainManagerContractAddress, input)
		if err != nil { panic(err) }
		caller = signers[1]
		contractRef = native.NewContractRef(sdb, caller, caller, big.NewInt(1), common.Hash{}, 10000000, nil)
		_, _, err = contractRef.NativeCall(caller, utils.SideChainManagerContractAddress, input)
		if err != nil { panic(err) }

		contract := native.NewNativeContract(sdb, contractRef)
		sideChain, err := side_chain_manager.GetSideChainObject(contract, param.ChainID)
		if err != nil { panic(err) }
		if sideChain == nil {
			panic("side chain not ready yet")
		}
	}
}

func TestImportOuterTransfer(t *testing.T) {
	param := new(scom.EntranceParam)
	param.SourceChainID = 8

	param1 := new(scom.EntranceParam)
	param1.SourceChainID = 9

	param2 := new(scom.EntranceParam)
	param2.SourceChainID = 10

	param3 := new(scom.EntranceParam)
	param3.SourceChainID = 11


	tr := native.NewTimer(scom.MethodImportOuterTransfer)
	for _, param := range []*scom.EntranceParam{param, param1, param2, param3} {
		digest, err := param.Digest()
		assert.Nil(t, err)
		param.Signature, err = crypto.Sign(digest, keys[0])
		assert.Nil(t, err)
		
		input, err := utils.PackMethodWithStruct(scom.ABI, cross_chain_manager_abi.MethodImportOuterTransfer, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		caller := common.Address{}
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, gasTable[cross_chain_manager_abi.MethodImportOuterTransfer]+extra, nil)
		tr.Start()
		ret, leftOverGas, err := contractRef.NativeCall(caller, utils.CrossChainManagerContractAddress, input)
		tr.Stop()
		assert.Nil(t, err)
		result, err := utils.PackOutputs(scom.ABI, cross_chain_manager_abi.MethodImportOuterTransfer, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, extra)
	}
	tr.Dump()
}

func TestWhiteChain(t *testing.T) {
	TestImportOuterTransfer(t)

	param := new(scom.BlackChainParam)
	param.ChainID = 8

	param1 := new(scom.BlackChainParam)
	param1.ChainID = 9

	tr := native.NewTimer(scom.MethodBlackChain)
	for _, param := range []*scom.BlackChainParam{param, param1} {
		input, err := utils.PackMethodWithStruct(scom.ABI, cross_chain_manager_abi.MethodWhiteChain, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		caller := signers[0]
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, gasTable[cross_chain_manager_abi.MethodWhiteChain]+extra, nil)
		tr.Start()
		ret, leftOverGas, err := contractRef.NativeCall(caller, utils.CrossChainManagerContractAddress, input)
		tr.Stop()
		assert.Nil(t, err)
		result, err := utils.PackOutputs(scom.ABI, cross_chain_manager_abi.MethodWhiteChain, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, extra)
	}
	tr.Dump()
}

func TestBlackChain(t *testing.T) {
	TestImportOuterTransfer(t)

	param := new(scom.BlackChainParam)
	param.ChainID = 8

	param1 := new(scom.BlackChainParam)
	param1.ChainID = 9

	tr := native.NewTimer(scom.MethodBlackChain)
	for _, param := range []*scom.BlackChainParam{param, param1} {
		input, err := utils.PackMethodWithStruct(scom.ABI, cross_chain_manager_abi.MethodBlackChain, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		caller := signers[0]
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, gasTable[cross_chain_manager_abi.MethodBlackChain]+extra, nil)
		tr.Start()
		ret, leftOverGas, err := contractRef.NativeCall(caller, utils.CrossChainManagerContractAddress, input)
		tr.Stop()
		assert.Nil(t, err)
		result, err := utils.PackOutputs(scom.ABI, cross_chain_manager_abi.MethodBlackChain, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, extra)
	}
	tr.Dump()
}


func TestCheckDone(t *testing.T) {
	TestImportOuterTransfer(t)

	param := new(scom.ReplenishParam)
	param.ChainID = 8

	param1 := new(scom.ReplenishParam)
	param1.ChainID = 9

	tr := native.NewTimer(scom.MethodReplenish)
	for _, param := range []*scom.ReplenishParam{param, param1} {
		input, err := utils.PackMethodWithStruct(scom.ABI, cross_chain_manager_abi.MethodReplenish, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		caller := common.Address{}
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, gasTable[cross_chain_manager_abi.MethodReplenish]+extra, nil)
		tr.Start()
		ret, leftOverGas, err := contractRef.NativeCall(caller, utils.CrossChainManagerContractAddress, input)
		tr.Stop()
		assert.Nil(t, err)
		result, err := utils.PackOutputs(scom.ABI, cross_chain_manager_abi.MethodReplenish, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, extra)
	}
	tr.Dump()
}

func TestReplenish(t *testing.T) {
	TestImportOuterTransfer(t)

	param := new(scom.ReplenishParam)
	param.ChainID = 8

	param1 := new(scom.ReplenishParam)
	param1.ChainID = 9

	tr := native.NewTimer(scom.MethodReplenish)
	for _, param := range []*scom.ReplenishParam{param, param1} {
		input, err := utils.PackMethodWithStruct(scom.ABI, cross_chain_manager_abi.MethodReplenish, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		caller := common.Address{}
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, gasTable[cross_chain_manager_abi.MethodReplenish]+extra, nil)
		tr.Start()
		ret, leftOverGas, err := contractRef.NativeCall(caller, utils.CrossChainManagerContractAddress, input)
		tr.Stop()
		assert.Nil(t, err)
		result, err := utils.PackOutputs(scom.ABI, cross_chain_manager_abi.MethodReplenish, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, extra)
	}
	tr.Dump()
}