// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package proposal_manager_abi

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
	MethodUpdateNodeManagerGlobalConfig = "updateNodeManagerGlobalConfig"

	EventUpdateNodeManagerGlobalConfig = "UpdateNodeManagerGlobalConfig"
)

// IProposalManagerABI is the input ABI used to generate the binding from.
const IProposalManagerABI = "[{\"anonymous\":false,\"inputs\":[],\"name\":\"UpdateNodeManagerGlobalConfig\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"maxCommissionChange\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"minInitialStake\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"maxDescLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPerEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"consensusValidatorNum\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"voterValidatorNum\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expireHeight\",\"type\":\"uint256\"}],\"name\":\"updateNodeManagerGlobalConfig\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IProposalManagerFuncSigs maps the 4-byte function signature to its string representation.
var IProposalManagerFuncSigs = map[string]string{
	"f42e78c0": "updateNodeManagerGlobalConfig(string,string,uint256,uint256,uint256,uint256,uint256)",
}

// IProposalManager is an auto generated Go binding around an Ethereum contract.
type IProposalManager struct {
	IProposalManagerCaller     // Read-only binding to the contract
	IProposalManagerTransactor // Write-only binding to the contract
	IProposalManagerFilterer   // Log filterer for contract events
}

// IProposalManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type IProposalManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IProposalManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IProposalManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IProposalManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IProposalManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IProposalManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IProposalManagerSession struct {
	Contract     *IProposalManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IProposalManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IProposalManagerCallerSession struct {
	Contract *IProposalManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// IProposalManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IProposalManagerTransactorSession struct {
	Contract     *IProposalManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// IProposalManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type IProposalManagerRaw struct {
	Contract *IProposalManager // Generic contract binding to access the raw methods on
}

// IProposalManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IProposalManagerCallerRaw struct {
	Contract *IProposalManagerCaller // Generic read-only contract binding to access the raw methods on
}

// IProposalManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IProposalManagerTransactorRaw struct {
	Contract *IProposalManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIProposalManager creates a new instance of IProposalManager, bound to a specific deployed contract.
func NewIProposalManager(address common.Address, backend bind.ContractBackend) (*IProposalManager, error) {
	contract, err := bindIProposalManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IProposalManager{IProposalManagerCaller: IProposalManagerCaller{contract: contract}, IProposalManagerTransactor: IProposalManagerTransactor{contract: contract}, IProposalManagerFilterer: IProposalManagerFilterer{contract: contract}}, nil
}

// NewIProposalManagerCaller creates a new read-only instance of IProposalManager, bound to a specific deployed contract.
func NewIProposalManagerCaller(address common.Address, caller bind.ContractCaller) (*IProposalManagerCaller, error) {
	contract, err := bindIProposalManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IProposalManagerCaller{contract: contract}, nil
}

// NewIProposalManagerTransactor creates a new write-only instance of IProposalManager, bound to a specific deployed contract.
func NewIProposalManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*IProposalManagerTransactor, error) {
	contract, err := bindIProposalManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IProposalManagerTransactor{contract: contract}, nil
}

// NewIProposalManagerFilterer creates a new log filterer instance of IProposalManager, bound to a specific deployed contract.
func NewIProposalManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*IProposalManagerFilterer, error) {
	contract, err := bindIProposalManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IProposalManagerFilterer{contract: contract}, nil
}

// bindIProposalManager binds a generic wrapper to an already deployed contract.
func bindIProposalManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IProposalManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IProposalManager *IProposalManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IProposalManager.Contract.IProposalManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IProposalManager *IProposalManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IProposalManager.Contract.IProposalManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IProposalManager *IProposalManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IProposalManager.Contract.IProposalManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IProposalManager *IProposalManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IProposalManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IProposalManager *IProposalManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IProposalManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IProposalManager *IProposalManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IProposalManager.Contract.contract.Transact(opts, method, params...)
}

// UpdateNodeManagerGlobalConfig is a paid mutator transaction binding the contract method 0xf42e78c0.
//
// Solidity: function updateNodeManagerGlobalConfig(string maxCommissionChange, string minInitialStake, uint256 maxDescLength, uint256 blockPerEpoch, uint256 consensusValidatorNum, uint256 voterValidatorNum, uint256 expireHeight) returns(bool success)
func (_IProposalManager *IProposalManagerTransactor) UpdateNodeManagerGlobalConfig(opts *bind.TransactOpts, maxCommissionChange string, minInitialStake string, maxDescLength *big.Int, blockPerEpoch *big.Int, consensusValidatorNum *big.Int, voterValidatorNum *big.Int, expireHeight *big.Int) (*types.Transaction, error) {
	return _IProposalManager.contract.Transact(opts, "updateNodeManagerGlobalConfig", maxCommissionChange, minInitialStake, maxDescLength, blockPerEpoch, consensusValidatorNum, voterValidatorNum, expireHeight)
}

