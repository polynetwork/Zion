package core

import "github.com/ethereum/go-ethereum/consensus/hotstuff"

func (c *core) sendCommitVote() {

}

func (c *core) handleCommitVote(msg *message, src hotstuff.Validator) error {
	return nil
}

func (c *core) acceptCommitVote(msg *message, src hotstuff.Validator) error {
	// todo: add log
	return c.current.AddCommitVote(msg)
}