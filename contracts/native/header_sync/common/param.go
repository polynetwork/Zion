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
package common

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
)

const (
	//key prefix
	CROSS_CHAIN_MSG             = "crossChainMsg"
	CURRENT_MSG_HEIGHT          = "currentMsgHeight"
	BLOCK_HEADER                = "blockHeader"
	CURRENT_HEADER_HEIGHT       = "currentHeaderHeight"
	HEADER_INDEX                = "headerIndex"
	CONSENSUS_PEER              = "consensusPeer"
	CONSENSUS_PEER_BLOCK_HEIGHT = "consensusPeerBlockHeight"
	KEY_HEIGHTS                 = "keyHeights"
	ETH_CACHE                   = "ethCaches"
	GENESIS_HEADER              = "genesisHeader"
	MAIN_CHAIN                  = "mainChain"
	EPOCH_SWITCH                = "epochSwitch"
	SYNC_HEADER_NAME_EVENT      = "syncHeader"
	SYNC_CROSSCHAIN_MSG         = "syncCrossChainMsg"
	POLYGON_SPAN                = "polygonSpan"
)

type HeaderSyncHandler interface {
	SyncGenesisHeader(service *native.NativeContract) error
	SyncBlockHeader(service *native.NativeContract) error
	SyncCrossChainMsg(service *native.NativeContract) error
}

type SyncGenesisHeaderParam struct {
	ChainID       uint64
	GenesisHeader []byte
}

func (m *SyncGenesisHeaderParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ChainID, m.GenesisHeader})
}
func (m *SyncGenesisHeaderParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ChainID       uint64
		GenesisHeader []byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.ChainID, m.GenesisHeader = data.ChainID, data.GenesisHeader
	return nil
}

type SyncBlockHeaderParam struct {
	ChainID uint64
	Address common.Address
	Headers [][]byte
}

func (m *SyncBlockHeaderParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ChainID, m.Address, m.Headers})
}
func (m *SyncBlockHeaderParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ChainID uint64
		Address common.Address
		Headers [][]byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.ChainID, m.Address, m.Headers = data.ChainID, data.Address, data.Headers
	return nil
}

type SyncCrossChainMsgParam struct {
	ChainID        uint64
	Address        common.Address
	CrossChainMsgs [][]byte
}

func (m *SyncCrossChainMsgParam) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.ChainID, m.Address, m.CrossChainMsgs})
}
func (m *SyncCrossChainMsgParam) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		ChainID        uint64
		Address        common.Address
		CrossChainMsgs [][]byte
	}

	if err := s.Decode(&data); err != nil {
		return err
	}

	m.ChainID, m.Address, m.CrossChainMsgs = data.ChainID, data.Address, data.CrossChainMsgs
	return nil
}

func NotifyPutHeader(native *native.NativeContract, chainID uint64, height uint64, blockHash string) {

	err := native.AddNotify(ABI, []string{SYNC_HEADER_NAME_EVENT}, chainID, height, blockHash, native.ContractRef().BlockHeight())
	if err != nil {
		panic(fmt.Sprintf("NotifyPutHeader failed: %v", err))
	}
}

func NotifyPutCrossChainMsg(native *native.NativeContract, chainID uint64, height uint32, hash string) {

	err := native.AddNotify(ABI, []string{SYNC_CROSSCHAIN_MSG}, chainID, height, hash, native.ContractRef().BlockHeight())
	if err != nil {
		panic(fmt.Sprintf("NotifyPutHeader failed: %v", err))
	}
}
