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
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math"
	"math/big"
	"sync/atomic"
)

type LockStatus uint8

const (
	Unspecified LockStatus = 0
	Unlock      LockStatus = 1
	Lock        LockStatus = 2
	Remove      LockStatus = 3
)

type AllValidators struct {
	AllValidators []string
}

type Validator struct {
	StakeAddress     common.Address
	ConsensusPubkey  string
	ConsensusAddress common.Address
	ProposalAddress  common.Address
	Commission       *Commission
	Status           LockStatus
	Jailed           bool
	UnlockHeight     *big.Int
	TotalStake       *big.Int
	SelfStake        *big.Int
	Desc             string
}

// IsLocked checks if the validator status equals Locked
func (m Validator) IsLocked() bool {
	return m.Status == Lock
}

// IsUnlocked checks if the validator status equals Unlocked
func (m Validator) IsUnlocked(height *big.Int) bool {
	return m.Status == Unlock && m.UnlockHeight.Cmp(height) <= 0
}

// IsUnlocking checks if the validator status equals Unlocking
func (m Validator) IsUnlocking(height *big.Int) bool {
	return m.Status == Unlock && m.UnlockHeight.Cmp(height) > 0
}

// IsRemoved checks if the validator status equals Unlocked
func (m Validator) IsRemoved(height *big.Int) bool {
	return m.Status == Remove && m.UnlockHeight.Cmp(height) <= 0
}

// IsRemoving checks if the validator status equals Unlocking
func (m Validator) IsRemoving(height *big.Int) bool {
	return m.Status == Remove && m.UnlockHeight.Cmp(height) > 0
}

type Commission struct {
	Rate         *big.Int
	UpdateHeight *big.Int
}

type GlobalConfig struct {
	MaxCommission         *big.Int
	MinInitialStake       *big.Int
	MaxDescLength         uint64
	BlockPerEpoch         *big.Int
	ConsensusValidatorNum uint64
	VoterValidatorNum     uint64
}

func (m *GlobalConfig) Decode(payload []byte) error {
	var data struct {
		GlobalConfig []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetGlobalConfig, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.GlobalConfig, m)
}

type StakeInfo struct {
	StakeAddress    common.Address
	ConsensusPubkey string
	Amount          *big.Int
}

type UnlockingInfo struct {
	StakeAddress   common.Address
	UnlockingStake []*UnlockingStake
}

type UnlockingStake struct {
	Height         *big.Int
	CompleteHeight *big.Int
	Amount         *big.Int
}

func (m *UnlockingStake) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Height, m.CompleteHeight, m.Amount})
}

func (m *UnlockingStake) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Height         *big.Int
		CompleteHeight *big.Int
		Amount         *big.Int
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.Height, m.CompleteHeight, m.Amount = data.Height, data.CompleteHeight, data.Amount
	return nil
}

type EpochInfo struct {
	ID          *big.Int
	Validators  []*Peer
	Voters      []*Peer
	StartHeight *big.Int
}

func (m *EpochInfo) Decode(payload []byte) error {
	var data struct {
		EpochInfo []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetCurrentEpochInfo, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.EpochInfo, m)
}

func (m *EpochInfo) ValidatorQuorumSize() int {
	if m == nil || m.Validators == nil {
		return 0
	}
	total := len(m.Validators)
	return int(math.Ceil(float64(2*total) / 3))
}

func (m *EpochInfo) VoterQuorumSize() int {
	if m == nil || m.Voters == nil {
		return 0
	}
	total := len(m.Voters)
	return int(math.Ceil(float64(2*total) / 3))
}

func (m *EpochInfo) MemberList() []common.Address {
	list := make([]common.Address, 0)
	if m == nil || m.Validators == nil || len(m.Validators) == 0 {
		return list
	}
	for _, v := range m.Validators {
		list = append(list, v.Address)
	}
	return list
}

type AccumulatedCommission struct {
	Amount *big.Int
}

type ValidatorAccumulatedRewards struct {
	Rewards *big.Int
	Period  uint64
}

type ValidatorOutstandingRewards struct {
	Rewards *big.Int
}

type OutstandingRewards struct {
	Rewards *big.Int
}

type ValidatorSnapshotRewards struct {
	AccumulatedRewardsRatio *big.Int // ratio already mul decimal
	ReferenceCount          uint64
}

type StakeStartingInfo struct {
	StartPeriod uint64
	Stake       *big.Int
	Height      *big.Int
}

type Peer struct {
	PubKey  string
	Address common.Address
}

type AddressList struct {
	List []common.Address
}

type ConsensusSign struct {
	Method string
	Input  []byte
	hash   atomic.Value
}

func (m *ConsensusSign) Hash() common.Hash {
	if hash := m.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	var inf = struct {
		Method string
		Input  []byte
	}{
		Method: m.Method,
		Input:  m.Input,
	}
	v := utils.RLPHash(inf)
	m.hash.Store(v)
	return v
}

type CommunityInfo struct {
	CommunityRate    *big.Int
	CommunityAddress common.Address
}

func (m *CommunityInfo) Decode(payload []byte) error {
	var data struct {
		CommunityInfo []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetCommunityInfo, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.CommunityInfo, m)
}
