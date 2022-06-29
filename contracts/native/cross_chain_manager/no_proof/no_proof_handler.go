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

package no_proof

import (
	"fmt"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
)

type NoProofHandler struct {
}

func NewNoProofHandler() *NoProofHandler {
	return &NoProofHandler{}
}

func (this *NoProofHandler) MakeDepositProposal(service *native.NativeContract) (*scom.MakeTxParam, error) {
	ctx := service.ContractRef().CurrentContext()
	params := &scom.EntranceParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodImportOuterTransfer, params, ctx.Payload); err != nil {
		return nil, err
	}

	//verify signature
	digest, err := params.Digest()
	if err != nil {
		return nil, fmt.Errorf("no proof MakeDepositProposal, digest input param error: %v", err)
	}
	if !crypto.VerifySignature(params.Pub, digest, params.Signature) {
		return nil, fmt.Errorf("no proof MakeDepositProposal, verify input signature error: %v", err)
	}
	pub, err := crypto.DecompressPubkey(params.Pub)
	if err != nil {
		return nil, fmt.Errorf("no proof MakeDepositProposal, crypto.DecompressPubkey error: %v", err)
	}
	addr := crypto.PubkeyToAddress(*pub)

	ok, err := node_manager.CheckVoterSigns(service, scom.MethodImportOuterTransfer, digest, addr)
	if err != nil {
		return nil, fmt.Errorf("no proof MakeDepositProposal, CheckVoterSigns error: %v", err)
	}

	if ok {
		makeTxParam, err := scom.DecodeTxParam(params.Extra)
		if err != nil {
			return nil, fmt.Errorf("no proof MakeDepositProposal, deserialize MakeTxParam error:%s", err)
		}
		if err := scom.CheckDoneTx(service, makeTxParam.CrossChainID, params.SourceChainID); err != nil {
			return nil, fmt.Errorf("no proof MakeDepositProposal, check done transaction error:%s", err)
		}
		if err := scom.PutDoneTx(service, makeTxParam.CrossChainID, params.SourceChainID); err != nil {
			return nil, fmt.Errorf("no proof MakeDepositProposal, PutDoneTx error:%s", err)
		}
		return makeTxParam, nil
	}
	return nil, nil
}
