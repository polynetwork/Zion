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
package cross_chain_manager

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const abijson = ``

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(abijson))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

type EntranceParam struct {
	SourceChainID         uint64 `json:"sourceChainId"`
	Height                uint32 `json:"height"`
	Proof                 []byte `json:"proof"`
	RelayerAddress        []byte `json:"relayerAddress"`
	Extra                 []byte `json:"extra"`
	HeaderOrCrossChainMsg []byte `json:"headerOrCrossChainMsg"`
}

type MultiSignParam struct {
	ChainID   uint64
	RedeemKey string
	TxHash    []byte
	Address   string
	Signs     [][]byte
}

type BlackChainParam struct {
	ChainID uint64
}
