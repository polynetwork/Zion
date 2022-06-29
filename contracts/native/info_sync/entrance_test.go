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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	iscom "github.com/ethereum/go-ethereum/contracts/native/info_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
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
	testGenesisPeers []*node_manager.Peer
	testGenesisPri   []*ecdsa.PrivateKey
)

const (
	CHAIN_ID uint64 = 1
)

func init() {
	key, _ = crypto.GenerateKey()
	pub = &key.PublicKey

	node_manager.InitNodeManager()
	side_chain_manager.InitSideChainManager()
	InitInfoSync()
	db := rawdb.NewMemoryDatabase()
	sdb, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	testGenesisPeers, testGenesisPri = generateTestPeers(testGenesisNum)
	node_manager.StoreGenesisEpoch(sdb, testGenesisPeers)

	putSideChain()
}

func putSideChain() {
	blockNumber := big.NewInt(1)
	extra := uint64(10)
	contractRef := native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
	contract := native.NewNativeContract(sdb, contractRef)

	err := side_chain_manager.PutSideChain(contract, &side_chain_manager.SideChain{
		Router:  utils.NO_PROOF_ROUTER,
		ChainId: CHAIN_ID,
	})
	if err != nil {
		log.Fatalf("PutSideChain fail:%v", err)
		return
	}
	sideChain, err := side_chain_manager.GetSideChain(contract, CHAIN_ID)
	if err != nil {
		log.Fatalf("PutSideChain fail:%v", err)
		return
	}

	if sideChain.ChainId != CHAIN_ID {
		log.Fatalf("GetSideChain mismatch")
	}
}

func TestNoAuthSyncRootInfo(t *testing.T) {
	var err error
	param := new(iscom.SyncRootInfoParam)
	param.ChainID = CHAIN_ID
	rootInfo1 := &iscom.RootInfo{Height: 2, Info: []byte{0x01, 0x02}}
	rootInfo2 := &iscom.RootInfo{Height: 3, Info: []byte{0x02, 0x03}}
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

	input, err := utils.PackMethodWithStruct(iscom.ABI, iscom.MethodSyncRootInfo, param)
	assert.Nil(t, err)

	caller := crypto.PubkeyToAddress(*pub)
	blockNumber := big.NewInt(1)
	extra := uint64(1000)
	contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, 0+extra, nil)
	_, _, err = contractRef.NativeCall(caller, utils.InfoSyncContractAddress, input)
	assert.NotNil(t, err)
}

func TestNormalSyncRootInfo(t *testing.T) {
	var err error
	param := new(iscom.SyncRootInfoParam)
	param.ChainID = CHAIN_ID
	rootInfo1 := &iscom.RootInfo{Height: 100, Info: []byte{0x01, 0x02}}
	rootInfo2 := &iscom.RootInfo{Height: 98, Info: []byte{0x02, 0x03}}
	b1, err := rlp.EncodeToBytes(rootInfo1)
	assert.Nil(t, err)
	b2, err := rlp.EncodeToBytes(rootInfo2)
	assert.Nil(t, err)
	param.RootInfos = [][]byte{b1, b2}

	for i := 0; i < testGenesisNum; i++ {
		caller := testGenesisPeers[i].Address
		blockNumber := big.NewInt(1)
		extra := uint64(1000)
		digest, err := param.Digest()
		assert.Nil(t, err)
		sig, err := crypto.Sign(digest, testGenesisPri[i])
		assert.Nil(t, err)
		param.Signature = sig

		input, err := utils.PackMethodWithStruct(iscom.ABI, iscom.MethodSyncRootInfo, param)
		assert.Nil(t, err)
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, 0+extra, nil)
		ret, _, err := contractRef.NativeCall(caller, utils.InfoSyncContractAddress, input)
		assert.Nil(t, err)
		result, err := utils.PackOutputs(iscom.ABI, iscom.MethodSyncRootInfo, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
	}
	blockNumber := big.NewInt(1)
	extra := uint64(10)
	contractRef := native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
	contract := native.NewNativeContract(sdb, contractRef)
	rootInfo, err := iscom.GetRootInfo(contract, CHAIN_ID, 100)
	assert.Nil(t, err)
	assert.Equal(t, rootInfo, []byte{0x01, 0x02})
	rootInfo, err = iscom.GetRootInfo(contract, CHAIN_ID, 98)
	assert.Nil(t, err)
	assert.Equal(t, rootInfo, []byte{0x02, 0x03})
	h, err := iscom.GetCurrentHeight(contract, CHAIN_ID)
	assert.Nil(t, err)
	assert.Equal(t, h, uint32(100))
}

// generateTestPeer ONLY used for testing
func generateTestPeer() (*node_manager.Peer, *ecdsa.PrivateKey) {
	pk, _ := crypto.GenerateKey()
	return &node_manager.Peer{
		PubKey:  hexutil.Encode(crypto.CompressPubkey(&pk.PublicKey)),
		Address: crypto.PubkeyToAddress(pk.PublicKey),
	}, pk
}

func generateTestPeers(n int) ([]*node_manager.Peer, []*ecdsa.PrivateKey) {
	peers := make([]*node_manager.Peer, n)
	pris := make([]*ecdsa.PrivateKey, n)
	for i := 0; i < n; i++ {
		peers[i], pris[i] = generateTestPeer()
	}
	return peers, pris
}
