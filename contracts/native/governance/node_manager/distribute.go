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

const (
	TokenPrecision   = 18
	PercentPrecision = 4
)

var (
	TokenDecimal   = new(big.Int).Exp(big.NewInt(10), big.NewInt(TokenPrecision), nil)
	PercentDecimal = new(big.Int).Exp(big.NewInt(10), big.NewInt(PercentPrecision), nil)
)

// IncreaseValidatorPeriod return the period just ended
func IncreaseValidatorPeriod(s *native.NativeContract, validator *Validator) (uint64, error) {
	// fetch current rewards
	validatorAccumulatedRewards, err := getValidatorAccumulatedRewards(s, validator.ConsensusAddress)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, getValidatorAccumulatedRewards error: %v", err)
	}

	// calculate current ratio
	// mul decimal
	ratio, err := validatorAccumulatedRewards.Rewards.DivWithTokenDecimal(validator.TotalStake)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, validatorAccumulatedRewards.Rewards.DivWithTokenDecimal error: %v", err)
	}

	// fetch snapshot rewards for last period
	validatorSnapshotRewards, err := getValidatorSnapshotRewards(s, validator.ConsensusAddress, validatorAccumulatedRewards.Period-1)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, getValidatorSnapshotRewards error: %v", err)
	}

	// decrement reference count
	err = decreaseReferenceCount(s, validator.ConsensusAddress, validatorAccumulatedRewards.Period-1)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, decreaseReferenceCount error: %v", err)
	}

	// set new snapshot rewards with reference count of 1
	newRatio, err := validatorSnapshotRewards.AccumulatedRewardsRatio.Add(ratio)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, validatorSnapshotRewards.AccumulatedRewardsRatio.Add error: %v", err)
	}
	newValidatorSnapshotRewards := &ValidatorSnapshotRewards{
		AccumulatedRewardsRatio: newRatio,
		ReferenceCount:          1,
	}
	err = setValidatorSnapshotRewards(s, validator.ConsensusAddress, validatorAccumulatedRewards.Period, newValidatorSnapshotRewards)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, setValidatorSnapshotRewards error: %v", err)
	}

	// set accumulate rewards, incrementing period by 1
	newValidatorAccumulatedRewards := &ValidatorAccumulatedRewards{
		Rewards: utils.NewDecFromBigInt(new(big.Int)),
		Period:  validatorAccumulatedRewards.Period + 1,
	}
	err = setValidatorAccumulatedRewards(s, validator.ConsensusAddress, newValidatorAccumulatedRewards)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, setValidatorAccumulatedRewards error: %v", err)
	}

	return validatorAccumulatedRewards.Period, nil
}

func withdrawStakeRewards(s *native.NativeContract, validator *Validator, stakeInfo *StakeInfo) (utils.Dec, error) {
	// end current period and calculate rewards
	endingPeriod, err := IncreaseValidatorPeriod(s, validator)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawStakeRewards, IncreaseValidatorPeriod error: %v", err)
	}
	rewards, err := CalculateStakeRewards(s, stakeInfo.StakeAddress, validator.ConsensusAddress, endingPeriod)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawStakeRewards, CalculateStakeRewards error: %v", err)
	}

	err = contract.NativeTransfer(s.StateDB(), this, stakeInfo.StakeAddress, rewards.BigInt())
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawStakeRewards, nativeTransfer error: %v", err)
	}

	// update the outstanding rewards
	outstanding, err := getOutstandingRewards(s)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawStakeRewards, getOutstandingRewards error: %v", err)
	}
	validatorOutstanding, err := getValidatorOutstandingRewards(s, validator.ConsensusAddress)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawStakeRewards, getValidatorOutstandingRewards error: %v", err)
	}
	newOutstandingRewards, err := outstanding.Rewards.Sub(rewards)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawStakeRewards, outstanding.Rewards.Sub error: %v", err)
	}
	err = setOutstandingRewards(s, &OutstandingRewards{Rewards: newOutstandingRewards})
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawStakeRewards, setOutstandingRewards error: %v", err)
	}
	newValidatorOutstandingRewards, err := validatorOutstanding.Rewards.Sub(rewards)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawStakeRewards, validatorOutstanding.Rewards.Sub error: %v", err)
	}
	err = setValidatorOutstandingRewards(s, validator.ConsensusAddress, &ValidatorOutstandingRewards{Rewards: newValidatorOutstandingRewards})
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawStakeRewards, setValidatorOutstandingRewards error: %v", err)
	}

	// decrement reference count of starting period
	startingInfo, err := getStakeStartingInfo(s, stakeInfo.StakeAddress, validator.ConsensusAddress)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawStakeRewards, getStakeStartingInfo error: %v", err)
	}
	startPeriod := startingInfo.StartPeriod
	err = decreaseReferenceCount(s, validator.ConsensusAddress, startPeriod)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawStakeRewards, decreaseReferenceCount error: %v", err)
	}

	// remove stake starting info
	delStakeStartingInfo(s, stakeInfo.StakeAddress, validator.ConsensusAddress)
	return rewards, nil
}

