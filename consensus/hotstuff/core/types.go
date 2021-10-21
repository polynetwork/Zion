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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	EmptyHash    = common.Hash{}
	EmptyAddress = common.Address{}
)

type MsgType uint64

const (
	MsgTypeNewView       MsgType = 1
	MsgTypePrepare       MsgType = 2
	MsgTypePrepareVote   MsgType = 3
	MsgTypePreCommit     MsgType = 4
	MsgTypePreCommitVote MsgType = 5
	MsgTypeCommit        MsgType = 6
	MsgTypeCommitVote    MsgType = 7
	MsgTypeDecide        MsgType = 8
)

func (m MsgType) String() string {
	switch m {
	case MsgTypeNewView:
		return "NEW_VIEW"
	case MsgTypePrepare:
		return "PREPARE"
	case MsgTypePrepareVote:
		return "PREPARE_VOTE"
	case MsgTypePreCommit:
		return "PRECOMMIT"
	case MsgTypePreCommitVote:
		return "PRECOMMIT_VOTE"
	case MsgTypeCommit:
		return "COMMIT"
	case MsgTypeCommitVote:
		return "COMMIT_VOTE"
	case MsgTypeDecide:
		return "DECIDE"
	default:
		return "UNKNOWN"
	}
}

func (m MsgType) Value() uint64 {
	return uint64(m)
}

type State uint64

const (
	StateAcceptRequest State = 1
	StateHighQC        State = 2
	StatePrepared      State = 3
	StatePreCommitted  State = 4
	StateCommitted     State = 5
)

func (s State) String() string {
	if s == StateAcceptRequest {
		return "StateAcceptRequest"
	} else if s == StateHighQC {
		return "StateHighQC"
	} else if s == StatePrepared {
		return "StatePrepared"
	} else if s == StatePreCommitted {
		return "StatePreCommitted"
	} else if s == StateCommitted {
		return "Committed"
	} else {
		return "Unknown"
	}
}

// Cmp compares s and y and returns:
//   -1 if s is the previous state of y
//    0 if s and y are the same state
//   +1 if s is the next state of y
func (s State) Cmp(y State) int {
	if uint64(s) < uint64(y) {
		return -1
	}
	if uint64(s) > uint64(y) {
		return 1
	}
	return 0
}

type MsgNewView struct {
	View      *hotstuff.View
	PrepareQC *hotstuff.QuorumCert
}

func (m *MsgNewView) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.View, m.PrepareQC})
}

func (m *MsgNewView) DecodeRLP(s *rlp.Stream) error {
	var proposal struct {
		View      *hotstuff.View
		PrepareQC *hotstuff.QuorumCert
	}

	if err := s.Decode(&proposal); err != nil {
		return err
	}
	m.View, m.PrepareQC = proposal.View, proposal.PrepareQC
	return nil
}

func (m *MsgNewView) String() string {
	return fmt.Sprintf("{NewView Height: %d Round: %d}", m.View.Height, m.View.Round)
}

type MsgPrepare struct {
	View     *hotstuff.View
	Proposal hotstuff.Proposal
	HighQC   *hotstuff.QuorumCert
}

func (m *MsgPrepare) EncodeRLP(w io.Writer) error {
	block, ok := m.Proposal.(*types.Block)
	if !ok {
		return errInvalidProposal
	}
	return rlp.Encode(w, []interface{}{m.View, block, m.HighQC})
}

func (m *MsgPrepare) DecodeRLP(s *rlp.Stream) error {
	var proposal struct {
		View     *hotstuff.View
		Proposal *types.Block
		HighQC   *hotstuff.QuorumCert
	}

	if err := s.Decode(&proposal); err != nil {
		return err
	}
	m.View, m.Proposal, m.HighQC = proposal.View, proposal.Proposal, proposal.HighQC
	return nil
}

func (m *MsgPrepare) String() string {
	return fmt.Sprintf("{NewProposal Height: %d Round: %d Hash: %s}", m.View.Height, m.View.Round, m.Proposal.Hash())
}

type MsgPreCommit struct {
	View      *hotstuff.View
	Proposal  hotstuff.Proposal
	PrepareQC *hotstuff.QuorumCert
}

func (m *MsgPreCommit) EncodeRLP(w io.Writer) error {
	block, ok := m.Proposal.(*types.Block)
	if !ok {
		return errInvalidProposal
	}
	return rlp.Encode(w, []interface{}{m.View, block, m.PrepareQC})
}

func (m *MsgPreCommit) DecodeRLP(s *rlp.Stream) error {
	var proposal struct {
		View      *hotstuff.View
		Proposal  *types.Block
		PrepareQC *hotstuff.QuorumCert
	}

	if err := s.Decode(&proposal); err != nil {
		return err
	}
	m.View, m.Proposal, m.PrepareQC = proposal.View, proposal.Proposal, proposal.PrepareQC
	return nil
}

func (m *MsgPreCommit) String() string {
	return fmt.Sprintf("{MsgPreCommit Height: %d Round: %d Hash: %s}", m.View.Height, m.View.Round, m.Proposal.Hash())
}

type Vote struct {
	View   *hotstuff.View
	Digest common.Hash // Digest of s.Announce.Proposal.Hash()
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (b *Vote) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{b.View, b.Digest})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (b *Vote) DecodeRLP(s *rlp.Stream) error {
	var subject struct {
		View   *hotstuff.View
		Digest common.Hash
	}

	if err := s.Decode(&subject); err != nil {
		return err
	}
	b.View, b.Digest = subject.View, subject.Digest
	return nil
}

func (b *Vote) String() string {
	return fmt.Sprintf("{View: %v, Digest: %v}", b.View, b.Digest.String())
}

type timeoutEvent struct{}
type backlogEvent struct {
	src hotstuff.Validator
	msg *hotstuff.Message
}

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}
