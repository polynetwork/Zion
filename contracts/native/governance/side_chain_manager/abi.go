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

package side_chain_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/side_chain_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"strings"
)

var (
	EventRegisterSideChain        = side_chain_manager_abi.EventRegisterSideChain
	EventApproveRegisterSideChain = side_chain_manager_abi.EventApproveRegisterSideChain
	EventUpdateSideChain          = side_chain_manager_abi.EventUpdateSideChain
	EventApproveUpdateSideChain   = side_chain_manager_abi.EventApproveUpdateSideChain
	EventQuitSideChain            = side_chain_manager_abi.EventQuitSideChain
	EventApproveQuitSideChain     = side_chain_manager_abi.EventApproveQuitSideChain
)

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(side_chain_manager_abi.ISideChainManagerABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

type RegisterSideChainParam struct {
	ChainID     uint64
	Router      uint64
	Name        string
	CCMCAddress []byte
	ExtraInfo   []byte
}

func (m *RegisterSideChainParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, side_chain_manager_abi.MethodRegisterSideChain, m)
}

type ChainIDParam struct {
	ChainID uint64
}

type UpdateFeeParam struct {
	ChainID   uint64
	ViewNum   uint64
	Fee       *big.Int
	Signature []byte
}

func (m *UpdateFeeParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, side_chain_manager_abi.MethodUpdateFee, m)
}

//Digest Digest calculate the hash of param input
func (m *UpdateFeeParam) Digest() ([]byte, error) {
	input := &UpdateFeeParam{
		ChainID: m.ChainID,
		ViewNum: m.ViewNum,
		Fee:     m.Fee,
	}
	msg, err := rlp.EncodeToBytes(input)
	if err != nil {
		return nil, fmt.Errorf("UpdateFeeParam, serialize input error: %v", err)
	}
	digest := crypto.Keccak256(msg)
	return digest, nil
}

type RegisterAssetParam struct {
	ChainID           uint64
	AssetMapKey       []uint64
	AssetMapValue     [][]byte
	LockProxyMapKey   []uint64
	LockProxyMapValue [][]byte
}

func (m *RegisterAssetParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, side_chain_manager_abi.MethodRegisterAsset, m)
}
