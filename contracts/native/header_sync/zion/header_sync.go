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
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	scom "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) SyncGenesisHeader(s *native.NativeContract) error {
	ctx := s.ContractRef().CurrentContext()
	msgSender := s.ContractRef().MsgSender()

	// main chain DONT need sync genesis header
	params := &scom.SyncGenesisHeaderParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodSyncGenesisHeader, params, ctx.Payload); err != nil {
		return fmt.Errorf("ZionHandler SyncGenesisHeader, contract params deserialize err: %v", err)
	}
	chainID := params.ChainID

	// Get current epoch operator
	ok, err := node_manager.CheckConsensusSigns(s, scom.MethodSyncGenesisHeader, ctx.Payload, msgSender)
	if err != nil {
		return fmt.Errorf("ZionHandler SyncGenesisHeader, CheckConsensusSigns err: %v", err)
	}
	if !ok {
		return nil
	}

	var header *types.Header
	if err := json.Unmarshal(params.GenesisHeader, &header); err != nil {
		return fmt.Errorf("ZionHandler SyncGenesisHeader, json.Unmarshal header err: %v", err)
	}

	if isGenesisStored(s, chainID) {
		return fmt.Errorf("ZionHandler SyncGenesisHeader, genesis header had been initialized")
	}

	//block header storage
	if err = storeGenesis(s, chainID, header); err != nil {
		return fmt.Errorf("ZionHandler SyncGenesisHeader, put blockHeader err: %v", err)
	}

	validators, err := getValidatorsFromHeader(header)
	if err != nil {
		return fmt.Errorf("ZionHandler SyncGenesisHeader, get validators from header err: %v", err)
	}

	height := header.Number.Uint64()
	if height != 0 {
		height += 1
	}
	if err := storeEpoch(s, chainID, height, validators); err != nil {
		return fmt.Errorf("ZionHandler SyncGenesisHeader, store epoch err: %v", err)
	}

	emitEpochChangeEvent(s, chainID, height, header.Hash())

	log.Debug("ZionHandler SyncGenesisHeader", "chainID", chainID, "height", height, "hash", header.Hash())
	return nil
}

func (h *Handler) SyncBlockHeader(s *native.NativeContract) error {
	ctx := s.ContractRef().CurrentContext()
	params := &scom.SyncBlockHeaderParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodSyncBlockHeader, params, ctx.Payload); err != nil {
		return err
	}

	if params.Headers == nil || len(params.Headers) < 2 {
		return fmt.Errorf("invalid params")
	}

	chainID := params.ChainID
	curEpochStartHeight, curEpochValidators, err := getEpoch(s, chainID)
	if err != nil {
		return fmt.Errorf("ZionHandler SynnBlockHeader, failed to get current epoch info, err: %v", err)
	}
	lastHeaderHeight := uint64(0)
	initStartHeight := curEpochStartHeight

	for i, v := range params.Headers {
		header := new(types.Header)
		if err := header.UnmarshalJSON(v); err != nil {
			return fmt.Errorf("ZionHandler SyncBlockHeader, deserialize No.%d header err: %v", i, err)
		}

		// check height
		h := header.Number.Uint64()
		if initStartHeight >= h {
			return fmt.Errorf("ZionHandler SyncBlockHeader, wrong height of No.%d header: (curr: %d, commit: %d)", i, initStartHeight, h)
		}
		if lastHeaderHeight > 0 && h != lastHeaderHeight+1 {
			return fmt.Errorf("ZionHandler SyncBlockHeader, should be continues block headers")
		} else {
			lastHeaderHeight = h
		}

		// validate header and epoch
		nextEpochStartHeight, nextEpochValidators, err := VerifyHeader(header, curEpochValidators, true)
		if err != nil {
			return fmt.Errorf("ZionHandler SyncBlockHeader, failed to verify No.%d quorum header %s: %v", i, header.Hash().Hex(), err)
		}

		// epoch changed
		if nextEpochStartHeight > 0 && len(nextEpochValidators) > 0 && nextEpochStartHeight > curEpochStartHeight {
			emitEpochChangeEvent(s, chainID, header.Number.Uint64(), header.Hash())
			if err := storeEpoch(s, chainID, nextEpochStartHeight, nextEpochValidators); err != nil {
				return fmt.Errorf("ZionHandler SyncBlockHeader, failed to store next epoch info, err: %v", err)
			}

			log.Debug("ZionHandler SyncBlockHeader", "chainID", chainID, "height", header.Number, "hash", header.Hash(),
				"current epoch start height", curEpochStartHeight, "current epoch validators", curEpochValidators,
				"next epoch start height", nextEpochStartHeight, "next epoch validators", nextEpochValidators)

			curEpochStartHeight = nextEpochStartHeight
			curEpochValidators = nextEpochValidators
		}
	}

	return nil
}

// todo(fuk): useless interface
func (h *Handler) SyncCrossChainMsg(native *native.NativeContract) error {
	return nil
}
