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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	zutils "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/utils"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/main_chain_lock_proxy"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/params"
)

const MIN_BALANCE = 1000000 // default minimum value of lock proxy balance

var (
	gasTable = map[string]uint64{
		MethodName: 0,
		MethodLock: 10000,
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
	s.Register(MethodLock, Lock)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return new(MethodContractNameOutput).Encode()
}

func Lock(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	sourceChainID := native.ZionMainChainID

	input := new(MethodLockInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to decode params, err: %v", err)
	}
	if input.FromAssetHash != common.EmptyAddress {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, only support native token")
	}
	if input.ToChainId == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, target chain id invalid")
	}
	if input.ToChainId == sourceChainID {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, target chain id wont be %d", sourceChainID)
	}
	if input.ToAddress == nil || len(input.ToAddress) == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, target address invalid")
	}
	if input.Amount == nil || input.Amount.Cmp(common.Big0) == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, amount invalid")
	}
	if input.Amount.Cmp(s.ContractRef().Value()) != 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, amount != tx.value")
	}

	// input fields alias, caller is proxy itself and `toContract` is `sideChain` proxy, which has the same address
	txOrigin := s.ContractRef().TxOrigin()
	msgSender := s.ContractRef().MsgSender()
	fromAsset := common.EmptyAddress
	toAsset := common.EmptyAddress.Bytes()
	toAddr := input.ToAddress
	amount := input.Amount
	toChainID := input.ToChainId
	caller := this
	toContract := this
	toMethod := "mint"

	// check side chain registered
	if sideChain, err := side_chain_manager.GetSideChain(s, toChainID); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to get side chain %d, err: %v", toChainID, err)
	} else if sideChain == nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, side chain %d is nil", toChainID)
	}

	// serialize tx args
	txData := zutils.EncodeTxArgs(toAsset, toAddr, amount)

	// get and set tx index
	txIndex, err := getNextTxIndex(s)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to get next tx index, err: %v", err)
	}
	paramTxHash := scom.Uint256ToBytes(txIndex)
	crossChainID := zutils.GenerateCrossChainID(this, paramTxHash)

	// check and store `doneTx`
	if err := scom.CheckDoneTx(s, crossChainID, sourceChainID); err != nil {
		return nil, fmt.Errorf("LockProxy.Lock, failed to check cross transaction, err: %v", err)
	}
	if err := scom.PutDoneTx(s, crossChainID, sourceChainID); err != nil {
		return nil, fmt.Errorf("LockProxy.Lock, faield to store cross transaction, err: %v", err)
	}

	// assemble tx, generate and store cross chain transaction proof
	txParams, rawTx, err := zutils.EncodeMakeTxParams(paramTxHash, crossChainID, caller[:], toChainID, toContract[:], toMethod, txData)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to encode `makeTxParams`, err: %v", err)
	}

	// emit event log
	if err := emitCrossChainEvent(s, txOrigin, paramTxHash, caller, toChainID, toAsset, rawTx); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to emit `CrossChainEvent` log, err: %v", err)
	}
	if err := emitLockEvent(s, fromAsset, msgSender, toChainID, toAsset, toAddr, amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to emit `LockEvent` log, err: %v", err)
	}

	// zion main chain DONT need a `relayer` to commit proof but directly stores the lock request.
	// but we should ensure that `relayer` of other chain can deserialize the main chain events correctly.
	if err := scom.MakeTransaction(s, txParams, sourceChainID); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to makeTransaction, err: %v", err)
	}
	return utils.ByteSuccess, nil
}

func Unlock(s *native.NativeContract, entranceParams *scom.EntranceParam, txParams *scom.MakeTxParam) error {
	args, err := zutils.DecodeTxArgs(txParams.Args)
	if err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to decode txArgs, err: %v", err)
	}

	// check params
	sourceChainID := entranceParams.SourceChainID
	if sourceChainID == native.ZionMainChainID || sourceChainID == 0 {
		return fmt.Errorf("LockProxy.Unlock, source chain id invalid")
	}
	if txParams.ToChainID != native.ZionMainChainID {
		return fmt.Errorf("LockProxy.Unlock, target chain id invalid")
	}
	if common.BytesToAddress(txParams.ToContractAddress) != this {
		return fmt.Errorf("LockProxy.Unlock, target contract is invalid")
	}
	if txParams.Method != "unlock" {
		return fmt.Errorf("LockProxy.Unlock, method is invalid")
	}
	toCaller := common.BytesToAddress(txParams.ToContractAddress)
	if toCaller != this {
		return fmt.Errorf("LockProxy.Unlock, target caller is invalid")
	}

	// filter duplicate tx
	if err := scom.CheckDoneTx(s, txParams.CrossChainID, sourceChainID); err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to check done transaction, err:%s", err)
	}

	// transfer asset
	toAddress := common.BytesToAddress(args.ToAddress)
	if err := transferFromContract(s, toAddress, args.Amount); err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to transfer asset, err: %v", err)
	}

	// store done tx
	if err := scom.PutDoneTx(s, txParams.CrossChainID, sourceChainID); err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to put done tx, err:%s", err)
	}

	// emit event logs
	toAsset := common.EmptyAddress
	if err := emitUnlockEvent(s, toAsset, toAddress, args.Amount); err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to emit `UnlockEvent`, err: %v", err)
	}

	crossChainTxHash := s.ContractRef().TxHash().Bytes()
	if err := emitVerifyHeaderAndExecuteTxEvent(s,
		sourceChainID,
		args.ToAssetHash,
		crossChainTxHash,
		txParams.TxHash,
	); err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to emit `VerifyHeaderAndExecuteTxEvent`, err: %v", err)
	}

	return nil
}
