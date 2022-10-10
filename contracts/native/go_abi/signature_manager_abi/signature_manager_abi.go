// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package signature_manager_abi

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
	MethodAddSignature = "addSignature"

	EventAddSignatureQuorumEvent = "AddSignatureQuorumEvent"
)

// ISignatureManagerABI is the input ABI used to generate the binding from.
const ISignatureManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"subject\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sideChainID\",\"type\":\"uint256\"}],\"name\":\"AddSignatureQuorumEvent\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"sideChainID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"subject\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"addSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ISignatureManagerFuncSigs maps the 4-byte function signature to its string representation.
var ISignatureManagerFuncSigs = map[string]string{
	"29d75da9": "addSignature(address,uint256,bytes,bytes)",
}

// ISignatureManager is an auto generated Go binding around an Ethereum contract.
type ISignatureManager struct {
	ISignatureManagerCaller     // Read-only binding to the contract
	ISignatureManagerTransactor // Write-only binding to the contract
	ISignatureManagerFilterer   // Log filterer for contract events
}

// ISignatureManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ISignatureManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISignatureManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ISignatureManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISignatureManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ISignatureManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISignatureManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ISignatureManagerSession struct {
	Contract     *ISignatureManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ISignatureManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ISignatureManagerCallerSession struct {
	Contract *ISignatureManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ISignatureManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ISignatureManagerTransactorSession struct {
	Contract     *ISignatureManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ISignatureManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ISignatureManagerRaw struct {
	Contract *ISignatureManager // Generic contract binding to access the raw methods on
}

// ISignatureManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ISignatureManagerCallerRaw struct {
	Contract *ISignatureManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ISignatureManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ISignatureManagerTransactorRaw struct {
	Contract *ISignatureManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewISignatureManager creates a new instance of ISignatureManager, bound to a specific deployed contract.
func NewISignatureManager(address common.Address, backend bind.ContractBackend) (*ISignatureManager, error) {
	contract, err := bindISignatureManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ISignatureManager{ISignatureManagerCaller: ISignatureManagerCaller{contract: contract}, ISignatureManagerTransactor: ISignatureManagerTransactor{contract: contract}, ISignatureManagerFilterer: ISignatureManagerFilterer{contract: contract}}, nil
}

// NewISignatureManagerCaller creates a new read-only instance of ISignatureManager, bound to a specific deployed contract.
func NewISignatureManagerCaller(address common.Address, caller bind.ContractCaller) (*ISignatureManagerCaller, error) {
	contract, err := bindISignatureManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ISignatureManagerCaller{contract: contract}, nil
}

// NewISignatureManagerTransactor creates a new write-only instance of ISignatureManager, bound to a specific deployed contract.
func NewISignatureManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ISignatureManagerTransactor, error) {
	contract, err := bindISignatureManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ISignatureManagerTransactor{contract: contract}, nil
}

// NewISignatureManagerFilterer creates a new log filterer instance of ISignatureManager, bound to a specific deployed contract.
func NewISignatureManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ISignatureManagerFilterer, error) {
	contract, err := bindISignatureManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ISignatureManagerFilterer{contract: contract}, nil
}

// bindISignatureManager binds a generic wrapper to an already deployed contract.
func bindISignatureManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ISignatureManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISignatureManager *ISignatureManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISignatureManager.Contract.ISignatureManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISignatureManager *ISignatureManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISignatureManager.Contract.ISignatureManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISignatureManager *ISignatureManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISignatureManager.Contract.ISignatureManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISignatureManager *ISignatureManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISignatureManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISignatureManager *ISignatureManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISignatureManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISignatureManager *ISignatureManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISignatureManager.Contract.contract.Transact(opts, method, params...)
}

// AddSignature is a paid mutator transaction binding the contract method 0x29d75da9.
//
// Solidity: function addSignature(address addr, uint256 sideChainID, bytes subject, bytes signature) returns(bool)
func (_ISignatureManager *ISignatureManagerTransactor) AddSignature(opts *bind.TransactOpts, addr common.Address, sideChainID *big.Int, subject []byte, signature []byte) (*types.Transaction, error) {
	return _ISignatureManager.contract.Transact(opts, "addSignature", addr, sideChainID, subject, signature)
}

