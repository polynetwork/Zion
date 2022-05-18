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
	CROSS_CHAIN_INFO            = "crossChainInfo"
	SYNC_CROSS_CHAIN_INFO_EVENT = "syncCrossChainInfo"
)

type SyncCrossChainInfoParam struct {
	ChainID         uint64
	CrossChainInfos []*CrossChainInfo
}

func (m *SyncCrossChainInfoParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ChainID, m.CrossChainInfos})
}
func (m *SyncCrossChainInfoParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ChainID         uint64
		CrossChainInfos []*CrossChainInfo
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.ChainID, m.CrossChainInfos = data.ChainID, data.CrossChainInfos
	return nil
}

type CrossChainInfo struct {
	Key   []byte
	Value []byte
}

func (m *CrossChainInfo) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Key, m.Value})
}
func (m *CrossChainInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Key   []byte
		Value []byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.Key, m.Value = data.Key, data.Value
	return nil
}

func NotifyPutCrossChainInfo(native *native.NativeContract, chainID uint64, key, value []byte) {
	err := native.AddNotify(ABI, []string{SYNC_CROSS_CHAIN_INFO_EVENT}, chainID, key, value, native.ContractRef().BlockHeight())
	if err != nil {
		panic(fmt.Sprintf("NotifyPutCrossChainInfo failed: %v", err))
	}
}
