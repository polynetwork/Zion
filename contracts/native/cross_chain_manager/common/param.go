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
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/rlp"
	polycomm "github.com/polynetwork/poly/common"
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

func (this *InitRedeemScriptParam) Serialization(sink *polycomm.ZeroCopySink) {
	sink.WriteString(this.RedeemScript)
}

func (this *InitRedeemScriptParam) Deserialization(source *polycomm.ZeroCopySource) error {
	redeemScript, eof := source.NextString()
	if eof {
		return fmt.Errorf("MultiSignParam deserialize redeemScript error")
	}

	this.RedeemScript = redeemScript
	return nil
}

type EntranceParam struct {
	SourceChainID         uint64 `json:"sourceChainId"`
	Height                uint32 `json:"height"`
	Proof                 []byte `json:"proof"`
	RelayerAddress        []byte `json:"relayerAddress"` //in zion can be empty because caller can get through ctx
	Extra                 []byte `json:"extra"`
	HeaderOrCrossChainMsg []byte `json:"headerOrCrossChainMsg"`
}

func (this *EntranceParam) Serialization(sink *polycomm.ZeroCopySink) {
	sink.WriteUint64(this.SourceChainID)
	sink.WriteUint32(this.Height)
	sink.WriteVarBytes(this.Proof)
	sink.WriteVarBytes(this.RelayerAddress)
	sink.WriteVarBytes(this.Extra)
	sink.WriteVarBytes(this.HeaderOrCrossChainMsg)
}

func (this *EntranceParam) Deserialization(source *polycomm.ZeroCopySource) error {
	sourceChainID, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("EntranceParam deserialize sourcechainid error")
	}

	height, eof := source.NextUint32()
	if eof {
		return fmt.Errorf("EntranceParam deserialize height error")
	}
	proof, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("EntranceParam deserialize proof error")
	}
	relayerAddr, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("EntranceParam deserialize relayerAddr error")
	}
	extra, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("EntranceParam deserialize txdata error")
	}
	headerOrCrossChainMsg, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("EntranceParam deserialize headerOrCrossChainMsg error")
	}
	this.SourceChainID = sourceChainID
	this.Height = height
	this.Proof = proof
	this.RelayerAddress = relayerAddr
	this.Extra = extra
	this.HeaderOrCrossChainMsg = headerOrCrossChainMsg
	return nil
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

func (this *MultiSignParam) Serialization(sink *polycomm.ZeroCopySink) {
	sink.WriteUint64(this.ChainID)
	sink.WriteString(this.RedeemKey)
	sink.WriteVarBytes(this.TxHash)
	sink.WriteVarBytes([]byte(this.Address))
	sink.WriteUint64(uint64(len(this.Signs)))
	for _, v := range this.Signs {
		sink.WriteVarBytes(v)
	}
}

func (this *MultiSignParam) Deserialization(source *polycomm.ZeroCopySource) error {
	chainID, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("MultiSignParam deserialize txHash error")
	}
	redeemKey, eof := source.NextString()
	if eof {
		return fmt.Errorf("MultiSignParam deserialize redeemKey error")
	}
	txHash, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("MultiSignParam deserialize txHash error")
	}
	address, eof := source.NextString()
	if eof {
		return fmt.Errorf("MultiSignParam deserialize address error")
	}
	n, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("MultiSignParam deserialize signs length error")
	}
	signs := make([][]byte, 0)
	for i := 0; uint64(i) < n; i++ {
		v, eof := source.NextVarBytes()
		if eof {
			return fmt.Errorf("deserialize Signs error")
		}
		signs = append(signs, v)
	}

	this.ChainID = chainID
	this.RedeemKey = redeemKey
	this.TxHash = txHash
	this.Address = address
	this.Signs = signs
	return nil
}
