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
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var (
	sdb              *state.StateDB
	testGenesisNum   = 4
	acct             *ecdsa.PublicKey
	testGenesisPeers []*Peer
)

func init() {
	key, _ := crypto.GenerateKey()
	acct = &key.PublicKey

	InitNodeManager()
	db := rawdb.NewMemoryDatabase()
	sdb, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	testGenesisPeers = generateTestPeers(testGenesisNum)
	StoreCommunityInfo(sdb, big.NewInt(20), common.EmptyAddress)
	StoreGenesisEpoch(sdb, testGenesisPeers)
	StoreGenesisGlobalConfig(sdb)
}

func TestCheckGenesis(t *testing.T) {
	// check get spec methodID
	m := GetSpecMethodID()
	assert.Equal(t, m["fe6f86f8"], true)
	assert.Equal(t, m["083c6323"], true)

	blockNumber := big.NewInt(1)
	extra := uint64(10)
	contractRef := native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
	contract := native.NewNativeContract(sdb, contractRef)

	globalConfig, err := getGlobalConfig(contract)
	assert.Nil(t, err)
	assert.Equal(t, globalConfig.MaxDescLength, GenesisMaxDescLength)
	assert.Equal(t, globalConfig.BlockPerEpoch, GenesisBlockPerEpoch)
	assert.Equal(t, globalConfig.MaxCommission, GenesisMaxCommission)
	assert.Equal(t, globalConfig.MinInitialStake, GenesisMinInitialStake)
	assert.Equal(t, globalConfig.VoterValidatorNum, GenesisVoterValidatorNum)
	assert.Equal(t, globalConfig.ConsensusValidatorNum, GenesisConsensusValidatorNum)

	communityInfo, err := getCommunityInfo(contract)
	assert.Nil(t, err)
	assert.Equal(t, communityInfo.CommunityRate, big.NewInt(20))
	assert.Equal(t, communityInfo.CommunityAddress, common.EmptyAddress)

	epochInfo, err := getCurrentEpochInfo(contract)
	assert.Nil(t, err)
	assert.Equal(t, epochInfo.ID, common.Big1)

	// check query method
	param1 := new(GetGlobalConfigParam)
	input, err := param1.Encode()
	assert.Nil(t, err)
	ret, _, err := contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	globalConfig2 := new(GlobalConfig)
	err = globalConfig2.Decode(ret)
	assert.Nil(t, err)
	assert.Equal(t, globalConfig2.MaxDescLength, GenesisMaxDescLength)
	assert.Equal(t, globalConfig2.BlockPerEpoch, GenesisBlockPerEpoch)
	assert.Equal(t, globalConfig2.MaxCommission, GenesisMaxCommission)
	assert.Equal(t, globalConfig2.MinInitialStake, GenesisMinInitialStake)
	assert.Equal(t, globalConfig2.VoterValidatorNum, GenesisVoterValidatorNum)
	assert.Equal(t, globalConfig2.ConsensusValidatorNum, GenesisConsensusValidatorNum)

	param2 := new(GetCommunityInfoParam)
	input, err = param2.Encode()
	assert.Nil(t, err)
	ret, _, err = contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	communityInfo2 := new(CommunityInfo)
	err = communityInfo2.Decode(ret)
	assert.Nil(t, err)
	assert.Equal(t, communityInfo2.CommunityRate, big.NewInt(20))
	assert.Equal(t, communityInfo2.CommunityAddress, common.EmptyAddress)

	param3 := new(GetCurrentEpochInfoParam)
	input, err = param3.Encode()
	assert.Nil(t, err)
	ret, _, err = contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	currentEpochInfo := new(EpochInfo)
	err = currentEpochInfo.Decode(ret)
	assert.Nil(t, err)
	assert.Equal(t, currentEpochInfo.ID, big.NewInt(1))
	assert.Equal(t, currentEpochInfo.StartHeight, big.NewInt(0))
	assert.Equal(t, uint64(len(currentEpochInfo.Validators)), GenesisConsensusValidatorNum)
	assert.Equal(t, uint64(len(currentEpochInfo.Voters)), GenesisVoterValidatorNum)
}

