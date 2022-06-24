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
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
)

// ===========================     utility function        ==========================

// callmsg implements core.Message to allow passing it as a transaction simulator.
type callmsg struct {
	ethereum.CallMsg
}

func (m callmsg) From() common.Address { return m.CallMsg.From }
func (m callmsg) Nonce() uint64        { return 0 }
func (m callmsg) CheckNonce() bool     { return false }
func (m callmsg) To() *common.Address  { return m.CallMsg.To }
func (m callmsg) GasPrice() *big.Int   { return m.CallMsg.GasPrice }
func (m callmsg) Gas() uint64          { return m.CallMsg.Gas }
func (m callmsg) Value() *big.Int      { return m.CallMsg.Value }
func (m callmsg) Data() []byte         { return m.CallMsg.Data }

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

// get system message
func (s *backend) getSystemMessage(from, toAddress common.Address, data []byte, value *big.Int) callmsg {
	return callmsg{
		ethereum.CallMsg{
			From:     from,
			Gas:      math.MaxUint64 / 2,
			GasPrice: big.NewInt(0), // consensus txs do not need to participate in gas price bidding
			Value:    value,
			To:       &toAddress,
			Data:     data,
		},
	}
}

// applyTransaction execute transaction without miner worker.
func (s *backend) applyTransaction(
	chain consensus.ChainHeaderReader,
	msg callmsg,
	state *state.StateDB,
	header *types.Header,
	chainContext core.ChainContext,
	txs *[]*types.Transaction, receipts *[]*types.Receipt,
	receivedTxs *[]*types.Transaction, usedGas *uint64, mining bool,
) (err error) {
	nonce := state.GetNonce(msg.From())

	expectedTx := types.NewTransaction(nonce, *msg.To(), msg.Value(), msg.Gas(), msg.GasPrice(), msg.Data())
	signer := types.NewEIP155Signer(chain.Config().ChainID)
	expectedHash := signer.Hash(expectedTx)

	// miner worker use finalizeAndAssemble in which the param of `mining` is true,  it's denote
	// that this tx comes from miner. `validator` send governance tx in the same nonce is forbidden.
	if msg.From() == s.signer.Address() && mining {
		expectedTx, err = s.signer.SignTx(expectedTx, signer)
		if err != nil {
			return err
		}
	} else {
		if receivedTxs == nil || len(*receivedTxs) == 0 || (*receivedTxs)[0] == nil {
			return errors.New("supposed to get a actual transaction, but get none")
		}
		actualTx := (*receivedTxs)[0]
		if !bytes.Equal(signer.Hash(actualTx).Bytes(), expectedHash.Bytes()) {
			return fmt.Errorf("expected tx hash %v, get %v, nonce %d, to %s, value %s, gas %d, gasPrice %s, data %s",
				expectedHash.String(),
				actualTx.Hash().String(),
				expectedTx.Nonce(),
				expectedTx.To().String(),
				expectedTx.Value().String(),
				expectedTx.Gas(),
				expectedTx.GasPrice().String(),
				hex.EncodeToString(expectedTx.Data()),
			)
		}
		expectedTx = actualTx
		// move to next
		*receivedTxs = (*receivedTxs)[1:]
	}
	state.Prepare(expectedTx.Hash(), common.Hash{}, len(*txs))
	gasUsed, err := applyMessage(msg, state, header, chain.Config(), chainContext)
	if err != nil {
		return err
	}
	*txs = append(*txs, expectedTx)
	var root []byte
	if chain.Config().IsByzantium(header.Number) {
		state.Finalise(true)
	} else {
		root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number)).Bytes()
	}
	*usedGas += gasUsed
	receipt := types.NewReceipt(root, false, *usedGas)
	receipt.TxHash = expectedTx.Hash()
	receipt.GasUsed = gasUsed

	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = state.GetLogs(expectedTx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	receipt.BlockHash = state.BlockHash()
	receipt.BlockNumber = header.Number
	receipt.TransactionIndex = uint(state.TxIndex())
	*receipts = append(*receipts, receipt)
	state.SetNonce(msg.From(), nonce+1)
	return nil
}

// apply message
func applyMessage(
	msg callmsg,
	state *state.StateDB,
	header *types.Header,
	chainConfig *params.ChainConfig,
	chainContext core.ChainContext,
) (uint64, error) {
	// Create a new context to be used in the EVM environment
	context := core.NewEVMBlockContext(header, chainContext, nil)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(context, vm.TxContext{Origin: msg.From(), GasPrice: big.NewInt(0)}, state, chainConfig, vm.Config{})
	// Apply the transaction to the current state (included in the env)
	ret, returnGas, err := vmenv.Call(
		vm.AccountRef(msg.From()),
		*msg.To(),
		msg.Data(),
		msg.Gas(),
		msg.Value(),
	)
	if err != nil {
		log.Error("apply message failed", "msg", string(ret), "err", err)
	}
	return msg.Gas() - returnGas, err
}

func quickPackBlock(state *state.StateDB, chain consensus.ChainHeaderReader,
	header *types.Header, txs []*types.Transaction, receipts []*types.Receipt) *types.Block {

	var (
		block *types.Block
		root  common.Hash
	)

	// perform root calculation and block reorganization at the same time which with a large number of memory copy.
	// and reset the header root after actions done.
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
		wg.Done()
	}()
	go func() {
		// the header uncle hash will be settle as EmptyUncleHash which as the same of `nilUncleHash`
		block = types.NewBlock(header, txs, nil, receipts, trie.NewStackTrie(nil))
		wg.Done()
	}()
	wg.Wait()

	block.SetRoot(root)
	return block
}