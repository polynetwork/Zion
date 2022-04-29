package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMaasGenesisFuncs(t *testing.T) {
	genesis := new(MaasGenesis)
	genesis.Default()
	genesisStr, err := genesis.Encode()
	t.Log(genesisStr)
	assert.NoError(t, err)
	got := new(MaasGenesis)
	err = got.Decode(genesisStr)
	assert.NoError(t, err)
	assert.Equal(t, genesis, got)
}

func TestDumpGenesis(t *testing.T) {
	filePaths := [3]string{
		"./temp/genesis.json",
		"./temp/static-nodes.json",
		"./temp/nodes.json",
	}

	contents := [3]string{}
	// not valid data, only for test
	contents[0] = "{\n    \"config\": {\n        \"chainId\": 10898, \n        \"homesteadBlock\": 0,\n        \"eip150Block\": 0,\n        \"eip155Block\": 0,\n        \"eip158Block\": 0,\n        \"byzantiumBlock\": 0,\n        \"constantinopleBlock\": 0,\n        \"petersburgBlock\": 0,\n        \"istanbulBlock\": 0,\n        \"hotstuff\": {\n            \"protocol\": \"basic\"\n        }\n    },\n    \"alloc\": {\n        \"0x0F45DAfCD39e1C59202e91306A1671cB7f8884be\": {\"publicKey\": \"0x039053297a9feb8c56ccd78592ed3d98dd6165131d3c9ee3398f720c9de65a7b74\", \"balance\": \"100000000000000000000000000000\"},\n        \"0x4ba9732E2358F41E682B7a7AAb71e614e08383dF\": {\"publicKey\": \"0x034b6bb61d6ab259d460701957784ba415ab05bf4b9093eab4eb64a2d4a588ffd2\", \"balance\": \"100000000000000000000000000000\"},\n        \"0x6519d82d761c275f01de6F197542dE924296928E\": {\"publicKey\": \"0x031edfa05742afbf0059141ebd0949e8f3766a6718e2d4b4df39ac7d96cda9bb46\", \"balance\": \"100000000000000000000000000000\"},\n        \"0x6cfd8d31A55CDf03303E280792C7D6cE855601f3\": {\"publicKey\": \"0x0270a210c0b6a437943423afef05403e1f9fbbe5dad4d9317233f4ad5e58ee17ff\", \"balance\": \"100000000000000000000000000000\"},\n        \"0x782833979973d83Cd48332E09938B95f9ba32B50\": {\"publicKey\": \"0x036f25d77c154017db892c763354e8cd7f54a91246f74c7a6400e3cba4f30e73b8\", \"balance\": \"100000000000000000000000000000\"},\n        \"0x990f7AaFA09FEa4583c0C72063b306cdE54e1e8F\": {\"publicKey\": \"0x03a9d10e5ca3475b47a7dc1a72679cb377b4dfb8dd5046b92bf8ba721f2985f053\", \"balance\": \"100000000000000000000000000000\"},\n        \"0x9BCd01b46b98254eE3611Eb1501A12780343f7D2\": {\"publicKey\": \"0x03f8e5aab12985e7e058d0b47afafa8ac5fb813fde8051354403c320d9cc38a3b1\", \"balance\": \"100000000000000000000000000000\"}\n    },\n    \"coinbase\": \"0x0000000000000000000000000000000000000000\",\n    \"difficulty\": \"0x1\",\n    \"extraData\": \"0x0000000000000000000000000000000000000000000000000000000000000000f8daf893940f45dafcd39e1c59202e91306a1671cb7f8884be944ba9732e2358f41e682b7a7aab71e614e08383df946519d82d761c275f01de6f197542de924296928e946cfd8d31a55cdf03303e280792c7d6ce855601f394782833979973d83cd48332e09938b95f9ba32b5094990f7aafa09fea4583c0c72063b306cde54e1e8f949bcd01b46b98254ee3611eb1501a12780343f7d2b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c080\",\n    \"gasLimit\": \"0xffffffff\",\n    \"nonce\": \"0x4510809143055965\",\n    \"mixhash\": \"0x0000000000000000000000000000000000000000000000000000000000000000\",\n    \"parentHash\": \"0x0000000000000000000000000000000000000000000000000000000000000000\",\n    \"timestamp\": \"0x00\"\n}"
	contents[1] = "[\n    \"enode://9053297a9feb8c56ccd78592ed3d98dd6165131d3c9ee3398f720c9de65a7b74fc18a9f870cf757d3aeed03c32dbacfc6aa41007b6b0192631859b7769d0689f@127.0.0.1:30300?discport=0\",\n    \"enode://4b6bb61d6ab259d460701957784ba415ab05bf4b9093eab4eb64a2d4a588ffd229fccde97d14c0c6f14ffcedd73d56e8be3a2cea49ec58243a01ad311edb943d@127.0.0.1:30300?discport=0\",\n    \"enode://1edfa05742afbf0059141ebd0949e8f3766a6718e2d4b4df39ac7d96cda9bb46a6a4e304c4cdd2d7bed2c29c7857b039d2bb08f4b81cee92cbbfe4155c76fc11@127.0.0.1:30300?discport=0\",\n    \"enode://70a210c0b6a437943423afef05403e1f9fbbe5dad4d9317233f4ad5e58ee17ff24b9f527b48ffdc5279ddc065e5810c8d647af85ab40a295676667ba085638c8@127.0.0.1:30300?discport=0\",\n    \"enode://6f25d77c154017db892c763354e8cd7f54a91246f74c7a6400e3cba4f30e73b81894600d3505c59df5ad9711621de07b80166dffacfcee9a2f4fff2b3f0c31dd@127.0.0.1:30300?discport=0\",\n    \"enode://a9d10e5ca3475b47a7dc1a72679cb377b4dfb8dd5046b92bf8ba721f2985f053538623149f3a7e59a78ae6f7c7d3305b76689b989863b12aab1c8b86837a3ef1@127.0.0.1:30300?discport=0\",\n    \"enode://f8e5aab12985e7e058d0b47afafa8ac5fb813fde8051354403c320d9cc38a3b1e34da681ffd1c41c99b50cdb04292d6228357686b264846f5b55e2ef67df547b@127.0.0.1:30300?discport=0\"\n]\n"
	contents[2] = "[\n        {\n                \"Address\": \"0x152e6e7C0d1637Cfc3C909852BD9914A67F91340\",\n                \"NodeKey\": \"0x22b2c3a6f4be86581dbc729acba0048d70e6a960847d61902dc54e0dad0b8a00\",\n                \"Static\": \"enode://13ec7f7dde00342f076bdbaafa1488b864fb7d34a744379a421d6546e4ea76f74d364fa3b1350622760674110ccae61c450631d69c20ac58b711f2e55c934bcf@127.0.0.1:30300?discport=0\"\n        },\n        {\n                \"Address\": \"0x1E306D7C8Ea7b042ecc974A7609a84b888453Cd8\",\n                \"NodeKey\": \"0xc9ac34375e9f739b5951058f17afb7264a77ad7decf086e85eafb1c9af8305a1\",\n                \"Static\": \"enode://c0e4ac8fa2daab3994b15ffa95284f62ae647e1a7b572f899d77eafd0b2f8be8e337260129acf98832a453ef25e749c10c13d45e5990fd27979c9c70357301f5@127.0.0.1:30300?discport=0\"\n        },\n        {\n                \"Address\": \"0x2B707ae427547Ae18c7bddB5FDbd78a14386E874\",\n                \"NodeKey\": \"0xf07946ed4c70a4c1bf7bdf33d0d7cd870b24f27c450c3a14cef824dde1e26b9c\",\n                \"Static\": \"enode://4deb701b087b33338e93bb87f6a843015a0d18987ff8f187f7f25dd822deedeb0bed2069070ace715e01bc9a0a0ae5016ae85e66f5293056571e3a0854aba26f@127.0.0.1:30300?discport=0\"\n        },\n        {\n                \"Address\": \"0x525b5500eE75fE2A84Cc879Cf55BB9e691A802EB\",\n                \"NodeKey\": \"0x2f09f88be6885546c2faba9d8835d772d8f2325d64273d7153432effdc78224f\",\n                \"Static\": \"enode://115ef399832f714c7723a42fb2d6707094020eb6446fdeeb72c0cb137e773984bdb4242b4c22691ed154129ef8f0216ee78d0df873a625aba156abf3628904b8@127.0.0.1:30300?discport=0\"\n        }\n]\n"
	err := DumpGenesis(filePaths, contents)
	assert.NoError(t, err)
	DeleteBasePath("./temp/")
}

func TestNodeFuncs(t *testing.T) {
	node := new(Node)
	node.Address = common.HexToAddress("0x123").String()
	node.NodeKey = "0x1231231asasd"
	node.PubKey = "0x123dsaad122"
	node.Static = "enode:1231asdas"
	nodeJson, err := node.Encode()
	assert.NoError(t, err)
	got := new(Node)
	err = got.Decode(nodeJson)
	assert.NoError(t, err)
	assert.Equal(t, got, node)
}

