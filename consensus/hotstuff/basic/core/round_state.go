package core

import (
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"math/big"
)

type roundState struct {
	vs hotstuff.ValidatorSet

	round  *big.Int
	height *big.Int
	state  State

	pendingRequest *hotstuff.Request // leader's pending request
	proposal       hotstuff.Proposal // Address's prepare proposal
	proposalLocked bool

	// o(4n)
	newViews       *messageSet
	prepareVotes   *messageSet
	preCommitVotes *messageSet
	commitVotes    *messageSet

	highQC      *hotstuff.QuorumCert // leader highQC
	prepareQC   *hotstuff.QuorumCert // prepareQC for repo and leader
	lockedQC    *hotstuff.QuorumCert // lockedQC for repo and pre-committedQC for leader
	committedQC *hotstuff.QuorumCert // committedQC for repo and leader
}

// newRoundState creates a new roundState instance with the given view and validatorSet
func newRoundState(view *hotstuff.View, validatorSet hotstuff.ValidatorSet, prepareQC *hotstuff.QuorumCert) *roundState {
	rs := &roundState{
		vs:             validatorSet,
		round:          view.Round,
		height:         view.Height,
		state:          StateAcceptRequest,
		newViews:       newMessageSet(validatorSet),
		prepareVotes:   newMessageSet(validatorSet),
		preCommitVotes: newMessageSet(validatorSet),
		commitVotes:    newMessageSet(validatorSet),
	}
	if prepareQC != nil {
		rs.prepareQC = prepareQC.Copy()
		rs.lockedQC = prepareQC.Copy()
		rs.committedQC = prepareQC.Copy()
	}
	return rs
}

func (s *roundState) Height() *big.Int {
	return s.height
}

func (s *roundState) Round() *big.Int {
	return s.round
}

func (s *roundState) View() *hotstuff.View {
	return &hotstuff.View{
		Round:  s.round,
		Height: s.height,
	}
}

func (s *roundState) SetState(state State) {
	s.state = state
}

func (s *roundState) State() State {
	return s.state
}

func (s *roundState) SetProposal(proposal hotstuff.Proposal) {
	s.proposal = proposal
}

func (s *roundState) Proposal() hotstuff.Proposal {
	return s.proposal
}

func (s *roundState) LockProposal() {
	if s.proposal != nil && !s.proposalLocked {
		s.proposalLocked = true
	}
}

func (s *roundState) UnLockProposal() {
	if s.proposal != nil && s.proposalLocked {
		s.proposalLocked = false
		s.proposal = nil
	}
}

func (s *roundState) IsProposalLocked() bool {
	return s.proposalLocked
}

func (s *roundState) LastLockedProposal() (bool, hotstuff.Proposal) {
	return s.proposalLocked, s.proposal
}

func (s *roundState) SetPendingRequest(req *hotstuff.Request) {
	s.pendingRequest = req
}

func (s *roundState) PendingRequest() *hotstuff.Request {
	return s.pendingRequest
}

func (s *roundState) Vote() *Vote {
	if s.proposal == nil || s.proposal.Hash() == EmptyHash {
		return nil
	}

	return &Vote{
		View: &hotstuff.View{
			Round:  new(big.Int).Set(s.round),
			Height: new(big.Int).Set(s.height),
		},
		Digest: s.proposal.Hash(),
	}
}

// AddNewViews all valid message, and invalid message would be ignore
func (s *roundState) AddNewViews(msg *message) error {
	return s.newViews.Add(msg)
}

func (s *roundState) NewViewSize() int {
	return s.newViews.Size()
}

func (s *roundState) NewViews() []*message {
	return s.newViews.Values()
}

func (s *roundState) AddPrepareVote(msg *message) error {
	return s.prepareVotes.Add(msg)
}

func (s *roundState) PrepareVotes() []*message {
	return s.prepareVotes.Values()
}

func (s *roundState) PrepareVoteSize() int {
	return s.prepareVotes.Size()
}

func (s *roundState) AddPreCommitVote(msg *message) error {
	return s.preCommitVotes.Add(msg)
}

func (s *roundState) PreCommitVoteSize() int {
	return s.preCommitVotes.Size()
}

func (s *roundState) AddCommitVote(msg *message) error {
	return s.commitVotes.Add(msg)
}

func (s *roundState) CommitVoteSize() int {
	return s.commitVotes.Size()
}

func (s *roundState) SetHighQC(qc *hotstuff.QuorumCert) {
	s.highQC = qc
}

func (s *roundState) HighQC() *hotstuff.QuorumCert {
	return s.highQC
}

func (s *roundState) SetPrepareQC(qc *hotstuff.QuorumCert) {
	s.prepareQC = qc
}

func (s *roundState) PrepareQC() *hotstuff.QuorumCert {
	return s.prepareQC
}

func (s *roundState) SetPreCommittedQC(qc *hotstuff.QuorumCert) {
	s.lockedQC = qc
}

func (s *roundState) PreCommittedQC() *hotstuff.QuorumCert {
	return s.lockedQC
}

func (s *roundState) SetCommittedQC(qc *hotstuff.QuorumCert) {
	s.committedQC = qc
}

func (s *roundState) CommittedQC() *hotstuff.QuorumCert {
	return s.committedQC
}
