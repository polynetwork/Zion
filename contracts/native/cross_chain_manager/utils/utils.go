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

package utils

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
	polycomm "github.com/polynetwork/poly/common"
)

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
