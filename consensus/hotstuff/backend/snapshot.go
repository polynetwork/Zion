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

package backend

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
)

type snapshot struct {
	ID     uint64 // id started from 1
	Start  uint64 // start block height
	End    uint64 // end block height
	ValSet hotstuff.ValidatorSet
}

func (e *snapshot) store(db ethdb.Database) error {
	blob, err := e.MarshalJSON()
	if err != nil {
		return err
	}
	return rawdb.WriteEpoch(db, e.ID, blob)
}

func (e *snapshot) load(db ethdb.Database, startHeight uint64) error {
	blob, err := rawdb.ReadEpoch(db, startHeight)
	if err != nil {
		return err
	}
	return e.UnmarshalJSON(blob)
}

func (e *snapshot) copy() *snapshot {
	return &snapshot{
		ID:     e.ID,
		Start:  e.Start,
		End:    e.End,
		ValSet: e.ValSet.Copy(),
	}
}

func (e *snapshot) String() string {
	return fmt.Sprintf("{ID: %d, Start: %d, End: %d, Size: %d, Valset: %v}", e.ID, e.Start, e.End, e.ValSet.Size(), e.ValSet.AddressList())
}

type epochJSON struct {
	ID         uint64           `json:"id"`
	Start      uint64           `json:"start"`
	End        uint64           `json:"end"`
	Validators []common.Address `json:"validators"`
}

func (e *snapshot) toJSONStruct() *epochJSON {
	return &epochJSON{
		ID:         e.ID,
		Start:      e.Start,
		End:        e.End,
		Validators: e.ValSet.AddressList(),
	}
}

// Unmarshal from a json byte array
func (e *snapshot) UnmarshalJSON(b []byte) error {
	var j epochJSON
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	e.ID = j.ID
	e.Start = j.Start
	e.End = j.End
	e.ValSet = NewDefaultValSet(j.Validators)
	return nil
}

// Marshal to a json byte array
func (e *snapshot) MarshalJSON() ([]byte, error) {
	j := e.toJSONStruct()
	return json.Marshal(j)
}

func NewDefaultValSet(list []common.Address) hotstuff.ValidatorSet {
	return validator.NewSet(list, hotstuff.RoundRobin)
}
