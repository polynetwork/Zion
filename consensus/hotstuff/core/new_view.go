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

func (c *core) sendNewView(view *View) {
	logger := c.newLogger()
	msgTyp := MsgTypeNewView

	curView := c.currentView()
	if curView.Cmp(view) > 0 {
		logger.Trace("Cannot send out the round change", "msg", msgTyp, "current round", curView.Round, "target round", view.Round)
		return
	}

	newViewMsg := c.current.PrepareQC()
	payload, err := Encode(newViewMsg)
	if err != nil {
		logger.Trace("Failed to encode", "msg", msgTyp, "err", err)
		return
	}
	c.broadcast(msgTyp, payload)

	logger.Trace("sendNewView", "msg", msgTyp)
}

func (c *core) handleNewView(data *Message) error {
	logger := c.newLogger()

	var (
		prepareQC *QuorumCert
		code      = MsgTypeNewView
		src       = data.address
	)

	if err := data.Decode(&prepareQC); err != nil {
		logger.Trace("Failed to decode", "msg", code, "src", src, "err", err)
		return errFailedDecodeNewView
	}
	if err := c.checkView(code, data.View); err != nil {
		logger.Trace("Failed to check view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkMsgToProposer(); err != nil {
		logger.Trace("Failed to check proposer", "msg", code, "src", src, "err", err)
		return err
	}

	//// todo(fuk): verify cross epoch qc
	//if err := c.verifyCrossEpochQC(msg.PrepareQC); err != nil {
	//	logger.Trace("Failed to verify highQC", "msg", msgTyp, "src", src, "err", err)
	//	return err
	//}

	if err := c.verifyVoteQC(prepareQC.hash, prepareQC); err != nil {
		logger.Trace("Failed to verify highQC", "msg", code, "src", src, "err", err)
		return err
	}

	// `checkView` ensure that +2/3 validators on the same view
	if err := c.current.AddNewViews(data); err != nil {
		logger.Trace("Failed to add new view", "msg", code, "src", src, "err", err)
		return errAddNewViews
	}

	logger.Trace("handleNewView", "msg", code, "src", src, "prepareQC", prepareQC.hash)

	if size := c.current.NewViewSize(); size >= c.Q() && c.currentState() < StateHighQC {
		if highQC, err := c.getHighQC(); err != nil || highQC == nil {
			logger.Trace("Failed to get highQC", "msg", code)
			return errGetHighQC
		} else {
			c.current.SetHighQC(highQC)
			c.setCurrentState(StateHighQC)

			logger.Trace("acceptHighQC", "msg", code, "prepareQC", prepareQC.hash, "msgSize", size)
			c.sendPrepare()
		}
	}

	return nil
}

func (c *core) getHighQC() (*QuorumCert, error) {
	var maxView *QuorumCert
	for _, data := range c.current.NewViews() {
		var qc *QuorumCert
		if err := data.Decode(&qc); err != nil {
			return nil, err
		}
		if maxView == nil || maxView.view.Cmp(qc.view) < 0 {
			maxView = qc
		}
	}
	return maxView, nil
}
