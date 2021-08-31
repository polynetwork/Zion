package node_manager

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
	polycomm "github.com/polynetwork/poly/common"
	cstates "github.com/polynetwork/poly/core/states"
)

func deleteConsensusSigns(native *native.NativeContract, key polycomm.Uint256) {
	contract := utils.NodeManagerContractAddress
	native.GetCacheDB().Delete(utils.ConcatKey(contract, []byte(CONSENSUS_SIGNS), key.ToArray()))
}

func putConsensusSigns(native *native.NativeContract, key polycomm.Uint256, consensusSigns *ConsensusSigns) {
	contract := utils.NodeManagerContractAddress
	sink := polycomm.NewZeroCopySink(nil)
	consensusSigns.Serialization(sink)
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(CONSENSUS_SIGNS), key.ToArray()), cstates.GenRawStorageItem(sink.Bytes()))
}

func CheckConsensusSigns(native *native.NativeContract, method string, input []byte, address common.Address) (bool, error) {
	message := append([]byte(method), input...)
	key := sha256.Sum256(message)
	consensusSigns, err := getConsensusSigns(native, key)
	if err != nil {
		return false, fmt.Errorf("CheckConsensusSigns, GetConsensusSigns error: %v", err)
	}
	consensusSigns.SignsMap[address] = true

	native.AddNotify(ABI, []string{"CheckConsensusSignsEvent"}, uint64(len(consensusSigns.SignsMap)))

	//check signs num
	//get view
	view, err := GetView(native)
	if err != nil {
		return false, fmt.Errorf("CheckConsensusSigns, GetView error: %v", err)
	}

	//get consensus peer
	peerPoolMap, err := GetPeerPoolMap(native, view)
	if err != nil {
		return false, fmt.Errorf("CheckConsensusSigns, GetPeerPoolMap error: %v", err)
	}
	num := 0
	sum := 0
	for key, v := range peerPoolMap.PeerPoolMap {
		if v.Status == ConsensusStatus {
			k, err := hex.DecodeString(key)
			if err != nil {
				return false, fmt.Errorf("CheckConsensusSigns, hex.DecodeString public key error: %v", err)
			}
			publicKey, err := crypto.UnmarshalPubkey(k)
			if err != nil {
				return false, fmt.Errorf("CheckConsensusSigns, keypair.DeserializePublicKey error: %v", err)
			}
			_, ok := consensusSigns.SignsMap[crypto.PubkeyToAddress(*publicKey)]
			if ok {
				num = num + 1
			}
			sum = sum + 1
		}
	}

	if num >= (2*sum+2)/3 {
		deleteConsensusSigns(native, key)
		return true, nil
	} else {
		putConsensusSigns(native, key, consensusSigns)
		return false, nil
	}
}

func GetPeerPoolMap(native *native.NativeContract, view uint32) (*PeerPoolMap, error) {
	contract := utils.NodeManagerContractAddress
	viewBytes := utils.GetUint32Bytes(view)
	peerPoolMap := &PeerPoolMap{
		PeerPoolMap: make(map[string]*PeerPoolItem),
	}
	peerPoolMapBytes, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(PEER_POOL), viewBytes))
	if err != nil {
		return nil, fmt.Errorf("getPeerPoolMap, get all peerPoolMap error: %v", err)
	}
	if peerPoolMapBytes == nil {
		return nil, fmt.Errorf("getPeerPoolMap, peerPoolMap is nil")
	}
	item := cstates.StorageItem{}
	err = item.Deserialize(bytes.NewBuffer(peerPoolMapBytes))
	if err != nil {
		return nil, fmt.Errorf("deserialize PeerPoolMap error:%v", err)
	}
	peerPoolMapStore := item.Value
	if err := peerPoolMap.Deserialization(polycomm.NewZeroCopySource(peerPoolMapStore)); err != nil {
		return nil, fmt.Errorf("deserialize, deserialize peerPoolMap error: %v", err)
	}
	return peerPoolMap, nil
}

func GetView(native *native.NativeContract) (uint32, error) {
	governanceView, err := GetGovernanceView(native)
	if err != nil {
		return 0, fmt.Errorf("getView, getGovernanceView error: %v", err)
	}
	return governanceView.View, nil
}

func GetGovernanceView(native *native.NativeContract) (*GovernanceView, error) {
	contract := utils.NodeManagerContractAddress
	governanceViewBytes, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(GOVERNANCE_VIEW)))
	if err != nil {
		return nil, fmt.Errorf("getGovernanceView, get governanceViewBytes error: %v", err)
	}
	if governanceViewBytes != nil {
		governanceView := new(GovernanceView)
		value, err := cstates.GetValueFromRawStorageItem(governanceViewBytes)
		if err != nil {
			return nil, fmt.Errorf("getGovernanceView, deserialize from raw storage item err:%v", err)
		}
		if err := governanceView.Deserialization(polycomm.NewZeroCopySource(value)); err != nil {
			return nil, fmt.Errorf("getGovernanceView, deserialize governanceView error: %v", err)
		}
		return governanceView, nil
	}
	return nil, fmt.Errorf("getGovernanceView, get nil governanceViewBytes")
}

func getConsensusSigns(native *native.NativeContract, key polycomm.Uint256) (*ConsensusSigns, error) {
	contract := utils.NodeManagerContractAddress
	consensusSignsStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(CONSENSUS_SIGNS), key.ToArray()))
	if err != nil {
		return nil, fmt.Errorf("GetConsensusSigns, get consensusSignsStore error: %v", err)
	}
	consensusSigns := &ConsensusSigns{
		SignsMap: make(map[common.Address]bool),
	}
	if consensusSignsStore != nil {
		consensusSignsBytes, err := cstates.GetValueFromRawStorageItem(consensusSignsStore)
		if err != nil {
			return nil, fmt.Errorf("getGovernanceView, deserialize from raw storage item err:%v", err)
		}
		if err := consensusSigns.Deserialization(polycomm.NewZeroCopySource(consensusSignsBytes)); err != nil {
			return nil, fmt.Errorf("getGovernanceView, deserialize governanceView error: %v", err)
		}
	}
	return consensusSigns, nil
}
