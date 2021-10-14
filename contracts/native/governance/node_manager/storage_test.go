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

package node_manager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestStorageEpoch(t *testing.T) {
	expect := generateTestEpochInfo(1, 12, 100)

	assert.NoError(t, storeEpoch(testEmptyCtx, expect))
	got, err := getEpoch(testEmptyCtx, expect.Hash())
	assert.NoError(t, err)

	// must calculate hash before compare `got` and `expect`
	assert.Equal(t, expect.Hash(), got.Hash())
	assert.Equal(t, expect, got)

	// epoch info should be nil after delete
	delEpoch(testEmptyCtx, expect.Hash())
	got, err = getEpoch(testEmptyCtx, expect.Hash())
	assert.NotNil(t, err)
	assert.Nil(t, got)
}

func TestStorageEpochProof(t *testing.T) {
	startEpochProofHash := EpochProofHash(StartEpochID)
	assert.NotEqual(t, EpochProofDigest, startEpochProofHash)
	t.Logf("start epoch proof hash is %s", startEpochProofHash.Hex())

	epochID := uint64(13)
	expect := generateTestHash(15732478)

	storeEpochProof(testEmptyCtx, epochID, expect)
	got, err := getEpochProof(testEmptyCtx, epochID)
	assert.NoError(t, err)

	assert.Equal(t, expect, got)
}

func TestStorageProposal(t *testing.T) {
	expectSize := int(100)
	expectHashList := generateTestHashList(expectSize).List
	expectEpochID := uint64(2)
	expect := make([]common.Hash, 0)

	// test store proposal
	for _, hash := range expectHashList {
		expect = append(expect, hash)
		assert.NoError(t, storeProposal(testEmptyCtx, expectEpochID, hash))
	}
	assert.Equal(t, expectSize, totalProposalsNum(testEmptyCtx, expectEpochID))

	// test get proposals
	got, err := getProposals(testEmptyCtx, expectEpochID)
	assert.NoError(t, err)
	assert.Equal(t, expect, got)

	// test find proposal
	for _, v := range expect {
		assert.True(t, findProposal(testEmptyCtx, expectEpochID, v))
		assert.True(t, checkProposal(testEmptyCtx, expectEpochID, v))
	}

	// test delete proposal
	for _, v := range expect {
		assert.NoError(t, delProposal(testEmptyCtx, expectEpochID, v))
	}
	assert.Equal(t, int(0), totalProposalsNum(testEmptyCtx, expectEpochID))
}

func TestProposalsNum(t *testing.T) {
	s := testEmptyCtx
	epochID := uint64(2)
	proposer := generateTestAddress(12334)
	list := []*EpochInfo{
		&EpochInfo{ID: epochID, Proposer: proposer, Peers: generateTestPeers(11), StartHeight: 13},
		&EpochInfo{ID: epochID, Proposer: proposer, Peers: generateTestPeers(9), StartHeight: 14},
		&EpochInfo{ID: epochID, Proposer: proposer, Peers: generateTestPeers(8), StartHeight: 15},
		&EpochInfo{ID: epochID, Proposer: proposer, Peers: generateTestPeers(18), StartHeight: 13},
	}
	for _, v := range list {
		assert.NoError(t, storeProposal(s, v.ID, v.Hash()))
		assert.NoError(t, storeEpoch(s, v))
	}
	assert.Equal(t, len(list), proposalsNum(s, epochID, proposer))
}

func TestStorageVote(t *testing.T) {
	epochHash := generateTestHash(13)
	expectSize := int(200)
	expect := generateTestAddressList(expectSize).List
	expectDelVoter := expect[0]
	expectSizeAfterDel := expectSize - 1

	// test store vote
	for _, voter := range expect {
		assert.NoError(t, storeVote(testEmptyCtx, epochHash, voter))
	}

	// test vote size
	assert.Equal(t, expectSize, voteSize(testEmptyCtx, epochHash))
	got, err := getVotes(testEmptyCtx, epochHash)

	// test get votes
	assert.NoError(t, err)
	assert.Equal(t, expect, got)

	// test find vote
	for _, voter := range expect {
		assert.True(t, findVote(testEmptyCtx, epochHash, voter))
	}

	// test del vote
	assert.NoError(t, deleteVote(testEmptyCtx, epochHash, expectDelVoter))
	assert.Equal(t, expectSizeAfterDel, voteSize(testEmptyCtx, epochHash))

	// test clear votes
	clearVotes(testEmptyCtx, epochHash)
	assert.Equal(t, int(0), voteSize(testEmptyCtx, epochHash))
}

func TestSigns(t *testing.T) {
	expectSign := &ConsensusSign{Method: "test", Input: []byte("jalksdfj")}
	assert.NoError(t, storeSign(testEmptyCtx, expectSign))
	gotSign, err := getSign(testEmptyCtx, expectSign.Hash())
	assert.NoError(t, err)
	assert.Equal(t, expectSign.Hash(), gotSign.Hash())

	expectVoteNum := 10
	expectSigners := make([]common.Address, expectVoteNum)
	for i := 0; i < expectVoteNum; i++ {
		expectSigners[i] = generateTestAddress(i + 10)
		assert.NoError(t, storeSigner(testEmptyCtx, expectSign.Hash(), expectSigners[i]))
	}
	gotVoteNum := getSignerSize(testEmptyCtx, expectSign.Hash())
	assert.Equal(t, expectVoteNum, gotVoteNum)

	gotSigners, err := getSigners(testEmptyCtx, expectSign.Hash())
	assert.NoError(t, err)
	assert.Equal(t, expectSigners, gotSigners)
}
