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
package relayer_manager

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const abijson = `[
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"ID","type":"uint64"}],"name":"` + MethodApproveRegisterRelayer + `","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"ID","type":"uint64"}],"name":"` + MethodApproveRemoveRelayer + `","type":"event"},
	{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"applyID","type":"uint64"}],"name":"MethodRegisterRelayerEvent","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"removeID","type":"uint64"}],"name":"MethodRemoveRelayerEvent","type":"event"},
    {"inputs":[{"internalType":"uint64","name":"ID","type":"uint64"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodApproveRegisterRelayer + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"uint64","name":"ID","type":"uint64"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodApproveRemoveRelayer + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[],"name":"` + MethodContractName + `","outputs":[{"internalType":"string","name":"Name","type":"string"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"address[]","name":"AddressList","type":"address[]"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodRegisterRelayer + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"address[]","name":"AddressList","type":"address[]"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodRemoveRelayer + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"}
	]`

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(abijson))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

type RelayerListParam struct {
	AddressList []common.Address
	Address     common.Address
}

type ApproveRelayerParam struct {
	ID      uint64
	Address common.Address
}
