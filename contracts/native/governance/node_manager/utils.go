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
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/core"
	"math/big"
)

func nativeTransfer(s *native.NativeContract, from, to common.Address, amount *big.Int) error {
	if amount.Sign() <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if !core.CanTransfer(s.StateDB(), from, amount) {
		return fmt.Errorf("%s insufficient balance", from.Hex())
	}
	core.Transfer(s.StateDB(), from, to, amount)
	return nil
}
