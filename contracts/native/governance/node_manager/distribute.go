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

const (
	TokenPrecision = 18
	PercentPrecision = 4
)

var (
	TokenDecimal = new(big.Int).Exp(big.NewInt(10), big.NewInt(TokenPrecision), nil)
	PercentDecimal = new(big.Int).Exp(big.NewInt(10), big.NewInt(PercentPrecision), nil)
)

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
	ratio, err := validatorAccumulatedRewards.Rewards.DivWithTokenDecimal(validator.TotalStake)

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
	newRatio, err := validatorSnapshotRewards.AccumulatedRewardsRatio.Add(ratio)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, validatorSnapshotRewards.AccumulatedRewardsRatio.Add error: %v", err)
	}
	newValidatorSnapshotRewards := &ValidatorSnapshotRewards{
		AccumulatedRewardsRatio: newRatio,
		ReferenceCount:          1,
	}
	err = setValidatorSnapshotRewards(s, dec, validatorAccumulatedRewards.Period, newValidatorSnapshotRewards)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, setValidatorSnapshotRewards error: %v", err)
	}

	// set accumulate rewards, incrementing period by 1
	newValidatorAccumulatedRewards := &ValidatorAccumulatedRewards{
		Rewards: NewDecFromBigInt(new(big.Int)),
		Period:  validatorAccumulatedRewards.Period + 1,
	}
	err = setValidatorAccumulatedRewards(s, dec, newValidatorAccumulatedRewards)
	if err != nil {
		return 0, fmt.Errorf("IncreaseValidatorPeriod, setValidatorAccumulatedRewards error: %v", err)
	}

	return validatorAccumulatedRewards.Period, nil
}

func withdrawStakeRewards(s *native.NativeContract, validator *Validator, stakeInfo *StakeInfo) (Dec, error) {
	dec, err := hexutil.Decode(validator.ConsensusPubkey)
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawStakeRewards, decode pubkey error: %v", err)
	}

	// end current period and calculate rewards
	endingPeriod, err := IncreaseValidatorPeriod(s, validator)
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawStakeRewards, IncreaseValidatorPeriod error: %v", err)
	}
	rewards, err := CalculateStakeRewards(s, stakeInfo.StakeAddress, dec, endingPeriod)
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawStakeRewards, CalculateStakeRewards error: %v", err)
	}

	err = nativeTransfer(s, this, stakeInfo.StakeAddress, rewards.BigInt())
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawStakeRewards, nativeTransfer error: %v", err)
	}

	// update the outstanding rewards
	outstanding, err := GetOutstandingRewards(s)
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawStakeRewards, GetOutstandingRewards error: %v", err)
	}
	validatorOutstanding, err := GetValidatorOutstandingRewards(s, dec)
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawStakeRewards, GetValidatorOutstandingRewards error: %v", err)
	}
	newOutstandingRewards, err := outstanding.Rewards.Sub(rewards)
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawStakeRewards, outstanding.Rewards.Sub error: %v", err)
	}
	err = setOutstandingRewards(s, &OutstandingRewards{Rewards: newOutstandingRewards})
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawStakeRewards, setOutstandingRewards error: %v", err)
	}
	newValidatorOutstandingRewards, err := validatorOutstanding.Rewards.Sub(rewards)
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawStakeRewards, validatorOutstanding.Rewards.Sub error: %v", err)
	}
	err = setValidatorOutstandingRewards(s, dec, &ValidatorOutstandingRewards{Rewards: newValidatorOutstandingRewards})
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawStakeRewards, setValidatorOutstandingRewards error: %v", err)
	}

	// decrement reference count of starting period
	startingInfo, err := GetStakeStartingInfo(s, stakeInfo.StakeAddress, dec)
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawStakeRewards, GetStakeStartingInfo error: %v", err)
	}
	startPeriod := startingInfo.StartPeriod
	err = decreaseReferenceCount(s, dec, startPeriod)
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawStakeRewards, decreaseReferenceCount error: %v", err)
	}

	// remove stake starting info
	delStakeStartingInfo(s, stakeInfo.StakeAddress, dec)
	return rewards, nil
}

