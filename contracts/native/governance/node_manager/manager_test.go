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
	"crypto/rand"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	testStateDB   *state.StateDB
	testCacheDB   *state.CacheDB
	testEmptyCtx  *native.NativeContract
	testSupplyGas uint64 = 100000000000000000
)

func TestMain(m *testing.M) {
	db := rawdb.NewMemoryDatabase()
	testStateDB, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	testCacheDB = (*state.CacheDB)(testStateDB)
	testEmptyCtx = native.NewNativeContract(testStateDB, nil)

	InitNodeManager()

	os.Exit(m.Run())
}

func TestPropose(t *testing.T) {
	caller := generateTestAddress(12)
	ctx := generateNativeContractRef(caller, 3)

	startHeight := uint64(2)
	peers := generateTestPeers(15)
	payload, err := new(MethodProposeInput).Encode(startHeight, peers)
	assert.NoError(t, err)

	ret, gasLeft, err := ctx.NativeCall(caller, this, payload)
	assert.NoError(t, err)
	assert.Equal(t, testSupplyGas-gasTable[MethodPropose], gasLeft)
	assert.Equal(t, utils.ByteSuccess, ret)
}

func generateNativeContractRef(origin common.Address, blockNum int) *native.ContractRef {
	token := make([]byte, common.HashLength)
	rand.Read(token)
	hash := common.BytesToHash(token)
	return native.NewContractRef(testStateDB, origin, origin, big.NewInt(int64(blockNum)), hash, testSupplyGas, nil)
}

// generateTestPeer ONLY used for testing
func generateTestPeer() *PeerInfo {
	pk, _ := crypto.GenerateKey()
	return &PeerInfo{
		PubKey:  common.Bytes2Hex(crypto.CompressPubkey(&pk.PublicKey)),
		Address: crypto.PubkeyToAddress(pk.PublicKey),
	}
}

func generateTestPeers(n int) *Peers {
	peers := &Peers{List: make([]*PeerInfo, n)}
	for i := 0; i < n; i++ {
		peers.List[i] = generateTestPeer()
	}
	return peers
}

func generateTestEpochInfo(id, height uint64, peersNum int) *EpochInfo {
	epoch := new(EpochInfo)
	epoch.ID = id
	epoch.StartHeight = height
	epoch.Peers = generateTestPeers(peersNum)
	return epoch
}

func generateTestHash(n int) common.Hash {
	data := big.NewInt(int64(n))
	return common.BytesToHash(data.Bytes())
}

func generateTestHashList(n int) *HashList {
	data := &HashList{List: make([]common.Hash, n)}
	for i := 0; i < n; i++ {
		data.List[i] = generateTestHash(i + 1)
	}
	return data
}

func generateTestAddress(n int) common.Address {
	data := big.NewInt(int64(n))
	return common.BytesToAddress(data.Bytes())
}

func generateTestAddressList(n int) *AddressList {
	data := &AddressList{List: make([]common.Address, n)}
	for i := 0; i < n; i++ {
		data.List[i] = generateTestAddress(i + 1)
	}
	return data
}
