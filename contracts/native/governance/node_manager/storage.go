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
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
)

var StartEpochID = common.Big1 // epoch started from 1, NOT 0!

var ErrEof = errors.New("EOF")

// storage key prefix
const (
	SKP_GLOBAL_CONFIG                 = "st_global_config"
	SKP_VALIDATOR                     = "st_validator"
	SKP_ALL_VALIDATOR                 = "st_all_validator"
	SKP_TOTAL_POOL                    = "st_lock_pool"
	SKP_STAKE_INFO                    = "st_stake_info"
	SKP_UNLOCK_INFO                   = "st_unlock_info"
	SKP_CURRENT_EPOCH                 = "st_current_epoch"
	SKP_EPOCH_INFO                    = "st_epoch_info"
	SKP_ACCUMULATED_COMMISSION        = "st_accumulated_commission"
	SKP_VALIDATOR_ACCUMULATED_REWARDS = "st_validator_accumulated_rewards"
	SKP_VALIDATOR_OUTSTANDING_REWARDS = "st_validator_outstanding_rewards"
	SKP_OUTSTANDING_REWARDS           = "st_outstanding_rewards"
	SKP_VALIDATOR_SNAPSHOT_REWARDS    = "st_validator_snapshot_rewards"
	SKP_STAKE_STARTING_INFO           = "st_stake_starting_info"
	SKP_SIGN                          = "st_sign"
	SKP_SIGNER                        = "st_signer"
	SKP_COMMUNITY_INFO                = "st_community_info"
)

func setAccumulatedCommission(s *native.NativeContract, dec []byte, accumulatedCommission *AccumulatedCommission) error {
	key := accumulatedCommissionKey(dec)
	store, err := rlp.EncodeToBytes(accumulatedCommission)
	if err != nil {
		return fmt.Errorf("setAccumulatedCommission, serialize accumulatedCommission error: %v", err)
	}
	set(s, key, store)
	return nil
}

func getAccumulatedCommission(s *native.NativeContract, dec []byte) (*AccumulatedCommission, error) {
	accumulatedCommission := &AccumulatedCommission{}
	key := accumulatedCommissionKey(dec)
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("getAccumulatedCommission, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, accumulatedCommission); err != nil {
		return nil, fmt.Errorf("getAccumulatedCommission, deserialize accumulatedCommission error: %v", err)
	}
	return accumulatedCommission, nil
}

func delAccumulatedCommission(s *native.NativeContract, dec []byte) {
	key := accumulatedCommissionKey(dec)
	del(s, key)
}

func setValidatorAccumulatedRewards(s *native.NativeContract, dec []byte, validatorAccumulatedRewards *ValidatorAccumulatedRewards) error {
	key := validatorAccumulatedRewardsKey(dec)
	store, err := rlp.EncodeToBytes(validatorAccumulatedRewards)
	if err != nil {
		return fmt.Errorf("setValidatorAccumulatedRewards, serialize validatorAccumulatedRewards error: %v", err)
	}
	set(s, key, store)
	return nil
}

func getValidatorAccumulatedRewards(s *native.NativeContract, dec []byte) (*ValidatorAccumulatedRewards, error) {
	validatorAccumulatedRewards := &ValidatorAccumulatedRewards{}
	key := validatorAccumulatedRewardsKey(dec)
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("GetValidatorAccumulatedRewards, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, validatorAccumulatedRewards); err != nil {
		return nil, fmt.Errorf("GetValidatorAccumulatedRewards, deserialize validatorAccumulatedRewards error: %v", err)
	}
	return validatorAccumulatedRewards, nil
}

func delValidatorAccumulatedRewards(s *native.NativeContract, dec []byte) {
	key := validatorAccumulatedRewardsKey(dec)
	del(s, key)
}

func setValidatorOutstandingRewards(s *native.NativeContract, dec []byte, validatorOutstandingRewards *ValidatorOutstandingRewards) error {
	key := validatorOutstandingRewardsKey(dec)
	store, err := rlp.EncodeToBytes(validatorOutstandingRewards)
	if err != nil {
		return fmt.Errorf("setValidatorOutstandingRewards, serialize validatorOutstandingRewards error: %v", err)
	}
	set(s, key, store)
	return nil
}

