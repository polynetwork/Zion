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
package boot

import (
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/economic"
	"github.com/ethereum/go-ethereum/contracts/native/governance/neo3_state_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/relayer_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/signature_manager"
	"github.com/ethereum/go-ethereum/contracts/native/info_sync"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/log"
)

func InitNativeContracts() {
	node_manager.InitNodeManager()
	economic.InitEconomic()
	info_sync.InitInfoSync()
	cross_chain_manager.InitCrossChainManager()
	side_chain_manager.InitSideChainManager()
	relayer_manager.InitRelayerManager()
	neo3_state_manager.InitNeo3StateManager()
	signature_manager.InitSignatureManager()

	log.Info("Initialize main chain native contracts",
		"node manager", utils.NodeManagerContractAddress.Hex(),
		"economic", utils.EconomicContractAddress.Hex(),
		"header sync", utils.InfoSyncContractAddress.Hex(),
		"cross chain manager", utils.CrossChainManagerContractAddress.Hex(),
		"side chain manager", utils.SideChainManagerContractAddress.Hex(),
		"relayer manager", utils.RelayerManagerContractAddress.Hex(),
		"neo3 state manager", utils.Neo3StateManagerContractAddress.Hex(),
		"signature manager", utils.SignatureManagerContractAddress.Hex(),
	)

}
