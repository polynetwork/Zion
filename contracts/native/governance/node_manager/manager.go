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
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/contract"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
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

// the real gas usage of `createValidator`,`changeEpoch`,`endBlock` are 1291500, 5087250 and 343875.
// in order to lower the total gas usage in an entire block, modify them to be 300000 and 200000, 150000.
var (
	gasTable = map[string]uint64{
		MethodCreateValidator:                300000,
		MethodUpdateValidator:                170625,
		MethodUpdateCommission:               126000,
		MethodStake:                          262500,
		MethodUnStake:                        824250,
		MethodWithdraw:                       349125,
		MethodCancelValidator:                333375,
		MethodWithdrawValidator:              328125,
		MethodChangeEpoch:                    200000,
		MethodWithdrawStakeRewards:           286125,
		MethodWithdrawCommission:             149625,
		MethodEndBlock:                       150000,
		MethodGetGlobalConfig:                91875,
		MethodGetCommunityInfo:               81375,
		MethodGetCurrentEpochInfo:            112875,
		MethodGetEpochInfo:                   86625,
		MethodGetAllValidators:               170625,
		MethodGetValidator:                   60375,
		MethodGetStakeInfo:                   76125,
		MethodGetUnlockingInfo:               357000,
		MethodGetStakeStartingInfo:           73500,
		MethodGetAccumulatedCommission:       63000,
		MethodGetValidatorSnapshotRewards:    128625,
		MethodGetValidatorAccumulatedRewards: 60375,
		MethodGetValidatorOutstandingRewards: 65625,
		MethodGetTotalPool:                   60375,
		MethodGetOutstandingRewards:          60375,
		MethodGetStakeRewards:                128625,
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
	s.Register(MethodGetEpochInfo, GetEpochInfo)
	s.Register(MethodGetAllValidators, GetAllValidators)
	s.Register(MethodGetValidator, GetValidator)
	s.Register(MethodGetStakeInfo, GetStakeInfo)
	s.Register(MethodGetUnlockingInfo, GetUnlockingInfo)
	s.Register(MethodGetStakeStartingInfo, GetStakeStartingInfo)
	s.Register(MethodGetAccumulatedCommission, GetAccumulatedCommission)
	s.Register(MethodGetValidatorSnapshotRewards, GetValidatorSnapshotRewards)
	s.Register(MethodGetValidatorAccumulatedRewards, GetValidatorAccumulatedRewards)
	s.Register(MethodGetValidatorOutstandingRewards, GetValidatorOutstandingRewards)
	s.Register(MethodGetTotalPool, GetTotalPool)
	s.Register(MethodGetOutstandingRewards, GetOutstandingRewards)
	s.Register(MethodGetStakeRewards, GetStakeRewards)
}

func CreateValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	height := s.ContractRef().BlockHeight()
	caller := ctx.Caller

	params := &CreateValidatorParam{}
	if err := utils.UnpackMethod(ABI, MethodCreateValidator, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("CreateValidator, unpack params error: %v", err)
	}

	// check consensus address
	if params.ConsensusAddress == common.EmptyAddress {
		return nil, fmt.Errorf("CreateValidator， invalid consensus address")
	}
	if params.SignerAddress == common.EmptyAddress {
		return nil, fmt.Errorf("CreateValidator， invalid signer address")
	}
	if params.ProposalAddress == common.EmptyAddress {
		return nil, fmt.Errorf("CreateValidator， invalid proposalAddress")
	}

	// check commission
	globalConfig, err := GetGlobalConfigImpl(s)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, GetGlobalConfig error: %v", err)
	}
	if params.Commission.Sign() == -1 {
		return nil, fmt.Errorf("CreateValidator, commission must be positive")
	}
	if params.Commission.Cmp(new(big.Int).SetUint64(10000)) == 1 {
		return nil, fmt.Errorf("CreateValidator, commission can not greater than 100 percent")
	}

	// check desc
	if len(params.Desc) > MaxDescLength {
		return nil, fmt.Errorf("CreateValidator, desc length more than limit %d", MaxDescLength)
	}

	// check to see if the pubkey has been registered before
	_, found, err := getValidator(s, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, getValidator error: %v", err)
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
		ConsensusAddress: params.ConsensusAddress,
		SignerAddress:    params.SignerAddress,
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
	err = addToAllValidators(s, validator.ConsensusAddress)
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

	err = s.AddNotify(ABI, []string{CREATE_VALIDATOR_EVENT}, params.ConsensusAddress.Hex())
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodCreateValidator, true)
}

func UpdateValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	params := &UpdateValidatorParam{}
	if err := utils.UnpackMethod(ABI, MethodUpdateValidator, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("UpdateValidator, unpack params error: %v", err)
	}

	validator, found, err := getValidator(s, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("UpdateValidator, get validator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("UpdateValidator, can not found record")
	}
	if validator.StakeAddress != caller {
		return nil, fmt.Errorf("UpdateValidator, stake address is not caller")
	}

	if params.ProposalAddress != common.EmptyAddress {
		validator.ProposalAddress = params.ProposalAddress
	}
	if params.SignerAddress != common.EmptyAddress {
		validator.SignerAddress = params.SignerAddress
	}

	if params.Desc != "" {
		// check desc
		if len(params.Desc) > MaxDescLength {
			return nil, fmt.Errorf("UpdateValidator, desc length more than limit %d", MaxDescLength)
		}
		validator.Desc = params.Desc
	}

	err = setValidator(s, validator)
	if err != nil {
		return nil, fmt.Errorf("UpdateValidator, setValidator error: %v", err)
	}

	err = s.AddNotify(ABI, []string{UPDATE_VALIDATOR_EVENT}, params.ConsensusAddress.Hex())
	if err != nil {
		return nil, fmt.Errorf("UpdateValidator, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodUpdateValidator, true)
}

func UpdateCommission(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	height := s.ContractRef().BlockHeight()
	caller := ctx.Caller

	params := &UpdateCommissionParam{}
	if err := utils.UnpackMethod(ABI, MethodUpdateCommission, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("UpdateCommission, unpack params error: %v", err)
	}

	validator, found, err := getValidator(s, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("UpdateCommission, get validator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("UpdateCommission, can not found record")
	}
	if validator.StakeAddress != caller {
		return nil, fmt.Errorf("UpdateCommission, stake address is not caller")
	}
	globalConfig, err := GetGlobalConfigImpl(s)
	if err != nil {
		return nil, fmt.Errorf("UpdateCommission, GetGlobalConfig error: %v", err)
	}

	// check commission
	if params.Commission.Sign() == -1 {
		return nil, fmt.Errorf("UpdateCommission, commission must be positive")
	}
	if params.Commission.Cmp(PercentDecimal) == 1 {
		return nil, fmt.Errorf("UpdateCommission, commission can not more than 100 percent")
	}
	// abs(old commission - new commission)
	if new(big.Int).Abs(new(big.Int).Sub(validator.Commission.Rate.BigInt(), params.Commission)).Cmp(globalConfig.MaxCommissionChange) == 1 {
		return nil, fmt.Errorf("UpdateCommission, commission change can not greater than globalConfig.MaxCommissionChange: %s",
			globalConfig.MaxCommissionChange.String())
	}
	if height.Cmp(new(big.Int).Add(validator.Commission.UpdateHeight, globalConfig.BlockPerEpoch)) < 0 {
		return nil, fmt.Errorf("UpdateCommission, commission can not changed in one epoch twice")
	}

	validator.Commission = &Commission{Rate: NewDecFromBigInt(params.Commission), UpdateHeight: height}

	err = setValidator(s, validator)
	if err != nil {
		return nil, fmt.Errorf("UpdateCommission, setValidator error: %v", err)
	}

	err = s.AddNotify(ABI, []string{UPDATE_COMMISSION_EVENT}, params.ConsensusAddress.Hex())
	if err != nil {
		return nil, fmt.Errorf("UpdateCommission, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodUpdateCommission, true)
}

func Stake(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	params := &StakeParam{}
	if err := utils.UnpackMethod(ABI, MethodStake, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("Stake, unpack params error: %v", err)
	}
	if params.Amount.Sign() <= 0 {
		return nil, fmt.Errorf("Stake, amount must be positive")
	}
	amount := NewDecFromBigInt(params.Amount)

	// check to see if the pubkey has been registered
	validator, found, err := getValidator(s, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("Stake, getValidator error: %v", err)
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
		validator.TotalStake, err = validator.TotalStake.Add(amount)
		if err != nil {
			return nil, fmt.Errorf("Stake, validator.TotalStake.Add error: %v", err)
		}
	} else {
		validator.TotalStake, err = validator.TotalStake.Add(amount)
		if err != nil {
			return nil, fmt.Errorf("Stake, validator.TotalStake.Add error: %v", err)
		}
		maxTotalStake, err := validator.SelfStake.Mul(MaxStakeRate)
		if err != nil {
			return nil, fmt.Errorf("Stake, validator.SelfStake.Mul error: %v", err)
		}
		if validator.TotalStake.GT(maxTotalStake) {
			return nil, fmt.Errorf("Stake, stake is more than max stake")
		}
	}
	err = setValidator(s, validator)
	if err != nil {
		return nil, fmt.Errorf("Stake, setValidator error: %v", err)
	}

	err = s.AddNotify(ABI, []string{STAKE_EVENT}, params.ConsensusAddress.Hex(), params.Amount.String())
	if err != nil {
		return nil, fmt.Errorf("Stake, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodStake, true)
}

func UnStake(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	params := &UnStakeParam{}
	if err := utils.UnpackMethod(ABI, MethodUnStake, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("UnStake, unpack params error: %v", err)
	}
	if params.Amount.Sign() <= 0 {
		return nil, fmt.Errorf("UnStake, amount must be positive")
	}
	amount := NewDecFromBigInt(params.Amount)

	// check to see if the pubkey has been registered
	validator, found, err := getValidator(s, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("UnStake, getValidator error: %v", err)
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
		delValidator(s, params.ConsensusAddress)
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

	err = s.AddNotify(ABI, []string{UNSTAKE_EVENT}, params.ConsensusAddress.Hex(), params.Amount.String())
	if err != nil {
		return nil, fmt.Errorf("UnStake, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodUnStake, true)
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
		err = contract.NativeTransfer(s, this, caller, amount.BigInt())
		if err != nil {
			return nil, fmt.Errorf("Withdraw, nativeTransfer error: %v", err)
		}
	} else {
		return nil, fmt.Errorf("Withdraw, no asset to withdraw")
	}

	err = s.AddNotify(ABI, []string{WITHDRAW_EVENT}, caller.Hex(), amount.BigInt().String())
	if err != nil {
		return nil, fmt.Errorf("Withdraw, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodWithdraw, true)
}

func CancelValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller
	height := s.ContractRef().BlockHeight()

	params := &CancelValidatorParam{}
	if err := utils.UnpackMethod(ABI, MethodCancelValidator, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("CancelValidator, unpack params error: %v", err)
	}

	globalConfig, err := GetGlobalConfigImpl(s)
	if err != nil {
		return nil, fmt.Errorf("CancelValidator, GetGlobalConfig error: %v", err)
	}
	allValidators, err := getAllValidators(s)
	if err != nil {
		return nil, fmt.Errorf("CancelValidator, getAllValidators error: %v", err)
	}
	if uint64(len(allValidators.AllValidators)) <= globalConfig.ConsensusValidatorNum {
		return nil, fmt.Errorf("CancelValidator, validator num is less than consensus node num error: %v", err)
	}

	validator, found, err := getValidator(s, params.ConsensusAddress)
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
	case validator.IsUnlocking(height), validator.IsUnlocked(height):
		validator.Status = Remove
	default:
		return nil, fmt.Errorf("CancelValidator, unsupported validator status")
	}
	err = setValidator(s, validator)
	if err != nil {
		return nil, fmt.Errorf("CancelValidator, setValidator error: %v", err)
	}
	err = removeFromAllValidators(s, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("CancelValidator, removeFromAllValidators error: %v", err)
	}

	err = s.AddNotify(ABI, []string{CANCEL_VALIDATOR_EVENT}, params.ConsensusAddress.Hex())
	if err != nil {
		return nil, fmt.Errorf("CancelValidator, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodCancelValidator, true)
}

func WithdrawValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller
	height := s.ContractRef().BlockHeight()

	params := &WithdrawValidatorParam{}
	if err := utils.UnpackMethod(ABI, MethodWithdrawValidator, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("WithdrawValidator, unpack params error: %v", err)
	}
	validator, found, err := getValidator(s, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("WithdrawValidator, getValidator error: %v", err)
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

	amount := validator.SelfStake
	// unStake native token
	err = unStake(s, caller, amount, validator)
	if err != nil {
		return nil, fmt.Errorf("WithdrawValidator, unStake error: %v", err)
	}

	_, err = withdrawCommission(s, caller, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("WithdrawValidator, withdrawCommission error: %v", err)
	}

	validator.TotalStake, err = validator.TotalStake.Sub(amount)
	if err != nil {
		return nil, fmt.Errorf("WithdrawValidator, validator.TotalStake.Sub error: %v", err)
	}
	validator.SelfStake = NewDecFromBigInt(new(big.Int))
	if validator.TotalStake.IsZero() {
		delValidator(s, params.ConsensusAddress)
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

	err = s.AddNotify(ABI, []string{WITHDRAW_VALIDATOR_EVENT}, params.ConsensusAddress.Hex(), amount.BigInt().String())
	if err != nil {
		return nil, fmt.Errorf("CancelValidator, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodWithdrawValidator, true)
}

func ChangeEpoch(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	if ctx.Caller != s.ContractRef().TxOrigin() || ctx.Caller != utils.SystemTxSender {
		return nil, fmt.Errorf("SystemTx authority failed")
	}

	endHeight := s.ContractRef().BlockHeight()
	startHeight := new(big.Int).Add(endHeight, common.Big1)

	currentEpochInfo, err := GetCurrentEpochInfoImpl(s)
	if err != nil {
		return nil, fmt.Errorf("ChangeEpoch, GetCurrentEpochInfoImpl error: %v", err)
	}
	globalConfig, err := GetGlobalConfigImpl(s)
	if err != nil {
		return nil, fmt.Errorf("ChangeEpoch, GetGlobalConfigImpl error: %v", err)
	}

	// anyone can call this if height reaches
	if startHeight.Cmp(currentEpochInfo.EndHeight) != 0 {
		return nil, fmt.Errorf("ChangeEpoch, block height does not reach, current epoch end at %s",
			currentEpochInfo.EndHeight.String())
	}

	epochInfo := &EpochInfo{
		ID:          new(big.Int).Add(currentEpochInfo.ID, common.Big1),
		Validators:  make([]common.Address, 0, globalConfig.ConsensusValidatorNum),
		Signers:     make([]common.Address, 0, globalConfig.ConsensusValidatorNum),
		Voters:      make([]common.Address, 0, globalConfig.VoterValidatorNum),
		Proposers:   make([]common.Address, 0, globalConfig.ConsensusValidatorNum),
		StartHeight: startHeight,
		EndHeight:   new(big.Int).Add(startHeight, globalConfig.BlockPerEpoch),
	}
	// get all validators
	allValidators, err := getAllValidators(s)
	if err != nil {
		return nil, fmt.Errorf("ChangeEpoch, getAllValidators error: %v", err)
	}
	if uint64(len(allValidators.AllValidators)) < globalConfig.ConsensusValidatorNum {
		epochInfo.Validators = currentEpochInfo.Validators
		epochInfo.Signers = currentEpochInfo.Signers
		epochInfo.Voters = currentEpochInfo.Voters
		epochInfo.Proposers = currentEpochInfo.Proposers
	} else {
		validatorList := make([]*Validator, 0, len(allValidators.AllValidators))
		for _, v := range allValidators.AllValidators {
			validator, found, err := getValidator(s, v)
			if err != nil {
				return nil, fmt.Errorf("ChangeEpoch, getValidator error: %v", err)
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

			epochInfo.Validators = append(epochInfo.Validators, validator.ConsensusAddress)
			epochInfo.Signers = append(epochInfo.Signers, validator.SignerAddress)
			epochInfo.Proposers = append(epochInfo.Proposers, validator.ProposalAddress)
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
			epochInfo.Voters = append(epochInfo.Voters, validator.SignerAddress)
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
	return utils.PackOutputs(ABI, MethodChangeEpoch, true)
}

func WithdrawStakeRewards(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	params := &WithdrawStakeRewardsParam{}
	if err := utils.UnpackMethod(ABI, MethodWithdrawStakeRewards, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, unpack params error: %v", err)
	}

	validator, found, err := getValidator(s, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, node_manager.getValidator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("WithdrawStakeRewards, validator not found")
	}
	stakeInfo, found, err := getStakeInfo(s, caller, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, getStakeInfo error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("WithdrawStakeRewards, stake info not found")
	}

	rewards, err := withdrawStakeRewards(s, validator, stakeInfo)
	if err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, withdrawStakeRewards error: %v", err)
	}

	// reinitialize the delegation
	err = initializeStake(s, stakeInfo, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, initializeStake error: %v", err)
	}

	err = s.AddNotify(ABI, []string{WITHDRAW_STAKE_REWARDS_EVENT}, params.ConsensusAddress.Hex(), caller.Hex(), rewards.BigInt().String())
	if err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodWithdrawStakeRewards, true)
}

func WithdrawCommission(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	params := &WithdrawCommissionParam{}
	if err := utils.UnpackMethod(ABI, MethodWithdrawCommission, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("WithdrawCommission, unpack params error: %v", err)
	}
	validator, found, err := getValidator(s, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("WithdrawCommission, getValidator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("WithdrawCommission, can not find validator")
	}
	if validator.StakeAddress != caller {
		return nil, fmt.Errorf("WithdrawCommission, caller is not stake address")
	}

	commission, err := withdrawCommission(s, caller, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("WithdrawCommission, withdrawCommission error: %v", err)
	}
	err = setAccumulatedCommission(s, params.ConsensusAddress, &AccumulatedCommission{NewDecFromBigInt(new(big.Int))})
	if err != nil {
		return nil, fmt.Errorf("WithdrawCommission, setAccumulatedCommission error: %v", err)
	}

	err = s.AddNotify(ABI, []string{WITHDRAW_COMMISSION_EVENT}, params.ConsensusAddress.Hex(), commission.BigInt().String())
	if err != nil {
		return nil, fmt.Errorf("WithdrawCommission, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodWithdrawCommission, true)
}

func EndBlock(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	if ctx.Caller != s.ContractRef().TxOrigin() || ctx.Caller != utils.SystemTxSender {
		return nil, fmt.Errorf("SystemTx authority failed")
	}

	// contract balance = totalpool + outstanding + reward
	balance := NewDecFromBigInt(s.StateDB().GetBalance(this))

	totalPool, err := getTotalPool(s)
	if err != nil {
		return nil, fmt.Errorf("EndBlock, getTotalPool error: %v", err)
	}
	outstanding, err := getOutstandingRewards(s)
	if err != nil {
		return nil, fmt.Errorf("EndBlock, getOutstandingRewards error: %v", err)
	}

	// cal rewards
	temp, err := outstanding.Rewards.Add(totalPool.TotalPool)
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
		validator, found, err := getValidator(s, v)
		if err != nil {
			return nil, fmt.Errorf("EndBlock, getValidator error: %v", err)
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

	return utils.PackOutputs(ABI, MethodEndBlock, true)
}

func GetGlobalConfig(s *native.NativeContract) ([]byte, error) {
	globalConfig, err := GetGlobalConfigImpl(s)
	if err != nil {
		return nil, fmt.Errorf("GetGlobalConfig, GetGlobalConfigImpl error: %v", err)
	}

	enc, err := rlp.EncodeToBytes(globalConfig)
	if err != nil {
		return nil, fmt.Errorf("GetGlobalConfig, serialize global config error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetGlobalConfig, enc)
}

func GetCommunityInfo(s *native.NativeContract) ([]byte, error) {
	communityInfo, err := GetCommunityInfoImpl(s)
	if err != nil {
		return nil, fmt.Errorf("GetCommunityInfo, GetCommunityInfoImpl error: %v", err)
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

func GetEpochInfo(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	params := &GetEpochInfoParam{}
	if err := utils.UnpackMethod(ABI, MethodGetEpochInfo, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("GetEpochInfo, unpack params error: %v", err)
	}

	epochInfo, err := getEpochInfo(s, params.ID)
	if err != nil {
		return nil, fmt.Errorf("GetEpochInfo, getEpochInfo error: %v", err)
	}
	enc, err := rlp.EncodeToBytes(epochInfo)
	if err != nil {
		return nil, fmt.Errorf("GetEpochInfo, serialize epoch info error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetEpochInfo, enc)
}

func GetAllValidators(s *native.NativeContract) ([]byte, error) {
	allValidators, err := getAllValidators(s)
	if err != nil {
		return nil, fmt.Errorf("GetAllValidators, getAllValidators error: %v", err)
	}
	enc, err := rlp.EncodeToBytes(allValidators)
	if err != nil {
		return nil, fmt.Errorf("GetAllValidators, serialize all validators error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetAllValidators, enc)
}

func GetValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	params := &GetValidatorParam{}
	if err := utils.UnpackMethod(ABI, MethodGetValidator, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("GetValidator, unpack params error: %v", err)
	}

	validator, found, err := getValidator(s, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("GetValidator, getValidator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("GetValidator, no record")
	}
	enc, err := rlp.EncodeToBytes(validator)
	if err != nil {
		return nil, fmt.Errorf("GetValidator, serialize validator error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetValidator, enc)
}

func GetStakeInfo(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	params := &GetStakeInfoParam{}
	if err := utils.UnpackMethod(ABI, MethodGetStakeInfo, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("GetStakeInfo, unpack params error: %v", err)
	}

	stakeInfo, found, err := getStakeInfo(s, params.StakeAddress, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("GetStakeInfo, getStakeInfo error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("GetStakeInfo, no record")
	}
	enc, err := rlp.EncodeToBytes(stakeInfo)
	if err != nil {
		return nil, fmt.Errorf("GetStakeInfo, serialize stakeInfo error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetStakeInfo, enc)
}

func GetUnlockingInfo(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	params := &GetUnlockingInfoParam{}
	if err := utils.UnpackMethod(ABI, MethodGetUnlockingInfo, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("GetUnlockingInfo, unpack params error: %v", err)
	}

	unlockingInfo, err := getUnlockingInfo(s, params.StakeAddress)
	if err != nil {
		return nil, fmt.Errorf("GetUnlockingInfo, getUnlockingInfo error: %v", err)
	}
	enc, err := rlp.EncodeToBytes(unlockingInfo)
	if err != nil {
		return nil, fmt.Errorf("GetUnlockingInfo, serialize unlockingInfo error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetUnlockingInfo, enc)
}

func GetStakeStartingInfo(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	params := &GetStakeStartingInfoParam{}
	if err := utils.UnpackMethod(ABI, MethodGetStakeStartingInfo, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("GetStakeStartingInfo, unpack params error: %v", err)
	}

	stakeStartingInfo, err := getStakeStartingInfo(s, params.StakeAddress, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("GetStakeStartingInfo, getStakeStartingInfo error: %v", err)
	}
	enc, err := rlp.EncodeToBytes(stakeStartingInfo)
	if err != nil {
		return nil, fmt.Errorf("GetStakeStartingInfo, serialize stakeStartingInfo error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetStakeStartingInfo, enc)
}

func GetAccumulatedCommission(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	params := &GetAccumulatedCommissionParam{}
	if err := utils.UnpackMethod(ABI, MethodGetAccumulatedCommission, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("GetAccumulatedCommission, unpack params error: %v", err)
	}

	accumulatedCommission, err := getAccumulatedCommission(s, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("GetAccumulatedCommission, getAccumulatedCommission error: %v", err)
	}
	enc, err := rlp.EncodeToBytes(accumulatedCommission)
	if err != nil {
		return nil, fmt.Errorf("GetAccumulatedCommission, serialize accumulatedCommission error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetAccumulatedCommission, enc)
}

func GetValidatorSnapshotRewards(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	params := &GetValidatorSnapshotRewardsParam{}
	if err := utils.UnpackMethod(ABI, MethodGetValidatorSnapshotRewards, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("GetValidatorSnapshotRewards, unpack params error: %v", err)
	}

	validatorSnapshotRewards, err := getValidatorSnapshotRewards(s, params.ConsensusAddress, params.Period)
	if err != nil {
		return nil, fmt.Errorf("GetValidatorSnapshotRewards, getValidatorSnapshotRewards error: %v", err)
	}
	enc, err := rlp.EncodeToBytes(validatorSnapshotRewards)
	if err != nil {
		return nil, fmt.Errorf("GetValidatorSnapshotRewards, serialize validatorSnapshotRewards error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetValidatorSnapshotRewards, enc)
}

func GetValidatorAccumulatedRewards(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	params := &GetValidatorAccumulatedRewardsParam{}
	if err := utils.UnpackMethod(ABI, MethodGetValidatorAccumulatedRewards, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("GetValidatorAccumulatedRewards, unpack params error: %v", err)
	}

	validatorAccumulatedRewards, err := getValidatorAccumulatedRewards(s, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("GetValidatorAccumulatedRewards, getValidatorAccumulatedRewards error: %v", err)
	}
	enc, err := rlp.EncodeToBytes(validatorAccumulatedRewards)
	if err != nil {
		return nil, fmt.Errorf("GetValidatorAccumulatedRewards, serialize validatorAccumulatedRewards error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetValidatorAccumulatedRewards, enc)
}

func GetValidatorOutstandingRewards(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	params := &GetValidatorOutstandingRewardsParam{}
	if err := utils.UnpackMethod(ABI, MethodGetValidatorOutstandingRewards, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("GetValidatorOutstandingRewards, unpack params error: %v", err)
	}

	validatorOutstandingRewards, err := getValidatorOutstandingRewards(s, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("GetValidatorOutstandingRewards, getValidatorOutstandingRewards error: %v", err)
	}
	enc, err := rlp.EncodeToBytes(validatorOutstandingRewards)
	if err != nil {
		return nil, fmt.Errorf("GetValidatorOutstandingRewards, serialize validatorOutstandingRewards error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetValidatorOutstandingRewards, enc)
}

func GetTotalPool(s *native.NativeContract) ([]byte, error) {
	totalPool, err := getTotalPool(s)
	if err != nil {
		return nil, fmt.Errorf("GetTotalPool, getTotalPool error: %v", err)
	}
	enc, err := rlp.EncodeToBytes(totalPool)
	if err != nil {
		return nil, fmt.Errorf("GetTotalPool, serialize totalPool error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetTotalPool, enc)
}

func GetOutstandingRewards(s *native.NativeContract) ([]byte, error) {
	outstandingRewards, err := getOutstandingRewards(s)
	if err != nil {
		return nil, fmt.Errorf("GetOutstandingRewards, getOutstandingRewards error: %v", err)
	}
	enc, err := rlp.EncodeToBytes(outstandingRewards)
	if err != nil {
		return nil, fmt.Errorf("GetOutstandingRewards, serialize outstandingRewards error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodGetOutstandingRewards, enc)
}

func GetStakeRewards(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	params := &GetStakeRewardsParam{}
	if err := utils.UnpackMethod(ABI, MethodGetStakeRewards, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("GetStakeRewards, unpack params error: %v", err)
	}

	validator, found, err := getValidator(s, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("GetStakeRewards, node_manager.getValidator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("GetStakeRewards, validator not found")
	}

	// fetch current rewards
	validatorAccumulatedRewards, err := getValidatorAccumulatedRewards(s, validator.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("GetStakeRewards, getValidatorAccumulatedRewards error: %v", err)
	}
	now, err := getValidatorSnapshotRewards(s, params.ConsensusAddress, validatorAccumulatedRewards.Period-1)
	if err != nil {
		return nil, fmt.Errorf("GetStakeRewards, getValidatorSnapshotRewards now error: %v", err)
	}
	// calculate current ratio
	// mul decimal
	ratio, err := validatorAccumulatedRewards.Rewards.DivWithTokenDecimal(validator.TotalStake)
	if err != nil {
		return nil, fmt.Errorf("GetStakeRewards, validatorAccumulatedRewards.Rewards.DivWithTokenDecimal error: %v", err)
	}
	newRatio, err := now.AccumulatedRewardsRatio.Add(ratio)
	if err != nil {
		return nil, fmt.Errorf("GetStakeRewards, now.AccumulatedRewardsRatio.Add error: %v", err)
	}
	// fetch starting info for delegation
	startingInfo, err := getStakeStartingInfo(s, params.StakeAddress, params.ConsensusAddress)
	if err != nil {
		return nil, fmt.Errorf("GetStakeRewards, getStakeStartingInfo error: %v", err)
	}

	startPeriod := startingInfo.StartPeriod
	stake := startingInfo.Stake

	// return staking * (ending - starting)
	starting, err := getValidatorSnapshotRewards(s, params.ConsensusAddress, startPeriod)
	if err != nil {
		return nil, fmt.Errorf("GetStakeRewards, getValidatorSnapshotRewards start error: %v", err)
	}
	difference, err := newRatio.Sub(starting.AccumulatedRewardsRatio)
	if err != nil {
		return nil, fmt.Errorf("GetStakeRewards, newRatio.Sub error: %v", err)
	}
	rewards, err := difference.MulWithTokenDecimal(stake)
	if err != nil {
		return nil, fmt.Errorf("GetStakeRewards, difference.MulWithTokenDecimal error: %v", err)
	}

	enc, err := rlp.EncodeToBytes(&StakeRewards{Rewards: rewards})
	if err != nil {
		return nil, fmt.Errorf("GetStakeRewards, serialize stake rewards error: %v", err)
	}

	return utils.PackOutputs(ABI, MethodGetStakeRewards, enc)
}
