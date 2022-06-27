/*
 * Copyright (C) 2022 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */
package signature_manager

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	SIG_INFO = "sigInfo"
)

func CheckSigns(native *native.NativeContract, id, sig []byte, address common.Address) (bool, error) {
	sigInfo, err := getSigInfo(native, id)
	if err != nil {
		return false, fmt.Errorf("CheckSigs, getSigInfo error: %v", err)
	}

	// get epoch info
	epochInfo, err := node_manager.GetCurrentEpochInfoImpl(native)
	if err != nil {
		return false, fmt.Errorf("CheckSigns, node_manager.GetCurrentEpochInfo error: %v", err)
	}

	ctx := native.ContractRef().CurrentContext()
	caller := ctx.Caller

	// check authority
	if err := node_manager.CheckVoterAuthority(caller, caller, epochInfo); err != nil {
		return false, fmt.Errorf("CheckSigns", "check authority failed", err)
	}

	//check signs num
	num := 0
	sum := len(epochInfo.Voters)
	for _, v := range epochInfo.Voters {
		address := v.Address
		_, ok := sigInfo.m[address.Hex()]
		if ok {
			num = num + 1
		}
	}
	//check if voted
	_, ok := sigInfo.m[address.Hex()]
	if !ok {
		sigInfo.m[address.Hex()] = sig
		num = num + 1
		if num < (2*sum+2)/3 {
			putSigInfo(native, id, sigInfo)
		}
	}
	if num >= (2*sum+2)/3 {
		shouldEmit := !sigInfo.Status
		sigInfo.Status = true
		putSigInfo(native, id, sigInfo)
		return shouldEmit, nil
	} else {
		return false, nil
	}
}

func getSigInfo(native *native.NativeContract, id []byte) (*SigInfo, error) {
	key := utils.ConcatKey(utils.SignatureManagerContractAddress, []byte(SIG_INFO), id)
	sigInfoBytes, err := native.GetCacheDB().Get(key)
	if err != nil {
		return nil, fmt.Errorf("getSigInfo, get getSigInfoStore error: %v", err)
	}

	sigInfo := &SigInfo{}
	if sigInfoBytes != nil {
		err = rlp.DecodeBytes(sigInfoBytes, sigInfo)
		if err != nil {
			return nil, fmt.Errorf("getSigInfo, deserialize SigInfo err:%v", err)
		}
	}
	return sigInfo, nil
}

func putSigInfo(native *native.NativeContract, id []byte, sigInfo *SigInfo) {
	contract := utils.SignatureManagerContractAddress

	sigInfoBytes, _ := rlp.EncodeToBytes(sigInfo)

	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(SIG_INFO), id), sigInfoBytes)
}
