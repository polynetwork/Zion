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
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "node manager"

const (
	MethodContractName        = "name"
	MethodRegisterCandidate   = "registerCandidate"
	MethodUnRegisterCandidate = "unRegisterCandidate"
	MethodApproveCandidate    = "approveCandidate"
	MethodBlackNode           = "blackNode"
	MethodWhiteNode           = "whiteNode"
	MethodQuitNode            = "quitNode"
	MethodCommitDpos          = "commitDpos"

	EventApproveCandidate    = "EventApproveCandidate"
	EventBlackNode           = "EventBlackNode"
	EventCommitDpos          = "EventCommitDpos"
	EventQuitNode            = "EventQuitNode"
	EventRegisterCandidate   = "EventRegisterCandidate"
	EventUnRegisterCandidate = "EventUnRegisterCandidate"
	EventWhiteNode           = "EventWhiteNode"
	EventCheckConsensusSigns = "CheckConsensusSignsEvent"
)

const abijson = `[
    {"type":"function","name":"` + MethodContractName + `","inputs":[],"outputs":[{"internalType":"string","name":"Name","type":"string"}],"stateMutability":"nonpayable"},
	{"type":"function","name":"` + MethodRegisterCandidate + `","inputs":[{"internalType":"string","name":"PeerPubkey","type":"string"},{"internalType":"address","name":"Address","type":"address"}],"outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable"},
    {"type":"function","name":"` + MethodUnRegisterCandidate + `","inputs":[{"internalType":"string","name":"PeerPubkey","type":"string"},{"internalType":"address","name":"Address","type":"address"}],"outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable"},
	{"type":"function","name":"` + MethodApproveCandidate + `","inputs":[{"internalType":"string","name":"PeerPubkey","type":"string"},{"internalType":"address","name":"Address","type":"address"}],"outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"ID","type":"uint64"}],"name":"` + EventApproveRegisterRelayer + `","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"ID","type":"uint64"}],"name":"` + EventApproveRemoveRelayer + `","type":"event"},
	{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"applyID","type":"uint64"}],"name":"` + EventRegisterRelayer + `","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"removeID","type":"uint64"}],"name":"` + EventRemoveRelayer + `","type":"event"}
]`

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(node_manager_abi.NodeManagerABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

var (
	ABI  *abi.ABI
	this = utils.NodeManagerContractAddress
)

type RegisterPeerParam struct {
	PeerPubkey string
	Address    common.Address
}

func (p *RegisterPeerParam) Serialization(sink *common.ZeroCopySink) {
	sink.WriteString(p.PeerPubkey)
	sink.WriteVarBytes(p.Address[:])
}

func (p *RegisterPeerParam) Deserialization(source *common.ZeroCopySource) error {
	peerPubkey, eof := source.NextString()
	if eof {
		return fmt.Errorf("source.NextString, deserialize peerPubkey error")
	}
	address, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("source.NextVarBytes, deserialize address error")
	}
	addr, err := common.AddressParseFromBytes(address)
	if err != nil {
		return fmt.Errorf("common.AddressParseFromBytes, deserialize address error: %s", err)
	}

	p.PeerPubkey = peerPubkey
	p.Address = addr
	return nil
}

type PeerParam struct {
	PeerPubkey string
	Address    common.Address
}

func (p *PeerParam) Serialization(sink *common.ZeroCopySink) {
	sink.WriteString(p.PeerPubkey)
	sink.WriteVarBytes(p.Address[:])
}

func (p *PeerParam) Deserialization(source *common.ZeroCopySource) error {
	peerPubkey, eof := source.NextString()
	if eof {
		return fmt.Errorf("source.NextString, deserialize peerPubkey error")
	}
	address, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("source.NextVarBytes, deserialize address error")
	}
	addr, err := common.AddressParseFromBytes(address)
	if err != nil {
		return fmt.Errorf("common.AddressParseFromBytes, deserialize address error: %s", err)
	}

	p.PeerPubkey = peerPubkey
	p.Address = addr
	return nil
}

type PeerListParam struct {
	PeerPubkeyList []string
	Address        common.Address
}

func (p *PeerListParam) Serialization(sink *common.ZeroCopySink) {
	sink.WriteVarUint(uint64(len(p.PeerPubkeyList)))
	for _, v := range p.PeerPubkeyList {
		sink.WriteString(v)
	}
	sink.WriteVarBytes(p.Address[:])
}

func (p *PeerListParam) Deserialization(source *common.ZeroCopySource) error {
	n, eof := source.NextVarUint()
	if eof {
		return fmt.Errorf("source.NextVarUint, deserialize PeerPubkeyList length error")
	}
	peerPubkeyList := make([]string, 0)
	for i := 0; uint64(i) < n; i++ {
		k, eof := source.NextString()
		if eof {
			return fmt.Errorf("source.NextString, deserialize peerPubkey error")
		}
		peerPubkeyList = append(peerPubkeyList, k)
	}

	address, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("source.NextVarBytes, deserialize address error")
	}
	addr := common.BytesToAddress(address)
	p.PeerPubkeyList = peerPubkeyList
	p.Address = addr
	return nil
}
