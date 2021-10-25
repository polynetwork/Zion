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
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
)

func init() {
	core.StoreGenesis = func(db ethdb.Database, header *types.Header) error {
		extra, err := types.ExtractHotstuffExtra(header)
		if err != nil {
			return err
		}
		epoch := &Epoch{
			BlockHash:   header.Hash(),
			StartHeight: 0,
			ValSet:      newValSet(extra.Validators),
		}
		return storeCurEpoch(db, epoch)
	}
}

func (s *backend) Validators() hotstuff.ValidatorSet {
	return s.epoch.ValSet.Copy()
}

func (s *backend) LoadEpoch() error {
	epoch, err := getCurEpoch(s.db)
	if err != nil {
		return err
	}
	s.epoch = epoch
	return nil
}

func (s *backend) UpdateEpoch(header *types.Header) error {
	extra, err := types.ExtractHotstuffExtra(header)
	if err != nil {
		return err
	}

	height := header.Number.Uint64()
	if extra.Validators != nil && len(extra.Validators) != 0 {
		return s.PrepareChangeEpoch(extra.Validators, height)
	}

	if s.nxtEpoch != nil && height == s.nxtEpoch.StartHeight+1 {
		return s.ChangeEpoch(height, extra.Validators)
	}
	return nil
}

func (s *backend) PrepareChangeEpoch(list []common.Address, startHeight uint64) error {
	if s.epoch == nil {
		return fmt.Errorf("current epoch is invalid")
	}

	// duplicate change
	if s.nxtEpoch != nil && s.nxtEpoch.StartHeight == startHeight {
		return nil
	}
	// already changed
	if s.epoch.StartHeight == startHeight {
		return nil
	}
	// after change or invalid start height
	if s.epoch.StartHeight > startHeight {
		return nil
	}

	cur := s.Validators()
	if s.epoch.ValSet.ParticipantsNumber(list) < cur.Q() {
		return fmt.Errorf("next val set not unreach current epoch's quorum size, (cur, next) (%v, %v)", cur, list)
	}

	ep := &Epoch{
		StartHeight: startHeight,
		ValSet:      newValSet(list),
	}
	s.nxtEpoch = ep
	return nil
}

func (s *backend) ChangeEpoch(epochStartHeight uint64, list []common.Address) error {
	if s.nxtEpoch == nil || s.epoch == nil {
		return fmt.Errorf("current epoch or next epoch is nil")
	}
	s.epoch = s.nxtEpoch.Copy()
	s.nxtEpoch = nil

	return storeCurEpoch(s.db, s.epoch)
}

func storeCurEpoch(db ethdb.Database, epoch *Epoch) error {
	blob, err := epoch.MarshalJSON()
	if err != nil {
		return err
	}
	if err := rawdb.WriteEpoch(db, epoch.BlockHash, blob); err != nil {
		return err
	}
	return rawdb.WriteCurrentEpochHash(db, epoch.BlockHash)
}

func getCurEpoch(db ethdb.Database) (*Epoch, error) {
	hash, err := rawdb.ReadCurrentEpochHash(db)
	if err != nil {
		return nil, err
	}
	blob, err := rawdb.ReadEpoch(db, hash)
	if err != nil {
		return nil, err
	}
	epoch := new(Epoch)
	if err := epoch.UnmarshalJSON(blob); err != nil {
		return nil, err
	}
	return epoch, nil
}

type Epoch struct {
	BlockHash   common.Hash
	StartHeight uint64
	ValSet      hotstuff.ValidatorSet
}

func (e *Epoch) Copy() *Epoch {
	return &Epoch{
		BlockHash:   e.BlockHash,
		StartHeight: e.StartHeight,
		ValSet:      e.ValSet.Copy(),
	}
}

type epochJSON struct {
	BlockHash   common.Hash      `json:"block_hash"`
	StartHeight uint64           `json:"start_height"`
	Validators  []common.Address `json:"validators"`
}

func (e *Epoch) toJSONStruct() *epochJSON {
	return &epochJSON{
		BlockHash:   e.BlockHash,
		StartHeight: e.StartHeight,
		Validators:  e.ValSet.AddressList(),
	}
}

// Unmarshal from a json byte array
func (e *Epoch) UnmarshalJSON(b []byte) error {
	var j epochJSON
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	e.BlockHash = j.BlockHash
	e.StartHeight = j.StartHeight
	e.ValSet = newValSet(j.Validators)
	return nil
}

// Marshal to a json byte array
func (e *Epoch) MarshalJSON() ([]byte, error) {
	j := e.toJSONStruct()
	return json.Marshal(j)
}

func newValSet(list []common.Address) hotstuff.ValidatorSet {
	return validator.NewSet(list, hotstuff.RoundRobin)
}
