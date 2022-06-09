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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"

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
