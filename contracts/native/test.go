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
	"fmt"
	"math/big"
	"testing"
	"time"

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
func GenerateTestContext(t *testing.T, value *big.Int, address common.Address, params ...interface{}) (*state.StateDB, *NativeContract) {
	var (
		block     = int(0)
		caller    = common.EmptyAddress
		sender    = common.EmptyAddress
		hash      = common.EmptyHash
		supplyGas = uint64(0)
	)

	var sdb *state.StateDB
	flag := false
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
			sdb = NewTestStateDB()
			fn := v.(func(*state.StateDB))
			fn(sdb)
			flag = true
		case *state.StateDB:
			sdb = v.(*state.StateDB)
			flag = true
		default:
			t.Fatal("invalid params type")
		}
	}

	if !flag {
		sdb = NewTestStateDB()
	}
	// set caller = sender if params only contains sender
	if sender != common.EmptyAddress && caller == common.EmptyAddress {
		caller = sender
	}

	blockHeight := new(big.Int).SetInt64(int64(block))
	contractRef := NewContractRef(sdb, sender, caller, blockHeight, hash, supplyGas, nil)
	contractRef.SetValue(value)
	contractRef.SetTo(address)
	ctx := NewNativeContract(sdb, contractRef)

	// need to break point at the next step, e.g: nativeCall or some contract function
	if DebugSpentOpen {
		ctx.BreakPoint()
	}
	return sdb, ctx
}

func TestNativeCall(t *testing.T, contract common.Address, name string, payload []byte, value *big.Int, params ...interface{}) ([]byte, error) {
	_, ctx := GenerateTestContext(t, value, contract, params...)
	ref := ctx.ContractRef()
	res, _, err := ref.NativeCall(ref.caller, contract, payload)
	if DebugSpentOpen {
		t.Logf("contract %s method %s execute time %v us", contract.Hex(), name, ctx.BreakPoint())
	}
	return res, err
}

type Timer struct {
	name  string
	since time.Time
	cases []time.Duration
}

func NewTimer(name string) *Timer {
	return &Timer{name: name}
}

func (t *Timer) Start() {
	if !t.since.IsZero() {
		panic("timer not clean")
	}
	t.since = time.Now()
}

func (t *Timer) Stop() {
	if t.since.IsZero() {
		panic("time was not started")
	}
	t.cases = append(t.cases, time.Since(t.since))
	t.since = time.Time{}
}

func (t *Timer) Add(duration time.Duration) {
	if t.since.IsZero() {
		panic("time was not started")
	}
	t.cases = append(t.cases, duration)
}

func (t *Timer) Dump() {
	var min, max, avg time.Duration
	for _, c := range t.cases {
		if min == 0 || c < min {
			min = c
		}
		if c > max {
			max = c
		}
		avg += c
	}
	if len(t.cases) > 0 {
		avg /= time.Duration(len(t.cases))
	}
	fmt.Printf("%s cases(%d), min: %d, max: %d, avg: %d\n", t.name, len(t.cases), min, max, avg)
}

func (t *Timer) ResetAs(name string) {
	if !t.since.IsZero() {
		panic("time is diry before reset")
	}
	t.name = name
	t.cases = []time.Duration{}
}
