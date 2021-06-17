package core

import (
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"reflect"
)


func (c *core) sendCommitVote() {
	logger := c.logger.New("state", c.state)

	sub := c.current.Subject()
	payload, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", MsgTypeCommitVote.String(), sub)
		return
	}
	c.broadcast(&message{
		Code:      MsgTypeCommitVote,
		Msg:       payload,
	})
}

func (c *core) handleCommitVote(msg *message, src hotstuff.Validator) error {
	var vote *hotstuff.Subject
	if err := msg.Decode(&vote); err != nil {
		return errFailedDecodeCommitVote
	}

	if err := c.checkMessage(MsgTypeCommitVote, vote.View); err != nil {
		return err
	}

	if err := c.verifyCommit(vote, src); err != nil {
		return err
	}

	c.acceptPrepareVote(msg, src)

	isHashLocked := c.current.IsHashLocked() && vote.Digest == c.current.GetLockedHash()
	isQuorum := c.current.PreCommitVoteSize() >= c.valSet.Q()
	if isHashLocked && isQuorum && c.state < StateCommitted {
		c.setState(StateCommitted)
		c.sendDecide()
	}
	return nil
}

func (c *core) verifyCommitVote(vote *hotstuff.Subject, src hotstuff.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)
	if c.IsProposer() {
		return errNotToProposer
	}

	sub := c.current.Subject()
	if !reflect.DeepEqual(vote, sub) {
		logger.Warn("Inconsistent votes between PREPARE and vote", "expected", sub, "got", vote)
		return errInconsistentSubject
	}
	return nil
}

func (c *core) acceptCommitVote(msg *message, src hotstuff.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)
	if err := c.current.AddCommitVote(msg); err != nil {
		logger.Error("Failed to add PREPARE vote message to round state", "msg", msg, "err", err)
		return err
	}
	return nil
}
