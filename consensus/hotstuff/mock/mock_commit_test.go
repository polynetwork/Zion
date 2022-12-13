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

package mock

import (
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCommitCase1
// net scale is 4, leader send fake message of preCommit with wrong height, repos change view.
func TestMockCommitCase1(t *testing.T) {
	H, R, fH := uint64(4), uint64(0), uint64(5)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if !node.IsProposer() {
				return data, true
			}

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeCommit {
				return data, true
			}
			msg := ori.Copy()
			msg.View.Height = new(big.Int).SetUint64(fH)
			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}
			view := &core.View{
				Round:  new(big.Int).SetUint64(r),
				Height: new(big.Int).SetUint64(h),
			}
			log.Info("-----fake message", "address", node.addr, "msg", msg.Code, "view", view, "msg", msg)
			return payload, true
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCommitCase2
// net scale is 4, leader send fake message of preCommit with wrong round, repos change view.
func TestMockCommitCase2(t *testing.T) {
	H, R, fR := uint64(4), uint64(0), uint64(1)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if !node.IsProposer() {
				return data, true
			}

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeCommit {
				return data, true
			}
			msg := ori.Copy()
			msg.View.Round = new(big.Int).SetUint64(fR)
			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}
			view := &core.View{
				Round:  new(big.Int).SetUint64(r),
				Height: new(big.Int).SetUint64(h),
			}
			log.Info("-----fake message", "address", node.addr, "msg", msg.Code, "view", view, "msg", msg)
			return payload, true
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCommitCase3
// net scale is 4, leader send fake message of preCommit with wrong qc.height, repos change view.
func TestMockCommitCase3(t *testing.T) {
	H, R, fH := uint64(4), uint64(0), uint64(3)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if !node.IsProposer() {
				return data, true
			}

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeCommit {
				return data, true
			}
			msg := ori.Copy()

			var qc QuorumCert
			if err := rlp.DecodeBytes(msg.Msg, &qc); err != nil {
				log.Error("failed to decode lockQC", "err", err)
				return data, true
			}
			qc.View.Height = new(big.Int).SetUint64(fH)
			if raw, err := rlp.EncodeToBytes(qc); err != nil {
				log.Error("failed to encode lockQC", "err", err)
				return data, true
			} else {
				msg.Msg = raw
			}

			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}

			view := &core.View{
				Round:  new(big.Int).SetUint64(r),
				Height: new(big.Int).SetUint64(h),
			}
			log.Info("-----fake message", "address", node.addr, "msg", msg.Code, "view", view, "msg", msg)
			return payload, true
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCommitCase4
// net scale is 4, leader send fake message of preCommit with wrong qc.round, repos change view.
func TestMockCommitCase4(t *testing.T) {
	H, R, fR := uint64(4), uint64(0), uint64(1)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if !node.IsProposer() {
				return data, true
			}

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeCommit {
				return data, true
			}
			msg := ori.Copy()

			var qc QuorumCert
			if err := rlp.DecodeBytes(msg.Msg, &qc); err != nil {
				log.Error("failed to decode lockQC", "err", err)
				return data, true
			}
			qc.View.Round = new(big.Int).SetUint64(fR)
			if raw, err := rlp.EncodeToBytes(qc); err != nil {
				log.Error("failed to encode lockQC", "err", err)
				return data, true
			} else {
				msg.Msg = raw
			}

			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}

			view := &core.View{
				Round:  new(big.Int).SetUint64(r),
				Height: new(big.Int).SetUint64(h),
			}
			log.Info("-----fake message", "address", node.addr, "msg", msg.Code, "view", view, "msg", msg)
			return payload, true
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCommitCase5
// net scale is 4, leader send fake message of preCommit with wrong qc.digest, repos change view.
func TestMockCommitCase5(t *testing.T) {
	H, R := uint64(4), uint64(0)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if !node.IsProposer() {
				return data, true
			}

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeCommit {
				return data, true
			}
			msg := ori.Copy()

			var qc QuorumCert
			if err := rlp.DecodeBytes(msg.Msg, &qc); err != nil {
				log.Error("failed to decode lockQC", "err", err)
				return data, true
			}
			qc.Node = common.HexToHash("0x123")
			if raw, err := rlp.EncodeToBytes(qc); err != nil {
				log.Error("failed to encode lockQC", "err", err)
				return data, true
			} else {
				msg.Msg = raw
			}

			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}

			view := &core.View{
				Round:  new(big.Int).SetUint64(r),
				Height: new(big.Int).SetUint64(h),
			}
			log.Info("-----fake message", "address", node.addr, "msg", msg.Code, "view", view, "msg", msg)
			return payload, true
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCommitCase6
// net scale is 4, leader send fake message of preCommit without enough qc.committedSeal, repos change view.
func TestMockCommitCase6(t *testing.T) {
	H, R := uint64(4), uint64(0)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if !node.IsProposer() {
				return data, true
			}

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeCommit {
				return data, true
			}
			msg := ori.Copy()

			var qc QuorumCert
			if err := rlp.DecodeBytes(msg.Msg, &qc); err != nil {
				log.Error("failed to decode lockQC", "err", err)
				return data, true
			}
			qc.CommittedSeal = qc.CommittedSeal[:len(qc.CommittedSeal)-1]
			if raw, err := rlp.EncodeToBytes(qc); err != nil {
				log.Error("failed to encode lockQC", "err", err)
				return data, true
			} else {
				msg.Msg = raw
			}

			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}

			view := &core.View{
				Round:  new(big.Int).SetUint64(r),
				Height: new(big.Int).SetUint64(h),
			}
			log.Info("-----fake message", "address", node.addr, "msg", msg.Code, "view", view, "msg", msg)
			return payload, true
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCommitCase7
// net scale is 4, leader send fake message of preCommit to some one repo, repos WONT change view.
func TestMockCommitCase7(t *testing.T) {
	H, R, fN := uint64(4), uint64(0), int32(1)

	var locked int32
	atomic.StoreInt32(&locked, 0)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if !node.IsProposer() {
				return data, true
			}

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeCommit {
				return data, true
			}
			msg := ori.Copy()

			var qc QuorumCert
			if err := rlp.DecodeBytes(msg.Msg, &qc); err != nil {
				log.Error("failed to decode lockQC", "err", err)
				return data, true
			}
			qc.CommittedSeal = qc.CommittedSeal[:len(qc.CommittedSeal)-1]
			if raw, err := rlp.EncodeToBytes(qc); err != nil {
				log.Error("failed to encode lockQC", "err", err)
				return data, true
			} else {
				msg.Msg = raw
			}

			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}

			if value := atomic.LoadInt32(&locked); value >= fN {
				return data, true
			} else {
				atomic.StoreInt32(&locked, value+1)
				view := &core.View{
					Round:  new(big.Int).SetUint64(r),
					Height: new(big.Int).SetUint64(h),
				}
				log.Info("-----fake message", "address", node.addr, "msg", msg.Code, "view", view, "msg", msg)
				return payload, true
			}
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCommitVoteCase1
// net scale is 4, leader send fake message of commitVote with wrong height. repos wont change view
func TestMockCommitVoteCase1(t *testing.T) {
	H, R, fH, fN := uint64(4), uint64(0), uint64(5), int32(1)

	var locked int32
	atomic.StoreInt32(&locked, 0)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
				return data, true
			}

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeCommitVote {
				return data, true
			}
			msg := ori.Copy()
			msg.View.Height = new(big.Int).SetUint64(fH)

			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message", "err", err)
				return data, true
			}

			if value := atomic.LoadInt32(&locked); value >= fN {
				return data, true
			} else {
				atomic.StoreInt32(&locked, value+1)
				view := &core.View{
					Round:  new(big.Int).SetUint64(r),
					Height: new(big.Int).SetUint64(h),
				}
				log.Info("-----fake message", "address", node.addr, "msg", msg.Code, "view", view, "msg", msg)
				return payload, true
			}
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCommitVoteCase2
// net scale is 4, leader send fake message of commitVote with wrong height. repos change view
func TestMockCommitVoteCase2(t *testing.T) {
	H, R, fH, fN := uint64(4), uint64(0), uint64(5), int32(2)

	var locked int32
	atomic.StoreInt32(&locked, 0)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
				return data, true
			}

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeCommitVote {
				return data, true
			}
			msg := ori.Copy()
			msg.View.Height = new(big.Int).SetUint64(fH)

			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message", "err", err)
				return data, true
			}

			if value := atomic.LoadInt32(&locked); value >= fN {
				return data, true
			} else {
				atomic.StoreInt32(&locked, value+1)
				view := &core.View{
					Round:  new(big.Int).SetUint64(r),
					Height: new(big.Int).SetUint64(h),
				}
				log.Info("-----fake message", "address", node.addr, "msg", msg.Code, "view", view, "msg", msg)
				return payload, true
			}
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCommitVoteCase3
// net scale is 4, leader send fake message of commitVote with wrong round. repos WONT change view
func TestMockCommitVoteCase3(t *testing.T) {
	H, R, fR, fN := uint64(4), uint64(0), uint64(1), int32(1)

	var locked int32
	atomic.StoreInt32(&locked, 0)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
				return data, true
			}

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeCommitVote {
				return data, true
			}
			msg := ori.Copy()
			msg.View.Round = new(big.Int).SetUint64(fR)

			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message", "err", err)
				return data, true
			}

			if value := atomic.LoadInt32(&locked); value >= fN {
				return data, true
			} else {
				atomic.StoreInt32(&locked, value+1)
				view := &core.View{
					Round:  new(big.Int).SetUint64(r),
					Height: new(big.Int).SetUint64(h),
				}
				log.Info("-----fake message", "address", node.addr, "msg", msg.Code, "view", view, "msg", msg)
				return payload, true
			}
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCommitVoteCase4
// net scale is 4, leader send fake message of commitVote with wrong round. repos change view
func TestMockCommitVoteCase4(t *testing.T) {
	H, R, fR, fN := uint64(4), uint64(0), uint64(1), int32(2)

	var locked int32
	atomic.StoreInt32(&locked, 0)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
				return data, true
			}

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeCommitVote {
				return data, true
			}
			msg := ori.Copy()
			msg.View.Round = new(big.Int).SetUint64(fR)

			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message", "err", err)
				return data, true
			}

			if value := atomic.LoadInt32(&locked); value >= fN {
				return data, true
			} else {
				atomic.StoreInt32(&locked, value+1)
				view := &core.View{
					Round:  new(big.Int).SetUint64(r),
					Height: new(big.Int).SetUint64(h),
				}
				log.Info("-----fake message", "address", node.addr, "msg", msg.Code, "view", view, "msg", msg)
				return payload, true
			}
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCommitVoteCase5
// net scale is 4, leader send fake message of commitVote with wrong digest. repos WONT change view
func TestMockCommitVoteCase5(t *testing.T) {
	H, R, fN := uint64(4), uint64(0), int32(1)

	var locked int32
	atomic.StoreInt32(&locked, 0)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
				return data, true
			}

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeCommitVote {
				return data, true
			}
			msg := ori.Copy()
			msg.Msg = common.HexToHash("0x12346").Bytes()

			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message", "err", err)
				return data, true
			}

			if value := atomic.LoadInt32(&locked); value >= fN {
				return data, true
			} else {
				atomic.StoreInt32(&locked, value+1)
				view := &core.View{
					Round:  new(big.Int).SetUint64(r),
					Height: new(big.Int).SetUint64(h),
				}
				log.Info("-----fake message", "address", node.addr, "msg", msg.Code, "view", view, "msg", msg)
				return payload, true
			}
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCommitVoteCase6
// net scale is 4, leader send fake message of preCommitVote with wrong digest. repos change view
func TestMockCommitVoteCase6(t *testing.T) {
	H, R, fN := uint64(4), uint64(0), int32(2)

	var locked int32
	atomic.StoreInt32(&locked, 0)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
				return data, true
			}

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeCommitVote {
				return data, true
			}
			msg := ori.Copy()
			msg.Msg = common.HexToHash("0x12346").Bytes()

			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message", "err", err)
				return data, true
			}

			if value := atomic.LoadInt32(&locked); value >= fN {
				return data, true
			} else {
				atomic.StoreInt32(&locked, value+1)
				view := &core.View{
					Round:  new(big.Int).SetUint64(r),
					Height: new(big.Int).SetUint64(h),
				}
				log.Info("-----fake message", "address", node.addr, "msg", msg.Code, "view", view, "msg", msg)
				return payload, true
			}
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}
