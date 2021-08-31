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
)

const (
	MethodContractName      = "name"
	MethodSyncGenesisHeader = "syncGenesisHeader"
	MethodSyncBlockHeader   = "syncBlockHeader"
	MethodSyncCrossChainMsg = "syncCrossChainMsg"
)

var GasTable = map[string]uint64{
	MethodContractName:      0,
	MethodSyncGenesisHeader: 0,
	MethodSyncBlockHeader:   100000,
	MethodSyncCrossChainMsg: 0,
}

const abijson = `[
	{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"chainID","type":"uint64"},{"indexed":false,"internalType":"string","name":"BlockHash","type":"string"},{"indexed":false,"internalType":"uint64","name":"Height","type":"uint64"},{"indexed":false,"internalType":"string","name":"NextValidatorsHash","type":"string"},{"indexed":false,"internalType":"string","name":"InfoChainID","type":"string"},{"indexed":false,"internalType":"uint64","name":"BlockHeight","type":"uint64"}],"name":"OKEpochSwitchInfoEvent","type":"event"},
	{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"chainID","type":"uint64"},{"indexed":false,"internalType":"uint64","name":"height","type":"uint64"},{"indexed":false,"internalType":"string","name":"blockHash","type":"string"},{"indexed":false,"internalType":"uint256","name":"BlockHeight","type":"uint256"}],"name":"` + SYNC_HEADER_NAME_EVENT + `","type":"event"},
    {"inputs":[],"name":"` + MethodContractName + `","outputs":[{"internalType":"string","name":"Name","type":"string"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"uint64","name":"ChainID","type":"uint64"},{"internalType":"address","name":"Address","type":"address"},{"internalType":"bytes[]","name":"Headers","type":"bytes[]"}],"name":"` + MethodSyncBlockHeader + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"uint64","name":"ChainID","type":"uint64"},{"internalType":"address","name":"Address","type":"address"},{"internalType":"bytes[]","name":"CrossChainMsgs","type":"bytes[]"}],"name":"` + MethodSyncCrossChainMsg + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"uint64","name":"ChainID","type":"uint64"},{"internalType":"bytes","name":"GenesisHeader","type":"bytes"}],"name":"` + MethodSyncGenesisHeader + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"}
]`

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(abijson))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

var ABI *abi.ABI
