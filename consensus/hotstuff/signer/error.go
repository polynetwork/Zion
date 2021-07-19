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
)
