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

package distribute

import (
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
)

type AccumulatedCommission struct {
	Amount *big.Int
}

func (m *AccumulatedCommission) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Amount})
}

func (m *AccumulatedCommission) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Amount *big.Int
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.Amount = data.Amount
	return nil
}

type ValidatorAccumulatedRewards struct {
	Rewards *big.Int
	Period  uint64
}

func (m *ValidatorAccumulatedRewards) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Rewards, m.Period})
}

func (m *ValidatorAccumulatedRewards) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Rewards *big.Int
		Period  uint64
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.Rewards, m.Period = data.Rewards, data.Period
	return nil
}

type ValidatorOutstandingRewards struct {
	Rewards *big.Int
}

func (m *ValidatorOutstandingRewards) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Rewards})
}

func (m *ValidatorOutstandingRewards) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Rewards *big.Int
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.Rewards = data.Rewards
	return nil
}

type ValidatorSnapshotRewards struct {
	AccumulatedRewardsRatio *big.Int // ratio already mul decimal
	ReferenceCount          uint64
}

func (m *ValidatorSnapshotRewards) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.AccumulatedRewardsRatio, m.ReferenceCount})
}

func (m *ValidatorSnapshotRewards) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		AccumulatedRewardsRatio *big.Int
		ReferenceCount          uint64
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.AccumulatedRewardsRatio, m.ReferenceCount = data.AccumulatedRewardsRatio, data.ReferenceCount
	return nil
}
