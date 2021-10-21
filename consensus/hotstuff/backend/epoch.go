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
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/ethereum/go-ethereum/core/types"
)

type epoch struct {
	startHeight uint64
	valset      hotstuff.ValidatorSet
}

func (e *epoch) Copy() *epoch {
	return &epoch{
		startHeight: e.startHeight,
		valset:      e.valset.Copy(),
	}
}

func (s *backend) Validators() hotstuff.ValidatorSet {
	return s.epoch.valset.Copy()
}

func (s *backend) UpdateEpoch(header *types.Header) error {
	extra, err := types.ExtractHotstuffExtra(header)
	if err != nil {
		return err
	}

	height := header.Number.Uint64()
	if height == 0 {
		s.SetGenesisEpoch(extra.Validators)
		return nil
	}

	if extra.Validators != nil && len(extra.Validators) != 0 {
		return s.PrepareChangeEpoch(extra.Validators, height)
	}

	if height == s.nxtEpoch.startHeight+1 {
		s.ChangeEpoch(height, extra.Validators)
	}
	return nil
}

func (s *backend) SetGenesisEpoch(list []common.Address) {
	if s.epoch != nil && s.epoch.startHeight == 0 {
		return
	}

	if list == nil || len(list) == 0 {
		panic("invalid validator set")
	}

	s.epoch = &epoch{
		startHeight: 0,
		valset:      s.newValSet(list),
	}
}

func (s *backend) PrepareChangeEpoch(list []common.Address, startHeight uint64) error {
	// duplicate change
	if s.nxtEpoch != nil && s.nxtEpoch.startHeight == startHeight {
		return nil
	}
	// already changed
	if s.epoch.startHeight == startHeight {
		return nil
	}
	// after change or invalid start height
	if startHeight <= s.epoch.startHeight {
		return nil
	}

	cur := s.Validators()
	if s.epoch.valset.ParticipantsNumber(list) < cur.Q() {
		return fmt.Errorf("next val set not unreach current epoch's quorum size, (cur, next) (%v, %v)", cur, list)
	}

	ep := &epoch{
		startHeight: startHeight,
		valset:      s.newValSet(list),
	}
	s.nxtEpoch = ep
	return nil
}

func (s *backend) ChangeEpoch(epochStartHeight uint64, list []common.Address) {
	s.epoch = s.nxtEpoch.Copy()
	s.nxtEpoch = nil
}

func (s *backend) newValSet(list []common.Address) hotstuff.ValidatorSet {
	return validator.NewSet(list, hotstuff.RoundRobin)
}
