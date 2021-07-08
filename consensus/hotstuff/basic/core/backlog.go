package core

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/prque"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) storeBacklog(msg *message, src hotstuff.Validator) {
	logger := c.newLogger()

	if src.Address() == c.Address() {
		logger.Trace("Backlog from self")
		return
	}
	if _, v := c.valSet.GetByAddress(src.Address()); v == nil {
		logger.Trace("Backlog from unknown validator", "address", src.Address())
		return
	}

	logger.Trace("Store backlog")
	logger.Debug("Retrieving backlog queue", "for", src.Address(), "backlogs_size", len(c.backlogs.queue))

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
			msg, ok := data.(*message)
			if !ok {
				logger.Trace("Skip the backlog, invalid message")
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
	mu    *sync.Mutex
	queue map[common.Address]*prque.Prque
}

func newBackLog() *backlog {
	return &backlog{
		mu:    new(sync.Mutex),
		queue: make(map[common.Address]*prque.Prque),
	}
}

func (b *backlog) Push(msg *message) {
	if msg == nil || msg.Address == EmptyAddress {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()

	addr := msg.Address
	if _, ok := b.queue[addr]; !ok {
		b.queue[addr] = prque.New(nil)
	}
	priority := b.toPriority(msg.Code, msg.View)
	b.queue[addr].Push(msg, priority)
}

func (b *backlog) Pop(addr common.Address) (data *message, priority int64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.queue[addr]; !ok {
		return
	} else {
		item, p := b.queue[addr].Pop()
		data = item.(*message)
		priority = p
		return
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

func (b *backlog) toPriority(msgCode MsgType, view *hotstuff.View) int64 {
	priority := -(view.Height.Int64()*100 + view.Round.Int64()*10 + int64(messagePriorityTable[msgCode]))
	return priority
}
