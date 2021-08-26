package governance

import (
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/stretchr/testify/assert"
)

var (
	testStateDB   *state.StateDB
	testEnv       *native.ContractRef
	testGasSupply = uint64(100000)

	testCaller = common.HexToAddress("0xab7ada5c57b3e796ec589bfe84bea3cb7ae47b63")
)

func TestMain(m *testing.M) {
	testStateDB = utils.NewTestStateDB()
	InitGovernance()

	blockNumber := big.NewInt(1)
	testEnv = native.NewContractRef(testStateDB, testCaller, blockNumber, common.Hash{}, testGasSupply, nil)
	os.Exit(m.Run())
}

func TestName(t *testing.T) {
	name := MethodContractName
	payload, err := utils.PackMethod(&ABI, name)
	assert.NoError(t, err)

	enc, gasLeft, err := testEnv.NativeCall(testCaller, this, payload)
	assert.NoError(t, err)

	expectGasLeft := uint64(testEnv.GasLeft() - gasTable[name])
	assert.Equal(t, expectGasLeft, gasLeft)

	output := new(MethodNameOutput)
	err = utils.UnpackOutputs(&ABI, name, output, enc)
	assert.NoError(t, err)
	assert.Equal(t, contractName, output.Name)
	t.Logf("left gas %d", gasLeft)
}

func TestEpoch(t *testing.T) {
	name := MethodGetEpoch

	payload, err := utils.PackMethod(&ABI, name)
	assert.NoError(t, err)

	enc, gasLeft, err := testEnv.NativeCall(testCaller, this, payload)
	assert.NoError(t, err)

	expectGasLeft := uint64(testEnv.GasLeft())
	assert.Equal(t, expectGasLeft, gasLeft)

	output := new(MethodEpochOutput)
	err = utils.UnpackOutputs(&ABI, name, output, enc)
	assert.NoError(t, err)

	assert.Equal(t, uint64(1), output.Epoch.Uint64())
}

func TestAddValidator(t *testing.T) {
	name := MethodAddValidator

	expectValidator := common.HexToAddress("0x12345")
	payload, err := utils.PackMethod(&ABI, name, expectValidator)
	assert.NoError(t, err)

	enc, gasLeft, err := testEnv.NativeCall(testCaller, this, payload)
	assert.NoError(t, err)

	expectGasLeft := uint64(testEnv.GasLeft())
	assert.Equal(t, expectGasLeft, gasLeft)

	output := new(MethodAddValidatorOutput)
	err = utils.UnpackOutputs(&ABI, name, output, enc)
	assert.NoError(t, err)
	assert.Equal(t, true, output.Succeed)

	hash := testEnv.StateDB().BlockHash()
	logs := testEnv.StateDB().GetLogs(hash)
	assert.Equal(t, len(logs), 1)

	event := logs[0]
	assert.Equal(t, 2, len(event.Topics))
	assert.Equal(t, ABI.Events[EventAddValidator].ID, event.Topics[0])
	assert.Equal(t, expectValidator, utils.Hash2Address(event.Topics[1]))
	assert.Equal(t, utils.Bool2Bytes(true), event.Data)
}
