package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) handlePreCommitVote(data *message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		vote   *hotstuff.Vote
		msgTyp = MsgTypePreCommitVote
	)
	if err := data.Decode(&vote); err != nil {
		logger.Trace("Failed to decode", "type", msgTyp, "err", err)
		return errFailedDecodePreCommitVote
	}
	if err := c.checkView(msgTyp, vote.View); err != nil {
		logger.Trace("Failed to check view", "type", msgTyp, "err", err)
		return err
	}
	if err := c.checkVote(vote); err != nil {
		logger.Trace("Failed to check vote", "type", msgTyp, "err", err)
		return err
	}
	if vote.Digest != c.current.Proposal().Hash() {
		logger.Trace("Failed to check hash", "type", msgTyp, "expect vote digest", c.current.Proposal().Hash().Hex(), vote.Digest.Hex())
		return errInvalidDigest
	}
	if err := c.checkMsgToProposer(); err != nil {
		logger.Trace("Failed to check proposal", "type", msgTyp, "err", err)
		return err
	}

	if err := c.current.AddPreCommitVote(data); err != nil {
		logger.Trace("Failed to add vote", "type", msgTyp, "err", err)
		return errAddPreCommitVote
	}

	logger.Trace("handlePreCommitVote", "src", src.Address(), "vote view", vote.View, "vote", vote.Digest)

	if c.current.PreCommitVoteSize() >= c.Q() && c.currentState() < StateLocked {
		c.acceptPreCommitVote(src, vote)
		c.sendCommit()
	}
	return nil
}

func (c *core) acceptPreCommitVote(src hotstuff.Validator, msg *hotstuff.Vote) {
	logger := c.newLogger()
	c.current.SetLockedQC(c.current.PrepareQC())
	c.current.SetState(StateLocked)
	logger.Trace("acceptPreCommitVote", "vote digest", msg.Digest.Hex())
}

func (c *core) sendCommit() {
	logger := c.newLogger()

	msgTyp := MsgTypeCommit
	sub := c.current.LockedQC()
	payload, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&message{Code: msgTyp, Msg: payload})
	logger.Trace("sendCommit", "msg view", sub.View, "proposal", sub.Hash)
}

func (c *core) handleCommit(data *message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		msg    *hotstuff.QuorumCert
		msgTyp = MsgTypeCommit
	)
	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to decode", "type", msgTyp, "err", err)
		return errFailedDecodeCommit
	}
	if err := c.checkView(MsgTypeCommit, msg.View); err != nil {
		logger.Trace("Failed to check view", "type", msgTyp, "err", err)
		return err
	}
	if err := c.checkMsgFromProposer(src); err != nil {
		logger.Trace("Failed to check proposer", "type", msgTyp, "err", err)
		return err
	}
	if err := c.checkPrepareQC(msg); err != nil {
		logger.Trace("Failed to check prepareQC", "type", msgTyp, "err", err)
		return err
	}
	if err := c.signer.VerifyQC(msg, c.valSet); err != nil {
		logger.Trace("Failed to check verify qc", "type", msgTyp, "err", err)
		return errVerifyQC
	}

	logger.Trace("handleCommit", "address", src.Address(), "msg view", msg.View, "proposal", msg.Hash)

	if !c.IsProposer() && c.current.State() < StateLocked {
		c.acceptCommit(src, msg)
	}
	c.sendCommitVote()
	return nil
}

func (c *core) acceptCommit(src hotstuff.Validator, msg *hotstuff.QuorumCert) {
	logger := c.newLogger()
	c.current.SetLockedQC(msg)
	c.current.SetState(StateLocked)
	logger.Trace("acceptCommit", "msg view", msg.View, "hash", msg.Hash.Hex())
}

func (c *core) sendCommitVote() {
	logger := c.newLogger()

	msgTyp := MsgTypeCommitVote
	vote := c.current.Vote()
	if vote == nil {
		logger.Error("Failed to send vote", "msg", msgTyp, "err", "current vote is nil")
		return
	}
	payload, err := Encode(vote)
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&message{Code: msgTyp, Msg: payload})
	logger.Trace("sendCommitVote", "vote view", vote.View, "vote", vote.Digest)

	if !c.IsProposer() {
		c.startNewRound(common.Big0)
	}
}
