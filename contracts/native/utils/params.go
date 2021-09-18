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
	HeaderSyncContractAddress        = common.HexToAddress("0xb2799bDE6831449d73C1F22CE815f773D0CafCc5")
	CrossChainManagerContractAddress = common.HexToAddress("0x5747C05FF236F8d18BB21Bc02ecc389deF853cae")
	SideChainManagerContractAddress  = common.HexToAddress("0x864Ff06eC5fFc75aB6eaf64263308ef5fa7d6637")
	NodeManagerContractAddress       = common.HexToAddress("0xA4Bf827047a08510722B2d62e668a72FCCFa232C")
	RelayerManagerContractAddress    = common.HexToAddress("0xA22f301D7Cb5b50dcA4a015b12EC0cc5f3971412")
	Neo3StateManagerContractAddress  = common.HexToAddress("0x5E839898821dB2A2F0eC9F8aAE7D7053744DB051")

	BTC_ROUTER              = uint64(1)
	ETH_ROUTER              = uint64(2)
	ONT_ROUTER              = uint64(3)
	NEO_ROUTER              = uint64(4)
	COSMOS_ROUTER           = uint64(5)
	BSC_ROUTER              = uint64(6)
	HECO_ROUTER             = uint64(7)
	QUORUM_ROUTER           = uint64(8)
	ZILLIQA_ROUTER          = uint64(9)
	MSC_ROUTER              = uint64(10)
	NEO3_LEGACY_ROUTER      = uint64(11)
	OKEX_ROUTER             = uint64(12)
	NEO3_ROUTER             = uint64(14)
	POLYGON_HEIMDALL_ROUTER = uint64(15)
	POLYGON_BOR_ROUTER      = uint64(16)
)
