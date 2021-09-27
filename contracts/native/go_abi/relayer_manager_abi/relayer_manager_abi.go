// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package relayer_manager_abi

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
	MethodApproveRegisterRelayer = "approveRegisterRelayer"

	MethodApproveRemoveRelayer = "approveRemoveRelayer"

	MethodName = "name"

	MethodRegisterRelayer = "registerRelayer"

	MethodRemoveRelayer = "removeRelayer"
)

// RelayerManagerABI is the input ABI used to generate the binding from.
const RelayerManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"}],\"name\":\"evtApproveRegisterRelayer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"}],\"name\":\"evtApproveRemoveRelayer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"applyID\",\"type\":\"uint64\"}],\"name\":\"evtRegisterRelayer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"removeID\",\"type\":\"uint64\"}],\"name\":\"evtRemoveRelayer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"approveRegisterRelayer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"approveRemoveRelayer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"AddressList\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"registerRelayer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"AddressList\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"removeRelayer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// RelayerManagerFuncSigs maps the 4-byte function signature to its string representation.
var RelayerManagerFuncSigs = map[string]string{
	"07b8ca31": "approveRegisterRelayer(uint64,address)",
	"2b1775dd": "approveRemoveRelayer(uint64,address)",
	"06fdde03": "name()",
	"d99802fe": "registerRelayer(address[],address)",
	"0cffb52a": "removeRelayer(address[],address)",
}

// RelayerManagerBin is the compiled bytecode used for deploying new contracts.
var RelayerManagerBin = "0x608060405234801561001057600080fd5b50610285806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806306fdde031461005c57806307b8ca31146100745780630cffb52a1461009a5780632b1775dd14610074578063d99802fe1461009a575b600080fd5b606060405161006b91906101e4565b60405180910390f35b61008a6100823660046101a2565b600092915050565b604051901515815260200161006b565b61008a6100823660046100c4565b80356001600160a01b03811681146100bf57600080fd5b919050565b600080604083850312156100d757600080fd5b823567ffffffffffffffff808211156100ef57600080fd5b818501915085601f83011261010357600080fd5b813560208282111561011757610117610239565b8160051b604051601f19603f8301168101818110868211171561013c5761013c610239565b604052838152828101945085830182870184018b101561015b57600080fd5b600096505b8487101561018557610171816100a8565b865260019690960195948301948301610160565b50965061019590508782016100a8565b9450505050509250929050565b600080604083850312156101b557600080fd5b823567ffffffffffffffff811681146101cd57600080fd5b91506101db602084016100a8565b90509250929050565b600060208083528351808285015260005b81811015610211578581018301518582016040015282016101f5565b81811115610223576000604083870101525b50601f01601f1916929092016040019392505050565b634e487b7160e01b600052604160045260246000fdfea2646970667358221220f8f573c9e4225c35f16c8b207ef9fe5a3541fabe38f3516118d08dbb43f7f78264736f6c63430008060033"

// DeployRelayerManager deploys a new Ethereum contract, binding an instance of RelayerManager to it.
func DeployRelayerManager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RelayerManager, error) {
	parsed, err := abi.JSON(strings.NewReader(RelayerManagerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RelayerManagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RelayerManager{RelayerManagerCaller: RelayerManagerCaller{contract: contract}, RelayerManagerTransactor: RelayerManagerTransactor{contract: contract}, RelayerManagerFilterer: RelayerManagerFilterer{contract: contract}}, nil
}

// RelayerManager is an auto generated Go binding around an Ethereum contract.
type RelayerManager struct {
	RelayerManagerCaller     // Read-only binding to the contract
	RelayerManagerTransactor // Write-only binding to the contract
	RelayerManagerFilterer   // Log filterer for contract events
}

// RelayerManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type RelayerManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayerManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RelayerManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayerManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RelayerManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayerManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RelayerManagerSession struct {
	Contract     *RelayerManager   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RelayerManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RelayerManagerCallerSession struct {
	Contract *RelayerManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// RelayerManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RelayerManagerTransactorSession struct {
	Contract     *RelayerManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// RelayerManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type RelayerManagerRaw struct {
	Contract *RelayerManager // Generic contract binding to access the raw methods on
}

// RelayerManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RelayerManagerCallerRaw struct {
	Contract *RelayerManagerCaller // Generic read-only contract binding to access the raw methods on
}

// RelayerManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RelayerManagerTransactorRaw struct {
	Contract *RelayerManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRelayerManager creates a new instance of RelayerManager, bound to a specific deployed contract.
func NewRelayerManager(address common.Address, backend bind.ContractBackend) (*RelayerManager, error) {
	contract, err := bindRelayerManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RelayerManager{RelayerManagerCaller: RelayerManagerCaller{contract: contract}, RelayerManagerTransactor: RelayerManagerTransactor{contract: contract}, RelayerManagerFilterer: RelayerManagerFilterer{contract: contract}}, nil
}

// NewRelayerManagerCaller creates a new read-only instance of RelayerManager, bound to a specific deployed contract.
func NewRelayerManagerCaller(address common.Address, caller bind.ContractCaller) (*RelayerManagerCaller, error) {
	contract, err := bindRelayerManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RelayerManagerCaller{contract: contract}, nil
}

// NewRelayerManagerTransactor creates a new write-only instance of RelayerManager, bound to a specific deployed contract.
func NewRelayerManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*RelayerManagerTransactor, error) {
	contract, err := bindRelayerManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RelayerManagerTransactor{contract: contract}, nil
}

