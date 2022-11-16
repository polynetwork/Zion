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
	"github.com/ethereum/go-ethereum/contracts/native/contract"
	"math/big"
)

func AfterValidatorCreated(s *native.NativeContract, validator *Validator) error {
	// set initial historical rewards (period 0) with reference count of 1
	err := setValidatorSnapshotRewards(s, validator.ConsensusAddress, 0, &ValidatorSnapshotRewards{NewDecFromBigInt(new(big.Int)), 1})
	if err != nil {
		return fmt.Errorf("AfterValidatorCreated, setValidatorSnapshotRewards error: %v", err)
	}

	// set accumulate rewards (starting at period 1)
	err = setValidatorAccumulatedRewards(s, validator.ConsensusAddress, &ValidatorAccumulatedRewards{NewDecFromBigInt(new(big.Int)), 1})
	if err != nil {
		return fmt.Errorf("AfterValidatorCreated, setValidatorAccumulatedRewards error: %v", err)
	}

	// set accumulated commission
	err = setAccumulatedCommission(s, validator.ConsensusAddress, &AccumulatedCommission{NewDecFromBigInt(new(big.Int))})
	if err != nil {
		return fmt.Errorf("AfterValidatorCreated, setAccumulatedCommission error: %v", err)
	}

	// set outstanding rewards
	err = setValidatorOutstandingRewards(s, validator.ConsensusAddress, &ValidatorOutstandingRewards{Rewards: NewDecFromBigInt(new(big.Int))})
	if err != nil {
		return fmt.Errorf("AfterValidatorCreated, setValidatorOutstandingRewards error: %v", err)
	}
	return nil
}

func AfterValidatorRemoved(s *native.NativeContract, validator *Validator) error {
	// fetch outstanding
	outstanding, err := getValidatorOutstandingRewards(s, validator.ConsensusAddress)
	if err != nil {
		return fmt.Errorf("AfterValidatorRemoved, getValidatorOutstandingRewards error: %v", err)
	}
	communityInfo, err := GetCommunityInfoImpl(s)
	if err != nil {
		return fmt.Errorf("AfterValidatorRemoved, GetCommunityInfoImpl error: %v", err)
	}
	// transfer outstanding dust to community pool
	err = contract.NativeTransfer(s.StateDB(), this, communityInfo.CommunityAddress, outstanding.Rewards.BigInt())
	if err != nil {
		return fmt.Errorf("AfterValidatorRemoved, nativeTransfer error: %v", err)
	}

	//delete signer and proposal address
	delSignerAddr(s, validator.SignerAddress)
	delProposalAddr(s, validator.ProposalAddress)

	// delete outstanding
	delValidatorOutstandingRewards(s, validator.ConsensusAddress)

	// remove commission record
	delAccumulatedCommission(s, validator.ConsensusAddress)

	validatorAccumulatedRewards, err := getValidatorAccumulatedRewards(s, validator.ConsensusAddress)
	if err != nil {
		return fmt.Errorf("AfterValidatorRemoved, getValidatorAccumulatedRewards error: %v", err)
	}

	// clear accumulate rewards
	delValidatorAccumulatedRewards(s, validator.ConsensusAddress)

	// clear snapshot rewards
	delValidatorSnapshotRewards(s, validator.ConsensusAddress, validatorAccumulatedRewards.Period-1)
	return nil
}

func BeforeStakeCreated(s *native.NativeContract, validator *Validator) error {
	_, err := IncreaseValidatorPeriod(s, validator)
	return err
}

func BeforeStakeModified(s *native.NativeContract, validator *Validator, stakeInfo *StakeInfo) error {
	if _, err := withdrawStakeRewards(s, validator, stakeInfo); err != nil {
		return err
	}
	return nil
}

func AfterStakeModified(s *native.NativeContract, stakeInfo *StakeInfo, consensusAddr common.Address) error {
	err := initializeStake(s, stakeInfo, consensusAddr)
	if err != nil {
		return fmt.Errorf("AfterStakeModified, initializeStake error: %v", err)
	}
	return nil
}
