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
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// todo: issue, add field `len` for stat addressList and hashList

const (
	StartEpoch uint64 = 1 // epoch started from 1, NOT 0!
)

// storage key prefix
const (
	SKP_EPOCH     = "st_epoch"
	SKP_PROOF     = "st_proof"
	SKP_PROPOSAL  = "st_proposal"
	SKP_VOTE      = "st_vote"
	SKP_CUR_EPOCH = "st_cur_epoch"
	SKP_SIGN      = "st_sign"
	SKP_SIGNER    = "st_signer"
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
func storeCurrentEpochHash(s *native.NativeContract, epochHash common.Hash) {
	key := curEpochKey()
	s.GetCacheDB().Put(key, epochHash.Bytes())
}

func getCurrentEpochHash(s *native.NativeContract) (common.Hash, error) {
	key := curEpochKey()
	value, err := s.GetCacheDB().Get(key)
	if err != nil {
		return common.EmptyHash, err
	}
	return common.BytesToHash(value), nil
}

// proof
func storeEpochProof(s *native.NativeContract, epochID uint64, epochHash common.Hash) {
	key := epochProofKey(EpochProofHash(epochID))
	s.GetCacheDB().Put(key, epochHash.Bytes())
}

func getEpochProof(s *native.NativeContract, epochID uint64) (common.Hash, error) {
	key := epochProofKey(EpochProofHash(epochID))
	value, err := s.GetCacheDB().Get(key)
	if err != nil {
		return common.EmptyHash, nil
	}
	return common.BytesToHash(value), nil
}

var EpochProofDigest = common.HexToHash("e4bf3526f07c80af3a5de1411dd34471c71bdd5d04eedbfa1040da2c96802041")

func EpochProofHash(epochID uint64) common.Hash {
	enc := EpochProofDigest.Bytes()
	enc = append(enc, utils.GetUint64Bytes(epochID)...)
	return crypto.Keccak256Hash(enc)
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
		if err.Error() == "EOF" {
			return []common.Address{}, nil
		} else {
			return nil, err
		}
	}

	var data *AddressList
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		return nil, err
	}
	return data.List, nil
}

func getVoteSize(s *native.NativeContract, epochHash common.Hash) int {
	votes, err := getVotes(s, epochHash)
	if err != nil {
		return 0
	}
	return len(votes)
}

// signatures
func storeSign(s *native.NativeContract, sign *ConsensusSign) error {
	key := signKey(sign.Hash())
	value, err := rlp.EncodeToBytes(sign)
	if err != nil {
		return err
	}
	s.GetCacheDB().Put(key, value)
	return nil
}

func delSign(s *native.NativeContract, hash common.Hash) {
	key := signKey(hash)
	s.GetCacheDB().Delete(key)
}

func getSign(s *native.NativeContract, hash common.Hash) (*ConsensusSign, error) {
	key := signKey(hash)
	value, err := s.GetCacheDB().Get(key)
	if err != nil {
		if err.Error() == "EOF" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	var sign *ConsensusSign
	if err := rlp.DecodeBytes(value, &sign); err != nil {
		return nil, err
	}
	return sign, nil
}

func storeSigner(s *native.NativeContract, hash common.Hash, signer common.Address) error {
	data, err := getSigners(s, hash)
	if err != nil {
		return err
	}
	data = append(data, signer)
	list := &AddressList{List: data}

	key := signerKey(hash)
	value, err := rlp.EncodeToBytes(list)
	if err != nil {
		return err
	}
	s.GetCacheDB().Put(key, value)

	return nil
}

func findSigner(s *native.NativeContract, hash common.Hash, signer common.Address) bool {
	list, err := getSigners(s, hash)
	if err != nil {
		return false
	}
	for _, v := range list {
		if v == signer {
			return true
		}
	}
	return false
}

func getSigners(s *native.NativeContract, hash common.Hash) ([]common.Address, error) {
	key := signerKey(hash)
	value, err := s.GetCacheDB().Get(key)
	if err != nil {
		if err.Error() == "EOF" {
			return []common.Address{}, nil
		} else {
			return nil, err
		}
	}

	var list *AddressList
	if err := rlp.DecodeBytes(value, &list); err != nil {
		return nil, err
	}
	return list.List, nil
}

func getSignerSize(s *native.NativeContract, hash common.Hash) int {
	list, err := getSigners(s, hash)
	if err != nil {
		return 0
	}
	return len(list)
}

func clearSigner(s *native.NativeContract, hash common.Hash) {
	key := signerKey(hash)
	s.GetCacheDB().Delete(key)
}

// keys
func epochKey(epochHash common.Hash) []byte {
	return utils.ConcatKey(this, []byte(SKP_EPOCH), epochHash.Bytes())
}

func epochProofKey(proofHashKey common.Hash) []byte {
	return utils.ConcatKey(this, []byte(SKP_PROOF), proofHashKey.Bytes())
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

func signKey(hash common.Hash) []byte {
	return utils.ConcatKey(this, []byte(SKP_SIGN), hash.Bytes())
}

func signerKey(hash common.Hash) []byte {
	return utils.ConcatKey(this, []byte(SKP_SIGNER), hash.Bytes())
}
