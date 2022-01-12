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
	assert.Equal(t, err, errUnauthorizedAddress, "error mismatch: have %v, want %v", err, errUnauthorizedAddress)

	emptyAddr := common.Address{}
	assert.Equal(t, emptyAddr, common.Address{}, "address mismatch: have %v, want %v", addr, emptyAddr)
}

func TestFillExtraAfterCommit(t *testing.T) {
	istRawData := hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000f90197c0b841964b5c6cd21696b9ce7dbc3d7d38b07ac50b341eb18f37c7fbf6d5a91e77c0e0138f471f77b5d1046d963f32ab9a7a2c9e28cc7bda97919421d3dc0b78f6267401f9014fb841d1c3f7987e8b114c7a8a13c0daae880a601f2d66d905f75a3964e9ebb8f353726108715e7639dd75ef6881a54c8d711ac75a876be1515d8805d8e99fdc83934701b84136fbe9504ae80b1dd591a2bbf8d8410bb66da71c19fab5d619cd68ce2f7d1ab66eae46da15c971177e8c0a856b8072775ad536f05522b7a9e056d8ed4cf5402200b841964b5c6cd21696b9ce7dbc3d7d38b07ac50b341eb18f37c7fbf6d5a91e77c0e0138f471f77b5d1046d963f32ab9a7a2c9e28cc7bda97919421d3dc0b78f6267401b8416c07659d165179a3b96d31bed33f3331e11188c12ffdafbcd8f7583b3aaa5ee25dce4ca928cef9d2a348f25d3a53715ea91b1f8d5dc0ed6956a2d7df6756f02f00b84196e6c7c76c28e72daa88ed2d379f220350bde45a2a946a9993692ccd9f7065f7618938d9f18f12e32041e0e5f393947609a278acda17e256b232a3e3270f60fe0180")
	extra, _ := types.ExtractHotstuffExtraPayload(istRawData)

	expectedCommittedSeal := append([]byte{1, 2, 3}, bytes.Repeat([]byte{0x00}, types.HotstuffExtraSeal-3)...)
	expectedIstExtra := &types.HotstuffExtra{
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
	assert.Equal(t, errInvalidCommittedSeals, emptySigner.SealAfterCommit(h, [][]byte{unexpectedCommittedSeal}))
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
