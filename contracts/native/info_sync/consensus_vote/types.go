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

package consensus_vote

import (
	"io"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

type SignerList struct {
	StartHeight uint64
	EndHeight   uint64
	SignerList  []*SignerInfo
}

func (m *SignerList) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.StartHeight, m.EndHeight, m.SignerList})
}

func (m *SignerList) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		StartHeight uint64
		EndHeight   uint64
		SignerMap   []*SignerInfo
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.StartHeight = data.StartHeight
	m.EndHeight = data.EndHeight
	m.SignerList = data.SignerMap
	return nil
}

type SignerInfo struct {
	Address    common.Address
	SignHeight uint64
}

func (m *SignerInfo) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Address, m.SignHeight})
}

func (m *SignerInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Address    common.Address
		SignHeight uint64
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.Address = data.Address
	m.SignHeight = data.SignHeight
	return nil
}

type VoteMessage struct {
	Input []byte
	hash  atomic.Value
}

func (m *VoteMessage) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Input})
}
func (m *VoteMessage) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Input []byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.Input = data.Input
	return nil
}
func (m *VoteMessage) Hash() common.Hash {
	if hash := m.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	var inf = struct {
		Input []byte
	}{
		Input: m.Input,
	}
	v := RLPHash(inf)
	m.hash.Store(v)
	return v
}

type RootInfoUnique struct {
	ChainID uint64
	Height  uint32
	Info    []byte
}

func (d *RootInfoUnique) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{d.ChainID, d.Height, d.Info})
}

func (d *RootInfoUnique) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ChainID uint64
		Height  uint32
		Info    []byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	d.ChainID, d.Height, d.Info = data.ChainID, data.Height, data.Info
	return nil
}

func RLPHash(v interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, v)
	hw.Sum(h[:0])
	return h
}
