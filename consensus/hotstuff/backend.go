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

package hotstuff

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Backend provides application specific functions for Istanbul core
type Backend interface {
	// Address returns the owner's address
	Address() common.Address

	// Validators returns current epoch participants
	Validators(height uint64) ValidatorSet

	// EventMux returns the event mux in backend
	EventMux() *event.TypeMux

	// Broadcast sends a message to all validators (include self)
	Broadcast(valSet ValidatorSet, payload []byte) error

	// Gossip sends a message to all validators (exclude self)
	Gossip(valSet ValidatorSet, payload []byte) error

	// Unicast send a message to single peer
	Unicast(valSet ValidatorSet, payload []byte) error

	// PreCommit write seal to header and assemble new qc
	PreCommit(proposal Proposal, seals [][]byte) (Proposal, error)

	// ForwardCommit assemble unsealed block and sealed extra into an new full block
	ForwardCommit(proposal Proposal, extra []byte) (Proposal, error)

	// Commit delivers an approved proposal to backend.
	// The delivered proposal will be put into blockchain.
	Commit(proposal Proposal) error

	// Verify verifies the proposal. If a consensus.ErrFutureBlock error is returned,
	// the time difference of the proposal and current time is also returned.
	Verify(Proposal) (time.Duration, error)

	// Verify verifies the proposal. If a consensus.ErrFutureBlock error is returned,
	// the time difference of the proposal and current time is also returned.
	VerifyUnsealedProposal(Proposal) (time.Duration, error)

	// LastProposal retrieves latest committed proposal and the address of proposer
	LastProposal() (Proposal, common.Address)

	// HasBadBlock returns whether the block with the hash is a bad block
	HasBadProposal(hash common.Hash) bool

	// ValidateBlock execute block which contained in prepare message, and validate block state
	ValidateBlock(block *types.Block) error

	// AskMiningProposalWithParent ask for mining proposal with parent block as parameters.
	AskMiningProposalWithParent(parent *types.Block)

	Close() error
}

type CoreEngine interface {
	Start(chain consensus.ChainReader) error

	Stop() error

	// IsProposer return true if self address equal leader/proposer address in current round/height
	IsProposer() bool

	// verify if a hash is the same as the proposed block in the current pending request
	//
	// this is useful when the engine is currently the speaker
	//
	// pending request is populated right at the request stage so this would give us the earliest verification
	// to avoid any race condition of coming propagated blocks
	IsCurrentProposal(blockHash common.Hash) bool
}

type HotstuffProtocol string

const (
	HOTSTUFF_PROTOCOL_BASIC        HotstuffProtocol = "basic"
	HOTSTUFF_PROTOCOL_EVENT_DRIVEN HotstuffProtocol = "event_driven"
)
