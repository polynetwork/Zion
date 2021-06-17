package core

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

// sendNextChangeView sends the CHANGE VIEW message with current round + 1
func (c *core) sendNextChangeView() {
	cv := c.currentView()
	c.sendChangeView(new(big.Int).Add(cv.Round, common.Big1))
}

func (c *core) sendChangeView(round *big.Int) {
	logger := c.logger.New("state", c.state)
	curView := c.currentView()
	if curView.Round.Cmp(round) >= 0 {
		logger.Error("Cannot send out the round change", "current round", curView.Round, "target round", round)
		return
	}

	c.catchUpRound(&hotstuff.View{
		Round: new(big.Int).Set(round),
		// todo: height should be curViewHeight - 1
		Height: new(big.Int).Set(curView.Height),
	})

	// get new view
	curView = c.currentView()
	changeView := hotstuff.Subject{
		View:   curView,
		Digest: EmptyHash,
	}
	payload, err := Encode(changeView)
	if err != nil {
		logger.Error("Failed to encode ROUND CHANGE", "changeView", changeView, "err", err)
		return
	}
	c.broadcast(&message{
		Code: MsgTypeChangeView,
		Msg:  payload,
	})
}

// only the next leader will receive the CHANGE VIEW message
func (c *core) handleChangeView(msg *message, src hotstuff.Validator) error {
	logger := c.logger.New("state", c.state, "from", src.Address().Hex())

	// Decode ROUND CHANGE message
	var changeView *hotstuff.Subject
	if err := msg.Decode(&changeView); err != nil {
		logger.Error("Failed to decode ROUND CHANGE", "err", err)
		return errInvalidMessage
	}

	// This make sure the view should be identical
	if err := c.checkMessage(MsgTypeChangeView, changeView.View); err != nil {
		return err
	}

	if !c.IsProposer() {
		logger.Error("ChangeView message not sent to proposer", "expect", c.Address(), "got", src.Address())
		return errNotToProposer
	}

	cv := c.currentView()
	roundView := changeView.View

	// Add the ROUND CHANGE message to its message set and return how many
	// messages we've got with the same round number and sequence number.
	num, err := c.changeViewSet.Add(roundView.Round, msg)
	if err != nil {
		logger.Warn("Failed to add round change message", "from", src, "msg", msg, "err", err)
		return err
	}

	// We've received n-(n-1)/3 ROUND CHANGE messages, start a new round immediately.
	if num == c.valSet.Q() && (c.waitingForRoundChange || cv.Round.Cmp(roundView.Round) < 0) {
		c.startNewRound(roundView.Round)
		return nil
	} else if cv.Round.Cmp(roundView.Round) < 0 {
		// Only gossip the message with current round to other validators.
		return errIgnored
	}
	return nil
}

// ==================================================================

type changeViewSet struct {
	vs  hotstuff.ValidatorSet  // validator set
	ms  map[uint64]*messageSet // change view vote message collection
	mtx *sync.Mutex
}

func newChangeViewSet(valSet hotstuff.ValidatorSet) *changeViewSet {
	return &changeViewSet{
		vs:  valSet,
		ms:  make(map[uint64]*messageSet),
		mtx: new(sync.Mutex),
	}
}

func (s *changeViewSet) Add(r *big.Int, msg *message) (int, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	round := r.Uint64()
	if s.ms[round] == nil {
		s.ms[round] = newMessageSet(s.vs)
	}
	if err := s.ms[round].Add(msg); err != nil {
		return 0, err
	}
	return s.ms[round].Size(), nil
}

func (s *changeViewSet) Clear(r *big.Int) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	round := r.Uint64()
	for k, ms := range s.ms {
		if len(ms.Values()) == 0 || k < round {
			delete(s.ms, k)
		}
	}
}

// MaxRound returns the max round which the number of messages is equal or larger than num
func (s *changeViewSet) MaxRound(num int) *big.Int {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	var maxRound *big.Int
	for k, ms := range s.ms {
		if ms.Size() < num {
			continue
		}
		round := new(big.Int).SetUint64(k)
		if maxRound == nil || maxRound.Cmp(round) < 0 {
			maxRound = round
		}
	}
	return maxRound
}
