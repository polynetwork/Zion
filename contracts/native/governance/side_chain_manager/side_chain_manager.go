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
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/side_chain_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	//key prefix
	SIDE_CHAIN_APPLY          = "sideChainApply"
	UPDATE_SIDE_CHAIN_REQUEST = "updateSideChainRequest"
	QUIT_SIDE_CHAIN_REQUEST   = "quitSideChainRequest"
	SIDE_CHAIN                = "sideChain"
	REDEEM_BIND               = "redeemBind"
	BIND_SIGN_INFO            = "bindSignInfo"
	BTC_TX_PARAM              = "btcTxParam"
	REDEEM_SCRIPT             = "redeemScript"
)

var (
	this     = native.NativeContractAddrMap[native.NativeSideChainManager]
	gasTable = map[string]uint64{
		side_chain_manager_abi.MethodGetSideChain: 0,
		side_chain_manager_abi.MethodRegisterSideChain: 100000,
		side_chain_manager_abi.MethodApproveRegisterSideChain: 100000,
		side_chain_manager_abi.MethodUpdateSideChain: 100000,
		side_chain_manager_abi.MethodApproveUpdateSideChain: 100000,
		side_chain_manager_abi.MethodQuitSideChain: 100000,
		side_chain_manager_abi.MethodApproveQuitSideChain: 100000,
		side_chain_manager_abi.MethodRegisterRedeem: 100000,
		side_chain_manager_abi.MethodSetBtcTxParam: 100000,
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
	s.Register(side_chain_manager_abi.MethodRegisterRedeem, RegisterRedeem)
	s.Register(side_chain_manager_abi.MethodSetBtcTxParam, SetBtcTxParam)
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
		Address:      s.ContractRef().TxOrigin(),
		ChainId:      params.ChainID,
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
	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodRegisterSideChain, true)
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
	if sideChain.Address != s.ContractRef().TxOrigin() {
		return nil, fmt.Errorf("UpdateSideChain, side chain owner is wrong")
	}
	updateSideChain := &SideChain{
		Address:      s.ContractRef().TxOrigin(),
		ChainId:      params.ChainID,
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

	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodUpdateSideChain, true)
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
	if sideChain.Address != s.ContractRef().TxOrigin() {
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
	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodQuitSideChain, true)
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

func RegisterRedeem(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &RegisterRedeemParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodRegisterRedeem, params, ctx.Payload); err != nil {
		return nil, err
	}

	ty, addrs, m, err := txscript.ExtractPkScriptAddrs(params.Redeem, netParam)
	if err != nil {
		return nil, fmt.Errorf("RegisterRedeem, failed to extract addrs: %v", err)
	}
	if ty != txscript.MultiSigTy {
		return nil, fmt.Errorf("RegisterRedeem, wrong type of redeem: %s", ty.String())
	}
	rk := btcutil.Hash160(params.Redeem)

	contract, err := GetContractBind(native, params.RedeemChainID, params.ContractChainID, rk)
	if err != nil {
		return nil, fmt.Errorf("RegisterRedeem, failed to get contract and version: %v", err)
	}
	if contract != nil && contract.Ver+1 != params.CVersion {
		return nil, fmt.Errorf("RegisterRedeem, previous version is %d and your version should "+
			"be %d not %d", contract.Ver, contract.Ver+1, params.CVersion)
	}
	verified, err := verifyRedeemRegister(params, addrs)
	if err != nil {
		return nil, fmt.Errorf("RegisterRedeem, failed to verify: %v", err)
	}
	key := append(append(append(rk, utils.GetUint64Bytes(params.RedeemChainID)...),
		params.ContractAddress...), utils.GetUint64Bytes(params.ContractChainID)...)
	bindSignInfo, err := getBindSignInfo(native, key)
	if err != nil {
		return nil, fmt.Errorf("RegisterRedeem, getBindSignInfo error: %v", err)
	}
	for k, v := range verified {
		bindSignInfo.BindSignInfo[k] = v
	}
	err = putBindSignInfo(native, key, bindSignInfo)
	if err != nil {
		return nil, fmt.Errorf("RegisterRedeem, failed to putBindSignInfo: %v", err)
	}

	if len(bindSignInfo.BindSignInfo) >= m {
		err = putContractBind(native, params.RedeemChainID, params.ContractChainID, rk, params.ContractAddress, params.CVersion)
		if err != nil {
			return nil, fmt.Errorf("RegisterRedeem, putContractBind error: %v", err)
		}
		if err = putBtcRedeemScript(native, hex.EncodeToString(rk), params.Redeem, params.RedeemChainID); err != nil {
			return nil, fmt.Errorf("RegisterRedeem, failed to save redeemscript %v with key %v, error: %v", hex.EncodeToString(params.Redeem), rk, err)
		}
		err = native.AddNotify(ABI, []string{EventRegisterRedeem}, hex.EncodeToString(rk), hex.EncodeToString(params.ContractAddress))
		if err != nil {
			return nil, fmt.Errorf("RegisterRedeem, AddNotify error: %v", err)
		}
	}

	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodRegisterRedeem, true)
}

