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
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/crypto/sha3"
)

const (
	inmemorySignatures = 4096 // Number of recent block signatures to keep in memory
)

type SignatureCache struct {
	Address common.Address
	Extra   *types.HotstuffExtra
}

type SignerImpl struct {
	address    common.Address
	privateKey *ecdsa.PrivateKey
	signatures *lru.ARCCache // Signatures of recent blocks to speed up mining
}

func NewSigner(privateKey *ecdsa.PrivateKey) hotstuff.Signer {
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

func (s *SignerImpl) Sign(data []byte) ([]byte, error) {
	if data == nil {
		return nil, errInvalidRawData
	}
	if s.privateKey == nil {
		return nil, errInvalidSigner
	}
	hashData := crypto.Keccak256(data)
	return crypto.Sign(hashData, s.privateKey)
}

func (s *SignerImpl) SignHash(hash common.Hash) ([]byte, error) {
	if hash == common.EmptyHash {
		return nil, errInvalidRawHash
	}
	wrapHash := s.wrapCommittedSeal(hash)
	return s.Sign(wrapHash)
}

func (s *SignerImpl) SignTx(tx *types.Transaction, signer types.Signer) (*types.Transaction, error) {
	if tx == nil {
		return nil, errInvalidRawData
	}
	if signer == nil {
		return nil, errInvalidSigner
	}
	return types.SignTx(tx, signer, s.privateKey)
}

// SigHash returns the hash which is used as input for the Hotstuff
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func (s *SignerImpl) SigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()

	// Clean seal is required for calculating proposer seal.
	rlp.Encode(hasher, types.HotstuffFilteredHeader(header))
	hasher.Sum(hash[:0])
	return hash
}

// Recover extracts the proposer address from a signed header.
func (s *SignerImpl) Recover(header *types.Header) (common.Address, *types.HotstuffExtra, error) {
	if header == nil {
		return common.EmptyAddress, nil, errInvalidHeader
	}

	hash := header.Hash()
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
		return common.EmptyAddress, nil, errInvalidExtraDataFormat
	}

	payload := s.SigHash(header).Bytes()
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

// SignerSeal proposer sign the header hash and fill extra seal with signature.
func (s *SignerImpl) SealBeforeCommit(h *types.Header) error {
	if h == nil {
		return errInvalidHeader
	}

	sigHash := s.SigHash(h)
	seal, err := s.Sign(sigHash.Bytes())
	if err != nil {
		return errInvalidSignature
	}

	if len(seal)%types.HotstuffExtraSeal != 0 {
		return errInvalidSignature
	}

	extra, err := types.ExtractHotstuffExtra(h)
	if err != nil {
		return err
	}
	extra.Seal = seal
	payload, err := rlp.EncodeToBytes(&extra)
	if err != nil {
		return err
	}
	h.Extra = append(h.Extra[:types.HotstuffExtraVanity], payload...)
	return nil
}

// SealAfterCommit writes the extra-data field of a block header with given committed seals.
func (s *SignerImpl) SealAfterCommit(h *types.Header, committedSeals [][]byte) error {
	if h == nil {
		return errInvalidHeader
	}
	if len(committedSeals) == 0 {
		return errInvalidCommittedSeals
	}

	for _, seal := range committedSeals {
		if len(seal) != types.HotstuffExtraSeal {
			return errInvalidCommittedSeals
		}
	}

	extra, err := types.ExtractHotstuffExtra(h)
	if err != nil {
		return err
	}

	extra.CommittedSeal = make([][]byte, len(committedSeals))
	copy(extra.CommittedSeal, committedSeals)

	payload, err := rlp.EncodeToBytes(&extra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.HotstuffExtraVanity], payload...)
	return nil
}

func (s *SignerImpl) VerifyHeader(header *types.Header, valSet hotstuff.ValidatorSet, seal bool) (*types.HotstuffExtra, error) {
	if header == nil {
		return nil, errInvalidHeader
	}
	if valSet == nil {
		return nil, errInvalidValset
	}

	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return nil, nil
	}

	// resolve the authorization key and check against signers
	signer, extra, err := s.Recover(header)
	if err != nil {
		return nil, err
	}
	if signer != header.Coinbase {
		return nil, errInvalidSigner
	}

	// Signer should be in the validator set of previous block's extraData.
	if _, v := valSet.GetByAddress(signer); v == nil {
		return nil, errUnauthorized
	}

	if seal {
		// The length of Committed seals should be larger than 0
		if len(extra.CommittedSeal) == 0 {
			return nil, errEmptyCommittedSeals
		}

		// Check whether the committed seals are generated by parent's validators
		committers, err := s.GetSignersFromCommittedSeals(header.Hash(), extra.CommittedSeal)
		if err != nil {
			return nil, err
		}
		return extra, checkValidatorQuorum(committers, valSet)
	}

	return extra, nil
}

