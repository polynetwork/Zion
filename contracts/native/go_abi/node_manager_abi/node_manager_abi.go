// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package node_manager_abi

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

var (
	MethodPropose = "propose"

	MethodVote = "vote"

	MethodEpoch = "epoch"

	MethodGetChangingEpoch = "getChangingEpoch"

	MethodGetEpochByID = "getEpochByID"

	MethodName = "name"

	MethodProof = "proof"

	EventConsensusSigned = "ConsensusSigned"

	EventEpochChanged = "EpochChanged"

	EventProposed = "Proposed"

	EventVoted = "Voted"
)

// INodeManagerABI is the input ABI used to generate the binding from.
const INodeManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"method\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"size\",\"type\":\"uint64\"}],\"name\":\"ConsensusSigned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"epoch\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"nextEpoch\",\"type\":\"bytes\"}],\"name\":\"EpochChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"epoch\",\"type\":\"bytes\"}],\"name\":\"Proposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"epochID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"epochHash\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"votedNumber\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"groupSize\",\"type\":\"uint64\"}],\"name\":\"Voted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"epoch\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChangingEpoch\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"epochID\",\"type\":\"uint64\"}],\"name\":\"getEpochByID\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"epochID\",\"type\":\"uint64\"}],\"name\":\"proof\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"startHeight\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"peers\",\"type\":\"bytes\"}],\"name\":\"propose\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"epochID\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"epochHash\",\"type\":\"bytes\"}],\"name\":\"vote\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// INodeManagerFuncSigs maps the 4-byte function signature to its string representation.
var INodeManagerFuncSigs = map[string]string{
	"900cf0cf": "epoch()",
	"76b85cd9": "getChangingEpoch()",
	"b9dda35e": "getEpochByID(uint64)",
	"06fdde03": "name()",
	"418f9899": "proof(uint64)",
	"bcc12328": "propose(uint64,bytes)",
	"08c16dbb": "vote(uint64,bytes)",
}

// INodeManager is an auto generated Go binding around an Ethereum contract.
type INodeManager struct {
	INodeManagerCaller     // Read-only binding to the contract
	INodeManagerTransactor // Write-only binding to the contract
	INodeManagerFilterer   // Log filterer for contract events
}

// INodeManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type INodeManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// INodeManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type INodeManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// INodeManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type INodeManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// INodeManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type INodeManagerSession struct {
	Contract     *INodeManager     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// INodeManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type INodeManagerCallerSession struct {
	Contract *INodeManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// INodeManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type INodeManagerTransactorSession struct {
	Contract     *INodeManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// INodeManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type INodeManagerRaw struct {
	Contract *INodeManager // Generic contract binding to access the raw methods on
}

// INodeManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type INodeManagerCallerRaw struct {
	Contract *INodeManagerCaller // Generic read-only contract binding to access the raw methods on
}

// INodeManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type INodeManagerTransactorRaw struct {
	Contract *INodeManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewINodeManager creates a new instance of INodeManager, bound to a specific deployed contract.
func NewINodeManager(address common.Address, backend bind.ContractBackend) (*INodeManager, error) {
	contract, err := bindINodeManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &INodeManager{INodeManagerCaller: INodeManagerCaller{contract: contract}, INodeManagerTransactor: INodeManagerTransactor{contract: contract}, INodeManagerFilterer: INodeManagerFilterer{contract: contract}}, nil
}

// NewINodeManagerCaller creates a new read-only instance of INodeManager, bound to a specific deployed contract.
func NewINodeManagerCaller(address common.Address, caller bind.ContractCaller) (*INodeManagerCaller, error) {
	contract, err := bindINodeManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &INodeManagerCaller{contract: contract}, nil
}

// NewINodeManagerTransactor creates a new write-only instance of INodeManager, bound to a specific deployed contract.
func NewINodeManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*INodeManagerTransactor, error) {
	contract, err := bindINodeManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &INodeManagerTransactor{contract: contract}, nil
}

