/*
 * Copyright (C) 2021 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */

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

	// SignHash returns an signature of wrapped proposal hash which used as an vote
	SignHash(hash common.Hash) ([]byte, error)

	// SignTx sign transaction and full fill it with signature
	SignTx(tx *types.Transaction, signer types.Signer) (*types.Transaction, error)

	// Recover extracts the proposer address from a signed header.
	Recover(h *types.Header) (common.Address, *types.HotstuffExtra, error)

	// SealBeforeCommit writes the extra-data field of a block header with given seal.
	SealBeforeCommit(h *types.Header) error

	// SealAfterCommit writes the extra-data field of a block header with given committed seals.
	SealAfterCommit(h *types.Header, committedSeals [][]byte) error

	// VerifyHeader verify proposer signature and committed seals
	VerifyHeader(header *types.Header, valSet ValidatorSet, seal bool) (*types.HotstuffExtra, error)

	// VerifyQC verify quorum cert
	VerifyQC(qc *QuorumCert, valSet ValidatorSet) error

	// CheckQCParticipant return nil if `signer` is qc proposer or committer
	CheckQCParticipant(qc *QuorumCert, signer common.Address) error

	// CheckSignature extract address from signature and check if the address exist in validator set
	CheckSignature(valSet ValidatorSet, data []byte, signature []byte) (common.Address, error)

	VerifyHash(valSet ValidatorSet, hash common.Hash, sig []byte) error

	VerifyCommittedSeal(valSet ValidatorSet, hash common.Hash, committedSeals [][]byte) error
}
