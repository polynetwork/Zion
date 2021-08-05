package core

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/prque"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) handleRequest(req *hotstuff.Request) error {
	logger := c.newLogger()

	if err := c.requests.checkRequest(c.currentView(), req); err != nil {
		if err == errFutureMessage {
			c.requests.StoreRequest(req)
			return nil
		} else {
			logger.Warn("receive request", "err", err)
			return err
		}
	} else {
		c.requests.StoreRequest(req)
	}

	if c.currentState() == StateAcceptRequest &&
		c.current.highQC != nil &&
		c.current.highQC.View.Cmp(c.currentView()) == 0 {
		c.sendPrepare()
	}

	logger.Trace("handleRequest", "height", req.Proposal.Number(), "proposal", req.Proposal.Hash())
	return nil
}

type requestSet struct {
	mtx *sync.RWMutex

	pendingRequest *prque.Prque
}

func newRequestSet() *requestSet {
	return &requestSet{
		mtx:            new(sync.RWMutex),
		pendingRequest: prque.New(nil),
	}
}

func (s *requestSet) checkRequest(view *hotstuff.View, req *hotstuff.Request) error {
	if req == nil || req.Proposal == nil {
		return errInvalidMessage
	}

	// todo(fuk): how to process future block, store or throw?
	if c := view.Height.Cmp(req.Proposal.Number()); c < 0 {
		return errFutureMessage
	} else if c > 0 {
		return errOldMessage
	} else {
		return nil
	}
}

func (s *requestSet) StoreRequest(req *hotstuff.Request) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	priority := -req.Proposal.Number().Int64()
	s.pendingRequest.Push(req, priority)
}

func (s *requestSet) GetRequest(view *hotstuff.View) *hotstuff.Request {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	maxRetry := 20
retry:
	for !s.pendingRequest.Empty() {
		m, prior := s.pendingRequest.Pop()
		req, ok := m.(*hotstuff.Request)
		if !ok {
			continue
		}

		// push back if it's future message
		if err := s.checkRequest(view, req); err != nil {
			if err == errFutureMessage {
				s.pendingRequest.Push(m, prior)
				// todo: 是否为continue
				break
			}
			continue
		}
		return req
	}
	if maxRetry -= 1; maxRetry > 0 {
		time.Sleep(500 * time.Millisecond)
		goto retry
	}
	return nil
}

func (s *requestSet) Size() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.pendingRequest.Size()
}
