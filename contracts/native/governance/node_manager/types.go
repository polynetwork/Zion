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

package node_manager

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
)

type LockStatus uint8

const (
	Unspecified LockStatus = 0
	Unlocked    LockStatus = 1
	Unlocking   LockStatus = 2
	Locked      LockStatus = 3
)

type AllValidators struct {
	AllValidators []string
}

func (m *AllValidators) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.AllValidators})
}

func (m *AllValidators) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		AllValidators []string
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.AllValidators = data.AllValidators
	return nil
}

type Validator struct {
	StakeAddress    common.Address
	ConsensusPubkey string
	ProposalAddress common.Address
	Commission      *Commission
	Status          LockStatus
	Jailed          bool
	UnlockTime      *big.Int
	TotalStake      *big.Int
	SelfStake       *big.Int
	Desc            string
}

func (m *Validator) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.StakeAddress, m.ConsensusPubkey, m.ProposalAddress, m.Commission, m.Status, m.Jailed, m.UnlockTime,
		m.TotalStake, m.SelfStake, m.Desc})
}

func (m *Validator) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		StakeAddress    common.Address
		ConsensusPubkey string
		ProposalAddress common.Address
		Commission      *Commission
		Status          LockStatus
		Jailed          bool
		UnlockTime      *big.Int
		TotalStake      *big.Int
		SelfStake       *big.Int
		Desc            string
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.StakeAddress, m.ConsensusPubkey, m.ProposalAddress, m.Commission, m.Status, m.Jailed, m.UnlockTime, m.TotalStake,
		m.SelfStake, m.Desc = data.StakeAddress, data.ConsensusPubkey, data.ProposalAddress, data.Commission,
		data.Status, data.Jailed, data.UnlockTime, data.TotalStake, data.SelfStake, data.Desc
	return nil
}

// IsLocked checks if the validator status equals Locked
func (m Validator) IsLocked() bool {
	return m.Status == Locked
}

// IsUnlocked checks if the validator status equals Unlocked
func (m Validator) IsUnlocked() bool {
	return m.Status == Unlocked
}

// IsUnlocking checks if the validator status equals Unlocking
func (m Validator) IsUnlocking() bool {
	return m.Status == Unlocking
}

type Commission struct {
	Rate         *big.Int
	UpdateHeight *big.Int
}

func (m *Commission) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Rate, m.UpdateHeight})
}

func (m *Commission) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Rate         *big.Int
		UpdateHeight *big.Int
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.Rate, m.UpdateHeight = data.Rate, data.UpdateHeight
	return nil
}

type GlobalConfig struct {
	MaxCommission   *big.Int
	MinInitialStake *big.Int
	MaxDescLength   uint64
}

func (m *GlobalConfig) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.MaxCommission, m.MinInitialStake, m.MaxDescLength})
}

func (m *GlobalConfig) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		MaxCommission   *big.Int
		MinInitialStake *big.Int
		MaxDescLength   uint64
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.MaxCommission, m.MinInitialStake, m.MaxDescLength = data.MaxCommission, data.MinInitialStake, data.MaxDescLength
	return nil
}

type StakeInfo struct {
	StakeAddress    common.Address
	ConsensusPubkey string
	Amount          *big.Int
}

func (m *StakeInfo) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.StakeAddress, m.ConsensusPubkey, m.Amount})
}

func (m *StakeInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		StakeAddress    common.Address
		ConsensusPubkey string
		Amount          *big.Int
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.StakeAddress, m.ConsensusPubkey, m.Amount = data.StakeAddress, data.ConsensusPubkey, data.Amount
	return nil
}
