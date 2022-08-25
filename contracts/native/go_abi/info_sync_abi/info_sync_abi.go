// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package info_sync_abi

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
	MethodReplenish = "replenish"

	MethodSyncRootInfo = "syncRootInfo"

	MethodGetInfo = "getInfo"

	MethodGetInfoHeight = "getInfoHeight"

	MethodName = "name"

	EventReplenishEvent = "ReplenishEvent"

	EventSyncRootInfoEvent = "SyncRootInfoEvent"
)

// IInfoSyncABI is the input ABI used to generate the binding from.
const IInfoSyncABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32[]\",\"name\":\"heights\",\"type\":\"uint32[]\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"name\":\"ReplenishEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"height\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"BlockHeight\",\"type\":\"uint256\"}],\"name\":\"SyncRootInfoEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"height\",\"type\":\"uint32\"}],\"name\":\"getInfo\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"name\":\"getInfoHeight\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"uint32[]\",\"name\":\"heights\",\"type\":\"uint32[]\"}],\"name\":\"replenish\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"bytes[]\",\"name\":\"rootInfos\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"syncRootInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IInfoSyncFuncSigs maps the 4-byte function signature to its string representation.
var IInfoSyncFuncSigs = map[string]string{
	"6a4a9f5e": "getInfo(uint64,uint32)",
	"16d80012": "getInfoHeight(uint64)",
	"06fdde03": "name()",
	"69ce93b4": "replenish(uint64,uint32[])",
	"1413cc01": "syncRootInfo(uint64,bytes[],bytes)",
}

// IInfoSync is an auto generated Go binding around an Ethereum contract.
type IInfoSync struct {
	IInfoSyncCaller     // Read-only binding to the contract
	IInfoSyncTransactor // Write-only binding to the contract
	IInfoSyncFilterer   // Log filterer for contract events
}

