package core

import "github.com/ethereum/go-ethereum/consensus"

// StateMachineRepo contains chainReader and block pool.
type StateMachineRepo struct {
	chain consensus.ChainReader
	pool *BlockPool
}


