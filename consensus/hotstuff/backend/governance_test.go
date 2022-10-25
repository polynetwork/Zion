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

package backend

import (
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
)

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestFillHeader
func TestFillHeader(t *testing.T) {
	chain, engine := singleNodeChain()
	statedb, err := chain.State()
	if err != nil {
		t.Error(err)
	}
	parent := chain.CurrentBlock()
	header := makeHeader(parent)

	//header.Number = node_manager.GenesisBlockPerEpoch
	if err := engine.FillHeader(statedb, header); err != nil {
		t.Error(err)
	}

	if extra, err := types.ExtractHotstuffExtra(header); err != nil {
		t.Error(err)
	} else {
		t.Logf("start height %v, end height %v", extra.StartHeight, extra.EndHeight)
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestCheckPoint
func TestCheckPoint(t *testing.T) {

}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestValidators
func TestValidators(t *testing.T) {

}
