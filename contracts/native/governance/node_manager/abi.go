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

const (
	EventRegisterCandidate   = "registerCandidate"
	EventUnRegisterCandidate = "unRegisterCandidate"
	EventApproveCandidate    = "approveCandidate"
	EventBlackNode           = "blackNode"
	EventWhiteNode           = "whiteNode"
	EventQuitNode            = "quitNode"
	EventUpdateConfig        = "updateConfig"
	EventCommitDpos          = "commitDpos"
)

const abijson = `[
	{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"signs","type":"uint64"}],"name":"CheckConsensusSignsEvent","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"Pubkey","type":"string"}],"name":"` + EventApproveCandidate + `","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"string[]","name":"PubkeyList","type":"string[]"}],"name":"` + EventBlackNode + `","type":"event"},
    {"anonymous":false,"inputs":[],"name":"` + EventCommitDpos + `","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"Pubkey","type":"string"}],"name":"` + EventQuitNode + `","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"Pubkey","type":"string"}],"name":"` + EventRegisterCandidate + `","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"Pubkey","type":"string"}],"name":"` + EventUnRegisterCandidate + `","type":"event"},
    {"anonymous":false,"inputs":[{"components":[{"internalType":"uint32","name":"BlockMsgDelay","type":"uint32"},{"internalType":"uint32","name":"HashMsgDelay","type":"uint32"},{"internalType":"uint32","name":"PeerHandshakeTimeout","type":"uint32"},{"internalType":"uint32","name":"MaxBlockChangeView","type":"uint32"}],"indexed":false,"internalType":"struct node_manager.Configuration","name":"Config","type":"tuple"}],"name":"` + EventUpdateConfig + `","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"Pubkey","type":"string"}],"name":"` + EventWhiteNode + `","type":"event"},
    {"inputs":[{"internalType":"string","name":"PeerPubkey","type":"string"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodApproveCandidate + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"string[]","name":"PeerPubkeyList","type":"string[]"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodBlackNode + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[],"name":"` + MethodCommitDpos + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[],"name":"` + MethodContractName + `","outputs":[{"internalType":"string","name":"Name","type":"string"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"uint32","name":"BlockMsgDelay","type":"uint32"},{"internalType":"uint32","name":"HashMsgDelay","type":"uint32"},{"internalType":"uint32","name":"PeerHandshakeTimeout","type":"uint32"},{"internalType":"uint32","name":"MaxBlockChangeView","type":"uint32"},{"internalType":"string","name":"VrfValue","type":"string"},{"internalType":"string","name":"VrfProof","type":"string"},{"components":[{"internalType":"uint32","name":"Index","type":"uint32"},{"internalType":"string","name":"PeerPubkey","type":"string"},{"internalType":"string","name":"Address","type":"string"}],"internalType":"struct node_manager.VBFTPeerInfo","name":"Peers","type":"tuple"}],"name":"` + MethodInitConfig + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"string","name":"PeerPubkey","type":"string"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodQuitNode + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"string","name":"PeerPubkey","type":"string"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodRegisterCandidate + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"string","name":"PeerPubkey","type":"string"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodUnRegisterCandidate + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"components":[{"components":[{"internalType":"uint32","name":"BlockMsgDelay","type":"uint32"},{"internalType":"uint32","name":"HashMsgDelay","type":"uint32"},{"internalType":"uint32","name":"PeerHandshakeTimeout","type":"uint32"},{"internalType":"uint32","name":"MaxBlockChangeView","type":"uint32"}],"internalType":"struct node_manager.Configuration","name":"Config","type":"tuple"}],"internalType":"struct node_manager.UpdateConfigParam","name":"ConfigParam","type":"tuple"}],"name":"` + MethodUpdateConfig + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"string","name":"PeerPubkey","type":"string"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodWhiteNode + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"}
]`

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
