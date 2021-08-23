package cross_chain_manager

import (
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	cstates "github.com/polynetwork/poly/core/states"
)

func PutBlackChain(native *native.NativeContract, chainID uint64) {
	contract := utils.CrossChainManagerContractAddress
	chainIDBytes := utils.GetUint64Bytes(chainID)
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(BLACKED_CHAIN), chainIDBytes),
		cstates.GenRawStorageItem(chainIDBytes))
}

func RemoveBlackChain(native *native.NativeContract, chainID uint64) {
	contract := utils.CrossChainManagerContractAddress
	chainIDBytes := utils.GetUint64Bytes(chainID)
	native.GetCacheDB().Delete(utils.ConcatKey(contract, []byte(BLACKED_CHAIN), chainIDBytes))
}

func CheckIfChainBlacked(native *native.NativeContract, chainID uint64) (bool, error) {
	contract := utils.CrossChainManagerContractAddress
	chainIDBytes := utils.GetUint64Bytes(chainID)
	chainIDStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(BLACKED_CHAIN), chainIDBytes))
	if err != nil {
		return true, fmt.Errorf("CheckBlackChain, get black chainIDStore error: %v", err)
	}
	if chainIDStore == nil {
		return false, nil
	}
	return true, nil
}
