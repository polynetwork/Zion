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

package lock_proxy

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	zutils "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/utils"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	midGetEthCrossChainManager = crypto.Keccak256(utils.EncodePacked([]byte("getEthCrossChainManager"), []byte("()")))[:4]
	midCrossChain              = crypto.Keccak256(utils.EncodePacked([]byte("crossChain"), []byte("(uint64,bytes,bytes,bytes)")))[:4]
	argsCrossChain             = abi.Arguments{
		{Type: zutils.Uint64Ty, Name: "_toChainId"},
		{Type: zutils.BytesTy, Name: "_toContract"},
		{Type: zutils.BytesTy, Name: "_method"},
		{Type: zutils.BytesTy, Name: "_txData"},
	}
)

func getEthCrossChainManager(s *native.NativeContract, ccmp common.Address) (common.Address, error) {
	gas := s.ContractRef().GasLeft()
	payload := midGetEthCrossChainManager
	if ret, _, err := s.ContractRef().EVMCall(this, ccmp, gas, payload); err != nil {
		return common.EmptyAddress, err
	} else {
		return common.BytesToAddress(ret), nil
	}
}

func crossChain(s *native.NativeContract, eccm, toProxy common.Address, toChainID uint64, method string, txData []byte) error {
	callData, err := argsCrossChain.Pack(toChainID, toProxy[:], []byte(method), txData)
	if err != nil {
		return err
	}
	payload := utils.EncodePacked(midCrossChain, callData)

	gas := s.ContractRef().GasLeft()
	if _, _, err := s.ContractRef().EVMCall(this, eccm, gas, payload); err != nil {
		return err
	}
	return nil
}

// compareVals return true if `src` equals to `cmp`
func compareVals(v1, v2 []common.Address) bool {
	exist := func(addr common.Address, list []common.Address) bool {
		for _, v := range list {
			if addr == v {
				return true
			}
		}
		return false
	}

	contain := func(l1, l2 []common.Address) bool {
		for _, v := range l1 {
			if !exist(v, l2) {
				return false
			}
		}
		return true
	}

	if !contain(v1, v2) {
		return false
	}
	if !contain(v2, v1) {
		return false
	}

	return true
}
