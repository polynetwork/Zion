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
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPrepareCase1
// net scale is 4, leader send fake message of newView with wrong height, repos change view.
func TestMockPrepareCase1(t *testing.T) {
	H, R, fH, fN := uint64(4), uint64(0), uint64(5), int(1)
	fakeNodes := make(map[common.Address]struct{})

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if !node.IsProposer() {
				return data, true
			}
			if _, ok := fakeNodes[node.addr]; ok {
				return data, true
			}
			if len(fakeNodes) >= fN {
				return data, true
			}
			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypePrepare {
				return data, true
			}
			msg := ori.Copy()
			msg.View.Height = new(big.Int).SetUint64(fH)
			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}
			fakeNodes[node.addr] = struct{}{}
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPrepareCase2
// net scale is 4, leader send fake message of newView with wrong height, repos change view.
func TestMockPrepareCase2(t *testing.T) {
	H, R, fR, fN := uint64(4), uint64(0), uint64(1), int(1)
	fakeNodes := make(map[common.Address]struct{})

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if !node.IsProposer() {
				return data, true
			}
			if _, ok := fakeNodes[node.addr]; ok {
				return data, true
			}
			if len(fakeNodes) >= fN {
				return data, true
			}
			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypePrepare {
				return data, true
			}
			msg := ori.Copy()
			msg.View.Round = new(big.Int).SetUint64(fR)
			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}
			fakeNodes[node.addr] = struct{}{}
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPrepareCase3
// net scale is 4, leader send fake message of newView with wrong qc.view.height, repos change view.
func TestMockPrepareCase3(t *testing.T) {
	H, R, fH, fN := uint64(4), uint64(0), uint64(4), int(1)
	fakeNodes := make(map[common.Address]struct{})

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if !node.IsProposer() {
				return data, true
			}
			if _, ok := fakeNodes[node.addr]; ok {
				return data, true
			}
			if len(fakeNodes) >= fN {
				return data, true
			}
			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypePrepare {
				return data, true
			}
			msg := ori.Copy()
			var sub core.Subject
			if err := rlp.DecodeBytes(msg.Msg, &sub); err != nil {
				log.Error("failed to decode subject", "err", err)
				return data, true
			}
			var qc QuorumCert
			if raw, err := rlp.EncodeToBytes(sub.QC); err != nil {
				log.Error("failed to encode prepareQC", "err", err)
				return data, true
			} else if err := rlp.DecodeBytes(raw, &qc); err != nil {
				log.Error("failed to decode prepareQC", "err", err)
				return data, true
			}
			qc.View.Height = new(big.Int).SetUint64(fH)
			var newSub = struct {
				Node *core.Node
				QC   *QuorumCert
			}{
				sub.Node,
				&qc,
			}
			if raw, err := rlp.EncodeToBytes(newSub); err != nil {
				log.Error("failed to encode new subject", "err", err)
				return data, true
			} else {
				msg.Msg = raw
			}

			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}
			fakeNodes[node.addr] = struct{}{}
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPrepareCase4
// net scale is 4, leader send fake message of newView with wrong qc.view.round, repos change view.
func TestMockPrepareCase4(t *testing.T) {
	H, R, fR, fN := uint64(4), uint64(0), uint64(1), int(1)
	fakeNodes := make(map[common.Address]struct{})

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if !node.IsProposer() {
				return data, true
			}
			if _, ok := fakeNodes[node.addr]; ok {
				return data, true
			}
			if len(fakeNodes) >= fN {
				return data, true
			}
			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypePrepare {
				return data, true
			}
			msg := ori.Copy()
			var sub core.Subject
			if err := rlp.DecodeBytes(msg.Msg, &sub); err != nil {
				log.Error("failed to decode subject", "err", err)
				return data, true
			}
			var qc QuorumCert
			if raw, err := rlp.EncodeToBytes(sub.QC); err != nil {
				log.Error("failed to encode prepareQC", "err", err)
				return data, true
			} else if err := rlp.DecodeBytes(raw, &qc); err != nil {
				log.Error("failed to decode prepareQC", "err", err)
				return data, true
			}
			qc.View.Round = new(big.Int).SetUint64(fR)
			var newSub = struct {
				Node *core.Node
				QC   *QuorumCert
			}{
				sub.Node,
				&qc,
			}
			if raw, err := rlp.EncodeToBytes(newSub); err != nil {
				log.Error("failed to encode new subject", "err", err)
				return data, true
			} else {
				msg.Msg = raw
			}

			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}
			fakeNodes[node.addr] = struct{}{}
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockPrepareCase5
// net scale is 4, leader send fake message of newView with wrong qc.hash, repos change view.
func TestMockPrepareCase5(t *testing.T) {
	H, R, fN := uint64(4), uint64(0), int(1)
	fakeNodes := make(map[common.Address]struct{})

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if !node.IsProposer() {
				return data, true
			}
			if _, ok := fakeNodes[node.addr]; ok {
				return data, true
			}
			if len(fakeNodes) >= fN {
				return data, true
			}
			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypePrepare {
				return data, true
			}
			msg := ori.Copy()
			var sub core.Subject
			if err := rlp.DecodeBytes(msg.Msg, &sub); err != nil {
				log.Error("failed to decode subject", "err", err)
				return data, true
			}
			var qc QuorumCert
			if raw, err := rlp.EncodeToBytes(sub.QC); err != nil {
				log.Error("failed to encode prepareQC", "err", err)
				return data, true
			} else if err := rlp.DecodeBytes(raw, &qc); err != nil {
				log.Error("failed to decode prepareQC", "err", err)
				return data, true
			}
			qc.Node = common.HexToHash("0x124")
			var newSub = struct {
				Node *core.Node
				QC   *QuorumCert
			}{
				sub.Node,
				&qc,
			}
			if raw, err := rlp.EncodeToBytes(newSub); err != nil {
				log.Error("failed to encode new subject", "err", err)
				return data, true
			} else {
				msg.Msg = raw
			}

			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}
			fakeNodes[node.addr] = struct{}{}
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
