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
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"sort"

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
	FEE                       = "fee"
	FEE_INFO                  = "feeInfo"
	ASSET_BIND                = "assetBind"

	UPDATE_FEE_TIMEOUT = 100
)

var (
	this     = native.NativeContractAddrMap[native.NativeSideChainManager]
	gasTable = map[string]uint64{
		side_chain_manager_abi.MethodGetSideChain:             9751875,
		side_chain_manager_abi.MethodRegisterSideChain:        3635625,
		side_chain_manager_abi.MethodApproveRegisterSideChain: 7263375,
		side_chain_manager_abi.MethodUpdateSideChain:          5599125,
		side_chain_manager_abi.MethodApproveUpdateSideChain:   1270500,
		side_chain_manager_abi.MethodQuitSideChain:            635250,
		side_chain_manager_abi.MethodApproveQuitSideChain:     223125,
		side_chain_manager_abi.MethodRegisterAsset:            10751875,
		side_chain_manager_abi.MethodUpdateFee:                5635625,
		side_chain_manager_abi.MethodGetFee:                   3751875,
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
	s.Register(side_chain_manager_abi.MethodRegisterAsset, RegisterAsset)
	s.Register(side_chain_manager_abi.MethodUpdateFee, UpdateFee)
	s.Register(side_chain_manager_abi.MethodGetFee, GetFee)
}

func GetSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ChainIDParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodGetSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}
	sideChain, err := GetSideChainObject(s, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("GetSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("GetSideChain error: side chain %v not exist", params.ChainID)
	}
	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodGetSideChain, sideChain)
}

func RegisterSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &RegisterSideChainParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodRegisterSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	if len(params.Name) > 100 {
		return nil, errors.New("param name too long, max is 100")
	}

	if len(params.ExtraInfo) > 1000000 {
		return nil, errors.New("param extra info too long, max is 1000000")
	}

	if len(params.CCMCAddress) > 1000 {
		return nil, errors.New("param extra info too long, max is 1000")
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
		Owner:       s.ContractRef().TxOrigin(),
		ChainID:     params.ChainID,
		Router:      params.Router,
		Name:        params.Name,
		CCMCAddress: params.CCMCAddress,
		ExtraInfo:   params.ExtraInfo,
	}
	err = putSideChainApply(s, sideChain)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, putRegisterSideChain error: %v", err)
	}

	err = s.AddNotify(ABI, []string{EventRegisterSideChain}, params.ChainID, params.Router, params.Name)
	if err != nil {
		return nil, fmt.Errorf("RegisterSideChain, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodRegisterSideChain)
}

func ApproveRegisterSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ChainIDParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodApproveRegisterSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	registerSideChain, err := GetSideChainApply(s, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("ApproveRegisterSideChain, getRegisterSideChain error: %v", err)
	}
	if registerSideChain == nil {
		return nil, fmt.Errorf("ApproveRegisterSideChain, chainid is not requested")
	}

	ok, err := node_manager.CheckConsensusSigns(s, side_chain_manager_abi.MethodApproveRegisterSideChain, utils.GetUint64Bytes(params.ChainID),
		s.ContractRef().TxOrigin(), node_manager.Signer)
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

	s.GetCacheDB().Delete(utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(SIDE_CHAIN_APPLY), utils.GetUint64Bytes(params.ChainID)))
	err = s.AddNotify(ABI, []string{EventApproveRegisterSideChain}, params.ChainID)
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

	if len(params.Name) > 100 {
		return nil, errors.New("param name too long, max is 100")
	}

	if len(params.ExtraInfo) > 1000000 {
		return nil, errors.New("param extra info too long, max is 1000000")
	}

	if len(params.CCMCAddress) > 1000 {
		return nil, errors.New("param extra info too long, max is 1000")
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
		Owner:       s.ContractRef().TxOrigin(),
		ChainID:     params.ChainID,
		Router:      params.Router,
		Name:        params.Name,
		CCMCAddress: params.CCMCAddress,
		ExtraInfo:   params.ExtraInfo,
	}
	err = putUpdateSideChain(s, updateSideChain)
	if err != nil {
		return nil, fmt.Errorf("UpdateSideChain, putUpdateSideChain error: %v", err)
	}
	err = s.AddNotify(ABI, []string{EventUpdateSideChain}, params.ChainID, params.Router, params.Name)
	if err != nil {
		return nil, fmt.Errorf("UpdateSideChain, AddNotify error: %v", err)
	}

	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodUpdateSideChain)
}

func ApproveUpdateSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ChainIDParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodApproveUpdateSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	sideChain, err := getUpdateSideChain(s, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("ApproveUpdateSideChain, getUpdateSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("ApproveUpdateSideChain, chainid is not requested update")
	}

	//check consensus signs
	ok, err := node_manager.CheckConsensusSigns(s, side_chain_manager_abi.MethodApproveUpdateSideChain, utils.GetUint64Bytes(params.ChainID),
		s.ContractRef().TxOrigin(), node_manager.Signer)
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

	chainidByte := utils.GetUint64Bytes(params.ChainID)
	s.GetCacheDB().Delete(utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(UPDATE_SIDE_CHAIN_REQUEST), chainidByte))

	err = s.AddNotify(ABI, []string{EventApproveUpdateSideChain}, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("ApproveUpdateSideChain, AddNotify error: %v", err)
	}

	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodApproveUpdateSideChain, true)
}

func QuitSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ChainIDParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodQuitSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	sideChain, err := GetSideChainObject(s, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("QuitSideChain, getSideChain error: %v", err)
	}
	if sideChain == nil {
		return nil, fmt.Errorf("QuitSideChain, side chain is not registered")
	}
	if sideChain.Owner != s.ContractRef().TxOrigin() {
		return nil, fmt.Errorf("QuitSideChain, side chain owner is wrong")
	}
	err = putQuitSideChain(s, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("QuitSideChain, putUpdateSideChain error: %v", err)
	}

	err = s.AddNotify(ABI, []string{EventQuitSideChain}, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("QuitSideChain, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodQuitSideChain)
}

func ApproveQuitSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ChainIDParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodApproveQuitSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	err := getQuitSideChain(s, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("ApproveQuitSideChain, getQuitSideChain error: %v", err)
	}

	//check consensus signs
	ok, err := node_manager.CheckConsensusSigns(s, side_chain_manager_abi.MethodApproveQuitSideChain, utils.GetUint64Bytes(params.ChainID),
		s.ContractRef().TxOrigin(), node_manager.Signer)
	if err != nil {
		return nil, fmt.Errorf("ApproveQuitSideChain, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.PackOutputs(ABI, side_chain_manager_abi.MethodApproveQuitSideChain, true)
	}

	chainidByte := utils.GetUint64Bytes(params.ChainID)
	s.GetCacheDB().Delete(utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(side_chain_manager_abi.MethodApproveQuitSideChain), chainidByte))
	s.GetCacheDB().Delete(utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(SIDE_CHAIN), chainidByte))

	err = s.AddNotify(ABI, []string{EventApproveQuitSideChain}, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("ApproveQuitSideChain, AddNotify error: %v", err)
	}
	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodApproveQuitSideChain, true)
}

func RegisterAsset(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &RegisterAssetParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodRegisterAsset, params, ctx.Payload); err != nil {
		return nil, err
	}

	if len(params.AssetMapKey) != len(params.AssetMapValue) {
		return nil, fmt.Errorf("invalid asset map length")
	}
	if len(params.LockProxyMapKey) != len(params.LockProxyMapValue) {
		return nil, fmt.Errorf("invalid lock proxy map length")
	}

	rippleExtraInfo, err := GetRippleExtraInfo(s, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("RegisterAsset, GetRippleExtraInfo error: %v", err)
	}
	if rippleExtraInfo.Operator != ctx.Caller {
		return nil, fmt.Errorf("RegisterAsset, caller is not operator")
	}

	assetBind, err := GetAssetBind(s, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("RegisterAsset, GetAssetBind error: %v", err)
	}
	for index, v := range params.AssetMapKey {
		assetBind.AssetMap[v] = params.AssetMapValue[index]
	}
	for index, v := range params.LockProxyMapKey {
		assetBind.LockProxyMap[v] = params.LockProxyMapValue[index]
	}

	if err := PutAssetBind(s, params.ChainID, assetBind); err != nil {
		return nil, fmt.Errorf("RegisterAsset, PutAssetBind error: %v", err)
	}
	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodRegisterAsset, true)
}

