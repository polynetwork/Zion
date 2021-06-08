package hotstuff

import (
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type MsgProposal struct {
	Proposal *types.Block
}

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *MsgProposal) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Proposal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *MsgProposal) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		Proposal *types.Block
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.Proposal = msg.Proposal
	return nil
}

func (m *MsgProposal) String() string {
	return fmt.Sprintf("{MsgType: %s, ViewNum: %d}", MsgTypeProposal.String(), m.Proposal.NumberU64())
}

func (s *roundState) sendPrepareMsg(block *types.Block) error {
	msg := &MsgProposal{Proposal: block}
	payload, err := s.finalizeMessage(msg, MsgTypeProposal)
	if err != nil {
		return err
	}
	return s.broadcast(payload)
}

func (s *roundState) handlePrepareMsg(msg *MsgProposal) error {
	block := msg.Proposal
	header := block.Header()

	// 1. store the block
	s.store.Put(block)

	// 2. check block height and vheight
	if block.NumberU64() <= s.vHeight {
		return fmt.Errorf("safety rule missing")
	}

	// 3. check block extend bLock
	parent := s.fetchParentHeader(header)
	if !s.extend(s.bLock.Header(), header) && (parent != nil && parent.Number.Uint64() > s.bLock.NumberU64()) {
		return fmt.Errorf("safety rule failed")
	}

	// 4. sign qc and send vote
	hash := PrepareVoteSeal(header.Hash())
	seal, err := s.Sign(hash)
	if err != nil {
		return err
	}
	vote := &MsgVote{
		QC:   header,
		Seal: seal,
	}
	return s.sendVoteMsg(vote)
}
