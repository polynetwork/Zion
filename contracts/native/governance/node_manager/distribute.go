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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"math/big"
)

var RatioDecimal = new(big.Int).SetUint64(1000000)

// IncreaseValidatorPeriod return the period just ended
func IncreaseValidatorPeriod(s *native.NativeContract, validator *Validator) (uint64, error) {
	dec, err := hexutil.Decode(validator.ConsensusPubkey)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, decode pubkey error: %v", err)
	}

	// fetch current rewards
	validatorAccumulatedRewards, err := GetValidatorAccumulatedRewards(s, dec)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, GetValidatorAccumulatedRewards error: %v", err)
	}

	// calculate current ratio
	// mul decimal
	rewardsD := new(big.Int).Mul(validatorAccumulatedRewards.Rewards, RatioDecimal)
	ratio := new(big.Int).Div(rewardsD, validator.TotalStake)

	// fetch snapshot rewards for last period
	validatorSnapshotRewards, err := GetValidatorSnapshotRewards(s, dec, validatorAccumulatedRewards.Period-1)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, GetValidatorSnapshotRewards error: %v", err)
	}

	// decrement reference count
	err = decreaseReferenceCount(s, dec, validatorAccumulatedRewards.Period-1)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, decreaseReferenceCount error: %v", err)
	}

	// set new snapshot rewards with reference count of 1
	newValidatorSnapshotRewards := &ValidatorSnapshotRewards{
		AccumulatedRewardsRatio: new(big.Int).Add(validatorSnapshotRewards.AccumulatedRewardsRatio, ratio),
		ReferenceCount:          1,
	}
	err = setValidatorSnapshotRewards(s, dec, validatorAccumulatedRewards.Period, newValidatorSnapshotRewards)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, setValidatorSnapshotRewards error: %v", err)
	}

	// set accumulate rewards, incrementing period by 1
	newValidatorAccumulatedRewards := &ValidatorAccumulatedRewards{
		Rewards: new(big.Int),
		Period:  validatorAccumulatedRewards.Period + 1,
	}
	err = setValidatorAccumulatedRewards(s, dec, newValidatorAccumulatedRewards)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, setValidatorAccumulatedRewards error: %v", err)
	}

	return validatorAccumulatedRewards.Period, nil
}

func withdrawStakeRewards(s *native.NativeContract, validator *Validator, stakeInfo *StakeInfo) (*big.Int, error) {
	dec, err := hexutil.Decode(validator.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("withdrawStakeRewards, decode pubkey error: %v", err)
	}

	// end current period and calculate rewards
	endingPeriod, err := IncreaseValidatorPeriod(s, validator)
	if err != nil {
		return nil, fmt.Errorf("withdrawStakeRewards, IncreaseValidatorPeriod error: %v", err)
	}
	rewards, err := CalculateStakeRewards(s, stakeInfo.StakeAddress, dec, endingPeriod)
	if err != nil {
		return nil, fmt.Errorf("withdrawStakeRewards, CalculateStakeRewards error: %v", err)
	}
	
	err = nativeTransfer(s, this, stakeInfo.StakeAddress, rewards)
	if err != nil {
		return nil, fmt.Errorf("withdrawStakeRewards, nativeTransfer error: %v", err)
	}

	// update the outstanding rewards
	outstanding, err := GetOutstandingRewards(s)
	if err != nil {
		return nil, fmt.Errorf("withdrawStakeRewards, GetOutstandingRewards error: %v", err)
	}
	validatorOutstanding, err := GetValidatorOutstandingRewards(s, dec)
	if err != nil {
		return nil, fmt.Errorf("withdrawStakeRewards, GetValidatorOutstandingRewards error: %v", err)
	}
	err = setOutstandingRewards(s, &OutstandingRewards{Rewards: new(big.Int).Sub(outstanding.Rewards, rewards)})
	if err != nil {
		return nil, fmt.Errorf("withdrawStakeRewards, setOutstandingRewards error: %v", err)
	}
	err = setValidatorOutstandingRewards(s, dec, &ValidatorOutstandingRewards{Rewards: new(big.Int).Sub(validatorOutstanding.Rewards, rewards)})
	if err != nil {
		return nil, fmt.Errorf("withdrawStakeRewards, setValidatorOutstandingRewards error: %v", err)
	}

	// decrement reference count of starting period
	startingInfo, err := GetStakeStartingInfo(s, stakeInfo.StakeAddress, dec)
	if err != nil {
		return nil, fmt.Errorf("withdrawStakeRewards, GetStakeStartingInfo error: %v", err)
	}
	startPeriod := startingInfo.StartPeriod
	err = decreaseReferenceCount(s, dec, startPeriod)
	if err != nil {
		return nil, fmt.Errorf("withdrawStakeRewards, decreaseReferenceCount error: %v", err)
	}

	// remove stake starting info
	delStakeStartingInfo(s, stakeInfo.StakeAddress, dec)
	return rewards, nil
}

