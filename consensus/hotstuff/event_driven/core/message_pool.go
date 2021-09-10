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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/message_set"
)

type MessagePool struct {
	valSet hotstuff.ValidatorSet

	vote    map[common.Hash]*message_set.MessageSet
	timeout map[uint64]*message_set.MessageSet
}

func NewMessagePool(valSet hotstuff.ValidatorSet) *MessagePool {
	return &MessagePool{
		valSet:  valSet.Copy(),
		vote:    make(map[common.Hash]*message_set.MessageSet),
		timeout: make(map[uint64]*message_set.MessageSet),
	}
}

func (mp *MessagePool) AddVote(hash common.Hash, msg *hotstuff.Message) error {
	if _, ok := mp.vote[hash]; !ok {
		mp.vote[hash] = message_set.NewMessageSet(mp.valSet)
		return mp.vote[hash].Add(msg)
	}

	if mp.vote[hash].Get(msg.Address) != nil {
		return fmt.Errorf("dumplicate message")
	}

	return mp.vote[hash].Add(msg)
}

func (mp *MessagePool) VoteSize(hash common.Hash) int {
	if set, ok := mp.vote[hash]; !ok {
		return 0
	} else {
		return set.Size()
	}
}

func (mp *MessagePool) Votes(hash common.Hash) []*hotstuff.Message {
	if set, ok := mp.vote[hash]; !ok {
		return nil
	} else {
		return set.Values()
	}
}

func (mp *MessagePool) AddTimeout(round uint64, msg *hotstuff.Message) error {
	if _, ok := mp.timeout[round]; !ok {
		mp.timeout[round] = message_set.NewMessageSet(mp.valSet)
		return mp.timeout[round].Add(msg)
	}

	if mp.timeout[round].Get(msg.Address) != nil {
		return fmt.Errorf("dumplicate message")
	}

	return mp.timeout[round].Add(msg)
}

func (mp *MessagePool) TimeoutSize(round uint64) int {
	if set, ok := mp.timeout[round]; !ok {
		return 0
	} else {
		return set.Size()
	}
}

func (mp *MessagePool) Timeouts(round uint64) []*hotstuff.Message {
	if set, ok := mp.timeout[round]; !ok {
		return nil
	} else {
		return set.Values()
	}
}
