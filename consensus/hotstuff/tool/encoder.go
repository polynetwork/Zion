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

package tool

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/backend"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// EncodeGenesisExtra generate hotstuff genesis extra
func EncodeGenesisExtra(validators []common.Address) (string, error) {

	// 1. sort validators
	valset := validator.NewSet(validators, hotstuff.RoundRobin)
	validators = valset.AddressList()

	// 2. set vanity
	var vanity []byte
	vanity = append(vanity, bytes.Repeat([]byte{0x00}, types.HotstuffExtraVanity)...)

	// 3. construct extra
	ist := &types.HotstuffExtra{
		StartHeight:   0,
		EndHeight:     nm.GenesisBlockPerEpoch.Uint64(),
		Validators:    validators,
		Seal:          make([]byte, types.HotstuffExtraSeal),
		CommittedSeal: [][]byte{},
	}

	// 4. serialization
	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return "", err
	}

	return "0x" + common.Bytes2Hex(append(vanity, payload...)), nil
}

type Node struct {
	once sync.Once
	addr common.Address
	pub  *ecdsa.PublicKey
	id   Discv5NodeID

	NodeKey string
	PubKey  string
	Static  string
}

func (n *Node) init() {
	if n.PubKey != "" {
		if !strings.Contains(n.PubKey, "0x") {
			n.PubKey = "0x" + n.PubKey
		}
		enc, err := hexutil.Decode(n.PubKey)
		if err != nil {
			panic(fmt.Sprintf("pubkey is not hex string, err: %v", err))
		}
		pubkey, err := crypto.DecompressPubkey(enc)
		if err != nil {
			panic(fmt.Sprintf("pubkey convert failed, err: %v", err))
		}
		n.pub = pubkey
	} else if n.NodeKey != "" {
		if !strings.Contains(n.NodeKey, "0x") {
			n.NodeKey = "0x" + n.NodeKey
		}

		enc, err := hexutil.Decode(n.NodeKey)
		if err != nil {
			panic(fmt.Sprintf("node key is not hex string, err: %v", err))
		}

		privKey, err := crypto.ToECDSA(enc)
		if err != nil {
			panic(fmt.Sprintf("ecdsa convert failed, err: %v", err))
		}
		n.pub = &privKey.PublicKey
	} else {
		panic("nodekey and pubkey are empty")
	}

	n.addr = crypto.PubkeyToAddress(*n.pub)
	n.id = PubkeyID(n.pub)
	n.Static = fmt.Sprintf("enode://%s@127.0.0.1:30300?discport=0", n.id)
}

func (n *Node) Address() string {
	if n.addr == common.EmptyAddress {
		n.once.Do(n.init)
	}
	return n.addr.Hex()
}

func (n *Node) Pubkey() string {
	if n.pub == nil {
		n.once.Do(n.init)
	}
	enc := crypto.CompressPubkey(n.pub)
	return hexutil.Encode(enc)
}

func SortNodes(src []*Node) []*Node {
	oriAddrs := make([]common.Address, len(src))
	idxMap := make(map[common.Address]int)
	for idx, v := range src {
		v.once.Do(v.init)
		oriAddrs[idx] = v.addr
		idxMap[v.addr] = idx
	}

	// sort address
	valset := backend.NewDefaultValSet(oriAddrs)

	list := make([]*Node, 0)
	for _, val := range valset.AddressList() {
		idx := idxMap[val]
		list = append(list, src[idx])
	}

	return list
}

func NodesAddress(src []*Node) []common.Address {
	list := make([]common.Address, 0)
	for _, v := range src {
		v.once.Do(v.init)
		list = append(list, v.addr)
	}
	return list
}

type Discv5NodeID [64]byte

func (n Discv5NodeID) String() string {
	return fmt.Sprintf("%x", n[:])
}

// PubkeyID returns a marshaled representation of the given public key.
func PubkeyID(pub *ecdsa.PublicKey) Discv5NodeID {
	var id Discv5NodeID
	pbytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	if len(pbytes)-1 != len(id) {
		panic(fmt.Errorf("need %d bit pubkey, got %d bits", (len(id)+1)*8, len(pbytes)))
	}
	copy(id[:], pbytes[1:])
	return id
}
