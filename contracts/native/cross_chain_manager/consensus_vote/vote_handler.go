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

package consensus_vote

import (
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	polycomm "github.com/polynetwork/poly/common"
)

type VoteHandler struct {
}

func NewVoteHandler() *VoteHandler {
	return &VoteHandler{}
}

func (this *VoteHandler) MakeDepositProposal(service *native.NativeContract) (*scom.MakeTxParam, error) {
	ctx := service.ContractRef().CurrentContext()
	params := &scom.EntranceParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodImportOuterTransfer, params, ctx.Payload); err != nil {
		return nil, err
	}

	//use sourcechainid, height, extra as unique id
	unique := &scom.EntranceParam{
		SourceChainID: params.SourceChainID,
		Height:        params.Height,
		Extra:         params.Extra,
	}
	sink := polycomm.NewZeroCopySink(nil)
	unique.Serialization(sink)

	ok, err := CheckConsensusSigns(service, sink.Bytes())
	if err != nil {
		return nil, fmt.Errorf("vote MakeDepositProposal, CheckConsensusSigns error: %v", err)
	}
	if ok {
		txParam, err := scom.DecodeTxParam(params.Extra)
		if err != nil {
			return nil, fmt.Errorf("vote MakeDepositProposal, deserialize MakeTxParam error:%s", err)
		}
		if err := scom.CheckDoneTx(service, txParam.CrossChainID, params.SourceChainID); err != nil {
			return nil, fmt.Errorf("vote MakeDepositProposal, check done transaction error:%s", err)
		}
		if err := scom.PutDoneTx(service, txParam.CrossChainID, params.SourceChainID); err != nil {
			return nil, fmt.Errorf("vote MakeDepositProposal, PutDoneTx error:%s", err)
		}
		return txParam, nil
	}
	return nil, nil
}
