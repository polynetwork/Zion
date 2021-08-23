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

package message_set

import (
	"crypto/ecdsa"
	"math/big"
	"sort"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/basic/core"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/ethereum/go-ethereum/crypto"
)

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

func TestMessageSetWithNewView(t *testing.T) {
	valSet, _ := newTestValidatorSet(4)

	ms := MewMessageSet(valSet)

	view := &hotstuff.View{
		Round:  new(big.Int),
		Height: new(big.Int),
	}
	pp := &core.MsgNewView{
		View: view,
	}
	payload, err := core.Encode(pp)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	msg := &hotstuff.Message{
		Code:    core.MsgTypeNewView,
		Msg:     payload,
		Address: valSet.GetProposer().Address(),
	}

	if err = ms.Add(msg); err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	if err = ms.Add(msg); err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	if ms.Size() != 1 {
		t.Errorf("the size of message set mismatch: have %v, want 1", ms.Size())
	}
}

func TestMessageSetVote(t *testing.T) {
	valSet, _ := newTestValidatorSet(4)

	ms := MewMessageSet(valSet)

	view := &hotstuff.View{
		Round:  new(big.Int),
		Height: new(big.Int),
	}

	sub := &core.Vote{
		View:   view,
		Digest: common.HexToHash("1234567890"),
	}

	payload, err := core.Encode(sub)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	msg := &hotstuff.Message{
		Code:    core.MsgTypePrepareVote,
		Msg:     payload,
		Address: valSet.GetProposer().Address(),
	}

	if err := ms.Add(msg); err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	if err := ms.Add(msg); err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	if ms.Size() != 1 {
		t.Errorf("the size of message set mismatch: have %v, want 1", ms.Size())
	}
}
