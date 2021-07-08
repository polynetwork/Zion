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
		c.current.SetState(StateCommitted)
		c.current.SetCommittedQC(c.current.LockedQC())
		if err := c.backend.Commit(c.current.Proposal()); err != nil {
			logger.Error("Failed to commit proposal", "err", err)
			return err
		}
		c.startNewRound(common.Big0)
	}

	logger.Trace("handleCommitVote", "src", src.Address(), "vote view", vote.View, "vote", vote.Digest)
	return nil
}

func (c *core) acceptDecide() {
}

func (c *core) handleFinalCommitted() error {
	logger := c.newLogger()
	logger.Trace("handleFinalCommitted")
	c.startNewRound(common.Big0)
	return nil
}
