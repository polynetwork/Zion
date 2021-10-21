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
