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
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
)

func NewTestStateDB() *state.StateDB {
	memdb := rawdb.NewMemoryDatabase()
	db := state.NewDatabase(memdb)
	stateDB, _ := state.New(common.Hash{}, db, nil)
	return stateDB
}

// generateTestPeer ONLY used for testing
func generateTestPeer() (common.Address, *ecdsa.PrivateKey) {
	pk, _ := crypto.GenerateKey()
	return crypto.PubkeyToAddress(pk.PublicKey), pk
}

func GenerateTestPeers(n int) ([]common.Address, []*ecdsa.PrivateKey) {
	peers := make([]common.Address, n)
	pris := make([]*ecdsa.PrivateKey, n)
	for i := 0; i < n; i++ {
		peers[i], pris[i] = generateTestPeer()
	}
	return peers, pris
}

// GenerateTestContext generate nativeContract with params, and the sequence of param contains:
// `blockHeight`, `caller`, `tx sender`, `tx hash`, `supply gas`. these params separated with
// each other by type, only the `caller` and `tx sender` are both of type of `common.Address`,
// and the first one is `sender` and the next is `caller`.
func GenerateTestContext(t *testing.T, params ...interface{}) (*state.StateDB, *NativeContract) {
	var (
		block     = int(0)
		caller    = common.EmptyAddress
		sender    = common.EmptyAddress
		hash      = common.EmptyHash
		supplyGas = uint64(0)
	)

	sdb := NewTestStateDB()
	for _, v := range params {
		switch v.(type) {
		case int:
			block = v.(int)
		case common.Address:
			if sender == common.EmptyAddress {
				sender = v.(common.Address)
			} else if caller == common.EmptyAddress {
				caller = v.(common.Address)
			}
		case common.Hash:
			hash = v.(common.Hash)
		case uint64:
			supplyGas = v.(uint64)
		case func(*state.StateDB):
			fn := v.(func(*state.StateDB))
			fn(sdb)
		default:
			t.Fatal("invalid params type")
		}
	}

	blockHeight := new(big.Int).SetInt64(int64(block))
	contractRef := NewContractRef(sdb, sender, caller, blockHeight, hash, supplyGas, nil)
	ctx := NewNativeContract(sdb, contractRef)

	// need to break point at the next step, e.g: nativeCall or some contract function
	ctx.BreakPoint()
	return sdb, ctx
}

func TestNativeCall(t *testing.T, contract common.Address, name string, payload []byte, params ...interface{}) ([]byte, error) {
	_, ctx := GenerateTestContext(t, params...)
	ref := ctx.ContractRef()
	res, _, err := ref.NativeCall(ref.caller, contract, payload)
	t.Logf("contract %s method %s execute time %v", contract.Hex(), name, ctx.BreakPoint())
	return res, err
}
