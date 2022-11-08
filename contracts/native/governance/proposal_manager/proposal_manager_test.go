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
	"github.com/ethereum/go-ethereum/contracts/native/contract"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	sdb              *state.StateDB
	testGenesisNum   = 4
	acct             *ecdsa.PublicKey
	testGenesisPeers []common.Address
)

func init() {
	key, _ := crypto.GenerateKey()
	acct = &key.PublicKey

	node_manager.InitNodeManager()
	InitProposalManager()
	sdb = native.NewTestStateDB()
	testGenesisPeers, _ = native.GenerateTestPeers(testGenesisNum)
	node_manager.StoreCommunityInfo(sdb, big.NewInt(2000), common.EmptyAddress)
	node_manager.StoreGenesisEpoch(sdb, testGenesisPeers, testGenesisPeers)
	node_manager.StoreGenesisGlobalConfig(sdb)
}

func TestProposalManager(t *testing.T) {
	extra := uint64(21000000000000)
	contractRef := native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, common.Big1, common.Hash{}, extra, nil)
	c := native.NewNativeContract(sdb, contractRef)

	globalConfig, err := node_manager.GetGlobalConfigImpl(c)
	assert.Nil(t, err)
	assert.Equal(t, globalConfig.BlockPerEpoch, node_manager.GenesisBlockPerEpoch)
	assert.Equal(t, globalConfig.MaxCommissionChange, node_manager.GenesisMaxCommissionChange)
	assert.Equal(t, globalConfig.MinInitialStake, node_manager.GenesisMinInitialStake)
	assert.Equal(t, globalConfig.VoterValidatorNum, node_manager.GenesisVoterValidatorNum)
	assert.Equal(t, globalConfig.ConsensusValidatorNum, node_manager.GenesisConsensusValidatorNum)
	assert.Equal(t, globalConfig.MinProposalStake, node_manager.GenesisMinProposalStake)

	communityInfo, err := node_manager.GetCommunityInfoImpl(c)
	assert.Nil(t, err)
	assert.Equal(t, communityInfo.CommunityRate, big.NewInt(2000))
	assert.Equal(t, communityInfo.CommunityAddress, common.EmptyAddress)

	sdb.SetBalance(common.EmptyAddress, new(big.Int).Mul(big.NewInt(100000), params.ZNT1))
	value := new(big.Int).Mul(big.NewInt(1000), params.ZNT1)
	// Propose
	for i := 0; i < ProposalListLen; i++ {
		param := new(ProposeParam)
		param.Content = make([]byte, 4000)
		input, err := param.Encode()
		assert.Nil(t, err)
		err = contract.NativeTransfer(contractRef.StateDB(), common.EmptyAddress, this, value)
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.ProposalManagerContractAddress, "Propose", input, value, common.EmptyAddress, common.EmptyAddress, 1, extra, sdb)
		assert.Nil(t, err)
	}

	// Propose config
	param2 := new(ProposeConfigParam)
	globalConfig.VoterValidatorNum = 2
	param2.Content, err = rlp.EncodeToBytes(globalConfig)
	assert.Nil(t, err)
	input, err := param2.Encode()
	assert.Nil(t, err)
	err = contract.NativeTransfer(contractRef.StateDB(), common.EmptyAddress, this, value)
	assert.Nil(t, err)
	_, err = native.TestNativeCall(t, utils.ProposalManagerContractAddress, "ProposeConfig", input, value, common.EmptyAddress, common.EmptyAddress, 1, extra, sdb)
	assert.Nil(t, err)

	for i := 0; i < ProposalListLen-1; i++ {
		param := new(ProposeConfigParam)
		globalConfig.VoterValidatorNum = 3
		param.Content, err = rlp.EncodeToBytes(globalConfig)
		assert.Nil(t, err)
		input, err := param.Encode()
		assert.Nil(t, err)
		err = contract.NativeTransfer(contractRef.StateDB(), common.EmptyAddress, this, value)
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.ProposalManagerContractAddress, "ProposeConfig", input, value, common.EmptyAddress, common.EmptyAddress, 1, extra, sdb)
		assert.Nil(t, err)
	}

	// Propose community
	param3 := new(ProposeCommunityParam)
	communityInfo.CommunityRate = new(big.Int).SetUint64(1000)
	param3.Content, err = rlp.EncodeToBytes(communityInfo)
	assert.Nil(t, err)
	input, err = param3.Encode()
	assert.Nil(t, err)
	err = contract.NativeTransfer(contractRef.StateDB(), common.EmptyAddress, this, value)
	assert.Nil(t, err)
	_, err = native.TestNativeCall(t, utils.ProposalManagerContractAddress, "ProposeCommunity", input, value, common.EmptyAddress, common.EmptyAddress, 1, extra, sdb)
	assert.Nil(t, err)

	for i := 0; i < ProposalListLen-1; i++ {
		param := new(ProposeCommunityParam)
		communityInfo.CommunityRate = new(big.Int).SetUint64(3000)
		param.Content, err = rlp.EncodeToBytes(communityInfo)
		assert.Nil(t, err)
		input, err := param.Encode()
		assert.Nil(t, err)
		err = contract.NativeTransfer(contractRef.StateDB(), common.EmptyAddress, this, value)
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.ProposalManagerContractAddress, "ProposeCommunity", input, value, common.EmptyAddress, common.EmptyAddress, 1, extra, sdb)
		assert.Nil(t, err)
	}

	// get proposal list
	param9 := new(GetProposalListParam)
	input, err = param9.Encode()
	assert.Nil(t, err)
	ret4, err := native.TestNativeCall(t, utils.ProposalManagerContractAddress, "GetProposalList", input, new(big.Int), common.EmptyAddress, common.EmptyAddress, 1, extra, sdb)
	assert.Nil(t, err)
	proposalList := new(ProposalList)
	err = proposalList.Decode(ret4)
	assert.Nil(t, err)
	assert.Equal(t, len(proposalList.ProposalList), ProposalListLen)
	param10 := new(GetConfigProposalListParam)
	input, err = param10.Encode()
	assert.Nil(t, err)
	ret5, err := native.TestNativeCall(t, utils.ProposalManagerContractAddress, "GetConfigProposalList", input, new(big.Int), common.EmptyAddress, common.EmptyAddress, 1, extra, sdb)
	assert.Nil(t, err)
	configProposalList := new(ConfigProposalList)
	err = configProposalList.Decode(ret5)
	assert.Nil(t, err)
	assert.Equal(t, len(configProposalList.ConfigProposalList), ProposalListLen)
	param11 := new(GetCommunityProposalListParam)
	input, err = param11.Encode()
	assert.Nil(t, err)
	ret6, err := native.TestNativeCall(t, utils.ProposalManagerContractAddress, "GetCommunityProposalList", input, new(big.Int), common.EmptyAddress, common.EmptyAddress, 1, extra, sdb)
	assert.Nil(t, err)
	communityProposalList := new(CommunityProposalList)
	err = communityProposalList.Decode(ret6)
	assert.Nil(t, err)
	assert.Equal(t, len(communityProposalList.CommunityProposalList), ProposalListLen)

	// vote
	param4 := new(VoteProposalParam)
	param4.ID = new(big.Int).SetUint64(0)
	assert.Nil(t, err)
	input, err = param4.Encode()
	assert.Nil(t, err)
	for i := 0; i < testGenesisNum; i++ {
		_, err = native.TestNativeCall(t, utils.ProposalManagerContractAddress, "VoteProposal", input, new(big.Int), testGenesisPeers[i], testGenesisPeers[i], 1, extra, sdb)
		assert.Nil(t, err)
	}
	param5 := new(VoteProposalParam)
	param5.ID = new(big.Int).SetUint64(20)
	assert.Nil(t, err)
	input, err = param5.Encode()
	assert.Nil(t, err)
	for i := 0; i < testGenesisNum; i++ {
		_, err = native.TestNativeCall(t, utils.ProposalManagerContractAddress, "VoteProposal", input, new(big.Int), testGenesisPeers[i], testGenesisPeers[i], 1, extra, sdb)
		assert.Nil(t, err)
	}
	param14 := new(VoteProposalParam)
	param14.ID = new(big.Int).SetUint64(40)
	assert.Nil(t, err)
	input, err = param14.Encode()
	assert.Nil(t, err)
	for i := 0; i < testGenesisNum; i++ {
		_, err = native.TestNativeCall(t, utils.ProposalManagerContractAddress, "VoteProposal", input, new(big.Int), testGenesisPeers[i], testGenesisPeers[i], 1, extra, sdb)
		assert.Nil(t, err)
	}

	// get proposal
	param6 := new(GetProposalParam)
	param6.ID = new(big.Int).SetUint64(0)
	assert.Nil(t, err)
	input, err = param6.Encode()
	assert.Nil(t, err)
	ret1, err := native.TestNativeCall(t, utils.ProposalManagerContractAddress, "GetProposal", input, new(big.Int), common.EmptyAddress, common.EmptyAddress, 1, extra, sdb)
	assert.Nil(t, err)
	proposal1 := new(Proposal)
	err = proposal1.Decode(ret1)
	assert.Nil(t, err)
	assert.Equal(t, proposal1.Status, PASS)
	param7 := new(GetProposalParam)
	param7.ID = new(big.Int).SetUint64(20)
	assert.Nil(t, err)
	input, err = param7.Encode()
	assert.Nil(t, err)
	ret2, err := native.TestNativeCall(t, utils.ProposalManagerContractAddress, "GetProposal", input, new(big.Int), common.EmptyAddress, common.EmptyAddress, 1, extra, sdb)
	assert.Nil(t, err)
	proposal2 := new(Proposal)
	err = proposal2.Decode(ret2)
	assert.Nil(t, err)
	assert.Equal(t, proposal2.Status, PASS)
	param8 := new(GetProposalParam)
	param8.ID = new(big.Int).SetUint64(22)
	assert.Nil(t, err)
	input, err = param8.Encode()
	assert.Nil(t, err)
	ret3, err := native.TestNativeCall(t, utils.ProposalManagerContractAddress, "GetProposal", input, new(big.Int), common.EmptyAddress, common.EmptyAddress, 1, extra, sdb)
	assert.Nil(t, err)
	proposal3 := new(Proposal)
	err = proposal3.Decode(ret3)
	assert.Nil(t, err)
	assert.Equal(t, proposal3.Status, FAIL)
	param12 := new(GetProposalParam)
	param12.ID = new(big.Int).SetUint64(40)
	assert.Nil(t, err)
	input, err = param12.Encode()
	assert.Nil(t, err)
	ret7, err := native.TestNativeCall(t, utils.ProposalManagerContractAddress, "GetProposal", input, new(big.Int), common.EmptyAddress, common.EmptyAddress, 1, extra, sdb)
	assert.Nil(t, err)
	proposal4 := new(Proposal)
	err = proposal4.Decode(ret7)
	assert.Nil(t, err)
	assert.Equal(t, proposal4.Status, PASS)
	param13 := new(GetProposalParam)
	param13.ID = new(big.Int).SetUint64(42)
	assert.Nil(t, err)
	input, err = param13.Encode()
	assert.Nil(t, err)
	ret8, err := native.TestNativeCall(t, utils.ProposalManagerContractAddress, "GetProposal", input, new(big.Int), common.EmptyAddress, common.EmptyAddress, 1, extra, sdb)
	assert.Nil(t, err)
	proposal5 := new(Proposal)
	err = proposal5.Decode(ret8)
	assert.Nil(t, err)
	assert.Equal(t, proposal5.Status, FAIL)

	// check
	globalConfig, err = node_manager.GetGlobalConfigImpl(c)
	assert.Nil(t, err)
	assert.Equal(t, globalConfig.VoterValidatorNum, uint64(2))
	communityInfo, err = node_manager.GetCommunityInfoImpl(c)
	assert.Nil(t, err)
	assert.Equal(t, communityInfo.CommunityRate, big.NewInt(1000))
	assert.Equal(t, sdb.GetBalance(common.EmptyAddress), new(big.Int).Mul(big.NewInt(81000), params.ZNT1))
	assert.Equal(t, sdb.GetBalance(communityInfo.CommunityAddress), new(big.Int).Mul(big.NewInt(81000), params.ZNT1))
}
