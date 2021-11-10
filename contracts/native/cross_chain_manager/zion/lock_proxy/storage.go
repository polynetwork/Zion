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
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const (
	SKP_PROXY     = "st_proxy"
	SKP_ASSET     = "st_asset"
	SKP_TX_HASH   = "st_tx_hash"
	SKP_TX_INDEX  = "st_tx_index"
	SKP_TX_PARAMS = "st_tx_params"
)

func getProxy(s *native.NativeContract, targetChainID uint64) ([]byte, error) {
	key := proxyKey(targetChainID)
	return s.GetCacheDB().Get(key)
}
func storeProxy(s *native.NativeContract, targetChainID uint64, proxyHash []byte) {
	key := proxyKey(targetChainID)
	s.GetCacheDB().Put(key, proxyHash)
}

func getAsset(s *native.NativeContract, fromAsset common.Address, targetChainID uint64) ([]byte, error) {
	key := assetKey(fromAsset, targetChainID)
	return s.GetCacheDB().Get(key)
}
func storeAsset(s *native.NativeContract, fromAsset common.Address, targetChainID uint64, toAssetHash []byte) {
	key := assetKey(fromAsset, targetChainID)
	s.GetCacheDB().Put(key, toAssetHash)
}

func getTxIndex(s *native.NativeContract) *big.Int {
	key := txIndexKey()
	blob, _ := s.GetCacheDB().Get(key)
	if blob == nil {
		return common.Big0
	}
	return new(big.Int).SetBytes(blob)
}

func storeTxIndex(s *native.NativeContract, txIndex *big.Int) {
	s.GetCacheDB().Put(txIndexKey(), txIndex.Bytes())
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

// storeTxParams store tx params and generate tx proof
func storeTxParams(s *native.NativeContract, hash common.Hash, params []byte) {
	key := txParamsKey(hash)
	s.GetCacheDB().Put(key, params)
}

// ====================================================================
//
// storage keys
//
// ====================================================================

func proxyKey(chainID uint64) []byte {
	return utils.ConcatKey(this, []byte(SKP_PROXY), utils.GetUint64Bytes(chainID))
}

func assetKey(fromAsset common.Address, chainID uint64) []byte {
	return utils.ConcatKey(this, []byte(SKP_ASSET), fromAsset[:], utils.GetUint64Bytes(chainID))
}

func txHashKey(paramTxHash []byte) []byte {
	return utils.ConcatKey(this, []byte(SKP_TX_HASH), paramTxHash)
}

func txIndexKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_TX_INDEX))
}

func txParamsKey(hash common.Hash) []byte {
	return utils.ConcatKey(this, []byte(SKP_TX_PARAMS), hash[:])
}
