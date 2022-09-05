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

	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	//"github.com/ethereum/go-ethereum/consensus/hotstuff/backend"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// Encode generate hotstuff genesis extra
func Encode(validators []common.Address) (string, error) {
	var vanity []byte
	vanity = append(vanity, bytes.Repeat([]byte{0x00}, types.HotstuffExtraVanity)...)

	ist := &types.HotstuffExtra{
		Validators:    validators,
		Seal:          make([]byte, types.HotstuffExtraSeal),
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return "", err
	}

	return "0x" + common.Bytes2Hex(append(vanity, payload...)), nil
}

type Node struct {
	Address  string
	NodeKey  string
	Static   string
	KeyStore string
}

func SortNodes(src []*Node) []*Node {
	oriAddrs := make([]common.Address, len(src))
	idxMap := make(map[common.Address]int)
	for idx, v := range src {
		addr := common.HexToAddress(v.Address)
		oriAddrs[idx] = addr
		idxMap[addr] = idx
	}

	// sort address
	valset := validator.NewSet(oriAddrs, hotstuff.RoundRobin)
	//valset := backend.NewDefaultValSet(oriAddrs)

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
		list = append(list, common.HexToAddress(v.Address))
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

func NodeKey2NodeInfo(key string) (string, error) {
	if !strings.Contains(key, "0x") {
		key = "0x" + key
	}

	enc, err := hexutil.Decode(key)
	if err != nil {
		return "", err
	}

	privKey, err := crypto.ToECDSA(enc)
	if err != nil {
		return "", err
	}

	id := PubkeyID(&privKey.PublicKey)
	return id.String(), nil
}

func NodeKey2PublicInfo(key string) (string, error) {
	if !strings.Contains(key, "0x") {
		key = "0x" + key
	}

	dec, err := hexutil.Decode(key)
	if err != nil {
		return "", err
	}

	privKey, err := crypto.ToECDSA(dec)
	if err != nil {
		return "", err
	}

	enc := crypto.CompressPubkey(&privKey.PublicKey)
	return hexutil.Encode(enc), nil
}

func NodeStaticInfoTemp(src string) string {
	return fmt.Sprintf("enode://%s@127.0.0.1:30300?discport=0", src)
}
