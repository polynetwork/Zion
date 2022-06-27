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

package common

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"io"

	"github.com/ethereum/go-ethereum/contracts/native"
)

const (
	//key prefix
	ROOT_INFO            = "rootInfo"
	CURRENT_HEIGHT       = "currentHeight"
	SYNC_ROOT_INFO_EVENT = "SyncRootInfoEvent"
)


type GetInfoParam struct {
	ChainID uint64
	Height  uint32
}

type GetInfoHeightParam struct {
	ChainID uint64
}

type SyncRootInfoParam struct {
	ChainID   uint64
	RootInfos [][]byte
}

func (m *SyncRootInfoParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ChainID, m.RootInfos})
}
func (m *SyncRootInfoParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ChainID   uint64
		RootInfos [][]byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.ChainID, m.RootInfos = data.ChainID, data.RootInfos
	return nil
}

type RootInfo struct {
	Height uint32
	Info   []byte
}

func (m *RootInfo) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Height, m.Info})
}
func (m *RootInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Height uint32
		Info   []byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.Height, m.Info = data.Height, data.Info
	return nil
}

func NotifyPutRootInfo(native *native.NativeContract, chainID uint64, height uint32) {
	err := native.AddNotify(ABI, []string{SYNC_ROOT_INFO_EVENT}, chainID, height, native.ContractRef().BlockHeight())
	if err != nil {
		panic(fmt.Sprintf("NotifyPutRootInfo failed: %v", err))
	}
}
