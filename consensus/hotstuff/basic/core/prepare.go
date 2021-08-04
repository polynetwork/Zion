package core

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *core) sendPrepare() {
	logger := c.newLogger()

	if !c.IsProposer() {
		return
	}

	// filter dump proposal
	if proposal := c.current.Proposal(); proposal != nil &&
		proposal.Number().Uint64() == c.currentView().Height.Uint64() &&
		proposal.Coinbase() == c.Address() {

		logger.Trace("Failed to send prepare", "err", "proposal already been sent")
		return
	}

	msgTyp := MsgTypePrepare
	if !c.current.IsProposalLocked() {
		proposal, err := c.createNewProposal()
		if err != nil {
			logger.Trace("Failed to create proposal", "err", err, "request set size", c.requests.Size(),
				"pendingRequest", c.current.PendingRequest(), "view", c.currentView())
			return
		}
		c.current.SetProposal(proposal)
	} else if c.current.Proposal() == nil {
		logger.Error("Failed to get locked proposal", "err", "locked proposal is nil")
		return
	}

	prepare := &MsgPrepare{
		View:     c.currentView(),
		Proposal: c.current.Proposal(),
		HighQC:   c.current.HighQC(),
	}
	payload, err := Encode(prepare)
	if err != nil {
		logger.Trace("Failed to encode", "msg", msgTyp, "err", err)
		return
	}

	// consensus spent time always less than a block period, waiting for `delay` time to catch up the system time.
	delay := time.Unix(int64(prepare.Proposal.Time()), 0).Sub(time.Now())
	time.Sleep(delay)
	logger.Trace("delay to broadcast proposal", "time", delay.Milliseconds())

	c.broadcast(&message{Code: msgTyp, Msg: payload})
	logger.Trace("sendPrepare", "prepare view", prepare.View, "proposal", prepare.Proposal.Hash())
}

func (c *core) handlePrepare(data *message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		msg    *MsgPrepare
		msgTyp = MsgTypePrepare
	)
	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to decode", "type", msgTyp, "err", err)
		return errFailedDecodePrepare
	}
	if err := c.checkView(msgTyp, msg.View); err != nil {
		logger.Trace("Failed to check view", "msg", msgTyp, "err", err)
		return err
	}
	if err := c.checkMsgFromProposer(src); err != nil {
		logger.Trace("Failed to check proposer", "msg", msgTyp, "err", err)
		return err
	}

	if _, err := c.backend.VerifyUnsealedProposal(msg.Proposal); err != nil {
		logger.Trace("Failed to verify unsealed proposal", "msg", msgTyp, "err", err)
		return errVerifyUnsealedProposal
	}
	if err := c.extend(msg.Proposal, msg.HighQC); err != nil {
		logger.Trace("Failed to check extend", "msg", msgTyp, "err", err)
		return errExtend
	}
	if err := c.safeNode(msg.Proposal, msg.HighQC); err != nil {
		logger.Trace("Failed to check safeNode", "msg", msgTyp, "err", err)
		return errSafeNode
	}
	if err := c.checkLockedProposal(msg.Proposal); err != nil {
		logger.Trace("Failed to check locked proposal", "msg", msgTyp, "err", err)
		return err
	}

	logger.Trace("handlePrepare", "msg", msgTyp, "src", src.Address(), "hash", msg.Proposal.Hash())

	if c.IsProposer() && c.currentState() < StatePrepared {
		c.sendPrepareVote()
	}
	if !c.IsProposer() && c.currentState() < StateHighQC {
		c.current.SetHighQC(msg.HighQC)
		c.current.SetProposal(msg.Proposal)
		c.current.SetState(StateHighQC)
		logger.Trace("acceptHighQC", "msg", msgTyp, "src", src.Address(), "highQC", msg.HighQC.Hash)

		c.sendPrepareVote()
	}

	return nil
}

func (c *core) sendPrepareVote() {
	logger := c.newLogger()

	msgTyp := MsgTypePrepareVote
	vote := c.current.Vote()
	if vote == nil {
		logger.Trace("Failed to send vote", "msg", msgTyp, "err", "current vote is nil")
		return
	}
	payload, err := Encode(vote)
	if err != nil {
		logger.Trace("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(&message{Code: msgTyp, Msg: payload})
	logger.Trace("sendPrepareVote", "vote view", vote.View, "vote", vote.Digest)
}

func (c *core) createNewProposal() (hotstuff.Proposal, error) {
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
	return req.Proposal, nil
}

func (c *core) extend(proposal hotstuff.Proposal, highQC *hotstuff.QuorumCert) error {
	block, ok := proposal.(*types.Block)
	if !ok {
		return fmt.Errorf("invalid proposal: hash %s", proposal.Hash())
	}
	if err := c.signer.VerifyQC(highQC, c.valSet); err != nil {
		return err
	}
	if highQC.Hash != block.ParentHash() {
		return fmt.Errorf("block %v (parent %v) not extend hiqhQC %v", block.Hash(), block.ParentHash(), highQC.Hash)
	}
	return nil
}

// proposal extend lockedQC `OR` hiqhQC.view > lockedQC.view
func (c *core) safeNode(proposal hotstuff.Proposal, highQC *hotstuff.QuorumCert) error {
	logger := c.newLogger()

	if proposal.Number().Uint64() == 1 {
		return nil
	}
	safety := false
	liveness := false
	if c.current.PreCommittedQC() == nil {
		logger.Trace("safeNodeChecking", "lockQC", "is nil")
		return errSafeNode
	}
	if err := c.extend(proposal, c.current.PreCommittedQC()); err == nil {
		safety = true
	} else {
		logger.Trace("safeNodeChecking", "extend err", err)
	}
	if highQC.View.Cmp(c.current.PreCommittedQC().View) > 0 {
		liveness = true
	}
	if safety || liveness {
		return nil
	} else {
		return errSafeNode
	}
}
