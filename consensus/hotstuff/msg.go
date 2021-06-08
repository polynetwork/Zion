package hotstuff

import (
	"bytes"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	P2PNewBlockMsg = 0x07
	P2PHotstuffMsg = 0x11
)

type MsgType uint8

const (
	MsgTypeNewView  MsgType = 1
	MsgTypeProposal MsgType = 2
	MsgTypeVote     MsgType = 3
	MsgTypeCommit   MsgType = 4
	//MsgTypePrepare       MsgType = 2
	//MsgTypePrepareVote   MsgType = 3
	//MsgTypePreCommit     MsgType = 4
	//MsgTypePreCommitVote MsgType = 5
	//MsgTypeCommit        MsgType = 6
	//MsgTypeCommitVote    MsgType = 7
	//MsgTypeDecide        MsgType = 8
)

func (m MsgType) String() string {
	switch m {
	case MsgTypeNewView:
		return "NEWVIEW"
	case MsgTypeProposal:
		return "PROPOSAL"
	case MsgTypeVote:
		return "VOTE"
	case MsgTypeCommit:
		return "COMMIT"
	//case MsgTypePreCommit:
	//	return "PRECOMMIT"
	//case MsgTypePreCommitVote:
	//	return "PRECOMMIT_VOTE"
	//case MsgTypeCommitVote:
	//	return "COMMIT_VOTE"
	//case MsgTypeDecide:
	//	return "DECIDE"
	default:
		panic("unknown msg type")
	}
}

type Message struct {
	Code      MsgType
	Msg       []byte
	Address   common.Address
	Signature []byte
}

// ==============================================
//
// define the functions that needs to be provided for rlp Encoder/Decoder.

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *Message) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Code, m.Msg, m.Address, m.Signature})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *Message) DecodeRLP(s *rlp.Stream) error {
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

func (m *Message) FromPayload(b []byte, validateFn func([]byte, []byte) (common.Address, error)) error {
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

func (m *Message) Payload() ([]byte, error) {
	return rlp.EncodeToBytes(m)
}

func (m *Message) PayloadNoSig() ([]byte, error) {
	return rlp.EncodeToBytes(&Message{
		Code:      m.Code,
		Msg:       m.Msg,
		Address:   m.Address,
		Signature: []byte{},
	})
}

func (m *Message) Decode(val interface{}) error {
	return rlp.DecodeBytes(m.Msg, val)
}

func (m *Message) String() string {
	return fmt.Sprintf("{MsgType: %d, Address: %s}", m.Code.String(), m.Address.Hex())
}
