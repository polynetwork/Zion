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

package cross_chain_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/bsc"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/cosmos"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/eth"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/eth_common"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/heco"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/msc"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/no_proof"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/okex"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/quorum"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zilliqa"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "cross chain manager"

const (
	BLACKED_CHAIN = "BlackedChain"
)

var (
	this     = native.NativeContractAddrMap[native.NativeCrossChain]
	gasTable = map[string]uint64{
		scom.MethodContractName:        0,
		scom.MethodImportOuterTransfer: 0,
		scom.MethodMultiSign:           100000,
		scom.MethodBlackChain:          0,
		scom.MethodWhiteChain:          0,
		scom.MethodCheckDone:           0,
	}
)

func InitCrossChainManager() {
	native.Contracts[this] = RegisterCrossChainManagerContract
}

func RegisterCrossChainManagerContract(s *native.NativeContract) {
	s.Prepare(scom.ABI, gasTable)

	s.Register(scom.MethodContractName, Name)
	s.Register(scom.MethodImportOuterTransfer, ImportOuterTransfer)
	s.Register(scom.MethodBlackChain, BlackChain)
	s.Register(scom.MethodWhiteChain, WhiteChain)
	s.Register(scom.MethodCheckDone, CheckDone)
	s.Register(scom.MethodReplenish, Replenish)
}

func GetChainHandler(router uint64) (scom.ChainHandler, error) {
	switch router {
	case utils.NO_PROOF_ROUTER:
		return no_proof.NewNoProofHandler(), nil
	case utils.ETH_COMMON_ROUTER:
		return eth_common.NewHandler(), nil
	case utils.BSC_ROUTER:
		return bsc.NewHandler(), nil
	case utils.ETH_ROUTER:
		return eth.NewETHHandler(), nil
	case utils.HECO_ROUTER:
		return heco.NewHecoHandler(), nil
	case utils.MSC_ROUTER:
		return msc.NewHandler(), nil
	case utils.OKEX_ROUTER:
		return okex.NewHandler(), nil
	case utils.QUORUM_ROUTER:
		return quorum.NewQuorumHandler(), nil
	case utils.COSMOS_ROUTER:
		return cosmos.NewCosmosHandler(), nil
	case utils.ZILLIQA_ROUTER:
		return zilliqa.NewHandler(), nil
	default:
		return nil, fmt.Errorf("not a supported router:%d", router)
	}
}

func Name(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(scom.ABI, scom.MethodContractName, contractName)
}

func CheckDone(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &scom.CheckDoneParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodCheckDone, params, ctx.Payload); err != nil {
		return nil, err
	}
	err := scom.CheckDoneTx(s, params.CrossChainID, params.SourceChainID)
	return utils.PackOutputs(scom.ABI, scom.MethodCheckDone, err == scom.ErrTxAlreadyImported)
}

func ImportOuterTransfer(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &scom.EntranceParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodImportOuterTransfer, params, ctx.Payload); err != nil {
		return nil, err
	}

	srcChainID := params.SourceChainID
	blacked, err := CheckIfChainBlacked(s, srcChainID)
	if err != nil {
		return nil, fmt.Errorf("ImportExTransfer, CheckIfChainBlacked err: %v", err)
	}
	if blacked {
		return nil, fmt.Errorf("ImportExTransfer, source chain is blacked")
	}

	srcChain, err := side_chain_manager.GetSideChain(s, srcChainID)
	if err != nil {
		return nil, fmt.Errorf("ImportExTransfer, side_chain_manager.GetSideChain err: %v", err)
	} else if srcChain == nil {
		return nil, fmt.Errorf("ImportExTransfer, side chain %d is not registered", srcChainID)
	}

	handler, err := GetChainHandler(srcChain.Router)
	if err != nil {
		return nil, err
	}
	if handler == nil {
		return nil, fmt.Errorf("ImportExTransfer, handler for side chain %d is not exist", srcChainID)
	}

	txParam, err := handler.MakeDepositProposal(s)
	if err != nil {
		return nil, err
	}

	//check target chain
	dstChainID := txParam.ToChainID
	if blacked, err = CheckIfChainBlacked(s, dstChainID); err != nil {
		return nil, fmt.Errorf("ImportExTransfer, CheckIfChainBlacked error: %v", err)
	}
	if blacked {
		return nil, fmt.Errorf("ImportExTransfer, target chain is blacked")
	}

	dstChain, err := side_chain_manager.GetSideChain(s, dstChainID)
	if err != nil {
		return nil, fmt.Errorf("ImportExTransfer, side_chain_manager.GetSideChain error: %v", err)
	}
	if dstChain == nil {
		return nil, fmt.Errorf("ImportExTransfer, side chain %d is not registered", dstChainID)
	}
	if dstChain.Router == utils.BTC_ROUTER {
		return nil, fmt.Errorf("btc is not supported")
	}

	//NOTE, you need to store the tx in this
	if err := scom.MakeTransaction(s, txParam, srcChainID); err != nil {
		return nil, err
	}

	return utils.PackOutputs(scom.ABI, scom.MethodImportOuterTransfer, true)
}

func BlackChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &scom.BlackChainParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodBlackChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	// Get current epoch operator
	ok, err := node_manager.CheckConsensusSigns(s, scom.MethodBlackChain, utils.GetUint64Bytes(params.ChainID), s.ContractRef().MsgSender())
	if err != nil {
		return nil, fmt.Errorf("BlackChain, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.PackOutputs(scom.ABI, scom.MethodBlackChain, true)
	}

	PutBlackChain(s, params.ChainID)
	return utils.PackOutputs(scom.ABI, scom.MethodBlackChain, true)
}

func WhiteChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &scom.BlackChainParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodWhiteChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	// Get current epoch operator
	ok, err := node_manager.CheckConsensusSigns(s, scom.MethodWhiteChain, ctx.Payload, s.ContractRef().MsgSender())
	if err != nil {
		return nil, fmt.Errorf("WhiteChain, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.PackOutputs(scom.ABI, scom.MethodWhiteChain, true)
	}

	RemoveBlackChain(s, params.ChainID)
	return utils.PackOutputs(scom.ABI, scom.MethodWhiteChain, true)
}

func Replenish(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &scom.ReplenishParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodReplenish, params, ctx.Payload); err != nil {
		return nil, fmt.Errorf("Replenish, unpack params error: %s", err)
	}

	err := scom.NotifyReplenish(s, params.TxHashes, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("Replenish, NotifyReplenish error: %s", err)
	}
	return utils.BYTE_TRUE, nil
}
