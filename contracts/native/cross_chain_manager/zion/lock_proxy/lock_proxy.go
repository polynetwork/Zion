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
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	cutils "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/utils"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/lock_proxy"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

var (
	gasTable = map[string]uint64{
		MethodName:          0,
		MethodBindProxyHash: 10000,
		MethodBindAssetHash: 10000,
		MethodLock:          10000,
	}
)

func InitLockProxy() {
	InitABI()
	native.Contracts[this] = RegisterLockProxyContract
}

func RegisterLockProxyContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodName, Name)
	s.Register(MethodBindProxyHash, BindProxy)
	s.Register(MethodBindAssetHash, BindAsset)
	s.Register(MethodLock, Lock)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return new(MethodContractNameOutput).Encode()
}

func BindProxy(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	input := new(MethodBindProxyInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindProxy, failed to decode params, err: %v", err)
	}
	if input.ToChainId == native.ZionMainChainID {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindProxy, bind self is illegal")
	}
	if input.TargetProxyHash == nil || len(input.TargetProxyHash) == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindProxy, target proxy address is invalid")
	}

	sender := s.ContractRef().TxOrigin()
	ok, err := nm.CheckConsensusSigns(s, MethodBindProxyHash, ctx.Payload, sender)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindProxy, failed to checkConsensusSigns, err: %v", err)
	}
	if !ok {
		return utils.ByteFailed, nil
	}

	gotTargetProxyHash, _ := getProxy(s, input.ToChainId)
	if gotTargetProxyHash != nil && bytes.Equal(gotTargetProxyHash, input.TargetProxyHash) {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindProxy, duplicate bindship, asset %s, chainID %d",
			hexutil.Encode(input.TargetProxyHash), input.ToChainId)
	}

	storeProxy(s, input.ToChainId, input.TargetProxyHash)
	if err := emitBindProxyEvent(s, input.ToChainId, input.TargetProxyHash); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindProxy, failed to emit event log, err: %v", err)
	}
	return utils.ByteSuccess, nil
}

func BindAsset(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	input := new(MethodBindAssetHashInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset, failed to decode params, err: %v", err)
	}

	if input.ToChainId == native.ZionMainChainID {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset, bind self is illegal")
	}
	if input.ToAssetHash == nil || len(input.ToAssetHash) == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset, target asset is invalid")
	}
	if onlySupportNativeToken(input.FromAssetHash) {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset, only support native token")
	}
	// allow the same address of `fromAsset` and `targetAsset`, different asset may have the same address in different chain

	sender := s.ContractRef().TxOrigin()
	ok, err := nm.CheckConsensusSigns(s, MethodBindAssetHash, ctx.Payload, sender)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset, failed to checkConsensusSigns, err: %v", err)
	}
	if !ok {
		return utils.ByteFailed, nil
	}

	gotTargetAssetHash, err := getAsset(s, input.FromAssetHash, input.ToChainId)
	if gotTargetAssetHash != nil && bytes.Equal(gotTargetAssetHash, input.ToAssetHash) {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset, duplicate bindship, from asset %s, target asset %s, target chain id %d",
			input.FromAssetHash.Hex(), hexutil.Encode(input.ToAssetHash), input.ToChainId)
	}

	currentBalance, err := getBalanceFor(s, input.FromAssetHash)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset")
	}

	storeAsset(s, input.FromAssetHash, input.ToChainId, input.ToAssetHash)
	if err := emitBindAssetEvent(s, input.FromAssetHash, input.ToChainId, input.ToAssetHash, currentBalance); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset, failed to emit event log, err: %v", err)
	}
	return utils.ByteSuccess, nil
}

func Lock(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	sender := s.ContractRef().TxOrigin()

	input := new(MethodLockInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to decode params, err: %v", err)
	}
	if input.Amount.Cmp(common.Big0) == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, amount should be greater than zero")
	}

	// transfer asset
	if err := transfer2Contract(s, input.FromAssetHash, sender, input.Amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, transfer to contract failed, err: %v", err)
	}

	// assemble tx data
	toAsset := common.EmptyAddress.Bytes()
	blob, _ := getAsset(s, input.FromAssetHash, input.ToChainId)
	if blob != nil {
		toAsset = blob
	}
	txData := EncodeTxArgs(toAsset, input.ToAddress, input.Amount)

	// get and set tx index
	lastTxIndex, _ := getTxIndex(s)
	storeTxIndex(s, lastTxIndex+1)
	txIndex, txIndexID := getTxIndex(s)

	// assemble tx, generate and store cross chain transaction proof
	txHash := s.ContractRef().TxHash()
	method := "unlock"
	txParams, txParamsEnc, proof := EncodeMakeTxParams(txHash, txIndex, sender, input.ToChainId, toAsset, method, txData)

	storeTxProof(s, txIndex, proof)
	storeTxParams(s, proof, txParamsEnc)

	// emit event log
	if err := emitCrossChainEvent(s, sender, txIndexID, sender, input.ToChainId, toAsset, method, txData); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to emit crossChainEvent log, err: %v", err)
	}
	if err := emitLockEvent(s, input.FromAssetHash, sender, input.ToChainId, toAsset, input.ToAddress, input.Amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to emit lockEvent log, err: %v", err)
	}

	if err := cutils.MakeTransaction(s, txParams, native.ZionMainChainID); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to makeTransaction, err: %v", err)
	}
	return utils.ByteSuccess, nil
}

func Unlock(s *native.NativeContract, txParams *scom.MakeTxParam) error {
	args, err := DecodeTxArgs(txParams.Args)
	if err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to decode txArgs, err: %v", err)
	}

	asset := common.BytesToAddress(args.ToAssetHash)
	toAddress := common.BytesToAddress(args.ToAddress)
	if err := transferFromContract(s, asset, toAddress, args.Amount); err != nil {
		return err
	}
	return emitUnlockEvent(s, asset, toAddress, args.Amount)
}
