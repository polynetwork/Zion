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

package eth

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

var (
	broadcastDuration   = 60 * time.Second // Time duration for broadcast static-nodes
	broadcastLastTime   = 24 * time.Hour   // Last time for broadcast static-nodes
	broadcastChCapacity = 10               // Capacity for broadcast channel, is a low frequency action
)

type haltTask struct {
	waiting chan struct{}
	done    chan struct{}
}

type nodeBroadcaster struct {
	handler *handler

	miner common.Address // miner address used to judge that whether this program need to broadcast static-nodes

	validators       map[common.Address]struct{}     // validators map for filter useless connection
	validatorsMu     sync.Mutex                      // mutex for validators map
	staticNodesCh    chan consensus.StaticNodesEvent // channel for listening `StaticNodeEvent`
	staticNodesSub   event.Subscription              // subscribe consensus `StaticNodeEvent`
	staticNodesMap   map[string]*enode.Node          // static nodes storage
	staticNodesMapMu sync.Mutex                      // mutex for static-nodes storage

	server staticNodeServer // interface of p2p server

	round         int
	haltCh        chan *haltTask // Signal for procedure halt
	quitBroadcast chan struct{}  // Signal for quit broadcasting
}

func newNodeBroadcaster(miner common.Address, manager staticNodeServer, handler *handler) *nodeBroadcaster {
	return &nodeBroadcaster{
		handler:       handler,
		miner:         miner,
		server:        manager,
		quitBroadcast: make(chan struct{}),
		haltCh:        make(chan *haltTask, 10),
	}
}

func (h *nodeBroadcaster) Start() {
	// only used for hotstuff
	handler := h.handler.engine.(consensus.Handler)
	h.staticNodesCh = make(chan consensus.StaticNodesEvent, broadcastChCapacity)
	h.staticNodesSub = handler.SubscribeNodes(h.staticNodesCh)
	go h.staticNodesBroadcastLoop()
}

func (h *nodeBroadcaster) Stop() {
	h.staticNodesSub.Unsubscribe() // quits staticNodesLoop
	close(h.quitBroadcast)
}

// staticNodesBroadcastLoop
func (h *nodeBroadcaster) staticNodesBroadcastLoop() {
	if h.server == nil || h.miner == common.EmptyAddress {
		return
	}

	for {
		select {
		case evt := <-h.staticNodesCh:
			if h.round > 0 {
				task := <-h.haltCh
				close(task.waiting)
				<-task.done
			}
			task := &haltTask{
				waiting: make(chan struct{}),
				done:    make(chan struct{}),
			}
			h.haltCh <- task
			go h.BroadcastNodes(evt.Validators, task)

		case <-h.staticNodesSub.Err():
			return

		case <-h.quitBroadcast:
			return
		}
	}
}

// BroadcastNodes will broadcast all validators' enode information in a fixed time,
// only working in fixed time.
func (h *nodeBroadcaster) BroadcastNodes(validators []common.Address, task *haltTask) {
	done := make(chan struct{})

	log.Debug("BroadcastNodes", "current round", h.round, "validators", validators)
	h.round += 1

	timer := time.NewTimer(broadcastDuration)
	time.AfterFunc(broadcastLastTime, func() {
		close(done)
	})

	defer func() {
		timer.Stop()
		task.done <- struct{}{}
	}()

	peers := h.server.Peers()
	local := h.server.LocalENode()

	// fulfill validators
	h.setValidators(validators)

	// do not need to broadcast static-nodes if miner is not validator
	if _, exist := h.validators[h.miner]; !exist {
		return
	}

	for {
		select {
		case <-timer.C:
			var (
				addinf   string
				addNodes = make([]*enode.Node, 0)
			)

			// only inbound connections' remote address contains its own p2p server listening tcp/udp port.
			for _, peer := range peers {
				if _, exist := h.isValidator(peer.Node()); exist && !peer.Inbound() {
					h.addStaticNode(peer.Node())
				}
			}

			// send all static-nodes to remote peer
			// todo(fuk): consider the packet size. todo, addInf change to id or address
			// todo(fuk): update log to local.ID
			for _, node := range h.staticNodesMap {
				addNodes = append(addNodes, node)
				addinf += fmt.Sprintf("%d,", node.TCP())
			}
			for _, peer := range h.handler.peers.peers {
				if err := peer.SendStaticNodes(local, addNodes); err != nil {
					log.Error("SendStaticNodes", "to", peer.Node().ID(), "err", err)
				}
			}
			log.Debug("SendStaticNodes", "local", local.TCP(), "add", addinf)

			// reset timer to start the next round
			timer.Reset(broadcastDuration)

		// task halt
		case <-task.waiting:
			return

		// task finished
		case <-done:
			return

		// system stop signal
		case <-h.quitBroadcast:
			return
		}
	}
}

