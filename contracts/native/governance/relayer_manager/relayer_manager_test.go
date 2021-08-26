package relayer_manager

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	polycomm "github.com/polynetwork/poly/common"
	cstates "github.com/polynetwork/poly/core/states"
	"github.com/stretchr/testify/assert"
)

func init() {
	InitRelayerManager()
	node_manager.InitNodeManager()
	db := rawdb.NewMemoryDatabase()
	sdb, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)

	cacheDB := (*state.CacheDB)(sdb)
	putPeerMapPoolAndView(cacheDB)
}

func putPeerMapPoolAndView(db *state.CacheDB) {
	key, _ := crypto.GenerateKey()
	acct = &key.PublicKey

	peerPoolMap := new(node_manager.PeerPoolMap)
	peerPoolMap.PeerPoolMap = make(map[string]*node_manager.PeerPoolItem)
	pkStr := hex.EncodeToString(crypto.FromECDSAPub(acct))
	peerPoolMap.PeerPoolMap[pkStr] = &node_manager.PeerPoolItem{
		Index:      uint32(0),
		PeerPubkey: pkStr,
		Address:    crypto.PubkeyToAddress(*acct),
		Status:     node_manager.ConsensusStatus,
	}
	view := uint32(0)
	viewBytes := utils.GetUint32Bytes(view)
	sink := polycomm.NewZeroCopySink(nil)
	peerPoolMap.Serialization(sink)
	db.Put(utils.ConcatKey(utils.NodeManagerContractAddress, []byte(node_manager.PEER_POOL), viewBytes), cstates.GenRawStorageItem(sink.Bytes()))

	sink.Reset()

	govView := node_manager.GovernanceView{
		View: view,
	}
	govView.Serialization(sink)
	db.Put(utils.ConcatKey(utils.NodeManagerContractAddress, []byte(node_manager.GOVERNANCE_VIEW)), cstates.GenRawStorageItem(sink.Bytes()))
}

var (
	sdb  *state.StateDB
	acct *ecdsa.PublicKey
)

func TestRegisterRelayer(t *testing.T) {
	{
		params := new(RelayerListParam)
		params.AddressList = []common.Address{{1, 2, 4, 6}, {1, 4, 5, 7}, {1, 3, 5, 7, 9}}

		input, err := utils.PackMethodWithStruct(ABI, MethodRegisterRelayer, params)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		contractRef := native.NewContractRef(sdb, common.Address{}, blockNumber, common.Hash{}, gasTable[MethodRegisterRelayer]+extra, nil)
		ret, leftOverGas, err := contractRef.NativeCall(common.Address{}, utils.RelayerManagerContractAddress, input)

		assert.Nil(t, err)

		result, err := utils.PackOutputs(ABI, MethodRegisterRelayer, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, extra)

		contract := native.NewNativeContract(sdb, contractRef)
		relayerListParam, err := getRelayerApply(contract, 0)
		assert.Nil(t, err)
		assert.Equal(t, params, relayerListParam)
	}

	// none consensus acct should not be able to approve register relayer
	{
		caller := crypto.PubkeyToAddress(*acct)
		param := new(ApproveRelayerParam)
		param.ID = 0
		param.Address = caller

		input, err := utils.PackMethodWithStruct(ABI, MethodApproveRegisterRelayer, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		contractRef := native.NewContractRef(sdb, caller, blockNumber, common.Hash{}, gasTable[MethodApproveRegisterRelayer]+extra, nil)
		ret, leftOverGas, err := contractRef.NativeCall(caller, utils.RelayerManagerContractAddress, input)

		assert.Nil(t, err)

		result, err := utils.PackOutputs(ABI, MethodApproveRegisterRelayer, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, extra)

		contract := native.NewNativeContract(sdb, contractRef)
		ok, err := node_manager.CheckConsensusSigns(contract, MethodApproveRegisterRelayer, utils.GetUint64Bytes(0), caller)
		assert.Nil(t, err)
		assert.Equal(t, true, ok)
	}

}

func TestRemoveRelayer(t *testing.T) {
	{
		params := new(RelayerListParam)
		params.AddressList = []common.Address{{1, 2, 4, 6}, {1, 4, 5, 7}}

		input, err := utils.PackMethodWithStruct(ABI, MethodRemoveRelayer, params)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		contractRef := native.NewContractRef(sdb, common.Address{}, blockNumber, common.Hash{}, gasTable[MethodRemoveRelayer]+extra, nil)
		ret, leftOverGas, err := contractRef.NativeCall(common.Address{}, utils.RelayerManagerContractAddress, input)

		assert.Nil(t, err)

		result, err := utils.PackOutputs(ABI, MethodRemoveRelayer, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, extra)

		contract := native.NewNativeContract(sdb, contractRef)
		relayerListParam, err := getRelayerRemove(contract, 0)
		assert.Nil(t, err)
		assert.Equal(t, params, relayerListParam)
	}

	{
		caller := crypto.PubkeyToAddress(*acct)
		param := new(ApproveRelayerParam)
		param.ID = 0
		param.Address = caller

		input, err := utils.PackMethodWithStruct(ABI, MethodApproveRemoveRelayer, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(10)
		contractRef := native.NewContractRef(sdb, caller, blockNumber, common.Hash{}, gasTable[MethodApproveRemoveRelayer]+extra, nil)
		ret, leftOverGas, err := contractRef.NativeCall(caller, utils.RelayerManagerContractAddress, input)

		assert.Nil(t, err)

		result, err := utils.PackOutputs(ABI, MethodApproveRemoveRelayer, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, extra)

		contract := native.NewNativeContract(sdb, contractRef)
		ok, err := node_manager.CheckConsensusSigns(contract, MethodApproveRemoveRelayer, utils.GetUint64Bytes(0), caller)
		assert.Nil(t, err)
		assert.Equal(t, true, ok)
	}
}
