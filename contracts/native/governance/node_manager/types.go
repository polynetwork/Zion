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
	"math"
	"strings"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

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

func (m *Peers) Len() int {
	if m == nil || m.List == nil {
		return 0
	}
	return len(m.List)
}

func (m *Peers) Less(i, j int) bool {
	return strings.Compare(m.List[i].Address.Hex(), m.List[j].Address.Hex()) < 0
}

func (m *Peers) Swap(i, j int) {
	m.List[i], m.List[j] = m.List[j], m.List[i]
}

func (m *Peers) Copy() *Peers {
	enc, err := rlp.EncodeToBytes(m)
	if err != nil {
		return nil
	}
	var cp *Peers
	if err := rlp.DecodeBytes(enc, &cp); err != nil {
		return nil
	}
	return cp
}

type ProposalStatusType uint8

const (
	ProposalStatusUnknown ProposalStatusType = 0
	ProposalStatusPropose ProposalStatusType = 1
	ProposalStatusPassed  ProposalStatusType = 2
)

func (p ProposalStatusType) String() string {
	switch p {
	case ProposalStatusPropose:
		return "STATUS_PROPOSE"
	case ProposalStatusPassed:
		return "STATUS_PASSED"
	default:
		return "STATUS_UNKNOWN"
	}
}

type EpochInfo struct {
	ID          uint64
	Peers       *Peers
	StartHeight uint64
	Proposer    common.Address // hash generating without fields of `Proposer` and `Status`
	Status      ProposalStatusType

	hash atomic.Value
}

func (m *EpochInfo) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ID, m.Peers, m.StartHeight, m.Proposer, uint8(m.Status)})
}

func (m *EpochInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ID          uint64
		Peers       *Peers
		StartHeight uint64
		Proposer    common.Address
		Status      uint8
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.ID, m.Peers, m.StartHeight, m.Proposer, m.Status = data.ID, data.Peers, data.StartHeight, data.Proposer, ProposalStatusType(data.Status)
	return nil
}

func (m *EpochInfo) String() string {
	pstr := ""
	if m.Peers != nil && m.Peers.List != nil {
		for _, v := range m.Peers.List {
			pstr += fmt.Sprintf("peer: %s, pubkey: %s\r\n", v.Address.Hex(), v.PubKey)
		}
	}
	return fmt.Sprintf("epochHash:%s\r\nepochId: %d\r\n%sstartHeight: %d\r\nproposer:%s\r\nstatus:%s",
		m.Hash().Hex(), m.ID, pstr, m.StartHeight, m.Proposer.Hex(), m.Status.String())
}

func (m *EpochInfo) Hash() common.Hash {
	if hash := m.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	var inf = struct {
		ID          uint64
		Peers       *Peers
		StartHeight uint64
	}{
		ID:          m.ID,
		Peers:       m.Peers,
		StartHeight: m.StartHeight,
	}
	v := RLPHash(inf)
	m.hash.Store(v)
	return v
}

func (m *EpochInfo) Members() map[common.Address]struct{} {
	if m == nil || m.Peers == nil || m.Peers.List == nil || len(m.Peers.List) == 0 {
		return nil
	}
	data := make(map[common.Address]struct{})
	for _, v := range m.Peers.List {
		data[v.Address] = struct{}{}
	}
	return data
}

func (m *EpochInfo) MemberList() []common.Address {
	list := make([]common.Address, 0)
	if m == nil || m.Peers == nil || m.Peers.List == nil || len(m.Peers.List) == 0 {
		return list
	}
	for _, v := range m.Peers.List {
		list = append(list, v.Address)
	}
	return list
}

func (m *EpochInfo) QuorumSize() int {
	if m == nil || m.Peers == nil {
		return 0
	}
	total := m.Peers.Len()
	return int(math.Ceil(float64(2*total) / 3))
}

func (m *EpochInfo) OldMemberNum(peers *Peers) int {
	if m == nil || m.Peers == nil || m.Peers.List == nil || len(m.Peers.List) == 0 {
		return 0
	}
	if peers == nil || peers.List == nil || len(peers.List) == 0 {
		return 0
	}

	isOldMember := func(addr common.Address) bool {
		for _, v := range m.Peers.List {
			if v.Address == addr {
				return true
			}
		}
		return false
	}

	num := 0
	for _, v := range peers.List {
		if isOldMember(v.Address) {
			num += 1
		}
	}
	return num
}

type HashList struct {
	List []common.Hash
}

func (m *HashList) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.List})
}

func (m *HashList) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		List []common.Hash
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.List = data.List
	return nil
}

type AddressList struct {
	List []common.Address
}

func (m *AddressList) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.List})
}

func (m *AddressList) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		List []common.Address
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.List = data.List
	return nil
}

type ConsensusSign struct {
	Method string
	Input  []byte
	hash   atomic.Value
}

func (m *ConsensusSign) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Method, m.Input})
}
func (m *ConsensusSign) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Method string
		Input  []byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.Method, m.Input = data.Method, data.Input
	return nil
}
func (m *ConsensusSign) Hash() common.Hash {
	if hash := m.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	var inf = struct {
		Method string
		Input  []byte
	}{
		Method: m.Method,
		Input:  m.Input,
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
