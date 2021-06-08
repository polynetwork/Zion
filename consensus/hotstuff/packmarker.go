package hotstuff

import (
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

type PaceMarker struct {
	core *roundState
	chain consensus.ChainHeaderReader
	istanbulEventMux *event.TypeMux
	commitCh chan *types.Block
}

func newPaceMaker(c consensus.ChainHeaderReader) *PaceMarker {
	p := new(PaceMarker)
	p.chain = c
	p.istanbulEventMux = new(event.TypeMux)
	p.commitCh = make(chan *types.Block, 1)

	return p
}

// UpdateQCHigh
func (p *PaceMarker) UpdateQCHigh(qc *types.Header) {
	block := p.core.fetchBlock(qc.Hash(), qc.Number.Uint64())
	if block == nil {
		return
	}
	oldQC := p.core.fetchBlock(p.core.qcHigh.Hash(), p.core.qcHigh.Number.Uint64())
	if oldQC == nil {
		return
	}
	if qc.Number.Uint64() > oldQC.NumberU64() {
		p.core.qcHigh = qc
		p.core.bLeaf = block
	}
}

// OnBeat if u = p.GetLeader then bleaf = onPropose(bleaf, cmd, qchigh)
// 调用onProposal, 在for循环内实现定时onBeat
func (p *PaceMarker) OnRequest(sealedBlock *types.Block) {
	p.core.store.Push(sealedBlock)
}

func (p *PaceMarker) OnBeat() {
	block := p.core.store.Pop()
	if block == nil {
		return
	}
	isLeader := p.core.snap.ValSet.IsProposer(p.core.address)
	if !isLeader {
		return
	}
	p.core.sendPrepareMsg(block)
}

// send msg to new leader with view number increase
// 接收来自consensus.prepare的block，
// 根据viewnumber计算出新的primary，设置coinbase
// 组装MsgNewView
// 发送给新的leader
// 定时器处理
func (p *PaceMarker) OnNextSyncView() {
	if p.core.isLeader() {
		block := p.core.store.Pop()
		if block == nil {
			return
		}
		p.core.sendPrepareMsg(block)
	} else {
		highQC := p.core.getHighQC()
		leader := p.core.getLeader()
		p.core.sendNewViewMsg(highQC, leader)
	}
	// todo: error log info and set timer
}
