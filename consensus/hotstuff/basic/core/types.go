package core

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type CoreEngine interface {
	Start() error

	Stop() error

	IsProposer() bool

	// verify if a hash is the same as the proposed block in the current pending request
	//
	// this is useful when the engine is currently the speaker
	//
	// pending request is populated right at the request stage so this would give us the earliest verification
	// to avoid any race condition of coming propagated blocks
	IsCurrentProposal(blockHash common.Hash) bool

	// CurrentRoundState() *roundState
}

type MsgType uint64

const (
	MsgTypeNewView       MsgType = 1 // newView or roundChange
	MsgTypePrepare       MsgType = 2
	MsgTypePrepareVote   MsgType = 3
	MsgTypePreCommit     MsgType = 4
	MsgTypePreCommitVote MsgType = 5
	MsgTypeCommit        MsgType = 6
)

func (m MsgType) String() string {
	switch m {
	case MsgTypeNewView:
		return "NEWVIEW"
	case MsgTypePrepare:
		return "PREPARE"
	case MsgTypePrepareVote:
		return "PREPAREVOTE"
	case MsgTypePreCommit:
		return "PRECOMMIT"
	case MsgTypePreCommitVote:
		return "PRECOMMITVOTE"
	case MsgTypeCommit:
		return "COMMIT"
	default:
		panic("unknown msg type")
	}
}

type message struct {
	Code      MsgType
	Msg       []byte
	Address   common.Address
	Signature []byte
}

// ==============================================
//
// define the functions that needs to be provided for rlp Encoder/Decoder.

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *message) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Code, m.Msg, m.Address, m.Signature})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *message) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		Code      MsgType
		Msg       []byte
		Address   common.Address
		Signature []byte
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.Code, m.Msg, m.Address, m.Signature = msg.Code, msg.Msg, msg.Address, msg.Signature
	return nil
}

// ==============================================
//
// define the functions that needs to be provided for core.

func (m *message) FromPayload(b []byte, validateFn func([]byte, []byte) (common.Address, error)) error {
	// Decode message
	err := rlp.DecodeBytes(b, &m)
	if err != nil {
		return err
	}

	// Validate message (on a message without Signature)
	if validateFn != nil {
		var payload []byte
		payload, err = m.PayloadNoSig()
		if err != nil {
			return err
		}

		signerAdd, err := validateFn(payload, m.Signature)
		if err != nil {
			return err
		}
		if !bytes.Equal(signerAdd.Bytes(), m.Address.Bytes()) {
			return errInvalidSigner
		}
	}
	return nil
}

func (m *message) Payload() ([]byte, error) {
	return rlp.EncodeToBytes(m)
}

func (m *message) PayloadNoSig() ([]byte, error) {
	return rlp.EncodeToBytes(&message{
		Code:      m.Code,
		Msg:       m.Msg,
		Address:   m.Address,
		Signature: []byte{},
	})
}

func (m *message) Decode(val interface{}) error {
	return rlp.DecodeBytes(m.Msg, val)
}

func (m *message) String() string {
	return fmt.Sprintf("{MsgType: %d, Address: %s}", m.Code.String(), m.Address.Hex())
}

type State uint64

const (
	StateNewRound  State = 1 // prepare to accept new view
	StatePrepared  State = 2
	StateCommitted State = 3
	StateDecide    State = 4
)

func (s State) String() string {
	if s == StateNewRound {
		return "NewRound"
	} else if s == StatePrepared {
		return "Prepared"
	} else if s == StateCommitted {
		return "Committed"
	} else if s == StateDecide {
		return "Decide"
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

// View includes a round number and a block height number.
// Height is the block height number we'd like to commit.

// If the given block is not accepted by validators, a round change will occur
// and the validators start a new round with round+1.
type View struct {
	Round  *big.Int
	Height *big.Int
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (v *View) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{v.Round, v.Height})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (v *View) DecodeRLP(s *rlp.Stream) error {
	var view struct {
		Round  *big.Int
		Height *big.Int
	}

	if err := s.Decode(&view); err != nil {
		return err
	}
	v.Round, v.Height = view.Round, view.Height
	return nil
}

func (v *View) String() string {
	return fmt.Sprintf("{Round: %d, Height: %d}", v.Round.Uint64(), v.Height.Uint64())
}

// Cmp compares v and y and returns:
//   -1 if v <  y
//    0 if v == y
//   +1 if v >  y
func (v *View) Cmp(y *View) int {
	if v.Height.Cmp(y.Height) != 0 {
		return v.Height.Cmp(y.Height)
	}
	if v.Round.Cmp(y.Round) != 0 {
		return v.Round.Cmp(y.Round)
	}
	return 0
}

type MsgNewView struct {
	View     *View
	Proposal hotstuff.Proposal
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (m *MsgNewView) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.View, m.Proposal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *MsgNewView) DecodeRLP(s *rlp.Stream) error {
	var msgNewView struct {
		View     *View
		Proposal *types.Block
	}

	if err := s.Decode(&msgNewView); err != nil {
		return err
	}
	m.View, m.Proposal = msgNewView.View, msgNewView.Proposal

	return nil
}

type Subject struct {
	View   *View
	Digest common.Hash
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (b *Subject) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{b.View, b.Digest})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (b *Subject) DecodeRLP(s *rlp.Stream) error {
	var subject struct {
		View   *View
		Digest common.Hash
	}

	if err := s.Decode(&subject); err != nil {
		return err
	}
	b.View, b.Digest = subject.View, subject.Digest
	return nil
}

func (b *Subject) String() string {
	return fmt.Sprintf("{View: %v, Digest: %v}", b.View, b.Digest.String())
}