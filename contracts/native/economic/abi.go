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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/economic_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "node manager"

func InitABI() {
	ab, err := abi.JSON(strings.NewReader(IEconomicABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	ABI = &ab
}

var (
	ABI  *abi.ABI
	this = utils.EconomicContractAddress
)

type MethodContractNameInput struct{}
func (m *MethodContractNameInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodName)
}
func (m *MethodContractNameInput) Decode(payload []byte) error { return nil }

type MethodContractNameOutput struct {
	Name string
}
func (m *MethodContractNameOutput) Encode() ([]byte, error) {
	m.Name = contractName
	return utils.PackOutputs(ABI, MethodName, m.Name)
}
func (m *MethodContractNameOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodName, m, payload)
}

type MethodTotalSupplyInput struct{}
func (m *MethodTotalSupplyInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodTotalSupply)
}
func (m *MethodTotalSupplyInput) Decode(payload []byte) error { return nil }

type MethodRewardInput struct{}
func (m *MethodRewardInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodReward)
}
func (m *MethodRewardInput) Decode(payload []byte) error { return nil }

type MethodRewardOutput struct {
	List []*RewardAmount
}
func (m *MethodRewardOutput) Encode() ([]byte, error) {
	enc, err := rlp.EncodeToBytes(m.List)
	if err != nil {
		return nil, err
	}
	return utils.PackOutputs(ABI, MethodReward, enc)
}
func (m *MethodRewardOutput) Decode(payload []byte) error {
	var data struct{
		List []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodReward, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.List, &m.List)
}
//
//type RewardAmountList struct {
//	List []*RewardAmount
//}
//
//func (m *RewardAmountList) EncodeRLP(w io.Writer) error {
//	return rlp.Encode(w, []interface{}{m.List})
//}
//
//func (m *RewardAmountList) DecodeRLP(s *rlp.Stream) error {
//	var data struct {
//		List []*RewardAmount
//	}
//
//	if err := s.Decode(&data); err != nil {
//		return err
//	}
//	m.List = data.List
//	return nil
//}

type RewardAmount struct {
	Address common.Address
	Amount *big.Int
}

func (m *RewardAmount) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Address, m.Amount})
}

func (m *RewardAmount) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Address common.Address
		Amount *big.Int
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.Address, m.Amount = data.Address, data.Amount
	return nil
}
