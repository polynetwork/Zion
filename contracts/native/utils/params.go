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
package utils

import "github.com/ethereum/go-ethereum/common"

type BtcNetType int

const (
	TyTestnet3 BtcNetType = iota
	TyRegtest
	TySimnet
	TyMainnet
)

var (
	BYTE_FALSE = []byte{0}
	BYTE_TRUE  = []byte{1}
)

var (
	NodeManagerContractAddress       = common.HexToAddress("0x0000000000000000000000000000000000001000")
	EconomicContractAddress          = common.HexToAddress("0x0000000000000000000000000000000000001001")
	InfoSyncContractAddress          = common.HexToAddress("0x0000000000000000000000000000000000001002")
	CrossChainManagerContractAddress = common.HexToAddress("0x0000000000000000000000000000000000001003")
	SideChainManagerContractAddress  = common.HexToAddress("0x0000000000000000000000000000000000001004")
	RelayerManagerContractAddress    = common.HexToAddress("0x0000000000000000000000000000000000001005")
	Neo3StateManagerContractAddress  = common.HexToAddress("0x0000000000000000000000000000000000001006")
	SignatureManagerContractAddress  = common.HexToAddress("0x0000000000000000000000000000000000001007")
	ProposalManagerContractAddress   = common.HexToAddress("0x0000000000000000000000000000000000001008")

	NO_PROOF_ROUTER    = uint64(0)
	BTC_ROUTER         = uint64(1)
	ETH_ROUTER         = uint64(2)
	ONT_ROUTER         = uint64(3)
	NEO_ROUTER         = uint64(4)
	COSMOS_ROUTER      = uint64(5)
	BSC_ROUTER         = uint64(6)
	HECO_ROUTER        = uint64(7)
	QUORUM_ROUTER      = uint64(8)
	ZILLIQA_ROUTER     = uint64(9)
	MSC_ROUTER         = uint64(10)
	NEO3_LEGACY_ROUTER = uint64(11)
	OKEX_ROUTER        = uint64(12)
	NEO3_ROUTER        = uint64(14)

	ETH_COMMON_ROUTER = uint64(15)
)