func TestStake(t *testing.T) {
	blockNumber := big.NewInt(400000)
	extra := uint64(10)
	contractRefQuery := native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
	contractQuery := native.NewNativeContract(sdb, contractRefQuery)

	type ValidatorKey struct {
		ConsensusPubkey string
		Dec             []byte
		Address         common.Address
	}
	// create validator
	loop := 6
	validatorsKey := make([]*ValidatorKey, 0, loop)
	for i := 0; i < loop; i++ {
		pk, _ := crypto.GenerateKey()
		caller := crypto.PubkeyToAddress(*acct)
		sdb.SetBalance(caller, new(big.Int).Mul(big.NewInt(1000000), params.ZNT1))
		param := new(CreateValidatorParam)
		param.ConsensusPubkey = hexutil.Encode(crypto.CompressPubkey(&pk.PublicKey))
		param.ProposalAddress = caller
		param.InitStake = new(big.Int).Mul(big.NewInt(100000), params.ZNT1)
		param.Commission = new(big.Int).SetUint64(20)
		param.Desc = "test"
		validatorsKey = append(validatorsKey, &ValidatorKey{param.ConsensusPubkey, crypto.CompressPubkey(&pk.PublicKey), caller})
		input, err := param.Encode()
		assert.Nil(t, err)
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
		_, _, err = contractRef.NativeCall(caller, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
	}
	// check
	allValidators, err := GetAllValidators(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, len(allValidators.AllValidators), loop)

	//stake
	pkStake, _ := crypto.GenerateKey()
	staker := &pkStake.PublicKey
	stakeAddress := crypto.PubkeyToAddress(*staker)
	sdb.SetBalance(stakeAddress, new(big.Int).Mul(big.NewInt(1000000), params.ZNT1))
	param1 := new(StakeParam)
	param1.ConsensusPubkey = validatorsKey[0].ConsensusPubkey
	param1.Amount = new(big.Int).Mul(big.NewInt(10000), params.ZNT1)
	input, err := param1.Encode()
	assert.Nil(t, err)
	contractRef := native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	// check
	validator, _, err := GetValidator(contractQuery, validatorsKey[0].Dec)
	assert.Nil(t, err)
	assert.Equal(t, validator.TotalStake, new(big.Int).Mul(big.NewInt(110000), params.ZNT1))
	assert.Equal(t, validator.SelfStake, new(big.Int).Mul(big.NewInt(100000), params.ZNT1))
	assert.Equal(t, validator.Status, Unlock)
	assert.Equal(t, validator.Commission.Rate, new(big.Int).SetUint64(20))
	assert.Equal(t, validator.UnlockHeight, common.Big0)
	totalPool, err := GetTotalPool(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, totalPool, new(big.Int).Mul(big.NewInt(610000), params.ZNT1))

	// unstake
	param2 := new(UnStakeParam)
	param2.ConsensusPubkey = validatorsKey[0].ConsensusPubkey
	param2.Amount = new(big.Int).Mul(big.NewInt(1000), params.ZNT1)
	input, err = param2.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	validator, _, err = GetValidator(contractQuery, validatorsKey[0].Dec)
	assert.Nil(t, err)
	assert.Equal(t, validator.TotalStake, new(big.Int).Mul(big.NewInt(109000), params.ZNT1))
	assert.Equal(t, validator.SelfStake, new(big.Int).Mul(big.NewInt(100000), params.ZNT1))
	assert.Equal(t, validator.Status, Unlock)
	assert.Equal(t, validator.UnlockHeight, common.Big0)
	totalPool, err = GetTotalPool(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, totalPool, new(big.Int).Mul(big.NewInt(609000), params.ZNT1))
	assert.Equal(t, sdb.GetBalance(stakeAddress), new(big.Int).Mul(big.NewInt(991000), params.ZNT1))

	// change epoch
	input, err = utils.PackMethod(ABI, MethodChangeEpoch)
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	epochInfo, err := getCurrentEpochInfo(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, epochInfo.ID, common.Big2)
	assert.Equal(t, epochInfo.StartHeight, new(big.Int).SetUint64(400000))
	assert.Equal(t, len(epochInfo.Validators), 4)
	assert.Equal(t, len(epochInfo.Voters), 4)
	validator, _, err = GetValidator(contractQuery, validatorsKey[0].Dec)
	assert.Nil(t, err)
	assert.Equal(t, validator.Status, Lock)
	totalPool, err = GetTotalPool(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, totalPool, new(big.Int).Mul(big.NewInt(609000), params.ZNT1))

	// unstake
	param3 := new(UnStakeParam)
	param3.ConsensusPubkey = validatorsKey[0].ConsensusPubkey
	param3.Amount = new(big.Int).Mul(big.NewInt(1000), params.ZNT1)
	input, err = param3.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	validator, _, err = GetValidator(contractQuery, validatorsKey[0].Dec)
	assert.Nil(t, err)
	assert.Equal(t, validator.TotalStake, new(big.Int).Mul(big.NewInt(108000), params.ZNT1))
	assert.Equal(t, validator.SelfStake, new(big.Int).Mul(big.NewInt(100000), params.ZNT1))
	assert.Equal(t, validator.Status, Lock)
	assert.Equal(t, validator.UnlockHeight, common.Big0)

	// withdraw
	input, err = utils.PackMethod(ABI, MethodWithdraw)
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	totalPool, err = GetTotalPool(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, totalPool, new(big.Int).Mul(big.NewInt(609000), params.ZNT1))
	blockNumber = big.NewInt(800000)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	totalPool, err = GetTotalPool(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, totalPool, new(big.Int).Mul(big.NewInt(608000), params.ZNT1))

	// update validator
	param4 := new(UpdateValidatorParam)
	param4.ConsensusPubkey = validatorsKey[0].ConsensusPubkey
	param4.Desc = "test2"
	input, err = param4.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, validatorsKey[0].Address, validatorsKey[0].Address, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(validatorsKey[0].Address, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	param5 := new(UpdateCommissionParam)
	param5.ConsensusPubkey = validatorsKey[0].ConsensusPubkey
	param5.Commission = new(big.Int).SetUint64(30)
	input, err = param5.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, validatorsKey[0].Address, validatorsKey[0].Address, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(validatorsKey[0].Address, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	validator, _, err = GetValidator(contractQuery, validatorsKey[0].Dec)
	assert.Nil(t, err)
	assert.Equal(t, validator.TotalStake, new(big.Int).Mul(big.NewInt(108000), params.ZNT1))
	assert.Equal(t, validator.SelfStake, new(big.Int).Mul(big.NewInt(100000), params.ZNT1))
	assert.Equal(t, validator.ProposalAddress, validatorsKey[0].Address)
	assert.Equal(t, validator.Status, Lock)
	assert.Equal(t, validator.Commission.Rate, new(big.Int).SetUint64(30))
	assert.Equal(t, validator.Commission.UpdateHeight, new(big.Int).SetUint64(800000))
	assert.Equal(t, validator.UnlockHeight, common.Big0)
	assert.Equal(t, validator.Desc, "test2")

	// cancel validator && unstake && withdraw validator
	param6 := new(CancelValidatorParam)
	param6.ConsensusPubkey = validatorsKey[0].ConsensusPubkey
	input, err = param6.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, validatorsKey[0].Address, validatorsKey[0].Address, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(validatorsKey[0].Address, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	param7 := new(UnStakeParam)
	param7.ConsensusPubkey = validatorsKey[0].ConsensusPubkey
	param7.Amount = new(big.Int).Mul(big.NewInt(1000), params.ZNT1)
	input, err = param7.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	blockNumber = new(big.Int).SetUint64(1000000)
	input, err = utils.PackMethod(ABI, MethodWithdraw)
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	param8 := new(WithdrawValidatorParam)
	param8.ConsensusPubkey = validatorsKey[0].ConsensusPubkey
	input, err = param8.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, validatorsKey[0].Address, validatorsKey[0].Address, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(validatorsKey[0].Address, utils.NodeManagerContractAddress, input)
	assert.NotNil(t, err)
	allValidators, err = GetAllValidators(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, len(allValidators.AllValidators), loop-1)

	// check
	validator, _, err = GetValidator(contractQuery, validatorsKey[0].Dec)
	assert.Nil(t, err)
	assert.Equal(t, validator.Status, Remove)
	assert.Equal(t, sdb.GetBalance(stakeAddress), new(big.Int).Mul(big.NewInt(992000), params.ZNT1))

	// add block num
	blockNumber = new(big.Int).SetUint64(1300000)
	input, err = utils.PackMethod(ABI, MethodWithdraw)
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	param9 := new(WithdrawValidatorParam)
	param9.ConsensusPubkey = validatorsKey[0].ConsensusPubkey
	input, err = param9.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, validatorsKey[0].Address, validatorsKey[0].Address, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(validatorsKey[0].Address, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	validator, _, err = GetValidator(contractQuery, validatorsKey[0].Dec)
	assert.Nil(t, err)
	assert.Equal(t, validator.Status, Remove)
	assert.Equal(t, sdb.GetBalance(stakeAddress), new(big.Int).Mul(big.NewInt(993000), params.ZNT1))
	assert.Equal(t, sdb.GetBalance(validatorsKey[0].Address), new(big.Int).Mul(big.NewInt(1000000), params.ZNT1))

	// unstake
	param10 := new(UnStakeParam)
	param10.ConsensusPubkey = validatorsKey[0].ConsensusPubkey
	param10.Amount = new(big.Int).Mul(big.NewInt(7000), params.ZNT1)
	input, err = param10.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	validator, found, err := GetValidator(contractQuery, validatorsKey[0].Dec)
	assert.Nil(t, err)
	assert.Equal(t, found, false)
	assert.Equal(t, sdb.GetBalance(stakeAddress), new(big.Int).Mul(big.NewInt(1000000), params.ZNT1))

	// change epoch
	input, err = utils.PackMethod(ABI, MethodChangeEpoch)
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	epochInfo, err = getCurrentEpochInfo(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, epochInfo.ID, common.Big3)
	assert.Equal(t, epochInfo.StartHeight, new(big.Int).SetUint64(1300000))
	assert.Equal(t, len(epochInfo.Validators), 4)
	assert.Equal(t, len(epochInfo.Voters), 4)
	validator, _, err = GetValidator(contractQuery, validatorsKey[4].Dec)
	assert.Nil(t, err)
	assert.Equal(t, validator.Status, Lock)
}

func TestDistribute(t *testing.T) {

}

// generateTestPeer ONLY used for testing
func generateTestPeer() *Peer {
	pk, _ := crypto.GenerateKey()
	return &Peer{
		PubKey:  hexutil.Encode(crypto.CompressPubkey(&pk.PublicKey)),
		Address: crypto.PubkeyToAddress(pk.PublicKey),
	}
}

func generateTestPeers(n int) []*Peer {
	peers := make([]*Peer, n)
	for i := 0; i < n; i++ {
		peers[i] = generateTestPeer()
	}
	return peers
}
