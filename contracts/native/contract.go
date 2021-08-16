package native

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
)

type (
	RegisterService func(native *NativeContract)
	MethodHandler   func(contract *NativeContract) ([]byte, error)
)

var (
	Contracts = make(map[common.Address]RegisterService)
)

type NativeContract struct {
	ref      *ContractRef
	db       *state.StateDB
	handlers map[string]MethodHandler // map method id to method handler
	gasTable map[string]uint64        // map method id to gas usage
	ab       abi.ABI
}

func NewNativeContract(db *state.StateDB, ref *ContractRef) *NativeContract {
	return &NativeContract{
		db:       db,
		ref:      ref,
		handlers: make(map[string]MethodHandler),
	}
}

func (s *NativeContract) ContractRef() *ContractRef {
	return s.ref
}

func (s *NativeContract) GetCacheDB() *state.CacheDB {
	return (*state.CacheDB)(s.db)
}

func (s *NativeContract) StateDB() *state.StateDB {
	return s.db
}

func (s *NativeContract) Prepare(ab abi.ABI, gasTb map[string]uint64) {
	s.ab = ab
	s.gasTable = make(map[string]uint64)
	for name, gas := range gasTb {
		id := utils.MethodID(s.ab, name)
		s.gasTable[id] = gas
	}
}

func (s *NativeContract) Register(name string, handler MethodHandler) {
	methodID := utils.MethodID(s.ab, name)
	s.handlers[methodID] = handler
}

// Invoke return execute ret and cost gas
func (s *NativeContract) Invoke() ([]byte, error) {
	// check context
	if !s.ref.CheckContexts() {
		return nil, fmt.Errorf("context error")
	}
	ctx := s.ref.CurrentContext()

	// find methodID
	if len(ctx.Payload) < 4 {
		return nil, fmt.Errorf("invalid input")
	}
	methodID := hexutil.Encode(ctx.Payload[:4])

	// register methods
	registerHandler, ok := Contracts[ctx.ContractAddress]
	if !ok {
		return nil, fmt.Errorf("failed to find contract: [%x]", ctx.ContractAddress)
	}
	registerHandler(s)

	// get method handler
	handler, ok := s.handlers[methodID]
	if !ok {
		return nil, fmt.Errorf("failed to find method: [%s]", methodID)
	}

	// check gasLeft
	needGas, ok := s.gasTable[methodID]
	if !ok {
		return nil, fmt.Errorf("failed to find method: [%s]", methodID)
	}
	gasLeft := s.ref.gasLeft
	if gasLeft < needGas && gasLeft < MinGasUsage {
		return nil, fmt.Errorf("gasLeft not enough, need %d", needGas)
	}

	// execute transaction and cost gas
	ret, err := handler(s)
	if err != nil && needGas > MinGasUsage {
		needGas = MinGasUsage
	}
	if needGas > 0 {
		s.ref.gasLeft -= needGas
	}

	return ret, err
}

func (s *NativeContract) AddNotify(topics []common.Hash, data []byte) {
	emitter := utils.NewEventEmitter(s.ref.CurrentContext().ContractAddress, s.ContractRef().BlockHeight().Uint64(), s.StateDB())
	emitter.Event(topics, data)
}
