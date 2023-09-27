// Copyright 2020 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"golang.org/x/crypto/sha3"
)

// TestStateProcessorErrors tests the output from the 'core' errors
// as defined in core/error.go. These errors are generated when the
// blockchain imports bad blocks, meaning blocks which have valid headers but
// contain invalid transactions
func TestStateProcessorErrors(t *testing.T) {
	var (
		signer     = types.HomesteadSigner{}
		testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		db         = rawdb.NewMemoryDatabase()
		gspec      = &Genesis{
			Config: params.TestChainConfig,
			Alloc:  GenesisAlloc{benchRootAddr: {Balance: benchRootFunds}},
		}
		genesis       = gspec.MustCommit(db)
		blockchain, _ = NewBlockChain(db, nil, gspec.Config, ethash.NewFaker(), vm.Config{}, nil, nil)
	)
	defer blockchain.Stop()
	var makeTx = func(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *types.Transaction {
		tx, _ := types.SignTx(types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data), signer, testKey)
		return tx
	}
	for i, tt := range []struct {
		txs  []*types.Transaction
		want string
	}{
		{
			txs: []*types.Transaction{
				makeTx(0, common.Address{}, big.NewInt(0), params.TxGas, big.NewInt(875000000), nil),
				makeTx(0, common.Address{}, big.NewInt(0), params.TxGas, big.NewInt(875000000), nil),
			},
			want: "could not apply tx 1 [0xf10328f0d8d80f42edac03c54540300eaf70bd555ec16b9ff4fcdbd703dd79b4]: nonce too low: address 0x71562b71999873DB5b286dF957af199Ec94617F7, tx: 0 state: 1",
		},
		{
			txs: []*types.Transaction{
				makeTx(100, common.Address{}, big.NewInt(0), params.TxGas, big.NewInt(875000000), nil),
			},
			want: "could not apply tx 0 [0xdc9395bdd9a34aa3e10cc781f9fcba36c47ab270519b571bb01ad9e99ae48d8d]: nonce too high: address 0x71562b71999873DB5b286dF957af199Ec94617F7, tx: 100 state: 0",
		},
		{
			txs: []*types.Transaction{
				makeTx(0, common.Address{}, big.NewInt(0), params.TxGas, big.NewInt(875000000), nil),
				makeTx(1, common.Address{}, big.NewInt(0), params.TxGas, big.NewInt(875000000), nil),
				makeTx(2, common.Address{}, big.NewInt(0), params.TxGas, big.NewInt(875000000), nil),
				makeTx(3, common.Address{}, big.NewInt(0), params.TxGas-1000, big.NewInt(875000000), nil),
			},
			want: "could not apply tx 3 [0xeab57cdda1267fc3b477ed95420885f8ad4fe895a034ace43f663aa0958cb3c3]: intrinsic gas too low: have 20000, want 21000",
		},
		// The last 'core' error is ErrGasUintOverflow: "gas uint64 overflow", but in order to
		// trigger that one, we'd have to allocate a _huge_ chunk of data, such that the
		// multiplication len(data) +gas_per_byte overflows uint64. Not testable at the moment
	} {
		block := GenerateBadBlock(genesis, ethash.NewFaker(), tt.txs, gspec.Config)
		_, err := blockchain.InsertChain(types.Blocks{block})
		if err == nil {
			t.Fatal("block imported without errors")
		}
		if have, want := err.Error(), tt.want; have != want {
			t.Errorf("test %d:\nhave \"%v\"\nwant \"%v\"\n", i, have, want)
		}
	}
}

// GenerateBadBlock constructs a "block" which contains the transactions. The transactions are not expected to be
// valid, and no proper post-state can be made. But from the perspective of the blockchain, the block is sufficiently
// valid to be considered for import:
// - valid pow (fake), ancestry, difficulty, gaslimit etc
func GenerateBadBlock(parent *types.Block, engine consensus.Engine, txs types.Transactions, config *params.ChainConfig) *types.Block {
	header := &types.Header{
		ParentHash: parent.Hash(),
		Coinbase:   parent.Coinbase(),
		Difficulty: engine.CalcDifficulty(&fakeChainReader{params.TestChainConfig}, parent.Time()+10, &types.Header{
			Number:     parent.Number(),
			Time:       parent.Time(),
			Difficulty: parent.Difficulty(),
			UncleHash:  parent.UncleHash(),
		}),
		GasLimit:  parent.GasLimit(),
		Number:    new(big.Int).Add(parent.Number(), common.Big1),
		Time:      parent.Time() + 10,
		UncleHash: types.EmptyUncleHash,
	}
	if config.IsLondon(header.Number) {
		header.BaseFee = misc.CalcBaseFee(config, parent.Header())
		parentGasLimit := parent.GasLimit()
		if !config.IsLondon(parent.Number()) {
			parentGasLimit = parent.GasLimit() * params.ElasticityMultiplier
		}
		header.GasLimit = CalcGasLimit1559(parentGasLimit, parentGasLimit)
	}
	var receipts []*types.Receipt

	// The post-state result doesn't need to be correct (this is a bad block), but we do need something there
	// Preferably something unique. So let's use a combo of blocknum + txhash
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(header.Number.Bytes())
	var cumulativeGas uint64
	for _, tx := range txs {
		txh := tx.Hash()
		hasher.Write(txh[:])
		receipt := types.NewReceipt(nil, false, cumulativeGas+tx.Gas())
		receipt.TxHash = tx.Hash()
		receipt.GasUsed = tx.Gas()
		receipts = append(receipts, receipt)
		cumulativeGas += tx.Gas()
	}
	header.Root = common.BytesToHash(hasher.Sum(nil))
	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts, trie.NewStackTrie(nil))
}
