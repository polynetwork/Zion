package neo3_state_manager

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "neo3 state manager"

const (
	MethodContractName                  = "name"
	MethodGetCurrentStateValidator      = "getCurrentStateValidator"
	MethodRegisterStateValidator        = "registerStateValidator"
	MethodApproveRegisterStateValidator = "approveRegisterStateValidator"
	MethodRemoveStateValidator          = "removeStateValidator"
	MethodApproveRemoveStateValidator   = "approveRemoveStateValidator"
)

var (
	this     = native.NativeContractAddrMap[native.NativeNeo3StateManager]
	gasTable = map[string]uint64{
		MethodContractName:                  0,
		MethodGetCurrentStateValidator:      0,
		MethodRegisterStateValidator:        100000,
		MethodApproveRegisterStateValidator: 0,
		MethodRemoveStateValidator:          0,
		MethodApproveRemoveStateValidator:   0,
	}

	ABI abi.ABI
)

func InitNeo3StateManager() {
	ABI = GetABI()
	native.Contracts[this] = RegisterNeo3StateManagerContract
}

func RegisterNeo3StateManagerContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodContractName, Name)
	s.Register(MethodGetCurrentStateValidator, GetCurrentStateValidator)
	s.Register(MethodRegisterStateValidator, RegisterStateValidator)
	s.Register(MethodApproveRegisterStateValidator, ApproveRegisterStateValidator)
	s.Register(MethodRemoveStateValidator, RemoveStateValidator)
	s.Register(MethodApproveRemoveStateValidator, ApproveRemoveStateValidator)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(ABI, MethodContractName, contractName)
}

func GetCurrentStateValidator(s *native.NativeContract) ([]byte, error) {

	return utils.PackOutputs(ABI, MethodGetCurrentStateValidator, []byte{})
}

func RegisterStateValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &StateValidatorListParam{}
	if err := utils.UnpackMethod(ABI, MethodRegisterStateValidator, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodRegisterStateValidator, true)
}

func ApproveRegisterStateValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ApproveStateValidatorParam{}
	if err := utils.UnpackMethod(ABI, MethodApproveRegisterStateValidator, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodApproveRegisterStateValidator, true)
}

func RemoveStateValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &StateValidatorListParam{}
	if err := utils.UnpackMethod(ABI, MethodRemoveStateValidator, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodRemoveStateValidator, true)
}

func ApproveRemoveStateValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &ApproveStateValidatorParam{}
	if err := utils.UnpackMethod(ABI, MethodApproveRemoveStateValidator, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodApproveRemoveStateValidator, true)
}
