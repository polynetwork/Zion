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
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

var (
	GenesisMaxCommission                = new(big.Int).SetUint64(50)
	GenesisMinInitialStake              = new(big.Int).Mul(big.NewInt(100000), params.ZNT1)
	GenesisMaxDescLength         uint64 = 2048
	GenesisBlockPerEpoch                = new(big.Int).SetUint64(400000)
	GenesisConsensusValidatorNum uint64 = 4
	GenesisVoterValidatorNum     uint64 = 4
)

func init() {
	// store data in genesis block
	core.RegGenesis = func(db *state.StateDB, genesis *core.Genesis) error {
		data := genesis.Alloc
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
		// the order of peer in the list is random, so we must sort the list before store.
		// btw, the mpt tree only needs the value of state_object to be deterministic.
		sort.Slice(peers, func(i, j int) bool {
			return peers[i].Address.Hex() < peers[j].Address.Hex()
		})
		if _, err := StoreCommunityInfo(db, genesis.CommunityRate, genesis.CommunityAddress); err != nil {
			return err
		}
		if _, err := StoreGenesisEpoch(db, peers); err != nil {
			return err
		}
		if err := StoreGenesisGlobalConfig(db); err != nil {
			return err
		}

		return nil
	}
}

func StoreCommunityInfo(s *state.StateDB, communityRate *big.Int, communityAddress common.Address) (*CommunityInfo, error) {
	cache := (*state.CacheDB)(s)
	communityInfo := &CommunityInfo{
		CommunityRate:    communityRate,
		CommunityAddress: communityAddress,
	}
	if err := setGenesisCommunityInfo(cache, communityInfo); err != nil {
		return nil, err
	}
	return communityInfo, nil
}

func StoreGenesisEpoch(s *state.StateDB, peers []*Peer) (*EpochInfo, error) {
	cache := (*state.CacheDB)(s)
	epoch := &EpochInfo{
		ID:          StartEpochID,
		Validators:  peers,
		Voters:      peers,
		StartHeight: new(big.Int),
	}

	// store current epoch and epoch info
	if err := setGenesisEpochInfo(cache, epoch); err != nil {
		return nil, err
	}
	return epoch, nil
}

func StoreGenesisGlobalConfig(s *state.StateDB) error {
	cache := (*state.CacheDB)(s)
	globalConfig := &GlobalConfig{
		MaxCommission:         GenesisMaxCommission,
		MinInitialStake:       GenesisMinInitialStake,
		MaxDescLength:         GenesisMaxDescLength,
		BlockPerEpoch:         GenesisBlockPerEpoch,
		ConsensusValidatorNum: GenesisConsensusValidatorNum,
		VoterValidatorNum:     GenesisVoterValidatorNum,
	}

	// store current epoch and epoch info
	if err := setGenesisGlobalConfig(cache, globalConfig); err != nil {
		return err
	}
	return nil
}
