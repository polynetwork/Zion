package neo3

import (
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	scom "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/joeqian10/neo3-gogogo/helper"
)

// Handler ...
type Handler struct {
}

// NewHandler ...
func NewHandler() *Handler {
	return &Handler{}
}

func (this *Handler) SyncGenesisHeader(native *native.NativeContract) error {
	ctx := native.ContractRef().CurrentContext()
	param := &scom.SyncGenesisHeaderParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodSyncGenesisHeader, param, ctx.Payload); err != nil {
		return fmt.Errorf("neo3 SyncGenesisHeader, contract params deserialize error: %v", err)
	}

	// check consensus signs
	ok, err := node_manager.CheckConsensusSigns(native, scom.MethodSyncGenesisHeader, ctx.Payload, native.ContractRef().MsgSender())
	if err != nil {
		return fmt.Errorf("neo3 SyncGenesisHeader, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return nil
	}

	// can only store once
	stored, err := isGenesisStored(native, param)
	if err != nil {
		return fmt.Errorf("neo3 SyncGenesisHeader, isGenesisStored error: %v", err)
	}
	if stored {
		return fmt.Errorf("neo3 SyncGenesisHeader, genesis header had been initialized")
	}

	// deserialize genesis header to verify
	genesis, err := DeserializeNeo3Header(param.GenesisHeader)
	if err != nil {
		return fmt.Errorf("neo3 SyncGenesisHeader, DeserializeNeo3Header error: %v", err)
	}

	err = putGenesis(native, param, genesis)
	if err != nil {
		return fmt.Errorf("neo3 SyncGenesisHeader, putGenesis error: %v", err)
	}

	scom.NotifyPutHeader(native, param.ChainID, uint64(genesis.GetIndex()), genesis.GetHashString())
	return nil
}

func (this *Handler) SyncBlockHeader(native *native.NativeContract) error {
	param := &scom.SyncBlockHeaderParam{}
	ctx := native.ContractRef().CurrentContext()
	if err := utils.UnpackMethod(scom.ABI, scom.MethodSyncBlockHeader, param, ctx.Payload); err != nil {
		return err
	}

	sideChain, err := side_chain_manager.GetSideChain(native, param.ChainID)
	if err != nil {
		return fmt.Errorf("neo3 SyncBlockHeader, GetSideChain error: %v", err)
	}
	magicNum := helper.BytesToUInt32(sideChain.ExtraInfo)

	for _, v := range param.Headers {
		header, err := DeserializeNeo3Header(v)
		if err != nil {
			return fmt.Errorf("neo3 SyncBlockHeader, DeserializeNeo3Header error: %v", err)
		}
		// get last stored header
		lastHeader, err := getLastStoredHeader(native, param.ChainID)
		if err != nil {
			return fmt.Errorf("neo3 SyncBlockHeader, getLastStoredHeader error: %v", err)
		}
		// verify new header
		if !header.GetNextConsensus().Equals(lastHeader.GetNextConsensus()) &&
			header.GetIndex() > lastHeader.GetIndex() {
			if err = verifyHeader(native, header, lastHeader.GetNextConsensus(), magicNum); err != nil {
				return fmt.Errorf("neo3 SyncBlockHeader, verifyHeader error: %v", err)
			}
		}
		// update the last stored header
		err = putHeader(native, param, header)

		scom.NotifyPutHeader(native, param.ChainID, uint64(header.GetIndex()), header.GetHash().String())
	}
	return nil
}

// SyncCrossChainMsg ...
func (h *Handler) SyncCrossChainMsg(native *native.NativeContract) error {
	return nil
}
