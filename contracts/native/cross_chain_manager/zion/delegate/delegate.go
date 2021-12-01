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

package delegate

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core"
)

var (
	this = utils.LockProxyContractAddress

	// allow user transfer `lock amount` to `lock proxy` without handing fee, because there are
	// always two txs in the entire cross chain procedure, the first step is an source tx in some
	// chain which only cost tiny gas usage across user's account. the next tx will be sent by an
	// relayer which should cost some handing fee.
	//
	// and the params of `allowNoDelegateContract` only used in debug mode.
	//
	allowNoDelegateContract = true
)

func Approve(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	input := new(MethodApproveInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Approve, failed to decode params, err: %v", err)
	}
	if input.Spender == common.EmptyAddress {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Approve, spender address invalid")
	}
	if input.Amount == nil || input.Amount.Cmp(common.Big0) <= 0 {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Approve, amount invalid")
	}

	owner := s.ContractRef().TxOrigin()
	spender := input.Spender
	amount := input.Amount

	setAllowance(s, owner, spender, amount)
	if err := emitApprovedEvent(s, owner, spender, amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Approve, failed to emit `ApprovedEvent`, err: %v", err)
	}
	return utils.PackOutputs(ABI, MethodApprove, true)
}

func Allowance(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	input := new(MethodAllowanceInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Allowance, failed to decode params, err: %v", err)
	}
	if input.Owner == common.EmptyAddress {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Allowance, owner address invalid")
	}
	if input.Spender == common.EmptyAddress {
		return utils.ByteFailed, fmt.Errorf("LockProxy.Allowance, spender address invalid")
	}

	data := getAllowance(s, input.Owner, input.Spender)
	return utils.PackOutputs(ABI, MethodAllowance, data)
}

func SafeTransfer2Contract(s *native.NativeContract, from common.Address, amount *big.Int) error {
	isWrapperCaller, err := checkOutcome(s, from, amount)
	if err != nil {
		return err
	}

	if isWrapperCaller {
		delegator := s.ContractRef().TxTo()
		if err := nativeTransfer(s, delegator, this, amount); err != nil {
			return err
		}
	}
	return nil
}

// SafeTransferFromContract entrance must be some cross chain manager contract address
func SafeTransferFromContract(s *native.NativeContract, entrance, to common.Address, amount *big.Int) error {
	ctx := s.ContractRef().CurrentContext()
	caller := ctx.Caller
	from := s.ContractRef().TxOrigin()
	txTo := s.ContractRef().TxTo()

	if caller == common.EmptyAddress || from == common.EmptyAddress || txTo == common.EmptyAddress {
		return fmt.Errorf("tx context params invalid")
	}
	if to == common.EmptyAddress {
		return fmt.Errorf("invalid dest account")
	}
	if amount == nil || amount.Cmp(common.Big0) <= 0 {
		return fmt.Errorf("invalid amount")
	}

	// caller MUST be from, e.g: some relayer
	if caller != from {
		return fmt.Errorf("caller must be equal to tx.from")
	}
	// tx must be sent to cross chain manager contract, and the native asset which locked before can be unlocked after proof validation.
	if txTo != entrance {
		return fmt.Errorf("the tx.to should be cross chain manager contract address")
	}

	return nativeTransfer(s, this, to, amount)
}

func SubBalance(s *native.NativeContract, from common.Address, amount *big.Int) error {
	isWrapperCaller, err := checkOutcome(s, from, amount)
	if err != nil {
		return err
	}

	delegator := s.ContractRef().TxTo()
	if isWrapperCaller {
		s.StateDB().SubBalance(delegator, amount)
	}
	return nil
}

// AddBalance entrance should be some eccm contract address
func AddBalance(s *native.NativeContract, entrance, to common.Address, amount *big.Int) error {
	ctx := s.ContractRef().CurrentContext()
	txTo := s.ContractRef().TxTo()

	if ctx.Caller != entrance {
		return fmt.Errorf("caller should be eccm contract address")
	}
	if txTo != entrance {
		return fmt.Errorf("tx.to must be eccm contract address")
	}

	if amount == nil || amount.Cmp(common.Big0) <= 0 {
		return fmt.Errorf("invalid amount")
	}
	s.StateDB().AddBalance(to, amount)
	return nil
}

// checkOutcome return isWrapper caller and error
func checkOutcome(s *native.NativeContract, from common.Address, amount *big.Int) (bool, error) {
	ctx := s.ContractRef().CurrentContext()
	owner := s.ContractRef().TxOrigin()
	caller := ctx.Caller
	spender := this
	delegator := s.ContractRef().TxTo()
	value := s.ContractRef().Value()

	if value == nil || value.Cmp(common.Big0) == 0 {
		return false, fmt.Errorf("tx value should be greater than zero")
	}
	if caller == common.EmptyAddress {
		return false, fmt.Errorf("caller is invalid")
	}
	if owner != from {
		return false, fmt.Errorf("only tx.from can sub balance")
	}

	// if user sender tx by himself, the context caller should be equal to `tx.From`, and the
	// `tx.Value` should be equal to transfer amount. and `tx.Value` will be sub from `tx.From`
	// in evm handler.
	// in this condition, the `tx.value` only contains user lock amount, gas cost is calculated
	// on the remaining part of the account balance.
	//
	// if user transfer native token with an wrapper/delegate contract, user should approve the
	// lock proxy contract with enough amount first.
	// in this condition, the context caller is not the `tx.From` but an wrapper/delegate contract
	// address. and there are 2 parts in `tx.value`, wrapper/delegate handling fee and the
	// `lock amount`. because the `tx.value` will be transferred to the wrapper/delegate contract
	// before native contract executing, so we need to transfer the `lock amount` from wrapper/delegate
	// contract to `lock proxy`.
	isWrapper := caller != owner
	if isWrapper {
		if delegator != caller {
			return isWrapper, fmt.Errorf("delegator should be some wrapper contract, and this contract is the only caller")
		}
		if value.Cmp(amount) < 0 {
			return isWrapper, fmt.Errorf("tx value is not enough, it should be greater than %v", amount)
		}
		allowance := getAllowance(s, owner, spender)
		if allowance.Cmp(amount) < 0 {
			return isWrapper, fmt.Errorf("allowance not enought, expect %v, got %v", amount, allowance)
		}
		resAllowance := new(big.Int).Sub(allowance, amount)
		setAllowance(s, owner, spender, resAllowance)
	} else {
		if !allowNoDelegateContract {
			return isWrapper, fmt.Errorf("sub balance without delegate contract is forbidden!")
		}
		if delegator != this {
			return isWrapper, fmt.Errorf("invalid delegator, tx.to should be lock proxy address")
		}
		if value.Cmp(amount) != 0 {
			return isWrapper, fmt.Errorf("transfer amount %v not equal to tx.value %v", amount, value)
		}
	}
	return isWrapper, nil
}

func nativeTransfer(s *native.NativeContract, from, to common.Address, amount *big.Int) error {
	if !core.CanTransfer(s.StateDB(), from, amount) {
		return fmt.Errorf("%s insufficient balance", from.Hex())
	}
	core.Transfer(s.StateDB(), from, to, amount)
	return nil
}
