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
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/rlp"
)

type Vote struct {
	Hash      common.Hash
	Round     *big.Int
	StateRoot common.Hash

	ParentHash  common.Hash
	ParentRound *big.Int
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (v *Vote) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{v.Hash, v.Round, v.StateRoot, v.ParentHash, v.ParentRound})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (v *Vote) DecodeRLP(s *rlp.Stream) error {
	var subject struct {
		Hash      common.Hash
		Round     *big.Int
		StateRoot common.Hash

		ParentHash  common.Hash
		ParentRound *big.Int
	}

	if err := s.Decode(&subject); err != nil {
		return err
	}

	v.Hash, v.Round, v.StateRoot, v.ParentHash, v.ParentRound = subject.Hash, subject.Round, subject.StateRoot, subject.ParentHash, subject.ParentRound
	return nil
}

func (v *Vote) String() string {
	return fmt.Sprintf("{Hash: %v, Round: %v, ParentHash: %v, ParentRound: %v}", v.Hash, v.Round, v.ParentHash, v.ParentRound)
}

type Timeout struct {
	Hash  common.Hash
	Round *big.Int
}

func (tm *Timeout) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{tm.Hash, tm.Round})
}

func (tm *Timeout) DecodeRLP(s *rlp.Stream) error {
	var subject struct {
		Hash  common.Hash
		Round *big.Int
	}

	if err := s.Decode(&subject); err != nil {
		return err
	}

	tm.Hash, tm.Round = subject.Hash, subject.Round
	return nil
}

func (tm *Timeout) String() string {
	return fmt.Sprintf("{Hash: %v, Round: %v}", tm.Hash, tm.Round)
}

type backlogEvent struct {
	src hotstuff.Validator
	msg *hotstuff.Message
}
