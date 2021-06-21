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
	StateLocked        State = 3
	StateCommitted     State = 4
)

func (s State) String() string {
	if s == StateAcceptRequest {
		return "AcceptRequest"
	} else if s == StatePrepared {
		return "Prepared"
	} else if s == StateLocked {
		return "StateLocked"
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
	View     *hotstuff.View
	Proposal hotstuff.Proposal
}

func (qc *QuorumCert) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{qc.View, qc.Proposal})
}

func (qc *QuorumCert) DecodeRLP(s *rlp.Stream) error {
	var quorumCert struct {
		View     *hotstuff.View
		Proposal hotstuff.Proposal
	}

	if err := s.Decode(&quorumCert); err != nil {
		return err
	}
	qc.View, qc.Proposal = quorumCert.View, quorumCert.Proposal
	return nil
}

func (qc *QuorumCert) String() string {
	return fmt.Sprintf("{QC Height: %d Round: %d Hash: %s}", qc.View.Height, qc.View.Round, qc.Proposal.Hash())
}

func (qc *QuorumCert) Copy() *QuorumCert {
	enc, err := Encode(qc)
	if err != nil {
		return nil
	}
	var newQC *QuorumCert
	if err := rlp.DecodeBytes(enc, &newQC); err != nil {
		return nil
	}
	return newQC
}

type MsgNewProposal struct {
	View     *hotstuff.View
	Proposal hotstuff.Proposal
	HighQC   *QuorumCert
}

func (m *MsgNewProposal) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.View, m.Proposal, m.HighQC})
}

func (m *MsgNewProposal) DecodeRLP(s *rlp.Stream) error {
	var proposal struct {
		View     *hotstuff.View
		Proposal hotstuff.Proposal
		HighQC   *QuorumCert
	}

	if err := s.Decode(&proposal); err != nil {
		return err
	}
	m.View, m.Proposal, m.HighQC = proposal.View, proposal.Proposal, proposal.HighQC
	return nil
}

func (m *MsgNewProposal) String() string {
	return fmt.Sprintf("{NewProposal Height: %d Round: %d Hash: %s}", m.View.Height, m.View.Round, m.Proposal.Hash())
}

type timeoutEvent struct{}
type backlogEvent struct {
	src hotstuff.Validator
	msg *message
}

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}
