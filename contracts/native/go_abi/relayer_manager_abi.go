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

// RelayerManagerABI is the input ABI used to generate the binding from.
const RelayerManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"}],\"name\":\"EventApproveRegisterRelayer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"}],\"name\":\"EventApproveRemoveRelayer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"applyID\",\"type\":\"uint64\"}],\"name\":\"EventRegisterRelayer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"removeID\",\"type\":\"uint64\"}],\"name\":\"EventRemoveRelayer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodApproveRegisterRelayer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"ID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodApproveRemoveRelayer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MethodContractName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"AddressList\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodRegisterRelayer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"AddressList\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodRemoveRelayer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// RelayerManagerFuncSigs maps the 4-byte function signature to its string representation.
var RelayerManagerFuncSigs = map[string]string{
	"2b999ab8": "MethodApproveRegisterRelayer(uint64,address)",
	"c8e4f8af": "MethodApproveRemoveRelayer(uint64,address)",
	"e50f8f44": "MethodContractName()",
	"8df11025": "MethodRegisterRelayer(address[],address)",
	"1791d10b": "MethodRemoveRelayer(address[],address)",
}

// RelayerManagerBin is the compiled bytecode used for deploying new contracts.
var RelayerManagerBin = "0x608060405234801561001057600080fd5b50610281806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80631791d10b1461005c5780632b999ab8146100875780638df110251461005c578063c8e4f8af14610087578063e50f8f4414610095575b600080fd5b61007261006a3660046100c0565b600092915050565b60405190151581526020015b60405180910390f35b61007261006a36600461019e565b606060405161007e91906101e0565b80356001600160a01b03811681146100bb57600080fd5b919050565b600080604083850312156100d357600080fd5b823567ffffffffffffffff808211156100eb57600080fd5b818501915085601f8301126100ff57600080fd5b813560208282111561011357610113610235565b8160051b604051601f19603f8301168101818110868211171561013857610138610235565b604052838152828101945085830182870184018b101561015757600080fd5b600096505b848710156101815761016d816100a4565b86526001969096019594830194830161015c565b50965061019190508782016100a4565b9450505050509250929050565b600080604083850312156101b157600080fd5b823567ffffffffffffffff811681146101c957600080fd5b91506101d7602084016100a4565b90509250929050565b600060208083528351808285015260005b8181101561020d578581018301518582016040015282016101f1565b8181111561021f576000604083870101525b50601f01601f1916929092016040019392505050565b634e487b7160e01b600052604160045260246000fdfea264697066735822122076502f86b779cc75290e2a2a1b0adc83e53813015f428eb07cc5347a0084669464736f6c63430008060033"

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

// MethodApproveRegisterRelayer is a paid mutator transaction binding the contract method 0x2b999ab8.
//
// Solidity: function MethodApproveRegisterRelayer(uint64 ID, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactor) MethodApproveRegisterRelayer(opts *bind.TransactOpts, ID uint64, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.contract.Transact(opts, "MethodApproveRegisterRelayer", ID, Address)
}

// MethodApproveRegisterRelayer is a paid mutator transaction binding the contract method 0x2b999ab8.
//
// Solidity: function MethodApproveRegisterRelayer(uint64 ID, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerSession) MethodApproveRegisterRelayer(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.MethodApproveRegisterRelayer(&_RelayerManager.TransactOpts, ID, Address)
}

// MethodApproveRegisterRelayer is a paid mutator transaction binding the contract method 0x2b999ab8.
//
// Solidity: function MethodApproveRegisterRelayer(uint64 ID, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactorSession) MethodApproveRegisterRelayer(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.MethodApproveRegisterRelayer(&_RelayerManager.TransactOpts, ID, Address)
}

// MethodApproveRemoveRelayer is a paid mutator transaction binding the contract method 0xc8e4f8af.
//
// Solidity: function MethodApproveRemoveRelayer(uint64 ID, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactor) MethodApproveRemoveRelayer(opts *bind.TransactOpts, ID uint64, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.contract.Transact(opts, "MethodApproveRemoveRelayer", ID, Address)
}

