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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	SKP_CURRENT_TX_INDEX = "st_tx_idx"
	SKP_TX_PROOF         = "st_tx_pf"
	SKP_TX_CONTENT       = "st_tx_dat"
	SKP_EPOCH            = "st_ep"
)

func getCrossTxIndex(s *native.NativeContract) uint64 {
	blob, _ := s.GetCacheDB().Get(currentTxIndexKey())
	if blob == nil {
		return 0
	}
	return utils.GetBytesUint64(blob)
}

func storeCrossTxIndex(s *native.NativeContract, index uint64) {
	s.GetCacheDB().Put(currentTxIndexKey(), utils.GetUint64Bytes(index))
}

func getCrossTxContent(s *native.NativeContract, hash common.Hash) (*CrossTx, error) {
	blob, err := s.GetCacheDB().Get(crossTxDataKey(hash))
	if err != nil {
		return nil, err
	}
	if blob == nil {
		return nil, nil
	}

	tx, err := DecodeCrossTx(blob)
	if err := rlp.DecodeBytes(blob, tx); err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil
	}

	if tx.Hash() != hash {
		return nil, fmt.Errorf("expect crossTx hash %s, got %s", hash.Hex(), tx.Hash().Hex())
	}
	return tx, nil
}

func storeCrossTxContent(s *native.NativeContract, tx *CrossTx) error {
	blob, err := EncodeCrossTx(tx)
	if err != nil {
		return err
	}

	s.GetCacheDB().Put(crossTxDataKey(tx.Hash()), blob)
	return nil
}

func getCrossTxProof(s *native.NativeContract, index uint64) (common.Hash, error) {
	blob, err := s.GetCacheDB().Get(crossTxProofKey(index))
	if err != nil {
		return common.EmptyHash, err
	}
	if blob == nil {
		return common.EmptyHash, nil
	}
	return common.BytesToHash(blob), nil
}

func storeCrossTxProof(s *native.NativeContract, index uint64, proof common.Hash) {
	s.GetCacheDB().Put(crossTxProofKey(index), proof[:])
}

func getEpoch(s *native.NativeContract) ([]byte, *nm.EpochInfo, error) {
	blob, err := s.GetCacheDB().Get(epochKey())
	if err != nil {
		return nil, nil, err
	}
	if blob == nil {
		return nil, nil, nil
	}
	epoch, err := DecodeEpoch(blob)
	if err != nil {
		return nil, nil, err
	}
	return blob, epoch, nil
}

func storeEpoch(s *native.NativeContract, epoch []byte) {
	s.GetCacheDB().Put(epochKey(), epoch)
}

// ====================================================================
//
// storage keys
//
// ====================================================================

func currentTxIndexKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_CURRENT_TX_INDEX))
}

func crossTxProofKey(index uint64) []byte {
	return utils.ConcatKey(this, []byte(SKP_TX_PROOF), utils.GetUint64Bytes(index))
}

func crossTxDataKey(hash common.Hash) []byte {
	return utils.ConcatKey(this, []byte(SKP_TX_CONTENT), hash[:])
}

func epochKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_EPOCH))
}
