package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	cstates "github.com/polynetwork/poly/core/states"
)

func PutDoneTx(native *native.NativeContract, crossChainID []byte, chainID uint64) error {
	contract := utils.CrossChainManagerContractAddress
	chainIDBytes := utils.GetUint64Bytes(chainID)
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(DONE_TX), chainIDBytes, crossChainID),
		cstates.GenRawStorageItem(crossChainID))
	return nil
}

func CheckDoneTx(native *native.NativeContract, crossChainID []byte, chainID uint64) error {
	contract := utils.CrossChainManagerContractAddress
	chainIDBytes := utils.GetUint64Bytes(chainID)
	value, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(DONE_TX), chainIDBytes, crossChainID))
	if err != nil {
		return fmt.Errorf("checkDoneTx, native.GetCacheDB().Get error: %v", err)
	}
	if value != nil {
		return fmt.Errorf("checkDoneTx, tx already done")
	}
	return nil
}
