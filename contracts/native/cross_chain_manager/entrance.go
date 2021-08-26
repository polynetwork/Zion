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
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/bsc"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/btc"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/eth"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/heco"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/msc"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/polygon"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/quorum"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	polycomm "github.com/polynetwork/poly/common"
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
	}
)

func InitCrossChainManager() {
	native.Contracts[this] = RegisterCrossChainManagerContract
}

func RegisterCrossChainManagerContract(s *native.NativeContract) {
	s.Prepare(scom.ABI, gasTable)

	s.Register(scom.MethodContractName, Name)
	s.Register(scom.MethodImportOuterTransfer, ImportOuterTransfer)
	s.Register(scom.MethodMultiSign, MultiSign)
	s.Register(scom.MethodBlackChain, BlackChain)
	s.Register(scom.MethodWhiteChain, WhiteChain)
}

func GetChainHandler(router uint64) (scom.ChainHandler, error) {
	switch router {
	case utils.BTC_ROUTER:
		return btc.NewBTCHandler(), nil
	case utils.BSC_ROUTER:
		return bsc.NewHandler(), nil
	case utils.ETH_ROUTER:
		return eth.NewETHHandler(), nil
	case utils.HECO_ROUTER:
		return heco.NewHecoHandler(), nil
	case utils.MSC_ROUTER:
		return msc.NewHandler(), nil
	case utils.QUORUM_ROUTER:
		return quorum.NewQuorumHandler(), nil
	case utils.POLYGON_BOR_ROUTER:
		return polygon.NewHandler(), nil
	default:
		return nil, fmt.Errorf("not a supported router:%d", router)
	}
}

func Name(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(scom.ABI, scom.MethodContractName, contractName)
}

func ImportOuterTransfer(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &scom.EntranceParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodImportOuterTransfer, params, ctx.Payload); err != nil {
		return nil, err
	}

	chainID := params.SourceChainID
	blacked, err := CheckIfChainBlacked(native, chainID)
	if err != nil {
		return nil, fmt.Errorf("ImportExTransfer, CheckIfChainBlacked error: %v", err)
	}
	if blacked {
		return nil, fmt.Errorf("ImportExTransfer, source chain is blacked")
	}

	//check if chainid exist
	sideChain, err := side_chain_manager.GetSideChain(native, chainID)
	if err != nil {
		return nil, fmt.Errorf("ImportExTransfer, side_chain_manager.GetSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("ImportExTransfer, side chain %d is not registered", chainID)
	}

	handler, err := GetChainHandler(sideChain.Router)
	if err != nil {
		return nil, err
	}
	//1. verify tx
	txParam, err := handler.MakeDepositProposal(native)
	if err != nil {
		return nil, err
	}

	//2. make target chain tx
	targetid := txParam.ToChainID
	blacked, err = CheckIfChainBlacked(native, targetid)
	if err != nil {
		return nil, fmt.Errorf("ImportExTransfer, CheckIfChainBlacked error: %v", err)
	}
	if blacked {
		return nil, fmt.Errorf("ImportExTransfer, target chain is blacked")
	}

	//check if chainid exist
	sideChain, err = side_chain_manager.GetSideChain(native, targetid)
	if err != nil {
		return nil, fmt.Errorf("ImportExTransfer, side_chain_manager.GetSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("ImportExTransfer, side chain %d is not registered", targetid)
	}
	if sideChain.Router == utils.BTC_ROUTER {
		err := btc.NewBTCHandler().MakeTransaction(native, txParam, chainID)
		if err != nil {
			return nil, err
		}
		return utils.PackOutputs(scom.ABI, scom.MethodImportOuterTransfer, true)
	}

	//NOTE, you need to store the tx in this
	err = MakeTransaction(native, txParam, chainID)
	if err != nil {
		return nil, err
	}

	return utils.PackOutputs(scom.ABI, scom.MethodImportOuterTransfer, true)
}

func MakeTransaction(service *native.NativeContract, params *scom.MakeTxParam, fromChainID uint64) error {

	txHash := service.ContractRef().TxHash()
	merkleValue := &scom.ToMerkleValue{
		TxHash:      txHash[:],
		FromChainID: fromChainID,
		MakeTxParam: params,
	}

	sink := polycomm.NewZeroCopySink(nil)
	merkleValue.Serialization(sink)
	err := PutRequest(service, merkleValue.TxHash, params.ToChainID, sink.Bytes())
	if err != nil {
		return fmt.Errorf("MakeTransaction, putRequest error:%s", err)
	}
	chainIDBytes := utils.GetUint64Bytes(params.ToChainID)
	key := hex.EncodeToString(utils.ConcatKey(utils.CrossChainManagerContractAddress, []byte(scom.REQUEST), chainIDBytes, merkleValue.TxHash))
	scom.NotifyMakeProof(service, hex.EncodeToString(sink.Bytes()), key)
	return nil
}

func PutRequest(native *native.NativeContract, txHash []byte, chainID uint64, request []byte) error {
	hash := crypto.Keccak256(request)
	contract := utils.CrossChainManagerContractAddress
	chainIDBytes := utils.GetUint64Bytes(chainID)
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(scom.REQUEST), chainIDBytes, txHash), hash)
	return nil
}

func MultiSign(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &scom.MultiSignParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodMultiSign, params, ctx.Payload); err != nil {
		return nil, err
	}

	handler := btc.NewBTCHandler()

	//1. multi sign
	err := handler.MultiSign(native, params)
	if err != nil {
		return nil, fmt.Errorf("MultiSign fail:%v", err)
	}

	return utils.PackOutputs(scom.ABI, scom.MethodMultiSign, true)
}

func BlackChain(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &scom.BlackChainParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodBlackChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	// Get current epoch operator
	ok, err := node_manager.CheckConsensusSigns(native, scom.MethodBlackChain, utils.GetUint64Bytes(params.ChainID), native.ContractRef().MsgSender())
	if err != nil {
		return nil, fmt.Errorf("BlackChain, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.PackOutputs(scom.ABI, scom.MethodBlackChain, true)
	}

	PutBlackChain(native, params.ChainID)
	return utils.PackOutputs(scom.ABI, scom.MethodBlackChain, true)
}

func WhiteChain(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &scom.BlackChainParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodWhiteChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	// Get current epoch operator
	ok, err := node_manager.CheckConsensusSigns(native, scom.MethodWhiteChain, ctx.Payload, native.ContractRef().MsgSender())
	if err != nil {
		return nil, fmt.Errorf("WhiteChain, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.PackOutputs(scom.ABI, scom.MethodWhiteChain, true)
	}

	RemoveBlackChain(native, params.ChainID)
	return utils.PackOutputs(scom.ABI, scom.MethodWhiteChain, true)
}
