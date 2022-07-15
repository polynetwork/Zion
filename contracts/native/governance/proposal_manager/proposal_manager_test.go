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

package proposal_manager

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var (
	sdb              *state.StateDB
	testGenesisNum   = 4
	acct             *ecdsa.PublicKey
	testGenesisPeers []*node_manager.Peer
)

func init() {
	key, _ := crypto.GenerateKey()
	acct = &key.PublicKey

	node_manager.InitNodeManager()
	InitProposalManager()
	db := rawdb.NewMemoryDatabase()
	sdb, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	testGenesisPeers = node_manager.GenerateTestPeers(testGenesisNum)
	node_manager.StoreCommunityInfo(sdb, big.NewInt(2000), common.EmptyAddress)
	node_manager.StoreGenesisEpoch(sdb, testGenesisPeers)
	node_manager.StoreGenesisGlobalConfig(sdb)
}

func TestUpdateNodeManagerGlobalConfig(t *testing.T) {
	blockNumber := big.NewInt(1)
	extra := uint64(10)
	contractRef := native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
	contract := native.NewNativeContract(sdb, contractRef)

	globalConfig, err := node_manager.GetGlobalConfigImpl(contract)
	assert.Nil(t, err)
	assert.Equal(t, globalConfig.MaxDescLength, node_manager.GenesisMaxDescLength)
	assert.Equal(t, globalConfig.BlockPerEpoch, node_manager.GenesisBlockPerEpoch)
	assert.Equal(t, globalConfig.MaxCommissionChange, node_manager.GenesisMaxCommissionChange)
	assert.Equal(t, globalConfig.MinInitialStake, node_manager.GenesisMinInitialStake)
	assert.Equal(t, globalConfig.VoterValidatorNum, node_manager.GenesisVoterValidatorNum)
	assert.Equal(t, globalConfig.ConsensusValidatorNum, node_manager.GenesisConsensusValidatorNum)

	// update global config
	for _, p := range testGenesisPeers {
		param1 := new(UpdateNodeManagerGlobalConfigParam)
		param1.ExpireHeight = 100
		param1.BlockPerEpoch = 10000
		input, err := param1.Encode()
		assert.Nil(t, err)
		contractRef := native.NewContractRef(sdb, p.Address, p.Address, blockNumber, common.Hash{}, extra, nil)
		_, _, err = contractRef.NativeCall(p.Address, utils.ProposalManagerContractAddress, input)
		assert.Nil(t, err)
	}

	globalConfig, err = node_manager.GetGlobalConfigImpl(contract)
	assert.Nil(t, err)
	assert.Equal(t, globalConfig.BlockPerEpoch, new(big.Int).SetUint64(10000))

	// not consensus address
	caller := crypto.PubkeyToAddress(*acct)
	param2 := new(UpdateNodeManagerGlobalConfigParam)
	param2.ExpireHeight = 100
	param2.BlockPerEpoch = 10000
	input, err := param2.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(caller, utils.ProposalManagerContractAddress, input)
	assert.NotNil(t, err)
}