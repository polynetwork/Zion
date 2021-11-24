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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/auth"
	zutils "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/utils"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/auth_abi"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/main_chain_lock_proxy_abi"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

var (
	gasTable = map[string]uint64{
		MethodName:                   0,
		MethodLock:                   10000,
		MethodGetSideChainLockAmount: 0,
		MethodApprove:                10000,
		MethodAllowance:              0,
	}
)

func InitLockProxy() {
	InitABI()
	native.Contracts[this] = RegisterLockProxyContract
}

func RegisterLockProxyContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodName, Name)
	s.Register(MethodLock, Lock)
	s.Register(MethodGetSideChainLockAmount, GetSideChainLockAmount)
	s.Register(MethodApprove, auth.Approve)
	s.Register(MethodAllowance, auth.Allowance)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return new(MethodContractNameOutput).Encode()
}

func Lock(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	sourceChainID := native.ZionMainChainID
	txOrigin := s.ContractRef().TxOrigin()
	msgSender := s.ContractRef().MsgSender()

	input := new(MethodLockInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to decode params, err: %v", err)
	}
	if input.FromAssetHash != common.EmptyAddress {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, only support native token")
	}
	if input.ToChainId == 0 || input.ToChainId == sourceChainID {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, target chain id invalid")
	}
	if input.ToAddress == nil || len(input.ToAddress) == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, target address invalid")
	}
	if common.BytesToAddress(input.ToAddress) != txOrigin {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, target address MUST be tx origin")
	}
	if input.Amount == nil || input.Amount.Cmp(common.Big0) == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, amount invalid")
	}

	// input fields alias, caller is proxy itself and `toContract` is `sideChain` proxy, which has the same address
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
	} else if sideChain.Router != utils.ZION_ROUTER {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, side chain %d router is not zion", toChainID)
	}

	// lock token into lock proxy
	if err := auth.SafeTransfer2Contract(s, amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Lock, failed to transfer token to lock proxy, err: %v", err)
	}

	// set total amount
	addTotalAmount(s, toChainID, amount)

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

	// check side chain registered
	if sideChain, err := side_chain_manager.GetSideChain(s, sourceChainID); err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to get side chain %d, err: %v", sourceChainID, err)
	} else if sideChain == nil {
		return fmt.Errorf("LockProxy.Unlock, side chain %d is nil", sourceChainID)
	} else if sideChain.Router != utils.ZION_ROUTER {
		return fmt.Errorf("LockProxy.Unlock, side chain %d router is not zion", sourceChainID)
	}

	// check contracts
	if txParams.ToContractAddress == nil || common.BytesToAddress(txParams.ToContractAddress) != this {
		return fmt.Errorf("LockProxy.Unlock, target contract is invalid")
	}
	if txParams.FromContractAddress == nil || common.BytesToAddress(txParams.FromContractAddress) != this {
		return fmt.Errorf("LockProxy.Unlock, source contract is invalid")
	}
	if txParams.Method != "unlock" {
		return fmt.Errorf("LockProxy.Unlock, method is invalid")
	}

	// check asset
	toAsset := common.BytesToAddress(args.ToAssetHash)
	if toAsset != common.EmptyAddress {
		return fmt.Errorf("LockProxy.Unlock, to asset invalid, %s", toAsset.Hex())
	}

	// do not need to check `request` and `DoneTx`, there were just settled at `entrance` while relayer send `commitProof`

	// transfer asset
	toAddress := common.BytesToAddress(args.ToAddress)
	entrance := s.ContractRef().CurrentContext().Caller
	if err := auth.SafeTransferFromContract(s, entrance, toAddress, args.Amount); err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to transfer native token, err: %v", err)
	}
	if err := subTotalAmount(s, sourceChainID, args.Amount); err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to sub total amount, err: %v", err)
	}

	// emit event logs
	if err := emitUnlockEvent(s, toAsset, toAddress, args.Amount); err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to emit `UnlockEvent`, err: %v", err)
	}

	crossChainTxHash := s.ContractRef().TxHash()
	if err := emitVerifyHeaderAndExecuteTxEvent(s,
		sourceChainID,
		toAsset[:],
		crossChainTxHash[:],
		txParams.TxHash,
	); err != nil {
		return fmt.Errorf("LockProxy.Unlock, failed to emit `VerifyHeaderAndExecuteTxEvent`, err: %v", err)
	}

	return nil
}

func GetSideChainLockAmount(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	failed := common.Big0.Bytes()

	input := new(MethodGetSideChainLockAmountInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return failed, fmt.Errorf("LockProxy.GetSideChainLockAmount, failed to decode params")
	}

	// ignore side chain register checking
	amount := getTotalAmount(s, input.ChainId)
	return amount.Bytes(), nil
}
