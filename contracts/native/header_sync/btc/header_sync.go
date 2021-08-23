package btc

import (
	"github.com/ethereum/go-ethereum/contracts/native"
	hscommon "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

type BTCHandler struct {
}

func NewBTCHandler() *BTCHandler {
	return &BTCHandler{}
}

func (this *BTCHandler) SyncGenesisHeader(native *native.NativeContract) error {
	ctx := native.ContractRef().CurrentContext()
	params := &hscommon.SyncCrossChainMsgParam{}
	if err := utils.UnpackMethod(hscommon.ABI, hscommon.MethodSyncCrossChainMsg, params, ctx.Payload); err != nil {
		return err
	}

	return nil

}

func (this *BTCHandler) SyncBlockHeader(native *native.NativeContract) error {
	return nil
}
