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

package node_manager

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	xctrs "github.com/ethereum/go-ethereum/contracts/native/contract"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/polynetwork/poly/common/config"
	"github.com/polynetwork/poly/core/genesis"
	cstates "github.com/polynetwork/poly/core/states"
)

const (
	//status
	CandidateStatus Status = iota
	ConsensusStatus
	QuitingStatus
	BlackStatus

	//function name
	REGISTER_CANDIDATE   = "registerCandidate"
	UNREGISTER_CANDIDATE = "unRegisterCandidate"
	APPROVE_CANDIDATE    = "approveCandidate"
	BLACK_NODE           = "blackNode"
	WHITE_NODE           = "whiteNode"
	QUIT_NODE            = "quitNode"
	UPDATE_CONFIG        = "updateConfig"
	COMMIT_DPOS          = "commitDpos"

	//key prefix
	GOVERNANCE_VIEW = "governanceView"
	VBFT_CONFIG     = "vbftConfig"
	CANDIDITE_INDEX = "candidateIndex"
	PEER_APPLY      = "peerApply"
	PEER_POOL       = "peerPool"
	PEER_INDEX      = "peerIndex"
	BLACK_LIST      = "blackList"
	CONSENSUS_SIGNS = "consensusSigns"

	//const
	MIN_PEER_NUM = 4
)

func InitNodeManager() {
	ABI = GetABI()
	native.Contracts[this] = RegisterNodeManagerContract
}

//Register methods of node_manager contract
func RegisterNodeManagerContract(native *native.NativeContract) {
	native.Register(genesis.INIT_CONFIG, InitConfig)
	native.Register(REGISTER_CANDIDATE, RegisterCandidate)
	native.Register(UNREGISTER_CANDIDATE, UnRegisterCandidate)
	native.Register(QUIT_NODE, QuitNode)
	native.Register(APPROVE_CANDIDATE, ApproveCandidate)
	native.Register(BLACK_NODE, BlackNode)
	native.Register(WHITE_NODE, WhiteNode)
	native.Register(UPDATE_CONFIG, UpdateConfig)
	native.Register(COMMIT_DPOS, CommitDpos)
}

//Init node_manager contract
func InitConfig(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	configuration := new(config.VBFTConfig)
	if err := utils.UnpackMethod(ABI, MethodInitConfig, configuration, ctx.Payload); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("initConfig, contract params deserialize error: %v", err)
	}

	// check if initConfig is already execute
	peerPoolMapBytes, err := native.GetCacheDB().Get(utils.ConcatKey(this, []byte(PEER_POOL)))
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("initConfig, get peerPoolMap error: %v", err)
	}
	if peerPoolMapBytes != nil {
		return utils.BYTE_FALSE, fmt.Errorf("initConfig. initConfig is already executed")
	}

	//check the configuration
	err = CheckVBFTConfig(configuration)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("initConfig, checkVBFTConfig failed: %v", err)
	}

	var view uint64 = 1
	var maxId uint64

	peerPoolMap := &PeerPoolMap{
		PeerPoolMap: make(map[string]*PeerPoolItem),
	}
	for _, peer := range configuration.Peers {
		if uint64(peer.Index) > maxId {
			maxId = uint64(peer.Index)
		}
		address := common.HexToAddress(peer.Address)

		peerPoolItem := new(PeerPoolItem)
		peerPoolItem.Index = uint64(peer.Index)
		peerPoolItem.PeerPubkey = peer.PeerPubkey
		peerPoolItem.Address = address
		peerPoolItem.Status = ConsensusStatus
		peerPoolMap.PeerPoolMap[peerPoolItem.PeerPubkey] = peerPoolItem

		peerPubkeyPrefix, err := hex.DecodeString(peerPoolItem.PeerPubkey)
		if err != nil {
			return utils.BYTE_FALSE, fmt.Errorf("initConfig, peerPubkey format error: %v", err)
		}
		index := peerPoolItem.Index
		indexBytes := utils.GetUint64Bytes(index)
		native.GetCacheDB().Put(utils.ConcatKey(this, []byte(PEER_INDEX), peerPubkeyPrefix), cstates.GenRawStorageItem(indexBytes))
	}

	//init peer pool
	putPeerPoolMap(native, peerPoolMap, 0)
	putPeerPoolMap(native, peerPoolMap, view)
	indexBytes := utils.GetUint64Bytes(maxId + 1)
	native.GetCacheDB().Put(utils.ConcatKey(this, []byte(CANDIDITE_INDEX)), cstates.GenRawStorageItem(indexBytes))

	//init governance view
	governanceView := &GovernanceView{
		View:   view,
		Height: native.ContractRef().BlockHeight().Uint64(),
		TxHash: native.ContractRef().TxHash(),
	}
	putGovernanceView(native, governanceView)

	//init config
	putConfig(native, &Configuration{
		BlockMsgDelay:        uint64(configuration.BlockMsgDelay),
		HashMsgDelay:         uint64(configuration.HashMsgDelay),
		PeerHandshakeTimeout: uint64(configuration.PeerHandshakeTimeout),
		MaxBlockChangeView:   uint64(configuration.MaxBlockChangeView),
	})

	return utils.BYTE_TRUE, nil
}

