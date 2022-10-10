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
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/assert"
)

var (
	sdb              *state.StateDB
	testGenesisNum   = 4
	acct             *ecdsa.PublicKey
	testGenesisPeers []common.Address
)

func Init() {
	key, _ := crypto.GenerateKey()
	acct = &key.PublicKey

	InitNodeManager()
	sdb = native.NewTestStateDB()
	testGenesisPeers, _ = native.GenerateTestPeers(testGenesisNum)
	StoreCommunityInfo(sdb, big.NewInt(2000), common.EmptyAddress)
	StoreGenesisEpoch(sdb, testGenesisPeers, testGenesisPeers)
	StoreGenesisGlobalConfig(sdb)
}

func TestCheckGenesis(t *testing.T) {
	Init()
	// check get spec methodID
	m := GetSpecMethodID()
	assert.Equal(t, m["fe6f86f8"], true)
	assert.Equal(t, m["083c6323"], true)

	blockNumber := big.NewInt(1)
	extra := uint64(21000000000000)
	contractRef := native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
	contract := native.NewNativeContract(sdb, contractRef)

	globalConfig, err := GetGlobalConfigImpl(contract)
	assert.Nil(t, err)
	assert.Equal(t, globalConfig.BlockPerEpoch, GenesisBlockPerEpoch)
	assert.Equal(t, globalConfig.MaxCommissionChange, GenesisMaxCommissionChange)
	assert.Equal(t, globalConfig.MinInitialStake, GenesisMinInitialStake)
	assert.Equal(t, globalConfig.VoterValidatorNum, GenesisVoterValidatorNum)
	assert.Equal(t, globalConfig.ConsensusValidatorNum, GenesisConsensusValidatorNum)

	communityInfo, err := GetCommunityInfoImpl(contract)
	assert.Nil(t, err)
	assert.Equal(t, communityInfo.CommunityRate, big.NewInt(2000))
	assert.Equal(t, communityInfo.CommunityAddress, common.EmptyAddress)

	epochInfo, err := GetCurrentEpochInfoImpl(contract)
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
	assert.Equal(t, globalConfig2.BlockPerEpoch, GenesisBlockPerEpoch)
	assert.Equal(t, globalConfig2.MaxCommissionChange, GenesisMaxCommissionChange)
	assert.Equal(t, globalConfig2.MinInitialStake, GenesisMinInitialStake)
	assert.Equal(t, globalConfig2.VoterValidatorNum, GenesisVoterValidatorNum)
	assert.Equal(t, globalConfig2.ConsensusValidatorNum, GenesisConsensusValidatorNum)
	assert.Equal(t, globalConfig2.MinProposalStake, GenesisMinProposalStake)

	param2 := new(GetCommunityInfoParam)
	input, err = param2.Encode()
	assert.Nil(t, err)
	ret, _, err = contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	communityInfo2 := new(CommunityInfo)
	err = communityInfo2.Decode(ret)
	assert.Nil(t, err)
	assert.Equal(t, communityInfo2.CommunityRate, big.NewInt(2000))
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

	currentEpochInfo, err = GetCurrentEpochInfoFromDB(sdb)
	assert.Nil(t, err)
	assert.Equal(t, currentEpochInfo.ID, big.NewInt(1))

	globalConfig, err = GetGlobalConfigFromDB(sdb)
	assert.Nil(t, err)
	assert.Equal(t, globalConfig.BlockPerEpoch, GenesisBlockPerEpoch)
}

