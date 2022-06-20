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
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
)

func init() {
	// store data in genesis block
	core.RegGenesis = func(db *state.StateDB, data core.GenesisAlloc) error {
		peers := make([]*Peer, 0, len(data))
		for addr, v := range data {
			pk := hexutil.Encode(v.PublicKey)
			pubkey, err := crypto.DecompressPubkey(v.PublicKey)
			if err != nil {
				return fmt.Errorf("store genesis peers, decompress pubkey failed, err: %v", err)
			}
			if got := crypto.PubkeyToAddress(*pubkey); got != addr {
				return fmt.Errorf("store genesis peers, expect address %s got %s", addr.Hex(), got.Hex())
			}
			peer := &Peer{PubKey: pk, Address: addr}
			peers = append(peers, peer)
		}
		if _, err := storeGenesisEpoch(db, peers); err != nil {
			return err
		} else {
			return nil
		}
	}
}

func storeGenesisEpoch(s *state.StateDB, peers []*Peer) (*EpochInfo, error) {
	cache := (*state.CacheDB)(s)
	epoch := &EpochInfo{
		ID:          StartEpochID,
		Validators:  peers,
		Voters:      peers,
		StartHeight: common.Big0,
	}

	// store current epoch and epoch info
	if err := setGenesisEpochInfo(cache, epoch); err != nil {
		return nil, err
	}

	return epoch, nil
}
