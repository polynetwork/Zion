package core

import "github.com/ethereum/go-ethereum/consensus/hotstuff"

// todo: 因为有viewChange，所以非leader节点的request也应该保存起来
func (c *core) handleRequest(req *hotstuff.Request) error {
	logger := c.logger.New("state", c.state, "height", c.current.Height())
	if err := c.checkRequest(req); err != nil {
		if err == errInvalidMessage {
			logger.Warn("invalid request")
			return err
		}
		logger.Warn("unexpected request", "number", req.Proposal.Number(), "hash", req.Proposal.Hash(), "err", err)
		return err
	}

	logger.Trace("handle request", "number", req.Proposal.Number(), "hash", req.Proposal.Hash())
	if c.state == StateNewRound {
		c.sendNewView(req)
	}
	return nil
}

func (c *core) checkRequest(req *hotstuff.Request) error {
	if req == nil || req.Proposal == nil {
		return errInvalidMessage
	}

	if c := c.current.Height().Cmp(req.Proposal.Number()); c < 0 {
		return errFutureMessage
	} else if c > 0 {
		return errOldMessage
	} else {
		return nil
	}
}

func (c *core) storeRequestMsg(req *hotstuff.Request) {
	logger := c.logger.New("state", c.state)
	logger.Trace("Store future request", "number", req.Proposal.Number(), "hash", req.Proposal.Hash())

	c.pendingRQMu.Lock()
	defer c.pendingRQMu.Unlock()

	c.pendingRQ.Push(req, float32(-req.Proposal.Number().Int64()))
}

func (c *core) processPendingRequest() {
	c.pendingRQMu.Lock()
	defer c.pendingRQMu.Unlock()

	for !(c.pendingRQ.Empty()) {
		m, prio := c.pendingRQ.Pop()
		req, ok := m.(*hotstuff.Request)
		if !ok {
			c.logger.Warn("Malformed request, skip", "msg", m)
			continue
		}

		// push back if it's future message
		if err := c.checkRequest(req); err != nil {
			if err == errFutureMessage {
				c.logger.Trace("Stop processing request", "number", req.Proposal.Number(), "hash", req.Proposal.Hash())
				c.pendingRQ.Push(m, prio)
				break
			}

			c.logger.Trace("Skip the pending request", "number", req.Proposal.Number(), "hash", req.Proposal.Hash(), "err", err)
			continue
		}

		// post valid request
		c.logger.Trace("Post pending request", "number", req.Proposal.Number(), "hash", req.Proposal.Hash())
		c.sendEvent(hotstuff.RequestEvent{
			Proposal: req.Proposal,
		})
	}
}
