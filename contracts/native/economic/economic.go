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

package economic

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/contracts/native"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/economic_abi"
	"github.com/ethereum/go-ethereum/contracts/native/governance/community"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

var (
	gasTable = map[string]uint64{
		MethodName:        39375,
		MethodTotalSupply: 23625,
		MethodReward:      73500,
	}
)

var (
	RewardPerBlock = params.ZNT1
	GenesisSupply  = params.GenesisSupply
)

func InitEconomic() {
	InitABI()
	native.Contracts[this] = RegisterEconomicContract
}

func RegisterEconomicContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodName, Name)
	s.Register(MethodTotalSupply, TotalSupply)
	s.Register(MethodReward, Reward)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return new(MethodContractNameOutput).Encode()
}

func TotalSupply(s *native.NativeContract) ([]byte, error) {
	height := s.ContractRef().BlockHeight()

	supply := GenesisSupply
	if height.Uint64() > 0 {
		reward := new(big.Int).Mul(height, RewardPerBlock)
		supply = new(big.Int).Add(supply, reward)
	}
	return utils.PackOutputs(ABI, MethodTotalSupply, supply)
}

func getBlockRewardList(s *native.NativeContract) ([]*RewardAmount, error) {
	community, err := community.GetCommunityInfoFromDB(s.StateDB())
	if err != nil {
		return nil, fmt.Errorf("GetCommunityInfo failed, err: %v", err)
	}

	// allow empty address as reward pool
	poolAddr := community.CommunityAddress
	rewardPerBlock := utils.NewDecFromBigInt(RewardPerBlock)
	rewardFactor := utils.NewDecFromBigInt(community.CommunityRate)
	poolRwdAmt, err := rewardPerBlock.MulWithPercentDecimal(rewardFactor)
	if err != nil {
		return nil, fmt.Errorf("calculate pool reward amount failed, err: %v ", err)
	}
	stakingRwdAmt, err := rewardPerBlock.Sub(poolRwdAmt)
	if err != nil {
		return nil, fmt.Errorf("calculate staking reward amount, failed, err: %v ", err)
	}

	poolRwd := &RewardAmount{
		Address: poolAddr,
		Amount:  poolRwdAmt.BigInt(),
	}
	stakingRwd := &RewardAmount{
		Address: utils.NodeManagerContractAddress,
		Amount:  stakingRwdAmt.BigInt(),
	}

	return []*RewardAmount{poolRwd, stakingRwd}, nil
}

func Reward(s *native.NativeContract) ([]byte, error) {
	list, err := getBlockRewardList(s)
	if err != nil {
		return nil, err
	}
	output := new(MethodRewardOutput)
	output.List = list
	return output.Encode()
}

func GenerateBlockReward(s *native.NativeContract) error {
	height := s.ContractRef().BlockHeight()
	// genesis block do not need to distribute reward
	if height.Uint64() == 0 {
		return nil
	}

	// get reward info list from native contract of `economic`
	list, err := getBlockRewardList(s)
	if err != nil {
		return err
	}

	// add balance to related addresses
	var sRwd string
	for _, v := range list {
		s.StateDB().AddBalance(v.Address, v.Amount)
		sRwd += fmt.Sprintf("address: %s, amount %v;", v.Address.Hex(), v.Amount)
	}
	log.Debug("reward", "num", height, "list", sRwd)

	return nil
}
