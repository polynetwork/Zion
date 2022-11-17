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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
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
	db ethdb.Database
	vs hotstuff.ValidatorSet

	round  *big.Int
	height *big.Int
	state  State

	pendingRequest *Request
	node           *Node        //
	lockedBlock    *types.Block // validator's prepare proposal
	proposalLocked bool

	// o(4n)
	newViews       *MessageSet // data set for newView message
	prepareVotes   *MessageSet // data set for prepareVote message
	preCommitVotes *MessageSet // data set for preCommitVote message
	commitVotes    *MessageSet // data set for commitVote message

	highQC      *QuorumCert // leader highQC
	preCommitQC *QuorumCert // leader preCommitQC
	prepareQC   *QuorumCert // prepareQC for repo and leader
	lockQC      *QuorumCert // lockQC for repo and leader
	committedQC *QuorumCert // committedQC for repo and leader
}

// newRoundState creates a new roundState instance with the given view and validatorSet
func newRoundState(view *View, validatorSet hotstuff.ValidatorSet, db ethdb.Database) *roundState {
	rs := &roundState{
		db:             db,
		vs:             validatorSet.Copy(),
		round:          view.Round,
		height:         view.Height,
		state:          StateAcceptRequest,
		newViews:       NewMessageSet(validatorSet),
		prepareVotes:   NewMessageSet(validatorSet),
		preCommitVotes: NewMessageSet(validatorSet),
		commitVotes:    NewMessageSet(validatorSet),
	}
	return rs
}

// clean all votes message set for new round
func (s *roundState) update(view *View, vs hotstuff.ValidatorSet) {
	s.vs = vs
	s.height = view.Height
	s.round = view.Round
	s.state = StateAcceptRequest
	s.newViews = NewMessageSet(vs)
	s.prepareVotes = NewMessageSet(vs)
	s.preCommitVotes = NewMessageSet(vs)
	s.commitVotes = NewMessageSet(vs)
}

func (s *roundState) Height() *big.Int {
	if s.height == nil {
		return big.NewInt(0)
	}

	return s.height
}

func (s *roundState) HeightU64() uint64 {
	return s.Height().Uint64()
}

func (s *roundState) Round() *big.Int {
	if s.round == nil {
		return big.NewInt(0)
	}
	return s.round
}

func (s *roundState) RoundU64() uint64 {
	return s.Round().Uint64()
}

