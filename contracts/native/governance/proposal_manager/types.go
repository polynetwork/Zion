package proposal_manager

import "math/big"

type ProposalType uint8

const (
	UpdateGlobalConfig ProposalType = 0
	Normal             ProposalType = 1

	ProposalListLen int = 20
)

type ProposalList struct {
	ProposalList []*Proposal // sorted list
}

type Proposal struct {
	ID        uint64
	Type      ProposalType
	Content   []byte
	EndHeight *big.Int
	Stake     *big.Int
}