func UpdateFee(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	blockHeight := s.ContractRef().BlockHeight().Uint64()
	params := &UpdateFeeParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodUpdateFee, params, ctx.Payload); err != nil {
		return nil, err
	}

	//get fee
	fee, err := GetFeeObj(s, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("UpdateFee, GetFeeObj error: %v", err)
	}
	if fee.View != params.ViewNum {
		return nil, fmt.Errorf("UpdateFee, poly view: %d, params view: %d not match",
			fee.View, params.ViewNum)
	}

	//add fee info
	feeInfo, err := GetFeeInfo(s, params.ChainID, fee.View)
	if err != nil {
		return nil, fmt.Errorf("UpdateFee, GetFeeInfo error: %v", err)
	}
	if feeInfo.StartHeight == 0 {
		feeInfo.StartHeight = blockHeight
	} else if blockHeight-feeInfo.StartHeight > UPDATE_FEE_TIMEOUT {
		// if time out view + 1
		fee.View = fee.View + 1
		if err := PutFee(s, params.ChainID, fee); err != nil {
			return nil, fmt.Errorf("UpdateFee, PutFee error: %v", err)
		}
		feeInfo = &FeeInfo{
			StartHeight: blockHeight,
			FeeInfo:     make(map[common.Address]*big.Int),
		}
	}
	feeInfo.FeeInfo[ctx.Caller] = params.Fee
	if err := PutFeeInfo(s, params.ChainID, fee.View, feeInfo); err != nil {
		return nil, fmt.Errorf("UpdateFee, PutFeeInfo error: %v", err)
	}

	//verify signature
	digest, err := params.Digest()
	if err != nil {
		return nil, fmt.Errorf("UpdateFee, digest input param error: %v", err)
	}
	pub, err := crypto.SigToPub(digest, params.Signature)
	if err != nil {
		return nil, fmt.Errorf("UpdateFee, crypto.SigToPub error: %v", err)
	}
	addr := crypto.PubkeyToAddress(*pub)

	//check consensus signs
	id := append(utils.GetUint64Bytes(params.ChainID), utils.GetUint64Bytes(fee.View)...)
	ok, err := node_manager.CheckConsensusSigns(s, side_chain_manager_abi.MethodUpdateFee, id, addr, node_manager.Voter)
	if err != nil {
		return nil, fmt.Errorf("UpdateFee, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return nil, nil
	}
	//vote enough
	feeInfoList := make([]*big.Int, 0, len(feeInfo.FeeInfo))
	for _, v := range feeInfo.FeeInfo {
		feeInfoList = append(feeInfoList, v)
	}
	sort.SliceStable(feeInfoList, func(i, j int) bool {
		return feeInfoList[i].Cmp(feeInfoList[j]) >= 1
	})
	l := len(feeInfoList)
	if l%2 == 0 {
		//even: (a + b)*5 / 2
		fee.Fee = new(big.Int).Div(new(big.Int).Mul(new(big.Int).Add(feeInfoList[l/2], feeInfoList[l/2-1]),
			new(big.Int).SetUint64(5)), new(big.Int).SetUint64(2))
	} else {
		//odd：a * 5
		fee.Fee = new(big.Int).Mul(feeInfoList[(l-1)/2], new(big.Int).SetUint64(5))
	}
	fee.View = fee.View + 1
	if err := PutFee(s, params.ChainID, fee); err != nil {
		return nil, fmt.Errorf("UpdateFee, PutFee error: %v", err)
	}
	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodUpdateFee, true)
}

func GetFee(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ChainIDParam{}
	if err := utils.UnpackMethod(ABI, side_chain_manager_abi.MethodGetFee, params, ctx.Payload); err != nil {
		return nil, err
	}
	//get fee
	fee, err := GetFeeObj(s, params.ChainID)
	if err != nil {
		return nil, fmt.Errorf("GetFee, GetFeeObj error: %v", err)
	}

	b, err := rlp.EncodeToBytes(fee)
	if err != nil {
		return nil, fmt.Errorf("GetFee, rlp encode fee error: %v", err)
	}
	return utils.PackOutputs(ABI, side_chain_manager_abi.MethodGetFee, b)
}
