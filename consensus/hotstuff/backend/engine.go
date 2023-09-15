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
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/contracts/native/governance"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	inmemorySnapshots = 128 // Number of recent epoch header
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

// Filter out system transactions from common transactions
// returns common transactions, system transactions and system transaction message provider
func (s *backend) BlockTransactions(block *types.Block, state *state.StateDB) (types.Transactions, types.Transactions,
	func(*types.Transaction, *big.Int) types.Message, error) {
	systemTransactions, err := governance.AssembleSystemTransactions(state, block.NumberU64())
	if err != nil {
		return nil, nil, nil, err
	}
	allTransactions := block.Transactions()
	commonTransactionCount := len(allTransactions) - len(systemTransactions)
	if commonTransactionCount < 0 {
		return nil, nil, nil, fmt.Errorf("missing required system transactions, count %v", len(systemTransactions))
	}

	signer := types.MakeSigner(s.chainConfig, block.Number())
	for i, tx := range systemTransactions {
		includedTx := allTransactions[commonTransactionCount + i]
		if signer.Hash(includedTx) != signer.Hash(tx) {
			return nil, nil, nil, fmt.Errorf("unexpected system tx hash detected, tx index %v, hash %s, expected: %s", commonTransactionCount + i, signer.Hash(includedTx), signer.Hash(tx))
		}
		from, err := signer.Sender(includedTx)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("check system tx signature failed, %w", err)
		}
		if from != block.Coinbase() {
			return nil, nil, nil, fmt.Errorf("check system tx signature failed, wrong signer %s", from)
		}
	}
	return allTransactions[:commonTransactionCount], systemTransactions, s.asSystemMessage, nil
}

// Change message from as valid system transaction sender
func(s *backend) asSystemMessage(tx *types.Transaction, baseFee *big.Int) types.Message {
	gasPrice := new(big.Int).Set(tx.GasPrice())
	if baseFee != nil {
		gasPrice = math.BigMin(gasPrice.Add(tx.GasTipCap(), baseFee), tx.GasFeeCap())
	}
	return types.NewMessage(utils.SystemTxSender, tx.To(), tx.Nonce(), tx.Value(), tx.Gas(), gasPrice,
		new(big.Int).Set(tx.GasFeeCap()), new(big.Int).Set(tx.GasTipCap()), tx.Data(), tx.AccessList(), true)
}

func (s *backend) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, uncles []*types.Header) error {
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash
	return nil
}

func (s *backend) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB,
	txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, []*types.Receipt, error) {

	// allow empty block in miner worker
	if txs == nil {
		txs = make([]*types.Transaction, 0)
	}
	if receipts == nil {
		receipts = make([]*types.Receipt, 0)
	}

	systemTransantions, err := governance.AssembleSystemTransactions(state, header.Number.Uint64())
	if err != nil {
		return nil, nil, err
	}

	for _, tx := range systemTransantions {
		chainContext := chainContext{Chain: chain, engine: s}
		gp := new(core.GasPool).AddGas(header.GasLimit)
		if err := gp.SubGas(header.GasUsed); err != nil { 
			return nil, nil, err
		}
		state.Prepare(tx.Hash(), common.Hash{}, len(txs))
		receipt, err := core.ApplyTransactionWithCustomMessageProvider(s.asSystemMessage, s.chainConfig, chainContext, nil, gp, state, header, tx, &header.GasUsed, vm.Config{})
		if err != nil {
			return nil, nil, err
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			return nil, nil, fmt.Errorf("unexpected reverted system transactions tx %d [%s], status %v", len(txs), tx.Hash(), receipt.Status)
		}
		signer := types.MakeSigner(s.chainConfig, header.Number)
		tx, err = s.signer.SignTx(tx, signer)
		if err != nil {
			return nil, nil, err
		}
		txs = append(txs, tx)
		receipts = append(receipts, receipt)
	}

	// Assemble and return the final block for sealing
	block := packBlock(state, chain, header, txs, receipts)
	return block, receipts, nil
}

