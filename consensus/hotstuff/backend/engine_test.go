package backend

import (
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	tu "github.com/ethereum/go-ethereum/consensus/hotstuff/testutils"
	"github.com/ethereum/go-ethereum/contracts/native/boot"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	boot.InitNativeContracts()
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestGenesisBlock
func TestGenesisBlock(t *testing.T) {
	memDB := rawdb.NewMemoryDatabase()
	g, _, err := tu.GenesisAndKeys(4)
	if err != nil {
		t.Error(err)
	}
	block := tu.GenesisBlock(g, memDB)
	raw, err := block.Header().MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	t.Logf("genesis block, header %s, mixDigest %v, number %d, hash %v, size %v",
		string(raw), block.MixDigest(), block.NumberU64(), block.Hash(), block.Size())
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestPrepare
func TestPrepare(t *testing.T) {
	chain, engine := singleNodeChain()
	defer engine.Stop()

	header := makeHeader(chain.Genesis())
	assert.NoError(t, engine.Prepare(chain, header))

	header.ParentHash = common.HexToHash("0x1234567890")
	assert.Error(t, engine.Prepare(chain, header), consensus.ErrUnknownAncestor)
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestSealStopChannel
// TestSealStopChannel stop seal and result channel should be empty
func TestSealStopChannel(t *testing.T) {
	chain, engine := singleNodeChain()
	defer engine.Stop()

	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	stop := make(chan struct{}, 1)
	eventSub := engine.EventMux().Subscribe(hotstuff.RequestEvent{})
	eventLoop := func() {
		ev := <-eventSub.Chan()
		if _, ok := ev.Data.(hotstuff.RequestEvent); !ok {
			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
		}
		stop <- struct{}{}
		eventSub.Unsubscribe()
	}
	resultCh := make(chan *types.Block, 10)
	go func() {
		if err := engine.Seal(chain, block, resultCh, stop); err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
	}()
	go eventLoop()

	finalBlock := <-resultCh
	if finalBlock != nil {
		t.Errorf("block mismatch: have %v, want nil", finalBlock)
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestSealCommittedOtherHash
// TestSealCommittedOtherHash result channel should be empty if engine commit another block before seal
func TestSealCommittedOtherHash(t *testing.T) {
	chain, engine := singleNodeChain()
	defer engine.Stop()
	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	otherBlock := makeBlockWithoutSeal(chain, engine, block)

	eventSub := engine.EventMux().Subscribe(hotstuff.RequestEvent{})
	blockOutputChannel := make(chan *types.Block)
	stopChannel := make(chan struct{})

	go func() {
		ev := <-eventSub.Chan()
		if _, ok := ev.Data.(hotstuff.RequestEvent); !ok {
			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
		}
		if err := engine.Commit(otherBlock); err != nil {
			t.Error(err.Error())
		}
		eventSub.Unsubscribe()
	}()

	go func() {
		if err := engine.Seal(chain, block, blockOutputChannel, stopChannel); err != nil {
			t.Error(err.Error())
		}
	}()

	select {
	case <-blockOutputChannel:
		t.Error("Wrong block found!")
	default:
		//no block found, stop the sealing
		close(stopChannel)
	}

	output := <-blockOutputChannel
	if output != nil {
		t.Error("Block not nil!")
	}
}

func updateTestBlock(block *types.Block, addr common.Address) *types.Block {
	header := block.Header()
	header.Coinbase = addr
	return block.WithSeal(header)
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestSealCommitted
// TestSealCommitted block hash WONT change after seal.
func TestSealCommitted(t *testing.T) {
	chain, engine := singleNodeChain()
	defer engine.Stop()

	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	expectedBlock := updateTestBlock(block, engine.Address())

	resultCh := make(chan *types.Block, 10)
	go func() {
		if err := engine.Seal(chain, block, resultCh, make(chan struct{})); err != nil {
			t.Errorf("error mismatch: have %v, want %v", err, expectedBlock)
		}
	}()

	finalBlock := <-resultCh
	if finalBlock.Hash() != expectedBlock.Hash() {
		t.Errorf("hash mismatch: have %v, want %v", finalBlock.Hash(), expectedBlock.Hash())
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestVerifyHeader
func TestVerifyHeader(t *testing.T) {
	chain, engine := singleNodeChain()
	defer engine.Stop()

	// errEmptyCommittedSeals case
	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header := engine.chain.GetHeader(block.ParentHash(), block.NumberU64()-1)
	block = updateTestBlock(block, engine.Address())
	if err := engine.VerifyHeader(chain, block.Header(), false); err != errEmptyCommittedSeals {
		t.Errorf("error mismatch: have %v, want %v", err, errEmptyCommittedSeals)
	}

	// short extra data
	header = block.Header()
	header.Extra = []byte{}
	if err := engine.VerifyHeader(chain, header, false); err != errInvalidExtraDataFormat {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidExtraDataFormat)
	}
	// incorrect extra format
	header.Extra = []byte("0000000000000000000000000000000012300000000000000000000000000000000000000000000000000000000000000000")
	if err := engine.VerifyHeader(chain, header, false); err != errInvalidExtraDataFormat {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidExtraDataFormat)
	}

	// non zero MixDigest
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.MixDigest = common.HexToHash("0x123456789")
	if err := engine.VerifyHeader(chain, header, false); err != errInvalidMixDigest {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidMixDigest)
	}

	// invalid uncles hash
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.UncleHash = common.HexToHash("0x123456789")
	if err := engine.VerifyHeader(chain, header, false); err != errInvalidUncleHash {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidUncleHash)
	}

	// invalid difficulty
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.Difficulty = big.NewInt(2)
	if err := engine.VerifyHeader(chain, header, false); err != errInvalidDifficulty {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidDifficulty)
	}

	// invalid timestamp
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.Time = chain.Genesis().Time() + (engine.config.BlockPeriod - 1)

	if err := engine.VerifyHeader(chain, header, false); err != errInvalidTimestamp {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidTimestamp)
	}

	// future block
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.Time = uint64(time.Now().Unix() + 10)
	if err := engine.VerifyHeader(chain, header, false); err != consensus.ErrFutureBlock {
		t.Errorf("error mismatch: have %v, want %v", err, consensus.ErrFutureBlock)
	}
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/backend -run TestVerifyHeaders
func TestVerifyHeaders(t *testing.T) {
	chain, engine := singleNodeChain()
	defer engine.Stop()
	genesis := chain.Genesis()

	// success case
	headers := []*types.Header{}
	blocks := []*types.Block{}
	size := 100

	for i := 0; i < size; i++ {
		var b *types.Block
		if i == 0 {
			b = makeBlockWithoutSeal(chain, engine, genesis)
			b = updateTestBlock(b, engine.Address())
		} else {
			b = makeBlockWithoutSeal(chain, engine, blocks[i-1])
			b = updateTestBlock(b, engine.Address())
		}
		blocks = append(blocks, b)
		headers = append(headers, blocks[i].Header())
	}

	_, results := engine.VerifyHeaders(chain, headers, nil)
	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	index := 0
OUT1:
	for {
		select {
		case err := <-results:
			if err != nil {
				if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals && err != consensus.ErrUnknownAncestor {
					t.Errorf("error mismatch: have %v, want istanbulcommon.ErrEmptyCommittedSeals|istanbulcommon.ErrInvalidCommittedSeals|ErrUnknownAncestor", err)
					break OUT1
				}
			}
			index++
			if index == size {
				break OUT1
			}
		case <-timeout.C:
			break OUT1
		}
	}
	_, results = engine.VerifyHeaders(chain, headers, nil)
	timeout = time.NewTimer(timeoutDura)
OUT2:
	for {
		select {
		case err := <-results:
			if err != nil {
				if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals && err != consensus.ErrUnknownAncestor {
					t.Errorf("error mismatch: have %v, want istanbulcommon.ErrEmptyCommittedSeals|istanbulcommon.ErrInvalidCommittedSeals|ErrUnknownAncestor", err)
					break OUT2
				}
			}
		case <-timeout.C:
			break OUT2
		}
	}
	// error header cases
	headers[2].Number = big.NewInt(100)
	_, results = engine.VerifyHeaders(chain, headers, nil)
	timeout = time.NewTimer(timeoutDura)
	index = 0
	errors := 0
	expectedErrors := 0
OUT3:
	for {
		select {
		case err := <-results:
			if err != nil {
				if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals && err != consensus.ErrUnknownAncestor {
					errors++
				}
			}
			index++
			if index == size {
				if errors != expectedErrors {
					t.Errorf("error mismatch: have %v, want %v", errors, expectedErrors)
				}
				break OUT3
			}
		case <-timeout.C:
			break OUT3
		}
	}
}
