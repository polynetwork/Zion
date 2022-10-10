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
	AllValidators []common.Address
}

func (m *AllValidators) Decode(payload []byte) error {
	var data struct {
		AllValidators []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetAllValidators, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.AllValidators, m)
}

type Validator struct {
	StakeAddress     common.Address
	ConsensusAddress common.Address
	SignerAddress    common.Address
	ProposalAddress  common.Address
	Commission       *Commission
	Status           LockStatus
	Jailed           bool
	UnlockHeight     *big.Int
	TotalStake       Dec
	SelfStake        Dec
	Desc             string
}

func (m *Validator) Decode(payload []byte) error {
	var data struct {
		Validator []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetValidator, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.Validator, m)
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
	Rate         Dec
	UpdateHeight *big.Int
}

type GlobalConfig struct {
	MaxCommissionChange   *big.Int
	MinInitialStake       *big.Int
	MinProposalStake      *big.Int
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
	StakeAddress  common.Address
	ConsensusAddr common.Address
	Amount        Dec
}

func (m *StakeInfo) Decode(payload []byte) error {
	var data struct {
		StakeInfo []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetStakeInfo, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.StakeInfo, m)
}

type UnlockingInfo struct {
	StakeAddress   common.Address
	UnlockingStake []*UnlockingStake
}

func (m *UnlockingInfo) Decode(payload []byte) error {
	var data struct {
		UnlockingInfo []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetUnlockingInfo, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.UnlockingInfo, m)
}

type UnlockingStake struct {
	Height           *big.Int
	CompleteHeight   *big.Int
	ConsensusAddress common.Address
	Amount           Dec
}

type EpochInfo struct {
	ID          *big.Int
	Validators  []common.Address
	Signers     []common.Address
	Voters      []common.Address
	Proposers   []common.Address
	StartHeight *big.Int
	EndHeight   *big.Int
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

func (m *EpochInfo) SignerQuorumSize() int {
	if m == nil || m.Signers == nil {
		return 0
	}
	total := len(m.Signers)
	return int(math.Ceil(float64(2*total) / 3))
}

func (m *EpochInfo) VoterQuorumSize() int {
	if m == nil || m.Voters == nil {
		return 0
	}
	total := len(m.Voters)
	return int(math.Ceil(float64(2*total) / 3))
}

func (m *EpochInfo) ProposerQuorumSize() int {
	if m == nil || m.Proposers == nil {
		return 0
	}
	total := len(m.Proposers)
	return int(math.Ceil(float64(2*total) / 3))
}

func (m *EpochInfo) MemberList() []common.Address {
	list := make([]common.Address, 0)
	if m == nil || m.Validators == nil || len(m.Validators) == 0 {
		return list
	}
	for _, v := range m.Validators {
		list = append(list, v)
	}
	return list
}

type AccumulatedCommission struct {
	Amount Dec
}

func (m *AccumulatedCommission) Decode(payload []byte) error {
	var data struct {
		AccumulatedCommission []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetAccumulatedCommission, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.AccumulatedCommission, m)
}

type ValidatorAccumulatedRewards struct {
	Rewards Dec
	Period  uint64
}

func (m *ValidatorAccumulatedRewards) Decode(payload []byte) error {
	var data struct {
		ValidatorAccumulatedRewards []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetValidatorAccumulatedRewards, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.ValidatorAccumulatedRewards, m)
}

type ValidatorOutstandingRewards struct {
	Rewards Dec
}

func (m *ValidatorOutstandingRewards) Decode(payload []byte) error {
	var data struct {
		ValidatorOutstandingRewards []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetValidatorOutstandingRewards, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.ValidatorOutstandingRewards, m)
}

type OutstandingRewards struct {
	Rewards Dec
}

func (m *OutstandingRewards) Decode(payload []byte) error {
	var data struct {
		OutstandingRewards []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetOutstandingRewards, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.OutstandingRewards, m)
}

type StakeRewards struct {
	Rewards Dec
}

func (m *StakeRewards) Decode(payload []byte) error {
	var data struct {
		StakeRewards []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetStakeRewards, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.StakeRewards, m)
}

type ValidatorSnapshotRewards struct {
	AccumulatedRewardsRatio Dec // ratio already mul decimal
	ReferenceCount          uint64
}

func (m *ValidatorSnapshotRewards) Decode(payload []byte) error {
	var data struct {
		ValidatorSnapshotRewards []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetValidatorSnapshotRewards, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.ValidatorSnapshotRewards, m)
}

type StakeStartingInfo struct {
	StartPeriod uint64
	Stake       Dec
	Height      *big.Int
}

func (m *StakeStartingInfo) Decode(payload []byte) error {
	var data struct {
		StakeStartingInfo []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetStakeStartingInfo, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.StakeStartingInfo, m)
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

type TotalPool struct {
	TotalPool Dec
}

func (m *TotalPool) Decode(payload []byte) error {
	var data struct {
		TotalPool []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodGetTotalPool, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.TotalPool, m)
}
