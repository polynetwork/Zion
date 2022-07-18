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
	"sync"
	"sync/atomic"
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
	nodeFetcherDuration   = 2 * time.Second
	nodeFetchingLastTime  = 1 * time.Minute
	nodeFetcherChCapacity = 10 // Capacity for broadcast channel, is a low frequency action
)

// staticNodeServer defines the methods need from a p2p server implementation to
// support operation needed by the `hotstuff` protocol.
type staticNodeServer interface {
	PeersInfo() []*p2p.PeerInfo
	Peers() []*p2p.Peer
	AddPeer(node *enode.Node)
	RemovePeer(node *enode.Node)
	Self() *enode.Node
	SeedNodes() []*enode.Node
}

type task struct {
	validators []common.Address
	halt       chan struct{}
	done       chan struct{}
}

type nodeFetcher struct {
	handler *handler         // eth handler
	server  staticNodeServer // interface of p2p server
	logger  log.Logger

	miner common.Address              // miner address used to judge that whether this program need to broadcast static-nodes
	local *enode.Node                 // local node
	seed  int32                       // flag of whether the node is seed node
	seeds map[common.Address]*ethPeer // seed node peer connections, identity is validator address

	validators   map[common.Address]*enode.Node  // validators map for filter useless connection
	validatorsMu sync.RWMutex                    // mutex for validators map
	notifyCh     chan consensus.StaticNodesEvent // channel for listening `StaticNodeEvent`
	notifySub    event.Subscription              // subscribe consensus `StaticNodeEvent`

	taskCh chan *task    // Signal for procedure halt
	quit   chan struct{} // Signal for quit broadcasting
}

func newNodeBroadcaster(miner common.Address, manager staticNodeServer, handler *handler) *nodeFetcher {
	return &nodeFetcher{
		handler:    handler,
		miner:      miner,
		server:     manager,
		logger:     log.New("address", miner.Hex()),
		validators: make(map[common.Address]*enode.Node),
		seeds:      make(map[common.Address]*ethPeer),
		taskCh:     make(chan *task, nodeFetcherChCapacity),
		quit:       make(chan struct{}),
	}
}

// only used for hotstuff
func (h *nodeFetcher) Start() {
	handler := h.handler.engine.(consensus.Handler)
	h.notifyCh = make(chan consensus.StaticNodesEvent, nodeFetcherChCapacity)
	h.notifySub = handler.SubscribeNodes(h.notifyCh)
	h.local = h.server.Self()
	go h.loop()
}

func (h *nodeFetcher) Stop() {
	if h.seed == 1 {
		atomic.StoreInt32(&h.seed, 0)
	}
	h.notifySub.Unsubscribe() // quits staticNodesLoop
	close(h.quit)
}

func (h *nodeFetcher) loop() {
	for {
		select {
		case evt := <-h.notifyCh:
			h.waitingLastTask()
			task := h.newTask(evt.Validators)
			go h.handleTask(task)

		case <-h.notifySub.Err():
			return

		case <-h.quit:
			h.logger.Trace("Node Fetcher end loop")
			return
		}
	}
}

func (h *nodeFetcher) newTask(validators []common.Address) *task {
	task := &task{
		validators: validators,
		halt:       make(chan struct{}),
		done:       make(chan struct{}),
	}
	h.taskCh <- task
	return task
}

func (h *nodeFetcher) waitingLastTask() {
	if len(h.taskCh) > 0 {
		h.logger.Trace("Node Fetcher waiting for last task end...")
		task := <-h.taskCh
		close(task.halt)
		<-task.done
		h.logger.Trace("Node Fetcher last task stopped")
	}
}

func (h *nodeFetcher) handleTask(task *task) {
	timer := time.NewTimer(0)
	done := make(chan struct{})

	time.AfterFunc(nodeFetchingLastTime, func() {
		close(done)
	})
	defer func() {
		// todo: remove all nodes which is not validator
		timer.Stop()
		task.done <- struct{}{}
	}()

	// fulfill validators
	h.resetValidators(task.validators)

	// sync node do not allow to ask validator address, and seed node has already persisted other seeds.
	if h.miner == common.EmptyAddress || !h.checkValidator(h.miner) || h.isSeedNode() {
		h.logger.Trace("Node Fetcher authority invalid")
		return
	}

	h.logger.Trace("Node Fetcher handle new task", "miner", h.miner, "validators", task.validators)

	for {
		select {
		case <-timer.C:
			if h.connectionQuorum() {
				h.logger.Trace("Node Fetcher full connected!")
				return
			}
			h.fullConnectSeedNodes()
			h.batchRequest()
			timer.Reset(nodeFetcherDuration)

		case <-task.halt:
			h.logger.Trace("Node Fetcher task halt")
			return

		case <-done:
			h.logger.Trace("Node Fetcher task done")
			return

		case <-h.quit:
			h.logger.Trace("Node Fetcher task quit")
			return
		}
	}
}

