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
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
)

var (
	startSnapID       = nm.StartEpochID.Uint64()
	genesisSnapStart  = uint64(0)
	genesisSnapLength = nm.GenesisBlockPerEpoch.Uint64()
)

func init() {
	core.StoreGenesis = func(db ethdb.Database, header *types.Header) error {
		extra, err := types.ExtractHotstuffExtra(header)
		if err != nil {
			return err
		}
		snap := &snapshot{
			ID:     startSnapID,
			Start:  genesisSnapStart,
			End:    calcEnd(genesisSnapStart, genesisSnapLength),
			ValSet: NewDefaultValSet(extra.Validators),
		}
		return snap.store(db)
	}
}

type snapshots struct {
	list []*snapshot
	mu   sync.Mutex
}

func newSnapshots() *snapshots {
	snaps := new(snapshots)
	snaps.list = make([]*snapshot, 0)
	return snaps
}

// load read epoch from database and append in slice cache while consensus engine started.
func (s *snapshots) load(db ethdb.Database) {
	id := startSnapID

	for {
		snap := new(snapshot)
		if err := snap.load(db, id); err != nil {
			return
		}
		if !s.append(snap) {
			return
		}
		id = s.nextId()
		log.Info("[snap]", "load snap", snap.String())
	}
}

// SyncEpoch light mode
func (s *snapshots) sync(db ethdb.Database, parent, header *types.Header) error {
	if parent.Number.Uint64() == 0 {
		return nil
	}

	parentExt, err := types.ExtractHotstuffExtra(parent)
	if err != nil {
		return err
	}
	if parentExt.Validators == nil || len(parentExt.Validators) == 0 {
		return nil
	}

	epoch := &snapshot{
		ID:     s.nextId(),
		Start:  header.Number.Uint64(),
		ValSet: NewDefaultValSet(parentExt.Validators),
	}

	if s.append(epoch) {
		if err := epoch.store(db); err != nil {
			return err
		}
		log.Info("[epoch]", "sync epoch", epoch.String())
	}

	return nil
}

func (s *snapshots) Dump() string {
	str := ""
	for _, v := range s.list {
		str += v.String() + "\r\n"
	}
	return str
}

func (s *snapshots) get(height uint64) *snapshot {
	for _, snap := range s.list {
		if height >= snap.Start {
			return snap
		}
	}
	// `height` is an uint number, and the min snap height is 0, so the function will return in above loop
	return nil
}

func (s *snapshots) append(snap *snapshot) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.list) == 0 {
		s.list = append(s.list, snap)
		return true
	}

	// already exist
	if snap.ID <= s.id() || snap.Start <= s.start() {
		return false
	}

	s.list = append(s.list, snap)
	s.list[0], s.list[len(s.list)-1] = s.list[len(s.list)-1], s.list[0]

	return true
}

// id retrieve current snapshot identity
func (s *snapshots) id() uint64 {
	if s.list == nil {
		return startSnapID
	}
	return s.list[0].ID
}

// start retrieve current snapshot start height
func (s *snapshots) start() uint64 {
	if s.list == nil {
		return genesisSnapStart
	}
	return s.list[0].Start
}

// end retrieve current snapshot end height
func (s *snapshots) end() uint64 {
	if s.list == nil {
		return genesisSnapLength
	}
	return s.list[0].End
}

// nextStart retrieve next snapshot start height, which should be equal to current end + 1
func (s *snapshots) nextStart() uint64 {
	if s.list == nil {
		return genesisSnapStart
	}
	return s.list[0].End + 1
}

// nextId retrieve next snapshot identity
func (s *snapshots) nextId() uint64 {
	if s.list == nil {
		return startSnapID
	}
	return s.id() + 1
}

type snapshot struct {
	ID     uint64 // id started from 1
	Start  uint64 // start block height
	End    uint64 // end block height
	ValSet hotstuff.ValidatorSet
}

func newSnapshot(id, start, end uint64, list []common.Address) *snapshot {
	return &snapshot{
		ID:     id,
		Start:  start,
		End:    end,
		ValSet: NewDefaultValSet(list),
	}
}

func (e *snapshot) store(db ethdb.Database) error {
	blob, err := e.MarshalJSON()
	if err != nil {
		return err
	}
	return rawdb.WriteEpoch(db, e.ID, blob)
}

func (e *snapshot) load(db ethdb.Database, id uint64) error {
	blob, err := rawdb.ReadEpoch(db, id)
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

// ----------------- utility functions ------------------

func calcEnd(start, length uint64) uint64 {
	return start + length - 1
}