// NewINodeManagerFilterer creates a new log filterer instance of INodeManager, bound to a specific deployed contract.
func NewINodeManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*INodeManagerFilterer, error) {
	contract, err := bindINodeManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &INodeManagerFilterer{contract: contract}, nil
}

// bindINodeManager binds a generic wrapper to an already deployed contract.
func bindINodeManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(INodeManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_INodeManager *INodeManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _INodeManager.Contract.INodeManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_INodeManager *INodeManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _INodeManager.Contract.INodeManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_INodeManager *INodeManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _INodeManager.Contract.INodeManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_INodeManager *INodeManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _INodeManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_INodeManager *INodeManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _INodeManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_INodeManager *INodeManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _INodeManager.Contract.contract.Transact(opts, method, params...)
}

// Epoch is a free data retrieval call binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() view returns(bytes)
func (_INodeManager *INodeManagerCaller) Epoch(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "epoch")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Epoch is a free data retrieval call binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() view returns(bytes)
func (_INodeManager *INodeManagerSession) Epoch() ([]byte, error) {
	return _INodeManager.Contract.Epoch(&_INodeManager.CallOpts)
}

// Epoch is a free data retrieval call binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) Epoch() ([]byte, error) {
	return _INodeManager.Contract.Epoch(&_INodeManager.CallOpts)
}

// GetChangingEpoch is a free data retrieval call binding the contract method 0x76b85cd9.
//
// Solidity: function getChangingEpoch() view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetChangingEpoch(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getChangingEpoch")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetChangingEpoch is a free data retrieval call binding the contract method 0x76b85cd9.
//
// Solidity: function getChangingEpoch() view returns(bytes)
func (_INodeManager *INodeManagerSession) GetChangingEpoch() ([]byte, error) {
	return _INodeManager.Contract.GetChangingEpoch(&_INodeManager.CallOpts)
}

// GetChangingEpoch is a free data retrieval call binding the contract method 0x76b85cd9.
//
// Solidity: function getChangingEpoch() view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetChangingEpoch() ([]byte, error) {
	return _INodeManager.Contract.GetChangingEpoch(&_INodeManager.CallOpts)
}

// GetEpochByID is a free data retrieval call binding the contract method 0xb9dda35e.
//
// Solidity: function getEpochByID(uint64 epochID) view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetEpochByID(opts *bind.CallOpts, epochID uint64) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getEpochByID", epochID)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetEpochByID is a free data retrieval call binding the contract method 0xb9dda35e.
//
// Solidity: function getEpochByID(uint64 epochID) view returns(bytes)
func (_INodeManager *INodeManagerSession) GetEpochByID(epochID uint64) ([]byte, error) {
	return _INodeManager.Contract.GetEpochByID(&_INodeManager.CallOpts, epochID)
}

// GetEpochByID is a free data retrieval call binding the contract method 0xb9dda35e.
//
// Solidity: function getEpochByID(uint64 epochID) view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetEpochByID(epochID uint64) ([]byte, error) {
	return _INodeManager.Contract.GetEpochByID(&_INodeManager.CallOpts, epochID)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_INodeManager *INodeManagerCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_INodeManager *INodeManagerSession) Name() (string, error) {
	return _INodeManager.Contract.Name(&_INodeManager.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_INodeManager *INodeManagerCallerSession) Name() (string, error) {
	return _INodeManager.Contract.Name(&_INodeManager.CallOpts)
}

// Proof is a free data retrieval call binding the contract method 0x418f9899.
//
// Solidity: function proof(uint64 epochID) view returns(bytes)
func (_INodeManager *INodeManagerCaller) Proof(opts *bind.CallOpts, epochID uint64) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "proof", epochID)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Proof is a free data retrieval call binding the contract method 0x418f9899.
//
// Solidity: function proof(uint64 epochID) view returns(bytes)
func (_INodeManager *INodeManagerSession) Proof(epochID uint64) ([]byte, error) {
	return _INodeManager.Contract.Proof(&_INodeManager.CallOpts, epochID)
}

