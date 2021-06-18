package core

import "github.com/ethereum/go-ethereum/consensus/hotstuff"

type backlogEvent struct {
	src hotstuff.Validator
	msg *message
}

type timeoutEvent struct{}
