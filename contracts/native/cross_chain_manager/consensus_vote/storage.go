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

package consensus_vote

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/rlp"
)

var ErrEof = errors.New("EOF")

const (
	SKP_VOTE_MESSAGE = "st_vote_message"
	SKP_SIGNER_MAP   = "st_signer_map"
)

func getVoteMessage(s *native.NativeContract, hash common.Hash) (*VoteMessage, error) {
	key := voteMessageKey(hash)
	value, err := get(s, key)
	if err != nil {
		return nil, err
	}
	var msg *VoteMessage
	if err := rlp.DecodeBytes(value, &msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func storeVoteMessage(s *native.NativeContract, sign *VoteMessage) error {
	key := voteMessageKey(sign.Hash())
	value, err := rlp.EncodeToBytes(sign)
	if err != nil {
		return err
	}
	set(s, key, value)
	return nil
}

func storeSignerAndCheckQuorum(s *native.NativeContract, hash common.Hash, signer common.Address, quorum int) (bool, error) {
	height := s.ContractRef().BlockHeight().Uint64()
	data, err := getSignerMap(s, hash)
	if err != nil {
		if err.Error() == ErrEof.Error() {
			data = &SignerMap{
				StartHeight: height,
				SignerMap: make(map[common.Address]*SignerInfo),
			}
		} else {
			return false, err
		}
	}
	data.SignerMap[signer] = &SignerInfo{height}

	flag := false
	//check quorum and store quorum height
	if data.EndHeight == 0 {
		size := len(data.SignerMap)
		if size >= quorum {
			data.EndHeight = height
			flag = true
		}
	}

	//store signer map
	key := signerMapKey(hash)
	value, err := rlp.EncodeToBytes(data)
	if err != nil {
		return false, err
	}
	set(s, key, value)

	return flag, nil
}

func findSigner(s *native.NativeContract, hash common.Hash, signer common.Address) bool {
	signerMap, err := getSignerMap(s, hash)
	if err != nil {
		return false
	}
	_, ok := signerMap.SignerMap[signer]
	if ok {
		return true
	}
	return false
}

func getSignerMap(s *native.NativeContract, hash common.Hash) (*SignerMap, error) {
	key := signerMapKey(hash)
	value, err := get(s, key)
	if err != nil {
		return nil, err
	}

	var signerMap *SignerMap
	if err := rlp.DecodeBytes(value, &signerMap); err != nil {
		return nil, err
	}
	return signerMap, nil
}

func getSignerSize(s *native.NativeContract, hash common.Hash) int {
	signerMap, err := getSignerMap(s, hash)
	if err != nil {
		return 0
	}
	return len(signerMap.SignerMap)
}

// ====================================================================
//
// storage basic operations
//
// ====================================================================

func get(s *native.NativeContract, key []byte) ([]byte, error) {
	return customGet(s.GetCacheDB(), key)
}

func set(s *native.NativeContract, key, value []byte) {
	customSet(s.GetCacheDB(), key, value)
}

func del(s *native.NativeContract, key []byte) {
	customDel(s.GetCacheDB(), key)
}

func customGet(db *state.CacheDB, key []byte) ([]byte, error) {
	value, err := db.Get(key)
	if err != nil {
		return nil, err
	} else if value == nil || len(value) == 0 {
		return nil, ErrEof
	} else {
		return value, nil
	}
}

func customSet(db *state.CacheDB, key, value []byte) {
	db.Put(key, value)
}

func customDel(db *state.CacheDB, key []byte) {
	db.Delete(key)
}

// ====================================================================
//
// storage keys
//
// ====================================================================

func voteMessageKey(hash common.Hash) []byte {
	return utils.ConcatKey(utils.CrossChainManagerContractAddress, []byte(SKP_VOTE_MESSAGE), hash.Bytes())
}

func signerMapKey(hash common.Hash) []byte {
	return utils.ConcatKey(utils.CrossChainManagerContractAddress, []byte(SKP_SIGNER_MAP), hash.Bytes())
}
