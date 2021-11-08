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

package lock_proxy

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
)

var ErrUnSupported = errors.New("erc20 asset cross chain unsupported yet")

func getBalanceFor(s *native.NativeContract, fromAsset common.Address) (*big.Int, error) {
	if fromAsset == common.EmptyAddress {
		return s.StateDB().GetBalance(this), nil
	} else {
		return erc20Balance(s, fromAsset, this)
	}
}

func onlySupportNativeToken(fromAsset common.Address) bool {
	if fromAsset == common.EmptyAddress {
		return true
	}
	return false
}

// todo: get erc20 balance
func erc20Balance(s *native.NativeContract, asset, user common.Address) (*big.Int, error) {
	return nil, ErrUnSupported
}
