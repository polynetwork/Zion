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
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

type miner struct {
	addr   common.Address
	chain  *core.BlockChain
	engine consensus.HotStuff
	geth   *Geth

	current *environment

	headCh       chan core.ChainHeadEvent
	chainHeadSub event.Subscription

	executedCh  chan consensus.ExecutedBlock
	executedSub event.Subscription

	nodesCh  chan consensus.StaticNodesEvent
	nodesSub event.Subscription

	pendingMu    sync.RWMutex
	pendingTasks map[common.Hash]*task

	exit chan struct{}
}

type environment struct {
	state    *state.StateDB
	header   *types.Header
	receipts types.Receipts
	logs     []*types.Log
}

type task struct {
	receipts []*types.Receipt
	state    *state.StateDB
	block    *types.Block
}

func makeMiner(address common.Address, chain *core.BlockChain, engine consensus.HotStuff) *miner {
	miner := &miner{
		addr:         address,
		chain:        chain,
		engine:       engine,
		headCh:       make(chan core.ChainHeadEvent, 1),
		nodesCh:      make(chan consensus.StaticNodesEvent, 1),
		executedCh:   make(chan consensus.ExecutedBlock, 1),
		pendingTasks: make(map[common.Hash]*task),
		exit:         make(chan struct{}),
	}

	handler := engine.(consensus.Handler)
	miner.chainHeadSub = chain.SubscribeChainHeadEvent(miner.headCh)
	miner.nodesSub = handler.SubscribeNodes(miner.nodesCh)
	miner.executedSub = handler.SubscribeBlock(miner.executedCh)
	return miner
}

func (m *miner) Start() {
	timer := time.NewTimer(0 * time.Second)

	for {
		select {
		case data := <-m.headCh:
			if h, ok := m.engine.(consensus.Handler); ok {
				h.NewChainHead(data.Block.Header())
			}
			m.newWork()

		case <-timer.C:
			m.newWork()
			timer.Reset(2 * time.Second)

		case data := <-m.executedCh:
			m.commit(&data)

			// ensure that backend nodes feed wont be blocked.
		case <-m.nodesCh:

		case <-m.exit:
			return
		}
	}
}

func (m *miner) Stop() {
	close(m.exit)
}

func (m *miner) newWork() {
	parent := m.chain.CurrentHeader().Copy()
	num := parent.Number
	timestamp := time.Now().Unix()
	log.Debug("Parent header", "hash", parent.Hash(), "num", num)

	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     num.Add(num, common.Big1),
		GasLimit:   math.MaxUint64,
		Time:       uint64(timestamp),
	}
	m.makeCurrent(header)

	// DONT read epoch info from native contracts, but use fixed big value as the first epoch. all mock tests running in this epoch.
	if err := types.HotstuffHeaderFillWithValidators(header, nil, EpochStart, EpochEnd); err != nil {
		log.Error("Failed to fill header", "err", err)
		return
	}

	if err := m.engine.Prepare(m.chain, header); err != nil {
		log.Error("Failed to prepare", "err", err)
		return
	}

	s := m.current.state.Copy()
	block, receipts, err := m.engine.FinalizeAndAssemble(m.chain, m.current.header, s, nil, nil, nil)
	if err != nil {
		log.Error("Failed to finalizeAndAssemble", "err", err)
		return
	}

	sealHash := m.engine.SealHash(block.Header())
	m.pendingMu.Lock()
	task := &task{receipts: receipts, state: s, block: block}
	m.pendingTasks[sealHash] = task
	m.pendingMu.Unlock()

	if err := m.engine.Seal(m.chain, task.block, nil, nil); err != nil {
		log.Error("Block sealing failed", "err", err)
		return
	}
}

func (m *miner) commit(data *consensus.ExecutedBlock) {
	block := data.Block
	if block == nil {
		log.Warn("Executed block is nil")
		return
	}
	if m.chain.HasBlock(block.Hash(), block.NumberU64()) {
		log.Debug("Block already exist", block.Number(), block.Hash())
		return
	}

	var (
		sealhash = m.engine.SealHash(block.Header())
		hash     = block.Hash()
		receipts []*types.Receipt
		logs     []*types.Log
		statedb  *state.StateDB
		task     *task
		exist    bool
	)

	if data.State != nil {
		receipts = data.Receipts
		logs = data.Logs
		statedb = data.State
	} else {
		m.pendingMu.RLock()
		task, exist = m.pendingTasks[sealhash]
		m.pendingMu.RUnlock()
		if !exist {
			log.Error("Failed to find local task", "hash", block.Hash())
			return
		}

		receipts = make([]*types.Receipt, len(task.receipts))
		statedb = task.state
		for i, receipt := range task.receipts {
			// add block location fields
			receipt.BlockHash = hash
			receipt.BlockNumber = block.Number()
			receipt.TransactionIndex = uint(i)

			receipts[i] = new(types.Receipt)
			*receipts[i] = *receipt
			// Update the block hash in all logs since it is now available and not when the
			// receipt/log of individual transactions were created.
			for _, log := range receipt.Logs {
				log.BlockHash = hash
			}
			logs = append(logs, receipt.Logs...)
		}
	}

	// Commit block and state to database.
	_, err := m.chain.WriteBlockWithState(block, receipts, logs, statedb, true)
	if err != nil {
		log.Error("Failed writing block to chain", "err", err)
		return
	}
	log.Info("Successfully sealed new block", "address", m.addr, "number", block.Number(), "sealhash", sealhash, "hash", hash)

	if m.geth != nil {
		go m.geth.broadcastBlock(block)
	}
}

func (m *miner) makeCurrent(header *types.Header) {
	block := m.chain.CurrentBlock().Copy()
	statedb, _ := m.chain.StateAt(block.Root())
	m.current = &environment{
		header: header,
		state:  statedb,
	}
}
