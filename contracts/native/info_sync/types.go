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

package info_sync

import (
	"github.com/ethereum/go-ethereum/rlp"
	"io"
)

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

type RootInfo struct {
	Height uint32
	Info   []byte
}

func (m *RootInfo) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Height, m.Info})
}
func (m *RootInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Height uint32
		Info   []byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.Height, m.Info = data.Height, data.Info
	return nil
}
