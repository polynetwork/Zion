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
package common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native"
)

const (
	REQUEST        = "request"
	DONE_TX        = "doneTx"
	MULTISIGN_INFO = "multisignInfo"
	RIPPLE_TX_INFO = "rippleTxInfo"

	NOTIFY_MAKE_PROOF_EVENT = "makeProof"
	REPLENISH_EVENT         = "ReplenishEvent"
)

type ChainHandler interface {
	MakeDepositProposal(service *native.NativeContract) (*MakeTxParam, error)
}

type MakeTxParam struct {
	TxHash              []byte
	CrossChainID        []byte
	FromContractAddress []byte
	ToChainID           uint64
	ToContractAddress   []byte
	Method              string
	Args                []byte
}

//used for param from evm contract
type MakeTxParamShim struct {
	TxHash              []byte
	CrossChainID        []byte
	FromContractAddress []byte
	ToChainID           *big.Int
	ToContractAddress   []byte
	Method              []byte
	Args                []byte
}

func DecodeTxParam(data []byte) (param *MakeTxParam, err error) {
	BytesTy, _ := abi.NewType("bytes", "", nil)
	IntTy, _ := abi.NewType("int", "", nil)
	// StringTy, _ := abi.NewType("string", "", nil)

	TxParam := abi.Arguments{
		{Type: BytesTy, Name: "txHash"},
		{Type: BytesTy, Name: "crossChainID"},
		{Type: BytesTy, Name: "fromContractAddress"},
		{Type: IntTy, Name: "toChainID"},
		{Type: BytesTy, Name: "toContractAddress"},
		{Type: BytesTy, Name: "method"},
		{Type: BytesTy, Name: "args"},
	}

	args, err := TxParam.Unpack(data)
	if err != nil {
		return
	}

	shim := new(MakeTxParamShim)
	err = TxParam.Copy(shim, args)
	if err != nil {
		return nil, err
	}
	param = &MakeTxParam{
		TxHash:              shim.TxHash,
		CrossChainID:        shim.CrossChainID,
		FromContractAddress: shim.FromContractAddress,
		ToChainID:           shim.ToChainID.Uint64(),
		ToContractAddress:   shim.ToContractAddress,
		Method:              string(shim.Method),
		Args:                shim.Args,
	}
	return
}

func EncodeTxParam(param *MakeTxParam) (data []byte, err error) {
	BytesTy, _ := abi.NewType("bytes", "", nil)
	IntTy, _ := abi.NewType("int", "", nil)

	TxParam := abi.Arguments{
		{Type: BytesTy, Name: "txHash"},
		{Type: BytesTy, Name: "crossChainID"},
		{Type: BytesTy, Name: "fromContractAddress"},
		{Type: IntTy, Name: "toChainID"},
		{Type: BytesTy, Name: "toContractAddress"},
		{Type: BytesTy, Name: "method"},
		{Type: BytesTy, Name: "args"},
	}

	shim := &MakeTxParamShim{
		TxHash:              param.TxHash,
		CrossChainID:        param.CrossChainID,
		FromContractAddress: param.FromContractAddress,
		ToChainID:           new(big.Int).SetUint64(param.ToChainID),
		ToContractAddress:   param.ToContractAddress,
		Method:              []byte(param.Method),
		Args:                param.Args,
	}

	data, err = TxParam.Pack(shim.TxHash, shim.CrossChainID, shim.FromContractAddress, shim.ToChainID, shim.ToContractAddress, shim.Method, shim.Args)
	return
}

type ToMerkleValue struct {
	TxHash      []byte
	FromChainID uint64
	MakeTxParam *MakeTxParam
}

type RippleTxArgs struct {
	ToAddress []byte
	Amount    *big.Int
}

func DecodeRippleTxArgs(data []byte) (param *RippleTxArgs, err error) {
	BytesTy, _ := abi.NewType("bytes", "", nil)
	IntTy, _ := abi.NewType("int", "", nil)
	// StringTy, _ := abi.NewType("string", "", nil)

	Args := abi.Arguments{
		{Type: BytesTy, Name: "toAddress"},
		{Type: IntTy, Name: "amount"},
	}

	args, err := Args.Unpack(data)
	if err != nil {
		return
	}

	param = new(RippleTxArgs)
	err = Args.Copy(param, args)
	if err != nil {
		return nil, err
	}
	return
}

func EncodeRippleTxArgs(args *RippleTxArgs) (data []byte, err error) {
	BytesTy, _ := abi.NewType("bytes", "", nil)
	IntTy, _ := abi.NewType("int", "", nil)

	Args := abi.Arguments{
		{Type: BytesTy, Name: "toAddress"},
		{Type: IntTy, Name: "amount"},
	}

	data, err = Args.Pack(args.ToAddress, args.Amount)
	return
}
