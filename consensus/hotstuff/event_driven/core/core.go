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
	logger.Trace("New round", "new_proposer", c.valset.GetProposer(), "valSet", c.valset.List(), "size", c.valset.Size(), "IsProposer", c.IsProposer())

	if !c.IsProposer() {
		return nil
	}

	c.processBacklog()
	//if err := c.forwardProposal(); err != nil {
	//	return c.sendRequest()
	//}
	//return nil
	return c.sendRequest()
}

func (c *core) handleTC(src hotstuff.Validator, data *hotstuff.Message) error {
	var (
		tc     *TimeoutCert
		msgTyp = MsgTypeTC
	)

	logger := c.newMsgLogger(msgTyp)

	if err := data.Decode(&tc); err != nil {
		logger.Trace("Failed to decode", "from", src.Address(), "err", err)
		return err
	}

	if tc == nil || tc.View == nil {
		logger.Trace("Invalid tc", "err", "tc is nil")
		return errInvalidTC
	}
	if tc.View.Cmp(c.currentView()) < 0 {
		return nil
	}

	if err := c.signer.VerifyCommittedSeal(c.valset, tc.Hash, tc.Seals); err != nil {
		logger.Trace("Failed to verify committed seal", "from", src.Address(), "err", err)
		return err
	}

	if err := c.advanceRoundByTC(tc, false); err != nil {
		if err == errOldMessage {
			return nil
		} else {
			logger.Trace("Failed to advance by tc", "from", src.Address(), "err", err)
			return err
		}
	}

	logger.Trace("Accept TC", "src", src.Address(), "tc", tc.Hash, "view", tc.View)
	return nil
}

func (c *core) commit3Chain() {
	lockQC := c.smr.LockQC()
	if lockQC == nil {
		return
	}

	c.logger.Trace("Try to Commit 3-chain block")

	committedBlock := c.blkPool.GetCommitBlock(lockQC.Hash)
	if committedBlock == nil {
		c.logger.Trace("Failed to get commit block","lockQC view", lockQC.View)
		return
	}

	c.backend.Commit(committedBlock)
	c.logger.Trace("Commit 3-chain block", "hash", committedBlock.Hash(), "number", committedBlock.Number())

	//// todo: 如果节点此时宕机怎么办？还是说允许所有的节点一起提交区块
	//if existProposal := c.backend.GetProposal(committedBlock.Hash()); existProposal == nil {
	//	//if c.isSelf(committedBlock.Coinbase()) {
	//	//
	//	//}
	//
	//} else {
	//	c.logger.Trace("block already synced to chain reader", "")
	//}

	c.blkPool.Pure(committedBlock.Hash())
}

//
//func (e *core) handleQC(src hotstuff.Validator, data *hotstuff.Message) error {
//	logger := e.newMsgLogger()
//
//	var (
//		qc     *hotstuff.QuorumCert
//		msgTyp = MsgTypeQC
//	)
//	if err := data.Decode(&qc); err != nil {
//		logger.Trace("Failed to decode", "msg", msgTyp, "from", src.Address(), "err", err)
//		return err
//	}
//
//	if err := e.signer.VerifyQC(qc, e.valset); err != nil {
//		logger.Trace("Failed to verify qc", "msg", msgTyp, "from", src.Address(), "err", err)
//		return err
//	}
//
//	if err := e.processQC(qc); err != nil {
//		logger.Trace("Failed to process qc", "msg", msgTyp, "from", src.Address(), "err", err)
//		return err
//	}
//
//	logger.Trace("Accept QC", "msg", msgTyp, "src", src.Address(), "qc", qc.Hash, "view", qc.View)
//	return nil
//}
