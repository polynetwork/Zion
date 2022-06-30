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

package info_sync

import (
	"fmt"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/info_sync_abi"
)

var (
	MethodContractName  = info_sync_abi.MethodName
	MethodSyncRootInfo  = info_sync_abi.MethodSyncRootInfo
	MethodReplenish     = info_sync_abi.MethodReplenish
	MethodGetInfoHeight = info_sync_abi.MethodGetInfoHeight
	MethodGetInfo       = info_sync_abi.MethodGetInfo
)

var GasTable = map[string]uint64{
	MethodContractName:  0,
	MethodGetInfoHeight: 0,
	MethodGetInfo:       0,
	MethodSyncRootInfo:  1000,
	MethodReplenish:     1000,
}

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(info_sync_abi.InfoSyncABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

var ABI *abi.ABI

type GetInfoParam struct {
	ChainID uint64
	Height  uint32
}

func (m *GetInfoParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodGetInfo, m)
}

type GetInfoOutput struct {
	Info []byte
}

func (m *GetInfoOutput) Decode(payload []byte) error {
	if err := utils.UnpackOutputs(ABI, MethodGetInfo, m, payload); err != nil {
		return err
	}
	return nil
}

type GetInfoHeightParam struct {
	ChainID uint64
}

func (m *GetInfoHeightParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodGetInfoHeight, m)
}

type SyncRootInfoParam struct {
	ChainID   uint64
	RootInfos [][]byte
	Signature []byte
}

func (m *SyncRootInfoParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodSyncRootInfo, m)
}

func (m *SyncRootInfoParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ChainID, m.RootInfos, m.Signature})
}
func (m *SyncRootInfoParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ChainID   uint64
		RootInfos [][]byte
		Signature []byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.ChainID, m.RootInfos, m.Signature = data.ChainID, data.RootInfos, data.Signature
	return nil
}

//Digest Digest calculate the hash of param input
func (m *SyncRootInfoParam) Digest() ([]byte, error) {
	input := &SyncRootInfoParam{
		ChainID:   m.ChainID,
		RootInfos: m.RootInfos,
	}
	msg, err := rlp.EncodeToBytes(input)
	if err != nil {
		return nil, fmt.Errorf("SyncRootInfoParam, serialize input error: %v", err)
	}
	digest := crypto.Keccak256(msg)
	return digest, nil
}

type ReplenishParam struct {
	ChainID  uint64
	Heights  []uint32
}

func (m *ReplenishParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodReplenish, m)
}