func (h *nodeFetcher) fullConnectSeedNodes() {
	// 如果种子节点已经全连接，则无需再重新连接
	if len(h.server.SeedNodes()) == len(h.seeds) {
		return
	}

	for _, seed := range h.server.SeedNodes() {
		addr := nodeAddress(seed)
		if peer := h.handler.FindPeer(addr); peer == nil {
			h.server.AddPeer(seed)
		} else {
			h.seeds[addr] = peer.(*ethPeer)
			h.logger.Trace("Node Fetcher add seed", "address", addr.Hex())
		}
	}
}

func (h *nodeFetcher) connectionQuorum() bool {
	if h.validators == nil {
		return false
	}
	for addr, v := range h.validators {
		if v == nil || h.handler.FindPeer(addr) == nil {
			return false
		}
	}
	return true
}

func (h *nodeFetcher) batchRequest() {
	for _, peer := range h.seeds {
		if err := peer.RequestStaticNodes(h.local); err != nil {
			h.logger.Trace("Node Fetcher", "request static node", identity(peer.Node()), "err", err)
		} else {
			h.logger.Trace("Node Fetcher", "request static node", identity(peer.Node()), "local", identity(h.local))
		}
	}
}

func (h *nodeFetcher) handleGetStaticNodesMsg(peer *eth.Peer, from *enode.Node) error {
	logger := h.logger.New("handle", "GetStaticNodesMsg", "local", identity(h.local), "from", identity(from))
	logger.Trace("Node Fetcher")

	if !h.isSeedNode() {
		logger.Trace("Failed to check seed authority", "err", "node is not seed node")
		return nil
	}

	validator := nodeAddress(peer.Node())
	if !h.checkValidator(validator) {
		logger.Trace("Failed to check validator", "err", "node is not validator")
		return nil
	}

	// todo: ensure that peer.node pubkey and ip is the same with `from` node
	node, err := enode.CopyUrlv4(from.URLv4(), from.IP(), from.TCP(), from.UDP())
	if err != nil {
		logger.Trace("Failed to regenerate new remote node", "err", err)
		return nil
	}
	h.setValidator(validator, node)

	list := h.validatorList()
	if err := peer.ReplyGetStaticNodes(list); err != nil {
		logger.Trace("Failed to reply `GetStaticNodes`", "err", err)
		return nil
	} else {
		logger.Trace("Reply `GetStaticNodes` succeed!", "length", len(list))
	}
	return nil
}

func (h *nodeFetcher) handleStaticNodesMsg(peer *eth.Peer, list []*enode.Node) error {
	logger := h.logger.New("handle", "StaticNodesMsg", "local", identity(h.local), "from", identity(peer.Node()))
	logger.Trace("Node Fetcher")

	if list == nil || len(list) == 0 {
		logger.Trace("Failed to check message packet", "err", "node list is empty")
		return nil
	}

	for _, node := range list {
		validator := nodeAddress(node)
		h.server.AddPeer(node)
		if h.setValidator(validator, node) {
		}
	}

	return nil
}

func (h *nodeFetcher) resetValidators(validators []common.Address) {
	h.validatorsMu.Lock()
	defer h.validatorsMu.Unlock()

	h.validators = make(map[common.Address]*enode.Node)
	for _, v := range validators {
		h.validators[v] = nil
	}
	h.logger.Trace("Node Fetcher", "reset validators", validators)
}

func (h *nodeFetcher) setValidator(validator common.Address, node *enode.Node) bool {
	h.validatorsMu.Lock()
	defer h.validatorsMu.Unlock()

	if data, exist := h.validators[validator]; exist && data == nil {
		h.validators[validator] = node
		h.logger.Trace("Node Fetcher", "set validator node", validator.Hex(), "node", identity(node))
		return true
	} else {
		h.logger.Trace("Node Fetcher", "set validator node", "failed", "exist", exist, "data == nil", data == nil)
	}
	return false
}

// checkValidator return true if the address exist in validator map
func (h *nodeFetcher) checkValidator(addr common.Address) bool {
	h.validatorsMu.RLock()
	defer h.validatorsMu.RUnlock()

	if h.validators == nil {
		return false
	}
	_, exist := h.validators[addr]
	return exist
}

func (h *nodeFetcher) validatorList() []*enode.Node {
	h.validatorsMu.RLock()
	defer h.validatorsMu.RUnlock()

	list := make([]*enode.Node, 0)
	for _, v := range h.validators {
		if v != nil {
			list = append(list, v)
		}
	}
	return list
}

func (h *nodeFetcher) isSeedNode() bool {
	return h.seed == 1
}

func nodeAddress(node *enode.Node) common.Address {
	return crypto.PubkeyToAddress(*node.Pubkey())
}

func identity(node *enode.Node) interface{} {
	if node == nil {
		return 0
	}
	return node.TCP()
}
