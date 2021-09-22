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
	"sort"

	"github.com/ethereum/go-ethereum/common"
)

type Status uint8

func (s *Status) Serialization(sink *common.ZeroCopySink) {
	sink.WriteUint8(uint8(*s))
}

func (s *Status) Deserialization(source *common.ZeroCopySource) error {
	status, eof := source.NextUint8()
	if eof {
		return fmt.Errorf("serialization.ReadUint8, deserialize status error: %v", io.ErrUnexpectedEOF)
	}
	*s = Status(status)
	return nil
}

type BlackListItem struct {
	PeerPubkey string         //peerPubkey in black list
	Address    common.Address //the owner of this peer
}

func (i *BlackListItem) Serialization(sink *common.ZeroCopySink) {
	sink.WriteString(i.PeerPubkey)
	sink.WriteVarBytes(i.Address[:])
}

func (i *BlackListItem) Deserialization(source *common.ZeroCopySource) error {
	peerPubkey, eof := source.NextString()
	if eof {
		return fmt.Errorf("source.NextString, deserialize peerPubkey error")
	}
	address, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("source.NextVarBytes, deserialize address error")
	}
	addr, err := common.AddressParseFromBytes(address)
	if err != nil {
		return fmt.Errorf("common.AddressParseFromBytes, deserialize address error: %s", err)
	}

	i.PeerPubkey = peerPubkey
	i.Address = addr
	return nil
}

type PeerPoolMap struct {
	PeerPoolMap map[string]*PeerPoolItem
}

func (m *PeerPoolMap) Serialization(sink *common.ZeroCopySink) {
	sink.WriteVarUint(uint64(len(m.PeerPoolMap)))
	var peerPoolItemList []*PeerPoolItem
	for _, v := range m.PeerPoolMap {
		peerPoolItemList = append(peerPoolItemList, v)
	}
	sort.SliceStable(peerPoolItemList, func(i, j int) bool {
		return peerPoolItemList[i].PeerPubkey > peerPoolItemList[j].PeerPubkey
	})
	for _, v := range peerPoolItemList {
		v.Serialization(sink)
	}
}

func (m *PeerPoolMap) Deserialization(source *common.ZeroCopySource) error {
	n, eof := source.NextVarUint()
	if eof {
		return fmt.Errorf("source.NextVarUint, deserialize PeerPoolMap length error")
	}
	peerPoolMap := make(map[string]*PeerPoolItem)
	for i := 0; uint64(i) < n; i++ {
		peerPoolItem := new(PeerPoolItem)
		if err := peerPoolItem.Deserialization(source); err != nil {
			return fmt.Errorf("deserialize peerPool error: %v", err)
		}
		peerPoolMap[peerPoolItem.PeerPubkey] = peerPoolItem
	}
	m.PeerPoolMap = peerPoolMap
	return nil
}

type PeerPoolItem struct {
	Index      uint64         //peer index
	PeerPubkey string         //peer pubkey
	Address    common.Address //peer owner
	Status     Status
}

func (i *PeerPoolItem) Serialization(sink *common.ZeroCopySink) {
	sink.WriteUint64(i.Index)
	sink.WriteString(i.PeerPubkey)
	sink.WriteVarBytes(i.Address[:])
	i.Status.Serialization(sink)
}

func (i *PeerPoolItem) Deserialization(source *common.ZeroCopySource) error {
	index, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("source.NextUint32, deserialize index error")
	}
	peerPubkey, eof := source.NextString()
	if eof {
		return fmt.Errorf("source.NextString, deserialize peerPubkey error")
	}
	address, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("source.NextVarBytes, deserialize address error")
	}
	status := new(Status)
	err := status.Deserialization(source)
	if err != nil {
		return fmt.Errorf("status.Deserialize. deserialize status error: %v", err)
	}
	addr := common.BytesToAddress(address)

	i.Index = index
	i.PeerPubkey = peerPubkey
	i.Address = addr
	i.Status = *status
	return nil
}

type GovernanceView struct {
	View   uint64
	Height uint64
	TxHash common.Hash
}

func (v *GovernanceView) Serialization(sink *common.ZeroCopySink) {
	sink.WriteUint64(v.View)
	sink.WriteUint64(v.Height)
	sink.WriteHash(v.TxHash)
}

func (v *GovernanceView) Deserialization(source *common.ZeroCopySource) error {
	view, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("source.NextUint32, deserialize view error")
	}
	height, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("source.NextUint32, deserialize height error")
	}
	txHash, eof := source.NextHash()
	if eof {
		return fmt.Errorf("source.NextHash, deserialize txHash error")
	}
	v.View = view
	v.Height = height
	v.TxHash = txHash
	return nil
}

type ConsensusSigns struct {
	SignsMap map[common.Address]bool
}

func (s *ConsensusSigns) Serialization(sink *common.ZeroCopySink) {
	sink.WriteVarUint(uint64(len(s.SignsMap)))
	var signsList []common.Address
	for k := range s.SignsMap {
		signsList = append(signsList, k)
	}
	sort.SliceStable(signsList, func(i, j int) bool {
		return signsList[i].Hex() > signsList[j].Hex()
	})
	for _, v := range signsList {
		sink.WriteVarBytes(v[:])
		sink.WriteBool(s.SignsMap[v])
	}
}

func (s *ConsensusSigns) Deserialization(source *common.ZeroCopySource) error {
	n, eof := source.NextVarUint()
	if eof {
		return fmt.Errorf("source.NextVarUint, deserialize length of signsMap error")
	}
	signsMap := make(map[common.Address]bool)
	for i := 0; uint64(i) < n; i++ {
		address, eof := source.NextVarBytes()
		if eof {
			return fmt.Errorf("source.NextVarBytes, deserialize address error")
		}
		v, eof := source.NextBool()
		if eof {
			return fmt.Errorf("source.NextBool, deserialize v error")
		}
		addr, err := common.AddressParseFromBytes(address)
		if err != nil {
			return fmt.Errorf("common.AddressParseFromBytes, deserialize address error")
		}
		signsMap[addr] = v
	}
	s.SignsMap = signsMap
	return nil
}

type Configuration struct {
	BlockMsgDelay        uint64
	HashMsgDelay         uint64
	PeerHandshakeTimeout uint64
	MaxBlockChangeView   uint64
}

func (c *Configuration) Serialization(sink *common.ZeroCopySink) {
	sink.WriteUint64(c.BlockMsgDelay)
	sink.WriteUint64(c.HashMsgDelay)
	sink.WriteUint64(c.PeerHandshakeTimeout)
	sink.WriteUint64(c.MaxBlockChangeView)
}

func (c *Configuration) Deserialization(source *common.ZeroCopySource) error {
	blockMsgDelay, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("source.NextUint32, deserialize blockMsgDelay error")
	}
	hashMsgDelay, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("source.NextUint32, deserialize hashMsgDelay error")
	}
	peerHandshakeTimeout, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("source.NextUint32, deserialize peerHandshakeTimeout error")
	}
	maxBlockChangeView, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("source.NextUint32, deserialize maxBlockChangeView error")
	}

	c.BlockMsgDelay = blockMsgDelay
	c.HashMsgDelay = hashMsgDelay
	c.PeerHandshakeTimeout = peerHandshakeTimeout
	c.MaxBlockChangeView = maxBlockChangeView
	return nil
}
