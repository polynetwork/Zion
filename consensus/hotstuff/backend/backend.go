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

package backend

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/core"
	snr "github.com/ethereum/go-ethereum/consensus/hotstuff/signer"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	lru "github.com/hashicorp/golang-lru"
)

// HotStuff is the scalable hotstuff consensus engine
type backend struct {
	db          ethdb.Database // Database to store and retrieve necessary information
	logger      log.Logger
	config      *hotstuff.Config
	chainConfig *params.ChainConfig
	core        hotstuff.CoreEngine
	signer      hotstuff.Signer
	chain       consensus.ChainReader
	vals        hotstuff.ValidatorSet // consensus participant collection

	recents        *lru.ARCCache // Snapshots for recent block to speed up reorgs
	recentMessages *lru.ARCCache // the cache of peer's messages
	knownMessages  *lru.ARCCache // the cache of self messages

	// signal for engine running status
	coreStarted bool
	coreMu      sync.RWMutex

	broadcaster consensus.Broadcaster // event subscription for ChainHeadEvent event
	nodesFeed   event.Feed            // event subscription for static nodes listen
	executeFeed event.Feed            // event subscription for executed state
	requestFeed, messageFeed, commitFeed    event.Feed        // message sender for engine

	epochMu int32 // check point mutex

	// closure for help
	currentBlock   func() *types.Block
	getBlockByHash func(hash common.Hash) *types.Block
	hasBadBlock    func(db ethdb.Reader, hash common.Hash) bool
}

func New(chainConfig *params.ChainConfig, config *hotstuff.Config, privateKey *ecdsa.PrivateKey, db ethdb.Database, mock bool) *backend {
	recents, _ := lru.NewARC(inmemorySnapshots)
	recentMessages, _ := lru.NewARC(inmemoryPeers)
	knownMessages, _ := lru.NewARC(inmemoryMessages)

	signer := snr.NewSigner(privateKey)
	backend := &backend{
		config:         config,
		chainConfig:    chainConfig,
		db:             db,
		logger:         log.New(),
		coreStarted:    false,
		signer:         signer,
		recentMessages: recentMessages,
		knownMessages:  knownMessages,
		recents:        recents,
	}

	if mock {
		backend.core = core.New(backend, config, signer, db, nil)
	} else {
		backend.core = core.New(backend, config, signer, db, backend.CheckPoint)
	}

	return backend
}

// Address implements hotstuff.Backend.Address
func (s *backend) Address() common.Address {
	return s.signer.Address()
}

func (s *backend) SubscribeEvent(ch interface{}) event.Subscription {
	switch c := ch.(type) {
	case chan hotstuff.RequestEvent:
		return s.requestFeed.Subscribe(c)
	case chan hotstuff.MessageEvent:
		return s.messageFeed.Subscribe(c)
	case chan hotstuff.FinalCommittedEvent:
		return s.commitFeed.Subscribe(c)
	default:
		panic(fmt.Sprintf("unexpected subscriber type %t", ch))
	}
}

func (s *backend) Send(ev interface{}) int {
	switch event := ev.(type) {
	case hotstuff.RequestEvent:
		return s.requestFeed.Send(event)
	case hotstuff.MessageEvent:
		return s.messageFeed.Send(event)
	case hotstuff.FinalCommittedEvent:
		return s.commitFeed.Send(event)
	default:
		panic(fmt.Sprintf("unexpected event type %t", ev))
	}
}

// Broadcast implements hotstuff.Backend.Broadcast
func (s *backend) Broadcast(valSet hotstuff.ValidatorSet, payload []byte) error {
	// send to others
	if err := s.Gossip(valSet, payload); err != nil {
		return err
	}
	// send to self
	msg := hotstuff.MessageEvent{
		Src:     s.Address(),
		Payload: payload,
	}
	go s.messageFeed.Send(msg)
	return nil
}

// Broadcast implements hotstuff.Backend.Gossip
func (s *backend) Gossip(valSet hotstuff.ValidatorSet, payload []byte) error {
	hash := hotstuff.RLPHash(payload)
	s.knownMessages.Add(hash, true)

	targets := make(map[common.Address]bool)
	for _, val := range valSet.List() { // hotstuff/validator/default.go - defaultValidator
		if val.Address() != s.Address() {
			targets[val.Address()] = true
		}
	}
	if s.broadcaster != nil && len(targets) > 0 {
		ps := s.broadcaster.FindPeers(targets)
		for addr, p := range ps {
			ms, ok := s.recentMessages.Get(addr)
			var m *lru.ARCCache
			if ok {
				m, _ = ms.(*lru.ARCCache)
				if _, k := m.Get(hash); k {
					// This peer had this event, skip it
					continue
				}
			} else {
				m, _ = lru.NewARC(inmemoryMessages)
			}

			m.Add(hash, true)
			s.recentMessages.Add(addr, m)
			go p.Send(hotstuffMsg, payload)
		}
	}
	return nil
}

