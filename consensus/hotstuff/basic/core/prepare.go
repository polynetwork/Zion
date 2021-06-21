package core

import (
	"fmt"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *core) sendPrepare() {
	logger := c.logger.New("state", c.current.State())

	if !c.IsProposer() {
		return
	}
	msgTyp := MsgTypePrepare

	proposal, err := c.createNewProposal()
	if err != nil {
		logger.Error("Failed to creat leaf", "err", err)
		return
	}

	payload, err := Encode(proposal)
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}

	c.broadcast(&message{Code: msgTyp, Msg: payload})
}

func (c *core) handlePrepare(data *message, src hotstuff.Validator) {
	logger := c.logger.New("state", c.currentState())

	var (
		msg    *MsgNewProposal
		msgTyp = MsgTypePrepare
	)
	if err := c.decodeAndCheckProposal(data, msgTyp, msg); err != nil {
		logger.Error("Failed to check msg", "type", msgTyp, "err", err)
		return
	}

	proposal := msg.Proposal
	highQC := msg.HighQC
	if _, err := c.backend.VerifyUnsealedProposal(proposal); err != nil {
		logger.Error("Failed to verify unsealed proposal", "err", err)
		return
	}
	if err := c.extend(proposal, highQC); err != nil {
		logger.Error("Failed to check extend", "err", err)
		return
	}
	if err := c.safeNode(proposal, highQC); err != nil {
		logger.Error("Failed to check safeNode", "err", err)
		return
	}

	c.current.SetProposal(proposal)
	c.sendPrepareVote()
}

func (c *core) sendPrepareVote() {
	logger := c.logger.New("state", c.current.State())

	msgTyp := MsgTypePrepareVote
	sub := c.current.Subject()
	payload, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&message{Code: msgTyp, Msg: payload})
}

func (c *core) createNewProposal() (*MsgNewProposal, error) {
	qc := c.current.PrepareQC()
	lastProposal, _ := c.backend.LastProposal()
	if lastProposal.Hash() != qc.Proposal.Hash() {
		return nil, fmt.Errorf("parent hash expect %s, got %s", lastProposal.Hash().Hex(), qc.Proposal.Hash().Hex())
	}

	req := c.requests.GetRequest(c.currentView())
	return &MsgNewProposal{
		View:     c.currentView(),
		Proposal: req.Proposal,
		HighQC:   c.getHighQC(),
	}, nil
}

func (c *core) extend(proposal hotstuff.Proposal, highQC *QuorumCert) error {
	if _, err := c.backend.Verify(highQC.Proposal); err != nil {
		return err
	}
	block, ok := proposal.(*types.Block)
	if !ok {
		return fmt.Errorf("invalid proposal: hash %s", proposal.Hash())
	}
	targetBlock, ok := highQC.Proposal.(*types.Block)
	if !ok {
		return fmt.Errorf("invalid highQC proposal: hash %s", highQC.Proposal.Hash())
	}
	if targetBlock.Hash() != block.ParentHash() {
		return fmt.Errorf("block %s not extend hiqhQC %s", block.Hash().String(), targetBlock.Hash().String())
	}
	return nil
}

// proposal extend lockedQC `OR` hiqhQC.view > lockedQC.view
func (c *core) safeNode(proposal hotstuff.Proposal, highQC *QuorumCert) error {
	safety := false
	liveness := false
	if err := c.extend(proposal, c.current.LockedQC()); err == nil {
		safety = true
	}
	if highQC.View.Cmp(c.current.LockedQC().View) > 0 {
		liveness = true
	}
	if safety || liveness {
		return nil
	} else {
		return errSafeNode
	}
}