// MethodApproveRemoveRelayer is a paid mutator transaction binding the contract method 0xc8e4f8af.
//
// Solidity: function MethodApproveRemoveRelayer(uint64 ID, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerSession) MethodApproveRemoveRelayer(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.MethodApproveRemoveRelayer(&_RelayerManager.TransactOpts, ID, Address)
}

// MethodApproveRemoveRelayer is a paid mutator transaction binding the contract method 0xc8e4f8af.
//
// Solidity: function MethodApproveRemoveRelayer(uint64 ID, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactorSession) MethodApproveRemoveRelayer(ID uint64, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.MethodApproveRemoveRelayer(&_RelayerManager.TransactOpts, ID, Address)
}

// MethodContractName is a paid mutator transaction binding the contract method 0xe50f8f44.
//
// Solidity: function MethodContractName() returns(string Name)
func (_RelayerManager *RelayerManagerTransactor) MethodContractName(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RelayerManager.contract.Transact(opts, "MethodContractName")
}

// MethodContractName is a paid mutator transaction binding the contract method 0xe50f8f44.
//
// Solidity: function MethodContractName() returns(string Name)
func (_RelayerManager *RelayerManagerSession) MethodContractName() (*types.Transaction, error) {
	return _RelayerManager.Contract.MethodContractName(&_RelayerManager.TransactOpts)
}

// MethodContractName is a paid mutator transaction binding the contract method 0xe50f8f44.
//
// Solidity: function MethodContractName() returns(string Name)
func (_RelayerManager *RelayerManagerTransactorSession) MethodContractName() (*types.Transaction, error) {
	return _RelayerManager.Contract.MethodContractName(&_RelayerManager.TransactOpts)
}

// MethodRegisterRelayer is a paid mutator transaction binding the contract method 0x8df11025.
//
// Solidity: function MethodRegisterRelayer(address[] AddressList, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactor) MethodRegisterRelayer(opts *bind.TransactOpts, AddressList []common.Address, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.contract.Transact(opts, "MethodRegisterRelayer", AddressList, Address)
}

// MethodRegisterRelayer is a paid mutator transaction binding the contract method 0x8df11025.
//
// Solidity: function MethodRegisterRelayer(address[] AddressList, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerSession) MethodRegisterRelayer(AddressList []common.Address, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.MethodRegisterRelayer(&_RelayerManager.TransactOpts, AddressList, Address)
}

// MethodRegisterRelayer is a paid mutator transaction binding the contract method 0x8df11025.
//
// Solidity: function MethodRegisterRelayer(address[] AddressList, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactorSession) MethodRegisterRelayer(AddressList []common.Address, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.MethodRegisterRelayer(&_RelayerManager.TransactOpts, AddressList, Address)
}

// MethodRemoveRelayer is a paid mutator transaction binding the contract method 0x1791d10b.
//
// Solidity: function MethodRemoveRelayer(address[] AddressList, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactor) MethodRemoveRelayer(opts *bind.TransactOpts, AddressList []common.Address, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.contract.Transact(opts, "MethodRemoveRelayer", AddressList, Address)
}

// MethodRemoveRelayer is a paid mutator transaction binding the contract method 0x1791d10b.
//
// Solidity: function MethodRemoveRelayer(address[] AddressList, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerSession) MethodRemoveRelayer(AddressList []common.Address, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.MethodRemoveRelayer(&_RelayerManager.TransactOpts, AddressList, Address)
}

// MethodRemoveRelayer is a paid mutator transaction binding the contract method 0x1791d10b.
//
// Solidity: function MethodRemoveRelayer(address[] AddressList, address Address) returns(bool success)
func (_RelayerManager *RelayerManagerTransactorSession) MethodRemoveRelayer(AddressList []common.Address, Address common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.MethodRemoveRelayer(&_RelayerManager.TransactOpts, AddressList, Address)
}

