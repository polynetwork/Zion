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

package info_sync

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

var (
	sdb              *state.StateDB
	testGenesisNum   = 4
	pub              *ecdsa.PublicKey
	key              *ecdsa.PrivateKey
	testGenesisPeers []common.Address
	testGenesisPri   []*ecdsa.PrivateKey
)

const (
	CHAIN_ID uint64 = 1
)

func Init() {
	key, _ = crypto.GenerateKey()
	pub = &key.PublicKey

	node_manager.InitNodeManager()
	side_chain_manager.InitSideChainManager()
	InitInfoSync()
	sdb = native.NewTestStateDB()
	testGenesisPeers, testGenesisPri = native.GenerateTestPeers(testGenesisNum)
	node_manager.StoreGenesisEpoch(sdb, testGenesisPeers, testGenesisPeers)

	putSideChain()
}

func putSideChain() {
	blockNumber := big.NewInt(1)
	extra := uint64(10)
	contractRef := native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
	contract := native.NewNativeContract(sdb, contractRef)

	err := side_chain_manager.PutSideChain(contract, &side_chain_manager.SideChain{
		Router:  utils.NO_PROOF_ROUTER,
		ChainID: CHAIN_ID,
	})
	if err != nil {
		log.Fatalf("PutSideChain fail:%v", err)
		return
	}
	sideChain, err := side_chain_manager.GetSideChainObject(contract, CHAIN_ID)
	if err != nil {
		log.Fatalf("PutSideChain fail:%v", err)
		return
	}

	if sideChain.ChainID != CHAIN_ID {
		log.Fatalf("GetSideChain mismatch")
	}
}

func TestNoAuthSyncRootInfo(t *testing.T) {
	Init()
	var err error
	param := new(SyncRootInfoParam)
	param.ChainID = CHAIN_ID
	rootInfo1 := &RootInfo{Height: 2, Info: []byte{0x01, 0x02}}
	rootInfo2 := &RootInfo{Height: 3, Info: []byte{0x02, 0x03}}
	b1, err := rlp.EncodeToBytes(rootInfo1)
	assert.Nil(t, err)
	b2, err := rlp.EncodeToBytes(rootInfo2)
	assert.Nil(t, err)
	param.RootInfos = [][]byte{b1, b2}
	digest, err := param.Digest()
	assert.Nil(t, err)
	sig, err := crypto.Sign(digest, key)
	assert.Nil(t, err)
	param.Signature = sig

	input, err := param.Encode()
	assert.Nil(t, err)

	caller := crypto.PubkeyToAddress(*pub)
	extra := uint64(21000000000000)
	_, err = native.TestNativeCall(t, utils.InfoSyncContractAddress, "SyncRootInfo", input, caller, caller, extra, sdb)
	assert.NotNil(t, err)
}

func TestNormalSyncRootInfo(t *testing.T) {
	Init()
	var err error
	param := new(SyncRootInfoParam)
	param.ChainID = CHAIN_ID
	rootInfo1 := &RootInfo{Height: 100, Info: []byte{0x01, 0x02}}
	rootInfo2 := &RootInfo{Height: 98, Info: []byte{0x02, 0x03}}
	b1, err := rlp.EncodeToBytes(rootInfo1)
	assert.Nil(t, err)
	b2, err := rlp.EncodeToBytes(rootInfo2)
	assert.Nil(t, err)
	param.RootInfos = [][]byte{b1, b2}

	for i := 0; i < testGenesisNum; i++ {
		caller := testGenesisPeers[i]
		extra := uint64(21000000000000)
		digest, err := param.Digest()
		assert.Nil(t, err)
		sig, err := crypto.Sign(digest, testGenesisPri[i])
		assert.Nil(t, err)
		param.Signature = sig

		input, err := param.Encode()
		assert.Nil(t, err)
		ret, err := native.TestNativeCall(t, utils.InfoSyncContractAddress, "SyncRootInfo", input, caller, caller, extra, sdb)
		assert.Nil(t, err)
		result, err := utils.PackOutputs(ABI, MethodSyncRootInfo, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
	}
	q1 := GetInfoParam{
		ChainID: CHAIN_ID,
		Height:  100,
	}
	input, err := q1.Encode()
	assert.Nil(t, err)
	extra := uint64(21000000000000)
	ret1, err := native.TestNativeCall(t, utils.InfoSyncContractAddress, "GetInfo", input, extra, sdb)
	rootInfo := new(GetInfoOutput)
	err = rootInfo.Decode(ret1)
	assert.Nil(t, err)
	assert.Equal(t, rootInfo.Info, []byte{0x01, 0x02})
	q2 := GetInfoParam{
		ChainID: CHAIN_ID,
		Height:  98,
	}
	input, err = q2.Encode()
	assert.Nil(t, err)
	ret2, err := native.TestNativeCall(t, utils.InfoSyncContractAddress, "GetInfo", input, extra, sdb)
	rootInfo = new(GetInfoOutput)
	err = rootInfo.Decode(ret2)
	assert.Nil(t, err)
	assert.Equal(t, rootInfo.Info, []byte{0x02, 0x03})
	q3 := &GetInfoHeightParam{CHAIN_ID}
	input, err = q3.Encode()
	assert.Nil(t, err)
	ret3, err := native.TestNativeCall(t, utils.InfoSyncContractAddress, "GetInfoHeight", input, extra, sdb)
	height := new(GetInfoHeightOutput)
	err = height.Decode(ret3)
	assert.Nil(t, err)
	assert.Equal(t, height.Height, uint32(100))
}

func TestReplenish(t *testing.T) {
	Init()
	param := &ReplenishParam{
		ChainID: CHAIN_ID,
		Heights: []uint32{100, 90},
	}
	extra := uint64(21000000000000)
	input, err := param.Encode()
	assert.Nil(t, err)
	_, err = native.TestNativeCall(t, utils.InfoSyncContractAddress, "Replenish", input, extra, sdb)
	assert.Nil(t, err)
}
