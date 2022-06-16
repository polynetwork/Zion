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

package distribute

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"math/big"
)

var RatioDecimal = new(big.Int).SetUint64(1000000)

// IncreaseValidatorPeriod return the period just ended
func IncreaseValidatorPeriod(s *native.NativeContract, validator *node_manager.Validator) (uint64, error) {
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

func withdrawDelegationRewards(s *native.NativeContract, validator *node_manager.Validator) error {

}