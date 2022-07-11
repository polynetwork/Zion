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
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestEncodeGenesisExtra(t *testing.T) {
	validators := []common.Address{
		common.HexToAddress("0x258af48e28e4a6846e931ddff8e1cdf8579821e5"),
		common.HexToAddress("0x6a708455c8777630aac9d1e7702d13f7a865b27c"),
		common.HexToAddress("0x8c09d936a1b408d6e0afaa537ba4e06c4504a0ae"),
		common.HexToAddress("0xad3bf5ed640cc72f37bd21d64a65c3c756e9c88c"),
	}
	enc, err := EncodeGenesisExtra(validators)
	assert.NoError(t, err)
	t.Log(enc)
}

var testOriginValAndNodeKeys = []*Node{
	// local side chain
	{
		Address: "0x09f4E484D43B3D6b20957F7E1760beE3C6F62186",
		NodeKey: "562aa98da69477996bd82422b97698541f25e71ba2f803970947b3ad8bdb7afa",
	},
	{
		Address: "0x294b8211E7010f457d85942aC874d076D739E32a",
		NodeKey: "8bea3ce27136df435ada62a40a4226404879b3c42e2e86ba9a236b4a61c99c26",
	},
	{
		Address: "0x9deAD91D8632DCEEC701710bAF7922324DD45F58",
		NodeKey: "53f7d9ec7657cdd3a3eaa8ddd126d36fbc60203448fca1bbfccec0d59d173da6",
	},
	{
		Address: "0xc5e2344b875e236b3475e9e4E70448525cA5210F",
		NodeKey: "305baf1e19a2da40b413dfb62b206b0ac74cb3d7e975cb70fe8391cbbe174f2a",
	},

	// local main chain
	{
		Address: "0x258af48e28e4a6846e931ddff8e1cdf8579821e5",
		NodeKey: "4b0c9b9d685db17ac9f295cb12f9d7d2369f5bf524b3ce52ce424031cafda1ae",
	},
	{
		Address: "0x6a708455c8777630aac9d1e7702d13f7a865b27c",
		NodeKey: "3d9c828244d3b2da70233a0a2aea7430feda17bded6edd7f0c474163802a431c",
	},
	{
		Address: "0x8c09d936a1b408d6e0afaa537ba4e06c4504a0ae",
		NodeKey: "cc69b13ca2c5cd4d76bb881f6ad18d93bd947042c0f3a7adc80bdd17dac68210",
	},
	{
		Address: "0xad3bf5ed640cc72f37bd21d64a65c3c756e9c88c",
		NodeKey: "018c71d5e3b245117ffba0975e46129371473c6a1d231c5eddf7a8364d704846",
	},
	{
		Address: "0xc095448424a5ecd5ca7ccdadfaad127a9d7e88ec",
		NodeKey: "49e26aa4d60196153153388a24538c2693d65f0010a3a488c0c4c2b2a64b2de4",
	},
	{
		Address: "0xd47a4e56e9262543db39d9203cf1a2e53735f834",
		NodeKey: "9fc1723cff3bc4c11e903a53edb3b31c57b604bfc88a5d16cfec6a64fbf3141c",
	},
	{
		Address: "0xbfb558f0dceb07fbb09e1c283048b551a4310921",
		NodeKey: "5555ebb339d3d5ed1efbf0ca96f5b145134e5ce8044fec693558056d268776ae",
	},
}

func TestEncodeSalt(t *testing.T) {
	dumpNodes(t, testOriginValAndNodeKeys)
}

func TestGenerateAndEncode(t *testing.T) {
	nodes := generateNodes(4)
	dumpNodes(t, nodes)
}

func dumpNodes(t *testing.T, nodes []*Node) {
	sortedNodes := SortNodes(nodes)
	staticNodes := make([]string, 0)
	for _, v := range sortedNodes {
		nodeInf, err := NodeKey2NodeInfo(v.NodeKey)
		if err != nil {
			t.Fatal(err)
		}
		pubInf, err := NodeKey2PublicInfo(v.NodeKey)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("addr: %s, pubKey: %s, nodeKey: %s, static-node-info:%s", v.Address, pubInf, v.NodeKey, nodeInf)
		staticNodes = append(staticNodes, NodeStaticInfoTemp(nodeInf))
	}
	t.Log("==================================================================")

	genesis, err := EncodeGenesisExtra(NodesAddress(sortedNodes))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("genesis extra %s", genesis)

	t.Log("==================================================================")

	staticNodesEnc, err := json.MarshalIndent(staticNodes, "", "\t")
	t.Log(string(staticNodesEnc))
}

func generateNodes(n int) []*Node {
	nodes := make([]*Node, 0)

	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)

		nodekey := hexutil.Encode(crypto.FromECDSA(key))
		nodeInf, _ := NodeKey2NodeInfo(nodekey)

		staticInf := NodeStaticInfoTemp(nodeInf)

		node := &Node{
			Address: addr.Hex(),
			NodeKey: nodekey,
			Static:  staticInf,
		}
		nodes = append(nodes, node)
	}

	return nodes
}
