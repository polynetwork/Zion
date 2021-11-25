// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package side_chain_lock_proxy_abi

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

	MethodMint = "mint"

	MethodName = "name"

	EventBurnEvent = "BurnEvent"

	EventMintEvent = "MintEvent"
)

// ISideChainLockProxyABI is the input ABI used to generate the binding from.
const ISideChainLockProxyABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fromAssetHash\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fromAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"toAssetHash\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"toAddress\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BurnEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"toAssetHash\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MintEvent\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"argsBs\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"fromContractAddr\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"fromChainId\",\"type\":\"uint64\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// ISideChainLockProxyFuncSigs maps the 4-byte function signature to its string representation.
var ISideChainLockProxyFuncSigs = map[string]string{
	"a4fa3313": "burn(uint64,address,uint256)",
	"48e6dbbb": "mint(bytes,bytes,uint64)",
	"06fdde03": "name()",
}

// ISideChainLockProxy is an auto generated Go binding around an Ethereum contract.
type ISideChainLockProxy struct {
	ISideChainLockProxyCaller     // Read-only binding to the contract
	ISideChainLockProxyTransactor // Write-only binding to the contract
	ISideChainLockProxyFilterer   // Log filterer for contract events
}

// ISideChainLockProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type ISideChainLockProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISideChainLockProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ISideChainLockProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISideChainLockProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ISideChainLockProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISideChainLockProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ISideChainLockProxySession struct {
	Contract     *ISideChainLockProxy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ISideChainLockProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ISideChainLockProxyCallerSession struct {
	Contract *ISideChainLockProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// ISideChainLockProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ISideChainLockProxyTransactorSession struct {
	Contract     *ISideChainLockProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ISideChainLockProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type ISideChainLockProxyRaw struct {
	Contract *ISideChainLockProxy // Generic contract binding to access the raw methods on
}

// ISideChainLockProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ISideChainLockProxyCallerRaw struct {
	Contract *ISideChainLockProxyCaller // Generic read-only contract binding to access the raw methods on
}

// ISideChainLockProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ISideChainLockProxyTransactorRaw struct {
	Contract *ISideChainLockProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewISideChainLockProxy creates a new instance of ISideChainLockProxy, bound to a specific deployed contract.
func NewISideChainLockProxy(address common.Address, backend bind.ContractBackend) (*ISideChainLockProxy, error) {
	contract, err := bindISideChainLockProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ISideChainLockProxy{ISideChainLockProxyCaller: ISideChainLockProxyCaller{contract: contract}, ISideChainLockProxyTransactor: ISideChainLockProxyTransactor{contract: contract}, ISideChainLockProxyFilterer: ISideChainLockProxyFilterer{contract: contract}}, nil
}

// NewISideChainLockProxyCaller creates a new read-only instance of ISideChainLockProxy, bound to a specific deployed contract.
func NewISideChainLockProxyCaller(address common.Address, caller bind.ContractCaller) (*ISideChainLockProxyCaller, error) {
	contract, err := bindISideChainLockProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ISideChainLockProxyCaller{contract: contract}, nil
}

// NewISideChainLockProxyTransactor creates a new write-only instance of ISideChainLockProxy, bound to a specific deployed contract.
func NewISideChainLockProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*ISideChainLockProxyTransactor, error) {
	contract, err := bindISideChainLockProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ISideChainLockProxyTransactor{contract: contract}, nil
}

// NewISideChainLockProxyFilterer creates a new log filterer instance of ISideChainLockProxy, bound to a specific deployed contract.
func NewISideChainLockProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*ISideChainLockProxyFilterer, error) {
	contract, err := bindISideChainLockProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ISideChainLockProxyFilterer{contract: contract}, nil
}

