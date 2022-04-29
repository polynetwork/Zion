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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

func getCurrentEpoch(s *native.NativeContract) (*EpochInfo, error) {
	// current epoch taking effective
	curEpochHash, err := getCurrentEpochHash(s)
	if err != nil {
		return nil, err
	}
	cur, err := getEpoch(s, curEpochHash)
	if err != nil {
		return nil, err
	}
	if cur.ID == StartEpochID {
		return cur, nil
	}

	height := s.ContractRef().BlockHeight().Uint64()
	if height >= cur.StartHeight {
		return cur, nil
	}

	if cur.ID-1 < StartEpochID {
		return nil, fmt.Errorf("epoch id should greater than %d", StartEpochID)
	} else {
		return getEffectiveEpochByID(s, cur.ID-1)
	}
}

func getEffectiveEpochByID(s *native.NativeContract, epochID uint64) (*EpochInfo, error) {
	if epochID < StartEpochID {
		return nil, fmt.Errorf("epoch %d not exist", epochID)
	}
	list, err := getProposals(s, epochID)
	if err != nil {
		return nil, fmt.Errorf("epoch %d has no proposal, err: %v", epochID, err)
	}
	if list == nil || len(list) == 0 {
		return nil, fmt.Errorf("epoch %d has no proposal", epochID)
	}
	if len(list) > 1 {
		return nil, fmt.Errorf("epoch %d has multi proposal", epochID)
	}
	lastEpochHash := list[0]
	return getEpoch(s, lastEpochHash)
}

func getEpochListByID(s *native.NativeContract, epochID uint64) ([]*EpochInfo, error) {
	if epochID < StartEpochID {
		return nil, fmt.Errorf("epoch %d not exist", epochID)
	}
	list, err := getProposals(s, epochID)
	if err != nil {
		return nil, fmt.Errorf("epoch %d has no proposal, err: %v", epochID, err)
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("epoch %d has no proposal", epochID)
	}
	var epochList []*EpochInfo
	for _, v := range list {
		item, err := getEpoch(s, v)
		if err != nil {
			return nil, fmt.Errorf("proposal %d not exists, err: %v", v, err)
		}
		epochList = append(epochList, item)
	}

	return epochList, nil
}

func getChangingEpoch(s *native.NativeContract) (*EpochInfo, error) {
	curEpochHash, err := getCurrentEpochHash(s)
	if err != nil {
		return nil, err
	}
	epoch, err := getEpoch(s, curEpochHash)
	if err != nil {
		return nil, err
	}

	height := s.ContractRef().BlockHeight().Uint64()
	if height > epoch.StartHeight {
		log.Warn("getChangingEpoch", "epoch changing invalidation, start height", epoch.StartHeight, "current height", height)
		return nil, fmt.Errorf("epoch invalid")
	}
	return epoch, nil
}

func CheckAuthority(origin, caller common.Address, epoch *EpochInfo) error {
	if epoch == nil || epoch.Peers == nil || epoch.Peers.List == nil {
		return fmt.Errorf("invalid epoch")
	}
	if origin == common.EmptyAddress || caller == common.EmptyAddress {
		return fmt.Errorf("origin/caller is empty address")
	}
	if origin != caller {
		return fmt.Errorf("origin must be caller")
	}
	for _, v := range epoch.Peers.List {
		if v.Address == origin {
			return nil
		}
	}
	return fmt.Errorf("tx origin %s is not valid validator", origin.Hex())
}

func checkPeer(peer *PeerInfo) error {
	if peer == nil || peer.Address == common.EmptyAddress || peer.PubKey == "" {
		return fmt.Errorf("invalid peer")
	}

	dec, err := hexutil.Decode(peer.PubKey)
	if err != nil {
		return err
	}
	pubkey, err := crypto.DecompressPubkey(dec)
	if err != nil {
		return err
	}
	addr := crypto.PubkeyToAddress(*pubkey)
	if addr == common.EmptyAddress {
		return fmt.Errorf("invalid pubkey")
	}
	if addr != peer.Address {
		return fmt.Errorf("pubkey not match address")
	}
	return nil
}

func generateEmptyContext(db *state.StateDB) *native.NativeContract {
	caller := common.EmptyAddress
	ref := native.NewContractRef(db, caller, caller, common.Big0, common.EmptyHash, 0, nil)
	ctx := native.NewNativeContract(db, ref)
	return ctx
}