func (s *backend) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) (err error) {
	// update the block header timestamp and signature and propose the block to core engine
	header := block.Header()

	// sign the sig hash and fill extra seal
	seal, err := s.signer.SignHash(s.SealHash(header))
	if err != nil {
		return err
	}
	if err := header.SetSeal(seal); err != nil {
		return err
	}
	block = block.WithSeal(header)

	go s.EventMux().Post(hotstuff.RequestEvent{Block: block})

	s.logger.Trace("WorkerSealNewBlock", "address", s.Address(), "hash", block.Hash(), "number", block.Number())
	return nil
}

func (s *backend) SealHash(header *types.Header) (hash common.Hash) {
	return types.SealHash(header)
}

// useless
func (s *backend) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return new(big.Int)
}

func (s *backend) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return []rpc.API{{
		Namespace: "hotstuff",
		Version:   "1.0",
		Service:   &API{chain: chain, hotstuff: s},
		Public:    true,
	}}
}

// Start implements consensus.Istanbul.Start
func (s *backend) Start(chain consensus.ChainReader, hasBadBlock func(db ethdb.Reader, hash common.Hash) bool) error {
	s.coreMu.Lock()
	defer s.coreMu.Unlock()

	if s.coreStarted {
		return ErrStartedEngine
	}

	s.chain = chain
	s.hasBadBlock = hasBadBlock

	// init validator set
	if next, err := s.newEpochValidators(); err != nil {
		return fmt.Errorf("get validators failed, err: %v", err)
	} else {
		s.vals = next.Copy()
	}

	// p2p module connect nodes directly
	s.nodesFeed.Send(consensus.StaticNodesEvent{Validators: s.vals.AddressList()})

	// MUST start in single goroutine because that the core.startNewRound need to request proposal in async mode.
	s.core.Start(chain)
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

	s.core.Stop()
	s.coreStarted = false
	return nil
}

func (s *backend) Close() error {
	return nil
}

func (s *backend) ReStart() {
	next, err := s.newEpochValidators()
	if err != nil {
		panic(fmt.Errorf("Restart consensus engine failed, err: %v ", err))
	}

	if next.Equal(s.vals.Copy()) {
		log.Trace("Restart Consensus engine, validators not changed.", "origin", s.vals.AddressList(), "current", next.AddressList())
		return
	}

	if s.coreStarted {
		s.Stop()
		// waiting for last engine instance free resource, e.g: unsubscribe...
		time.Sleep(2 * time.Second)
		log.Debug("Restart consensus engine...")
		s.Start(s.chain, s.hasBadBlock)
	}
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (s *backend) verifyHeader(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header, seal bool) error {
	if header.Number == nil {
		return errUnknownBlock
	}
	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != types.HotstuffDigest {
		return errInvalidMixDigest
	}
	// Ensure that extra info is not nil
	if header.Extra == nil {
		return errUnknownBlock
	}
	// Ensure that the block doesn't contain any uncles which are meaningless in Istanbul
	if header.UncleHash != nilUncleHash {
		return errInvalidUncleHash
	}
	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if header.Difficulty == nil || header.Difficulty.Cmp(defaultDifficulty) != 0 {
		return errInvalidDifficulty
	}

	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}

	// Ensure that the block's timestamp isn't less than it's parent
	var (
		parent *types.Header
	)
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	if header.Time < parent.Time {
		return errInvalidTimestamp
	}
	if header.Time > uint64(now().Unix()) {
		return consensus.ErrFutureBlock
	}

	// Verify the block's gas usage and (if applicable) verify the base fee.
	if !chain.Config().IsLondon(header.Number) {
		// Verify BaseFee not present before EIP-1559 fork.
		if header.BaseFee != nil {
			return fmt.Errorf("invalid baseFee before fork: have %d, expected 'nil'", header.BaseFee)
		}
		if err := misc.VerifyGaslimit(parent.GasLimit, header.GasLimit); err != nil {
			return err
		}
	} else if err := misc.VerifyEip1559Header(chain.Config(), parent, header); err != nil {
		// Verify the header's EIP-1559 attributes.
		return err
	}

	// Get validator set
	isEpoch, vals, err := s.getValidatorsByHeader(header, parent, chain)
	if err != nil {
		return err
	}

	// recover and verify signatures
	if _, err := s.signer.VerifyHeader(header, vals, seal); err != nil {
		return err
	}

	// save validators in lru cache
	if isEpoch {
		s.saveRecentHeader(header)
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
