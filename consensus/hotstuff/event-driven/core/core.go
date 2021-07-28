package core

// core event-driven consensus protocol impl.
type core struct {
	// backend
}

type protocol interface {
	processProposal()

	processVote()

	processQC()

	processCertificates()

	processLocalTimeout()
}