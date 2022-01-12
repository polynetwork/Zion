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

package zion

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func storeEpoch(s *native.NativeContract, chainID, height uint64, validators []common.Address) error {
	storeHeight(s, chainID, height)
	return storeValSet(s, chainID, validators)
}

func getEpoch(s *native.NativeContract, chainID uint64) (height uint64, valset []common.Address, err error) {
	if height, err = getHeight(s, chainID); err != nil {
		return
	}
	valset, err = getValSet(s, chainID)
	return
}

func storeValSet(s *native.NativeContract, chainID uint64, validators []common.Address) error {
	blob, err := rlp.EncodeToBytes([]interface{}{validators})
	if err != nil {
		return err
	}

	key := valsetKey(chainID)
	s.GetCacheDB().Put(key, blob)
	return nil
}

func getValSet(s *native.NativeContract, chainID uint64) ([]common.Address, error) {
	key := valsetKey(chainID)
	blob, err := s.GetCacheDB().Get(key)
	if err != nil {
		return nil, err
	}
	if blob == nil {
		return nil, fmt.Errorf("storage data is nil")
	}

	var valset struct {
		List []common.Address
	}
	if err := rlp.DecodeBytes(blob, &valset); err != nil {
		return nil, err
	}
	return valset.List, nil
}

func storeHeight(s *native.NativeContract, chainID uint64, height uint64) {
	key := heightKey(chainID)
	s.GetCacheDB().Put(key, utils.GetUint64Bytes(height))
}

func getHeight(s *native.NativeContract, chainID uint64) (uint64, error) {
	key := heightKey(chainID)
	blob, err := s.GetCacheDB().Get(key)
	if err != nil {
		return 0, err
	}
	return utils.GetBytesUint64(blob), nil
}

func isGenesisStored(s *native.NativeContract, chainID uint64) bool {
	blob, err := scom.GetGenesisHeader(s, chainID)
	if blob != nil && len(blob) > 0 && err == nil {
		return true
	}
	return false
}

func storeGenesis(s *native.NativeContract, chainID uint64, genesisHeader *types.Header) error {
	blob, err := genesisHeader.MarshalJSON()
	if err != nil {
		return err
	}

	scom.SetGenesisHeader(s, chainID, blob)
	return nil
}

func getGenesisHeader(s *native.NativeContract, chainID uint64) (*types.Header, error) {
	blob, err := scom.GetGenesisHeader(s, chainID)
	if err != nil {
		return nil, err
	}
	header := new(types.Header)
	if err := header.UnmarshalJSON(blob); err != nil {
		return nil, err
	}
	return header, nil
}

////////////////////////////////////////////////////////////////////////////////////
//
// emit event logs
//
////////////////////////////////////////////////////////////////////////////////////

func emitEpochChangeEvent(s *native.NativeContract, chainID, height uint64, hash common.Hash) {
	scom.NotifyPutHeader(s, chainID, height, hash.Hex())
}

////////////////////////////////////////////////////////////////////////////////////
//
// storage keys
//
////////////////////////////////////////////////////////////////////////////////////
func valsetKey(chainID uint64) []byte {
	return utils.ConcatKey(utils.HeaderSyncContractAddress, []byte(scom.CONSENSUS_PEER), utils.GetUint64Bytes(chainID))
}

func heightKey(chainID uint64) []byte {
	return utils.ConcatKey(utils.HeaderSyncContractAddress, []byte(scom.CONSENSUS_PEER_BLOCK_HEIGHT), utils.GetUint64Bytes(chainID))
}
