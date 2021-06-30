package backend

import (
	"bytes"
	"crypto/ecdsa"
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

type SignerImpl struct {
	address       common.Address
	privateKey    *ecdsa.PrivateKey
	signatures    *lru.ARCCache // Signatures of recent blocks to speed up mining
	commitSigSalt byte          //
}

func NewSigner(privateKey *ecdsa.PrivateKey, commitSigSalt byte) hotstuff.Signer {
	signatures, _ := lru.NewARC(inmemorySignatures)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return &SignerImpl{
		address:       address,
		privateKey:    privateKey,
		signatures:    signatures,
		commitSigSalt: commitSigSalt,
	}
}

func (s *SignerImpl) Address() common.Address {
	return s.address
}

func (s *SignerImpl) Sign(data []byte) ([]byte, error) {
	hashData := crypto.Keccak256(data)
	return crypto.Sign(hashData, s.privateKey)
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
	rlp.Encode(hasher, types.HotstuffFilteredHeader(header, false))
	hasher.Sum(hash[:0])
	return hash
}

// Recover extracts the proposer address from a signed header.
func (s *SignerImpl) Recover(header *types.Header) (common.Address, error) {
	hash := header.Hash()
	if addr, ok := s.signatures.Get(hash); ok {
		return addr.(common.Address), nil
	}

	// Retrieve the signature from the header extra-data
	extra, err := types.ExtractHotstuffExtra(header)
	if err != nil {
		return common.Address{}, err
	}

	payload := s.SigHash(header).Bytes()
	addr, err := getSignatureAddress(payload, extra.Seal)
	if err != nil {
		return addr, err
	}

	s.signatures.Add(hash, addr)
	return addr, nil
}

func (s *SignerImpl) PrepareExtra(header *types.Header, valSet hotstuff.ValidatorSet) ([]byte, error) {
	var (
		buf  bytes.Buffer
		vals = valSet.AddressList()
	)

	// compensate the lack bytes if header.Extra is not enough IstanbulExtraVanity bytes.
	if len(header.Extra) < types.HotstuffExtraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, types.HotstuffExtraVanity-len(header.Extra))...)
	}
	buf.Write(header.Extra[:types.HotstuffExtraVanity])

	ist := &types.HotstuffExtra{
		Validators:    vals,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return nil, err
	}

	return append(buf.Bytes(), payload...), nil
}

func (s *SignerImpl) FillExtraBeforeCommit(h *types.Header) error {
	// sign the hash
	seal, err := s.Sign(s.SigHash(h).Bytes())
	if err != nil {
		return err
	}

	// generate extra
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

// FillExtraAfterCommit writes the extra-data field of a block header with given committed seals.
func (s *SignerImpl) FillExtraAfterCommit(h *types.Header, committedSeals [][]byte) error {
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

func (s *SignerImpl) VerifyHeader(header *types.Header, valSet hotstuff.ValidatorSet, seal bool) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}

	// resolve the authorization key and check against signers
	signer, err := s.Recover(header)
	if err != nil {
		return err
	}

	// Signer should be in the validator set of previous block's extraData.
	if _, v := valSet.GetByAddress(signer); v == nil {
		return errUnauthorized
	}

	if seal {
		extra, err := types.ExtractHotstuffExtra(header)
		if err != nil {
			return err
		}
		// The length of Committed seals should be larger than 0
		if len(extra.CommittedSeal) == 0 {
			return errEmptyCommittedSeals
		}

		// Check whether the committed seals are generated by parent's validators
		committers, err := s.getSignersFromCommittedSeals(header.Hash(), extra.CommittedSeal)
		if err != nil {
			return err
		}
		return checkValidatorQuorum(committers, valSet)
	}

	return nil
}

func (s *SignerImpl) VerifyQC(qc *hotstuff.QuorumCert, valSet hotstuff.ValidatorSet) error {
	if qc.View.Height.Uint64() == 0 {
		return nil
	}
	extra, err := types.ExtractHotstuffExtraPayload(qc.Extra)
	if err != nil {
		return err
	}

	// check proposer signature
	addr, err := getSignatureAddress(qc.SealHash[:], extra.Seal)
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
	committers, err := s.getSignersFromCommittedSeals(qc.Hash, extra.CommittedSeal)
	if err != nil {
		return err
	}
	if err := checkValidatorQuorum(committers, valSet); err != nil {
		return err
	}
	return nil
}

func (s *SignerImpl) CheckSignature(valSet hotstuff.ValidatorSet, data []byte, sig []byte) (common.Address, error) {
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

// WrapCommittedSeal returns a committed seal for the given hash
func (s *SignerImpl) WrapCommittedSeal(hash common.Hash) []byte {
	var (
		buf bytes.Buffer
	)
	buf.Write(hash.Bytes())
	buf.Write([]byte{s.commitSigSalt})
	return buf.Bytes()
}

func checkValidatorQuorum(committers []common.Address, validators hotstuff.ValidatorSet) error {
	validSeal := 0
	for _, addr := range committers {
		if validators.RemoveValidator(addr) {
			validSeal++
			continue
		}
		return errInvalidCommittedSeals
	}

	// The length of validSeal should be larger than number of faulty node + 1
	if validSeal <= validators.Q() {
		return errInvalidCommittedSeals
	}
	return nil
}

func (s *SignerImpl) getSignersFromCommittedSeals(hash common.Hash, seals [][]byte) ([]common.Address, error) {
	var addrs []common.Address
	proposalSeal := s.WrapCommittedSeal(hash)

	// 1. Get committed seals from current header
	for _, seal := range seals {
		// 2. Get the original address by seal and parent block hash
		addr, err := getSignatureAddress(proposalSeal, seal)
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

//func (s *SignerImpl) GetSignerFromCommittedSeal(hash common.Hash, sig []byte) (common.Address, error) {
//	proposalSeal := s.WrapCommittedSeal(hash)
//	return s.GetSignatureAddress(proposalSeal, sig)
//}
//
//func (s *SignerImpl) GetSignersFromCommittedSeal(hash common.Hash, sigs [][]byte) ([]common.Address, error) {
//	var addrs []common.Address
//	proposalSeal := s.WrapCommittedSeal(hash)
//
//	// 1. Get committed seals from current header
//	for _, sig := range sigs {
//		// 2. Get the original address by seal and parent block hash
//		addr, err := s.GetSignatureAddress(proposalSeal, sig)
//		if err != nil {
//			return nil, err
//		}
//		addrs = append(addrs, addr)
//	}
//	return addrs, nil
//}
