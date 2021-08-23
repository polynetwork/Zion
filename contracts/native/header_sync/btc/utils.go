package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	polycomm "github.com/polynetwork/poly/common"
	cstates "github.com/polynetwork/poly/core/states"
)

func GetBlockHashByHeight(native *native.NativeContract, chainID uint64, height uint32) (*chainhash.Hash, error) {
	contract := utils.HeaderSyncContractAddress

	hashStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(scom.HEADER_INDEX), utils.GetUint64Bytes(chainID), utils.GetUint32Bytes(height)))
	if err != nil {
		return nil, fmt.Errorf("GetBlockHashByHeight, get heightBlockHashStore error: %v", err)
	}
	if hashStore == nil {
		return nil, fmt.Errorf("GetBlockHashByHeight, can not find any index records")
	}
	hashBs, err := cstates.GetValueFromRawStorageItem(hashStore)
	if err != nil {
		return nil, fmt.Errorf("GetBlockHashByHeight, deserialize blockHashBytes from raw storage item err:%v", err)
	}

	hash := new(chainhash.Hash)
	err = hash.SetBytes(hashBs)
	if err != nil {
		return nil, fmt.Errorf("GetBlockHashByHeight at height = %d, error:%v", height, err)
	}
	return hash, nil
}

func GetHeaderByHash(native *native.NativeContract, chainID uint64, hash chainhash.Hash) (*StoredHeader, error) {
	contract := utils.HeaderSyncContractAddress

	headerStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(scom.BLOCK_HEADER), utils.GetUint64Bytes(chainID), hash.CloneBytes()))
	if err != nil {
		return nil, fmt.Errorf("GetHeaderByHash, get hashBlockHeaderStore error: %v", err)
	}
	if headerStore == nil {
		return nil, fmt.Errorf("GetHeaderByHash, can not find any index records")
	}
	shBs, err := cstates.GetValueFromRawStorageItem(headerStore)
	if err != nil {
		return nil, fmt.Errorf("GetHeaderByHash, deserialize blockHashBytes from raw storage item err: %v", err)
	}

	sh := new(StoredHeader)
	if err := sh.Deserialization(polycomm.NewZeroCopySource(shBs)); err != nil {
		return nil, fmt.Errorf("GetStoredHeader, deserializeHeader error: %v", err)
	}

	return sh, nil
}

func GetHeaderByHeight(native *native.NativeContract, chainID uint64, height uint32) (*StoredHeader, error) {
	blockHash, err := GetBlockHashByHeight(native, chainID, height)
	if err != nil {
		return nil, fmt.Errorf("GetHeaderByHeight, error: %v", err)
	}
	storedHeader, err := GetHeaderByHash(native, chainID, *blockHash)
	if err != nil {
		return nil, fmt.Errorf("GetHeaderByHeight, error: %v", err)
	}
	return storedHeader, nil
}

func GetBestBlockHeader(native *native.NativeContract, chainID uint64) (*StoredHeader, error) {
	contract := utils.HeaderSyncContractAddress

	bestBlockHeaderStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, []byte(scom.CURRENT_HEADER_HEIGHT), utils.GetUint64Bytes(chainID)))
	if err != nil {
		return nil, fmt.Errorf("GetBestBlockHeader, get BestBlockHeader error: %v", err)
	}
	if bestBlockHeaderStore == nil {
		return nil, fmt.Errorf("GetBestBlockHeader, can not find any index records")
	}
	bestBlockHeaderBs, err := cstates.GetValueFromRawStorageItem(bestBlockHeaderStore)
	if err != nil {
		return nil, fmt.Errorf("GetBestBlockHeader, deserialize bestBlockHeaderBytes from raw storage item err: %v", err)
	}
	bestBlockHeader := new(StoredHeader)
	err = bestBlockHeader.Deserialization(polycomm.NewZeroCopySource(bestBlockHeaderBs))
	if err != nil {
		return nil, fmt.Errorf("GetBestBlockHeader, deserialize storedHeader error: %v", err)
	}
	return bestBlockHeader, nil
}
