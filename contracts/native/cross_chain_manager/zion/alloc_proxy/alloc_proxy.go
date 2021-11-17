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
	if input.Proof == nil || input.Header == nil || input.Epoch == nil || input.Extra == nil ||
		len(input.Proof) == 0 || len(input.Header) == 0 || len(input.Epoch) == 0 || len(input.Extra) == 0 {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.InitGenesisHeader, invalid params")
	}

	if _, exist, _ := getEpoch(s); exist != nil {
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
	if err := emitInitGenesisBlockEvent(s, header.Number, input.Header, input.Epoch); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.InitGenesisHeader, failed to emit `InitGenesisBlockEvent`, err: %v", err)
	}
	return utils.ByteSuccess, nil
}

// ChangeEpoch update bookeepers for main chain. in hotstuff consensus, epoch changed at the height of
// `epoch.StartHeight` - 1 denotes that if the epoch start height is 1000, the block header of 999 carry
// bookeepers addresses and these new bookeepers will participant in consensus after block 999(current
// block of 999 still verified by old bookeepers). so we should ensure that new epoch's `startHeight` is
// higher than last epoch's `startHeight` and current header which carry proof of `epochChange`.
func ChangeEpoch(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()

	input := new(MethodChangeEpochInput)
	if err := input.Decode(ctx.Payload); err != nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, failed to decode params, err: %v", err)
	}
	if input.Extra == nil || input.Epoch == nil || input.Header == nil || input.Proof == nil ||
		len(input.Proof) == 0 || len(input.Header) == 0 || len(input.Epoch) == 0 || len(input.Extra) == 0 {
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
	if lastEpoch == nil {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, last epoch is nil")
	}
	if lastEpoch.Hash() == epoch.Hash() {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, duplicate epoch")
	}
	if header.Number.Uint64()+1 != epoch.StartHeight {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, header height %v + 1 should be equals to %d",
			header.Number, epoch.StartHeight)
	}
	if lstEpID := lastEpoch.ID; epoch.ID <= lstEpID {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, epoch ID should be greater than %d", lstEpID)
	}
	if lstEpStartNo := lastEpoch.StartHeight; header.Number.Uint64() <= lstEpStartNo || epoch.StartHeight <= lstEpStartNo {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.ChangeEpoch, header number should > %d", lstEpStartNo)
	}

	nextEpochStartHeight, nextEpochVals, err := zion.VerifyHeader(header, lastEpoch.MemberList(), true)
	if err != nil {
		return nil, fmt.Errorf("AllocProxy.ChangeEpoch, failed to verify header, err: %v", err)
	}
	if nextEpochStartHeight != epoch.StartHeight {
		return nil, fmt.Errorf("AllocProxy.ChangeEpoch, failed to verify header, err: epoch start height expect %d got %d",
			nextEpochStartHeight, epoch.StartHeight)
	}
	if curEpochVals := epoch.MemberList(); !compareVals(nextEpochVals, curEpochVals) {
		return nil, fmt.Errorf("AllocProxy.ChangeEpoch, failed to verify header, err: vals expect %v got %v",
			nextEpochVals, curEpochVals)
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
	if input.ToAddress == common.EmptyAddress {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, invaild to address")
	}
	if from != input.ToAddress {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, only allow self tx cross chain")
	}
	if input.Amount == nil || input.Amount.Cmp(common.Big0) <= 0 {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, invalid amount")
	}
	if input.ToChainId == native.ZionMainChainID || input.ToChainId == 0 {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, dest chain id invalid")
	}

	// check side chain
	if srcChain, err := side_chain_manager.GetSideChain(s, input.ToChainId); err != nil {
		return nil, fmt.Errorf("AllocProxy.Burn, failed to get side chain, err: %v", err)
	} else if srcChain == nil {
		return nil, fmt.Errorf("AllocProxy.Burn, side chain %d is not registered", input.ToChainId)
	}

	// check and sub balance
	if s.StateDB().GetBalance(from).Cmp(input.Amount) <= 0 {
		return utils.ByteFailed, fmt.Errorf("AllocProxy.Burn, Insufficient balance")
	}
	s.StateDB().SubBalance(from, input.Amount)

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

	input := new(MethodVerifyHeaderAndMintInput)
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