//Register a candidate node, used by users.
func RegisterCandidate(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := new(RegisterPeerParam)
	if err := utils.UnpackMethod(ABI, MethodRegisterCandidate, params, ctx.Payload); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("registerCandidate, contract params deserialize error: %v", err)
	}

	//check witness
	err := xctrs.ValidateOwner(native, params.Address)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("registerCandidate, checkWitness error: %v", err)
	}

	//check peerPubkey
	if err := utils.ValidatePeerPubKeyFormat(params.PeerPubkey); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("registerCandidate, invalid peer pubkey")
	}

	peerPubkeyPrefix, err := hex.DecodeString(params.PeerPubkey)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("registerCandidate, peerPubkey format error: %v", err)
	}
	//get black list
	blackList, err := native.GetCacheDB().Get(utils.ConcatKey(this, []byte(BLACK_LIST), peerPubkeyPrefix))
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("registerCandidate, get BlackList error: %v", err)
	}
	if blackList != nil {
		return utils.BYTE_FALSE, fmt.Errorf("registerCandidate, this Peer is in BlackList")
	}

	//check if applied
	peer, err := GetPeerApply(native, params.PeerPubkey)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("registerCandidate, GetPeerApply error: %v", err)
	}
	if peer != nil {
		return utils.BYTE_FALSE, fmt.Errorf("registerCandidate, peer already applied")
	}

	//get current view
	view, err := GetView(native)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("registerCandidate, get view error: %v", err)
	}
	//get peerPoolMap
	peerPoolMap, err := GetPeerPoolMap(native, view)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("registerCandidate, get peerPoolMap error: %v", err)
	}
	//check if exist in PeerPool
	_, ok := peerPoolMap.PeerPoolMap[params.PeerPubkey]
	if ok {
		return utils.BYTE_FALSE, fmt.Errorf("registerCandidate, peerPubkey is already in peerPoolMap")
	}

	if err := putPeerApply(native, params); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("registerCandidate, put putPeerApply error: %v", err)
	}

	if err := native.AddNotify(ABI, []string{EventRegisterCandidate}, params.PeerPubkey); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("registerCandidate, add notify error: %v", err)
	}
	return utils.BYTE_TRUE, nil
}

//Unregister a registered candidate node, will remove node from pool
func UnRegisterCandidate(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := new(PeerParam)
	if err := utils.UnpackMethod(ABI, MethodUnRegisterCandidate, params, ctx.Payload); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("unRegisterCandidate, contract params deserialize error: %v", err)
	}
	contract := utils.NodeManagerContractAddress

	//check witness
	if err := xctrs.ValidateOwner(native, params.Address); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("unRegisterCandidate, checkWitness error: %v", err)
	}

	//check if applied
	peer, err := GetPeerApply(native, params.PeerPubkey)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("unRegisterCandidate, GetPeerApply error: %v", err)
	}
	if peer == nil {
		return utils.BYTE_FALSE, fmt.Errorf("unRegisterCandidate, peer is not applied")
	}
	//check owner address
	if peer.Address != params.Address {
		return utils.BYTE_FALSE, fmt.Errorf("unRegisterCandidate, address is not peer owner")
	}

	peerPubkeyPrefix, err := hex.DecodeString(params.PeerPubkey)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("unRegisterCandidate, peerPubkey format error: %v", err)
	}
	native.GetCacheDB().Delete(utils.ConcatKey(contract, []byte(PEER_APPLY), peerPubkeyPrefix))

	if err := native.AddNotify(ABI, []string{EventUnRegisterCandidate}, params.PeerPubkey); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("unRegisterCandiate, add notify error: %v", err)
	}
	return utils.BYTE_TRUE, nil
}

