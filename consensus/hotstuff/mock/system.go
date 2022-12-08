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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/backend"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/crypto"
)

type Geth struct {
	miner  *miner
	chain  *core.BlockChain
	engine Engine
	api    *backend.API
}

func MakeGeth(privateKey *ecdsa.PrivateKey, vals []common.Address) *Geth {
	db := rawdb.NewMemoryDatabase()
	engine := makeEngine(privateKey, db)
	chain := makeChain(db, engine, vals)
	miner := makeMiner(chain, engine.(consensus.HotStuff))
	api := engine.APIs(chain)[0].Service.(*backend.API)
	return &Geth{
		miner:  miner,
		chain:  chain,
		engine: engine,
		api:    api,
	}
}

func (g *Geth) Start() {
	g.engine.(consensus.HotStuff).Start(g.chain, nil)
	g.miner.Start()
}

func (g *Geth) Stop() {
	g.engine.(consensus.Handler).GetBroadcaster().Stop()
	g.engine.(consensus.HotStuff).Stop()
	g.miner.Stop()
}

func (g *Geth) Sequence() (uint64, uint64) {
	return g.api.CurrentSequence()
}

func (g *Geth) IsProposer() bool {
	return g.api.IsProposer()
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
