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

package native

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/stretchr/testify/assert"
)

func TestAddNotify(t *testing.T) {
	// event CrossChainEvent(address indexed sender, bytes txId, address proxyOrAssetContract);
	abiJsonStr := `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"bytes","name":"txId","type":"bytes"},{"indexed":false,"internalType":"address","name":"proxyOrAssetContract","type":"address"}],"name":"CrossChainEvent","type":"event"}]`
	topic := "CrossChainEvent"
	sender := common.HexToHash("0x123")
	txId := []byte{'1', 'a'}
	proxy := common.HexToAddress("0x3a")

	ab, _ := abi.JSON(strings.NewReader(abiJsonStr))
	db := rawdb.NewMemoryDatabase()
	sdb, _ := state.New(common.Hash{}, state.NewDatabase(db), nil)
	ctx := NewNativeContract(sdb, nil)
	ref := NewContractRef(sdb, common.Address{}, common.Address{}, big.NewInt(1), common.Hash{}, 0, nil)
	ref.PushContext(&Context{
		Caller:          common.Address{},
		ContractAddress: common.Address{},
		Payload:         nil,
	})
	ctx.ref = ref

	assert.NoError(t, ctx.AddNotify(&ab, []string{topic}, sender, txId, proxy))
}
