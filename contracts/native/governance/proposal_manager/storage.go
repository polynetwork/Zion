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

package proposal_manager

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"sort"
)

var ErrEof = errors.New("EOF")

const (
	SKP_PROPOSAL_ID   = "st_proposal_id"
	SKP_PROPOSAL_LIST = "st_proposal_list"
)

func getProposalID(s *native.NativeContract) (*big.Int, error) {
	proposalID := new(big.Int)
	key := proposalIDKey()
	store, err := get(s, key)
	if err == ErrEof {
		return proposalID, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getProposalID, get store error: %v", err)
	}
	return new(big.Int).SetBytes(store), nil
}

func setProposalID(s *native.NativeContract, proposalID *big.Int) {
	key := proposalIDKey()
	set(s, key, proposalID.Bytes())
}

func getProposalList(s *native.NativeContract) (*ProposalList, error) {
	proposalList := &ProposalList{
		make([]*Proposal, 0),
	}
	key := proposalListKey()
	store, err := get(s, key)
	if err == ErrEof {
		return proposalList, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getProposalList, get store error: %v", err)
	}
	if err := rlp.DecodeBytes(store, proposalList); err != nil {
		return nil, fmt.Errorf("getProposalList, deserialize proposal list error: %v", err)
	}
	return proposalList, nil
}

func setProposalList(s *native.NativeContract, proposalList *ProposalList) error {
	// sort
	sort.SliceStable(proposalList.ProposalList, func(i, j int) bool {
		if proposalList.ProposalList[i].Status == Active {
			return true
		} else if proposalList.ProposalList[i].Stake.Cmp(proposalList.ProposalList[j].Stake) > 0 {
			return true
		}
		return false
	})

	key := proposalListKey()
	store, err := rlp.EncodeToBytes(proposalList)
	if err != nil {
		return fmt.Errorf("setProposalList, serialize proposalList error: %v", err)
	}
	set(s, key, store)
	return nil
}

func getActiveProposal(s *native.NativeContract) (*Proposal, bool, error) {
	flag := false
	proposalList, err := getProposalList(s)
	if err != nil {
		return nil, false, fmt.Errorf("getActiveProposal, getProposalList error: %v", err)
	}
	if len(proposalList.ProposalList) == 0 {
		return nil, false, fmt.Errorf("getActiveProposal, there is no proposal")
	}
	proposal := proposalList.ProposalList[0]
	if proposal.Status != Active || proposal.EndHeight.Cmp(s.ContractRef().BlockHeight()) > 0 {
		if len(proposalList.ProposalList) == 1 {
			return nil, false, fmt.Errorf("getActiveProposal, there is no active proposal")
		}
		proposal = proposalList.ProposalList[1]
		flag = true
	}
	return proposal, flag, nil
}

func setActiveProposal(s *native.NativeContract) error {
	proposalList, err := getProposalList(s)
	if err != nil {
		return fmt.Errorf("setActiveProposal, getProposalList error: %v", err)
	}
	if len(proposalList.ProposalList) == 0 {
		return nil
	}
	if len(proposalList.ProposalList) == 1 {
		proposalList.ProposalList = make([]*Proposal, 0)
	} else {
		oldActive := proposalList.ProposalList[0]
		proposalList.ProposalList = proposalList.ProposalList[1:]
		proposalList.ProposalList[0].Status = Active
		proposalList.ProposalList[0].EndHeight = new(big.Int).Add(oldActive.EndHeight, s.ContractRef().BlockHeight())
	}

	err = setProposalList(s, proposalList)
	if err != nil {
		return fmt.Errorf("setActiveProposal, setProposalList error: %v", err)
	}
	return nil
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

func proposalIDKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_PROPOSAL_ID))
}

func proposalListKey() []byte {
	return utils.ConcatKey(this, []byte(SKP_PROPOSAL_LIST))
}