func getValidatorOutstandingRewards(s *native.NativeContract, dec []byte) (*ValidatorOutstandingRewards, error) {
	validatorOutstandingRewards := &ValidatorOutstandingRewards{}
	key := validatorOutstandingRewardsKey(dec)
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("getValidatorOutstandingRewards, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, validatorOutstandingRewards); err != nil {
		return nil, fmt.Errorf("getValidatorOutstandingRewards, deserialize validatorOutstandingRewards error: %v", err)
	}
	return validatorOutstandingRewards, nil
}

func delValidatorOutstandingRewards(s *native.NativeContract, dec []byte) {
	key := validatorOutstandingRewardsKey(dec)
	del(s, key)
}

func setOutstandingRewards(s *native.NativeContract, outstandingRewards *OutstandingRewards) error {
	key := outstandingRewardsKey()
	store, err := rlp.EncodeToBytes(outstandingRewards)
	if err != nil {
		return fmt.Errorf("setOutstandingRewards, serialize outstandingRewards error: %v", err)
	}
	set(s, key, store)
	return nil
}

func getOutstandingRewards(s *native.NativeContract) (*OutstandingRewards, error) {
	outstandingRewards := &OutstandingRewards{
		Rewards: NewDecFromBigInt(new(big.Int)),
	}
	key := outstandingRewardsKey()
	store, err := get(s, key)
	if err == ErrEof {
		return outstandingRewards, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getOutstandingRewards, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, outstandingRewards); err != nil {
		return nil, fmt.Errorf("getOutstandingRewards, deserialize outstandingRewards error: %v", err)
	}
	return outstandingRewards, nil
}

func increaseReferenceCount(s *native.NativeContract, dec []byte, period uint64) error {
	validatorSnapshotRewards, err := getValidatorSnapshotRewards(s, dec, period)
	if err != nil {
		return fmt.Errorf("increaseReferenceCount, getValidatorSnapshotRewards error: %v", err)
	}
	if validatorSnapshotRewards.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	validatorSnapshotRewards.ReferenceCount++
	err = setValidatorSnapshotRewards(s, dec, period, validatorSnapshotRewards)
	if err != nil {
		return fmt.Errorf("increaseReferenceCount, setValidatorSnapshotRewards error: %v", err)
	}
	return nil
}

func decreaseReferenceCount(s *native.NativeContract, dec []byte, period uint64) error {
	validatorSnapshotRewards, err := getValidatorSnapshotRewards(s, dec, period)
	if err != nil {
		return fmt.Errorf("decreaseReferenceCount, getValidatorSnapshotRewards error: %v", err)
	}
	if validatorSnapshotRewards.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	validatorSnapshotRewards.ReferenceCount--
	if validatorSnapshotRewards.ReferenceCount == 0 {
		delValidatorSnapshotRewards(s, dec, period)
	} else {
		err = setValidatorSnapshotRewards(s, dec, period, validatorSnapshotRewards)
		if err != nil {
			return fmt.Errorf("decreaseReferenceCount, setValidatorSnapshotRewards error: %v", err)
		}
	}
	return nil
}

func setValidatorSnapshotRewards(s *native.NativeContract, dec []byte, period uint64, validatorSnapshotRewards *ValidatorSnapshotRewards) error {
	key := validatorSnapshotRewardsKey(dec, period)
	store, err := rlp.EncodeToBytes(validatorSnapshotRewards)
	if err != nil {
		return fmt.Errorf("setValidatorSnapshotRewards, serialize validatorSnapshotRewards error: %v", err)
	}
	set(s, key, store)
	return nil
}

