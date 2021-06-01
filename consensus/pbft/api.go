package pbft

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

func (b *backend) Author(header *types.Header) (common.Address, error) {
	return common.Address{}, nil
}

func (b *backend) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	return nil
}

func (b *backend) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	return nil, nil
}

func (b *backend) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	return nil
}

func (b *backend) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	return nil
}

func (b *backend) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header) {

}

func (b *backend) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	return nil, nil
}

func (b *backend) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	return nil
}

func (b *backend) SealHash(header *types.Header) common.Hash {
	return common.Hash{}
}

func (b *backend) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return nil
}

func (b *backend) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return nil
}

func (b *backend) Close() error {
	return nil
}
