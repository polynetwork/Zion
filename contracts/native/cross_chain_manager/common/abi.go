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
package common

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/cross_chain_manager_abi"
)

var (
	MethodContractName        = cross_chain_manager_abi.MethodName
	MethodImportOuterTransfer = cross_chain_manager_abi.MethodImportOuterTransfer
	MethodMultiSign           = cross_chain_manager_abi.MethodMultiSign
	MethodBlackChain          = cross_chain_manager_abi.MethodBlackChain
	MethodWhiteChain          = cross_chain_manager_abi.MethodWhiteChain
)

var ABI *abi.ABI

func init() {
	ABI = GetABI()
}

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(cross_chain_manager_abi.CrossChainManagerABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

type BlackChainParam struct {
	ChainID uint64
}
