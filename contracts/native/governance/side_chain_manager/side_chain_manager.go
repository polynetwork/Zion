package side_chain_manager

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "side chain manager"

const (
	MethodContractName             = "name"
	MethodRegisterSideChain        = "registerSideChain"
	MethodApproveRegisterSideChain = "approveRegisterSideChain"
	MethodUpdateSideChain          = "updateSideChain"
	MethodApproveUpdateSideChain   = "approveUpdateSideChain"
	MethodQuitSideChain            = "quitSideChain"
	MethodApproveQuitSideChain     = "approveQuitSideChain"
	MethodRegisterRedeem           = "registerRedeem"
	MethodSetBtcTxParam            = "setBtcTxParam"
)

var (
	this     = native.NativeContractAddrMap[native.NativeSideChainManager]
	gasTable = map[string]uint64{
		MethodContractName:             0,
		MethodRegisterSideChain:        0,
		MethodApproveRegisterSideChain: 100000,
		MethodUpdateSideChain:          0,
		MethodApproveUpdateSideChain:   0,
		MethodQuitSideChain:            0,
		MethodApproveQuitSideChain:     0,
		MethodRegisterRedeem:           0,
		MethodSetBtcTxParam:            0,
	}

	ABI abi.ABI
)

func InitSideChainManager() {
	ABI = GetABI()
	native.Contracts[this] = RegisterSideChainManagerContract
}

func RegisterSideChainManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodContractName, Name)
	s.Register(MethodRegisterSideChain, RegisterSideChain)
	s.Register(MethodApproveRegisterSideChain, ApproveRegisterSideChain)
	s.Register(MethodUpdateSideChain, UpdateSideChain)
	s.Register(MethodApproveUpdateSideChain, ApproveUpdateSideChain)
	s.Register(MethodQuitSideChain, QuitSideChain)
	s.Register(MethodApproveQuitSideChain, ApproveQuitSideChain)
	s.Register(MethodRegisterRedeem, RegisterRedeem)
	s.Register(MethodSetBtcTxParam, SetBtcTxParam)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(ABI, MethodContractName, contractName)
}

func RegisterSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &RegisterSideChainParam{}
	if err := utils.UnpackMethod(ABI, MethodRegisterSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodRegisterSideChain, true)
}

func ApproveRegisterSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ChainidParam{}
	if err := utils.UnpackMethod(ABI, MethodApproveRegisterSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodApproveRegisterSideChain, true)
}

func UpdateSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &RegisterSideChainParam{}
	if err := utils.UnpackMethod(ABI, MethodUpdateSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodUpdateSideChain, true)
}

func ApproveUpdateSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ChainidParam{}
	if err := utils.UnpackMethod(ABI, MethodApproveUpdateSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodApproveUpdateSideChain, true)
}

func QuitSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ChainidParam{}
	if err := utils.UnpackMethod(ABI, MethodQuitSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodQuitSideChain, true)
}

func ApproveQuitSideChain(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ChainidParam{}
	if err := utils.UnpackMethod(ABI, MethodApproveQuitSideChain, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodApproveQuitSideChain, true)
}

func RegisterRedeem(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &RegisterRedeemParam{}
	if err := utils.UnpackMethod(ABI, MethodRegisterRedeem, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodRegisterRedeem, true)
}

func SetBtcTxParam(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &BtcTxParam{}
	if err := utils.UnpackMethod(ABI, MethodRegisterRedeem, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodRegisterRedeem, true)
}
