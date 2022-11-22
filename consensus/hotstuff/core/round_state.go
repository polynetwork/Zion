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
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/log"

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

type roundNode struct {
	temp *Node // cache node before `prepared`
	node *Node
}

type roundState struct {
	db     ethdb.Database
	logger log.Logger
	vs     hotstuff.ValidatorSet

	round  *big.Int
	height *big.Int
	state  State

	lastChainedBlock *types.Block
	pendingRequest   *Request
	node             *roundNode
	lockedBlock      *types.Block // validator's prepare proposal
	proposalLocked   bool

	// todo(fuk): need temp nodes queue

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
func newRoundState(db ethdb.Database, logger log.Logger, validatorSet hotstuff.ValidatorSet, lastChainedBlock *types.Block, view *View) *roundState {
	rs := &roundState{
		db:               db,
		logger:           logger,
		vs:               validatorSet.Copy(),
		round:            view.Round,
		height:           view.Height,
		state:            StateAcceptRequest,
		node:             new(roundNode),
		lastChainedBlock: lastChainedBlock,
		newViews:         NewMessageSet(validatorSet),
		prepareVotes:     NewMessageSet(validatorSet),
		preCommitVotes:   NewMessageSet(validatorSet),
		commitVotes:      NewMessageSet(validatorSet),
	}
	return rs
}

// clean all votes message set for new round
func (s *roundState) update(vs hotstuff.ValidatorSet, lastChainedBlock *types.Block, view *View) {
	if view.HeightU64() > s.HeightU64() {
		if lastChainedBlock.NumberU64() == s.lastChainedBlock.NumberU64() {
			panic(fmt.Sprintf("block fork choice, last view %v, current view %v,"+
				" last round hash %v, current round hash %v",
				s.View(), view, s.lastChainedBlock.Hash(), lastChainedBlock.Hash()))
		} else if lastChainedBlock.NumberU64() != s.lastChainedBlock.NumberU64()+1 {
			panic(fmt.Sprintf("invalid `lastProposal` height, expect %v, got %v",
				s.lastChainedBlock.NumberU64()+1, lastChainedBlock.NumberU64()))
		} else {
			s.lastChainedBlock = lastChainedBlock
		}
	}

	s.vs = vs.Copy()
	s.height = view.Height
	s.round = view.Round
	s.newViews = NewMessageSet(vs)
	s.prepareVotes = NewMessageSet(vs)
	s.preCommitVotes = NewMessageSet(vs)
	s.commitVotes = NewMessageSet(vs)
}

func (s *roundState) View() *View {
	return &View{
		Round:  s.round,
		Height: s.height,
	}
}

func (s *roundState) Height() *big.Int {
	return s.height
}

func (s *roundState) HeightU64() uint64 {
	return s.height.Uint64()
}

func (s *roundState) Round() *big.Int {
	return s.round
}

func (s *roundState) RoundU64() uint64 {
	return s.round.Uint64()
}

func (s *roundState) SetState(state State) {
	s.state = state
}

func (s *roundState) State() State {
	return s.state
}

func (s *roundState) LastChainedBlock() *types.Block {
	return s.lastChainedBlock
}

// accept pending request from miner only for once.
func (s *roundState) SetPendingRequest(req *Request) {
	if s.pendingRequest == nil {
		s.pendingRequest = req
	}
}

func (s *roundState) PendingRequest() *Request {
	return s.pendingRequest
}

func (s *roundState) SetNode(node *Node) error {
	if temp := s.node.temp; temp == nil {
		s.node.temp = node
		return nil
	} else if temp.Hash() != node.Hash() {
		return fmt.Errorf("expect node %v, got %v", temp.Hash(), node.Hash())
	} else {
		if err := s.storeNode(node); err != nil {
			return err
		}
		s.node.node = node
		s.node.temp = nil
	}
	return nil
}

func (s *roundState) Node() *Node {
	if temp := s.node.temp; temp != nil {
		return temp
	} else {
		return s.node.node
	}
}

func (s *roundState) Lock(qc *QuorumCert) error {
	if s.node == nil || s.node.node == nil {
		return errInvalidNode
	}

	if err := s.storeLockQC(qc); err != nil {
		return err
	}
	if err := s.storeNode(s.node.node); err != nil {
		return err
	}

	s.lockQC = qc
	s.lockedBlock = s.node.node.Block
	s.proposalLocked = true
	return nil
}

// Unlock it's happened at the start of new round, new state is `StateAcceptRequest`, and `lockQC` keep to judge safety rule
func (s *roundState) Unlock() error {
	s.pendingRequest = nil
	s.proposalLocked = false
	s.lockedBlock = nil
	s.node.temp = nil
	return nil
}

func (s *roundState) LockedBlock() *types.Block {
	if s.proposalLocked && s.lockedBlock != nil {
		return s.lockedBlock
	}
	return nil
}

func (s *roundState) SetSealedBlock(block *types.Block) error {
	if s.node.node == nil || s.node.node.Block == nil {
		return fmt.Errorf("locked block is nil")
	}
	if s.node.node.Block.Hash() != block.Hash() {
		return fmt.Errorf("node block not equal to multi-seal block")
	}
	s.node.node.Block = block
	if err := s.storeNode(s.node.node); err != nil {
		return err
	}
	s.lockedBlock = block
	return nil
}

func (s *roundState) Vote() common.Hash {
	if node := s.Node(); node == nil {
		return common.EmptyHash
	} else {
		return node.Hash()
	}
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

// -----------------------------------------------------------------------
//
// leader collect votes
//
// -----------------------------------------------------------------------
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
// state是否需要reload???
func (s *roundState) reload(view *View) {
	var (
		err      error
		printErr = s.logger != nil && s.height.Uint64() > 1
	)

	if err = s.loadView(view); err != nil && printErr {
		s.logger.Warn("Load view failed", "err", err)
	}
	if err = s.loadPrepareQC(); err != nil && printErr {
		s.logger.Warn("Load prepareQC failed", "err", err)
	}
	if err = s.loadLockQC(); err != nil && printErr {
		s.logger.Warn("Load lockQC failed", "err", err)
	}
	if err = s.loadCommitQC(); err != nil && printErr {
		s.logger.Warn("Load commitQC failed", "err", err)
	}
	if err = s.loadNode(); err != nil && printErr {
		s.logger.Warn("Load node failed", "err", err)
	}

	// reset locked node
	if s.lockQC != nil && s.node.node != nil && s.node.node.Block != nil && s.lockQC.node == s.node.node.Hash() {
		s.lockedBlock = s.node.node.Block
		s.proposalLocked = true
	}
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
	s.node.node = data
	return nil
}

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
