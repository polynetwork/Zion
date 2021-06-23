package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) handleCommitVote(data *message, src hotstuff.Validator) {
	logger := c.logger.New("state", c.currentState())

	var (
		vote   *hotstuff.Subject
		msgTyp = MsgTypeCommitVote
	)
	if err := c.decodeAndCheckVote(data, msgTyp, vote); err != nil {
		logger.Error("Failed to check vote", "vote", msgTyp, "err", err)
		return
	}

	if size := c.current.CommitVoteSize(); size >= c.Q() && c.currentState() < StateCommitted {
		c.current.SetState(StateCommitted)
		c.current.SetCommittedQC(c.current.LockedQC())
		if err := c.backend.Commit(c.current.Proposal()); err != nil {
			logger.Error("Failed to commit proposal", "err", err)
			return
		}
	}
}

func (c *core) handleFinalCommitted() {
	logger := c.logger.New("state", c.currentState())
	logger.Trace("Received a final committed proposal")
	c.startNewRound(common.Big0)
}
