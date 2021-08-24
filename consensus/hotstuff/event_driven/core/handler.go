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
	"github.com/ethereum/go-ethereum/event"
)

type EventSender struct {
	eventMtx *event.TypeMux
}

func NewEventSender(backend hotstuff.Backend) *EventSender {
	return &EventSender{eventMtx: backend.EventMux()}
}

func (s *EventSender) sendEvent(val interface{}) {
	s.eventMtx.Post(val)
}

func (e *EventDrivenEngine) subscribe() {

}

func (e *EventDrivenEngine) Handler() {

}

