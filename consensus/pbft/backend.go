package pbft

import "crypto/ecdsa"

type backend struct {
	privateKey *ecdsa.PrivateKey
	vs *VoteSet
}

func NewBackend() *backend {
	return nil
}

func (b *backend) enterNewHeight() error {
	return nil
}

func (b *backend) enterProposal() error {
	return nil
}

func (b *backend) enterPreVote() error {
	return nil
}

func (b *backend) enterPreVoteWait() error {
	return nil
}

func (b *backend) enterPreCommit() error {
	return nil
}

func (b *backend) enterCommit() error {
	return nil
}