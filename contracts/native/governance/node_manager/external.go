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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/params"
)

var (
	GenesisMaxCommissionChange, _        = new(big.Int).SetString("500", 10) // 50%
	GenesisMinInitialStake               = new(big.Int).Mul(big.NewInt(100000), params.ZNT1)
	GenesisMinProposalStake              = new(big.Int).Mul(big.NewInt(1000), params.ZNT1)
	// TODO: change GenesisBlockPerEpoch to 400000
	GenesisBlockPerEpoch                 = new(big.Int).SetUint64(10000)
	GenesisConsensusValidatorNum  uint64 = 4
	GenesisVoterValidatorNum      uint64 = 4

	MaxDescLength   int = 4000
	MaxValidatorNum int = 300
	MaxUnlockingNum int = 100
)

func init() {
	// store data in genesis block
	core.RegGenesis = func(db *state.StateDB, genesis *core.Genesis) error {
		data := genesis.Governance
		peers := make([]common.Address, 0, len(data))
		signers := make([]common.Address, 0, len(data))
		for _, v := range data {
			peers = append(peers, v.Validator)
			signers = append(signers, v.Signer)
		}
		if _, err := StoreCommunityInfo(db, genesis.CommunityRate, genesis.CommunityAddress); err != nil {
			return err
		}
		if _, err := StoreGenesisEpoch(db, peers, signers); err != nil {
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

func StoreGenesisEpoch(s *state.StateDB, peers []common.Address, signers []common.Address) (*EpochInfo, error) {
	cache := (*state.CacheDB)(s)
	epoch := &EpochInfo{
		ID:          StartEpochID,
		Validators:  peers,
		Signers:     signers,
		Voters:      signers,
		Proposers:   signers,
		StartHeight: new(big.Int),
		EndHeight:   GenesisBlockPerEpoch,
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
		MaxCommissionChange:   GenesisMaxCommissionChange,
		MinInitialStake:       GenesisMinInitialStake,
		MinProposalStake:      GenesisMinProposalStake,
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
