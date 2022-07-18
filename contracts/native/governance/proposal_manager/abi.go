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

package proposal_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/proposal_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"math/big"
	"strings"
)

const contractName = "proposal manager"

func InitABI() {
	ab, err := abi.JSON(strings.NewReader(IProposalManagerABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	ABI = &ab
}

var (
	ABI  *abi.ABI
	this = utils.ProposalManagerContractAddress
)

type ProposeParam struct {
	PType    ProposalType
	Content []byte
	Stake   *big.Int
}

func (m *ProposeParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodPropose, m)
}

type VoteActiveProposalParam struct {
	ID *big.Int
}

func (m *VoteActiveProposalParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodVoteActiveProposal, m)
}
