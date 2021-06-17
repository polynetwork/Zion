package core

import (
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"reflect"
)

func (c *core) sendCommit() {
	if !c.IsProposer() {
		return
	}
	logger := c.logger.New("state", c.state)

	sub := c.current.Subject()
	payload, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", MsgTypeCommit.String(), sub)
		return
	}
	c.broadcast(&message{
		Code:      MsgTypeCommit,
		Msg:       payload,
	})
}

func (c *core) handleCommit(data *message, src hotstuff.Validator) error {
	var msg *hotstuff.Subject
	if err := data.Decode(&msg); err != nil {
		return errFailedDecodeCommit
	}

	if err := c.checkMessage(MsgTypeCommit, msg.View); err != nil {
		return err
	}

	if err := c.verifyCommit(msg, src); err != nil {
		return err
	}

	c.sendCommitVote()
	return nil
}

func (c *core) verifyCommit(msg *hotstuff.Subject, src hotstuff.Validator) error {
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
