package core

import (
	"bytes"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/rlp"
)

var EmptyHash = common.Hash{}

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
	MsgTypeRoundChange   MsgType = 1
	MsgTypePrepare       MsgType = 2
	MsgTypePrepareVote   MsgType = 3
	MsgTypePreCommit     MsgType = 4
	MsgTypePreCommitVote MsgType = 5
	MsgTypeCommit        MsgType = 6
	MsgTypeCommitVote    MsgType = 7
)

func (m MsgType) String() string {
	switch m {
	case MsgTypeRoundChange:
		return "ROUND_CHANGE"
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
	default:
		panic("unknown msg type")
	}
}

func (m MsgType) Value() uint64 {
	return uint64(m)
}

type State uint64

const (
	StateAcceptRequest State = 1
	StatePrepared      State = 2
	StatePreCommitted  State = 3
	StateCommitted     State = 4
)

func (s State) String() string {
	if s == StateAcceptRequest {
		return "AcceptRequest"
	} else if s == StatePrepared {
		return "Prepared"
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

type message struct {
	Code          MsgType
	Msg           []byte
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte
}

// ==============================================
//
// define the functions that needs to be provided for rlp Encoder/Decoder.

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *message) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Code, m.Msg, m.Address, m.Signature, m.CommittedSeal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *message) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		Code          MsgType
		Msg           []byte
		Address       common.Address
		Signature     []byte
		CommittedSeal []byte
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.Code, m.Msg, m.Address, m.Signature, m.CommittedSeal = msg.Code, msg.Msg, msg.Address, msg.Signature, msg.CommittedSeal
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
	return fmt.Sprintf("{MsgType: %s, Address: %s}", m.Code.String(), m.Address.Hex())
}

type QuorumCert struct {
	Type     MsgType
	Proposal hotstuff.Proposal
}

func (qc *QuorumCert) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{qc.Type, qc.Proposal})
}

func (qc *QuorumCert) DecodeRLP(s *rlp.Stream) error {
	var quorumCert struct {
		Type     MsgType
		Proposal hotstuff.Proposal
	}

	if err := s.Decode(&quorumCert); err != nil {
		return err
	}
	qc.Type, qc.Proposal = quorumCert.Type, quorumCert.Proposal
	return nil
}

func (qc *QuorumCert) String() string {
	return fmt.Sprintf("{QC Type:%s, Hash: %s}", qc.Type.String(), qc.Proposal.Hash())
}

type MsgNewView struct {
	View *hotstuff.View
	QC   *QuorumCert
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (m *MsgNewView) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.View, m.QC})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *MsgNewView) DecodeRLP(s *rlp.Stream) error {
	var newView struct {
		View *hotstuff.View
		QC   *QuorumCert
	}

	if err := s.Decode(&newView); err != nil {
		return err
	}
	m.View, m.QC = newView.View, newView.QC

	return nil
}

func (m *MsgNewView) String() string {
	return fmt.Sprintf("{MsgType: %s, Number:%d, Hash: %s}", MsgTypeRoundChange.String(), m.QC.Proposal.Number(), m.QC.Proposal.Hash())
}

type MsgPrepare struct {
	View     *hotstuff.View
	Proposal hotstuff.Proposal
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (m *MsgPrepare) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.View, m.Proposal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *MsgPrepare) DecodeRLP(s *rlp.Stream) error {
	var prepare struct {
		View     *hotstuff.View
		Proposal hotstuff.Proposal
	}

	if err := s.Decode(&prepare); err != nil {
		return err
	}
	m.View, m.Proposal = prepare.View, prepare.Proposal

	return nil
}

func (m *MsgPrepare) String() string {
	return fmt.Sprintf("{MsgType: %s, Number:%d, Hash: %s}", MsgTypePrepare.String(), m.Proposal.Number(), m.Proposal.Hash())
}

type MsgPrepareVote struct {
	View      *hotstuff.View
	BlockHash common.Hash
	Signature []byte
}

func (m *MsgPrepareVote) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.View, m.BlockHash, m.Signature})
}

func (m *MsgPrepareVote) DecodeRLP(s *rlp.Stream) error {
	var vote struct {
		View      *hotstuff.View
		BlockHash common.Hash
		Signature []byte
	}

	if err := s.Decode(&vote); err != nil {
		return err
	}
	m.View, m.BlockHash, m.Signature = vote.View, vote.BlockHash, vote.Signature

	return nil
}

func (m *MsgPrepareVote) String() string {
	return fmt.Sprintf("{MsgType: %s, Number: %d, Hash: %s}", MsgTypePrepareVote, m.View.Height, m.BlockHash)
}

type MsgPreCommit struct {
	View     *hotstuff.View
	Proposal hotstuff.Proposal
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (m *MsgPreCommit) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.View, m.Proposal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *MsgPreCommit) DecodeRLP(s *rlp.Stream) error {
	var precommit struct {
		View     *hotstuff.View
		Proposal hotstuff.Proposal
	}

	if err := s.Decode(&precommit); err != nil {
		return err
	}
	m.View, m.Proposal = precommit.View, precommit.Proposal

	return nil
}

func (m *MsgPreCommit) String() string {
	return fmt.Sprintf("{MsgType: %s, Number:%d, Hash: %s}", MsgTypePreCommit.String(), m.Proposal.Number(), m.Proposal.Hash())
}

type MsgPreCommitVote struct {
	View      *hotstuff.View
	BlockHash common.Hash
	Signature []byte
}

func (m *MsgPreCommitVote) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.View, m.BlockHash, m.Signature})
}

func (m *MsgPreCommitVote) DecodeRLP(s *rlp.Stream) error {
	var vote struct {
		View      *hotstuff.View
		BlockHash common.Hash
		Signature []byte
	}

	if err := s.Decode(&vote); err != nil {
		return err
	}
	m.View, m.BlockHash, m.Signature = vote.View, vote.BlockHash, vote.Signature

	return nil
}

func (m *MsgPreCommitVote) String() string {
	return fmt.Sprintf("{MsgType: %s, Number: %d, Hash: %s}", MsgTypePreCommitVote, m.View.Height, m.BlockHash)
}

//
type MsgCommit struct {
	View     *hotstuff.View
	Proposal hotstuff.Proposal
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (m *MsgCommit) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.View, m.Proposal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *MsgCommit) DecodeRLP(s *rlp.Stream) error {
	var commit struct {
		View     *hotstuff.View
		Proposal hotstuff.Proposal
	}

	if err := s.Decode(&commit); err != nil {
		return err
	}
	m.View, m.Proposal = commit.View, commit.Proposal

	return nil
}

func (m *MsgCommit) String() string {
	return fmt.Sprintf("{MsgType: %s, Number:%d, Hash: %s}", MsgTypeCommit.String(), m.Proposal.Number(), m.Proposal.Hash())
}

type MsgCommitVote struct {
	View      *hotstuff.View
	BlockHash common.Hash
	Signature []byte
}

func (m *MsgCommitVote) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.View, m.BlockHash, m.Signature})
}

func (m *MsgCommitVote) DecodeRLP(s *rlp.Stream) error {
	var vote struct {
		View      *hotstuff.View
		BlockHash common.Hash
		Signature []byte
	}

	if err := s.Decode(&vote); err != nil {
		return err
	}
	m.View, m.BlockHash, m.Signature = vote.View, vote.BlockHash, vote.Signature

	return nil
}

func (m *MsgCommitVote) String() string {
	return fmt.Sprintf("{MsgType: %s, Number: %d, Hash: %s}", MsgTypePreCommitVote, m.View.Height, m.BlockHash)
}

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}
