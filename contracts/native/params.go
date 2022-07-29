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

const (
	NativeGovernance         = "governance"
	NativeSyncCrossChainInfo = "sync_cross_chain_info"
	NativeCrossChain         = "cross_chain"
	NativeNeo3StateManager   = "neo3_state_manager"
	NativeNodeManager        = "node_manager"
	NativeRelayerManager     = "relayer_manager"
	NativeSideChainManager   = "side_chain_manager"
	NativeEconomic           = "economic"
	NativeSignatureManager   = "signature_manager"
	NativeProposalManager    = "proposal_manager"

	// native backup contracts
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
	NativeNodeManager:        utils.NodeManagerContractAddress,
	NativeEconomic:           utils.EconomicContractAddress,
	NativeSyncCrossChainInfo: utils.InfoSyncContractAddress,
	NativeCrossChain:         utils.CrossChainManagerContractAddress,
	NativeSideChainManager:   utils.SideChainManagerContractAddress,
	NativeRelayerManager:     utils.RelayerManagerContractAddress,
	NativeNeo3StateManager:   utils.Neo3StateManagerContractAddress,
	NativeSignatureManager:   utils.SignatureManagerContractAddress,
	NativeProposalManager:    utils.ProposalManagerContractAddress,
	NativeExtra6:             common.HexToAddress("0x0000000000000000000000000000000000001009"),
	NativeExtra7:             common.HexToAddress("0x000000000000000000000000000000000000100a"),
	NativeExtra8:             common.HexToAddress("0x000000000000000000000000000000000000100b"),
	NativeExtra9:             common.HexToAddress("0x000000000000000000000000000000000000100c"),
	NativeExtra10:            common.HexToAddress("0x000000000000000000000000000000000000100d"),
	NativeExtra11:            common.HexToAddress("0x000000000000000000000000000000000000100e"),
	NativeExtra12:            common.HexToAddress("0x000000000000000000000000000000000000100f"),
	NativeExtra13:            common.HexToAddress("0x0000000000000000000000000000000000001010"),
	NativeExtra14:            common.HexToAddress("0x0000000000000000000000000000000000001011"),
	NativeExtra15:            common.HexToAddress("0x0000000000000000000000000000000000001012"),
	NativeExtra16:            common.HexToAddress("0x0000000000000000000000000000000000001013"),
	NativeExtra17:            common.HexToAddress("0x0000000000000000000000000000000000001014"),
	NativeExtra18:            common.HexToAddress("0x0000000000000000000000000000000000001015"),
	NativeExtra19:            common.HexToAddress("0x0000000000000000000000000000000000001016"),
}

func IsNativeContract(addr common.Address) bool {
	for _, v := range NativeContractAddrMap {
		if v == addr {
			return true
		}
	}
	return false
}
