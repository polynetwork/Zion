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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	zutils "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/utils"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/lock_proxy"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/params"
)

const MIN_BALANCE = 1000000 // default minimum value of lock proxy balance

var (
	gasTable = map[string]uint64{
		MethodName:          0,
		MethodBindProxyHash: 10000,
		MethodGetProxyHash:  0,
		MethodBindAssetHash: 10000,
		MethodGetAssetHash:  0,
		MethodBindCaller:    10000,
		MethodGetCaller:     0,
		MethodLock:          10000,
	}

	minBalance = new(big.Int).Mul(big.NewInt(MIN_BALANCE), params.OneEth)
)

func InitLockProxy() {
	InitABI()
	native.Contracts[this] = RegisterLockProxyContract
}

func RegisterLockProxyContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodName, Name)
	s.Register(MethodBindProxyHash, BindProxy)
	s.Register(MethodGetProxyHash, GetProxy)
	s.Register(MethodBindAssetHash, BindAsset)
	s.Register(MethodGetAssetHash, GetAsset)
	s.Register(MethodBindCaller, BindCaller)
	s.Register(MethodGetCaller, GetCaller)
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
	if input.ToChainId == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindProxy, target chain id invalid")
	}
	if input.ToChainId == native.ZionMainChainID {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindProxy, bind self is illegal")
	}
	if input.TargetProxyHash == nil || len(input.TargetProxyHash) == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindProxy, target proxy address is invalid")
	}

	// filter duplicate proxy
	gotTargetProxyHash, _ := getProxy(s, input.ToChainId)
	if gotTargetProxyHash != nil && bytes.Equal(gotTargetProxyHash, input.TargetProxyHash) {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindProxy, duplicate proxy, proxy %s, chainID %d",
			hexutil.Encode(input.TargetProxyHash), input.ToChainId)
	}

	// vote and check quorum size
	sender := s.ContractRef().TxOrigin()
	ok, err := nm.CheckConsensusSigns(s, MethodBindProxyHash, ctx.Payload, sender)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindProxy, failed to checkConsensusSigns, err: %v", err)
	}
	if !ok {
		return utils.ByteFailed, nil
	}

	// bind success
	storeProxy(s, input.ToChainId, input.TargetProxyHash)
	if err := emitBindProxyEvent(s, input.ToChainId, input.TargetProxyHash); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindProxy, failed to emit `BindProxyEvent`, err: %v", err)
	}
	return utils.ByteSuccess, nil
}

func GetProxy(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	input := new(MethodGetProxyInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.GetProxy, failed to decode params, err: %v", err)
	}
	if input.ToChainId == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.GetProxy, target chain id invalid")
	}
	if input.ToChainId == native.ZionMainChainID {
		return utils.ByteFailed, fmt.Errorf("LockProxy.GetProxy, target chain id wont be 1")
	}

	proxy, err := getProxy(s, input.ToChainId)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.GetProxy, get bound proxy failed, err: %v", err)
	}
	return proxy, nil
}

func BindAsset(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	input := new(MethodBindAssetHashInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset, failed to decode params, err: %v", err)
	}

	if input.ToChainId == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset, target chain id invalid")
	}
	if input.ToChainId == native.ZionMainChainID {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset, bind self is illegal")
	}
	if input.ToAssetHash == nil || len(input.ToAssetHash) == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset, target asset invalid")
	}
	if !onlySupportNativeToken(input.FromAssetHash) {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset, only support native token")
	}

	// filter duplicate asset
	gotTargetAssetHash, err := getAsset(s, input.FromAssetHash, input.ToChainId)
	if gotTargetAssetHash != nil && bytes.Equal(gotTargetAssetHash, input.ToAssetHash) {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset, duplicate asset, from asset %s, target asset %s, target chain id %d",
			input.FromAssetHash.Hex(), hexutil.Encode(input.ToAssetHash), input.ToChainId)
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

	currentBalance, err := getBalanceFor(s, input.FromAssetHash, this)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset")
	}
	if currentBalance == nil || currentBalance.Cmp(minBalance) < 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset, contract insufficient balance")
	}

	storeAsset(s, input.FromAssetHash, input.ToChainId, input.ToAssetHash)
	if err := emitBindAssetEvent(s, input.FromAssetHash, input.ToChainId, input.ToAssetHash, currentBalance); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindAsset, failed to emit `BindAssetEvent`, err: %v", err)
	}
	return utils.ByteSuccess, nil
}

