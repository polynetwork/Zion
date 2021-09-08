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
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

// core implement event-driven hotstuff protocol, it obtains:
// 1.validator set which represent consensus participants
type core struct {
	config  *hotstuff.Config
	logger  log.Logger
	db      ethdb.Database
	chain   consensus.ChainReader
	backend hotstuff.Backend

	started bool
	address common.Address
	signer  hotstuff.Signer
	valset  hotstuff.ValidatorSet

	smr      *SMR
	messages *MessagePool
	blkPool  *BlockPool
	backlogs *backlog
	timer    *time.Timer // drive consensus round

	feed              event.Feed // request feed
	events            *event.TypeMuxSubscription
	finalCommittedSub *event.TypeMuxSubscription

	validateFn func([]byte, []byte) (common.Address, error)
}

// todo: chain
func New(
	backend hotstuff.Backend,
	c *hotstuff.Config,
	db ethdb.Database,
	signer hotstuff.Signer,
	valset hotstuff.ValidatorSet,
) hotstuff.CoreEngine {

	addr := signer.Address()
	engine := &core{
		config:  c,
		db:      db,
		backend: backend,
		logger:  log.New("address", addr),
	}

	engine.smr = newSMR()
	engine.address = addr
	engine.valset = valset
	engine.signer = signer
	engine.started = false
	engine.backlogs = newBackLog()
	engine.messages = NewMessagePool(valset)
	engine.validateFn = engine.checkValidatorSignature

	return engine
}

// handleNewRound proposer at this round get an new proposal and broadcast to all validators.
func (c *core) handleNewRound() error {
	if !c.started {
		return nil
	}

	logger := c.newSenderLogger("MSG_NEW_ROUND")
	logger.Trace("[New round]", "new_proposer", c.valset.GetProposer(), "valSet", c.valset.List(), "size", c.valset.Size(), "IsProposer", c.IsProposer())

	if !c.IsProposer() {
		return nil
	}

	c.processBacklog()
	return c.sendRequest()
}

func (c *core) handleTC(src hotstuff.Validator, data *hotstuff.Message) error {
	var (
		tc     *TimeoutCert
		msgTyp = MsgTypeTC
	)

	logger := c.newMsgLogger(msgTyp)

	if err := data.Decode(&tc); err != nil {
		logger.Trace("[Handle TC], failed to decode", "from", src.Address(), "err", err)
		return err
	}

	if tc == nil || tc.View == nil {
		logger.Trace("[Handle TC], invalid tc", "err", "tc is nil")
		return errInvalidTC
	}
	if tc.View.Cmp(c.currentView()) < 0 {
		return nil
	}

	if err := c.signer.VerifyCommittedSeal(c.valset, tc.Hash, tc.Seals); err != nil {
		logger.Trace("[Handle TC], failed to verify committed seal", "from", src.Address(), "err", err)
		return err
	}

	if err := c.advanceRoundByTC(tc, false); err != nil {
		if err == errOldMessage {
			return nil
		} else {
			logger.Trace("[Handle TC], failed to advance by tc", "from", src.Address(), "err", err)
			return err
		}
	}

	logger.Trace("[Handle TC], accept TC", "src", src.Address(), "tc", tc.Hash, "view", tc.View)
	return nil
}

func (c *core) commit3Chain() {
	if !c.IsProposer() {
		return
	}

	lockQC := c.smr.LockQC()
	if lockQC == nil {
		return
	}
	//
	//committedBlock := c.blkPool.GetCommitBlock(lockQC.Hash)
	//if committedBlock == nil {
	//	c.logger.Trace("[Commit 3-Chain], failed to get commit block", "lockQC view", lockQC.View)
	//	return
	//}
	//
	//round := lockQC.Round()
	//if existProposal := c.backend.GetProposal(committedBlock.Hash()); c.IsProposer() && existProposal == nil {
	//	if err := c.backend.Commit(committedBlock); err != nil {
	//		c.logger.Trace("[Commit 3-Chain], failed to commit", "err", err)
	//	} else {
	//		c.logger.Trace("[Commit 3-Chain], leader commit", "address", c.address, "hash", committedBlock.Hash(), "number", committedBlock.Number())
	//	}
	//}
	//
	//c.updateHighestCommittedRound(round)
	//puredBlocks := c.blkPool.Pure(committedBlock.Hash())
	//c.logger.Trace("[Commit 3-Chain], pured blocks", "hash lists", puredBlocks)

	branch := c.blkPool.GetCommitBranch(lockQC.Hash)
	if branch == nil {
		c.logger.Trace("[Commit 3-Chain], failed to get branch", "lockQC view", lockQC.View, "lockQC hash", lockQC.Hash)
		return
	}

	for _, block := range branch {
		if exist := c.backend.GetProposal(block.Hash()); exist != nil {
			continue
		}
		if err := c.backend.Commit(block); err != nil {
			c.logger.Trace("[Commit 3-Chain], failed to commit", "err", err)
		} else {
			c.logger.Trace("[Commit 3-Chain], leader commit", "address", c.address, "hash", block.Hash(), "number", block.Number())
		}
	}

	c.updateHighestCommittedRound(lockQC.Round())
	c.blkPool.Pure(lockQC.Hash)
}
