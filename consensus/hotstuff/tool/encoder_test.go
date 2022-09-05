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
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

/*
static nodes
[
	"enode://44e509103445d5e8fd290608308d16d08c739655d6994254e413bc1a067838564f7a32ed8fed182450ec2841856c0cc0cd313588a6e25002071596a7363e84b6@127.0.0.1:30300?discport=0",
	"enode://3884de29148505a8d862992e5721767d4b47ff52ffab4c2d2527182d812a6d95d2049e00b7c5579ca7b86b3dba8c935e742d2dfde9ae16abb5e3265e33a6d472@127.0.0.1:30301?discport=0",
	"enode://c07fb7d48eac559a2483e249d27841c18c7ce5dbbbf2796a6963cc9cef27cabd2e1bc9c456a83f0777a98dfd6e7baf272739b9e5f8febf0077dc09509c2dfa48@127.0.0.1:30302?discport=0",
	"enode://ecac0ebe7224cfd04056c940605a4a9d4cb0367cf5819bf7e5502bf44f68bdd471a6b215c733f4a4ab6a1b417ec18b2e382e83d2e1a4d7936b437e8c047b41f5@127.0.0.1:30303?discport=0",
]
对应地址为:
0xc095448424a5ecd5ca7ccdadfaad127a9d7e88ec
0xd47a4e56e9262543db39d9203cf1a2e53735f834
0x258af48e28e4a6846e931ddff8e1cdf8579821e5
0x8c09d936a1b408d6e0afaa537ba4e06c4504a0ae
*/

/*
原始extra
0x0000000000000000000000000000000000000000000000000000000000000000f89af854 94 c095448424a5ecd5ca7ccdadfaad127a9d7e88ec 94 d47a4e56e9262543db39d9203cf1a2e53735f834 94 258af48e28e4a6846e931ddff8e1cdf8579821e5 94 8c09d936a1b408d6e0afaa537ba4e06c4504a0ae b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0
重新排序后生成的extra
0x0000000000000000000000000000000000000000000000000000000000000000f89af854 94 258af48e28e4a6846e931ddff8e1cdf8579821e5 94 8c09d936a1b408d6e0afaa537ba4e06c4504a0ae 94 c095448424a5ecd5ca7ccdadfaad127a9d7e88ec 94 d47a4e56e9262543db39d9203cf1a2e53735f834 b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0
*/

func TestEncode(t *testing.T) {
	//validators := []common.Address{
	//	common.HexToAddress("0xc095448424a5ecd5ca7ccdadfaad127a9d7e88ec"),
	//	common.HexToAddress("0xd47a4e56e9262543db39d9203cf1a2e53735f834"),
	//	common.HexToAddress("0x258af48e28e4a6846e931ddff8e1cdf8579821e5"),
	//	common.HexToAddress("0x8c09d936a1b408d6e0afaa537ba4e06c4504a0ae"),
	//}
	validators := []common.Address{
		common.HexToAddress("0x24e7d00243aa0fb83a398d687c2951ad4c9bc288"),
		common.HexToAddress("0x0662c575eaa19c168dc8bcb83121dcf132f87b53"),
		common.HexToAddress("0x55547a2c919b9a9a84b0dd280f551ccc1316b22e"),
		common.HexToAddress("0x5055f522105732392b57bea3ba3aaaff69dece08"),
		common.HexToAddress("0x16bc0237a18dd154a74dc42db458b79b328c3e27"),
		common.HexToAddress("0x9a4e1e4e1662eddf557936cdabec398d9c856e91"),
		common.HexToAddress("0x2b06ef09277ed35eb83d10421e55e0b5ac6d8bf5"),
		common.HexToAddress("0xba49b8ea949d5c5c0f4d9281a989deacfb38d6c8"),
		common.HexToAddress("0xeb85568b5ba73e4eb6fc8e59e6c72a2dbd8b02fe"),
		common.HexToAddress("0x35f9783875c34ec9e18897c32b5ce74a98332eb3"),
		common.HexToAddress("0xe9fb4465894997c5e68944cd92a250bfe6e52ac3"),
		common.HexToAddress("0x324d0370899309aee59dc435698a6670015562d8"),
		common.HexToAddress("0x1ec3992eb7f1bfa545a092a5bb53008628b01801"),
	}
	valset := validator.NewSet(validators, hotstuff.RoundRobin)
	validators = valset.AddressList()
	enc, err := Encode(validators)
	assert.NoError(t, err)
	t.Log(enc)
}

