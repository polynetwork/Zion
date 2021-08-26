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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type MsgType uint64

const (
	MsgTypeProposal MsgType = 1
	MsgTypeVote     MsgType = 2
	MsgTypeTimeout  MsgType = 3
	MsgTypeQC       MsgType = 4
)

func (m MsgType) String() string {
	switch m {
	case MsgTypeProposal:
		return "MSG_PROPOSAL"
	case MsgTypeVote:
		return "MSG_VOTE"
	case MsgTypeTimeout:
		return "MSG_TIMEOUT"
	case MsgTypeQC:
		return "MSG_QC"
	default:
		return "MSG_UNKNOWN"
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

type MsgProposal struct {
	Epoch     uint64
	View      *hotstuff.View
	Proposal  *types.Block
	JustifyQC *hotstuff.QuorumCert
}

func (m *MsgProposal) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Epoch, m.View, m.Proposal, m.JustifyQC})
}

func (m *MsgProposal) DecodeRLP(s *rlp.Stream) error {
	var subject struct {
		Epoch     uint64
		View      *hotstuff.View
		Proposal  *types.Block
		JustifyQC *hotstuff.QuorumCert
	}

	if err := s.Decode(&subject); err != nil {
		return err
	}
	m.Epoch, m.View, m.Proposal, m.JustifyQC = subject.Epoch, subject.View, subject.Proposal, subject.JustifyQC
	return nil
}

func (m *MsgProposal) String() string {
	return fmt.Sprintf("{NewView Height: %v Round: %v Hash: %v}", m.View.Height, m.View.Round, m.Proposal.Hash())
}

type Vote struct {
	Epoch          uint64
	Hash           common.Hash
	View           *hotstuff.View
	ParentHash     common.Hash
	ParentView     *hotstuff.View
	GrandHash      common.Hash
	GrandView      *hotstuff.View
	GreatGrandHash common.Hash
	GreatGrandView *hotstuff.View
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (v *Vote) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{v.Epoch, v.Hash, v.View, v.ParentHash, v.View, v.GrandHash, v.View, v.GreatGrandHash, v.GreatGrandView})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (v *Vote) DecodeRLP(s *rlp.Stream) error {
	var subject struct {
		Epoch          uint64
		Hash           common.Hash
		View           *hotstuff.View
		ParentHash     common.Hash
		ParentView     *hotstuff.View
		GrandHash      common.Hash
		GrandView      *hotstuff.View
		GreatGrandHash common.Hash
		GreatGrandView *hotstuff.View
	}

	if err := s.Decode(&subject); err != nil {
		return err
	}

	v.Epoch, v.Hash, v.View, v.ParentHash, v.ParentView,
		v.GrandHash, v.GrandView, v.GreatGrandHash, v.GreatGrandView =
		subject.Epoch, subject.Hash, subject.View, subject.ParentHash, subject.ParentView,
		subject.GrandHash, subject.GrandView, subject.GreatGrandHash, subject.GreatGrandView
	return nil
}

func (v *Vote) String() string {
	if v.GrandView == nil {
		return fmt.Sprintf("{Epoch: %v, Hash: %v, View: %v, ParentHash: %v, ParentView: %v}", v.Epoch, v.Hash, v.View, v.ParentHash, v.ParentView)
	} else if v.GreatGrandView == nil {
		return fmt.Sprintf("{Epoch: %v, Hash: %v, View: %v, ParentHash: %v, ParentView: %v, GrandHash: %v, GrandView: %v}", v.Epoch, v.Hash, v.View, v.ParentHash, v.ParentView, v.GrandHash, v.GrandView)
	} else {
		return fmt.Sprintf("{Epoch: %v, Hash: %v, View: %v, ParentHash: %v, ParentView: %v, GrandHash: %v, GrandView: %v, GreateGrandHash: %v, GreateGrandView: %v}", v.Epoch, v.Hash, v.View, v.ParentHash, v.ParentView, v.GrandHash, v.GrandView, v.GreatGrandHash, v.GreatGrandView)
	}
}

type TimeoutEvent struct {
	Epoch uint64
	Round *big.Int
}

func (tm *TimeoutEvent) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{tm.Epoch, tm.Round})
}

func (tm *TimeoutEvent) DecodeRLP(s *rlp.Stream) error {
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

func (tm *TimeoutEvent) String() string {
	return fmt.Sprintf("{Epoch: %v, Round: %v}", tm.Epoch, tm.Round)
}

type CertificateEvent struct {
	Cert *hotstuff.QuorumCert
}

func (ce *CertificateEvent) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{ce.Cert})
}

func (ce *CertificateEvent) DecodeRLP(s *rlp.Stream) error {
	var subject struct {
		Cert *hotstuff.QuorumCert
	}

	if err := s.Decode(&subject); err != nil {
		return err
	}

	ce.Cert = subject.Cert
	return nil
}

func (ce *CertificateEvent) String() string {
	return fmt.Sprintf("{Hash: %v, View: %v, Proposer: %v}", ce.Cert.Hash, ce.Cert.View, ce.Cert.Proposer)
}

type ExtraSalt struct {
	Epoch uint64
	Round *big.Int
}

func (es *ExtraSalt) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{es.Epoch, es.Round})
}

func (es *ExtraSalt) DecodeRLP(s *rlp.Stream) error {
	var subject struct {
		Epoch uint64
		Round *big.Int
	}

	if err := s.Decode(&subject); err != nil {
		return err
	}

	es.Epoch, es.Round = subject.Epoch, subject.Round
	return nil
}

func (es *ExtraSalt) String() string {
	return fmt.Sprintf("{Epoch: %v, Round: %v}", es.Epoch, es.Round)
}

type backlogEvent struct {
	src hotstuff.Validator
	msg *hotstuff.Message
}
