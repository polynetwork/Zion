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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const (
	SKP_TX_HASH      = "st_tx_hash"
	SKP_TX_INDEX     = "st_tx_index"
	SKP_TX_PARAMS    = "st_tx_params"
	SKP_TOTAL_AMOUNT = "st_amt"
)

func getNextTxIndex(s *native.NativeContract) (*big.Int, error) {
	lastTxIndex := getTxIndex(s)
	storeTxIndex(s, new(big.Int).Add(lastTxIndex, common.Big1))
	txIndex := getTxIndex(s)
	if txIndex.Cmp(common.Big0) <= 0 {
		return nil, fmt.Errorf("txIndex invalid")
	}
	return txIndex, nil
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
	s.GetCacheDB().Put(txIndexKey(), scom.Uint256ToBytes(txIndex))
}

func getTotalAmount(s *native.NativeContract, sideChainID uint64) *big.Int {
	key := totalAmountKey(sideChainID)
	blob, _ := s.GetCacheDB().Get(key)
	if blob == nil {
		return common.Big0
	}
	return new(big.Int).SetBytes(blob)
}

func addTotalAmount(s *native.NativeContract, sideChainID uint64, amount *big.Int) {
	total := getTotalAmount(s, sideChainID)
	total = new(big.Int).Add(total, amount)
	storeTotalAmount(s, sideChainID, total)
}

func subTotalAmount(s *native.NativeContract, sideChainID uint64, amount *big.Int) error {
	total := getTotalAmount(s, sideChainID)
	if total.Cmp(amount) < 0 {
		return fmt.Errorf("side chain %d only locked %v native token", sideChainID, total)
	} else {
		total = new(big.Int).Sub(total, amount)
	}
	storeTotalAmount(s, sideChainID, total)
	return nil
}

func storeTotalAmount(s *native.NativeContract, sideChainID uint64, amount *big.Int) {
	s.GetCacheDB().Put(totalAmountKey(sideChainID), amount.Bytes())
}

// ====================================================================
//
// storage keys
//
// ====================================================================
func txIndexKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_TX_INDEX))
}

func totalAmountKey(chainID uint64) []byte {
	return utils.ConcatKey(this, []byte(SKP_TOTAL_AMOUNT), utils.Uint64Bytes(chainID))
}
