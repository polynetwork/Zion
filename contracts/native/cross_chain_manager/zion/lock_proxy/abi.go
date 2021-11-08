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

package lock_proxy

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/lock_proxy"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "zion main chain cross chain manager"

func InitABI() {
	ab, err := abi.JSON(strings.NewReader(ILockProxyABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	ABI = &ab
}

var (
	ABI  *abi.ABI
	this = utils.MainChainECCMContractAddress
)

type MethodContractNameInput struct{}

func (m *MethodContractNameInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodName)
}
func (m *MethodContractNameInput) Decode(payload []byte) error { return nil }

type MethodContractNameOutput struct {
	Name string
}

func (m *MethodContractNameOutput) Encode() ([]byte, error) {
	m.Name = contractName
	return utils.PackOutputs(ABI, MethodName, m.Name)
}
func (m *MethodContractNameOutput) Decode(payload []byte) error {
	return utils.UnpackOutputs(ABI, MethodName, m, payload)
}

type MethodCrossChainInput struct {
	ToChainID  uint64
	ToContract []byte
	Method     []byte
	TxData     []byte
}

func (i *MethodCrossChainInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodCrossChain, i.ToChainID, i.ToContract, i.Method, i.TxData)
}
func (i *MethodCrossChainInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodCrossChain, &i, payload)
}

type MethodVerifyHeaderAndExecuteTxInput struct {
	Header []byte
	Proof  []byte
}

func (i *MethodVerifyHeaderAndExecuteTxInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodVerifyHeaderAndExecuteTx, i.Header, i.Proof)
}
func (i *MethodVerifyHeaderAndExecuteTxInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodVerifyHeaderAndExecuteTx, i, payload)
}

func emitCrossChainEvent(s *native.NativeContract,
	sender common.Address,
	txID []byte,
	proxyOrAssetContract common.Address,
	toChainID uint64,
	toContract []byte,
	method string,
	rawData []byte,
) error {
	return s.AddNotify(ABI, []string{EventCrossChainEvent, sender.Hex()}, txID, proxyOrAssetContract, toChainID, toContract, method, rawData)
}

func emitVerifyHeaderAndExecuteTxEvent(s *native.NativeContract,
	fromChainID uint64,
	toContract []byte,
	crossChainTxHash []byte,
	fromChainTxHash []byte,
) error {
	return s.AddNotify(ABI, []string{EventVerifyHeaderAndExecuteTxEvent}, fromChainID, toContract, crossChainTxHash, fromChainTxHash)
}
