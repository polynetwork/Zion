package core

import (
	"reflect"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) sendPrepareVote() {
	logger := c.logger.New("state", c.state)

	sub := c.current.Subject()
	payload, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", MsgTypePrepareVote.String(), sub)
		return
	}

	c.broadcast(&message{
		Code: MsgTypePrepareVote,
		Msg:  payload,
	})
}

func (c *core) handlePrepareVote(msg *message, src hotstuff.Validator) error {
	var vote *hotstuff.Subject
	if err := msg.Decode(&vote); err != nil {
		return errFailedDecodePrepareVote
	}

	if err := c.checkMessage(MsgTypePrepareVote, vote.View); err != nil {
		return err
	}

	if err := c.verifyPrepareVote(vote, src); err != nil {
		return err
	}

	c.acceptPrepareVote(msg, src)

	if c.current.PrepareVoteSize() == c.valSet.Q() {
		c.setState(StatePrepared)
		c.sendPreCommit()
	}
	return nil
}

func (c *core) verifyPrepareVote(vote *hotstuff.Subject, src hotstuff.Validator) error {
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

func (c *core) acceptPrepareVote(msg *message, src hotstuff.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)
	if err := c.current.AddPrepareVote(msg); err != nil {
		logger.Error("Failed to add PREPARE vote message to round state", "msg", msg, "err", err)
		return err
	}
	return nil
}
