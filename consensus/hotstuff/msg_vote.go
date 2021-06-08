package hotstuff

import (
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type MsgVote struct {
	QC   *types.Header
	Seal []byte
}

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *MsgVote) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.QC, m.Seal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *MsgVote) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		QC   *types.Header
		Seal []byte
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.QC, m.Seal = msg.QC, msg.Seal
	return nil
}

func (m *MsgVote) String() string {
	return fmt.Sprintf("{MsgType: %s, ViewNum: %d}", MsgTypeVote.String(), m.QC.Number.Uint64())
}

func (s *roundState) sendVoteMsg(msg *MsgVote) error {
	leader := s.getLeader()
	payload, err := s.finalizeMessage(msg, MsgTypeVote)
	if err != nil {
		return err
	}
	return s.unicast(leader, payload)
}

func (s *roundState) handleVoteMsg(m *Message) error {
	// 1. calculate hash with leader's highQC
	msg := new(MsgVote)
	if err := m.Decode(msg); err != nil {
		return err
	}
	qc := s.qcHigh.Copy()
	hash := PrepareVoteSeal(qc.Hash())

	// 2. get signer address and validate
	val, err := GetSignatureAddress(hash, msg.Seal)
	if err != nil {
		return err
	}
	if !s.isValidator(val) {
		return errInvalidSigner
	}

	// 3. collect vote msg and remove duplicate votes
	if _, ok := s.curRnd.votes[val]; ok {
		return errDuplicateVote
	}
	s.curRnd.votes[val] = msg

	// 4. check votes number
	alreadyVotes := len(s.curRnd.votes)
	major := s.snap.major()
	if alreadyVotes > major || alreadyVotes < major {
		return nil
	}

	// 5. assemble committed seal
	committedSeal := make([][]byte, major)
	for _, v := range s.curRnd.votes {
		committedSeal = append(committedSeal, v.Seal)
	}
	if err := writeCommittedSeals(qc, committedSeal); err != nil {
		return err
	}

	// 6. update
	s.Update(qc)
	return nil
}
