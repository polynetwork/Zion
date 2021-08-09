package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
)

func NewTestStateDB() *state.StateDB {
	memdb := rawdb.NewMemoryDatabase()
	db := state.NewDatabase(memdb)
	stateDB, _ := state.New(common.Hash{}, db, nil)
	return stateDB
}
