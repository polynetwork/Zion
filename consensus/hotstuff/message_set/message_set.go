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

package message_set

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
		view: &hotstuff.View{
			Round:  new(big.Int),
			Height: new(big.Int),
		},
		mtx:  new(sync.Mutex),
		msgs: make(map[common.Address]*hotstuff.Message),
		vs:   valSet,
	}
}

type MessageSet struct {
	view *hotstuff.View
	vs   hotstuff.ValidatorSet
	mtx  *sync.Mutex
	msgs map[common.Address]*hotstuff.Message
}

func (s *MessageSet) View() *hotstuff.View {
	return s.view
}

func (s *MessageSet) Add(msg *hotstuff.Message) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if err := s.verify(msg); err != nil {
		return err
	}

	s.msgs[msg.Address] = msg
	return nil
}

func (s *MessageSet) Values() (result []*hotstuff.Message) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, v := range s.msgs {
		result = append(result, v)
	}
	return
}

func (s *MessageSet) Size() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return len(s.msgs)
}

func (s *MessageSet) Get(addr common.Address) *hotstuff.Message {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.msgs[addr]
}

func (s *MessageSet) String() string {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	addresses := make([]string, 0, len(s.msgs))
	for _, v := range s.msgs {
		addresses = append(addresses, v.Address.Hex())
	}
	return fmt.Sprintf("[%v]", strings.Join(addresses, ", "))
}

// verify if the message comes from one of the validators
func (s *MessageSet) verify(msg *hotstuff.Message) error {
	if _, v := s.vs.GetByAddress(msg.Address); v == nil {
		return fmt.Errorf("unauthorized address")
	}

	return nil
}