func (s *SignerImpl) VerifyQC(qc *hotstuff.QuorumCert, valSet hotstuff.ValidatorSet) error {
	if qc == nil {
		return errInvalidQC
	}
	if valSet == nil {
		return errInvalidValset
	}

	if qc.View.Height.Uint64() == 0 {
		return nil
	}
	extra, err := types.ExtractHotstuffExtraPayload(qc.Extra)
	if err != nil {
		return err
	}

	// check proposer signature
	addr, err := getSignatureAddress(qc.Hash.Bytes(), extra.Seal)
	if err != nil {
		return err
	}
	if addr != qc.Proposer {
		return errInvalidSigner
	}
	if idx, _ := valSet.GetByAddress(addr); idx < 0 {
		return errInvalidSigner
	}

	// check committed seals
	committers, err := s.GetSignersFromCommittedSeals(qc.Hash, extra.CommittedSeal)
	if err != nil {
		return err
	}
	if err := checkValidatorQuorum(committers, valSet); err != nil {
		return err
	}
	return nil
}

func (s *SignerImpl) CheckQCParticipant(qc *hotstuff.QuorumCert, signer common.Address) error {
	if qc == nil {
		return errInvalidQC
	}

	if qc.View.Height.Uint64() == 0 {
		return nil
	}
	extra, err := types.ExtractHotstuffExtraPayload(qc.Extra)
	if err != nil {
		return err
	}

	// check proposer signature
	proposer, err := getSignatureAddress(qc.Hash.Bytes(), extra.Seal)
	if err != nil {
		return err
	}
	if signer == qc.Proposer && signer == proposer {
		return nil
	}

	// check committed seals
	committers, err := s.GetSignersFromCommittedSeals(qc.Hash, extra.CommittedSeal)
	if err != nil {
		return err
	}
	for _, committer := range committers {
		if signer == committer {
			return nil
		}
	}
	return fmt.Errorf("address %s is not proposer or committer", signer.Hex())
}

func (s *SignerImpl) CheckSignature(valSet hotstuff.ValidatorSet, data []byte, sig []byte) (common.Address, error) {
	if valSet == nil {
		return common.EmptyAddress, errInvalidValset
	}
	if data == nil {
		return common.EmptyAddress, errInvalidRawData
	}
	if sig == nil {
		return common.EmptyAddress, errInvalidSignature
	}

	// 1. Get signature address
	signer, err := getSignatureAddress(data, sig)
	if err != nil {
		return common.Address{}, err
	}

	// 2. Check validator
	if _, val := valSet.GetByAddress(signer); val != nil {
		return val.Address(), nil
	}

	return common.Address{}, errUnauthorizedAddress
}

func (s *SignerImpl) VerifyHash(valSet hotstuff.ValidatorSet, hash common.Hash, sig []byte) error {
	if valSet == nil {
		return errInvalidValset
	}
	if hash == common.EmptyHash {
		return errInvalidRawHash
	}
	if sig == nil {
		return errInvalidSignature
	}

	data := s.wrapCommittedSeal(hash)
	signer, err := getSignatureAddress(data, sig)
	if err != nil {
		return err
	}

	if _, val := valSet.GetByAddress(signer); val == nil {
		return errUnauthorizedAddress
	}

	return nil
}

func (s *SignerImpl) VerifyCommittedSeal(valSet hotstuff.ValidatorSet, hash common.Hash, committedSeals [][]byte) error {
	if valSet == nil {
		return errInvalidValset
	}
	if hash == common.EmptyHash {
		return errInvalidRawHash
	}
	if committedSeals == nil {
		return errInvalidCommittedSeals
	}

	signers, err := s.GetSignersFromCommittedSeals(hash, committedSeals)
	if err != nil {
		return err
	}
	return checkValidatorQuorum(signers, valSet)
}

// todo: useless of wrap committed seal. field of `commitSigSalt` used only to approve that participants sign block header hash
// at the `commit` step in consensus.
// wrapCommittedSeal returns a committed seal for the given hash
func (s *SignerImpl) wrapCommittedSeal(hash common.Hash) []byte {
	//var (
	//	buf bytes.Buffer
	//)
	//buf.Write(hash.Bytes())
	//buf.Write([]byte{s.commitSigSalt})
	//return buf.Bytes()
	return hash.Bytes()
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
		return errInvalidCommittedSeals
	}
	return nil
}

func (s *SignerImpl) GetSignersFromCommittedSeals(hash common.Hash, seals [][]byte) ([]common.Address, error) {
	var addrs []common.Address
	sealHash := s.wrapCommittedSeal(hash)

	// 1. Get committed seals from current header
	for _, seal := range seals {
		// 2. Get the original address by seal and parent block hash
		addr, err := getSignatureAddress(sealHash, seal)
		if err != nil {
			return nil, errInvalidSignature
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

// GetSignatureAddress gets the address address from the signature
func getSignatureAddress(data []byte, sig []byte) (common.Address, error) {
	// 1. Keccak data
	hashData := crypto.Keccak256(data)
	// 2. Recover public key
	pubkey, err := crypto.SigToPub(hashData, sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}
