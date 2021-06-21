package core

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

type roundState struct {
	vs hotstuff.ValidatorSet

	round  *big.Int
	height *big.Int
	state  State

	proposal hotstuff.Proposal // leader's pending request

	newViews       *messageSet
	prepareVotes   *messageSet
	preCommitVotes *messageSet
	commitVotes    *messageSet

	highQC      *QuorumCert // leader highQC
	prepareQC   *QuorumCert // repo and leader's prepareQC
	lockedQC    *QuorumCert // repo's lockedQC or leader's pre-committed QC
	committedQC *QuorumCert // repo and leader's committedQC

	mtx *sync.RWMutex
}

// newRoundState creates a new roundState instance with the given view and validatorSet
func newRoundState(view *hotstuff.View, validatorSet hotstuff.ValidatorSet, prepareQC *QuorumCert) *roundState {
	return &roundState{
		vs:             validatorSet,
		round:          view.Round,
		height:         view.Height,
		state:          StateAcceptRequest,
		newViews:       newMessageSet(validatorSet),
		prepareVotes:   newMessageSet(validatorSet),
		preCommitVotes: newMessageSet(validatorSet),
		commitVotes:    newMessageSet(validatorSet),
		prepareQC:      prepareQC,
		mtx:            new(sync.RWMutex),
	}
}

func (s *roundState) Spawn(view *hotstuff.View) *roundState {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	nrs := new(roundState)
	nrs.vs = s.vs
	nrs.height = view.Height
	nrs.round = view.Round
	nrs.state = StateAcceptRequest

	nrs.newViews = newMessageSet(nrs.vs)
	nrs.prepareVotes = newMessageSet(nrs.vs)
	nrs.preCommitVotes = newMessageSet(nrs.vs)
	nrs.commitVotes = newMessageSet(nrs.vs)

	nrs.highQC = s.highQC
	nrs.prepareQC = s.prepareQC
	nrs.lockedQC = s.lockedQC
	nrs.committedQC = s.committedQC

	nrs.mtx = new(sync.RWMutex)

	return nrs
}

func (s *roundState) Height() *big.Int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.height
}

func (s *roundState) Round() *big.Int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.round
}

func (s *roundState) NextRound() *big.Int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return new(big.Int).Add(s.round, common.Big1)
}

func (s *roundState) View() *hotstuff.View {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return &hotstuff.View{
		Round:  s.round,
		Height: s.height,
	}
}

func (s *roundState) SetState(state State) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.state = state
}

func (s *roundState) State() State {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.state
}

func (s *roundState) SetProposal(proposal hotstuff.Proposal) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.proposal = proposal
}

func (s *roundState) Proposal() hotstuff.Proposal {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.proposal
}

func (s *roundState) Subject() *hotstuff.Subject {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if s.proposal.Hash() == EmptyHash {
		return nil
	}

	return &hotstuff.Subject{
		View: &hotstuff.View{
			Round:  new(big.Int).Set(s.round),
			Height: new(big.Int).Set(s.height),
		},
		Digest: s.proposal.Hash(),
	}
}

func (s *roundState) AddNewViews(msg *message) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.newViews.Add(msg)
}

func (s *roundState) NewViewSize() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.newViews.Size()
}

func (s *roundState) NewViews() []*message {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.newViews.Values()
}

func (s *roundState) AddPrepareVote(msg *message) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.prepareVotes.Add(msg)
}

func (s *roundState) PrepareVotes() []*message {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.prepareVotes.Values()
}

func (s *roundState) PrepareVoteSize() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.prepareVotes.Size()
}

func (s *roundState) AddPreCommitVote(msg *message) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.preCommitVotes.Add(msg)
}

func (s *roundState) PreCommitVoteSize() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.preCommitVotes.Size()
}

func (s *roundState) AddCommitVote(msg *message) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.commitVotes.Add(msg)
}

func (s *roundState) CommitVoteSize() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.commitVotes.Size()
}

func (s *roundState) SetHighQC(qc *QuorumCert) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.highQC = qc
}

func (s *roundState) HighQC() *QuorumCert {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.highQC
}

func (s *roundState) SetPrepareQC(qc *QuorumCert) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.prepareQC = qc
}

func (s *roundState) PrepareQC() *QuorumCert {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.prepareQC
}

func (s *roundState) SetLockedQC(qc *QuorumCert) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.lockedQC = qc
}

func (s *roundState) LockedQC() *QuorumCert {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.lockedQC
}

func (s *roundState) SetCommittedQC(qc *QuorumCert) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.committedQC = qc
}

func (s *roundState) CommittedQC() *QuorumCert {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.committedQC
}
