// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package maas_config_abi

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
	MethodBlockAccount = "blockAccount"

	MethodChangeOwner = "changeOwner"

	MethodGetBlacklist = "getBlacklist"

	MethodGetOwner = "getOwner"

	MethodIsBlocked = "isBlocked"

	MethodName = "name"

	EventBlockAccount = "BlockAccount"

	EventChangeOwner = "ChangeOwner"
)

// MaasConfigABI is the input ABI used to generate the binding from.
const MaasConfigABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"doBlock\",\"type\":\"bool\"}],\"name\":\"BlockAccount\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"ChangeOwner\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"doBlock\",\"type\":\"bool\"}],\"name\":\"blockAccount\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlacklist\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isBlocked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// MaasConfig is an auto generated Go binding around an Ethereum contract.
type MaasConfig struct {
	MaasConfigCaller     // Read-only binding to the contract
	MaasConfigTransactor // Write-only binding to the contract
	MaasConfigFilterer   // Log filterer for contract events
}

// MaasConfigCaller is an auto generated read-only Go binding around an Ethereum contract.
type MaasConfigCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MaasConfigTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MaasConfigTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MaasConfigFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MaasConfigFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MaasConfigSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MaasConfigSession struct {
	Contract     *MaasConfig       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MaasConfigCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MaasConfigCallerSession struct {
	Contract *MaasConfigCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// MaasConfigTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MaasConfigTransactorSession struct {
	Contract     *MaasConfigTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// MaasConfigRaw is an auto generated low-level Go binding around an Ethereum contract.
type MaasConfigRaw struct {
	Contract *MaasConfig // Generic contract binding to access the raw methods on
}

// MaasConfigCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MaasConfigCallerRaw struct {
	Contract *MaasConfigCaller // Generic read-only contract binding to access the raw methods on
}

// MaasConfigTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MaasConfigTransactorRaw struct {
	Contract *MaasConfigTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMaasConfig creates a new instance of MaasConfig, bound to a specific deployed contract.
func NewMaasConfig(address common.Address, backend bind.ContractBackend) (*MaasConfig, error) {
	contract, err := bindMaasConfig(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MaasConfig{MaasConfigCaller: MaasConfigCaller{contract: contract}, MaasConfigTransactor: MaasConfigTransactor{contract: contract}, MaasConfigFilterer: MaasConfigFilterer{contract: contract}}, nil
}

// NewMaasConfigCaller creates a new read-only instance of MaasConfig, bound to a specific deployed contract.
func NewMaasConfigCaller(address common.Address, caller bind.ContractCaller) (*MaasConfigCaller, error) {
	contract, err := bindMaasConfig(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MaasConfigCaller{contract: contract}, nil
}

// NewMaasConfigTransactor creates a new write-only instance of MaasConfig, bound to a specific deployed contract.
func NewMaasConfigTransactor(address common.Address, transactor bind.ContractTransactor) (*MaasConfigTransactor, error) {
	contract, err := bindMaasConfig(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MaasConfigTransactor{contract: contract}, nil
}

// NewMaasConfigFilterer creates a new log filterer instance of MaasConfig, bound to a specific deployed contract.
func NewMaasConfigFilterer(address common.Address, filterer bind.ContractFilterer) (*MaasConfigFilterer, error) {
	contract, err := bindMaasConfig(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MaasConfigFilterer{contract: contract}, nil
}

// bindMaasConfig binds a generic wrapper to an already deployed contract.
func bindMaasConfig(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MaasConfigABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MaasConfig *MaasConfigRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MaasConfig.Contract.MaasConfigCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MaasConfig *MaasConfigRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MaasConfig.Contract.MaasConfigTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MaasConfig *MaasConfigRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MaasConfig.Contract.MaasConfigTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MaasConfig *MaasConfigCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MaasConfig.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MaasConfig *MaasConfigTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MaasConfig.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MaasConfig *MaasConfigTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MaasConfig.Contract.contract.Transact(opts, method, params...)
}

// GetBlacklist is a free data retrieval call binding the contract method 0x338d6c30.
//
// Solidity: function getBlacklist() view returns(string)
func (_MaasConfig *MaasConfigCaller) GetBlacklist(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MaasConfig.contract.Call(opts, &out, "getBlacklist")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetBlacklist is a free data retrieval call binding the contract method 0x338d6c30.
//
// Solidity: function getBlacklist() view returns(string)
func (_MaasConfig *MaasConfigSession) GetBlacklist() (string, error) {
	return _MaasConfig.Contract.GetBlacklist(&_MaasConfig.CallOpts)
}

// GetBlacklist is a free data retrieval call binding the contract method 0x338d6c30.
//
// Solidity: function getBlacklist() view returns(string)
func (_MaasConfig *MaasConfigCallerSession) GetBlacklist() (string, error) {
	return _MaasConfig.Contract.GetBlacklist(&_MaasConfig.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() view returns(address)
func (_MaasConfig *MaasConfigCaller) GetOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MaasConfig.contract.Call(opts, &out, "getOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() view returns(address)
func (_MaasConfig *MaasConfigSession) GetOwner() (common.Address, error) {
	return _MaasConfig.Contract.GetOwner(&_MaasConfig.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() view returns(address)
func (_MaasConfig *MaasConfigCallerSession) GetOwner() (common.Address, error) {
	return _MaasConfig.Contract.GetOwner(&_MaasConfig.CallOpts)
}

// IsBlocked is a free data retrieval call binding the contract method 0xfbac3951.
//
// Solidity: function isBlocked(address addr) view returns(bool)
func (_MaasConfig *MaasConfigCaller) IsBlocked(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _MaasConfig.contract.Call(opts, &out, "isBlocked", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBlocked is a free data retrieval call binding the contract method 0xfbac3951.
//
// Solidity: function isBlocked(address addr) view returns(bool)
func (_MaasConfig *MaasConfigSession) IsBlocked(addr common.Address) (bool, error) {
	return _MaasConfig.Contract.IsBlocked(&_MaasConfig.CallOpts, addr)
}

// IsBlocked is a free data retrieval call binding the contract method 0xfbac3951.
//
// Solidity: function isBlocked(address addr) view returns(bool)
func (_MaasConfig *MaasConfigCallerSession) IsBlocked(addr common.Address) (bool, error) {
	return _MaasConfig.Contract.IsBlocked(&_MaasConfig.CallOpts, addr)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MaasConfig *MaasConfigCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MaasConfig.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MaasConfig *MaasConfigSession) Name() (string, error) {
	return _MaasConfig.Contract.Name(&_MaasConfig.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MaasConfig *MaasConfigCallerSession) Name() (string, error) {
	return _MaasConfig.Contract.Name(&_MaasConfig.CallOpts)
}

// BlockAccount is a paid mutator transaction binding the contract method 0x52c163bb.
//
// Solidity: function blockAccount(address addr, bool doBlock) returns(bool)
func (_MaasConfig *MaasConfigTransactor) BlockAccount(opts *bind.TransactOpts, addr common.Address, doBlock bool) (*types.Transaction, error) {
	return _MaasConfig.contract.Transact(opts, "blockAccount", addr, doBlock)
}

// BlockAccount is a paid mutator transaction binding the contract method 0x52c163bb.
//
// Solidity: function blockAccount(address addr, bool doBlock) returns(bool)
func (_MaasConfig *MaasConfigSession) BlockAccount(addr common.Address, doBlock bool) (*types.Transaction, error) {
	return _MaasConfig.Contract.BlockAccount(&_MaasConfig.TransactOpts, addr, doBlock)
}

// BlockAccount is a paid mutator transaction binding the contract method 0x52c163bb.
//
// Solidity: function blockAccount(address addr, bool doBlock) returns(bool)
func (_MaasConfig *MaasConfigTransactorSession) BlockAccount(addr common.Address, doBlock bool) (*types.Transaction, error) {
	return _MaasConfig.Contract.BlockAccount(&_MaasConfig.TransactOpts, addr, doBlock)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address addr) returns(bool)
func (_MaasConfig *MaasConfigTransactor) ChangeOwner(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _MaasConfig.contract.Transact(opts, "changeOwner", addr)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address addr) returns(bool)
func (_MaasConfig *MaasConfigSession) ChangeOwner(addr common.Address) (*types.Transaction, error) {
	return _MaasConfig.Contract.ChangeOwner(&_MaasConfig.TransactOpts, addr)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address addr) returns(bool)
func (_MaasConfig *MaasConfigTransactorSession) ChangeOwner(addr common.Address) (*types.Transaction, error) {
	return _MaasConfig.Contract.ChangeOwner(&_MaasConfig.TransactOpts, addr)
}

// MaasConfigBlockAccountIterator is returned from FilterBlockAccount and is used to iterate over the raw logs and unpacked data for BlockAccount events raised by the MaasConfig contract.
type MaasConfigBlockAccountIterator struct {
	Event *MaasConfigBlockAccount // Event containing the contract specifics and raw log

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
func (it *MaasConfigBlockAccountIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaasConfigBlockAccount)
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
		it.Event = new(MaasConfigBlockAccount)
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
func (it *MaasConfigBlockAccountIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MaasConfigBlockAccountIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MaasConfigBlockAccount represents a BlockAccount event raised by the MaasConfig contract.
type MaasConfigBlockAccount struct {
	Addr    common.Address
	DoBlock bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBlockAccount is a free log retrieval operation binding the contract event 0x977826a31e63a99f714f2677060d8f5d42a578272b31da3a8088f758ca915fdf.
//
// Solidity: event BlockAccount(address addr, bool doBlock)
func (_MaasConfig *MaasConfigFilterer) FilterBlockAccount(opts *bind.FilterOpts) (*MaasConfigBlockAccountIterator, error) {

	logs, sub, err := _MaasConfig.contract.FilterLogs(opts, "BlockAccount")
	if err != nil {
		return nil, err
	}
	return &MaasConfigBlockAccountIterator{contract: _MaasConfig.contract, event: "BlockAccount", logs: logs, sub: sub}, nil
}

// WatchBlockAccount is a free log subscription operation binding the contract event 0x977826a31e63a99f714f2677060d8f5d42a578272b31da3a8088f758ca915fdf.
//
// Solidity: event BlockAccount(address addr, bool doBlock)
func (_MaasConfig *MaasConfigFilterer) WatchBlockAccount(opts *bind.WatchOpts, sink chan<- *MaasConfigBlockAccount) (event.Subscription, error) {

	logs, sub, err := _MaasConfig.contract.WatchLogs(opts, "BlockAccount")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MaasConfigBlockAccount)
				if err := _MaasConfig.contract.UnpackLog(event, "BlockAccount", log); err != nil {
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

// ParseBlockAccount is a log parse operation binding the contract event 0x977826a31e63a99f714f2677060d8f5d42a578272b31da3a8088f758ca915fdf.
//
// Solidity: event BlockAccount(address addr, bool doBlock)
func (_MaasConfig *MaasConfigFilterer) ParseBlockAccount(log types.Log) (*MaasConfigBlockAccount, error) {
	event := new(MaasConfigBlockAccount)
	if err := _MaasConfig.contract.UnpackLog(event, "BlockAccount", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MaasConfigChangeOwnerIterator is returned from FilterChangeOwner and is used to iterate over the raw logs and unpacked data for ChangeOwner events raised by the MaasConfig contract.
type MaasConfigChangeOwnerIterator struct {
	Event *MaasConfigChangeOwner // Event containing the contract specifics and raw log

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
func (it *MaasConfigChangeOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaasConfigChangeOwner)
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
		it.Event = new(MaasConfigChangeOwner)
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
func (it *MaasConfigChangeOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MaasConfigChangeOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MaasConfigChangeOwner represents a ChangeOwner event raised by the MaasConfig contract.
type MaasConfigChangeOwner struct {
	OldOwner common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterChangeOwner is a free log retrieval operation binding the contract event 0x9aecf86140d81442289f667eb72e1202a8fbb3478a686659952e145e85319656.
//
// Solidity: event ChangeOwner(address oldOwner, address newOwner)
func (_MaasConfig *MaasConfigFilterer) FilterChangeOwner(opts *bind.FilterOpts) (*MaasConfigChangeOwnerIterator, error) {

	logs, sub, err := _MaasConfig.contract.FilterLogs(opts, "ChangeOwner")
	if err != nil {
		return nil, err
	}
	return &MaasConfigChangeOwnerIterator{contract: _MaasConfig.contract, event: "ChangeOwner", logs: logs, sub: sub}, nil
}

// WatchChangeOwner is a free log subscription operation binding the contract event 0x9aecf86140d81442289f667eb72e1202a8fbb3478a686659952e145e85319656.
//
// Solidity: event ChangeOwner(address oldOwner, address newOwner)
func (_MaasConfig *MaasConfigFilterer) WatchChangeOwner(opts *bind.WatchOpts, sink chan<- *MaasConfigChangeOwner) (event.Subscription, error) {

	logs, sub, err := _MaasConfig.contract.WatchLogs(opts, "ChangeOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MaasConfigChangeOwner)
				if err := _MaasConfig.contract.UnpackLog(event, "ChangeOwner", log); err != nil {
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

// ParseChangeOwner is a log parse operation binding the contract event 0x9aecf86140d81442289f667eb72e1202a8fbb3478a686659952e145e85319656.
//
// Solidity: event ChangeOwner(address oldOwner, address newOwner)
func (_MaasConfig *MaasConfigFilterer) ParseChangeOwner(log types.Log) (*MaasConfigChangeOwner, error) {
	event := new(MaasConfigChangeOwner)
	if err := _MaasConfig.contract.UnpackLog(event, "ChangeOwner", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
