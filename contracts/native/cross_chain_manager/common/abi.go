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
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/cross_chain_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	MethodContractName        = cross_chain_manager_abi.MethodName
	MethodImportOuterTransfer = cross_chain_manager_abi.MethodImportOuterTransfer
	MethodMultiSignRipple     = cross_chain_manager_abi.MethodMultiSignRipple
	MethodReconstructRippleTx = cross_chain_manager_abi.MethodReconstructRippleTx
	MethodCheckDone           = cross_chain_manager_abi.MethodCheckDone
	MethodBlackChain          = cross_chain_manager_abi.MethodBlackChain
	MethodWhiteChain          = cross_chain_manager_abi.MethodWhiteChain
	MethodReplenish           = cross_chain_manager_abi.MethodReplenish
)

var ABI *abi.ABI

func init() {
	ABI = GetABI()
}

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(cross_chain_manager_abi.ICrossChainManagerABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

type InitRedeemScriptParam struct {
	RedeemScript string
}

type CheckDoneParam struct {
	ChainID      uint64
	CrossChainID []byte
}

func (m *CheckDoneParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodCheckDone, m)
}

type CheckDoneOutput struct {
	Done bool
}

func (m *CheckDoneOutput) Decode(payload []byte) error {
	if err := utils.UnpackOutputs(ABI, MethodCheckDone, m, payload); err != nil {
		return err
	}
	return nil
}

type EntranceParam struct {
	SourceChainID uint64
	Height        uint32
	Proof         []byte
	Extra         []byte
	Signature     []byte
}

func (m *EntranceParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodImportOuterTransfer, m)
}

func (m *EntranceParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.SourceChainID, m.Height, m.Proof, m.Extra, m.Signature})
}
func (m *EntranceParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		SourceChainID uint64
		Height        uint32
		Proof         []byte
		Extra         []byte
		Signature     []byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.SourceChainID, m.Height, m.Proof, m.Extra, m.Signature = data.SourceChainID, data.Height,
		data.Proof, data.Extra, data.Signature
	return nil
}

//Digest Digest calculate the hash of param input
func (m *EntranceParam) Digest() ([]byte, error) {
	input := &EntranceParam{
		SourceChainID: m.SourceChainID,
		Height:        m.Height,
		Proof:         m.Proof,
		Extra:         m.Extra,
	}
	msg, err := rlp.EncodeToBytes(input)
	if err != nil {
		return nil, fmt.Errorf("EntranceParam, serialize input error: %v", err)
	}
	digest := crypto.Keccak256(msg)
	return digest, nil
}

func (this *EntranceParam) String() string {
	str := "{"
	str += fmt.Sprintf("source chain id: %d,", this.SourceChainID)
	str += fmt.Sprintf("height: %d,", this.Height)
	if this.Proof != nil && len(this.Proof) > 0 {
		str += fmt.Sprintf("proof: %s,", hexutil.Encode(this.Proof))
	}
	if this.Extra != nil && len(this.Extra) > 0 {
		str += fmt.Sprintf("extra: %s,", hexutil.Encode(this.Extra))
	}
	str += "}"
	return str
}

type MultiSignParam struct {
	ToChainId    uint64
	AssetAddress []byte
	FromChainId  uint64
	TxHash       []byte
	TxJson       string
}

func (m *MultiSignParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodMultiSignRipple, m)
}

func (m *MultiSignParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ToChainId, m.AssetAddress, m.FromChainId, m.TxHash, m.TxJson})
}

func (m *MultiSignParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ToChainId    uint64
		AssetAddress []byte
		FromChainId  uint64
		TxHash       []byte
		TxJson       string
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.ToChainId, m.AssetAddress, m.FromChainId, m.TxHash, m.TxJson = data.ToChainId, data.AssetAddress,
		data.FromChainId, data.TxHash, data.TxJson
	return nil
}

type ReconstructTxParam struct {
	FromChainId uint64
	TxHash      []byte
	ToChainId   uint64
}

func (m *ReconstructTxParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodReconstructRippleTx, m)
}

func (m *ReconstructTxParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.FromChainId, m.TxHash, m.ToChainId})
}

func (m *ReconstructTxParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		FromChainId uint64
		TxHash      []byte
		ToChainId   uint64
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.FromChainId, m.TxHash, m.ToChainId = data.FromChainId, data.TxHash, data.ToChainId
	return nil
}

type BlackChainParam struct {
	ChainID uint64
}

type ReplenishParam struct {
	ChainID  uint64
	TxHashes []string
}

func (m *ReplenishParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, MethodReplenish, m)
}
