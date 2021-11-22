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
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/core"
)

func getBalanceFor(s *native.NativeContract, owner common.Address) *big.Int {
	return s.StateDB().GetBalance(owner)
}

func transfer2Contract(s *native.NativeContract, amount *big.Int) error {
	if amount.Cmp(common.Big0) == 0 {
		return fmt.Errorf("transferred native token cannot be zero!")
	}
	if s.ContractRef().Value() == nil || s.ContractRef().Value().Cmp(amount) != 0 {
		return fmt.Errorf("transferred native token is not equal to amount!")
	}
	return nil
}

func transferFromContract(s *native.NativeContract, toAddress common.Address, amount *big.Int) error {
	return nativeTransfer(s, this, toAddress, amount)
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
