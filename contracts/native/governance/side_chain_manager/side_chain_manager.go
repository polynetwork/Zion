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
	"github.com/ethereum/go-ethereum/contracts/native/contract"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/polynetwork/poly/common"
)

const contractName = "side chain manager"

const (
	//function name
	MethodContractName             = "name"
	MethodRegisterSideChain        = "registerSideChain"
	MethodApproveRegisterSideChain = "approveRegisterSideChain"
	MethodUpdateSideChain          = "updateSideChain"
	MethodApproveUpdateSideChain   = "approveUpdateSideChain"
	MethodQuitSideChain            = "quitSideChain"
	MethodApproveQuitSideChain     = "approveQuitSideChain"
	MethodRegisterRedeem           = "registerRedeem"
	MethodSetBtcTxParam            = "setBtcTxParam"

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
		// MethodContractName:             0,
		MethodRegisterSideChain:        0,
		MethodApproveRegisterSideChain: 100000,
		MethodUpdateSideChain:          0,
		MethodApproveUpdateSideChain:   0,
		MethodQuitSideChain:            0,
		MethodApproveQuitSideChain:     0,
		MethodRegisterRedeem:           0,
		MethodSetBtcTxParam:            0,
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
	s.Register(MethodRegisterSideChain, RegisterSideChain)
	s.Register(MethodApproveRegisterSideChain, ApproveRegisterSideChain)
	s.Register(MethodUpdateSideChain, UpdateSideChain)
	s.Register(MethodApproveUpdateSideChain, ApproveUpdateSideChain)
	s.Register(MethodQuitSideChain, QuitSideChain)
	s.Register(MethodApproveQuitSideChain, ApproveQuitSideChain)
	s.Register(MethodRegisterRedeem, RegisterRedeem)
	s.Register(MethodSetBtcTxParam, SetBtcTxParam)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(ABI, MethodContractName, contractName)
}

func RegisterSideChain(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &RegisterSideChainParam{}
	if err := utils.UnpackMethod(ABI, MethodRegisterSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	//check witness
	err := contract.ValidateOwner(native, params.Address)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, checkWitness error: %v", err)
	}
	registerSideChain, err := getSideChainApply(native, params.ChainId)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, getRegisterSideChain error: %v", err)
	}
	if registerSideChain != nil {
		return nil, fmt.Errorf("RegisterSideChain, chainid already requested")
	}
	sideChain, err := GetSideChain(native, params.ChainId)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, getSideChain error: %v", err)
	}
	if sideChain != nil {
		return nil, fmt.Errorf("RegisterSideChain, chainid already registered")
	}

	sideChain = &SideChain{
		Address:      params.Address,
		ChainId:      params.ChainId,
		Router:       params.Router,
		Name:         params.Name,
		BlocksToWait: params.BlocksToWait,
		CCMCAddress:  params.CCMCAddress,
		ExtraInfo:    params.ExtraInfo,
	}
	err = putSideChainApply(native, sideChain)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, putRegisterSideChain error: %v", err)
	}

	err = native.AddNotify(ABI, []string{EventRegisterSideChain}, params.ChainId, params.Router, params.Name, params.BlocksToWait)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodRegisterSideChain, true)
}

func ApproveRegisterSideChain(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &ChainidParam{}
	if err := utils.UnpackMethod(ABI, MethodApproveRegisterSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	//check witness
	err := contract.ValidateOwner(native, params.Address)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterSideChain, checkWitness error: %v", err)
	}

	registerSideChain, err := getSideChainApply(native, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterSideChain, getRegisterSideChain error: %v", err)
	}
	if registerSideChain == nil {
		return nil, fmt.Errorf("ApproveRegisterSideChain, chainid is not requested")
	}

	ok, err := node_manager.CheckConsensusSigns(native, MethodApproveRegisterSideChain, utils.GetUint64Bytes(params.Chainid),
		params.Address)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterSideChain, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.PackOutputs(ABI, MethodApproveRegisterSideChain, true)
	}

	err = PutSideChain(native, registerSideChain)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterSideChain, putSideChain error: %v", err)
	}

	native.GetCacheDB().Delete(utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(SIDE_CHAIN_APPLY), utils.GetUint64Bytes(params.Chainid)))
	err = native.AddNotify(ABI, []string{EventApproveRegisterSideChain}, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterSideChain, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodApproveRegisterSideChain, true)
}

