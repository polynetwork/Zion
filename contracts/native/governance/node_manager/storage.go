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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
)

const (
	StartEpochID uint64 = 1 // epoch started from 1, NOT 0!
)

var ErrEof = errors.New("EOF")

// storage key prefix
const (
	SKP_GLOBAL_CONFIG = "st_global_config"
	SKP_VALIDATOR     = "st_validator"
	SKP_ALL_VALIDATOR = "st_all_validator"
	SKP_LOCK_POOL     = "st_lock_pool"
	SKP_UNLOCK_POOL   = "st_unlock_pool"
	SKP_STAKE_INFO    = "st_stake_info"
	SKP_UNLOCK_INFO   = "st_unlock_info"
	SKP_CURRENT_EPOCH = "st_current_epoch"
	SKP_EPOCH_INFO    = "st_epoch_info"
)

func setGlobalConfig(s *native.NativeContract, globalConfig *GlobalConfig) error {
	key := globalConfigKey()
	store, err := rlp.EncodeToBytes(globalConfig)
	if err != nil {
		return fmt.Errorf("setGlobalConfig, serialize globalConfig error: %v", err)
	}
	set(s, key, store)
	return nil
}

func GetGlobalConfig(s *native.NativeContract) (*GlobalConfig, error) {
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

func addToAllValidators(s *native.NativeContract, consensusPk string) error {
	allValidators, err := GetAllValidators(s)
	if err != nil {
		return fmt.Errorf("addToAllValidators, GetAllValidators error: %v", err)
	}
	allValidators.AllValidators = append(allValidators.AllValidators, consensusPk)
	err = setAllValidators(s, allValidators)
	if err != nil {
		return fmt.Errorf("addToAllValidators, set all validators error: %v", err)
	}
	return nil
}

func removeFromAllValidators(s *native.NativeContract, consensusPk string) error {
	allValidators, err := GetAllValidators(s)
	if err != nil {
		return fmt.Errorf("removeFromAllValidators, GetAllValidators error: %v", err)
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

func GetValidator(s *native.NativeContract, dec []byte) (*Validator, bool, error) {
	key := validatorKey(dec)
	store, err := get(s, key)
	if err == ErrEof {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("GetValidator, get store error: %v", err)
	}
	validator := new(Validator)
	if err := rlp.DecodeBytes(store, validator); err != nil {
		return nil, false, fmt.Errorf("GetValidator, deserialize validator error: %v", err)
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

func GetAllValidators(s *native.NativeContract) (*AllValidators, error) {
	allValidators := &AllValidators{
		AllValidators: make([]string, 0),
	}
	key := allValidatorKey()
	store, err := get(s, key)
	if err == ErrEof {
		return allValidators, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetAllValidators, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, allValidators); err != nil {
		return nil, fmt.Errorf("GetAllValidators, deserialize all validators error: %v", err)
	}
	return allValidators, nil
}

func depositLockPool(s *native.NativeContract, amount *big.Int) error {
	lockPool, err := GetLockPool(s)
	if err != nil {
		return fmt.Errorf("depositLockPool, get lock pool error: %v", err)
	}
	lockPool = new(big.Int).Add(lockPool, amount)
	setLockPool(s, lockPool)
	return nil
}

func withdrawLockPool(s *native.NativeContract, amount *big.Int) error {
	lockPool, err := GetLockPool(s)
	if err != nil {
		return fmt.Errorf("withdrawLockPool, get lock pool error: %v", err)
	}
	lockPool = new(big.Int).Sub(lockPool, amount)
	if lockPool.Sign() < 0 {
		return fmt.Errorf("withdrawLockPool, lock pool is less than amount, please check")
	}
	setLockPool(s, lockPool)
	return nil
}

func depositUnlockPool(s *native.NativeContract, amount *big.Int) error {
	unlockPool, err := GetUnlockPool(s)
	if err != nil {
		return fmt.Errorf("depositUnlockPool, get unlock pool error: %v", err)
	}
	unlockPool = new(big.Int).Add(unlockPool, amount)
	setUnlockPool(s, unlockPool)
	return nil
}

func withdrawUnlockPool(s *native.NativeContract, amount *big.Int) error {
	unlockPool, err := GetUnlockPool(s)
	if err != nil {
		return fmt.Errorf("withdrawUnlockPool, get lock pool error: %v", err)
	}
	unlockPool = new(big.Int).Sub(unlockPool, amount)
	if unlockPool.Sign() < 0 {
		return fmt.Errorf("withdrawUnlockPool, unlock pool is less than amount, please check")
	}
	setUnlockPool(s, unlockPool)
	return nil
}

func setLockPool(s *native.NativeContract, amount *big.Int) {
	key := lockPoolKey()
	set(s, key, amount.Bytes())
}

func GetLockPool(s *native.NativeContract) (*big.Int, error) {
	key := lockPoolKey()
	store, err := get(s, key)
	if err == ErrEof {
		return common.Big0, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetLockPool, get store error: %v", err)
	}
	return new(big.Int).SetBytes(store), nil
}

func setUnlockPool(s *native.NativeContract, amount *big.Int) {
	key := unlockPoolKey()
	set(s, key, amount.Bytes())
}

func GetUnlockPool(s *native.NativeContract) (*big.Int, error) {
	key := unlockPoolKey()
	store, err := get(s, key)
	if err == ErrEof {
		return common.Big0, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetUnlockPool, get store error: %v", err)
	}
	return new(big.Int).SetBytes(store), nil
}

func depositStakeInfo(s *native.NativeContract, stakeAddress common.Address, consensusPk string, amount *big.Int) error {
	stakeInfo, err := GetStakeInfo(s, stakeAddress, consensusPk)
	if err != nil {
		return fmt.Errorf("depositStakeInfo, get stake info error: %v", err)
	}
	stakeInfo.Amount = new(big.Int).Add(stakeInfo.Amount, amount)
	err = setStakeInfo(s, stakeInfo)
	if err != nil {
		return fmt.Errorf("depositStakeInfo, set stake info error: %v", err)
	}
	return nil
}

func withdrawStakeInfo(s *native.NativeContract, stakeAddress common.Address, consensusPk string, amount *big.Int) error {
	stakeInfo, err := GetStakeInfo(s, stakeAddress, consensusPk)
	if err != nil {
		return fmt.Errorf("withdrawStakeInfo, get stake info error: %v", err)
	}
	if stakeInfo.Amount.Cmp(amount) == -1 {
		return fmt.Errorf("withdrawStakeInfo, stake info is less than amount")
	}
	stakeInfo.Amount = new(big.Int).Sub(stakeInfo.Amount, amount)
	if stakeInfo.Amount.Sign() == 0 {
		err = delStakeInfo(s, stakeAddress, consensusPk)
		if err != nil {
			return fmt.Errorf("withdrawStakeInfo, delete stake info error: %v", err)
		}
	} else {
		err = setStakeInfo(s, stakeInfo)
		if err != nil {
			return fmt.Errorf("withdrawStakeInfo, set stake info error: %v", err)
		}
	}
	return nil
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

func GetStakeInfo(s *native.NativeContract, stakeAddress common.Address, consensusPk string) (*StakeInfo, error) {
	stakeInfo := &StakeInfo{
		StakeAddress:    stakeAddress,
		ConsensusPubkey: consensusPk,
	}
	dec, err := hexutil.Decode(stakeInfo.ConsensusPubkey)
	if err != nil {
		return nil, err
	}
	key := stakeInfoKey(stakeAddress, dec)
	store, err := get(s, key)
	if err == ErrEof {
		return stakeInfo, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetStakeInfo, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, stakeInfo); err != nil {
		return nil, fmt.Errorf("GetStakeInfo, deserialize stakeInfo error: %v", err)
	}
	return stakeInfo, nil
}

func addUnlockingInfo(s *native.NativeContract, stakeAddress common.Address, unlockingStake *UnlockingStake) error {
	unlockingInfo, err := GetUnlockingInfo(s, stakeAddress)
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

func filterExpiredUnlockingInfo(s *native.NativeContract, stakeAddress common.Address) (*big.Int, error) {
	height := s.ContractRef().BlockHeight()
	unlockingInfo, err := GetUnlockingInfo(s, stakeAddress)
	if err != nil {
		return nil, fmt.Errorf("filterExpiredUnlockingInfo, GetUnlockingInfo error: %v", err)
	}
	j := 0
	var expiredSum *big.Int
	for _, unlockingStake := range unlockingInfo.UnlockingStake {
		if unlockingStake.CompleteHeight.Cmp(height) == 1 {
			unlockingInfo.UnlockingStake[j] = unlockingStake
			j++
		} else {
			expiredSum = new(big.Int).Add(expiredSum, unlockingStake.Amount)
		}
	}
	unlockingInfo.UnlockingStake = unlockingInfo.UnlockingStake[:j]
	if len(unlockingInfo.UnlockingStake) == 0 {
		delUnlockingInfo(s, stakeAddress)
	} else {
		err = setUnlockingInfo(s, unlockingInfo)
		if err != nil {
			return nil, fmt.Errorf("filterExpiredUnlockingInfo, setUnlockingInfo error: %v", err)
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

func GetUnlockingInfo(s *native.NativeContract, stakeAddress common.Address) (*UnlockingInfo, error) {
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

func GetCurrentEpoch(s *native.NativeContract) (*big.Int, error) {
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

func GetCurrentEpochInfo(s *native.NativeContract) (*EpochInfo, error) {
	ID, err := GetCurrentEpoch(s)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpochInfo, GetCurrentEpochInfo error: %v", err)
	}
	epochInfo, err := GetEpochInfo(s, ID)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpochInfo, GetEpochInfo error: %v", err)
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

func GetEpochInfo(s *native.NativeContract, ID *big.Int) (*EpochInfo, error) {
	epochInfo := &EpochInfo{
		Validators: make([]*Validator, 0),
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

func lockPoolKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_LOCK_POOL))
}

func unlockPoolKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_UNLOCK_POOL))
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