//Approve a registered candidate node
func ApproveCandidate(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := new(PeerParam)
	if err := utils.UnpackMethod(ABI, MethodApproveCandidate, params, ctx.Payload); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("approveCandidate, contract params deserialize error: %v", err)
	}

	//check witness
	err := xctrs.ValidateOwner(native, params.Address)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("approveCandidate, checkWitness error: %v", err)
	}

	//check if applied
	peer, err := GetPeerApply(native, params.PeerPubkey)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("approveCandidate, GetPeerApply error: %v", err)
	}
	if peer == nil {
		return utils.BYTE_FALSE, fmt.Errorf("approveCandidate, peer is not applied")
	}

	//check consensus signs
	ok, err := CheckConsensusSigns(native, APPROVE_CANDIDATE, []byte(params.PeerPubkey), params.Address)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("approveCandidate, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.BYTE_TRUE, nil
	}

	peerPoolItem := &PeerPoolItem{
		PeerPubkey: peer.PeerPubkey,
		Address:    peer.Address,
	}

	//check if has index
	peerPubkeyPrefix, err := hex.DecodeString(peer.PeerPubkey)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("approveCandidate, peerPubkey format error: %v", err)
	}
	indexBytes, err := native.GetCacheDB().Get(utils.ConcatKey(this, []byte(PEER_INDEX), peerPubkeyPrefix))
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("approveCandidate, get indexBytes error: %v", err)
	}
	if indexBytes != nil {
		value, err := cstates.GetValueFromRawStorageItem(indexBytes)
		if err != nil {
			return nil, fmt.Errorf("approveCandidate, get value from raw storage item error:%v", err)
		}
		index := utils.GetBytesUint64(value)
		peerPoolItem.Index = index
	} else {
		//get candidate index
		candidateIndex, err := getCandidateIndex(native)
		if err != nil {
			return nil, fmt.Errorf("approveCandidate, get candidateIndex error: %v", err)
		}
		peerPoolItem.Index = candidateIndex

		//update candidateIndex
		newCandidateIndex := candidateIndex + 1
		putCandidateIndex(native, newCandidateIndex)

		indexBytes := utils.GetUint64Bytes(peerPoolItem.Index)
		native.GetCacheDB().Put(utils.ConcatKey(this, []byte(PEER_INDEX), peerPubkeyPrefix), cstates.GenRawStorageItem(indexBytes))
	}

	//get current view
	view, err := GetView(native)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("approveCandidate, get view error: %v", err)
	}
	//get peerPoolMap
	peerPoolMap, err := GetPeerPoolMap(native, view)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("approveCandidate, get peerPoolMap error: %v", err)
	}

	peerPoolItem.Status = CandidateStatus
	peerPoolMap.PeerPoolMap[params.PeerPubkey] = peerPoolItem
	putPeerPoolMap(native, peerPoolMap, view)

	native.GetCacheDB().Delete(utils.ConcatKey(this, []byte(PEER_APPLY), peerPubkeyPrefix))

	if err := native.AddNotify(ABI, []string{EventApproveCandidate}, params.PeerPubkey); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("approveCandidate, add notify error: %v", err)
	}
	return utils.BYTE_TRUE, nil
}

