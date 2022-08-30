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

package maas_config

import (
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
)

// storage key prefix
const (
	BLACKLIST         = "blacklist"
	OWNER             = "owner"
	NODE_WHITE_ENABLE = "node_white_enable"
	GAS_MANAGE_ENABLE = "gas_manage_enable"
	NODE_WHITELIST    = "node_whitelist"
	GAS_MANAGER_LIST  = "gas_manager_list"
	GAS_USER_LIST     = "gas_user_list"
	GAS_ADMIN_LIST    = "gas_admin_list"
)

var (
	blacklistKey       = utils.ConcatKey(this, []byte(BLACKLIST))
	ownerKey           = utils.ConcatKey(this, []byte(OWNER))
	gasManageEnableKey = utils.ConcatKey(this, []byte(GAS_MANAGE_ENABLE))
	gasManagerListKey  = utils.ConcatKey(this, []byte(GAS_MANAGER_LIST))
	gasUserListKey     = utils.ConcatKey(this, []byte(GAS_USER_LIST))
	gasAdminListKey    = utils.ConcatKey(this, []byte(GAS_ADMIN_LIST))
)

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
		// } else if value == nil || len(value) == 0 {
		// 	return nil, ErrEof
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
