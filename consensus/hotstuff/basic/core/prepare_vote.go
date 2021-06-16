package core

import (
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) sendPrepareVote() {
	logger := c.logger.New("state", c.state)

	// allow leader to send prepareVote to itself
	qc := c.current.prepareQC
	curView := c.currentView()
	if curView.Height.Cmp(qc.Proposal.Number()) == 0 {
		vote := &MsgPrepareVote{
			View:     curView,
			BlockHash: qc.Proposal.Hash(),
		}
		if sig, err := c.backend.Sign(vote.BlockHash[:]); err != nil {
			logger.Error("Failed to sign prepare vote", "view", curView, "err", err)
			return
		} else {
			vote.Signature = sig
		}

		payload, err := Encode(vote)
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}
		c.broadcast(&message{
			Code: MsgTypePrepareVote,
			Msg:  payload,
		}, curView.Round)
	}
}

func (c *core) handlePrepareVote(msg *message, src hotstuff.Validator) error {
	var vote *MsgPrepareVote
	if err := msg.Decode(&vote); err != nil {
		return errFailedDecodePrepareVote
	}

	// todo check message
	if vote.View.Cmp(c.currentView()) != 0 {
		return errInvalidMessage
	}

	if c.current.PrepareQC().Proposal.Hash() != vote.BlockHash {
		return errInvalidMessage
	}

	if err := c.backend.CheckSignature(vote.BlockHash[:], src.Address(), vote.Signature); err != nil {
		return err
	}

	if err := c.acceptPrepareVote(msg, src); err != nil {
		return err
	}

	if c.current.PrepareVoteSize() == c.valSet.Q() {
		c.setState(StatePrepared)
		// todo: copy prepare qc and set preCommitType
		c.current.SetPreCommitQC(c.current.prepareQC)
		c.sendPreCommit()
	}
	return nil
}

func (c *core) acceptPrepareVote(msg *message, src hotstuff.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Add the PREPARE message to current round state
	if err := c.current.AddPrepareVote(msg); err != nil {
		logger.Error("Failed to add PREPARE vote message to round state", "msg", msg, "err", err)
		return err
	}

	return nil
}
