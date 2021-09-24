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
	"io"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

const PeerInfoLength int = 91

type PeerInfo struct {
	PubKey  string
	Address common.Address
}

func (m *PeerInfo) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.PubKey, m.Address})
}

func (m *PeerInfo) DecodeRLP(s *rlp.Stream) error {
	var peer struct {
		PubKey  string
		Address common.Address
	}

	if err := s.Decode(&peer); err != nil {
		return err
	}
	m.PubKey, m.Address = peer.PubKey, peer.Address
	return nil
}

func (m *PeerInfo) String() string {
	return fmt.Sprintf("{Address: %s PubKey: %s}", m.Address.Hex(), m.PubKey)
}

type Peers struct {
	List []*PeerInfo
}

func (m *Peers) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.List})
}

func (m *Peers) DecodeRLP(s *rlp.Stream) error {
	var peers struct {
		List []*PeerInfo
	}

	if err := s.Decode(&peers); err != nil {
		return err
	}
	m.List = peers.List
	return nil
}

type EpochInfo struct {
	ID          uint64
	Peers       *Peers
	StartHeight uint64
	EndHeight   uint64

	hash atomic.Value
}

func (m *EpochInfo) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ID, m.Peers, m.StartHeight, m.EndHeight})
}

func (m *EpochInfo) DecodeRLP(s *rlp.Stream) error {
	var inf struct {
		ID          uint64
		Peers       *Peers
		StartHeight uint64
		EndHeight   uint64
	}

	if err := s.Decode(&inf); err != nil {
		return err
	}
	m.ID, m.Peers, m.StartHeight, m.EndHeight = inf.ID, inf.Peers, inf.StartHeight, inf.EndHeight
	return nil
}

func (m *EpochInfo) String() string {
	return fmt.Sprintf("{ID: %d Start Height: %d}", m.ID, m.StartHeight)
}

func (m *EpochInfo) Hash() common.Hash {
	if hash := m.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	var inf = struct {
		ID          uint64
		Peers       *Peers
		StartHeight uint64
		EndHeight   uint64
	}{
		ID:          m.ID,
		Peers:       m.Peers,
		StartHeight: m.StartHeight,
		EndHeight:   m.EndHeight,
	}
	v := RLPHash(inf)
	m.hash.Store(v)
	return v
}

func RLPHash(v interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, v)
	hw.Sum(h[:0])
	return h
}
