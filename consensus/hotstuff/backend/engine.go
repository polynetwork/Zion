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
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	inmemorySnapshots = 128 // Number of recent vote snapshots to keep in memory
	inmemoryPeers     = 1000
	inmemoryMessages  = 1024
)

// HotStuff protocol constants.
var (
	defaultDifficulty = big.NewInt(1)
	nilUncleHash      = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
	emptyNonce        = types.BlockNonce{}
	now               = time.Now
)

func (s *backend) Author(header *types.Header) (common.Address, error) {
	signer, _, err := s.signer.Recover(header)
	return signer, err
}

func (s *backend) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	return s.verifyHeader(chain, header, nil, seal)
}

func (s *backend) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))
	go func() {
		for i, header := range headers {
			seal := false
			if seals != nil && len(seals) > i {
				seal = seals[i]
			}
			err := s.verifyHeader(chain, header, headers[:i], seal)

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

func (s *backend) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errInvalidUncleHash
	}
	return nil
}

func (s *backend) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	// unused fields, force to set to empty
	header.Coinbase = s.Address()
	header.Nonce = emptyNonce
	header.MixDigest = types.HotstuffDigest

	// copy the parent extra data as the header extra data
	parent, err := s.getPendingParentHeader(chain, header)
	if err != nil {
		return err
	}

	// use the same difficulty for all blocks
	header.Difficulty = defaultDifficulty

	// set header's timestamp
	header.Time = parent.Time + s.config.BlockPeriod
	if header.Time < uint64(time.Now().Unix()) {
		header.Time = uint64(time.Now().Unix())
	}

	return nil
}

func (s *backend) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs *[]*types.Transaction,
	uncles []*types.Header, receipts *[]*types.Receipt, systemTxs *[]*types.Transaction, usedGas *uint64) error {

	if err := s.reward(state, header.Number); err != nil {
		return err
	}

	if err := s.execEndBlock(&systemTxContext{
		chain:    chain,
		state:    state,
		header:   header,
		chainCtx: chainContext{Chain: chain, engine: s},
		txs:      txs,
		sysTxs:   systemTxs,
		receipts: receipts,
		usedGas:  usedGas,
		mining:   true,
	}); err != nil {
		s.logger.Debug("FinalizeAndAssemble", "hash", header.Hash(), "execute `endBlock` failed", err)
	}

	if beforeChange, _, _ := s.CheckPoint(header.Number.Uint64()); beforeChange {
		ctx := &systemTxContext{
			chain:    chain,
			state:    state,
			header:   header,
			chainCtx: chainContext{Chain: chain, engine: s},
			txs:      txs,
			sysTxs:   systemTxs,
			receipts: receipts,
			usedGas:  usedGas,
			mining:   true,
		}
		if err := s.execEpochChange(ctx); err != nil {
			s.logger.Debug("FinalizeAndAssemble", "hash", header.Hash(), "execute `epochChange` failed", err)
		}
		if err := s.applySnapshot(state, header.Number, false); err != nil {
			return err
		}
	}

	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash
	return nil
}

func (s *backend) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB,
	txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, []*types.Receipt, error) {

	if err := s.reward(state, header.Number); err != nil {
		return nil, nil, err
	}

	if txs == nil {
		txs = make([]*types.Transaction, 0)
	}
	if receipts == nil {
		receipts = make([]*types.Receipt, 0)
	}

	if err := s.execEndBlock(&systemTxContext{
		chain:    chain,
		state:    state,
		header:   header,
		chainCtx: chainContext{Chain: chain, engine: s},
		txs:      &txs,
		sysTxs:   nil,
		receipts: &receipts,
		usedGas:  &header.GasUsed,
		mining:   true,
	}); err != nil {
		// todo(fuk): governance forbid execute error
		s.logger.Debug("FinalizeAndAssemble", "hash", header.Hash(), "execute `endBlock` failed", err)
	}

	if beforeChange, _, _ := s.CheckPoint(header.Number.Uint64()); beforeChange {
		if err := s.execEpochChange(&systemTxContext{
			chain:    chain,
			state:    state,
			header:   header,
			chainCtx: chainContext{Chain: chain, engine: s},
			txs:      &txs,
			sysTxs:   nil,
			receipts: &receipts,
			usedGas:  &header.GasUsed,
			mining:   true,
		}); err == nil {
			// todo(fuk): governance forbid execute error
			s.logger.Debug("FinalizeAndAssemble", "hash", header.Hash(), "execute `epochChange` failed", err)
		}
		if err := s.applySnapshot(state, header.Number, true); err != nil {
			return nil, nil, err
		}
	}

	// Assemble and return the final block for sealing
	block := packBlock(state, chain, header, txs, receipts)
	return block, receipts, nil
}

