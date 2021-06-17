package core

import "github.com/ethereum/go-ethereum/core/types"

func (c *core) sendDecide() {
	if !c.IsProposer() {
		return
	}

	if proposal := c.current.Proposal(); proposal != nil {
		committedSeals := make([][]byte, c.current.CommitVoteSize())
		for i, v := range c.current.commitVotes.Values() {
			committedSeals[i] = make([]byte, types.HotstuffExtraSeal)
			copy(committedSeals[i][:], v.CommittedSeal[:])
		}

		if err := c.backend.Commit(proposal, committedSeals); err != nil {
			c.current.UnlockHash() //Unlock block when insertion fails
			c.sendNextRoundChange()
			return
		}
	}
}
