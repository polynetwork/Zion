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

package auth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const (
	SKP_AUTH = "st_auth"
)

func getAllowance(s *native.NativeContract, owner, spender common.Address) *big.Int {
	key := allowanceKey(owner, spender)
	blob, _ := s.GetCacheDB().Get(key)
	if blob == nil {
		return common.Big0
	}
	return new(big.Int).SetBytes(blob)
}

func setAllowance(s *native.NativeContract, owner, spender common.Address, amount *big.Int) {
	key := allowanceKey(owner, spender)
	s.GetCacheDB().Put(key, amount.Bytes())
}

// ====================================================================
//
// storage keys
//
// ====================================================================
func allowanceKey(owner, spender common.Address) []byte {
	return utils.ConcatKey(this, []byte(SKP_AUTH), owner[:], spender[:])
}
