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

// sendNewView, repo send message of new-view, formula as follow:
// 	MSG(new-view, _, prepareQC)
// the field of view will be packaged in message before broadcast.
func (c *core) sendNewView() {
	logger := c.newLogger()
	code := MsgTypeNewView

	prepareQC := c.current.PrepareQC()
	payload, err := Encode(prepareQC)
	if err != nil {
		logger.Trace("Failed to encode", "msg", code, "err", err)
		return
	}

	c.broadcast(code, payload)
	logger.Trace("sendNewView", "msg", code)
}

// handleNewView, leader gather new-view messages and pick the max `prepareQC` to be `highQC` by view sequence.
// `stateHighQC` denote that node is ready to pack block to send the `prepare` message.
func (c *core) handleNewView(data *Message) error {
	var (
		logger    = c.newLogger()
		prepareQC *QuorumCert
		code      = data.Code
		src       = data.address
	)

	// check message
	if err := data.Decode(&prepareQC); err != nil {
		logger.Trace("Failed to decode", "msg", code, "src", src, "err", err)
		return errFailedDecodeNewView
	}
	if err := c.checkView(data.View); err != nil {
		logger.Trace("Failed to check view", "msg", code, "src", src, "err", err)
		return err
	}
	if err := c.checkMsgDest(); err != nil {
		logger.Trace("Failed to check proposer", "msg", code, "src", src, "err", err)
		return err
	}

	// ensure remote `prepareQC` is legal.
	if err := c.verifyQC(data, prepareQC); err != nil {
		logger.Trace("Failed to verify prepareQC", "msg", code, "src", src, "err", err)
		return err
	}
	// messages queued in messageSet to ensure there will be at least 2/3 validators on the same step
	if err := c.current.AddNewViews(data); err != nil {
		logger.Trace("Failed to add new view", "msg", code, "src", src, "err", err)
		return errAddNewViews
	}

	logger.Trace("handleNewView", "msg", code, "src", src, "prepareQC", prepareQC.node)

	if size := c.current.NewViewSize(); size >= c.Q() && c.currentState() < StateHighQC {
		highQC, err := c.getHighQC()
		if err != nil {
			logger.Trace("Failed to get highQC", "msg", code, "err", err)
			return err
		}
		c.current.SetHighQC(highQC)
		c.setCurrentState(StateHighQC)

		logger.Trace("acceptHighQC", "msg", code, "prepareQC", prepareQC.node, "msgSize", size)
		c.sendPrepare()
	}

	return nil
}

// getHighQC leader find the highest `prepareQC` as highQC by `view` sequence.
func (c *core) getHighQC() (*QuorumCert, error) {
	var highQC *QuorumCert
	for _, data := range c.current.NewViews() {
		var qc *QuorumCert
		if err := data.Decode(&qc); err != nil {
			return nil, err
		}
		if highQC == nil || highQC.view.Cmp(qc.view) < 0 {
			highQC = qc
		}
	}
	if highQC == nil {
		return nil, errNilHighQC
	}
	return highQC, nil
}
