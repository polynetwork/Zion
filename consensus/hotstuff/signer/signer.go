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
	"fmt"

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

func NewSigner(privateKey *ecdsa.PrivateKey) *SignerImpl {
	signatures, _ := lru.NewARC(inmemorySignatures)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return &SignerImpl{
		address:    address,
		privateKey: privateKey,
		signatures: signatures,
	}
}

func (s *SignerImpl) Address() common.Address {
	return s.address
}

func (s *SignerImpl) SignHash(hash common.Hash) ([]byte, error) {
	if hash == common.EmptyHash {
		return nil, ErrInvalidRawHash
	}
	if s.privateKey == nil {
		return nil, ErrInvalidSigner
	}
	return crypto.Sign(hash.Bytes(), s.privateKey)
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

	signer, err := getSignatureAddress(hash, sig)
	if err != nil {
		return common.Address{}, err
	}
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

	hash := types.SealHash(header)
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

	addr, err := getSignatureAddress(hash, extra.Seal)
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
		sealHash := types.SealHash(header)
		if err := s.VerifyCommittedSeal(valSet, sealHash, extra.CommittedSeal); err != nil {
			return extra, err
		}
	}

	return extra, nil
}

func (s *SignerImpl) VerifyQC(qc hotstuff.QC, valSet hotstuff.ValidatorSet, epoch bool) error {
	if qc == nil {
		return fmt.Errorf("qc is nil")
	}
	if valSet == nil {
		return fmt.Errorf("valset is nil")
	}

	hash := qc.SealHash()
	seal := qc.Seal()
	committedSeal := qc.CommittedSeal()
	if len(seal) < 65 || len(committedSeal) == 0 {
		return fmt.Errorf("seal length %v, committed seal lenght %v not enough", len(seal), len(committedSeal))
	}
	if hash == common.EmptyHash {
		return fmt.Errorf("seal hash is empty")
	}
	if epoch {
		hash = qc.NodeHash()
	}

	addr, err := getSignatureAddress(hash, seal)
	if err != nil {
		return err
	}
	if addr != qc.Proposer() {
		return fmt.Errorf("proposer expect %v, got %v", qc.Proposer(), addr)
	}
	if idx, _ := valSet.GetByAddress(addr); idx < 0 {
		return fmt.Errorf("proposer not in validator set")
	}

	return s.checkQuorum(valSet, hash, committedSeal)
}

func (s *SignerImpl) VerifyCommittedSeal(valset hotstuff.ValidatorSet, hash common.Hash, committedSeal [][]byte) error {
	if hash == common.EmptyHash {
		return ErrInvalidRawHash
	}
	if committedSeal == nil {
		return ErrInvalidCommittedSeals
	}

	return s.checkQuorum(valset, hash, committedSeal)
}

func (s *SignerImpl) checkQuorum(valset hotstuff.ValidatorSet, hash common.Hash, seals [][]byte) error {
	var addrs []common.Address

	for _, seal := range seals {
		addr, err := getSignatureAddress(hash, seal)
		if err != nil {
			return err
		}
		addrs = append(addrs, addr)
	}

	return valset.CheckQuorum(addrs)
}

// getSignatureAddress gets the address address from the signature
func getSignatureAddress(hash common.Hash, sig []byte) (common.Address, error) {
	if hash == common.EmptyHash {
		return common.EmptyAddress, fmt.Errorf("invalid hash")
	}
	if sig == nil {
		return common.EmptyAddress, fmt.Errorf("invalid sig")
	}

	pubkey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}
