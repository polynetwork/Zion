package eth

import (
	"fmt"

	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

type ETHHandler struct {
}

func NewETHHandler() *ETHHandler {
	return &ETHHandler{}
}

func (this *ETHHandler) MakeDepositProposal(service *native.NativeContract) (*scom.MakeTxParam, error) {
	ctx := service.ContractRef().CurrentContext()
	params := &scom.EntranceParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodImportOuterTransfer, params, ctx.Payload); err != nil {
		return nil, err
	}

	sideChain, err := side_chain_manager.GetSideChain(service, params.SourceChainID)
	if err != nil {
		return nil, fmt.Errorf("eth MakeDepositProposal, side_chain_manager.GetSideChain error: %v", err)
	}

	value, err := verifyFromEthTx(service, params.Proof, params.Extra, params.SourceChainID, params.Height, sideChain)
	if err != nil {
		return nil, fmt.Errorf("eth MakeDepositProposal, verifyFromEthTx error: %s", err)
	}
	if err := scom.CheckDoneTx(service, value.CrossChainID, params.SourceChainID); err != nil {
		return nil, fmt.Errorf("eth MakeDepositProposal, check done transaction error:%s", err)
	}
	if err := scom.PutDoneTx(service, value.CrossChainID, params.SourceChainID); err != nil {
		return nil, fmt.Errorf("eth MakeDepositProposal, PutDoneTx error:%s", err)
	}
	return value, nil
}