// IInfoSyncCaller is an auto generated read-only Go binding around an Ethereum contract.
type IInfoSyncCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IInfoSyncTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IInfoSyncTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IInfoSyncFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IInfoSyncFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IInfoSyncSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IInfoSyncSession struct {
	Contract     *IInfoSync        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IInfoSyncCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IInfoSyncCallerSession struct {
	Contract *IInfoSyncCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// IInfoSyncTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IInfoSyncTransactorSession struct {
	Contract     *IInfoSyncTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IInfoSyncRaw is an auto generated low-level Go binding around an Ethereum contract.
type IInfoSyncRaw struct {
	Contract *IInfoSync // Generic contract binding to access the raw methods on
}

// IInfoSyncCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IInfoSyncCallerRaw struct {
	Contract *IInfoSyncCaller // Generic read-only contract binding to access the raw methods on
}

// IInfoSyncTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IInfoSyncTransactorRaw struct {
	Contract *IInfoSyncTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIInfoSync creates a new instance of IInfoSync, bound to a specific deployed contract.
func NewIInfoSync(address common.Address, backend bind.ContractBackend) (*IInfoSync, error) {
	contract, err := bindIInfoSync(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IInfoSync{IInfoSyncCaller: IInfoSyncCaller{contract: contract}, IInfoSyncTransactor: IInfoSyncTransactor{contract: contract}, IInfoSyncFilterer: IInfoSyncFilterer{contract: contract}}, nil
}

// NewIInfoSyncCaller creates a new read-only instance of IInfoSync, bound to a specific deployed contract.
func NewIInfoSyncCaller(address common.Address, caller bind.ContractCaller) (*IInfoSyncCaller, error) {
	contract, err := bindIInfoSync(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IInfoSyncCaller{contract: contract}, nil
}

// NewIInfoSyncTransactor creates a new write-only instance of IInfoSync, bound to a specific deployed contract.
func NewIInfoSyncTransactor(address common.Address, transactor bind.ContractTransactor) (*IInfoSyncTransactor, error) {
	contract, err := bindIInfoSync(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IInfoSyncTransactor{contract: contract}, nil
}

// NewIInfoSyncFilterer creates a new log filterer instance of IInfoSync, bound to a specific deployed contract.
func NewIInfoSyncFilterer(address common.Address, filterer bind.ContractFilterer) (*IInfoSyncFilterer, error) {
	contract, err := bindIInfoSync(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IInfoSyncFilterer{contract: contract}, nil
}

// bindIInfoSync binds a generic wrapper to an already deployed contract.
func bindIInfoSync(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IInfoSyncABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IInfoSync *IInfoSyncRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IInfoSync.Contract.IInfoSyncCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IInfoSync *IInfoSyncRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IInfoSync.Contract.IInfoSyncTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IInfoSync *IInfoSyncRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IInfoSync.Contract.IInfoSyncTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IInfoSync *IInfoSyncCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IInfoSync.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IInfoSync *IInfoSyncTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IInfoSync.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IInfoSync *IInfoSyncTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IInfoSync.Contract.contract.Transact(opts, method, params...)
}

// GetInfo is a free data retrieval call binding the contract method 0x6a4a9f5e.
//
// Solidity: function getInfo(uint64 chainID, uint32 height) view returns(bytes)
func (_IInfoSync *IInfoSyncCaller) GetInfo(opts *bind.CallOpts, chainID uint64, height uint32) ([]byte, error) {
	var out []interface{}
	err := _IInfoSync.contract.Call(opts, &out, "getInfo", chainID, height)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetInfo is a free data retrieval call binding the contract method 0x6a4a9f5e.
//
// Solidity: function getInfo(uint64 chainID, uint32 height) view returns(bytes)
func (_IInfoSync *IInfoSyncSession) GetInfo(chainID uint64, height uint32) ([]byte, error) {
	return _IInfoSync.Contract.GetInfo(&_IInfoSync.CallOpts, chainID, height)
}

// GetInfo is a free data retrieval call binding the contract method 0x6a4a9f5e.
//
// Solidity: function getInfo(uint64 chainID, uint32 height) view returns(bytes)
func (_IInfoSync *IInfoSyncCallerSession) GetInfo(chainID uint64, height uint32) ([]byte, error) {
	return _IInfoSync.Contract.GetInfo(&_IInfoSync.CallOpts, chainID, height)
}

// GetInfoHeight is a free data retrieval call binding the contract method 0x16d80012.
//
// Solidity: function getInfoHeight(uint64 chainID) view returns(uint32)
func (_IInfoSync *IInfoSyncCaller) GetInfoHeight(opts *bind.CallOpts, chainID uint64) (uint32, error) {
	var out []interface{}
	err := _IInfoSync.contract.Call(opts, &out, "getInfoHeight", chainID)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetInfoHeight is a free data retrieval call binding the contract method 0x16d80012.
//
// Solidity: function getInfoHeight(uint64 chainID) view returns(uint32)
func (_IInfoSync *IInfoSyncSession) GetInfoHeight(chainID uint64) (uint32, error) {
	return _IInfoSync.Contract.GetInfoHeight(&_IInfoSync.CallOpts, chainID)
}

// GetInfoHeight is a free data retrieval call binding the contract method 0x16d80012.
//
// Solidity: function getInfoHeight(uint64 chainID) view returns(uint32)
func (_IInfoSync *IInfoSyncCallerSession) GetInfoHeight(chainID uint64) (uint32, error) {
	return _IInfoSync.Contract.GetInfoHeight(&_IInfoSync.CallOpts, chainID)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_IInfoSync *IInfoSyncCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IInfoSync.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_IInfoSync *IInfoSyncSession) Name() (string, error) {
	return _IInfoSync.Contract.Name(&_IInfoSync.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_IInfoSync *IInfoSyncCallerSession) Name() (string, error) {
	return _IInfoSync.Contract.Name(&_IInfoSync.CallOpts)
}

// Replenish is a paid mutator transaction binding the contract method 0x69ce93b4.
//
// Solidity: function replenish(uint64 chainID, uint32[] heights) returns(bool)
func (_IInfoSync *IInfoSyncTransactor) Replenish(opts *bind.TransactOpts, chainID uint64, heights []uint32) (*types.Transaction, error) {
	return _IInfoSync.contract.Transact(opts, "replenish", chainID, heights)
}

// Replenish is a paid mutator transaction binding the contract method 0x69ce93b4.
//
// Solidity: function replenish(uint64 chainID, uint32[] heights) returns(bool)
func (_IInfoSync *IInfoSyncSession) Replenish(chainID uint64, heights []uint32) (*types.Transaction, error) {
	return _IInfoSync.Contract.Replenish(&_IInfoSync.TransactOpts, chainID, heights)
}

// Replenish is a paid mutator transaction binding the contract method 0x69ce93b4.
//
// Solidity: function replenish(uint64 chainID, uint32[] heights) returns(bool)
func (_IInfoSync *IInfoSyncTransactorSession) Replenish(chainID uint64, heights []uint32) (*types.Transaction, error) {
	return _IInfoSync.Contract.Replenish(&_IInfoSync.TransactOpts, chainID, heights)
}

// SyncRootInfo is a paid mutator transaction binding the contract method 0x1413cc01.
//
// Solidity: function syncRootInfo(uint64 chainID, bytes[] rootInfos, bytes signature) returns(bool)
func (_IInfoSync *IInfoSyncTransactor) SyncRootInfo(opts *bind.TransactOpts, chainID uint64, rootInfos [][]byte, signature []byte) (*types.Transaction, error) {
	return _IInfoSync.contract.Transact(opts, "syncRootInfo", chainID, rootInfos, signature)
}

// SyncRootInfo is a paid mutator transaction binding the contract method 0x1413cc01.
//
// Solidity: function syncRootInfo(uint64 chainID, bytes[] rootInfos, bytes signature) returns(bool)
func (_IInfoSync *IInfoSyncSession) SyncRootInfo(chainID uint64, rootInfos [][]byte, signature []byte) (*types.Transaction, error) {
	return _IInfoSync.Contract.SyncRootInfo(&_IInfoSync.TransactOpts, chainID, rootInfos, signature)
}

// SyncRootInfo is a paid mutator transaction binding the contract method 0x1413cc01.
//
// Solidity: function syncRootInfo(uint64 chainID, bytes[] rootInfos, bytes signature) returns(bool)
func (_IInfoSync *IInfoSyncTransactorSession) SyncRootInfo(chainID uint64, rootInfos [][]byte, signature []byte) (*types.Transaction, error) {
	return _IInfoSync.Contract.SyncRootInfo(&_IInfoSync.TransactOpts, chainID, rootInfos, signature)
}

// IInfoSyncReplenishEventIterator is returned from FilterReplenishEvent and is used to iterate over the raw logs and unpacked data for ReplenishEvent events raised by the IInfoSync contract.
type IInfoSyncReplenishEventIterator struct {
	Event *IInfoSyncReplenishEvent // Event containing the contract specifics and raw log

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
func (it *IInfoSyncReplenishEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IInfoSyncReplenishEvent)
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
		it.Event = new(IInfoSyncReplenishEvent)
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
func (it *IInfoSyncReplenishEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IInfoSyncReplenishEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IInfoSyncReplenishEvent represents a ReplenishEvent event raised by the IInfoSync contract.
type IInfoSyncReplenishEvent struct {
	Heights []uint32
	ChainID uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterReplenishEvent is a free log retrieval operation binding the contract event 0x21af63b15a3c411234e8b1e975c9a0ea23ac5f43bca2f0a451042297560cd973.
//
// Solidity: event ReplenishEvent(uint32[] heights, uint64 chainID)
func (_IInfoSync *IInfoSyncFilterer) FilterReplenishEvent(opts *bind.FilterOpts) (*IInfoSyncReplenishEventIterator, error) {

	logs, sub, err := _IInfoSync.contract.FilterLogs(opts, "ReplenishEvent")
	if err != nil {
		return nil, err
	}
	return &IInfoSyncReplenishEventIterator{contract: _IInfoSync.contract, event: "ReplenishEvent", logs: logs, sub: sub}, nil
}

// WatchReplenishEvent is a free log subscription operation binding the contract event 0x21af63b15a3c411234e8b1e975c9a0ea23ac5f43bca2f0a451042297560cd973.
//
// Solidity: event ReplenishEvent(uint32[] heights, uint64 chainID)
func (_IInfoSync *IInfoSyncFilterer) WatchReplenishEvent(opts *bind.WatchOpts, sink chan<- *IInfoSyncReplenishEvent) (event.Subscription, error) {

	logs, sub, err := _IInfoSync.contract.WatchLogs(opts, "ReplenishEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IInfoSyncReplenishEvent)
				if err := _IInfoSync.contract.UnpackLog(event, "ReplenishEvent", log); err != nil {
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

// ParseReplenishEvent is a log parse operation binding the contract event 0x21af63b15a3c411234e8b1e975c9a0ea23ac5f43bca2f0a451042297560cd973.
//
// Solidity: event ReplenishEvent(uint32[] heights, uint64 chainID)
func (_IInfoSync *IInfoSyncFilterer) ParseReplenishEvent(log types.Log) (*IInfoSyncReplenishEvent, error) {
	event := new(IInfoSyncReplenishEvent)
	if err := _IInfoSync.contract.UnpackLog(event, "ReplenishEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IInfoSyncSyncRootInfoEventIterator is returned from FilterSyncRootInfoEvent and is used to iterate over the raw logs and unpacked data for SyncRootInfoEvent events raised by the IInfoSync contract.
type IInfoSyncSyncRootInfoEventIterator struct {
	Event *IInfoSyncSyncRootInfoEvent // Event containing the contract specifics and raw log

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
func (it *IInfoSyncSyncRootInfoEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IInfoSyncSyncRootInfoEvent)
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
		it.Event = new(IInfoSyncSyncRootInfoEvent)
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
func (it *IInfoSyncSyncRootInfoEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IInfoSyncSyncRootInfoEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IInfoSyncSyncRootInfoEvent represents a SyncRootInfoEvent event raised by the IInfoSync contract.
type IInfoSyncSyncRootInfoEvent struct {
	ChainID     uint64
	Height      uint32
	BlockHeight *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterSyncRootInfoEvent is a free log retrieval operation binding the contract event 0x24cd873c72270d09dcf64a7c56ff70129c04f398e38116d7ad773d43340509e2.
//
// Solidity: event SyncRootInfoEvent(uint64 chainID, uint32 height, uint256 BlockHeight)
func (_IInfoSync *IInfoSyncFilterer) FilterSyncRootInfoEvent(opts *bind.FilterOpts) (*IInfoSyncSyncRootInfoEventIterator, error) {

	logs, sub, err := _IInfoSync.contract.FilterLogs(opts, "SyncRootInfoEvent")
	if err != nil {
		return nil, err
	}
	return &IInfoSyncSyncRootInfoEventIterator{contract: _IInfoSync.contract, event: "SyncRootInfoEvent", logs: logs, sub: sub}, nil
}

// WatchSyncRootInfoEvent is a free log subscription operation binding the contract event 0x24cd873c72270d09dcf64a7c56ff70129c04f398e38116d7ad773d43340509e2.
//
// Solidity: event SyncRootInfoEvent(uint64 chainID, uint32 height, uint256 BlockHeight)
func (_IInfoSync *IInfoSyncFilterer) WatchSyncRootInfoEvent(opts *bind.WatchOpts, sink chan<- *IInfoSyncSyncRootInfoEvent) (event.Subscription, error) {

	logs, sub, err := _IInfoSync.contract.WatchLogs(opts, "SyncRootInfoEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IInfoSyncSyncRootInfoEvent)
				if err := _IInfoSync.contract.UnpackLog(event, "SyncRootInfoEvent", log); err != nil {
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

// ParseSyncRootInfoEvent is a log parse operation binding the contract event 0x24cd873c72270d09dcf64a7c56ff70129c04f398e38116d7ad773d43340509e2.
//
// Solidity: event SyncRootInfoEvent(uint64 chainID, uint32 height, uint256 BlockHeight)
func (_IInfoSync *IInfoSyncFilterer) ParseSyncRootInfoEvent(log types.Log) (*IInfoSyncSyncRootInfoEvent, error) {
	event := new(IInfoSyncSyncRootInfoEvent)
	if err := _IInfoSync.contract.UnpackLog(event, "SyncRootInfoEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
