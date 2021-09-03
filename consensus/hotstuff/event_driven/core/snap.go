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

package core

import (
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type Snapshot struct {
	Epoch uint64

	CurrentRound       *big.Int
	CurrentHeight      *big.Int
	LastVoteRound      *big.Int
	HighestCommitRound *big.Int

	Blocks []*types.Block
}

// storeSnapshot generate snapshot with consensus info and marshal and store the structure.
func (e *core) storeSnapshot() ([]byte, error) {
	return nil, nil
}

// loadSnapshot load snapshot info from db and unmarshal to structure.
func (e *core) loadSnapshot() (*Snapshot, error) {
	return nil, nil
}