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
	"crypto/rand"
	"encoding/json"
	"math/big"
	"os"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	testStateDB  *state.StateDB
	testEmptyCtx *native.NativeContract

	testSupplyGas    uint64 = 100000000000000000
	testGenesisNum   int    = 4
	testCaller       common.Address
	testGenesisEpoch *EpochInfo
)

func TestMain(m *testing.M) {
	db := rawdb.NewMemoryDatabase()
	testStateDB, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	testEmptyCtx = native.NewNativeContract(testStateDB, nil)
	testGenesisPeers := generateTestPeers(testGenesisNum)
	testGenesisEpoch, _ = storeGenesisEpoch(testStateDB, testGenesisPeers)
	InitNodeManager()

	os.Exit(m.Run())
}

// go test -v -count=1 github.com/ethereum/go-ethereum/contracts/native/governance/node_manager -run TestPropose
func TestPropose(t *testing.T) {
	type TestCase struct {
		BlockNum      int
		StartHeight   uint64
		Payload       []byte
		BeforeHandler func(c *TestCase, ctx *native.NativeContract)
		AfterHandler  func(c *TestCase, ctx *native.NativeContract)
		Index         int
		Expect        error
	}

	cases := []*TestCase{
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       1,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := generateTestPeers(2)
				input := &MethodProposeInput{StartHeight: c.StartHeight, Peers: peers}
				c.Payload, _ = input.Encode()
				delEpoch(ctx, testGenesisEpoch.Hash())
			},
			Expect: ErrEpochNotExist,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       2,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := generateTestPeers(2)
				input := &MethodProposeInput{StartHeight: c.StartHeight, Peers: peers}
				c.Payload, _ = input.Encode()
				testCaller = generateTestAddress(78)
			},
			Expect: ErrInvalidAuthority,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       3,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				input := &MethodProposeInput{StartHeight: 0, Peers: nil}
				payload, _ := input.Encode()
				c.Payload = payload[0 : len(payload)-2]
			},
			Expect: ErrInvalidInput,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       4,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := generateTestPeers(3)
				peers.List = nil
				input := &MethodProposeInput{StartHeight: 0, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrInvalidPeers,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       5,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := generateTestPeers(MinProposalPeersLen - 1)
				input := &MethodProposeInput{StartHeight: 0, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrPeersNum,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       6,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := generateTestPeers(MaxProposalPeersLen + 1)
				input := &MethodProposeInput{StartHeight: 0, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrPeersNum,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       7,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := generateTestPeers(MinProposalPeersLen + 1)
				peers.List[0].PubKey = "0ruf8nkj"
				input := &MethodProposeInput{StartHeight: 0, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrInvalidPubKey,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       8,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := testGenesisEpoch.Peers.Copy()
				newPeers := generateTestPeers(len(peers.List))
				for i, _ := range peers.List {
					if i%2 == 0 {
						peers.List[i] = newPeers.List[i]
					}
				}
				input := &MethodProposeInput{StartHeight: 0, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrOldParticipantsNumber,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       9,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := testGenesisEpoch.Peers.Copy()
				peers.List = append(peers.List, generateTestPeers(1).List...)
				c.StartHeight += MinEpochValidPeriod - 1
				input := &MethodProposeInput{StartHeight: c.StartHeight, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrProposalStartHeight,
		},
		{
			BlockNum:    3,
			StartHeight: 2,
			Index:       10,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := testGenesisEpoch.Peers.Copy()
				peers.List = append(peers.List, generateTestPeers(1).List...)
				c.StartHeight += MaxEpochValidPeriod + 10
				input := &MethodProposeInput{StartHeight: c.StartHeight, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrProposalStartHeight,
		},
		{
			BlockNum:    3,
			StartHeight: MinEpochValidPeriod + 10,
			Index:       11,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := testGenesisEpoch.Peers.Copy()
				peers.List = append(peers.List, generateTestPeers(1).List...)
				input := &MethodProposeInput{StartHeight: c.StartHeight, Peers: peers}
				sort.Sort(peers)
				epoch := &EpochInfo{
					ID:          testGenesisEpoch.ID + 1,
					StartHeight: c.StartHeight,
					Peers:       peers,
					Proposer:    ctx.ContractRef().TxOrigin(),
				}
				storeProposal(ctx, epoch.ID, epoch.Hash())
				c.Payload, _ = input.Encode()
			},
			Expect: ErrDuplicateProposal,
		},
		{
			BlockNum:    3,
			StartHeight: MinEpochValidPeriod + 10,
			Index:       12,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := testGenesisEpoch.Peers.Copy()
				peers.List = append(peers.List, generateTestPeers(1).List...)
				input := &MethodProposeInput{StartHeight: c.StartHeight, Peers: peers}
				sort.Sort(peers)
				proposer := ctx.ContractRef().TxOrigin()
				for i := 0; i < MaxProposalNumPerEpoch; i++ {
					epoch := &EpochInfo{
						ID:          testGenesisEpoch.ID + 1,
						StartHeight: c.StartHeight + 2 + uint64(i),
						Peers:       peers,
						Proposer:    proposer,
					}
					storeProposal(ctx, epoch.ID, epoch.Hash())
					storeEpoch(ctx, epoch)
				}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrProposalsNum,
		},
		{
			BlockNum:    3,
			StartHeight: MinEpochValidPeriod + 10,
			Index:       13,
			BeforeHandler: func(c *TestCase, ctx *native.NativeContract) {
				peers := testGenesisEpoch.Peers.Copy()
				peers.List = append(peers.List, generateTestPeers(1).List...)
				input := &MethodProposeInput{StartHeight: c.StartHeight, Peers: peers}
				c.Payload, _ = input.Encode()
			},
			Expect: nil,
		},
	}

	for _, v := range cases {
		resetTestContext()
		ctx := generateNativeContract(testCaller, v.BlockNum)
		if v.BeforeHandler != nil {
			v.BeforeHandler(v, ctx)
		}
		_, _, err := ctx.ContractRef().NativeCall(testCaller, this, v.Payload)
		assert.Equal(t, v.Expect, err)
		if v.AfterHandler != nil {
			v.AfterHandler(v, ctx)
		}
	}
}

// go test -v -count=1 github.com/ethereum/go-ethereum/contracts/native/governance/node_manager -run TestVote
func TestVote(t *testing.T) {
	epochID := uint64(2)

	type TestCase struct {
		ProposeBlockNum     int
		ProposalStartHeight uint64
		NeedProposal        bool
		Epoch               *EpochInfo
		OldMembers          []common.Address
		NewMembers          []common.Address
		Caller              common.Address
		Ctx                 *native.NativeContract
		VoteBlockNum        int
		Payload             []byte
		BeforeHandler       func(c *TestCase)
		AfterHandler        func()
		Index               int
		Expect              error
	}

	proposal := func(c *TestCase) {
		peers := testGenesisEpoch.Peers.Copy()
		c.OldMembers = make([]common.Address, 0)
		c.NewMembers = make([]common.Address, 0)
		for _, v := range peers.List {
			c.OldMembers = append(c.OldMembers, v.Address)
		}
		newList := generateTestPeers(1)
		for _, v := range newList.List {
			c.NewMembers = append(c.NewMembers, v.Address)
		}
		peers.List = append(peers.List, newList.List...)
		sort.Sort(peers)

		c.Epoch = &EpochInfo{StartHeight: c.ProposalStartHeight, Peers: peers, ID: epochID}
		input := &MethodProposeInput{StartHeight: c.Epoch.StartHeight, Peers: c.Epoch.Peers}
		payload, err := input.Encode()
		if err != nil {
			t.Fatal(err)
		}
		proposer := c.OldMembers[0]
		ctx := generateNativeContract(proposer, c.ProposeBlockNum)
		if _, _, err := ctx.ContractRef().NativeCall(proposer, this, payload); err != nil {
			t.Fatal(err)
		}
	}

	cases := []*TestCase{
		{
			ProposeBlockNum:     0,
			ProposalStartHeight: 0,
			NeedProposal:        false,
			VoteBlockNum:        3,
			Index:               1,
			BeforeHandler: func(c *TestCase) {
				c.Ctx = generateNativeContract(c.Caller, c.VoteBlockNum)
				input := &MethodVoteInput{EpochID: epochID, EpochHash: generateTestHash(12)}
				c.Payload, _ = input.Encode()

				delEpoch(c.Ctx, testGenesisEpoch.Hash())
			},
			Expect: ErrEpochNotExist,
		},
		{
			ProposeBlockNum:     3,
			ProposalStartHeight: MinEpochValidPeriod + 10,
			NeedProposal:        true,
			VoteBlockNum:        3,
			Index:               2,
			BeforeHandler: func(c *TestCase) {
				c.Caller = generateTestAddress(1)
				c.Ctx = generateNativeContract(c.Caller, c.VoteBlockNum)
				input := &MethodVoteInput{EpochID: c.Epoch.ID, EpochHash: c.Epoch.Hash()}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrInvalidAuthority,
		},
		{
			ProposeBlockNum:     3,
			ProposalStartHeight: MinEpochValidPeriod + 10,
			NeedProposal:        true,
			VoteBlockNum:        3,
			Index:               3,
			BeforeHandler: func(c *TestCase) {
				c.Caller = c.NewMembers[0]
				c.Ctx = generateNativeContract(c.Caller, c.VoteBlockNum)
				input := &MethodVoteInput{EpochID: c.Epoch.ID, EpochHash: c.Epoch.Hash()}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrInvalidAuthority,
		},
		{
			ProposeBlockNum:     3,
			ProposalStartHeight: MinEpochValidPeriod + 10,
			NeedProposal:        true,
			VoteBlockNum:        0,
			Index:               4,
			BeforeHandler: func(c *TestCase) {
				c.Caller = c.OldMembers[0]
				c.VoteBlockNum = int(c.ProposalStartHeight - MinVoteEffectivePeriod - 10)
				c.Ctx = generateNativeContract(c.Caller, c.VoteBlockNum)
				input := &MethodVoteInput{EpochID: c.Epoch.ID, EpochHash: c.Epoch.Hash()}
				payload, _ := input.Encode()
				c.Payload = payload[0 : len(payload)-2]
			},
			Expect: ErrInvalidInput,
		},
		{
			ProposeBlockNum:     3,
			ProposalStartHeight: MinEpochValidPeriod + 10,
			NeedProposal:        true,
			VoteBlockNum:        0,
			Index:               5,
			BeforeHandler: func(c *TestCase) {
				c.Caller = c.OldMembers[0]
				c.VoteBlockNum = int(c.ProposalStartHeight - MinVoteEffectivePeriod - 10)
				c.Ctx = generateNativeContract(c.Caller, c.VoteBlockNum)
				input := &MethodVoteInput{EpochID: c.Epoch.ID + 1, EpochHash: c.Epoch.Hash()}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrInvalidInput,
		},
		{
			ProposeBlockNum:     3,
			ProposalStartHeight: MinEpochValidPeriod + 10,
			NeedProposal:        true,
			VoteBlockNum:        0,
			Index:               6,
			BeforeHandler: func(c *TestCase) {
				c.Caller = c.OldMembers[0]
				c.VoteBlockNum = int(c.ProposalStartHeight - MinVoteEffectivePeriod - 10)
				c.Ctx = generateNativeContract(c.Caller, c.VoteBlockNum)
				input := &MethodVoteInput{EpochID: c.Epoch.ID, EpochHash: c.Epoch.Hash()}
				c.Payload, _ = input.Encode()

				delProposal(c.Ctx, c.Epoch.ID, c.Epoch.Hash())

			},
			Expect: ErrProposalNotExist,
		},
		{
			ProposeBlockNum:     3,
			ProposalStartHeight: MinEpochValidPeriod + 10,
			NeedProposal:        true,
			VoteBlockNum:        0,
			Index:               7,
			BeforeHandler: func(c *TestCase) {
				c.Caller = c.OldMembers[0]
				c.VoteBlockNum = int(c.ProposalStartHeight - MinVoteEffectivePeriod - 10)
				c.Ctx = generateNativeContract(c.Caller, c.VoteBlockNum)
				input := &MethodVoteInput{EpochID: c.Epoch.ID, EpochHash: c.Epoch.Hash()}
				c.Payload, _ = input.Encode()

				epoch, _ := getEpoch(c.Ctx, c.Epoch.Hash())
				epoch.Status = ProposalStatusPassed
				storeEpoch(c.Ctx, epoch)
			},
			Expect: ErrProposalPassed,
		},
		{
			ProposeBlockNum:     3,
			ProposalStartHeight: MinEpochValidPeriod + 10,
			NeedProposal:        true,
			VoteBlockNum:        0,
			Index:               8,
			BeforeHandler: func(c *TestCase) {
				c.Caller = c.OldMembers[0]
				c.VoteBlockNum = int(c.ProposalStartHeight - MinVoteEffectivePeriod + 1)
				c.Ctx = generateNativeContract(c.Caller, c.VoteBlockNum)
				input := &MethodVoteInput{EpochID: c.Epoch.ID, EpochHash: c.Epoch.Hash()}
				c.Payload, _ = input.Encode()
			},
			Expect: ErrVoteHeight,
		},
		{
			ProposeBlockNum:     3,
			ProposalStartHeight: MinEpochValidPeriod + 10,
			NeedProposal:        true,
			VoteBlockNum:        0,
			Index:               8,
			BeforeHandler: func(c *TestCase) {
				c.Caller = c.OldMembers[0]
				c.VoteBlockNum = int(c.ProposalStartHeight - MinVoteEffectivePeriod - 10)
				c.Ctx = generateNativeContract(c.Caller, c.VoteBlockNum)
				input := &MethodVoteInput{EpochID: c.Epoch.ID, EpochHash: c.Epoch.Hash()}
				c.Payload, _ = input.Encode()
			},
			Expect: nil,
		},
	}

	for _, v := range cases {
		resetTestContext()

		if v.NeedProposal {
			proposal(v)
		}

		if v.BeforeHandler != nil {
			v.BeforeHandler(v)
		}
		_, _, err := v.Ctx.ContractRef().NativeCall(v.Caller, this, v.Payload)
		assert.Equal(t, v.Expect, err)
		if v.AfterHandler != nil {
			v.AfterHandler()
		}
	}
}

func TestProposalPassed(t *testing.T) {
	resetTestContext()

	// prepare propose data
	peers := testGenesisEpoch.Peers.Copy()
	oldMembers := make([]common.Address, 0)
	for _, v := range peers.List {
		oldMembers = append(oldMembers, v.Address)
	}
	newList := generateTestPeers(1)
	peers.List = append(peers.List, newList.List...)
	sort.Sort(peers)

	epochID := uint64(2)
	proposeBlockNum := 9
	proposalStartHeight := uint64(proposeBlockNum) + MinEpochValidPeriod + 1
	epoch := &EpochInfo{StartHeight: proposalStartHeight, Peers: peers, ID: epochID, Status: ProposalStatusPropose}
	input := &MethodProposeInput{StartHeight: epoch.StartHeight, Peers: epoch.Peers}
	payload, _ := input.Encode()

	// propose
	proposer := oldMembers[0]
	ctx := generateNativeContract(proposer, proposeBlockNum)
	_, _, _ = ctx.ContractRef().NativeCall(proposer, this, payload)

	curEpoch, err := getCurrentEpoch(ctx)
	assert.NoError(t, err)

	// prepare vote data
	n := curEpoch.QuorumSize()
	voteBlockNum := proposeBlockNum + 1
	voteInput := &MethodVoteInput{EpochID: epoch.ID, EpochHash: epoch.Hash()}
	votePayload, err := voteInput.Encode()
	assert.NoError(t, err)

	// proposal not passed
	for i := 0; i < n-1; i++ {
		voter := oldMembers[i]
		ctx = generateNativeContract(voter, voteBlockNum)
		_, _, err := ctx.ContractRef().NativeCall(voter, this, votePayload)
		assert.NoError(t, err)
	}
	curEpoch, err = getEpoch(ctx, epoch.Hash())
	assert.NoError(t, err)
	assert.Equal(t, ProposalStatusPropose, curEpoch.Status)

	// proposal passed
	voter := oldMembers[n-1]
	ctx = generateNativeContract(voter, voteBlockNum)
	_, _, err = ctx.ContractRef().NativeCall(voter, this, votePayload)
	assert.NoError(t, err)
	curEpoch, err = getEpoch(ctx, epoch.Hash())
	assert.NoError(t, err)
	assert.Equal(t, ProposalStatusPassed, curEpoch.Status)

	// get changing epoch
	changingInput := &MethodGetChangingEpochInput{}
	changingEpochPayload, err := changingInput.Encode()
	assert.Nil(t, err)
	ctx = generateNativeContract(common.EmptyAddress, int(proposalStartHeight-1))
	enc, _, err := ctx.ContractRef().NativeCall(common.EmptyAddress, this, changingEpochPayload)
	changingOutPut := &MethodEpochOutput{}
	changingOutPut.Decode(enc)
	assert.NoError(t, err)
	assert.Equal(t, curEpoch.Hash(), changingOutPut.Epoch.Hash())
}

func TestDirtyJob(t *testing.T) {
	resetTestContext()

	s := testEmptyCtx
	epochID := uint64(2)
	peers := generateTestPeers(12)
	voters := []common.Address{peers.List[2].Address, peers.List[3].Address}

	// store last epoch
	lastEpoch := &EpochInfo{ID: epochID - 1, Proposer: peers.List[0].Address, Peers: peers, StartHeight: 60}
	assert.NoError(t, storeEpoch(s, lastEpoch))
	assert.NoError(t, storeProposal(s, lastEpoch.ID, lastEpoch.Hash()))

	// store current useless epoch and votes
	eps := []*EpochInfo{
		{ID: epochID, Proposer: peers.List[0].Address, Peers: &Peers{List: peers.List[:5]}, StartHeight: 270},
		{ID: epochID, Proposer: peers.List[1].Address, Peers: &Peers{List: peers.List[:6]}, StartHeight: 290},
	}
	for i, v := range eps {
		assert.NoError(t, storeEpoch(s, v))
		assert.NoError(t, storeProposal(s, v.ID, v.Hash()))
		assert.NoError(t, storeVote(s, v.Hash(), voters[i]))
		storeVoteTo(s, v.ID, voters[i], v.Hash())
	}

	curEpoch := generateTestEpochInfo(epochID, 270, 13)
	assert.NoError(t, storeEpoch(s, curEpoch))
	assert.NoError(t, storeProposal(s, curEpoch.ID, curEpoch.Hash()))

	// before dirty job
	list, err := getProposals(s, epochID)
	assert.NoError(t, err)
	assert.Equal(t, 1+len(eps), len(list))
	for _, v := range eps {
		assert.Equal(t, 1, voteSize(s, v.Hash()))
	}

	dirtyJob(s, lastEpoch, curEpoch)

	// after dirty job
	// proposal number should be only 1 with same epochID
	list, err = getProposals(s, epochID)
	assert.NoError(t, err)
	assert.Equal(t, int(1), len(list))
	assert.Equal(t, curEpoch.Hash(), list[0])

	for _, v := range eps {
		inf, _ := getEpoch(s, v.Hash())
		assert.Nil(t, inf)
		assert.Equal(t, 0, voteSize(s, v.Hash()))
	}
}

func TestGetEpochByID(t *testing.T) {
	resetTestContext()

	s := testEmptyCtx

	block := uint64(100)
	epoch := generateTestEpochInfo(2, block, 7)

	assert.NoError(t, storeEpoch(s, epoch))
	assert.NoError(t, storeProposal(s, epoch.ID, epoch.Hash()))
	storeEpochProof(s, epoch.ID, epoch.Hash())

	input := new(MethodGetEpochByIDInput)
	input.EpochID = epoch.ID
	payload, err := input.Encode()
	assert.NoError(t, err)
	ctx := generateNativeContract(common.EmptyAddress, int(block+1))
	enc, _, err := ctx.ContractRef().NativeCall(common.EmptyAddress, this, payload)
	assert.NoError(t, err)

	output := new(MethodEpochOutput)
	assert.NoError(t, output.Decode(enc))
	assert.Equal(t, epoch.Hash(), output.Epoch.Hash())
}

func TestGetEpochListJson(t *testing.T) {
	resetTestContext()

	s := testEmptyCtx
	epochID := uint64(2)
	peers := generateTestPeers(12)
	voters := []common.Address{peers.List[2].Address, peers.List[3].Address}

	// store last epoch
	lastEpoch := &EpochInfo{ID: epochID - 1, Proposer: peers.List[0].Address, Peers: peers, StartHeight: 60}
	assert.NoError(t, storeEpoch(s, lastEpoch))
	assert.NoError(t, storeProposal(s, lastEpoch.ID, lastEpoch.Hash()))

	// store current useless epoch and votes
	eps := []*EpochInfo{
		{ID: epochID, Proposer: peers.List[0].Address, Peers: &Peers{List: peers.List[:5]}, StartHeight: 270},
		{ID: epochID, Proposer: peers.List[1].Address, Peers: &Peers{List: peers.List[:6]}, StartHeight: 290},
	}
	for i, v := range eps {
		assert.NoError(t, storeEpoch(s, v))
		assert.NoError(t, storeProposal(s, v.ID, v.Hash()))
		assert.NoError(t, storeVote(s, v.Hash(), voters[i]))
		storeVoteTo(s, v.ID, voters[i], v.Hash())
	}

	input := new(MethodGetEpochListJsonInput)
	input.EpochID = epochID
	payload, err := input.Encode()
	assert.NoError(t, err)
	ctx := generateNativeContract(common.EmptyAddress, int(61))
	enc, _, err := ctx.ContractRef().NativeCall(common.EmptyAddress, this, payload)
	assert.NoError(t, err)

	output := new(MethodGetJsonOutput)
	assert.NoError(t, output.Decode(enc, node_manager_abi.MethodGetEpochListJson))
	var outputEpochs []*EpochInfo
	assert.NoError(t, json.Unmarshal([]byte(output.Result), &outputEpochs))
	assert.Equal(t, eps[0].Hash(), outputEpochs[0].Hash())
	assert.Equal(t, eps[1].Hash(), outputEpochs[1].Hash())
	assert.Equal(t, 2, len(outputEpochs))
}

func TestGetProofByID(t *testing.T) {
	resetTestContext()

	s := testEmptyCtx

	block := uint64(100)
	epoch := generateTestEpochInfo(2, block, 7)
	storeEpochProof(s, epoch.ID, epoch.Hash())

	input := new(MethodProofInput)
	input.EpochID = epoch.ID
	payload, err := input.Encode()
	assert.NoError(t, err)
	ctx := generateNativeContract(common.EmptyAddress, int(block+1))
	enc, _, err := ctx.ContractRef().NativeCall(common.EmptyAddress, this, payload)
	assert.NoError(t, err)

	output := new(MethodProofOutput)
	assert.NoError(t, output.Decode(enc))
	assert.Equal(t, epoch.Hash(), output.Hash)
}

func generateNativeContractRef(origin common.Address, blockNum int) *native.ContractRef {
	token := make([]byte, common.HashLength)
	rand.Read(token)
	hash := common.BytesToHash(token)
	return native.NewContractRef(testStateDB, origin, origin, big.NewInt(int64(blockNum)), hash, testSupplyGas, nil)
}

func generateNativeContract(origin common.Address, blockNum int) *native.NativeContract {
	ref := generateNativeContractRef(origin, blockNum)
	return native.NewNativeContract(testStateDB, ref)
}

func resetTestContext() {
	db := rawdb.NewMemoryDatabase()
	testStateDB, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	testEmptyCtx = native.NewNativeContract(testStateDB, nil)
	testGenesisPeers := generateTestPeers(testGenesisNum)
	testGenesisEpoch, _ = storeGenesisEpoch(testStateDB, testGenesisPeers)
	testCaller = testGenesisEpoch.Peers.List[0].Address
}

// generateTestPeer ONLY used for testing
func generateTestPeer() *PeerInfo {
	pk, _ := crypto.GenerateKey()
	return &PeerInfo{
		PubKey:  hexutil.Encode(crypto.CompressPubkey(&pk.PublicKey)),
		Address: crypto.PubkeyToAddress(pk.PublicKey),
	}
}

func generateTestPeers(n int) *Peers {
	peers := &Peers{List: make([]*PeerInfo, n)}
	for i := 0; i < n; i++ {
		peers.List[i] = generateTestPeer()
	}
	return peers
}

func generateTestEpochInfo(id, height uint64, peersNum int) *EpochInfo {
	epoch := new(EpochInfo)
	epoch.ID = id
	epoch.StartHeight = height
	epoch.Peers = generateTestPeers(peersNum)
	epoch.Proposer = epoch.Peers.List[0].Address
	return epoch
}

func generateTestHash(n int) common.Hash {
	data := big.NewInt(int64(n))
	return common.BytesToHash(data.Bytes())
}

func generateTestHashList(n int) *HashList {
	data := &HashList{List: make([]common.Hash, n)}
	for i := 0; i < n; i++ {
		data.List[i] = generateTestHash(i + 1)
	}
	return data
}

func generateTestAddress(n int) common.Address {
	data := big.NewInt(int64(n))
	return common.BytesToAddress(data.Bytes())
}

func generateTestAddressList(n int) *AddressList {
	data := &AddressList{List: make([]common.Address, n)}
	for i := 0; i < n; i++ {
		data.List[i] = generateTestAddress(i + 1)
	}
	return data
}
