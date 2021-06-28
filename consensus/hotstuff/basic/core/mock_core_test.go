package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"testing"
	"time"
)

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

	// leader := sys.getLeader()
	block := makeBlockWithParentHash(1, genesisBlock.Hash())
	for _, backend := range sys.backends {
		//go backend.EventMux().Post(hotstuff.RequestEvent{
		//	Proposal: block,
		//})
		backend.core().requests.StoreRequest(&hotstuff.Request{
			Proposal: block,
		})
	}

	close := sys.Run(true)
	defer close()

	<-time.After(2 * time.Second)

	for _, v := range sys.backends {
		if len(v.committedMsgs) > 1 {
			block := v.committedMsgs[1].commitProposal
			t.Logf("proposer %s committed block %d hash %s", v.address.Hex(), block.Number().Uint64(), block.Hash().Hex())
		}
	}
}
