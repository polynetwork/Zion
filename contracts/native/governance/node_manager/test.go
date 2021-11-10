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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/crypto"
)

// GenerateTestPeer ONLY used for testing
func GenerateTestPeer() *PeerInfo {
	pk, _ := crypto.GenerateKey()
	return &PeerInfo{
		PubKey:  hexutil.Encode(crypto.CompressPubkey(&pk.PublicKey)),
		Address: crypto.PubkeyToAddress(pk.PublicKey),
	}
}

func GenerateTestPeers(n int) *Peers {
	peers := &Peers{List: make([]*PeerInfo, n)}
	for i := 0; i < n; i++ {
		peers.List[i] = GenerateTestPeer()
	}
	return peers
}

func GenerateTestEpochInfo(id, height uint64, peersNum int) *EpochInfo {
	epoch := new(EpochInfo)
	epoch.ID = id
	epoch.StartHeight = height
	epoch.Peers = GenerateTestPeers(peersNum)
	epoch.Proposer = epoch.Peers.List[0].Address
	return epoch
}

func GenerateTestHash(n int) common.Hash {
	data := big.NewInt(int64(n))
	return common.BytesToHash(data.Bytes())
}

func GenerateTestHashList(n int) *HashList {
	data := &HashList{List: make([]common.Hash, n)}
	for i := 0; i < n; i++ {
		data.List[i] = GenerateTestHash(i + 1)
	}
	return data
}

func GenerateTestAddress(n int) common.Address {
	data := big.NewInt(int64(n))
	return common.BytesToAddress(data.Bytes())
}

func GenerateTestAddressList(n int) *AddressList {
	data := &AddressList{List: make([]common.Address, n)}
	for i := 0; i < n; i++ {
		data.List[i] = GenerateTestAddress(i + 1)
	}
	return data
}

func StoreTestEpoch(s *native.NativeContract, epoch *EpochInfo) error {
	if err := storeProposal(s, epoch.ID, epoch.Hash()); err != nil {
		return err
	}
	epoch.Status = ProposalStatusPassed
	storeCurrentEpochHash(s, epoch.Hash())
	storeEpochProof(s, epoch.ID, epoch.Hash())
	return storeEpoch(s, epoch)
}
