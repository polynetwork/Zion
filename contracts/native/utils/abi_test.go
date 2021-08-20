package utils

import (
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestABIMethod(t *testing.T) {
	name := "addValidator"
	abijson := `[
	{"type":"function","constant":true,"name":"` + name + `","inputs":[{"name":"validator","type":"address"},{"name":"stakeAccount","type":"address"},{"name":"revoke","type":"bool"}],"outputs":[{"name":"list","type":"address[]"},{"name":"number","type":"uint256"}, {"name":"succeed","type":"bool"}]}
]`

	ab, err := abi.JSON(strings.NewReader(abijson))
	assert.NoError(t, err)

	type Input struct {
		Validator    common.Address
		StakeAccount common.Address
		Revoke       bool
	}
	type Output struct {
		List    []common.Address
		Number  *big.Int
		Succeed bool
	}
	expectInput := &Input{
		Validator:    common.HexToAddress("0x02"),
		StakeAccount: common.HexToAddress("0x03"),
		Revoke:       true,
	}
	expectOutput := &Output{
		List: []common.Address{
			common.HexToAddress("0x23"),
			common.HexToAddress("0x25"),
			common.HexToAddress("0x37"),
		},
		Number:  big.NewInt(123456789),
		Succeed: true,
	}
	// test input
	payload, err := PackMethod(&ab, name, expectInput.Validator, expectInput.StakeAccount, expectInput.Revoke)
	assert.NoError(t, err)

	inputData := &Input{}
	err = UnpackMethod(&ab, name, inputData, payload)
	assert.NoError(t, err)

	assert.True(t, reflect.DeepEqual(expectInput, inputData))

	payload, err = PackOutputs(&ab, name, expectOutput.List, expectOutput.Number, expectOutput.Succeed)
	assert.NoError(t, err)

	outputData := &Output{}
	err = UnpackOutputs(&ab, name, outputData, payload)
	assert.NoError(t, err)

	assert.True(t, reflect.DeepEqual(expectOutput, outputData))
}
