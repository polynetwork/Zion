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
	errInvalidSignature = errors.New("invalid signature")

	// errUnauthorized is returned if a header is signed by a non authorized entity.
	errUnauthorized = errors.New("unauthorized")

	// errInvalidExtraDataFormat is returned when the extra data format is incorrect
	errInvalidExtraDataFormat = errors.New("invalid extra data format")

	// errInvalidCommittedSeals is returned if the committed seal is not signed by any of parent validators.
	errInvalidCommittedSeals = errors.New("invalid committed seals")

	// errEmptyCommittedSeals is returned if the field of committed seals is zero.
	errEmptyCommittedSeals = errors.New("zero committed seals")

	// errUnauthorizedAddress is returned when given address cannot be found in
	// current validator set.
	errUnauthorizedAddress = errors.New("unauthorized address")

	// errInvalidSigner is returned if the msg is unsigned
	errInvalidSigner = errors.New("message not signed by the sender")

	// errInvalidRawData is returned if the raw input is nil
	errInvalidRawData = errors.New("raw input is invalid")

	// errInvalidRawHash is returned if the raw hash is nil
	errInvalidRawHash = errors.New("raw hash is invalid")

	// errInvalidHeader is returned if the raw header is nil
	errInvalidHeader = errors.New("raw header is invalid")

	// errInvalidValset is returned if the validator set is nil
	errInvalidValset = errors.New("valset is nil")

	// errInvalidQC is returned if the quorum cert is nil
	errInvalidQC = errors.New("qc is nil")
)
