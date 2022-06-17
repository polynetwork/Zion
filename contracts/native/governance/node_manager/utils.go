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
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/core"
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

func withdrawDelegationRewards(s *native.NativeContract, validator *Validator, stakeInfo *StakeInfo) error {
	dec, err := hexutil.Decode(validator.ConsensusPubkey)
	if err != nil {
		return fmt.Errorf("withdrawDelegationRewards, decode pubkey error: %v", err)
	}

	// get stake starting info
	stakeStaringInfo, err := GetStakeStartingInfo(s, stakeInfo.StakeAddress)
	if err != nil {
		return fmt.Errorf("withdrawDelegationRewards, GetStakeStartingInfo error: %v", err)
	}

	// end current period and calculate rewards
	endingPeriod, err := IncreaseValidatorPeriod(s, validator)
	if err != nil {
		return fmt.Errorf("withdrawDelegationRewards, IncreaseValidatorPeriod error: %v", err)
	}
	rewardsRaw, err := CalculateStakeRewards(s, stakeInfo.StakeAddress, dec, endingPeriod)
	if err != nil {
		return fmt.Errorf("withdrawDelegationRewards, CalculateStakeRewards error: %v", err)
	}
	outstanding, err := GetValidatorOutstandingRewards(s, dec)
	if err != nil {
		return fmt.Errorf("withdrawDelegationRewards, GetValidatorOutstandingRewards error: %v", err)
	}

}

func CalculateStakeRewards(s *native.NativeContract, stakeAddress common.Address, dec []byte, endPeriod uint64) (*big.Int, error) {
	height := s.ContractRef().BlockHeight()
	// fetch starting info for delegation
	startingInfo, err := GetStakeStartingInfo(s, stakeAddress)
	if err != nil {
		return nil, fmt.Errorf("CalculateStakeRewards, GetStakeStartingInfo error: %v", err)
	}

	if startingInfo.Height.Cmp(height) == 0 {
		// started this height, no rewards yet
		return common.Big0, nil
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
