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
	"bytes"
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type MsgType uint64

const (
	MsgTypeUnknown       MsgType = 0
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
	StateUnknown       State = 0
	StateAcceptRequest State = 1
	StateHighQC        State = 2
	StatePrepared      State = 3
	StatePreCommitted  State = 4
	StateCommitted     State = 5
)

func (s State) String() string {
	if s == StateUnknown {
		return "StateUnknown"
	} else if s == StateAcceptRequest {
		return "StateAcceptRequest"
	} else if s == StateHighQC {
		return "StateHighQC"
	} else if s == StatePrepared {
		return "StatePrepared"
	} else if s == StatePreCommitted {
		return "StatePreCommitted"
	} else if s == StateCommitted {
		return "StateCommitted"
	} else {
		return "Invalid"
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

// view includes a round number and a block height number.
// Height is the block height number we'd like to commit.
//
// If the given block is not accepted by validators, a round change will occur
// and the validators start a new round with round+1.
//
type View struct {
	Round  *big.Int
	Height *big.Int
}

// Cmp compares v and y and returns:
//   -1 if v <  y
//    0 if v == y
//   +1 if v >  y
func (v *View) Cmp(y *View) int {
	if hdiff := v.Height.Cmp(y.Height); hdiff != 0 {
		return hdiff
	}
	if rdiff := v.Round.Cmp(y.Round); rdiff != 0 {
		return rdiff
	}
	return 0
}

func (v *View) Sub(y *View) (int64, int64) {
	h := new(big.Int).Sub(v.Height, y.Height).Int64()
	r := new(big.Int).Sub(v.Round, y.Round).Int64()
	return h, r
}

func (v *View) String() string {
	return fmt.Sprintf("{Round: %d, Height: %d}", v.Round.Uint64(), v.Height.Uint64())
}

type MsgNewView struct {
	PrepareQC *QuorumCert
}

type MsgPrepare struct {
	Proposal hotstuff.Proposal
	HighQC   *QuorumCert
}

func (m *MsgPrepare) EncodeRLP(w io.Writer) error {
	block, ok := m.Proposal.(*types.Block)
	if !ok {
		return errInvalidProposal
	}
	return rlp.Encode(w, []interface{}{block, m.HighQC})
}

func (m *MsgPrepare) DecodeRLP(s *rlp.Stream) error {
	var proposal struct {
		Proposal *types.Block
		HighQC   *QuorumCert
	}

	if err := s.Decode(&proposal); err != nil {
		return err
	}
	m.Proposal, m.HighQC = proposal.Proposal, proposal.HighQC
	return nil
}

func (m *MsgPrepare) String() string {
	return fmt.Sprintf("{NewProposal Hash: %s}", m.Proposal.Hash())
}

type MsgPreCommit struct {
	Proposal  hotstuff.Proposal
	PrepareQC *QuorumCert
}

func (m *MsgPreCommit) EncodeRLP(w io.Writer) error {
	block, ok := m.Proposal.(*types.Block)
	if !ok {
		return errInvalidProposal
	}
	return rlp.Encode(w, []interface{}{block, m.PrepareQC})
}

func (m *MsgPreCommit) DecodeRLP(s *rlp.Stream) error {
	var proposal struct {
		Proposal  *types.Block
		PrepareQC *QuorumCert
	}

	if err := s.Decode(&proposal); err != nil {
		return err
	}
	m.Proposal, m.PrepareQC = proposal.Proposal, proposal.PrepareQC
	return nil
}

func (m *MsgPreCommit) String() string {
	return fmt.Sprintf("{MsgPreCommit Hash: %s}", m.Proposal.Hash())
}

type MsgDecide struct {
	Proposal hotstuff.Proposal
	CommitQC *QuorumCert
}

func (m *MsgDecide) EncodeRLP(w io.Writer) error {
	block, ok := m.Proposal.(*types.Block)
	if !ok {
		return errInvalidProposal
	}
	return rlp.Encode(w, []interface{}{block, m.CommitQC})
}

func (m *MsgDecide) DecodeRLP(s *rlp.Stream) error {
	var proposal struct {
		Proposal *types.Block
		CommitQC *QuorumCert
	}

	if err := s.Decode(&proposal); err != nil {
		return err
	}
	m.Proposal, m.CommitQC = proposal.Proposal, proposal.CommitQC
	return nil
}

type QuorumCert struct {
	view          *View
	code          MsgType
	hash          common.Hash // proposal hash but not message hash
	proposer      common.Address
	seal          []byte
	committedSeal [][]byte
}

func (qc *QuorumCert) Height() *big.Int {
	if qc.view == nil {
		return common.Big0
	}
	return qc.view.Height
}

func (qc *QuorumCert) HeightU64() uint64 {
	return qc.Height().Uint64()
}

func (qc *QuorumCert) Round() *big.Int {
	if qc.view == nil {
		return common.Big0
	}
	return qc.view.Round
}

func (qc *QuorumCert) RoundU64() uint64 {
	return qc.Round().Uint64()
}

// Hash retrieve message hash but not proposal hash
func (qc *QuorumCert) SealHash() common.Hash {
	msg := NewCleanMessage(qc.view, qc.code, qc.hash.Bytes())
	msg.PayloadNoSig()
	return msg.hash
}

func (qc *QuorumCert) Proposer() common.Address {
	return qc.proposer
}

func (qc *QuorumCert) Seal() []byte {
	return qc.seal
}

func (qc *QuorumCert) CommittedSeal() [][]byte {
	return qc.committedSeal
}

func (qc *QuorumCert) String() string {
	return fmt.Sprintf("{QuorumCert View: %v, Hash: %v, Proposer: %v}", qc.view, qc.hash, qc.proposer)
}

func (qc *QuorumCert) Copy() *QuorumCert {
	enc, err := rlp.EncodeToBytes(qc)
	if err != nil {
		return nil
	}
	newQC := new(QuorumCert)
	if err := rlp.DecodeBytes(enc, &newQC); err != nil {
		return nil
	}
	return newQC
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (qc *QuorumCert) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{qc.view, qc.code, qc.hash, qc.proposer, qc.seal, qc.committedSeal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (qc *QuorumCert) DecodeRLP(s *rlp.Stream) error {
	var cert struct {
		View          *View
		Code          MsgType
		Hash          common.Hash
		Proposer      common.Address
		Seal          []byte
		CommittedSeal [][]byte
	}

	if err := s.Decode(&cert); err != nil {
		return err
	}
	qc.view, qc.code, qc.hash, qc.proposer, qc.seal, qc.committedSeal = cert.View, cert.Code, cert.Hash, cert.Proposer, cert.Seal, cert.CommittedSeal
	return nil
}

type Message struct {
	address common.Address
	hash    common.Hash

	Code          MsgType
	View          *View
	Msg           []byte
	Signature     []byte
	CommittedSeal []byte
}

func NewCleanMessage(view *View, code MsgType, payload []byte) *Message {
	return &Message{
		View: view,
		Code: code,
		Msg:  payload,
	}
}

// ==============================================
//
// define the functions that needs to be provided for rlp Encoder/Decoder.

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *Message) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Code.Value(), m.View, m.Msg, m.Signature, m.CommittedSeal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *Message) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		Code          uint64
		View          *View
		Msg           []byte
		Signature     []byte
		CommittedSeal []byte
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}

	m.Code, m.View, m.Msg, m.Signature, m.CommittedSeal = MsgType(msg.Code), msg.View, msg.Msg, msg.Signature, msg.CommittedSeal
	return nil
}

