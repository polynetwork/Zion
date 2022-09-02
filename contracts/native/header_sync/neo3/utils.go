package neo3

import (
	"fmt"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/tx"
)

func isGenesisStored(native *native.NativeContract, param *scom.SyncGenesisHeaderParam) (bool, error) {
	genesis, err := getGenesis(native, param.ChainID)
	if err != nil {
		return false, fmt.Errorf("getGenesis error: %v", err)
	}

	return genesis != nil, nil
}

func getGenesis(native *native.NativeContract, chainID uint64) (*BlockHeader, error) {
	genesisBytes, err := scom.GetGenesisHeader(native, chainID)
	if err != nil {
		return nil, fmt.Errorf("GetGenesisHeader error: %v", err)
	}

	if genesisBytes == nil {
		return nil, nil
	}

	genesisHeader, err := DeserializeNeo3Header(genesisBytes)
	if err != nil {
		return nil, fmt.Errorf("getGenesis, DeserializeNeo3Header error: %v", err)
	}

	return genesisHeader, nil
}

func putGenesis(native *native.NativeContract, param *scom.SyncGenesisHeaderParam, genesis *BlockHeader) error {
	raw, err := SerializeNeo3Header(genesis)
	if err != nil {
		return fmt.Errorf("SerializeNeo3Header error: %v", err)
	}

	scom.SetGenesisHeader(native, param.ChainID, raw)
	scom.SetHeaderIndex(native, param.ChainID, genesis.GetHash().ToByteArray(), raw)
	scom.SetCurrentHeight(native, param.ChainID, utils.GetUint64Bytes(uint64(genesis.GetIndex())))
	scom.SetMainChain(native, param.ChainID, uint64(genesis.GetIndex()), raw)

	return nil
}

func getLastStoredHeader(native *native.NativeContract, chainID uint64) (*BlockHeader, error) {
	currentHeightBytes, err := scom.GetCurrentHeight(native, chainID)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentHeight error: %v", err)
	}

	currentHeight := utils.GetBytesUint64(currentHeightBytes)
	raw, err := scom.GetMainChain(native, chainID, currentHeight)
	if err != nil {
		return nil, fmt.Errorf("GetMainChain error: %v", err)
	}

	header, err := DeserializeNeo3Header(raw)
	if err != nil {
		return nil, fmt.Errorf("DeserializeNeo3Header error: %v", err)
	}
	return header, nil
}

func putHeader(native *native.NativeContract, param *scom.SyncBlockHeaderParam, header *BlockHeader) error {
	raw, err := SerializeNeo3Header(header)
	if err != nil {
		return fmt.Errorf("SerializeNeo3Header error: %v", err)
	}

	scom.SetHeaderIndex(native, param.ChainID, header.GetHash().ToByteArray(), raw)
	scom.SetCurrentHeight(native, param.ChainID, utils.GetUint64Bytes(uint64(header.GetIndex())))
	scom.SetMainChain(native, param.ChainID, uint64(header.GetIndex()), raw)

	return nil
}

func verifyHeader(native *native.NativeContract, header *BlockHeader, lastNeoConsensus *helper.UInt160, magic uint32) error {
	if !header.GetNextConsensus().Equals(lastNeoConsensus) {
		return fmt.Errorf("invalid witness script hash, expected: %s, got: %s", lastNeoConsensus.String(), header.GetNextConsensus().String())
	}
	// get hash
	msg, err := header.GetMessage(magic)
	if err != nil {
		return fmt.Errorf("header.GetMessage error: %v", err)
	}
	// verify witness
	if verified := tx.VerifyMultiSignatureWitness(msg, header.Witness); !verified {
		return fmt.Errorf("VerifyMultiSignatureWitness error: %v, height: %d", err, header.GetIndex())
	}
	return nil
}
