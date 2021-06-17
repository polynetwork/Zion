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

		//newViews:       newMessageSet(validatorSet),
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

	prepare    *MsgPrepare
	lockedHash common.Hash

	prepareVotes   *messageSet
	preCommitVotes *messageSet
	commitVotes    *messageSet

	//highQC      *QuorumCert
	//prepareQC   *QuorumCert
	//preCommitQC *QuorumCert
	//lockedQC    *QuorumCert
	//commitQC    *QuorumCert

	mtx            *sync.RWMutex
	hasBadProposal func(hash common.Hash) bool
}

//func (r *roundState) SetState(st State) {
//	r.mtx.Lock()
//	defer r.mtx.Unlock()
//
//	if st != r.state {
//		r.state = st
//	}
//}

func (s *roundState) Subject() *hotstuff.Subject {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if s.prepare == nil {
		return nil
	}

	return &hotstuff.Subject{
		View: &hotstuff.View{
			Round:  new(big.Int).Set(s.round),
			Height: new(big.Int).Set(s.height),
		},
		Digest: s.prepare.Proposal.Hash(),
	}
}

func (s *roundState) SetPrepare(prepare *MsgPrepare) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.prepare = prepare
}

func (s *roundState) Proposal() hotstuff.Proposal {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if s.prepare != nil {
		return s.prepare.Proposal
	}

	return nil
}

func (s *roundState) SetPendingRequest(req *hotstuff.Request) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.pendingRequest = req
}

func (s *roundState) PendingRequest() *hotstuff.Request {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.pendingRequest
}

//
//func (r *roundState) SetHighQC(qc *QuorumCert) {
//	r.mtx.Lock()
//	defer r.mtx.Unlock()
//	r.highQC = qc
//}
//
//func (r *roundState) HighQC() *QuorumCert {
//	r.mtx.RLock()
//	defer r.mtx.RUnlock()
//	return r.highQC
//}
//
//func (r *roundState) SetPrepareQC(qc *QuorumCert) {
//	r.mtx.Lock()
//	defer r.mtx.Unlock()
//	r.prepareQC = qc
//}
//
//func (r *roundState) PrepareQC() *QuorumCert {
//	r.mtx.RLock()
//	defer r.mtx.RUnlock()
//	return r.preCommitQC
//}
//
//func (r *roundState) SetPreCommitQC(qc *QuorumCert) {
//	r.mtx.Lock()
//	defer r.mtx.Unlock()
//	r.preCommitQC = qc
//}
//
//func (r *roundState) PreCommitQC() *QuorumCert {
//	r.mtx.RLock()
//	defer r.mtx.RUnlock()
//	return r.preCommitQC
//}
//
//func (r *roundState) SetLockedQC(qc *QuorumCert) {
//	r.mtx.Lock()
//	defer r.mtx.Unlock()
//
//	r.lockedQC = qc
//}
//
//func (r *roundState) LockedQC() *QuorumCert {
//	r.mtx.RLock()
//	defer r.mtx.RUnlock()
//	return r.lockedQC
//}
//
//func (r *roundState) SetCommitQC(qc *QuorumCert) {
//	r.mtx.Lock()
//	defer r.mtx.Unlock()
//
//	r.commitQC = qc
//}
//
//func (r *roundState) CommitQC() *QuorumCert {
//	r.mtx.RLock()
//	defer r.mtx.RUnlock()
//	return r.commitQC
//}

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

func (s *roundState) View() *hotstuff.View {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return &hotstuff.View{
		Round:  s.round,
		Height: s.height,
	}
}

//func (r *roundState) Proposal() hotstuff.Proposal {
//	r.mtx.RLock()
//	defer r.mtx.RUnlock()
//	if r.lockedQC.Proposal != nil {
//		return r.lockedQC.Proposal
//	}
//	return nil
//}
//
//func (r *roundState) AddNewView(msg *message) error {
//	r.mtx.Lock()
//	defer r.mtx.Unlock()
//	return r.newViews.Add(msg)
//}
//
//func (r *roundState) NewViewSize() int {
//	r.mtx.RLock()
//	defer r.mtx.RUnlock()
//	return r.newViews.Size()
//}

func (s *roundState) LockHash() {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.prepare != nil {
		s.lockedHash = s.prepare.Proposal.Hash()
	}
}

func (s *roundState) UnlockHash() {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.lockedHash = EmptyHash
}

func (s *roundState) IsHashLocked() bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if s.lockedHash == EmptyHash {
		return false
	}
	return !s.hasBadProposal(s.GetLockedHash())
}

func (s *roundState) GetLockedHash() common.Hash {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return s.lockedHash
}

func (s *roundState) AddPrepareVote(msg *message) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.prepareVotes.Add(msg)
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
