/*
 * Copyright (C) 2021 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */

package core

import (
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) sendNewView(view *hotstuff.View) {
	logger := c.newLogger()

	curView := c.currentView()
	if curView.Cmp(view) > 0 {
		logger.Trace("Cannot send out the round change", "current round", curView.Round, "target round", view.Round)
		return
	}

	prepareQC := c.current.PrepareQC().Copy()
	newViewMsg := &MsgNewView{
		View:      curView,
		PrepareQC: prepareQC,
	}
	payload, err := Encode(newViewMsg)
	if err != nil {
		logger.Trace("Failed to encode", "msg", MsgTypeNewView, "err", err)
		return
	}
	c.broadcast(&hotstuff.Message{
		Code: MsgTypeNewView,
		Msg:  payload,
	})

	logger.Trace("sendNewView", "prepareQC", prepareQC.Hash)
}

func (c *core) handleNewView(data *hotstuff.Message, src hotstuff.Validator) error {
	logger := c.newLogger()

	var (
		msg    *MsgNewView
		msgTyp = MsgTypeNewView
	)

	if err := data.Decode(&msg); err != nil {
		logger.Trace("Failed to decode", "msg", msgTyp, "err", err)
		return errFailedDecodeNewView
	}
	if err := c.checkView(msgTyp, msg.View); err != nil {
		logger.Trace("Failed to check view", "msg", msgTyp, "err", err)
		return err
	}
	if err := c.checkMsgToProposer(); err != nil {
		logger.Trace("Failed to check proposer", "msg", msgTyp, "err", err)
		return err
	}

	if err := c.verifyCrossEpochQC(msg.PrepareQC); err != nil {
		logger.Trace("Failed to verify highQC", "msg", msgTyp, "err", err)
		return err
	}

	if err := c.current.AddNewViews(data); err != nil {
		logger.Trace("Failed to add new view", "msg", msgTyp, "err", err)
		return errAddNewViews
	}

	logger.Trace("handleNewView", "msg", msgTyp, "src", src.Address(), "prepareQC", msg.PrepareQC.Hash)

	if size := c.current.NewViewSize(); size >= c.Q() && c.currentState() < StateHighQC {
		highQC := c.getHighQC()
		c.current.SetHighQC(highQC)
		c.current.SetState(StateHighQC)
		logger.Trace("acceptHighQC", "msg", msgTyp, "src", src.Address(), "prepareQC", msg.PrepareQC.Hash, "msgSize", size)

		c.sendPrepare()
	}

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
