package core

import "github.com/ethereum/go-ethereum/common"

type CoreEngine interface {
	Start() error
	Stop() error

	IsProposer() bool

	// verify if a hash is the same as the proposed block in the current pending request
	//
	// this is useful when the engine is currently the speaker
	//
	// pending request is populated right at the request stage so this would give us the earliest verification
	// to avoid any race condition of coming propagated blocks
	IsCurrentProposal(blockHash common.Hash) bool

	// CurrentRoundState() *roundState
}
