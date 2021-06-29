package core

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/prque"
)

func (c *core) handleRequest(req *hotstuff.Request) error {
	logger := c.logger.New("handle request, state", c.currentState(), "view", c.currentView())

	if err := c.requests.checkRequest(c.currentView(), req); err != nil {
		if err == errFutureMessage {
			c.requests.StoreRequest(req)
			return nil
		} else {
			logger.Warn("receive request", "err", err)
			return err
		}
	}

	if c.currentState() == StateAcceptRequest && c.current.PendingRequest() == nil {
		c.current.SetPendingRequest(req)
	}

	logger.Trace("store request", "validator", c.address.Hex())
	return nil
}

type requestSet struct {
	mtx *sync.RWMutex

	pendingRequest *prque.Prque
}

func newRequestSet() *requestSet {
	return &requestSet{
		mtx:            new(sync.RWMutex),
		pendingRequest: prque.New(),
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

	s.pendingRequest.Push(req, float32(-req.Proposal.Number().Int64()))
}

func (s *requestSet) GetRequest(view *hotstuff.View) *hotstuff.Request {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	maxRetry := 3
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
				break
			}
			continue
		}
		return req
	}
	if maxRetry -= 1; maxRetry > 0 {
		time.Sleep(1 * time.Second)
		goto retry
	}
	return nil
}

func (s *requestSet) Size() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.pendingRequest.Size()
}
