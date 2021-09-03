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
	engine.messages = NewMessagePool(valset)
	engine.validateFn = engine.checkValidatorSignature

	return engine
}

// handleNewRound proposer at this round get an new proposal and broadcast to all validators.
func (e *core) handleNewRound() error {
	if !e.started {
		return nil
	}

	logger := e.newLogger()
	msgTyp := MsgTypeNewRound
	logger.Trace("New round", "msg", msgTyp, "new_proposer", e.valset.GetProposer(), "valSet", e.valset.List(), "size", e.valset.Size(), "IsProposer", e.IsProposer())

	if !e.IsProposer() {
		return nil
	} else {
		return e.sendRequest()
	}
}

func (e *core) handleTC(src hotstuff.Validator, data *hotstuff.Message) error {
	logger := e.newLogger()

	var (
		tc     *TimeoutCert
		msgTyp = MsgTypeTC
	)
	if err := data.Decode(&tc); err != nil {
		logger.Trace("Failed to decode", "msg", msgTyp, "from", src.Address(), "err", err)
		return err
	}

	if err := e.signer.VerifyCommittedSeal(e.valset, tc.Hash, tc.Seals); err != nil {
		logger.Trace("Failed to verify committed seal", "msg", msgTyp, "from", src.Address(), "err", err)
		return err
	}

	if err := e.advanceRoundByTC(tc, false); err != nil {
		logger.Trace("Failed to advance by tc", "msg", msgTyp, "from", src.Address(), "err", err)
		return err
	}

	logger.Trace("Accept TC", "msg", msgTyp, "src", src.Address(), "tc", tc.Hash, "view", tc.View)
	return nil
}

func (e *core) commit3Chain() {
	highQC := e.smr.HighQC()
	lockQC := e.smr.LockQC()
	if highQC == nil || lockQC == nil {
		return
	}

	committedBlock := e.blkPool.GetCommitBlock(highQC.Hash, lockQC.Hash)
	if committedBlock == nil {
		return
	}

	// todo: 如果节点此时宕机怎么办？还是说允许所有的节点一起提交区块
	if existProposal := e.backend.GetProposal(committedBlock.Hash()); existProposal == nil {
		if e.isSelf(committedBlock.Coinbase()) {
			e.backend.Commit(committedBlock)
		}
	}

	e.blkPool.Pure(committedBlock.Hash())
}

//
//func (e *core) handleQC(src hotstuff.Validator, data *hotstuff.Message) error {
//	logger := e.newLogger()
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
