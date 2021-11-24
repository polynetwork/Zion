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
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/auth"
	zutils "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/utils"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/auth_abi"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/side_chain_lock_proxy_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

var (
	gasTable = map[string]uint64{
		MethodName:      0,
		MethodBurn:      10000,
		MethodMint:      10000,
		MethodApprove:   10000,
		MethodAllowance: 0,
	}

	ccmp = common.HexToAddress("")
)

func InitLockProxy() {
	InitABI()
	native.Contracts[this] = RegisterLockProxyContract
}

func RegisterLockProxyContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodName, Name)
	s.Register(MethodBurn, Burn)
	s.Register(MethodMint, Mint)
	s.Register(MethodApprove, auth.Approve)
	s.Register(MethodAllowance, auth.Allowance)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return new(MethodContractNameOutput).Encode()
}

func Burn(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	from := s.ContractRef().TxOrigin()

	input := new(MethodBurnInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Burn, failed to decode params, err: %v", err)
	}
	if input.ToAddress == common.EmptyAddress {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Burn, invaild to address")
	}
	if from != input.ToAddress {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Burn, only allow self tx cross chain")
	}
	if input.Amount == nil || input.Amount.Cmp(common.Big0) <= 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Burn, invalid amount")
	}
	if input.ToChainId != native.ZionMainChainID || input.ToChainId == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Burn, dest chain id invalid")
	}

	// check and sub balance
	if err := auth.SubBalance(s, from, input.Amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Burn, failed to sub balance, err: %v", err)
	}

	rawArgs := zutils.EncodeTxArgs(common.EmptyAddress[:], input.ToAddress[:], input.Amount)
	eccm, err := getEthCrossChainManager(s, ccmp)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Burn, failed to get eccm address, err: %v", err)
	}
	if err := crossChain(s, eccm, this, input.ToChainId, "unlock", rawArgs); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Burn, failed to ")
	}

	asset := common.EmptyAddress
	if err := emitBurnEvent(s, asset, from, input.ToChainId, asset[:], input.ToAddress[:], input.Amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Burn, emit `BurnEvent` failed, err: %v", err)
	}

	return utils.ByteSuccess, nil
}

func Mint(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller

	input := new(MethodMintInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, failed to decode params, err: %v", err)
	}
	if input.ArgsBs == nil || input.FromContractAddr == nil || input.FromChainId == 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, invalid params")
	}

	eccm, err := getEthCrossChainManager(s, ccmp)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, failed to get eccm contract address, err: %v", err)
	}
	if caller != eccm {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, caller authority invalid, must be eccm")
	}

	args, err := zutils.DecodeTxArgs(input.ArgsBs)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, failed to decode args, err: %v", err)
	}
	if args.ToAssetHash == nil || args.ToAddress == nil || args.Amount == nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, args field invalid")
	}

	toAddr := common.BytesToAddress(args.ToAddress)
	asset := common.EmptyAddress
	if common.BytesToAddress(args.ToAssetHash) != asset {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, target asset invalid")
	}
	amount := args.Amount
	if amount.Cmp(common.Big0) <= 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, source amount invalid")
	}

	if err := auth.AddBalance(s, eccm, toAddr, amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, failed to add balance, err: %v", err)
	}

	if err := emitMintEvent(s, asset, toAddr, amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, failed to emit `MintEvent`, err: %v", err)
	}
	return utils.ByteSuccess, nil
}
