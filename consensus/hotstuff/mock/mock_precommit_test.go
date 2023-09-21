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

	"github.com/ethereum/go-ethereum/contracts/native/boot"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPreCommitCase1
// net scale is 4, leader send fake message of preCommit with wrong height, repos change view.
func TestMockPreCommitCase1(t *testing.T) {
	node_manager.InitABI()
	boot.InitNativeContracts()
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
			if ori.Code != core.MsgTypePreCommit {
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPreCommitCase2
// net scale is 4, leader send fake message of preCommit with wrong round, repos change view.
func TestMockPreCommitCase2(t *testing.T) {
	node_manager.InitABI()
	boot.InitNativeContracts()
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
			if ori.Code != core.MsgTypePreCommit {
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
			log.Info("fake message", "address", node.addr, "msg", msg.Code, "view", view, "msg", msg)
			return payload, true
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPreCommitCase3
// net scale is 4, leader send fake message of preCommit with wrong qc.height, repos change view.
func TestMockPreCommitCase3(t *testing.T) {
	node_manager.InitABI()
	boot.InitNativeContracts()
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
			if ori.Code != core.MsgTypePreCommit {
				return data, true
			}
			msg := ori.Copy()

			var qc QuorumCert
			if err := rlp.DecodeBytes(msg.Msg, &qc); err != nil {
				log.Error("failed to decode prepareQC", "err", err)
				return data, true
			}
			qc.View.Height = new(big.Int).SetUint64(fH)
			if raw, err := rlp.EncodeToBytes(qc); err != nil {
				log.Error("failed to encode prepareQC", "err", err)
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPreCommitCase4
// net scale is 4, leader send fake message of preCommit with wrong qc.height, repos change view.
func TestMockPreCommitCase4(t *testing.T) {
	node_manager.InitABI()
	boot.InitNativeContracts()
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
			if ori.Code != core.MsgTypePreCommit {
				return data, true
			}
			msg := ori.Copy()

			var qc QuorumCert
			if err := rlp.DecodeBytes(msg.Msg, &qc); err != nil {
				log.Error("failed to decode prepareQC", "err", err)
				return data, true
			}
			qc.View.Height = new(big.Int).SetUint64(fH)
			if raw, err := rlp.EncodeToBytes(qc); err != nil {
				log.Error("failed to encode prepareQC", "err", err)
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPreCommitCase5
// net scale is 4, leader send fake message of preCommit with wrong qc.round, repos change view.
func TestMockPreCommitCase5(t *testing.T) {
	node_manager.InitABI()
	boot.InitNativeContracts()
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
			if ori.Code != core.MsgTypePreCommit {
				return data, true
			}
			msg := ori.Copy()

			var qc QuorumCert
			if err := rlp.DecodeBytes(msg.Msg, &qc); err != nil {
				log.Error("failed to decode prepareQC", "err", err)
				return data, true
			}
			qc.View.Round = new(big.Int).SetUint64(fR)
			if raw, err := rlp.EncodeToBytes(qc); err != nil {
				log.Error("failed to encode prepareQC", "err", err)
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPreCommitCase6
// net scale is 4, leader send fake message of preCommit with wrong qc.digest, repos change view.
func TestMockPreCommitCase6(t *testing.T) {
	node_manager.InitABI()
	boot.InitNativeContracts()
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
			if ori.Code != core.MsgTypePreCommit {
				return data, true
			}
			msg := ori.Copy()

			var qc QuorumCert
			if err := rlp.DecodeBytes(msg.Msg, &qc); err != nil {
				log.Error("failed to decode prepareQC", "err", err)
				return data, true
			}
			qc.Node = common.HexToHash("0x12345")
			if raw, err := rlp.EncodeToBytes(qc); err != nil {
				log.Error("failed to encode prepareQC", "err", err)
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPreCommitCase7
// net scale is 4, leader send fake message of preCommit without enough qc.committedSeal, repos change view.
func TestMockPreCommitCase7(t *testing.T) {
	node_manager.InitABI()
	boot.InitNativeContracts()
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
			if ori.Code != core.MsgTypePreCommit {
				return data, true
			}
			msg := ori.Copy()

			var qc QuorumCert
			if err := rlp.DecodeBytes(msg.Msg, &qc); err != nil {
				log.Error("failed to decode prepareQC", "err", err)
				return data, true
			}
			qc.CommittedSeal = qc.CommittedSeal[:len(qc.CommittedSeal)-1]
			if raw, err := rlp.EncodeToBytes(qc); err != nil {
				log.Error("failed to encode prepareQC", "err", err)
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPreCommitCase8
// net scale is 4, leader send fake message of preCommit to some one repo, repos WONT change view.
func TestMockPreCommitCase8(t *testing.T) {
	node_manager.InitABI()
	boot.InitNativeContracts()
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
			if ori.Code != core.MsgTypePreCommit {
				return data, true
			}
			msg := ori.Copy()

			var qc QuorumCert
			if err := rlp.DecodeBytes(msg.Msg, &qc); err != nil {
				log.Error("failed to decode prepareQC", "err", err)
				return data, true
			}
			qc.CommittedSeal = qc.CommittedSeal[:len(qc.CommittedSeal)-1]
			if raw, err := rlp.EncodeToBytes(qc); err != nil {
				log.Error("failed to encode prepareQC", "err", err)
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPreCommitVoteCase1
// net scale is 4, leader send fake message of preCommitVote with wrong height. repos wont change view
func TestMockPreCommitVoteCase1(t *testing.T) {
	node_manager.InitABI()
	boot.InitNativeContracts()
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
			if ori.Code != core.MsgTypePreCommitVote {
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPreCommitVoteCase2
// net scale is 4, leader send fake message of preCommitVote with wrong height. repos change view
func TestMockPreCommitVoteCase2(t *testing.T) {
	node_manager.InitABI()
	boot.InitNativeContracts()
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
			if ori.Code != core.MsgTypePreCommitVote {
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPreCommitVoteCase3
// net scale is 4, leader send fake message of preCommitVote with wrong round. repos WONT change view
func TestMockPreCommitVoteCase3(t *testing.T) {
	node_manager.InitABI()
	boot.InitNativeContracts()
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
			if ori.Code != core.MsgTypePreCommitVote {
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPreCommitVoteCase4
// net scale is 4, leader send fake message of preCommitVote with wrong round. repos change view
func TestMockPreCommitVoteCase4(t *testing.T) {
	node_manager.InitABI()
	boot.InitNativeContracts()
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
			if ori.Code != core.MsgTypePreCommitVote {
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPreCommitVoteCase5
// net scale is 4, leader send fake message of preCommitVote with wrong digest. repos WONT change view
func TestMockPreCommitVoteCase5(t *testing.T) {
	node_manager.InitABI()
	boot.InitNativeContracts()
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
			if ori.Code != core.MsgTypePreCommitVote {
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPreCommitVoteCase6
// net scale is 4, leader send fake message of preCommitVote with wrong digest. repos change view
func TestMockPreCommitVoteCase6(t *testing.T) {
	node_manager.InitABI()
	boot.InitNativeContracts()
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
			if ori.Code != core.MsgTypePreCommitVote {
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
