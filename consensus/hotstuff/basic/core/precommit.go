package core

import (
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) handlePrepareVote(data *message, src hotstuff.Validator) error {
	logger := c.logger.New("state", c.currentState())

	var (
		vote   *hotstuff.Vote
		msgTyp = MsgTypePrepareVote
	)
	if err := data.Decode(&vote); err != nil {
		logger.Error("Failed to check vote", "type", msgTyp, "err", err)
		return errFailedDecodePrepareVote
	}
	if err := c.checkView(msgTyp, vote.View); err != nil {
		return err
	}
	if err := c.checkVote(vote); err != nil {
		return err
	}
	if err := c.checkMsgToProposer(); err != nil {
		return err
	}
	if err := c.current.AddPrepareVote(data); err != nil {
		logger.Error("Failed to add vote", "type", msgTyp, "err", err)
		return errAddPrepareVote
	}

	if size := c.current.PrepareVoteSize(); size >= c.Q() && c.current.state < StatePrepared {
		seal := c.getSeals(size)
		newProposal, prepareQC, err := c.backend.PreCommit(c.currentView(), c.current.Proposal(), seal)
		if err != nil {
			logger.Error("Failed to assemble committed seal", "err", err)
			return err
		}
		c.current.SetProposal(newProposal)
		c.current.SetPrepareQC(prepareQC)
		c.current.SetState(StatePrepared)
		c.sendPreCommit()
	}

	return nil
}

func (c *core) getSeals(n int) [][]byte {
	seals := make([][]byte, n)
	for i, data := range c.current.PrepareVotes() {
		if i < n {
			seals[i] = data.Signature
		}
	}
	return seals
}

func (c *core) sendPreCommit() {
	logger := c.logger.New("state", c.current.State())

	msgTyp := MsgTypePreCommit
	payload, err := Encode(c.current.PrepareQC())
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&message{Code: msgTyp, Msg: payload})
}

func (c *core) handlePreCommit(data *message, src hotstuff.Validator) error {
	logger := c.logger.New("state", c.currentState())

	var (
		msg    *hotstuff.QuorumCert
		msgTyp = MsgTypePreCommit
	)
	if err := data.Decode(&msg); err != nil {
		logger.Error("Failed to check msg", "type", msgTyp, "err", err)
		return errFailedDecodePreCommit
	}
	if err := c.checkView(MsgTypePreCommit, msg.View); err != nil {
		return err
	}
	if err := c.checkMsgFromProposer(src); err != nil {
		return err
	}
	// todo compare high qc
	if err := c.backend.VerifyQuorumCert(msg); err != nil {
		logger.Error("Failed to verify proposal", "err", err)
		return errVerifyQC
	}
	// todo: compare state in other steps
	if c.current.state >= StatePrepared {
		return errState
	}

	c.current.SetPrepareQC(msg)
	c.current.SetState(StatePrepared)
	c.sendPreCommitVote()

	return nil
}

func (c *core) sendPreCommitVote() {
	logger := c.logger.New("state", c.current.State())

	msgTyp := MsgTypePreCommitVote
	sub := c.current.Vote()
	payload, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&message{Code: msgTyp, Msg: payload})
}