func (s *roundState) View() *View {
	return &View{
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

func (s *roundState) SetNode(node *Node) error {
	if err := s.storeNode(node); err != nil {
		return err
	}
	s.node = node
	return nil
}

func (s *roundState) SetNodeWithSealBlock(block *types.Block) error {
	s.node.Block = block
	if err := s.storeNode(s.node); err != nil {
		return err
	}
	s.lockedBlock = block
	return nil
}

func (s *roundState) Node() *Node {
	return s.node
}

//func (s *roundState) Proposal() hotstuff.Proposal {
//	return s.proposal
//}

func (s *roundState) LockNode() error {
	if s.node == nil || s.node.Block == nil {
		return errInvalidNode
	}
	s.lockedBlock = s.node.Block
	s.proposalLocked = true
	return nil
}

func (s *roundState) LockedNode() *Node {
	if s.proposalLocked && s.node != nil {
		return s.node
	}
	return nil
}

func (s *roundState) LockedBlock() *types.Block {
	return s.lockedBlock
}

// 如果在prepare阶段就开始锁，会导致更多问题，prepareQC。
func (s *roundState) IsProposalLocked() bool {
	return s.proposalLocked
}

func (s *roundState) SetPendingRequest(req *Request) {
	s.pendingRequest = req
}

func (s *roundState) PendingRequest() *Request {
	return s.pendingRequest
}

func (s *roundState) Vote() common.Hash {
	return s.node.Hash()
}

func (s *roundState) SetHighQC(qc *QuorumCert) {
	s.highQC = qc
}

func (s *roundState) HighQC() *QuorumCert {
	return s.highQC
}

func (s *roundState) SetPrepareQC(qc *QuorumCert) error {
	if err := s.storePrepareQC(qc); err != nil {
		return err
	}
	s.prepareQC = qc
	return nil
}

func (s *roundState) PrepareQC() *QuorumCert {
	return s.prepareQC
}

func (s *roundState) SetPreCommittedQC(qc *QuorumCert) {
	s.preCommitQC = qc
}

func (s *roundState) PreCommittedQC() *QuorumCert {
	return s.preCommitQC
}

func (s *roundState) SetLockQC(qc *QuorumCert) error {
	if err := s.storeLockQC(qc); err != nil {
		return err
	}
	s.lockQC = qc
	return nil
}

func (s *roundState) LockQC() *QuorumCert {
	return s.lockQC
}

func (s *roundState) SetCommittedQC(qc *QuorumCert) error {
	if err := s.storeCommitQC(qc); err != nil {
		return err
	}
	s.committedQC = qc
	return nil
}

func (s *roundState) CommittedQC() *QuorumCert {
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

func (s *roundState) PreCommitVotes() []*Message {
	return s.preCommitVotes.Values()
}

func (s *roundState) PreCommitVoteSize() int {
	return s.preCommitVotes.Size()
}

func (s *roundState) AddCommitVote(msg *Message) error {
	return s.commitVotes.Add(msg)
}

func (s *roundState) CommitVotes() []*Message {
	return s.commitVotes.Values()
}

func (s *roundState) CommitVoteSize() int {
	return s.commitVotes.Size()
}

func (s *roundState) GetCommittedSeals(n int) [][]byte {
	seals := make([][]byte, n)
	for i, data := range s.commitVotes.Values() {
		if i < n {
			seals[i] = data.CommittedSeal
		}
	}
	return seals
}

// -----------------------------------------------------------------------
//
// store round state as snapshot
//
// -----------------------------------------------------------------------

const (
	dbRoundStatePrefix = "round-state-"
	viewSuffix         = "view"
	prepareQCSuffix    = "prepareQC"
	lockQCSuffix       = "lockQC"
	commitQCSuffix     = "commitQC"
	nodeSuffix         = "node"
	blockSuffix        = "block"
)

// todo(fuk): 不能返回error，这里需要考虑到两种情况，一种是节点半路加入共识，此时其所有的存储状态为空，也就是之前的qc都没有存储过
// 此外就是对于block1，可能存在几轮都失败的情况
func (s *roundState) reload(view *View) {
	_ = s.loadView(view)
	_ = s.loadPrepareQC()
	_ = s.loadLockQC()
	_ = s.loadCommitQC()
	_ = s.loadNode()
}

func (s *roundState) storeView(view *View) error {
	if s.db == nil {
		return nil
	}

	raw, err := Encode(view)
	if err != nil {
		return err
	}
	return s.db.Put(viewKey(), raw)
}

func (s *roundState) loadView(cur *View) error {
	if s.db == nil {
		return nil
	}

	view := new(View)
	raw, err := s.db.Get(viewKey())
	if err != nil {
		return err
	}
	if err = rlp.DecodeBytes(raw, view); err != nil {
		return err
	}
	if view.Cmp(cur) > 0 {
		s.height = view.Height
		s.round = view.Round
	}
	return nil
}

func (s *roundState) storePrepareQC(qc *QuorumCert) error {
	if s.db == nil {
		return nil
	}

	raw, err := Encode(qc)
	if err != nil {
		return err
	}
	return s.db.Put(prepareQCKey(), raw)
}

func (s *roundState) loadPrepareQC() error {
	if s.db == nil {
		return nil
	}

	data := new(QuorumCert)
	raw, err := s.db.Get(prepareQCKey())
	if err != nil {
		return err
	}
	if err = rlp.DecodeBytes(raw, data); err != nil {
		return err
	}
	s.prepareQC = data
	return nil
}

func (s *roundState) storeLockQC(qc *QuorumCert) error {
	if s.db == nil {
		return nil
	}

	raw, err := Encode(qc)
	if err != nil {
		return err
	}
	return s.db.Put(lockQCKey(), raw)
}

func (s *roundState) loadLockQC() error {
	if s.db == nil {
		return nil
	}

	data := new(QuorumCert)
	raw, err := s.db.Get(lockQCKey())
	if err != nil {
		return err
	}
	if err = rlp.DecodeBytes(raw, data); err != nil {
		return err
	}
	s.lockQC = data
	return nil
}

func (s *roundState) storeCommitQC(qc *QuorumCert) error {
	if s.db == nil {
		return nil
	}

	raw, err := Encode(qc)
	if err != nil {
		return err
	}
	return s.db.Put(commitQCKey(), raw)
}

func (s *roundState) loadCommitQC() error {
	if s.db == nil {
		return nil
	}

	data := new(QuorumCert)
	raw, err := s.db.Get(commitQCKey())
	if err != nil {
		return err
	}
	if err = rlp.DecodeBytes(raw, data); err != nil {
		return err
	}
	s.committedQC = data
	return nil
}

func (s *roundState) storeNode(node *Node) error {
	if s.db == nil {
		return nil
	}

	raw, err := Encode(node)
	if err != nil {
		return err
	}
	return s.db.Put(nodeKey(), raw)
}

func (s *roundState) loadNode() error {
	if s.db == nil {
		return nil
	}

	data := new(Node)
	raw, err := s.db.Get(nodeKey())
	if err != nil {
		return err
	}
	if err = rlp.DecodeBytes(raw, data); err != nil {
		return err
	}
	s.node = data
	return nil
}

// todo(fuk): delete after test
//func (s *roundState) storeBlock(block *types.Block) error {
//	if s.db == nil {
//		return nil
//	}
//
//	raw, err := Encode(block)
//	if err != nil {
//		return err
//	}
//	return s.db.Put(nodeKey(), raw)
//}
//
//func (s *roundState) loadBlock() error {
//	if s.db == nil {
//		return nil
//	}
//
//	data := new(Node)
//	raw, err := s.db.Get(nodeKey())
//	if err != nil {
//		return err
//	}
//	if err = rlp.DecodeBytes(raw, data); err != nil {
//		return err
//	}
//	s.node = data
//	return nil
//}

func viewKey() []byte {
	return append([]byte(dbRoundStatePrefix), []byte(viewSuffix)...)
}

func prepareQCKey() []byte {
	return append([]byte(dbRoundStatePrefix), []byte(prepareQCSuffix)...)
}

func lockQCKey() []byte {
	return append([]byte(dbRoundStatePrefix), []byte(lockQCSuffix)...)
}

func commitQCKey() []byte {
	return append([]byte(dbRoundStatePrefix), []byte(commitQCSuffix)...)
}

func nodeKey() []byte {
	return append([]byte(dbRoundStatePrefix), []byte(nodeSuffix)...)
}

func blockKey() []byte {
	return append([]byte(dbRoundStatePrefix), []byte(blockSuffix)...)
}
