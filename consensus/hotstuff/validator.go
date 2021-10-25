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
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Validator interface {
	// Address returns address
	Address() common.Address

	// String representation of Validator
	String() string
}

// ----------------------------------------------------------------------------

type Validators []Validator

func (slice Validators) Len() int {
	return len(slice)
}

func (slice Validators) Less(i, j int) bool {
	return strings.Compare(slice[i].String(), slice[j].String()) < 0
}

func (slice Validators) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ----------------------------------------------------------------------------

type ValidatorSet interface {
	// Calculate the proposer
	CalcProposer(lastProposer common.Address, round uint64)
	// Calculate the proposer with index
	CalcProposerByIndex(index uint64)
	// Return the validator size
	Size() int
	// Return the validator array
	List() []Validator
	// Return the validator address array
	AddressList() []common.Address
	// Get validator by index
	GetByIndex(i uint64) Validator
	// Get validator by given address
	GetByAddress(addr common.Address) (int, Validator)
	// Get current proposer
	GetProposer() Validator
	// Check whether the validator with given address is a proposer
	IsProposer(address common.Address) bool
	// Add validator
	AddValidator(address common.Address) bool
	// Remove validator
	RemoveValidator(address common.Address) bool
	// Copy validator set
	Copy() ValidatorSet
	// ParticipantsNumber calculate invalid validator size
	ParticipantsNumber(list []common.Address) int
	// CheckQuorum check committers
	CheckQuorum(committers []common.Address) error
	// Get the maximum number of faulty nodes
	F() int
	// Get the minimum number of quorum nodes
	Q() int
	// Get speaker policy
	Policy() SelectProposerPolicy
	// Cmp compare with another validator set, return false if not the same
	Cmp(src ValidatorSet) bool
}

// ----------------------------------------------------------------------------

type ProposalSelector func(ValidatorSet, common.Address, uint64) Validator
