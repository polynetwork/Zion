package hotstuff

import (
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	P2PNewBlockMsg = 0x07
	P2PHotstuffMsg = 0x11
)

type MsgType uint8

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
		return "NEWVIEW"
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

type QuorumCert struct {
	BlockHash common.Hash
	ViewNum   uint64
	Type      MsgType
	Signature []byte
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (q *QuorumCert) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{q.BlockHash, q.ViewNum, q.Type, q.Signature})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (q *QuorumCert) DecodeRLP(s *rlp.Stream) error {
	var cert struct {
		BlockHash common.Hash
		ViewNum   uint64
		Type      MsgType
		Signature []byte
	}

	if err := s.Decode(&cert); err != nil {
		return err
	}
	q.BlockHash, q.ViewNum, q.Type, q.Signature = cert.BlockHash, cert.ViewNum, cert.Type, cert.Signature
	return nil
}

func (q *QuorumCert) String() string {
	return fmt.Sprintf("{BlockHash: %s, ViewNum: %d, MsgType: %s}", q.BlockHash.Hex(), q.ViewNum, q.Type.String())
}

type Message struct {
	Code          MsgType
	Payload       []byte
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte
}

// ==============================================
//
// define the functions that needs to be provided for rlp Encoder/Decoder.

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *Message) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Code, m.Payload, m.Address, m.Signature, m.CommittedSeal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *Message) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		Code          MsgType
		Payload       []byte
		Address       common.Address
		Signature     []byte
		CommittedSeal []byte
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.Code, m.Payload, m.Address, m.Signature, m.CommittedSeal = msg.Code, msg.Payload, msg.Address, msg.Signature, msg.CommittedSeal
	return nil
}

func (m *Message) String() string {
	return fmt.Sprintf("{MsgType: %d, Address: %s}", m.Code.String(), m.Address.Hex())
}

type MsgNewView struct {
	PrepareQC *QuorumCert
	ViewNum   uint64
}

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *MsgNewView) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.PrepareQC, m.ViewNum})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *MsgNewView) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		PrepareQC *QuorumCert
		ViewNum   uint64
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.PrepareQC, m.ViewNum = msg.PrepareQC, msg.ViewNum
	return nil
}

func (m *MsgNewView) String() string {
	return fmt.Sprintf("{MsgType: %s, ViewNum: %d}", MsgTypeNewView.String(), m.ViewNum)
}

type MsgPrepare struct {
	CurProposal *types.Block
	HighQC      *QuorumCert
	ViewNum     uint64
}

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *MsgPrepare) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.CurProposal, m.HighQC, m.ViewNum})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *MsgPrepare) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		CurProposal *types.Block
		HighQC      *QuorumCert
		ViewNum     uint64
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.CurProposal, m.HighQC, m.ViewNum = msg.CurProposal, msg.HighQC, msg.ViewNum
	return nil
}

func (m *MsgPrepare) String() string {
	return fmt.Sprintf("{MsgType: %s, ViewNum: %d}", MsgTypePrepare.String(), m.ViewNum)
}

type MsgPrepareVote struct {
	BlockHash  common.Hash
	QC         *QuorumCert
	PartialSig []byte
	ViewNum    uint64
}

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *MsgPrepareVote) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.BlockHash, m.QC, m.PartialSig, m.ViewNum})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *MsgPrepareVote) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		BlockHash  common.Hash
		QC         *QuorumCert
		PartialSig []byte
		ViewNum    uint64
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.BlockHash, m.QC, m.PartialSig, m.ViewNum = msg.BlockHash, msg.QC, msg.PartialSig, msg.ViewNum
	return nil
}

func (m *MsgPrepareVote) String() string {
	return fmt.Sprintf("{MsgType: %s, ViewNum: %d}", MsgTypePrepareVote.String(), m.ViewNum)
}

type MsgDecide struct {
}

// MessageEvent is posted for Istanbul engine communication
type MessageEvent struct {
	Payload []byte
}