func getValidatorSnapshotRewards(s *native.NativeContract, dec []byte, period uint64) (*ValidatorSnapshotRewards, error) {
	validatorSnapshotRewards := &ValidatorSnapshotRewards{}
	key := validatorSnapshotRewardsKey(dec, period)
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("getValidatorSnapshotRewards, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, validatorSnapshotRewards); err != nil {
		return nil, fmt.Errorf("getValidatorSnapshotRewards, deserialize validatorSnapshotRewards error: %v", err)
	}
	return validatorSnapshotRewards, nil
}

func delValidatorSnapshotRewards(s *native.NativeContract, dec []byte, period uint64) {
	key := validatorSnapshotRewardsKey(dec, period)
	del(s, key)
}

func setStakeStartingInfo(s *native.NativeContract, stakeAddress common.Address, dec []byte, stakeStartingInfo *StakeStartingInfo) error {
	key := stakeStartingInfoKey(stakeAddress, dec)
	store, err := rlp.EncodeToBytes(stakeStartingInfo)
	if err != nil {
		return fmt.Errorf("setStakeStartingInfo, serialize stakeStartingInfo error: %v", err)
	}
	set(s, key, store)
	return nil
}

func getStakeStartingInfo(s *native.NativeContract, stakeAddress common.Address, dec []byte) (*StakeStartingInfo, error) {
	stakeStartingInfo := &StakeStartingInfo{}
	key := stakeStartingInfoKey(stakeAddress, dec)
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("getStakeStartingInfo, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, stakeStartingInfo); err != nil {
		return nil, fmt.Errorf("getStakeStartingInfo, deserialize stakeStartingInfo error: %v", err)
	}
	return stakeStartingInfo, nil
}

func delStakeStartingInfo(s *native.NativeContract, stakeAddress common.Address, dec []byte) {
	key := stakeStartingInfoKey(stakeAddress, dec)
	del(s, key)
}

func setGlobalConfig(s *native.NativeContract, globalConfig *GlobalConfig) error {
	key := globalConfigKey()
	store, err := rlp.EncodeToBytes(globalConfig)
	if err != nil {
		return fmt.Errorf("setGlobalConfig, serialize globalConfig error: %v", err)
	}
	set(s, key, store)
	return nil
}

func setGenesisGlobalConfig(s *state.CacheDB, globalConfig *GlobalConfig) error {
	key := globalConfigKey()
	store, err := rlp.EncodeToBytes(globalConfig)
	if err != nil {
		return fmt.Errorf("setGenesisGlobalConfig, serialize globalConfig error: %v", err)
	}
	customSet(s, key, store)
	return nil
}

func getGlobalConfig(s *native.NativeContract) (*GlobalConfig, error) {
	key := globalConfigKey()
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("GetGlobalConfig, get store error: %v", err)
	}
	globalConfig := new(GlobalConfig)
	if err := rlp.DecodeBytes(store, globalConfig); err != nil {
		return nil, fmt.Errorf("GetGlobalConfig, deserialize globalConfig error: %v", err)
	}
	return globalConfig, nil
}

