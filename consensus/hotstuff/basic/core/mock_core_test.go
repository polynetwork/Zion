package core

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/basic/core -run TestNewRound
func TestNewRound(t *testing.T) {
	N := uint64(4)
	F := uint64(1)
	H := uint64(1)
	R := uint64(0)

	needBroadCast = true
	sys := NewTestSystemWithBackend(N, F, H, R)

	// prepare genesis qc
	lastLeader := sys.getLeaderByRound(EmptyAddress, common.Big0)
	genesisBlock, _ := newProposalAndQC(lastLeader, 0, 0)
	for _, v := range sys.backends {
		v.committedMsgs = append(v.committedMsgs, testCommittedMsgs{
			commitProposal: genesisBlock,
		})
		v.core().current = nil
	}

	close := sys.Run(true)
	defer close()

	block := makeBlockWithParentHash(1, genesisBlock.Hash())
	for _, backend := range sys.backends {
		go backend.EventMux().Post(hotstuff.RequestEvent{
			Proposal: block,
		})
	}

	<-time.After(10 * time.Second)

	for _, v := range sys.backends {
		if len(v.committedMsgs) > 1 {
			block := v.committedMsgs[1].commitProposal
			t.Logf("proposer %s committed block %d hash %s", v.address.Hex(), block.Number().Uint64(), block.Hash().Hex())
		}
	}
}
