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
	"math/big"
)

func deposit(s *native.NativeContract, from common.Address, amount *big.Int, validator *Validator) error {
	// store deposit info
	err := depositStakeInfo(s, from, validator.ConsensusPubkey, amount)
	if err != nil {
		return fmt.Errorf("deposit, depositStakeInfo error: %v", err)
	}

	// transfer native token
	err = nativeTransfer(s, from, this, amount)
	if err != nil {
		return fmt.Errorf("deposit, nativeTransfer error: %v", err)
	}

	// update lock and unlock token pool
	switch {
	case validator.IsLocked():
		err := depositLockPool(s, amount)
		if err != nil {
			return fmt.Errorf("deposit, depositLockPool error: %v", err)
		}
	case validator.IsUnlocking(), validator.IsUnlocked():
		err := depositUnlockPool(s, amount)
		if err != nil {
			return fmt.Errorf("deposit, depositUnlockPool error: %v", err)
		}
	default:
		return fmt.Errorf("invalid status")
	}

	return nil
}