// Unicast implements hotstuff.Backend.Unicast
func (s *backend) Unicast(valSet hotstuff.ValidatorSet, payload []byte) error {
	msg := hotstuff.MessageEvent{Src: s.Address(), Payload: payload}
	leader := valSet.GetProposer()
	target := leader.Address()
	hash := hotstuff.RLPHash(payload)
	s.knownMessages.Add(hash, true)

	// send to self
	if s.Address() == target {
		go s.messageFeed.Send(msg)
		return nil
	}

	// send to other peer
	if s.broadcaster != nil {
		if p := s.broadcaster.FindPeer(target); p != nil {
			ms, ok := s.recentMessages.Get(target)
			var m *lru.ARCCache
			if ok {
				m, _ = ms.(*lru.ARCCache)
				if _, k := m.Get(hash); k {
					return nil
				}
			} else {
				m, _ = lru.NewARC(inmemoryMessages)
			}
			m.Add(hash, true)
			s.recentMessages.Add(target, m)
			go func() {
				if err := p.Send(hotstuffMsg, payload); err != nil {
					s.logger.Error("unicast message failed", "err", err)
				}
			}()
		}
	}
	return nil
}

// SealBlock fullfill multi signatures into block header
func (s *backend) SealBlock(block *types.Block, seals [][]byte) (*types.Block, error) {

	// check proposal
	if block.Header() == nil {
		s.logger.Error("Invalid proposal precommit")
		return nil, errInvalidProposal
	}

	// check seals
	if len(seals) == 0 {
		return nil, errInvalidCommittedSeals
	}
	for _, seal := range seals {
		if len(seal) != types.HotstuffExtraSeal {
			return nil, errInvalidCommittedSeals
		}
	}

	// Append seals into extra-data and update block's header
	h := block.Header()
	if err := h.SetCommittedSeal(seals); err != nil {
		return nil, err
	}
	return block.WithSeal(h), nil
}

// Commit for most pos and pow chain, the local block should be write in state database directly,
// and the remote block only need to broadcast to other nodes. in hotstuff consensus,
// sent the finalized block to miner.worker although it may be an remote block.
func (s *backend) Commit(executed *consensus.ExecutedBlock) error {
	if executed == nil || executed.Block == nil {
		return fmt.Errorf("invalid executed block")
	}
	s.executeFeed.Send(*executed)
	s.logger.Info("Committed", "address", s.Address(), "hash", executed.Block.Hash(), "number", executed.Block.Number())
	return nil
}

// Verify implements hotstuff.Backend.Verify
func (s *backend) Verify(block *types.Block, seal bool) (time.Duration, error) {
	// check bad block
	if s.HasBadProposal(block.Hash()) {
		return 0, errBADProposal
	}

	// check block body
	txnHash := types.DeriveSha(block.Transactions(), trie.NewStackTrie(nil))
	uncleHash := types.CalcUncleHash(block.Uncles())
	if txnHash != block.Header().TxHash {
		return 0, errMismatchTxhashes
	}
	if uncleHash != nilUncleHash {
		return 0, errInvalidUncleHash
	}

	// verify the header of proposed block
	err := s.VerifyHeader(s.chain, block.Header(), seal)
	if err == nil {
		return 0, nil
	} else if err == consensus.ErrFutureBlock {
		return time.Unix(int64(block.Header().Time), 0).Sub(now()), consensus.ErrFutureBlock
	}
	return 0, err
}

func (s *backend) LastProposal() (*types.Block, common.Address) {
	var (
		proposer common.Address
		err      error
	)

	block := s.chain.CurrentBlock()
	if block.Number().Cmp(common.Big0) > 0 {
		if proposer, err = s.Author(block.Header()); err != nil {
			s.logger.Error("Failed to get block proposer", "err", err)
			return nil, common.Address{}
		}
	}

	// Return header only block here since we don't need block body
	return block, proposer
}

func (s *backend) GetProposal(hash common.Hash) hotstuff.Proposal {
	return s.chain.GetBlockByHash(hash)
}

// HasProposal implements hotstuff.Backend.HashBlock
func (s *backend) HasProposal(hash common.Hash, number *big.Int) bool {
	return s.chain.GetHeader(hash, number.Uint64()) != nil
}

// GetSpeaker implements hotstuff.Backend.GetProposer
func (s *backend) GetProposer(number uint64) common.Address {
	if header := s.chain.GetHeaderByNumber(number); header != nil {
		a, _ := s.Author(header)
		return a
	}
	return common.Address{}
}

func (s *backend) HasBadProposal(hash common.Hash) bool {
	if s.hasBadBlock == nil {
		return false
	}
	return s.hasBadBlock(s.db, hash)
}

func (s *backend) ExecuteBlock(block *types.Block) (*consensus.ExecutedBlock, error) {
	state, receipts, allLogs, err := s.chain.ExecuteBlock(block)
	if err != nil {
		return nil, err
	}
	return &consensus.ExecutedBlock{
		State:    state,
		Block:    block,
		Receipts: receipts,
		Logs:     allLogs,
	}, nil
}