// NewRelayerManagerFilterer creates a new log filterer instance of RelayerManager, bound to a specific deployed contract.
func NewRelayerManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*RelayerManagerFilterer, error) {
	contract, err := bindRelayerManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RelayerManagerFilterer{contract: contract}, nil
}

// bindRelayerManager binds a generic wrapper to an already deployed contract.
func bindRelayerManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RelayerManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RelayerManager *RelayerManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RelayerManager.Contract.RelayerManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RelayerManager *RelayerManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RelayerManager.Contract.RelayerManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RelayerManager *RelayerManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RelayerManager.Contract.RelayerManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RelayerManager *RelayerManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RelayerManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RelayerManager *RelayerManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RelayerManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RelayerManager *RelayerManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RelayerManager.Contract.contract.Transact(opts, method, params...)
}

// ApproveRegisterRelayer is a paid mutator transaction binding the contract method 0x07b8ca31.
//
// Solidity: function approveRegisterRelayer(uint64 ID, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactor) ApproveRegisterRelayer(opts *bind.TransactOpts, ID uint64, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.contract.Transact(opts, "approveRegisterRelayer", ID, Address)
}

// ApproveRegisterRelayer is a paid mutator transaction binding the contract method 0x07b8ca31.
//
// Solidity: function approveRegisterRelayer(uint64 ID, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerSession) ApproveRegisterRelayer(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.ApproveRegisterRelayer(&_RelayerManager.TransactOpts, ID, Address)
}

// ApproveRegisterRelayer is a paid mutator transaction binding the contract method 0x07b8ca31.
//
// Solidity: function approveRegisterRelayer(uint64 ID, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactorSession) ApproveRegisterRelayer(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.ApproveRegisterRelayer(&_RelayerManager.TransactOpts, ID, Address)
}

// ApproveRemoveRelayer is a paid mutator transaction binding the contract method 0x2b1775dd.
//
// Solidity: function approveRemoveRelayer(uint64 ID, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactor) ApproveRemoveRelayer(opts *bind.TransactOpts, ID uint64, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.contract.Transact(opts, "approveRemoveRelayer", ID, Address)
}

// ApproveRemoveRelayer is a paid mutator transaction binding the contract method 0x2b1775dd.
//
// Solidity: function approveRemoveRelayer(uint64 ID, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerSession) ApproveRemoveRelayer(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.ApproveRemoveRelayer(&_RelayerManager.TransactOpts, ID, Address)
}

// ApproveRemoveRelayer is a paid mutator transaction binding the contract method 0x2b1775dd.
//
// Solidity: function approveRemoveRelayer(uint64 ID, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactorSession) ApproveRemoveRelayer(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.ApproveRemoveRelayer(&_RelayerManager.TransactOpts, ID, Address)
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string Name)
func (_RelayerManager *RelayerManagerTransactor) Name(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RelayerManager.contract.Transact(opts, "name")
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string Name)
func (_RelayerManager *RelayerManagerSession) Name() (*types.Transaction, error) {
	return _RelayerManager.Contract.Name(&_RelayerManager.TransactOpts)
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string Name)
func (_RelayerManager *RelayerManagerTransactorSession) Name() (*types.Transaction, error) {
	return _RelayerManager.Contract.Name(&_RelayerManager.TransactOpts)
}

