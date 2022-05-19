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
	MethodSyncRootInfo = "syncRootInfo"

	MethodName = "name"
)

// InfoSyncABI is the input ABI used to generate the binding from.
const InfoSyncABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"syncRootInfo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// InfoSyncFuncSigs maps the 4-byte function signature to its string representation.
var InfoSyncFuncSigs = map[string]string{
	"06fdde03": "name()",
	"48c8f119": "syncRootInfo(uint64,bytes[])",
}

// InfoSync is an auto generated Go binding around an Ethereum contract.
type InfoSync struct {
	InfoSyncCaller     // Read-only binding to the contract
	InfoSyncTransactor // Write-only binding to the contract
	InfoSyncFilterer   // Log filterer for contract events
}

// InfoSyncCaller is an auto generated read-only Go binding around an Ethereum contract.
type InfoSyncCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InfoSyncTransactor is an auto generated write-only Go binding around an Ethereum contract.
type InfoSyncTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InfoSyncFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type InfoSyncFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InfoSyncSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type InfoSyncSession struct {
	Contract     *InfoSync         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// InfoSyncCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type InfoSyncCallerSession struct {
	Contract *InfoSyncCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// InfoSyncTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type InfoSyncTransactorSession struct {
	Contract     *InfoSyncTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// InfoSyncRaw is an auto generated low-level Go binding around an Ethereum contract.
type InfoSyncRaw struct {
	Contract *InfoSync // Generic contract binding to access the raw methods on
}

// InfoSyncCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type InfoSyncCallerRaw struct {
	Contract *InfoSyncCaller // Generic read-only contract binding to access the raw methods on
}

// InfoSyncTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type InfoSyncTransactorRaw struct {
	Contract *InfoSyncTransactor // Generic write-only contract binding to access the raw methods on
}

// NewInfoSync creates a new instance of InfoSync, bound to a specific deployed contract.
func NewInfoSync(address common.Address, backend bind.ContractBackend) (*InfoSync, error) {
	contract, err := bindInfoSync(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &InfoSync{InfoSyncCaller: InfoSyncCaller{contract: contract}, InfoSyncTransactor: InfoSyncTransactor{contract: contract}, InfoSyncFilterer: InfoSyncFilterer{contract: contract}}, nil
}

// NewInfoSyncCaller creates a new read-only instance of InfoSync, bound to a specific deployed contract.
func NewInfoSyncCaller(address common.Address, caller bind.ContractCaller) (*InfoSyncCaller, error) {
	contract, err := bindInfoSync(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &InfoSyncCaller{contract: contract}, nil
}

// NewInfoSyncTransactor creates a new write-only instance of InfoSync, bound to a specific deployed contract.
func NewInfoSyncTransactor(address common.Address, transactor bind.ContractTransactor) (*InfoSyncTransactor, error) {
	contract, err := bindInfoSync(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &InfoSyncTransactor{contract: contract}, nil
}

// NewInfoSyncFilterer creates a new log filterer instance of InfoSync, bound to a specific deployed contract.
func NewInfoSyncFilterer(address common.Address, filterer bind.ContractFilterer) (*InfoSyncFilterer, error) {
	contract, err := bindInfoSync(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &InfoSyncFilterer{contract: contract}, nil
}

// bindInfoSync binds a generic wrapper to an already deployed contract.
func bindInfoSync(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(InfoSyncABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InfoSync *InfoSyncRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _InfoSync.Contract.InfoSyncCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InfoSync *InfoSyncRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InfoSync.Contract.InfoSyncTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InfoSync *InfoSyncRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InfoSync.Contract.InfoSyncTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InfoSync *InfoSyncCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _InfoSync.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InfoSync *InfoSyncTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InfoSync.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InfoSync *InfoSyncTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InfoSync.Contract.contract.Transact(opts, method, params...)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_InfoSync *InfoSyncCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _InfoSync.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_InfoSync *InfoSyncSession) Name() (string, error) {
	return _InfoSync.Contract.Name(&_InfoSync.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_InfoSync *InfoSyncCallerSession) Name() (string, error) {
	return _InfoSync.Contract.Name(&_InfoSync.CallOpts)
}

// SyncRootInfo is a paid mutator transaction binding the contract method 0x48c8f119.
//
// Solidity: function syncRootInfo(uint64 chainID, bytes[] data) returns()
func (_InfoSync *InfoSyncTransactor) SyncRootInfo(opts *bind.TransactOpts, chainID uint64, data [][]byte) (*types.Transaction, error) {
	return _InfoSync.contract.Transact(opts, "syncRootInfo", chainID, data)
}

// SyncRootInfo is a paid mutator transaction binding the contract method 0x48c8f119.
//
// Solidity: function syncRootInfo(uint64 chainID, bytes[] data) returns()
func (_InfoSync *InfoSyncSession) SyncRootInfo(chainID uint64, data [][]byte) (*types.Transaction, error) {
	return _InfoSync.Contract.SyncRootInfo(&_InfoSync.TransactOpts, chainID, data)
}

// SyncRootInfo is a paid mutator transaction binding the contract method 0x48c8f119.
//
// Solidity: function syncRootInfo(uint64 chainID, bytes[] data) returns()
func (_InfoSync *InfoSyncTransactorSession) SyncRootInfo(chainID uint64, data [][]byte) (*types.Transaction, error) {
	return _InfoSync.Contract.SyncRootInfo(&_InfoSync.TransactOpts, chainID, data)
}

