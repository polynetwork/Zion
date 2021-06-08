package hotstuff

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"io"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type MsgNewView struct {
	PrepareQC *types.Header
}

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *MsgNewView) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.PrepareQC})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *MsgNewView) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		PrepareQC *types.Header
		ViewNum   uint64
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.PrepareQC = msg.PrepareQC
	return nil
}

func (m *MsgNewView) String() string {
	return fmt.Sprintf("{MsgType: %s, ViewNum: %d}", MsgTypeNewView.String(), m.PrepareQC.Number.Uint64())
}

func (s *roundState) sendNewViewMsg(qc *types.Header, leader common.Address) error {
	msg := &MsgNewView{PrepareQC: qc}
	payload, err := s.finalizeMessage(msg, MsgTypeNewView)
	if err != nil {
		return err
	}
	return s.unicast(leader, payload)
}

func (s *roundState) handleNewViewMsg(msg *MsgNewView) error {
	qc := msg.PrepareQC
	if err := s.verifyCommittedSeals(s.chain, qc, nil); err != nil {
		return err
	}

	s.curRnd.highQC = append(s.curRnd.highQC, qc)
	if len(s.curRnd.highQC) >= s.snap.major() {
		max := s.qcHigh
		for _, cert := range s.curRnd.highQC {
			if cert.Number.Uint64() > max.Number.Uint64() {
				max = cert
			}
		}
		s.pace.UpdateQCHigh(max)
	}
	return nil
}