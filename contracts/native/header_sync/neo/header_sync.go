package neo

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	hscommon "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/polynetwork/poly/common"
)

type NEOHandler struct {
}

func NewNEOHandler() *NEOHandler {
	return &NEOHandler{}
}

func (this *NEOHandler) SyncGenesisHeader(native *native.NativeContract) error {
	ctx := native.ContractRef().CurrentContext()
	params := &hscommon.SyncGenesisHeaderParam{}
	if err := utils.UnpackMethod(hscommon.ABI, hscommon.MethodSyncGenesisHeader, params, ctx.Payload); err != nil {
		return fmt.Errorf("SyncGenesisHeader, contract params deserialize error: %v", err)
	}
	// Get current epoch operator
	ok, err := node_manager.CheckConsensusSigns(native, hscommon.MethodSyncGenesisHeader, ctx.Payload, native.ContractRef().MsgSender())
	if err != nil {
		return fmt.Errorf("SyncGenesisHeader, CheckConsensusSigns error: %v", err)
	}
	if !ok {
		return nil
	}

	// Deserialize neo block header
	var header NeoBlockHeader
	err = json.Unmarshal(params.GenesisHeader, &header)
	if err != nil {
		return fmt.Errorf("SyncGenesisHeader, json.Unmarshal header err: %v", err)
	}
	if neoConsensus, _ := getConsensusValByChainId(native, params.ChainID); neoConsensus == nil {
		// Put NeoConsensus.NextConsensus into storage
		if err = putConsensusValByChainId(native, &NeoConsensus{
			ChainID:       params.ChainID,
			Height:        header.Index,
			NextConsensus: header.NextConsensus,
		}); err != nil {
			return fmt.Errorf("NeoHandler SyncGenesisHeader, update ConsensusPeer error: %v", err)
		}
	}
	return nil
}

func (this *NEOHandler) SyncBlockHeader(native *native.NativeContract) error {
	params := &hscommon.SyncBlockHeaderParam{}
	{
		ctx := native.ContractRef().CurrentContext()
		if err := utils.UnpackMethod(hscommon.ABI, hscommon.MethodSyncBlockHeader, params, ctx.Payload); err != nil {
			return err
		}
	}
	neoConsensus, err := getConsensusValByChainId(native, params.ChainID)
	if err != nil {
		return fmt.Errorf("SyncBlockHeader, the consensus validator has not been initialized, chainId: %d", params.ChainID)
	}
	var newNeoConsensus *NeoConsensus
	for _, v := range params.Headers {
		header := new(NeoBlockHeader)
		if err := header.Deserialization(common.NewZeroCopySource(v)); err != nil {
			return fmt.Errorf("SyncBlockHeader, NeoBlockHeaderFromBytes error: %v", err)
		}
		if !header.NextConsensus.Equals(neoConsensus.NextConsensus) && header.Index > neoConsensus.Height {
			if err = verifyHeader(native, params.ChainID, header); err != nil {
				return fmt.Errorf("SyncBlockHeader, verifyHeader error: %v", err)
			}
			newNeoConsensus = &NeoConsensus{
				ChainID:       neoConsensus.ChainID,
				Height:        header.Index,
				NextConsensus: header.NextConsensus,
			}
		}
	}
	if newNeoConsensus != nil {
		if err = putConsensusValByChainId(native, newNeoConsensus); err != nil {
			return fmt.Errorf("SyncBlockHeader, update ConsensusPeer error: %v", err)
		}
	}
	return nil
}

func (this *NEOHandler) SyncCrossChainMsg(native *native.NativeContract) error {
	return nil
}
