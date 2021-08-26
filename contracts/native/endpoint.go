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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

// support native functions to evm functions.
type EVMHandler func(caller, addr common.Address, input []byte) ([]byte, uint64, error)

type ContractRef struct {
	contexts []*Context

	stateDB     *state.StateDB
	blockHeight *big.Int
	txHash      common.Hash
	msgSender   common.Address
	evmHandler  EVMHandler
	gasLeft     uint64
}

func NewContractRef(
	db *state.StateDB,
	sender common.Address,
	blockHeight *big.Int,
	txHash common.Hash,
	suppliedGas uint64,
	evmHandler EVMHandler) *ContractRef {

	return &ContractRef{
		contexts:    make([]*Context, 0),
		stateDB:     db,
		msgSender:   sender,
		blockHeight: blockHeight,
		txHash:      txHash,
		gasLeft:     suppliedGas,
		evmHandler:  evmHandler,
	}
}

func (s *ContractRef) NativeCall(
	caller,
	contractAddr common.Address,
	payload []byte,
) (ret []byte, gasLeft uint64, err error) {

	s.PushContext(&Context{
		Caller:          caller,
		ContractAddress: contractAddr,
		Payload:         payload,
	})
	defer s.PopContext()

	contract := NewNativeContract(s.stateDB, s)
	ret, err = contract.Invoke()
	gasLeft = s.gasLeft
	return
}

func (s *ContractRef) EVMCall(caller, contractAddr common.Address, input []byte) ([]byte, uint64, error) {
	if s.evmHandler == nil {
		return nil, 0, nil
	}
	return s.evmHandler(caller, contractAddr, input)
}

func (s *ContractRef) StateDB() *state.StateDB {
	return s.stateDB
}

func (s *ContractRef) BlockHeight() *big.Int {
	return s.blockHeight
}

func (s *ContractRef) TxHash() common.Hash {
	return s.txHash
}

func (s *ContractRef) MsgSender() common.Address {
	return s.msgSender
}

func (s *ContractRef) GasLeft() uint64 {
	return s.gasLeft
}

const (
	MAX_EXECUTE_CONTEXT = 128
)

type Context struct {
	Caller          common.Address
	ContractAddress common.Address
	Payload         []byte
}

// PushContext push current context to smart contract
func (s *ContractRef) PushContext(context *Context) {
	s.contexts = append(s.contexts, context)
}

// CurrentContext return smart contract current context
func (s *ContractRef) CurrentContext() *Context {
	if len(s.contexts) < 1 {
		return nil
	}
	return s.contexts[len(s.contexts)-1]
}

// PopContext pop smart contract current context
func (s *ContractRef) PopContext() {
	if len(s.contexts) > 1 {
		s.contexts = s.contexts[:len(s.contexts)-1]
	}
}

// CallingContext return smart contract caller context
func (s *ContractRef) CallingContext() *Context {
	if len(s.contexts) < 2 {
		return nil
	}
	return s.contexts[len(s.contexts)-2]
}

// EntryContext return smart contract entry entrance context
func (s *ContractRef) EntryContext() *Context {
	if len(s.contexts) < 1 {
		return nil
	}
	return s.contexts[0]
}

func (s *ContractRef) CheckContexts() bool {
	if len(s.contexts) == 0 {
		return false
	}
	if len(s.contexts) > MAX_EXECUTE_CONTEXT {
		return false
	}
	return true
}

func (s *ContractRef) CheckWitness(address common.Address) bool {
	if s.CheckAccountAddress(address) || s.CheckContractAddress(address) {
		return true
	}
	return false
}

func (s *ContractRef) CheckAccountAddress(address common.Address) bool {
	if s.msgSender == address {
		return true
	}
	return false
}

func (s *ContractRef) CheckContractAddress(address common.Address) bool {
	if s.CallingContext() != nil && s.CallingContext().ContractAddress == address {
		return true
	}
	return false
}