// RelayerManagerEventApproveRegisterRelayerIterator is returned from FilterEventApproveRegisterRelayer and is used to iterate over the raw logs and unpacked data for EventApproveRegisterRelayer events raised by the RelayerManager contract.
type RelayerManagerEventApproveRegisterRelayerIterator struct {
	Event *RelayerManagerEventApproveRegisterRelayer // Event containing the contract specifics and raw log

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
func (it *RelayerManagerEventApproveRegisterRelayerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerManagerEventApproveRegisterRelayer)
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
		it.Event = new(RelayerManagerEventApproveRegisterRelayer)
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
func (it *RelayerManagerEventApproveRegisterRelayerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerManagerEventApproveRegisterRelayerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerManagerEventApproveRegisterRelayer represents a EventApproveRegisterRelayer event raised by the RelayerManager contract.
type RelayerManagerEventApproveRegisterRelayer struct {
	ID  uint64
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEventApproveRegisterRelayer is a free log retrieval operation binding the contract event 0x1ed91da95b8b4cb6c350992aedc9d61a9062bff0544bf96e58457d34f1212022.
//
// Solidity: event EventApproveRegisterRelayer(uint64 ID)
func (_RelayerManager *RelayerManagerFilterer) FilterEventApproveRegisterRelayer(opts *bind.FilterOpts) (*RelayerManagerEventApproveRegisterRelayerIterator, error) {

	logs, sub, err := _RelayerManager.contract.FilterLogs(opts, "EventApproveRegisterRelayer")
	if err != nil {
		return nil, err
	}
	return &RelayerManagerEventApproveRegisterRelayerIterator{contract: _RelayerManager.contract, event: "EventApproveRegisterRelayer", logs: logs, sub: sub}, nil
}

// WatchEventApproveRegisterRelayer is a free log subscription operation binding the contract event 0x1ed91da95b8b4cb6c350992aedc9d61a9062bff0544bf96e58457d34f1212022.
//
// Solidity: event EventApproveRegisterRelayer(uint64 ID)
func (_RelayerManager *RelayerManagerFilterer) WatchEventApproveRegisterRelayer(opts *bind.WatchOpts, sink chan<- *RelayerManagerEventApproveRegisterRelayer) (event.Subscription, error) {

	logs, sub, err := _RelayerManager.contract.WatchLogs(opts, "EventApproveRegisterRelayer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerManagerEventApproveRegisterRelayer)
				if err := _RelayerManager.contract.UnpackLog(event, "EventApproveRegisterRelayer", log); err != nil {
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

// ParseEventApproveRegisterRelayer is a log parse operation binding the contract event 0x1ed91da95b8b4cb6c350992aedc9d61a9062bff0544bf96e58457d34f1212022.
//
// Solidity: event EventApproveRegisterRelayer(uint64 ID)
func (_RelayerManager *RelayerManagerFilterer) ParseEventApproveRegisterRelayer(log types.Log) (*RelayerManagerEventApproveRegisterRelayer, error) {
	event := new(RelayerManagerEventApproveRegisterRelayer)
	if err := _RelayerManager.contract.UnpackLog(event, "EventApproveRegisterRelayer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayerManagerEventApproveRemoveRelayerIterator is returned from FilterEventApproveRemoveRelayer and is used to iterate over the raw logs and unpacked data for EventApproveRemoveRelayer events raised by the RelayerManager contract.
type RelayerManagerEventApproveRemoveRelayerIterator struct {
	Event *RelayerManagerEventApproveRemoveRelayer // Event containing the contract specifics and raw log

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
func (it *RelayerManagerEventApproveRemoveRelayerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerManagerEventApproveRemoveRelayer)
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
		it.Event = new(RelayerManagerEventApproveRemoveRelayer)
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
func (it *RelayerManagerEventApproveRemoveRelayerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerManagerEventApproveRemoveRelayerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerManagerEventApproveRemoveRelayer represents a EventApproveRemoveRelayer event raised by the RelayerManager contract.
type RelayerManagerEventApproveRemoveRelayer struct {
	ID  uint64
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEventApproveRemoveRelayer is a free log retrieval operation binding the contract event 0x1d5e8420b7108839cc5f3e9e1e3bfbf2ddb512b2a5e0506dedbba7200b012e28.
//
// Solidity: event EventApproveRemoveRelayer(uint64 ID)
func (_RelayerManager *RelayerManagerFilterer) FilterEventApproveRemoveRelayer(opts *bind.FilterOpts) (*RelayerManagerEventApproveRemoveRelayerIterator, error) {

	logs, sub, err := _RelayerManager.contract.FilterLogs(opts, "EventApproveRemoveRelayer")
	if err != nil {
		return nil, err
	}
	return &RelayerManagerEventApproveRemoveRelayerIterator{contract: _RelayerManager.contract, event: "EventApproveRemoveRelayer", logs: logs, sub: sub}, nil
}

// WatchEventApproveRemoveRelayer is a free log subscription operation binding the contract event 0x1d5e8420b7108839cc5f3e9e1e3bfbf2ddb512b2a5e0506dedbba7200b012e28.
//
// Solidity: event EventApproveRemoveRelayer(uint64 ID)
func (_RelayerManager *RelayerManagerFilterer) WatchEventApproveRemoveRelayer(opts *bind.WatchOpts, sink chan<- *RelayerManagerEventApproveRemoveRelayer) (event.Subscription, error) {

	logs, sub, err := _RelayerManager.contract.WatchLogs(opts, "EventApproveRemoveRelayer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerManagerEventApproveRemoveRelayer)
				if err := _RelayerManager.contract.UnpackLog(event, "EventApproveRemoveRelayer", log); err != nil {
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

// ParseEventApproveRemoveRelayer is a log parse operation binding the contract event 0x1d5e8420b7108839cc5f3e9e1e3bfbf2ddb512b2a5e0506dedbba7200b012e28.
//
// Solidity: event EventApproveRemoveRelayer(uint64 ID)
func (_RelayerManager *RelayerManagerFilterer) ParseEventApproveRemoveRelayer(log types.Log) (*RelayerManagerEventApproveRemoveRelayer, error) {
	event := new(RelayerManagerEventApproveRemoveRelayer)
	if err := _RelayerManager.contract.UnpackLog(event, "EventApproveRemoveRelayer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayerManagerEventRegisterRelayerIterator is returned from FilterEventRegisterRelayer and is used to iterate over the raw logs and unpacked data for EventRegisterRelayer events raised by the RelayerManager contract.
type RelayerManagerEventRegisterRelayerIterator struct {
	Event *RelayerManagerEventRegisterRelayer // Event containing the contract specifics and raw log

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
func (it *RelayerManagerEventRegisterRelayerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerManagerEventRegisterRelayer)
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
		it.Event = new(RelayerManagerEventRegisterRelayer)
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
func (it *RelayerManagerEventRegisterRelayerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerManagerEventRegisterRelayerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerManagerEventRegisterRelayer represents a EventRegisterRelayer event raised by the RelayerManager contract.
type RelayerManagerEventRegisterRelayer struct {
	ApplyID uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterEventRegisterRelayer is a free log retrieval operation binding the contract event 0x413d9d4e393cff2595025983ee75b5ec5144d4301a803414d3760e490acaa4ee.
//
// Solidity: event EventRegisterRelayer(uint64 applyID)
func (_RelayerManager *RelayerManagerFilterer) FilterEventRegisterRelayer(opts *bind.FilterOpts) (*RelayerManagerEventRegisterRelayerIterator, error) {

	logs, sub, err := _RelayerManager.contract.FilterLogs(opts, "EventRegisterRelayer")
	if err != nil {
		return nil, err
	}
	return &RelayerManagerEventRegisterRelayerIterator{contract: _RelayerManager.contract, event: "EventRegisterRelayer", logs: logs, sub: sub}, nil
}

// WatchEventRegisterRelayer is a free log subscription operation binding the contract event 0x413d9d4e393cff2595025983ee75b5ec5144d4301a803414d3760e490acaa4ee.
//
// Solidity: event EventRegisterRelayer(uint64 applyID)
func (_RelayerManager *RelayerManagerFilterer) WatchEventRegisterRelayer(opts *bind.WatchOpts, sink chan<- *RelayerManagerEventRegisterRelayer) (event.Subscription, error) {

	logs, sub, err := _RelayerManager.contract.WatchLogs(opts, "EventRegisterRelayer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerManagerEventRegisterRelayer)
				if err := _RelayerManager.contract.UnpackLog(event, "EventRegisterRelayer", log); err != nil {
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

// ParseEventRegisterRelayer is a log parse operation binding the contract event 0x413d9d4e393cff2595025983ee75b5ec5144d4301a803414d3760e490acaa4ee.
//
// Solidity: event EventRegisterRelayer(uint64 applyID)
func (_RelayerManager *RelayerManagerFilterer) ParseEventRegisterRelayer(log types.Log) (*RelayerManagerEventRegisterRelayer, error) {
	event := new(RelayerManagerEventRegisterRelayer)
	if err := _RelayerManager.contract.UnpackLog(event, "EventRegisterRelayer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayerManagerEventRemoveRelayerIterator is returned from FilterEventRemoveRelayer and is used to iterate over the raw logs and unpacked data for EventRemoveRelayer events raised by the RelayerManager contract.
type RelayerManagerEventRemoveRelayerIterator struct {
	Event *RelayerManagerEventRemoveRelayer // Event containing the contract specifics and raw log

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
func (it *RelayerManagerEventRemoveRelayerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerManagerEventRemoveRelayer)
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
		it.Event = new(RelayerManagerEventRemoveRelayer)
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
func (it *RelayerManagerEventRemoveRelayerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerManagerEventRemoveRelayerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerManagerEventRemoveRelayer represents a EventRemoveRelayer event raised by the RelayerManager contract.
type RelayerManagerEventRemoveRelayer struct {
	RemoveID uint64
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterEventRemoveRelayer is a free log retrieval operation binding the contract event 0x5346c918bfcf400258e68a4fe4ddfdf407402418b51c535edc529deaade004d8.
//
// Solidity: event EventRemoveRelayer(uint64 removeID)
func (_RelayerManager *RelayerManagerFilterer) FilterEventRemoveRelayer(opts *bind.FilterOpts) (*RelayerManagerEventRemoveRelayerIterator, error) {

	logs, sub, err := _RelayerManager.contract.FilterLogs(opts, "EventRemoveRelayer")
	if err != nil {
		return nil, err
	}
	return &RelayerManagerEventRemoveRelayerIterator{contract: _RelayerManager.contract, event: "EventRemoveRelayer", logs: logs, sub: sub}, nil
}

// WatchEventRemoveRelayer is a free log subscription operation binding the contract event 0x5346c918bfcf400258e68a4fe4ddfdf407402418b51c535edc529deaade004d8.
//
// Solidity: event EventRemoveRelayer(uint64 removeID)
func (_RelayerManager *RelayerManagerFilterer) WatchEventRemoveRelayer(opts *bind.WatchOpts, sink chan<- *RelayerManagerEventRemoveRelayer) (event.Subscription, error) {

	logs, sub, err := _RelayerManager.contract.WatchLogs(opts, "EventRemoveRelayer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerManagerEventRemoveRelayer)
				if err := _RelayerManager.contract.UnpackLog(event, "EventRemoveRelayer", log); err != nil {
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

// ParseEventRemoveRelayer is a log parse operation binding the contract event 0x5346c918bfcf400258e68a4fe4ddfdf407402418b51c535edc529deaade004d8.
//
// Solidity: event EventRemoveRelayer(uint64 removeID)
func (_RelayerManager *RelayerManagerFilterer) ParseEventRemoveRelayer(log types.Log) (*RelayerManagerEventRemoveRelayer, error) {
	event := new(RelayerManagerEventRemoveRelayer)
	if err := _RelayerManager.contract.UnpackLog(event, "EventRemoveRelayer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
