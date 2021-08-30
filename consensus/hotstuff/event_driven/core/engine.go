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
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
)

func (e *EventDrivenEngine) Start() error {
	e.handleNewRound()

	// Tests will handle events itself, so we have to make subscribeEvents()
	// be able to call in test.
	e.subscribeEvents()
	go e.handleEvents()
	return nil
}

func (e *EventDrivenEngine) Stop() error {
	e.stopTimer()
	e.unsubscribeEvents()
	return nil
}

func (e *EventDrivenEngine) IsProposer() bool {
	if e.valset.IsProposer(e.address()) {
		return true
	}
	return false
}

// verify if a hash is the same as the proposed block in the current pending request
//
// this is useful when the engine is currently the speaker
//
// pending request is populated right at the request stage so this would give us the earliest verification
// to avoid any race condition of coming propagated blocks
// 判断是否已经提交或者正在提交, 这样一来，request必须在一开始就写入到blockTree
func (e *EventDrivenEngine) IsCurrentProposal(blockHash common.Hash) bool {
	block := e.blkPool.GetBlockByHash(blockHash)
	if block == nil {
		return false
	}

	if block.NumberU64() != e.curHeight.Uint64() {
		return false
	}

	return true
}

func (e *EventDrivenEngine) PrepareExtra(header *types.Header, valSet hotstuff.ValidatorSet) ([]byte, error) {
	return generateExtra(header, valSet, e.epoch, e.curRound)
}

func (e *EventDrivenEngine) Address() common.Address {
	return e.signer.Address()
}
