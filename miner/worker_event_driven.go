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

package miner

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

type eventDrivenWorker struct {
}

func newEventDrivenWorker() *eventDrivenWorker {
	return nil
}

func (w *eventDrivenWorker) Start() {}

func (w *eventDrivenWorker) Stop() {}

func (w *eventDrivenWorker) Close() {}

func (w *eventDrivenWorker) IsRunning() bool { return false }

func (w *eventDrivenWorker) SetExtra(extra []byte) {}

func (w *eventDrivenWorker) SetRecommitInterval(interval time.Duration) {}
func (w *eventDrivenWorker) Pending() (*types.Block, *state.StateDB)    { return nil, nil }

func (w *eventDrivenWorker) PendingBlock() *types.Block       { return nil }
func (w *eventDrivenWorker) SetEtherbase(addr common.Address) {}
func (w *eventDrivenWorker) EnablePreseal()                   {}
func (w *eventDrivenWorker) DisablePreseal()                  {}
func (w *eventDrivenWorker) SubscribePendingLogs(ch chan<- []*types.Log) event.Subscription {
	return nil
}
