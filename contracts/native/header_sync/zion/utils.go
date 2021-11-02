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

package zion

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/backend"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/signer"
	"github.com/ethereum/go-ethereum/core/types"
)

var verifier = &signer.SignerImpl{}

func VerifyHeader(header *types.Header, validators []common.Address, checkEpochChange bool) (
	nextEpochStartHeight uint64, nextEpochValidators []common.Address, err error) {

	if err = backend.CustomVerifyHeader(header); err != nil {
		return
	}

	var extra *types.HotstuffExtra
	valset := backend.NewDefaultValSet(validators)
	if extra, err = verifier.VerifyHeader(header, valset, true); err != nil {
		return
	}

	// DONT need to check epoch change
	if !checkEpochChange {
		return
	}

	// epoch NOT changed
	if len(extra.Validators) == 0 {
		return
	}

	// new validators taking effective at the next block, current block header ONLY record
	// the next epoch validators, and new epoch validators should be 2/3 of old validators.
	if err = valset.CheckQuorum(extra.Validators); err != nil {
		return
	}
	nextEpochStartHeight = header.Number.Uint64() + 1
	nextEpochValidators = extra.Validators
	return
}

func getValidatorsFromHeader(header *types.Header) ([]common.Address, error) {
	extra, err := types.ExtractHotstuffExtra(header)
	if err != nil {
		return nil, err
	}
	if extra.Validators == nil || len(extra.Validators) == 0 {
		return nil, fmt.Errorf("not epoch change header, validators is empty")
	}
	return extra.Validators, nil
}
