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
	"github.com/ethereum/go-ethereum/crypto"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
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

const (
	// The minimum distance between two adjacent epochs is 60 blocks
	MinEpochValidPeriod uint64 = 60
	// The default value of distance for two adjacent epochs
	DefaultEpochValidPeriod uint64 = 86400
	// The max distance between two adjacent epochs
	MaxEpochValidPeriod uint64 = 86400 * 10
	// Consensus engine allows at least 4 validators, this denote F >= 1
	MinProposalPeersLen int = 4
	// Consensus engine allows at most 100 validators, this denote F <= 33
	MaxProposalPeersLen int = 100
	// Every validator can propose at most 6 proposals in an epoch
	MaxProposalNumPerEpoch int = 6
	// Proposal should be voted and passed in period
	MinVoteEffectivePeriod uint64 = 10
)

func InitNodeManager() {
	InitABI()
	native.Contracts[this] = RegisterNodeManagerContract
}

func RegisterNodeManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodCreateValidator, CreateValidator)
	s.Register(MethodUpdateValidator, UpdateValidator)
	s.Register(MethodStake, Stake)
	s.Register(MethodUnStake, UnStake)
	s.Register(MethodWithdraw, Withdraw)
	s.Register(MethodChangeEpoch, ChangeEpoch)
}

func CreateValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	height := s.ContractRef().BlockHeight()
	caller := ctx.Caller

	params := &CreateValidatorParam{}
	if err := utils.UnpackMethod(ABI, MethodCreateValidator, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("CreateValidator, unpack params error: %v", err)
	}

	// check pub key
	dec, err := hexutil.Decode(params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, decode pubkey error: %v", err)
	}
	pubkey, err := crypto.DecompressPubkey(dec)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, decompress pubkey error: %v", err)
	}
	addr := crypto.PubkeyToAddress(*pubkey)
	if addr == common.EmptyAddress {
		return nil, fmt.Errorf("invalid pubkey")
	}
	if addr == caller || addr == params.ProposalAddress {
		return nil, fmt.Errorf("stake, consensus and proposal address can not be duplicate")
	}

	// check commission
	globalConfig, err := GetGlobalConfig(s)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, GetGlobalConfig error: %v", err)
	}
	if params.Commission.Sign() == -1 {
		return nil, fmt.Errorf("CreateValidator, commission must be positive")
	}
	if params.Commission.Cmp(globalConfig.MaxCommission) == 1 {
		return nil, fmt.Errorf("CreateValidator, commission can not greater than globalConfig.MaxCommission: %s",
			globalConfig.MaxCommission.String())
	}

	// check desc
	if uint64(len(params.Desc)) > globalConfig.MaxDescLength {
		return nil, fmt.Errorf("CreateValidator, desc length more than limit %d", globalConfig.MaxDescLength)
	}

	// check to see if the pubkey has been registered before
	_, found, err := GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, GetValidator error: %v", err)
	}
	if found {
		return nil, fmt.Errorf("CreateValidator, validator already exist")
	}

	// check initial stake
	if globalConfig.MinInitialStake.Cmp(params.InitStake) == 1 {
		return nil, fmt.Errorf("CreateValidator, initial stake %s is less than min initial stake %s",
			params.InitStake.String(), globalConfig.MinInitialStake.String())
	}

	// store validator
	validator := &Validator{
		StakeAddress:    caller,
		ConsensusPubkey: params.ConsensusPubkey,
		ProposalAddress: params.ProposalAddress,
		Commission:      &Commission{Rate: params.Commission, UpdateHeight: height},
		Status:          Unlocked,
		Jailed:          false,
		UnlockTime:      common.Big0,
		TotalStake:      params.InitStake,
		SelfStake:       params.InitStake,
		Desc:            params.Desc,
	}
	err = setValidator(s, validator)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, setValidator error: %v", err)
	}
	// add validator to all validators pool
	err = addToAllValidators(s, validator.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, addToAllValidators error: %v", err)
	}

	// deposit native token
	err = deposit(s, caller, params.InitStake, validator)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, deposit error: %v", err)
	}

	err = s.AddNotify(ABI, []string{MethodCreateValidator}, params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, AddNotify error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func UpdateValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	height := s.ContractRef().BlockHeight()
	caller := ctx.Caller

	params := &UpdateValidatorParam{}
	if err := utils.UnpackMethod(ABI, MethodUpdateValidator, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("UpdateValidator, unpack params error: %v", err)
	}

	dec, err := hexutil.Decode(params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("UpdateValidator, decode pubkey error: %v", err)
	}
	validator, found, err := GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("UpdateValidator, get validator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("UpdateValidator, can not found record")
	}
	if validator.StakeAddress != caller {
		return nil, fmt.Errorf("UpdateValidator, stake address is not caller")
	}
	globalConfig, err := GetGlobalConfig(s)
	if err != nil {
		return nil, fmt.Errorf("UpdateValidator, GetGlobalConfig error: %v", err)
	}

	if params.ProposalAddress != common.EmptyAddress {
		validator.ProposalAddress = params.ProposalAddress
	}

	if params.Commission != common.Big0 {
		// check commission
		if params.Commission.Sign() == -1 {
			return nil, fmt.Errorf("UpdateValidator, commission must be positive")
		}
		if params.Commission.Cmp(globalConfig.MaxCommission) == 1 {
			return nil, fmt.Errorf("UpdateValidator, commission can not greater than globalConfig.MaxCommission: %s",
				globalConfig.MaxCommission.String())
		}
		validator.Commission = &Commission{Rate: params.Commission, UpdateHeight: height}
	}

	if params.Desc != "" {
		// check desc
		if uint64(len(params.Desc)) > globalConfig.MaxDescLength {
			return nil, fmt.Errorf("UpdateValidator, desc length more than limit %d", globalConfig.MaxDescLength)
		}
		validator.Desc = params.Desc
	}

	err = setValidator(s, validator)
	if err != nil {
		return nil, fmt.Errorf("UpdateValidator, setValidator error: %v", err)
	}

	err = s.AddNotify(ABI, []string{MethodUpdateValidator}, params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("UpdateValidator, AddNotify error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func Stake(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	params := &StakeParam{}
	if err := utils.UnpackMethod(ABI, MethodStake, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("Stake, unpack params error: %v", err)
	}
	dec, err := hexutil.Decode(params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("Stake, decode pubkey error: %v", err)
	}

	// check to see if the pubkey has been registered
	validator, found, err := GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("Stake, GetValidator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("Stake, validator is not exist")
	}

	// deposit native token
	err = deposit(s, caller, params.Amount, validator)
	if err != nil {
		return nil, fmt.Errorf("Stake, deposit error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func UnStake(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	params := &UnStakeParam{}
	if err := utils.UnpackMethod(ABI, MethodUnStake, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("UnStake, unpack params error: %v", err)
	}
	dec, err := hexutil.Decode(params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("UnStake, decode pubkey error: %v", err)
	}

	// check to see if the pubkey has been registered
	validator, found, err := GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("UnStake, GetValidator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("UnStake, validator is not exist")
	}

	// unStake native token
	err = unStake(s, caller, params.Amount, validator)
	if err != nil {
		return nil, fmt.Errorf("UnStake, deposit error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func Withdraw(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	amount, err := filterExpiredUnlockingInfo(s, caller)
	if err != nil {
		return nil, fmt.Errorf("Withdraw, filterExpiredUnlockingInfo error: %v", err)
	}

	if amount.Sign() == 1 {
		err = withdrawUnlockPool(s, amount)
		if err != nil {
			return nil, fmt.Errorf("Withdraw, withdrawUnlockPool error: %v", err)
		}
		err = nativeTransfer(s, this, caller, amount)
		if err != nil {
			return nil, fmt.Errorf("Withdraw, nativeTransfer error: %v", err)
		}
	}
	return utils.ByteSuccess, nil
}

func ChangeEpoch(s *native.NativeContract) ([]byte, error) {

}
