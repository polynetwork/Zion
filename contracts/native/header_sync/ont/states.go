/*
 * Copyright (C) 2021 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package ont

import (
	"io"
	"sort"

	"github.com/ethereum/go-ethereum/rlp"
)

type Peer struct {
	Index      uint32
	PeerPubkey string
}

func (this *Peer) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{this.Index, this.PeerPubkey})
}

func (this *Peer) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Index      uint32
		PeerPubkey string
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	this.Index = data.Index
	this.PeerPubkey = data.PeerPubkey
	return nil
}

type KeyHeights struct {
	HeightList []uint32
}

func (this *KeyHeights) EncodeRLP(w io.Writer) error {
	//first sort the list  (big -> small)
	sort.SliceStable(this.HeightList, func(i, j int) bool {
		return this.HeightList[i] > this.HeightList[j]
	})
	return rlp.Encode(w, []interface{}{this.HeightList})
}

func (this *KeyHeights) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		HeightList []uint32
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	this.HeightList = data.HeightList
	return nil
}

type ConsensusPeers struct {
	ChainID uint64
	Height  uint32
	PeerMap map[string]*Peer
}

func (this *ConsensusPeers) EncodeRLP(w io.Writer) error {
	var peerList []*Peer
	for _, v := range this.PeerMap {
		peerList = append(peerList, v)
	}
	sort.SliceStable(peerList, func(i, j int) bool {
		return peerList[i].PeerPubkey > peerList[j].PeerPubkey
	})

	return rlp.Encode(w, []interface{}{this.ChainID, this.Height, peerList})
}

func (this *ConsensusPeers) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ChainID  uint64
		Height   uint32
		PeerList []*Peer
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	peerMap := make(map[string]*Peer)
	for _, peer := range data.PeerList {
		peerMap[peer.PeerPubkey] = peer
	}

	this.ChainID = data.ChainID
	this.Height = data.Height
	this.PeerMap = peerMap
	return nil
}
