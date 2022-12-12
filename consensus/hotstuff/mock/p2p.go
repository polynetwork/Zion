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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
)

type broadcaster struct {
	addr  common.Address
	eng   Engine
	peers map[common.Address]*MockPeer
	geth  *Geth
}

func makeBroadcaster(addr common.Address, engine Engine) *broadcaster {
	return &broadcaster{
		addr:  addr,
		eng:   engine,
		peers: make(map[common.Address]*MockPeer),
	}
}

func (b *broadcaster) FindPeers(targets map[common.Address]bool) map[common.Address]consensus.Peer {
	m := make(map[common.Address]consensus.Peer)
	for addr, p := range b.peers {
		if targets[addr] {
			m[addr] = p
		}
	}
	return m
}

func (b *broadcaster) FindPeer(target common.Address) consensus.Peer {
	for addr, p := range b.peers {
		if addr == target {
			return p
		}
	}
	return nil
}

func (b *broadcaster) Stop() {
	for _, peer := range b.peers {
		peer.Close()
	}
}

func (b *broadcaster) Connect(b2 *broadcaster) {
	if _, exist := b.peers[b2.addr]; exist {
		return
	}

	rw1, rw2 := p2p.MsgPipe()
	b.add(b2.addr, rw1)
	b2.add(b.addr, rw2)
}

func (b *broadcaster) add(remote common.Address, rw *p2p.MsgPipeRW) {
	peer := &MockPeer{rw: rw, local: b.addr, remote: remote, geth: b.geth}
	b.peers[remote] = peer
	handler := b.eng.(consensus.Handler)

	log.Debug("connect", "local", b.addr, "remote", remote)
	go func() {
		defer peer.Close()

		for {
			msg, err := peer.ReadMsg()
			if err != nil {
				log.Error("Failed to read message", "err", err)
				return
			}
			if _, err := handler.HandleMsg(remote, msg); err != nil {
				log.Error("Failed to handle message", "err", err)
				return
			}
			if msg.Code == eth.NewBlockMsg {
				b.geth.handleBlock(msg)
			}
		}
	}()
}

type MockPeer struct {
	local, remote common.Address
	rw            *p2p.MsgPipeRW
	geth          *Geth
}

func (p *MockPeer) SendNewBlock(block *types.Block, td *big.Int) error {
	return p2p.Send(p.rw, eth.NewBlockMsg, &eth.NewBlockPacket{
		Block: block,
		TD:    td,
	})
}

func (p *MockPeer) Send(msgcode uint64, data interface{}) error {
	send := true

	if p.geth.hook != nil && msgcode == hotstuffMsg {
		if raw, ok := data.([]byte); !ok {
			panic("Send hotstuff message data convert failed")
		} else {
			data, send = p.geth.hook(p.geth, raw)
		}
	}
	if send {
		if err := p2p.Send(p.rw, msgcode, data); err != nil {
			log.Error("Failed to send msg", "local", p.local, "remote", p.remote, "err", err)
			return err
		}
	}

	return nil
}

func (p *MockPeer) ReadMsg() (p2p.Msg, error) {
	return p.rw.ReadMsg()
}

func (p *MockPeer) Close() error {
	return p.rw.Close()
}
