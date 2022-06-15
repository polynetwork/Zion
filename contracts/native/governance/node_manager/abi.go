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
	"github.com/ethereum/go-ethereum/rlp"
	"io"
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

func (m *CreateValidatorParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ConsensusPubkey, m.ProposalAddress, m.Commission, m.InitStake, m.Desc})
}

func (m *CreateValidatorParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ConsensusPubkey string
		ProposalAddress common.Address
		Commission      *big.Int
		InitStake       *big.Int
		Desc            string
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.ConsensusPubkey, m.ProposalAddress, m.Commission, m.InitStake, m.Desc = data.ConsensusPubkey,
		data.ProposalAddress, data.Commission, data.InitStake, data.Desc
	return nil
}

type UpdateValidatorParam struct {
	ConsensusPubkey string
	ProposalAddress common.Address
	Commission      *big.Int
	Desc            string
}

func (m *UpdateValidatorParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ConsensusPubkey, m.ProposalAddress, m.Commission, m.Desc})
}

func (m *UpdateValidatorParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ConsensusPubkey string
		ProposalAddress common.Address
		Commission      *big.Int
		Desc            string
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.ConsensusPubkey, m.ProposalAddress, m.Commission, m.Desc = data.ConsensusPubkey, data.ProposalAddress,
		data.Commission, data.Desc
	return nil
}

type StakeParam struct {
	ConsensusPubkey string
	Amount          *big.Int
}

func (m *StakeParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ConsensusPubkey, m.Amount})
}

func (m *StakeParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ConsensusPubkey string
		Amount          *big.Int
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.ConsensusPubkey, m.Amount = data.ConsensusPubkey, data.Amount
	return nil
}

type UnStakeParam struct {
	ConsensusPubkey string
	Amount          *big.Int
}

func (m *UnStakeParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ConsensusPubkey, m.Amount})
}

func (m *UnStakeParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ConsensusPubkey string
		Amount          *big.Int
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.ConsensusPubkey, m.Amount = data.ConsensusPubkey, data.Amount
	return nil
}

type CancelValidatorParam struct {
	ConsensusPubkey string
}

func (m *CancelValidatorParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ConsensusPubkey})
}

func (m *CancelValidatorParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ConsensusPubkey string
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.ConsensusPubkey = data.ConsensusPubkey
	return nil
}

type WithdrawValidatorParam struct {
	ConsensusPubkey string
}

func (m *WithdrawValidatorParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ConsensusPubkey})
}

func (m *WithdrawValidatorParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ConsensusPubkey string
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.ConsensusPubkey = data.ConsensusPubkey
	return nil
}