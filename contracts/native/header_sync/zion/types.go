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

package zion

import (
	"io"

	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type HeaderWithEpoch struct {
	Header *types.Header
	Epoch  *nm.EpochInfo
	Proof  []byte
}

func (h *HeaderWithEpoch) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{h.Header, h.Epoch, h.Proof})
}

func (h *HeaderWithEpoch) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Header *types.Header
		Epoch  *nm.EpochInfo
		Proof  []byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	h.Header, h.Epoch, h.Proof = data.Header, data.Epoch, data.Proof
	return nil
}

func (h *HeaderWithEpoch) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(h)
}

func (h *HeaderWithEpoch) Decode(payload []byte) error {
	return rlp.DecodeBytes(payload, h)
}
