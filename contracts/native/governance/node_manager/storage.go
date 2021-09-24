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

package node_manager

import (
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/core/state"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
)

// todo: issue, add field `len` for stat addressList and hashList

// storage key prefix
const (
	SKP_EPOCH     = "st_epoch"
	SKP_PROPOSAL  = "st_proposal"
	SKP_VOTE      = "st_vote"
	SKP_CUR_EPOCH = "st_cur_epoch"
)

// epoch storage
func storeEpoch(s *native.NativeContract, epoch *EpochInfo) error {
	return setEpoch(s.GetCacheDB(), epoch)
}

func getEpoch(s *native.NativeContract, epochHash common.Hash) (*EpochInfo, error) {
	key := epochKey(epochHash)
	enc, err := s.GetCacheDB().Get(key)
	if err != nil {
		return nil, err
	}

	epoch := new(EpochInfo)
	if err := rlp.DecodeBytes(enc, epoch); err != nil {
		return nil, err
	}

	return epoch, nil
}

func delEpoch(s *native.NativeContract, epochHash common.Hash) {
	key := epochKey(epochHash)
	s.GetCacheDB().Delete(key)
}

func setEpoch(s *state.CacheDB, epoch *EpochInfo) error {
	hash := epoch.Hash()
	key := epochKey(hash)

	value, err := rlp.EncodeToBytes(epoch)
	if err != nil {
		return err
	}

	s.Put(key, value)
	return nil
}

// current epoch
func storeCurrentEpoch(s *native.NativeContract, epochHash common.Hash) {
	key := curEpochKey()
	s.GetCacheDB().Put(key, epochHash.Bytes())
}

func getCurrentEpoch(s *native.NativeContract) (common.Hash, error) {
	key := curEpochKey()
	value, err := s.GetCacheDB().Get(key)
	if err != nil {
		return common.EmptyHash, err
	}
	return common.BytesToHash(value), nil
}

// proposal storage
func storeProposal(s *native.NativeContract, epochHash common.Hash) error {
	list, err := getProposals(s)
	if err != nil {
		if err.Error() == "EOF" {
			list = make([]common.Hash, 0)
		} else {
			return err
		}
	}
	list = append(list, epochHash)
	return setProposals(s, list)
}

func firstProposal(s *native.NativeContract) (common.Hash, error) {
	list, err := getProposals(s)
	if err != nil {
		return common.EmptyHash, nil
	}
	if list == nil || len(list) == 0 {
		return common.EmptyHash, nil
	}
	return list[0], nil
}

func findProposal(s *native.NativeContract, epochHash common.Hash) bool {
	list, err := getProposals(s)
	if err != nil {
		return false
	}
	for _, hash := range list {
		if hash == epochHash {
			return true
		}
	}
	return false
}

func proposalsNum(s *native.NativeContract) int {
	list, err := getProposals(s)
	if err != nil || list == nil {
		return 0
	}
	return len(list)
}

func delProposal(s *native.NativeContract, epochHash common.Hash) error {
	list, err := getProposals(s)
	if err != nil {
		return err
	}
	if list == nil || len(list) == 0 {
		return fmt.Errorf("proposal %s not exist", epochHash.Hex())
	}
	dst := make([]common.Hash, 0)
	for _, hash := range list {
		if hash == epochHash {
			continue
		} else {
			dst = append(dst, hash)
		}
	}

	return setProposals(s, dst)
}

func setProposals(s *native.NativeContract, list []common.Hash) error {
	value, err := rlp.EncodeToBytes(&HashList{List: list})
	if err != nil {
		return err
	}

	key := proposalsKey()
	s.GetCacheDB().Put(key, value)
	return nil
}

func getProposals(s *native.NativeContract) ([]common.Hash, error) {
	key := proposalsKey()
	enc, err := s.GetCacheDB().Get(key)
	if err != nil {
		return nil, err
	}

	var data *HashList
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		return nil, err
	}
	return data.List, nil
}

// vote storage
func storeVote(s *native.NativeContract, epochHash common.Hash, voter common.Address) error {
	list, err := getVotes(s, epochHash)
	if err != nil {
		if err.Error() == "EOF" {
			list = make([]common.Address, 0)
		} else {
			return err
		}
	}
	list = append(list, voter)
	return setVotes(s, epochHash, list)
}

func voteSize(s *native.NativeContract, epochHash common.Hash) int {
	list, err := getVotes(s, epochHash)
	if err != nil {
		return 0
	}
	return len(list)
}

func findVote(s *native.NativeContract, epochHash common.Hash, voter common.Address) bool {
	list, err := getVotes(s, epochHash)
	if err != nil {
		return false
	}
	for _, addr := range list {
		if addr == voter {
			return true
		}
	}
	return false
}

func deleteVote(s *native.NativeContract, epochHash common.Hash, voter common.Address) error {
	list, err := getVotes(s, epochHash)
	if err != nil {
		return err
	}
	dst := make([]common.Address, 0)
	for _, addr := range list {
		if addr == voter {
			continue
		} else {
			dst = append(dst, addr)
		}
	}
	return setVotes(s, epochHash, dst)
}

func clearVotes(s *native.NativeContract, epochHash common.Hash) {
	key := voteKey(epochHash)
	s.GetCacheDB().Put(key, common.EmptyHash.Bytes())
}

func setVotes(s *native.NativeContract, epochHash common.Hash, list []common.Address) error {
	key := voteKey(epochHash)

	value, err := rlp.EncodeToBytes(&AddressList{List: list})
	if err != nil {
		return err
	}

	s.GetCacheDB().Put(key, value)
	return nil
}

func getVotes(s *native.NativeContract, epochHash common.Hash) ([]common.Address, error) {
	key := voteKey(epochHash)
	enc, err := s.GetCacheDB().Get(key)
	if err != nil {
		return nil, err
	}

	var data *AddressList
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		return nil, err
	}
	return data.List, nil
}

// keys
func epochKey(epochHash common.Hash) []byte {
	return utils.ConcatKey(this, []byte(SKP_EPOCH), epochHash.Bytes())
}

func curEpochKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_CUR_EPOCH), []byte("1"))
}

func proposalsKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_PROPOSAL), []byte("2"))
}

func voteKey(epochHash common.Hash) []byte {
	return utils.ConcatKey(this, []byte(SKP_VOTE), epochHash.Bytes())
}
