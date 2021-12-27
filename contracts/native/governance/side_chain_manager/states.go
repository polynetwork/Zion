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
	return rlp.Encode(w, []interface{}{m.BindSignInfo})
}

func (m *BindSignInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		BindSignInfo map[string][]byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.BindSignInfo = data.BindSignInfo
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
