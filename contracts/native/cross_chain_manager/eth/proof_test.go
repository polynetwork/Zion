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

package eth

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestMappingKeyAt(t *testing.T) {
	txIndex := big.NewInt(0)

	encodeToString := func(b *big.Int) string {
		if b.Uint64() == 0 {
			return "00"
		}
		return hex.EncodeToString(b.Bytes())
	}

	pos1 := encodeToString(txIndex)
	pos2 := "01"

	enc, err := MappingKeyAt(pos1, pos2)
	assert.NoError(t, err)
	t.Logf(hexutil.Encode(enc))

	raw := "1d77fbbbe01a87b553ea16f9a97eaa0b80ff7b99de9c93da44a78a522d27a383"
	hash := crypto.Keccak256(common.HexToHash(raw).Bytes())
	t.Log(hexutil.Encode(hash))
}
