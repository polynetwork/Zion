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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestSimple
func TestSimple(t *testing.T) {
	sys := makeSystem(7)
	sys.Start()
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCase1
// net scale is 7, 2 of them send fake message of newView with wrong height.
func TestMockCase1(t *testing.T) {
	H, R, fR, fN := uint64(4), uint64(0), uint64(1), int(1)
	fakeNodes := make(map[common.Address]struct{})

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCase2
// net scale is 4, one of them send fake message of newView with wrong node. err should be "failed to verify prepareQC"
func TestMockCase2(t *testing.T) {
	H, R, fN := uint64(4), uint64(0), 1
	fakeNodes := make(map[common.Address]struct{})

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCase3
// net scale is 4, one of them send message of newView to wrong leader
func TestMockCase3(t *testing.T) {
	H, R, fN := uint64(4), uint64(0), 1
	fakeNodes := make(map[common.Address]struct{})

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
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
			if ori.Code != core.MsgTypeNewView {
				return data, true
			}

			// send to other repo
			for _, peer := range node.broadcaster.peers {
				if !peer.geth.IsProposer() && peer.geth.addr != node.addr {
					peer.Send(hotstuffMsg, data)
				}
			}

			fakeNodes[node.addr] = struct{}{}
			view := &core.View{
				Round:  new(big.Int).SetUint64(r),
				Height: new(big.Int).SetUint64(h),
			}
			log.Info("fake message", "address", node.addr, "msg", ori.Code, "view", view)
			return data, false
		}
		return data, true
	}

	for _, node := range sys.nodes {
		node.setHook(hook)
	}
	sys.Close(10)
}

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCase4
// net scale is 4, one of them send fake message of newView with wrong height. err should be "failed to verify prepareQC"
func TestMockCase4(t *testing.T) {
	H, R, fH, fN := uint64(4), uint64(0), uint64(5), 1
	fakeNodes := make(map[common.Address]struct{})

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
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

// go test -v -count=1 github.com/ethereum/go-ethereum/consensus/hotstuff/mock -run TestMockCase5
// net scale is 4, one of them send fake message of newView with wrong round. err should be "failed to verify prepareQC"
func TestMockCase5(t *testing.T) {
	H, R, fR, fN := uint64(4), uint64(0), uint64(1), 1
	fakeNodes := make(map[common.Address]struct{})

	sys := makeSystem(4)
	sys.Start()
	time.Sleep(2 * time.Second)

	hook := func(node *Geth, data []byte) ([]byte, bool) {
		if h, r := node.api.CurrentSequence(); h == H && r == R {
			if node.IsProposer() {
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