func CalculateStakeRewards(s *native.NativeContract, stakeAddress common.Address, dec []byte, endPeriod uint64) (Dec, error) {
	// fetch starting info for delegation
	startingInfo, err := GetStakeStartingInfo(s, stakeAddress, dec)
	if err != nil {
		return Dec{nil}, fmt.Errorf("CalculateStakeRewards, GetStakeStartingInfo error: %v", err)
	}

	startPeriod := startingInfo.StartPeriod
	stake := startingInfo.Stake

	// sanity check
	if startPeriod > endPeriod {
		panic("startPeriod cannot be greater than endPeriod")
	}

	// return staking * (ending - starting)
	starting, err := GetValidatorSnapshotRewards(s, dec, startPeriod)
	if err != nil {
		return Dec{nil}, fmt.Errorf("CalculateStakeRewards, GetValidatorSnapshotRewards start error: %v", err)
	}
	ending, err := GetValidatorSnapshotRewards(s, dec, endPeriod)
	if err != nil {
		return Dec{nil}, fmt.Errorf("CalculateStakeRewards, GetValidatorSnapshotRewards end error: %v", err)
	}
	difference, err := ending.AccumulatedRewardsRatio.Sub(starting.AccumulatedRewardsRatio)
	if err != nil {
		return Dec{nil}, fmt.Errorf("CalculateStakeRewards error: %v", err)
	}
	rewards, err := difference.MulWithTokenDecimal(stake)
	if err != nil {
		return Dec{nil}, fmt.Errorf("CalculateStakeRewards error: %v", err)
	}
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

func withdrawCommission(s *native.NativeContract, stakeAddress common.Address, dec []byte) (Dec, error) {
	accumulatedCommission, err := GetAccumulatedCommission(s, dec)
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawCommission, GetAccumulatedCommission error: %v", err)
	}

	// update the outstanding rewards
	outstanding, err := GetOutstandingRewards(s)
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawCommission, GetOutstandingRewards error: %v", err)
	}
	validatorOutstanding, err := GetValidatorOutstandingRewards(s, dec)
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawCommission, GetValidatorOutstandingRewards error: %v", err)
	}
	newOutstandingRewards, err := outstanding.Rewards.Sub(accumulatedCommission.Amount)
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawCommission, outstanding.Rewards.Sub error: %v", err)
	}
	err = setOutstandingRewards(s, &OutstandingRewards{Rewards: newOutstandingRewards})
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawCommission, setOutstandingRewards error: %v", err)
	}
	newValidatorOutstandingRewards, err := validatorOutstanding.Rewards.Sub(accumulatedCommission.Amount)
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawCommission, validatorOutstanding.Rewards.Sub error: %v", err)
	}
	err = setValidatorOutstandingRewards(s, dec, &ValidatorOutstandingRewards{Rewards: newValidatorOutstandingRewards})
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawCommission, setValidatorOutstandingRewards error: %v", err)
	}

	err = nativeTransfer(s, this, stakeAddress, accumulatedCommission.Amount.BigInt())
	if err != nil {
		return Dec{nil}, fmt.Errorf("withdrawCommission, nativeTransfer commission error: %v", err)
	}
	return accumulatedCommission.Amount, nil
}

func allocateRewardsToValidator(s *native.NativeContract, validator *Validator, rewards Dec) error {
	dec, err := hexutil.Decode(validator.ConsensusPubkey)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, decode pubkey error: %v", err)
	}

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
	currentCommission, err := GetAccumulatedCommission(s, dec)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, GetAccumulatedCommission error: %v", err)
	}
	currentCommission.Amount, err = currentCommission.Amount.Add(commission)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, currentCommission.Amount.Add error: %v", err)
	}
	err = setAccumulatedCommission(s, dec, currentCommission)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, setAccumulatedCommission error: %v", err)
	}

	// update accumulate rewards
	validatorAccumulatedRewards, err := GetValidatorAccumulatedRewards(s, dec)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, GetValidatorAccumulatedRewards error: %v", err)
	}
	validatorAccumulatedRewards.Rewards, err = validatorAccumulatedRewards.Rewards.Add(stakeRewards)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, validatorAccumulatedRewards.Rewards.Add error: %v", err)
	}
	err = setValidatorAccumulatedRewards(s, dec, validatorAccumulatedRewards)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, setValidatorAccumulatedRewards error: %v", err)
	}

	// update validator outstanding
	outstanding, err := GetValidatorOutstandingRewards(s, dec)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, GetValidatorOutstandingRewards error: %v", err)
	}
	outstanding.Rewards, err = outstanding.Rewards.Add(rewards)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, outstanding.Rewards.Add error: %v", err)
	}
	err = setValidatorOutstandingRewards(s, dec, outstanding)
	if err != nil {
		return fmt.Errorf("allocateRewardsToValidator, setValidatorOutstandingRewards error: %v", err)
	}
	return nil
}
