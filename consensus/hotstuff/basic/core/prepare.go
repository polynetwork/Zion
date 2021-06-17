package core

import (
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) sendPrepare(request *hotstuff.Request) {
	if c.current.Height().Cmp(request.Proposal.Number()) != 0 || !c.IsProposer() {
		return
	}

	logger := c.logger.New("state", c.state)
	curView := c.currentView()
	prepare, err := Encode(&MsgPrepare{
		View:     curView,
		Proposal: request.Proposal,
	})
	if err != nil {
		logger.Error("Failed to encode", MsgTypePrepare.String(), curView)
		return
	}
	c.broadcast(&message{
		Code:      MsgTypePrepare,
		Msg:       prepare,
	})
}

func (c *core) handlePrepare(msg *message, src hotstuff.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	var prepare *MsgPrepare
	if err := msg.Decode(&prepare); err != nil {
		return errFailedDecodePrepare
	}

	if err := c.checkMessage(MsgTypePrepare, prepare.View); err != nil {
		return err
	}

	if err := c.verifyPrepare(prepare, src); err != nil {
		logger.Warn("err", err)
		return err
	}

	c.acceptPrepare(prepare)
	c.sendPrepareVote()

	return nil
}

func (c *core) verifyPrepare(msg *MsgPrepare, src hotstuff.Validator) error {
	if !c.valSet.IsProposer(src.Address()) {
		return errNotFromProposer
	}
	if c.current.IsHashLocked() {
		return errHashAlreayLocked
	}
	if _, err := c.backend.Verify(msg.Proposal); err != nil {
		return err
	}
	return nil
}

func (c *core) acceptPrepare(prepare *MsgPrepare) {
	//c.consensusTimestamp = time.Now()
	c.current.SetPrepare(prepare)
}