func (s *backend) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) (err error) {
	// update the block header timestamp and signature and propose the block to core engine
	header := block.Header()

	// sign the sig hash and fill extra seal
	if err = s.signer.SealBeforeCommit(header); err != nil {
		return err
	}
	block = block.WithSeal(header)

	go func() {
		// get the proposed block hash and clear it if the seal() is completed.
		s.sealMu.Lock()
		s.proposedBlockHash = block.Hash()
		s.logger.Trace("WorkerSealNewBlock", "hash", block.Hash(), "number", block.Number())

		defer func() {
			s.proposedBlockHash = common.EmptyHash
			s.sealMu.Unlock()
		}()

		// post block into Istanbul engine
		go s.EventMux().Post(hotstuff.RequestEvent{
			Proposal: block,
		})
		for {
			select {
			case result := <-s.commitCh:
				// if the block hash and the hash from channel are the same,
				// return the result. Otherwise, keep waiting the next hash.
				if result != nil && block.Hash() == result.Hash() {
					results <- result
					return
				}
			case <-stop:
				s.logger.Trace("Stop seal, check miner status!")
				results <- nil
				return
			}
		}
	}()
	return nil
}

func (s *backend) SealHash(header *types.Header) common.Hash {
	return s.signer.SigHash(header)
}

func (s *backend) ValidateBlock(block *types.Block) error {
	return s.chain.PreExecuteBlock(block)
}

// useless
func (s *backend) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return new(big.Int)
}

func (s *backend) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return []rpc.API{{
		Namespace: "istanbul",
		Version:   "1.0",
		Service:   &API{chain: chain, hotstuff: s},
		Public:    true,
	}}
}

// Start implements consensus.Istanbul.Start
func (s *backend) Start(chain consensus.ChainReader, currentBlock func() *types.Block, getBlockByHash func(hash common.Hash) *types.Block, hasBadBlock func(hash common.Hash) bool) error {
	s.coreMu.Lock()
	defer s.coreMu.Unlock()

	if s.coreStarted {
		return ErrStartedEngine
	}

	// clear previous data
	if s.commitCh != nil {
		close(s.commitCh)
	}
	s.commitCh = make(chan *types.Block, 1)

	s.chain = chain
	s.currentBlock = currentBlock
	s.getBlockByHash = getBlockByHash
	s.hasBadBlock = hasBadBlock

	if err := s.core.Start(chain); err != nil {
		return err
	}

	s.coreStarted = true
	return nil
}

// Stop implements consensus.Istanbul.Stop
func (s *backend) Stop() error {
	s.coreMu.Lock()
	defer s.coreMu.Unlock()
	if !s.coreStarted {
		return nil
	}
	if err := s.core.Stop(); err != nil {
		return err
	}
	s.coreStarted = false
	return nil
}

func (s *backend) Close() error {
	return nil
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (s *backend) verifyHeader(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header, seal bool) error {
	if err := CustomVerifyHeader(header); err != nil {
		return err
	}

	// verifyCascadingFields verifies all the header fields that are not standalone,
	// rather depend on a batch of previous headers. The caller may optionally pass
	// in a batch of parents (ascending order) to avoid looking those up from the
	// database. This is useful for concurrently verifying a batch of new headers.
	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}

	// Ensure that the block's timestamp isn't too close to it's parent
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	if header.Time > parent.Time+s.config.BlockPeriod && header.Time > uint64(now().Unix()) {
		return errInvalidTimestamp
	}

	if err := s.snaps.sync(s.db, parent, header); err != nil {
		return err
	}

	vals := s.Validators(number)
	if _, err := s.signer.VerifyHeader(header, vals, seal); err != nil {
		return err
	}
	return nil
}

func CustomVerifyHeader(header *types.Header) error {
	if header.Number == nil {
		return errUnknownBlock
	}

	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != types.HotstuffDigest {
		return errInvalidMixDigest
	}
	// Ensure that the block doesn't contain any uncles which are meaningless in Istanbul
	if header.UncleHash != nilUncleHash {
		return errInvalidUncleHash
	}
	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if header.Difficulty == nil || header.Difficulty.Cmp(defaultDifficulty) != 0 {
		return errInvalidDifficulty
	}

	return nil
}

func (s *backend) getPendingParentHeader(chain consensus.ChainHeaderReader, header *types.Header) (*types.Header, error) {
	number := header.Number.Uint64()
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return nil, consensus.ErrUnknownAncestor
	}
	return parent, nil
}