var testOriginValAndNodeKeys = []*Node{
	// local side chain
	//{
	//	Address: "0x09f4E484D43B3D6b20957F7E1760beE3C6F62186",
	//	NodeKey: "562aa98da69477996bd82422b97698541f25e71ba2f803970947b3ad8bdb7afa",
	//},
	//{
	//	Address: "0x294b8211E7010f457d85942aC874d076D739E32a",
	//	NodeKey: "8bea3ce27136df435ada62a40a4226404879b3c42e2e86ba9a236b4a61c99c26",
	//},
	//{
	//	Address: "0x9deAD91D8632DCEEC701710bAF7922324DD45F58",
	//	NodeKey: "53f7d9ec7657cdd3a3eaa8ddd126d36fbc60203448fca1bbfccec0d59d173da6",
	//},
	//{
	//	Address: "0xc5e2344b875e236b3475e9e4E70448525cA5210F",
	//	NodeKey: "305baf1e19a2da40b413dfb62b206b0ac74cb3d7e975cb70fe8391cbbe174f2a",
	//},

	// local main chain
	//{
	//	Address: "0x258af48e28e4a6846e931ddff8e1cdf8579821e5",
	//	NodeKey: "4b0c9b9d685db17ac9f295cb12f9d7d2369f5bf524b3ce52ce424031cafda1ae",
	//},
	//{
	//	Address: "0x6a708455c8777630aac9d1e7702d13f7a865b27c",
	//	NodeKey: "3d9c828244d3b2da70233a0a2aea7430feda17bded6edd7f0c474163802a431c",
	//},
	//{
	//	Address: "0x8c09d936a1b408d6e0afaa537ba4e06c4504a0ae",
	//	NodeKey: "cc69b13ca2c5cd4d76bb881f6ad18d93bd947042c0f3a7adc80bdd17dac68210",
	//},
	//{
	//	Address: "0xad3bf5ed640cc72f37bd21d64a65c3c756e9c88c",
	//	NodeKey: "018c71d5e3b245117ffba0975e46129371473c6a1d231c5eddf7a8364d704846",
	//},
	//{
	//	Address: "0xc095448424a5ecd5ca7ccdadfaad127a9d7e88ec",
	//	NodeKey: "49e26aa4d60196153153388a24538c2693d65f0010a3a488c0c4c2b2a64b2de4",
	//},
	//{
	//	Address: "0xd47a4e56e9262543db39d9203cf1a2e53735f834",
	//	NodeKey: "9fc1723cff3bc4c11e903a53edb3b31c57b604bfc88a5d16cfec6a64fbf3141c",
	//},
	//{
	//	Address: "0xbfb558f0dceb07fbb09e1c283048b551a4310921",
	//	NodeKey: "5555ebb339d3d5ed1efbf0ca96f5b145134e5ce8044fec693558056d268776ae",
	//},
}

func TestEncodeSalt(t *testing.T) {
	dumpNodes(t, testOriginValAndNodeKeys)
}

func TestGenerateAndEncode(t *testing.T) {
	nodes := generateNodes(7)
	dumpNodes(t, nodes)
}

func TestGenerateKeyStore(t *testing.T) {
	nodekey := "26cc96a0d256d45e1515bf325bec1925746d796b3637b147f35a01d6c2d6399b"
	passphrase := "Onchain@Maas"

	privateKey, err := crypto.HexToECDSA(nodekey)
	if err != nil {
		log.Fatal(err)
	}

	keyjson, err := keystore.GenerateKeyJson(privateKey, passphrase)
	if err != nil {
		t.Fatalf("Error encrypting key: %v", err)
	}
	log.Println(string(keyjson))
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

	genesis, err := Encode(NodesAddress(sortedNodes))
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