func TestStake(t *testing.T) {
	Init()
	blockNumber := big.NewInt(399999)
	extra := uint64(21000000000000)
	contractRefQuery := native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
	contractQuery := native.NewNativeContract(sdb, contractRefQuery)

	type ValidatorKey struct {
		ConsensusAddr common.Address
		StakeAddress  common.Address
	}
	// create validator
	loop := 6
	validatorsKey := make([]*ValidatorKey, 0, loop)
	for i := 0; i < loop; i++ {
		pk, _ := crypto.GenerateKey()
		consensusAddr := crypto.PubkeyToAddress(pk.PublicKey)
		caller := crypto.PubkeyToAddress(*acct)
		sdb.SetBalance(caller, new(big.Int).Mul(big.NewInt(1000000), params.ZNT1))
		param := new(CreateValidatorParam)
		param.ConsensusAddress = consensusAddr
		param.SignerAddress = consensusAddr
		param.ProposalAddress = caller
		param.InitStake = new(big.Int).Mul(big.NewInt(100000), params.ZNT1)
		param.Commission = new(big.Int).SetUint64(2000)
		param.Desc = "test"
		validatorsKey = append(validatorsKey, &ValidatorKey{param.ConsensusAddress, caller})
		input, err := param.Encode()
		assert.Nil(t, err)
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
		_, _, err = contractRef.NativeCall(caller, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
	}
	// check
	allValidators, err := getAllValidators(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, len(allValidators.AllValidators), loop)
	validator, _, err := getValidator(contractQuery, validatorsKey[0].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, validator.TotalStake.BigInt(), new(big.Int).Mul(big.NewInt(100000), params.ZNT1))
	assert.Equal(t, validator.SelfStake.BigInt(), new(big.Int).Mul(big.NewInt(100000), params.ZNT1))
	assert.Equal(t, validator.Status, Unlock)
	assert.Equal(t, validator.Commission.Rate.BigInt(), new(big.Int).SetUint64(2000))
	assert.Equal(t, validator.UnlockHeight, new(big.Int))

	//stake
	pkStake, _ := crypto.GenerateKey()
	staker := &pkStake.PublicKey
	stakeAddress := crypto.PubkeyToAddress(*staker)
	sdb.SetBalance(stakeAddress, new(big.Int).Mul(big.NewInt(1000000), params.ZNT1))
	param1 := new(StakeParam)
	param1.ConsensusAddress = validatorsKey[0].ConsensusAddr
	param1.Amount = new(big.Int).Mul(big.NewInt(10000), params.ZNT1)
	input, err := param1.Encode()
	assert.Nil(t, err)
	contractRef := native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	// check
	validator, _, err = getValidator(contractQuery, validatorsKey[0].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, validator.TotalStake.BigInt(), new(big.Int).Mul(big.NewInt(110000), params.ZNT1))
	assert.Equal(t, validator.SelfStake.BigInt(), new(big.Int).Mul(big.NewInt(100000), params.ZNT1))
	assert.Equal(t, validator.Status, Unlock)
	assert.Equal(t, validator.Commission.Rate.BigInt(), new(big.Int).SetUint64(2000))
	assert.Equal(t, validator.UnlockHeight, new(big.Int))
	totalPool, err := getTotalPool(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, totalPool.TotalPool.BigInt(), new(big.Int).Mul(big.NewInt(610000), params.ZNT1))

	// unstake
	param2 := new(UnStakeParam)
	param2.ConsensusAddress = validatorsKey[0].ConsensusAddr
	param2.Amount = new(big.Int).Mul(big.NewInt(1000), params.ZNT1)
	input, err = param2.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	validator, _, err = getValidator(contractQuery, validatorsKey[0].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, validator.TotalStake.BigInt(), new(big.Int).Mul(big.NewInt(109000), params.ZNT1))
	assert.Equal(t, validator.SelfStake.BigInt(), new(big.Int).Mul(big.NewInt(100000), params.ZNT1))
	assert.Equal(t, validator.Status, Unlock)
	assert.Equal(t, validator.UnlockHeight, new(big.Int))
	totalPool, err = getTotalPool(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, totalPool.TotalPool.BigInt(), new(big.Int).Mul(big.NewInt(609000), params.ZNT1))
	assert.Equal(t, sdb.GetBalance(stakeAddress), new(big.Int).Mul(big.NewInt(991000), params.ZNT1))

	// change epoch
	input, err = utils.PackMethod(ABI, MethodChangeEpoch)
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, utils.SystemTxSender, utils.SystemTxSender, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(utils.SystemTxSender, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	epochInfo, err := GetCurrentEpochInfoImpl(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, epochInfo.ID, common.Big2)
	assert.Equal(t, epochInfo.StartHeight, new(big.Int).SetUint64(400000))
	assert.Equal(t, len(epochInfo.Validators), 4)
	assert.Equal(t, len(epochInfo.Voters), 4)
	validator, _, err = getValidator(contractQuery, validatorsKey[0].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, validator.Status, Lock)
	totalPool, err = getTotalPool(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, totalPool.TotalPool.BigInt(), new(big.Int).Mul(big.NewInt(609000), params.ZNT1))

	// unstake
	param3 := new(UnStakeParam)
	param3.ConsensusAddress = validatorsKey[0].ConsensusAddr
	param3.Amount = new(big.Int).Mul(big.NewInt(1000), params.ZNT1)
	input, err = param3.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	validator, _, err = getValidator(contractQuery, validatorsKey[0].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, validator.TotalStake.BigInt(), new(big.Int).Mul(big.NewInt(108000), params.ZNT1))
	assert.Equal(t, validator.SelfStake.BigInt(), new(big.Int).Mul(big.NewInt(100000), params.ZNT1))
	assert.Equal(t, validator.Status, Lock)
	assert.Equal(t, validator.UnlockHeight, new(big.Int))

	// withdraw
	input, err = utils.PackMethod(ABI, MethodWithdraw)
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.NotNil(t, err)
	totalPool, err = getTotalPool(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, totalPool.TotalPool.BigInt(), new(big.Int).Mul(big.NewInt(609000), params.ZNT1))
	blockNumber = big.NewInt(800000)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	totalPool, err = getTotalPool(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, totalPool.TotalPool.BigInt(), new(big.Int).Mul(big.NewInt(608000), params.ZNT1))

	// update validator
	param4 := new(UpdateValidatorParam)
	param4.ConsensusAddress = validatorsKey[0].ConsensusAddr
	param4.Desc = "test2"
	input, err = param4.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, validatorsKey[0].StakeAddress, validatorsKey[0].StakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(validatorsKey[0].StakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	param5 := new(UpdateCommissionParam)
	param5.ConsensusAddress = validatorsKey[0].ConsensusAddr
	param5.Commission = new(big.Int).SetUint64(2500)
	input, err = param5.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, validatorsKey[0].StakeAddress, validatorsKey[0].StakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(validatorsKey[0].StakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	validator, _, err = getValidator(contractQuery, validatorsKey[0].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, validator.TotalStake.BigInt(), new(big.Int).Mul(big.NewInt(108000), params.ZNT1))
	assert.Equal(t, validator.SelfStake.BigInt(), new(big.Int).Mul(big.NewInt(100000), params.ZNT1))
	assert.Equal(t, validator.ProposalAddress, validatorsKey[0].StakeAddress)
	assert.Equal(t, validator.Status, Lock)
	assert.Equal(t, validator.Commission.Rate.BigInt(), new(big.Int).SetUint64(2500))
	assert.Equal(t, validator.Commission.UpdateHeight, new(big.Int).SetUint64(800000))
	assert.Equal(t, validator.UnlockHeight, new(big.Int))
	assert.Equal(t, validator.Desc, "test2")

	// cancel validator && unstake && withdraw validator
	param6 := new(CancelValidatorParam)
	param6.ConsensusAddress = validatorsKey[0].ConsensusAddr
	input, err = param6.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, validatorsKey[0].StakeAddress, validatorsKey[0].StakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(validatorsKey[0].StakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	param7 := new(UnStakeParam)
	param7.ConsensusAddress = validatorsKey[0].ConsensusAddr
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
	assert.NotNil(t, err)
	param8 := new(WithdrawValidatorParam)
	param8.ConsensusAddress = validatorsKey[0].ConsensusAddr
	input, err = param8.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, validatorsKey[0].StakeAddress, validatorsKey[0].StakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(validatorsKey[0].StakeAddress, utils.NodeManagerContractAddress, input)
	assert.NotNil(t, err)
	allValidators, err = getAllValidators(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, len(allValidators.AllValidators), loop-1)

	// check
	validator, _, err = getValidator(contractQuery, validatorsKey[0].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, validator.Status, Remove)
	assert.Equal(t, sdb.GetBalance(stakeAddress), new(big.Int).Mul(big.NewInt(992000), params.ZNT1))

	blockNumber = new(big.Int).SetUint64(799999)
	// change epoch
	input, err = utils.PackMethod(ABI, MethodChangeEpoch)
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, utils.SystemTxSender, utils.SystemTxSender, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(utils.SystemTxSender, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	blockNumber = new(big.Int).SetUint64(1199999)
	// change epoch
	input, err = utils.PackMethod(ABI, MethodChangeEpoch)
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, utils.SystemTxSender, utils.SystemTxSender, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(utils.SystemTxSender, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// add block num
	blockNumber = new(big.Int).SetUint64(1599999)
	input, err = utils.PackMethod(ABI, MethodWithdraw)
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	param9 := new(WithdrawValidatorParam)
	param9.ConsensusAddress = validatorsKey[0].ConsensusAddr
	input, err = param9.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, validatorsKey[0].StakeAddress, validatorsKey[0].StakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(validatorsKey[0].StakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	validator, _, err = getValidator(contractQuery, validatorsKey[0].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, validator.Status, Remove)
	assert.Equal(t, sdb.GetBalance(stakeAddress), new(big.Int).Mul(big.NewInt(993000), params.ZNT1))
	assert.Equal(t, sdb.GetBalance(validatorsKey[0].StakeAddress), new(big.Int).Mul(big.NewInt(1000000), params.ZNT1))

	// unstake
	param10 := new(UnStakeParam)
	param10.ConsensusAddress = validatorsKey[0].ConsensusAddr
	param10.Amount = new(big.Int).Mul(big.NewInt(7000), params.ZNT1)
	input, err = param10.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	validator, found, err := getValidator(contractQuery, validatorsKey[0].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, found, false)
	assert.Equal(t, sdb.GetBalance(stakeAddress), new(big.Int).Mul(big.NewInt(1000000), params.ZNT1))

	// change epoch
	input, err = utils.PackMethod(ABI, MethodChangeEpoch)
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, utils.SystemTxSender, utils.SystemTxSender, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(utils.SystemTxSender, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	epochInfo, err = GetCurrentEpochInfoImpl(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, epochInfo.ID, new(big.Int).SetUint64(5))
	assert.Equal(t, epochInfo.StartHeight, new(big.Int).SetUint64(1600000))
	assert.Equal(t, len(epochInfo.Validators), 4)
	assert.Equal(t, len(epochInfo.Voters), 4)
	validator, _, err = getValidator(contractQuery, validatorsKey[4].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, validator.Status, Lock)
}

func TestChangeEpoch(t *testing.T) {
	Init()
	blockNumber := big.NewInt(0)
	extra := uint64(21000000000000)
	contractRefQuery := native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
	contractQuery := native.NewNativeContract(sdb, contractRefQuery)

	type ValidatorKey struct {
		ConsensusAddr common.Address
		StakeAddress  common.Address
	}
	// create validator
	loop := 6
	validatorsKey := make([]*ValidatorKey, 0, loop)
	for i := 0; i < loop; i++ {
		pk, _ := crypto.GenerateKey()
		consensusAddr := crypto.PubkeyToAddress(pk.PublicKey)
		caller := crypto.PubkeyToAddress(*acct)
		sdb.SetBalance(caller, new(big.Int).Mul(big.NewInt(1000000), params.ZNT1))
		param := new(CreateValidatorParam)
		param.ConsensusAddress = consensusAddr
		param.SignerAddress = consensusAddr
		param.ProposalAddress = caller
		param.InitStake = new(big.Int).Mul(big.NewInt(100000), params.ZNT1)
		param.Commission = new(big.Int).SetUint64(2000)
		param.Desc = "test"
		validatorsKey = append(validatorsKey, &ValidatorKey{param.ConsensusAddress, caller})
		input, err := param.Encode()
		assert.Nil(t, err)
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
		_, _, err = contractRef.NativeCall(caller, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
	}

	blockNumber = new(big.Int).SetUint64(399999)
	// change epoch
	input, err := utils.PackMethod(ABI, MethodChangeEpoch)
	assert.Nil(t, err)
	contractRef := native.NewContractRef(sdb, utils.SystemTxSender, utils.SystemTxSender, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(utils.SystemTxSender, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	epochInfo, err := GetCurrentEpochInfoImpl(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, epochInfo.ID, common.Big2)
	fmt.Println(epochInfo.Validators)

	//stake
	pkStake, _ := crypto.GenerateKey()
	staker := &pkStake.PublicKey
	stakeAddress := crypto.PubkeyToAddress(*staker)
	sdb.SetBalance(stakeAddress, new(big.Int).Mul(big.NewInt(1000000), params.ZNT1))
	param1 := new(StakeParam)
	param1.ConsensusAddress = validatorsKey[4].ConsensusAddr
	param1.Amount = new(big.Int).Mul(big.NewInt(10000), params.ZNT1)
	input, err = param1.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	blockNumber = new(big.Int).SetUint64(799999)
	// change epoch
	input, err = utils.PackMethod(ABI, MethodChangeEpoch)
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, utils.SystemTxSender, utils.SystemTxSender, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(utils.SystemTxSender, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	epochInfo, err = GetCurrentEpochInfoImpl(contractQuery)
	assert.Nil(t, err)
	assert.Equal(t, epochInfo.ID, common.Big3)
	fmt.Println(epochInfo.Validators)
}

func TestDistribute(t *testing.T) {
	Init()
	blockNumber := big.NewInt(399999)
	extra := uint64(21000000000000)
	contractRefQuery := native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
	contractQuery := native.NewNativeContract(sdb, contractRefQuery)

	type ValidatorKey struct {
		ConsensusAddr common.Address
		StakeAddress  common.Address
	}
	// create validator
	// 6 address with 1000000 token and  create 6 validators with 100000 init stake
	loop := 6
	validatorsKey := make([]*ValidatorKey, 0, loop)
	for i := 0; i < loop; i++ {
		pk, _ := crypto.GenerateKey()
		consensusAddr := crypto.PubkeyToAddress(pk.PublicKey)
		caller := crypto.PubkeyToAddress(*acct)
		sdb.SetBalance(caller, new(big.Int).Mul(big.NewInt(1000000), params.ZNT1))
		param := new(CreateValidatorParam)
		param.ConsensusAddress = consensusAddr
		param.SignerAddress = consensusAddr
		param.ProposalAddress = caller
		param.InitStake = new(big.Int).Mul(big.NewInt(100000), params.ZNT1)
		param.Commission = new(big.Int).SetUint64(2000)
		param.Desc = "test"
		validatorsKey = append(validatorsKey, &ValidatorKey{param.ConsensusAddress, caller})
		input, err := param.Encode()
		assert.Nil(t, err)
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
		_, _, err = contractRef.NativeCall(caller, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
	}

	// stake
	// 2 address with 1000000 token and stake 10000, 20000 to first validator
	pkStake, _ := crypto.GenerateKey()
	staker := &pkStake.PublicKey
	stakeAddress := crypto.PubkeyToAddress(*staker)
	sdb.SetBalance(stakeAddress, new(big.Int).Mul(big.NewInt(1000000), params.ZNT1))
	param1 := new(StakeParam)
	param1.ConsensusAddress = validatorsKey[0].ConsensusAddr
	param1.Amount = new(big.Int).Mul(big.NewInt(10000), params.ZNT1)
	input, err := param1.Encode()
	assert.Nil(t, err)
	contractRef := native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	pkStake2, _ := crypto.GenerateKey()
	staker2 := &pkStake2.PublicKey
	stakeAddress2 := crypto.PubkeyToAddress(*staker2)
	sdb.SetBalance(stakeAddress2, new(big.Int).Mul(big.NewInt(1000000), params.ZNT1))
	param2 := new(StakeParam)
	param2.ConsensusAddress = validatorsKey[0].ConsensusAddr
	param2.Amount = new(big.Int).Mul(big.NewInt(20000), params.ZNT1)
	input, err = param2.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress2, stakeAddress2, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress2, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// change epoch
	input, err = utils.PackMethod(ABI, MethodChangeEpoch)
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, utils.SystemTxSender, utils.SystemTxSender, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(utils.SystemTxSender, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// here we have 4 validators with 100000 self stake, and validator 1 have 10000, 20000 user stake, and commission is 20%
	// first add 1000 balance of node_manager contract to distribute
	sdb.AddBalance(utils.NodeManagerContractAddress, new(big.Int).Mul(big.NewInt(1000), params.ZNT1))
	// call endblock
	param3 := new(EndBlockParam)
	input, err = param3.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, utils.SystemTxSender, utils.SystemTxSender, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(utils.SystemTxSender, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	accumulatedCommission, err := getAccumulatedCommission(contractQuery, validatorsKey[0].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, accumulatedCommission.Amount.BigInt(), new(big.Int).Mul(big.NewInt(50), params.ZNT1))
	validatorAccumulatedRewards, err := getValidatorAccumulatedRewards(contractQuery, validatorsKey[0].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, validatorAccumulatedRewards.Rewards.BigInt(), new(big.Int).Mul(big.NewInt(200), params.ZNT1))

	accumulatedCommission2, err := getAccumulatedCommission(contractQuery, validatorsKey[1].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, accumulatedCommission2.Amount.BigInt(), new(big.Int).Mul(big.NewInt(50), params.ZNT1))
	validatorAccumulatedRewards2, err := getValidatorAccumulatedRewards(contractQuery, validatorsKey[1].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, validatorAccumulatedRewards2.Rewards.BigInt(), new(big.Int).Mul(big.NewInt(200), params.ZNT1))

	// test query method
	{
		p1 := &GetEpochInfoParam{
			ID: common.Big1,
		}
		input, err = p1.Encode()
		assert.Nil(t, err)
		contractRef = native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
		ret, _, err := contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
		epochInfo := new(EpochInfo)
		err = epochInfo.Decode(ret)
		assert.Nil(t, err)
		assert.Equal(t, epochInfo.ID, common.Big1)

		p2 := &GetAllValidatorsParam{}
		input, err = p2.Encode()
		assert.Nil(t, err)
		contractRef = native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
		ret, _, err = contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
		allValidators := new(AllValidators)
		err = allValidators.Decode(ret)
		assert.Nil(t, err)
		assert.Equal(t, len(allValidators.AllValidators), loop)

		p3 := &GetValidatorParam{
			ConsensusAddress: validatorsKey[0].ConsensusAddr,
		}
		input, err = p3.Encode()
		assert.Nil(t, err)
		contractRef = native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
		ret, _, err = contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
		validator := new(Validator)
		err = validator.Decode(ret)
		assert.Nil(t, err)
		assert.Equal(t, validator.ConsensusAddress, validatorsKey[0].ConsensusAddr)

		p4 := &GetStakeInfoParam{
			ConsensusAddress: validatorsKey[0].ConsensusAddr,
			StakeAddress:     stakeAddress,
		}
		input, err = p4.Encode()
		assert.Nil(t, err)
		contractRef = native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
		ret, _, err = contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
		stakeInfo := new(StakeInfo)
		err = stakeInfo.Decode(ret)
		assert.Nil(t, err)
		assert.Equal(t, stakeInfo.ConsensusAddr, validatorsKey[0].ConsensusAddr)

		p5 := &GetUnlockingInfoParam{
			StakeAddress: stakeAddress,
		}
		input, err = p5.Encode()
		assert.Nil(t, err)
		contractRef = native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
		ret, _, err = contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
		unlockingInfo := new(UnlockingInfo)
		err = unlockingInfo.Decode(ret)
		assert.Nil(t, err)
		assert.Equal(t, unlockingInfo.StakeAddress, stakeAddress)

		p6 := &GetStakeStartingInfoParam{
			ConsensusAddress: validatorsKey[0].ConsensusAddr,
			StakeAddress:     stakeAddress,
		}
		input, err = p6.Encode()
		assert.Nil(t, err)
		contractRef = native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
		ret, _, err = contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
		stakeStartingInfo := new(StakeStartingInfo)
		err = stakeStartingInfo.Decode(ret)
		assert.Nil(t, err)
		assert.Equal(t, stakeStartingInfo.Stake.BigInt(), new(big.Int).Mul(big.NewInt(10000), params.ZNT1))

		p7 := &GetAccumulatedCommissionParam{
			ConsensusAddress: validatorsKey[0].ConsensusAddr,
		}
		input, err = p7.Encode()
		assert.Nil(t, err)
		contractRef = native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
		ret, _, err = contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
		accumulatedCommission = new(AccumulatedCommission)
		err = accumulatedCommission.Decode(ret)
		assert.Nil(t, err)
		assert.Equal(t, accumulatedCommission.Amount.BigInt(), new(big.Int).Mul(big.NewInt(50), params.ZNT1))

		p8 := &GetValidatorSnapshotRewardsParam{
			ConsensusAddress: validatorsKey[0].ConsensusAddr,
			Period:           3,
		}
		input, err = p8.Encode()
		assert.Nil(t, err)
		contractRef = native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
		ret, _, err = contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
		validatorSnapshotRewards := new(ValidatorSnapshotRewards)
		err = validatorSnapshotRewards.Decode(ret)
		assert.Nil(t, err)
		assert.Equal(t, validatorSnapshotRewards.ReferenceCount, uint64(2))

		p9 := &GetValidatorAccumulatedRewardsParam{
			ConsensusAddress: validatorsKey[0].ConsensusAddr,
		}
		input, err = p9.Encode()
		assert.Nil(t, err)
		contractRef = native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
		ret, _, err = contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
		validatorAccumulatedRewards = new(ValidatorAccumulatedRewards)
		err = validatorAccumulatedRewards.Decode(ret)
		assert.Nil(t, err)
		assert.Equal(t, validatorAccumulatedRewards.Period, uint64(4))

		p10 := &GetValidatorOutstandingRewardsParam{
			ConsensusAddress: validatorsKey[0].ConsensusAddr,
		}
		input, err = p10.Encode()
		assert.Nil(t, err)
		contractRef = native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
		ret, _, err = contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
		validatorOutstandingRewards := new(ValidatorOutstandingRewards)
		err = validatorOutstandingRewards.Decode(ret)
		assert.Nil(t, err)
		assert.Equal(t, validatorOutstandingRewards.Rewards.BigInt(), new(big.Int).Mul(big.NewInt(250), params.ZNT1))

		p11 := &GetTotalPoolParam{}
		input, err = p11.Encode()
		assert.Nil(t, err)
		contractRef = native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
		ret, _, err = contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
		totalPool := new(TotalPool)
		err = totalPool.Decode(ret)
		assert.Nil(t, err)
		assert.Equal(t, totalPool.TotalPool.BigInt(), new(big.Int).Mul(big.NewInt(630000), params.ZNT1))

		p12 := &GetOutstandingRewardsParam{}
		input, err = p12.Encode()
		assert.Nil(t, err)
		contractRef = native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
		ret, _, err = contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
		assert.Nil(t, err)
		outstandingRewards := new(OutstandingRewards)
		err = outstandingRewards.Decode(ret)
		assert.Nil(t, err)
		assert.Equal(t, outstandingRewards.Rewards.BigInt(), new(big.Int).Mul(big.NewInt(1000), params.ZNT1))
	}

	// check getStakeRewards
	p13 := &GetStakeRewardsParam{
		ConsensusAddress: validatorsKey[0].ConsensusAddr,
		StakeAddress:     stakeAddress,
	}
	input, err = p13.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, common.EmptyAddress, common.EmptyAddress, blockNumber, common.Hash{}, extra, nil)
	ret, _, err := contractRef.NativeCall(common.EmptyAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	stakeRewards := new(StakeRewards)
	err = stakeRewards.Decode(ret)
	assert.Nil(t, err)
	s, _ := new(big.Int).SetString("15384615384615380000", 10)
	assert.Equal(t, stakeRewards.Rewards.BigInt(), s)

	// withdraw stake rewards and commission
	param4 := new(WithdrawStakeRewardsParam)
	param4.ConsensusAddress = validatorsKey[0].ConsensusAddr
	input, err = param4.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	param5 := new(WithdrawStakeRewardsParam)
	param5.ConsensusAddress = validatorsKey[0].ConsensusAddr
	input, err = param5.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress2, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress2, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	param6 := new(WithdrawCommissionParam)
	param6.ConsensusAddress = validatorsKey[0].ConsensusAddr
	input, err = param6.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, validatorsKey[0].StakeAddress, validatorsKey[0].StakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(validatorsKey[0].StakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check balance
	b1, _ := new(big.Int).SetString("990015384615384615380000", 10)
	b2, _ := new(big.Int).SetString("980030769230769230760000", 10)
	assert.Equal(t, sdb.GetBalance(stakeAddress), b1)
	assert.Equal(t, sdb.GetBalance(stakeAddress2), b2)
	assert.Equal(t, sdb.GetBalance(validatorsKey[0].StakeAddress), new(big.Int).Mul(big.NewInt(900050), params.ZNT1))

	// check states
	validatorAccumulatedRewards, err = getValidatorAccumulatedRewards(contractQuery, validatorsKey[0].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, validatorAccumulatedRewards.Rewards.BigInt(), new(big.Int))
	assert.Equal(t, validatorAccumulatedRewards.Period, uint64(6))
	_, err = getValidatorSnapshotRewards(contractQuery, validatorsKey[0].ConsensusAddr, 0)
	assert.NotNil(t, err)
	_, err = getValidatorSnapshotRewards(contractQuery, validatorsKey[0].ConsensusAddr, 1)
	assert.Nil(t, err)
	_, err = getValidatorSnapshotRewards(contractQuery, validatorsKey[0].ConsensusAddr, 2)
	assert.NotNil(t, err)
	_, err = getValidatorSnapshotRewards(contractQuery, validatorsKey[0].ConsensusAddr, 3)
	assert.NotNil(t, err)
	_, err = getValidatorSnapshotRewards(contractQuery, validatorsKey[0].ConsensusAddr, 4)
	assert.Nil(t, err)
	_, err = getValidatorSnapshotRewards(contractQuery, validatorsKey[0].ConsensusAddr, 5)
	assert.Nil(t, err)

	// withdraw validator stake rewards
	param7 := new(WithdrawStakeRewardsParam)
	param7.ConsensusAddress = validatorsKey[0].ConsensusAddr
	input, err = param7.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, validatorsKey[0].StakeAddress, validatorsKey[0].StakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(validatorsKey[0].StakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	param8 := new(WithdrawStakeRewardsParam)
	param8.ConsensusAddress = validatorsKey[0].ConsensusAddr
	input, err = param8.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, validatorsKey[0].StakeAddress, validatorsKey[0].StakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(validatorsKey[0].StakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	b3, _ := new(big.Int).SetString("900203846153846153800000", 10)
	assert.Equal(t, sdb.GetBalance(validatorsKey[0].StakeAddress), b3)
	validatorOutstandingRewards, err := getValidatorOutstandingRewards(contractQuery, validatorsKey[0].ConsensusAddr)
	assert.Nil(t, err)
	b4, _ := new(big.Int).SetString("60000", 10)
	assert.Equal(t, validatorOutstandingRewards.Rewards.BigInt(), b4)
	validatorAccumulatedRewards, err = getValidatorAccumulatedRewards(contractQuery, validatorsKey[0].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, validatorAccumulatedRewards.Rewards.BigInt(), new(big.Int))
	assert.Equal(t, validatorAccumulatedRewards.Period, uint64(8))
	_, err = getValidatorSnapshotRewards(contractQuery, validatorsKey[0].ConsensusAddr, 0)
	assert.NotNil(t, err)
	_, err = getValidatorSnapshotRewards(contractQuery, validatorsKey[0].ConsensusAddr, 1)
	assert.NotNil(t, err)
	_, err = getValidatorSnapshotRewards(contractQuery, validatorsKey[0].ConsensusAddr, 2)
	assert.NotNil(t, err)
	_, err = getValidatorSnapshotRewards(contractQuery, validatorsKey[0].ConsensusAddr, 3)
	assert.NotNil(t, err)
	_, err = getValidatorSnapshotRewards(contractQuery, validatorsKey[0].ConsensusAddr, 4)
	assert.Nil(t, err)
	_, err = getValidatorSnapshotRewards(contractQuery, validatorsKey[0].ConsensusAddr, 5)
	assert.Nil(t, err)
	_, err = getValidatorSnapshotRewards(contractQuery, validatorsKey[0].ConsensusAddr, 6)
	assert.NotNil(t, err)
	_, err = getValidatorSnapshotRewards(contractQuery, validatorsKey[0].ConsensusAddr, 7)
	assert.Nil(t, err)

	// here we have 4 validators with 100000 self stake, and validator 1 have 10000, 20000 user stake, and commission is 20%
	// add 2000 balance of node_manager contract to distribute
	sdb.AddBalance(utils.NodeManagerContractAddress, new(big.Int).Mul(big.NewInt(1000), params.ZNT1))
	// call endblock
	param9 := new(EndBlockParam)
	input, err = param9.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, utils.SystemTxSender, utils.SystemTxSender, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(utils.SystemTxSender, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	sdb.AddBalance(utils.NodeManagerContractAddress, new(big.Int).Mul(big.NewInt(1000), params.ZNT1))
	// call endblock
	param10 := new(EndBlockParam)
	input, err = param10.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, utils.SystemTxSender, utils.SystemTxSender, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(utils.SystemTxSender, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// cancel validator
	param11 := new(CancelValidatorParam)
	param11.ConsensusAddress = validatorsKey[0].ConsensusAddr
	input, err = param11.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, validatorsKey[0].StakeAddress, validatorsKey[0].StakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(validatorsKey[0].StakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	blockNumber = big.NewInt(799999)
	// change epoch
	input, err = utils.PackMethod(ABI, MethodChangeEpoch)
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, utils.SystemTxSender, utils.SystemTxSender, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(utils.SystemTxSender, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	blockNumber = big.NewInt(900000)
	// withdraw validator
	param12 := new(WithdrawValidatorParam)
	param12.ConsensusAddress = validatorsKey[0].ConsensusAddr
	input, err = param12.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, validatorsKey[0].StakeAddress, validatorsKey[0].StakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(validatorsKey[0].StakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	b5, _ := new(big.Int).SetString("1000611538461538461400000", 10) // include commission and stake rewards
	assert.Equal(t, sdb.GetBalance(validatorsKey[0].StakeAddress), b5)

	// unstake
	param13 := new(UnStakeParam)
	param13.ConsensusAddress = validatorsKey[0].ConsensusAddr
	param13.Amount = new(big.Int).Mul(big.NewInt(10000), params.ZNT1)
	input, err = param13.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress, stakeAddress, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)
	param14 := new(UnStakeParam)
	param14.ConsensusAddress = validatorsKey[0].ConsensusAddr
	param14.Amount = new(big.Int).Mul(big.NewInt(20000), params.ZNT1)
	input, err = param14.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, stakeAddress2, stakeAddress2, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(stakeAddress2, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// check
	b6, _ := new(big.Int).SetString("1000046153846153846140000", 10)
	b7, _ := new(big.Int).SetString("1000092307692307692280000", 10)
	assert.Equal(t, sdb.GetBalance(stakeAddress), b6)
	assert.Equal(t, sdb.GetBalance(stakeAddress2), b7)
	assert.Equal(t, sdb.GetBalance(common.EmptyAddress), new(big.Int).SetUint64(180000))
	_, found, err := getValidator(contractQuery, validatorsKey[0].ConsensusAddr)
	assert.Nil(t, err)
	assert.Equal(t, found, false)
}

func TestPerformance(t *testing.T) {
	Init()
	blockNumber := 399999
	extra := uint64(21000000000000)

	type ValidatorKey struct {
		ConsensusAddr common.Address
		StakeAddress  common.Address
	}
	// create validator
	loop := 300
	validatorsKey := make([]*ValidatorKey, 0, loop)
	for i := 0; i < loop; i++ {
		pk, _ := crypto.GenerateKey()
		consensusAddr := crypto.PubkeyToAddress(pk.PublicKey)
		caller := crypto.PubkeyToAddress(*acct)
		sdb.SetBalance(caller, new(big.Int).Mul(big.NewInt(1000000), params.ZNT1))
		param := new(CreateValidatorParam)
		param.ConsensusAddress = consensusAddr
		param.SignerAddress = consensusAddr
		param.ProposalAddress = caller
		param.InitStake = new(big.Int).Mul(big.NewInt(1000000), params.ZNT1)
		param.Commission = new(big.Int).SetUint64(2000)
		param.Desc = "test"
		validatorsKey = append(validatorsKey, &ValidatorKey{param.ConsensusAddress, caller})
		input, err := param.Encode()
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "CreateValidator", input, caller, caller, blockNumber, extra, sdb)
		assert.Nil(t, err)
	}

	loop = 10000
	stakeAddressList := make([]common.Address, 0, loop)
	for i := 0; i < loop; i++ {
		//stake
		pkStake, _ := crypto.GenerateKey()
		staker := &pkStake.PublicKey
		stakeAddress := crypto.PubkeyToAddress(*staker)
		stakeAddressList = append(stakeAddressList, stakeAddress)
		sdb.SetBalance(stakeAddress, new(big.Int).Mul(big.NewInt(1000000), params.ZNT1))
		param1 := new(StakeParam)
		param1.ConsensusAddress = validatorsKey[0].ConsensusAddr
		param1.Amount = new(big.Int).Mul(big.NewInt(200), params.ZNT1)
		input, err := param1.Encode()
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "Stake", input, stakeAddress, stakeAddress, blockNumber, extra, sdb)
		assert.Nil(t, err)
	}

	// change epoch
	input, err := utils.PackMethod(ABI, MethodChangeEpoch)
	assert.Nil(t, err)
	_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "ChangeEpoch", input, utils.SystemTxSender, utils.SystemTxSender, blockNumber, extra, sdb)
	assert.Nil(t, err)

	// call endblock
	sdb.AddBalance(utils.NodeManagerContractAddress, new(big.Int).Mul(big.NewInt(10000000), params.ZNT1))
	param := new(EndBlockParam)
	input, err = param.Encode()
	assert.Nil(t, err)
	_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "EndBlock", input, utils.SystemTxSender, utils.SystemTxSender, blockNumber, extra, sdb)
	assert.Nil(t, err)

	// withdraw stake rewards and commission
	param1 := new(WithdrawStakeRewardsParam)
	param1.ConsensusAddress = validatorsKey[0].ConsensusAddr
	input, err = param1.Encode()
	assert.Nil(t, err)
	_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "WithdrawStakeRewards", input, stakeAddressList[0], stakeAddressList[0], blockNumber, extra, sdb)
	assert.Nil(t, err)
	param2 := new(WithdrawCommissionParam)
	param2.ConsensusAddress = validatorsKey[0].ConsensusAddr
	input, err = param2.Encode()
	assert.Nil(t, err)
	_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "WithdrawCommission", input, validatorsKey[0].StakeAddress, validatorsKey[0].StakeAddress, blockNumber, extra, sdb)
	assert.Nil(t, err)

	loop = 100
	for i := 0; i < loop; i++ {
		// unstake
		param3 := new(UnStakeParam)
		param3.ConsensusAddress = validatorsKey[0].ConsensusAddr
		param3.Amount = new(big.Int).Mul(big.NewInt(1), params.ZNT1)
		input, err = param3.Encode()
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "UnStake", input, stakeAddressList[0], stakeAddressList[0], blockNumber, extra, sdb)
		assert.Nil(t, err)
	}

	// query
	{
		p1 := &GetEpochInfoParam{
			ID: common.Big1,
		}
		input, err = p1.Encode()
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "GetEpochInfo", input, common.EmptyAddress, common.EmptyAddress, blockNumber, extra, sdb)
		assert.Nil(t, err)

		p2 := &GetAllValidatorsParam{}
		input, err = p2.Encode()
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "GetAllValidators", input, common.EmptyAddress, common.EmptyAddress, blockNumber, extra, sdb)
		assert.Nil(t, err)

		p3 := &GetValidatorParam{
			ConsensusAddress: validatorsKey[0].ConsensusAddr,
		}
		input, err = p3.Encode()
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "GetValidator", input, common.EmptyAddress, common.EmptyAddress, blockNumber, extra, sdb)
		assert.Nil(t, err)

		p4 := &GetStakeInfoParam{
			ConsensusAddress: validatorsKey[0].ConsensusAddr,
			StakeAddress:     stakeAddressList[0],
		}
		input, err = p4.Encode()
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "GetStakeInfo", input, common.EmptyAddress, common.EmptyAddress, blockNumber, extra, sdb)
		assert.Nil(t, err)

		p5 := &GetUnlockingInfoParam{
			StakeAddress: stakeAddressList[0],
		}
		input, err = p5.Encode()
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "GetUnlockingInfo", input, common.EmptyAddress, common.EmptyAddress, blockNumber, extra, sdb)
		assert.Nil(t, err)

		p6 := &GetStakeStartingInfoParam{
			ConsensusAddress: validatorsKey[0].ConsensusAddr,
			StakeAddress:     stakeAddressList[0],
		}
		input, err = p6.Encode()
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "GetStakeStartingInfo", input, common.EmptyAddress, common.EmptyAddress, blockNumber, extra, sdb)
		assert.Nil(t, err)

		p7 := &GetAccumulatedCommissionParam{
			ConsensusAddress: validatorsKey[0].ConsensusAddr,
		}
		input, err = p7.Encode()
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "GetAccumulatedCommission", input, common.EmptyAddress, common.EmptyAddress, blockNumber, extra, sdb)
		assert.Nil(t, err)

		p8 := &GetValidatorSnapshotRewardsParam{
			ConsensusAddress: validatorsKey[0].ConsensusAddr,
			Period:           3,
		}
		input, err = p8.Encode()
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "GetValidatorSnapshotRewards", input, common.EmptyAddress, common.EmptyAddress, blockNumber, extra, sdb)
		assert.Nil(t, err)

		p9 := &GetValidatorAccumulatedRewardsParam{
			ConsensusAddress: validatorsKey[0].ConsensusAddr,
		}
		input, err = p9.Encode()
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "GetValidatorAccumulatedRewards", input, common.EmptyAddress, common.EmptyAddress, blockNumber, extra, sdb)
		assert.Nil(t, err)

		p10 := &GetValidatorOutstandingRewardsParam{
			ConsensusAddress: validatorsKey[0].ConsensusAddr,
		}
		input, err = p10.Encode()
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "GetValidatorOutstandingRewards", input, common.EmptyAddress, common.EmptyAddress, blockNumber, extra, sdb)
		assert.Nil(t, err)

		p11 := &GetTotalPoolParam{}
		input, err = p11.Encode()
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "GetTotalPool", input, common.EmptyAddress, common.EmptyAddress, blockNumber, extra, sdb)
		assert.Nil(t, err)

		p12 := &GetOutstandingRewardsParam{}
		input, err = p12.Encode()
		assert.Nil(t, err)
		_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "GetOutstandingRewards", input, common.EmptyAddress, common.EmptyAddress, blockNumber, extra, sdb)
		assert.Nil(t, err)
	}

	// withdraw
	input, err = utils.PackMethod(ABI, MethodWithdraw)
	assert.Nil(t, err)
	_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "Withdraw", input, stakeAddressList[0], stakeAddressList[0], blockNumber, extra, sdb)
	assert.NotNil(t, err)

	// cancel validator && withdraw validator
	param4 := new(CancelValidatorParam)
	param4.ConsensusAddress = validatorsKey[0].ConsensusAddr
	input, err = param4.Encode()
	assert.Nil(t, err)
	_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "CancelValidator", input, validatorsKey[0].StakeAddress, validatorsKey[0].StakeAddress, blockNumber, extra, sdb)
	assert.Nil(t, err)

	blockNumber = 799999
	// change epoch
	input, err = utils.PackMethod(ABI, MethodChangeEpoch)
	assert.Nil(t, err)
	_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "ChangeEpoch", input, utils.SystemTxSender, utils.SystemTxSender, blockNumber, extra, sdb)
	assert.Nil(t, err)

	blockNumber = 900000
	// withdraw validator
	param5 := new(WithdrawValidatorParam)
	param5.ConsensusAddress = validatorsKey[0].ConsensusAddr
	input, err = param5.Encode()
	assert.Nil(t, err)
	_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "WithdrawValidator", input, validatorsKey[0].StakeAddress, validatorsKey[0].StakeAddress, blockNumber, extra, sdb)
	assert.Nil(t, err)

	// update validator
	param6 := new(UpdateValidatorParam)
	param6.ConsensusAddress = validatorsKey[0].ConsensusAddr
	param6.Desc = "test2"
	input, err = param6.Encode()
	assert.Nil(t, err)
	_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "UpdateValidator", input, validatorsKey[0].StakeAddress, validatorsKey[0].StakeAddress, blockNumber, extra, sdb)
	assert.Nil(t, err)
	param7 := new(UpdateCommissionParam)
	param7.ConsensusAddress = validatorsKey[0].ConsensusAddr
	param7.Commission = new(big.Int).SetUint64(2500)
	input, err = param7.Encode()
	assert.Nil(t, err)
	_, err = native.TestNativeCall(t, utils.NodeManagerContractAddress, "UpdateCommission", input, validatorsKey[0].StakeAddress, validatorsKey[0].StakeAddress, blockNumber, extra, sdb)
	assert.Nil(t, err)
}

