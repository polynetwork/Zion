// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package lock_proxy

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
	MethodBindAssetHash = "bindAssetHash"

	MethodBindProxyHash = "bindProxyHash"

	MethodCrossChain = "crossChain"

	MethodLock = "lock"

	MethodUnlock = "unlock"

	MethodVerifyHeaderAndExecuteTx = "verifyHeaderAndExecuteTx"

	MethodName = "name"

	EventBindAssetEvent = "BindAssetEvent"

	EventBindProxyEvent = "BindProxyEvent"

	EventCrossChainEvent = "CrossChainEvent"

	EventLockEvent = "LockEvent"

	EventUnlockEvent = "UnlockEvent"

	EventVerifyHeaderAndExecuteTxEvent = "VerifyHeaderAndExecuteTxEvent"
)

// ILockProxyABI is the input ABI used to generate the binding from.
const ILockProxyABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fromAssetHash\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"targetProxyHash\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialAmount\",\"type\":\"uint256\"}],\"name\":\"BindAssetEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"targetProxyHash\",\"type\":\"bytes\"}],\"name\":\"BindProxyEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"txId\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proxyOrAssetContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"toContract\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"method\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"rawdata\",\"type\":\"bytes\"}],\"name\":\"CrossChainEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fromAssetHash\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fromAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"toAssetHash\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"toAddress\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LockEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"toAssetHash\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"UnlockEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"fromChainID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"toContract\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"crossChainTxHash\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"fromChainTxHash\",\"type\":\"bytes\"}],\"name\":\"VerifyHeaderAndExecuteTxEvent\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"fromAssetHash\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"toAssetHash\",\"type\":\"bytes\"}],\"name\":\"bindAssetHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"targetProxyHash\",\"type\":\"bytes\"}],\"name\":\"bindProxyHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"toContract\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"method\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"txData\",\"type\":\"bytes\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"fromAssetHash\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"toAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"lock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"argsBs\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"fromContractAddr\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"fromChainId\",\"type\":\"uint64\"}],\"name\":\"unlock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"header\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"verifyHeaderAndExecuteTx\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ILockProxyFuncSigs maps the 4-byte function signature to its string representation.
var ILockProxyFuncSigs = map[string]string{
	"3348f63b": "bindAssetHash(address,uint64,bytes)",
	"379b98f6": "bindProxyHash(uint64,bytes)",
	"bd5cf625": "crossChain(uint64,bytes,bytes,bytes)",
	"84a6d055": "lock(address,uint64,bytes,uint256)",
	"06fdde03": "name()",
	"06af4b9f": "unlock(bytes,bytes,uint64)",
	"a828f302": "verifyHeaderAndExecuteTx(bytes,bytes)",
}

// ILockProxy is an auto generated Go binding around an Ethereum contract.
type ILockProxy struct {
	ILockProxyCaller     // Read-only binding to the contract
	ILockProxyTransactor // Write-only binding to the contract
	ILockProxyFilterer   // Log filterer for contract events
}

// ILockProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type ILockProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ILockProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ILockProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ILockProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ILockProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ILockProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ILockProxySession struct {
	Contract     *ILockProxy       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ILockProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ILockProxyCallerSession struct {
	Contract *ILockProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ILockProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ILockProxyTransactorSession struct {
	Contract     *ILockProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ILockProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type ILockProxyRaw struct {
	Contract *ILockProxy // Generic contract binding to access the raw methods on
}

// ILockProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ILockProxyCallerRaw struct {
	Contract *ILockProxyCaller // Generic read-only contract binding to access the raw methods on
}

// ILockProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ILockProxyTransactorRaw struct {
	Contract *ILockProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewILockProxy creates a new instance of ILockProxy, bound to a specific deployed contract.
func NewILockProxy(address common.Address, backend bind.ContractBackend) (*ILockProxy, error) {
	contract, err := bindILockProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ILockProxy{ILockProxyCaller: ILockProxyCaller{contract: contract}, ILockProxyTransactor: ILockProxyTransactor{contract: contract}, ILockProxyFilterer: ILockProxyFilterer{contract: contract}}, nil
}

// NewILockProxyCaller creates a new read-only instance of ILockProxy, bound to a specific deployed contract.
func NewILockProxyCaller(address common.Address, caller bind.ContractCaller) (*ILockProxyCaller, error) {
	contract, err := bindILockProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ILockProxyCaller{contract: contract}, nil
}

// NewILockProxyTransactor creates a new write-only instance of ILockProxy, bound to a specific deployed contract.
func NewILockProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*ILockProxyTransactor, error) {
	contract, err := bindILockProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ILockProxyTransactor{contract: contract}, nil
}