// RegisterRelayer is a paid mutator transaction binding the contract method 0xd99802fe.
//
// Solidity: function registerRelayer(address[] AddressList, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactor) RegisterRelayer(opts *bind.TransactOpts, AddressList []common.Address, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.contract.Transact(opts, "registerRelayer", AddressList, Address)
}

// RegisterRelayer is a paid mutator transaction binding the contract method 0xd99802fe.
//
// Solidity: function registerRelayer(address[] AddressList, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerSession) RegisterRelayer(AddressList []common.Address, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.RegisterRelayer(&_RelayerManager.TransactOpts, AddressList, Address)
}

// RegisterRelayer is a paid mutator transaction binding the contract method 0xd99802fe.
//
// Solidity: function registerRelayer(address[] AddressList, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactorSession) RegisterRelayer(AddressList []common.Address, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.RegisterRelayer(&_RelayerManager.TransactOpts, AddressList, Address)
}

// RemoveRelayer is a paid mutator transaction binding the contract method 0x0cffb52a.
//
// Solidity: function removeRelayer(address[] AddressList, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactor) RemoveRelayer(opts *bind.TransactOpts, AddressList []common.Address, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.contract.Transact(opts, "removeRelayer", AddressList, Address)
}

// RemoveRelayer is a paid mutator transaction binding the contract method 0x0cffb52a.
//
// Solidity: function removeRelayer(address[] AddressList, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerSession) RemoveRelayer(AddressList []common.Address, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.RemoveRelayer(&_RelayerManager.TransactOpts, AddressList, Address)
}

// RemoveRelayer is a paid mutator transaction binding the contract method 0x0cffb52a.
//
// Solidity: function removeRelayer(address[] AddressList, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactorSession) RemoveRelayer(AddressList []common.Address, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.RemoveRelayer(&_RelayerManager.TransactOpts, AddressList, Address)
}

