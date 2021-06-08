package hotstuff

import (
	"context"
	"crypto/ecdsa"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
)

type currentProposal struct {
	highQC []*types.Header
	votes  map[common.Address]*MsgVote
}

type roundState struct {
	privateKey *ecdsa.PrivateKey
	address    common.Address

	snap   *Snapshot
	curRnd *currentProposal
	store  *Storage
	chain  consensus.ChainReader
	msgCh chan *InnerMsg

	bLeaf,
	bLock,
	bExec *types.Block

	vHeight uint64
	qcHigh  *types.Header

	pace *PaceMarker

	mtx          *sync.Mutex
	waitProposal *sync.Cond
	broadcaster  consensus.Broadcaster

	started bool
}

// todo
func newRoundState(ctx context.Context, privateKey *ecdsa.PrivateKey, chain consensus.ChainReader) *roundState {
	s := new(roundState)
	s.privateKey = privateKey
	s.address = crypto.PubkeyToAddress(s.privateKey.PublicKey)

	s.curRnd = &currentProposal{
		highQC: []*types.Header{},
	}

	s.mtx = new(sync.Mutex)
	s.waitProposal = sync.NewCond(s.mtx)

	s.store = newStorage(chain)

	// todo: genesis header or persist into snapshot
	s.vHeight = s.store.CurrentHeader().Number.Uint64()
	s.qcHigh = s.store.CurrentHeader()
	s.msgCh = make(chan *InnerMsg, 1)

	s.started = true

	return s
}

// 在创建consensus engine时使用
func (s *roundState) SetBroadcaster(bc consensus.Broadcaster) {
	s.broadcaster = bc
}

// 在worker中调用，用来启动新的一轮共识, todo: 无需commit
func (s *roundState) NewChainHead() error {
	//if !s.started {
	//	return errStoppedEngine
	//}
	//go s.msgEvtFeed.Send(MsgDecide{})
	//return nil
	return nil
}

func (s *roundState) getHighQC() *types.Header {
	if s.qcHigh == nil {
		header := s.store.GetHeaderByNumber(0)
		return header
	}
	return s.qcHigh
}

func (s *roundState) HandleMsg(address common.Address, m p2p.Msg) (bool, error) {
	if m.Code != P2PHotstuffMsg {
		return false, errUnknownMsgType
	}
	if !s.started {
		return true, errStoppedEngine
	}

	payload, err := s.decode(m)
	if err != nil {
		return true, errDecodeFailed
	}
	if err := s.handlePayload(payload); err != nil {
		return true, err
	}
	return true, nil
}

func (s *roundState) handleSelfMsg(ctx context.Context) {
	for {
		select {
		case data := <-s.msgCh:
			s.handlePayload(data.Payload)
		case <- ctx.Done():
			break
		}
	}
}

func (s *roundState) handlePayload(payload []byte) error {
	// Decode message and check its signature
	msg := new(Message)
	if err := msg.FromPayload(payload, s.checkValidatorSignature); err != nil {
		log.Error("Failed to decode message from payload", "err", err)
		return err
	}

	// Only accept message if the address is valid
	_, src := s.snap.ValSet.GetByAddress(msg.Address)
	if src == nil {
		log.Error("Invalid address in message", "msg", msg)
		return errUnauthorizedAddress
	}

	if err := s.handleCheckedMsg(msg); err != nil {
		log.Debug("handle msg failed, ", "type", msg.Code.String(), "error", err)
		return err
	}

	return nil
}

func (s *roundState) handleCheckedMsg(msg *Message) error {
	var err error
	switch msg.Code {
	case MsgTypeNewView:
		err = s.handleNewViewMsg(msg)
	case MsgTypeProposal:
		err = s.handleProposalMsg(msg)
	case MsgTypeVote:
		err = s.handleVoteMsg(msg)
	}
	return err
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
