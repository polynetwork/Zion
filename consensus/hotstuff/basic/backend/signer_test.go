package backend

import (
	"bytes"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/basic/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func newTestSigner() hotstuff.Signer {
	key, _ := generatePrivateKey()
	return NewSigner(key, 3)
}

func TestSign(t *testing.T) {
	s := newTestSigner()
	data := []byte("Here is a string....")
	sig, err := s.Sign(data)
	assert.NoError(t, err, "error mismatch: have %v, want nil", err)

	//Check signature recover
	hashData := crypto.Keccak256(data)
	pubkey, _ := crypto.Ecrecover(hashData, sig)
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])
	assert.Equal(t, signer, getAddress(), "address mismatch: have %v, want %s", signer.Hex(), getAddress().Hex())
}

func TestCheckValidatorSignature(t *testing.T) {
	vset, keys := newTestValidatorSet(5)

	// 1. Positive test: sign with validator's key should succeed
	data := []byte("dummy data")
	hashData := crypto.Keccak256([]byte(data))
	for i, k := range keys {
		// Sign
		sig, err := crypto.Sign(hashData, k)
		assert.NoError(t, err, "error mismatch: have %v, want nil", err)

		// CheckValidatorSignature should succeed
		signer := NewSigner(k, 3)
		addr, err := signer.CheckSignature(vset, data, sig)
		assert.NoError(t, err, "error mismatch: have %v, want nil", err)

		val := vset.GetByIndex(uint64(i))
		assert.Equal(t, addr, val.Address(), "validator address mismatch: have %v, want %v", addr, val.Address())
	}

	// 2. Negative test: sign with any key other than validator's key should return error
	key, err := crypto.GenerateKey()
	assert.NoError(t, err, "error mismatch: have %v, want nil", err)

	// Sign
	sig, err := crypto.Sign(hashData, key)
	assert.NoError(t, err, "error mismatch: have %v, want nil", err)

	// CheckValidatorSignature should return ErrUnauthorizedAddress
	signer := NewSigner(key, byte(core.MsgTypePrepareVote))
	addr, err := signer.CheckSignature(vset, data, sig)
	assert.Equal(t, err, hotstuff.ErrUnauthorizedAddress, "error mismatch: have %v, want %v", err, hotstuff.ErrUnauthorizedAddress)

	emptyAddr := common.Address{}
	assert.Equal(t, emptyAddr, common.Address{}, "address mismatch: have %v, want %v", addr, emptyAddr)
}

func TestPrepare(t *testing.T) {
	chain, engine := singleNodeChain()
	header := makeHeader(chain.Genesis(), engine.config)
	assert.NoError(t, engine.Prepare(chain, header))

	header.ParentHash = common.HexToHash("1234567890")
	assert.Equal(t, consensus.ErrUnknownAncestor, engine.Prepare(chain, header))
}

func TestVerifyHeader(t *testing.T) {
	chain, engine := singleNodeChain()

	// errEmptyCommittedSeals case
	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	block, _ = engine.UpdateBlock(block)
	err := engine.VerifyHeader(chain, block.Header(), true)
	assert.Equal(t, errEmptyCommittedSeals, err, "error mismatch")

	// short extra data
	header := block.Header()
	header.Extra = []byte{}
	err = engine.VerifyHeader(chain, header, false)
	assert.Equal(t, errInvalidExtraDataFormat, err, "error mismatch")

	// incorrect extra format
	header.Extra = []byte("0000000000000000000000000000000012300000000000000000000000000000000000000000000000000000000000000000")
	err = engine.VerifyHeader(chain, header, false)
	assert.Equal(t, errInvalidExtraDataFormat, err, "error mismatch")

	// non zero MixDigest
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.MixDigest = common.HexToHash("123456789")
	err = engine.VerifyHeader(chain, header, false)
	assert.Equal(t, errInvalidMixDigest, err, "error mismatch")

	// invalid uncles hash
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.UncleHash = common.HexToHash("123456789")
	err = engine.VerifyHeader(chain, header, false)
	assert.Equal(t, errInvalidUncleHash, err, "error mismatch")

	// invalid difficulty
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.Difficulty = big.NewInt(2)
	err = engine.VerifyHeader(chain, header, false)
	assert.Equal(t, errInvalidDifficulty, err, "error mismatch")

	// invalid timestamp
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.Time = chain.Genesis().Time() + (engine.config.BlockPeriod - 1)
	err = engine.VerifyHeader(chain, header, false)
	assert.Equal(t, errInvalidTimestamp, err, "error mismatch")
}

