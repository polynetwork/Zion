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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func init() {
	// store data in genesis block
	core.RegGenesis = func(db *state.StateDB, data core.GenesisAlloc) error {
		peers := &Peers{List: make([]*PeerInfo, 0)}
		for addr, v := range data {
			pubkey, err := crypto.DecompressPubkey(v.PublicKey)
			if err != nil {
				return fmt.Errorf("store genesis peers, decompress pubkey failed, err: %v", err)
			}
			if got := crypto.PubkeyToAddress(*pubkey); got != addr {
				return fmt.Errorf("store genesis peers, expect address %s got %s", addr.Hex(), got.Hex())
			}
			peer := &PeerInfo{Address: addr, PubKey: hexutil.Encode(v.PublicKey)}
			peers.List = append(peers.List, peer)
		}
		sort.Sort(peers)
		if _, err := storeGenesisEpoch(db, peers); err != nil {
			return err
		} else {
			return nil
		}
	}
}

func storeGenesisEpoch(s *state.StateDB, peers *Peers) (*EpochInfo, error) {
	cache := (*state.CacheDB)(s)
	epoch := &EpochInfo{
		ID:          StartEpochID,
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

	// store genesis epoch id to list
	value, err := rlp.EncodeToBytes(&HashList{List: []common.Hash{epoch.Hash()}})
	if err != nil {
		return nil, err
	}
	proposalKey := proposalsKey(epoch.ID)
	cache.Put(proposalKey, value)

	// store genesis epoch proof
	key := EpochProofKey(EpochProofHash(epoch.ID))
	cache.Put(key, epoch.Hash().Bytes())

	return epoch, nil
}