//Put a node into black list, remove node from pool
//Node in black list can't be registered.
func BlackNode(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := new(PeerListParam)
	if err := utils.UnpackMethod(ABI, MethodBlackNode, params, ctx.Payload); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("blackNode, contract params deserialize error: %v", err)
	}

	//check witness
	err := xctrs.ValidateOwner(native, params.Address)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("blackNode, checkWitness error: %v", err)
	}

	//get current view
	view, err := GetView(native)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("blackNode, get view error: %v", err)
	}
	//get peerPoolMap
	peerPoolMap, err := GetPeerPoolMap(native, view)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("blackNode, get peerPoolMap error: %v", err)
	}

	//check peers num
	num := 0
	for _, peerPoolItem := range peerPoolMap.PeerPoolMap {
		if peerPoolItem.Status == CandidateStatus || peerPoolItem.Status == ConsensusStatus {
			num = num + 1
		}
	}
	if num <= MIN_PEER_NUM+len(params.PeerPubkeyList)-1 {
		return utils.BYTE_FALSE, fmt.Errorf("blackNode, num of peers is less than 4")
	}
	for _, peerPubkey := range params.PeerPubkeyList {
		peerPoolItem, ok := peerPoolMap.PeerPoolMap[peerPubkey]
		if !ok {
			return utils.BYTE_FALSE, fmt.Errorf("blockNode, peerPubkey: %s is not in peerPoolMap", peerPubkey)
		}
		if peerPoolItem.Status == BlackStatus {
			return utils.BYTE_FALSE, fmt.Errorf("blackNode, peerPubkey: %s is already blacked", peerPubkey)
		}
	}

	input := []byte{}
	for _, v := range params.PeerPubkeyList {
		input = append(input, []byte(v)...)
	}
	//check consensus signs
	ok, err := CheckConsensusSigns(native, BLACK_NODE, input, params.Address)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("blackNode, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.BYTE_TRUE, nil
	}

	commit := false
	for _, peerPubkey := range params.PeerPubkeyList {
		peerPubkeyPrefix, err := hex.DecodeString(peerPubkey)
		if err != nil {
			return utils.BYTE_FALSE, fmt.Errorf("blackNode, peerPubkey format error: %v", err)
		}
		peerPoolItem, ok := peerPoolMap.PeerPoolMap[peerPubkey]
		if !ok {
			return utils.BYTE_FALSE, fmt.Errorf("blackNode, peerPubkey is not in peerPoolMap")
		}

		blackListItem := &BlackListItem{
			PeerPubkey: peerPoolItem.PeerPubkey,
			Address:    peerPoolItem.Address,
		}
		sink := common.NewZeroCopySink(nil)
		blackListItem.Serialization(sink)
		//put peer into black list
		native.GetCacheDB().Put(utils.ConcatKey(this, []byte(BLACK_LIST), peerPubkeyPrefix), cstates.GenRawStorageItem(sink.Bytes()))

		//change peerPool status
		if peerPoolItem.Status == ConsensusStatus {
			commit = true
		}
		peerPoolItem.Status = BlackStatus
		peerPoolMap.PeerPoolMap[peerPubkey] = peerPoolItem
	}
	putPeerPoolMap(native, peerPoolMap, view)

	//commitDpos
	if commit {
		err = executeCommitDpos(native)
		if err != nil {
			return utils.BYTE_FALSE, fmt.Errorf("blackNode, executeCommitDpos error: %v", err)
		}
	}
	if err := native.AddNotify(ABI, []string{EventBlackNode}, params.PeerPubkeyList); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("blackNode, add notify error: %v", err)
	}
	return utils.BYTE_TRUE, nil
}

//Remove a node from black list, allow it to be registered
func WhiteNode(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := new(PeerParam)
	if err := utils.UnpackMethod(ABI, MethodWhiteNode, params, ctx.Payload); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("whiteNode, contract params deserialize error: %v", err)
	}
	contract := utils.NodeManagerContractAddress

	//check witness
	err := xctrs.ValidateOwner(native, params.Address)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("whiteNode, checkWitness error: %v", err)
	}

	peerPubkeyPrefix, err := hex.DecodeString(params.PeerPubkey)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("whiteNode, peerPubkey format error: %v", err)
	}
	//check black list
	blackListBytes, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(BLACK_LIST), peerPubkeyPrefix))
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("whiteNode, get BlackList error: %v", err)
	}
	if blackListBytes == nil {
		return utils.BYTE_FALSE, fmt.Errorf("whiteNode, this Peer is not in BlackList: %v", err)
	}

	//check consensus signs
	ok, err := CheckConsensusSigns(native, WHITE_NODE, []byte(params.PeerPubkey), params.Address)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("whiteNode, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return utils.BYTE_TRUE, nil
	}

	//remove peer from black list
	native.GetCacheDB().Delete(utils.ConcatKey(contract, []byte(BLACK_LIST), peerPubkeyPrefix))
	if err := native.AddNotify(ABI, []string{EventWhiteNode}, params.PeerPubkey); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("whiteNode, add notify error: %v", err)
	}
	return utils.BYTE_TRUE, nil
}