func GetAsset(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	input := new(MethodGetAssetInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.GetAsset, failed to decode params, err: %v", err)
	}
	if input.ToChainId == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.GetAsset, target chain id invalid")
	}
	if input.ToChainId == native.ZionMainChainID {
		return utils.ByteFailed, fmt.Errorf("LockProxy.GetAsset, target chain id wont be 1")
	}
	if !onlySupportNativeToken(input.FromAssetHash) {
		return utils.ByteFailed, fmt.Errorf("LockProxy.GetAsset, only support native token")
	}

	asset, err := getAsset(s, input.FromAssetHash, input.ToChainId)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.GetAsset, failed to get bound asset, err: %v", err)
	}

	return asset, nil
}

func BindCaller(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	input := new(MethodBindCallerInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindCaller, failed to decode params, err: %v", err)
	}
	if input.ToChainId == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindCaller, target chain id invalid")
	}
	if input.ToChainId == native.ZionMainChainID {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindCaller, bind self is illegal")
	}
	if input.Caller == nil || len(input.Caller) == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindCaller, target caller invalid")
	}

	gotTargetCaller, err := getCaller(s, input.ToChainId)
	if gotTargetCaller != nil && bytes.Equal(gotTargetCaller, input.Caller) {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindCaller, duplicate caller,target caller %s, target chain id %d",
			hexutil.Encode(input.Caller), input.ToChainId)
	}

	sender := s.ContractRef().TxOrigin()
	ok, err := nm.CheckConsensusSigns(s, MethodBindCaller, ctx.Payload, sender)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindCaller, failed to checkConsensusSigns, err: %v", err)
	}
	if !ok {
		return utils.ByteFailed, nil
	}

	storeCaller(s, input.ToChainId, input.Caller)
	if err := emitBindCallerEvent(s, input.ToChainId, input.Caller); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.BindCaller, failed to emit `BindCallerEvent`, err: %v", err)
	}
	return utils.ByteSuccess, nil
}

func GetCaller(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	input := new(MethodGetCallerInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.GetCaller, failed to decode params, err: %v", err)
	}
	if input.ToChainId == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.GetCaller, target chain id invalid")
	}
	if input.ToChainId == native.ZionMainChainID {
		return utils.ByteFailed, fmt.Errorf("LockProxy.GetCaller, target chain id wont be 1")
	}

	caller, err := getCaller(s, input.ToChainId)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.GetCaller, failed to get bound caller, err: %v", err)
	}

	return caller, nil
}

