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
	mlp "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/mainchain/lock_proxy"
	slp "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/sidechain/lock_proxy"
	"github.com/ethereum/go-ethereum/contracts/native/governance/neo3_state_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/relayer_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/log"
)

func InitMainChainNativeContracts() {
	header_sync.InitHeaderSync()
	cross_chain_manager.InitCrossChainManager()
	neo3_state_manager.InitNeo3StateManager()
	node_manager.InitNodeManager()
	relayer_manager.InitRelayerManager()
	side_chain_manager.InitSideChainManager()
	mlp.InitLockProxy()

	log.Info("Initialize main chain native contracts",
		"header sync", utils.HeaderSyncContractAddress.Hex(),
		"cross chain manager", utils.CrossChainManagerContractAddress.Hex(),
		"neo3 state manager", utils.Neo3StateManagerContractAddress.Hex(),
		"node manager", utils.NodeManagerContractAddress.Hex(),
		"relayer manager", utils.RelayerManagerContractAddress.Hex(),
		"side chain manager", utils.SideChainManagerContractAddress.Hex(),
		"native lock proxy", utils.LockProxyContractAddress.Hex())
}

func InitSideChainNativeContracts() {
	node_manager.InitNodeManager()
	slp.InitLockProxy()
	log.Info("Initialize side chain native contracts",
		"node manager", utils.NodeManagerContractAddress.Hex(),
		"native lock proxy", utils.LockProxyContractAddress.Hex())
}
