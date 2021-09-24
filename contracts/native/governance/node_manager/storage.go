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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
)

// storage key prefix
const (
	SKP_EPOCH    = "st_epoch"
	SKP_PROPOSAL = "st_proposal"
	SKP_VOTE     = "st_vote"
)

// epoch storage
func storeEpoch(s *native.NativeContract, epoch *EpochInfo) error {
	hash := epoch.Hash()
	key := epochKey(hash)

	value, err := rlp.EncodeToBytes(epoch)
	if err != nil {
		return err
	}

	s.GetCacheDB().Put(key, value)
	return nil
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

// proposal storage
func storeProposal(s *native.NativeContract, epochHash common.Hash) error {
	list, err := getProposals(s)
	if err != nil {
		return err
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

	key := proposalKey()
	s.GetCacheDB().Put(key, value)
	return nil
}

func getProposals(s *native.NativeContract) ([]common.Hash, error) {
	key := proposalKey()
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
		return err
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

func proposalKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_PROPOSAL), []byte("1"))
}

func voteKey(epochHash common.Hash) []byte {
	return utils.ConcatKey(this, []byte(SKP_VOTE), epochHash.Bytes())
}
