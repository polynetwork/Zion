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

package zion

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestVerifyHeader(t *testing.T) {
	rawHeader := `{"parentHash":"0x8adbb7aa118074c58ce20966b19734a4a0cfa2898f5bfcd01b086b068351ff5a","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","miner":"0x8c09d936a1b408d6e0afaa537ba4e06c4504a0ae","stateRoot":"0x8a5364fbb3e7d3c5076c9c887d84e9569d277975ab9b428046a311a60543639b","transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","difficulty":"0x1","number":"0x44b","gasLimit":"0x577ae927","gasUsed":"0x0","timestamp":"0x6177784f","extraData":"0x0000000000000000000000000000000000000000000000000000000000000000f901e8f89394258af48e28e4a6846e931ddff8e1cdf8579821e5946a708455c8777630aac9d1e7702d13f7a865b27c948c09d936a1b408d6e0afaa537ba4e06c4504a0ae94ad3bf5ed640cc72f37bd21d64a65c3c756e9c88c94c095448424a5ecd5ca7ccdadfaad127a9d7e88ec94d47a4e56e9262543db39d9203cf1a2e53735f83494bfb558f0dceb07fbb09e1c283048b551a4310921b8418aece0db7a7534f6d8f3b2b49c5bd5aae27c20fe8637f57670c24c8255a7119d1d70eb279d1dceb397514f6c7182981aad0b20febd18d00653056fb837d205cc01f9010cb841e065730b6ab96d2172eb9c34c7b89dcec57a4933ffdc0e941ac0beb91a69861604007302525252a94919e0f302830cd444ab43b52f593a8a287afc82aa5ef08f00b841e3696e39bb78b65fb9a9af5960e95147a8a2b54fdb78876a1ac5adeb2d09497c33bbf05bd0cdcc921e16ab66a3b409268b551ee9f805db587e19557d3e5bad7701b8417ed6a92254435c47ce41c8119621a1b9a2e1eda12241912c8127b5fd3ee76d9d5b9326dc50ae6e4f9d142086a7eb4b99f5bdd270d86d22adbe4591c33c6f5c4700b8418aece0db7a7534f6d8f3b2b49c5bd5aae27c20fe8637f57670c24c8255a7119d1d70eb279d1dceb397514f6c7182981aad0b20febd18d00653056fb837d205cc0180","mixHash":"0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365","nonce":"0x0000000000000000","hash":"0xefd31bec3a0056825006b3c68d1589b5c5ba2a5a19fe02360e84236ea6682415"}`
	header := new(types.Header)
	assert.NoError(t, header.UnmarshalJSON([]byte(rawHeader)))

	/*
		peer: 0x258af48e28E4A6846E931dDfF8e1Cdf8579821e5, pubkey: 0x02c07fb7d48eac559a2483e249d27841c18c7ce5dbbbf2796a6963cc9cef27cabd
		peer: 0x6A708455C8777630AaC9d1e7702d13F7a865b27C, pubkey: 0x02f5135ae0853af71f017a8ecb68e720b729ab92c7123c686e75b7487d4a57ae07
		peer: 0x8c09D936a1B408D6e0afAA537ba4E06c4504a0AE, pubkey: 0x03ecac0ebe7224cfd04056c940605a4a9d4cb0367cf5819bf7e5502bf44f68bdd4
		peer: 0xAd3Bf5eD640CC72f37BD21d64A65C3C756e9C88C, pubkey: 0x03d0ecfd09db6b1e4f59da7ebde8f6c3ea3ed09f06f5190477ae4ee528ec692fa8
		peer: 0xC095448424A5ECd5cA7CcDaDFaAD127a9d7E88ec, pubkey: 0x0244e509103445d5e8fd290608308d16d08c739655d6994254e413bc1a06783856
		peer: 0xD47a4e56e9262543Db39d9203CF1a2e53735f834, pubkey: 0x023884de29148505a8d862992e5721767d4b47ff52ffab4c2d2527182d812a6d95
		peer: 0xbfB558F0dceb07Fbb09E1C283048B551A4310921, pubkey: 0x03b838fa2387beb3a56aed86e447309f8844cb208387c63af64ad740729b5c0a27
	*/
	valsets := []common.Address{
		common.HexToAddress("0x258af48e28E4A6846E931dDfF8e1Cdf8579821e5"),
		common.HexToAddress("0x6A708455C8777630AaC9d1e7702d13F7a865b27C"),
		common.HexToAddress("0x8c09D936a1B408D6e0afAA537ba4E06c4504a0AE"),
		common.HexToAddress("0xAd3Bf5eD640CC72f37BD21d64A65C3C756e9C88C"),
		common.HexToAddress("0xC095448424A5ECd5cA7CcDaDFaAD127a9d7E88ec"),
	}
	nextEpochStartHeight, nextEpochVals, err := VerifyHeader(header, valsets, true)
	assert.NoError(t, err)
	t.Logf("next epoch start height %d", nextEpochStartHeight)
	t.Logf("next epoch validators %v", nextEpochVals)
}