// addStaticNode add node and retrieve the node id and existence
func (h *nodeBroadcaster) addStaticNode(node *enode.Node) (id string, exist bool) {
	h.staticNodesMapMu.Lock()
	defer h.staticNodesMapMu.Unlock()

	if h.staticNodesMap == nil {
		h.staticNodesMap = make(map[string]*enode.Node)
	}

	id = node.ID().String()
	if _, exist = h.staticNodesMap[id]; !exist {
		h.staticNodesMap[id] = node
	}
	return
}

// remStaticNode remove static node from node map and retrieve node id and existence
func (h *nodeBroadcaster) remStaticNode(node *enode.Node) (id string, exist bool) {
	h.staticNodesMapMu.Lock()
	defer h.staticNodesMapMu.Unlock()

	id = node.ID().String()
	if h.staticNodesMap == nil {
		return
	}

	if _, exist = h.staticNodesMap[id]; exist {
		delete(h.staticNodesMap, id)
	}
	return
}

func (h *nodeBroadcaster) isValidator(node *enode.Node) (addr common.Address, exist bool) {
	addr = crypto.PubkeyToAddress(*node.Pubkey())
	if h.validators == nil {
		return
	}
	_, exist = h.validators[addr]
	return
}

func (h *nodeBroadcaster) setValidators(validators []common.Address) {
	h.validatorsMu.Lock()
	h.validators = make(map[common.Address]struct{})
	for _, v := range validators {
		h.validators[v] = struct{}{}
	}
	h.validatorsMu.Unlock()
}

// handleStaticNodesMsg is invoked from a peer's message handler when it transmits a
// static-nodes broadcast for the local node to process.
func (h *ethHandler) handleStaticNodesMsg(peer *eth.Peer, packet *eth.StaticNodesPacket) error {
	broadcaster := (*handler)(h).nodeBroadcaster

	if broadcaster.validators == nil || len(broadcaster.validators) < 1 {
		return nil
	}

	// filter validator's message
	if _, exist := broadcaster.isValidator(peer.Node()); !exist {
		log.Debug("handleStaticNodesMsg", "from", peer.Node().TCP(), "err", "node is not validator",
			"list", broadcaster.validators)
		return nil
	}

	// `packet.Local` is the `urlv4` string of the message sender, which needs to be parsed into a `enode` structure.
	// in addition, this node must be placed at the end of the `packet.remote` list, because there is deduplication
	// in the process of traversing the list. take the pre-existing node that does not need to be resolved.
	from := packet.Local
	node, err := enode.CopyUrlv4(from.URLv4(), from.IP(), from.TCP(), from.UDP())
	if err != nil {
		return err
	}
	packet.Remotes = append(packet.Remotes, node)

	// collect all static nodes
	for _, node := range packet.Remotes {
		log.Debug("handleStaticNodesMsg", "from", peer.Node().TCP(), "receive tcp", node.TCP())
		if _, exist := broadcaster.isValidator(node); exist {
			broadcaster.addStaticNode(node)
			log.Debug("handleStaticNodesMsg", "from", peer.Node().TCP(), "add tcp", node.TCP())
		}
	}

	// filter old validators
	olds := make(map[enode.ID]*enode.Node)
	for _, node := range broadcaster.staticNodesMap {
		if _, exist := broadcaster.isValidator(node); !exist {
			olds[node.ID()] = node
			broadcaster.remStaticNode(node)
			log.Debug("handleStaticNodesMsg", "from", peer.Node().TCP(), "remove old tcp", node.TCP())
		}
	}

	// add new peer for p2p server
	for _, node := range broadcaster.staticNodesMap {
		if _, exist := h.peers.peers[node.ID().String()]; !exist {
			broadcaster.server.AddPeer(node)
		}
	}

	// remove old peer in p2p server
	for _, node := range olds {
		if _, exist := h.peers.peers[node.ID().String()]; exist {
			broadcaster.server.RemovePeer(node)
		}
	}

	return nil
}
