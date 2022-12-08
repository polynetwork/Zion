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

package mock

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/backend"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
)

type Geth struct {
	miner       *miner
	chain       *core.BlockChain
	engine      Engine
	hotstuff    consensus.HotStuff
	api         *backend.API
	broadcaster *broadcaster
}

func MakeGeth(privateKey *ecdsa.PrivateKey, vals []common.Address) *Geth {
	db := rawdb.NewMemoryDatabase()
	engine := makeEngine(privateKey, db)
	chain := makeChain(db, engine, vals)
	hotstuffEngine := engine.(consensus.HotStuff)
	broadcaster := engine.(consensus.Handler).GetBroadcaster().(*broadcaster)
	api := engine.APIs(chain)[0].Service.(*backend.API)
	miner := makeMiner(broadcaster.addr, chain, hotstuffEngine)
	geth := &Geth{
		miner:       miner,
		chain:       chain,
		engine:      engine,
		api:         api,
		hotstuff:    hotstuffEngine,
		broadcaster: broadcaster,
	}
	miner.geth = geth
	broadcaster.geth = geth
	return geth
}

func (g *Geth) Start() {
	g.hotstuff.Start(g.chain, nil)
	g.miner.Start()
}

func (g *Geth) Stop() {
	g.broadcaster.Stop()
	g.hotstuff.Stop()
	g.miner.Stop()
}

func (g *Geth) Sequence() (uint64, uint64) {
	return g.api.CurrentSequence()
}

func (g *Geth) IsProposer() bool {
	return g.api.IsProposer()
}

func (g *Geth) broadcastBlock(block *types.Block) {
	var (
		td   *big.Int
		hash = block.Hash()
	)
	if parent := g.chain.GetBlock(block.ParentHash(), block.NumberU64()-1); parent != nil {
		td = new(big.Int).Add(block.Difficulty(), g.chain.GetTd(block.ParentHash(), block.NumberU64()-1))
	} else {
		log.Error("Propagating dangling block", "number", block.Number(), "hash", hash)
		return
	}

	// Send the block to a subset of our peers
	for _, peer := range g.broadcaster.peers {
		go func(p *MockPeer) {
			if err := p.SendNewBlock(block, td); err != nil {
				log.Error("SendNewBlock", "to", peer.remote, "err", err)
			}
		}(peer)
	}
}

func (g *Geth) handleBlock(msg p2p.Msg) {
	ann := new(eth.NewBlockPacket)
	if err := msg.Decode(ann); err != nil {
		log.Error("decode newBlockPacket failed", "err", err)
		return
	}
	block := ann.Block
	curBlock := g.chain.CurrentBlock().Copy()
	if block.NumberU64() == curBlock.NumberU64()+1 && block.ParentHash() == curBlock.Hash() {
		statedb, receipts, allLogs, err := g.chain.ExecuteBlock(block)
		if err != nil {
			log.Error("failed to execute block", "err", err)
			return
		}
		if _, err := g.chain.WriteBlockWithState(block, receipts, allLogs, statedb, false); err != nil {
			log.Error("failed to writeBlockWithState", "err", err)
		}
	}
	////g.chain.GetBlockByHash(block.Hash())
	//if _, err := g.chain.InsertChain([]*types.Block{block}); err != nil {
	//	log.Error("insert chain failed", "num", block.Number(), "hash", block.Hash(), "err", err)
	//}
}

type System struct {
	nodes []*Geth
}

func makeSystem(n int) *System {
	pks, addrs := newAccountLists(n)
	nodes := make([]*Geth, n)

	for i := 0; i < n; i++ {
		nodes[i] = MakeGeth(pks[i], addrs)
	}

	return &System{nodes: nodes}
}

func (s *System) Start() {
	for i := 0; i < len(s.nodes); i++ {
		for j := 0; j < len(s.nodes); j++ {
			if i != j {
				src := s.nodes[i].engine.(consensus.Handler).GetBroadcaster()
				dst := s.nodes[j].engine.(consensus.Handler).GetBroadcaster()
				bsrc := src.(*broadcaster)
				bdst := dst.(*broadcaster)
				bsrc.Connect(bdst)
			}
		}
	}
	for _, node := range s.nodes {
		go node.Start()
	}
}

func (s *System) Stop() {
	for _, cli := range s.nodes {
		go cli.Stop()
	}
}

func (s *System) Leader() *Geth {
	for _, node := range s.nodes {
		if node.IsProposer() {
			return node
		}
	}
	return nil
}

func newAccountLists(n int) ([]*ecdsa.PrivateKey, []common.Address) {
	pks := make([]*ecdsa.PrivateKey, n)
	addrs := make([]common.Address, n)
	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey()
		pks[i] = key
		addrs[i] = crypto.PubkeyToAddress(key.PublicKey)
	}
	return pks, addrs
}
