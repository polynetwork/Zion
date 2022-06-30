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
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"math/big"
	"strings"
)

const contractName = "node manager"

func InitABI() {
	ab, err := abi.JSON(strings.NewReader(INodeManagerABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	ABI = &ab
}

var (
	ABI  *abi.ABI
	this = utils.NodeManagerContractAddress
)

type CreateValidatorParam struct {
	ConsensusPubkey string
	ProposalAddress common.Address
	Commission      *big.Int
	InitStake       *big.Int
	Desc            string
}

func (m *CreateValidatorParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodCreateValidator, m)
}

type UpdateValidatorParam struct {
	ConsensusPubkey string
	ProposalAddress common.Address
	Desc            string
}

func (m *UpdateValidatorParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodUpdateValidator, m)
}

type UpdateCommissionParam struct {
	ConsensusPubkey string
	Commission      *big.Int
}

func (m *UpdateCommissionParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodUpdateCommission, m)
}

type StakeParam struct {
	ConsensusPubkey string
	Amount          *big.Int
}

func (m *StakeParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodStake, m)
}

type UnStakeParam struct {
	ConsensusPubkey string
	Amount          *big.Int
}

func (m *UnStakeParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodUnStake, m)
}

type CancelValidatorParam struct {
	ConsensusPubkey string
}

func (m *CancelValidatorParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodCancelValidator, m)
}

type WithdrawValidatorParam struct {
	ConsensusPubkey string
}

func (m *WithdrawValidatorParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodWithdrawValidator, m)
}

type WithdrawStakeRewardsParam struct {
	ConsensusPubkey string
}

func (m *WithdrawStakeRewardsParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodWithdrawStakeRewards, m)
}

type WithdrawCommissionParam struct {
	ConsensusPubkey string
}

func (m *WithdrawCommissionParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodWithdrawCommission, m)
}

type ChangeEpochParam struct{}

func (m *ChangeEpochParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodChangeEpoch)
}

type WithdrawParam struct{}

func (m *WithdrawParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodWithdraw)
}

type EndBlockParam struct{}

func (m *EndBlockParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodEndBlock)
}

type GetGlobalConfigParam struct {}

func (m *GetGlobalConfigParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodGetGlobalConfig)
}

type GetCommunityInfoParam struct {}

func (m *GetCommunityInfoParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodGetCommunityInfo)
}

type GetCurrentEpochInfoParam struct {}

func (m *GetCurrentEpochInfoParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodGetCurrentEpochInfo)
}
