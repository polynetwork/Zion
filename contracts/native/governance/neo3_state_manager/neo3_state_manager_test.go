package neo3_state_manager

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func init() {
	InitNeo3StateManager()
	node_manager.InitNodeManager()
	db := rawdb.NewMemoryDatabase()
	sdb, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)

	putPeerMapPoolAndView(sdb)
}

func putPeerMapPoolAndView(db *state.StateDB) {
	height := uint64(120)
	epoch := node_manager.GenerateTestEpochInfo(1, height, 4)
	peer := epoch.Peers.List[0]
	rawPubKey, _ := hexutil.Decode(peer.PubKey)
	pubkey, _ := crypto.DecompressPubkey(rawPubKey)
	acct = pubkey
	caller := peer.Address

	txhash := common.HexToHash("0x1234")
	ref := native.NewContractRef(db, caller, caller, new(big.Int).SetUint64(height), txhash, 0, nil)
	s := native.NewNativeContract(db, ref)
	node_manager.StoreTestEpoch(s, epoch)
}

var (
	sdb  *state.StateDB
	acct *ecdsa.PublicKey
)

func TestRegisterStateValidator(t *testing.T) {
	{
		param := new(StateValidatorListParam)
		param.StateValidators = []string{
			"039b45040cc529966165ef5dff3d046a4960520ce616ae170e265d669e0e2de7f4",
			"0345e2bbda8d3d9e24d1e9ee61df15d4f435f69a44fe012d86e9cf9377baaa42cd",
			"023ccd59ec0fda27844984876ef2d440eca88e45c7401110210f7760cdcc73b5f7",
			"0392fbd1d809a3c62f7dcde8f25454a1570830a21e4b014b3f362a79baf413e115",
		}

		input, err := utils.PackMethodWithStruct(ABI, MethodRegisterStateValidator, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		contractRef := native.NewContractRef(sdb, common.Address{}, common.Address{}, blockNumber,
			common.Hash{}, gasTable[MethodRegisterStateValidator]+extra, nil)

		ret, leftOverGas, err := contractRef.NativeCall(common.Address{},
			utils.Neo3StateManagerContractAddress, input)
		assert.Nil(t, err)

		result, err := utils.PackOutputs(ABI, MethodRegisterStateValidator, true)
		assert.Nil(t, err)
		assert.Equal(t, result, ret)
		assert.Equal(t, extra, leftOverGas)

		contract := native.NewNativeContract(sdb, contractRef)
		svListParam, err := getStateValidatorApply(contract, 0)
		assert.Nil(t, err)
		assert.Equal(t, param.StateValidators[0], svListParam.StateValidators[0])
	}

	// none consensus acct should not be able to approve register sv
	{
		caller := crypto.PubkeyToAddress(*acct)
		param := new(ApproveStateValidatorParam)
		param.ID = 0
		param.Address = caller

		input, err := utils.PackMethodWithStruct(ABI, MethodApproveRegisterStateValidator, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber,
			common.Hash{}, gasTable[MethodApproveRegisterStateValidator]+extra, nil)

		ret, leftOverGas, err := contractRef.NativeCall(caller,
			utils.Neo3StateManagerContractAddress, input)
		assert.Nil(t, err)

		result, err := utils.PackOutputs(ABI, MethodApproveRegisterStateValidator, true)
		assert.Nil(t, err)
		assert.Equal(t, result, ret)
		assert.Equal(t, extra, leftOverGas)

		contract := native.NewNativeContract(sdb, contractRef)
		ok, err := node_manager.CheckConsensusSigns(contract, MethodApproveRegisterStateValidator, utils.GetUint64Bytes(0), caller)
		assert.NotNil(t, err, "duplicate signer") // the signer is already stored
		assert.Equal(t, false, ok)
	}
}

func TestRemoveStateValidator(t *testing.T) {
	{
		param := new(StateValidatorListParam)
		param.StateValidators = []string{
			"039b45040cc529966165ef5dff3d046a4960520ce616ae170e265d669e0e2de7f4",
			"0345e2bbda8d3d9e24d1e9ee61df15d4f435f69a44fe012d86e9cf9377baaa42cd",
			"023ccd59ec0fda27844984876ef2d440eca88e45c7401110210f7760cdcc73b5f7",
			"0392fbd1d809a3c62f7dcde8f25454a1570830a21e4b014b3f362a79baf413e115",
		}

		input, err := utils.PackMethodWithStruct(ABI, MethodRemoveStateValidator, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		contractRef := native.NewContractRef(sdb, common.Address{}, common.Address{}, blockNumber,
			common.Hash{}, gasTable[MethodRemoveStateValidator]+extra, nil)

		ret, leftOverGas, err := contractRef.NativeCall(common.Address{},
			utils.Neo3StateManagerContractAddress, input)
		assert.Nil(t, err)

		result, err := utils.PackOutputs(ABI, MethodRemoveStateValidator, true)
		assert.Nil(t, err)
		assert.Equal(t, result, ret)
		assert.Equal(t, extra, leftOverGas)

		contract := native.NewNativeContract(sdb, contractRef)
		svListParam, err := getStateValidatorRemove(contract, 0)
		assert.Nil(t, err)
		assert.Equal(t, param.StateValidators[0], svListParam.StateValidators[0])
	}

	{
		caller := crypto.PubkeyToAddress(*acct)
		param := new(ApproveStateValidatorParam)
		param.ID = 0
		param.Address = caller

		input, err := utils.PackMethodWithStruct(ABI, MethodApproveRemoveStateValidator, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber,
			common.Hash{}, gasTable[MethodApproveRegisterStateValidator]+extra, nil)

		ret, leftOverGas, err := contractRef.NativeCall(caller,
			utils.Neo3StateManagerContractAddress, input)
		assert.Nil(t, err)

		result, err := utils.PackOutputs(ABI, MethodApproveRemoveStateValidator, true)
		assert.Nil(t, err)
		assert.Equal(t, result, ret)
		assert.Equal(t, extra, leftOverGas)

		contract := native.NewNativeContract(sdb, contractRef)
		ok, err := node_manager.CheckConsensusSigns(contract, MethodApproveRemoveStateValidator, utils.GetUint64Bytes(0), caller)
		assert.NotNil(t, err, "duplicate signer") // the signer is already stored
		assert.Equal(t, false, ok)
	}
}
