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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type AddSignatureParam struct {
	Address     common.Address
	SideChainID uint64
	Subject     []byte
	Signature   []byte
}

func (p *AddSignatureParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{p.Address, p.SideChainID, p.Subject, p.Signature})
}
func (p *AddSignatureParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Address     common.Address
		SideChainID uint64
		Subject     []byte
		Signature   []byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	p.Address, p.SideChainID, p.Subject, p.Signature = data.Address, data.SideChainID, data.Subject, data.Signature
	return nil
}
