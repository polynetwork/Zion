package governance

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const (
	EventAddValidator = "addValidator"
)

const abijson = `[
	{"type":"function","constant":true,"name":"` + MethodContractName + `","inputs":[],"outputs":[{"name":"_name","type":"string"}],"payable":false,"stateMutability":"view"},
	{"type":"function","constant":true,"name":"` + MethodGetEpoch + `","inputs":[],"outputs":[{"name":"_epoch","type":"uint256"}],"payable":false,"stateMutability":"view"},
	{"type":"function","constant":true,"name":"` + MethodAddValidator + `","inputs":[{"name":"validator","type":"address"}],"outputs":[{"name":"succeed","type":"bool"}]},
	{"type":"function","constant":true,"name":"` + MethodGetValidators + `","inputs":[],"outputs":[{"name":"list","type":"address[]"}]},
	{"type":"event","anonymous":false,"name":"` + EventAddValidator + `","inputs":[{"indexed":false,"name":"validator","type":"address"},{"indexed":false,"name":"succeed","type":"bool"}]}
]`

func GetABI() abi.ABI {
	ab, err := abi.JSON(strings.NewReader(abijson))
	if err != nil {
		panic(fmt.Sprintf("failed to load PLT abi json string: [%v]", err))
	}
	return ab
}

type MethodNameInput struct{}
type MethodNameOutput struct {
	Name string
}

type MethodEpochInput struct{}
type MethodEpochOutput struct {
	Epoch *big.Int
}

// validators
type MethodAddValidatorInput struct {
	Validator common.Address
}
type MethodAddValidatorOutput struct {
	Succeed bool
}

type MethodGetValidatorsInput struct{}
type MethodGetValidatorsOutput struct {
	List []common.Address
}