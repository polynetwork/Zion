package core

import (
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) sendNewView(view *hotstuff.View) {
	logger := c.newLogger()

	curView := c.currentView()
	if curView.Cmp(view) > 0 {
		logger.Error("Cannot send out the round change", "current round", curView.Round, "target round", view.Round)
		return
	}

	prepareQC := c.current.PrepareQC().Copy()
	newViewMsg := &MsgNewView{
		View:      curView,
		PrepareQC: prepareQC,
	}
	payload, err := Encode(newViewMsg)
	if err != nil {
		logger.Error("Failed to encode", "msg", MsgTypeNewView, "err", err)
		return
	}
	c.broadcast(&message{
		Code: MsgTypeNewView,
		Msg:  payload,
	})

	logger.Info("sendNewView", "prepareQC", prepareQC.Hash)
}

func (c *core) handleNewView(data *message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		msg    *MsgNewView
		msgTyp = MsgTypeNewView
	)

	if err := data.Decode(&msg); err != nil {
		return errFailedDecodeNewView
	}
	if err := c.checkView(msgTyp, msg.View); err != nil {
		return err
	}
	if err := c.checkMsgToProposer(); err != nil {
		return err
	}

	if err := c.signer.VerifyQC(msg.PrepareQC, c.valSet); err != nil {
		logger.Error("Failed to verify highQC", "err", err)
		return errVerifyQC
	}

	if err := c.current.AddNewViews(data); err != nil {
		logger.Error("Failed to add new view", "err", err)
		return errAddNewViews
	}

	if c.current.NewViewSize() == c.Q() {
		highQC := c.getHighQC()
		c.current.SetHighQC(highQC)
		c.sendPrepare()
	}

	logger.Trace("handleNewView", "src", src.Address(), "msg view", msg.View, "prepareQC", msg.PrepareQC.Hash)
	return nil
}

func (c *core) getHighQC() *hotstuff.QuorumCert {
	var maxView *hotstuff.QuorumCert
	for _, data := range c.current.NewViews() {
		var msg *MsgNewView
		data.Decode(&msg)
		if maxView == nil || maxView.View.Cmp(msg.PrepareQC.View) < 0 {
			maxView = msg.PrepareQC
		}
	}
	return maxView
}
