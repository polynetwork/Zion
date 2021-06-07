package hotstuff

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

type PaceMarker struct {
	core *HotStuffService
	qc   *QuorumCert
	istanbulEventMux *event.TypeMux
	commitCh chan *types.Block
}

// GetLeader get leader with view number approach round robin
func (p *PaceMarker) GetLeader() common.Address {
	return p.core.vs.GetLeader(p.core.qcHigh.ViewNum)
}

// UpdateQCHigh
func (p *PaceMarker) UpdateQCHigh(qc *QuorumCert) {
	block := p.core.fetchBlock(qc.BlockHash, qc.ViewNum)
	if block == nil {
		return
	}
	oldQC := p.core.fetchBlock(p.core.qcHigh.BlockHash, p.core.qcHigh.ViewNum)
	if oldQC == nil {
		return
	}

	if qc.ViewNum > oldQC.NumberU64() {
		p.core.qcHigh = qc
		p.core.bLeaf = block
	}
}

// OnBeat if u = p.GetLeader then bleaf = onPropose(bleaf, cmd, qchigh)
// 调用onProposal, 在for循环内实现定时onBeat
func (p *PaceMarker) OnBeat() {

}

// send msg to new leader with view number increase
// 接收来自consensus.prepare的block，
// 根据viewnumber计算出新的primary，设置coinbase
// 组装MsgNewView
// 发送给新的leader
// 定时器处理
func (p *PaceMarker) OnNextSyncView() {

	//// view change
	//p.ehs.View.ViewNum++
	//p.ehs.View.Primary = p.ehs.GetLeader()
	//// create a dummyNode
	//dummyBlock := p.ehs.CreateLeaf(p.ehs.GetLeaf().Hash, nil, nil)
	//p.ehs.SetLeaf(dummyBlock)
	//dummyBlock.Committed = true
	//_ = p.ehs.BlockStorage.Put(dummyBlock)
	//// create a new view msg
	//newViewMsg := p.ehs.Msg(pb.MsgType_NEWVIEW, nil, p.ehs.GetHighQC())
	//// send msg
	//if p.ehs.ID != p.ehs.GetLeader() {
	//	_ = p.ehs.Unicast(p.ehs.GetNetworkInfo()[p.ehs.GetLeader()], newViewMsg)
	//}
	//// clean the current proposal
	//p.ehs.CurExec = consensus.NewCurProposal()
	//p.ehs.TimeChan.HardStartTimer()
}

// updateQCHigh with new qc
// 使用msg的qc更新updateHighQC
func (p *PaceMarker) OnReceiveNewView(msg *MsgNewView) {

}
