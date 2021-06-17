package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

// sendNextRoundChange sends the ROUND CHANGE message with current round + 1
func (c *core) sendNextRoundChange() {
	cv := c.currentView()
	c.sendRoundChange(new(big.Int).Add(cv.Round, common.Big1))
}

func (c *core) sendRoundChange(round *big.Int) {
	logger := c.logger.New("state", c.state)

	cv := c.currentView()
	if cv.Round.Cmp(round) >= 0 {
		logger.Error("Cannot send out the round change", "current round", cv.Round, "target round", round)
		return
	}

	// Reset ROUND CHANGE timeout timer with new round number
	c.catchUpRound(&hotstuff.View{
		// The round number we'd like to transfer to.
		Round: new(big.Int).Set(round),
		// TODO: Need to check if (height - 1) is right
		// Height: new(big.Int).Set(cv.Height),
		Height: new(big.Int).Sub(cv.Height, common.Big1),
	})

	// allow leader send self prepareQC

	// Now we have the new round number and block height number
	cv = c.currentView()
	payload, err := Encode(&MsgNewView{
		View: c.currentView(),
		QC:   c.current.PrepareQC(),
	})
	if err != nil {
		logger.Error("Failed to encode ROUND CHANGE", "rc", rc, "err", err)
		return
	}

	c.broadcast(&message{
		Code: MsgTypeRoundChange,
		Msg:  payload,
	}, round)
}

func (c *core) handleNewView(msg *message, src hotstuff.Validator) error {
	logger := c.logger.New("state", c.state, "from", src.Address().Hex())

	if c.IsProposer() {
		return nil
	}

	// Decode ROUND CHANGE message
	var rc *MsgNewView
	if err := msg.Decode(&rc); err != nil {
		logger.Error("Failed to decode ROUND CHANGE", "err", err)
		return errInvalidMessage
	}

	// todo: check message
	// This make sure the view should be identical
	//if err := c.checkMessage(msgRoundChange, rc.View); err != nil {
	//	return err
	//}

	// todo:
	round := c.current.Round()

	// Add the ROUND CHANGE message to its message set and return how many
	// messages we've got with the same round number and sequence number.
	if err := c.acceptNewView(msg, round); err != nil {
		logger.Warn("Failed to add round change message", "from", src, "msg", msg, "err", err)
		return err
	}

	if c.current.NewViewSize() == c.valSet.Q() {
		list := c.current.newViews.Values()
		var maxView *MsgNewView
		for _, v := range list {
			var nv *MsgNewView
			v.Decode(&nv)
			if maxView == nil || nv.View.Cmp(maxView.View) >= 0 {
				maxView = nv
			}
		}
		// todo: round state high qc
		c.current.SetHighQC(maxView.QC)
	}

	return nil
}

func (c *core) acceptNewView(msg *message, round *big.Int) error {
	if err := c.current.AddNewView(msg); err != nil {
		return err
	}
	return nil
}