// RelayerManagerApproveRegisterRelayerIterator is returned from FilterApproveRegisterRelayer and is used to iterate over the raw logs and unpacked data for ApproveRegisterRelayer events raised by the RelayerManager contract.
type RelayerManagerApproveRegisterRelayerIterator struct {
	Event *RelayerManagerApproveRegisterRelayer // Event containing the contract specifics and raw log

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
func (it *RelayerManagerApproveRegisterRelayerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerManagerApproveRegisterRelayer)
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
		it.Event = new(RelayerManagerApproveRegisterRelayer)
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
func (it *RelayerManagerApproveRegisterRelayerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerManagerApproveRegisterRelayerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerManagerApproveRegisterRelayer represents a ApproveRegisterRelayer event raised by the RelayerManager contract.
type RelayerManagerApproveRegisterRelayer struct {
	ID  uint64
	Raw types.Log // Blockchain specific contextual infos
}

// FilterApproveRegisterRelayer is a free log retrieval operation binding the contract event 0x85f1458b6de2aa7d6fb453821800c90d0ed717c565c2d5cdb73699b3ba657570.
//
// Solidity: event evtApproveRegisterRelayer(uint64 ID)
func (_RelayerManager *RelayerManagerFilterer) FilterApproveRegisterRelayer(opts *bind.FilterOpts) (*RelayerManagerApproveRegisterRelayerIterator, error) {

	logs, sub, err := _RelayerManager.contract.FilterLogs(opts, "evtApproveRegisterRelayer")
	if err != nil {
		return nil, err
	}
	return &RelayerManagerApproveRegisterRelayerIterator{contract: _RelayerManager.contract, event: "evtApproveRegisterRelayer", logs: logs, sub: sub}, nil
}

// WatchApproveRegisterRelayer is a free log subscription operation binding the contract event 0x85f1458b6de2aa7d6fb453821800c90d0ed717c565c2d5cdb73699b3ba657570.
//
// Solidity: event evtApproveRegisterRelayer(uint64 ID)
func (_RelayerManager *RelayerManagerFilterer) WatchApproveRegisterRelayer(opts *bind.WatchOpts, sink chan<- *RelayerManagerApproveRegisterRelayer) (event.Subscription, error) {

	logs, sub, err := _RelayerManager.contract.WatchLogs(opts, "evtApproveRegisterRelayer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerManagerApproveRegisterRelayer)
				if err := _RelayerManager.contract.UnpackLog(event, "evtApproveRegisterRelayer", log); err != nil {
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

// ParseApproveRegisterRelayer is a log parse operation binding the contract event 0x85f1458b6de2aa7d6fb453821800c90d0ed717c565c2d5cdb73699b3ba657570.
//
// Solidity: event evtApproveRegisterRelayer(uint64 ID)
func (_RelayerManager *RelayerManagerFilterer) ParseApproveRegisterRelayer(log types.Log) (*RelayerManagerApproveRegisterRelayer, error) {
	event := new(RelayerManagerApproveRegisterRelayer)
	if err := _RelayerManager.contract.UnpackLog(event, "evtApproveRegisterRelayer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayerManagerApproveRemoveRelayerIterator is returned from FilterApproveRemoveRelayer and is used to iterate over the raw logs and unpacked data for ApproveRemoveRelayer events raised by the RelayerManager contract.
type RelayerManagerApproveRemoveRelayerIterator struct {
	Event *RelayerManagerApproveRemoveRelayer // Event containing the contract specifics and raw log

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
func (it *RelayerManagerApproveRemoveRelayerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerManagerApproveRemoveRelayer)
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
		it.Event = new(RelayerManagerApproveRemoveRelayer)
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
func (it *RelayerManagerApproveRemoveRelayerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerManagerApproveRemoveRelayerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerManagerApproveRemoveRelayer represents a ApproveRemoveRelayer event raised by the RelayerManager contract.
type RelayerManagerApproveRemoveRelayer struct {
	ID  uint64
	Raw types.Log // Blockchain specific contextual infos
}

// FilterApproveRemoveRelayer is a free log retrieval operation binding the contract event 0xff19e858c848163b5e0f8038dfb0870427c43f36f7446878ee73f9bd2ec491f4.
//
// Solidity: event evtApproveRemoveRelayer(uint64 ID)
func (_RelayerManager *RelayerManagerFilterer) FilterApproveRemoveRelayer(opts *bind.FilterOpts) (*RelayerManagerApproveRemoveRelayerIterator, error) {

	logs, sub, err := _RelayerManager.contract.FilterLogs(opts, "evtApproveRemoveRelayer")
	if err != nil {
		return nil, err
	}
	return &RelayerManagerApproveRemoveRelayerIterator{contract: _RelayerManager.contract, event: "evtApproveRemoveRelayer", logs: logs, sub: sub}, nil
}

// WatchApproveRemoveRelayer is a free log subscription operation binding the contract event 0xff19e858c848163b5e0f8038dfb0870427c43f36f7446878ee73f9bd2ec491f4.
//
// Solidity: event evtApproveRemoveRelayer(uint64 ID)
func (_RelayerManager *RelayerManagerFilterer) WatchApproveRemoveRelayer(opts *bind.WatchOpts, sink chan<- *RelayerManagerApproveRemoveRelayer) (event.Subscription, error) {

	logs, sub, err := _RelayerManager.contract.WatchLogs(opts, "evtApproveRemoveRelayer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerManagerApproveRemoveRelayer)
				if err := _RelayerManager.contract.UnpackLog(event, "evtApproveRemoveRelayer", log); err != nil {
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

// ParseApproveRemoveRelayer is a log parse operation binding the contract event 0xff19e858c848163b5e0f8038dfb0870427c43f36f7446878ee73f9bd2ec491f4.
//
// Solidity: event evtApproveRemoveRelayer(uint64 ID)
func (_RelayerManager *RelayerManagerFilterer) ParseApproveRemoveRelayer(log types.Log) (*RelayerManagerApproveRemoveRelayer, error) {
	event := new(RelayerManagerApproveRemoveRelayer)
	if err := _RelayerManager.contract.UnpackLog(event, "evtApproveRemoveRelayer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayerManagerRegisterRelayerIterator is returned from FilterRegisterRelayer and is used to iterate over the raw logs and unpacked data for RegisterRelayer events raised by the RelayerManager contract.
type RelayerManagerRegisterRelayerIterator struct {
	Event *RelayerManagerRegisterRelayer // Event containing the contract specifics and raw log

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
func (it *RelayerManagerRegisterRelayerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerManagerRegisterRelayer)
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
		it.Event = new(RelayerManagerRegisterRelayer)
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
func (it *RelayerManagerRegisterRelayerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerManagerRegisterRelayerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerManagerRegisterRelayer represents a RegisterRelayer event raised by the RelayerManager contract.
type RelayerManagerRegisterRelayer struct {
	ApplyID uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRegisterRelayer is a free log retrieval operation binding the contract event 0xcde26c3470cafd5b40f89e103f29443624c7648f3751a6fb8b56fb8b25059430.
//
// Solidity: event evtRegisterRelayer(uint64 applyID)
func (_RelayerManager *RelayerManagerFilterer) FilterRegisterRelayer(opts *bind.FilterOpts) (*RelayerManagerRegisterRelayerIterator, error) {

	logs, sub, err := _RelayerManager.contract.FilterLogs(opts, "evtRegisterRelayer")
	if err != nil {
		return nil, err
	}
	return &RelayerManagerRegisterRelayerIterator{contract: _RelayerManager.contract, event: "evtRegisterRelayer", logs: logs, sub: sub}, nil
}

// WatchRegisterRelayer is a free log subscription operation binding the contract event 0xcde26c3470cafd5b40f89e103f29443624c7648f3751a6fb8b56fb8b25059430.
//
// Solidity: event evtRegisterRelayer(uint64 applyID)
func (_RelayerManager *RelayerManagerFilterer) WatchRegisterRelayer(opts *bind.WatchOpts, sink chan<- *RelayerManagerRegisterRelayer) (event.Subscription, error) {

	logs, sub, err := _RelayerManager.contract.WatchLogs(opts, "evtRegisterRelayer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerManagerRegisterRelayer)
				if err := _RelayerManager.contract.UnpackLog(event, "evtRegisterRelayer", log); err != nil {
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

// ParseRegisterRelayer is a log parse operation binding the contract event 0xcde26c3470cafd5b40f89e103f29443624c7648f3751a6fb8b56fb8b25059430.
//
// Solidity: event evtRegisterRelayer(uint64 applyID)
func (_RelayerManager *RelayerManagerFilterer) ParseRegisterRelayer(log types.Log) (*RelayerManagerRegisterRelayer, error) {
	event := new(RelayerManagerRegisterRelayer)
	if err := _RelayerManager.contract.UnpackLog(event, "evtRegisterRelayer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayerManagerRemoveRelayerIterator is returned from FilterRemoveRelayer and is used to iterate over the raw logs and unpacked data for RemoveRelayer events raised by the RelayerManager contract.
type RelayerManagerRemoveRelayerIterator struct {
	Event *RelayerManagerRemoveRelayer // Event containing the contract specifics and raw log

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
func (it *RelayerManagerRemoveRelayerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerManagerRemoveRelayer)
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
		it.Event = new(RelayerManagerRemoveRelayer)
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
func (it *RelayerManagerRemoveRelayerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerManagerRemoveRelayerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerManagerRemoveRelayer represents a RemoveRelayer event raised by the RelayerManager contract.
type RelayerManagerRemoveRelayer struct {
	RemoveID uint64
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRemoveRelayer is a free log retrieval operation binding the contract event 0x8f7667d2acd70ab373b61b6e8f28ca5259bfa2d53467f35f65bcf19041b5ba7d.
//
// Solidity: event evtRemoveRelayer(uint64 removeID)
func (_RelayerManager *RelayerManagerFilterer) FilterRemoveRelayer(opts *bind.FilterOpts) (*RelayerManagerRemoveRelayerIterator, error) {

	logs, sub, err := _RelayerManager.contract.FilterLogs(opts, "evtRemoveRelayer")
	if err != nil {
		return nil, err
	}
	return &RelayerManagerRemoveRelayerIterator{contract: _RelayerManager.contract, event: "evtRemoveRelayer", logs: logs, sub: sub}, nil
}

// WatchRemoveRelayer is a free log subscription operation binding the contract event 0x8f7667d2acd70ab373b61b6e8f28ca5259bfa2d53467f35f65bcf19041b5ba7d.
//
// Solidity: event evtRemoveRelayer(uint64 removeID)
func (_RelayerManager *RelayerManagerFilterer) WatchRemoveRelayer(opts *bind.WatchOpts, sink chan<- *RelayerManagerRemoveRelayer) (event.Subscription, error) {

	logs, sub, err := _RelayerManager.contract.WatchLogs(opts, "evtRemoveRelayer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerManagerRemoveRelayer)
				if err := _RelayerManager.contract.UnpackLog(event, "evtRemoveRelayer", log); err != nil {
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

// ParseRemoveRelayer is a log parse operation binding the contract event 0x8f7667d2acd70ab373b61b6e8f28ca5259bfa2d53467f35f65bcf19041b5ba7d.
//
// Solidity: event evtRemoveRelayer(uint64 removeID)
func (_RelayerManager *RelayerManagerFilterer) ParseRemoveRelayer(log types.Log) (*RelayerManagerRemoveRelayer, error) {
	event := new(RelayerManagerRemoveRelayer)
	if err := _RelayerManager.contract.UnpackLog(event, "evtRemoveRelayer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
