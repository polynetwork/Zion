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

package mock

import (
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

const (
	EpochStart = uint64(0)
	EpochEnd   = uint64(10000000000)
)

func init() {
	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.LvlTrace)
	log.Root().SetHandler(glogger)
}

func makeGenesis(vals []common.Address) *core.Genesis {
	core.RegGenesis = nil

	genesis := &core.Genesis{
		Config: &params.ChainConfig{
			ChainID:             big.NewInt(60801),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			BerlinBlock:         big.NewInt(0),
			LondonBlock:         big.NewInt(0),
			HotStuff:            &params.HotStuffConfig{Protocol: "basic"},
		},
		CommunityRate:    big.NewInt(2000),
		CommunityAddress: common.HexToAddress("0x79ad3ca3faa0F30f4A0A2839D2DaEb4Eb6B6820D"),
		Coinbase:         common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Difficulty:       big.NewInt(1),
		GasLimit:         2097151,
		Nonce:            4976618949627435365,
		Mixhash:          common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		ParentHash:       common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		Timestamp:        0,
	}

	govAccs := make([]core.GovernanceAccount, len(vals))
	for i := 0; i < len(vals); i++ {
		govAccs[i] = core.GovernanceAccount{
			Validator: vals[i],
			Signer:    common.EmptyAddress,
		}
	}
	genesis.Governance = govAccs

	valset := validator.NewSet(vals, hotstuff.RoundRobin)
	genesis.ExtraData, _ = types.GenerateExtraWithSignature(EpochStart, EpochEnd, valset.AddressList(), []byte{}, [][]byte{})
	return genesis
}