func CalculateStakeRewards(s *native.NativeContract, stakeAddress common.Address, dec []byte, endPeriod uint64) (*big.Int, error) {
	// fetch starting info for delegation
	startingInfo, err := GetStakeStartingInfo(s, stakeAddress, dec)
	if err != nil {
		return nil, fmt.Errorf("CalculateStakeRewards, GetStakeStartingInfo error: %v", err)
	}

	startPeriod := startingInfo.StartPeriod
	stake := startingInfo.Stake

	// sanity check
	if startPeriod > endPeriod {
		panic("startPeriod cannot be greater than endPeriod")
	}
	if stake.Sign() < 0 {
		panic("stake should not be negative")
	}

	// return staking * (ending - starting)
	starting, err := GetValidatorSnapshotRewards(s, dec, startPeriod)
	if err != nil {
		return nil, fmt.Errorf("CalculateStakeRewards, GetValidatorSnapshotRewards start error: %v", err)
	}
	ending, err := GetValidatorSnapshotRewards(s, dec, endPeriod)
	if err != nil {
		return nil, fmt.Errorf("CalculateStakeRewards, GetValidatorSnapshotRewards end error: %v", err)
	}
	difference := new(big.Int).Sub(ending.AccumulatedRewardsRatio, starting.AccumulatedRewardsRatio)
	if difference.Sign() < 0 {
		panic("negative rewards should not be possible")
	}
	rewardsD := new(big.Int).Mul(difference, stake)
	rewards := new(big.Int).Div(rewardsD, RatioDecimal)
	return rewards, nil
}

func initializeStake(s *native.NativeContract, stakeInfo *StakeInfo, dec []byte) error {
	// period has already been incremented
	validatorAccumulatedRewards, err := GetValidatorAccumulatedRewards(s, dec)
	if err != nil {
		return fmt.Errorf("initializeStake, GetValidatorAccumulatedRewards start error: %v", err)
	}
	previousPeriod := validatorAccumulatedRewards.Period - 1

	// increment reference count for the period we're going to track
	err = increaseReferenceCount(s, dec, previousPeriod)
	if err != nil {
		return fmt.Errorf("initializeStake, increaseReferenceCount start error: %v", err)
	}

	stake := stakeInfo.Amount
	err = setStakeStartingInfo(s, stakeInfo.StakeAddress, dec, &StakeStartingInfo{previousPeriod,
		stake, s.ContractRef().BlockHeight()})
	if err != nil {
		return fmt.Errorf("initializeStake, setStakeStartingInfo error: %v", err)
	}
	return nil
}

func withdrawCommission(s *native.NativeContract, stakeAddress common.Address, dec []byte) (*big.Int, error) {
	accumulatedCommission, err := GetAccumulatedCommission(s, dec)
	if err != nil {
		return nil, fmt.Errorf("withdrawCommission, GetAccumulatedCommission error: %v", err)
	}

	// update the outstanding rewards
	outstanding, err := GetOutstandingRewards(s)
	if err != nil {
		return nil, fmt.Errorf("withdrawCommission, GetOutstandingRewards error: %v", err)
	}
	validatorOutstanding, err := GetValidatorOutstandingRewards(s, dec)
	if err != nil {
		return nil, fmt.Errorf("withdrawCommission, GetValidatorOutstandingRewards error: %v", err)
	}
	err = setOutstandingRewards(s, &OutstandingRewards{Rewards: new(big.Int).Sub(outstanding.Rewards, accumulatedCommission.Amount)})
	if err != nil {
		return nil, fmt.Errorf("withdrawCommission, setOutstandingRewards error: %v", err)
	}
	err = setValidatorOutstandingRewards(s, dec, &ValidatorOutstandingRewards{Rewards: new(big.Int).Sub(validatorOutstanding.Rewards, accumulatedCommission.Amount)})
	if err != nil {
		return nil, fmt.Errorf("withdrawCommission, setValidatorOutstandingRewards error: %v", err)
	}

	err = nativeTransfer(s, this, stakeAddress, accumulatedCommission.Amount)
	if err != nil {
		return nil, fmt.Errorf("withdrawCommission, nativeTransfer commission error: %v", err)
	}
	return accumulatedCommission.Amount, nil
}

func allocateRewardsToValidator(s *native.NativeContract, validator *Validator, rewards *big.Int) error {
	dec, err := hexutil.Decode(validator.ConsensusPubkey)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, decode pubkey error: %v", err)
	}

	// commission = commission*reward/100
	commission := new(big.Int).Div(new(big.Int).Mul(validator.Commission.Rate, rewards), new(big.Int).SetUint64(100))
	// stake reward = reward-commission
	stakeRewards := new(big.Int).Sub(rewards, commission)

	// update accumulate commission
	currentCommission, err := GetAccumulatedCommission(s, dec)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, GetAccumulatedCommission error: %v", err)
	}
	currentCommission.Amount = new(big.Int).Add(currentCommission.Amount, commission)
	err = setAccumulatedCommission(s, dec, currentCommission)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, setAccumulatedCommission error: %v", err)
	}

	// update accumulate rewards
	validatorAccumulatedRewards, err := GetValidatorAccumulatedRewards(s, dec)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, GetValidatorAccumulatedRewards error: %v", err)
	}
	validatorAccumulatedRewards.Rewards = new(big.Int).Add(validatorAccumulatedRewards.Rewards, stakeRewards)
	err = setValidatorAccumulatedRewards(s, dec, validatorAccumulatedRewards)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, setValidatorAccumulatedRewards error: %v", err)
	}

	// update validator outstanding
	outstanding, err := GetValidatorOutstandingRewards(s, dec)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, GetValidatorOutstandingRewards error: %v", err)
	}
	outstanding.Rewards = new(big.Int).Add(outstanding.Rewards, rewards)
	err = setValidatorOutstandingRewards(s, dec, outstanding)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, setValidatorOutstandingRewards error: %v", err)
	}
	return nil
}
