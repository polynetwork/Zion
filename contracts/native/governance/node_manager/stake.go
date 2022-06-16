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
	height := s.ContractRef().BlockHeight()
	// get deposit info
	stakeInfo, found, err := GetStakeInfo(s, from, validator.ConsensusPubkey)
	if err != nil {
		return fmt.Errorf("deposit, GetStakeInfo error: %v", err)
	}
	// call the appropriate hook if present
	if found {
		err = BeforeStakeModified(s, from, validator)
	} else {
		err = BeforeStakeCreated(s, validator)
	}
	// update stake info
	stakeInfo.Amount = new(big.Int).Add(stakeInfo.Amount, amount)
	err = setStakeInfo(s, stakeInfo)
	if err != nil {
		return fmt.Errorf("deposit, setStakeInfo error: %v", err)
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
	case validator.IsUnlocking(height), validator.IsUnlocked(height):
		err := depositUnlockPool(s, amount)
		if err != nil {
			return fmt.Errorf("deposit, depositUnlockPool error: %v", err)
		}
	default:
		return fmt.Errorf("invalid status")
	}

	return nil
}

func unStake(s *native.NativeContract, from common.Address, amount *big.Int, validator *Validator) error {
	height := s.ContractRef().BlockHeight()
	globalConfig, err := GetGlobalConfig(s)
	if err != nil {
		return fmt.Errorf("unStake, GetGlobalConfig error: %v", err)
	}

	// store deposit info
	err = withdrawStakeInfo(s, from, validator.ConsensusPubkey, amount)
	if err != nil {
		return fmt.Errorf("unStake, depositStakeInfo error: %v", err)
	}

	// update lock and unlock token pool
	if validator.IsLocked() {
		err = withdrawLockPool(s, amount)
		if err != nil {
			return fmt.Errorf("unStake, withdrawLockPool error: %v", err)
		}
		err = depositUnlockPool(s, amount)
		if err != nil {
			return fmt.Errorf("unStake, depositUnlockPool error: %v", err)
		}
		unlockingStake := &UnlockingStake{
			Height:         height,
			CompleteHeight: new(big.Int).Add(height, globalConfig.BlockPerEpoch),
			Amount:         amount,
		}
		err = addUnlockingInfo(s, from, unlockingStake)
		if err != nil {
			return fmt.Errorf("unStake, addUnlockingInfo error: %v", err)
		}
	}

	if validator.IsUnlocked(height) || validator.IsRemoved(height) {
		err = withdrawUnlockPool(s, amount)
		if err != nil {
			return fmt.Errorf("unStake, withdrawUnlockPool error: %v", err)
		}
		// transfer native token
		err = nativeTransfer(s, this, from, amount)
		if err != nil {
			return fmt.Errorf("unStake, nativeTransfer error: %v", err)
		}
	}

	if validator.IsUnlocking(height) || validator.IsRemoving(height) {
		unlockingStake := &UnlockingStake{
			Height:         height,
			CompleteHeight: validator.UnlockHeight,
			Amount:         amount,
		}
		err = addUnlockingInfo(s, from, unlockingStake)
		if err != nil {
			return fmt.Errorf("unStake, addUnlockingInfo error: %v", err)
		}
	}

	return nil
}
