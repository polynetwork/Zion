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
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	testStateDB  *state.StateDB
	testEmptyCtx *native.NativeContract

	testSupplyGas    uint64 = 100000000000000000
	testGenesisNum   int    = 4
	testCaller       common.Address
	testGenesisEpoch *EpochInfo
)

func TestMain(m *testing.M) {
	db := rawdb.NewMemoryDatabase()
	testStateDB, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	testEmptyCtx = native.NewNativeContract(testStateDB, nil)
	testGenesisPeers := generateTestPeers(testGenesisNum)
	testGenesisEpoch, _ = StoreGenesisEpoch(testStateDB, testGenesisPeers)
	InitNodeManager()

	os.Exit(m.Run())
}

func TestPropose(t *testing.T) {
	type TestCase struct {
		BlockNum      int
		StartHeight   uint64
		Payload       []byte
		BeforeHandler func(c *TestCase, ctx *native.NativeContract)
		AfterHandler  func(c *TestCase, ctx *native.NativeContract)
		Index         int
		Expect        error
	}

	cases := []*TestCase{
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       1,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := generateTestPeers(2)
				input := &MethodProposeInput{StartHeight: c.StartHeight, Peers: peers}
				c.Payload, _ = input.Encode()
				delEpoch(ctx, testGenesisEpoch.Hash())
			},
			Expect: ErrEpochNotExist,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       2,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := generateTestPeers(2)
				input := &MethodProposeInput{StartHeight: c.StartHeight, Peers: peers}
				c.Payload, _ = input.Encode()
				testCaller = generateTestAddress(78)
			},
			Expect: ErrInvalidAuthority,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				input := &MethodProposeInput{StartHeight: 0, Peers: nil}
				payload, _ := input.Encode()
				c.Payload = payload[0 : len(payload)-2]
			},
			Expect: ErrInvalidInput,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       4,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := generateTestPeers(3)
				peers.List = nil
				input := &MethodProposeInput{StartHeight: 0, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrInvalidPeers,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       5,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := generateTestPeers(MinProposalPeersLen - 1)
				input := &MethodProposeInput{StartHeight: 0, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrPeersNum,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       6,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := generateTestPeers(MaxProposalPeersLen + 1)
				input := &MethodProposeInput{StartHeight: 0, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrPeersNum,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       7,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := generateTestPeers(MinProposalPeersLen + 1)
				peers.List[0].PubKey = "0ruf8nkj"
				input := &MethodProposeInput{StartHeight: 0, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrInvalidPubKey,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       8,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := testGenesisEpoch.Peers.Copy()
				newPeers := generateTestPeers(len(peers.List))
				for i, _ := range peers.List {
					if i%2 == 0 {
						peers.List[i] = newPeers.List[i]
					}
				}
				input := &MethodProposeInput{StartHeight: 0, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrOldParticipantsNumber,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       9,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := testGenesisEpoch.Peers.Copy()
				peers.List = append(peers.List, generateTestPeers(1).List...)
				c.StartHeight += MinEpochValidPeriod - 1
				input := &MethodProposeInput{StartHeight: c.StartHeight, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrProposalStartHeight,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       10,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := testGenesisEpoch.Peers.Copy()
				peers.List = append(peers.List, generateTestPeers(1).List...)
				c.StartHeight += MaxEpochValidPeriod + 10
				input := &MethodProposeInput{StartHeight: c.StartHeight, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrProposalStartHeight,
		},
		{
			BlockNum:    3,
			StartHeight: MinEpochValidPeriod + 10,
			Index:       11,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := testGenesisEpoch.Peers.Copy()
				peers.List = append(peers.List, generateTestPeers(1).List...)
				input := &MethodProposeInput{StartHeight: c.StartHeight, Peers: peers}
				sort.Sort(peers)
				epoch := &EpochInfo{
					ID:          testGenesisEpoch.ID + 1,
					StartHeight: c.StartHeight,
					Peers:       peers,
				}
				storeProposal(ctx, epoch.ID, ctx.ContractRef().TxOrigin(), epoch.Hash())
				c.Payload, _ = input.Encode()
			},
			Expect: ErrDuplicateProposal,
		},
		{
			BlockNum:    3,
			StartHeight: MinEpochValidPeriod + 10,
			Index:       12,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := testGenesisEpoch.Peers.Copy()
				peers.List = append(peers.List, generateTestPeers(1).List...)
				input := &MethodProposeInput{StartHeight: c.StartHeight, Peers: peers}
				sort.Sort(peers)
				proposer := ctx.ContractRef().TxOrigin()
				epoch1 := &EpochInfo{
					ID:          testGenesisEpoch.ID + 1,
					StartHeight: c.StartHeight + 2,
					Peers:       peers,
				}
				epoch2 := &EpochInfo{
					ID:          testGenesisEpoch.ID + 1,
					StartHeight: c.StartHeight + 3,
					Peers:       peers,
				}
				epoch3 := &EpochInfo{
					ID:          testGenesisEpoch.ID + 1,
					StartHeight: c.StartHeight + 4,
					Peers:       peers,
				}
				storeProposal(ctx, epoch1.ID, proposer, epoch1.Hash())
				storeProposal(ctx, epoch2.ID, proposer, epoch2.Hash())
				storeProposal(ctx, epoch3.ID, proposer, epoch3.Hash())
				c.Payload, _ = input.Encode()
			},
			Expect: ErrProposalsNum,
		},
		{
			BlockNum:    3,
			StartHeight: MinEpochValidPeriod + 10,
			Index:       13,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := testGenesisEpoch.Peers.Copy()
				peers.List = append(peers.List, generateTestPeers(1).List...)
				input := &MethodProposeInput{StartHeight: c.StartHeight, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: nil,
		},
	}

	for _, v := range cases {
		resetTestContext()
		ctx := generateNativeContract(testCaller, v.BlockNum)
		if v.BeforeHandler != nil {
			v.BeforeHandler(v, ctx)
		}
		_, _, err := ctx.ContractRef().NativeCall(testCaller, this, v.Payload)
		assert.Equal(t, v.Expect, err)
		if v.AfterHandler != nil {
			v.AfterHandler(v, ctx)
		}
	}
}

func generateNativeContractRef(origin common.Address, blockNum int) *native.ContractRef {
	token := make([]byte, common.HashLength)
	rand.Read(token)
	hash := common.BytesToHash(token)
	return native.NewContractRef(testStateDB, origin, origin, big.NewInt(int64(blockNum)), hash, testSupplyGas, nil)
}

func generateNativeContract(origin common.Address, blockNum int) *native.NativeContract {
	ref := generateNativeContractRef(origin, blockNum)
	return native.NewNativeContract(testStateDB, ref)
}

func resetTestContext() {
	db := rawdb.NewMemoryDatabase()
	testStateDB, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	testEmptyCtx = native.NewNativeContract(testStateDB, nil)
	testGenesisPeers := generateTestPeers(testGenesisNum)
	testGenesisEpoch, _ = StoreGenesisEpoch(testStateDB, testGenesisPeers)
	testCaller = testGenesisEpoch.Peers.List[0].Address
}

// generateTestPeer ONLY used for testing
func generateTestPeer() *PeerInfo {
	pk, _ := crypto.GenerateKey()
	return &PeerInfo{
		PubKey:  hexutil.Encode(crypto.CompressPubkey(&pk.PublicKey)),
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
