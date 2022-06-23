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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/economic_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/params"
)

var (
	gasTable = map[string]uint64{
		MethodName:        0,
		MethodTotalSupply: 0,
		MethodReward:      0,
	}
)

var (
	RewardPerBlock = params.ZNT1
	GenesisSupply  = params.GenesisSupply

	// rewardPoolFactor default value should be 0.2 = 2000/10000
	defaultPoolRewardFactor = big.NewInt(2000)
	totalRewardFactor       = big.NewInt(10000)
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

// todo(fuk): getPoolRewardFactor from governance contract
// todo(fuk): get reward pool address from governance
func Reward(s *native.NativeContract) ([]byte, error) {
	data := new(big.Int).Mul(RewardPerBlock, defaultPoolRewardFactor)
	poolRwdAmt := new(big.Int).Div(data, totalRewardFactor)
	stakingRwdAmt := new(big.Int).Sub(RewardPerBlock, poolRwdAmt)

	poolAddr := common.HexToAddress("0x0000000000000000000000000000000001000000")
	poolRwd := &RewardAmount{
		Address: poolAddr,
		Amount:  poolRwdAmt,
	}
	stakingRwd := &RewardAmount{
		Address: utils.GovernanceContractAddress,
		Amount:  stakingRwdAmt,
	}

	output := new(MethodRewardOutput)
	output.List = []*RewardAmount{poolRwd, stakingRwd}
	return output.Encode()
}
