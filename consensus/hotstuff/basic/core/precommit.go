package core

import (
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) handlePrepareVote(data *message, src hotstuff.Validator) error {
	logger := c.newLogger()

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
		seals := c.getMessageSeals(size)
		newProposal, err := c.backend.PreCommit(c.current.Proposal(), seals)
		if err != nil {
			logger.Error("Failed to assemble committed seal", "err", err)
			return err
		}
		prepareQC := proposal2QC(newProposal, c.current.Round())
		c.current.SetProposal(newProposal)
		c.current.SetPrepareQC(prepareQC)
		c.current.SetState(StatePrepared)
		c.sendPreCommit()
	}

	logger.Trace("handlePrepareVote", "src", src.Address(), "msg view", vote.View, "vote", vote.Digest)
	return nil
}

func (c *core) sendPreCommit() {
	logger := c.newLogger()

	msgTyp := MsgTypePreCommit
	msg := &MsgPreCommit{
		View:      c.currentView(),
		Proposal:  c.current.Proposal(),
		PrepareQC: c.current.PrepareQC(),
	}
	payload, err := Encode(msg)
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&message{Code: msgTyp, Msg: payload})
	logger.Trace("sendPreCommit", "msg view", msg.View, "proposal", msg.Proposal.Hash())
}

func (c *core) handlePreCommit(data *message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		msg    *MsgPreCommit
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
	if msg.Proposal.Hash() != msg.PrepareQC.Hash {
		return errInvalidProposal
	}
	if _, err := c.backend.Verify(msg.Proposal); err != nil {
		return err
	}
	if err := c.signer.VerifyQC(msg.PrepareQC, c.valSet); err != nil {
		logger.Error("Failed to verify prepareQC", "err", err)
		return errVerifyQC
	}

	if !c.IsProposer() && c.current.state < StatePrepared {
		c.current.SetPrepareQC(msg.PrepareQC)
		c.current.SetProposal(msg.Proposal)
		c.current.SetState(StatePrepared)
	}
	c.sendPreCommitVote()

	logger.Trace("handlePreCommit", "src", src.Address(),"msg view", msg.View, "proposal", msg.Proposal.Hash())
	return nil
}

func (c *core) sendPreCommitVote() {
	logger := c.newLogger()

	msgTyp := MsgTypePreCommitVote
	vote := c.current.Vote()
	if vote == nil {
		logger.Error("proposal is nil")
		return
	}
	payload, err := Encode(vote)
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&message{Code: msgTyp, Msg: payload})
	logger.Trace("sendPreCommitVote", "vote view", vote.View, "vote", vote.Digest)
}
