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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"math/big"
	"sync"
	"time"
)

var once sync.Once

func (e *core) Start() error {
	once.Do(func() {
		hotstuff.RegisterMsgTypeConvertHandler(func(data interface{}) hotstuff.MsgType {
			code := data.(uint64)
			return MsgType(code)
		})
	})

	if err := e.initialize(); err != nil {
		return err
	}

	time.Sleep(8 * time.Second)

	// Tests will handle events itself, so we have to make subscribeEvents()
	// be able to call in test.
	e.subscribeEvents()
	go e.handleEvents()

	// engine is started after this step, DONT allow to return err to miner worker, this may cause worker invalid
	e.started = true
	highQC := e.blkPool.GetHighQC()
	e.advanceRoundByQC(highQC, false)
	return nil
}

func (e *core) Stop() error {
	e.stopTimer()
	e.unsubscribeEvents()
	e.started = false
	return nil
}

func (e *core) IsProposer() bool {
	if e.valset.IsProposer(e.address()) {
		return true
	}
	return false
}

func (e *core) Address() common.Address {
	return e.signer.Address()
}

// verify if a hash is the same as the proposed block in the current pending request
//
// this is useful when the engine is currently the speaker
//
// pending request is populated right at the request stage so this would give us the earliest verification
// to avoid any race condition of coming propagated blocks
// 判断是否已经提交或者正在提交, 这样一来，request必须在一开始就写入到blockTree
func (e *core) IsCurrentProposal(blockHash common.Hash) bool {
	block := e.blkPool.GetBlockByHash(blockHash)
	if block == nil {
		return false
	}

	if block.NumberU64() != e.curHeight.Uint64() {
		return false
	}

	return true
}

func (e *core) PrepareExtra(header *types.Header, valSet hotstuff.ValidatorSet) ([]byte, error) {
	return generateExtra(header, valSet, e.epoch, e.curRound)
}

func (e *core) GetHeader(hash common.Hash, number uint64) *types.Header {
	block := e.blkPool.GetBlockAndCheckHeight(hash, new(big.Int).SetUint64(number))
	if block == nil {
		return nil
	} else {
		return block.Header()
	}
}

func (e *core) SubscribeRequest(ch chan<- consensus.AskRequest) event.Subscription {
	return e.feed.Subscribe(ch)
}
