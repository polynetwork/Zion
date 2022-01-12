package cross_chain_manager

import (
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

func PutBlackChain(native *native.NativeContract, chainID uint64) error {
	native.GetCacheDB().Put(blackChainKey(chainID), utils.GetUint64Bytes(chainID))
	return nil
}

func RemoveBlackChain(native *native.NativeContract, chainID uint64) {
	native.GetCacheDB().Delete(blackChainKey(chainID))
}

func CheckIfChainBlacked(native *native.NativeContract, chainID uint64) (bool, error) {
	chainIDStore, err := native.GetCacheDB().Get(blackChainKey(chainID))
	if err != nil {
		return true, fmt.Errorf("CheckBlackChain, get black chainIDStore error: %v", err)
	}
	if chainIDStore == nil {
		return false, nil
	}
	return true, nil
}

func blackChainKey(chainID uint64) []byte {
	contract := utils.CrossChainManagerContractAddress
	return utils.ConcatKey(contract, []byte(BLACKED_CHAIN), utils.GetUint64Bytes(chainID))
}
