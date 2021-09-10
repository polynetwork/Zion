package tool

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/hotstuff"
	"github.com/ethereum/go-ethereum/consensus/hotstuff/validator"
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

func TestEncodeSalt(t *testing.T) {
	validators := []common.Address{
		common.HexToAddress("0xc095448424a5ecd5ca7ccdadfaad127a9d7e88ec"),
		common.HexToAddress("0xd47a4e56e9262543db39d9203cf1a2e53735f834"),
		common.HexToAddress("0x258af48e28e4a6846e931ddff8e1cdf8579821e5"),
		common.HexToAddress("0x8c09d936a1b408d6e0afaa537ba4e06c4504a0ae"),
	}
	valset := validator.NewSet(validators, hotstuff.RoundRobin)
	validators = valset.AddressList()
	enc, err := Encode(validators)
	assert.NoError(t, err)
	t.Log(enc)
}