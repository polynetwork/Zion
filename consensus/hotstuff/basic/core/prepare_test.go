package core

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/stretchr/testify/assert"
)

func TestHandlePrepare(t *testing.T) {
	N := uint64(4)
	F := uint64(1)
	H := uint64(5)
	R := uint64(1)

	newPrepareMsg := func(c *core) *MsgPrepare {
		coreView := c.currentView()
		h := coreView.Height.Uint64()
		r := coreView.Round.Uint64()

		view := makeView(h, r)
		highQC := newTestQC(c, h-1, r)
		proposal := makeBlockWithParentHash(int64(h), highQC.Hash)
		return &MsgPrepare{
			View:     view,
			Proposal: proposal,
			HighQC:   highQC,
		}
	}
	newP2PMsg := func(msg *MsgPrepare) *message {
		payload, _ := Encode(msg)
		return &message{
			Code: MsgTypePrepare,
			Msg:  payload,
		}
	}

	type testcase struct {
		Sys       *testSystem
		Msg       *message
		Leader    hotstuff.Validator
		ExpectErr error
	}
	testcases := []*testcase{
		// normal case
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			leader := sys.getLeader()
			val := validator.New(leader.Address())
			var data *MsgPrepare
			for _, backend := range sys.backends {
				core := backend.core()
				data = newPrepareMsg(core)
				core.current.SetPreCommittedQC(data.HighQC)
			}
			msg := newP2PMsg(data)
			return &testcase{
				Sys:       sys,
				Msg:       msg,
				Leader:    val,
				ExpectErr: nil,
			}
		}(),

		// errMsgOld
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			var data *MsgPrepare
			for _, backend := range sys.backends {
				core := backend.core()
				data = newPrepareMsg(core)
				core.current.height = new(big.Int).SetUint64(H + 1)
			}
			msg := newP2PMsg(data)
			leader := validator.New(sys.getLeader().Address())
			return &testcase{
				Sys:       sys,
				Msg:       msg,
				Leader:    leader,
				ExpectErr: errOldMessage,
			}
		}(),

		// errFutureMsg
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			var data *MsgPrepare
			for _, backend := range sys.backends {
				core := backend.core()
				data = newPrepareMsg(core)
				core.current.round = new(big.Int).SetUint64(R - 1)
			}
			msg := newP2PMsg(data)
			leader := validator.New(sys.getLeader().Address())
			return &testcase{
				Sys:       sys,
				Msg:       msg,
				Leader:    leader,
				ExpectErr: errFutureMessage,
			}
		}(),

		// errNotFromProposer
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			var data *MsgPrepare
			for _, backend := range sys.backends {
				core := backend.core()
				data = newPrepareMsg(core)
			}
			msg := newP2PMsg(data)
			wrongLeader := validator.New(sys.getRepos()[0].Address())
			return &testcase{
				Sys:       sys,
				Msg:       msg,
				Leader:    wrongLeader,
				ExpectErr: errNotFromProposer,
			}
		}(),

		// errExtend
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			leader := sys.getLeader()
			val := validator.New(leader.Address())
			var data *MsgPrepare
			for _, backend := range sys.backends {
				core := backend.core()
				data = newPrepareMsg(core)
				core.current.SetPreCommittedQC(data.HighQC)
			}
			// msg.proposal.parentHash not equal to the field of `lockedQC.Hash`
			data.HighQC.Hash = common.HexToHash("0x124")
			msg := newP2PMsg(data)
			return &testcase{
				Sys:       sys,
				Msg:       msg,
				Leader:    val,
				ExpectErr: errExtend,
			}
		}(),

		// errSafeNode
		func() *testcase {
			sys := NewTestSystemWithBackend(N, F, H, R)
			leader := sys.getLeader()
			val := validator.New(leader.Address())
			var data *MsgPrepare
			for _, backend := range sys.backends {
				core := backend.core()
				data = newPrepareMsg(core)
				// safety is false, and liveness false:
				// msg.proposal is not extend lockedQC
				// msg.highQC.view is smaller than lockedQC.view
				// or just set lockedQC is nil
				//core.current.SetPreCommittedQC(data.HighQC)
			}
			msg := newP2PMsg(data)
			return &testcase{
				Sys:       sys,
				Msg:       msg,
				Leader:    val,
				ExpectErr: errSafeNode,
			}
		}(),
	}

	for _, c := range testcases {
		for _, backend := range c.Sys.backends {
			core := backend.core()
			assert.Equal(t, c.ExpectErr, core.handlePrepare(c.Msg, c.Leader))
		}
	}
}