// AddSignature is a paid mutator transaction binding the contract method 0x29d75da9.
//
// Solidity: function addSignature(address addr, uint256 sideChainID, bytes subject, bytes signature) returns(bool)
func (_ISignatureManager *ISignatureManagerSession) AddSignature(addr common.Address, sideChainID *big.Int, subject []byte, signature []byte) (*types.Transaction, error) {
	return _ISignatureManager.Contract.AddSignature(&_ISignatureManager.TransactOpts, addr, sideChainID, subject, signature)
}

// AddSignature is a paid mutator transaction binding the contract method 0x29d75da9.
//
// Solidity: function addSignature(address addr, uint256 sideChainID, bytes subject, bytes signature) returns(bool)
func (_ISignatureManager *ISignatureManagerTransactorSession) AddSignature(addr common.Address, sideChainID *big.Int, subject []byte, signature []byte) (*types.Transaction, error) {
	return _ISignatureManager.Contract.AddSignature(&_ISignatureManager.TransactOpts, addr, sideChainID, subject, signature)
}

// ISignatureManagerAddSignatureQuorumEventIterator is returned from FilterAddSignatureQuorumEvent and is used to iterate over the raw logs and unpacked data for AddSignatureQuorumEvent events raised by the ISignatureManager contract.
type ISignatureManagerAddSignatureQuorumEventIterator struct {
	Event *ISignatureManagerAddSignatureQuorumEvent // Event containing the contract specifics and raw log

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
func (it *ISignatureManagerAddSignatureQuorumEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ISignatureManagerAddSignatureQuorumEvent)
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
		it.Event = new(ISignatureManagerAddSignatureQuorumEvent)
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
func (it *ISignatureManagerAddSignatureQuorumEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ISignatureManagerAddSignatureQuorumEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ISignatureManagerAddSignatureQuorumEvent represents a AddSignatureQuorumEvent event raised by the ISignatureManager contract.
type ISignatureManagerAddSignatureQuorumEvent struct {
	Id          []byte
	Subject     []byte
	SideChainID *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAddSignatureQuorumEvent is a free log retrieval operation binding the contract event 0xc347f5ccf7997409d4ef282db7e4c041bde177c835151101af49171c43651d42.
//
// Solidity: event AddSignatureQuorumEvent(bytes id, bytes subject, uint256 sideChainID)
func (_ISignatureManager *ISignatureManagerFilterer) FilterAddSignatureQuorumEvent(opts *bind.FilterOpts) (*ISignatureManagerAddSignatureQuorumEventIterator, error) {

	logs, sub, err := _ISignatureManager.contract.FilterLogs(opts, "AddSignatureQuorumEvent")
	if err != nil {
		return nil, err
	}
	return &ISignatureManagerAddSignatureQuorumEventIterator{contract: _ISignatureManager.contract, event: "AddSignatureQuorumEvent", logs: logs, sub: sub}, nil
}

// WatchAddSignatureQuorumEvent is a free log subscription operation binding the contract event 0xc347f5ccf7997409d4ef282db7e4c041bde177c835151101af49171c43651d42.
//
// Solidity: event AddSignatureQuorumEvent(bytes id, bytes subject, uint256 sideChainID)
func (_ISignatureManager *ISignatureManagerFilterer) WatchAddSignatureQuorumEvent(opts *bind.WatchOpts, sink chan<- *ISignatureManagerAddSignatureQuorumEvent) (event.Subscription, error) {

	logs, sub, err := _ISignatureManager.contract.WatchLogs(opts, "AddSignatureQuorumEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ISignatureManagerAddSignatureQuorumEvent)
				if err := _ISignatureManager.contract.UnpackLog(event, "AddSignatureQuorumEvent", log); err != nil {
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

// ParseAddSignatureQuorumEvent is a log parse operation binding the contract event 0xc347f5ccf7997409d4ef282db7e4c041bde177c835151101af49171c43651d42.
//
// Solidity: event AddSignatureQuorumEvent(bytes id, bytes subject, uint256 sideChainID)
func (_ISignatureManager *ISignatureManagerFilterer) ParseAddSignatureQuorumEvent(log types.Log) (*ISignatureManagerAddSignatureQuorumEvent, error) {
	event := new(ISignatureManagerAddSignatureQuorumEvent)
	if err := _ISignatureManager.contract.UnpackLog(event, "AddSignatureQuorumEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
