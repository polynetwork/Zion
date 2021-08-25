package common

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	cstates "github.com/polynetwork/poly/core/states"
)

func Replace0x(s string) string {
	return strings.Replace(strings.ToLower(s), "0x", "", 1)
}

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

func PutBytes(native *native.NativeContract, key []byte, value []byte) {
	native.GetCacheDB().Put(key, cstates.GenRawStorageItem(value))
}

func NotifyMakeProof(native *native.NativeContract, fromChainID, toChainID uint64, txHash string, key string) {

	native.AddNotify(ABI, []string{NOTIFY_MAKE_PROOF}, fromChainID, toChainID, txHash, native.ContractRef().BlockHeight(), key)
}
