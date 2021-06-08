package hotstuff

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
)

type PaceMarker struct {
	core                *roundState
	chain               consensus.ChainHeaderReader
	requestCh, commitCh chan *types.Block
	timer               *time.Timer
}

func newPaceMaker(c consensus.ChainHeaderReader) *PaceMarker {
	p := new(PaceMarker)
	p.chain = c
	p.commitCh = make(chan *types.Block, 1)
	p.requestCh = make(chan *types.Block, 1)
	p.timer = time.NewTimer(time.Duration(blockPeriod) * time.Second)
	return p
}

func (p *PaceMarker) Start(ctx context.Context) {
	for {
		select {
		case <-p.timer.C:
			p.OnNextSyncView()
		case sealedBlock := <-p.requestCh:
			p.core.store.Push(sealedBlock)
		case <- ctx.Done():
			break
		}
	}
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