func GetGlobalConfigFromDB(s *state.StateDB) (*GlobalConfig, error) {
	cache := (*state.CacheDB)(s)

	globalConfig := new(GlobalConfig)
	key := globalConfigKey()
	store, err := customGet(cache, key)
	if err != nil {
		return nil, fmt.Errorf("GetGlobalConfigFromDB, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, globalConfig); err != nil {
		return nil, fmt.Errorf("GetGlobalConfigFromDB, deserialize globalConfig error: %v", err)
	}
	return globalConfig, nil
}

func addToAllValidators(s *native.NativeContract, consensusPk string) error {
	allValidators, err := getAllValidators(s)
	if err != nil {
		return fmt.Errorf("addToAllValidators, getAllValidators error: %v", err)
	}
	allValidators.AllValidators = append(allValidators.AllValidators, consensusPk)
	err = setAllValidators(s, allValidators)
	if err != nil {
		return fmt.Errorf("addToAllValidators, set all validators error: %v", err)
	}
	return nil
}

func removeFromAllValidators(s *native.NativeContract, consensusPk string) error {
	allValidators, err := getAllValidators(s)
	if err != nil {
		return fmt.Errorf("removeFromAllValidators, getAllValidators error: %v", err)
	}
	j := 0
	for _, validator := range allValidators.AllValidators {
		if validator != consensusPk {
			allValidators.AllValidators[j] = validator
			j++
		}
	}
	allValidators.AllValidators = allValidators.AllValidators[:j]
	err = setAllValidators(s, allValidators)
	if err != nil {
		return fmt.Errorf("removeFromAllValidators, set all validators error: %v", err)
	}
	return nil
}

func setValidator(s *native.NativeContract, validator *Validator) error {
	dec, err := hexutil.Decode(validator.ConsensusPubkey)
	if err != nil {
		return err
	}
	key := validatorKey(dec)
	store, err := rlp.EncodeToBytes(validator)
	if err != nil {
		return fmt.Errorf("setValidator, serialize validator error: %v", err)
	}
	set(s, key, store)
	return nil
}

func delValidator(s *native.NativeContract, consensusPk string) error {
	dec, err := hexutil.Decode(consensusPk)
	if err != nil {
		return err
	}
	key := validatorKey(dec)
	del(s, key)
	return nil
}

func getValidator(s *native.NativeContract, dec []byte) (*Validator, bool, error) {
	key := validatorKey(dec)
	store, err := get(s, key)
	if err == ErrEof {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("getValidator, get store error: %v", err)
	}
	validator := new(Validator)
	if err := rlp.DecodeBytes(store, validator); err != nil {
		return nil, false, fmt.Errorf("getValidator, deserialize validator error: %v", err)
	}
	return validator, true, nil
}

func setAllValidators(s *native.NativeContract, allValidators *AllValidators) error {
	key := allValidatorKey()
	store, err := rlp.EncodeToBytes(allValidators)
	if err != nil {
		return fmt.Errorf("setAllValidators, serialize all validators error: %v", err)
	}
	set(s, key, store)
	return nil
}

func getAllValidators(s *native.NativeContract) (*AllValidators, error) {
	allValidators := &AllValidators{
		AllValidators: make([]string, 0),
	}
	key := allValidatorKey()
	store, err := get(s, key)
	if err == ErrEof {
		return allValidators, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getAllValidators, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, allValidators); err != nil {
		return nil, fmt.Errorf("getAllValidators, deserialize all validators error: %v", err)
	}
	return allValidators, nil
}

func depositTotalPool(s *native.NativeContract, amount Dec) error {
	totalPool, err := getTotalPool(s)
	if err != nil {
		return fmt.Errorf("depositTotalPool, get total pool error: %v", err)
	}
	totalPool.TotalPool, err = totalPool.TotalPool.Add(amount)
	if err != nil {
		return fmt.Errorf("depositTotalPool, totalPool.TotalPool.Add error: %v", err)
	}
	err = setTotalPool(s, totalPool)
	if err != nil {
		return fmt.Errorf("depositTotalPool, setTotalPool error: %v", err)
	}
	return nil
}

func withdrawTotalPool(s *native.NativeContract, amount Dec) error {
	totalPool, err := getTotalPool(s)
	if err != nil {
		return fmt.Errorf("withdrawTotalPool, get total pool error: %v", err)
	}
	totalPool.TotalPool, err = totalPool.TotalPool.Sub(amount)
	if err != nil {
		return fmt.Errorf("withdrawTotalPool, totalPool.Sub error: %v", err)
	}
	err = setTotalPool(s, totalPool)
	if err != nil {
		return fmt.Errorf("withdrawTotalPool, setTotalPool error: %v", err)
	}
	return nil
}

func setTotalPool(s *native.NativeContract, totalPool *TotalPool) error {
	key := totalPoolKey()
	store, err := rlp.EncodeToBytes(totalPool)
	if err != nil {
		return fmt.Errorf("setStakeInfo, serialize stake info error: %v", err)
	}
	set(s, key, store)
	return nil
}

func getTotalPool(s *native.NativeContract) (*TotalPool, error) {
	totalPool := &TotalPool{NewDecFromBigInt(new(big.Int))}
	key := totalPoolKey()
	store, err := get(s, key)
	if err == ErrEof {
		return totalPool, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getTotalPool, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, totalPool); err != nil {
		return nil, fmt.Errorf("getTotalPool, deserialize totalPool error: %v", err)
	}
	return totalPool, nil
}

func setStakeInfo(s *native.NativeContract, stakeInfo *StakeInfo) error {
	dec, err := hexutil.Decode(stakeInfo.ConsensusPubkey)
	if err != nil {
		return err
	}
	key := stakeInfoKey(stakeInfo.StakeAddress, dec)
	store, err := rlp.EncodeToBytes(stakeInfo)
	if err != nil {
		return fmt.Errorf("setStakeInfo, serialize stake info error: %v", err)
	}
	set(s, key, store)
	return nil
}

func delStakeInfo(s *native.NativeContract, stakeAddress common.Address, consensusPk string) error {
	dec, err := hexutil.Decode(consensusPk)
	if err != nil {
		return err
	}
	key := stakeInfoKey(stakeAddress, dec)
	del(s, key)
	return nil
}

func getStakeInfo(s *native.NativeContract, stakeAddress common.Address, consensusPk string) (*StakeInfo, bool, error) {
	stakeInfo := &StakeInfo{
		StakeAddress:    stakeAddress,
		ConsensusPubkey: consensusPk,
		Amount:          NewDecFromBigInt(new(big.Int)),
	}
	dec, err := hexutil.Decode(stakeInfo.ConsensusPubkey)
	if err != nil {
		return nil, false, err
	}
	key := stakeInfoKey(stakeAddress, dec)
	store, err := get(s, key)
	if err == ErrEof {
		return stakeInfo, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("getStakeInfo, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, stakeInfo); err != nil {
		return nil, false, fmt.Errorf("getStakeInfo, deserialize stakeInfo error: %v", err)
	}
	return stakeInfo, true, nil
}

func addUnlockingInfo(s *native.NativeContract, stakeAddress common.Address, unlockingStake *UnlockingStake) error {
	unlockingInfo, err := getUnlockingInfo(s, stakeAddress)
	if err != nil {
		return fmt.Errorf("addUnlockingInfo, GetUnlockingInfo error: %v", err)
	}
	unlockingInfo.UnlockingStake = append(unlockingInfo.UnlockingStake, unlockingStake)
	err = setUnlockingInfo(s, unlockingInfo)
	if err != nil {
		return fmt.Errorf("addUnlockingInfo, setUnlockingInfo error: %v", err)
	}
	return nil
}

func filterExpiredUnlockingInfo(s *native.NativeContract, stakeAddress common.Address) (Dec, error) {
	height := s.ContractRef().BlockHeight()
	unlockingInfo, err := getUnlockingInfo(s, stakeAddress)
	if err != nil {
		return Dec{nil}, fmt.Errorf("filterExpiredUnlockingInfo, GetUnlockingInfo error: %v", err)
	}
	j := 0
	expiredSum := NewDecFromBigInt(new(big.Int))
	for _, unlockingStake := range unlockingInfo.UnlockingStake {
		if unlockingStake.CompleteHeight.Cmp(height) == 1 {
			unlockingInfo.UnlockingStake[j] = unlockingStake
			j++
		} else {
			expiredSum, err = expiredSum.Add(unlockingStake.Amount)
			if err != nil {
				return Dec{nil}, fmt.Errorf("filterExpiredUnlockingInfo, expiredSum.Add error: %v", err)
			}
		}
	}
	unlockingInfo.UnlockingStake = unlockingInfo.UnlockingStake[:j]
	if len(unlockingInfo.UnlockingStake) == 0 {
		delUnlockingInfo(s, stakeAddress)
	} else {
		err = setUnlockingInfo(s, unlockingInfo)
		if err != nil {
			return Dec{nil}, fmt.Errorf("filterExpiredUnlockingInfo, setUnlockingInfo error: %v", err)
		}
	}
	return expiredSum, nil
}

func setUnlockingInfo(s *native.NativeContract, unlockingInfo *UnlockingInfo) error {
	key := unlockingInfoKey(unlockingInfo.StakeAddress)
	store, err := rlp.EncodeToBytes(unlockingInfo)
	if err != nil {
		return fmt.Errorf("setUnlockingInfo, serialize unlock info error: %v", err)
	}
	set(s, key, store)
	return nil
}

func delUnlockingInfo(s *native.NativeContract, stakeAddress common.Address) {
	key := unlockingInfoKey(stakeAddress)
	del(s, key)
}

func getUnlockingInfo(s *native.NativeContract, stakeAddress common.Address) (*UnlockingInfo, error) {
	unlockingInfo := &UnlockingInfo{
		StakeAddress:   stakeAddress,
		UnlockingStake: make([]*UnlockingStake, 0),
	}
	key := unlockingInfoKey(stakeAddress)
	store, err := get(s, key)
	if err == ErrEof {
		return unlockingInfo, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetUnlockingInfo, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, unlockingInfo); err != nil {
		return nil, fmt.Errorf("GetUnlockingInfo, deserialize unlocking info error: %v", err)
	}
	return unlockingInfo, nil
}

func setCurrentEpoch(s *native.NativeContract, ID *big.Int) {
	key := currentEpochKey()
	set(s, key, ID.Bytes())
}

func getCurrentEpoch(s *native.NativeContract) (*big.Int, error) {
	key := currentEpochKey()
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpoch, get store error: %v", err)
	}
	return new(big.Int).SetBytes(store), nil
}

func setCurrentEpochInfo(s *native.NativeContract, epochInfo *EpochInfo) error {
	// set current epoch
	setCurrentEpoch(s, epochInfo.ID)
	//set epoch info
	err := setEpochInfo(s, epochInfo)
	if err != nil {
		return fmt.Errorf("setCurrentEpochInfo, setEpochInfo error: %v", err)
	}
	return nil
}

func setGenesisEpochInfo(s *state.CacheDB, epochInfo *EpochInfo) error {
	// set current epoch
	key1 := currentEpochKey()
	customSet(s, key1, epochInfo.ID.Bytes())
	//set epoch info
	key2 := epochInfoKey(epochInfo.ID)
	store, err := rlp.EncodeToBytes(epochInfo)
	if err != nil {
		return fmt.Errorf("setGenesisEpochInfo, serialize epoch info error: %v", err)
	}
	customSet(s, key2, store)
	return nil
}

func GetCurrentEpochInfoImpl(s *native.NativeContract) (*EpochInfo, error) {
	ID, err := getCurrentEpoch(s)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpochInfoImpl, getCurrentEpochInfo error: %v", err)
	}
	epochInfo, err := getEpochInfo(s, ID)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpochInfoImpl, getEpochInfo error: %v", err)
	}
	return epochInfo, nil
}

