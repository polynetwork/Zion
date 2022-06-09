/*
 * Copyright (C) 2022 The poly network Authors
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
package signature_manager

import (
	"io"
	"sort"

	"github.com/ethereum/go-ethereum/rlp"
)

type SigInfo struct {
	Status  bool
	SigInfo []Signature
	m       map[string][]byte
}

type Signature struct {
	Addr    string
	Content []byte
}

func (this *SigInfo) init(full bool) {
	this.m = make(map[string][]byte)
	for _, sig := range this.SigInfo {
		this.m[sig.Addr] = sig.Content
	}

	if !full {
		return
	}

	sigInfoList := make([]Signature, 0, len(this.m))
	for k, v := range this.m {
		sigInfoList = append(sigInfoList, Signature{Addr: k, Content: v})
	}
	sort.SliceStable(sigInfoList, func(i, j int) bool {
		return sigInfoList[i].Addr > sigInfoList[j].Addr
	})

	this.SigInfo = sigInfoList
}

func (this *SigInfo) EncodeRLP(w io.Writer) error {

	sigInfoList := make([]Signature, 0, len(this.m))
	for k, v := range this.m {
		sigInfoList = append(sigInfoList, Signature{Addr: k, Content: v})
	}
	sort.SliceStable(sigInfoList, func(i, j int) bool {
		return sigInfoList[i].Addr > sigInfoList[j].Addr
	})
	return rlp.Encode(w, []interface{}{this.Status, sigInfoList})
}

func (this *SigInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Status  bool
		SigInfo []Signature
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	this.Status, this.SigInfo = data.Status, data.SigInfo

	this.init(false)

	return nil
}
