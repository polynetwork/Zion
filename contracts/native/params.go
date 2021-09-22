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
package native

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

// MinGasUsage tx's gas usage should not be greater than an minimum fixed value if it execute failed.
const MinGasUsage = uint64(100)

const (
	NativeGovernance       = "governance"
	NativeSyncHeader       = "sync_header"
	NativeCrossChain       = "cross_chain"
	NativeNeo3StateManager = "neo3_state_manager"
	NativeNodeManager      = "node_manager"
	NativeRelayerManager   = "relayer_manager"
	NativeSideChainManager = "side_chain_manager"
	// native backup contracts
	NativeExtra4  = "extra4"
	NativeExtra5  = "extra5"
	NativeExtra6  = "extra6"
	NativeExtra7  = "extra7"
	NativeExtra8  = "extra8"
	NativeExtra9  = "extra9"
	NativeExtra10 = "extra10"
	NativeExtra11 = "extra11"
	NativeExtra12 = "extra12"
	NativeExtra13 = "extra13"
	NativeExtra14 = "extra14"
	NativeExtra15 = "extra15"
	NativeExtra16 = "extra16"
	NativeExtra17 = "extra17"
	NativeExtra18 = "extra18"
	NativeExtra19 = "extra19"
)

var NativeContractAddrMap = map[string]common.Address{
	NativeGovernance:       common.HexToAddress("0x4600691499997fCc224425ba5C93EebC57f3615b"),
	NativeSyncHeader:       utils.HeaderSyncContractAddress,
	NativeCrossChain:       utils.CrossChainManagerContractAddress,
	NativeNeo3StateManager: utils.CrossChainManagerContractAddress,
	NativeNodeManager:      utils.NodeManagerContractAddress,
	NativeRelayerManager:   utils.RelayerManagerContractAddress,
	NativeSideChainManager: utils.SideChainManagerContractAddress,
	NativeExtra4:           common.HexToAddress("0x7d79D936DA7833c7fe056eB450064f34A327DcA8"),
	NativeExtra5:           common.HexToAddress("0xD37F626c9E007DdD244E5Cbee0C223fec6D11289"),
	NativeExtra6:           common.HexToAddress("0x33463b771Da32D450723C7C23a2240dE223b53bd"),
	NativeExtra7:           common.HexToAddress("0x0F257CD338Fa8F1Af3D31b16C1fBddae2Dc96D41"),
	NativeExtra8:           common.HexToAddress("0x4479AcbCeA458Badf21dbEC7Db6fC236Bf08fbb9"),
	NativeExtra9:           common.HexToAddress("0xc204aDF052C52F74863d76c94a311b82D98d87AE"),
	NativeExtra10:          common.HexToAddress("0xD62B67170A6bb645f1c59601FbC6766940ee12e5"),
	NativeExtra11:          common.HexToAddress("0xf7EBd79DB6240b9A85571f61b543425e2A7045Fb"),
	NativeExtra12:          common.HexToAddress("0x20B019ea369923eF1971A30f1974003051f1863C"),
	NativeExtra13:          common.HexToAddress("0x2951b823F25344797D9294634F44e867490B86c9"),
	NativeExtra14:          common.HexToAddress("0x370f0dDA62BDc610d8FFE8c71882D27d2a26648f"),
	NativeExtra15:          common.HexToAddress("0xC782D7244bdd2ebeb56ac87A65c4873B6c4D427D"),
	NativeExtra16:          common.HexToAddress("0x90dc8B0B8625DD3Fa33eBd5E502D6c908EFB68Fe"),
	NativeExtra17:          common.HexToAddress("0x40E25A4c3316F54c913542Ad293420cF3c6C2Ff3"),
	NativeExtra18:          common.HexToAddress("0x5e66f4D53236348334E13F1d5F83b48083a4ADd0"),
	NativeExtra19:          common.HexToAddress("0x0763E5717f8bD8C710E0d38a21e224D8C560e597"),
}

func IsNativeContract(addr common.Address) bool {
	for _, v := range NativeContractAddrMap {
		if v == addr {
			return true
		}
	}
	return false
}
