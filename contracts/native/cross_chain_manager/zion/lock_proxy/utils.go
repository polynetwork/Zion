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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
	polycomm "github.com/polynetwork/poly/common"
)

func EncodeTxArgs(toAssetHash, toAddress []byte, amount *big.Int) []byte {
	sink := polycomm.NewZeroCopySink(nil)
	args := &scom.TxArgs{
		ToAssetHash: toAssetHash,
		ToAddress:   toAddress,
		Amount:      amount,
	}
	args.Serialization(sink)
	return sink.Bytes()
}

func DecodeTxArgs(payload []byte) (*scom.TxArgs, error) {
	source := polycomm.NewZeroCopySource(payload)
	args := new(scom.TxArgs)
	if err := args.Deserialization(source); err != nil {
		return nil, err
	}
	return args, nil
}

func EncodeMakeTxParams(tx common.Hash, txIndexID uint64, caller common.Address,
	toChainID uint64, toContract []byte, method string, args []byte) (
	*scom.MakeTxParam, []byte, common.Hash) {

	txParams := &scom.MakeTxParam{
		TxHash:              tx[:],
		CrossChainID:        utils.Uint64Bytes(txIndexID),
		FromContractAddress: caller[:],
		ToChainID:           toChainID,
		ToContractAddress:   toContract,
		Method:              method,
		Args:                args,
	}

	sink := polycomm.NewZeroCopySink(nil)
	txParams.Serialization(sink)
	txProof := crypto.Keccak256Hash(sink.Bytes())
	return txParams, sink.Bytes(), txProof
}

func DecodeMakeTxParams(blob []byte) (*scom.MakeTxParam, error) {
	source := polycomm.NewZeroCopySource(blob)
	txParams := new(scom.MakeTxParam)
	if err := txParams.Deserialization(source); err != nil {
		return nil, err
	}
	return txParams, nil
}
