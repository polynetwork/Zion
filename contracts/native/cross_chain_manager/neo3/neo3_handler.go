package neo3

import (
	"fmt"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
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

// MakeDepositProposal ...
func (h *Handler) MakeDepositProposal(native *native.NativeContract) (*scom.MakeTxParam, error) {
	ctx := native.ContractRef().CurrentContext()
	param := &scom.EntranceParam{}
	if err := utils.UnpackMethod(scom.ABI, scom.MethodImportOuterTransfer, param, ctx.Payload); err != nil {
		return nil, err
	}

	// deserialize neo3 state root
	ccm, err := DeserializeCrossChainMsg(param.HeaderOrCrossChainMsg)
	if err != nil {
		return nil, fmt.Errorf("neo3 MakeDepositProposal, DeserializeCrossChainMsg error: %v", err)
	}
	sideChain, err := side_chain_manager.GetSideChain(native, param.SourceChainID)
	if err != nil {
		return nil, fmt.Errorf("neo3 MakeDepositProposal, side_chain_manager.GetSideChain error: %v", err)
	}
	magicNum := helper.BytesToUInt32(sideChain.ExtraInfo)
	// verify its signature
	err = verifyCrossChainMsg(native, ccm, magicNum)
	if err != nil {
		return nil, fmt.Errorf("neo3 MakeDepositProposal, VerifyCrossChainMsg error: %v", err)
	}

	// when register neo N3, convert ccmc id to []byte
	// change neo3 contract address bytes to id, it is different from other chains
	// need to store "int" in a []byte, contract id is available from "getcontractstate" api
	// neo3 native contracts have negative ids, while custom contracts have positive ones
	id := int(int32(helper.BytesToUInt32(sideChain.CCMCAddress)))

	// get cross chain param
	makeTxParam, err := verifyFromNeoTx(param.Proof, ccm, id)
	if err != nil {
		return nil, fmt.Errorf("neo3 MakeDepositProposal, verifyFromNeoTx error: %v", err)
	}
	// check done tx
	if err := scom.CheckDoneTx(native, makeTxParam.CrossChainID, param.SourceChainID); err != nil {
		return nil, fmt.Errorf("neo3 MakeDepositProposal, CheckDoneTx error:%s", err)
	}
	if err = scom.PutDoneTx(native, makeTxParam.CrossChainID, param.SourceChainID); err != nil {
		return nil, fmt.Errorf("neo3 MakeDepositProposal, PutDoneTx error:%s", err)
	}
	return makeTxParam, nil
}
