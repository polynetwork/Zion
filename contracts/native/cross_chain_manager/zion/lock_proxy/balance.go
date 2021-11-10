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
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/core"
)

var ErrUnSupported = errors.New("erc20 asset cross chain unsupported yet")

func getBalanceFor(s *native.NativeContract, fromAsset, owner common.Address) (*big.Int, error) {
	if fromAsset == common.EmptyAddress {
		return s.StateDB().GetBalance(owner), nil
	} else {
		return erc20Balance(s, fromAsset, owner)
	}
}

func transfer2Contract(s *native.NativeContract, fromAsset, txOrigin common.Address, amount *big.Int) error {
	if fromAsset == common.EmptyAddress {
		if amount.Cmp(common.Big0) == 0 {
			return fmt.Errorf("transferred native token cannot be zero!")
		}
		if s.ContractRef().Value() == nil || s.ContractRef().Value().Cmp(amount) != 0 {
			return fmt.Errorf("transferred native token is not equal to amount!")
		}
	} else {
		if s.ContractRef().Value() != nil && s.ContractRef().Value().Cmp(common.Big0) != 0 {
			return fmt.Errorf("there should be no native token transfer!")
		}
		return erc20Transfer(s, fromAsset, txOrigin, this, amount)
	}
	return nil
}

func transferFromContract(s *native.NativeContract, toAsset, toAddress common.Address, amount *big.Int) error {
	if toAsset == common.EmptyAddress {
		return nativeTransfer(s, this, toAddress, amount)
	} else {
		return erc20TransferFrom(s, toAsset, this, toAddress, amount)
	}
}

func nativeTransfer(s *native.NativeContract, from, to common.Address, amount *big.Int) error {
	if !core.CanTransfer(s.StateDB(), from, amount) {
		return fmt.Errorf("Insufficient balance")
	}
	core.Transfer(s.StateDB(), from, to, amount)
	return nil
}

func onlySupportNativeToken(fromAsset common.Address) bool {
	if fromAsset == common.EmptyAddress {
		return true
	}
	return false
}

// todo: implement erc20 interfaces
func erc20Balance(s *native.NativeContract, asset, owner common.Address) (*big.Int, error) {
	return nil, ErrUnSupported
}

func erc20Transfer(s *native.NativeContract, asset, from, to common.Address, amount *big.Int) error {
	return ErrUnSupported
}

func erc20TransferFrom(s *native.NativeContract, asset, from, to common.Address, amount *big.Int) error {
	return ErrUnSupported
}
