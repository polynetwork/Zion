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
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/tool -run TestIDAndPubkey
func TestIDAndPubkey(t *testing.T) {
	pk, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(pk.PublicKey)
	id := PubkeyID(&pk.PublicKey)
	t.Logf("origin address %s, node %s", addr.Hex(), id)

	pubKey := ID2PubKey(id[:])
	got := crypto.PubkeyToAddress(*pubKey)
	t.Logf("recover address %v", got.Hex())

	assert.Equal(t, addr, got)
}

// go test -v github.com/ethereum/go-ethereum/consensus/hotstuff/tool -run TestGenesisExtra
func TestGenesisExtra(t *testing.T) {
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

func TestEncode(t *testing.T) {
	var testNodeKeys = []*Node{
		{
			NodeKey: "562aa98da69477996bd82422b97698541f25e71ba2f803970947b3ad8bdb7afa",
		},
		{
			NodeKey: "8bea3ce27136df435ada62a40a4226404879b3c42e2e86ba9a236b4a61c99c26",
		},
		{
			NodeKey: "53f7d9ec7657cdd3a3eaa8ddd126d36fbc60203448fca1bbfccec0d59d173da6",
		},
		{
			NodeKey: "305baf1e19a2da40b413dfb62b206b0ac74cb3d7e975cb70fe8391cbbe174f2a",
		},
	}

	dumpNodes(t, testNodeKeys)
}


func TestEncodeSeeds(t *testing.T) {
	var addrs []common.Address
	for _, str := range []string {
		"04c1a1927b9a506ece82ed9db1bfc1da854ddc40e0a4a6a618bf9c5ca671ac893cd3653248c2853c6786ac6d03964717ae1a0a318de18a06a5ef83e8145f9daab7",
		"04b8e55b48e89532efb956e6b2732ba5e124fed65b96092fce32a89080411ac0f2eefafae3020e97e248571103b0c1906de9a7ef641dc9037f067f775646cba2a3",
		"04f9b194f397426f6540741114fb72cbf4d531fdb1e68a20e6ccc5d4d86e9a93014057e4bb0104e6cb0d3a37efdfe934f77b1c680663ad0e31b99878e2a8246bc2",
		"041d64c9eac537ea8536622c4a28ad2fde1869b036241264c0fe5090d495e9db0c355c53c7b03f29845396a067fd9ee6c73722c419eecd77bbe3008c562174bf91",
	} {
		data, err := hex.DecodeString(str)
		if err != nil {
			t.Fatal(err)
		}
		pub, err := crypto.UnmarshalPubkey(data)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(PubkeyID(pub).String())
		addrs = append(addrs, crypto.PubkeyToAddress(*pub))
	}

	valset := validator.NewSet(addrs, hotstuff.RoundRobin)
	genesis, err := EncodeGenesisExtra(valset.AddressList())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("genesis extra %s", genesis)
}

func TestGenerateAndEncode(t *testing.T) {
	nodes := generateNodes(4)
	dumpNodes(t, nodes)
}

func dumpNodes(t *testing.T, nodes []*Node) {
	sortedNodes := SortNodes(nodes)
	staticNodes := make([]string, 0)
	for _, v := range sortedNodes {
		t.Logf("addr: %s, pubKey: %s, nodeKey: %s", v.Address(), v.Pubkey(), v.NodeKey)
		staticNodes = append(staticNodes, v.Static)
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
		nodekey := hexutil.Encode(crypto.FromECDSA(key))
		node := &Node{
			NodeKey: nodekey,
		}
		nodes = append(nodes, node)
	}
	return nodes
}
