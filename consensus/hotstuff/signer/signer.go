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

package signer

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	lru "github.com/hashicorp/golang-lru"
)

const (
	// todo(fuk): observe the sync procedure in light mode
	inmemorySignatures = 16 // Number of recent block signatures to keep in memory
)

type SignatureCache struct {
	Address common.Address
	Extra   *types.HotstuffExtra
}

type SignerImpl struct {
	address       common.Address
	privateKey    *ecdsa.PrivateKey
	signatures    *lru.ARCCache // Signatures of recent blocks to speed up mining
	commitSigSalt []byte
}

func NewSigner(privateKey *ecdsa.PrivateKey) hotstuff.Signer {
	signatures, _ := lru.NewARC(inmemorySignatures)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return &SignerImpl{
		address:       address,
		privateKey:    privateKey,
		signatures:    signatures,
		commitSigSalt: []byte("commit"),
	}
}

func (s *SignerImpl) Address() common.Address {
	return s.address
}

func (s *SignerImpl) SignHash(hash common.Hash, seal bool) ([]byte, error) {
	if hash == common.EmptyHash {
		return nil, ErrInvalidRawHash
	}
	if s.privateKey == nil {
		return nil, ErrInvalidSigner
	}

	raw := hash.Bytes()
	if seal {
		wrapHash := s.wrapCommittedHash(hash)
		raw = wrapHash.Bytes()
	}
	return crypto.Sign(raw, s.privateKey)
}

func (s *SignerImpl) SignTx(tx *types.Transaction, signer types.Signer) (*types.Transaction, error) {
	if tx == nil {
		return nil, ErrInvalidRawData
	}
	if signer == nil {
		return nil, ErrInvalidSigner
	}
	return types.SignTx(tx, signer, s.privateKey)
}

func (s *SignerImpl) CheckSignature(valSet hotstuff.ValidatorSet, hash common.Hash, sig []byte) (common.Address, error) {
	if valSet == nil {
		return common.EmptyAddress, ErrInvalidValset
	}
	if hash == common.EmptyHash {
		return common.EmptyAddress, ErrInvalidRawData
	}
	if sig == nil {
		return common.EmptyAddress, ErrInvalidSignature
	}

	// 1. Get signature address
	signer, err := getSignatureAddress(hash, sig)
	if err != nil {
		return common.Address{}, err
	}

	// 2. Check validator
	if _, val := valSet.GetByAddress(signer); val != nil {
		return val.Address(), nil
	}

	return common.Address{}, ErrUnauthorizedAddress
}

// Recover extracts the proposer address from a signed header.
func (s *SignerImpl) Recover(header *types.Header) (common.Address, *types.HotstuffExtra, error) {
	if header == nil {
		return common.EmptyAddress, nil, ErrInvalidHeader
	}

	hash := header.Hash()
	if hash == common.EmptyHash {
		return common.EmptyAddress, nil, ErrInvalidHeader
	}

	if s.signatures != nil {
		if data, ok := s.signatures.Get(hash); ok {
			if cache, ok := data.(*SignatureCache); ok {
				if cache.Extra != nil && cache.Extra.CommittedSeal != nil && len(cache.Extra.CommittedSeal) > 0 {
					return cache.Address, cache.Extra, nil
				}
			}
		}
	}

	// Retrieve the signature from the header extra-data
	extra, err := types.ExtractHotstuffExtra(header)
	if err != nil {
		return common.EmptyAddress, nil, ErrInvalidExtraDataFormat
	}

	payload := s.wrapCommittedHash(hash)
	addr, err := getSignatureAddress(payload, extra.Seal)
	if err != nil {
		return common.EmptyAddress, nil, err
	}

	if s.signatures != nil {
		s.signatures.Add(hash, &SignatureCache{
			Address: addr,
			Extra:   extra,
		})
	}
	return addr, extra, nil
}

func (s *SignerImpl) VerifyHeader(header *types.Header, valSet hotstuff.ValidatorSet, seal bool) (*types.HotstuffExtra, error) {
	if header == nil {
		return nil, ErrInvalidHeader
	}
	if valSet == nil {
		return nil, ErrInvalidValset
	}

	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return nil, nil
	}

	// resolve the authorization key and check against signers
	signer, extra, err := s.Recover(header)
	if err != nil {
		return nil, ErrInvalidSignature
	}
	if signer != header.Coinbase || signer == common.EmptyAddress {
		return nil, ErrInvalidSigner
	}

	// Signer should be in the validator set of previous block's extraData.
	if _, v := valSet.GetByAddress(signer); v == nil {
		return nil, ErrUnauthorized
	}

	if seal {
		// The length of Committed seals should be larger than 0
		if len(extra.CommittedSeal) == 0 {
			return nil, ErrEmptyCommittedSeals
		}

		// Check whether the committed seals are generated by parent's validators
		commiters, err := s.getSignersFromCommittedSeals(header.Hash(), extra.CommittedSeal, true)
		if err != nil {
			return nil, err
		}
		return extra, checkValidatorQuorum(commiters, valSet)
	}

	return extra, nil
}

func (s *SignerImpl) VerifyQC(qc hotstuff.QC, valSet hotstuff.ValidatorSet) error {
	if qc == nil {
		return ErrInvalidQC
	}
	if valSet == nil {
		return ErrInvalidValset
	}

	var (
		hash          = qc.SealHash()
		seal          = qc.Seal()
		committedSeal = qc.CommittedSeal()
	)

	addr, err := getSignatureAddress(hash, seal)
	if err != nil {
		return err
	}
	if addr != qc.Proposer() {
		return ErrInvalidSigner
	}
	if idx, _ := valSet.GetByAddress(addr); idx < 0 {
		return ErrInvalidSigner
	}

	// check committed seals
	committers, err := s.getSignersFromCommittedSeals(hash, committedSeal, false)
	if err != nil {
		return err
	}
	if err := checkValidatorQuorum(committers, valSet); err != nil {
		return err
	}

	return nil
}

type WrapHash struct {
	Hash common.Hash
	Salt []byte
}

// at the `commit` step in consensus.
// wrapCommittedHash returns a committed seal for the given hash
func (s *SignerImpl) wrapCommittedHash(hash common.Hash) common.Hash {
	return hotstuff.RLPHash(WrapHash{
		Hash: hash,
		Salt: s.commitSigSalt,
	})
}

func checkValidatorQuorum(committers []common.Address, valSet hotstuff.ValidatorSet) error {
	validators := valSet.Copy()
	validSeal := 0
	for _, addr := range committers {
		if validators.RemoveValidator(addr) {
			validSeal++
			continue
		}
	}

	// The length of validSeal should be larger than number of faulty node + 1
	if validSeal < validators.Q() {
		return ErrInvalidCommittedSeals
	}
	return nil
}

func (s *SignerImpl) getSignersFromCommittedSeals(hash common.Hash, seals [][]byte, salt bool) ([]common.Address, error) {
	var addrs []common.Address

	sealHash := hash
	if salt {
		sealHash = s.wrapCommittedHash(hash)
	}

	for _, seal := range seals {
		addr, err := getSignatureAddress(sealHash, seal)
		if err != nil {
			return nil, ErrInvalidSignature
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

// getSignatureAddress gets the address address from the signature
func getSignatureAddress(hash common.Hash, sig []byte) (common.Address, error) {
	pubkey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}
