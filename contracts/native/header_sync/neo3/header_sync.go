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

package neo3

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	hscommon "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/joeqian10/neo3-gogogo/helper"

	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/polynetwork/poly/common"
)

type Neo3Handler struct {
}

func NewNeo3Handler() *Neo3Handler {
	return &Neo3Handler{}
}

func (this *Neo3Handler) SyncGenesisHeader(native *native.NativeContract) error {
	ctx := native.ContractRef().CurrentContext()
	params := &hscommon.SyncGenesisHeaderParam{}
	if err := utils.UnpackMethod(hscommon.ABI, hscommon.MethodSyncGenesisHeader, params, ctx.Payload); err != nil {
		return fmt.Errorf("SyncGenesisHeader, contract params deserialize error: %v", err)
	}
	// Get current epoch operator
	ok, err := node_manager.CheckConsensusSigns(native, hscommon.MethodSyncGenesisHeader, ctx.Payload, native.ContractRef().MsgSender())
	if err != nil {
		return fmt.Errorf("SyncGenesisHeader, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return nil
	}
	// Deserialize neo block header
	var header NeoBlockHeader
	err = json.Unmarshal(params.GenesisHeader, &header)
	if err != nil {
		return fmt.Errorf("SyncGenesisHeader, json.Unmarshal header err: %v", err)
	}

	if neoConsensus, _ := getConsensusValByChainId(native, params.ChainID); neoConsensus == nil {
		// Put NeoConsensus.NextConsensus into storage
		if err = putConsensusValByChainId(native, &NeoConsensus{
			ChainID:       params.ChainID,
			Height:        header.GetIndex(),
			NextConsensus: header.GetNextConsensus(),
		}); err != nil {
			return fmt.Errorf("Neo3Handler SyncGenesisHeader, update ConsensusPeer error: %v", err)
		}
	}
	return nil
}

func (this *Neo3Handler) SyncBlockHeader(native *native.NativeContract) error {
	params := &hscommon.SyncBlockHeaderParam{}
	{
		ctx := native.ContractRef().CurrentContext()
		if err := utils.UnpackMethod(hscommon.ABI, hscommon.MethodSyncBlockHeader, params, ctx.Payload); err != nil {
			return err
		}
	}
	neoConsensus, err := getConsensusValByChainId(native, params.ChainID)
	if err != nil {
		return fmt.Errorf("Neo3Handler SyncBlockHeader, the consensus validator has not been initialized, chainId: %d", params.ChainID)
	}
	sideChain, err := side_chain_manager.GetSideChain(native, params.ChainID)
	if err != nil {
		return fmt.Errorf("neo3 MakeDepositProposal, side_chain_manager.GetSideChain error: %v", err)
	}
	var newNeoConsensus *NeoConsensus
	for _, v := range params.Headers {
		header := new(NeoBlockHeader)
		if err := header.Deserialization(common.NewZeroCopySource(v)); err != nil {
			return fmt.Errorf("Neo3Handler SyncBlockHeader, NeoBlockHeaderFromBytes error: %v", err)
		}
		if !header.GetNextConsensus().Equals(neoConsensus.NextConsensus) && header.GetIndex() > neoConsensus.Height {
			if err = verifyHeader(native, params.ChainID, header, helper.BytesToUInt32(sideChain.ExtraInfo)); err != nil {
				return fmt.Errorf("Neo3Handler SyncBlockHeader, verifyHeader error: %v", err)
			}
			newNeoConsensus = &NeoConsensus{
				ChainID:       neoConsensus.ChainID,
				Height:        header.GetIndex(),
				NextConsensus: header.GetNextConsensus(),
			}
		}
	}
	if newNeoConsensus != nil {
		if err = putConsensusValByChainId(native, newNeoConsensus); err != nil {
			return fmt.Errorf("Neo3Handler SyncBlockHeader, update ConsensusPeer error: %v", err)
		}
	}
	return nil
}

func (this *Neo3Handler) SyncCrossChainMsg(native *native.NativeContract) error {
	return nil
}
