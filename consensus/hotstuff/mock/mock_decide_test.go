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

	"github.com/ethereum/go-ethereum/contracts/native/boot"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockDecideCase1
// net scale is 4, leader send fake message of decide with wrong height, repos change view.
func TestMockDecideCase1(t *testing.T) {
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
			if ori.Code != core.MsgTypeDecide {
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockDecideCase2
// net scale is 4, leader send fake message of decide with wrong round, repos change view.
func TestMockDecideCase2(t *testing.T) {
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
			if ori.Code != core.MsgTypeDecide {
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockDecideCase3
// net scale is 4, leader send fake message of decide with wrong block hash, repos change view.
func TestMockDecideCase3(t *testing.T) {
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
			if ori.Code != core.MsgTypeDecide {
				return data, true
			}
			msg := ori.Copy()

			var diploma core.Diploma
			if err := rlp.DecodeBytes(msg.Msg, &diploma); err != nil {
				log.Error("failed to decode diploma", "err", err)
				return data, true
			}
			diploma.BlockHash = common.HexToHash("0x123")
			if raw, err := rlp.EncodeToBytes(diploma); err != nil {
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
	sys.Close(20)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockDecideCase4
// net scale is 4, leader send fake message of decide with wrong qc.node, repos change view.
func TestMockDecideCase4(t *testing.T) {
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
			if ori.Code != core.MsgTypeDecide {
				return data, true
			}
			msg := ori.Copy()

			var diploma core.Diploma
			if err := rlp.DecodeBytes(msg.Msg, &diploma); err != nil {
				log.Error("failed to decode diploma", "err", err)
				return data, true
			}
			if raw, err := rlp.EncodeToBytes(diploma.CommitQC); err != nil {
				log.Error("failed to encode diploma.commitQC", "err", err)
				return data, true
			} else {
				var qc QuorumCert
				if err = rlp.DecodeBytes(raw, &qc); err != nil {
					log.Error("failed to decode diploma.commitQC", "err", err)
					return data, true
				} else {
					qc.Node = common.HexToHash("0x123")
				}
				var newDiploma = struct {
					CommitQC       *QuorumCert
					BlockHash      common.Hash
					CommittedSeals [][]byte
				}{
					&qc,
					diploma.BlockHash,
					diploma.CommittedSeals,
				}
				if raw, err := rlp.EncodeToBytes(newDiploma); err != nil {
					log.Error("failed to encode new diploma", "err", err)
					return data, true
				} else {
					msg.Msg = raw
				}
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
	sys.Close(20)
}
