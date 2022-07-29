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

package backend

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/contracts/native/economic"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
)

// reward distribute native token to `governance contract` and `reward pool`
func (s *backend) reward(state *state.StateDB, height *big.Int) error {
	// genesis block do not need to distribute reward
	if height.Uint64() == 0 {
		return nil
	}

	list, err := s.getRewardList(state, height)
	if err != nil {
		return err
	}

	var sRwd string
	for _, v := range list {
		state.AddBalance(v.Address, v.Amount)
		sRwd += fmt.Sprintf("address: %s, amount %v;", v.Address.Hex(), v.Amount)
	}
	log.Debug("reward", "list", sRwd)

	return nil
}

// prepare for slashing...
// todo(fuk): slash for governance
func (s *backend) slash(ctx *systemTxContext) error {
	return s.executeTransaction(ctx, contractAddr, nil)
}

func (s *backend) getRewardList(state *state.StateDB, height *big.Int) ([]*economic.RewardAmount, error) {
	caller := s.signer.Address()
	ref := s.getSystemCaller(state, height)
	payload, err := new(economic.MethodRewardInput).Encode()
	if err != nil {
		return nil, fmt.Errorf("encode reward input failed: %v", err)
	}
	enc, _, err := ref.NativeCall(caller, utils.EconomicContractAddress, payload)
	if err != nil {
		return nil, fmt.Errorf("reward native call failed: %v", err)
	}
	output := new(economic.MethodRewardOutput)
	if err := output.Decode(enc); err != nil {
		return nil, fmt.Errorf("reward output decode failed: %v", err)
	}
	return output.List, nil
}