// NewILockProxyFilterer creates a new log filterer instance of ILockProxy, bound to a specific deployed contract.
func NewILockProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*ILockProxyFilterer, error) {
	contract, err := bindILockProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ILockProxyFilterer{contract: contract}, nil
}

// bindILockProxy binds a generic wrapper to an already deployed contract.
func bindILockProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ILockProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ILockProxy *ILockProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ILockProxy.Contract.ILockProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ILockProxy *ILockProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ILockProxy.Contract.ILockProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ILockProxy *ILockProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ILockProxy.Contract.ILockProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ILockProxy *ILockProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ILockProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ILockProxy *ILockProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ILockProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ILockProxy *ILockProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ILockProxy.Contract.contract.Transact(opts, method, params...)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ILockProxy *ILockProxyCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ILockProxy.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ILockProxy *ILockProxySession) Name() (string, error) {
	return _ILockProxy.Contract.Name(&_ILockProxy.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ILockProxy *ILockProxyCallerSession) Name() (string, error) {
	return _ILockProxy.Contract.Name(&_ILockProxy.CallOpts)
}

// BindAssetHash is a paid mutator transaction binding the contract method 0x3348f63b.
//
// Solidity: function bindAssetHash(address fromAssetHash, uint64 toChainId, bytes toAssetHash) returns(bool)
func (_ILockProxy *ILockProxyTransactor) BindAssetHash(opts *bind.TransactOpts, fromAssetHash common.Address, toChainId uint64, toAssetHash []byte) (*types.Transaction, error) {
	return _ILockProxy.contract.Transact(opts, "bindAssetHash", fromAssetHash, toChainId, toAssetHash)
}

// BindAssetHash is a paid mutator transaction binding the contract method 0x3348f63b.
//
// Solidity: function bindAssetHash(address fromAssetHash, uint64 toChainId, bytes toAssetHash) returns(bool)
func (_ILockProxy *ILockProxySession) BindAssetHash(fromAssetHash common.Address, toChainId uint64, toAssetHash []byte) (*types.Transaction, error) {
	return _ILockProxy.Contract.BindAssetHash(&_ILockProxy.TransactOpts, fromAssetHash, toChainId, toAssetHash)
}

// BindAssetHash is a paid mutator transaction binding the contract method 0x3348f63b.
//
// Solidity: function bindAssetHash(address fromAssetHash, uint64 toChainId, bytes toAssetHash) returns(bool)
func (_ILockProxy *ILockProxyTransactorSession) BindAssetHash(fromAssetHash common.Address, toChainId uint64, toAssetHash []byte) (*types.Transaction, error) {
	return _ILockProxy.Contract.BindAssetHash(&_ILockProxy.TransactOpts, fromAssetHash, toChainId, toAssetHash)
}

// BindProxyHash is a paid mutator transaction binding the contract method 0x379b98f6.
//
// Solidity: function bindProxyHash(uint64 toChainId, bytes targetProxyHash) returns(bool)
func (_ILockProxy *ILockProxyTransactor) BindProxyHash(opts *bind.TransactOpts, toChainId uint64, targetProxyHash []byte) (*types.Transaction, error) {
	return _ILockProxy.contract.Transact(opts, "bindProxyHash", toChainId, targetProxyHash)
}

// BindProxyHash is a paid mutator transaction binding the contract method 0x379b98f6.
//
// Solidity: function bindProxyHash(uint64 toChainId, bytes targetProxyHash) returns(bool)
func (_ILockProxy *ILockProxySession) BindProxyHash(toChainId uint64, targetProxyHash []byte) (*types.Transaction, error) {
	return _ILockProxy.Contract.BindProxyHash(&_ILockProxy.TransactOpts, toChainId, targetProxyHash)
}

// BindProxyHash is a paid mutator transaction binding the contract method 0x379b98f6.
//
// Solidity: function bindProxyHash(uint64 toChainId, bytes targetProxyHash) returns(bool)
func (_ILockProxy *ILockProxyTransactorSession) BindProxyHash(toChainId uint64, targetProxyHash []byte) (*types.Transaction, error) {
	return _ILockProxy.Contract.BindProxyHash(&_ILockProxy.TransactOpts, toChainId, targetProxyHash)
}

// CrossChain is a paid mutator transaction binding the contract method 0xbd5cf625.
//
// Solidity: function crossChain(uint64 toChainId, bytes toContract, bytes method, bytes txData) returns(bool)
func (_ILockProxy *ILockProxyTransactor) CrossChain(opts *bind.TransactOpts, toChainId uint64, toContract []byte, method []byte, txData []byte) (*types.Transaction, error) {
	return _ILockProxy.contract.Transact(opts, "crossChain", toChainId, toContract, method, txData)
}

// CrossChain is a paid mutator transaction binding the contract method 0xbd5cf625.
//
// Solidity: function crossChain(uint64 toChainId, bytes toContract, bytes method, bytes txData) returns(bool)
func (_ILockProxy *ILockProxySession) CrossChain(toChainId uint64, toContract []byte, method []byte, txData []byte) (*types.Transaction, error) {
	return _ILockProxy.Contract.CrossChain(&_ILockProxy.TransactOpts, toChainId, toContract, method, txData)
}

// CrossChain is a paid mutator transaction binding the contract method 0xbd5cf625.
//
// Solidity: function crossChain(uint64 toChainId, bytes toContract, bytes method, bytes txData) returns(bool)
func (_ILockProxy *ILockProxyTransactorSession) CrossChain(toChainId uint64, toContract []byte, method []byte, txData []byte) (*types.Transaction, error) {
	return _ILockProxy.Contract.CrossChain(&_ILockProxy.TransactOpts, toChainId, toContract, method, txData)
}

// Lock is a paid mutator transaction binding the contract method 0x84a6d055.
//
// Solidity: function lock(address fromAssetHash, uint64 toChainId, bytes toAddress, uint256 amount) payable returns(bool)
func (_ILockProxy *ILockProxyTransactor) Lock(opts *bind.TransactOpts, fromAssetHash common.Address, toChainId uint64, toAddress []byte, amount *big.Int) (*types.Transaction, error) {
	return _ILockProxy.contract.Transact(opts, "lock", fromAssetHash, toChainId, toAddress, amount)
}

// Lock is a paid mutator transaction binding the contract method 0x84a6d055.
//
// Solidity: function lock(address fromAssetHash, uint64 toChainId, bytes toAddress, uint256 amount) payable returns(bool)
func (_ILockProxy *ILockProxySession) Lock(fromAssetHash common.Address, toChainId uint64, toAddress []byte, amount *big.Int) (*types.Transaction, error) {
	return _ILockProxy.Contract.Lock(&_ILockProxy.TransactOpts, fromAssetHash, toChainId, toAddress, amount)
}

// Lock is a paid mutator transaction binding the contract method 0x84a6d055.
//
// Solidity: function lock(address fromAssetHash, uint64 toChainId, bytes toAddress, uint256 amount) payable returns(bool)
func (_ILockProxy *ILockProxyTransactorSession) Lock(fromAssetHash common.Address, toChainId uint64, toAddress []byte, amount *big.Int) (*types.Transaction, error) {
	return _ILockProxy.Contract.Lock(&_ILockProxy.TransactOpts, fromAssetHash, toChainId, toAddress, amount)
}

// Unlock is a paid mutator transaction binding the contract method 0x06af4b9f.
//
// Solidity: function unlock(bytes argsBs, bytes fromContractAddr, uint64 fromChainId) returns(bool)
func (_ILockProxy *ILockProxyTransactor) Unlock(opts *bind.TransactOpts, argsBs []byte, fromContractAddr []byte, fromChainId uint64) (*types.Transaction, error) {
	return _ILockProxy.contract.Transact(opts, "unlock", argsBs, fromContractAddr, fromChainId)
}

// Unlock is a paid mutator transaction binding the contract method 0x06af4b9f.
//
// Solidity: function unlock(bytes argsBs, bytes fromContractAddr, uint64 fromChainId) returns(bool)
func (_ILockProxy *ILockProxySession) Unlock(argsBs []byte, fromContractAddr []byte, fromChainId uint64) (*types.Transaction, error) {
	return _ILockProxy.Contract.Unlock(&_ILockProxy.TransactOpts, argsBs, fromContractAddr, fromChainId)
}

// Unlock is a paid mutator transaction binding the contract method 0x06af4b9f.
//
// Solidity: function unlock(bytes argsBs, bytes fromContractAddr, uint64 fromChainId) returns(bool)
func (_ILockProxy *ILockProxyTransactorSession) Unlock(argsBs []byte, fromContractAddr []byte, fromChainId uint64) (*types.Transaction, error) {
	return _ILockProxy.Contract.Unlock(&_ILockProxy.TransactOpts, argsBs, fromContractAddr, fromChainId)
}

// VerifyHeaderAndExecuteTx is a paid mutator transaction binding the contract method 0xa828f302.
//
// Solidity: function verifyHeaderAndExecuteTx(bytes header, bytes proof) returns(bool)
func (_ILockProxy *ILockProxyTransactor) VerifyHeaderAndExecuteTx(opts *bind.TransactOpts, header []byte, proof []byte) (*types.Transaction, error) {
	return _ILockProxy.contract.Transact(opts, "verifyHeaderAndExecuteTx", header, proof)
}

// VerifyHeaderAndExecuteTx is a paid mutator transaction binding the contract method 0xa828f302.
//
// Solidity: function verifyHeaderAndExecuteTx(bytes header, bytes proof) returns(bool)
func (_ILockProxy *ILockProxySession) VerifyHeaderAndExecuteTx(header []byte, proof []byte) (*types.Transaction, error) {
	return _ILockProxy.Contract.VerifyHeaderAndExecuteTx(&_ILockProxy.TransactOpts, header, proof)
}

// VerifyHeaderAndExecuteTx is a paid mutator transaction binding the contract method 0xa828f302.
//
// Solidity: function verifyHeaderAndExecuteTx(bytes header, bytes proof) returns(bool)
func (_ILockProxy *ILockProxyTransactorSession) VerifyHeaderAndExecuteTx(header []byte, proof []byte) (*types.Transaction, error) {
	return _ILockProxy.Contract.VerifyHeaderAndExecuteTx(&_ILockProxy.TransactOpts, header, proof)
}

// ILockProxyBindAssetEventIterator is returned from FilterBindAssetEvent and is used to iterate over the raw logs and unpacked data for BindAssetEvent events raised by the ILockProxy contract.
type ILockProxyBindAssetEventIterator struct {
	Event *ILockProxyBindAssetEvent // Event containing the contract specifics and raw log

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
func (it *ILockProxyBindAssetEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ILockProxyBindAssetEvent)
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
		it.Event = new(ILockProxyBindAssetEvent)
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
func (it *ILockProxyBindAssetEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ILockProxyBindAssetEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ILockProxyBindAssetEvent represents a BindAssetEvent event raised by the ILockProxy contract.
type ILockProxyBindAssetEvent struct {
	FromAssetHash   common.Address
	ToChainId       uint64
	TargetProxyHash []byte
	InitialAmount   *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBindAssetEvent is a free log retrieval operation binding the contract event 0x1628c8374c1bdfeb2275fd9f4c90562fd3fae974783dc522c8234e36abcfc58e.
//
// Solidity: event BindAssetEvent(address fromAssetHash, uint64 toChainId, bytes targetProxyHash, uint256 initialAmount)
func (_ILockProxy *ILockProxyFilterer) FilterBindAssetEvent(opts *bind.FilterOpts) (*ILockProxyBindAssetEventIterator, error) {

	logs, sub, err := _ILockProxy.contract.FilterLogs(opts, "BindAssetEvent")
	if err != nil {
		return nil, err
	}
	return &ILockProxyBindAssetEventIterator{contract: _ILockProxy.contract, event: "BindAssetEvent", logs: logs, sub: sub}, nil
}

// WatchBindAssetEvent is a free log subscription operation binding the contract event 0x1628c8374c1bdfeb2275fd9f4c90562fd3fae974783dc522c8234e36abcfc58e.
//
// Solidity: event BindAssetEvent(address fromAssetHash, uint64 toChainId, bytes targetProxyHash, uint256 initialAmount)
func (_ILockProxy *ILockProxyFilterer) WatchBindAssetEvent(opts *bind.WatchOpts, sink chan<- *ILockProxyBindAssetEvent) (event.Subscription, error) {

	logs, sub, err := _ILockProxy.contract.WatchLogs(opts, "BindAssetEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ILockProxyBindAssetEvent)
				if err := _ILockProxy.contract.UnpackLog(event, "BindAssetEvent", log); err != nil {
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

// ParseBindAssetEvent is a log parse operation binding the contract event 0x1628c8374c1bdfeb2275fd9f4c90562fd3fae974783dc522c8234e36abcfc58e.
//
// Solidity: event BindAssetEvent(address fromAssetHash, uint64 toChainId, bytes targetProxyHash, uint256 initialAmount)
func (_ILockProxy *ILockProxyFilterer) ParseBindAssetEvent(log types.Log) (*ILockProxyBindAssetEvent, error) {
	event := new(ILockProxyBindAssetEvent)
	if err := _ILockProxy.contract.UnpackLog(event, "BindAssetEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ILockProxyBindProxyEventIterator is returned from FilterBindProxyEvent and is used to iterate over the raw logs and unpacked data for BindProxyEvent events raised by the ILockProxy contract.
type ILockProxyBindProxyEventIterator struct {
	Event *ILockProxyBindProxyEvent // Event containing the contract specifics and raw log

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
func (it *ILockProxyBindProxyEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ILockProxyBindProxyEvent)
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
		it.Event = new(ILockProxyBindProxyEvent)
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
func (it *ILockProxyBindProxyEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ILockProxyBindProxyEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ILockProxyBindProxyEvent represents a BindProxyEvent event raised by the ILockProxy contract.
type ILockProxyBindProxyEvent struct {
	ToChainId       uint64
	TargetProxyHash []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBindProxyEvent is a free log retrieval operation binding the contract event 0xdacd7d303272a3b58aec6620d6d1fb588f4996a5b46858ed437f1c34348f2d0f.
//
// Solidity: event BindProxyEvent(uint64 toChainId, bytes targetProxyHash)
func (_ILockProxy *ILockProxyFilterer) FilterBindProxyEvent(opts *bind.FilterOpts) (*ILockProxyBindProxyEventIterator, error) {

	logs, sub, err := _ILockProxy.contract.FilterLogs(opts, "BindProxyEvent")
	if err != nil {
		return nil, err
	}
	return &ILockProxyBindProxyEventIterator{contract: _ILockProxy.contract, event: "BindProxyEvent", logs: logs, sub: sub}, nil
}

// WatchBindProxyEvent is a free log subscription operation binding the contract event 0xdacd7d303272a3b58aec6620d6d1fb588f4996a5b46858ed437f1c34348f2d0f.
//
// Solidity: event BindProxyEvent(uint64 toChainId, bytes targetProxyHash)
func (_ILockProxy *ILockProxyFilterer) WatchBindProxyEvent(opts *bind.WatchOpts, sink chan<- *ILockProxyBindProxyEvent) (event.Subscription, error) {

	logs, sub, err := _ILockProxy.contract.WatchLogs(opts, "BindProxyEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ILockProxyBindProxyEvent)
				if err := _ILockProxy.contract.UnpackLog(event, "BindProxyEvent", log); err != nil {
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

// ParseBindProxyEvent is a log parse operation binding the contract event 0xdacd7d303272a3b58aec6620d6d1fb588f4996a5b46858ed437f1c34348f2d0f.
//
// Solidity: event BindProxyEvent(uint64 toChainId, bytes targetProxyHash)
func (_ILockProxy *ILockProxyFilterer) ParseBindProxyEvent(log types.Log) (*ILockProxyBindProxyEvent, error) {
	event := new(ILockProxyBindProxyEvent)
	if err := _ILockProxy.contract.UnpackLog(event, "BindProxyEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ILockProxyCrossChainEventIterator is returned from FilterCrossChainEvent and is used to iterate over the raw logs and unpacked data for CrossChainEvent events raised by the ILockProxy contract.
type ILockProxyCrossChainEventIterator struct {
	Event *ILockProxyCrossChainEvent // Event containing the contract specifics and raw log

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
func (it *ILockProxyCrossChainEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ILockProxyCrossChainEvent)
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
		it.Event = new(ILockProxyCrossChainEvent)
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
func (it *ILockProxyCrossChainEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ILockProxyCrossChainEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ILockProxyCrossChainEvent represents a CrossChainEvent event raised by the ILockProxy contract.
type ILockProxyCrossChainEvent struct {
	Sender               common.Address
	TxId                 []byte
	ProxyOrAssetContract common.Address
	ToChainId            uint64
	ToContract           []byte
	Method               string
	Rawdata              []byte
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterCrossChainEvent is a free log retrieval operation binding the contract event 0xa50a39568d86d718c11b0750147606498d4b06fbabc1bf01b0fca5c426f5deb9.
//
// Solidity: event CrossChainEvent(address indexed sender, bytes txId, address proxyOrAssetContract, uint64 toChainId, bytes toContract, string method, bytes rawdata)
func (_ILockProxy *ILockProxyFilterer) FilterCrossChainEvent(opts *bind.FilterOpts, sender []common.Address) (*ILockProxyCrossChainEventIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ILockProxy.contract.FilterLogs(opts, "CrossChainEvent", senderRule)
	if err != nil {
		return nil, err
	}
	return &ILockProxyCrossChainEventIterator{contract: _ILockProxy.contract, event: "CrossChainEvent", logs: logs, sub: sub}, nil
}

// WatchCrossChainEvent is a free log subscription operation binding the contract event 0xa50a39568d86d718c11b0750147606498d4b06fbabc1bf01b0fca5c426f5deb9.
//
// Solidity: event CrossChainEvent(address indexed sender, bytes txId, address proxyOrAssetContract, uint64 toChainId, bytes toContract, string method, bytes rawdata)
func (_ILockProxy *ILockProxyFilterer) WatchCrossChainEvent(opts *bind.WatchOpts, sink chan<- *ILockProxyCrossChainEvent, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ILockProxy.contract.WatchLogs(opts, "CrossChainEvent", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ILockProxyCrossChainEvent)
				if err := _ILockProxy.contract.UnpackLog(event, "CrossChainEvent", log); err != nil {
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

// ParseCrossChainEvent is a log parse operation binding the contract event 0xa50a39568d86d718c11b0750147606498d4b06fbabc1bf01b0fca5c426f5deb9.
//
// Solidity: event CrossChainEvent(address indexed sender, bytes txId, address proxyOrAssetContract, uint64 toChainId, bytes toContract, string method, bytes rawdata)
func (_ILockProxy *ILockProxyFilterer) ParseCrossChainEvent(log types.Log) (*ILockProxyCrossChainEvent, error) {
	event := new(ILockProxyCrossChainEvent)
	if err := _ILockProxy.contract.UnpackLog(event, "CrossChainEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ILockProxyLockEventIterator is returned from FilterLockEvent and is used to iterate over the raw logs and unpacked data for LockEvent events raised by the ILockProxy contract.
type ILockProxyLockEventIterator struct {
	Event *ILockProxyLockEvent // Event containing the contract specifics and raw log

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
func (it *ILockProxyLockEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ILockProxyLockEvent)
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
		it.Event = new(ILockProxyLockEvent)
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
func (it *ILockProxyLockEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ILockProxyLockEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ILockProxyLockEvent represents a LockEvent event raised by the ILockProxy contract.
type ILockProxyLockEvent struct {
	FromAssetHash common.Address
	FromAddress   common.Address
	ToChainId     uint64
	ToAssetHash   []byte
	ToAddress     []byte
	Amount        *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterLockEvent is a free log retrieval operation binding the contract event 0x8636abd6d0e464fe725a13346c7ac779b73561c705506044a2e6b2cdb1295ea5.
//
// Solidity: event LockEvent(address fromAssetHash, address fromAddress, uint64 toChainId, bytes toAssetHash, bytes toAddress, uint256 amount)
func (_ILockProxy *ILockProxyFilterer) FilterLockEvent(opts *bind.FilterOpts) (*ILockProxyLockEventIterator, error) {

	logs, sub, err := _ILockProxy.contract.FilterLogs(opts, "LockEvent")
	if err != nil {
		return nil, err
	}
	return &ILockProxyLockEventIterator{contract: _ILockProxy.contract, event: "LockEvent", logs: logs, sub: sub}, nil
}

// WatchLockEvent is a free log subscription operation binding the contract event 0x8636abd6d0e464fe725a13346c7ac779b73561c705506044a2e6b2cdb1295ea5.
//
// Solidity: event LockEvent(address fromAssetHash, address fromAddress, uint64 toChainId, bytes toAssetHash, bytes toAddress, uint256 amount)
func (_ILockProxy *ILockProxyFilterer) WatchLockEvent(opts *bind.WatchOpts, sink chan<- *ILockProxyLockEvent) (event.Subscription, error) {

	logs, sub, err := _ILockProxy.contract.WatchLogs(opts, "LockEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ILockProxyLockEvent)
				if err := _ILockProxy.contract.UnpackLog(event, "LockEvent", log); err != nil {
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

// ParseLockEvent is a log parse operation binding the contract event 0x8636abd6d0e464fe725a13346c7ac779b73561c705506044a2e6b2cdb1295ea5.
//
// Solidity: event LockEvent(address fromAssetHash, address fromAddress, uint64 toChainId, bytes toAssetHash, bytes toAddress, uint256 amount)
func (_ILockProxy *ILockProxyFilterer) ParseLockEvent(log types.Log) (*ILockProxyLockEvent, error) {
	event := new(ILockProxyLockEvent)
	if err := _ILockProxy.contract.UnpackLog(event, "LockEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ILockProxyUnlockEventIterator is returned from FilterUnlockEvent and is used to iterate over the raw logs and unpacked data for UnlockEvent events raised by the ILockProxy contract.
type ILockProxyUnlockEventIterator struct {
	Event *ILockProxyUnlockEvent // Event containing the contract specifics and raw log

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
func (it *ILockProxyUnlockEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ILockProxyUnlockEvent)
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
		it.Event = new(ILockProxyUnlockEvent)
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
func (it *ILockProxyUnlockEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ILockProxyUnlockEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ILockProxyUnlockEvent represents a UnlockEvent event raised by the ILockProxy contract.
type ILockProxyUnlockEvent struct {
	ToAssetHash common.Address
	ToAddress   common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUnlockEvent is a free log retrieval operation binding the contract event 0xd90288730b87c2b8e0c45bd82260fd22478aba30ae1c4d578b8daba9261604df.
//
// Solidity: event UnlockEvent(address toAssetHash, address toAddress, uint256 amount)
func (_ILockProxy *ILockProxyFilterer) FilterUnlockEvent(opts *bind.FilterOpts) (*ILockProxyUnlockEventIterator, error) {

	logs, sub, err := _ILockProxy.contract.FilterLogs(opts, "UnlockEvent")
	if err != nil {
		return nil, err
	}
	return &ILockProxyUnlockEventIterator{contract: _ILockProxy.contract, event: "UnlockEvent", logs: logs, sub: sub}, nil
}

// WatchUnlockEvent is a free log subscription operation binding the contract event 0xd90288730b87c2b8e0c45bd82260fd22478aba30ae1c4d578b8daba9261604df.
//
// Solidity: event UnlockEvent(address toAssetHash, address toAddress, uint256 amount)
func (_ILockProxy *ILockProxyFilterer) WatchUnlockEvent(opts *bind.WatchOpts, sink chan<- *ILockProxyUnlockEvent) (event.Subscription, error) {

	logs, sub, err := _ILockProxy.contract.WatchLogs(opts, "UnlockEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ILockProxyUnlockEvent)
				if err := _ILockProxy.contract.UnpackLog(event, "UnlockEvent", log); err != nil {
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

// ParseUnlockEvent is a log parse operation binding the contract event 0xd90288730b87c2b8e0c45bd82260fd22478aba30ae1c4d578b8daba9261604df.
//
// Solidity: event UnlockEvent(address toAssetHash, address toAddress, uint256 amount)
func (_ILockProxy *ILockProxyFilterer) ParseUnlockEvent(log types.Log) (*ILockProxyUnlockEvent, error) {
	event := new(ILockProxyUnlockEvent)
	if err := _ILockProxy.contract.UnpackLog(event, "UnlockEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ILockProxyVerifyHeaderAndExecuteTxEventIterator is returned from FilterVerifyHeaderAndExecuteTxEvent and is used to iterate over the raw logs and unpacked data for VerifyHeaderAndExecuteTxEvent events raised by the ILockProxy contract.
type ILockProxyVerifyHeaderAndExecuteTxEventIterator struct {
	Event *ILockProxyVerifyHeaderAndExecuteTxEvent // Event containing the contract specifics and raw log

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
func (it *ILockProxyVerifyHeaderAndExecuteTxEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ILockProxyVerifyHeaderAndExecuteTxEvent)
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
		it.Event = new(ILockProxyVerifyHeaderAndExecuteTxEvent)
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
func (it *ILockProxyVerifyHeaderAndExecuteTxEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ILockProxyVerifyHeaderAndExecuteTxEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ILockProxyVerifyHeaderAndExecuteTxEvent represents a VerifyHeaderAndExecuteTxEvent event raised by the ILockProxy contract.
type ILockProxyVerifyHeaderAndExecuteTxEvent struct {
	FromChainID      uint64
	ToContract       []byte
	CrossChainTxHash []byte
	FromChainTxHash  []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterVerifyHeaderAndExecuteTxEvent is a free log retrieval operation binding the contract event 0x8a4a2663ce60ce4955c595da2894de0415240f1ace024cfbff85f513b656bdae.
//
// Solidity: event VerifyHeaderAndExecuteTxEvent(uint64 fromChainID, bytes toContract, bytes crossChainTxHash, bytes fromChainTxHash)
func (_ILockProxy *ILockProxyFilterer) FilterVerifyHeaderAndExecuteTxEvent(opts *bind.FilterOpts) (*ILockProxyVerifyHeaderAndExecuteTxEventIterator, error) {

	logs, sub, err := _ILockProxy.contract.FilterLogs(opts, "VerifyHeaderAndExecuteTxEvent")
	if err != nil {
		return nil, err
	}
	return &ILockProxyVerifyHeaderAndExecuteTxEventIterator{contract: _ILockProxy.contract, event: "VerifyHeaderAndExecuteTxEvent", logs: logs, sub: sub}, nil
}

// WatchVerifyHeaderAndExecuteTxEvent is a free log subscription operation binding the contract event 0x8a4a2663ce60ce4955c595da2894de0415240f1ace024cfbff85f513b656bdae.
//
// Solidity: event VerifyHeaderAndExecuteTxEvent(uint64 fromChainID, bytes toContract, bytes crossChainTxHash, bytes fromChainTxHash)
func (_ILockProxy *ILockProxyFilterer) WatchVerifyHeaderAndExecuteTxEvent(opts *bind.WatchOpts, sink chan<- *ILockProxyVerifyHeaderAndExecuteTxEvent) (event.Subscription, error) {

	logs, sub, err := _ILockProxy.contract.WatchLogs(opts, "VerifyHeaderAndExecuteTxEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ILockProxyVerifyHeaderAndExecuteTxEvent)
				if err := _ILockProxy.contract.UnpackLog(event, "VerifyHeaderAndExecuteTxEvent", log); err != nil {
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

// ParseVerifyHeaderAndExecuteTxEvent is a log parse operation binding the contract event 0x8a4a2663ce60ce4955c595da2894de0415240f1ace024cfbff85f513b656bdae.
//
// Solidity: event VerifyHeaderAndExecuteTxEvent(uint64 fromChainID, bytes toContract, bytes crossChainTxHash, bytes fromChainTxHash)
func (_ILockProxy *ILockProxyFilterer) ParseVerifyHeaderAndExecuteTxEvent(log types.Log) (*ILockProxyVerifyHeaderAndExecuteTxEvent, error) {
	event := new(ILockProxyVerifyHeaderAndExecuteTxEvent)
	if err := _ILockProxy.contract.UnpackLog(event, "VerifyHeaderAndExecuteTxEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

