package core

import (
	"reflect"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) sendPreCommit() {
	if !c.IsProposer() {
		return
	}

	logger := c.logger.New("state", c.state)
	sub := c.current.Subject()
	payload, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", MsgTypePreCommit.String(), sub)
		return
	}
	c.broadcast(&message{
		Code: MsgTypePreCommit,
		Msg:  payload,
	})
}

func (c *core) handlePreCommit(data *message, src hotstuff.Validator) error {
	var msg *hotstuff.Subject
	if err := data.Decode(&msg); err != nil {
		return errFailedDecodePreCommit
	}

	if err := c.checkMessage(MsgTypePreCommit, msg.View); err != nil {
		return err
	}

	if err := c.verifyPreCommit(msg, src); err != nil {
		return err
	}

	// validator lock proposal hash and set state as locked
	if !c.IsProposer() {
		c.current.LockHash()
		c.setState(StateLocked)
	}

	// allow leader or validator send vote
	c.sendPreCommitVote()

	return nil
}

func (c *core) verifyPreCommit(msg *hotstuff.Subject, src hotstuff.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	if !c.valSet.IsProposer(src.Address()) {
		return errNotFromProposer
	}

	sub := c.current.Subject()
	if !reflect.DeepEqual(sub, msg) {
		logger.Warn("Inconsistent votes between PREPARE and vote", "expected", sub, "got", msg)
		return errInconsistentSubject
	}

	return nil
}
