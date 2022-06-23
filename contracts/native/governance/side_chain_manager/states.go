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

package side_chain_manager

import (
	"fmt"
	"io"

	ethcomm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type SideChain struct {
	Address      ethcomm.Address
	ChainId      uint64
	Router       uint64
	Name         string
	BlocksToWait uint64
	CCMCAddress  []byte
	ExtraInfo    []byte
}

func (m *SideChain) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Address, m.ChainId, m.Router, m.Name, m.BlocksToWait, m.CCMCAddress, m.ExtraInfo})
}

func (m *SideChain) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Address      ethcomm.Address
		ChainId      uint64
		Router       uint64
		Name         string
		BlocksToWait uint64
		CCMCAddress  []byte
		ExtraInfo    []byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.Address, m.ChainId, m.Router, m.Name, m.BlocksToWait, m.CCMCAddress, m.ExtraInfo =
		data.Address, data.ChainId, data.Router, data.Name, data.BlocksToWait, data.CCMCAddress, data.ExtraInfo
	return nil
}

type BindSignInfo struct {
	BindSignInfo map[string][]byte
}

func (m *BindSignInfo) EncodeRLP(w io.Writer) error {
	keys := make([]string, 0)
	values := make([][]byte, 0)

	if m.BindSignInfo == nil || len(m.BindSignInfo) == 0 {
		return fmt.Errorf("invalid BindSignInfo")
	}

	for k, v := range m.BindSignInfo {
		if v == nil {
			return fmt.Errorf("BindSignInfo value can be empty slice but not nil")
		}
		keys = append(keys, k)
		values = append(values, v)
	}
	return rlp.Encode(w, []interface{}{keys, values})
}

func (m *BindSignInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Keys   []string
		Values [][]byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	if data.Keys == nil || data.Values == nil ||
		len(data.Keys) == 0 || len(data.Keys) != len(data.Values) {
		return fmt.Errorf("invalid bindSignInfo")
	}

	m.BindSignInfo = make(map[string][]byte)
	for i := 0; i < len(data.Keys); i++ {
		m.BindSignInfo[data.Keys[i]] = data.Values[i]
	}
	return nil
}

type ContractBinded struct {
	Contract []byte
	Ver      uint64
}

func (m *ContractBinded) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Contract, m.Ver})
}

func (m *ContractBinded) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Contract []byte
		Ver      uint64
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.Contract, m.Ver = data.Contract, data.Ver
	return nil
}
