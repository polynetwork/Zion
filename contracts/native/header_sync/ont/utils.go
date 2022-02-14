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

	"github.com/ethereum/go-ethereum/contracts/native"
	hscommon "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ontio/ontology-crypto/keypair"
	ocommon "github.com/ontio/ontology/common"
	vconfig "github.com/ontio/ontology/consensus/vbft/config"
	"github.com/ontio/ontology/core/signature"
	otypes "github.com/ontio/ontology/core/types"
)

func PutCrossChainMsg(native *native.NativeContract, chainID uint64, crossChainMsg *otypes.CrossChainMsg) error {
	contract := utils.HeaderSyncContractAddress
	sink := ocommon.NewZeroCopySink(nil)
	crossChainMsg.Serialization(sink)
	chainIDBytes := utils.GetUint64Bytes(chainID)
	heightBytes := utils.GetUint32Bytes(crossChainMsg.Height)

	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(hscommon.CROSS_CHAIN_MSG), chainIDBytes, heightBytes),
		sink.Bytes())
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(hscommon.CURRENT_MSG_HEIGHT), chainIDBytes),
		heightBytes)
	hash := crossChainMsg.Hash()
	hscommon.NotifyPutCrossChainMsg(native, chainID, crossChainMsg.Height, hash.ToHexString())
	return nil
}

func GetCrossChainMsg(native *native.NativeContract, chainID uint64, height uint32) (*otypes.CrossChainMsg, error) {
	contract := utils.HeaderSyncContractAddress
	chainIDBytes := utils.GetUint64Bytes(chainID)
	heightBytes := utils.GetUint32Bytes(height)

	crossChainMsgStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(hscommon.CROSS_CHAIN_MSG),
		chainIDBytes, heightBytes))
	if err != nil {
		return nil, fmt.Errorf("GetCrossChainMsg, get headerStore error: %v", err)
	}
	if crossChainMsgStore == nil {
		return nil, fmt.Errorf("GetCrossChainMsg, can not find any header records")
	}
	crossChainMsg := new(otypes.CrossChainMsg)
	if err := crossChainMsg.Deserialization(ocommon.NewZeroCopySource(crossChainMsgStore)); err != nil {
		return nil, fmt.Errorf("GetCrossChainMsg, deserialize header error: %v", err)
	}
	return crossChainMsg, nil
}

func PutBlockHeader(native *native.NativeContract, chainID uint64, blockHeader *otypes.Header) error {
	contract := utils.HeaderSyncContractAddress
	sink := ocommon.NewZeroCopySink(nil)
	blockHeader.Serialization(sink)
	chainIDBytes := utils.GetUint64Bytes(chainID)
	heightBytes := utils.GetUint32Bytes(blockHeader.Height)

	blockHash := blockHeader.Hash()
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(hscommon.BLOCK_HEADER), chainIDBytes, blockHash.ToArray()),
		sink.Bytes())
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(hscommon.HEADER_INDEX), chainIDBytes, heightBytes),
		blockHash.ToArray())
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(hscommon.CURRENT_HEADER_HEIGHT), chainIDBytes),
		heightBytes)
	hscommon.NotifyPutHeader(native, chainID, uint64(blockHeader.Height), blockHash.ToHexString())
	return nil
}

func GetHeaderByHeight(native *native.NativeContract, chainID uint64, height uint32) (*otypes.Header, error) {
	contract := utils.HeaderSyncContractAddress
	chainIDBytes := utils.GetUint64Bytes(chainID)
	heightBytes := utils.GetUint32Bytes(height)

	blockHashStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(hscommon.HEADER_INDEX),
		chainIDBytes, heightBytes))
	if err != nil {
		return nil, fmt.Errorf("GetHeaderByHeight, get blockHashStore error: %v", err)
	}
	if blockHashStore == nil {
		return nil, fmt.Errorf("GetHeaderByHeight, can not find any index records")
	}
	headerStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(hscommon.BLOCK_HEADER),
		chainIDBytes, blockHashStore))
	if err != nil {
		return nil, fmt.Errorf("GetHeaderByHeight, get headerStore error: %v", err)
	}
	if headerStore == nil {
		return nil, fmt.Errorf("GetHeaderByHeight, can not find any header records")
	}
	header := new(otypes.Header)
	if err := header.Deserialization(ocommon.NewZeroCopySource(headerStore)); err != nil {
		return nil, fmt.Errorf("GetHeaderByHeight, deserialize header error: %v", err)
	}
	return header, nil
}

