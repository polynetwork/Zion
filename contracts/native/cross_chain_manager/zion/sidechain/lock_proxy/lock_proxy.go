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
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/delegate"
	zutils "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/utils"
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

	ccmp = common.HexToAddress("0xc6195336878Fc34B1b5A13895015a97c1aD9cc25")
)

func InitLockProxy() {
	InitABI()
	delegate.InitABI(ABI)

	native.Contracts[this] = RegisterLockProxyContract
}

func RegisterLockProxyContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodName, Name)
	s.Register(MethodBurn, Burn)
	s.Register(MethodMint, Mint)
	s.Register(MethodApprove, delegate.Approve)
	s.Register(MethodAllowance, delegate.Allowance)
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
	if input.Amount == nil || input.Amount.Cmp(common.Big0) <= 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Burn, invalid amount")
	}
	if input.ToChainId != native.ZionMainChainID {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Burn, dest chain id invalid")
	}

	amount := input.Amount
	asset := common.EmptyAddress
	toAddr := input.ToAddress[:]
	toChainID := input.ToChainId

	// check and sub balance
	if err := delegate.SubBalance(s, from, amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Burn, failed to sub balance, err: %v", err)
	}

	rawArgs, err := zutils.EncodeTxArgs(asset[:], toAddr, amount)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Burn, failed to encode txArgs, err: %v", err)
	}
	eccm, err := getEthCrossChainManager(s, ccmp)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Burn, failed to get eccm address, err: %v", err)
	}
	if err := crossChain(s, eccm, this, toChainID, "unlock", rawArgs); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Burn, failed to ")
	}

	if err := emitBurnEvent(s, asset, from, toChainID, asset[:], toAddr, amount); err != nil {
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
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, DANGER! caller is not eccm!")
	}

	args, err := zutils.DecodeTxArgs(input.ArgsBs)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, failed to decode args, err: %v", err)
	}
	if args == nil || args.ToAssetHash == nil || args.ToAddress == nil || args.Amount == nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, args field invalid")
	}

	toAddr := common.BytesToAddress(args.ToAddress)
	asset := common.EmptyAddress
	amount := args.Amount
	if common.BytesToAddress(args.ToAssetHash) != asset {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, target asset invalid")
	}
	if amount.Cmp(common.Big0) <= 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, source amount invalid")
	}

	if err := delegate.AddBalance(s, eccm, toAddr, amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, failed to add balance, err: %v", err)
	}

	if err := emitMintEvent(s, asset, toAddr, amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Mint, failed to emit `MintEvent`, err: %v", err)
	}
	return utils.ByteSuccess, nil
}
