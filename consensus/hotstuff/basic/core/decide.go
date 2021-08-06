package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) handleCommitVote(data *message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		vote   *Vote
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
	if vote.Digest != c.current.PreCommittedQC().Hash {
		logger.Trace("Failed to check hash", "msg", msgTyp, "expect vote", c.current.PreCommittedQC().Hash, "got", vote.Digest)
		return errInvalidDigest
	}
	if err := c.checkMsgToProposer(); err != nil {
		logger.Trace("Failed to check proposer", "msg", msgTyp, "err", err)
		return err
	}

	if err := c.current.AddCommitVote(data); err != nil {
		logger.Trace("Failed to add vote", "msg", msgTyp, "err", err)
		return errAddPreCommitVote
	}

	logger.Trace("handleCommitVote", "msg", msgTyp, "src", src.Address(), "hash", vote.Digest)

	if size := c.current.CommitVoteSize(); size >= c.Q() && c.currentState() < StateCommitted {
		c.current.SetState(StateCommitted)
		c.current.SetCommittedQC(c.current.PreCommittedQC())
		logger.Trace("acceptCommit", "msg", msgTyp, "src", src.Address(), "hash", vote.Digest, "msgSize", size)
		if err := c.backend.Commit(c.current.Proposal()); err != nil {
			logger.Trace("Failed to commit proposal", "err", err)
			return err
		}
		c.startNewRound(common.Big0)
	}

	return nil
}

// handleFinalCommitted start new round if consensus engine accept notify signal from miner.worker.
// signals should be related with sync header or body. in fact, we DONT need this function to start an new round,
// because that the function `startNewRound` will sync header to preparing new consensus round args.
// we just kept it here for backup.
func (c *core) handleFinalCommitted() error {
	logger := c.newLogger()
	logger.Trace("handleFinalCommitted")
	c.startNewRound(common.Big0)
	return nil
}
