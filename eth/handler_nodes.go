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
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

// hotstuff is the master-slave network, so that unconnected nodes in devp2p cannot receive
// consensus messages. It is necessary to implement an new protocol to achieve direct network
// connection of validator nodes.
//
// implement an loop to monitor and process the `epochChange` event, miner should parse and
// resemble it's own enode info and broadcast it and remote enodes to other validators. After
// the node receives node info, it adds connections and deletes unnecessary links.
//

var (
	// // todo(fuk): update parameters as below, it should be a bit longer in mainnet.
	broadcastDuration   = 3 * time.Second // Time duration for broadcast static-nodes
	broadcastLastTime   = 10 * time.Minute  // Last time for broadcast static-nodes
	broadcastChCapacity = 10              // Capacity for broadcast channel, is a low frequency action
)

// staticNodeServer defines the methods need from a p2p server implementation to
// support operation needed by the `hotstuff` protocol.
type staticNodeServer interface {
	PeersInfo() []*p2p.PeerInfo
	Peers() []*p2p.Peer
	AddPeer(node *enode.Node)
	RemovePeer(node *enode.Node)
	LocalENode() *enode.Node
}

type haltTask struct {
	waiting chan struct{}
	done    chan struct{}
}

type nodeBroadcaster struct {
	handler *handler

	miner common.Address // miner address used to judge that whether this program need to broadcast static-nodes

	validators   map[common.Address]struct{}     // validators map for filter useless connection
	validatorsMu sync.Mutex                      // mutex for validators map
	nodesCh      chan consensus.StaticNodesEvent // channel for listening `StaticNodeEvent`
	nodesSub     event.Subscription              // subscribe consensus `StaticNodeEvent`
	nodesMap     map[string]*enode.Node          // static nodes storage
	nodesMapMu   sync.Mutex                      // mutex for static-nodes storage

	server staticNodeServer // interface of p2p server

	round  int
	haltCh chan *haltTask // Signal for procedure halt
	quit   chan struct{}  // Signal for quit broadcasting
}

func newNodeBroadcaster(miner common.Address, manager staticNodeServer, handler *handler) *nodeBroadcaster {
	return &nodeBroadcaster{
		handler: handler,
		miner:   miner,
		server:  manager,
		quit:    make(chan struct{}),
		haltCh:  make(chan *haltTask, 10),
	}
}

// only used for hotstuff
func (h *nodeBroadcaster) Start() {
	handler := h.handler.engine.(consensus.Handler)
	h.nodesCh = make(chan consensus.StaticNodesEvent, broadcastChCapacity)
	h.nodesSub = handler.SubscribeNodes(h.nodesCh)
	go h.loop()
}

func (h *nodeBroadcaster) Stop() {
	h.nodesSub.Unsubscribe() // quits staticNodesLoop
	close(h.quit)
}

// loop
func (h *nodeBroadcaster) loop() {
	if h.server == nil || h.miner == common.EmptyAddress {
		return
	}

	for {
		select {
		case evt := <-h.nodesCh:
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
			go h.handleTask(evt.Validators, task)

		case <-h.nodesSub.Err():
			return

		case <-h.quit:
			return
		}
	}
}

// handleTask will broadcast all validators' enode information in fixed time
func (h *nodeBroadcaster) handleTask(validators []common.Address, task *haltTask) {
	done := make(chan struct{})

	log.Debug("handleTask", "current round", h.round, "validators", validators)
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
	//if _, exist := h.validators[h.miner]; !exist {
	//	return
	//}

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
					h.addNode(peer.Node())
				}
			}

			// send all static-nodes to remote peer
			// todo(fuk): consider the packet size. todo, addInf change to id or address
			// todo(fuk): update log to local.ID
			for _, node := range h.nodesMap {
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
		case <-h.quit:
			return
		}
	}
}

// addNode add node and retrieve the node id and existence
func (h *nodeBroadcaster) addNode(node *enode.Node) (id string, exist bool) {
	h.nodesMapMu.Lock()
	defer h.nodesMapMu.Unlock()

	if h.nodesMap == nil {
		h.nodesMap = make(map[string]*enode.Node)
	}

	id = node.ID().String()
	if _, exist = h.nodesMap[id]; !exist {
		h.nodesMap[id] = node
		log.Trace("add node", "")
	}
	return
}

// delNode remove static node from node map and retrieve node id and existence
func (h *nodeBroadcaster) delNode(node *enode.Node) (id string, exist bool) {
	h.nodesMapMu.Lock()
	defer h.nodesMapMu.Unlock()

	id = node.ID().String()
	if h.nodesMap == nil {
		return
	}

	if _, exist = h.nodesMap[id]; exist {
		delete(h.nodesMap, id)
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
	// 允许老validator将新validator信息转发给节点
	//if _, exist := broadcaster.isValidator(peer.Node()); !exist {
	//	log.Trace("handleStaticNodesMsg", "from", peer.Node().TCP(), "err", "node is not validator",
	//		"list", broadcaster.validators)
	//	return nil
	//}

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
		log.Trace("handleStaticNodesMsg", "from", peer.Node().TCP(), "receive tcp", node.TCP())
		if _, exist := broadcaster.isValidator(node); exist {
			broadcaster.addNode(node)
			log.Trace("handleStaticNodesMsg", "from", peer.Node().TCP(), "add tcp", node.TCP())
		}
	}

	// add new peer for p2p server
	for _, node := range broadcaster.nodesMap {
		if _, exist := h.peers.peers[node.ID().String()]; !exist {
			broadcaster.server.AddPeer(node)
		}
	}

	// filter old validators
	//olds := make(map[enode.ID]*enode.Node)
	//for _, node := range broadcaster.nodesMap {
	//	if _, exist := broadcaster.isValidator(node); !exist {
	//		olds[node.ID()] = node
	//		broadcaster.delNode(node)
	//		log.Debug("handleStaticNodesMsg", "from", peer.Node().TCP(), "remove old tcp", node.TCP())
	//	}
	//}

	//// remove old peer in p2p server
	//for _, node := range olds {
	//	if _, exist := h.peers.peers[node.ID().String()]; exist {
	//		broadcaster.server.RemovePeer(node)
	//	}
	//}

	return nil
}
