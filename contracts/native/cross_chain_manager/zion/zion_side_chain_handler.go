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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/eth/types"
	"github.com/ethereum/go-ethereum/contracts/native/header_sync/quorum"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	pcom "github.com/polynetwork/poly/common"
)

type ZionHandler struct{}

func NewHandler() *ZionHandler {
	return &ZionHandler{}
}

// todo
func (this *ZionHandler) MakeDepositProposal(ns *native.NativeContract) (*common.MakeTxParam, error) {
	ctx := ns.ContractRef().CurrentContext()
	params := &common.EntranceParam{}
	if err := utils.UnpackMethod(common.ABI, common.MethodImportOuterTransfer, params, ctx.Payload); err != nil {
		return nil, err
	}

	sideChain, err := side_chain_manager.GetSideChain(ns, params.SourceChainID)
	if err != nil {
		return nil, fmt.Errorf("Zion MakeDepositProposal, side_chain_manager.GetSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, errors.New("Zion MakeDepositProposal, side chain not found")
	}

	val := &common.MakeTxParam{}
	if err := val.Deserialization(pcom.NewZeroCopySource(params.Extra)); err != nil {
		return nil, fmt.Errorf("Zion MakeDepositProposal, failed to deserialize MakeTxParam: %v", err)
	}
	if err := common.CheckDoneTx(ns, val.CrossChainID, params.SourceChainID); err != nil {
		return nil, fmt.Errorf("Zion MakeDepositProposal, check done transaction error: %v", err)
	}
	if err := common.PutDoneTx(ns, val.CrossChainID, params.SourceChainID); err != nil {
		return nil, fmt.Errorf("Zion MakeDepositProposal, PutDoneTx error: %v", err)
	}

	header := &types.Header{}
	if err := json.Unmarshal(params.HeaderOrCrossChainMsg, header); err != nil {
		return nil, fmt.Errorf("Zion MakeDepositProposal, deserialize header err: %v", err)
	}
	valh, err := quorum.GetCurrentValHeight(ns, params.SourceChainID)
	if err != nil {
		return nil, fmt.Errorf("Zion MakeDepositProposal, failed to get current validators height: %v", err)
	}
	if header.Number.Uint64() < valh {
		return nil, fmt.Errorf("Zion MakeDepositProposal, height of header %d is less than epoch height %d", header.Number.Uint64(), valh)
	}
	vs, err := quorum.GetValSet(ns, params.SourceChainID)
	if err != nil {
		return nil, fmt.Errorf("Zion MakeDepositProposal, failed to get quorum validators: %v", err)
	}
	if _, err := quorum.VerifyQuorumHeader(vs, header, false); err != nil {
		return nil, fmt.Errorf("Zion MakeDepositProposal, failed to verify quorum header %s: %v", header.Hash().String(), err)
	}

	if err := verifyFromQuorumTx(params.Proof, params.Extra, header, sideChain); err != nil {
		return nil, fmt.Errorf("Zion MakeDepositProposal, verifyFromEthTx error: %s", err)
	}

	return val, nil
}
