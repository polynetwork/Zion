/*
 * Copyright (C) 2022 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package signature_manager

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

func TestSigInfoRLP(t *testing.T) {
	si := &SigInfo{Status: true, SigInfo: []Signature{
		{Addr: common.HexToAddress("0x01").Hex(), Content: []byte("sig")},
		{Addr: common.HexToAddress("0x02").Hex(), Content: []byte("sig2")},
		{Addr: common.HexToAddress("0x03").Hex(), Content: []byte("sig3")},
	}}
	si.init(true)

	result, err := rlp.EncodeToBytes(si)
	assert.Nil(t, err, "EncodeToBytes ng")

	decode := &SigInfo{}
	err = rlp.DecodeBytes(result, decode)
	assert.Nil(t, err, "DecodeBytes ng")

	assert.True(t, si.Status == decode.Status && reflect.DeepEqual(si.SigInfo, decode.SigInfo))
}
