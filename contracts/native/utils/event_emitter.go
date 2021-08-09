package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
)

type EventEmitter struct {
	contract common.Address
	state    *state.StateDB
	block    uint64
}

func NewEventEmitter(contract common.Address, block uint64, state *state.StateDB) *EventEmitter {
	emitter := &EventEmitter{
		contract: contract,
		state:    state,
		block:    block,
	}
	return emitter
}

func (emitter *EventEmitter) Event(topics []common.Hash, data []byte) {
	event := &types.Log{
		Address:     emitter.contract,
		Topics:      topics,
		Data:        data,
		BlockNumber: emitter.block,
	}
	emitter.state.AddLog(event)
}
