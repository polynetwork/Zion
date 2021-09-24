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
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
)

const contractName = "node manager"

const (
	MethodContractName = "name"
	MethodPropose      = "propose"
	MethodVote         = "vote"
	MethodEpoch        = "epoch"
	MethodNextEpoch    = "nextEpoch"

	EventPropose      = "proposed"
	EventVote         = "voted"
	EventEpochChange = "epochChanged"
)

const abijson = `[
    {"type":"function","name":"` + MethodContractName + `","inputs":[],"outputs":[{"internalType":"string","name":"Name","type":"string"}],"stateMutability":"nonpayable"},
	{"type":"function","name":"` + MethodPropose + `","inputs":[{"internalType":"uint64","name":"StartHeight","type":"uint64"},{"internalType":"bytes","name":"Peers","type":"bytes"}],"outputs":[{"internalType":"bool","name":"Success","type":"bool"}],"stateMutability":"nonpayable"},
    {"type":"function","name":"` + MethodVote + `","inputs":[{"internalType":"uint64","name":"EpochID","type":"uint64"},{"internalType":"bytes","name":"Hash","type":"bytes"}],"outputs":[{"internalType":"bool","name":"Success","type":"bool"}],"stateMutability":"nonpayable"},
	{"type":"function","name":"` + MethodEpoch + `","inputs":[],"outputs":[{"internalType":"bytes","name":"EpochInfo","type":"bytes"}],"stateMutability":"nonpayable"},
	{"type":"function","name":"` + MethodNextEpoch + `","inputs":[],"outputs":[{"internalType":"bytes","name":"EpochInfo","type":"bytes"}],"stateMutability":"nonpayable"},
    {"type":"event","name":"` + EventPropose + `","anonymous":false,"inputs":[{"internalType":"bytes","name":"EpochInfo","type":"bytes"}]},
	{"type":"event","name":"` + EventVote + `","anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"EpochID","type":"uint64"},{"indexed":false,"internalType":"bytes","name":"Hash","type":"bytes"},{"indexed":false,"internalType":"uint64","name":"VotedNumber","type":"uint64"},{"indexed":false,"internalType":"uint64","name":"GroupSize","type":"uint64"}]},
	{"type":"event","name":"` + EventEpochChange + `","anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"EpochInfo","type":"bytes"},{"indexed":false,"internalType":"bytes","name":"NextEpochInfo","type":"bytes"}]}
]`

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(abijson))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

var (
	ABI  *abi.ABI
	this = utils.NodeManagerContractAddress
)

type MethodContractNameOutput struct {
	Name string
}

func (m *MethodContractNameOutput) Encode(name string) ([]byte, error) {
	m.Name = name
	return utils.PackOutputs(ABI, MethodContractName, m.Name)
}
func (m *MethodContractNameOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodContractName, m, payload)
}

type MethodProposeInput struct {
	StartHeight uint64
	Peers       []byte
}

func (m *MethodProposeInput) Encode(epochID uint64, peers *Peers) ([]byte, error) {
	enc, err := rlp.EncodeToBytes(peers)
	if err != nil {
		return nil, err
	}
	m.StartHeight = epochID
	m.Peers = enc
	return utils.PackMethod(ABI, MethodPropose, m.StartHeight, m.Peers)
}
func (m *MethodProposeInput) Decode(payload []byte) (startHeight uint64, peers *Peers, err error) {
	if err = utils.UnpackMethod(ABI, MethodPropose, m, payload); err != nil {
		return
	}
	startHeight = m.StartHeight
	err = rlp.DecodeBytes(m.Peers, &peers)
	return
}

type MethodProposeOutput struct {
	Success bool
}

func (m *MethodProposeOutput) Encode(succeed bool) ([]byte, error) {
	m.Success = succeed
	return utils.PackOutputs(ABI, MethodPropose, m.Success)
}
func (m *MethodProposeOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodPropose, m, payload)
}

type MethodVoteInput struct {
	EpochID uint64
	Hash    []byte
}

func (m *MethodVoteInput) Encode(epochID uint64, hash common.Hash) ([]byte, error) {
	m.EpochID = epochID
	m.Hash = hash.Bytes()
	return utils.PackMethod(ABI, MethodVote, m.EpochID, m.Hash)
}
func (m *MethodVoteInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodVote, m, payload)
}

type MethodVoteOutput struct {
	Success bool
}

func (m *MethodVoteOutput) Encode(succeed bool) ([]byte, error) {
	m.Success = succeed
	return utils.PackOutputs(ABI, MethodVote, m.Success)
}
func (m *MethodVoteOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodVote, m, payload)
}

// useless input
//type MethodEpochInput struct{}
type MethodEpochOutput struct {
	EpochInfo []byte
}

func (m *MethodEpochOutput) Encode(epoch *EpochInfo) ([]byte, error) {
	enc, err := rlp.EncodeToBytes(epoch)
	if err != nil {
		return nil, err
	}
	m.EpochInfo = enc
	return utils.PackOutputs(ABI, MethodEpoch, m.EpochInfo)
}
func (m *MethodEpochOutput) Decode(payload []byte) (epoch *EpochInfo, err error) {
	if err = utils.UnpackOutputs(ABI, MethodEpoch, m, payload); err != nil {
		return
	}
	err = rlp.DecodeBytes(m.EpochInfo, &epoch)
	return
}

// useless input
//type MethodNextEpochInput struct{}
type MethodNextEpochOutput struct {
	EpochInfo []byte
}

func (m *MethodNextEpochOutput) Encode(epoch *EpochInfo) ([]byte, error) {
	enc, err := rlp.EncodeToBytes(epoch)
	if err != nil {
		return nil, err
	}
	m.EpochInfo = enc
	return utils.PackOutputs(ABI, MethodNextEpoch, m.EpochInfo)
}
func (m *MethodNextEpochOutput) Decode(payload []byte) (epoch *EpochInfo, err error) {
	if err = utils.UnpackOutputs(ABI, MethodNextEpoch, m, payload); err != nil {
		return
	}
	err = rlp.DecodeBytes(m.EpochInfo, &epoch)
	return
}

type EventProposed struct {
	EpochInfo []byte
}

func emitEventProposed(s *native.NativeContract, epoch *EpochInfo) error {
	enc, err := rlp.EncodeToBytes(epoch)
	if err != nil {
		return err
	}
	return s.AddNotify(ABI, []string{EventPropose}, enc)
}

type EventVoted struct {
	EpochID     uint64
	Hash        common.Hash
	VotedNumber uint64
	GroupSize   uint64
}

type EventEpochChanged struct {
	EpochInfo     []byte
	NextEpochInfo []byte
}
