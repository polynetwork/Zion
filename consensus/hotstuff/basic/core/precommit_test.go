package core

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/stretchr/testify/assert"
)

func TestHandlePrepareVote(t *testing.T) {
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
			Code: MsgTypePrepareVote,
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
	}

	for _, v := range testcases {
		leader := v.Sys.getLeader()
		for src, vote := range v.Votes {
			assert.Equal(t, v.ExpectErr, leader.handlePrepareVote(vote, src))
		}
		if v.ExpectErr == nil {
			assert.Equal(t, StatePrepared, leader.current.State())
			assert.Equal(t, int(N), leader.current.PrepareVoteSize())
		}
	}
}

func TestHandlePreCommit(t *testing.T) {
	N := uint64(4)
	F := uint64(1)
	H := uint64(5)
	R := uint64(1)

	newPreCommitMsg := func(c *core) (hotstuff.Proposal, *hotstuff.QuorumCert) {
		coreView := c.currentView()
		h := coreView.Height.Uint64()
		r := coreView.Round.Uint64()
		return newProposalAndQC(c, h, r)
	}
	newP2PMsg := func(msg *hotstuff.QuorumCert) *message {
		payload, _ := Encode(msg)
		return &message{
			Code: MsgTypePreCommit,
			Msg:  payload,
		}
	}

	type testcase struct {
		Sys       *testSystem
		Msg       *message
		Leader    hotstuff.Validator
		ExpectErr error
	}
	testcases := []*testcase{
		// normal case
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			leader := sys.getLeader()
			val := validator.New(leader.Address())
			var (
				proposal hotstuff.Proposal
				qc       *hotstuff.QuorumCert
			)
			for _, backend := range sys.backends {
				core := backend.core()
				proposal, qc = newPreCommitMsg(core)
				core.current.SetProposal(proposal)
			}
			msg := newP2PMsg(qc)
			return &testcase{
				Sys:       sys,
				Msg:       msg,
				Leader:    val,
				ExpectErr: nil,
			}
		}(),

		// errOldMsg
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			leader := sys.getLeader()
			val := validator.New(leader.Address())
			var (
				proposal hotstuff.Proposal
				qc       *hotstuff.QuorumCert
			)
			for _, backend := range sys.backends {
				core := backend.core()
				proposal, qc = newPreCommitMsg(core)
				core.current.SetProposal(proposal)
			}
			qc.View.Height = new(big.Int).SetUint64(H - 1)
			msg := newP2PMsg(qc)
			return &testcase{
				Sys:       sys,
				Msg:       msg,
				Leader:    val,
				ExpectErr: errOldMessage,
			}
		}(),

		// errFutureMsg
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			leader := sys.getLeader()
			val := validator.New(leader.Address())
			var (
				proposal hotstuff.Proposal
				qc       *hotstuff.QuorumCert
			)
			for _, backend := range sys.backends {
				core := backend.core()
				proposal, qc = newPreCommitMsg(core)
				core.current.SetProposal(proposal)
			}
			qc.View.Round = new(big.Int).SetUint64(R + 1)
			msg := newP2PMsg(qc)
			return &testcase{
				Sys:       sys,
				Msg:       msg,
				Leader:    val,
				ExpectErr: errFutureMessage,
			}
		}(),

		// errNotFromProposal
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			var (
				proposal hotstuff.Proposal
				qc       *hotstuff.QuorumCert
			)
			for _, backend := range sys.backends {
				core := backend.core()
				proposal, qc = newPreCommitMsg(core)
				core.current.SetProposal(proposal)
			}
			msg := newP2PMsg(qc)
			val := validator.New(sys.getRepos()[0].Address())
			return &testcase{
				Sys:       sys,
				Msg:       msg,
				Leader:    val,
				ExpectErr: errNotFromProposer,
			}
		}(),

		// errState
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			var (
				proposal hotstuff.Proposal
				qc       *hotstuff.QuorumCert
			)
			for _, backend := range sys.backends {
				core := backend.core()
				proposal, qc = newPreCommitMsg(core)
				core.current.SetProposal(proposal)
				core.current.SetState(StatePrepared)
			}
			msg := newP2PMsg(qc)
			val := validator.New(sys.getLeader().Address())
			return &testcase{
				Sys:       sys,
				Msg:       msg,
				Leader:    val,
				ExpectErr: errState,
			}
		}(),
	}

	for _, c := range testcases {
		for _, backend := range c.Sys.backends {
			core := backend.core()
			assert.Equal(t, c.ExpectErr, core.handlePreCommit(c.Msg, c.Leader))
		}
	}
}
