package governance

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

// todo: design and implement.

const contractName = "Zion governance"

const (
	MethodContractName  = "name"
	MethodGetEpoch      = "epoch"
	MethodAddValidator  = "addValidator"
	MethodGetValidators = "validators"
)

var (
	this = native.NativeContractAddrMap[native.NativeGovernance]

	gasTable = map[string]uint64{
		MethodContractName:  0,
		MethodGetEpoch:      0,
		MethodAddValidator:  100000,
		MethodGetValidators: 0,
	}

	ABI abi.ABI
)

func InitGovernance() {
	ABI = GetABI()
	native.Contracts[this] = RegisterGovernanceContract
}

func RegisterGovernanceContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodContractName, Name)
	s.Register(MethodGetEpoch, GetEpoch)
	s.Register(MethodAddValidator, AddValidator)
	s.Register(MethodGetValidators, GetValidators)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(ABI, MethodContractName, contractName)
}

func GetEpoch(s *native.NativeContract) ([]byte, error) {
	testEpoch := big.NewInt(1)
	return utils.PackOutputs(ABI, MethodGetEpoch, testEpoch)
}

func AddValidator(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &MethodAddValidatorInput{}
	if err := utils.UnpackMethod(ABI, MethodAddValidator, params, ctx.Payload); err != nil {
		return nil, err
	}

	emitAddValidator(s, params.Validator, true)
	return utils.PackOutputs(ABI, MethodAddValidator, true)
}

// todo: genesis nodes as validators in the first epoch
func GetValidators(s *native.NativeContract) ([]byte, error) {
	return nil, nil
}

func emitAddValidator(s *native.NativeContract, validator common.Address, succeed bool) {
	topics := make([]common.Hash, 2)
	topics[0] = ABI.Events[EventAddValidator].ID
	topics[1] = utils.Address2Hash(validator)
	data := utils.Bool2Bytes(succeed)
	emitter := utils.NewEventEmitter(this, s.ContractRef().BlockHeight().Uint64(), s.StateDB())
	emitter.Event(topics, data)
}