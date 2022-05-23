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
	acct             *ecdsa.PublicKey
	testGenesisPeers *node_manager.Peers
)

const (
	CHAIN_ID uint64 = 1
)

func init() {
	key, _ := crypto.GenerateKey()
	acct = &key.PublicKey

	node_manager.InitNodeManager()
	side_chain_manager.InitSideChainManager()
	InitInfoSync()
	db := rawdb.NewMemoryDatabase()
	sdb, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	testGenesisPeers = generateTestPeers(testGenesisNum)
	storeGenesisEpoch(sdb, testGenesisPeers)

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

	input, err := utils.PackMethodWithStruct(iscom.ABI, iscom.MethodSyncRootInfo, param)
	assert.Nil(t, err)

	caller := crypto.PubkeyToAddress(*acct)
	blockNumber := big.NewInt(1)
	extra := uint64(1000)
	contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, 0+extra, nil)
	_, _, err = contractRef.NativeCall(caller, utils.InfoSyncContractAddress, input)
	assert.Equal(t,"SyncRootInfo, CheckConsensusSigns error: invalid authority", err.Error())
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

	input, err := utils.PackMethodWithStruct(iscom.ABI, iscom.MethodSyncRootInfo, param)
	assert.Nil(t, err)

	for i := 0; i < testGenesisNum; i++ {
		caller := testGenesisPeers.List[i].Address
		blockNumber := big.NewInt(1)
		extra := uint64(1000)
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
func generateTestPeer() *node_manager.PeerInfo {
	pk, _ := crypto.GenerateKey()
	return &node_manager.PeerInfo{
		PubKey:  hexutil.Encode(crypto.CompressPubkey(&pk.PublicKey)),
		Address: crypto.PubkeyToAddress(pk.PublicKey),
	}
}

func generateTestPeers(n int) *node_manager.Peers {
	peers := &node_manager.Peers{List: make([]*node_manager.PeerInfo, n)}
	for i := 0; i < n; i++ {
		peers.List[i] = generateTestPeer()
	}
	return peers
}

func storeGenesisEpoch(s *state.StateDB, peers *node_manager.Peers) (*node_manager.EpochInfo, error) {
	cache := (*state.CacheDB)(s)
	epoch := &node_manager.EpochInfo{
		ID:          node_manager.StartEpochID,
		Peers:       peers,
		StartHeight: 0,
	}

	// store current epoch and epoch info
	if err := setEpoch(cache, epoch); err != nil {
		return nil, err
	}

	// store current hash
	curKey := curEpochKey()
	cache.Put(curKey, epoch.Hash().Bytes())

	// store genesis epoch id to list
	value, err := rlp.EncodeToBytes(&node_manager.HashList{List: []common.Hash{epoch.Hash()}})
	if err != nil {
		return nil, err
	}
	proposalKey := proposalsKey(epoch.ID)
	cache.Put(proposalKey, value)

	// store genesis epoch proof
	key := epochProofKey(node_manager.EpochProofHash(epoch.ID))
	cache.Put(key, epoch.Hash().Bytes())

	return epoch, nil
}

func setEpoch(s *state.CacheDB, epoch *node_manager.EpochInfo) error {
	hash := epoch.Hash()
	key := epochKey(hash)

	value, err := rlp.EncodeToBytes(epoch)
	if err != nil {
		return err
	}

	s.Put(key, value)
	return nil
}

func epochKey(epochHash common.Hash) []byte {
	return utils.ConcatKey(utils.NodeManagerContractAddress, []byte("st_epoch"), epochHash.Bytes())
}

func curEpochKey() []byte {
	return utils.ConcatKey(utils.NodeManagerContractAddress, []byte("st_cur_epoch"), []byte("1"))
}

func epochProofKey(proofHashKey common.Hash) []byte {
	return utils.ConcatKey(utils.NodeManagerContractAddress, []byte("st_proof"), proofHashKey.Bytes())
}

func proposalsKey(epochID uint64) []byte {
	return utils.ConcatKey(utils.NodeManagerContractAddress, []byte("st_proposal"), utils.GetUint64Bytes(epochID))
}
