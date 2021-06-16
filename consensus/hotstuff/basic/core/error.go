package core

import "errors"

var (
	// errInconsistentSubject is returned when received subject is different from
	// current subject.
	errInconsistentSubject = errors.New("inconsistent subjects")

	// errNotFromProposer is returned when received message is supposed to be from
	// speaker.
	errNotFromSpeaker = errors.New("message does not come from speaker")

	// errIgnored is returned when a message was ignored.
	errIgnored = errors.New("message is ignored")

	// errFutureMessage is returned when current view is earlier than the
	// view of the received message.
	errFutureMessage = errors.New("future message")

	// errOldMessage is returned when the received message's view is earlier
	// than current view.
	errOldMessage = errors.New("old message")

	// errInvalidMessage is returned when the message is malformed.
	errInvalidMessage = errors.New("invalid message")

	// errFailedDecodeAnnounce is returned when the ANNOUNCE message is malformed.
	errFailedDecodeAnnounce = errors.New("failed to decode ANNOUNCE")

	// errFailedDecodeRequest is returned when the RESPONSE message is malformed.
	errFailedDecodeResponse = errors.New("failed to decode RESPONSE")

	// // errFailedDecodeCommit is returned when the COMMIT message is malformed.
	// errFailedDecodeCommit = errors.New("failed to decode COMMIT")

	// errFailedDecodeSendPub is returned when the SENDPUB message is malformed.
	errFailedDecodeSendPub = errors.New("failed to decode SENDPUB")

	// errFailedDecodeMessageSet is returned when the message set is malformed.
	// errFailedDecodeMessageSet = errors.New("failed to decode message set")

	// errInvalidSigner is returned when the message is signed by a validator different than message sender
	errInvalidSigner = errors.New("message not signed by the sender")

	// errInsufficientPub is returned when the local node has no enough aggregated public key collection
	errInsufficientPub = errors.New("not enough aggregated public key collection")
)
