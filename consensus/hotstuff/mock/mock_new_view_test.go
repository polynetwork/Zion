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
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestSimple
func TestSimple(t *testing.T) {
	node_manager.InitABI()
	sys := makeSystem(7)
	sys.Start()
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockNewViewCase1
// net scale is 7, 2 of them send fake message of newView with wrong height.
func TestMockNewViewCase1(t *testing.T) {
	node_manager.InitABI()
	H, R, fR, fN := uint64(4), uint64(0), uint64(1), int(1)
	fakeNodes := make(map[common.Address]struct{})
	mu := new(sync.Mutex)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
				return data, true
			}

			mu.Lock()
			if _, ok := fakeNodes[node.addr]; ok {
				mu.Unlock()
				return data, true
			}
			mu.Unlock()

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeNewView {
				return data, true
			}
			msg := ori.Copy()
			msg.View.Round = new(big.Int).SetUint64(fR)
			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}

			mu.Lock()
			fakeNodes[node.addr] = struct{}{}
			if len(fakeNodes) > fN {
				mu.Unlock()
				return data, true
			}
			mu.Unlock()

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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockNewViewCase2
// net scale is 4, one of them send fake message of newView with wrong node. err should be "failed to verify prepareQC"
func TestMockNewViewCase2(t *testing.T) {
	node_manager.InitABI()
	H, R, fN := uint64(4), uint64(0), 1
	fakeNodes := make(map[common.Address]struct{})
	mu := new(sync.Mutex)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
				return data, true
			}

			mu.Lock()
			if _, ok := fakeNodes[node.addr]; ok {
				mu.Unlock()
				return data, true
			}
			mu.Unlock()

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeNewView {
				return data, true
			}
			msg := ori.Copy()
			var qc QuorumCert
			if err := rlp.DecodeBytes(msg.Msg, &qc); err != nil {
				log.Error("failed to decode prepareQC", "err", err)
				return data, true
			}
			qc.Node = common.HexToHash("0x123")
			raw, err := rlp.EncodeToBytes(qc)
			if err != nil {
				log.Error("encode prepareQC failed", "err", err)
				return data, true
			}
			msg.Msg = raw
			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}

			mu.Lock()
			fakeNodes[node.addr] = struct{}{}
			if len(fakeNodes) > fN {
				mu.Unlock()
				return data, true
			}
			mu.Unlock()

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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockNewViewCase3
// net scale is 4, one of them send message of newView to wrong leader
func TestMockNewViewCase3(t *testing.T) {
	node_manager.InitABI()
	H, R, fN := uint64(4), uint64(0), 1
	fakeNodes := make(map[common.Address]struct{})
	mu := new(sync.Mutex)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
				return data, true
			}

			mu.Lock()
			if _, ok := fakeNodes[node.addr]; ok {
				mu.Unlock()
				return data, true
			}
			mu.Unlock()

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeNewView {
				return data, true
			}

			// send to other repo
			for _, peer := range node.broadcaster.peers {
				if !peer.geth.IsProposer() && peer.geth.addr != node.addr {
					peer.Send(hotstuffMsg, data)
				}
			}

			mu.Lock()
			fakeNodes[node.addr] = struct{}{}
			if len(fakeNodes) > fN {
				mu.Unlock()
				return data, true
			}
			mu.Unlock()

			view := &core.View{
				Round:  new(big.Int).SetUint64(r),
				Height: new(big.Int).SetUint64(h),
			}
			log.Info("-----fake message", "address", node.addr, "msg", ori.Code, "view", view)
			return data, false
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockNewViewCase4
// net scale is 4, one of them send fake message of newView with wrong height. err should be "failed to verify prepareQC"
func TestMockNewViewCase4(t *testing.T) {
	node_manager.InitABI()
	H, R, fH, fN := uint64(4), uint64(0), uint64(5), 1
	fakeNodes := make(map[common.Address]struct{})
	mu := new(sync.Mutex)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
				return data, true
			}

			mu.Lock()
			if _, ok := fakeNodes[node.addr]; ok {
				mu.Unlock()
				return data, true
			}
			mu.Unlock()

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeNewView {
				return data, true
			}
			msg := ori.Copy()
			var qc QuorumCert
			if err := rlp.DecodeBytes(msg.Msg, &qc); err != nil {
				log.Error("failed to decode prepareQC", "err", err)
				return data, true
			}
			qc.View.Height = new(big.Int).SetUint64(fH)
			raw, err := rlp.EncodeToBytes(qc)
			if err != nil {
				log.Error("encode prepareQC failed", "err", err)
				return data, true
			}
			msg.Msg = raw
			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}

			mu.Lock()
			fakeNodes[node.addr] = struct{}{}
			if len(fakeNodes) > fN {
				mu.Unlock()
				return data, true
			}
			mu.Unlock()

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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockNewViewCase5
// net scale is 4, one of them send fake message of newView with wrong round. err should be "failed to verify prepareQC"
func TestMockNewViewCase5(t *testing.T) {
	node_manager.InitABI()
	H, R, fR, fN := uint64(4), uint64(0), uint64(1), 1
	fakeNodes := make(map[common.Address]struct{})
	mu := new(sync.Mutex)
	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
				return data, true
			}

			mu.Lock()
			if _, ok := fakeNodes[node.addr]; ok {
				mu.Unlock()
				return data, true
			}
			mu.Unlock()

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeNewView {
				return data, true
			}
			msg := ori.Copy()
			var qc QuorumCert
			if err := rlp.DecodeBytes(msg.Msg, &qc); err != nil {
				log.Error("failed to decode prepareQC", "err", err)
				return data, true
			}
			qc.View.Round = new(big.Int).SetUint64(fR)
			raw, err := rlp.EncodeToBytes(qc)
			if err != nil {
				log.Error("encode prepareQC failed", "err", err)
				return data, true
			}
			msg.Msg = raw
			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}

			mu.Lock()
			fakeNodes[node.addr] = struct{}{}
			if len(fakeNodes) > fN {
				mu.Unlock()
				return data, true
			}
			mu.Unlock()

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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockNewViewCase6
// net scale is 4, one of them send fake message of newView without enough signatures. err should be "failed to verify prepareQC"
func TestMockNewViewCase6(t *testing.T) {
	node_manager.InitABI()
	H, R, fN := uint64(4), uint64(0), 1
	fakeNodes := make(map[common.Address]struct{})
	mu := new(sync.Mutex)

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
				return data, true
			}

			mu.Lock()
			if _, ok := fakeNodes[node.addr]; ok {
				mu.Unlock()
				return data, true
			}
			mu.Unlock()

			var ori core.Message
			if err := rlp.DecodeBytes(data, &ori); err != nil {
				log.Error("failed to decode message", "err", err)
				return data, true
			}
			if ori.Code != core.MsgTypeNewView {
				return data, true
			}
			msg := ori.Copy()
			var qc QuorumCert
			if err := rlp.DecodeBytes(msg.Msg, &qc); err != nil {
				log.Error("failed to decode prepareQC", "err", err)
				return data, true
			}
			qc.CommittedSeal = qc.CommittedSeal[:len(qc.CommittedSeal)-1]
			raw, err := rlp.EncodeToBytes(qc)
			if err != nil {
				log.Error("encode prepareQC failed", "err", err)
				return data, true
			}
			msg.Msg = raw
			payload, err := node.resignMsg(msg)
			if err != nil {
				log.Error("failed to resign message")
				return data, true
			}

			mu.Lock()
			fakeNodes[node.addr] = struct{}{}
			if len(fakeNodes) > fN {
				mu.Unlock()
				return data, true
			}
			mu.Unlock()

			view := &core.View{
				Round:  new(big.Int).SetUint64(r),
				Height: new(big.Int).SetUint64(h),
			}
			log.Info("-----fake message", "address", node.addr, "msg", msg.Code, "view", view, "msg", msg, "qc.length", len(qc.CommittedSeal))
			return payload, true
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}
