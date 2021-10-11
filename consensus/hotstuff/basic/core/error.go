package core

import "errors"

var (
	// ErrUnauthorizedAddress is returned when given address cannot be found in
	// current validator set.
	ErrUnauthorizedAddress = errors.New("unauthorized address")
	// errInconsistentVote is returned when received subject is different from
	// current subject.
	errInconsistentVote = errors.New("inconsistent vote")
	errInvalidDigest    = errors.New("invalid digest")
	// errNotFromProposer is returned when received Message is supposed to be from proposer.
	errNotFromProposer = errors.New("Message does not come from proposer")
	errNotToProposer   = errors.New("Message does not send to proposer")
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
	errFailedDecodePrepare     = errors.New("failed to decode PREPARE")
	errFailedDecodePrepareVote = errors.New("failed to decode PREPARE_VOTE")
	// errFailedDecodePreCommit is returned when the PRECOMMIT Message is malformed.
	errFailedDecodePreCommit     = errors.New("failed to decode PRECOMMIT")
	errFailedDecodePreCommitVote = errors.New("faild to decode PRECOMMIT_VOTE")
	// errFailedDecodeCommit is returned when the COMMIT Message is malformed.
	errFailedDecodeCommit     = errors.New("failed to decode COMMIT")
	errFailedDecodeCommitVote = errors.New("failed to decode COMMIT_VOTE")
	// errInvalidSigner is returned when the Message is signed by a validator different than Message sender
	errInvalidSigner          = errors.New("Message not signed by the sender")
	errState                  = errors.New("error state")
	errNoRequest              = errors.New("no valid request")
	errInvalidProposal        = errors.New("invalid proposal")
	errVerifyUnsealedProposal = errors.New("verify unsealed proposal failed")
	errExtend                 = errors.New("proposal extend relationship error")
	errSafeNode               = errors.New("safeNode checking failed")
	errAddNewViews            = errors.New("add new view error")
	errAddPrepareVote         = errors.New("add prepare vote error")
	errAddPreCommitVote       = errors.New("add pre commit vote error")
	errBadEpochValidators     = errors.New("last epoch validator set is empty")
)