// Proof is a free data retrieval call binding the contract method 0x418f9899.
//
// Solidity: function proof(uint64 epochID) view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) Proof(epochID uint64) ([]byte, error) {
	return _INodeManager.Contract.Proof(&_INodeManager.CallOpts, epochID)
}

// Propose is a paid mutator transaction binding the contract method 0xbcc12328.
//
// Solidity: function propose(uint64 startHeight, bytes peers) returns(bool)
func (_INodeManager *INodeManagerTransactor) Propose(opts *bind.TransactOpts, startHeight uint64, peers []byte) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "propose", startHeight, peers)
}

// Propose is a paid mutator transaction binding the contract method 0xbcc12328.
//
// Solidity: function propose(uint64 startHeight, bytes peers) returns(bool)
func (_INodeManager *INodeManagerSession) Propose(startHeight uint64, peers []byte) (*types.Transaction, error) {
	return _INodeManager.Contract.Propose(&_INodeManager.TransactOpts, startHeight, peers)
}

// Propose is a paid mutator transaction binding the contract method 0xbcc12328.
//
// Solidity: function propose(uint64 startHeight, bytes peers) returns(bool)
func (_INodeManager *INodeManagerTransactorSession) Propose(startHeight uint64, peers []byte) (*types.Transaction, error) {
	return _INodeManager.Contract.Propose(&_INodeManager.TransactOpts, startHeight, peers)
}

// Vote is a paid mutator transaction binding the contract method 0x08c16dbb.
//
// Solidity: function vote(uint64 epochID, bytes epochHash) returns(bool)
func (_INodeManager *INodeManagerTransactor) Vote(opts *bind.TransactOpts, epochID uint64, epochHash []byte) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "vote", epochID, epochHash)
}

// Vote is a paid mutator transaction binding the contract method 0x08c16dbb.
//
// Solidity: function vote(uint64 epochID, bytes epochHash) returns(bool)
func (_INodeManager *INodeManagerSession) Vote(epochID uint64, epochHash []byte) (*types.Transaction, error) {
	return _INodeManager.Contract.Vote(&_INodeManager.TransactOpts, epochID, epochHash)
}

// Vote is a paid mutator transaction binding the contract method 0x08c16dbb.
//
// Solidity: function vote(uint64 epochID, bytes epochHash) returns(bool)
func (_INodeManager *INodeManagerTransactorSession) Vote(epochID uint64, epochHash []byte) (*types.Transaction, error) {
	return _INodeManager.Contract.Vote(&_INodeManager.TransactOpts, epochID, epochHash)
}

