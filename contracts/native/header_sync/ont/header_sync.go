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

package ont

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"

	"github.com/ethereum/go-ethereum/contracts/native"
	hscommon "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ontio/ontology-crypto/keypair"
	ocommon "github.com/ontio/ontology/common"
	otypes "github.com/ontio/ontology/core/types"
)

type ONTHandler struct {
}

func NewONTHandler() *ONTHandler {
	return &ONTHandler{}
}

func (this *ONTHandler) SyncGenesisHeader(native *native.NativeContract) error {
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

	var header *otypes.Header
	err = json.Unmarshal(params.GenesisHeader, &header)
	if err != nil {
		return fmt.Errorf("SyncGenesisHeader, json.Unmarshal header err: %v", err)
	}
	//block header storage
	err = PutBlockHeader(native, params.ChainID, header)
	if err != nil {
		return fmt.Errorf("ONTHandler SyncGenesisHeader, put blockHeader error: %v", err)
	}

	//consensus node pk storage
	err = UpdateConsensusPeer(native, params.ChainID, header)
	if err != nil {
		return fmt.Errorf("SyncGenesisHeader, update ConsensusPeer error: %v", err)
	}
	return nil
}

func (this *ONTHandler) SyncBlockHeader(native *native.NativeContract) error {
	params := &hscommon.SyncBlockHeaderParam{}
	{
		ctx := native.ContractRef().CurrentContext()
		if err := utils.UnpackMethod(hscommon.ABI, hscommon.MethodSyncBlockHeader, params, ctx.Payload); err != nil {
			return err
		}
	}
	for _, v := range params.Headers {
		header, err := otypes.HeaderFromRawBytes(v)
		if err != nil {
			return fmt.Errorf("SyncBlockHeader, otypes.HeaderFromRawBytes error: %v", err)
		}
		_, err = GetHeaderByHeight(native, params.ChainID, header.Height)
		if err == nil {
			continue
		}
		err = verifyHeader(native, params.ChainID, header)
		if err != nil {
			return fmt.Errorf("SyncBlockHeader, verifyHeader error: %v", err)
		}
		err = PutBlockHeader(native, params.ChainID, header)
		if err != nil {
			return fmt.Errorf("SyncBlockHeader, put BlockHeader error: %v", err)
		}
		err = UpdateConsensusPeer(native, params.ChainID, header)
		if err != nil {
			return fmt.Errorf("SyncBlockHeader, update ConsensusPeer error: %v", err)
		}
	}
	return nil
}

func (this *ONTHandler) SyncCrossChainMsg(native *native.NativeContract) error {
	params := &hscommon.SyncCrossChainMsgParam{}
	{
		ctx := native.ContractRef().CurrentContext()
		if err := utils.UnpackMethod(hscommon.ABI, hscommon.MethodSyncBlockHeader, params, ctx.Payload); err != nil {
			return err
		}
	}
	for _, v := range params.CrossChainMsgs {
		source := ocommon.NewZeroCopySource(v)
		crossChainMsg := new(otypes.CrossChainMsg)
		err := crossChainMsg.Deserialization(source)
		if err != nil {
			return fmt.Errorf("SyncCrossChainMsg, deserialize crossChainMsg error: %v", err)
		}
		n, _, irr, eof := source.NextVarUint()
		if irr || eof {
			return fmt.Errorf("SyncCrossChainMsg, deserialization bookkeeper length error")
		}
		var bookkeepers []keypair.PublicKey
		for i := 0; uint64(i) < n; i++ {
			v, _, irr, eof := source.NextVarBytes()
			if irr || eof {
				return fmt.Errorf("SyncCrossChainMsg, deserialization bookkeeper error")
			}
			bookkeeper, err := keypair.DeserializePublicKey(v)
			if err != nil {
				return fmt.Errorf("SyncCrossChainMsg, keypair.DeserializePublicKey error: %v", err)
			}
			bookkeepers = append(bookkeepers, bookkeeper)
		}
		_, err = GetCrossChainMsg(native, params.ChainID, crossChainMsg.Height)
		if err == nil {
			continue
		}
		err = VerifyCrossChainMsg(native, params.ChainID, crossChainMsg, bookkeepers)
		if err != nil {
			return fmt.Errorf("SyncCrossChainMsg, VerifyCrossChainMsg error: %v", err)
		}
		err = PutCrossChainMsg(native, params.ChainID, crossChainMsg)
		if err != nil {
			return fmt.Errorf("SyncCrossChainMsg, put PutCrossChainMsg error: %v", err)
		}
	}
	return nil
}