func GetCurrentEpochInfoFromDB(s *state.StateDB) (*EpochInfo, error) {
	cache := (*state.CacheDB)(s)
	key := currentEpochKey()
	store, err := customGet(cache, key)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpochInfoFromDB, get store error: %v", err)
	}
	ID := new(big.Int).SetBytes(store)

	epochInfo := &EpochInfo{
		Validators: make([]*Peer, 0),
	}
	key = epochInfoKey(ID)
	store, err = customGet(cache, key)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpochInfoFromDB, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, epochInfo); err != nil {
		return nil, fmt.Errorf("GetCurrentEpochInfoFromDB, deserialize epoch info error: %v", err)
	}
	return epochInfo, nil
}

func setEpochInfo(s *native.NativeContract, epochInfo *EpochInfo) error {
	key := epochInfoKey(epochInfo.ID)
	store, err := rlp.EncodeToBytes(epochInfo)
	if err != nil {
		return fmt.Errorf("setEpochInfo, serialize epoch info error: %v", err)
	}
	set(s, key, store)
	return nil
}

func getEpochInfo(s *native.NativeContract, ID *big.Int) (*EpochInfo, error) {
	epochInfo := &EpochInfo{
		Validators: make([]*Peer, 0),
	}
	key := epochInfoKey(ID)
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("GetEpochInfo, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, epochInfo); err != nil {
		return nil, fmt.Errorf("GetEpochInfo, deserialize epoch info error: %v", err)
	}
	return epochInfo, nil
}

