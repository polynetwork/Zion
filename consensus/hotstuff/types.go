// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package hotstuff

import (
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

// Proposal supports retrieving height and serialized block to be used during HotStuff consensus.
// It is the interface that abstracts different message structure. (consensus/hotstuff/core/core.go)
type Proposal interface {
	// Number retrieves the block height number of this proposal.
	Number() *big.Int

	// Hash retrieves the hash of this proposal.
	Hash() common.Hash

	Coinbase() common.Address

	Time() uint64

	EncodeRLP(w io.Writer) error

	DecodeRLP(s *rlp.Stream) error
}

type Request struct {
	Proposal Proposal
}

// View includes a round number and a block height number.
// Height is the block height number we'd like to commit.
//
// If the given block is not accepted by validators, a round change will occur
// and the validators start a new round with round+1.
//
type View struct {
	Round  *big.Int
	Height *big.Int
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (v *View) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{v.Round, v.Height})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (v *View) DecodeRLP(s *rlp.Stream) error {
	var view struct {
		Round  *big.Int
		Height *big.Int
	}

	if err := s.Decode(&view); err != nil {
		return err
	}
	v.Round, v.Height = view.Round, view.Height
	return nil
}

func (v *View) String() string {
	return fmt.Sprintf("{Round: %d, Height: %d}", v.Round.Uint64(), v.Height.Uint64())
}

// Cmp compares v and y and returns:
//   -1 if v <  y
//    0 if v == y
//   +1 if v >  y
func (v *View) Cmp(y *View) int {
	if v.Height.Cmp(y.Height) != 0 {
		return v.Height.Cmp(y.Height)
	}
	if v.Round.Cmp(y.Round) != 0 {
		return v.Round.Cmp(y.Round)
	}
	return 0
}

func (v *View) Sub(y *View) (int64, int64) {
	h := new(big.Int).Sub(v.Height, y.Height).Int64()
	r := new(big.Int).Sub(v.Round, y.Round).Int64()
	return h, r
}

type QuorumCert struct {
	View     *View
	Hash     common.Hash // block header sig hash
	Proposer common.Address
	Extra    []byte
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (qc *QuorumCert) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{qc.View, qc.Hash, qc.Proposer, qc.Extra})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (qc *QuorumCert) DecodeRLP(s *rlp.Stream) error {
	var cert struct {
		View     *View
		Hash     common.Hash
		Proposer common.Address
		Extra    []byte
	}

	if err := s.Decode(&cert); err != nil {
		return err
	}
	qc.View, qc.Hash, qc.Proposer, qc.Extra = cert.View, cert.Hash, cert.Proposer, cert.Extra
	return nil
}

func (qc *QuorumCert) String() string {
	return fmt.Sprintf("{QuorumCert View: %v, Hash: %v, Proposer: %v}", qc.View, qc.Hash.String(), qc.Proposer.Hex())
}

func (qc *QuorumCert) Copy() *QuorumCert {
	enc, err := rlp.EncodeToBytes(qc)
	if err != nil {
		return nil
	}
	newQC := new(QuorumCert)
	if err := rlp.DecodeBytes(enc, &newQC); err != nil {
		return nil
	}
	return newQC
}

func RLPHash(v interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, v)
	hw.Sum(h[:0])
	return h
}
