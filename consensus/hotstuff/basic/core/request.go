package core

import (
	"sync"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/prque"
)

func (c *core) handleRequest(req *hotstuff.Request) {
	logger := c.logger.New("state", c.currentState(), "height", c.current.Height())
	if err := c.requests.checkRequest(c.currentView(), req); err != nil {
		if err == errInvalidMessage {
			logger.Warn("invalid request")
			return
		}
		logger.Warn("unexpected request", "number", req.Proposal.Number(), "hash", req.Proposal.Hash(), "err", err)
		return
	}

	logger.Trace("handle request", "number", req.Proposal.Number(), "hash", req.Proposal.Hash())
	c.requests.AddRequest(req)
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

	if c := view.Height.Cmp(req.Proposal.Number()); c < 0 {
		return errFutureMessage
	} else if c > 0 {
		return errOldMessage
	} else {
		return nil
	}
}

func (s *requestSet) AddRequest(req *hotstuff.Request) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.pendingRequest.Push(req, float32(-req.Proposal.Number().Int64()))
}

func (s *requestSet) GetRequest(view *hotstuff.View) *hotstuff.Request {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for !s.pendingRequest.Empty() {
		m, prior := s.pendingRequest.Pop()
		req, ok := m.(*hotstuff.Request)
		if !ok {
			// todo: add log
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
	return nil
}