func CalculateStakeRewards(s *native.NativeContract, stakeAddress common.Address, consensusAddr common.Address, endPeriod uint64) (utils.Dec, error) {
	// fetch starting info for delegation
	startingInfo, err := getStakeStartingInfo(s, stakeAddress, consensusAddr)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("CalculateStakeRewards, getStakeStartingInfo error: %v", err)
	}

	startPeriod := startingInfo.StartPeriod
	stake := startingInfo.Stake

	// sanity check
	if startPeriod > endPeriod {
		panic("startPeriod cannot be greater than endPeriod")
	}

	// return staking * (ending - starting)
	starting, err := getValidatorSnapshotRewards(s, consensusAddr, startPeriod)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("CalculateStakeRewards, getValidatorSnapshotRewards start error: %v", err)
	}
	ending, err := getValidatorSnapshotRewards(s, consensusAddr, endPeriod)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("CalculateStakeRewards, getValidatorSnapshotRewards end error: %v", err)
	}
	difference, err := ending.AccumulatedRewardsRatio.Sub(starting.AccumulatedRewardsRatio)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("CalculateStakeRewards error: %v", err)
	}
	rewards, err := difference.MulWithTokenDecimal(stake)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("CalculateStakeRewards error: %v", err)
	}
	return rewards, nil
}

func initializeStake(s *native.NativeContract, stakeInfo *StakeInfo, consensusAddr common.Address) error {
	// period has already been incremented
	validatorAccumulatedRewards, err := getValidatorAccumulatedRewards(s, consensusAddr)
	if err != nil {
		return fmt.Errorf("initializeStake, getValidatorAccumulatedRewards start error: %v", err)
	}
	previousPeriod := validatorAccumulatedRewards.Period - 1

	// increment reference count for the period we're going to track
	err = increaseReferenceCount(s, consensusAddr, previousPeriod)
	if err != nil {
		return fmt.Errorf("initializeStake, increaseReferenceCount start error: %v", err)
	}

	stake := stakeInfo.Amount
	err = setStakeStartingInfo(s, stakeInfo.StakeAddress, consensusAddr, &StakeStartingInfo{previousPeriod,
		stake, s.ContractRef().BlockHeight()})
	if err != nil {
		return fmt.Errorf("initializeStake, setStakeStartingInfo error: %v", err)
	}
	return nil
}

func withdrawCommission(s *native.NativeContract, stakeAddress common.Address, consensusAddr common.Address) (utils.Dec, error) {
	accumulatedCommission, err := getAccumulatedCommission(s, consensusAddr)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawCommission, getAccumulatedCommission error: %v", err)
	}

	// update the outstanding rewards
	outstanding, err := getOutstandingRewards(s)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawCommission, getOutstandingRewards error: %v", err)
	}
	validatorOutstanding, err := getValidatorOutstandingRewards(s, consensusAddr)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawCommission, getValidatorOutstandingRewards error: %v", err)
	}
	newOutstandingRewards, err := outstanding.Rewards.Sub(accumulatedCommission.Amount)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawCommission, outstanding.Rewards.Sub error: %v", err)
	}
	err = setOutstandingRewards(s, &OutstandingRewards{Rewards: newOutstandingRewards})
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawCommission, setOutstandingRewards error: %v", err)
	}
	newValidatorOutstandingRewards, err := validatorOutstanding.Rewards.Sub(accumulatedCommission.Amount)
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawCommission, validatorOutstanding.Rewards.Sub error: %v", err)
	}
	err = setValidatorOutstandingRewards(s, consensusAddr, &ValidatorOutstandingRewards{Rewards: newValidatorOutstandingRewards})
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawCommission, setValidatorOutstandingRewards error: %v", err)
	}

	err = contract.NativeTransfer(s.StateDB(), this, stakeAddress, accumulatedCommission.Amount.BigInt())
	if err != nil {
		return utils.Dec{}, fmt.Errorf("withdrawCommission, nativeTransfer commission error: %v", err)
	}
	return accumulatedCommission.Amount, nil
}

func allocateRewardsToValidator(s *native.NativeContract, validator *Validator, rewards utils.Dec) error {
	commission, err := validator.Commission.Rate.MulWithPercentDecimal(rewards)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, validator.Commission.Rate.Mul error: %v", err)
	}
	// stake reward = reward-commission
	stakeRewards, err := rewards.Sub(commission)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, rewards.Sub error: %v", err)
	}

	// update accumulate commission
	currentCommission, err := getAccumulatedCommission(s, validator.ConsensusAddress)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, getAccumulatedCommission error: %v", err)
	}
	currentCommission.Amount, err = currentCommission.Amount.Add(commission)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, currentCommission.Amount.Add error: %v", err)
	}
	err = setAccumulatedCommission(s, validator.ConsensusAddress, currentCommission)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, setAccumulatedCommission error: %v", err)
	}

	// update accumulate rewards
	validatorAccumulatedRewards, err := getValidatorAccumulatedRewards(s, validator.ConsensusAddress)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, getValidatorAccumulatedRewards error: %v", err)
	}
	validatorAccumulatedRewards.Rewards, err = validatorAccumulatedRewards.Rewards.Add(stakeRewards)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, validatorAccumulatedRewards.Rewards.Add error: %v", err)
	}
	err = setValidatorAccumulatedRewards(s, validator.ConsensusAddress, validatorAccumulatedRewards)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, setValidatorAccumulatedRewards error: %v", err)
	}

	// update validator outstanding
	outstanding, err := getValidatorOutstandingRewards(s, validator.ConsensusAddress)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, getValidatorOutstandingRewards error: %v", err)
	}
	outstanding.Rewards, err = outstanding.Rewards.Add(rewards)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, outstanding.Rewards.Add error: %v", err)
	}
	err = setValidatorOutstandingRewards(s, validator.ConsensusAddress, outstanding)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, setValidatorOutstandingRewards error: %v", err)
	}
	return nil
}
