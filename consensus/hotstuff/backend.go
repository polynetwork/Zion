// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

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

	// Validators returns the validator set
	Validators() ValidatorSet

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

	GetProposal(hash common.Hash) Proposal

	// HasProposal checks if the combination of the given hash and height matches any existing blocks
	//HasProposal(hash common.Hash, number *big.Int) bool

	// GetProposer returns the proposer of the given block height
	//GetProposer(number uint64) common.Address

	// ParentValidators returns the validator set of the given proposal's parent block
	ParentValidators(proposal Proposal) ValidatorSet

	// HasBadBlock returns whether the block with the hash is a bad block
	HasBadProposal(hash common.Hash) bool

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

	// PrepareExtra generate header extra field with validator set
	PrepareExtra(header *types.Header, valSet ValidatorSet) ([]byte, error)

	// GetHeader get block header with hash and correct block height
	GetHeader(hash common.Hash, number uint64) *types.Header

	// SubscribeRequest notify to miner worker that event-driven engine need an new proposal
	SubscribeRequest(ch chan<- consensus.AskRequest) event.Subscription

	ResetValidators(valset ValidatorSet)
}

type HotstuffProtocol string

const (
	HOTSTUFF_PROTOCOL_BASIC        HotstuffProtocol = "basic"
	HOTSTUFF_PROTOCOL_EVENT_DRIVEN HotstuffProtocol = "event_driven"
)
