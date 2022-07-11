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

package common

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/core/state"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func MakeTransaction(service *native.NativeContract, params *MakeTxParam, fromChainID uint64) error {

	txHash := service.ContractRef().TxHash()
	merkleValue := &ToMerkleValue{
		TxHash:      txHash[:],
		FromChainID: fromChainID,
		MakeTxParam: params,
	}

	value, err := rlp.EncodeToBytes(merkleValue)
	if err != nil {
		return fmt.Errorf("MakeTransaction, rlp.EncodeToBytes merkle value error:%s", err)
	}
	err = PutRequest(service, merkleValue.TxHash, params.ToChainID, value)
	if err != nil {
		return fmt.Errorf("MakeTransaction, putRequest error:%s", err)
	}
	chainIDBytes := utils.GetUint64Bytes(params.ToChainID)
	key := state.Key2Slot(append([]byte(REQUEST), append(chainIDBytes, merkleValue.TxHash...)...)).String()
	if err := NotifyMakeProof(service, hex.EncodeToString(value), key); err != nil {
		return fmt.Errorf("MakeTransaction, NotifyMakeProof error:%s", err)
	}
	return nil
}

func PutRequest(native *native.NativeContract, txHash []byte, chainID uint64, request []byte) error {
	hash := crypto.Keccak256(request)
	contract := utils.CrossChainManagerContractAddress
	chainIDBytes := utils.GetUint64Bytes(chainID)
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(REQUEST), chainIDBytes, txHash), hash)
	return nil
}
