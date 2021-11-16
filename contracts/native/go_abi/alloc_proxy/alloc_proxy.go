// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package alloc_proxy

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
	MethodBurn = "burn"

	MethodChangeEpoch = "changeEpoch"

	MethodInitGenesisHeader = "initGenesisHeader"

	MethodVerifyHeaderAndMint = "verifyHeaderAndMint"

	MethodName = "name"

	EventBurnEvent = "BurnEvent"

	EventChangeEpochEvent = "ChangeEpochEvent"

	EventInitGenesisBlockEvent = "InitGenesisBlockEvent"

	EventMintEvent = "MintEvent"
)

// IAllocProxyABI is the input ABI used to generate the binding from.
const IAllocProxyABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fromAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"crossTxId\",\"type\":\"bytes\"}],\"name\":\"BurnEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"header\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"oldEpoch\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newEpoch\",\"type\":\"bytes\"}],\"name\":\"ChangeEpochEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"header\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"epoch\",\"type\":\"bytes\"}],\"name\":\"InitGenesisBlockEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fromAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MintEvent\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"header\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extra\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"epoch\",\"type\":\"bytes\"}],\"name\":\"changeEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"header\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extra\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"epoch\",\"type\":\"bytes\"}],\"name\":\"initGenesisHeader\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"header\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rawCrossTx\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extra\",\"type\":\"bytes\"}],\"name\":\"verifyHeaderAndMint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IAllocProxyFuncSigs maps the 4-byte function signature to its string representation.
var IAllocProxyFuncSigs = map[string]string{
	"a4fa3313": "burn(uint64,address,uint256)",
	"c5107f3a": "changeEpoch(bytes,bytes,bytes,bytes)",
	"d4e025ef": "initGenesisHeader(bytes,bytes,bytes,bytes)",
	"06fdde03": "name()",
	"12979116": "verifyHeaderAndMint(bytes,bytes,bytes,bytes)",
}

// IAllocProxy is an auto generated Go binding around an Ethereum contract.
type IAllocProxy struct {
	IAllocProxyCaller     // Read-only binding to the contract
	IAllocProxyTransactor // Write-only binding to the contract
	IAllocProxyFilterer   // Log filterer for contract events
}

// IAllocProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type IAllocProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAllocProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IAllocProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAllocProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IAllocProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAllocProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IAllocProxySession struct {
	Contract     *IAllocProxy      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IAllocProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IAllocProxyCallerSession struct {
	Contract *IAllocProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// IAllocProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IAllocProxyTransactorSession struct {
	Contract     *IAllocProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// IAllocProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type IAllocProxyRaw struct {
	Contract *IAllocProxy // Generic contract binding to access the raw methods on
}

// IAllocProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IAllocProxyCallerRaw struct {
	Contract *IAllocProxyCaller // Generic read-only contract binding to access the raw methods on
}

// IAllocProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IAllocProxyTransactorRaw struct {
	Contract *IAllocProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIAllocProxy creates a new instance of IAllocProxy, bound to a specific deployed contract.
func NewIAllocProxy(address common.Address, backend bind.ContractBackend) (*IAllocProxy, error) {
	contract, err := bindIAllocProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAllocProxy{IAllocProxyCaller: IAllocProxyCaller{contract: contract}, IAllocProxyTransactor: IAllocProxyTransactor{contract: contract}, IAllocProxyFilterer: IAllocProxyFilterer{contract: contract}}, nil
}

// NewIAllocProxyCaller creates a new read-only instance of IAllocProxy, bound to a specific deployed contract.
func NewIAllocProxyCaller(address common.Address, caller bind.ContractCaller) (*IAllocProxyCaller, error) {
	contract, err := bindIAllocProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAllocProxyCaller{contract: contract}, nil
}

// NewIAllocProxyTransactor creates a new write-only instance of IAllocProxy, bound to a specific deployed contract.
func NewIAllocProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*IAllocProxyTransactor, error) {
	contract, err := bindIAllocProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAllocProxyTransactor{contract: contract}, nil
}

// NewIAllocProxyFilterer creates a new log filterer instance of IAllocProxy, bound to a specific deployed contract.
func NewIAllocProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*IAllocProxyFilterer, error) {
	contract, err := bindIAllocProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAllocProxyFilterer{contract: contract}, nil
}

