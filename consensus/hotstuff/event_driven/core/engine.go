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

package core

import (
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

var once sync.Once

func (c *core) Start(chain consensus.ChainReader) error {
	once.Do(func() {
		hotstuff.RegisterMsgTypeConvertHandler(convertUpMsgType)
	})

	c.chain = chain
	if err := c.initialize(); err != nil {
		return err
	}

	// todo: start all nodes at the same time
	time.Sleep(15 * time.Second)

	// Tests will handle events itself, so we have to make subscribeEvents()
	// be able to call in test.
	c.subscribeEvents()
	go c.handleEvents()

	// engine is started after this step, DONT allow to return err to miner worker, this may cause worker invalid
	c.started = true
	highQC := c.smr.HighQC()
	c.advanceRoundByQC(highQC)
	return nil
}

func (c *core) Stop() error {
	c.stopTimer()
	c.unsubscribeEvents()
	c.started = false
	return nil
}

func (c *core) IsProposer() bool {
	if c.valset.IsProposer(c.address) {
		return true
	}
	return false
}

func (c *core) Address() common.Address {
	return c.address
}

// verify if a hash is the same as the proposed block in the current pending request
//
// this is useful when the engine is currently the speaker
//
// pending request is populated right at the request stage so this would give us the earliest verification
// to avoid any race condition of coming propagated blocks
// 判断是否已经提交或者正在提交, 这样一来，request必须在一开始就写入到blockTree
func (c *core) IsCurrentProposal(blockHash common.Hash) bool {
	block := c.blkPool.GetBlockByHash(blockHash)
	if block == nil {
		return false
	}

	if block.NumberU64() != c.smr.Height().Uint64() {
		return false
	}
	return true
}

func (c *core) PrepareExtra(header *types.Header, valSet hotstuff.ValidatorSet) ([]byte, error) {
	return generateExtra(header, valSet, c.smr.Epoch(), c.smr.Round())
}

func (c *core) GetHeader(hash common.Hash, number uint64) *types.Header {
	block := c.blkPool.GetBlockAndCheckHeight(hash, new(big.Int).SetUint64(number))
	if block == nil {
		return nil
	} else {
		return block.Header()
	}
}

func (c *core) SubscribeRequest(ch chan<- consensus.AskRequest) event.Subscription {
	return c.feed.Subscribe(ch)
}

func (c *core) InitValidators(valset hotstuff.ValidatorSet) {
	c.valset = valset
	c.messages = NewMessagePool(valset)
}

func (c *core) ChangeEpoch(epochStartHeight uint64, valset hotstuff.ValidatorSet) error {
	return nil
}
