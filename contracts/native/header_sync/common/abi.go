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
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/header_sync_abi"
)

var (
	MethodContractName      = header_sync_abi.MethodName
	MethodSyncGenesisHeader = header_sync_abi.MethodSyncGenesisHeader
	MethodSyncBlockHeader   = header_sync_abi.MethodSyncBlockHeader
	MethodSyncCrossChainMsg = header_sync_abi.MethodSyncCrossChainMsg
)

var GasTable = map[string]uint64{
	MethodContractName:      0,
	MethodSyncGenesisHeader: 0,
	MethodSyncBlockHeader:   10000,
	MethodSyncCrossChainMsg: 0,
}

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(header_sync_abi.HeaderSyncABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

var ABI *abi.ABI
