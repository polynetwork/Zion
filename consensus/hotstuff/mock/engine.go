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
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/backend"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

type Engine consensus.Engine

// backend is engine but also hotstuff engine and consensus handler.
func makeEngine(privateKey *ecdsa.PrivateKey, db ethdb.Database) Engine {
	config := hotstuff.DefaultBasicConfig
	chainConfig := &params.ChainConfig{
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
	}
	engine := backend.New(chainConfig, config, privateKey, db, true)
	broadcaster := makeBroadcaster(engine.Address(), engine)
	engine.SetBroadcaster(broadcaster)
	return engine
}
