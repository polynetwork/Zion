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
	xutils "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/utils"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/lock_proxy"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/zion"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	gasTable = map[string]uint64{
		MethodName:                     0,
		MethodBindProxyHash:            30000,
		MethodBindAssetHash:            30000,
		MethodLock:                     30000,
		MethodVerifyHeaderAndExecuteTx: 30000,
	}
)

func InitECCM() {
	InitABI()
	native.Contracts[this] = RegisterLockProxyContract
}

func RegisterLockProxyContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodName, Name)
	s.Register(MethodBindProxyHash, BindProxy)
	s.Register(MethodBindAssetHash, BindAsset)
	s.Register(MethodLock, Lock)
	s.Register(MethodVerifyHeaderAndExecuteTx, VerifyHeaderAndExecuteTx)
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
	input := new(MethodLockInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("ZionMainChain Lock proxy, failed to decode params, err: %v", err)
	}
	if input.Amount.Cmp(common.Big0) <= 0 {
		return utils.ByteFailed, fmt.Errorf("ZionMainChain lock proxy, amount should be greater than zero")
	}
	return nil, nil
}

func Unlock(s *native.NativeContract) ([]byte, error) {
	return nil, nil
}

func CrossChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := s.ContractRef().MsgSender()

	// check authority
	if caller != utils.MainChainLockProxyContractAddress {
		return utils.ByteFailed, fmt.Errorf("ZionMainChainECCM crossChain, caller MUST be `lock proxy`")
	}

	input := new(MethodCrossChainInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("ZionMainChainECCM crossChain, failed to decode params, err: %v", err)
	}

	// get and set tx index
	lastTxIndex, _ := getTxIndex(s)
	storeTxIndex(s, lastTxIndex+1)
	txIndex, txIndexID := getTxIndex(s)

	// assemble tx, generate and store cross chain transaction proof
	sender := s.ContractRef().TxOrigin()
	txHash := s.ContractRef().TxHash()
	method := string(input.Method)
	args := input.TxData
	blob, proof := encodeMakeTxParams(txHash, txIndex, caller, input.ToChainID, input.ToContract, method, args)

	storeTxProof(s, txIndex, proof)
	storeTxParams(s, proof, blob)

	// emit event log
	if err := emitCrossChainEvent(s, sender, txIndexID, caller, input.ToChainID, input.ToContract, method, args); err != nil {
		return utils.ByteFailed, fmt.Errorf("ZionMainChainECCM crossChain, failed to emit event log, err: %v", err)
	}
	return utils.ByteSuccess, nil
}

func VerifyHeaderAndExecuteTx(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	input := new(MethodVerifyHeaderAndExecuteTxInput)

	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("ZionMainChainECCM verifyHeaderAndExecuteTx, failed to decode params, err: %v", err)
	}

	header := new(types.Header)
	if err := header.UnmarshalJSON(input.Header); err != nil {
		return utils.ByteFailed, fmt.Errorf("ZionMainChainECCM verifyHeaderAndExecuteTx, failed to umarshal header, err: %v", err)
	}
	height := header.Number.Uint64()

	epoch, err := nm.GetEpochByHeight(s.StateDB(), height)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("ZionMainChainECCM verifyHeaderAndExecuteTx, failed to get epoch by height, err: %v", err)
	}

	if _, _, err := zion.VerifyHeader(header, epoch.MemberList(), false); err != nil {
		return utils.ByteFailed, fmt.Errorf("ZionMainChainECCM verifyHeaderAndExecuteTx, failed to verify header, err: %v", err)
	}

	proofResult, err := xutils.VerifyTx(input.Proof, header, this, nil, false)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("ZionMainChainECCM verifyHeaderAndExecuteTx, failed to verify tx, err: %v", err)
	}

	// todo
	txParams, err := decodeMakeTxParams(proofResult)
	if err != nil {
		return nil, err
	}
	return txParams.Args, nil
	return nil, nil
}

func executeCrossChainTx(s *native.NativeContract) ([]byte, error) {
	return nil, nil
}
