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

package alloc_proxy

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	zutils "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/utils"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/alloc_proxy"
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

var (
	gasTable = map[string]uint64{
		MethodName:                0,
		MethodInitGenesisHeader:   10000,
		MethodChangeEpoch:         10000,
		MethodBurn:                10000,
		MethodVerifyHeaderAndMint: 10000,
	}

	IsMainChain bool
)

func init() {
	core.SetMainChain = func(_isMainChain bool) {
		IsMainChain = _isMainChain
	}
}

func InitAllocProxy() {
	InitABI()
	native.Contracts[this] = RegisterAllocProxyContract
}

func RegisterAllocProxyContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodName, Name)
	if IsMainChain {
		s.Register(MethodBurn, Burn)
	} else {
		s.Register(MethodInitGenesisHeader, InitGenesisHeader)
		s.Register(MethodChangeEpoch, ChangeEpoch)
		s.Register(MethodVerifyHeaderAndMint, Mint)
	}
}

func Name(s *native.NativeContract) ([]byte, error) {
	return new(MethodContractNameOutput).Encode()
}

// InitGenesisHeader store first header and epoch, this epoch should contains consensus participants.
func InitGenesisHeader(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	input := new(MethodInitGenesisHeaderInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.InitGenesisHeader, failed to decode params, err: %v", err)
	}
	if input.Proof == nil || input.Header == nil || input.Epoch == nil || input.Extra == nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.InitGenesisHeader, invalid params")
	}

	_, existEpoch, _ := getEpoch(s)
	if existEpoch != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.InitGenesisHeader, genesis header already exist")
	}

	header, err := DecodeHeader(input.Header)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.InitGenesisHeader, failed to unmarshal header, err: %v", err)
	}
	epoch, err := DecodeEpoch(input.Epoch)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.InitGenesisHeader, failed to decode epoch, err: %v", err)
	}
	if _, _, err := zion.VerifyHeader(header, epoch.MemberList(), false); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.InitGenesisHeader, failed to verify header, err: %v", err)
	}
	if _, err := zutils.VerifyTx(input.Proof, header, utils.NodeManagerContractAddress, input.Extra, true); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.InitGenesisHeader, failed to verify proof, err: %v", err)
	}

	storeEpoch(s, input.Epoch)
	if err := emitInitGenesisBlockEvent(s, header.Number.Uint64(), input.Header, input.Epoch); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.InitGenesisHeader, failed to emit `InitGenesisBlockEvent`, err: %v", err)
	}
	return utils.ByteSuccess, nil
}

func ChangeEpoch(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	input := new(MethodChangeEpochInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, failed to decode params, err: %v", err)
	}
	if input.Extra == nil || input.Epoch == nil || input.Header == nil || input.Proof == nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, invalid params")
	}

	header, err := DecodeHeader(input.Header)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, failed to unmarshal header, err: %v", err)
	}
	epoch, err := DecodeEpoch(input.Epoch)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, failed to decode epoch, err: %v", err)
	}

	lastEpochEnc, lastEpoch, err := getEpoch(s)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, failed to get last epoch, err: %v", err)
	}
	if lastEpoch.Hash() == epoch.Hash() {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, duplicate epoch")
	}
	if epoch.ID <= lastEpoch.ID {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, epoch ID should be greater than %d", lastEpoch.ID)
	}
	if header.Number.Uint64() <= lastEpoch.StartHeight || epoch.StartHeight <= lastEpoch.StartHeight {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, invalid header height, require greater than %d", lastEpoch.StartHeight)
	}

	nextEpochStartHeight, nextEpochVals, err := zion.VerifyHeader(header, lastEpoch.MemberList(), true)
	if err != nil {
		return nil, fmt.Errorf("AllocProxy.ChangeEpoch, failed to verify header, err: %v", err)
	}
	if nextEpochStartHeight != epoch.StartHeight {
		return nil, fmt.Errorf("AllocProxy.ChangeEpoch, failed to verify header, err: epoch start height expect %d got %d", nextEpochStartHeight, epoch.StartHeight)
	}
	if curEpochVals := epoch.MemberList(); !isSameVals(nextEpochVals, curEpochVals) {
		return nil, fmt.Errorf("AllocProxy.ChangeEpoch, failed to verify header, err: vals expect %v got %v", nextEpochVals, curEpochVals)
	}
	if _, err := zutils.VerifyTx(input.Proof, header, utils.NodeManagerContractAddress, input.Extra, true); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, failed to verify proof, err: %v", err)
	}

	storeEpoch(s, input.Epoch)
	if err := emitChangeEpochEvent(s, header.Number, input.Header, lastEpochEnc, input.Epoch); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, failed to emit `ChangeEpochEvent`, err: %v", err)
	}
	return utils.ByteSuccess, nil
}

func Burn(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	from := s.ContractRef().TxOrigin()

	input := new(MethodBurnInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, failed to decode params, err: %v", err)
	}
	// todo: get side chain
	if input.ToChainId == native.ZionMainChainID || input.ToChainId == 0 {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, dest chain id invalid")
	}
	if input.Amount == nil || input.Amount.Cmp(common.Big0) <= 0 {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, invalid amount")
	}

	// check and sub balance
	if s.StateDB().GetBalance(from).Cmp(input.Amount) <= 0 {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, Insufficient balance")
	}
	s.StateDB().SubBalance(from, input.Amount)

	lastCrossTxIndex := getCrossTxIndex(s)
	crossTxIndex := lastCrossTxIndex + 1
	crossTx := &CrossTx{
		ToChainId:   input.ToChainId,
		FromAddress: from,
		ToAddress:   input.ToAddress,
		Amount:      input.Amount,
		Index:       crossTxIndex,
	}

	if err := storeCrossTxContent(s, crossTx); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, failed to store cross tx content, err: %v", err)
	}
	if err := emitBurnEvent(s, input.ToChainId, from, input.ToAddress, input.Amount, crossTxIndex); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, emit `BurnEvent` failed, err: %v", err)
	}
	return utils.ByteSuccess, nil
}

func Mint(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	input := new(MethodVerifyHeaderAndMintInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, failed to decode params, err: %v", err)
	}

	if input.Header == nil || input.Proof == nil || input.RowCrossTx == nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, invalid params")
	}

	header, err := DecodeHeader(input.Header)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, failed to decode header, err: %v", err)
	}

	_, epoch, err := getEpoch(s)
	if err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, failed to get epoch, err: %v", err)
	}
	if epoch.StartHeight < header.Number.Uint64() {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, header number should be greater than %d", epoch.StartHeight)
	}
	if _, _, err := zion.VerifyHeader(header, epoch.MemberList(), false); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, failed to verify header, err: %v", err)
	}
	if _, err := zutils.VerifyTx(input.Proof, header, this, input.RowCrossTx, true); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Mint, failed to verify proof, err: %v", err)
	}
	return utils.ByteSuccess, nil
}
