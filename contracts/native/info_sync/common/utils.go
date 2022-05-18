package common

import (
	"fmt"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

func PutCrossChainInfo(native *native.NativeContract, chainID uint64, key, value []byte) error {
	contract := utils.InfoSyncContractAddress
	chainIDBytes := utils.GetUint64Bytes(chainID)

	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(CROSS_CHAIN_INFO), chainIDBytes, key),
		value)
	NotifyPutCrossChainInfo(native, chainID, key, value)
	return nil
}

func GetCrossChainInfo(native *native.NativeContract, chainID uint64, key []byte) ([]byte, error) {
	contract := utils.InfoSyncContractAddress
	chainIDBytes := utils.GetUint64Bytes(chainID)

	r, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(CROSS_CHAIN_INFO), chainIDBytes, key))
	if err != nil {
		return nil, fmt.Errorf("GetCrossChainInfo, native.GetCacheDB().Get error: %v", err)
	}
	return r, nil
}