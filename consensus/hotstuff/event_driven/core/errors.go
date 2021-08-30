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

import "errors"

var (
	// errNotFromProposer is returned when received Message is supposed to be from proposer.
	errNotFromProposer = errors.New("Message does not come from proposer")
	// errFutureMessage is returned when current view is earlier than the
	// view of the received Message.
	errFutureMessage        = errors.New("future Message")
	errFarAwayFutureMessage = errors.New("far away future Message")
	// errOldMessage is returned when the received Message's view is earlier
	// than current view.
	errOldMessage = errors.New("old Message")
	// errInvalidMessage is returned when the Message is malformed.
	errInvalidMessage = errors.New("invalid Message")
	// errFailedDecodeNewView is returned when the NEWVIEW Message is malformed.
	errFailedDecodeNewView = errors.New("failed to decode NEWVIEW")
	// errFailedDecodePrepare is returned when the PREPARE Message is malformed.
	errFailedDecodePrepare = errors.New("failed to decode PREPARE")
	errProposalConvert     = errors.New("failed to convert proposal to types.block")
	errExtraHeader         = errors.New("failed to extract header")
	errInvalidEpoch        = errors.New("invalid epoch")
	errInvalidVote         = errors.New("invalid vote")
	errInvalidQC           = errors.New("invalid qc")
	errInvalidHighQC       = errors.New("invalid highQC")

	// errInvalidSigner is returned when the Message is signed by a validator different than Message sender
	errInvalidSigner = errors.New("Message not signed by the sender")
)
