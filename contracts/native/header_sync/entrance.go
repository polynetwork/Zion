package header_sync

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

const contractName = "header sync"

const (
	MethodContractName      = "name"
	MethodSyncGenesisHeader = "syncGenesisHeader"
	MethodSyncBlockHeader   = "syncBlockHeader"
	MethodSyncCrossChainMsg = "syncCrossChainMsg"
)

var (
	this     = native.NativeContractAddrMap[native.NativeSyncHeader]
	gasTable = map[string]uint64{
		MethodContractName:      0,
		MethodSyncGenesisHeader: 0,
		MethodSyncBlockHeader:   100000,
		MethodSyncCrossChainMsg: 0,
	}

	ABI abi.ABI
)

func InitHeaderSync() {
	ABI = GetABI()
	native.Contracts[this] = RegisterHeaderSyncContract
}

func RegisterHeaderSyncContract(s *native.NativeContract) {
	s.Prepare(ABI, gasTable)

	s.Register(MethodContractName, Name)
	s.Register(MethodSyncGenesisHeader, SyncGenesisHeader)
	s.Register(MethodSyncBlockHeader, SyncBlockHeader)
	s.Register(MethodSyncCrossChainMsg, SyncCrossChainMsg)
}

func Name(s *native.NativeContract) ([]byte, error) {
	return utils.PackOutputs(ABI, MethodContractName, contractName)
}

func SyncGenesisHeader(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &SyncGenesisHeaderParam{}
	if err := utils.UnpackMethod(ABI, MethodSyncGenesisHeader, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodSyncGenesisHeader, true)
}

func SyncBlockHeader(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &SyncBlockHeaderParam{}
	if err := utils.UnpackMethod(ABI, MethodSyncBlockHeader, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodSyncBlockHeader, true)
}

func SyncCrossChainMsg(s *native.NativeContract) ([]byte, error) {
	ctx := s.ContractRef().CurrentContext()
	params := &SyncCrossChainMsgParam{}
	if err := utils.UnpackMethod(ABI, MethodSyncCrossChainMsg, params, ctx.Payload); err != nil {
		return nil, err
	}

	return utils.PackOutputs(ABI, MethodSyncCrossChainMsg, true)
}
