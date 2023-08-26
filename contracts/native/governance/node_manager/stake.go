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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/contract"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

func deposit(s *native.NativeContract, from common.Address, amount utils.Dec, validator *Validator) error {
	// get deposit info
	stakeInfo, found, err := getStakeInfo(s, from, validator.ConsensusAddress)
	if err != nil {
		return fmt.Errorf("deposit, getStakeInfo error: %v", err)
	}
	// call the appropriate hook if present
	if found {
		err = BeforeStakeModified(s, validator, stakeInfo)
		if err != nil {
			return fmt.Errorf("deposit, BeforeStakeModified error: %v", err)
		}
	} else {
		err = BeforeStakeCreated(s, validator)
		if err != nil {
			return fmt.Errorf("deposit, BeforeStakeCreated error: %v", err)
		}
	}
	// update stake info
	stakeInfo.Amount, err = stakeInfo.Amount.Add(amount)
	if err != nil {
		return fmt.Errorf("deposit, stakeInfo.Amount.Add error: %v", err)
	}
	err = setStakeInfo(s, stakeInfo)
	if err != nil {
		return fmt.Errorf("deposit, setStakeInfo error: %v", err)
	}

	// do not transfer native token, already transfered by value

	// update total token pool
	err = depositTotalPool(s, amount)
	if err != nil {
		return fmt.Errorf("deposit, depositTotalPool error: %v", err)
	}

	// Call the after-stake hook
	if err = AfterStakeModified(s, stakeInfo, validator.ConsensusAddress); err != nil {
		return err
	}

	return nil
}

func unStake(s *native.NativeContract, from common.Address, amount utils.Dec, validator *Validator) error {
	height := s.ContractRef().BlockHeight()
	globalConfig, err := GetGlobalConfigImpl(s)
	if err != nil {
		return fmt.Errorf("unStake, GetGlobalConfig error: %v", err)
	}

	stakeInfo, found, err := getStakeInfo(s, from, validator.ConsensusAddress)
	if err != nil {
		return fmt.Errorf("unStake, get stake info error: %v", err)
	}
	if !found {
		return fmt.Errorf("unStake, stake info nit exist")
	}

	err = BeforeStakeModified(s, validator, stakeInfo)
	if err != nil {
		return fmt.Errorf("unStake, BeforeStakeModified error: %v", err)
	}

	// update lock and unlock token pool
	if validator.IsLocked() {
		unlockingStake := &UnlockingStake{
			Height:           height,
			CompleteHeight:   new(big.Int).Add(height, globalConfig.BlockPerEpoch),
			ConsensusAddress: validator.ConsensusAddress,
			Amount:           amount,
		}
		err = addUnlockingInfo(s, from, unlockingStake)
		if err != nil {
			return fmt.Errorf("unStake, addUnlockingInfo error: %v", err)
		}
	}

	if validator.IsUnlocked(height) || validator.IsRemoved(height) {
		err = withdrawTotalPool(s, amount)
		if err != nil {
			return fmt.Errorf("unStake, withdrawTotalPool error: %v", err)
		}
		// transfer native token
		err = contract.NativeTransfer(s.StateDB(), this, from, amount.BigInt())
		if err != nil {
			return fmt.Errorf("unStake, nativeTransfer error: %v", err)
		}
	}

	if validator.IsUnlocking(height) || validator.IsRemoving(height) {
		unlockingStake := &UnlockingStake{
			Height:           height,
			CompleteHeight:   validator.UnlockHeight,
			ConsensusAddress: validator.ConsensusAddress,
			Amount:           amount,
		}
		err = addUnlockingInfo(s, from, unlockingStake)
		if err != nil {
			return fmt.Errorf("unStake, addUnlockingInfo error: %v", err)
		}
	}

	stakeInfo.Amount, err = stakeInfo.Amount.Sub(amount)
	if err != nil {
		return fmt.Errorf("unStake, stakeInfo.Amount.Sub error: %v", err)
	}
	if stakeInfo.Amount.IsZero() {
		delStakeInfo(s, from, validator.ConsensusAddress)
	} else {
		err = setStakeInfo(s, stakeInfo)
		if err != nil {
			return fmt.Errorf("unStake, set stake info error: %v", err)
		}
		// Call the after-stake hook
		if err = AfterStakeModified(s, stakeInfo, validator.ConsensusAddress); err != nil {
			return err
		}
	}

	return nil
}
