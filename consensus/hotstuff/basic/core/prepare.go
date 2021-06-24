package core

import (
	"fmt"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *core) sendPrepare() {
	logger := c.logger.New("state", c.current.State())

	msgTyp := MsgTypePrepare

	prepare, err := c.createNewProposal()
	if err != nil {
		logger.Error("Failed to creat leaf", "err", err)
		return
	}

	payload, err := Encode(prepare)
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}

	c.broadcast(&message{Code: msgTyp, Msg: payload})
}

func (c *core) handlePrepare(data *message, src hotstuff.Validator) error {
	logger := c.logger.New("state", c.currentState())

	var (
		msg    *MsgPrepare
		msgTyp = MsgTypePrepare
	)
	if err := data.Decode(&msg); err != nil {
		logger.Error("Failed to check msg", "type", msgTyp, "err", err)
		return errFailedDecodePrepare
	}
	if err := c.checkView(msgTyp, msg.View); err != nil {
		return err
	}
	if err := c.checkMsgFromProposer(src); err != nil {
		return err
	}

	proposal := msg.Proposal
	highQC := msg.HighQC
	if _, err := c.backend.VerifyUnsealedProposal(proposal); err != nil {
		logger.Error("Failed to verify unsealed proposal", "err", err)
		return errVerifyUnsealedProposal
	}
	if err := c.extend(proposal, highQC); err != nil {
		logger.Error("Failed to check extend", "err", err)
		return errExtend
	}
	if err := c.safeNode(proposal, highQC); err != nil {
		logger.Error("Failed to check safeNode", "err", err)
		return errSafeNode
	}

	c.current.SetProposal(proposal)
	c.sendPrepareVote()
	return nil
}

func (c *core) sendPrepareVote() {
	logger := c.logger.New("state", c.current.State())

	msgTyp := MsgTypePrepareVote
	sub := c.current.Vote()
	payload, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&message{Code: msgTyp, Msg: payload})
}

func (c *core) createNewProposal() (*MsgPrepare, error) {
	qc := c.current.PrepareQC()
	lastProposal, _ := c.backend.LastProposal()
	if lastProposal.Hash() != qc.Hash {
		return nil, fmt.Errorf("parent hash expect %s, got %s", lastProposal.Hash().Hex(), qc.Hash.Hex())
	}

	var req *hotstuff.Request
	if c.current.PendingRequest() != nil && c.current.PendingRequest().Proposal.Number().Cmp(c.current.Height()) == 0 {
		req = c.current.PendingRequest()
	} else {
		if req = c.requests.GetRequest(c.currentView()); req != nil {
			c.current.SetPendingRequest(req)
		} else {
			return nil, errNoRequest
		}
	}

	return &MsgPrepare{
		View:     c.currentView(),
		Proposal: req.Proposal,
		HighQC:   c.getHighQC(),
	}, nil
}

func (c *core) extend(proposal hotstuff.Proposal, highQC *hotstuff.QuorumCert) error {
	if err := c.backend.VerifyQuorumCert(highQC); err != nil {
		return err
	}
	block, ok := proposal.(*types.Block)
	if !ok {
		return fmt.Errorf("invalid proposal: hash %s", proposal.Hash())
	}
	if highQC.Hash != block.ParentHash() {
		return fmt.Errorf("block %s not extend hiqhQC %s", block.Hash().String(), highQC.Hash)
	}
	return nil
}

// proposal extend lockedQC `OR` hiqhQC.view > lockedQC.view
func (c *core) safeNode(proposal hotstuff.Proposal, highQC *hotstuff.QuorumCert) error {
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
