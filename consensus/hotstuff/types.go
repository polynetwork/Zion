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
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

// Proposal supports retrieving height and serialized block to be used during HotStuff consensus.
// It is the interface that abstracts different message structure. (consensus/hotstuff/core/core.go)
type Proposal interface {
	// Number retrieves the block height number of this proposal.
	Number() *big.Int

	NumberU64() uint64

	// Hash retrieves the hash of this proposal.
	Hash() common.Hash

	ParentHash() common.Hash

	Coinbase() common.Address

	Time() uint64

	Copy() *types.Block

	EncodeRLP(w io.Writer) error

	DecodeRLP(s *rlp.Stream) error
}

type QC interface {
	Height() *big.Int
	HeightU64() uint64
	NodeHash() common.Hash
	SealHash() common.Hash
	Proposer() common.Address
	Seal() []byte
	CommittedSeal() [][]byte
}

func RLPHash(v interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, v)
	hw.Sum(h[:0])
	return h
}
