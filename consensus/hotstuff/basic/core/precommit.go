package core

import (
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) handlePrepareVote(data *message, src hotstuff.Validator) {
	logger := c.logger.New("state", c.currentState())

	var (
		vote   *hotstuff.Subject
		msgTyp = MsgTypePrepareVote
	)
	if err := c.decodeAndCheckVote(data, msgTyp, vote); err != nil {
		logger.Error("Failed to check vote", "type", msgTyp, "err", err)
		return
	}

	if err := c.current.AddPreCommitVote(data); err != nil {
		logger.Error("Failed to add vote", "type", msgTyp, "err", err)
		return
	}

	if size := c.current.PreCommitVoteSize(); size >= c.Q() && c.current.state < StatePrepared {
		seal := c.getSeals(size)
		newProposal, err := c.backend.PreCommit(c.current.Proposal(), seal)
		if err != nil {
			logger.Error("Failed to pre-commit", "err", err)
			return
		}
		c.current.SetPrepareQC(&QuorumCert{
			View:     c.currentView(),
			Proposal: newProposal,
		})
		c.current.SetState(StatePrepared)
		c.sendPreCommit()
	}
}

func (c *core) getSeals(n int) [][]byte {
	seals := make([][]byte, n)
	for i, data := range c.current.PrepareVotes() {
		if i < n {
			seals[i] = data.Signature
		}
	}
	return seals
}

func (c *core) sendPreCommit() {
	logger := c.logger.New("state", c.current.State())

	msgTyp := MsgTypePreCommit
	payload, err := Encode(c.current.PrepareQC())
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&message{Code: msgTyp, Msg: payload})
}

func (c *core) handlePreCommit(data *message, src hotstuff.Validator) {
	logger := c.logger.New("state", c.currentState())

	var (
		msg    *QuorumCert
		msgTyp = MsgTypePreCommit
	)
	if err := c.decodeAndCheckMessage(data, msgTyp, msg); err != nil {
		logger.Error("Failed to check msg", "type", msgTyp, "err", err)
		return
	}

	if _, err := c.backend.Verify(msg.Proposal); err != nil {
		logger.Error("Failed to verify proposal", "err", err)
		return
	}

	c.current.SetPrepareQC(msg)
	c.current.SetState(StatePrepared)
}

func (c *core) sendPreCommitVote() {
	logger := c.logger.New("state", c.current.State())

	msgTyp := MsgTypePreCommitVote
	sub := c.current.Subject()
	payload, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&message{Code: msgTyp, Msg: payload})
}
