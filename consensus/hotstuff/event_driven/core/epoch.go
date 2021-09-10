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

	"github.com/ethereum/go-ethereum/core/types"
)

//todo: epoch manager
func (c *core) initialize() error {
	c.smr.SetEpoch(0)
	c.smr.SetEpochStart(big.NewInt(0))
	c.smr.SetEpochEnd(big.NewInt(10000000))

	lastBlock, _ := c.backend.LastProposal()
	if lastBlock == nil {
		return fmt.Errorf("initialize event-driven engine with first block failed!")
	}
	salt, qc, err := extraProposal(lastBlock)
	if err != nil {
		return err
	}

	c.smr.SetHighCommitRound(salt.Round)
	c.smr.SetRound(c.smr.HighCommitRound())
	c.smr.SetHeight(lastBlock.Number())

	proposal := c.backend.GetProposal(lastBlock.Hash())
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

	c.blkPool = NewBlockPool(blktr)
	c.blkPool.AddQC(qc)
	if err := c.updateHighQCAndProposal(highQC, rootBlock); err != nil {
		c.logger.Trace("[Initialize Engine], update high qc and proposal", "err", err)
	}
	c.smr.SetLatestVoteRound(salt.Round)

	c.logger.Trace("[Initialize Engine]", "view", c.currentView())
	return nil
}
