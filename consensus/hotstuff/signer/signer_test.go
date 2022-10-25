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
	"bytes"
	"crypto/ecdsa"
	"sort"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/signer -run TestSign
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

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/signer -run TestCheckValidatorSignature
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
		signer := NewSigner(k)
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
	signer := NewSigner(key)
	addr, err := signer.CheckSignature(vset, data, sig)
	assert.Equal(t, err, ErrUnauthorizedAddress, "error mismatch: have %v, want %v", err, ErrUnauthorizedAddress)

	emptyAddr := common.Address{}
	assert.Equal(t, emptyAddr, common.Address{}, "address mismatch: have %v, want %v", addr, emptyAddr)
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/signer -run TestFillExtraAfterCommit
func TestFillExtraAfterCommit(t *testing.T) {
	istRawData := hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000f89d801ef85494258af48e28e4a6846e931ddff8e1cdf8579821e5946a708455c8777630aac9d1e7702d13f7a865b27c948c09d936a1b408d6e0afaa537ba4e06c4504a0ae94ad3bf5ed640cc72f37bd21d64a65c3c756e9c88cb8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c080")
	extra, _ := types.ExtractHotstuffExtraPayload(istRawData)

	expectedCommittedSeal := append([]byte{1, 2, 3}, bytes.Repeat([]byte{0x00}, types.HotstuffExtraSeal-3)...)
	expectedIstExtra := &types.HotstuffExtra{
		StartHeight:   0,
		EndHeight:     30,
		Validators:    extra.Validators,
		Seal:          extra.Seal,
		CommittedSeal: [][]byte{expectedCommittedSeal},
		Salt:          extra.Salt,
	}
	h := &types.Header{
		Extra: istRawData,
	}

	// normal case
	assert.NoError(t, emptySigner.SealAfterCommit(h, [][]byte{expectedCommittedSeal}))

	// verify istanbul extra-data
	istExtra, err := types.ExtractHotstuffExtra(h)
	assert.NoError(t, err)
	assert.Equal(t, expectedIstExtra, istExtra)

	// invalid seal
	unexpectedCommittedSeal := append(expectedCommittedSeal, make([]byte, 1)...)
	assert.Equal(t, ErrInvalidCommittedSeals, emptySigner.SealAfterCommit(h, [][]byte{unexpectedCommittedSeal}))
}

var emptySigner = &SignerImpl{}

type Keys []*ecdsa.PrivateKey

func (slice Keys) Len() int {
	return len(slice)
}

func (slice Keys) Less(i, j int) bool {
	return strings.Compare(crypto.PubkeyToAddress(slice[i].PublicKey).String(), crypto.PubkeyToAddress(slice[j].PublicKey).String()) < 0
}

func (slice Keys) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func newTestValidatorSet(n int) (hotstuff.ValidatorSet, []*ecdsa.PrivateKey) {
	// generate validators
	keys := make(Keys, n)
	addrs := make([]common.Address, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		keys[i] = privateKey
		addrs[i] = crypto.PubkeyToAddress(privateKey.PublicKey)
	}
	vset := validator.NewSet(addrs, hotstuff.RoundRobin)
	sort.Sort(keys) //Keys need to be sorted by its public key address
	return vset, keys
}

func getAddress() common.Address {
	return common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
}

func getInvalidAddress() common.Address {
	return common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
}

func generatePrivateKey() (*ecdsa.PrivateKey, error) {
	key := "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
	return crypto.HexToECDSA(key)
}

func newTestSigner() hotstuff.Signer {
	key, _ := generatePrivateKey()
	return NewSigner(key)
}
