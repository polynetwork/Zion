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
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func StoreGenesisEpoch(s *state.StateDB, peers *Peers) (*EpochInfo, error) {
	cache := (*state.CacheDB)(s)
	epoch := &EpochInfo{
		ID:          StartEpoch,
		Peers:       peers,
		StartHeight: 0,
	}

	// store current epoch and epoch info
	if err := setEpoch(cache, epoch); err != nil {
		return nil, err
	}

	// store current hash
	curKey := curEpochKey()
	cache.Put(curKey, epoch.Hash().Bytes())
	return epoch, nil
}

func GetCurrentEpoch(s *state.StateDB) (*EpochInfo, error) {
	cache := (*state.CacheDB)(s)
	return getCurEpoch(cache)
}

func getCurEpoch(cache *state.CacheDB) (*EpochInfo, error) {
	// get current hash
	curKey := curEpochKey()
	enc, err := cache.Get(curKey)
	if err != nil {
		return nil, err
	}
	hash := common.BytesToHash(enc)

	// get current epoch info
	key := epochKey(hash)
	value, err := cache.Get(key)
	if err != nil {
		return nil, err
	}

	var epoch *EpochInfo
	if err := rlp.DecodeBytes(value, &epoch); err != nil {
		return nil, err
	}
	return epoch, nil
}

func checkAuthority(origin common.Address, epoch *EpochInfo) error {
	if epoch == nil || epoch.Peers == nil || epoch.Peers.List == nil {
		return fmt.Errorf("invalid epoch")
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
