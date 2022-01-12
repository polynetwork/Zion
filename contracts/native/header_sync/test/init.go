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
package test

import (
	"crypto/ecdsa"
	"encoding/json"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"

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
)

func init() {
	node_manager.InitNodeManager()
	side_chain_manager.InitSideChainManager()
	header_sync.InitHeaderSync()
	db := rawdb.NewMemoryDatabase()
	sdb, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)

	putPeerMapPoolAndView(sdb)
	putSideChain()
}

func putPeerMapPoolAndView(db *state.StateDB) {
	height := uint64(120)
	epoch := node_manager.GenerateTestEpochInfo(1, height, 4)

	peer := epoch.Peers.List[0]
	rawPubKey, _ := hexutil.Decode(peer.PubKey)
	pubkey, _ := crypto.DecompressPubkey(rawPubKey)
	acct = pubkey
	caller := peer.Address

	txhash := common.HexToHash("0x123")
	ref := native.NewContractRef(db, caller, caller, new(big.Int).SetUint64(height), txhash, 0, nil)
	s := native.NewNativeContract(db, ref)
	node_manager.StoreTestEpoch(s, epoch)
}

func putSideChain() {
	caller := crypto.PubkeyToAddress(*acct)
	blockNumber := big.NewInt(1)
	extra := uint64(10)
	contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, extra, nil)
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
