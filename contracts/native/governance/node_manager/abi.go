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
	MethodProof        = "proof"
	MethodNextEpoch    = "nextEpoch"

	EventPropose     = "proposed"
	EventVote        = "voted"
	EventEpochChange = "epochChanged"
)

const abijson = `[
    {"type":"function","name":"` + MethodContractName + `","inputs":[],"outputs":[{"internalType":"string","name":"Name","type":"string"}],"stateMutability":"nonpayable"},
	{"type":"function","name":"` + MethodPropose + `","inputs":[{"internalType":"uint64","name":"StartHeight","type":"uint64"},{"internalType":"bytes","name":"Peers","type":"bytes"}],"outputs":[{"internalType":"bool","name":"Success","type":"bool"}],"stateMutability":"nonpayable"},
    {"type":"function","name":"` + MethodVote + `","inputs":[{"internalType":"uint64","name":"EpochID","type":"uint64"},{"internalType":"bytes","name":"Hash","type":"bytes"}],"outputs":[{"internalType":"bool","name":"Success","type":"bool"}],"stateMutability":"nonpayable"},
	{"type":"function","name":"` + MethodEpoch + `","inputs":[],"outputs":[{"internalType":"bytes","name":"Epoch","type":"bytes"}],"stateMutability":"nonpayable"},
	{"type":"function","name":"` + MethodNextEpoch + `","inputs":[],"outputs":[{"internalType":"bytes","name":"Epoch","type":"bytes"}],"stateMutability":"nonpayable"},
	{"type":"function","name":"` + MethodProof + `","inputs":[],"outputs":[{"internalType":"bytes","name":"Hash","type":"bytes"}],"stateMutability":"nonpayable"},
    {"type":"event","name":"` + EventPropose + `","anonymous":false,"inputs":[{"internalType":"bytes","name":"Epoch","type":"bytes"}]},
	{"type":"event","name":"` + EventVote + `","anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"EpochID","type":"uint64"},{"indexed":false,"internalType":"bytes","name":"Hash","type":"bytes"},{"indexed":false,"internalType":"uint64","name":"VotedNumber","type":"uint64"},{"indexed":false,"internalType":"uint64","name":"GroupSize","type":"uint64"}]},
	{"type":"event","name":"` + EventEpochChange + `","anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"Epoch","type":"bytes"},{"indexed":false,"internalType":"bytes","name":"NextEpoch","type":"bytes"}]}
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

func (m *MethodContractNameOutput) Encode() ([]byte, error) {
	m.Name = contractName
	return utils.PackOutputs(ABI, MethodContractName, m.Name)
}
func (m *MethodContractNameOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodContractName, m, payload)
}

type MethodProposeInput struct {
	StartHeight uint64
	Peers       *Peers
}

func (m *MethodProposeInput) Encode() ([]byte, error) {
	enc, err := rlp.EncodeToBytes(m.Peers)
	if err != nil {
		return nil, err
	}
	return utils.PackMethod(ABI, MethodPropose, m.StartHeight, enc)
}
func (m *MethodProposeInput) Decode(payload []byte) error {
	var data struct {
		StartHeight uint64
		Peers       []byte
	}
	if err := utils.UnpackMethod(ABI, MethodPropose, &data, payload); err != nil {
		return err
	}
	m.StartHeight = data.StartHeight
	return rlp.DecodeBytes(data.Peers, &m.Peers)
}

type MethodProposeOutput struct {
	Success bool
}

func (m *MethodProposeOutput) Encode() ([]byte, error) {
	return utils.PackOutputs(ABI, MethodPropose, m.Success)
}
func (m *MethodProposeOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodPropose, m, payload)
}

type MethodVoteInput struct {
	EpochID uint64
	Hash    common.Hash
}

func (m *MethodVoteInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodVote, m.EpochID, m.Hash.Bytes())
}
func (m *MethodVoteInput) Decode(payload []byte) error {
	var data struct {
		EpochID uint64
		Hash    []byte
	}
	if err := utils.UnpackMethod(ABI, MethodVote, &data, payload); err != nil {
		return err
	}

	m.EpochID = data.EpochID
	m.Hash = common.BytesToHash(data.Hash)
	return nil
}

type MethodVoteOutput struct {
	Success bool
}

func (m *MethodVoteOutput) Encode() ([]byte, error) {
	return utils.PackOutputs(ABI, MethodVote, m.Success)
}
func (m *MethodVoteOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodVote, m, payload)
}

// useless input
//type MethodEpochInput struct{}
type MethodEpochOutput struct {
	Epoch *EpochInfo
}

func (m *MethodEpochOutput) Encode() ([]byte, error) {
	enc, err := rlp.EncodeToBytes(m.Epoch)
	if err != nil {
		return nil, err
	}
	return utils.PackOutputs(ABI, MethodEpoch, enc)
}
func (m *MethodEpochOutput) Decode(payload []byte) error {
	var data struct {
		Epoch []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodEpoch, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.Epoch, &m.Epoch)
}

// useless input
//type MethodNextEpochInput struct{}
type MethodNextEpochOutput struct {
	Epoch *EpochInfo
}

func (m *MethodNextEpochOutput) Encode() ([]byte, error) {
	enc, err := rlp.EncodeToBytes(m.Epoch)
	if err != nil {
		return nil, err
	}
	return utils.PackOutputs(ABI, MethodNextEpoch, enc)
}
func (m *MethodNextEpochOutput) Decode(payload []byte) error {
	var data struct {
		Epoch []byte
	}
	if err := utils.UnpackOutputs(ABI, MethodNextEpoch, &data, payload); err != nil {
		return err
	}
	return rlp.DecodeBytes(data.Epoch, &m.Epoch)
}

type MethodProofOutput struct {
	Hash common.Hash
}

func (m *MethodProofOutput) Encode() ([]byte, error) {
	return utils.PackOutputs(ABI, MethodProof, m.Hash.Bytes())
}
func (m *MethodProofOutput) Decode(payload []byte) error {
	var data common.Hash
	if err := utils.UnpackOutputs(ABI, MethodProof, &data, payload); err != nil {
		return err
	}
	m.Hash = data
	return nil
}

func emitEventProposed(s *native.NativeContract, epoch *EpochInfo) error {
	enc, err := rlp.EncodeToBytes(epoch)
	if err != nil {
		return err
	}
	return s.AddNotify(ABI, []string{EventPropose}, enc)
}

func emitEventVoted(s *native.NativeContract, epochID uint64, hash common.Hash, curVotedNum int, groupSize int) error {
	return s.AddNotify(ABI, []string{EventVote}, epochID, hash, uint64(curVotedNum), uint64(groupSize))
}

func emitEpochChange(s *native.NativeContract, curEpoch, nextEpoch *EpochInfo) error {
	curEnc, err := rlp.EncodeToBytes(curEpoch)
	if err != nil {
		return err
	}
	nextEnc, err := rlp.EncodeToBytes(nextEpoch)
	if err != nil {
		return err
	}
	return s.AddNotify(ABI, []string{EventEpochChange}, curEnc, nextEnc)
}