func SetBtcTxParam(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &BtcTxParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodSetBtcTxParam, params, ctx.Payload); err != nil {
		return nil, err
	}

	if params.Detial.FeeRate == 0 {
		return nil, fmt.Errorf("SetBtcTxParam, fee rate can't be zero")
	}
	if params.Detial.MinChange < 2000 {
		return nil, fmt.Errorf("SetBtcTxParam, min-change can't less than 2000")
	}
	cls, addrs, m, err := txscript.ExtractPkScriptAddrs(params.Redeem, netParam)
	if err != nil {
		return nil, fmt.Errorf("SetBtcTxParam, extract addrs from redeem %v", err)
	}
	if cls != txscript.MultiSigTy {
		return nil, fmt.Errorf("SetBtcTxParam, redeem script is not multisig script: %s", cls.String())
	}
	rk := btcutil.Hash160(params.Redeem)
	prev, err := GetBtcTxParam(native, rk, params.RedeemChainId)
	if err != nil {
		return nil, fmt.Errorf("SetBtcTxParam, get previous param error: %v", err)
	}
	if prev != nil && params.Detial.PVersion != prev.PVersion+1 {
		return nil, fmt.Errorf("SetBtcTxParam, previous version is %d and your version should "+
			"be %d not %d", prev.PVersion, prev.PVersion+1, params.Detial.PVersion)
	}

	blob, err := rlp.EncodeToBytes(params.Detial)
	if err != nil {
		return nil, fmt.Errorf("SetBtcTxParam, EncodeToBytes error: %v", err)
	}
	key := append(append(rk, utils.GetUint64Bytes(params.RedeemChainId)...), blob...)
	info, err := getBindSignInfo(native, key)
	if err != nil {
		return nil, fmt.Errorf("SetBtcTxParam, getBindSignInfo error: %v", err)
	}
	if len(info.BindSignInfo) >= m {
		return nil, fmt.Errorf("SetBtcTxParam, the signatures are already enough")
	}
	verified, err := verifyBtcTxParam(params, addrs)
	if err != nil {
		return nil, fmt.Errorf("SetBtcTxParam, failed to verify: %v", err)
	}
	for k, v := range verified {
		info.BindSignInfo[k] = v
	}
	if err = putBindSignInfo(native, key, info); err != nil {
		return nil, fmt.Errorf("SetBtcTxParam, failed to put bindSignInfo: %v", err)
	}
	if len(info.BindSignInfo) >= m {
		if err = putBtcTxParam(native, rk, params.RedeemChainId, params.Detial); err != nil {
			return nil, fmt.Errorf("SetBtcTxParam, failed to put btcTxParam: %v", err)
		}
		native.AddNotify(
			ABI, []string{side_chain_manager_abi.MethodSetBtcTxParam}, hex.EncodeToString(rk), params.RedeemChainId,
			params.Detial.FeeRate, params.Detial.MinChange)
	}

	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodSetBtcTxParam, true)
}
