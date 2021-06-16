package core

import "math/big"

func (c *core) sendPreCommit() {
	logger := c.logger.New("state", c.state)

	prepareQC := c.current.PrepareQC()
	curView := c.currentView()

	if prepareQC != nil && prepareQC.Proposal.Number().Cmp(curView.Height) == 0 && c.IsProposer() {
		payload, err := Encode(&MsgPreCommit{
			View:     curView,
			Proposal: prepareQC.Proposal,
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

func (c *core) handlePreCommit(msg *message) error {
	var precm *MsgPreCommit
	if err := msg.Decode(&precm); err != nil {
		return errFailedDecodePreCommit
	}

	if _, err := c.backend.Verify(precm.Proposal); err != nil {
		return err
	}

	c.current.SetPrepareQC(&QuorumCert{
		Type:     MsgTypePreCommit,
		Proposal: precm.Proposal,
	})

	c.setState(StateLocked)

	if !c.IsProposer() {
		c.sendPreCommitVote()
	}
	return nil
}