//Quit a registered node, used by node owner.
//Remove node from pool
func QuitNode(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := new(PeerParam)
	if err := utils.UnpackMethod(ABI, MethodQuitNode, params, ctx.Payload); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("quitNode, contract params deserialize error: %v", err)
	}

	//check witness
	err := xctrs.ValidateOwner(native, params.Address)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("quitNode, checkWitness error: %v", err)
	}

	//get current view
	view, err := GetView(native)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("quitNode, get view error: %v", err)
	}
	//get peerPoolMap
	peerPoolMap, err := GetPeerPoolMap(native, view)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("quitNode, get peerPoolMap error: %v", err)
	}

	peerPoolItem, ok := peerPoolMap.PeerPoolMap[params.PeerPubkey]
	if !ok {
		return utils.BYTE_FALSE, fmt.Errorf("quitNode, peerPubkey is not in peerPoolMap")
	}
	if peerPoolItem.Status != ConsensusStatus && peerPoolItem.Status != CandidateStatus {
		return utils.BYTE_FALSE, fmt.Errorf("quitNode, peerPubkey is not CandidateStatus or ConsensusStatus")
	}
	if params.Address != peerPoolItem.Address {
		return utils.BYTE_FALSE, fmt.Errorf("quitNode, peerPubkey is not registered by this address")
	}

	//check peers num
	num := 0
	for _, peerPoolItem := range peerPoolMap.PeerPoolMap {
		if peerPoolItem.Status == CandidateStatus || peerPoolItem.Status == ConsensusStatus {
			num = num + 1
		}
	}
	if num <= MIN_PEER_NUM {
		return utils.BYTE_FALSE, fmt.Errorf("quitNode, num of peers is less than 4")
	}

	//change peerPool status
	peerPoolItem.Status = QuitingStatus

	peerPoolMap.PeerPoolMap[params.PeerPubkey] = peerPoolItem
	putPeerPoolMap(native, peerPoolMap, view)
	if err := native.AddNotify(ABI, []string{EventQuitNode}, params.PeerPubkey); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("quitNode, add notify error: %v", err)
	}
	return utils.BYTE_TRUE, nil
}

//Go to next consensus epoch
func CommitDpos(native *native.NativeContract) ([]byte, error) {
	// get config
	config, err := GetConfig(native)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("commitDpos, get config error: %v", err)
	}

	//get governance view
	governanceView, err := GetGovernanceView(native)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("commitDpos, get GovernanceView error: %v", err)
	}

	// Get current epoch operator
	operatorAddress, err := GetCurConOperator(native)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("commitDPos, get current consensus operator address error: %v", err)
	}

	//check witness
	height := native.ContractRef().BlockHeight().Uint64()
	err = xctrs.ValidateOwner(native, operatorAddress)
	if err != nil {
		cycle := (height - governanceView.Height) >= config.MaxBlockChangeView
		if !cycle {
			return utils.BYTE_FALSE, fmt.Errorf("commitDpos, authentication Failed")
		}
	}

	if err := executeCommitDpos(native); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("executeCommitDpos, executeCommitDpos error: %v", err)
	}

	if err := native.AddNotify(ABI, []string{EventCommitDpos}); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("executeCommitDpos, add notify error: %v", err)
	}
	return utils.BYTE_TRUE, nil
}

//Update VBFT config
func UpdateConfig(native *native.NativeContract) ([]byte, error) {
	ctx := native.ContractRef().CurrentContext()
	params := new(UpdateConfigParam)
	if err := utils.UnpackMethod(ABI, MethodUpdateConfig, params, ctx.Payload); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("updateConfig, deserialize configuration error: %v", err)
	}
	sink := common.NewZeroCopySink(nil)
	params.Configuration.Serialization(sink)

	// Get current epoch operator
	operatorAddress, err := GetCurConOperator(native)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("updateConfig, get current consensus operator address error: %v", err)
	}
	//check witness
	err = xctrs.ValidateOwner(native, operatorAddress)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("updateConfig, checkWitness error: %v", err)
	}

	if params.Configuration.BlockMsgDelay < 5000 {
		return utils.BYTE_FALSE, fmt.Errorf("updateConfig. BlockMsgDelay must >= 5000")
	}
	if params.Configuration.HashMsgDelay < 5000 {
		return utils.BYTE_FALSE, fmt.Errorf("updateConfig. HashMsgDelay must >= 5000")
	}
	if params.Configuration.PeerHandshakeTimeout < 10 {
		return utils.BYTE_FALSE, fmt.Errorf("updateConfig. PeerHandshakeTimeout must >= 10")
	}
	if params.Configuration.MaxBlockChangeView < 10000 {
		return utils.BYTE_FALSE, fmt.Errorf("updateConfig. MaxBlockChangeView must >= 10000")
	}

	putConfig(native, params.Configuration)
	if err := native.AddNotify(ABI, []string{EventUpdateConfig}, params.Configuration); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("updateConfig, add notify error: %v", err)
	}
	return utils.BYTE_TRUE, nil
}
