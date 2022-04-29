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
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
)

const contractName = "node manager"

func InitABI() {
	ab, err := abi.JSON(strings.NewReader(INodeManagerABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	ABI = &ab
}

var (
	ABI  *abi.ABI
	this = utils.NodeManagerContractAddress
)

type MethodContractNameInput struct{}

func (m *MethodContractNameInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodName)
}
func (m *MethodContractNameInput) Decode(payload []byte) error { return nil }

type MethodContractNameOutput struct {
	Name string
}

func (m *MethodContractNameOutput) Encode() ([]byte, error) {
	m.Name = contractName
	return utils.PackOutputs(ABI, MethodName, m.Name)
}
func (m *MethodContractNameOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodName, m, payload)
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
	EpochID   uint64
	EpochHash common.Hash
}

func (m *MethodVoteInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodVote, m.EpochID, m.EpochHash.Bytes())
}
func (m *MethodVoteInput) Decode(payload []byte) error {
	var data struct {
		EpochID   uint64
		EpochHash []byte
	}
	if err := utils.UnpackMethod(ABI, MethodVote, &data, payload); err != nil {
		return err
	}

	m.EpochID = data.EpochID
	m.EpochHash = common.BytesToHash(data.EpochHash)
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
type MethodEpochInput struct{}

func (m *MethodEpochInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodEpoch)
}
func (m *MethodEpochInput) Decode(payload []byte) error { return nil }

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
type MethodGetChangingEpochInput struct{}

func (m *MethodGetChangingEpochInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodGetChangingEpoch)
}
func (m *MethodGetChangingEpochInput) Decode(payload []byte) error { return nil }

type MethodGetEpochByIDInput struct {
	EpochID uint64
}

func (m *MethodGetEpochByIDInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodGetEpochByID, m.EpochID)
}
func (m *MethodGetEpochByIDInput) Decode(payload []byte) error {
	var data struct {
		EpochID uint64
	}
	if err := utils.UnpackMethod(ABI, MethodGetEpochByID, &data, payload); err != nil {
		return err
	}
	m.EpochID = data.EpochID
	return nil
}

type MethodProofInput struct {
	EpochID uint64
}

func (m *MethodProofInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodProof, m.EpochID)
}
func (m *MethodProofInput) Decode(payload []byte) error {
	var data struct {
		EpochID uint64
	}
	if err := utils.UnpackMethod(ABI, MethodProof, &data, payload); err != nil {
		return err
	}
	m.EpochID = data.EpochID
	return nil
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
	return s.AddNotify(ABI, []string{EventProposed}, enc)
}

func emitEventVoted(s *native.NativeContract, epochID uint64, hash common.Hash, curVotedNum int, groupSize int) error {
	return s.AddNotify(ABI, []string{EventVoted}, epochID, hash.Bytes(), uint64(curVotedNum), uint64(groupSize))
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
	return s.AddNotify(ABI, []string{EventEpochChanged}, curEnc, nextEnc)
}

func emitConsensusSign(s *native.NativeContract, sign *ConsensusSign, signer common.Address, num int) error {
	return s.AddNotify(ABI, []string{EventConsensusSigned}, sign.Method, sign.Input, signer, uint64(num))
}

type MethodGetEpochListJsonInput struct {
	EpochID uint64
}

func (m *MethodGetEpochListJsonInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodGetEpochListJson, m.EpochID)
}
func (m *MethodGetEpochListJsonInput) Decode(payload []byte) error {
	var data struct {
		EpochID uint64
	}
	if err := utils.UnpackMethod(ABI, MethodGetEpochListJson, &data, payload); err != nil {
		return err
	}
	m.EpochID = data.EpochID
	return nil
}

type MethodGetJsonOutput struct {
	Result string
}

func (m *MethodGetJsonOutput) Encode(methodName string) ([]byte, error) {
	return utils.PackOutputs(ABI, methodName, m.Result)
}

func (m *MethodGetJsonOutput) Decode(payload []byte, methodName string) error {
	return utils.UnpackOutputs(ABI, methodName, m, payload)
}
