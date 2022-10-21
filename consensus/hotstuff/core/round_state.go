/*
 * Copyright (C) 2021 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */

package core

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

func (c *core) currentView() *View {
	return &View{
		Height: new(big.Int).Set(c.current.Height()),
		Round:  new(big.Int).Set(c.current.Round()),
	}
}

func (c *core) currentState() State {
	return c.current.State()
}

func (c *core) currentProposer() hotstuff.Validator {
	return c.valSet.GetProposer()
}

type roundState struct {
	vs hotstuff.ValidatorSet

	round  *big.Int
	height *big.Int
	state  State

	pendingRequest *Request
	proposal       hotstuff.Proposal // validator's prepare proposal
	proposalLocked bool

	// o(4n)
	newViews       *MessageSet // data set for newView message
	prepareVotes   *MessageSet // data set for prepareVote message
	preCommitVotes *MessageSet // data set for preCommitVote message
	commitVotes    *MessageSet // data set for commitVote message

	highQC      *QuorumCert // leader highQC
	prepareQC   *QuorumCert // prepareQC for repo and leader
	lockedQC    *QuorumCert // lockedQC for repo and pre-committedQC for leader
	committedQC *QuorumCert // committedQC for repo and leader

	mu *sync.RWMutex // mutex for fields except message set.
}

// newRoundState creates a new roundState instance with the given view and validatorSet
func newRoundState(view *View, validatorSet hotstuff.ValidatorSet, prepareQC *QuorumCert) *roundState {
	rs := &roundState{
		vs:             validatorSet,
		round:          view.Round,
		height:         view.Height,
		state:          StateAcceptRequest,
		newViews:       NewMessageSet(validatorSet),
		prepareVotes:   NewMessageSet(validatorSet),
		preCommitVotes: NewMessageSet(validatorSet),
		commitVotes:    NewMessageSet(validatorSet),
		mu:             new(sync.RWMutex),
	}
	if prepareQC != nil {
		rs.prepareQC = prepareQC.Copy()
		rs.lockedQC = prepareQC.Copy()
		rs.committedQC = prepareQC.Copy()
	}
	return rs
}

func (s *roundState) Height() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.height
}

func (s *roundState) HeightU64() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.height.Uint64()
}

func (s *roundState) Round() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.round
}

func (s *roundState) View() *View {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &View{
		Round:  s.round,
		Height: s.height,
	}
}

func (s *roundState) SetState(state State) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state = state
}

func (s *roundState) State() State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.state
}

func (s *roundState) SetProposal(proposal hotstuff.Proposal) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.proposal = proposal
	s.proposalLocked = false
}

func (s *roundState) Proposal() hotstuff.Proposal {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.proposal
}

func (s *roundState) LockProposal() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.proposal != nil && !s.proposalLocked {
		s.proposalLocked = true
	}
}

func (s *roundState) SetPendingRequest(req *Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pendingRequest = req
}

func (s *roundState) PendingRequest() *Request {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.pendingRequest
}

func (s *roundState) Vote() *Vote {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.proposal == nil || s.proposal.Hash() == EmptyHash {
		return nil
	}

	return &Vote{
		View: &View{
			Round:  new(big.Int).Set(s.round),
			Height: new(big.Int).Set(s.height),
		},
		Digest: s.proposal.Hash(),
	}
}

func (s *roundState) SetHighQC(qc *QuorumCert) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.highQC = qc
}

func (s *roundState) HighQC() *QuorumCert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.highQC
}

func (s *roundState) SetPrepareQC(qc *QuorumCert) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.prepareQC = qc
}

func (s *roundState) PrepareQC() *QuorumCert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.prepareQC
}

func (s *roundState) SetPreCommittedQC(qc *QuorumCert) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lockedQC = qc
}

func (s *roundState) PreCommittedQC() *QuorumCert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lockedQC
}

func (s *roundState) SetCommittedQC(qc *QuorumCert) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.committedQC = qc
}

func (s *roundState) CommittedQC() *QuorumCert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.committedQC
}

// message set has it's own mutex, do not lock or unlock with roundState mutex.
func (s *roundState) AddNewViews(msg *Message) error {
	return s.newViews.Add(msg)
}

func (s *roundState) NewViewSize() int {
	return s.newViews.Size()
}

func (s *roundState) NewViews() []*Message {
	return s.newViews.Values()
}

func (s *roundState) AddPrepareVote(msg *Message) error {
	return s.prepareVotes.Add(msg)
}

func (s *roundState) PrepareVotes() []*Message {
	return s.prepareVotes.Values()
}

func (s *roundState) PrepareVoteSize() int {
	return s.prepareVotes.Size()
}

func (s *roundState) AddPreCommitVote(msg *Message) error {
	return s.preCommitVotes.Add(msg)
}

func (s *roundState) PreCommitVoteSize() int {
	return s.preCommitVotes.Size()
}

func (s *roundState) AddCommitVote(msg *Message) error {
	return s.commitVotes.Add(msg)
}

func (s *roundState) CommitVoteSize() int {
	return s.commitVotes.Size()
}
