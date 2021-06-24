package core

import (
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) handlePreCommitVote(data *message, src hotstuff.Validator) error {
	logger := c.logger.New("state", c.currentState())

	var (
		vote   *hotstuff.Vote
		msgTyp = MsgTypePreCommitVote
	)
	if err := data.Decode(&vote); err != nil {
		logger.Error("Failed to check vote", "type", msgTyp, "err", err)
		return errFailedDecodePreCommitVote
	}
	if err := c.checkView(msgTyp, vote.View); err != nil {
		return err
	}
	if err := c.checkVote(vote); err != nil {
		return err
	}
	if vote.Digest != c.current.PrepareQC().Hash {
		return errInvalidDigest
	}
	if err := c.checkMsgToProposer(); err != nil {
		return err
	}

	if err := c.current.AddPreCommitVote(data); err != nil {
		logger.Error("Failed to add vote", "type", msgTyp, "err", err)
		return errAddPreCommitVote
	}

	if size := c.current.PreCommitVoteSize(); size >= c.Q() && c.currentState() < StateLocked {
		c.current.SetLockedQC(c.current.PrepareQC())
		c.current.SetState(StateLocked)
		c.sendCommit()
	}
	return nil
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

func (c *core) handleCommit(data *message, src hotstuff.Validator) error {
	logger := c.logger.New("state", c.currentState())

	var (
		msg    *hotstuff.QuorumCert
		msgTyp = MsgTypeCommit
	)
	if err := data.Decode(&msg); err != nil {
		logger.Error("Failed to check msg", "type", msgTyp, "err", err)
		return errFailedDecodeCommit
	}
	if err := c.checkView(MsgTypeCommit, msg.View); err != nil {
		return err
	}
	if err := c.checkMsgFromProposer(src); err != nil {
		return err
	}
	if err := c.checkPrepareQC(msg); err != nil {
		return err
	}
	if err := c.backend.VerifyQuorumCert(msg); err != nil {
		return errVerifyQC
	}
	if c.current.State() >= StateLocked {
		return errState
	}

	c.current.SetLockedQC(msg)
	c.current.SetState(StateLocked)
	c.sendCommitVote()
	return nil
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