func TestCreateValidatorParam(t *testing.T) {
	Init()

	blockNumber := common.Big1
	extra := uint64(21000000000000)

	pk, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(pk.PublicKey)
	caller := crypto.PubkeyToAddress(*acct)
	sdb.SetBalance(caller, new(big.Int).Mul(big.NewInt(100000000), params.ZNT1))

	tests := []struct {
		name   string
		params *CreateValidatorParam
	}{
		{
			name: "invalid consensus address",
			params: &CreateValidatorParam{ConsensusAddress: common.Address{}, SignerAddress: addr, ProposalAddress: addr,
				Commission: new(big.Int).SetUint64(3000), InitStake: new(big.Int).Mul(big.NewInt(100000), params.ZNT1), Desc: "test"},
		},
		{
			name: "invalid signer address",
			params: &CreateValidatorParam{ConsensusAddress: addr, SignerAddress: common.Address{}, ProposalAddress: addr,
				Commission: new(big.Int).SetUint64(3000), InitStake: new(big.Int).Mul(big.NewInt(100000), params.ZNT1), Desc: "test"},
		},
		{
			name: "invalid proposal address",
			params: &CreateValidatorParam{ConsensusAddress: addr, SignerAddress: addr, ProposalAddress: common.Address{},
				Commission: new(big.Int).SetUint64(3000), InitStake: new(big.Int).Mul(big.NewInt(100000), params.ZNT1), Desc: "test"},
		},
		{
			name: "negative commission",
			params: &CreateValidatorParam{ConsensusAddress: addr, SignerAddress: addr, ProposalAddress: addr,
				Commission: new(big.Int).SetInt64(-3000), InitStake: new(big.Int).Mul(big.NewInt(100000), params.ZNT1), Desc: "test"},
		},
		{
			name: "large commission",
			params: &CreateValidatorParam{ConsensusAddress: addr, SignerAddress: addr, ProposalAddress: addr,
				Commission: new(big.Int).SetInt64(30000), InitStake: new(big.Int).Mul(big.NewInt(100000), params.ZNT1), Desc: "test"},
		},
		{
			name: "negative init stake",
			params: &CreateValidatorParam{ConsensusAddress: addr, SignerAddress: addr, ProposalAddress: addr,
				Commission: new(big.Int).SetInt64(3000), InitStake: new(big.Int).Mul(big.NewInt(-100000), params.ZNT1), Desc: "test"},
		},
		{
			name: "small init stake",
			params: &CreateValidatorParam{ConsensusAddress: addr, SignerAddress: addr, ProposalAddress: addr,
				Commission: new(big.Int).SetInt64(3000), InitStake: new(big.Int).Mul(big.NewInt(1000), params.ZNT1), Desc: "test"},
		},
		{
			name: "large desc",
			params: &CreateValidatorParam{ConsensusAddress: addr, SignerAddress: addr, ProposalAddress: addr,
				Commission: new(big.Int).SetInt64(3000), InitStake: new(big.Int).Mul(big.NewInt(100000), params.ZNT1), Desc: "Welcome to the Poly Network cross chain technology center!\n\nBlock technology has now undergone over 10 years of gradual development. Even though many public chain systems have emerged, most of the existing blockchain architecture systems focus on experimenting with the scalability and performance of a single blockchain, while for a single blockchain, due to the constraints that come into picture due to the architecture system and development direction, it is difficult to meet all the needs. Therefore, we expect that the future blockchain ecology must be a pattern in which multiple blockchains coexist. Different chains with different characteristics can become part of the blockchain infrastructure, but the chains are still limited to themselves, each forming its own island of value. Chains that are different in nature lack quick interoperability and convenient means of value circulation. In order to build a better next-generation internet infrastructure, we have launched a new cross-chain technology, the Poly Network.\n\nCross-chain technology is a new technical method to allow inter-chain interactions via cross-chain interoperability protocols that are based on the existing single-blockchain architecture design. The cross-chain interoperability protocol was first proposed by Vitalik Buterin in September 2016. He divided the cross-chain interoperability protocol into three modes: Notary schemes, side chains/relays, and hash locking. The notary mechanism refers to the interaction through a trusted third party as part of the process of cross-chain interaction, which is mainly implemented through mechanisms such as \"single sign, multi-sign notary mechanism\"; the side chain/relay mode is a protocol that is currently being used, for example: Cosmos and Polkadot and other popular cross-chain projects all use this protocol. It refers to a technical methodology that can rely on a trusted third party for cross-chain transaction verification and can conduct cross-chain transaction verification on its own; Hash time lock is the underlying technology of the earliest lightning network. The protocol consists of time lock and hash lock. Time lock refers to the agreement between the two parties that the transaction must be submitted within a certain time to be valid. Hash lock refers to a hash value H, where if the original image R is provided such that Hash (R) = H, the promise is deemed valid, otherwise it is invalid. If the two parties involved in the transaction could not succeed for any reason, time lock can allow the parties to the transaction to get their assets back to avoid losses due to fraud or transaction failure."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := tt.params.Encode()
			assert.Nil(t, err)
			contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
			_, _, err = contractRef.NativeCall(caller, utils.NodeManagerContractAddress, input)
			fmt.Println("#######", err)
			assert.NotNil(t, err)
		})
	}
}

