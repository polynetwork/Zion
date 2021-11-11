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

package utils

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCrossChainID(t *testing.T) {
	paramTxHash := scom.Uint256ToBytes(big.NewInt(3))
	address := common.HexToAddress("0xcbc84f846c4afabd5a8adcef92b40c1c4448f31a")
	expect := "0x75c015c7cc2df8003a206a18f71db0cc2670515f0bf701132d38a8b4deb2ea39"

	blob := GenerateCrossChainID(address, paramTxHash)

	got := hexutil.Encode(blob)
	assert.Equal(t, expect, got)
}

func TestTunnelData(t *testing.T) {
	expect := &TunnelData{
		Caller:     common.HexToAddress("0x1232342"),
		ToContract: []byte{'1', 'a', '3'},
		Method:     []byte("unlock"),
		TxData:     []byte{'a', 'c', 'd', '5'},
	}

	payload, err := expect.Encode()
	assert.NoError(t, err)

	got := new(TunnelData)
	assert.NoError(t, got.Decode(payload))

	assert.Equal(t, expect, got)
}