func UpdateSideChain(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &RegisterSideChainParam{}
	if err := utils.UnpackMethod(ABI, MethodUpdateSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	//check witness
	err := contract.ValidateOwner(native, params.Address)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, checkWitness error: %v", err)
	}

	sideChain, err := GetSideChain(native, params.ChainId)
	if err != nil {
		return nil, fmt.Errorf("UpdateSideChain, getSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("UpdateSideChain, side chain is not registered")
	}
	if sideChain.Address != params.Address {
		return nil, fmt.Errorf("UpdateSideChain, side chain owner is wrong")
	}
	updateSideChain := &SideChain{
		Address:      params.Address,
		ChainId:      params.ChainId,
		Router:       params.Router,
		Name:         params.Name,
		BlocksToWait: params.BlocksToWait,
		CCMCAddress:  params.CCMCAddress,
		ExtraInfo:    params.ExtraInfo,
	}
	err = putUpdateSideChain(native, updateSideChain)
	if err != nil {
		return nil, fmt.Errorf("UpdateSideChain, putUpdateSideChain error: %v", err)
	}
	err = native.AddNotify(ABI, []string{EventUpdateSideChain}, params.ChainId, params.Router, params.Name, params.BlocksToWait)
	if err != nil {
		return nil, fmt.Errorf("UpdateSideChain, AddNotify error: %v", err)
	}

	return utils.PackOutputs(ABI, MethodUpdateSideChain, true)
}

func ApproveUpdateSideChain(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &ChainidParam{}
	if err := utils.UnpackMethod(ABI, MethodApproveUpdateSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	//check witness
	err := contract.ValidateOwner(native, params.Address)
	if err != nil {
		return nil, fmt.Errorf("ApproveUpdateSideChain, checkWitness error: %v", err)
	}
	sideChain, err := getUpdateSideChain(native, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("ApproveUpdateSideChain, getUpdateSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("ApproveUpdateSideChain, chainid is not requested update")
	}

	//check consensus signs
	ok, err := node_manager.CheckConsensusSigns(native, MethodApproveUpdateSideChain, utils.GetUint64Bytes(params.Chainid),
		params.Address)
	if err != nil {
		return nil, fmt.Errorf("ApproveUpdateSideChain, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		utils.PackOutputs(ABI, MethodApproveUpdateSideChain, true)
	}

	err = PutSideChain(native, sideChain)
	if err != nil {
		return nil, fmt.Errorf("ApproveUpdateSideChain, putSideChain error: %v", err)
	}

	chainidByte := utils.GetUint64Bytes(params.Chainid)
	native.GetCacheDB().Delete(utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(UPDATE_SIDE_CHAIN_REQUEST), chainidByte))

	err = native.AddNotify(ABI, []string{EventApproveUpdateSideChain}, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("ApproveUpdateSideChain, AddNotify error: %v", err)
	}

	return utils.PackOutputs(ABI, MethodApproveUpdateSideChain, true)
}

func QuitSideChain(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &ChainidParam{}
	if err := utils.UnpackMethod(ABI, MethodQuitSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	//check witness
	err := contract.ValidateOwner(native, params.Address)
	if err != nil {
		return nil, fmt.Errorf("ApproveUpdateSideChain, checkWitness error: %v", err)
	}

	sideChain, err := GetSideChain(native, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("QuitSideChain, getSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("QuitSideChain, side chain is not registered")
	}
	if sideChain.Address != params.Address {
		return nil, fmt.Errorf("QuitSideChain, side chain owner is wrong")
	}
	err = putQuitSideChain(native, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("QuitSideChain, putUpdateSideChain error: %v", err)
	}

	err = native.AddNotify(ABI, []string{EventQuitSideChain}, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("QuitSideChain, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodQuitSideChain, true)
}

func ApproveQuitSideChain(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &ChainidParam{}
	if err := utils.UnpackMethod(ABI, MethodApproveQuitSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	//check witness
	err := contract.ValidateOwner(native, params.Address)
	if err != nil {
		return nil, fmt.Errorf("ApproveQuitSideChain, checkWitness error: %v", err)
	}

	err = getQuitSideChain(native, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("ApproveQuitSideChain, getQuitSideChain error: %v", err)
	}

	//check consensus signs
	ok, err := node_manager.CheckConsensusSigns(native, MethodQuitSideChain, utils.GetUint64Bytes(params.Chainid),
		params.Address)
	if err != nil {
		return nil, fmt.Errorf("ApproveQuitSideChain, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.PackOutputs(ABI, MethodApproveQuitSideChain, true)
	}

	chainidByte := utils.GetUint64Bytes(params.Chainid)
	native.GetCacheDB().Delete(utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(MethodQuitSideChain), chainidByte))
	native.GetCacheDB().Delete(utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(SIDE_CHAIN), chainidByte))

	err = native.AddNotify(ABI, []string{EventApproveQuitSideChain}, params.Chainid)
	if err != nil {
		return nil, fmt.Errorf("ApproveQuitSideChain, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, MethodApproveQuitSideChain, true)
}

func RegisterRedeem(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &RegisterRedeemParam{}
	if err := utils.UnpackMethod(ABI, MethodRegisterRedeem, params, ctx.Payload); err != nil {
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

	return utils.PackOutputs(ABI, MethodRegisterRedeem, true)
}

func SetBtcTxParam(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := &BtcTxParam{}
	if err := utils.UnpackMethod(ABI, MethodRegisterRedeem, params, ctx.Payload); err != nil {
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
	sink := common.NewZeroCopySink(nil)
	params.Detial.Serialization(sink)
	key := append(append(rk, utils.GetUint64Bytes(params.RedeemChainId)...), sink.Bytes()...)
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
			ABI, []string{MethodRegisterRedeem}, hex.EncodeToString(rk), params.RedeemChainId,
			params.Detial.FeeRate, params.Detial.MinChange)
	}

	return utils.PackOutputs(ABI, MethodRegisterRedeem, true)
}
