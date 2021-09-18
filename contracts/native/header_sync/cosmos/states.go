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
	"fmt"

	"github.com/polynetwork/poly/common"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/types"
)

type CosmosEpochSwitchInfo struct {
	// The height where validators set changed last time. Poly only accept
	// header and proof signed by new validators. That means the header
	// can not be lower than this height.
	Height int64

	// Hash of the block at `Height`. Poly don't save the whole header.
	// So we can identify the content of this block by `BlockHash`.
	BlockHash bytes.HexBytes

	// The hash of new validators set which used to verify validators set
	// committed with proof.
	NextValidatorsHash bytes.HexBytes

	// The cosmos chain-id of this chain basing Cosmos-sdk.
	ChainID string
}

func (info *CosmosEpochSwitchInfo) Serialization(sink *common.ZeroCopySink) {
	sink.WriteInt64(info.Height)
	sink.WriteVarBytes(info.BlockHash)
	sink.WriteVarBytes(info.NextValidatorsHash)
	sink.WriteString(info.ChainID)
}

func (info *CosmosEpochSwitchInfo) Deserialization(source *common.ZeroCopySource) error {
	var eof bool
	info.Height, eof = source.NextInt64()
	if eof {
		return fmt.Errorf("deserialize height of CosmosEpochSwitchInfo failed")
	}
	info.BlockHash, eof = source.NextVarBytes()
	if eof {
		return fmt.Errorf("deserialize BlockHash of CosmosEpochSwitchInfo failed")
	}
	info.NextValidatorsHash, eof = source.NextVarBytes()
	if eof {
		return fmt.Errorf("deserialize NextValidatorsHash of CosmosEpochSwitchInfo failed")
	}
	info.ChainID, eof = source.NextString()
	if eof {
		return fmt.Errorf("deserialize ChainID of CosmosEpochSwitchInfo failed")
	}
	return nil
}

type CosmosHeader struct {
	Header  types.Header
	Commit  *types.Commit
	Valsets []*types.Validator
}
