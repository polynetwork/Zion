package info_sync

import (
	"fmt"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const (
	//key prefix
	ROOT_INFO            = "rootInfo"
	CURRENT_HEIGHT       = "currentHeight"
	SYNC_ROOT_INFO_EVENT = "SyncRootInfoEvent"
	REPLENISH_EVENT      = "ReplenishEvent"
)

func PutRootInfo(native *native.NativeContract, chainID uint64, height uint32, info []byte) error {
	contract := utils.InfoSyncContractAddress
	chainIDBytes := utils.GetUint64Bytes(chainID)
	heightBytes := utils.GetUint32Bytes(height)

	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(ROOT_INFO), chainIDBytes, heightBytes),
		info)
	currentHeight, err := GetCurrentHeight(native, chainID)
	if err != nil {
		return fmt.Errorf("PutRootInfo, GetCurrentHeight error: %v", err)
	}
	if currentHeight < height {
		native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(CURRENT_HEIGHT), chainIDBytes), heightBytes)
	}
	err = NotifyPutRootInfo(native, chainID, height)
	if err != nil {
		return fmt.Errorf("PutRootInfo, NotifyPutRootInfo error: %v", err)
	}
	return nil
}

func GetRootInfo(native *native.NativeContract, chainID uint64, height uint32) ([]byte, error) {
	contract := utils.InfoSyncContractAddress
	chainIDBytes := utils.GetUint64Bytes(chainID)
	heightBytes := utils.GetUint32Bytes(height)

	r, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(ROOT_INFO), chainIDBytes, heightBytes))
	if err != nil {
		return nil, fmt.Errorf("GetRootInfo, native.GetCacheDB().Get error: %v", err)
	}
	return r, nil
}

func GetCurrentHeight(native *native.NativeContract, chainID uint64) (uint32, error) {
	contract := utils.InfoSyncContractAddress
	chainIDBytes := utils.GetUint64Bytes(chainID)

	r, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(CURRENT_HEIGHT), chainIDBytes))
	if err != nil {
		return 0, fmt.Errorf("GetCurrentHeight, native.GetCacheDB().Get error: %v", err)
	}
	return utils.GetBytesUint32(r), nil
}

func NotifyPutRootInfo(native *native.NativeContract, chainID uint64, height uint32) error {
	err := native.AddNotify(ABI, []string{SYNC_ROOT_INFO_EVENT}, chainID, height, native.ContractRef().BlockHeight())
	if err != nil {
		return fmt.Errorf("NotifyPutRootInfo failed: %v", err)
	}
	return nil
}

func NotifyReplenish(native *native.NativeContract, txHashes []string, chainId uint64) error {
	err := native.AddNotify(ABI, []string{REPLENISH_EVENT}, txHashes, chainId)
	if err != nil {
		return fmt.Errorf("NotifyReplenish failed: %v", err)
	}
	return nil
}