func TestVerifyHeaders(t *testing.T) {
	chain, engine := singleNodeChain()
	genesis := chain.Genesis()

	// success case
	headers := []*types.Header{}
	blocks := []*types.Block{}
	size := 100

	for i := 0; i < size; i++ {
		var b *types.Block
		if i == 0 {
			b = makeBlockWithoutSeal(chain, engine, genesis)
			b, _ = engine.UpdateBlock(b)
		} else {
			b = makeBlockWithoutSeal(chain, engine, blocks[i-1])
			b, _ = engine.UpdateBlock(b)
		}
		blocks = append(blocks, b)
		headers = append(headers, blocks[i].Header())
	}
	now = func() time.Time {
		return time.Unix(int64(headers[size-1].Time), 0)
	}
	_, results := engine.VerifyHeaders(chain, headers, nil)
	const timeoutDuration = 2 * time.Second
	timeout := time.NewTimer(timeoutDuration)
	index := 0
OUT1:
	for {
		select {
		case err := <-results:
			if err != nil {
				if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals {
					t.Errorf("error mismatch: have %v, want errEmptyCommittedSeals|errInvalidCommittedSeals", err)
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
	// abort cases
	abort, results := engine.VerifyHeaders(chain, headers, nil)
	timeout = time.NewTimer(timeoutDuration)
	index = 0
OUT2:
	for {
		select {
		case err := <-results:
			if err != nil {
				if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals {
					t.Errorf("error mismatch: have %v, want errEmptyCommittedSeals|errInvalidCommittedSeals", err)
					break OUT2
				}
			}
			index++
			if index == 5 {
				abort <- struct{}{}
			}
			if index >= size {
				t.Errorf("verifyheaders should be aborted")
				break OUT2
			}
		case <-timeout.C:
			break OUT2
		}
	}
	// error header cases
	headers[2].Number = big.NewInt(100)
	abort, results = engine.VerifyHeaders(chain, headers, nil)
	timeout = time.NewTimer(timeoutDuration)
	index = 0
	errors := 0
	expectedErrors := 2
OUT3:
	for {
		select {
		case err := <-results:
			if err != nil {
				if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals {
					errors++
				}
			}
			index++
			if index == size {
				if errors != expectedErrors {
					t.Errorf("error mismatch: have %v, want %v", err, expectedErrors)
				}
				break OUT3
			}
		case <-timeout.C:
			break OUT3
		}
	}
}

func TestPrepareExtra(t *testing.T) {
	validators := make([]common.Address, 4)

	// validators will be sorted asc
	validators[0] = common.BytesToAddress(hexutil.MustDecode("0x44add0ec310f115a0e603b2d7db9f067778eaf8a"))
	validators[1] = common.BytesToAddress(hexutil.MustDecode("0x294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212"))
	validators[2] = common.BytesToAddress(hexutil.MustDecode("0x6beaaed781d2d2ab6350f5c4566a2c6eaac407a6"))
	validators[3] = common.BytesToAddress(hexutil.MustDecode("0x8be76812f765c24641ec63dc2852b378aba2b440"))

	vanity := make([]byte, types.HotstuffExtraVanity)
	expectedResult := append(vanity, hexutil.MustDecode("0xf858f85494294fc7e8f22b3bcdcf955dd7ff3ba2ed833f82129444add0ec310f115a0e603b2d7db9f067778eaf8a946beaaed781d2d2ab6350f5c4566a2c6eaac407a6948be76812f765c24641ec63dc2852b378aba2b44080c0")...)
	h := &types.Header{
		Extra: vanity,
	}
	valSet := makeValSet(validators)
	payload, err := emptySigner.PrepareExtra(h, valSet)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, payload)

	// append useless information to extra-data
	h.Extra = append(vanity, make([]byte, 15)...)
	payload, err = emptySigner.PrepareExtra(h, valSet)
	assert.Equal(t, expectedResult, payload)
}

func TestFillExtraBeforeCommit(t *testing.T) {
	vanity := bytes.Repeat([]byte{0x00}, types.HotstuffExtraVanity)
	istRawData := hexutil.MustDecode("0xf858f8549444add0ec310f115a0e603b2d7db9f067778eaf8a94294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212946beaaed781d2d2ab6350f5c4566a2c6eaac407a6948be76812f765c24641ec63dc2852b378aba2b44080c0")
	expectedSeal := append([]byte{1, 2, 3}, bytes.Repeat([]byte{0x00}, types.HotstuffExtraSeal-3)...)
	expectedIstExtra := &types.HotstuffExtra{
		Validators: []common.Address{
			common.BytesToAddress(hexutil.MustDecode("0x44add0ec310f115a0e603b2d7db9f067778eaf8a")),
			common.BytesToAddress(hexutil.MustDecode("0x294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212")),
			common.BytesToAddress(hexutil.MustDecode("0x6beaaed781d2d2ab6350f5c4566a2c6eaac407a6")),
			common.BytesToAddress(hexutil.MustDecode("0x8be76812f765c24641ec63dc2852b378aba2b440")),
		},
		Seal:          expectedSeal,
		CommittedSeal: [][]byte{},
	}
	h := &types.Header{
		Extra: append(vanity, istRawData...),
	}

	// normal case
	assert.NoError(t, emptySigner.FillExtraBeforeCommit(h, expectedSeal))

	// verify istanbul extra-data
	istExtra, err := types.ExtractHotstuffExtra(h)
	assert.NoError(t, err)
	assert.Equal(t, expectedIstExtra, istExtra)

	// invalid seal
	unexpectedSeal := append(expectedSeal, make([]byte, 1)...)
	assert.Equal(t, errInvalidSignature, emptySigner.FillExtraBeforeCommit(h, unexpectedSeal))
}

func TestFillExtraAfterCommit(t *testing.T) {
	vanity := bytes.Repeat([]byte{0x00}, types.HotstuffExtraVanity)
	istRawData := hexutil.MustDecode("0xf858f8549444add0ec310f115a0e603b2d7db9f067778eaf8a94294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212946beaaed781d2d2ab6350f5c4566a2c6eaac407a6948be76812f765c24641ec63dc2852b378aba2b44080c0")
	expectedCommittedSeal := append([]byte{1, 2, 3}, bytes.Repeat([]byte{0x00}, types.HotstuffExtraSeal-3)...)
	expectedIstExtra := &types.HotstuffExtra{
		Validators: []common.Address{
			common.BytesToAddress(hexutil.MustDecode("0x44add0ec310f115a0e603b2d7db9f067778eaf8a")),
			common.BytesToAddress(hexutil.MustDecode("0x294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212")),
			common.BytesToAddress(hexutil.MustDecode("0x6beaaed781d2d2ab6350f5c4566a2c6eaac407a6")),
			common.BytesToAddress(hexutil.MustDecode("0x8be76812f765c24641ec63dc2852b378aba2b440")),
		},
		Seal:          []byte{},
		CommittedSeal: [][]byte{expectedCommittedSeal},
	}
	h := &types.Header{
		Extra: append(vanity, istRawData...),
	}

	// normal case
	assert.NoError(t, emptySigner.FillExtraAfterCommit(h, [][]byte{expectedCommittedSeal}))

	// verify istanbul extra-data
	istExtra, err := types.ExtractHotstuffExtra(h)
	assert.NoError(t, err)
	assert.Equal(t, expectedIstExtra, istExtra)

	// invalid seal
	unexpectedCommittedSeal := append(expectedCommittedSeal, make([]byte, 1)...)
	assert.Equal(t, errInvalidCommittedSeals, emptySigner.FillExtraAfterCommit(h, [][]byte{unexpectedCommittedSeal}))
}
