/*
 * Copyright (C) 2021 The poly network Authors
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
package quorum

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/eth/types"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	common2 "github.com/ethereum/go-ethereum/contracts/native/info_sync"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

type QuorumHandler struct{}

func NewQuorumHandler() *QuorumHandler {
	return &QuorumHandler{}
}

func (this *QuorumHandler) MakeDepositProposal(ns *native.NativeContract) (*common.MakeTxParam, error) {
	ctx := ns.ContractRef().CurrentContext()
	params := &common.EntranceParam{}
	if err := utils.UnpackMethod(common.ABI, common.MethodImportOuterTransfer, params, ctx.Payload); err != nil {
		return nil, err
	}

	sideChain, err := side_chain_manager.GetSideChain(ns, params.SourceChainID)
	if err != nil {
		return nil, fmt.Errorf("Quorum MakeDepositProposal, side_chain_manager.GetSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, errors.New("Quorum MakeDepositProposal, side chain not found")
	}

	val, err := scom.DecodeTxParam(params.Extra)
	if err != nil {
		return nil, fmt.Errorf("Quorum MakeDepositProposal, failed to deserialize MakeTxParam: %v", err)
	}
	if err := common.CheckDoneTx(ns, val.CrossChainID, params.SourceChainID); err != nil {
		return nil, fmt.Errorf("Quorum MakeDepositProposal, check done transaction error: %v", err)
	}
	if err := common.PutDoneTx(ns, val.CrossChainID, params.SourceChainID); err != nil {
		return nil, fmt.Errorf("Quorum MakeDepositProposal, PutDoneTx error: %v", err)
	}

	value, err := common2.GetRootInfo(ns, params.SourceChainID, params.Height)
	if err != nil {
		return nil, fmt.Errorf("Quorum MakeDepositProposal, GetCrossChainInfo error:%s", err)
	}
	header := &types.Header{}
	err = json.Unmarshal(value, header)
	if err != nil {
		return nil, fmt.Errorf("Quorum MakeDepositProposal, json unmarshal header error:%s", err)
	}

	if err := verifyFromQuorumTx(params.Proof, params.Extra, header, sideChain); err != nil {
		return nil, fmt.Errorf("Quorum MakeDepositProposal, verifyFromEthTx error: %s", err)
	}

	return val, nil
}
