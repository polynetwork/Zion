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
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

func GetGenesisHeader(s *native.NativeContract, chainID uint64) ([]byte, error) {
	return s.GetCacheDB().Get(genesisHeaderKey(chainID))
}

func SetGenesisHeader(s *native.NativeContract, chainID uint64, raw []byte) {
	s.GetCacheDB().Put(genesisHeaderKey(chainID), raw)
}

func GetHeaderIndex(s *native.NativeContract, chainID uint64, hash []byte) ([]byte, error) {
	return s.GetCacheDB().Get(headerIndexKey(chainID, hash))
}

func SetHeaderIndex(s *native.NativeContract, chainID uint64, hash []byte, raw []byte) {
	s.GetCacheDB().Put(headerIndexKey(chainID, hash), raw)
}

func GetMainChain(s *native.NativeContract, chainID, blockNum uint64) ([]byte, error) {
	return s.GetCacheDB().Get(mainChainKey(chainID, blockNum))
}

func SetMainChain(s *native.NativeContract, chainID, blockNum uint64, raw []byte) {
	s.GetCacheDB().Put(mainChainKey(chainID, blockNum), raw)
}

func DelMainChain(s *native.NativeContract, chainID, blockNum uint64) {
	s.GetCacheDB().Delete(mainChainKey(chainID, blockNum))
}

func GetCurrentHeight(s *native.NativeContract, chainID uint64) ([]byte, error) {
	return s.GetCacheDB().Get(currentHeightKey(chainID))
}

func SetCurrentHeight(s *native.NativeContract, chainID uint64, raw []byte) {
	s.GetCacheDB().Put(currentHeightKey(chainID), raw)
}

////////////////////////////////////////////////////////////////////////////////////////////////
//
// storage keys
//
////////////////////////////////////////////////////////////////////////////////////////////////

var this = utils.HeaderSyncContractAddress

func genesisHeaderKey(chainID uint64) []byte {
	return utils.ConcatKey(this, []byte(GENESIS_HEADER), utils.GetUint64Bytes(chainID))
}

func headerIndexKey(chainID uint64, blockHash []byte) []byte {
	return utils.ConcatKey(this, []byte(HEADER_INDEX), utils.GetUint64Bytes(chainID), blockHash)
}

func mainChainKey(chainID, blockNum uint64) []byte {
	return utils.ConcatKey(this, []byte(MAIN_CHAIN), utils.GetUint64Bytes(chainID), utils.GetUint64Bytes(blockNum))
}

func currentHeightKey(chainID uint64) []byte {
	return utils.ConcatKey(this, []byte(CURRENT_HEADER_HEIGHT), utils.GetUint64Bytes(chainID))
}
