package hotstuff

import (
	"bufio"
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	// ErrUnauthorizedAddress is returned when given address cannot be found in
	// current validator set.
	ErrUnauthorizedAddress = errors.New("unauthorized address")
	// ErrStoppedEngine is returned if the engine is stopped
	ErrStoppedEngine = errors.New("stopped engine")
	// ErrStartedEngine is returned if the engine is already started
	ErrStartedEngine = errors.New("started engine")
)

type HotStuffService struct {
	privateKey *ecdsa.PrivateKey
	address  common.Address
	vs    Validators
	votes map[common.Hash]VoteSet

	bLeaf,
	bLock,
	bExec *types.Block

	vHeight uint64
	qcHigh  *QuorumCert

	pace *PaceMarker

	mtx          *sync.Mutex
	waitProposal *sync.Cond
	chain        consensus.ChainReader
	broadcaster  consensus.Broadcaster

	msgEvtCh   chan *MessageEvent
	msgEvtFeed event.Feed

	started bool
}

// todo
func NewHotstuffService() *HotStuffService {
	s := new(HotStuffService)

	s.mtx = new(sync.Mutex)
	s.waitProposal = sync.NewCond(s.mtx)

	s.msgEvtCh = make(chan *MessageEvent, 100)
	s.msgEvtFeed.Subscribe(s.msgEvtCh)

	return s
}

// 在创建consensus engine时使用
func (s *HotStuffService) SetBroadcaster(bc consensus.Broadcaster) {
	s.broadcaster = bc
}

// 在worker中调用，用来启动新的一轮共识
func (s *HotStuffService) NewChainHead() error {
	if !s.started {
		return ErrStoppedEngine
	}
	go s.msgEvtFeed.Send(MsgDecide{})
	return nil
}

func (s *HotStuffService) HandleMsg(address common.Address, msg p2p.Msg) (bool, error) {
	if msg.Code == P2PHotstuffMsg {
		if !s.started {
			return true, ErrStoppedEngine
		}

		// todo 消息转发过程
		//data, hash, err := s.decode(msg)
		//if err != nil {
		//	return true, errDecodeFailed
		//}
		//// Mark peer's message
		//ms, ok := sb.recentMessages.Get(addr)
		//var m *lru.ARCCache
		//if ok {
		//	m, _ = ms.(*lru.ARCCache)
		//} else {
		//	m, _ = lru.NewARC(inmemoryMessages)
		//	sb.recentMessages.Add(addr, m)
		//}
		//m.Add(hash, true)
		//
		//// Mark self known message
		//if _, ok := sb.knownMessages.Get(hash); ok {
		//	return true, nil
		//}
		//sb.knownMessages.Add(hash, true)
		//
		//go sb.istanbulEventMux.Post(istanbul.MessageEvent{
		//	Payload: data,
		//})
		//return true, nil
	}

	if msg.Code == P2PNewBlockMsg && s.vs.IsProposer(s.addr) { // eth.NewBlockMsg: import cycle
		// this case is to safeguard the race of similar block which gets propagated from other node while this node is proposing
		// as p2p.Msg can only be decoded once (get EOF for any subsequence read), we need to make sure the payload is restored after we decode it
		//log.Debug("Proposer received NewBlockMsg", "size", msg.Size, "payload.type", reflect.TypeOf(msg.Payload), "sender", addr)
		if reader, ok := msg.Payload.(*bytes.Reader); ok {
			payload, err := ioutil.ReadAll(reader)
			if err != nil {
				return true, err
			}
			reader.Reset(payload)       // ready to be decoded
			defer reader.Reset(payload) // restore so main eth/handler can decode
			var request struct {        // this has to be same as eth/protocol.go#newBlockData as we are reading NewBlockMsg
				Block *types.Block
				TD    *big.Int
			}
			if err := msg.Decode(&request); err != nil {
				//log.Debug("Proposer was unable to decode the NewBlockMsg", "error", err)
				return false, nil
			}
			// todo: free comment
			//newRequestedBlock := request.Block
			//if newRequestedBlock.Header().MixDigest == types.HotstuffDigest && sb.core.IsCurrentProposal(newRequestedBlock.Hash()) {
			//	//log.Debug("Proposer already proposed this block", "hash", newRequestedBlock.Hash(), "sender", addr)
			//	return true, nil
			//}
		}
	}
	return false, nil
}

func (s *HotStuffService) handleMessage() {
	for {
		select {
		case msg := <-s.msgEvtCh:
			println(msg)
		}
	}
}

// 应该在worker.prepare(创建区块头) -> finalizeAndAssemble(组装区块头以及), commitTransactions之后，进入到commit才seal
// 这里主要是将quorumCert的签名内容转换成extra
func (s *HotStuffService) CreateLeaf(parentHash common.Hash, block *types.Block, justify *QuorumCert) *types.Block {
	header := block.Header()
	if header != nil && justify != nil {
		//block.Header().Extra = justify.Signature.Extra()
	}
	return block
}

