package core

import "errors"

var (
	// errInconsistentSubject is returned when received subject is different from
	// current subject.
	errInconsistentSubject = errors.New("inconsistent subjects")

	// errNotFromProposer is returned when received message is supposed to be from proposer.
	errNotFromProposer = errors.New("message does not come from proposer")

	errNotToProposer = errors.New("message does not send to proposer")

	errNotToValidator = errors.New("message does not send to repo validator")

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

	// errFailedDecodeNewView is returned when the NEWVIEW message is malformed.
	errFailedDecodeNewView = errors.New("failed to decode NEWVIEW")

	// errFailedDecodePrepare is returned when the PREPARE message is malformed.
	errFailedDecodePrepare = errors.New("failed to decode PREPARE")
	errFailedDecodePrepareVote = errors.New("failed to decode PREPARE_VOTE")

	// errFailedDecodePreCommit is returned when the PRECOMMIT message is malformed.
	errFailedDecodePreCommit = errors.New("failed to decode PRECOMMIT")

	errFailedDecodePreCommitVote = errors.New("faild to decode PRECOMMIT_VOTE")

	// errFailedDecodeCommit is returned when the COMMIT message is malformed.
	errFailedDecodeCommit = errors.New("failed to decode COMMIT")

	errFailedDecodeCommitVote = errors.New("failed to decode COMMIT_VOTE")

	// errFailedDecodeMessageSet is returned when the message set is malformed.
	errFailedDecodeMessageSet = errors.New("failed to decode message set")

	// errInvalidSigner is returned when the message is signed by a validator different than message sender
	errInvalidSigner = errors.New("message not signed by the sender")

	errHashAlreayLocked = errors.New("proposal hash already locked")
)
