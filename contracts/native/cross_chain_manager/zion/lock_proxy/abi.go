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
	"math/big"
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

// function name
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

// function bindProxy
type MethodBindProxyInput struct {
	ToChainId       uint64
	TargetProxyHash []byte
}

func (i *MethodBindProxyInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodBindProxyHash, i.ToChainId, i.TargetProxyHash)
}
func (i *MethodBindProxyInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodBindProxyHash, i, payload)
}

// function bindAsset
type MethodBindAssetHashInput struct {
	FromAssetHash common.Address
	ToChainId     uint64
	ToAssetHash   []byte
}

func (i *MethodBindAssetHashInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodBindAssetHash, i.FromAssetHash, i.ToChainId, i.ToAssetHash)
}
func (i *MethodBindAssetHashInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodBindAssetHash, i, payload)
}

// function lock
type MethodLockInput struct {
	FromAssetHash common.Address
	ToChainId     uint64
	ToAddress     []byte
	Amount        *big.Int
}

func (i *MethodLockInput) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, MethodLock, i.FromAssetHash, i.ToChainId, i.ToAddress, i.Amount)
}
func (i *MethodLockInput) Decode(payload []byte) error {
	return utils.UnpackMethod(ABI, MethodLock, i, payload)
}

// function verifyHeaderAndExecuteTx
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

//type MethodCrossChainInput struct {
//	ToChainID  uint64
//	ToContract []byte
//	Method     []byte
//	TxData     []byte
//}
//
//func (i *MethodCrossChainInput) Encode() ([]byte, error) {
//	return utils.PackMethod(ABI, MethodCrossChain, i.ToChainID, i.ToContract, i.Method, i.TxData)
//}
//func (i *MethodCrossChainInput) Decode(payload []byte) error {
//	return utils.UnpackMethod(ABI, MethodCrossChain, &i, payload)
//}

//event BindProxyEvent(uint64 toChainId, bytes targetProxyHash);
func emitBindProxyEvent(s *native.NativeContract, toChainID uint64, targetProxyHash []byte) error {
	return s.AddNotify(ABI, []string{EventBindProxyEvent}, toChainID, targetProxyHash)
}

//event BindAssetEvent(address fromAssetHash, uint64 toChainId, bytes targetProxyHash, uint initialAmount);
func emitBindAssetEvent(s *native.NativeContract,
	fromAsset common.Address,
	toChainID uint64,
	targetProxyHash []byte,
	initialAmount *big.Int) error {
	return s.AddNotify(ABI, []string{EventBindAssetEvent}, fromAsset, toChainID, targetProxyHash, initialAmount)
}

//event LockEvent(address fromAssetHash, address fromAddress, uint64 toChainId, bytes toAssetHash, bytes toAddress, uint256 amount);
func emitLockEvent(s *native.NativeContract,
	fromAssetHash, fromAddress common.Address,
	toChainID uint64,
	toAssetHash, toAddress []byte,
	amount *big.Int) error {
	return s.AddNotify(ABI, []string{EventLockEvent}, fromAssetHash, fromAddress, toChainID, toAssetHash, toAddress, amount)
}

//event UnlockEvent(address toAssetHash, address toAddress, uint256 amount);
func emitUnlockEvent(s *native.NativeContract, toAssetHash, toAddress common.Address, amount *big.Int) error {
	return s.AddNotify(ABI, []string{EventUnlockEvent}, toAssetHash, toAddress, amount)
}

//event CrossChainEvent(address indexed sender, bytes txId, address proxyOrAssetContract, uint64 toChainId, bytes toContract, string method, bytes rawdata);
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

//event VerifyHeaderAndExecuteTxEvent(uint64 fromChainID, bytes toContract, bytes crossChainTxHash, bytes fromChainTxHash);
func emitVerifyHeaderAndExecuteTxEvent(s *native.NativeContract,
	fromChainID uint64,
	toContract []byte,
	crossChainTxHash []byte,
	fromChainTxHash []byte,
) error {
	return s.AddNotify(ABI, []string{EventVerifyHeaderAndExecuteTxEvent}, fromChainID, toContract, crossChainTxHash, fromChainTxHash)
}