package backend

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/basic/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func newTestSigner() hotstuff.Signer {
	key, _ := generatePrivateKey()
	return NewSigner(key, 3)
}

func TestSign(t *testing.T) {
	s := newTestSigner()
	data := []byte("Here is a string....")
	sig, err := s.Sign(data)
	assert.NoError(t, err, "error mismatch: have %v, want nil", err)

	//Check signature recover
	hashData := crypto.Keccak256(data)
	pubkey, _ := crypto.Ecrecover(hashData, sig)
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])
	assert.Equal(t, signer, getAddress(), "address mismatch: have %v, want %s", signer.Hex(), getAddress().Hex())
}

func TestCheckValidatorSignature(t *testing.T) {
	vset, keys := newTestValidatorSet(5)

	// 1. Positive test: sign with validator's key should succeed
	data := []byte("dummy data")
	hashData := crypto.Keccak256([]byte(data))
	for i, k := range keys {
		// Sign
		sig, err := crypto.Sign(hashData, k)
		assert.NoError(t, err, "error mismatch: have %v, want nil", err)

		// CheckValidatorSignature should succeed
		signer := NewSigner(k, 3)
		addr, err := signer.CheckSignature(vset, data, sig)
		assert.NoError(t, err, "error mismatch: have %v, want nil", err)

		val := vset.GetByIndex(uint64(i))
		assert.Equal(t, addr, val.Address(), "validator address mismatch: have %v, want %v", addr, val.Address())
	}

	// 2. Negative test: sign with any key other than validator's key should return error
	key, err := crypto.GenerateKey()
	assert.NoError(t, err, "error mismatch: have %v, want nil", err)

	// Sign
	sig, err := crypto.Sign(hashData, key)
	assert.NoError(t, err, "error mismatch: have %v, want nil", err)

	// CheckValidatorSignature should return ErrUnauthorizedAddress
	signer := NewSigner(key, byte(core.MsgTypePrepareVote))
	addr, err := signer.CheckSignature(vset, data, sig)
	assert.Equal(t, err, hotstuff.ErrUnauthorizedAddress, "error mismatch: have %v, want %v", err, hotstuff.ErrUnauthorizedAddress)

	emptyAddr := common.Address{}
	assert.Equal(t, emptyAddr, common.Address{}, "address mismatch: have %v, want %v", addr, emptyAddr)
}
