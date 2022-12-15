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

package core

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestMessageWithoutSeal
func TestMessage(t *testing.T) {

	var (
		view = &View{
			Round:  big.NewInt(1),
			Height: big.NewInt(1),
		}
		code       = MsgTypeCommitVote
		msg        = []byte{'a', 'b'}
		key, _     = crypto.GenerateKey()
		signer     = crypto.PubkeyToAddress(key.PublicKey)
		validateFn = func(hash common.Hash, sig []byte) (addresses common.Address, e error) {
			pubkey, err := crypto.SigToPub(hash.Bytes(), sig)
			if err != nil {
				return common.Address{}, err
			}
			addr := crypto.PubkeyToAddress(*pubkey)
			t.Logf("%s verify hash %s", addr.Hex(), hash.Hex())
			return addr, nil
		}
	)

	expect := NewCleanMessage(view, code, msg)

	// generate and check message hash
	_, err := expect.PayloadNoSig()
	assert.NoError(t, err)
	assert.NotEqual(t, common.EmptyHash, expect.hash)
	t.Logf("message hash %s", expect.hash.Hex())

	{
		t.Log("-----test message with signature only-----")
		// generate signer and sign message
		sig, _ := crypto.Sign(expect.hash.Bytes(), key)
		expect.Signature = sig
		payload, err := expect.Payload()
		assert.NoError(t, err)
		t.Logf("%s sign msg %s", signer.Hex(), expect.hash.Hex())

		// verify message
		got := new(Message)
		err = got.FromPayload(signer, payload, validateFn)
		assert.NoError(t, err)
	}

	{
		t.Log("-----test message with committed seal-----")
		proposalHash := common.HexToHash("0xab12ba3")
		sig, _ := crypto.Sign(expect.hash.Bytes(), key)
		t.Logf("%s sign msg %s", signer.Hex(), expect.hash.Hex())
		seal, _ := crypto.Sign(proposalHash.Bytes(), key)
		t.Logf("%s sign proposal %s", signer.Hex(), proposalHash.Hex())
		expect.Signature = sig
		expect.CommittedSeal = seal

		payload, err := expect.Payload()
		assert.NoError(t, err)
		t.Logf("%s sign msg %s", signer.Hex(), expect.hash.Hex())

		// verify message
		got := new(Message)
		err = got.FromPayload(signer, payload, validateFn)
		assert.NoError(t, err)
		proposer, err := validateFn(proposalHash, got.CommittedSeal)
		assert.NoError(t, err)
		assert.Equal(t, signer, proposer)
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestProposal
func TestProposal(t *testing.T) {
	block := makeBlock(1)
	payload, err := Encode(block)
	assert.NoError(t, err)

	addr := singerAddress()
	view := &View{
		Round:  big.NewInt(0),
		Height: big.NewInt(1),
	}
	msg := &Message{
		address:   addr,
		View:      view,
		Code:      MsgTypeNewView,
		Msg:       payload,
		Signature: make([]byte, 0),
	}
	msgPayload, err := msg.Payload()
	assert.NoError(t, err)

	var (
		decodedProposal *types.Block
		decodedMsg      = new(Message)
	)
	assert.NoError(t, decodedMsg.FromPayload(addr, msgPayload, nil))
	assert.NoError(t, decodedMsg.Decode(&decodedProposal))
	assert.Equal(t, block.Hash(), decodedProposal.Hash())

	t.Logf("block hash %s, msg hash %s, proposal hash %s", block.Hash().Hex(), decodedMsg.hash.Hex(), decodedProposal.Hash().Hex())
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestQuorumCert
func TestQuorumCert(t *testing.T) {

	//1.missing fields
	{
		expect := &QuorumCert{
			view: &View{
				Round:  big.NewInt(0),
				Height: big.NewInt(1),
			},
		}
		payload, err := Encode(expect)
		assert.NoError(t, err)

		got := new(QuorumCert)
		assert.NoError(t, rlp.DecodeBytes(payload, got))
		assert.NotEqual(t, expect, got)
	}

	// 2. full qc
	{
		view := &View{
			Round:  big.NewInt(0),
			Height: big.NewInt(1),
		}
		block := makeBlock(1)
		hash := block.Hash()
		proposer := common.HexToAddress("0xabc")
		expect := &QuorumCert{
			view:          view,
			code:          MsgTypePrepareVote,
			node:          hash,
			proposer:      proposer,
			seal:          []byte{'a', 'b'},
			committedSeal: [][]byte{[]byte{'a', '1'}, []byte{'c', 'b'}},
		}
		payload, err := Encode(expect)
		assert.NoError(t, err)

		got := new(QuorumCert)
		assert.NoError(t, rlp.DecodeBytes(payload, got))
		assert.Equal(t, expect, got)

		t.Logf("proposer %s, view %v, hash %s, seal %s", got.proposer.Hex(), got.view, got.node.Hex(), string(got.seal))
	}

}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestVerifyCommittedSeal
func TestVerifyCommittedSeal(t *testing.T) {
	hash := common.HexToHash("0xf313870c2e28a138dfcb1d4c0150265a881374bd13870df5c216fe01ce4c631a")
	data := "0x0000000000000000000000000000000000000000000000000000000000000000f901168201868201a4c0b84172ddcc9dc9855cfc5f737fc0438d0a654f54580d2ccc8f7fda4b034683de419f50df9d4a729d2e6300e9c545c065aea32bb88699ad6bf40e002f5f1f7d93fc3b00f8c9b841fb9baab3b5b464f6d0a6bff202e6c89bdcf25487a192a31a730204c4ab7f595a1b67696d3dab895e1cbc9689a13d2248e65255b0464e1510cb936ebe4664b98901b84172ddcc9dc9855cfc5f737fc0438d0a654f54580d2ccc8f7fda4b034683de419f50df9d4a729d2e6300e9c545c065aea32bb88699ad6bf40e002f5f1f7d93fc3b00b841228f6495b05e8d4fe1e5b5173e01f5554666f319f66eb1d28198b353ff9edf8c23a6ccbe4b8d2af952b9d9c47c7c7ae1181070d8a64208491e66f44ad65ddc700180"
	raw, err := hexutil.Decode(data)
	assert.NoError(t, err)

	extra, err := types.ExtractHotstuffExtraPayload(raw)
	assert.NoError(t, err)

	validateFn := func(hash common.Hash, sig []byte) error {
		pubkey, err := crypto.SigToPub(hash.Bytes(), sig)
		if err != nil {
			return err
		}
		addr := crypto.PubkeyToAddress(*pubkey)
		t.Logf("%s verify hash %s", addr.Hex(), hash.Hex())
		return nil
	}

	assert.NoError(t, validateFn(hash, extra.Seal))
	for _, v := range extra.CommittedSeal {
		assert.NoError(t, validateFn(hash, v))
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestNewView
func TestNewView(t *testing.T) {
	pp := &QuorumCert{
		code: MsgTypePrepareVote,
		view: makeView(1, 0),
		node: makeBlock(0).Hash(),
	}
	payload, err := Encode(pp)
	assert.NoError(t, err)

	addr := singerAddress()
	m := &Message{
		View: &View{
			Round:  big.NewInt(0),
			Height: big.NewInt(1),
		},
		Code:    MsgTypeNewView,
		Msg:     payload,
		address: addr,
	}

	msgPayload, err := m.Payload()
	assert.NoError(t, err)

	decodedMsg := new(Message)
	if err = decodedMsg.FromPayload(addr, msgPayload, nil); err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	var decodedPP *QuorumCert
	if err = decodedMsg.Decode(&decodedPP); err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// if block is encoded/decoded by rlp, we cannot to compare interface data type using reflect.DeepEqual. (like istanbul.Proposal)
	// so individual comparison here.
	if !reflect.DeepEqual(pp.SealHash(), decodedPP.SealHash()) {
		t.Errorf("proposal hash mismatch: have %v, want %v", decodedPP.SealHash(), pp.SealHash())
	}

	if !reflect.DeepEqual(pp.view, decodedPP.view) {
		t.Errorf("view mismatch: have %v, want %v", decodedPP.view, pp.view)
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/core -run TestQuorumCertWithSig
func TestQuorumCertWithSig(t *testing.T) {
	s := &QuorumCert{
		code: MsgTypePrepareVote,
		view: makeView(2, 1),
		node: common.HexToHash("1234567890"),
	}
	expectedSig := []byte{0x01}

	subjectPayload, _ := Encode(s)
	// 1. Encode test
	address := common.HexToAddress("0x1234567890")
	m := &Message{
		View:          s.view,
		Code:          s.code,
		Msg:           subjectPayload,
		address:       address,
		Signature:     expectedSig,
		CommittedSeal: []byte{},
	}

	msgPayload, err := m.Payload()
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// 2. Decode test
	// 2.1 Test normal validate func
	decodedMsg := new(Message)
	if err := decodedMsg.FromPayload(address, msgPayload, func(data common.Hash, sig []byte) (common.Address, error) {
		return address, nil
	}); err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	_, err = m.PayloadNoSig()
	assert.NoError(t, err)
	assert.Equal(t, decodedMsg.hash, m.hash)

	// 2.2 Test nil validate func
	decodedMsg = new(Message)
	if err := decodedMsg.FromPayload(address, msgPayload, nil); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(decodedMsg, m) {
		t.Errorf("Message mismatch: have %v, want %v", decodedMsg, m)
	}

	// 2.3 Test failed validate func
	decodedMsg = new(Message)
	if err := decodedMsg.FromPayload(address, msgPayload, func(data common.Hash, sig []byte) (common.Address, error) {
		return common.Address{}, errInvalidSigner
	}); err != errInvalidSigner {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidSigner)
	}
}
