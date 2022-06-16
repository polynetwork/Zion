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

func AfterValidatorCreated(s *native.NativeContract, validator *node_manager.Validator) error {
	dec, err := hexutil.Decode(validator.ConsensusPubkey)
	if err != nil {
		return fmt.Errorf("AfterValidatorCreated, decode pubkey error: %v", err)
	}

	// set initial historical rewards (period 0) with reference count of 1
	err = setValidatorSnapshotRewards(s, dec, 0, &ValidatorSnapshotRewards{new(big.Int), 1})
	if err != nil {
		return fmt.Errorf("AfterValidatorCreated, setValidatorSnapshotRewards error: %v", err)
	}

	// set accumulate rewards (starting at period 1)
	err = setValidatorAccumulatedRewards(s, dec, &ValidatorAccumulatedRewards{new(big.Int), 1})
	if err != nil {
		return fmt.Errorf("AfterValidatorCreated, setValidatorAccumulatedRewards error: %v", err)
	}

	// set accumulated commission
	err = setAccumulatedCommission(s, dec, &AccumulatedCommission{new(big.Int)})
	if err != nil {
		return fmt.Errorf("AfterValidatorCreated, setAccumulatedCommission error: %v", err)
	}

	// set outstanding rewards
	err = setValidatorOutstandingRewards(s, dec, &ValidatorOutstandingRewards{Rewards: new(big.Int)})
	if err != nil {
		return fmt.Errorf("AfterValidatorCreated, setValidatorOutstandingRewards error: %v", err)
	}
	return nil
}

func BeforeStakeCreated(s *native.NativeContract, validator *node_manager.Validator) error {
	_, err := IncreaseValidatorPeriod(s, validator)
	return err
}

func BeforeStakeModified(s *native.NativeContract, validator *node_manager.Validator) error {

}