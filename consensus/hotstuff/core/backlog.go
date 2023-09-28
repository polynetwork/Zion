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
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/prque"
)

func (c *core) storeBacklog(msg *Message) {
	logger := c.newLogger()

	src := msg.address
	if src == c.Address() {
		logger.Trace("Backlog from self")
		return
	}
	if _, v := c.valSet.GetByAddress(src); v == nil {
		logger.Trace("Backlog from unknown validator", "address", src)
		return
	}

	logger.Trace("Retrieving backlog queue", "msg", msg.Code, "src", src, "backlogs_size", c.backlogs.Size(src))

	c.backlogs.Push(msg)
}

func (c *core) processBacklog() {
	logger := c.newLogger()

	c.backlogs.mu.Lock()
	defer c.backlogs.mu.Unlock()

	for addr, queue := range c.backlogs.queue {
		if queue == nil {
			continue
		}
		_, src := c.valSet.GetByAddress(addr)
		if src == nil {
			logger.Trace("Skip the backlog", "unknown validator", addr)
			continue
		}

		isFuture := false
		for !(queue.Empty() || isFuture) {
			data, priority := queue.Pop()
			msg, ok := data.(*Message)
			if !ok {
				logger.Trace("Skip the backlog, invalid Message")
				continue
			}
			if err := c.checkView(msg.View); err != nil {
				if err == errFutureMessage {
					queue.Push(data, priority)
					isFuture = true
					break
				}
				logger.Trace("Skip the backlog", "msg view", msg.View, "err", err)
				continue
			}

			logger.Trace("Replay the backlog", "msg", msg)
			go c.backlogFeed.Send(backlogEvent{src: src, msg: msg})
		}
	}
}

type backlog struct {
	mu    *sync.RWMutex
	queue map[common.Address]*prque.Prque
}

func newBackLog() *backlog {
	return &backlog{
		mu:    new(sync.RWMutex),
		queue: make(map[common.Address]*prque.Prque),
	}
}

func (b *backlog) Push(msg *Message) {
	if msg == nil || msg.address == common.EmptyAddress ||
		msg.View == nil || msg.Code > MsgTypeDecide {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	addr := msg.address
	if _, ok := b.queue[addr]; !ok {
		b.queue[addr] = prque.New(nil)
	}
	priority := b.toPriority(msg.Code, msg.View)
	b.queue[addr].Push(msg, priority)
}

func (b *backlog) Pop(addr common.Address) (data *Message, priority int64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.queue[addr]; !ok {
		return
	} else {
		item, p := b.queue[addr].Pop()
		data = item.(*Message)
		priority = p
		return
	}
}

func (b *backlog) Size(addr common.Address) int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if que, ok := b.queue[addr]; !ok {
		return 0
	} else {
		return que.Size()
	}
}

var messagePriorityTable = map[MsgType]int64{
	MsgTypeNewView:       1,
	MsgTypePrepare:       2,
	MsgTypePrepareVote:   3,
	MsgTypePreCommit:     4,
	MsgTypePreCommitVote: 5,
	MsgTypeCommit:        6,
	MsgTypeCommitVote:    7,
	MsgTypeDecide:        8,
}

func (b *backlog) toPriority(msgCode MsgType, view *View) int64 {
	priority := -(view.Height.Int64()*100 + view.Round.Int64()*10 + messagePriorityTable[msgCode])
	return priority
}