func GetEpochInfoFromDB(s *state.StateDB, ID *big.Int) (*EpochInfo, error) {
	cache := (*state.CacheDB)(s)

	epochInfo := &EpochInfo{
		Validators: make([]*Peer, 0),
	}
	key := epochInfoKey(ID)
	store, err := customGet(cache, key)
	if err != nil {
		return nil, fmt.Errorf("GetEpochInfoFromDB, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, epochInfo); err != nil {
		return nil, fmt.Errorf("GetEpochInfoFromDB, deserialize epoch info error: %v", err)
	}
	return epochInfo, nil
}

func setGenesisCommunityInfo(s *state.CacheDB, communityInfo *CommunityInfo) error {
	key := communityInfoKey()
	store, err := rlp.EncodeToBytes(communityInfo)
	if err != nil {
		return fmt.Errorf("setCommunityInfo, serialize community info error: %v", err)
	}
	customSet(s, key, store)
	return nil
}

func setCommunityInfo(s *native.NativeContract, communityInfo *CommunityInfo) error {
	key := communityInfoKey()
	store, err := rlp.EncodeToBytes(communityInfo)
	if err != nil {
		return fmt.Errorf("setCommunityInfo, serialize community info error: %v", err)
	}
	set(s, key, store)
	return nil
}

func getCommunityInfo(s *native.NativeContract) (*CommunityInfo, error) {
	communityInfo := new(CommunityInfo)
	key := communityInfoKey()
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("GetCommunityInfo, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, communityInfo); err != nil {
		return nil, fmt.Errorf("GetCommunityInfo, deserialize community info error: %v", err)
	}
	return communityInfo, nil
}

func GetCommunityInfoFromDB(s *state.StateDB) (*CommunityInfo, error) {
	cache := (*state.CacheDB)(s)
	communityInfo := new(CommunityInfo)
	key := communityInfoKey()
	store, err := customGet(cache, key)
	if err != nil {
		return nil, fmt.Errorf("GetCommunityInfoFromDB, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, communityInfo); err != nil {
		return nil, fmt.Errorf("GetCommunityInfoFromDB, deserialize community info error: %v", err)
	}
	return communityInfo, nil
}

// ====================================================================
//
// `consensus sign` storage
//
// ====================================================================
func storeSign(s *native.NativeContract, sign *ConsensusSign) error {
	key := signKey(sign.Hash())
	value, err := rlp.EncodeToBytes(sign)
	if err != nil {
		return err
	}
	set(s, key, value)
	return nil
}

func delSign(s *native.NativeContract, hash common.Hash) {
	key := signKey(hash)
	del(s, key)
}

func getSign(s *native.NativeContract, hash common.Hash) (*ConsensusSign, error) {
	key := signKey(hash)
	value, err := get(s, key)
	if err != nil {
		return nil, err
	}
	var sign *ConsensusSign
	if err := rlp.DecodeBytes(value, &sign); err != nil {
		return nil, err
	}
	return sign, nil
}

func storeSigner(s *native.NativeContract, hash common.Hash, signer common.Address) error {
	data, err := getSigners(s, hash)
	if err != nil {
		if err.Error() == ErrEof.Error() {
			data = make([]common.Address, 0)
		} else {
			return err
		}
	}
	data = append(data, signer)
	list := &AddressList{List: data}

	key := signerKey(hash)
	value, err := rlp.EncodeToBytes(list)
	if err != nil {
		return err
	}
	set(s, key, value)

	return nil
}

func findSigner(s *native.NativeContract, hash common.Hash, signer common.Address) bool {
	list, err := getSigners(s, hash)
	if err != nil {
		return false
	}
	for _, v := range list {
		if v == signer {
			return true
		}
	}
	return false
}

func getSigners(s *native.NativeContract, hash common.Hash) ([]common.Address, error) {
	key := signerKey(hash)
	value, err := get(s, key)
	if err != nil {
		return nil, err
	}

	var list *AddressList
	if err := rlp.DecodeBytes(value, &list); err != nil {
		return nil, err
	}
	return list.List, nil
}

func getSignerSize(s *native.NativeContract, hash common.Hash) int {
	list, err := getSigners(s, hash)
	if err != nil {
		return 0
	}
	return len(list)
}

func clearSigner(s *native.NativeContract, hash common.Hash) {
	key := signerKey(hash)
	del(s, key)
}

// ====================================================================
//
// storage basic operations
//
// ====================================================================

func get(s *native.NativeContract, key []byte) ([]byte, error) {
	return customGet(s.GetCacheDB(), key)
}

func set(s *native.NativeContract, key, value []byte) {
	customSet(s.GetCacheDB(), key, value)
}

func del(s *native.NativeContract, key []byte) {
	customDel(s.GetCacheDB(), key)
}

func customGet(db *state.CacheDB, key []byte) ([]byte, error) {
	value, err := db.Get(key)
	if err != nil {
		return nil, err
	} else if value == nil || len(value) == 0 {
		return nil, ErrEof
	} else {
		return value, nil
	}
}

func customSet(db *state.CacheDB, key, value []byte) {
	db.Put(key, value)
}

func customDel(db *state.CacheDB, key []byte) {
	db.Delete(key)
}

// ====================================================================
//
// storage keys
//
// ====================================================================

func globalConfigKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_GLOBAL_CONFIG))
}

