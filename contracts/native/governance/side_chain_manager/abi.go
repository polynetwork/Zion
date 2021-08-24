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
package side_chain_manager

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	polycomm "github.com/polynetwork/poly/common"
)

const (
	EventRegisterSideChain        = "registerSideChain"
	EventApproveRegisterSideChain = "approveRegisterSideChain"
	EventUpdateSideChain          = "updateSideChain"
	EventApproveUpdateSideChain   = "approveUpdateSideChain"
	EventQuitSideChain            = "quitSideChain"
	EventApproveQuitSideChain     = "approveQuitSideChain"
	EventRegisterRedeem           = "registerRedeem"
	EventSetBtcTxParam            = "setBtcTxParam"
)

const abijson = `[
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"ChainId","type":"uint64"}],"name":"` + EventApproveQuitSideChain + `","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"ChainId","type":"uint64"}],"name":"` + EventApproveRegisterSideChain + `","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"ChainId","type":"uint64"}],"name":"` + EventApproveUpdateSideChain + `","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"ChainId","type":"uint64"}],"name":"` + EventQuitSideChain + `","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"rk","type":"string"},{"indexed":false,"internalType":"string","name":"ContractAddress","type":"string"}],"name":"` + EventRegisterRedeem + `","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"ChainId","type":"uint64"},{"indexed":false,"internalType":"uint64","name":"Router","type":"uint64"},{"indexed":false,"internalType":"string","name":"Name","type":"string"},{"indexed":false,"internalType":"uint64","name":"BlocksToWait","type":"uint64"}],"name":"` + EventRegisterSideChain + `","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"rk","type":"string"},{"indexed":false,"internalType":"uint64","name":"RedeemChainId","type":"uint64"},{"indexed":false,"internalType":"uint64","name":"FeeRate","type":"uint64"},{"indexed":false,"internalType":"uint64","name":"MinChange","type":"uint64"}],"name":"` + EventSetBtcTxParam + `","type":"event"},
    {"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"ChainId","type":"uint64"},{"indexed":false,"internalType":"uint64","name":"Router","type":"uint64"},{"indexed":false,"internalType":"string","name":"Name","type":"string"},{"indexed":false,"internalType":"uint64","name":"BlocksToWait","type":"uint64"}],"name":"` + EventUpdateSideChain + `","type":"event"},
    {"inputs":[{"internalType":"uint64","name":"Chainid","type":"uint64"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodApproveQuitSideChain + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"uint64","name":"Chainid","type":"uint64"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodApproveRegisterSideChain + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"uint64","name":"Chainid","type":"uint64"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodApproveUpdateSideChain + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"uint64","name":"Chainid","type":"uint64"},{"internalType":"address","name":"Address","type":"address"}],"name":"` + MethodQuitSideChain + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"uint64","name":"RedeemChainID","type":"uint64"},{"internalType":"uint64","name":"ContractChainID","type":"uint64"},{"internalType":"bytes","name":"Redeem","type":"bytes"},{"internalType":"uint64","name":"CVersion","type":"uint64"},{"internalType":"bytes","name":"ContractAddress","type":"bytes"},{"internalType":"bytes[]","name":"Signs","type":"bytes[]"}],"name":"` + MethodRegisterRedeem + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"address","name":"Address","type":"address"},{"internalType":"uint64","name":"ChainId","type":"uint64"},{"internalType":"uint64","name":"Router","type":"uint64"},{"internalType":"string","name":"Name","type":"string"},{"internalType":"uint64","name":"BlocksToWait","type":"uint64"},{"internalType":"bytes","name":"CCMCAddress","type":"bytes"},{"internalType":"bytes","name":"ExtraInfo","type":"bytes"}],"name":"` + MethodRegisterSideChain + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"bytes","name":"Redeem","type":"bytes"},{"internalType":"uint64","name":"RedeemChainId","type":"uint64"},{"internalType":"bytes[]","name":"Sigs","type":"bytes[]"},{"components":[{"internalType":"uint64","name":"PVersion","type":"uint64"},{"internalType":"uint64","name":"FeeRate","type":"uint64"},{"internalType":"uint64","name":"MinChange","type":"uint64"}],"internalType":"struct side_chain_manager.BtcTxParamDetial","name":"Detial","type":"tuple"}],"name":"` + MethodSetBtcTxParam + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"internalType":"address","name":"Address","type":"address"},{"internalType":"uint64","name":"ChainId","type":"uint64"},{"internalType":"uint64","name":"Router","type":"uint64"},{"internalType":"string","name":"Name","type":"string"},{"internalType":"uint64","name":"BlocksToWait","type":"uint64"},{"internalType":"bytes","name":"CCMCAddress","type":"bytes"},{"internalType":"bytes","name":"ExtraInfo","type":"bytes"}],"name":"` + MethodUpdateSideChain + `","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"}
]`

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(abijson))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

type RegisterSideChainParam struct {
	Address      common.Address
	ChainId      uint64
	Router       uint64
	Name         string
	BlocksToWait uint64
	CCMCAddress  []byte
	ExtraInfo    []byte
}

type ChainidParam struct {
	Chainid uint64
	Address common.Address
}

type RegisterRedeemParam struct {
	RedeemChainID   uint64
	ContractChainID uint64
	Redeem          []byte
	CVersion        uint64
	ContractAddress []byte
	Signs           [][]byte
}

type BtcTxParam struct {
	Redeem        []byte
	RedeemChainId uint64
	Sigs          [][]byte
	Detial        *BtcTxParamDetial
}

type BtcTxParamDetial struct {
	PVersion  uint64
	FeeRate   uint64
	MinChange uint64
}

func (this *BtcTxParamDetial) Serialization(sink *polycomm.ZeroCopySink) {
	sink.WriteVarUint(this.PVersion)
	sink.WriteVarUint(this.FeeRate)
	sink.WriteVarUint(this.MinChange)
}

func (this *BtcTxParamDetial) Deserialization(source *polycomm.ZeroCopySource) error {
	var eof bool
	this.PVersion, eof = source.NextVarUint()
	if eof {
		return fmt.Errorf("BtcTxParamDetial deserialize version error")
	}
	this.FeeRate, eof = source.NextVarUint()
	if eof {
		return fmt.Errorf("BtcTxParamDetial deserialize fee rate error")
	}
	this.MinChange, eof = source.NextVarUint()
	if eof {
		return fmt.Errorf("BtcTxParamDetial deserialize min-change error")
	}
	return nil
}
