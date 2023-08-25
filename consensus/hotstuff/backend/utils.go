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
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
)

// ===========================     utility function        ==========================

// chain context
type chainContext struct {
	Chain  consensus.ChainHeaderReader
	engine consensus.Engine
}

func (c chainContext) Engine() consensus.Engine {
	return c.engine
}

// GetHeader blockContext need this function
func (c chainContext) GetHeader(hash common.Hash, number uint64) *types.Header {
	return c.Chain.GetHeader(hash, number)
}

const (
	systemGas      = math.MaxUint64 / 2 // system tx will be executed in evm, and gas calculating is needed.
	systemGasPrice = int64(0)           // consensus txs do not need to participate in gas price bidding
)

// getSystemCaller use fixed systemCaller as contract caller, and tx hash is useless in contract call.
func (s *backend) getSystemCaller(state *state.StateDB, height *big.Int) *native.ContractRef {
	caller := utils.SystemTxSender
	hash := common.EmptyHash
	return native.NewContractRef(state, caller, caller, height, hash, systemGas, nil)
}

func packBlock(state *state.StateDB, chain consensus.ChainHeaderReader,
	header *types.Header, txs []*types.Transaction, receipts []*types.Receipt) *types.Block {
	// perform root calculation and block reorganization at the same time which with a large number of memory copy.
	// and reset the header root after actions done.
	root := state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	// the header uncle hash will be settle as EmptyUncleHash which as the same of `nilUncleHash`
	block := types.NewBlock(header, txs, nil, receipts, trie.NewStackTrie(nil))

	block.SetRoot(root)
	return block
}

func NewDefaultValSet(list []common.Address) hotstuff.ValidatorSet {
	return validator.NewSet(list, hotstuff.RoundRobin)
}
