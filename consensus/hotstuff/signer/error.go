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

package signer

import "errors"

var (
	ErrInvalidSignature = errors.New("invalid signature")

	// ErrUnauthorized is returned if a header is signed by a non authorized entity.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrInvalidExtraDataFormat is returned when the extra data format is incorrect
	ErrInvalidExtraDataFormat = errors.New("invalid extra data format")

	// ErrInvalidCommittedSeals is returned if the committed seal is not signed by any of parent validators.
	ErrInvalidCommittedSeals = errors.New("invalid committed seals")

	// ErrEmptyCommittedSeals is returned if the field of committed seals is zero.
	ErrEmptyCommittedSeals = errors.New("zero committed seals")

	// ErrUnauthorizedAddress is returned when given address cannot be found in
	// current validator set.
	ErrUnauthorizedAddress = errors.New("unauthorized address")

	// ErrInvalidSigner is returned if the msg is unsigned
	ErrInvalidSigner = errors.New("message not signed by the sender")

	// ErrInvalidRawData is returned if the raw input is nil
	ErrInvalidRawData = errors.New("raw input is invalid")

	// ErrInvalidRawHash is returned if the raw hash is nil
	ErrInvalidRawHash = errors.New("raw hash is invalid")

	// ErrInvalidHeader is returned if the raw header is nil
	ErrInvalidHeader = errors.New("raw header is invalid")

	// ErrInvalidValset is returned if the validator set is nil
	ErrInvalidValset = errors.New("valset is nil")

	ErrNilQC = errors.New("qc is nil")

	// ErrInvalidQC is returned if the quorum cert is nil
	ErrInvalidQC = errors.New("qc is invalid")
)
