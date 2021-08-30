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
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// todo:
func (e *EventDrivenEngine) initialize() error {
	lastBlock, _ := e.backend.LastProposal()
	if lastBlock == nil {
		return fmt.Errorf("initialize event-driven engine with first block failed!")
	}

	//todo:
	e.epoch = 0
	e.epochHeightStart = big.NewInt(0)
	e.epochHeightEnd = big.NewInt(100)

	e.curHeight = new(big.Int).Add(lastBlock.Number(), common.Big1)

	salt, qc, err := extraProposal(lastBlock)
	if err != nil {
		return err
	}
	e.highestCommitRound = salt.Round
	e.curRound = new(big.Int).Add(e.highestCommitRound, common.Big1)

	if e.epochHeightStart.Cmp(e.highestCommitRound) > 0 {
		// todo
	}
	if e.highestCommitRound.Cmp(common.Big0) == 0 {
		proposal := e.backend.GetProposal(lastBlock.Hash())
		if proposal == nil {
			return fmt.Errorf("Can't get block %v", lastBlock.Hash())
		}
		rootBlock := proposal.(*types.Block)
		rootSalt, highQC, err := extraHeader(rootBlock.Header())
		if err != nil {
			return err
		}
		blktr, err := NewBlockTree(rootBlock, rootSalt.Round.Uint64(), 100)
		if err != nil {
			return err
		}
		e.blkPool = NewBlockPool(highQC, blktr)
	}

	// todo:
	e.lastVoteRound = salt.Round
	e.lockQC = qc

	return nil
}