// INodeManagerConsensusSignedIterator is returned from FilterConsensusSigned and is used to iterate over the raw logs and unpacked data for ConsensusSigned events raised by the INodeManager contract.
type INodeManagerConsensusSignedIterator struct {
	Event *INodeManagerConsensusSigned // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *INodeManagerConsensusSignedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(INodeManagerConsensusSigned)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(INodeManagerConsensusSigned)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *INodeManagerConsensusSignedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *INodeManagerConsensusSignedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// INodeManagerConsensusSigned represents a ConsensusSigned event raised by the INodeManager contract.
type INodeManagerConsensusSigned struct {
	Method string
	Input  []byte
	Signer common.Address
	Size   uint64
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterConsensusSigned is a free log retrieval operation binding the contract event 0x0061afebf4fdedb651e1607bf3b25a3b5073565ab6424ca51a4e66bd632b15ce.
//
// Solidity: event ConsensusSigned(string method, bytes input, address signer, uint64 size)
func (_INodeManager *INodeManagerFilterer) FilterConsensusSigned(opts *bind.FilterOpts) (*INodeManagerConsensusSignedIterator, error) {

	logs, sub, err := _INodeManager.contract.FilterLogs(opts, "ConsensusSigned")
	if err != nil {
		return nil, err
	}
	return &INodeManagerConsensusSignedIterator{contract: _INodeManager.contract, event: "ConsensusSigned", logs: logs, sub: sub}, nil
}

// WatchConsensusSigned is a free log subscription operation binding the contract event 0x0061afebf4fdedb651e1607bf3b25a3b5073565ab6424ca51a4e66bd632b15ce.
//
// Solidity: event ConsensusSigned(string method, bytes input, address signer, uint64 size)
func (_INodeManager *INodeManagerFilterer) WatchConsensusSigned(opts *bind.WatchOpts, sink chan<- *INodeManagerConsensusSigned) (event.Subscription, error) {

	logs, sub, err := _INodeManager.contract.WatchLogs(opts, "ConsensusSigned")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(INodeManagerConsensusSigned)
				if err := _INodeManager.contract.UnpackLog(event, "ConsensusSigned", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConsensusSigned is a log parse operation binding the contract event 0x0061afebf4fdedb651e1607bf3b25a3b5073565ab6424ca51a4e66bd632b15ce.
//
// Solidity: event ConsensusSigned(string method, bytes input, address signer, uint64 size)
func (_INodeManager *INodeManagerFilterer) ParseConsensusSigned(log types.Log) (*INodeManagerConsensusSigned, error) {
	event := new(INodeManagerConsensusSigned)
	if err := _INodeManager.contract.UnpackLog(event, "ConsensusSigned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// INodeManagerEpochChangedIterator is returned from FilterEpochChanged and is used to iterate over the raw logs and unpacked data for EpochChanged events raised by the INodeManager contract.
type INodeManagerEpochChangedIterator struct {
	Event *INodeManagerEpochChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *INodeManagerEpochChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(INodeManagerEpochChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(INodeManagerEpochChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *INodeManagerEpochChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *INodeManagerEpochChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// INodeManagerEpochChanged represents a EpochChanged event raised by the INodeManager contract.
type INodeManagerEpochChanged struct {
	Epoch     []byte
	NextEpoch []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterEpochChanged is a free log retrieval operation binding the contract event 0x5bad90814f3890720ae7e64921fad1d8441d38991c42dd2958fba86418bc120c.
//
// Solidity: event EpochChanged(bytes epoch, bytes nextEpoch)
func (_INodeManager *INodeManagerFilterer) FilterEpochChanged(opts *bind.FilterOpts) (*INodeManagerEpochChangedIterator, error) {

	logs, sub, err := _INodeManager.contract.FilterLogs(opts, "EpochChanged")
	if err != nil {
		return nil, err
	}
	return &INodeManagerEpochChangedIterator{contract: _INodeManager.contract, event: "EpochChanged", logs: logs, sub: sub}, nil
}

// WatchEpochChanged is a free log subscription operation binding the contract event 0x5bad90814f3890720ae7e64921fad1d8441d38991c42dd2958fba86418bc120c.
//
// Solidity: event EpochChanged(bytes epoch, bytes nextEpoch)
func (_INodeManager *INodeManagerFilterer) WatchEpochChanged(opts *bind.WatchOpts, sink chan<- *INodeManagerEpochChanged) (event.Subscription, error) {

	logs, sub, err := _INodeManager.contract.WatchLogs(opts, "EpochChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(INodeManagerEpochChanged)
				if err := _INodeManager.contract.UnpackLog(event, "EpochChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEpochChanged is a log parse operation binding the contract event 0x5bad90814f3890720ae7e64921fad1d8441d38991c42dd2958fba86418bc120c.
//
// Solidity: event EpochChanged(bytes epoch, bytes nextEpoch)
func (_INodeManager *INodeManagerFilterer) ParseEpochChanged(log types.Log) (*INodeManagerEpochChanged, error) {
	event := new(INodeManagerEpochChanged)
	if err := _INodeManager.contract.UnpackLog(event, "EpochChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// INodeManagerProposedIterator is returned from FilterProposed and is used to iterate over the raw logs and unpacked data for Proposed events raised by the INodeManager contract.
type INodeManagerProposedIterator struct {
	Event *INodeManagerProposed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *INodeManagerProposedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(INodeManagerProposed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(INodeManagerProposed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *INodeManagerProposedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *INodeManagerProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// INodeManagerProposed represents a Proposed event raised by the INodeManager contract.
type INodeManagerProposed struct {
	Epoch []byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterProposed is a free log retrieval operation binding the contract event 0x10b2060c55406ea48522476f67fd813d4984b12078555d3e2a377e35839d7d01.
//
// Solidity: event Proposed(bytes epoch)
func (_INodeManager *INodeManagerFilterer) FilterProposed(opts *bind.FilterOpts) (*INodeManagerProposedIterator, error) {

	logs, sub, err := _INodeManager.contract.FilterLogs(opts, "Proposed")
	if err != nil {
		return nil, err
	}
	return &INodeManagerProposedIterator{contract: _INodeManager.contract, event: "Proposed", logs: logs, sub: sub}, nil
}

// WatchProposed is a free log subscription operation binding the contract event 0x10b2060c55406ea48522476f67fd813d4984b12078555d3e2a377e35839d7d01.
//
// Solidity: event Proposed(bytes epoch)
func (_INodeManager *INodeManagerFilterer) WatchProposed(opts *bind.WatchOpts, sink chan<- *INodeManagerProposed) (event.Subscription, error) {

	logs, sub, err := _INodeManager.contract.WatchLogs(opts, "Proposed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(INodeManagerProposed)
				if err := _INodeManager.contract.UnpackLog(event, "Proposed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseProposed is a log parse operation binding the contract event 0x10b2060c55406ea48522476f67fd813d4984b12078555d3e2a377e35839d7d01.
//
// Solidity: event Proposed(bytes epoch)
func (_INodeManager *INodeManagerFilterer) ParseProposed(log types.Log) (*INodeManagerProposed, error) {
	event := new(INodeManagerProposed)
	if err := _INodeManager.contract.UnpackLog(event, "Proposed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// INodeManagerVotedIterator is returned from FilterVoted and is used to iterate over the raw logs and unpacked data for Voted events raised by the INodeManager contract.
type INodeManagerVotedIterator struct {
	Event *INodeManagerVoted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *INodeManagerVotedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(INodeManagerVoted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(INodeManagerVoted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *INodeManagerVotedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *INodeManagerVotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// INodeManagerVoted represents a Voted event raised by the INodeManager contract.
type INodeManagerVoted struct {
	EpochID     uint64
	EpochHash   []byte
	VotedNumber uint64
	GroupSize   uint64
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterVoted is a free log retrieval operation binding the contract event 0x2c3818322730b87b7a73fc5d82355f808c91822a7c4244a5ddd54e8cfd5903cd.
//
// Solidity: event Voted(uint64 epochID, bytes epochHash, uint64 votedNumber, uint64 groupSize)
func (_INodeManager *INodeManagerFilterer) FilterVoted(opts *bind.FilterOpts) (*INodeManagerVotedIterator, error) {

	logs, sub, err := _INodeManager.contract.FilterLogs(opts, "Voted")
	if err != nil {
		return nil, err
	}
	return &INodeManagerVotedIterator{contract: _INodeManager.contract, event: "Voted", logs: logs, sub: sub}, nil
}

// WatchVoted is a free log subscription operation binding the contract event 0x2c3818322730b87b7a73fc5d82355f808c91822a7c4244a5ddd54e8cfd5903cd.
//
// Solidity: event Voted(uint64 epochID, bytes epochHash, uint64 votedNumber, uint64 groupSize)
func (_INodeManager *INodeManagerFilterer) WatchVoted(opts *bind.WatchOpts, sink chan<- *INodeManagerVoted) (event.Subscription, error) {

	logs, sub, err := _INodeManager.contract.WatchLogs(opts, "Voted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(INodeManagerVoted)
				if err := _INodeManager.contract.UnpackLog(event, "Voted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVoted is a log parse operation binding the contract event 0x2c3818322730b87b7a73fc5d82355f808c91822a7c4244a5ddd54e8cfd5903cd.
//
// Solidity: event Voted(uint64 epochID, bytes epochHash, uint64 votedNumber, uint64 groupSize)
func (_INodeManager *INodeManagerFilterer) ParseVoted(log types.Log) (*INodeManagerVoted, error) {
	event := new(INodeManagerVoted)
	if err := _INodeManager.contract.UnpackLog(event, "Voted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
