package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) handleCommitVote(data *message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		vote   *hotstuff.Vote
		msgTyp = MsgTypeCommitVote
	)
	if err := data.Decode(&vote); err != nil {
		logger.Error("Failed to check vote", "vote", msgTyp, "err", err)
		return errFailedDecodeCommitVote
	}
	if err := c.checkView(msgTyp, vote.View); err != nil {
		return err
	}
	if err := c.checkVote(vote); err != nil {
		return err
	}
	if vote.Digest != c.current.LockedQC().Hash {
		return errInvalidDigest
	}
	if err := c.checkMsgToProposer(); err != nil {
		return err
	}

	if err := c.current.AddCommitVote(data); err != nil {
		logger.Error("Failed to add vote", "type", msgTyp, "err", err)
		return errAddPreCommitVote
	}

	if size := c.current.CommitVoteSize(); size >= c.Q() && c.currentState() < StateCommitted {
		c.acceptDecide()
		// c.sendDecide()
		if err := c.backend.Commit(c.current.Proposal()); err != nil {
			logger.Error("Failed to commit proposal", "err", err)
			return err
		}
		c.startNewRound(common.Big0)
	}

	logger.Trace("handleCommitVote", "src", src.Address(), "vote view", vote.View, "vote", vote.Digest)
	return nil
}

func (c *core) sendDecide() {
	logger := c.newLogger()

	msgTyp := MsgTypeDecide
	qc := c.current.CommittedQC()
	payload, err := Encode(qc)
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&message{Code: msgTyp, Msg: payload})
	logger.Trace("sendDecide", "msg view", qc.View, "proposal", qc.Hash)
}

func (c *core) handleDecide(data *message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		msg    *hotstuff.QuorumCert
		msgTyp = MsgTypeDecide
	)
	if err := data.Decode(&msg); err != nil {
		logger.Error("Failed to check msg", "type", msgTyp, "err", err)
		return errFailedDecodeDecide
	}
	if err := c.checkView(MsgTypeDecide, msg.View); err != nil {
		return err
	}
	if err := c.checkMsgFromProposer(src); err != nil {
		return err
	}
	if !c.IsProposer() {
		c.acceptDecide()
	}
	logger.Trace("handleDecide", "src", src.Address(), "msg view", msg.View, "proposal", msg.Hash)
	c.startNewRound(common.Big0)
	return nil
}

func (c *core) acceptDecide() {
	c.current.SetState(StateCommitted)
	c.current.SetCommittedQC(c.current.LockedQC())
}

func (c *core) handleFinalCommitted() error {
	logger := c.newLogger()
	logger.Trace("handleFinalCommitted")
	c.startNewRound(common.Big0)
	return nil
}
