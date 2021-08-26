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
package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
)

type EventEmitter struct {
	contract common.Address
	state    *state.StateDB
	block    uint64
}

func NewEventEmitter(contract common.Address, block uint64, state *state.StateDB) *EventEmitter {
	emitter := &EventEmitter{
		contract: contract,
		state:    state,
		block:    block,
	}
	return emitter
}

func (emitter *EventEmitter) Event(topics []common.Hash, data []byte) {
	event := &types.Log{
		Address:     emitter.contract,
		Topics:      topics,
		Data:        data,
		BlockNumber: emitter.block,
	}
	emitter.state.AddLog(event)
}
