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

package auth

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/auth_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

func InitABI() {
	ab, err := abi.JSON(strings.NewReader(IAuthABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	ABI = &ab
}

var (
	ABI *abi.ABI
)

type MethodApproveInput struct {
	Spender common.Address
	Amount  *big.Int
}

func (i *MethodApproveInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodApprove, i.Spender, i.Amount)
}

func (i *MethodApproveInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodApprove, i, payload)
}

type MethodAllowanceInput struct {
	Owner   common.Address
	Spender common.Address
}

func (i *MethodAllowanceInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodAllowance, i.Owner, i.Spender)
}

func (i *MethodAllowanceInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodAllowance, i, payload)
}

// event Approval(address indexed owner, address indexed spender, uint256 value);
func emitApprovedEvent(s *native.NativeContract, owner, spender common.Address, value *big.Int) error {
	return s.AddNotify(ABI, []string{EventApproval}, owner, spender, value)
}