func Lock(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	input := new(MethodLockInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to decode params, err: %v", err)
	}
	if !onlySupportNativeToken(input.FromAssetHash) {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, only support native token")
	}
	if input.ToChainId == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, target chain id invalid")
	}
	if input.ToChainId == native.ZionMainChainID {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, target chain id wont be 1")
	}
	if input.ToAddress == nil || len(input.ToAddress) == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, target address invalid")
	}
	if input.Amount == nil || input.Amount.Cmp(common.Big0) == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, amount invalid")
	}

	// check target caller
	targetCaller, err := getCaller(s, input.ToChainId)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to get bound caller, err: %v", err)
	}
	if targetCaller == nil || len(targetCaller) == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, bound caller not exist")
	}

	// check target proxy
	targetProxy, err := getProxy(s, input.ToChainId)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to get bound proxy, err: %v", err)
	}
	if targetProxy == nil || len(targetProxy) == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, bound proxy not exist")
	}

	// all asset MUST be bind first, include native token
	toAsset, err := getAsset(s, input.FromAssetHash, input.ToChainId)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to get bound asset, err: %v", err)
	}
	if toAsset == nil || len(toAsset) == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, bound asset not exit")
	}

	// transfer asset
	txOrigin := s.ContractRef().TxOrigin()
	msgSender := s.ContractRef().MsgSender()
	if err := transfer2Contract(s, input.FromAssetHash, txOrigin, input.Amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, transfer to contract failed, err: %v", err)
	}

	// generate tunnel data
	txData := zutils.EncodeTxArgs(toAsset, input.ToAddress, input.Amount)
	tunnel := &zutils.TunnelData{
		Caller:     this,
		ToContract: targetProxy,
		Method:     []byte("unlock"),
		TxData:     txData,
	}
	tunnelData, err := tunnel.Encode()
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to encode tunnel data, err: %v", err)
	}

	// get and set tx index
	lastTxIndex := getTxIndex(s)
	storeTxIndex(s, new(big.Int).Add(lastTxIndex, common.Big1))
	txIndex := getTxIndex(s)
	if txIndex.Cmp(common.Big0) <= 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, tx index can not be zero")
	}
	paramTxHash := scom.Uint256ToBytes(txIndex)
	crossChainID := zutils.GenerateCrossChainID(this, paramTxHash)

	// assemble tx, generate and store cross chain transaction proof
	txParams, txParamsEnc, proof := zutils.EncodeMakeTxParams(
		paramTxHash,
		crossChainID,
		this[:],
		input.ToChainId,
		targetCaller,
		"unwrap",
		tunnelData,
	)

	storeTxProof(s, paramTxHash, proof)
	storeTxParams(s, paramTxHash, txParamsEnc)

	// emit event log
	if err := emitCrossChainEvent(s, txOrigin, paramTxHash, this, input.ToChainId, toAsset, "unlock", txData); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to emit `CrossChainEvent` log, err: %v", err)
	}
	if err := emitLockEvent(s, input.FromAssetHash, msgSender, input.ToChainId, toAsset, input.ToAddress, input.Amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to emit `LockEvent` log, err: %v", err)
	}

	// zion main chain DONT need a `relayer` to commit proof but directly stores the lock request.
	// but we should ensure that `relayer` of other chain can deserialize the main chain events correctly.
	if err := scom.MakeTransaction(s, txParams, native.ZionMainChainID); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to makeTransaction, err: %v", err)
	}
	return utils.ByteSuccess, nil
}

func Unlock(s *native.NativeContract, entranceParams *scom.EntranceParam, txParams *scom.MakeTxParam) error {
	args, err := zutils.DecodeTxArgs(txParams.Args)
	if err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to decode txArgs, err: %v", err)
	}

	// filter duplicate tx
	if err := scom.CheckDoneTx(s, txParams.CrossChainID, entranceParams.SourceChainID); err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to check done transaction, err:%s", err)
	}

	// transfer asset
	toAsset := common.BytesToAddress(args.ToAssetHash)
	toAddress := common.BytesToAddress(args.ToAddress)
	if err := transferFromContract(s, toAsset, toAddress, args.Amount); err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to transfer asset, err: %v", err)
	}

	// store done tx
	if err := scom.PutDoneTx(s, txParams.CrossChainID, entranceParams.SourceChainID); err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to put done tx, err:%s", err)
	}

	// emit event logs
	if err := emitUnlockEvent(s, toAsset, toAddress, args.Amount); err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to emit `UnlockEvent`, err: %v", err)
	}

	crossChainTxHash := s.ContractRef().TxHash().Bytes()
	if err := emitVerifyHeaderAndExecuteTxEvent(s,
		entranceParams.SourceChainID,
		args.ToAssetHash,
		crossChainTxHash,
		txParams.TxHash,
	); err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to emit `VerifyHeaderAndExecuteTxEvent`, err: %v", err)
	}

	return nil
}
