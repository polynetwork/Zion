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

package backend

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	nmabi "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/rlp"
)

var contractAddr = utils.NodeManagerContractAddress

func (s *backend) GetEpochChangeInfo(state *state.StateDB, height *big.Int) (
	beforeChange, changing bool, validators []common.Address, err error) {

	var (
		config *nm.GlobalConfig
		epoch *nm.EpochInfo
	)

	if config, err = s.getGlobalConfig(state, height); err != nil {
		return
	}
	if epoch, err = s.getEpoch(state, height); err != nil {
		return
	}

	epochHeight := new(big.Int).Add(config.BlockPerEpoch, epoch.StartHeight)
	if diff := new(big.Int).Sub(epochHeight, height); diff.Cmp(big.NewInt(1)) == 0 {
		beforeChange = true
		validators = epoch.MemberList()
	} else if diff.Sign() == 0 {
		changing = true
		validators = epoch.MemberList()
	}

	return
}

func (s *backend) getGlobalConfig(state *state.StateDB, height *big.Int) (*nm.GlobalConfig, error) {

	caller := s.signer.Address()
	ref := native.NewContractRef(state, caller, caller, height, common.EmptyHash, 0, nil)
	payload, err := new(nm.GetGlobalConfigParam).Encode()
	if err != nil {
		return nil, fmt.Errorf("encode GetGlobalConfig input failed: %v", err)
	}
	output, _, err := ref.NativeCall(caller, contractAddr, payload)
	if err != nil {
		return nil, fmt.Errorf("GetGlobalConfig native call failed: %v", err)
	}

	var (
		raw    []byte
		config = new(nm.GlobalConfig)
	)
	if err := utils.UnpackOutputs(nm.ABI, nmabi.MethodGetGlobalConfig, &raw, output); err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(raw, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (s *backend) getEpoch(state *state.StateDB, height *big.Int) (*nm.EpochInfo, error) {
	caller := s.signer.Address()
	ref := native.NewContractRef(state, caller, caller, height, common.EmptyHash, 0, nil)
	payload, err := new(nm.GetCurrentEpochInfoParam).Encode()
	if err != nil {
		return nil, fmt.Errorf("encode GetGlobalConfig input failed: %v", err)
	}
	output, _, err := ref.NativeCall(caller, utils.NodeManagerContractAddress, payload)
	if err != nil {
		return nil, fmt.Errorf("GetGlobalConfig native call failed: %v", err)
	}

	var (
		raw   []byte
		epoch = new(nm.EpochInfo)
	)
	if err := utils.UnpackOutputs(nm.ABI, nmabi.MethodGetCurrentEpochInfo, &raw, output); err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(raw, epoch); err != nil {
		return nil, err
	}
	return epoch, nil
}

func (s *backend) execEpochChange(ctx *systemTxContext) error {
	payload, err := new(nm.ChangeEpochParam).Encode()
	if err != nil {
		return err
	}
	return s.executeSystemTx(ctx, contractAddr, payload)
}

func (s *backend) execEndBlock(ctx *systemTxContext) error {
	payload, err := new(nm.ChangeEpochParam).Encode()
	if err != nil {
		return err
	}

	return s.executeSystemTx(ctx, contractAddr, payload)
}
