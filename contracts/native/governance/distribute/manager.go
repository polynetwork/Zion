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
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

var (
	gasTable = map[string]uint64{
		MethodName:             0,
		MethodPropose:          30000,
		MethodVote:             30000,
		MethodEpoch:            0,
		MethodGetEpochByID:     0,
		MethodProof:            0,
		MethodGetChangingEpoch: 0,
	}
)

func InitNodeManager() {
	InitABI()
	native.Contracts[this] = RegisterDistributeContract
}

func RegisterDistributeContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodWithdrawStakeRewards, WithdrawStakeRewards)
	s.Register(MethodWithdrawCommssion, WithdrawCommssion)
}

func WithdrawStakeRewards(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	params := &WithdrawStakeRewardsParam{}
	if err := utils.UnpackMethod(ABI, MethodWithdrawStakeRewards, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("CreateValidator, unpack params error: %v", err)
	}
	dec, err := hexutil.Decode(params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, decode pubkey error: %v", err)
	}

	validator, found, err := node_manager.GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, node_manager.GetValidator error: %v", err)
	}
}
