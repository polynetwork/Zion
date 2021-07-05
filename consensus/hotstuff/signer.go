package hotstuff

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Signer interface {
	Address() common.Address

	// Sign generate signature
	Sign(data []byte) ([]byte, error)

	// SigHash generate header hash without signature
	SigHash(header *types.Header) (hash common.Hash)

	// SignVote returns an signature of wrapped proposal hash which used as an vote
	SignVote(proposal Proposal) ([]byte, error)

	// Recover extracts the proposer address from a signed header.
	Recover(h *types.Header) (common.Address, error)

	// PrepareExtra returns a extra-data of the given header and validators, without `Seal` and `CommittedSeal`
	PrepareExtra(header *types.Header, valSet ValidatorSet) ([]byte, error)

	// SealBeforeCommit writes the extra-data field of a block header with given seal.
	SealBeforeCommit(h *types.Header) error

	// SealAfterCommit writes the extra-data field of a block header with given committed seals.
	SealAfterCommit(h *types.Header, committedSeals [][]byte) error

	// VerifyHeader verify proposer signature and committed seals
	VerifyHeader(header *types.Header, valSet ValidatorSet, seal bool) error

	// VerifyQC verify quorum cert
	VerifyQC(qc *QuorumCert, valSet ValidatorSet) error

	// CheckSignature extract address from signature and check if the address exist in validator set
	CheckSignature(valSet ValidatorSet, data []byte, signature []byte) (common.Address, error)
}