// UpdateNodeManagerGlobalConfig is a paid mutator transaction binding the contract method 0xf42e78c0.
//
// Solidity: function updateNodeManagerGlobalConfig(string maxCommissionChange, string minInitialStake, uint256 maxDescLength, uint256 blockPerEpoch, uint256 consensusValidatorNum, uint256 voterValidatorNum, uint256 expireHeight) returns(bool success)
func (_IProposalManager *IProposalManagerSession) UpdateNodeManagerGlobalConfig(maxCommissionChange string, minInitialStake string, maxDescLength *big.Int, blockPerEpoch *big.Int, consensusValidatorNum *big.Int, voterValidatorNum *big.Int, expireHeight *big.Int) (*types.Transaction, error) {
	return _IProposalManager.Contract.UpdateNodeManagerGlobalConfig(&_IProposalManager.TransactOpts, maxCommissionChange, minInitialStake, maxDescLength, blockPerEpoch, consensusValidatorNum, voterValidatorNum, expireHeight)
}

// UpdateNodeManagerGlobalConfig is a paid mutator transaction binding the contract method 0xf42e78c0.
//
// Solidity: function updateNodeManagerGlobalConfig(string maxCommissionChange, string minInitialStake, uint256 maxDescLength, uint256 blockPerEpoch, uint256 consensusValidatorNum, uint256 voterValidatorNum, uint256 expireHeight) returns(bool success)
func (_IProposalManager *IProposalManagerTransactorSession) UpdateNodeManagerGlobalConfig(maxCommissionChange string, minInitialStake string, maxDescLength *big.Int, blockPerEpoch *big.Int, consensusValidatorNum *big.Int, voterValidatorNum *big.Int, expireHeight *big.Int) (*types.Transaction, error) {
	return _IProposalManager.Contract.UpdateNodeManagerGlobalConfig(&_IProposalManager.TransactOpts, maxCommissionChange, minInitialStake, maxDescLength, blockPerEpoch, consensusValidatorNum, voterValidatorNum, expireHeight)
}

// IProposalManagerUpdateNodeManagerGlobalConfigIterator is returned from FilterUpdateNodeManagerGlobalConfig and is used to iterate over the raw logs and unpacked data for UpdateNodeManagerGlobalConfig events raised by the IProposalManager contract.
type IProposalManagerUpdateNodeManagerGlobalConfigIterator struct {
	Event *IProposalManagerUpdateNodeManagerGlobalConfig // Event containing the contract specifics and raw log

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
func (it *IProposalManagerUpdateNodeManagerGlobalConfigIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IProposalManagerUpdateNodeManagerGlobalConfig)
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
		it.Event = new(IProposalManagerUpdateNodeManagerGlobalConfig)
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
func (it *IProposalManagerUpdateNodeManagerGlobalConfigIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IProposalManagerUpdateNodeManagerGlobalConfigIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IProposalManagerUpdateNodeManagerGlobalConfig represents a UpdateNodeManagerGlobalConfig event raised by the IProposalManager contract.
type IProposalManagerUpdateNodeManagerGlobalConfig struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUpdateNodeManagerGlobalConfig is a free log retrieval operation binding the contract event 0x34862e998fa4124289470e7e6e3904de60f85a5f7616cc8dcdce7b2a693584d1.
//
// Solidity: event UpdateNodeManagerGlobalConfig()
func (_IProposalManager *IProposalManagerFilterer) FilterUpdateNodeManagerGlobalConfig(opts *bind.FilterOpts) (*IProposalManagerUpdateNodeManagerGlobalConfigIterator, error) {

	logs, sub, err := _IProposalManager.contract.FilterLogs(opts, "UpdateNodeManagerGlobalConfig")
	if err != nil {
		return nil, err
	}
	return &IProposalManagerUpdateNodeManagerGlobalConfigIterator{contract: _IProposalManager.contract, event: "UpdateNodeManagerGlobalConfig", logs: logs, sub: sub}, nil
}

// WatchUpdateNodeManagerGlobalConfig is a free log subscription operation binding the contract event 0x34862e998fa4124289470e7e6e3904de60f85a5f7616cc8dcdce7b2a693584d1.
//
// Solidity: event UpdateNodeManagerGlobalConfig()
func (_IProposalManager *IProposalManagerFilterer) WatchUpdateNodeManagerGlobalConfig(opts *bind.WatchOpts, sink chan<- *IProposalManagerUpdateNodeManagerGlobalConfig) (event.Subscription, error) {

	logs, sub, err := _IProposalManager.contract.WatchLogs(opts, "UpdateNodeManagerGlobalConfig")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IProposalManagerUpdateNodeManagerGlobalConfig)
				if err := _IProposalManager.contract.UnpackLog(event, "UpdateNodeManagerGlobalConfig", log); err != nil {
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

// ParseUpdateNodeManagerGlobalConfig is a log parse operation binding the contract event 0x34862e998fa4124289470e7e6e3904de60f85a5f7616cc8dcdce7b2a693584d1.
//
// Solidity: event UpdateNodeManagerGlobalConfig()
func (_IProposalManager *IProposalManagerFilterer) ParseUpdateNodeManagerGlobalConfig(log types.Log) (*IProposalManagerUpdateNodeManagerGlobalConfig, error) {
	event := new(IProposalManagerUpdateNodeManagerGlobalConfig)
	if err := _IProposalManager.contract.UnpackLog(event, "UpdateNodeManagerGlobalConfig", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

