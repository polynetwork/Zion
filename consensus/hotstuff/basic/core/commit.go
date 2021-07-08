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
		logger.Error("Failed to check vote", "type", msgTyp, "err", err)
		return errFailedDecodePreCommitVote
	}
	if err := c.checkView(msgTyp, vote.View); err != nil {
		return err
	}
	if err := c.checkVote(vote); err != nil {
		return err
	}
	if vote.Digest != c.current.Proposal().Hash() {
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
	logger.Trace("handlePreCommitVote", "src", src.Address(), "vote view", vote.View, "vote", vote.Digest)
	return nil
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
	if err := c.signer.VerifyQC(msg, c.valSet); err != nil {
		return errVerifyQC
	}
	if !c.IsProposer() && c.current.State() < StateLocked {
		c.current.SetLockedQC(msg)
		c.current.SetState(StateLocked)
	}
	c.sendCommitVote()
	logger.Trace("handleCommit", "address", src.Address(), "msg view", msg.View, "proposal", msg.Hash)
	return nil
}

func (c *core) sendCommitVote() {
	logger := c.newLogger()

	msgTyp := MsgTypeCommitVote
	vote := c.current.Vote()
	if vote == nil {
		logger.Error("proposal is nil")
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