func VerifyCrossChainMsg(native *native.NativeContract, chainID uint64, crossChainMsg *otypes.CrossChainMsg,
	bookkeepers []keypair.PublicKey) error {
	height := crossChainMsg.Height
	//search consensus peer
	keyHeight, err := FindKeyHeight(native, height, chainID)
	if err != nil {
		return fmt.Errorf("verifyCrossChainMsg, findKeyHeight error:%v", err)
	}

	consensusPeer, err := getConsensusPeersByHeight(native, chainID, keyHeight)
	if err != nil {
		return fmt.Errorf("verifyCrossChainMsg, get ConsensusPeer error:%v", err)
	}
	if len(bookkeepers)*3 < len(consensusPeer.PeerMap) {
		return fmt.Errorf("verifyCrossChainMsg, header Bookkeepers num %d must more than 2/3 consensus node num %d",
			len(bookkeepers), len(consensusPeer.PeerMap))
	}
	for _, bookkeeper := range bookkeepers {
		pubkey := vconfig.PubkeyID(bookkeeper)
		_, present := consensusPeer.PeerMap[pubkey]
		if !present {
			return fmt.Errorf("verifyCrossChainMsg, invalid pubkey error:%v", pubkey)
		}
	}
	hash := crossChainMsg.Hash()
	err = signature.VerifyMultiSignature(hash[:], bookkeepers, len(bookkeepers),
		crossChainMsg.SigData)
	if err != nil {
		return fmt.Errorf("verifyCrossChainMsg, VerifyMultiSignature error:%s, heigh:%d", err,
			crossChainMsg.Height)
	}
	return nil
}

//verify header of any height
//find key height and get consensus peer first, then check the sign
func verifyHeader(native *native.NativeContract, chainID uint64, header *otypes.Header) error {
	height := header.Height
	//search consensus peer
	keyHeight, err := FindKeyHeight(native, height, chainID)
	if err != nil {
		return fmt.Errorf("verifyHeader, findKeyHeight error:%v", err)
	}

	consensusPeer, err := getConsensusPeersByHeight(native, chainID, keyHeight)
	if err != nil {
		return fmt.Errorf("verifyHeader, get ConsensusPeer error:%v", err)
	}
	if len(header.Bookkeepers)*3 < len(consensusPeer.PeerMap) {
		return fmt.Errorf("verifyHeader, header Bookkeepers num %d must more than 2/3 consensus node num %d", len(header.Bookkeepers), len(consensusPeer.PeerMap))
	}
	for _, bookkeeper := range header.Bookkeepers {
		pubkey := vconfig.PubkeyID(bookkeeper)
		_, present := consensusPeer.PeerMap[pubkey]
		if !present {
			return fmt.Errorf("verifyHeader, invalid pubkey error:%v", pubkey)
		}
	}
	hash := header.Hash()
	err = signature.VerifyMultiSignature(hash[:], header.Bookkeepers, len(header.Bookkeepers), header.SigData)
	if err != nil {
		return fmt.Errorf("verifyHeader, VerifyMultiSignature error:%s, heigh:%d", err, header.Height)
	}
	return nil
}

func GetKeyHeights(native *native.NativeContract, chainID uint64) (*KeyHeights, error) {
	contract := utils.HeaderSyncContractAddress
	value, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(hscommon.KEY_HEIGHTS), utils.GetUint64Bytes(chainID)))
	if err != nil {
		return nil, fmt.Errorf("GetKeyHeights, get keyHeights value error: %v", err)
	}
	keyHeights := new(KeyHeights)
	if value != nil {
		if err := rlp.DecodeBytes(value, keyHeights); err != nil {
			return nil, err
		}
	}
	return keyHeights, nil
}

