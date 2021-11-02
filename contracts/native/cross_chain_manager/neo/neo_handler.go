/*
 * Copyright (C) 2020 The poly network Authors
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

package neo

import (
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/neo"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/polynetwork/poly/common"
)

type NEOHandler struct {
}

func NewNEOHandler() *NEOHandler {
	return &NEOHandler{}
}

func (this *NEOHandler) MakeDepositProposal(service *native.NativeContract) (*scom.MakeTxParam, error) {
	ctx := service.ContractRef().CurrentContext()
	params := &scom.EntranceParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodImportOuterTransfer, params, ctx.Payload); err != nil {
		return nil, err
	}
	// Deserialize neo cross chain msg and verify its signature
	crossChainMsg := new(neo.NeoCrossChainMsg)
	if err := crossChainMsg.Deserialization(common.NewZeroCopySource(params.HeaderOrCrossChainMsg)); err != nil {
		return nil, fmt.Errorf("neo MakeDepositProposal, deserialize crossChainMsg error: %v", err)
	}
	if err := neo.VerifyCrossChainMsgSig(service, params.SourceChainID, crossChainMsg); err != nil {
		return nil, fmt.Errorf("neo MakeDepositProposal, VerifyCrossChainMsg error: %v", err)
	}
	// Verify the validity of proof with the help of state root in verified neo cross chain msg
	sideChain, err := side_chain_manager.GetSideChain(service, params.SourceChainID)
	if err != nil {
		return nil, fmt.Errorf("neo MakeDepositProposal, side_chain_manager.GetSideChain error: %v", err)
	}
	value, err := verifyFromNeoTx(params.Proof, crossChainMsg, sideChain.CCMCAddress)
	if err != nil {
		return nil, fmt.Errorf("neo MakeDepositProposal, VerifyFromNeoTx error: %v", err)
	}
	// Ensure the tx has not been processed before, and mark the tx as processed
	if err := scom.CheckDoneTx(service, value.CrossChainID, params.SourceChainID); err != nil {
		return nil, fmt.Errorf("neo MakeDepositProposal, check done transaction error: %s", err)
	}
	if err = scom.PutDoneTx(service, value.CrossChainID, params.SourceChainID); err != nil {
		return nil, fmt.Errorf("neo MakeDepositProposal, putDoneTx error: %s", err)
	}
	return value, nil
}
