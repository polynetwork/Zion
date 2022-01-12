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
package cosmos

import (
	"io"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/types"
)

type CosmosEpochSwitchInfo struct {
	// The height where validators set changed last time. Poly only accept
	// header and proof signed by new validators. That means the header
	// can not be lower than this height.
	Height uint64

	// Hash of the block at `Height`. Poly don't save the whole header.
	// So we can identify the content of this block by `BlockHash`.
	BlockHash bytes.HexBytes

	// The hash of new validators set which used to verify validators set
	// committed with proof.
	NextValidatorsHash bytes.HexBytes

	// The cosmos chain-id of this chain basing Cosmos-sdk.
	ChainID string
}

func (m *CosmosEpochSwitchInfo) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Height, m.BlockHash, m.NextValidatorsHash, m.ChainID})
}
func (m *CosmosEpochSwitchInfo) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		Height uint64
		BlockHash bytes.HexBytes
		NextValidatorsHash bytes.HexBytes
		ChainID string
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.Height, m.BlockHash, m.NextValidatorsHash, m.ChainID = data.Height, data.BlockHash, data.NextValidatorsHash, data.ChainID
	return nil
}

type CosmosHeader struct {
	Header  types.Header
	Commit  *types.Commit
	Valsets []*types.Validator
}
