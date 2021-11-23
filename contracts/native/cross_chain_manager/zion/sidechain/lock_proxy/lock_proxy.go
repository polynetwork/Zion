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

	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/auth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	zutils "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/utils"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/auth_abi"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/side_chain_lock_proxy_abi"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/zion"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core"
)

// zion `alloc` module used to implement native asset cross chain between main/side chain.
// this module implement interface of `mint` and `burn` of native token between zion main chain and side chain.
// the native asset only alloc in genesis block on the main chain, and the side chain do not allow to do this.
//
// main chain only use the interface of `burn`, and side chain only use the interface of `mint`.
// the amount of `mint` and `burn` should always be the same.
//
// user get zion native token at main chain and only allow to burn/mint asset for it's self.
//
// sequence of logic:
// main chain alloc -> exchange -> user burn on main net -> side chain relayer `verifyHeaderAndMint` -> user add balance
//

var (
	gasTable = map[string]uint64{
		MethodName:                     0,
		MethodBurn:                     10000,
		MethodVerifyHeaderAndExecuteTx: 10000,
		MethodApprove:                  10000,
		MethodAllowance:                0,
	}

	IsMainChain bool
)

func init() {
	core.SetMainChain = func(_isMainChain bool) {
		IsMainChain = _isMainChain
	}
}

func InitLockProxy() {
	InitABI()
	native.Contracts[this] = RegisterLockProxyContract
}

func RegisterLockProxyContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodName, Name)
	s.Register(MethodBurn, Burn)
	s.Register(MethodVerifyHeaderAndExecuteTx, Mint)
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
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, failed to decode params, err: %v", err)
	}
	if input.ToAddress == common.EmptyAddress {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, invaild to address")
	}
	if from != input.ToAddress {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, only allow self tx cross chain")
	}
	if input.Amount == nil || input.Amount.Cmp(common.Big0) <= 0 {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, invalid amount")
	}
	if input.ToChainId != native.ZionMainChainID || input.ToChainId == 0 {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, dest chain id invalid")
	}

	// check side chain
	if srcChain, err := side_chain_manager.GetSideChain(s, input.ToChainId); err != nil {
		return nil, fmt.Errorf("AllocProxy.Burn, failed to get side chain, err: %v", err)
	} else if srcChain == nil {
		return nil, fmt.Errorf("AllocProxy.Burn, side chain %d is not registered", input.ToChainId)
	}

	// check and sub balance
	if err := auth.SubBalance(s, from, input.Amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, Insufficient balance")
	}

	// store cross tx data
	lastCrossTxIndex := getCrossTxIndex(s)
	nextCrossTxIndex := lastCrossTxIndex + 1
	storeCrossTxIndex(s, nextCrossTxIndex)
	crossTxIndex := getCrossTxIndex(s)
	if crossTxIndex != nextCrossTxIndex {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, store cross tx failed, expect tx index %d, got %d", nextCrossTxIndex, crossTxIndex)
	}

	crossTx := &CrossTx{
		ToChainId:   input.ToChainId,
		FromAddress: from,
		ToAddress:   input.ToAddress,
		Amount:      input.Amount,
		Index:       crossTxIndex,
	}
	proof := crossTx.Proof()
	storeCrossTxProof(s, crossTxIndex, proof)

	if err := emitBurnEvent(s, input.ToChainId, from, input.ToAddress, input.Amount, crossTxIndex); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, emit `BurnEvent` failed, err: %v", err)
	}
	return utils.ByteSuccess, nil
}

func Mint(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	input := new(MethodVerifyHeaderAndExecuteTxInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, failed to decode params, err: %v", err)
	}
	if input.Header == nil || input.Proof == nil || input.RawCrossTx == nil || input.Extra == nil ||
		len(input.Proof) == 0 || len(input.Header) == 0 || len(input.RawCrossTx) == 0 || len(input.Extra) == 0 {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, invalid params")
	}

	// deserialize parameters
	header, err := DecodeHeader(input.Header)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, failed to decode header, err: %v", err)
	}
	crossTx, err := DecodeCrossTx(input.RawCrossTx)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, failed to decode cross tx, err: %v", err)
	}
	if crossTx.ToChainId == native.ZionMainChainID || crossTx.ToChainId == 0 {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, target chain id invalid")
	}
	if crossTx.Amount == nil || crossTx.Amount.Uint64() == 0 {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, amount invalid")
	}
	if crossTx.ToAddress == common.EmptyAddress {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, to address invalid")
	}
	if crossTx.Index == 0 {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, invalid cross tx index")
	}
	if crossTx.FromAddress == common.EmptyAddress {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, invalid cross tx from address")
	}
	if crossTx.FromAddress != crossTx.ToAddress {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, only allow self tx cross chain")
	}

	// check and store cross tx
	if exist, _ := getCrossTxContent(s, crossTx.Hash()); exist != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, duplicate cross tx %s", exist.Hash().Hex())
	}
	if err := storeCrossTxContent(s, crossTx); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, failed to store cross tx, err: %v", err)
	}

	// get and check epoch start height
	_, epoch, err := getEpoch(s)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, failed to get epoch, err: %v", err)
	}
	if epoch == nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, current epoch is nil")
	}
	// allow epoch.startHeight == header.Num, there may be some cross tx happened at the first block of new epoch
	if epoch.StartHeight < header.Number.Uint64() {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, header number should be greater than %d", epoch.StartHeight)
	}

	if _, _, err := zion.VerifyHeader(header, epoch.MemberList(), false); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, failed to verify header, err: %v", err)
	}
	if _, err := zutils.VerifyTx(input.Proof, header, this, input.Extra, true); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, failed to verify proof, err: %v", err)
	}

	s.StateDB().AddBalance(crossTx.ToAddress, crossTx.Amount)

	if err := emitMintEvent(s, crossTx.ToChainId, crossTx.FromAddress, crossTx.ToAddress, crossTx.Amount); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, failed to emit `MintEvent`, err: %v", err)
	}
	return utils.ByteSuccess, nil
}
