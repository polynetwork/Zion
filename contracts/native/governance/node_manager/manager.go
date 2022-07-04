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
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	CREATE_VALIDATOR_EVENT       = "CreateValidator"
	UPDATE_VALIDATOR_EVENT       = "UpdateValidator"
	UPDATE_COMMISSION_EVENT      = "UpdateCommission"
	STAKE_EVENT                  = "Stake"
	UNSTAKE_EVENT                = "UnStake"
	WITHDRAW_EVENT               = "Withdraw"
	CANCEL_VALIDATOR_EVENT       = "CancelValidator"
	WITHDRAW_VALIDATOR_EVENT     = "WithdrawValidator"
	CHANGE_EPOCH_EVENT           = "ChangeEpoch"
	WITHDRAW_STAKE_REWARDS_EVENT = "WithdrawStakeRewards"
	WITHDRAW_COMMISSION_EVENT    = "WithdrawCommission"
)

var (
	gasTable = map[string]uint64{
		MethodCreateValidator:      0,
		MethodUpdateValidator:      0,
		MethodUpdateCommission:     0,
		MethodStake:                0,
		MethodUnStake:              0,
		MethodWithdraw:             0,
		MethodCancelValidator:      0,
		MethodWithdrawValidator:    0,
		MethodChangeEpoch:          0,
		MethodWithdrawStakeRewards: 0,
		MethodWithdrawCommission:   0,
		MethodEndBlock:             0,
		MethodGetGlobalConfig:      0,
		MethodGetCommunityInfo:     0,
		MethodGetCurrentEpochInfo:  0,
	}
)

func InitNodeManager() {
	InitABI()
	native.Contracts[this] = RegisterNodeManagerContract
}

func RegisterNodeManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodCreateValidator, CreateValidator)
	s.Register(MethodUpdateValidator, UpdateValidator)
	s.Register(MethodUpdateCommission, UpdateCommission)
	s.Register(MethodStake, Stake)
	s.Register(MethodUnStake, UnStake)
	s.Register(MethodWithdraw, Withdraw)
	s.Register(MethodCancelValidator, CancelValidator)
	s.Register(MethodWithdrawValidator, WithdrawValidator)
	s.Register(MethodChangeEpoch, ChangeEpoch)

	s.Register(MethodWithdrawStakeRewards, WithdrawStakeRewards)
	s.Register(MethodWithdrawCommission, WithdrawCommission)
	s.Register(MethodEndBlock, EndBlock)

	// Query
	s.Register(MethodGetGlobalConfig, GetGlobalConfig)
	s.Register(MethodGetCommunityInfo, GetCommunityInfo)
	s.Register(MethodGetCurrentEpochInfo, GetCurrentEpochInfo)
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
		return nil, fmt.Errorf("CreateValidator， invalid pubkey")
	}
	if params.ProposalAddress == common.EmptyAddress {
		return nil, fmt.Errorf("CreateValidator， invalid proposalAddress")
	}

	// check commission
	globalConfig, err := getGlobalConfig(s)
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
		StakeAddress:     caller,
		ConsensusPubkey:  params.ConsensusPubkey,
		ConsensusAddress: addr,
		ProposalAddress:  params.ProposalAddress,
		Commission:       &Commission{Rate: NewDecFromBigInt(params.Commission), UpdateHeight: height},
		Status:           Unlock,
		Jailed:           false,
		UnlockHeight:     new(big.Int),
		TotalStake:       NewDecFromBigInt(params.InitStake),
		SelfStake:        NewDecFromBigInt(params.InitStake),
		Desc:             params.Desc,
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

	// call distrubute hook
	err = AfterValidatorCreated(s, validator)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, distribute.AfterValidatorCreated error: %v", err)
	}

	// deposit native token
	err = deposit(s, caller, NewDecFromBigInt(params.InitStake), validator)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, deposit error: %v", err)
	}

	err = s.AddNotify(ABI, []string{CREATE_VALIDATOR_EVENT}, params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, AddNotify error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func UpdateValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
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
	globalConfig, err := getGlobalConfig(s)
	if err != nil {
		return nil, fmt.Errorf("UpdateValidator, GetGlobalConfig error: %v", err)
	}

	if params.ProposalAddress != common.EmptyAddress {
		validator.ProposalAddress = params.ProposalAddress
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

	err = s.AddNotify(ABI, []string{UPDATE_VALIDATOR_EVENT}, params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("UpdateValidator, AddNotify error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func UpdateCommission(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	height := s.ContractRef().BlockHeight()
	caller := ctx.Caller

	params := &UpdateCommissionParam{}
	if err := utils.UnpackMethod(ABI, MethodUpdateCommission, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("UpdateCommission, unpack params error: %v", err)
	}

	dec, err := hexutil.Decode(params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("UpdateCommission, decode pubkey error: %v", err)
	}
	validator, found, err := GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("UpdateCommission, get validator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("UpdateCommission, can not found record")
	}
	if validator.StakeAddress != caller {
		return nil, fmt.Errorf("UpdateCommission, stake address is not caller")
	}
	globalConfig, err := getGlobalConfig(s)
	if err != nil {
		return nil, fmt.Errorf("UpdateCommission, GetGlobalConfig error: %v", err)
	}

	// check commission
	if params.Commission.Sign() == -1 {
		return nil, fmt.Errorf("UpdateCommission, commission must be positive")
	}
	if params.Commission.Cmp(new(big.Int).SetUint64(100)) == 1 {
		return nil, fmt.Errorf("UpdateCommission, commission can not more than 100")
	}
	if params.Commission.Cmp(globalConfig.MaxCommission) == 1 {
		return nil, fmt.Errorf("UpdateCommission, commission can not greater than globalConfig.MaxCommission: %s",
			globalConfig.MaxCommission.String())
	}
	if height.Cmp(new(big.Int).Add(validator.Commission.UpdateHeight, globalConfig.BlockPerEpoch)) < 0 {
		return nil, fmt.Errorf("UpdateCommission, commission can not changed in one epoch twice")
	}

	validator.Commission = &Commission{Rate: NewDecFromBigInt(params.Commission), UpdateHeight: height}

	err = setValidator(s, validator)
	if err != nil {
		return nil, fmt.Errorf("UpdateCommission, setValidator error: %v", err)
	}

	err = s.AddNotify(ABI, []string{UPDATE_COMMISSION_EVENT}, params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("UpdateCommission, AddNotify error: %v", err)
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
	amount := NewDecFromBigInt(params.Amount)

	// check to see if the pubkey has been registered
	validator, found, err := GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("Stake, GetValidator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("Stake, validator is not exist")
	}

	// deposit native token
	err = deposit(s, caller, amount, validator)
	if err != nil {
		return nil, fmt.Errorf("Stake, deposit error: %v", err)
	}

	// update validator
	if validator.StakeAddress == caller {
		validator.SelfStake, err = validator.SelfStake.Add(amount)
		if err != nil {
			return nil, fmt.Errorf("Stake, validator.SelfStake.Add error: %v", err)
		}
	} else {
		validator.TotalStake, err = validator.TotalStake.Add(amount)
		if err != nil {
			return nil, fmt.Errorf("Stake, validator.TotalStake.Add error: %v", err)
		}
	}
	err = setValidator(s, validator)
	if err != nil {
		return nil, fmt.Errorf("Stake, setValidator error: %v", err)
	}

	err = s.AddNotify(ABI, []string{STAKE_EVENT}, params.ConsensusPubkey, params.Amount.String())
	if err != nil {
		return nil, fmt.Errorf("Stake, AddNotify error: %v", err)
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
	amount := NewDecFromBigInt(params.Amount)

	// check to see if the pubkey has been registered
	validator, found, err := GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("UnStake, GetValidator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("UnStake, validator is not exist")
	}

	// unStake native token
	err = unStake(s, caller, amount, validator)
	if err != nil {
		return nil, fmt.Errorf("UnStake, unStake error: %v", err)
	}

	// update validator
	if validator.StakeAddress == caller {
		return nil, fmt.Errorf("UnStake, stake address can not unstake")
	} else {
		validator.TotalStake, err = validator.TotalStake.Sub(amount)
		if err != nil {
			return nil, fmt.Errorf("UnStake, validator.TotalStake.Sub error: %v", err)
		}
	}
	if validator.TotalStake.IsZero() && validator.SelfStake.IsZero() {
		err = delValidator(s, params.ConsensusPubkey)
		if err != nil {
			return nil, fmt.Errorf("UnStake, delValidator error: %v", err)
		}
		err = AfterValidatorRemoved(s, validator)
		if err != nil {
			return nil, fmt.Errorf("UnStake, AfterValidatorRemoved error: %v", err)
		}
	} else {
		err = setValidator(s, validator)
		if err != nil {
			return nil, fmt.Errorf("UnStake, setValidator error: %v", err)
		}
	}

	err = s.AddNotify(ABI, []string{UNSTAKE_EVENT}, params.ConsensusPubkey, params.Amount.String())
	if err != nil {
		return nil, fmt.Errorf("UnStake, AddNotify error: %v", err)
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

	if amount.IsPositive() {
		err = withdrawTotalPool(s, amount)
		if err != nil {
			return nil, fmt.Errorf("Withdraw, withdrawTotalPool error: %v", err)
		}
		err = nativeTransfer(s, this, caller, amount.BigInt())
		if err != nil {
			return nil, fmt.Errorf("Withdraw, nativeTransfer error: %v", err)
		}
	}

	err = s.AddNotify(ABI, []string{WITHDRAW_EVENT}, caller.Hex(), amount.BigInt().String())
	if err != nil {
		return nil, fmt.Errorf("Withdraw, AddNotify error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func CancelValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller
	height := s.ContractRef().BlockHeight()

	params := &CancelValidatorParam{}
	if err := utils.UnpackMethod(ABI, MethodCancelValidator, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("CancelValidator, unpack params error: %v", err)
	}
	dec, err := hexutil.Decode(params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("CancelValidator, decode pubkey error: %v", err)
	}

	globalConfig, err := getGlobalConfig(s)
	if err != nil {
		return nil, fmt.Errorf("CancelValidator, GetGlobalConfig error: %v", err)
	}

	validator, found, err := GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("CancelValidator, GetValidator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("CancelValidator, validator is not created")
	}
	if validator.StakeAddress != caller {
		return nil, fmt.Errorf("CancelValidator, stake address %s is not caller", validator.StakeAddress.Hex())
	}

	switch {
	case validator.IsLocked():
		validator.Status = Remove
		validator.UnlockHeight = new(big.Int).Add(height, globalConfig.BlockPerEpoch)
		err = setValidator(s, validator)
		if err != nil {
			return nil, fmt.Errorf("CancelValidator, setValidator error: %v", err)
		}
	case validator.IsUnlocking(height), validator.IsUnlocked(height):
		validator.Status = Remove
	default:
		return nil, fmt.Errorf("CancelValidator, unsupported validator status")
	}
	err = removeFromAllValidators(s, params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("CancelValidator, removeFromAllValidators error: %v", err)
	}

	err = s.AddNotify(ABI, []string{CANCEL_VALIDATOR_EVENT}, params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("CancelValidator, AddNotify error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func WithdrawValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller
	height := s.ContractRef().BlockHeight()

	params := &CancelValidatorParam{}
	if err := utils.UnpackMethod(ABI, MethodWithdrawValidator, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("WithdrawValidator, unpack params error: %v", err)
	}
	dec, err := hexutil.Decode(params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("WithdrawValidator, decode pubkey error: %v", err)
	}
	validator, found, err := GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("WithdrawValidator, GetValidator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("WithdrawValidator, validator is not created")
	}
	if validator.StakeAddress != caller {
		return nil, fmt.Errorf("WithdrawValidator, stake address %s is not caller", validator.StakeAddress.Hex())
	}
	if !validator.IsRemoved(height) {
		return nil, fmt.Errorf("WithdrawValidator, validator is not removed")
	}

	// unStake native token
	err = unStake(s, caller, validator.SelfStake, validator)
	if err != nil {
		return nil, fmt.Errorf("WithdrawValidator, unStake error: %v", err)
	}

	validator.TotalStake, err = validator.TotalStake.Sub(validator.SelfStake)
	if err != nil {
		return nil, fmt.Errorf("WithdrawValidator, validator.TotalStake.Sub error: %v", err)
	}
	validator.SelfStake = NewDecFromBigInt(new(big.Int))
	if validator.TotalStake.IsZero() {
		err = delValidator(s, params.ConsensusPubkey)
		if err != nil {
			return nil, fmt.Errorf("WithdrawValidator, delValidator error: %v", err)
		}
		err = AfterValidatorRemoved(s, validator)
		if err != nil {
			return nil, fmt.Errorf("WithdrawValidator, AfterValidatorRemoved error: %v", err)
		}
	} else {
		err = setValidator(s, validator)
		if err != nil {
			return nil, fmt.Errorf("WithdrawValidator, setValidator error: %v", err)
		}
	}

	_, err = withdrawCommission(s, caller, dec)
	if err != nil {
		return nil, fmt.Errorf("WithdrawValidator, withdrawCommission error: %v", err)
	}
	delAccumulatedCommission(s, dec)

	err = s.AddNotify(ABI, []string{WITHDRAW_VALIDATOR_EVENT}, params.ConsensusPubkey, validator.SelfStake.BigInt().String())
	if err != nil {
		return nil, fmt.Errorf("CancelValidator, AddNotify error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func ChangeEpoch(s *native.NativeContract) ([]byte, error) {
	endHeight := s.ContractRef().BlockHeight()
	startHeight := new(big.Int).Add(endHeight, common.Big1)

	currentEpochInfo, err := GetCurrentEpochInfoImpl(s)
	if err != nil {
		return nil, fmt.Errorf("ChangeEpoch, GetCurrentEpochInfoImpl error: %v", err)
	}
	globalConfig, err := getGlobalConfig(s)

	// anyone can call this if height reaches
	if new(big.Int).Sub(startHeight, currentEpochInfo.StartHeight).Cmp(globalConfig.BlockPerEpoch) == -1 {
		return nil, fmt.Errorf("ChangeEpoch, block height does not reach, current epoch start at %s",
			currentEpochInfo.StartHeight.String())
	}

	epochInfo := &EpochInfo{
		ID:          new(big.Int).Add(currentEpochInfo.ID, common.Big1),
		Validators:  make([]*Peer, 0, globalConfig.ConsensusValidatorNum),
		Voters:      make([]*Peer, 0, globalConfig.VoterValidatorNum),
		StartHeight: startHeight,
	}
	// get all validators
	allValidators, err := GetAllValidators(s)
	if err != nil {
		return nil, fmt.Errorf("ChangeEpoch, GetAllValidators error: %v", err)
	}
	if uint64(len(allValidators.AllValidators)) < globalConfig.ConsensusValidatorNum {
		epochInfo.Validators = currentEpochInfo.Validators
		epochInfo.Voters = currentEpochInfo.Voters
	} else {
		validatorList := make([]*Validator, 0, len(allValidators.AllValidators))
		for _, v := range allValidators.AllValidators {
			dec, err := hexutil.Decode(v)
			if err != nil {
				return nil, fmt.Errorf("ChangeEpoch, decode pubkey error: %v", err)
			}
			validator, found, err := GetValidator(s, dec)
			if err != nil {
				return nil, fmt.Errorf("ChangeEpoch, GetValidator error: %v", err)
			}
			if !found {
				return nil, fmt.Errorf("ChangeEpoch, validator %s not found", v)
			}
			validatorList = append(validatorList, validator)
		}

		// sort by total stake desc, if equal, use the old slice order
		sort.SliceStable(validatorList, func(i, j int) bool {
			return validatorList[i].TotalStake.GT(validatorList[j].TotalStake)
		})
		// update validator status
		for i := 0; uint64(i) < globalConfig.ConsensusValidatorNum; i++ {
			validator := validatorList[i]
			switch {
			case validator.IsLocked():
			case validator.IsUnlocking(endHeight), validator.IsUnlocked(endHeight):
				validator.Status = Lock
			}

			peer := &Peer{
				PubKey:  validator.ConsensusPubkey,
				Address: validator.ConsensusAddress,
			}
			epochInfo.Validators = append(epochInfo.Validators, peer)
			err = setValidator(s, validator)
			if err != nil {
				return nil, fmt.Errorf("ChangeEpoch, set lock validator error: %v", err)
			}
		}
		for i := globalConfig.ConsensusValidatorNum; i < uint64(len(validatorList)); i++ {
			validator := validatorList[i]
			switch {
			case validator.IsLocked():
				validator.Status = Unlock
				validator.UnlockHeight = new(big.Int).Add(startHeight, globalConfig.BlockPerEpoch)
			case validator.IsUnlocking(endHeight), validator.IsUnlocked(endHeight):
			}
			err = setValidator(s, validator)
			if err != nil {
				return nil, fmt.Errorf("ChangeEpoch, set unlock validator error: %v", err)
			}
		}
		//update voters
		for i := 0; uint64(i) < globalConfig.VoterValidatorNum; i++ {
			validator := validatorList[i]
			peer := &Peer{
				PubKey:  validator.ConsensusPubkey,
				Address: validator.ConsensusAddress,
			}
			epochInfo.Voters = append(epochInfo.Voters, peer)
		}
	}

	// update epoch info
	err = setCurrentEpochInfo(s, epochInfo)
	if err != nil {
		return nil, fmt.Errorf("ChangeEpoch, setEpochInfo error: %v", err)
	}
	err = s.AddNotify(ABI, []string{CHANGE_EPOCH_EVENT}, epochInfo.ID.String())
	if err != nil {
		return nil, fmt.Errorf("ChangeEpoch, AddNotify error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func WithdrawStakeRewards(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	params := &WithdrawStakeRewardsParam{}
	if err := utils.UnpackMethod(ABI, MethodWithdrawStakeRewards, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, unpack params error: %v", err)
	}
	dec, err := hexutil.Decode(params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, decode pubkey error: %v", err)
	}

	validator, found, err := GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, node_manager.GetValidator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("WithdrawStakeRewards, validator not found")
	}
	stakeInfo, found, err := GetStakeInfo(s, caller, params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, GetStakeInfo error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("WithdrawStakeRewards, stake info not found")
	}

	rewards, err := withdrawStakeRewards(s, validator, stakeInfo)
	if err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, withdrawStakeRewards error: %v", err)
	}

	// reinitialize the delegation
	err = initializeStake(s, stakeInfo, dec)
	if err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, initializeStake error: %v", err)
	}

	err = s.AddNotify(ABI, []string{WITHDRAW_STAKE_REWARDS_EVENT}, rewards.BigInt().String())
	if err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, AddNotify error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func WithdrawCommission(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	params := &WithdrawCommissionParam{}
	if err := utils.UnpackMethod(ABI, MethodWithdrawCommission, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("WithdrawCommission, unpack params error: %v", err)
	}
	dec, err := hexutil.Decode(params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("WithdrawCommission, decode pubkey error: %v", err)
	}
	validator, found, err := GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("WithdrawCommission, GetValidator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("WithdrawCommission, can not find validator")
	}
	if validator.StakeAddress != caller {
		return nil, fmt.Errorf("WithdrawCommission, caller is not stake address")
	}

	commission, err := withdrawCommission(s, caller, dec)
	if err != nil {
		return nil, fmt.Errorf("WithdrawCommission, withdrawCommission error: %v", err)
	}
	err = setAccumulatedCommission(s, dec, &AccumulatedCommission{NewDecFromBigInt(new(big.Int))})
	if err != nil {
		return nil, fmt.Errorf("WithdrawCommission, setAccumulatedCommission error: %v", err)
	}

	err = s.AddNotify(ABI, []string{WITHDRAW_COMMISSION_EVENT}, params.ConsensusPubkey, commission.BigInt().String())
	if err != nil {
		return nil, fmt.Errorf("WithdrawCommission, AddNotify error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func EndBlock(s *native.NativeContract) ([]byte, error) {
	// contract balance = totalpool + outstanding + new block reward
	balance := NewDecFromBigInt(s.StateDB().GetBalance(this))

	totalPool, err := GetTotalPool(s)
	if err != nil {
		return nil, fmt.Errorf("EndBlock, GetTotalPool error: %v", err)
	}
	outstanding, err := GetOutstandingRewards(s)
	if err != nil {
		return nil, fmt.Errorf("EndBlock, GetOutstandingRewards error: %v", err)
	}

	// cal rewards
	temp, err := outstanding.Rewards.Add(totalPool)
	if err != nil {
		return nil, fmt.Errorf("EndBlock, outstanding.Rewards.Add error: %v", err)
	}
	newRewards, err := balance.Sub(temp)
	if err != nil {
		return nil, fmt.Errorf("EndBlock, balance.Sub error: %v", err)
	}

	epochInfo, err := GetCurrentEpochInfoImpl(s)
	if err != nil {
		return nil, fmt.Errorf("EndBlock, GetCurrentEpochInfoImpl error: %v", err)
	}
	validatorRewards, err := newRewards.DivUint64(uint64(len(epochInfo.Validators)))
	if err != nil {
		return nil, fmt.Errorf("EndBlock, newRewards.DivUint64 error: %v", err)
	}
	allocateSum := NewDecFromBigInt(new(big.Int))
	for _, v := range epochInfo.Validators {
		dec, err := hexutil.Decode(v.PubKey)
		if err != nil {
			return nil, fmt.Errorf("EndBlock, decode pubkey error: %v", err)
		}
		validator, found, err := GetValidator(s, dec)
		if err != nil {
			return nil, fmt.Errorf("EndBlock, GetValidator error: %v", err)
		}
		if found {
			err = allocateRewardsToValidator(s, validator, validatorRewards)
			if err != nil {
				return nil, fmt.Errorf("EndBlock, allocateRewardsToValidator error: %v", err)
			}
			allocateSum, err = allocateSum.Add(validatorRewards)
			if err != nil {
				return nil, fmt.Errorf("EndBlock, allocateSum.Add error: %v", err)
			}
		}
	}

	// update outstanding rewards
	outstanding.Rewards, err = outstanding.Rewards.Add(allocateSum)
	if err != nil {
		return nil, fmt.Errorf("EndBlock, outstanding.Rewards.Add error: %v", err)
	}
	err = setOutstandingRewards(s, outstanding)
	if err != nil {
		return nil, fmt.Errorf("EndBlock, setOutstandingRewards error: %v", err)
	}

	return utils.ByteSuccess, nil
}

func GetGlobalConfig(s *native.NativeContract) ([]byte, error) {
	globalConfig, err := getGlobalConfig(s)
	if err != nil {
		return nil, fmt.Errorf("GetGlobalConfig, getGlobalConfig error: %v", err)
	}

	enc, err := rlp.EncodeToBytes(globalConfig)
	if err != nil {
		return nil, fmt.Errorf("GetGlobalConfig, serialize global config error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetGlobalConfig, enc)
}

func GetCommunityInfo(s *native.NativeContract) ([]byte, error) {
	communityInfo, err := getCommunityInfo(s)
	if err != nil {
		return nil, fmt.Errorf("GetCommunityInfo, getCommunityInfo error: %v", err)
	}

	enc, err := rlp.EncodeToBytes(communityInfo)
	if err != nil {
		return nil, fmt.Errorf("GetCommunityInfo, serialize community info error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetCommunityInfo, enc)
}

func GetCurrentEpochInfo(s *native.NativeContract) ([]byte, error) {
	currentEpochInfo, err := GetCurrentEpochInfoImpl(s)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpochInfo, GetCurrentEpochInfoImpl error: %v", err)
	}
	enc, err := rlp.EncodeToBytes(currentEpochInfo)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpochInfo, serialize current epoch info error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetCurrentEpochInfo, enc)
}
