package hotstuff

import (
	"bytes"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"math/big"
	"sync"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/p2p"
)

type currentProposal struct {
	highQC []*types.Header
	votes  map[common.Address]*MsgVote
}

type roundState struct {
	privateKey *ecdsa.PrivateKey
	address  common.Address

	snap *Snapshot
	curRnd *currentProposal
	store *Storage
	chain consensus.ChainReader

	bLeaf,
	bLock,
	bExec *types.Block

	vHeight uint64
	qcHigh  *types.Header

	pace *PaceMarker

	mtx          *sync.Mutex
	waitProposal *sync.Cond
	broadcaster  consensus.Broadcaster

	msgEvtCh   chan *MessageEvent
	msgEvtFeed event.Feed

	started bool
}

// todo
func newRoundState(privateKey *ecdsa.PrivateKey, chain consensus.ChainReader) *roundState {
	s := new(roundState)
	s.privateKey = privateKey
	s.address = crypto.PubkeyToAddress(s.privateKey.PublicKey)

	s.curRnd = &currentProposal{
		highQC: []*types.Header{},
	}

	s.mtx = new(sync.Mutex)
	s.waitProposal = sync.NewCond(s.mtx)

	s.msgEvtCh = make(chan *MessageEvent, 100)
	s.msgEvtFeed.Subscribe(s.msgEvtCh)

	s.store = newStorage(chain)

	// todo: genesis header or persist into snapshot
	s.vHeight = s.store.CurrentHeader().Number.Uint64()
	s.qcHigh = s.store.CurrentHeader()
	return s
}

// 在创建consensus engine时使用
func (s *roundState) SetBroadcaster(bc consensus.Broadcaster) {
	s.broadcaster = bc
}

// 在worker中调用，用来启动新的一轮共识
func (s *roundState) NewChainHead() error {
	if !s.started {
		return errStoppedEngine
	}
	go s.msgEvtFeed.Send(MsgDecide{})
	return nil
}

func (s *roundState) getHighQC() *types.Header {
	if s.qcHigh == nil {
		header := s.store.GetHeaderByNumber(0)
		return header
	}
	return s.qcHigh
}

func (s *roundState) HandleMsg(address common.Address, msg p2p.Msg) (bool, error) {
	if msg.Code == P2PHotstuffMsg {
		if !s.started {
			return true, errStoppedEngine
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

func (s *roundState) handleMessage() {
	for {
		select {
		case msg := <-s.msgEvtCh:
			println(msg)
		}
	}
}

// block(b*), b'' <- b*.justify.node; b' <- b''.justify.node; b <- b'.justify.node;
// steps as follow:
// pre-commit step on b''-> s.paceMarker.UpdateQCHigh(b*) <--> paceMarker.QCHigh = b*
// commit step on b' -> if b'.height > b.height then s.block = b'
// decide step on b -> if (b''.parent == b') && (b'.parent == b) then onCommit(b) and s.exec = b
//
func (s *roundState) Update(qc *types.Header) {
	// block1 = b'', block2 = b', block3 = b
	b1 := s.fetchParentHeader(qc)
	if b1 == nil {
		return
	}

	s.pace.UpdateQCHigh(qc)

	b2 := s.fetchParentBlockWithHeader(b1)
	if b2 == nil {
		return
	}
	if b2.NumberU64() > s.bLock.NumberU64() {
		s.bLock = b2
	}

	b3 := s.fetchParentBlock(b2)
	if b3 == nil {
		return
	}

	s.pace.commitCh <- b3
	s.bExec = b3
}
