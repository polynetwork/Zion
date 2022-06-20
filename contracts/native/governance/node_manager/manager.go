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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
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
	native.Contracts[this] = RegisterNodeManagerContract
}

func RegisterNodeManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodCreateValidator, CreateValidator)
	s.Register(MethodUpdateValidator, UpdateValidator)
	s.Register(MethodStake, Stake)
	s.Register(MethodUnStake, UnStake)
	s.Register(MethodWithdraw, Withdraw)
	s.Register(MethodCancelValidator, CancelValidator)
	s.Register(MethodWithdrawValidator, WithdrawValidator)
	s.Register(MethodChangeEpoch, ChangeEpoch)

	s.Register(MethodWithdrawStakeRewards, WithdrawStakeRewards)
	s.Register(MethodWithdrawCommission, WithdrawCommission)
	s.Register(MethodBeginBlock, BeginBlock)
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
	if addr == caller || addr == params.ProposalAddress {
		return nil, fmt.Errorf("CreateValidator，stake, consensus and proposal address can not be duplicate")
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
		StakeAddress:     caller,
		ConsensusPubkey:  params.ConsensusPubkey,
		ConsensusAddress: addr,
		ProposalAddress:  params.ProposalAddress,
		Commission:       &Commission{Rate: params.Commission, UpdateHeight: height},
		Status:           Unlock,
		Jailed:           false,
		UnlockHeight:     common.Big0,
		TotalStake:       params.InitStake,
		SelfStake:        params.InitStake,
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

	// deposit native token
	err = deposit(s, caller, params.InitStake, validator)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, deposit error: %v", err)
	}

	// call distrubute hook
	err = AfterValidatorCreated(s, validator)
	if err != nil {
		return nil, fmt.Errorf("CreateValidator, distribute.AfterValidatorCreated error: %v", err)
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
		if params.Commission.Cmp(new(big.Int).SetUint64(100)) == 1 {
			return nil, fmt.Errorf("UpdateValidator, commission can not more than 100")
		}
		if params.Commission.Cmp(globalConfig.MaxCommission) == 1 {
			return nil, fmt.Errorf("UpdateValidator, commission can not greater than globalConfig.MaxCommission: %s",
				globalConfig.MaxCommission.String())
		}
		if height.Cmp(new(big.Int).Add(validator.Commission.UpdateHeight, globalConfig.BlockPerEpoch)) < 0 {
			return nil, fmt.Errorf("UpdateValidator, commission can not changed in one epoch twice")
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

	// update validator
	if validator.StakeAddress == caller {
		validator.SelfStake = new(big.Int).Add(validator.SelfStake, params.Amount)
	} else {
		validator.TotalStake = new(big.Int).Add(validator.TotalStake, params.Amount)
	}
	err = setValidator(s, validator)
	if err != nil {
		return nil, fmt.Errorf("Stake, setValidator error: %v", err)
	}

	err = s.AddNotify(ABI, []string{MethodStake}, params.ConsensusPubkey, params.Amount.String())
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

	// check to see if the pubkey has been registered
	validator, found, err := GetValidator(s, dec)
	if err != nil {
		return nil, fmt.Errorf("UnStake, GetValidator error: %v", err)
	}
	if !found {
		return nil, fmt.Errorf("UnStake, validator is not exist")
	}

	// update validator
	if validator.StakeAddress == caller {
		return nil, fmt.Errorf("UnStake, stake address can not unstake")
	} else {
		if validator.TotalStake.Cmp(params.Amount) == -1 {
			return nil, fmt.Errorf("UnStake, total stake of validator is less than amount")
		}
		validator.TotalStake = new(big.Int).Sub(validator.TotalStake, params.Amount)
	}
	if validator.TotalStake == common.Big0 && validator.SelfStake == common.Big0 {
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

	// unStake native token
	err = unStake(s, caller, params.Amount, validator)
	if err != nil {
		return nil, fmt.Errorf("UnStake, unStake error: %v", err)
	}

	err = s.AddNotify(ABI, []string{MethodUnStake}, params.ConsensusPubkey, params.Amount.String())
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

	err = s.AddNotify(ABI, []string{MethodWithdraw}, caller.Hex(), amount.String())
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

	globalConfig, err := GetGlobalConfig(s)
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
		err = withdrawLockPool(s, validator.TotalStake)
		if err != nil {
			return nil, fmt.Errorf("CancelValidator, withdrawLockPool error: %v", err)
		}
		err = depositUnlockPool(s, validator.TotalStake)
		if err != nil {
			return nil, fmt.Errorf("CancelValidator, depositUnlockPool error: %v", err)
		}
		validator.Status = Remove
		validator.UnlockHeight = new(big.Int).Add(height, globalConfig.BlockPerEpoch)
		err = setValidator(s, validator)
		if err != nil {
			return nil, fmt.Errorf("CancelValidator, setValidator error: %v", err)
		}
	case validator.IsUnlocking(height):
		validator.Status = Remove
	case validator.IsUnlocked(height):
		validator.Status = Remove
	default:
		return nil, fmt.Errorf("CancelValidator, unsupported validator status")
	}
	err = removeFromAllValidators(s, params.ConsensusPubkey)
	if err != nil {
		return nil, fmt.Errorf("CancelValidator, removeFromAllValidators error: %v", err)
	}

	err = s.AddNotify(ABI, []string{MethodCancelValidator}, params.ConsensusPubkey)
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

	_, err = withdrawCommission(s, caller, dec)
	if err != nil {
		return nil, fmt.Errorf("WithdrawValidator, withdrawCommission error: %v", err)
	}
	delAccumulatedCommission(s, dec)

	validator.SelfStake = common.Big0
	if validator.TotalStake == common.Big0 {
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

	err = s.AddNotify(ABI, []string{MethodWithdrawValidator}, params.ConsensusPubkey, validator.SelfStake.String())
	if err != nil {
		return nil, fmt.Errorf("CancelValidator, AddNotify error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func ChangeEpoch(s *native.NativeContract) ([]byte, error) {
	height := s.ContractRef().BlockHeight()

	currentEpochInfo, err := GetCurrentEpochInfo(s)
	if err != nil {
		return nil, fmt.Errorf("ChangeEpoch, GetCurrentEpochInfo error: %v", err)
	}
	globalConfig, err := GetGlobalConfig(s)

	// anyone can call this if height reaches
	if new(big.Int).Sub(height, currentEpochInfo.StartHeight).Cmp(globalConfig.BlockPerEpoch) == -1 {
		return nil, fmt.Errorf("ChangeEpoch, block height does not reach, current epoch start at %s",
			currentEpochInfo.StartHeight.String())
	}

	// get all validators
	allValidators, err := GetAllValidators(s)
	if err != nil {
		return nil, fmt.Errorf("ChangeEpoch, GetAllValidators error: %v", err)
	}
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
		return validatorList[i].TotalStake.Cmp(validatorList[j].TotalStake) == 1
	})
	epochInfo := &EpochInfo{
		ID:          new(big.Int).Add(currentEpochInfo.ID, common.Big1),
		Validators:  make([]*Peer, 0, globalConfig.ConsensusValidatorNum),
		Voters:      make([]*Peer, 0, globalConfig.VoterValidatorNum),
		StartHeight: height,
	}
	// update validator status
	for i := 0; uint64(i) < globalConfig.ConsensusValidatorNum; i++ {
		validator := validatorList[i]
		switch {
		case validator.IsLocked():
		case validator.IsUnlocking(height), validator.IsUnlocked(height):
			validator.Status = Lock
		}

		peer := &Peer{
			PubKey:  validator.ConsensusPubkey,
			Address: validator.ConsensusAddress,
		}
		epochInfo.Validators = append(epochInfo.Validators, peer)
	}
	for i := globalConfig.ConsensusValidatorNum; i < uint64(len(validatorList)); i++ {
		validator := validatorList[i]
		switch {
		case validator.IsLocked():
			validator.Status = Unlock
			validator.UnlockHeight = new(big.Int).Add(height, globalConfig.BlockPerEpoch)
		case validator.IsUnlocking(height), validator.IsUnlocked(height):
		}
	}
	//update voters
	for i := 0; uint64(i) < globalConfig.VoterValidatorNum; i++ {
		validator := validatorList[i]
		peer := &Peer{
			PubKey:  validator.ConsensusPubkey,
			Address: validator.ConsensusAddress,
		}
		epochInfo.Validators = append(epochInfo.Voters, peer)
	}

	// update epoch info
	err = setCurrentEpochInfo(s, epochInfo)
	if err != nil {
		return nil, fmt.Errorf("ChangeEpoch, setEpochInfo error: %v", err)
	}

	err = s.AddNotify(ABI, []string{MethodChangeEpoch}, epochInfo.ID.String())
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
	if caller != validator.StakeAddress {
		return nil, fmt.Errorf("WithdrawStakeRewards, caller is not stake address error: %v", err)
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
		return nil, fmt.Errorf("WithdrawStakeRewards, withdrawDelegationRewards error: %v", err)
	}

	// reinitialize the delegation
	err = initializeStake(s, stakeInfo, dec)
	if err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, initializeStake error: %v", err)
	}

	err = s.AddNotify(ABI, []string{MethodWithdrawStakeRewards}, rewards.String())
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
		return nil, fmt.Errorf("WithdrawValidator, withdrawCommission error: %v", err)
	}
	err = setAccumulatedCommission(s, dec, &AccumulatedCommission{common.Big0})
	if err != nil {
		return nil, fmt.Errorf("WithdrawValidator, setAccumulatedCommission error: %v", err)
	}

	err = s.AddNotify(ABI, []string{MethodWithdrawCommission}, commission.String())
	if err != nil {
		return nil, fmt.Errorf("WithdrawStakeRewards, AddNotify error: %v", err)
	}
	return utils.ByteSuccess, nil
}

