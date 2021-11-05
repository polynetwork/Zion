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

package eccm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/eccm"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
	polycomm "github.com/polynetwork/poly/common"
)

var (
	gasTable = map[string]uint64{
		MethodName:                     0,
		MethodCrossChain:               30000,
		MethodVerifyHeaderAndExecuteTx: 30000,
	}
)

func InitECCM() {
	InitABI()
	native.Contracts[this] = RegisterECCMContract
}

func RegisterECCMContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodName, Name)
	s.Register(MethodCrossChain, CrossChain)
	s.Register(MethodVerifyHeaderAndExecuteTx, VerifyHeaderAndExecuteTx)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return new(MethodContractNameOutput).Encode()
}

func CrossChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	sender := s.ContractRef().TxOrigin()
	caller := s.ContractRef().MsgSender()

	// check witness
	if caller != utils.MainChainLockProxyContractAddress {
		return utils.ByteFailed, fmt.Errorf("ZionMainChainECCM cross chain, caller MUST be lock proxy")
	}

	input := new(MethodCrossChainInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("ZionMainChainECCM cross chain, failed to decode params, err: %v", err)
	}

	// get and set tx index
	lastTxIndex, _ := getTxIndex(s)
	storeTxIndex(s, lastTxIndex+1)
	txIndex, txIndexProof := getTxIndex(s)

	// assemble tx, generate and store cross chain transaction proof
	txHash := s.ContractRef().TxHash()
	toChainID := input.ToChainID
	method := string(input.Method)
	toContract := input.ToContract
	args := input.TxData

	txParams := &scom.MakeTxParam{
		TxHash:              txHash[:],
		CrossChainID:        txIndexProof,
		FromContractAddress: caller[:],
		ToChainID:           toChainID,
		ToContractAddress:   toContract,
		Method:              method,
		Args:                args,
	}

	sink := polycomm.NewZeroCopySink(nil)
	txParams.Serialization(sink)
	txProof := crypto.Keccak256Hash(sink.Bytes())
	storeTxProof(s, txIndex, txProof)

	// emit event log
	if err := emitCrossChainEvent(s, sender, txIndexProof, caller, toChainID, toContract, args); err != nil {
		return utils.ByteFailed, fmt.Errorf("ZionMainChainECCM cross chain, failed to emit event log, err: %v", err)
	}
	return utils.ByteSuccess, nil
}

func VerifyHeaderAndExecuteTx(s *native.NativeContract) ([]byte, error) {
	return nil, nil
}
