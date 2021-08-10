package boot

import (
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance"
	"github.com/ethereum/go-ethereum/contracts/native/governance/neo3_state_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/relayer_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync"
)

func InitialNativeContracts() {
	governance.InitGovernance()
	header_sync.InitHeaderSync()
	cross_chain_manager.InitCrossChainManager()
	neo3_state_manager.InitNeo3StateManager()
	node_manager.InitNodeManager()
	relayer_manager.InitRelayerManager()
	side_chain_manager.InitSideChainManager()

}
