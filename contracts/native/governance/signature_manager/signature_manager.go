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
	"crypto/sha256"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/contract"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/signature_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

var (
	this     = native.NativeContractAddrMap[native.NativeSignatureManager]
	gasTable = map[string]uint64{
		// MethodContractName:             0,
		signature_manager_abi.MethodAddSignature: 100000,
	}

	ABI *abi.ABI
)

func InitSignatureManager() {
	ABI = GetABI()
	native.Contracts[this] = RegisterSignatureManagerContract
}

//Register methods of signature_manager contract
func RegisterSignatureManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(signature_manager_abi.MethodAddSignature, AddSignature)

}

func AddSignature(s *native.NativeContract) ([]byte, error) {

	ctx := s.ContractRef().CurrentContext()
	params := &AddSignatureParam{}
	if err := utils.UnpackMethod(ABI, signature_manager_abi.MethodAddSignature, params, ctx.Payload); err != nil {
		return nil, err
	}
	//check witness
	if err := contract.ValidateOwner(s, params.Address); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("AddSignature, checkWitness: %s, error: %v", params.Address, err)
	}

	temp := sha256.Sum256(params.Subject)
	id := temp[:]
	//check consensus signs
	ok, err := CheckSigns(s, id, params.Signature, params.Address)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("AddSignature, CheckSigns error: %v", err)
	}
	if !ok {
		return utils.BYTE_TRUE, nil
	}

	s.AddNotify(
		ABI, []string{signature_manager_abi.EventAddSignatureQuorumEvent}, id, params.Subject, params.SideChainID)
	return utils.BYTE_TRUE, nil

}
