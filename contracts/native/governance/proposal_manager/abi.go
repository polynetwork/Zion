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
	Content []byte
}

func (m *ProposeParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodPropose, m)
}

type ProposeConfigParam struct {
	Content []byte
}

func (m *ProposeConfigParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodProposeConfig, m)
}

type ProposeCommunityParam struct {
	Content []byte
}

func (m *ProposeCommunityParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodProposeCommunity, m)
}

type VoteProposalParam struct {
	ID *big.Int
}

func (m *VoteProposalParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodVoteProposal, m)
}

type GetProposalParam struct {
	ID *big.Int
}

func (m *GetProposalParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodGetProposal, m)
}

type GetProposalListParam struct{}

func (m *GetProposalListParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodGetProposalList)
}

type GetConfigProposalListParam struct{}

func (m *GetConfigProposalListParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodGetConfigProposalList)
}

type GetCommunityProposalListParam struct{}

func (m *GetCommunityProposalListParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodGetCommunityProposalList)
}
