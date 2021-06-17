package core

import (
	"reflect"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) sendPreCommitVote() {
	logger := c.logger.New("state", c.state)

	sub := c.current.Subject()
	payload, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", MsgTypePreCommitVote.String(), sub)
		return
	}

	c.broadcast(&message{
		Code: MsgTypePreCommitVote,
		Msg:  payload,
	})
}

func (c *core) handlePreCommitVote(msg *message, src hotstuff.Validator) error {
	var vote *hotstuff.Subject
	if err := msg.Decode(&vote); err != nil {
		return errFailedDecodePreCommitVote
	}

	if err := c.checkMessage(MsgTypePreCommitVote, vote.View); err != nil {
		return err
	}

	if err := c.verifyPrepareVote(vote, src); err != nil {
		return err
	}

	c.acceptPrepareVote(msg, src)

	isHashLocked := c.current.IsHashLocked() && vote.Digest == c.current.GetLockedHash()
	isQuorum := c.current.PreCommitVoteSize() > c.valSet.Q()
	if (isHashLocked || isQuorum) && c.state.Cmp(StatePreCommitted) < 0 {
		c.current.LockHash()
		c.setState(StatePreCommitted)
		c.sendCommit()
	}

	return nil
}

func (c *core) verifyPreCommitVote(vote *hotstuff.Subject, src hotstuff.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	if !c.IsProposer() {
		return errNotToProposer
	}

	sub := c.current.Subject()
	if !reflect.DeepEqual(sub, vote) {
		logger.Warn("Inconsistent votes between PREPARE and vote", "expected", sub, "got", vote)
		return errInconsistentSubject
	}
	return nil
}

func (c *core) acceptPreCommitVote(msg *message, src hotstuff.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)
	// Add the PREPARE message to current round state
	if err := c.current.AddPreCommitVote(msg); err != nil {
		logger.Error("Failed to add PREPARE vote message to round state", "msg", msg, "err", err)
		return err
	}
	return nil
}
