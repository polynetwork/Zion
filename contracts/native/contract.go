/*
 * Copyright (C) 2021 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */
package native

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
)

type (
	RegisterService func(native *NativeContract)
	MethodHandler   func(contract *NativeContract) ([]byte, error)
)

var (
	Contracts = make(map[common.Address]RegisterService)
)

type NativeContract struct {
	ref      *ContractRef
	db       *state.StateDB
	handlers map[string]MethodHandler // map method id to method handler
	gasTable map[string]uint64        // map method id to gas usage
	ab       *abi.ABI
}

func NewNativeContract(db *state.StateDB, ref *ContractRef) *NativeContract {
	return &NativeContract{
		db:       db,
		ref:      ref,
		handlers: make(map[string]MethodHandler),
	}
}

func (s *NativeContract) ContractRef() *ContractRef {
	return s.ref
}

func (s *NativeContract) GetCacheDB() *state.CacheDB {
	return (*state.CacheDB)(s.db)
}

func (s *NativeContract) StateDB() *state.StateDB {
	return s.db
}

func (s *NativeContract) Prepare(ab *abi.ABI, gasTb map[string]uint64) {
	s.ab = ab
	s.gasTable = make(map[string]uint64)
	for name, gas := range gasTb {
		id := utils.MethodID(s.ab, name)
		s.gasTable[id] = gas
	}
}

func (s *NativeContract) Register(name string, handler MethodHandler) {
	methodID := utils.MethodID(s.ab, name)
	s.handlers[methodID] = handler
}

// Invoke return execute ret and cost gas
func (s *NativeContract) Invoke() ([]byte, error) {
	// check context
	if !s.ref.CheckContexts() {
		return nil, fmt.Errorf("context error")
	}
	ctx := s.ref.CurrentContext()

	// find methodID
	if len(ctx.Payload) < 4 {
		return nil, fmt.Errorf("invalid input")
	}
	methodID := hexutil.Encode(ctx.Payload[:4])

	// register methods
	registerHandler, ok := Contracts[ctx.ContractAddress]
	if !ok {
		return nil, fmt.Errorf("failed to find contract: [%x]", ctx.ContractAddress)
	}
	registerHandler(s)

	// get method handler
	handler, ok := s.handlers[methodID]
	if !ok {
		return nil, fmt.Errorf("failed to find method: [%s]", methodID)
	}

	// check gasLeft
	needGas, ok := s.gasTable[methodID]
	if !ok {
		return nil, fmt.Errorf("failed to find method: [%s]", methodID)
	}
	gasLeft := s.ref.gasLeft
	if gasLeft < needGas && gasLeft < MinGasUsage {
		return nil, fmt.Errorf("gasLeft not enough, need %d", needGas)
	}

	// execute transaction and cost gas
	ret, err := handler(s)
	if err != nil && needGas > MinGasUsage {
		needGas = MinGasUsage
	}
	if needGas > 0 {
		s.ref.gasLeft -= needGas
	}

	return ret, err
}

func (s *NativeContract) AddNotify(abi *abi.ABI, topics []string, data ...interface{}) (err error) {

	var topicIDs []common.Hash
	for _, topic := range topics {
		eventInfo, ok := abi.Events[topic]
		if !ok {
			err = fmt.Errorf("topic %s not exists", topic)
			return
		}
		topicIDs = append(topicIDs, eventInfo.ID)
	}

	packedData, err := utils.PackEvents(abi, topics[0], data...)
	if err != nil {
		err = fmt.Errorf("AddNotify, PackEvents error: %v", err)
		return
	}
	emitter := utils.NewEventEmitter(s.ref.CurrentContext().ContractAddress, s.ContractRef().BlockHeight().Uint64(), s.StateDB())
	emitter.Event(topicIDs, packedData)
	return
}
