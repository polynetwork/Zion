package core

import (
	"math/big"
)

func (c *core) sendPrepare(msg *MsgPrepare) {
	logger := c.logger.New("state", c.state)

	if c.current.Height().Cmp(msg.Proposal.Number()) == 0 && c.IsProposer() {
		curView := c.currentView()
		payload, err := Encode(&MsgPrepare{
			View:     curView,
			Proposal: msg.Proposal,
		})
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}
		c.broadcast(&message{
			Code: MsgTypePrepare,
			Msg:  payload,
		}, new(big.Int))
	}
}

func (c *core) handlePrepare(msg *message) error {
	// logger := c.logger.New("state", c.state)

	var prepare *MsgPrepare
	if err := msg.Decode(&prepare); err != nil {
		return errFailedDecodePrepare
	}

	if _, err := c.backend.Verify(prepare.Proposal); err != nil {
		return err
	}

	c.current.SetPrepareQC(&QuorumCert{
		Type:     MsgTypePrepare,
		Proposal: prepare.Proposal,
	})

	c.setState(StatePrepared)

	if !c.IsProposer() {

	}
	return nil
}