func validatorKey(dec []byte) []byte {
	return utils.ConcatKey(this, []byte(SKP_VALIDATOR), dec)
}

func allValidatorKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_ALL_VALIDATOR))
}

func totalPoolKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_TOTAL_POOL))
}

func stakeInfoKey(stakeAddress common.Address, dec []byte) []byte {
	return utils.ConcatKey(this, []byte(SKP_STAKE_INFO), stakeAddress[:], dec)
}

func unlockingInfoKey(stakeAddress common.Address) []byte {
	return utils.ConcatKey(this, []byte(SKP_UNLOCK_INFO), stakeAddress[:])
}

func currentEpochKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_CURRENT_EPOCH))
}

func epochInfoKey(ID *big.Int) []byte {
	return utils.ConcatKey(this, []byte(SKP_EPOCH_INFO), ID.Bytes())
}

func accumulatedCommissionKey(dec []byte) []byte {
	return utils.ConcatKey(this, []byte(SKP_ACCUMULATED_COMMISSION), dec)
}

func validatorAccumulatedRewardsKey(dec []byte) []byte {
	return utils.ConcatKey(this, []byte(SKP_VALIDATOR_ACCUMULATED_REWARDS), dec)
}

func validatorOutstandingRewardsKey(dec []byte) []byte {
	return utils.ConcatKey(this, []byte(SKP_VALIDATOR_OUTSTANDING_REWARDS), dec)
}

func outstandingRewardsKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_OUTSTANDING_REWARDS))
}

func validatorSnapshotRewardsKey(dec []byte, period uint64) []byte {
	return utils.ConcatKey(this, []byte(SKP_VALIDATOR_SNAPSHOT_REWARDS), dec, utils.Uint64Bytes(period))
}

func stakeStartingInfoKey(stakeAddress common.Address, dec []byte) []byte {
	return utils.ConcatKey(this, []byte(SKP_STAKE_STARTING_INFO), stakeAddress[:], dec)
}

func signKey(hash common.Hash) []byte {
	return utils.ConcatKey(this, []byte(SKP_SIGN), hash.Bytes())
}

func signerKey(hash common.Hash) []byte {
	return utils.ConcatKey(this, []byte(SKP_SIGNER), hash.Bytes())
}

func communityInfoKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_COMMUNITY_INFO))
}
