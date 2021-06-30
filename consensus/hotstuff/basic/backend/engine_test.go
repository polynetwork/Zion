// Copyright 2017 The go-ethereum Authors
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

package backend
//
//import (
//	"bytes"
//	"math/big"
//	"reflect"
//	"testing"
//	"time"
//
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/common/hexutil"
//	"github.com/ethereum/go-ethereum/consensus"
//	"github.com/ethereum/go-ethereum/consensus/hotstuff"
//	"github.com/ethereum/go-ethereum/core/types"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestPrepare(t *testing.T) {
//	chain, engine := newBlockChain(1)
//	header := makeHeader(chain.Genesis(), engine.config)
//	assert.NoError(t, engine.Prepare(chain, header))
//
//	header.ParentHash = common.HexToHash("1234567890")
//	assert.Equal(t, consensus.ErrUnknownAncestor, engine.Prepare(chain, header))
//}
//
//func TestVerifyHeader(t *testing.T) {
//	chain, engine := newBlockChain(1)
//
//	// errEmptyCommittedSeals case
//	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
//	block, _ = engine.UpdateBlock(block)
//	err := engine.VerifyHeader(chain, block.Header(), false)
//	assert.Equal(t, errEmptyCommittedSeals, err, "error mismatch")
//
//	// short extra data
//	header := block.Header()
//	header.Extra = []byte{}
//	err = engine.VerifyHeader(chain, header, false)
//	assert.Equal(t, errInvalidExtraDataFormat, err, "error mismatch")
//
//	// incorrect extra format
//	header.Extra = []byte("0000000000000000000000000000000012300000000000000000000000000000000000000000000000000000000000000000")
//	err = engine.VerifyHeader(chain, header, false)
//	assert.Equal(t, errInvalidExtraDataFormat, err, "error mismatch")
//
//	// non zero MixDigest
//	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
//	header = block.Header()
//	header.MixDigest = common.HexToHash("123456789")
//	err = engine.VerifyHeader(chain, header, false)
//	assert.Equal(t, errInvalidMixDigest, err, "error mismatch")
//
//	// invalid uncles hash
//	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
//	header = block.Header()
//	header.UncleHash = common.HexToHash("123456789")
//	err = engine.VerifyHeader(chain, header, false)
//	assert.Equal(t, errInvalidUncleHash, err, "error mismatch")
//
//	// invalid difficulty
//	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
//	header = block.Header()
//	header.Difficulty = big.NewInt(2)
//	err = engine.VerifyHeader(chain, header, false)
//	assert.Equal(t, errInvalidDifficulty, err, "error mismatch")
//
//	// invalid timestamp
//	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
//	header = block.Header()
//	header.Time = chain.Genesis().Time() + (engine.config.BlockPeriod - 1)
//	err = engine.VerifyHeader(chain, header, false)
//	assert.Equal(t, errInvalidTimestamp, err, "error mismatch")
//}
//
//func TestVerifyHeaders(t *testing.T) {
//	chain, engine := newBlockChain(1)
//	genesis := chain.Genesis()
//
//	// success case
//	headers := []*types.Header{}
//	blocks := []*types.Block{}
//	size := 100
//
//	for i := 0; i < size; i++ {
//		var b *types.Block
//		if i == 0 {
//			b = makeBlockWithoutSeal(chain, engine, genesis)
//			b, _ = engine.UpdateBlock(b)
//		} else {
//			b = makeBlockWithoutSeal(chain, engine, blocks[i-1])
//			b, _ = engine.UpdateBlock(b)
//		}
//		blocks = append(blocks, b)
//		headers = append(headers, blocks[i].Header())
//	}
//	now = func() time.Time {
//		return time.Unix(int64(headers[size-1].Time), 0)
//	}
//	_, results := engine.VerifyHeaders(chain, headers, nil)
//	const timeoutDuration = 2 * time.Second
//	timeout := time.NewTimer(timeoutDuration)
//	index := 0
//OUT1:
//	for {
//		select {
//		case err := <-results:
//			if err != nil {
//				if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals {
//					t.Errorf("error mismatch: have %v, want errEmptyCommittedSeals|errInvalidCommittedSeals", err)
//					break OUT1
//				}
//			}
//			index++
//			if index == size {
//				break OUT1
//			}
//		case <-timeout.C:
//			break OUT1
//		}
//	}
//	// abort cases
//	abort, results := engine.VerifyHeaders(chain, headers, nil)
//	timeout = time.NewTimer(timeoutDuration)
//	index = 0
//OUT2:
//	for {
//		select {
//		case err := <-results:
//			if err != nil {
//				if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals {
//					t.Errorf("error mismatch: have %v, want errEmptyCommittedSeals|errInvalidCommittedSeals", err)
//					break OUT2
//				}
//			}
//			index++
//			if index == 5 {
//				abort <- struct{}{}
//			}
//			if index >= size {
//				t.Errorf("verifyheaders should be aborted")
//				break OUT2
//			}
//		case <-timeout.C:
//			break OUT2
//		}
//	}
//	// error header cases
//	headers[2].Number = big.NewInt(100)
//	abort, results = engine.VerifyHeaders(chain, headers, nil)
//	timeout = time.NewTimer(timeoutDuration)
//	index = 0
//	errors := 0
//	expectedErrors := 2
//OUT3:
//	for {
//		select {
//		case err := <-results:
//			if err != nil {
//				if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals {
//					errors++
//				}
//			}
//			index++
//			if index == size {
//				if errors != expectedErrors {
//					t.Errorf("error mismatch: have %v, want %v", err, expectedErrors)
//				}
//				break OUT3
//			}
//		case <-timeout.C:
//			break OUT3
//		}
//	}
//}
//
//func TestPrepareExtra(t *testing.T) {
//	validators := make([]common.Address, 4)
//	validators[0] = common.BytesToAddress(hexutil.MustDecode("0x44add0ec310f115a0e603b2d7db9f067778eaf8a"))
//	validators[1] = common.BytesToAddress(hexutil.MustDecode("0x294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212"))
//	validators[2] = common.BytesToAddress(hexutil.MustDecode("0x6beaaed781d2d2ab6350f5c4566a2c6eaac407a6"))
//	validators[3] = common.BytesToAddress(hexutil.MustDecode("0x8be76812f765c24641ec63dc2852b378aba2b440"))
//
//	vanity := make([]byte, types.HotstuffExtraVanity)
//	expectedResult := append(vanity, hexutil.MustDecode("0xf858f8549444add0ec310f115a0e603b2d7db9f067778eaf8a94294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212946beaaed781d2d2ab6350f5c4566a2c6eaac407a6948be76812f765c24641ec63dc2852b378aba2b44080c0")...)
//
//	h := &types.Header{
//		Extra: vanity,
//	}
//
//	payload, err := prepareExtra(h, validators)
//	assert.NoError(t, err)
//	assert.Equal(t, expectedResult, payload)
//
//	// append useless information to extra-data
//	h.Extra = append(vanity, make([]byte, 15)...)
//	payload, err = prepareExtra(h, validators)
//	assert.Equal(t, expectedResult, payload)
//}
//
//func TestWriteSeal(t *testing.T) {
//	vanity := bytes.Repeat([]byte{0x00}, types.HotstuffExtraVanity)
//	istRawData := hexutil.MustDecode("0xf858f8549444add0ec310f115a0e603b2d7db9f067778eaf8a94294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212946beaaed781d2d2ab6350f5c4566a2c6eaac407a6948be76812f765c24641ec63dc2852b378aba2b44080c0")
//	expectedSeal := append([]byte{1, 2, 3}, bytes.Repeat([]byte{0x00}, types.HotstuffExtraSeal-3)...)
//	expectedIstExtra := &types.HotstuffExtra{
//		Validators: []common.Address{
//			common.BytesToAddress(hexutil.MustDecode("0x44add0ec310f115a0e603b2d7db9f067778eaf8a")),
//			common.BytesToAddress(hexutil.MustDecode("0x294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212")),
//			common.BytesToAddress(hexutil.MustDecode("0x6beaaed781d2d2ab6350f5c4566a2c6eaac407a6")),
//			common.BytesToAddress(hexutil.MustDecode("0x8be76812f765c24641ec63dc2852b378aba2b440")),
//		},
//		Seal:          expectedSeal,
//		CommittedSeal: [][]byte{},
//	}
//	h := &types.Header{
//		Extra: append(vanity, istRawData...),
//	}
//
//	// normal case
//	assert.NoError(t, writeSeal(h, expectedSeal))
//
//	// verify istanbul extra-data
//	istExtra, err := types.ExtractHotstuffExtra(h)
//	assert.NoError(t, err)
//	assert.Equal(t, expectedIstExtra, istExtra)
//
//	// invalid seal
//	unexpectedSeal := append(expectedSeal, make([]byte, 1)...)
//	assert.Equal(t, errInvalidSignature, writeSeal(h, unexpectedSeal))
//}
//
//func TestWriteCommittedSeals(t *testing.T) {
//	vanity := bytes.Repeat([]byte{0x00}, types.HotstuffExtraVanity)
//	istRawData := hexutil.MustDecode("0xf858f8549444add0ec310f115a0e603b2d7db9f067778eaf8a94294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212946beaaed781d2d2ab6350f5c4566a2c6eaac407a6948be76812f765c24641ec63dc2852b378aba2b44080c0")
//	expectedCommittedSeal := append([]byte{1, 2, 3}, bytes.Repeat([]byte{0x00}, types.HotstuffExtraSeal-3)...)
//	expectedIstExtra := &types.HotstuffExtra{
//		Validators: []common.Address{
//			common.BytesToAddress(hexutil.MustDecode("0x44add0ec310f115a0e603b2d7db9f067778eaf8a")),
//			common.BytesToAddress(hexutil.MustDecode("0x294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212")),
//			common.BytesToAddress(hexutil.MustDecode("0x6beaaed781d2d2ab6350f5c4566a2c6eaac407a6")),
//			common.BytesToAddress(hexutil.MustDecode("0x8be76812f765c24641ec63dc2852b378aba2b440")),
//		},
//		Seal:          []byte{},
//		CommittedSeal: [][]byte{expectedCommittedSeal},
//	}
//	h := &types.Header{
//		Extra: append(vanity, istRawData...),
//	}
//
//	// normal case
//	assert.NoError(t, writeCommittedSeals(h, [][]byte{expectedCommittedSeal}))
//
//	// verify istanbul extra-data
//	istExtra, err := types.ExtractHotstuffExtra(h)
//	assert.NoError(t, err)
//	assert.Equal(t, expectedIstExtra, istExtra)
//
//	// invalid seal
//	unexpectedCommittedSeal := append(expectedCommittedSeal, make([]byte, 1)...)
//	assert.Equal(t, errInvalidCommittedSeals, writeCommittedSeals(h, [][]byte{unexpectedCommittedSeal}))
//}
//
//// TestSealStopChannel stop consensus before result committed
//func TestSealStopChannel(t *testing.T) {
//	chain, engine := newBlockChain(4)
//	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
//	stop := make(chan struct{}, 1)
//	eventSub := engine.EventMux().Subscribe(hotstuff.RequestEvent{})
//	eventLoop := func() {
//		ev := <-eventSub.Chan()
//		_, ok := ev.Data.(hotstuff.RequestEvent)
//		if !ok {
//			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
//		}
//		stop <- struct{}{}
//		eventSub.Unsubscribe()
//	}
//	go eventLoop()
//	resultCh := make(chan *types.Block, 10)
//	go func() {
//		err := engine.Seal(chain, block, resultCh, stop)
//		if err != nil {
//			t.Errorf("error mismatch: have %v, want nil", err)
//		}
//	}()
//
//	finalBlock := <-resultCh
//	assert.Nil(t, finalBlock)
//}
//
//func TestSealPreCommitOtherHash(t *testing.T) {
//	chain, engine := newBlockChain(4)
//	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
//	otherBlock := makeBlockWithoutSeal(chain, engine, block)
//	expectedCommittedSeal := append([]byte{1, 2, 3}, bytes.Repeat([]byte{0x00}, types.HotstuffExtraSeal-3)...)
//
//	eventSub := engine.EventMux().Subscribe(hotstuff.RequestEvent{})
//	blockOutputChannel := make(chan *types.Block)
//	stopChannel := make(chan struct{})
//
//	go func() {
//		ev := <-eventSub.Chan()
//		if _, ok := ev.Data.(hotstuff.RequestEvent); !ok {
//			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
//		}
//		view := &hotstuff.View{
//			Round:  new(big.Int).SetUint64(0),
//			Height: otherBlock.Number(),
//		}
//		if _, _, err := engine.PreCommit(view, otherBlock, [][]byte{expectedCommittedSeal}); err != nil {
//			t.Error(err.Error())
//		}
//		eventSub.Unsubscribe()
//	}()
//
//	go func() {
//		if err := engine.Seal(chain, block, blockOutputChannel, stopChannel); err != nil {
//			t.Error(err.Error())
//		}
//	}()
//
//	select {
//	case <-blockOutputChannel:
//		t.Error("Wrong block found!")
//	default:
//		//no block found, stop the sealing
//		close(stopChannel)
//	}
//
//	output := <-blockOutputChannel
//	if output != nil {
//		t.Error("Block not nil!")
//	}
//}
//
//func TestSealCommitted(t *testing.T) {
//	chain, engine := newBlockChain(1)
//	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
//	expectedBlock, _ := engine.UpdateBlock(block)
//	resultCh := make(chan *types.Block, 10)
//	go func() {
//		err := engine.Seal(chain, block, resultCh, make(chan struct{}))
//
//		if err != nil {
//			t.Errorf("error mismatch: have %v, want %v", err, expectedBlock)
//		}
//	}()
//
//	finalBlock := <-resultCh
//	t.Logf("final blcok %s", finalBlock.Hash())
//	// todo
//	//if finalBlock.Hash() != expectedBlock.Hash() {
//	//	t.Errorf("hash mismatch: have %v, want %v", finalBlock.Hash(), expectedBlock.Hash())
//	//}
//}
