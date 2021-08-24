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
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
)

// EventDrivenEngine implement event-driven hotstuff protocol, it obtains:
// 1.validator set which represent consensus participants
// 2.timer which used to

type EventDrivenEngine struct {
	valset hotstuff.ValidatorSet
}

func NewEventDrivenEngine(valset hotstuff.ValidatorSet) *EventDrivenEngine {
	return nil
}

// ProcessCertificates validate and handle QC/TC
func (e *EventDrivenEngine) ProcessCertificates(qc *hotstuff.QuorumCert) error {
	return nil
}

// ProcessProposal check proposal info and vote to the next leader if the proposal is valid
func (e *EventDrivenEngine) ProcessProposal(proposal *types.Block) error {
	return nil
}

// ProcessVoteMsg validate vote message and try to assemble qc
func (e *EventDrivenEngine) ProcessVoteMsg(vote *VoteMsg) error {
	return nil
}

// ProcessNewRoundEvent generate new proposal and broadcast to all validators.
func (e *EventDrivenEngine) ProcessNewRoundEvent() error {
	return nil
}
