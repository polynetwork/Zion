package core

import (
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) sendNewView(view *hotstuff.View) {
	logger := c.logger.New("state", c.currentState())

	curView := c.currentView()
	if curView.Cmp(view) >= 0 {
		logger.Error("Cannot send out the round change", "current round", curView.Round, "target round", view.Round)
		return
	}

	newViewMsg := c.current.PrepareQC().Copy()
	payload, err := Encode(newViewMsg)
	if err != nil {
		logger.Error("Failed to encode", "msg", MsgTypeNewView, "err", err)
		return
	}
	c.broadcast(&message{
		Code: MsgTypeNewView,
		Msg:  payload,
	})
}

func (c *core) handleNewView(data *message, src hotstuff.Validator) {
	logger := c.logger.New("state", c.currentState())

	var (
		msg    *hotstuff.QuorumCert
		msgTyp = MsgTypeNewView
	)
	if err := c.decodeAndCheckMessage(data, msgTyp, msg); err != nil {
		logger.Error("Failed to check msg", "type", msgTyp, "err", err)
		return
	}

	if err := c.backend.VerifyQuorumCert(msg); err != nil {
		logger.Error("Failed to verify proposal", "err", err)
		return
	}

	if err := c.current.AddNewViews(data); err != nil {
		logger.Error("Failed to add new view", "err", err)
		return
	}

	if c.current.NewViewSize() == c.Q() {
		highQC := c.getHighQC()
		c.current.SetHighQC(highQC)
	}
}

func (c *core) getHighQC() *hotstuff.QuorumCert {
	var maxView *hotstuff.QuorumCert
	for _, data := range c.current.NewViews() {
		var msg *hotstuff.QuorumCert
		data.Decode(&msg)
		if maxView == nil || maxView.View.Cmp(msg.View) < 0 {
			maxView = msg
		}
	}
	return maxView
}
