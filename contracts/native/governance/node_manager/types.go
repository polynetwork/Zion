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

func (m *Validator) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.StakeAddress, m.ConsensusPubkey, m.ConsensusAddress,
		m.ProposalAddress, m.Commission, m.Status, m.Jailed, m.UnlockHeight, m.TotalStake, m.SelfStake, m.Desc})
}

func (m *Validator) DecodeRLP(s *rlp.Stream) error {
	var data struct {
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

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.StakeAddress, m.ConsensusPubkey, m.ConsensusAddress, m.ProposalAddress, m.Commission, m.Status, m.Jailed,
		m.UnlockHeight, m.TotalStake, m.SelfStake, m.Desc = data.StakeAddress, data.ConsensusPubkey, data.ConsensusAddress,
		data.ProposalAddress, data.Commission, data.Status, data.Jailed, data.UnlockHeight, data.TotalStake,
		data.SelfStake, data.Desc
	return nil
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
	MaxCommission         *big.Int
	MinInitialStake       *big.Int
	MaxDescLength         uint64
	BlockPerEpoch         *big.Int
	ConsensusValidatorNum uint64
	VoterValidatorNum     uint64
}

func (m *GlobalConfig) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.MaxCommission, m.MinInitialStake, m.MaxDescLength, m.BlockPerEpoch,
		m.ConsensusValidatorNum, m.VoterValidatorNum})
}

func (m *GlobalConfig) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		MaxCommission         *big.Int
		MinInitialStake       *big.Int
		MaxDescLength         uint64
		BlockPerEpoch         *big.Int
		ConsensusValidatorNum uint64
		VoterValidatorNum     uint64
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.MaxCommission, m.MinInitialStake, m.MaxDescLength, m.BlockPerEpoch, m.ConsensusValidatorNum,
		m.VoterValidatorNum = data.MaxCommission, data.MinInitialStake, data.MaxDescLength, data.BlockPerEpoch,
		data.ConsensusValidatorNum, data.VoterValidatorNum
	return nil
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

type UnlockingInfo struct {
	StakeAddress   common.Address
	UnlockingStake []*UnlockingStake
}

func (m *UnlockingInfo) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.StakeAddress, m.UnlockingStake})
}

func (m *UnlockingInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		StakeAddress   common.Address
		UnlockingStake []*UnlockingStake
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.StakeAddress, m.UnlockingStake = data.StakeAddress, data.UnlockingStake
	return nil
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

func (m *EpochInfo) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ID, m.Validators, m.Voters, m.StartHeight})
}
func (m *EpochInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ID          *big.Int
		Validators  []*Peer
		Voters      []*Peer
		StartHeight *big.Int
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.ID, m.Validators, m.Voters, m.StartHeight = data.ID, data.Validators, data.Voters, data.StartHeight
	return nil
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

type OutstandingRewards struct {
	Rewards *big.Int
}

func (m *OutstandingRewards) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Rewards})
}

func (m *OutstandingRewards) DecodeRLP(s *rlp.Stream) error {
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

type StakeStartingInfo struct {
	StartPeriod uint64
	Stake       *big.Int
	Height      *big.Int
}

func (m *StakeStartingInfo) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.StartPeriod, m.Stake, m.Height})
}

func (m *StakeStartingInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		StartPeriod uint64
		Stake       *big.Int
		Height      *big.Int
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.StartPeriod, m.Stake, m.Height = data.StartPeriod, data.Stake, data.Height
	return nil
}

type Peer struct {
	PubKey  string
	Address common.Address
}

func (m *Peer) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.PubKey, m.Address})
}

func (m *Peer) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ConsensusPubkey  string
		ConsensusAddress common.Address
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.PubKey, m.Address = data.ConsensusPubkey, data.ConsensusAddress
	return nil
}

type AddressList struct {
	List []common.Address
}

func (m *AddressList) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.List})
}

func (m *AddressList) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		List []common.Address
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.List = data.List
	return nil
}

type ConsensusSign struct {
	Method string
	Input  []byte
	hash   atomic.Value
}

func (m *ConsensusSign) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Method, m.Input})
}
func (m *ConsensusSign) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Method string
		Input  []byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.Method, m.Input = data.Method, data.Input
	return nil
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

func (m *CommunityInfo) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.CommunityRate, m.CommunityAddress})
}
func (m *CommunityInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		CommunityRate    *big.Int
		CommunityAddress common.Address
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.CommunityRate, m.CommunityAddress = data.CommunityRate, data.CommunityAddress
	return nil
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
