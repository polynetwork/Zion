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
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	KEY_PREFIX_BTC = "btc"

	KEY_PREFIX_BTC_VOTE = "btcVote"
	REQUEST             = "request"
	DONE_TX             = "doneTx"

	NOTIFY_MAKE_PROOF_EVENT = "makeProof"
)

type ChainHandler interface {
	MakeDepositProposal(service *native.NativeContract) (*MakeTxParam, error)
}

type InitRedeemScriptParam struct {
	RedeemScript string
}

type EntranceParam struct {
	SourceChainID         uint64 `json:"sourceChainId"`
	Height                uint32 `json:"height"`
	Proof                 []byte `json:"proof"`
	RelayerAddress        []byte `json:"relayerAddress"` //in zion can be empty because caller can get through ctx
	Extra                 []byte `json:"extra"`
	HeaderOrCrossChainMsg []byte `json:"headerOrCrossChainMsg"`
}

func (this *EntranceParam) String() string {
	str := "{"
	str += fmt.Sprintf("source chain id: %d,", this.SourceChainID)
	str += fmt.Sprintf("height: %d,", this.Height)
	if this.Proof != nil && len(this.Proof) > 0 {
		str += fmt.Sprintf("proof: %s,", hexutil.Encode(this.Proof))
	}
	if this.RelayerAddress != nil && len(this.RelayerAddress) > 0 {
		str += fmt.Sprintf("relayer address: %s,", hexutil.Encode(this.RelayerAddress))
	}
	if this.Extra != nil && len(this.Extra) > 0 {
		str += fmt.Sprintf("extra: %s,", hexutil.Encode(this.Extra))
	}
	if this.HeaderOrCrossChainMsg != nil && len(this.HeaderOrCrossChainMsg) > 0 {
		str += fmt.Sprintf("header or cross chain msg: %s", hexutil.Encode(this.HeaderOrCrossChainMsg))
	}
	str += "}"
	return str
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

func (m *MakeTxParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.TxHash, m.CrossChainID, m.FromContractAddress, m.ToChainID,
		m.ToContractAddress, m.Method, m.Args})
}
func (m *MakeTxParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		TxHash              []byte
		CrossChainID        []byte
		FromContractAddress []byte
		ToChainID           uint64
		ToContractAddress   []byte
		Method              string
		Args                []byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.TxHash, m.CrossChainID, m.FromContractAddress, m.ToChainID, m.ToContractAddress, m.Method, m.Args =
		data.TxHash, data.CrossChainID, data.FromContractAddress, data.ToChainID, data.ToContractAddress, data.Method, data.Args
	return nil
}

type MakeTxParamWithSender struct {
	Sender common.Address
	MakeTxParam
}

//used for param from evm contract
type MakeTxParamWithSenderShim struct {
	Sender      common.Address
	MakeTxParam []byte
}

func (this *MakeTxParamWithSender) Deserialization(data []byte) (err error) {

	BytesTy, _ := abi.NewType("bytes", "", nil)
	AddrTy, _ := abi.NewType("address", "", nil)
	// StringTy, _ := abi.NewType("string", "", nil)

	TxParam := abi.Arguments{
		{Type: AddrTy, Name: "sender"},
		{Type: BytesTy, Name: "makeTxParam"},
	}

	args, err := TxParam.Unpack(data)
	if err != nil {
		return
	}

	shim := new(MakeTxParamWithSenderShim)
	err = TxParam.Copy(shim, args)
	if err != nil {
		return
	}

	this.Sender = shim.Sender
	makeTxParam, err := DecodeTxParam(shim.MakeTxParam)
	if err != nil {
		return
	}
	this.MakeTxParam = *makeTxParam
	return
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

func (m *ToMerkleValue) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.TxHash, m.FromChainID, m.MakeTxParam})
}

func (m *ToMerkleValue) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		TxHash      []byte
		FromChainID uint64
		MakeTxParam *MakeTxParam
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.TxHash, m.FromChainID, m.MakeTxParam = data.TxHash, data.FromChainID, data.MakeTxParam
	return nil
}

type MultiSignParam struct {
	ChainID   uint64
	RedeemKey string
	TxHash    []byte
	Address   string
	Signs     [][]byte
}

type TxArgs struct {
	ToAssetHash []byte
	ToAddress   []byte
	Amount      *big.Int
}

func (tx *TxArgs) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{tx.ToAssetHash, tx.ToAddress, tx.Amount})
}

func (tx *TxArgs) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ToAssetHash []byte
		ToAddress   []byte
		Amount      *big.Int
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	tx.ToAssetHash, tx.ToAddress, tx.Amount = data.ToAssetHash, data.ToAddress, data.Amount
	return nil
}
