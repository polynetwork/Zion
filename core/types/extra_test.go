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

package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

func TestExtraMissingField(t *testing.T) {
	var expect = &HotstuffExtra{
		Validators: []common.Address{},
		Seal:       []byte("111"),
		CommittedSeal: [][]byte{
			[]byte(""),
			[]byte("12"),
			[]byte("13"),
		},
		Salt: []byte{},
	}

	enc, err := rlp.EncodeToBytes(expect)
	assert.NoError(t, err)

	var got *HotstuffExtra
	assert.NoError(t, rlp.DecodeBytes(enc, &got))
	assert.Equal(t, expect, got)
}

// go test -v github.com/ethereum/go-ethereum/core/types -run TestSimple
func TestSimple(t *testing.T) {
	extraData, err := GenerateExtraWithSignature(0, 1, nil, []byte{}, [][]byte{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%x", extraData)

	list := []string{
		"0x0000000000000000000000000000000000000000000000000000000000000000f85e8083061a80f85494258af48e28e4a6846e931ddff8e1cdf8579821e5946a708455c8777630aac9d1e7702d13f7a865b27c948c09d936a1b408d6e0afaa537ba4e06c4504a0ae94ad3bf5ed640cc72f37bd21d64a65c3c756e9c88c80c080",
		"0x0000000000000000000000000000000000000000000000000000000000000000c68080c080c080",
		"0x0000000000000000000000000000000000000000000000000000000000000000f85e83061a8080f85494258af48e28e4a6846e931ddff8e1cdf8579821e5946a708455c8777630aac9d1e7702d13f7a865b27c948c09d936a1b408d6e0afaa537ba4e06c4504a0ae94ad3bf5ed640cc72f37bd21d64a65c3c756e9c88c80c080",
	}

	for _, v := range list {
		raw, err := hexutil.Decode(v)
		if err != nil {
			t.Error(err)
		}
		extra, err := ExtractHotstuffExtraPayload(raw)
		if err != nil {
			t.Error(err)
		}
		t.Logf("start %v, end %v, valset %v", extra.StartHeight, extra.EndHeight, extra.Validators)
	}
}