// block(b*), b'' <- b*.justify.node; b' <- b''.justify.node; b <- b'.justify.node;
// steps as follow:
// pre-commit step on b''-> s.paceMarker.UpdateQCHigh(b*) <--> paceMarker.QCHigh = b*
// commit step on b' -> if b'.height > b.height then s.block = b'
// decide step on b -> if (b''.parent == b') && (b'.parent == b) then onCommit(b) and s.exec = b
//
func (s *HotStuffService) Update(block *types.Block) {
	// block1 = b'', block2 = b', block3 = b
	b1 := s.fetchParent(block)
	if b1 == nil {
		return
	}

	// todo
	newqcFromBlock := &QuorumCert{
		BlockHash: block.Hash(),
		ViewNum:   block.NumberU64(),
		Type:      0,
		Signature: block.Extra(),
	}
	s.pace.UpdateQCHigh(newqcFromBlock)

	b2 := s.fetchParent(b1)
	if b2 == nil {
		return
	}
	if b2.NumberU64() > s.bLock.NumberU64() {
		s.bLock = b2
	}

	b3 := s.fetchParent(b2)
	if b3 == nil {
		return
	}
	s.OnCommit(b3)
	s.bExec = b3
}

// 落账，在worker中实现，这里应该要做的是将seal的结果放入到worker.resultCh
// todo: 问题在于如果这里递归的话，worker应该没法接收多个seal过的block
func (s *HotStuffService) OnCommit(block *types.Block) {
	if s.bExec.NumberU64() < block.NumberU64() {
		// 或者说如何追块
	}
}

// Msg(prepare, block, nil)
// if block.height > vheight &&
// (block extend s.block || block.justify.node.height > s.block.height) then
// vheight = block.height; send(getLeader, voteMsg)
// todo: timer 定时器重置
func (s *HotStuffService) OnReceiveProposal(msg *MsgPrepare) error {
	newBlock := msg.CurProposal
	if newBlock.NumberU64() <= s.vHeight {
		return fmt.Errorf("block height %d less than vHeight %d", msg.ViewNum, s.vHeight)
	}

	// s.extendBlock(s.bLock, newBlock)
	parent := s.fetchParent(newBlock)
	if (parent != nil && parent.NumberU64() > s.bLock.NumberU64()) || s.extendBlock(s.bLock, newBlock) {
		s.Update(newBlock)
		s.vHeight = newBlock.NumberU64()
		// todo: bls or multi sig,
		vote := &MsgPrepareVote{
			BlockHash:  common.Hash{},
			QC:         nil,
			PartialSig: nil,
			ViewNum:    0,
		}
		payload, _ := rlp.EncodeToBytes(vote)
		// todo: 确定是发送给当前newBlock的viewNumber对应的proposer？
		leader := s.vs.GetLeader(vote.ViewNum)
		s.unicast(leader, payload)
	}
	return nil
}

// collect votes and avoid dump votes
// if votes number > n - f or number >= 2f + 1 qc = QC; s.paceMarker.UpdateQCHigh
func (s *HotStuffService) OnReceiveVote(msg *MsgPrepareVote) {
	if _, ok := s.votes[msg.BlockHash]; !ok {
		s.votes[msg.BlockHash] = VoteSet{}
	}
	votes := s.votes[msg.BlockHash]
	votes.Add(msg)

	if votes.Marjor() {
		qc := &QuorumCert{
			BlockHash: msg.BlockHash,
			ViewNum:   msg.ViewNum,
			Type:      MsgTypePrepareVote,
			Signature: []byte{},
		}
		s.pace.UpdateQCHigh(qc)
	}
}

// create leaf and broadcast
func (s *HotStuffService) OnPropose() error {
	leaf := s.CreateLeaf(common.Hash{}, nil, nil)
	msg := new(MsgNewView)
	msg.ViewNum = leaf.NumberU64()

	var buf bytes.Buffer
	if err := msg.EncodeRLP(bufio.NewWriter(&buf)); err != nil {
		return err
	}
	return s.broadcast(buf.Bytes())
}

func (s *HotStuffService) extendBlock(ancestor, block *types.Block) bool {
	b := block
	for b = s.fetchParent(b); b != nil; {
		if b.Hash() == ancestor.Hash() {
			return true
		}
	}
	return false
}

func (s *HotStuffService) fetchParent(block *types.Block) *types.Block {
	return s.fetchBlock(block.ParentHash(), block.NumberU64()-1)
}

// todo: 从chainReader中拿到block，阻塞式等待
func (s *HotStuffService) fetchBlock(hash common.Hash, view uint64) *types.Block {
	return s.chain.GetBlock(hash, view)
	//if block != nil {
	//	return block, nil
	//}
	//
	//s.waitProposal.Wait()
	//block = s.chain.GetBlock(hash, view)
	//if block == nil {
	//	return nil, fmt.Errorf("block (%s %d) not arrived", hash.Hex(), view)
	//}
	//return block, nil
}

func (s *HotStuffService) broadcast(payload []byte) error {
	// send to others
	s.gossip(payload)

	// send to self
	msg := &MessageEvent{
		Payload: payload,
	}
	go s.msgEvtFeed.Send(msg)
	return nil
}

// Broadcast implements istanbul.Backend.Gossip
// todo:
//  1. record message
//  3. peer msg lru
func (s *HotStuffService) gossip(payload []byte) error {
	targets := make(map[common.Address]bool)
	for _, v := range s.vs.Get() {
		if v != s.addr {
			targets[v] = true
		}
	}
	if s.broadcaster != nil && len(targets) > 0 {
		ps := s.broadcaster.FindPeers(targets)
		for _, p := range ps {
			go p.Send(P2PHotstuffMsg, payload)
		}
	}
	return nil
}

func (s *HotStuffService) unicast(to common.Address, payload []byte) error {
	if to == s.addr {
		return nil
	}
	peer := s.broadcaster.FindPeer(to)
	if peer == nil {
		return fmt.Errorf("can't find p2p peer of %s", to.Hex())
	}
	go peer.Send(P2PHotstuffMsg, payload)
	return nil
}
