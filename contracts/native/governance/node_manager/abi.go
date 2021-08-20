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
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const abijson = ``

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(abijson))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

type VBFTConfig struct {
	BlockMsgDelay        uint32          `json:"block_msg_delay"`
	HashMsgDelay         uint32          `json:"hash_msg_delay"`
	PeerHandshakeTimeout uint32          `json:"peer_handshake_timeout"`
	MaxBlockChangeView   uint32          `json:"max_block_change_view"`
	VrfValue             string          `json:"vrf_value"`
	VrfProof             string          `json:"vrf_proof"`
	Peers                []*VBFTPeerInfo `json:"peers"`
}

type VBFTPeerInfo struct {
	Index      uint32 `json:"index"`
	PeerPubkey string `json:"peerPubkey"`
	Address    string `json:"address"`
}

type RegisterPeerParam struct {
	PeerPubkey string
	Address    common.Address
}

type PeerParam struct {
	PeerPubkey string
	Address    common.Address
}

type PeerListParam struct {
	PeerPubkeyList []string
	Address        common.Address
}

type UpdateConfigParam struct {
	Configuration *Configuration
}

type Configuration struct {
	BlockMsgDelay        uint32
	HashMsgDelay         uint32
	PeerHandshakeTimeout uint32
	MaxBlockChangeView   uint32
}
