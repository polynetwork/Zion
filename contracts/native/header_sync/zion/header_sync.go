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
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	scom "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
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

	if isGenesisStored(s, chainID) {
		return fmt.Errorf("ZionHandler SyncGenesisHeader, genesis header had been initialized")
	}

	header := new(types.Header)
	if err := header.UnmarshalJSON(params.GenesisHeader); err != nil {
		return fmt.Errorf("ZionHandler SyncGenesisHeader, json.Unmarshal header err: %v", err)
	}

	//block header storage
	if err = storeGenesis(s, chainID, header); err != nil {
		return fmt.Errorf("ZionHandler SyncGenesisHeader, put blockHeader err: %v", err)
	}

	hash := header.Hash()
	height := header.Number.Uint64() + 1
	validators, err := getValidatorsFromHeader(header)
	if err != nil {
		return fmt.Errorf("ZionHandler SyncGenesisHeader, get validators from header err: %v", err)
	}

	if err := storeEpoch(s, chainID, height, validators); err != nil {
		return fmt.Errorf("ZionHandler SyncGenesisHeader, store epoch err: %v", err)
	}

	emitEpochChangeEvent(s, chainID, height, hash)

	log.Debug("ZionHandler SyncGenesisHeader", "chainID", chainID, "height", height, "hash", hash)
	return nil
}

func (h *Handler) SyncBlockHeader(s *native.NativeContract) error {
	ctx := s.ContractRef().CurrentContext()
	params := &scom.SyncBlockHeaderParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodSyncBlockHeader, params, ctx.Payload); err != nil {
		return err
	}

	chainID := params.ChainID
	curEpochStartHeight, curEpochValidators, err := getEpoch(s, chainID)
	if err != nil {
		return fmt.Errorf("ZionHandler SynnBlockHeader, failed to get current epoch info, err: %v", err)
	}

	for i, v := range params.Headers {
		hd := new(types.Header)
		if err := hd.UnmarshalJSON(v); err != nil {
			return fmt.Errorf("ZionHandler SyncBlockHeader, deserialize No.%d header err: %v", i, err)
		}

		// check height
		hash := hd.Hash()
		h := hd.Number.Uint64()

		if curEpochStartHeight >= h {
			continue
		}

		nextEpochStartHeight, nextEpochValidators, err := VerifyHeader(hd, curEpochValidators, true)
		if err != nil {
			return fmt.Errorf("ZionHandler SyncBlockHeader, verify No.%d header err: %v", i, err)
		}

		if err := storeEpoch(s, chainID, nextEpochStartHeight, nextEpochValidators); err != nil {
			return fmt.Errorf("ZionHandler SyncBlockHeader, store No.%d epoch err: %v", i, err)
		}

		emitEpochChangeEvent(s, chainID, h, hash)
		log.Debug("ZionHandler SyncBlockHeader", "chainID", chainID, "height", h, "hash", hash,
			"current epoch start height", curEpochStartHeight, "current epoch validators", curEpochValidators,
			"next epoch start height", nextEpochStartHeight, "next epoch validators", nextEpochValidators)

		curEpochStartHeight = nextEpochStartHeight
		curEpochValidators = nextEpochValidators
	}

	return nil
}

// todo(fuk): useless interface
func (h *Handler) SyncCrossChainMsg(native *native.NativeContract) error {
	return nil
}

func checkProof(ep *node_manager.EpochInfo, proofResult []byte) error {
	hash := ep.Hash()
	data := append([]byte{1}, hash[:]...)[:common.HashLength]
	value, err := rlp.EncodeToBytes(data[:])
	if err != nil {
		return err
	}
	if !bytes.Equal(value, proofResult) {
		return fmt.Errorf("expect %s, got %s", hexutil.Encode(value), hexutil.Encode(proofResult))
	}

	return nil
}