func BeginBlock(s *native.NativeContract) ([]byte, error) {
	// contract balance = lockpool + unlockpool + outstanding + new block reward
	balance := s.StateDB().GetBalance(this)

	lockPool, err := GetLockPool(s)
	if err != nil {
		return nil, fmt.Errorf("BeginBlock, GetLockPool error: %v", err)
	}
	unlockPool, err := GetUnlockPool(s)
	if err != nil {
		return nil, fmt.Errorf("BeginBlock, GetUnlockPool error: %v", err)
	}
	outstanding, err := GetOutstandingRewards(s)
	if err != nil {
		return nil, fmt.Errorf("BeginBlock, GetOutstandingRewards error: %v", err)
	}

	// cal rewards
	temp := new(big.Int).Add(outstanding.Rewards, new(big.Int).Add(lockPool, unlockPool))
	newRewards := new(big.Int).Sub(balance, temp)
	if newRewards.Sign() < 0 {
		panic("new block rewards is negative")
	}

	epochInfo, err := GetCurrentEpochInfo(s)
	if err != nil {
		return nil, fmt.Errorf("BeginBlock, GetCurrentEpochInfo error: %v", err)
	}
	validatorRewards := new(big.Int).Div(newRewards, new(big.Int).SetUint64(uint64(len(epochInfo.Validators))))
	for _, v := range epochInfo.Validators {
		dec, err := hexutil.Decode(v.PubKey)
		if err != nil {
			return nil, fmt.Errorf("BeginBlock, decode pubkey error: %v", err)
		}
		validator, found, err := GetValidator(s, dec)
		if err != nil {
			return nil, fmt.Errorf("BeginBlock, GetValidator error: %v", err)
		}
		if !found {
			panic("validator is not found")
		}
		err = allocateRewardsToValidator(s, validator, validatorRewards)
		if err != nil {
			return nil, fmt.Errorf("BeginBlock, allocateRewardsToValidator error: %v", err)
		}
	}

	// update outstanding rewards
	outstanding.Rewards = new(big.Int).Add(outstanding.Rewards, newRewards)
	err = setOutstandingRewards(s, outstanding)
	if err != nil {
		return nil, fmt.Errorf("BeginBlock, setOutstandingRewards error: %v", err)
	}

	return utils.ByteSuccess, nil
}