// bindIAllocProxy binds a generic wrapper to an already deployed contract.
func bindIAllocProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IAllocProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAllocProxy *IAllocProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAllocProxy.Contract.IAllocProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAllocProxy *IAllocProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAllocProxy.Contract.IAllocProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAllocProxy *IAllocProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAllocProxy.Contract.IAllocProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAllocProxy *IAllocProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAllocProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAllocProxy *IAllocProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAllocProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAllocProxy *IAllocProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAllocProxy.Contract.contract.Transact(opts, method, params...)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_IAllocProxy *IAllocProxyCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IAllocProxy.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_IAllocProxy *IAllocProxySession) Name() (string, error) {
	return _IAllocProxy.Contract.Name(&_IAllocProxy.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_IAllocProxy *IAllocProxyCallerSession) Name() (string, error) {
	return _IAllocProxy.Contract.Name(&_IAllocProxy.CallOpts)
}

// Burn is a paid mutator transaction binding the contract method 0xa4fa3313.
//
// Solidity: function burn(uint64 toChainId, address toAddress, uint256 amount) returns(bool)
func (_IAllocProxy *IAllocProxyTransactor) Burn(opts *bind.TransactOpts, toChainId uint64, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAllocProxy.contract.Transact(opts, "burn", toChainId, toAddress, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xa4fa3313.
//
// Solidity: function burn(uint64 toChainId, address toAddress, uint256 amount) returns(bool)
func (_IAllocProxy *IAllocProxySession) Burn(toChainId uint64, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAllocProxy.Contract.Burn(&_IAllocProxy.TransactOpts, toChainId, toAddress, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xa4fa3313.
//
// Solidity: function burn(uint64 toChainId, address toAddress, uint256 amount) returns(bool)
func (_IAllocProxy *IAllocProxyTransactorSession) Burn(toChainId uint64, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAllocProxy.Contract.Burn(&_IAllocProxy.TransactOpts, toChainId, toAddress, amount)
}

// ChangeEpoch is a paid mutator transaction binding the contract method 0xc5107f3a.
//
// Solidity: function changeEpoch(bytes header, bytes proof, bytes extra, bytes epoch) returns(bool)
func (_IAllocProxy *IAllocProxyTransactor) ChangeEpoch(opts *bind.TransactOpts, header []byte, proof []byte, extra []byte, epoch []byte) (*types.Transaction, error) {
	return _IAllocProxy.contract.Transact(opts, "changeEpoch", header, proof, extra, epoch)
}

// ChangeEpoch is a paid mutator transaction binding the contract method 0xc5107f3a.
//
// Solidity: function changeEpoch(bytes header, bytes proof, bytes extra, bytes epoch) returns(bool)
func (_IAllocProxy *IAllocProxySession) ChangeEpoch(header []byte, proof []byte, extra []byte, epoch []byte) (*types.Transaction, error) {
	return _IAllocProxy.Contract.ChangeEpoch(&_IAllocProxy.TransactOpts, header, proof, extra, epoch)
}

// ChangeEpoch is a paid mutator transaction binding the contract method 0xc5107f3a.
//
// Solidity: function changeEpoch(bytes header, bytes proof, bytes extra, bytes epoch) returns(bool)
func (_IAllocProxy *IAllocProxyTransactorSession) ChangeEpoch(header []byte, proof []byte, extra []byte, epoch []byte) (*types.Transaction, error) {
	return _IAllocProxy.Contract.ChangeEpoch(&_IAllocProxy.TransactOpts, header, proof, extra, epoch)
}

// InitGenesisHeader is a paid mutator transaction binding the contract method 0xd4e025ef.
//
// Solidity: function initGenesisHeader(bytes header, bytes proof, bytes extra, bytes epoch) returns(bool)
func (_IAllocProxy *IAllocProxyTransactor) InitGenesisHeader(opts *bind.TransactOpts, header []byte, proof []byte, extra []byte, epoch []byte) (*types.Transaction, error) {
	return _IAllocProxy.contract.Transact(opts, "initGenesisHeader", header, proof, extra, epoch)
}

// InitGenesisHeader is a paid mutator transaction binding the contract method 0xd4e025ef.
//
// Solidity: function initGenesisHeader(bytes header, bytes proof, bytes extra, bytes epoch) returns(bool)
func (_IAllocProxy *IAllocProxySession) InitGenesisHeader(header []byte, proof []byte, extra []byte, epoch []byte) (*types.Transaction, error) {
	return _IAllocProxy.Contract.InitGenesisHeader(&_IAllocProxy.TransactOpts, header, proof, extra, epoch)
}

// InitGenesisHeader is a paid mutator transaction binding the contract method 0xd4e025ef.
//
// Solidity: function initGenesisHeader(bytes header, bytes proof, bytes extra, bytes epoch) returns(bool)
func (_IAllocProxy *IAllocProxyTransactorSession) InitGenesisHeader(header []byte, proof []byte, extra []byte, epoch []byte) (*types.Transaction, error) {
	return _IAllocProxy.Contract.InitGenesisHeader(&_IAllocProxy.TransactOpts, header, proof, extra, epoch)
}

// VerifyHeaderAndMint is a paid mutator transaction binding the contract method 0x12979116.
//
// Solidity: function verifyHeaderAndMint(bytes header, bytes rawCrossTx, bytes proof, bytes extra) returns(bool)
func (_IAllocProxy *IAllocProxyTransactor) VerifyHeaderAndMint(opts *bind.TransactOpts, header []byte, rawCrossTx []byte, proof []byte, extra []byte) (*types.Transaction, error) {
	return _IAllocProxy.contract.Transact(opts, "verifyHeaderAndMint", header, rawCrossTx, proof, extra)
}

// VerifyHeaderAndMint is a paid mutator transaction binding the contract method 0x12979116.
//
// Solidity: function verifyHeaderAndMint(bytes header, bytes rawCrossTx, bytes proof, bytes extra) returns(bool)
func (_IAllocProxy *IAllocProxySession) VerifyHeaderAndMint(header []byte, rawCrossTx []byte, proof []byte, extra []byte) (*types.Transaction, error) {
	return _IAllocProxy.Contract.VerifyHeaderAndMint(&_IAllocProxy.TransactOpts, header, rawCrossTx, proof, extra)
}

// VerifyHeaderAndMint is a paid mutator transaction binding the contract method 0x12979116.
//
// Solidity: function verifyHeaderAndMint(bytes header, bytes rawCrossTx, bytes proof, bytes extra) returns(bool)
func (_IAllocProxy *IAllocProxyTransactorSession) VerifyHeaderAndMint(header []byte, rawCrossTx []byte, proof []byte, extra []byte) (*types.Transaction, error) {
	return _IAllocProxy.Contract.VerifyHeaderAndMint(&_IAllocProxy.TransactOpts, header, rawCrossTx, proof, extra)
}

// IAllocProxyBurnEventIterator is returned from FilterBurnEvent and is used to iterate over the raw logs and unpacked data for BurnEvent events raised by the IAllocProxy contract.
type IAllocProxyBurnEventIterator struct {
	Event *IAllocProxyBurnEvent // Event containing the contract specifics and raw log

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
func (it *IAllocProxyBurnEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAllocProxyBurnEvent)
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
		it.Event = new(IAllocProxyBurnEvent)
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
func (it *IAllocProxyBurnEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAllocProxyBurnEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAllocProxyBurnEvent represents a BurnEvent event raised by the IAllocProxy contract.
type IAllocProxyBurnEvent struct {
	ToChainId   uint64
	FromAddress common.Address
	ToAddress   common.Address
	Amount      *big.Int
	CrossTxId   []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBurnEvent is a free log retrieval operation binding the contract event 0x6f8d47349ac6a71905a129c9f07798e8cd37575aeadfe3af6a18948232d7878f.
//
// Solidity: event BurnEvent(uint64 toChainId, address fromAddress, address toAddress, uint256 amount, bytes crossTxId)
func (_IAllocProxy *IAllocProxyFilterer) FilterBurnEvent(opts *bind.FilterOpts) (*IAllocProxyBurnEventIterator, error) {

	logs, sub, err := _IAllocProxy.contract.FilterLogs(opts, "BurnEvent")
	if err != nil {
		return nil, err
	}
	return &IAllocProxyBurnEventIterator{contract: _IAllocProxy.contract, event: "BurnEvent", logs: logs, sub: sub}, nil
}

// WatchBurnEvent is a free log subscription operation binding the contract event 0x6f8d47349ac6a71905a129c9f07798e8cd37575aeadfe3af6a18948232d7878f.
//
// Solidity: event BurnEvent(uint64 toChainId, address fromAddress, address toAddress, uint256 amount, bytes crossTxId)
func (_IAllocProxy *IAllocProxyFilterer) WatchBurnEvent(opts *bind.WatchOpts, sink chan<- *IAllocProxyBurnEvent) (event.Subscription, error) {

	logs, sub, err := _IAllocProxy.contract.WatchLogs(opts, "BurnEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAllocProxyBurnEvent)
				if err := _IAllocProxy.contract.UnpackLog(event, "BurnEvent", log); err != nil {
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

// ParseBurnEvent is a log parse operation binding the contract event 0x6f8d47349ac6a71905a129c9f07798e8cd37575aeadfe3af6a18948232d7878f.
//
// Solidity: event BurnEvent(uint64 toChainId, address fromAddress, address toAddress, uint256 amount, bytes crossTxId)
func (_IAllocProxy *IAllocProxyFilterer) ParseBurnEvent(log types.Log) (*IAllocProxyBurnEvent, error) {
	event := new(IAllocProxyBurnEvent)
	if err := _IAllocProxy.contract.UnpackLog(event, "BurnEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAllocProxyChangeEpochEventIterator is returned from FilterChangeEpochEvent and is used to iterate over the raw logs and unpacked data for ChangeEpochEvent events raised by the IAllocProxy contract.
type IAllocProxyChangeEpochEventIterator struct {
	Event *IAllocProxyChangeEpochEvent // Event containing the contract specifics and raw log

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
func (it *IAllocProxyChangeEpochEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAllocProxyChangeEpochEvent)
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
		it.Event = new(IAllocProxyChangeEpochEvent)
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
func (it *IAllocProxyChangeEpochEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAllocProxyChangeEpochEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAllocProxyChangeEpochEvent represents a ChangeEpochEvent event raised by the IAllocProxy contract.
type IAllocProxyChangeEpochEvent struct {
	Height   *big.Int
	Header   []byte
	OldEpoch []byte
	NewEpoch []byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterChangeEpochEvent is a free log retrieval operation binding the contract event 0xef6229e7d9f33cdb9aafa4a1ec2deb98baa7fd31ffb3e15a6371979f49e954b5.
//
// Solidity: event ChangeEpochEvent(uint256 height, bytes header, bytes oldEpoch, bytes newEpoch)
func (_IAllocProxy *IAllocProxyFilterer) FilterChangeEpochEvent(opts *bind.FilterOpts) (*IAllocProxyChangeEpochEventIterator, error) {

	logs, sub, err := _IAllocProxy.contract.FilterLogs(opts, "ChangeEpochEvent")
	if err != nil {
		return nil, err
	}
	return &IAllocProxyChangeEpochEventIterator{contract: _IAllocProxy.contract, event: "ChangeEpochEvent", logs: logs, sub: sub}, nil
}

// WatchChangeEpochEvent is a free log subscription operation binding the contract event 0xef6229e7d9f33cdb9aafa4a1ec2deb98baa7fd31ffb3e15a6371979f49e954b5.
//
// Solidity: event ChangeEpochEvent(uint256 height, bytes header, bytes oldEpoch, bytes newEpoch)
func (_IAllocProxy *IAllocProxyFilterer) WatchChangeEpochEvent(opts *bind.WatchOpts, sink chan<- *IAllocProxyChangeEpochEvent) (event.Subscription, error) {

	logs, sub, err := _IAllocProxy.contract.WatchLogs(opts, "ChangeEpochEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAllocProxyChangeEpochEvent)
				if err := _IAllocProxy.contract.UnpackLog(event, "ChangeEpochEvent", log); err != nil {
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

// ParseChangeEpochEvent is a log parse operation binding the contract event 0xef6229e7d9f33cdb9aafa4a1ec2deb98baa7fd31ffb3e15a6371979f49e954b5.
//
// Solidity: event ChangeEpochEvent(uint256 height, bytes header, bytes oldEpoch, bytes newEpoch)
func (_IAllocProxy *IAllocProxyFilterer) ParseChangeEpochEvent(log types.Log) (*IAllocProxyChangeEpochEvent, error) {
	event := new(IAllocProxyChangeEpochEvent)
	if err := _IAllocProxy.contract.UnpackLog(event, "ChangeEpochEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAllocProxyInitGenesisBlockEventIterator is returned from FilterInitGenesisBlockEvent and is used to iterate over the raw logs and unpacked data for InitGenesisBlockEvent events raised by the IAllocProxy contract.
type IAllocProxyInitGenesisBlockEventIterator struct {
	Event *IAllocProxyInitGenesisBlockEvent // Event containing the contract specifics and raw log

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
func (it *IAllocProxyInitGenesisBlockEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAllocProxyInitGenesisBlockEvent)
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
		it.Event = new(IAllocProxyInitGenesisBlockEvent)
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
func (it *IAllocProxyInitGenesisBlockEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAllocProxyInitGenesisBlockEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAllocProxyInitGenesisBlockEvent represents a InitGenesisBlockEvent event raised by the IAllocProxy contract.
type IAllocProxyInitGenesisBlockEvent struct {
	Height *big.Int
	Header []byte
	Epoch  []byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterInitGenesisBlockEvent is a free log retrieval operation binding the contract event 0x0dab1016941cb99833e4adee14a432336685eb90c538bff97cfe2052be2de2c7.
//
// Solidity: event InitGenesisBlockEvent(uint256 height, bytes header, bytes epoch)
func (_IAllocProxy *IAllocProxyFilterer) FilterInitGenesisBlockEvent(opts *bind.FilterOpts) (*IAllocProxyInitGenesisBlockEventIterator, error) {

	logs, sub, err := _IAllocProxy.contract.FilterLogs(opts, "InitGenesisBlockEvent")
	if err != nil {
		return nil, err
	}
	return &IAllocProxyInitGenesisBlockEventIterator{contract: _IAllocProxy.contract, event: "InitGenesisBlockEvent", logs: logs, sub: sub}, nil
}

// WatchInitGenesisBlockEvent is a free log subscription operation binding the contract event 0x0dab1016941cb99833e4adee14a432336685eb90c538bff97cfe2052be2de2c7.
//
// Solidity: event InitGenesisBlockEvent(uint256 height, bytes header, bytes epoch)
func (_IAllocProxy *IAllocProxyFilterer) WatchInitGenesisBlockEvent(opts *bind.WatchOpts, sink chan<- *IAllocProxyInitGenesisBlockEvent) (event.Subscription, error) {

	logs, sub, err := _IAllocProxy.contract.WatchLogs(opts, "InitGenesisBlockEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAllocProxyInitGenesisBlockEvent)
				if err := _IAllocProxy.contract.UnpackLog(event, "InitGenesisBlockEvent", log); err != nil {
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

// ParseInitGenesisBlockEvent is a log parse operation binding the contract event 0x0dab1016941cb99833e4adee14a432336685eb90c538bff97cfe2052be2de2c7.
//
// Solidity: event InitGenesisBlockEvent(uint256 height, bytes header, bytes epoch)
func (_IAllocProxy *IAllocProxyFilterer) ParseInitGenesisBlockEvent(log types.Log) (*IAllocProxyInitGenesisBlockEvent, error) {
	event := new(IAllocProxyInitGenesisBlockEvent)
	if err := _IAllocProxy.contract.UnpackLog(event, "InitGenesisBlockEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAllocProxyMintEventIterator is returned from FilterMintEvent and is used to iterate over the raw logs and unpacked data for MintEvent events raised by the IAllocProxy contract.
type IAllocProxyMintEventIterator struct {
	Event *IAllocProxyMintEvent // Event containing the contract specifics and raw log

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
func (it *IAllocProxyMintEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAllocProxyMintEvent)
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
		it.Event = new(IAllocProxyMintEvent)
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
func (it *IAllocProxyMintEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAllocProxyMintEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAllocProxyMintEvent represents a MintEvent event raised by the IAllocProxy contract.
type IAllocProxyMintEvent struct {
	ToChainId   uint64
	FromAddress common.Address
	ToAddress   common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMintEvent is a free log retrieval operation binding the contract event 0x78ee4538b11131291669cab42be6d75e08b37b8f3dedbead5b1f23c753e4bd12.
//
// Solidity: event MintEvent(uint64 toChainId, address fromAddress, address toAddress, uint256 amount)
func (_IAllocProxy *IAllocProxyFilterer) FilterMintEvent(opts *bind.FilterOpts) (*IAllocProxyMintEventIterator, error) {

	logs, sub, err := _IAllocProxy.contract.FilterLogs(opts, "MintEvent")
	if err != nil {
		return nil, err
	}
	return &IAllocProxyMintEventIterator{contract: _IAllocProxy.contract, event: "MintEvent", logs: logs, sub: sub}, nil
}

// WatchMintEvent is a free log subscription operation binding the contract event 0x78ee4538b11131291669cab42be6d75e08b37b8f3dedbead5b1f23c753e4bd12.
//
// Solidity: event MintEvent(uint64 toChainId, address fromAddress, address toAddress, uint256 amount)
func (_IAllocProxy *IAllocProxyFilterer) WatchMintEvent(opts *bind.WatchOpts, sink chan<- *IAllocProxyMintEvent) (event.Subscription, error) {

	logs, sub, err := _IAllocProxy.contract.WatchLogs(opts, "MintEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAllocProxyMintEvent)
				if err := _IAllocProxy.contract.UnpackLog(event, "MintEvent", log); err != nil {
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

// ParseMintEvent is a log parse operation binding the contract event 0x78ee4538b11131291669cab42be6d75e08b37b8f3dedbead5b1f23c753e4bd12.
//
// Solidity: event MintEvent(uint64 toChainId, address fromAddress, address toAddress, uint256 amount)
func (_IAllocProxy *IAllocProxyFilterer) ParseMintEvent(log types.Log) (*IAllocProxyMintEvent, error) {
	event := new(IAllocProxyMintEvent)
	if err := _IAllocProxy.contract.UnpackLog(event, "MintEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

