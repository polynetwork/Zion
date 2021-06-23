package core

import "github.com/ethereum/go-ethereum/consensus/hotstuff"

func (c *core) handlePreCommitVote(data *message, src hotstuff.Validator) {
	logger := c.logger.New("state", c.currentState())

	var (
		msg    *hotstuff.Subject
		msgTyp = MsgTypePreCommitVote
	)
	if err := c.decodeAndCheckVote(data, msgTyp, msg); err != nil {
		logger.Error("Failed to check vote", "type", msgTyp, "err", err)
		return
	}

	if err := c.current.AddPreCommitVote(data); err != nil {
		logger.Error("Failed to add vote", "type", msgTyp, "err", err)
		return
	}

	if size := c.current.PreCommitVoteSize(); size >= c.Q() && c.currentState() < StateLocked {
		c.current.SetLockedQC(c.current.PrepareQC())
		c.current.SetState(StateLocked)
		c.sendCommit()
	}
}

func (c *core) sendCommit() {
	logger := c.logger.New("state", c.currentState())

	msgTyp := MsgTypeCommit
	payload, err := Encode(c.current.LockedQC())
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&message{Code: msgTyp, Msg: payload})
}

func (c *core) handleCommit(data *message, src hotstuff.Validator) {
	logger := c.logger.New("state", c.currentState())

	var (
		msg    *hotstuff.QuorumCert
		msgTyp = MsgTypeCommit
	)
	if err := c.decodeAndCheckMessage(data, msgTyp, msg); err != nil {
		logger.Error("Failed to check msg", "type", msgTyp, "err", err)
		return
	}

	c.current.SetLockedQC(msg)
	c.current.SetState(StateLocked)
}

func (c *core) sendCommitVote() {
	logger := c.logger.New("state", c.currentState())

	msgTyp := MsgTypeCommitVote
	payload, err := Encode(c.current.LockedQC())
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&message{Code: msgTyp, Msg: payload})
}
