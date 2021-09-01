package test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/msc"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/polygon"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	polycomm "github.com/polynetwork/poly/common"
	cstates "github.com/polynetwork/poly/core/states"
)

func init() {
	node_manager.InitNodeManager()
	side_chain_manager.InitSideChainManager()
	header_sync.InitHeaderSync()
	db := rawdb.NewMemoryDatabase()
	sdb, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)

	cacheDB := (*state.CacheDB)(sdb)
	putPeerMapPoolAndView(cacheDB)
	putSideChain()
}

func putPeerMapPoolAndView(db *state.CacheDB) {
	key, _ := crypto.GenerateKey()
	acct = &key.PublicKey

	peerPoolMap := new(node_manager.PeerPoolMap)
	peerPoolMap.PeerPoolMap = make(map[string]*node_manager.PeerPoolItem)
	pkStr := hex.EncodeToString(crypto.FromECDSAPub(acct))
	peerPoolMap.PeerPoolMap[pkStr] = &node_manager.PeerPoolItem{
		Index:      uint32(0),
		PeerPubkey: pkStr,
		Address:    crypto.PubkeyToAddress(*acct),
		Status:     node_manager.ConsensusStatus,
	}
	view := uint32(0)
	viewBytes := utils.GetUint32Bytes(view)
	sink := polycomm.NewZeroCopySink(nil)
	peerPoolMap.Serialization(sink)
	db.Put(utils.ConcatKey(utils.NodeManagerContractAddress, []byte(node_manager.PEER_POOL), viewBytes), cstates.GenRawStorageItem(sink.Bytes()))

	sink.Reset()

	govView := node_manager.GovernanceView{
		View: view,
	}
	govView.Serialization(sink)
	db.Put(utils.ConcatKey(utils.NodeManagerContractAddress, []byte(node_manager.GOVERNANCE_VIEW)), cstates.GenRawStorageItem(sink.Bytes()))
}

func putSideChain() {
	caller := crypto.PubkeyToAddress(*acct)
	blockNumber := big.NewInt(1)
	extra := uint64(10)
	contractRef := native.NewContractRef(sdb, caller, blockNumber, common.Hash{}, extra, nil)
	contract := native.NewNativeContract(sdb, contractRef)

	err := side_chain_manager.PutSideChain(contract, &side_chain_manager.SideChain{
		Router:  utils.POLYGON_HEIMDALL_ROUTER,
		ChainId: heimdalChainID,
	})
	if err != nil {
		log.Fatalf("PutSideChain fail:%v", err)
		return
	}
	sideChain, err := side_chain_manager.GetSideChain(contract, heimdalChainID)
	if err != nil {
		log.Fatalf("PutSideChain fail:%v", err)
		return
	}

	if sideChain.ChainId != heimdalChainID {
		log.Fatalf("GetSideChain mismatch")
	}
	extraInfo := polygon.ExtraInfo{
		Sprint:              64,
		Period:              2,
		ProducerDelay:       6,
		BackupMultiplier:    2,
		HeimdallPolyChainID: heimdalChainID,
	}
	extraInfoBytes, _ := json.Marshal(extraInfo)
	err = side_chain_manager.PutSideChain(contract, &side_chain_manager.SideChain{
		Router:    utils.POLYGON_BOR_ROUTER,
		ChainId:   borChainID,
		ExtraInfo: extraInfoBytes,
	})
	if err != nil {
		log.Fatalf("PutSideChain fail:%v", err)
		return
	}

	sideChain, err = side_chain_manager.GetSideChain(contract, borChainID)
	if err != nil {
		log.Fatalf("PutSideChain fail:%v", err)
		return
	}

	if sideChain.ChainId != borChainID {
		log.Fatalf("GetSideChain mismatch")
	}

	{
		err = side_chain_manager.PutSideChain(contract, &side_chain_manager.SideChain{
			Router:       utils.QUORUM_ROUTER,
			ChainId:      quorumChainID,
			BlocksToWait: 1,
			Address:      caller,
			CCMCAddress:  common.Hex2Bytes("0x0000000000000000000000000000000000000105"),
		})
		if err != nil {
			log.Fatalf("PutSideChain fail:%v", err)
			return
		}
	}

	{
		// add sidechain info
		extra := msc.ExtraInfo{
			// test id 97
			ChainID: big.NewInt(97),
			Period:  1,
			Epoch:   200,
		}
		extraBytes, _ := json.Marshal(extra)
		side_chain_manager.PutSideChain(contract, &side_chain_manager.SideChain{
			ExtraInfo: extraBytes,
			Router:    utils.MSC_ROUTER,
			ChainId:   mscChainID,
		})
	}

	{
		side_chain_manager.PutSideChain(contract, &side_chain_manager.SideChain{
			Router:  utils.ETH_ROUTER,
			ChainId: ethChainID,
		})
	}
}

var (
	sdb            *state.StateDB
	acct           *ecdsa.PublicKey
	heimdalChainID = uint64(2)
	borChainID     = uint64(3)
	quorumChainID  = uint64(4)
	mscChainID     = uint64(5)
	ethChainID     = uint64(6)
)
