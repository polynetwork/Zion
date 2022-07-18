package proposal_manager

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/proposal_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

type ProposalType uint8
type Status uint8

const (
	UpdateGlobalConfig ProposalType = 0
	Normal             ProposalType = 1

	Active   Status = 1
	UnActive Status = 0

	ProposalListLen int = 20
)

type ProposalList struct {
	ProposalList []*Proposal // sorted list
}

func (m *ProposalList) Decode(payload []byte) error {
	var data struct {
		ProposalList []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetProposalList, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.ProposalList, m)
}

type Proposal struct {
	ID        *big.Int
	Address   common.Address
	PType     ProposalType
	Content   []byte
	EndHeight *big.Int
	Stake     *big.Int
	Status    Status
}

func (m *Proposal) Decode(payload []byte) error {
	var data struct {
		Proposal []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetActiveProposal, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.Proposal, m)
}
