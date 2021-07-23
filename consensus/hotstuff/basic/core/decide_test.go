package core

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/stretchr/testify/assert"
)

func TestHandleCommitVote(t *testing.T) {
	N := uint64(4)
	F := uint64(1)
	H := uint64(5)
	R := uint64(1)

	newVote := func(c *core, hash common.Hash) *Vote {
		view := c.currentView()
		return &Vote{
			View:   view,
			Digest: hash,
		}
	}
	newVoteMsg := func(vote *Vote) *message {
		payload, _ := Encode(vote)
		return &message{
			Code: MsgTypeCommitVote,
			Msg:  payload,
		}
	}

	type testcase struct {
		Sys       *testSystem
		Votes     map[hotstuff.Validator]*message
		ExpectErr error
	}

	testcases := []*testcase{

		// normal case
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			proposal := makeBlock(int64(H))
			votes := make(map[hotstuff.Validator]*message)
			for _, v := range sys.backends {
				core := v.core()
				core.current.SetProposal(proposal)
				core.current.SetPreCommittedQC(&hotstuff.QuorumCert{Hash: proposal.Hash()})

				vote := newVote(core, proposal.Hash())
				msg := newVoteMsg(vote)
				msg.Address = core.Address()
				val := validator.New(msg.Address)

				votes[val] = msg
			}
			return &testcase{
				Sys:       sys,
				Votes:     votes,
				ExpectErr: nil,
			}
		}(),

		// errOldMessage
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			proposal := makeBlock(int64(H))
			votes := make(map[hotstuff.Validator]*message)
			for _, v := range sys.backends {
				core := v.core()
				core.current.SetProposal(proposal)
				core.current.SetPreCommittedQC(&hotstuff.QuorumCert{Hash: proposal.Hash()})

				vote := newVote(core, proposal.Hash())
				vote.View.Height = new(big.Int).SetUint64(H - 1)

				msg := newVoteMsg(vote)
				msg.Address = core.Address()
				val := validator.New(msg.Address)

				votes[val] = msg
			}
			return &testcase{
				Sys:       sys,
				Votes:     votes,
				ExpectErr: errOldMessage,
			}
		}(),

		// errFutureMessage
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			proposal := makeBlock(int64(H))
			votes := make(map[hotstuff.Validator]*message)
			for _, v := range sys.backends {
				core := v.core()
				core.current.SetProposal(proposal)
				core.current.SetPreCommittedQC(&hotstuff.QuorumCert{Hash: proposal.Hash()})

				vote := newVote(core, proposal.Hash())
				vote.View.Round = new(big.Int).SetUint64(R + 1)

				msg := newVoteMsg(vote)
				msg.Address = core.Address()
				val := validator.New(msg.Address)

				votes[val] = msg
			}
			return &testcase{
				Sys:       sys,
				Votes:     votes,
				ExpectErr: errFutureMessage,
			}
		}(),

		// errInconsistentVote
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			proposal := makeBlock(int64(H))
			votes := make(map[hotstuff.Validator]*message)
			for _, v := range sys.backends {
				core := v.core()
				core.current.SetProposal(proposal)
				core.current.SetPreCommittedQC(&hotstuff.QuorumCert{Hash: proposal.Hash()})

				vote := newVote(core, proposal.Hash())
				vote.Digest = common.HexToHash("0x1234")
				msg := newVoteMsg(vote)
				msg.Address = core.Address()
				val := validator.New(msg.Address)

				votes[val] = msg
			}
			return &testcase{
				Sys:       sys,
				Votes:     votes,
				ExpectErr: errInconsistentVote,
			}
		}(),

		// errInvalidDigest
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			proposal := makeBlock(int64(H))
			votes := make(map[hotstuff.Validator]*message)
			for _, v := range sys.backends {
				core := v.core()
				core.current.SetProposal(proposal)
				core.current.SetPreCommittedQC(&hotstuff.QuorumCert{Hash: common.HexToHash("0x124")})

				vote := newVote(core, proposal.Hash())
				msg := newVoteMsg(vote)
				msg.Address = core.Address()
				val := validator.New(msg.Address)

				votes[val] = msg
			}
			return &testcase{
				Sys:       sys,
				Votes:     votes,
				ExpectErr: errInvalidDigest,
			}
		}(),
	}

	for _, v := range testcases {
		leader := v.Sys.getLeader()
		for src, vote := range v.Votes {
			assert.Equal(t, v.ExpectErr, leader.handleCommitVote(vote, src))
		}
		if v.ExpectErr == nil {
			assert.Equal(t, StateCommitted, leader.current.State())
			assert.Equal(t, int(N), leader.current.CommitVoteSize())
		}
	}
}
