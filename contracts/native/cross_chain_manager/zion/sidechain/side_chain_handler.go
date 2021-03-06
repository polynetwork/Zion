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

package sidechain

import (
	"encoding/json"
	"fmt"

	ecom "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	xutils "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/utils"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/zion"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) MakeDepositProposal(s *native.NativeContract) (*common.MakeTxParam, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &common.EntranceParam{}
	if err := utils.UnpackMethod(common.ABI, common.MethodImportOuterTransfer, params, ctx.Payload); err != nil {
		return nil, err
	}

	log.Trace("ZionSideChainHandler", "MakeDepositProposal", params.String())

	sideChain, err := side_chain_manager.GetSideChain(s, params.SourceChainID)
	if err != nil {
		return nil, fmt.Errorf("ZionSideChainHandler MakeDepositProposal, side_chain_manager.GetSideChain err: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("ZionSideChainHandler MakeDepositProposal, side chain not found")
	}

	val, err := common.DecodeTxParam(params.Extra)
	if err != nil {
		return nil, fmt.Errorf("ZionSideChainHandler MakeDepositProposal, failed to deserialize MakeTxParam: %v", err)
	}
	if err := common.CheckDoneTx(s, val.CrossChainID, params.SourceChainID); err != nil {
		return nil, fmt.Errorf("ZionSideChainHandler MakeDepositProposal, check done transaction error: %v", err)
	}

	header := &types.Header{}
	if err := json.Unmarshal(params.HeaderOrCrossChainMsg, header); err != nil {
		return nil, fmt.Errorf("ZionSideChainHandler MakeDepositProposal, deserialize header err: %v", err)
	}

	curEpochStartHeight, curEpochValidators, err := zion.GetEpoch(s, params.SourceChainID)
	if err != nil {
		return nil, fmt.Errorf("ZionSideChainHandler MakeDepositProposal, failed to get current validators height: %v", err)
	}
	if header.Number.Uint64() < curEpochStartHeight {
		return nil, fmt.Errorf("ZionSideChainHandler MakeDepositProposal, height of header %d is less than epoch height %d", header.Number.Uint64(), curEpochStartHeight)
	}

	if _, _, err := zion.VerifyHeader(header, curEpochValidators, false); err != nil {
		return nil, fmt.Errorf("ZionSideChainHandler MakeDepositProposal, failed to verify quorum header %s: %v", header.Hash().String(), err)
	}

	if _, err := xutils.VerifyTx(params.Proof, header, ecom.BytesToAddress(sideChain.CCMCAddress), params.Extra, true); err != nil {
		return nil, fmt.Errorf("ZionSideChainHandler MakeDepositProposal, verifyFromEthTx err: %v", err)
	}

	if err := common.PutDoneTx(s, val.CrossChainID, params.SourceChainID); err != nil {
		return nil, fmt.Errorf("ZionSideChainHandler MakeDepositProposal, PutDoneTx error: %v", err)
	}

	return val, nil
}
