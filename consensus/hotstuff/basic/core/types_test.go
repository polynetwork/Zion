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

package core

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestProposal(t *testing.T) {
	block := makeBlock(1)
	payload, err := Encode(block)
	assert.NoError(t, err)

	addr := makeAddress(1)
	msg := &message{
		Code:    MsgTypeNewView,
		Msg:     payload,
		Address: addr,
	}
	msgPayload, err := msg.Payload()
	assert.NoError(t, err)

	var decodedProposal *types.Block
	decodedMsg := new(message)
	assert.NoError(t, decodedMsg.FromPayload(msgPayload, nil))
	assert.NoError(t, decodedMsg.Decode(&decodedProposal))
	assert.Equal(t, block.Hash(), decodedProposal.Hash())

	//t.Logf("expect block %s, got %s", block.Hash().Hex(), decodedProposal.Hash().Hex())
}

func TestQuorumCert(t *testing.T) {
	qc := &hotstuff.QuorumCert{
		View: &hotstuff.View{
			Round:  big.NewInt(0),
			Height: big.NewInt(1),
		},
		Hash: makeBlock(1).Hash(),
	}
	payload, err := Encode(qc)
	assert.NoError(t, err)

	msg := &message{
		Code:    MsgTypeNewView,
		Msg:     payload,
		Address: makeAddress(1),
	}
	msgPayload, err := msg.Payload()
	assert.NoError(t, err)

	var decodedQC *hotstuff.QuorumCert
	decodedMsg := new(message)
	assert.NoError(t, decodedMsg.FromPayload(msgPayload, nil))
	assert.NoError(t, decodedMsg.Decode(&decodedQC))
	assert.Equal(t, qc.Hash, decodedQC.Hash)
}

func TestNewView(t *testing.T) {
	pp := &MsgNewView{
		View: &hotstuff.View{
			Round:  big.NewInt(1),
			Height: big.NewInt(2),
		},
		PrepareQC: &hotstuff.QuorumCert{
			View: &hotstuff.View{
				Round:  big.NewInt(0),
				Height: big.NewInt(1),
			},
			Hash: makeBlock(0).Hash(),
		},
	}
	payload, err := Encode(pp)
	assert.NoError(t, err)

	addr := makeAddress(1)
	m := &message{
		Code:    MsgTypeNewView,
		Msg:     payload,
		Address: addr,
	}

	msgPayload, err := m.Payload()
	assert.NoError(t, err)

	decodedMsg := new(message)
	if err = decodedMsg.FromPayload(msgPayload, nil); err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	var decodedPP *MsgNewView
	if err = decodedMsg.Decode(&decodedPP); err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// if block is encoded/decoded by rlp, we cannot to compare interface data type using reflect.DeepEqual. (like istanbul.Proposal)
	// so individual comparison here.
	if !reflect.DeepEqual(pp.PrepareQC.Hash, decodedPP.PrepareQC.Hash) {
		t.Errorf("proposal hash mismatch: have %v, want %v", decodedPP.PrepareQC.Hash, pp.PrepareQC.Hash)
	}

	if !reflect.DeepEqual(pp.View, decodedPP.View) {
		t.Errorf("view mismatch: have %v, want %v", decodedPP.View, pp.View)
	}
}

func TestQuorumCertWithSig(t *testing.T) {
	s := &hotstuff.QuorumCert{
		View: &hotstuff.View{
			Round:  big.NewInt(1),
			Height: big.NewInt(2),
		},
		Hash: common.HexToHash("1234567890"),
	}
	expectedSig := []byte{0x01}

	subjectPayload, _ := Encode(s)
	// 1. Encode test
	address := common.HexToAddress("0x1234567890")
	m := &message{
		Code:          MsgTypePrepareVote,
		Msg:           subjectPayload,
		Address:       address,
		Signature:     expectedSig,
		CommittedSeal: []byte{},
	}

	msgPayload, err := m.Payload()
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// 2. Decode test
	// 2.1 Test normal validate func
	decodedMsg := new(message)
	err = decodedMsg.FromPayload(msgPayload, func(data []byte, sig []byte) (common.Address, error) {
		return address, nil
	})
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	if !reflect.DeepEqual(decodedMsg, m) {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// 2.2 Test nil validate func
	decodedMsg = new(message)
	err = decodedMsg.FromPayload(msgPayload, nil)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(decodedMsg, m) {
		t.Errorf("message mismatch: have %v, want %v", decodedMsg, m)
	}

	// 2.3 Test failed validate func
	decodedMsg = new(message)
	err = decodedMsg.FromPayload(msgPayload, func(data []byte, sig []byte) (common.Address, error) {
		return common.Address{}, ErrUnauthorizedAddress
	})
	if err != ErrUnauthorizedAddress {
		t.Errorf("error mismatch: have %v, want %v", err, ErrUnauthorizedAddress)
	}
}
