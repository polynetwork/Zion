package core
//
//import (
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/consensus/hotstuff"
//	"github.com/ethereum/go-ethereum/common/prque"
//	"sync"
//)
//
//func (c *core) storeBacklog(msg *message, src hotstuff.Validator) {
//	logger := c.newLogger()
//
//	if src.Address() == c.Address() {
//		logger.Warn("Backlog from self")
//		return
//	}
//
//	logger.Trace("Store future message")
//
//	c.backlogsMu.Lock()
//	defer c.backlogsMu.Unlock()
//
//	logger.Debug("Retrieving backlog queue", "for", src.Address(), "backlogs_size", len(c.backlogs))
//	backlog := c.backlogs[src.Address()]
//	if backlog == nil {
//		backlog = prque.New(nil)
//	}
//	switch msg.Code {
//	case msgPreprepare:
//		var p *istanbul.Preprepare
//		err := msg.Decode(&p)
//		if err == nil {
//			backlog.Push(msg, toPriority(msg.Code, p.View))
//		}
//		// for msgRoundChange, msgPrepare and msgCommit cases
//	default:
//		var p *istanbul.Subject
//		err := msg.Decode(&p)
//		if err == nil {
//			backlog.Push(msg, toPriority(msg.Code, p.View))
//		}
//	}
//	c.backlogs[src.Address()] = backlog
//}
//
//func (c *core) processBacklog() {
//	c.backlogsMu.Lock()
//	defer c.backlogsMu.Unlock()
//
//	for srcAddress, backlog := range c.backlogs {
//		if backlog == nil {
//			continue
//		}
//		_, src := c.valSet.GetByAddress(srcAddress)
//		if src == nil {
//			// validator is not available
//			delete(c.backlogs, srcAddress)
//			continue
//		}
//		logger := c.logger.New("from", src, "state", c.state)
//		isFuture := false
//
//		// We stop processing if
//		//   1. backlog is empty
//		//   2. The first message in queue is a future message
//		for !(backlog.Empty() || isFuture) {
//			m, prio := backlog.Pop()
//			msg := m.(*message)
//			var view *istanbul.View
//			switch msg.Code {
//			case msgPreprepare:
//				var m *istanbul.Preprepare
//				err := msg.Decode(&m)
//				if err == nil {
//					view = m.View
//				}
//				// for msgRoundChange, msgPrepare and msgCommit cases
//			default:
//				var sub *istanbul.Subject
//				err := msg.Decode(&sub)
//				if err == nil {
//					view = sub.View
//				}
//			}
//			if view == nil {
//				logger.Debug("Nil view", "msg", msg)
//				continue
//			}
//			// Push back if it's a future message
//			err := c.checkMessage(msg.Code, view)
//			if err != nil {
//				if err == errFutureMessage {
//					logger.Trace("Stop processing backlog", "msg", msg)
//					backlog.Push(msg, prio)
//					isFuture = true
//					break
//				}
//				logger.Trace("Skip the backlog event", "msg", msg, "err", err)
//				continue
//			}
//			logger.Trace("Post backlog event", "msg", msg)
//
//			go c.sendEvent(backlogEvent{
//				src: src,
//				msg: msg,
//			})
//		}
//	}
//}
//
//func toPriority(msgCode uint64, view *hotstuff.View) int64 {
//	if msgCode == msgRoundChange {
//		// For msgRoundChange, set the message priority based on its sequence
//		return -int64(view.Height.Uint64() * 1000)
//	}
//	// FIXME: round will be reset as 0 while new sequence
//	// 10 * Round limits the range of message code is from 0 to 9
//	// 1000 * Sequence limits the range of round is from 0 to 99
//	return -float32(view.Sequence.Uint64()*1000 + view.Round.Uint64()*10 + uint64(msgPriority[msgCode]))
//}
//
//type backlog struct {
//	mu *sync.RWMutex
//	backlogs   map[common.Address]*prque.Prque
//}
//
//func newBlocklog() *backlog {
//
//}