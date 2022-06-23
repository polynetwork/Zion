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
	"io"
	"strings"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/side_chain_manager_abi"
)

var (
	EventRegisterSideChain        = side_chain_manager_abi.MethodRegisterSideChain
	EventApproveRegisterSideChain = side_chain_manager_abi.MethodApproveRegisterSideChain
	EventUpdateSideChain          = side_chain_manager_abi.MethodUpdateSideChain
	EventApproveUpdateSideChain   = side_chain_manager_abi.MethodApproveUpdateSideChain
	EventQuitSideChain            = side_chain_manager_abi.MethodQuitSideChain
	EventApproveQuitSideChain     = side_chain_manager_abi.MethodApproveQuitSideChain
	EventRegisterRedeem           = side_chain_manager_abi.MethodRegisterRedeem
)

func GetABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(side_chain_manager_abi.SideChainManagerABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load abi json string: [%v]", err))
	}
	return &ab
}

type RegisterSideChainParam struct {
	Address      common.Address
	ChainId      uint64
	Router       uint64
	Name         string
	BlocksToWait uint64
	CCMCAddress  []byte
	ExtraInfo    []byte
}

type ChainidParam struct {
	Chainid uint64
	Address common.Address
}

type RegisterRedeemParam struct {
	RedeemChainID   uint64
	ContractChainID uint64
	Redeem          []byte
	CVersion        uint64
	ContractAddress []byte
	Signs           [][]byte
}

type BtcTxParam struct {
	Redeem        []byte
	RedeemChainId uint64
	Sigs          [][]byte
	Detial        *BtcTxParamDetial
}

type BtcTxParamDetial struct {
	PVersion  uint64
	FeeRate   uint64
	MinChange uint64
}

func (m *BtcTxParamDetial) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.PVersion, m.FeeRate, m.MinChange})
}
func (m *BtcTxParamDetial) DecodeRLP(s *rlp.Stream) error {
	var data struct {
		PVersion  uint64
		FeeRate   uint64
		MinChange uint64
	}

	if err := s.Decode(&data); err != nil {
		return err
	}
	m.PVersion, m.FeeRate, m.MinChange = data.PVersion, data.FeeRate, data.MinChange
	return nil
}
