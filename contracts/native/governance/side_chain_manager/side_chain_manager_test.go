package side_chain_manager

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/stretchr/testify/assert"
)

func init() {
	InitSideChainManager()
}

func TestRegisterSideChainManager(t *testing.T) {
	param := new(RegisterSideChainParam)
	param.BlocksToWait = 4
	param.ChainId = 8
	param.Name = "mychain"
	param.Router = 3

	input, err := utils.PackMethodWithStruct(ABI, MethodRegisterSideChain, param)
	assert.Nil(t, err)
	db := rawdb.NewMemoryDatabase()
	sdb, _ := state.New(common.Hash{}, state.NewDatabase(db), nil)
	blockNumber := big.NewInt(1)
	extra := uint64(10)
	contractRef := native.NewContractRef(sdb, common.Address{}, blockNumber, gasTable[MethodRegisterSideChain]+extra, nil)
	ret, leftOverGas, err := contractRef.NativeCall(common.Address{}, utils.SideChainManagerContractAddress, input)

	assert.Nil(t, err)

	result, err := utils.PackOutputs(ABI, MethodRegisterSideChain, true)
	assert.Nil(t, err)
	assert.Equal(t, ret, result)
	assert.Equal(t, leftOverGas, extra)

	contract := native.NewNativeContract(sdb, contractRef)
	sideChain, err := getSideChainApply(contract, 8)
	assert.Equal(t, sideChain.Name, "mychain")
	assert.Nil(t, err)

	_, _, err = contractRef.NativeCall(common.Address{}, utils.SideChainManagerContractAddress, input)
	assert.NotNil(t, err)
}

func TestApproveRegisterSideChain(t *testing.T) {
	param := new(ChainidParam)
	param.Chainid = 8

	input, err := utils.PackMethodWithStruct(ABI, MethodApproveRegisterSideChain, param)
	assert.Nil(t, err)

	db := rawdb.NewMemoryDatabase()
	sdb, _ := state.New(common.Hash{}, state.NewDatabase(db), nil)
	blockNumber := big.NewInt(1)
	extra := uint64(10)
	contractRef := native.NewContractRef(sdb, common.Address{}, blockNumber, gasTable[MethodApproveRegisterSideChain]+extra, nil)
	ret, leftOverGas, err := contractRef.NativeCall(common.Address{}, utils.SideChainManagerContractAddress, input)

	assert.Nil(t, err)

	result, err := utils.PackOutputs(ABI, MethodApproveRegisterSideChain, true)
	assert.Nil(t, err)
	assert.Equal(t, ret, result)
	assert.Equal(t, leftOverGas, extra)
}
