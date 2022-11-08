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

package ripple

import (
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"sort"
)

type MultisignInfo struct {
	Status bool
	SigMap map[string]bool
}

func (this *MultisignInfo) EncodeRLP(w io.Writer) error {
	sigList := make([]string, 0, len(this.SigMap))
	for k := range this.SigMap {
		sigList = append(sigList, k)
	}
	sort.SliceStable(sigList, func(i, j int) bool {
		return sigList[i] > sigList[j]
	})
	return rlp.Encode(w, []interface{}{this.Status, sigList})
}

func (this *MultisignInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Status  bool
		SigList []string
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	this.Status = data.Status

	sigMap := make(map[string]bool, len(data.SigList))
	for _, v := range data.SigList {
		sigMap[v] = true
	}
	this.SigMap = sigMap

	return nil
}

type Signer struct {
	Account       []byte
	TxnSignature  []byte
	SigningPubKey []byte
}