func TestUpdateValidatorParam(t *testing.T) {
	Init()

	blockNumber := common.Big1
	extra := uint64(21000000000000)

	pk, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(pk.PublicKey)
	pk1, _ := crypto.GenerateKey()
	addr1 := crypto.PubkeyToAddress(pk1.PublicKey)
	caller := crypto.PubkeyToAddress(*acct)

	// create validator
	sdb.SetBalance(caller, new(big.Int).Mul(big.NewInt(100000000), params.ZNT1))
	param := new(CreateValidatorParam)
	param.ConsensusAddress = addr
	param.SignerAddress = addr
	param.ProposalAddress = addr
	param.InitStake = new(big.Int).Mul(big.NewInt(100000), params.ZNT1)
	param.Commission = new(big.Int).SetUint64(2000)
	param.Desc = "test"
	input, err := param.Encode()
	assert.Nil(t, err)
	contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(caller, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	tests := []struct {
		name   string
		params *UpdateValidatorParam
	}{
		{
			name:   "large desc",
			params: &UpdateValidatorParam{ConsensusAddress: addr, SignerAddress: addr1, ProposalAddress: common.Address{}, Desc: "Welcome to the Poly Network cross chain technology center!\n\nBlock technology has now undergone over 10 years of gradual development. Even though many public chain systems have emerged, most of the existing blockchain architecture systems focus on experimenting with the scalability and performance of a single blockchain, while for a single blockchain, due to the constraints that come into picture due to the architecture system and development direction, it is difficult to meet all the needs. Therefore, we expect that the future blockchain ecology must be a pattern in which multiple blockchains coexist. Different chains with different characteristics can become part of the blockchain infrastructure, but the chains are still limited to themselves, each forming its own island of value. Chains that are different in nature lack quick interoperability and convenient means of value circulation. In order to build a better next-generation internet infrastructure, we have launched a new cross-chain technology, the Poly Network.\n\nCross-chain technology is a new technical method to allow inter-chain interactions via cross-chain interoperability protocols that are based on the existing single-blockchain architecture design. The cross-chain interoperability protocol was first proposed by Vitalik Buterin in September 2016. He divided the cross-chain interoperability protocol into three modes: Notary schemes, side chains/relays, and hash locking. The notary mechanism refers to the interaction through a trusted third party as part of the process of cross-chain interaction, which is mainly implemented through mechanisms such as \"single sign, multi-sign notary mechanism\"; the side chain/relay mode is a protocol that is currently being used, for example: Cosmos and Polkadot and other popular cross-chain projects all use this protocol. It refers to a technical methodology that can rely on a trusted third party for cross-chain transaction verification and can conduct cross-chain transaction verification on its own; Hash time lock is the underlying technology of the earliest lightning network. The protocol consists of time lock and hash lock. Time lock refers to the agreement between the two parties that the transaction must be submitted within a certain time to be valid. Hash lock refers to a hash value H, where if the original image R is provided such that Hash (R) = H, the promise is deemed valid, otherwise it is invalid. If the two parties involved in the transaction could not succeed for any reason, time lock can allow the parties to the transaction to get their assets back to avoid losses due to fraud or transaction failure."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := tt.params.Encode()
			assert.Nil(t, err)
			contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
			_, _, err = contractRef.NativeCall(caller, utils.NodeManagerContractAddress, input)
			fmt.Println("#######", err)
			assert.NotNil(t, err)
		})
	}
}

func TestUpdateCommissionParam(t *testing.T) {
	Init()

	blockNumber := common.Big1
	extra := uint64(21000000000000)

	pk, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(pk.PublicKey)
	caller := crypto.PubkeyToAddress(*acct)

	// create validator
	sdb.SetBalance(caller, new(big.Int).Mul(big.NewInt(100000000), params.ZNT1))
	param := new(CreateValidatorParam)
	param.ConsensusAddress = addr
	param.SignerAddress = addr
	param.ProposalAddress = addr
	param.InitStake = new(big.Int).Mul(big.NewInt(100000), params.ZNT1)
	param.Commission = new(big.Int).SetUint64(2000)
	param.Desc = "test"
	input, err := param.Encode()
	assert.Nil(t, err)
	contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(caller, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	tests := []struct {
		name   string
		params *UpdateCommissionParam
	}{
		{
			name:   "negative commission",
			params: &UpdateCommissionParam{ConsensusAddress: addr, Commission: new(big.Int).SetInt64(-3000)},
		},
		{
			name:   "too large commission",
			params: &UpdateCommissionParam{ConsensusAddress: addr, Commission: new(big.Int).SetInt64(30000)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := tt.params.Encode()
			assert.Nil(t, err)
			contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
			_, _, err = contractRef.NativeCall(caller, utils.NodeManagerContractAddress, input)
			fmt.Println("#######", err)
			assert.NotNil(t, err)
		})
	}
}

func TestStakeParam(t *testing.T) {
	Init()

	blockNumber := common.Big1
	extra := uint64(21000000000000)

	pk, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(pk.PublicKey)
	caller := crypto.PubkeyToAddress(*acct)

	// create validator
	sdb.SetBalance(caller, new(big.Int).Mul(big.NewInt(100000000), params.ZNT1))
	param := new(CreateValidatorParam)
	param.ConsensusAddress = addr
	param.SignerAddress = addr
	param.ProposalAddress = addr
	param.InitStake = new(big.Int).Mul(big.NewInt(100000), params.ZNT1)
	param.Commission = new(big.Int).SetUint64(2000)
	param.Desc = "test"
	input, err := param.Encode()
	assert.Nil(t, err)
	contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(caller, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	tests := []struct {
		name   string
		params *StakeParam
	}{
		{
			name:   "negative amount",
			params: &StakeParam{ConsensusAddress: addr, Amount: new(big.Int).Mul(big.NewInt(-10), params.ZNT1)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := tt.params.Encode()
			assert.Nil(t, err)
			contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
			_, _, err = contractRef.NativeCall(caller, utils.NodeManagerContractAddress, input)
			fmt.Println("#######", err)
			assert.NotNil(t, err)
		})
	}
}

func TestUnStakeParam(t *testing.T) {
	Init()

	blockNumber := common.Big1
	extra := uint64(21000000000000)

	pk, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(pk.PublicKey)
	caller := crypto.PubkeyToAddress(*acct)

	// create validator
	sdb.SetBalance(caller, new(big.Int).Mul(big.NewInt(100000000), params.ZNT1))
	param := new(CreateValidatorParam)
	param.ConsensusAddress = addr
	param.SignerAddress = addr
	param.ProposalAddress = addr
	param.InitStake = new(big.Int).Mul(big.NewInt(100000), params.ZNT1)
	param.Commission = new(big.Int).SetUint64(2000)
	param.Desc = "test"
	input, err := param.Encode()
	assert.Nil(t, err)
	contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(caller, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	// stake
	param2 := new(StakeParam)
	param2.ConsensusAddress = addr
	param2.Amount = new(big.Int).Mul(big.NewInt(10), params.ZNT1)
	input, err = param2.Encode()
	assert.Nil(t, err)
	contractRef = native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
	_, _, err = contractRef.NativeCall(caller, utils.NodeManagerContractAddress, input)
	assert.Nil(t, err)

	tests := []struct {
		name   string
		params *UnStakeParam
	}{
		{
			name:   "negative amount",
			params: &UnStakeParam{ConsensusAddress: addr, Amount: new(big.Int).Mul(big.NewInt(-10), params.ZNT1)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := tt.params.Encode()
			assert.Nil(t, err)
			contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
			_, _, err = contractRef.NativeCall(caller, utils.NodeManagerContractAddress, input)
			fmt.Println("#######", err)
			assert.NotNil(t, err)
		})
	}
}