// ==============================================
//
// define the functions that needs to be provided for core.

func (m *Message) FromPayload(src common.Address, payload []byte, validateFn func(common.Hash, []byte) (common.Address, error)) error {
	// Decode Message
	if err := rlp.DecodeBytes(payload, &m); err != nil {
		return err
	}

	// check that msg fields should NOT be nil or empty
	if m.View == nil || src == common.EmptyAddress || m.Msg == nil {
		return errInvalidMessage
	}

	// Validate Message (on a Message without Signature)
	if _, err := m.PayloadNoSig(); err != nil {
		return err
	}
	if validateFn != nil {
		signer, err := validateFn(m.hash, m.Signature)
		if err != nil {
			return err
		}
		if !bytes.Equal(signer.Bytes(), src.Bytes()) {
			return errInvalidSigner
		}
	}

	m.address = src
	return nil
}

func (m *Message) Payload() ([]byte, error) {
	return Encode(m)
}

func (m *Message) PayloadNoSig() ([]byte, error) {
	data, err := rlp.EncodeToBytes(&Message{
		Code:      m.Code,
		View:      m.View,
		Msg:       m.Msg,
		Signature: []byte{},
	})
	if err != nil {
		return nil, err
	}

	m.hash = crypto.Keccak256Hash(data)
	return data, nil
}

func (m *Message) Decode(val interface{}) error {
	return rlp.DecodeBytes(m.Msg, val)
}

func (m *Message) String() string {
	return fmt.Sprintf("{MsgType: %v, view: %v}", m.Code, m.View)
}

type timeoutEvent struct{}
type backlogEvent struct {
	src hotstuff.Validator
	msg *Message
}

type Request struct {
	Proposal hotstuff.Proposal
}

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}
