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
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/side_chain_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const (
	//key prefix
	SIDE_CHAIN_APPLY          = "sideChainApply"
	UPDATE_SIDE_CHAIN_REQUEST = "updateSideChainRequest"
	QUIT_SIDE_CHAIN_REQUEST   = "quitSideChainRequest"
	SIDE_CHAIN                = "sideChain"
)

var (
	this     = native.NativeContractAddrMap[native.NativeSideChainManager]
	gasTable = map[string]uint64{
		side_chain_manager_abi.MethodGetSideChain:             0,
		side_chain_manager_abi.MethodRegisterSideChain:        100000,
		side_chain_manager_abi.MethodApproveRegisterSideChain: 100000,
		side_chain_manager_abi.MethodUpdateSideChain:          100000,
		side_chain_manager_abi.MethodApproveUpdateSideChain:   100000,
		side_chain_manager_abi.MethodQuitSideChain:            100000,
		side_chain_manager_abi.MethodApproveQuitSideChain:     100000,
	}

	ABI *abi.ABI
)

func InitSideChainManager() {
	ABI = GetABI()
	native.Contracts[this] = RegisterSideChainManagerContract
}

func RegisterSideChainManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	// s.Register(MethodContractName, Name)
	s.Register(side_chain_manager_abi.MethodGetSideChain, GetSideChain)
	s.Register(side_chain_manager_abi.MethodRegisterSideChain, RegisterSideChain)
	s.Register(side_chain_manager_abi.MethodApproveRegisterSideChain, ApproveRegisterSideChain)
	s.Register(side_chain_manager_abi.MethodUpdateSideChain, UpdateSideChain)
	s.Register(side_chain_manager_abi.MethodApproveUpdateSideChain, ApproveUpdateSideChain)
	s.Register(side_chain_manager_abi.MethodQuitSideChain, QuitSideChain)
	s.Register(side_chain_manager_abi.MethodApproveQuitSideChain, ApproveQuitSideChain)
}

func GetSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ChainidParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodGetSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}
	sideChain, err := GetSideChainObject(s, params.Chainid)
	if err != nil { return nil, fmt.Errorf("GetSideChain error: %v", err) }

	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodGetSideChain, sideChain)
}

func RegisterSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &RegisterSideChainParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodRegisterSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	registerSideChain, err := GetSideChainApply(s, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, getRegisterSideChain error: %v", err)
	}
	if registerSideChain != nil {
		return nil, fmt.Errorf("RegisterSideChain, chainid already requested")
	}
	sideChain, err := GetSideChainObject(s, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, getSideChain error: %v", err)
	}
	if sideChain != nil {
		return nil, fmt.Errorf("RegisterSideChain, chainid already registered")
	}

	sideChain = &SideChain{
		Owner:        s.ContractRef().TxOrigin(),
		ChainID:      params.ChainID,
		Router:       params.Router,
		Name:         params.Name,
		BlocksToWait: params.BlocksToWait,
		CCMCAddress:  params.CCMCAddress,
		ExtraInfo:    params.ExtraInfo,
	}
	err = putSideChainApply(s, sideChain)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, putRegisterSideChain error: %v", err)
	}

	err = s.AddNotify(ABI, []string{EventRegisterSideChain}, params.ChainID, params.Router, params.Name, params.BlocksToWait)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodRegisterSideChain)
}

func ApproveRegisterSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ChainidParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodApproveRegisterSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	registerSideChain, err := GetSideChainApply(s, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterSideChain, getRegisterSideChain error: %v", err)
	}
	if registerSideChain == nil {
		return nil, fmt.Errorf("ApproveRegisterSideChain, chainid is not requested")
	}

	ok, err := node_manager.CheckConsensusSigns(s, side_chain_manager_abi.MethodApproveRegisterSideChain, utils.GetUint64Bytes(params.Chainid),
	s.ContractRef().TxOrigin())
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterSideChain, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.PackOutputs(ABI, side_chain_manager_abi.MethodApproveRegisterSideChain, true)
	}

	err = PutSideChain(s, registerSideChain)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterSideChain, putSideChain error: %v", err)
	}

	s.GetCacheDB().Delete(utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(SIDE_CHAIN_APPLY), utils.GetUint64Bytes(params.Chainid)))
	err = s.AddNotify(ABI, []string{EventApproveRegisterSideChain}, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterSideChain, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodApproveRegisterSideChain, true)
}

func UpdateSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &RegisterSideChainParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodUpdateSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	sideChain, err := GetSideChainObject(s, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("UpdateSideChain, getSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("UpdateSideChain, side chain is not registered")
	}
	if sideChain.Owner != s.ContractRef().TxOrigin() {
		return nil, fmt.Errorf("UpdateSideChain, side chain owner is wrong")
	}
	updateSideChain := &SideChain{
		Owner:        s.ContractRef().TxOrigin(),
		ChainID:      params.ChainID,
		Router:       params.Router,
		Name:         params.Name,
		BlocksToWait: params.BlocksToWait,
		CCMCAddress:  params.CCMCAddress,
		ExtraInfo:    params.ExtraInfo,
	}
	err = putUpdateSideChain(s, updateSideChain)
	if err != nil {
		return nil, fmt.Errorf("UpdateSideChain, putUpdateSideChain error: %v", err)
	}
	err = s.AddNotify(ABI, []string{EventUpdateSideChain}, params.ChainID, params.Router, params.Name, params.BlocksToWait)
	if err != nil {
		return nil, fmt.Errorf("UpdateSideChain, AddNotify error: %v", err)
	}

	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodUpdateSideChain)
}

func ApproveUpdateSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ChainidParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodApproveUpdateSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	sideChain, err := getUpdateSideChain(s, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("ApproveUpdateSideChain, getUpdateSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("ApproveUpdateSideChain, chainid is not requested update")
	}

	//check consensus signs
	ok, err := node_manager.CheckConsensusSigns(s, side_chain_manager_abi.MethodApproveUpdateSideChain, utils.GetUint64Bytes(params.Chainid),
	s.ContractRef().TxOrigin())
	if err != nil {
		return nil, fmt.Errorf("ApproveUpdateSideChain, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.PackOutputs(ABI, side_chain_manager_abi.MethodApproveUpdateSideChain, true)
	}

	err = PutSideChain(s, sideChain)
	if err != nil {
		return nil, fmt.Errorf("ApproveUpdateSideChain, putSideChain error: %v", err)
	}

	chainidByte := utils.GetUint64Bytes(params.Chainid)
	s.GetCacheDB().Delete(utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(UPDATE_SIDE_CHAIN_REQUEST), chainidByte))

	err = s.AddNotify(ABI, []string{EventApproveUpdateSideChain}, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("ApproveUpdateSideChain, AddNotify error: %v", err)
	}

	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodApproveUpdateSideChain, true)
}

func QuitSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ChainidParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodQuitSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	sideChain, err := GetSideChainObject(s, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("QuitSideChain, getSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("QuitSideChain, side chain is not registered")
	}
	if sideChain.Owner != s.ContractRef().TxOrigin() {
		return nil, fmt.Errorf("QuitSideChain, side chain owner is wrong")
	}
	err = putQuitSideChain(s, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("QuitSideChain, putUpdateSideChain error: %v", err)
	}

	err = s.AddNotify(ABI, []string{EventQuitSideChain}, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("QuitSideChain, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodQuitSideChain)
}

func ApproveQuitSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ChainidParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodApproveQuitSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	err := getQuitSideChain(s, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("ApproveQuitSideChain, getQuitSideChain error: %v", err)
	}

	//check consensus signs
	ok, err := node_manager.CheckConsensusSigns(s, side_chain_manager_abi.MethodApproveQuitSideChain, utils.GetUint64Bytes(params.Chainid),
		s.ContractRef().TxOrigin())
	if err != nil {
		return nil, fmt.Errorf("ApproveQuitSideChain, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.PackOutputs(ABI, side_chain_manager_abi.MethodApproveQuitSideChain, true)
	}

	chainidByte := utils.GetUint64Bytes(params.Chainid)
	s.GetCacheDB().Delete(utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(side_chain_manager_abi.MethodApproveQuitSideChain), chainidByte))
	s.GetCacheDB().Delete(utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(SIDE_CHAIN), chainidByte))

	err = s.AddNotify(ABI, []string{EventApproveQuitSideChain}, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("ApproveQuitSideChain, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodApproveQuitSideChain, true)
}
