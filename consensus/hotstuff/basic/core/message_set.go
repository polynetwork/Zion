package core

import (
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
)

// Construct a new message set to accumulate messages for given height/view number.
func newMessageSet(valSet hotstuff.ValidatorSet) *messageSet {
	return &messageSet{
		view: &hotstuff.View{
			Round:  new(big.Int),
			Height: new(big.Int),
		},
		mtx:  new(sync.Mutex),
		msgs: make(map[common.Address]*message),
		vs:   valSet,
	}
}

type messageSet struct {
	view *hotstuff.View
	vs   hotstuff.ValidatorSet
	mtx  *sync.Mutex
	msgs map[common.Address]*message
}

func (s *messageSet) View() *hotstuff.View {
	return s.view
}

func (s *messageSet) Add(msg *message) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if err := s.verify(msg); err != nil {
		return err
	}

	s.msgs[msg.Address] = msg
	return nil
}

func (s *messageSet) Values() (result []*message) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, v := range s.msgs {
		result = append(result, v)
	}
	return
}

func (s *messageSet) Size() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return len(s.msgs)
}

func (s *messageSet) Get(addr common.Address) *message {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.msgs[addr]
}

func (s *messageSet) String() string {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	addresses := make([]string, 0, len(s.msgs))
	for _, v := range s.msgs {
		addresses = append(addresses, v.Address.Hex())
	}
	return fmt.Sprintf("[%v]", strings.Join(addresses, ", "))
}

// verify if the message comes from one of the validators
func (s *messageSet) verify(msg *message) error {
	if _, v := s.vs.GetByAddress(msg.Address); v == nil {
		return ErrUnauthorizedAddress
	}

	return nil
}