// bindISideChainLockProxy binds a generic wrapper to an already deployed contract.
func bindISideChainLockProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ISideChainLockProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISideChainLockProxy *ISideChainLockProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISideChainLockProxy.Contract.ISideChainLockProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISideChainLockProxy *ISideChainLockProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISideChainLockProxy.Contract.ISideChainLockProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISideChainLockProxy *ISideChainLockProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISideChainLockProxy.Contract.ISideChainLockProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISideChainLockProxy *ISideChainLockProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISideChainLockProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISideChainLockProxy *ISideChainLockProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISideChainLockProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISideChainLockProxy *ISideChainLockProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISideChainLockProxy.Contract.contract.Transact(opts, method, params...)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ISideChainLockProxy *ISideChainLockProxyCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ISideChainLockProxy.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ISideChainLockProxy *ISideChainLockProxySession) Name() (string, error) {
	return _ISideChainLockProxy.Contract.Name(&_ISideChainLockProxy.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ISideChainLockProxy *ISideChainLockProxyCallerSession) Name() (string, error) {
	return _ISideChainLockProxy.Contract.Name(&_ISideChainLockProxy.CallOpts)
}

// Burn is a paid mutator transaction binding the contract method 0xa4fa3313.
//
// Solidity: function burn(uint64 toChainId, address toAddress, uint256 amount) returns(bool)
func (_ISideChainLockProxy *ISideChainLockProxyTransactor) Burn(opts *bind.TransactOpts, toChainId uint64, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ISideChainLockProxy.contract.Transact(opts, "burn", toChainId, toAddress, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xa4fa3313.
//
// Solidity: function burn(uint64 toChainId, address toAddress, uint256 amount) returns(bool)
func (_ISideChainLockProxy *ISideChainLockProxySession) Burn(toChainId uint64, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ISideChainLockProxy.Contract.Burn(&_ISideChainLockProxy.TransactOpts, toChainId, toAddress, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xa4fa3313.
//
// Solidity: function burn(uint64 toChainId, address toAddress, uint256 amount) returns(bool)
func (_ISideChainLockProxy *ISideChainLockProxyTransactorSession) Burn(toChainId uint64, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ISideChainLockProxy.Contract.Burn(&_ISideChainLockProxy.TransactOpts, toChainId, toAddress, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x48e6dbbb.
//
// Solidity: function mint(bytes argsBs, bytes fromContractAddr, uint64 fromChainId) returns(bool)
func (_ISideChainLockProxy *ISideChainLockProxyTransactor) Mint(opts *bind.TransactOpts, argsBs []byte, fromContractAddr []byte, fromChainId uint64) (*types.Transaction, error) {
	return _ISideChainLockProxy.contract.Transact(opts, "mint", argsBs, fromContractAddr, fromChainId)
}

// Mint is a paid mutator transaction binding the contract method 0x48e6dbbb.
//
// Solidity: function mint(bytes argsBs, bytes fromContractAddr, uint64 fromChainId) returns(bool)
func (_ISideChainLockProxy *ISideChainLockProxySession) Mint(argsBs []byte, fromContractAddr []byte, fromChainId uint64) (*types.Transaction, error) {
	return _ISideChainLockProxy.Contract.Mint(&_ISideChainLockProxy.TransactOpts, argsBs, fromContractAddr, fromChainId)
}

// Mint is a paid mutator transaction binding the contract method 0x48e6dbbb.
//
// Solidity: function mint(bytes argsBs, bytes fromContractAddr, uint64 fromChainId) returns(bool)
func (_ISideChainLockProxy *ISideChainLockProxyTransactorSession) Mint(argsBs []byte, fromContractAddr []byte, fromChainId uint64) (*types.Transaction, error) {
	return _ISideChainLockProxy.Contract.Mint(&_ISideChainLockProxy.TransactOpts, argsBs, fromContractAddr, fromChainId)
}

// ISideChainLockProxyBurnEventIterator is returned from FilterBurnEvent and is used to iterate over the raw logs and unpacked data for BurnEvent events raised by the ISideChainLockProxy contract.
type ISideChainLockProxyBurnEventIterator struct {
	Event *ISideChainLockProxyBurnEvent // Event containing the contract specifics and raw log

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
func (it *ISideChainLockProxyBurnEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ISideChainLockProxyBurnEvent)
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
		it.Event = new(ISideChainLockProxyBurnEvent)
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
func (it *ISideChainLockProxyBurnEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ISideChainLockProxyBurnEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ISideChainLockProxyBurnEvent represents a BurnEvent event raised by the ISideChainLockProxy contract.
type ISideChainLockProxyBurnEvent struct {
	FromAssetHash common.Address
	FromAddress   common.Address
	ToChainId     uint64
	ToAssetHash   []byte
	ToAddress     []byte
	Amount        *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterBurnEvent is a free log retrieval operation binding the contract event 0x9f6f5896351abb9d6af7998e408c5b94b906038aaac69f1d6da63d395f2a2ab3.
//
// Solidity: event BurnEvent(address fromAssetHash, address fromAddress, uint64 toChainId, bytes toAssetHash, bytes toAddress, uint256 amount)
func (_ISideChainLockProxy *ISideChainLockProxyFilterer) FilterBurnEvent(opts *bind.FilterOpts) (*ISideChainLockProxyBurnEventIterator, error) {

	logs, sub, err := _ISideChainLockProxy.contract.FilterLogs(opts, "BurnEvent")
	if err != nil {
		return nil, err
	}
	return &ISideChainLockProxyBurnEventIterator{contract: _ISideChainLockProxy.contract, event: "BurnEvent", logs: logs, sub: sub}, nil
}

// WatchBurnEvent is a free log subscription operation binding the contract event 0x9f6f5896351abb9d6af7998e408c5b94b906038aaac69f1d6da63d395f2a2ab3.
//
// Solidity: event BurnEvent(address fromAssetHash, address fromAddress, uint64 toChainId, bytes toAssetHash, bytes toAddress, uint256 amount)
func (_ISideChainLockProxy *ISideChainLockProxyFilterer) WatchBurnEvent(opts *bind.WatchOpts, sink chan<- *ISideChainLockProxyBurnEvent) (event.Subscription, error) {

	logs, sub, err := _ISideChainLockProxy.contract.WatchLogs(opts, "BurnEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ISideChainLockProxyBurnEvent)
				if err := _ISideChainLockProxy.contract.UnpackLog(event, "BurnEvent", log); err != nil {
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

// ParseBurnEvent is a log parse operation binding the contract event 0x9f6f5896351abb9d6af7998e408c5b94b906038aaac69f1d6da63d395f2a2ab3.
//
// Solidity: event BurnEvent(address fromAssetHash, address fromAddress, uint64 toChainId, bytes toAssetHash, bytes toAddress, uint256 amount)
func (_ISideChainLockProxy *ISideChainLockProxyFilterer) ParseBurnEvent(log types.Log) (*ISideChainLockProxyBurnEvent, error) {
	event := new(ISideChainLockProxyBurnEvent)
	if err := _ISideChainLockProxy.contract.UnpackLog(event, "BurnEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ISideChainLockProxyMintEventIterator is returned from FilterMintEvent and is used to iterate over the raw logs and unpacked data for MintEvent events raised by the ISideChainLockProxy contract.
type ISideChainLockProxyMintEventIterator struct {
	Event *ISideChainLockProxyMintEvent // Event containing the contract specifics and raw log

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
func (it *ISideChainLockProxyMintEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ISideChainLockProxyMintEvent)
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
		it.Event = new(ISideChainLockProxyMintEvent)
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
func (it *ISideChainLockProxyMintEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ISideChainLockProxyMintEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ISideChainLockProxyMintEvent represents a MintEvent event raised by the ISideChainLockProxy contract.
type ISideChainLockProxyMintEvent struct {
	ToAssetHash common.Address
	ToAddress   common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMintEvent is a free log retrieval operation binding the contract event 0xa185a288bfeb0bc3ac58fe6994088736867f1a8ca58eecf2cd37978d51b9de6b.
//
// Solidity: event MintEvent(address toAssetHash, address toAddress, uint256 amount)
func (_ISideChainLockProxy *ISideChainLockProxyFilterer) FilterMintEvent(opts *bind.FilterOpts) (*ISideChainLockProxyMintEventIterator, error) {

	logs, sub, err := _ISideChainLockProxy.contract.FilterLogs(opts, "MintEvent")
	if err != nil {
		return nil, err
	}
	return &ISideChainLockProxyMintEventIterator{contract: _ISideChainLockProxy.contract, event: "MintEvent", logs: logs, sub: sub}, nil
}

// WatchMintEvent is a free log subscription operation binding the contract event 0xa185a288bfeb0bc3ac58fe6994088736867f1a8ca58eecf2cd37978d51b9de6b.
//
// Solidity: event MintEvent(address toAssetHash, address toAddress, uint256 amount)
func (_ISideChainLockProxy *ISideChainLockProxyFilterer) WatchMintEvent(opts *bind.WatchOpts, sink chan<- *ISideChainLockProxyMintEvent) (event.Subscription, error) {

	logs, sub, err := _ISideChainLockProxy.contract.WatchLogs(opts, "MintEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ISideChainLockProxyMintEvent)
				if err := _ISideChainLockProxy.contract.UnpackLog(event, "MintEvent", log); err != nil {
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

// ParseMintEvent is a log parse operation binding the contract event 0xa185a288bfeb0bc3ac58fe6994088736867f1a8ca58eecf2cd37978d51b9de6b.
//
// Solidity: event MintEvent(address toAssetHash, address toAddress, uint256 amount)
func (_ISideChainLockProxy *ISideChainLockProxyFilterer) ParseMintEvent(log types.Log) (*ISideChainLockProxyMintEvent, error) {
	event := new(ISideChainLockProxyMintEvent)
	if err := _ISideChainLockProxy.contract.UnpackLog(event, "MintEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
