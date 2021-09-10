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
	"github.com/ethereum/go-ethereum/common/prque"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) storeBacklog(msg *hotstuff.Message, src hotstuff.Validator) {
	logger := c.newLogger()

	if src.Address() == c.Address() {
		logger.Trace("Backlog from self")
		return
	}
	if _, v := c.valset.GetByAddress(src.Address()); v == nil {
		logger.Trace("Backlog from unknown validator", "address", src.Address())
		return
	}

	logger.Debug("Retrieving backlog queue", "for", src.Address(), "backlogs_size", len(c.backlogs.queue))

	c.backlogs.Push(msg)
}

func (c *core) processBacklog() {
	logger := c.newLogger()

	for addr, queue := range c.backlogs.queue {
		if queue == nil {
			continue
		}
		_, src := c.valset.GetByAddress(addr)
		if src == nil {
			logger.Trace("Skip the backlog", "unknown validator", addr)
			continue
		}

		isFuture := false
		for !(queue.Empty() || isFuture) {
			data, priority := queue.Pop()
			msg, ok := data.(*hotstuff.Message)
			if !ok {
				logger.Trace("Skip the backlog, invalid Message")
				continue
			}
			if err := c.checkView(msg.Code, msg.View); err != nil {
				if err == errFutureMessage {
					queue.Push(data, priority)
					isFuture = true
					break
				}
				logger.Trace("Skip the backlog", "msg view", msg.View, "err", err)
				continue
			}

			logger.Trace("Replay the backlog", "msg", msg)
			go c.sendEvent(backlogEvent{src: src, msg: msg})
		}
	}
}

type backlog struct {
	queue map[common.Address]*prque.Prque
}

func newBackLog() *backlog {
	return &backlog{
		queue: make(map[common.Address]*prque.Prque),
	}
}

func (b *backlog) Push(msg *hotstuff.Message) {
	if msg == nil || msg.Address == common.EmptyAddress {
		return
	}

	addr := msg.Address
	if _, ok := b.queue[addr]; !ok {
		b.queue[addr] = prque.New(nil)
	}
	priority := b.toPriority(msg.Code, msg.View)
	b.queue[addr].Push(msg, priority)
}

func (b *backlog) Pop(addr common.Address) (data *hotstuff.Message, priority int64) {
	if _, ok := b.queue[addr]; !ok {
		return
	} else {
		item, p := b.queue[addr].Pop()
		data = item.(*hotstuff.Message)
		priority = p
		return
	}
}

var messagePriorityTable = map[hotstuff.MsgType]int64{
	MsgTypeProposal: 1,
	MsgTypeVote:     2,
	MsgTypeTimeout:  3,
	MsgTypeTC:       4,
}

func (b *backlog) toPriority(msgCode hotstuff.MsgType, view *hotstuff.View) int64 {
	priority := -(view.Height.Int64()*100 + view.Round.Int64()*10 + int64(messagePriorityTable[msgCode]))
	return priority
}
