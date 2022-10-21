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
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

// Construct a new message set to accumulate messages for given height/view number.
func NewMessageSet(valSet hotstuff.ValidatorSet) *MessageSet {
	return &MessageSet{
		view: &View{
			Round:  new(big.Int),
			Height: new(big.Int),
		},
		mu:   new(sync.RWMutex),
		msgs: make(map[common.Address]*Message),
		vs:   valSet,
	}
}

type MessageSet struct {
	view *View
	vs   hotstuff.ValidatorSet
	mu   *sync.RWMutex
	msgs map[common.Address]*Message
}

func (s *MessageSet) View() *View {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.view
}

func (s *MessageSet) Add(msg *Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if index, v := s.vs.GetByAddress(msg.Address); index < 0 || v == nil {
		return fmt.Errorf("unauthorized address")
	}

	s.msgs[msg.Address] = msg
	return nil
}

func (s *MessageSet) Values() (result []*Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, v := range s.msgs {
		result = append(result, v)
	}
	return
}

func (s *MessageSet) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.msgs)
}

func (s *MessageSet) Get(addr common.Address) *Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.msgs[addr]
}

func (s *MessageSet) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	addresses := make([]string, 0, len(s.msgs))
	for _, v := range s.msgs {
		addresses = append(addresses, v.Address.Hex())
	}
	return fmt.Sprintf("[%v]", strings.Join(addresses, ", "))
}
