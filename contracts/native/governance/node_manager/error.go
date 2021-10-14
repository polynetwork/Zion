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

package node_manager

import "errors"

var (
	ErrInvalidEpoch = errors.New("invalid epoch")

	ErrProposalNotExist = errors.New("proposal not exist")

	ErrEpochNotExist = errors.New("epoch not exist")

	ErrEpochProofNotExist = errors.New("epoch proof not exist")

	ErrConsensusSignNotExist = errors.New("consensus sign not exist")

	ErrInvalidAuthority = errors.New("invalid authority")

	ErrInvalidInput = errors.New("decode input params failed")

	ErrInvalidPeers = errors.New("invalid peers")

	ErrInvalidSign = errors.New("sign invalid")

	ErrInvalidPubKey = errors.New("invalid public key")

	ErrDuplicateProposal = errors.New("duplicate proposal")

	ErrDuplicateVote = errors.New("duplicate vote")

	ErrDuplicateSigner = errors.New("duplicate signer")

	ErrProposalsNum = errors.New("proposals number out of range")

	ErrPeersNum = errors.New("proposal peers out of range")

	ErrProposalPassed = errors.New("proposal already passed")

	ErrOldParticipantsNumber = errors.New("old participants should >= 2/3")

	ErrProposalStartHeight = errors.New("proposal start height invalid")

	ErrVoteHeight = errors.New("too late to vote")

	ErrStorage = errors.New("store key value failed")

	ErrEmitLog = errors.New("emit log failed")
)
