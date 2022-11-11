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
	//var (
	//	buf bytes.Buffer
	//)
	//buf.Write(hash.Bytes())
	//buf.Write(s.commitSigSalt)
	//return buf.Bytes()
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
	//// 1. Keccak data
	//hashData := crypto.Keccak256(data)
	//// 2. Recover public key
	pubkey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}

//func (s *SignerImpl) SignData(data []byte) ([]byte, error) {
//	if data == nil {
//		return nil, ErrInvalidRawData
//	}
//	if s.privateKey == nil {
//		return nil, ErrInvalidSigner
//	}
//	hashData := crypto.Keccak256(data)
//	return crypto.Sign(hashData, s.privateKey)
//}
// todo(fuk): delete after test
//// SignerSeal proposer sign the header hash and fill extra seal with signature.
//func (s *SignerImpl) SealBeforeCommit(h *types.Header) error {
//	if h == nil {
//		return ErrInvalidHeader
//	}
//
//	sigHash := h.Hash() //s.SigHash(h)
//	seal, err := s.SignHash(sigHash)
//	if err != nil {
//		return ErrInvalidSignature
//	}
//
//	if len(seal)%types.HotstuffExtraSeal != 0 {
//		return ErrInvalidSignature
//	}
//
//	return h.SetSeal(seal)
//}
//
//// SealAfterCommit writes the extra-data field of a block header with given committed seals.
//func (s *SignerImpl) SealAfterCommit(h *types.Header, committedSeals [][]byte) error {
//	if h == nil {
//		return ErrInvalidHeader
//	}
//	if len(committedSeals) == 0 {
//		return ErrInvalidCommittedSeals
//	}
//
//	for _, seal := range committedSeals {
//		if len(seal) != types.HotstuffExtraSeal {
//			return ErrInvalidCommittedSeals
//		}
//	}
//
//	return h.SetCommittedSeal(committedSeals)
//}
//
//// SigHash returns the hash which is used as input for the Hotstuff
//// signing. It is the hash of the entire header apart from the 65 byte signature
//// contained at the end of the extra data.
////
//// Note, the method requires the extra data to be at least 65 bytes, otherwise it
//// panics. This is done to avoid accidentally using both forms (signature present
//// or not), which could be abused to produce different hashes for the same header.
//func (s *SignerImpl) SigHash(header *types.Header) (hash common.Hash) {
//	hasher := sha3.NewLegacyKeccak256()
//
//	// Clean seal is required for calculating proposer seal.
//	rlp.Encode(hasher, types.HotstuffFilteredHeader(header))
//	hasher.Sum(hash[:0])
//	return hash
//}
//func (s *SignerImpl) VerifyHash(valSet hotstuff.ValidatorSet, hash common.Hash, sig []byte) error {
//	if valSet == nil {
//		return ErrInvalidValset
//	}
//	if hash == common.EmptyHash {
//		return ErrInvalidRawHash
//	}
//	if sig == nil {
//		return ErrInvalidSignature
//	}
//
//	data := s.wrapCommittedHash(hash)
//	signer, err := getSignatureAddress(data, sig)
//	if err != nil {
//		return err
//	}
//
//	if _, val := valSet.GetByAddress(signer); val == nil {
//		return ErrUnauthorizedAddress
//	}
//
//	return nil
//}

//func (s *SignerImpl) VerifyCommittedSeal(valSet hotstuff.ValidatorSet, hash common.Hash, committedSeals [][]byte) error {
//	if valSet == nil {
//		return ErrInvalidValset
//	}
//	if hash == common.EmptyHash {
//		return ErrInvalidRawHash
//	}
//	if committedSeals == nil {
//		return ErrInvalidCommittedSeals
//	}
//
//	signers, err := s.getSignersFromCommittedSeals(hash, committedSeals)
//	if err != nil {
//		return err
//	}
//	return checkValidatorQuorum(signers, valSet)
//}
//
//func (s *SignerImpl) CheckQCParticipant(qc hotstuff.QC, signer common.Address) error {
//	if qc == nil {
//		return ErrInvalidQC
//	}
//
//	if qc.HeightU64() == 0 {
//		return nil
//	}
//	// todo(fuk): should not be deserialize
//	//extra, err := types.ExtractHotstuffExtraPayload(qc.Extra())
//	//if err != nil {
//	//	return err
//	//}
//
//	// check proposer signature
//	sealHash := s.wrapCommittedHash(qc.Hash())
//	proposer, err := getSignatureAddress(sealHash, qc.Seal())
//	if err != nil {
//		return err
//	}
//	if signer == qc.Proposer() && signer == proposer {
//		return nil
//	}
//
//	// check committed seals
//	committers, err := s.getSignersFromCommittedSeals(qc.Hash(), qc.CommittedSeal())
//	if err != nil {
//		return err
//	}
//	for _, committer := range committers {
//		if signer == committer {
//			return nil
//		}
//	}
//	return fmt.Errorf("address %s is not proposer or committer", signer.Hex())
//}
