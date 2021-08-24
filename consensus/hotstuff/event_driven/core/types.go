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
	"github.com/ethereum/go-ethereum/core/types"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	EmptyHash    = common.Hash{}
	EmptyAddress = common.Address{}
)

type MsgType uint64

const (
	MsgTypeNewView MsgType = 1
	MsgTypeVote    MsgType = 2
	MsgTypeTimeout MsgType = 3
	MsgTypeCommit  MsgType = 4
)

func (m MsgType) String() string {
	switch m {
	case MsgTypeNewView:
		return "NEW_VIEW"
	case MsgTypeVote:
		return "VOTE"
	case MsgTypeTimeout:
		return "TIMEOUT"
	case MsgTypeCommit:
		return "COMMIT"
	default:
		return "UNKNOWN"
	}
}

func (m MsgType) Value() uint64 {
	return uint64(m)
}

// todo: set in start function
func init() {
	hotstuff.RegisterMsgTypeConvertHandler(func(data interface{}) hotstuff.MsgType {
		code := data.(uint64)
		return MsgType(code)
	})
}

type MsgNewView struct {
	View     *hotstuff.View
	Proposal hotstuff.Proposal
}

func (m *MsgNewView) EncodeRLP(w io.Writer) error {
	block, ok := m.Proposal.(*types.Block)
	if !ok {
		return errInvalidProposal
	}
	return rlp.Encode(w, []interface{}{m.View, block})
}

func (m *MsgNewView) DecodeRLP(s *rlp.Stream) error {
	var proposal struct {
		View     *hotstuff.View
		Proposal *types.Block
	}

	if err := s.Decode(&proposal); err != nil {
		return err
	}
	m.View, m.Proposal = proposal.View, proposal.Proposal
	return nil
}

func (m *MsgNewView) String() string {
	return fmt.Sprintf("{NewView Height: %v Round: %v Hash: %v}", m.View.Height, m.View.Round, m.Proposal.Hash())
}

type Vote struct {
	Epoch     uint64
	Hash      common.Hash
	Round     *big.Int
	StateRoot common.Hash

	ParentHash  common.Hash
	ParentRound *big.Int
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (v *Vote) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{v.Epoch, v.Hash, v.Round, v.StateRoot, v.ParentHash, v.ParentRound})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (v *Vote) DecodeRLP(s *rlp.Stream) error {
	var subject struct {
		Epoch     uint64
		Hash      common.Hash
		Round     *big.Int
		StateRoot common.Hash

		ParentHash  common.Hash
		ParentRound *big.Int
	}

	if err := s.Decode(&subject); err != nil {
		return err
	}

	v.Epoch, v.Hash, v.Round, v.StateRoot, v.ParentHash, v.ParentRound = subject.Epoch, subject.Hash, subject.Round, subject.StateRoot, subject.ParentHash, subject.ParentRound
	return nil
}

func (v *Vote) String() string {
	return fmt.Sprintf("{Epoch: %v, Hash: %v, Round: %v, ParentHash: %v, ParentRound: %v}", v.Epoch, v.Hash, v.Round, v.ParentHash, v.ParentRound)
}

type Timeout struct {
	Epoch uint64
	Round *big.Int
}

func (tm *Timeout) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{tm.Epoch, tm.Round})
}

func (tm *Timeout) DecodeRLP(s *rlp.Stream) error {
	var subject struct {
		Epoch uint64
		Round *big.Int
	}

	if err := s.Decode(&subject); err != nil {
		return err
	}

	tm.Epoch, tm.Round = subject.Epoch, subject.Round
	return nil
}

func (tm *Timeout) String() string {
	return fmt.Sprintf("{Epoch: %v, Round: %v}", tm.Epoch, tm.Round)
}

type backlogEvent struct {
	src hotstuff.Validator
	msg *hotstuff.Message
}
