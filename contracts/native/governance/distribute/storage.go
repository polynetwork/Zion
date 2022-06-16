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

package distribute

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/rlp"
)

var ErrEof = errors.New("EOF")

// storage key prefix
const (
	SKP_ACCUMULATED_COMMISSION        = "st_accumulated_commission"
	SKP_VALIDATOR_ACCUMULATED_REWARDS = "st_validator_accumulated_rewards"
	SKP_VALIDATOR_OUTSTANDING_REWARDS = "st_validator_outstanding_rewards"
	SKP_VALIDATOR_SNAPSHOT_REWARDS    = "st_validator_snapshot_rewards"
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

func GetAccumulatedCommission(s *native.NativeContract, dec []byte) (*AccumulatedCommission, error) {
	accumulatedCommission := &AccumulatedCommission{}
	key := accumulatedCommissionKey(dec)
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("GetAccumulatedCommission, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, accumulatedCommission); err != nil {
		return nil, fmt.Errorf("GetAccumulatedCommission, deserialize accumulatedCommission error: %v", err)
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

func GetValidatorAccumulatedRewards(s *native.NativeContract, dec []byte) (*ValidatorAccumulatedRewards, error) {
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

func GetValidatorOutstandingRewards(s *native.NativeContract, dec []byte) (*ValidatorOutstandingRewards, error) {
	validatorOutstandingRewards := &ValidatorOutstandingRewards{}
	key := validatorOutstandingRewardsKey(dec)
	store, err := get(s, key)
	if err != nil {
		return nil, fmt.Errorf("GetValidatorOutstandingRewards, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, validatorOutstandingRewards); err != nil {
		return nil, fmt.Errorf("GetValidatorOutstandingRewards, deserialize validatorOutstandingRewards error: %v", err)
	}
	return validatorOutstandingRewards, nil
}

func delValidatorOutstandingRewards(s *native.NativeContract, dec []byte) {
	key := validatorOutstandingRewardsKey(dec)
	del(s, key)
}

func increaseReferenceCount(s *native.NativeContract, dec []byte, period uint64) error {
	validatorSnapshotRewards, err := GetValidatorSnapshotRewards(s, dec, period)
	if err != nil {
		return fmt.Errorf("increaseReferenceCount, GetValidatorSnapshotRewards error: %v", err)
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
	validatorSnapshotRewards, err := GetValidatorSnapshotRewards(s, dec, period)
	if err != nil {
		return fmt.Errorf("decreaseReferenceCount, GetValidatorSnapshotRewards error: %v", err)
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

func GetValidatorSnapshotRewards(s *native.NativeContract, dec []byte, period uint64) (*ValidatorSnapshotRewards, error) {
	validatorSnapshotRewards := &ValidatorSnapshotRewards{}
	key := validatorSnapshotRewardsKey(dec, period)
	store, err := get(s, key)
	if err == ErrEof {
		return validatorSnapshotRewards, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetValidatorSnapshotRewards, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, validatorSnapshotRewards); err != nil {
		return nil, fmt.Errorf("GetValidatorSnapshotRewards, deserialize validatorSnapshotRewards error: %v", err)
	}
	return validatorSnapshotRewards, nil
}

func delValidatorSnapshotRewards(s *native.NativeContract, dec []byte, period uint64) {
	key := validatorSnapshotRewardsKey(dec, period)
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

func accumulatedCommissionKey(dec []byte) []byte {
	return utils.ConcatKey(this, []byte(SKP_ACCUMULATED_COMMISSION), dec)
}

func validatorAccumulatedRewardsKey(dec []byte) []byte {
	return utils.ConcatKey(this, []byte(SKP_VALIDATOR_ACCUMULATED_REWARDS), dec)
}

func validatorOutstandingRewardsKey(dec []byte) []byte {
	return utils.ConcatKey(this, []byte(SKP_VALIDATOR_OUTSTANDING_REWARDS), dec)
}

func validatorSnapshotRewardsKey(dec []byte, period uint64) []byte {
	return utils.ConcatKey(this, []byte(SKP_VALIDATOR_SNAPSHOT_REWARDS), dec, utils.Uint64Bytes(period))
}
