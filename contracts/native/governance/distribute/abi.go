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

package distribute

import (
	"fmt"
	"io"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/distribute_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
)

const contractName = "distribute"

func InitABI() {
	ab, err := abi.JSON(strings.NewReader(IDistributeABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	ABI = &ab
}

var (
	ABI  *abi.ABI
	this = utils.DistributeContractAddress
)

type WithdrawStakeRewardsParam struct {
	ConsensusPubkey string
}

func (m *WithdrawStakeRewardsParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ConsensusPubkey})
}

func (m *WithdrawStakeRewardsParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ConsensusPubkey string
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.ConsensusPubkey = data.ConsensusPubkey
	return nil
}
