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
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const (
	SKP_TX_HASH   = "st_tx_hash"
	SKP_TX_INDEX  = "st_tx_index"
	SKP_TX_PARAMS = "st_tx_params"
)

func getTxIndex(s *native.NativeContract) *big.Int {
	key := txIndexKey()
	blob, _ := s.GetCacheDB().Get(key)
	if blob == nil {
		return common.Big0
	}
	return new(big.Int).SetBytes(blob)
}

func storeTxIndex(s *native.NativeContract, txIndex *big.Int) {
	s.GetCacheDB().Put(txIndexKey(), scom.Uint256ToBytes(txIndex))
}

func getTxProof(s *native.NativeContract, paramTxHash []byte) common.Hash {
	blob, _ := s.GetCacheDB().Get(txHashKey(paramTxHash))
	if blob == nil {
		return common.EmptyHash
	}
	return common.BytesToHash(blob)
}

func storeTxProof(s *native.NativeContract, paramTxHash []byte, proof common.Hash) {
	s.GetCacheDB().Put(txHashKey(paramTxHash), proof[:])
}

func getTxParams(s *native.NativeContract, paramTxHash []byte) ([]byte, error) {
	return s.GetCacheDB().Get(txHashKey(paramTxHash))
}

func storeTxParams(s *native.NativeContract, paramTxHash []byte, params []byte) {
	key := txParamsKey(paramTxHash)
	s.GetCacheDB().Put(key, params)
}

// ====================================================================
//
// storage keys
//
// ====================================================================
func txHashKey(paramTxHash []byte) []byte {
	return utils.ConcatKey(this, []byte(SKP_TX_HASH), paramTxHash)
}

func txIndexKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_TX_INDEX))
}

func txParamsKey(hash []byte) []byte {
	return utils.ConcatKey(this, []byte(SKP_TX_PARAMS), hash)
}
