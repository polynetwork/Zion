package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) handleCommitVote(data *message, src hotstuff.Validator) error {
	logger := c.logger.New("state", c.currentState())

	var (
		vote   *hotstuff.Vote
		msgTyp = MsgTypeCommitVote
	)
	if err := data.Decode(&vote); err != nil {
		logger.Error("Failed to check vote", "vote", msgTyp, "err", err)
		return errFailedDecodeCommitVote
	}
	if err := c.checkVote(vote); err != nil {
		return err
	}
	if err := c.checkView(msgTyp, vote.View); err != nil {
		return err
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
	}

	return nil
}

func (c *core) handleFinalCommitted() error {
	logger := c.logger.New("state", c.currentState())
	logger.Trace("Received a final committed proposal")
	c.startNewRound(common.Big0)
	return nil
}
