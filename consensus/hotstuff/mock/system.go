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
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/backend"
	hcore "github.com/ethereum/go-ethereum/consensus/hotstuff/core"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/signer"
	_ "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
)

type Geth struct {
	addr        common.Address
	miner       *miner
	chain       *core.BlockChain
	engine      Engine
	hotstuff    consensus.HotStuff
	api         *backend.API
	broadcaster *broadcaster
	signer      hotstuff.Signer
	hook        func(node *Geth, raw []byte) ([]byte, bool)
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
		signer:      signer.NewSigner(privateKey),
	}
	geth.addr = geth.signer.Address()
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
	time.Sleep(10 * time.Millisecond)
	g.hotstuff.Stop()
	g.miner.Stop()
	g.chain.Stop()
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
}

func (g *Geth) setHook(hook func(node *Geth, data []byte) ([]byte, bool)) {
	g.hook = hook
}

func (g *Geth) resignMsg(msg *hcore.Message) ([]byte, error) {
	hash, err := msg.Hash()
	if err != nil {
		return nil, err
	}
	sig, err := g.signer.SignHash(hash)
	if err != nil {
		return nil, err
	}
	msg.Signature = sig
	return msg.Payload()
}

type System struct {
	nodes []*Geth
	exit  chan struct{}
}

func makeSystem(n int) *System {
	pks, addrs := newAccountLists(n)
	nodes := make([]*Geth, n)

	for i := 0; i < n; i++ {
		nodes[i] = MakeGeth(pks[i], addrs)
	}

	return &System{nodes: nodes, exit: make(chan struct{})}
}

func (s *System) Start() {
	for i := 0; i < len(s.nodes); i++ {
		for j := 0; j < len(s.nodes); j++ {
			if j > i {
				s.nodes[i].broadcaster.Connect(s.nodes[j].broadcaster)
			}
		}
	}

	for _, node := range s.nodes {
		go node.Start()
	}

	go func() {
		for {
			select {
			case <-s.exit:
				for _, cli := range s.nodes {
					cli.Stop()
				}
				log.Info("-----System stopped!-----")
				return
			}
		}
	}()
}

func (s *System) Stop() {
	close(s.exit)
}

func (s *System) Close(n int) {
	timer := time.NewTimer(time.Duration(n) * time.Second)
	select {
	case <-timer.C:
		s.Stop()
		time.Sleep(1 * time.Second)
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
