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
		logger.Trace("Failed to decode", "msg", msgTyp, "err", err)
		return errFailedDecodeCommitVote
	}
	if err := c.checkView(msgTyp, vote.View); err != nil {
		logger.Trace("Failed to check view", "msg", msgTyp, "err", err)
		return err
	}
	if err := c.checkVote(vote); err != nil {
		logger.Trace("Failed to check vote", "msg", msgTyp, "err", err)
		return err
	}
	if vote.Digest != c.current.LockedQC().Hash {
		logger.Trace("Failed to check hash", "msg", msgTyp, "expect vote", c.current.LockedQC().Hash.Hex(), "got", vote.Digest.Hex())
		return errInvalidDigest
	}
	if err := c.checkMsgToProposer(); err != nil {
		logger.Trace("Failed to check proposer", "msg", msgTyp, "err", err)
		return err
	}

	if err := c.current.AddCommitVote(data); err != nil {
		logger.Trace("Failed to add vote", "type", msgTyp, "err", err)
		return errAddPreCommitVote
	}

	logger.Trace("handleCommitVote", "src", src.Address(), "vote view", vote.View, "vote", vote.Digest)

	if size := c.current.CommitVoteSize(); size >= c.Q() && c.currentState() < StateCommitted {
		if c.acceptCommitVote(src, vote) {
			c.startNewRound(common.Big0)
		}
	}

	return nil
}

func (c *core) acceptCommitVote(src hotstuff.Validator, msg *hotstuff.Vote) bool {
	logger := c.newLogger()

	c.current.SetState(StateCommitted)
	c.current.SetCommittedQC(c.current.LockedQC())
	if err := c.backend.Commit(c.current.Proposal()); err != nil {
		logger.Trace("Failed to commit proposal", "err", err)
		return false
	}
	return true
}

func (c *core) handleFinalCommitted() error {
	logger := c.newLogger()
	logger.Trace("handleFinalCommitted")
	c.startNewRound(common.Big0)
	return nil
}
