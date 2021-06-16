package core

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

// newRoundState creates a new roundState instance with the given view and validatorSet
func newRoundState(view *hotstuff.View, validatorSet hotstuff.ValidatorSet, pendingRequest *hotstuff.Request, hasBadProposal func(hash common.Hash) bool) *roundState {
	return &roundState{
		round:  view.Round,
		height: view.Height,

		pendingRequest: pendingRequest,

		newViews:       newMessageSet(validatorSet),
		prepareVotes:   newMessageSet(validatorSet),
		preCommitVotes: newMessageSet(validatorSet),
		commitVotes:    newMessageSet(validatorSet),

		mtx:            new(sync.RWMutex),
		hasBadProposal: hasBadProposal,
	}
}

type roundState struct {
	round  *big.Int
	height *big.Int
	state  State

	pendingRequest *hotstuff.Request

	newViews       *messageSet
	prepareVotes   *messageSet
	preCommitVotes *messageSet
	commitVotes    *messageSet

	highQC      *QuorumCert
	prepareQC   *QuorumCert
	preCommitQC *QuorumCert
	lockedQC    *QuorumCert
	commitQC    *QuorumCert

	mtx            *sync.RWMutex
	hasBadProposal func(hash common.Hash) bool
}

func (r *roundState) SetState(st State) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if st != r.state {
		r.state = st
	}
}

func (r *roundState) SetPendingRequest(req *hotstuff.Request) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.pendingRequest = req
}

func (r *roundState) PendingRequest() *hotstuff.Request {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.pendingRequest
}

func (r *roundState) SetHighQC(qc *QuorumCert) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.highQC = qc
}

func (r *roundState) HighQC() *QuorumCert {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.highQC
}

func (r *roundState) SetPrepareQC(qc *QuorumCert) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.prepareQC = qc
}

func (r *roundState) PrepareQC() *QuorumCert {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.preCommitQC
}

func (r *roundState) SetPreCommitQC(qc *QuorumCert) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.preCommitQC = qc
}

func (r *roundState) PreCommitQC() *QuorumCert {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.preCommitQC
}

func (r *roundState) SetLockedQC(qc *QuorumCert) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.lockedQC = qc
}

func (r *roundState) LockedQC() *QuorumCert {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.lockedQC
}

func (r *roundState) SetCommitQC(qc *QuorumCert) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.commitQC = qc
}

func (r *roundState) CommitQC() *QuorumCert {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.commitQC
}

func (r *roundState) Height() *big.Int {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.height
}

func (r *roundState) Round() *big.Int {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.round
}

func (r *roundState) View() *hotstuff.View {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return &hotstuff.View{
		Round:  r.round,
		Height: r.height,
	}
}

func (r *roundState) Proposal() hotstuff.Proposal {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if r.lockedQC.Proposal != nil {
		return r.lockedQC.Proposal
	}
	return nil
}

func (r *roundState) AddNewView(msg *message) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.newViews.Add(msg)
}

func (r *roundState) NewViewSize() int {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.newViews.Size()
}

func (r *roundState) AddPrepareVote(msg *message) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.prepareVotes.Add(msg)
}

func (r *roundState) PrepareVoteSize() int {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.prepareVotes.Size()
}

func (r *roundState) AddPreCommitVote(msg *message) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.preCommitVotes.Add(msg)
}

func (r *roundState) PreCommitVoteSize() int {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.preCommitVotes.Size()
}

func (r *roundState) AddCommitVote(msg *message) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.commitVotes.Add(msg)
}

func (r *roundState) CommitVoteSize() int {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.commitVotes.Size()
}