func PutKeyHeights(native *native.NativeContract, chainID uint64, keyHeights *KeyHeights) error {
	contract := utils.HeaderSyncContractAddress
	value, err := rlp.EncodeToBytes(keyHeights)
	if err != nil {
		return err
	}
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(hscommon.KEY_HEIGHTS), utils.GetUint64Bytes(chainID)), value)
	return nil
}

func getConsensusPeersByHeight(native *native.NativeContract, chainID uint64, height uint32) (*ConsensusPeers, error) {
	contract := utils.HeaderSyncContractAddress
	heightBytes := utils.GetUint32Bytes(height)
	chainIDBytes := utils.GetUint64Bytes(chainID)
	consensusPeerStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(hscommon.CONSENSUS_PEER), chainIDBytes, heightBytes))
	if err != nil {
		return nil, fmt.Errorf("getConsensusPeerByHeight, get consensusPeerStore error: %v", err)
	}

	if consensusPeerStore == nil {
		return nil, fmt.Errorf("getConsensusPeerByHeight, can not find any record")
	}
	consensusPeers := new(ConsensusPeers)
	if err := rlp.DecodeBytes(consensusPeerStore, consensusPeers); err != nil {
		return nil, err
	}
	return consensusPeers, nil
}

func putConsensusPeers(native *native.NativeContract, consensusPeers *ConsensusPeers) error {
	contract := utils.HeaderSyncContractAddress
	value, err := rlp.EncodeToBytes(consensusPeers)
	if err != nil {
		return err
	}
	chainIDBytes := utils.GetUint64Bytes(consensusPeers.ChainID)
	heightBytes := utils.GetUint32Bytes(consensusPeers.Height)
	blockHeightBytes := native.ContractRef().BlockHeight().Bytes()

	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(hscommon.CONSENSUS_PEER), chainIDBytes, heightBytes), value)
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(hscommon.CONSENSUS_PEER_BLOCK_HEIGHT), chainIDBytes, heightBytes),
		blockHeightBytes)

	//update key heights
	keyHeights, err := GetKeyHeights(native, consensusPeers.ChainID)
	if err != nil {
		return fmt.Errorf("putConsensusPeer, GetKeyHeights error: %v", err)
	}

	keyHeights.HeightList = append(keyHeights.HeightList, consensusPeers.Height)
	err = PutKeyHeights(native, consensusPeers.ChainID, keyHeights)
	if err != nil {
		return fmt.Errorf("putConsensusPeer, putKeyHeights error: %v", err)
	}
	return nil
}

func UpdateConsensusPeer(native *native.NativeContract, chainID uint64, header *otypes.Header) error {
	blkInfo := &vconfig.VbftBlockInfo{}
	if err := json.Unmarshal(header.ConsensusPayload, blkInfo); err != nil {
		return fmt.Errorf("updateConsensusPeer, unmarshal blockInfo error: %s", err)
	}
	if blkInfo.NewChainConfig != nil {
		consensusPeers := &ConsensusPeers{
			ChainID: chainID,
			Height:  header.Height,
			PeerMap: make(map[string]*Peer),
		}
		for _, p := range blkInfo.NewChainConfig.Peers {
			consensusPeers.PeerMap[p.ID] = &Peer{Index: p.Index, PeerPubkey: p.ID}
		}
		err := putConsensusPeers(native, consensusPeers)
		if err != nil {
			return fmt.Errorf("updateConsensusPeer, put ConsensusPeer error: %s", err)
		}
	}
	return nil
}

func FindKeyHeight(native *native.NativeContract, height uint32, chainID uint64) (uint32, error) {
	keyHeights, err := GetKeyHeights(native, chainID)
	if err != nil {
		return 0, fmt.Errorf("findKeyHeight, GetKeyHeights error: %v", err)
	}
	for _, v := range keyHeights.HeightList {
		if height > v {
			return v, nil
		}
	}
	return 0, fmt.Errorf("findKeyHeight, can not find key height with height %d", height)
}
