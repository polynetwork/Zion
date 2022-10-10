// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package economic

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
	MethodName = "name"

	MethodReward = "reward"

	MethodTotalSupply = "totalSupply"
)

// IEconomicABI is the input ABI used to generate the binding from.
const IEconomicABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"reward\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// IEconomicFuncSigs maps the 4-byte function signature to its string representation.
var IEconomicFuncSigs = map[string]string{
	"06fdde03": "name()",
	"228cb733": "reward()",
	"18160ddd": "totalSupply()",
}

// IEconomic is an auto generated Go binding around an Ethereum contract.
type IEconomic struct {
	IEconomicCaller     // Read-only binding to the contract
	IEconomicTransactor // Write-only binding to the contract
	IEconomicFilterer   // Log filterer for contract events
}

// IEconomicCaller is an auto generated read-only Go binding around an Ethereum contract.
type IEconomicCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IEconomicTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IEconomicTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IEconomicFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IEconomicFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IEconomicSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IEconomicSession struct {
	Contract     *IEconomic        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IEconomicCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IEconomicCallerSession struct {
	Contract *IEconomicCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// IEconomicTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IEconomicTransactorSession struct {
	Contract     *IEconomicTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IEconomicRaw is an auto generated low-level Go binding around an Ethereum contract.
type IEconomicRaw struct {
	Contract *IEconomic // Generic contract binding to access the raw methods on
}

// IEconomicCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IEconomicCallerRaw struct {
	Contract *IEconomicCaller // Generic read-only contract binding to access the raw methods on
}

// IEconomicTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IEconomicTransactorRaw struct {
	Contract *IEconomicTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIEconomic creates a new instance of IEconomic, bound to a specific deployed contract.
func NewIEconomic(address common.Address, backend bind.ContractBackend) (*IEconomic, error) {
	contract, err := bindIEconomic(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IEconomic{IEconomicCaller: IEconomicCaller{contract: contract}, IEconomicTransactor: IEconomicTransactor{contract: contract}, IEconomicFilterer: IEconomicFilterer{contract: contract}}, nil
}

// NewIEconomicCaller creates a new read-only instance of IEconomic, bound to a specific deployed contract.
func NewIEconomicCaller(address common.Address, caller bind.ContractCaller) (*IEconomicCaller, error) {
	contract, err := bindIEconomic(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IEconomicCaller{contract: contract}, nil
}

// NewIEconomicTransactor creates a new write-only instance of IEconomic, bound to a specific deployed contract.
func NewIEconomicTransactor(address common.Address, transactor bind.ContractTransactor) (*IEconomicTransactor, error) {
	contract, err := bindIEconomic(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IEconomicTransactor{contract: contract}, nil
}

// NewIEconomicFilterer creates a new log filterer instance of IEconomic, bound to a specific deployed contract.
func NewIEconomicFilterer(address common.Address, filterer bind.ContractFilterer) (*IEconomicFilterer, error) {
	contract, err := bindIEconomic(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IEconomicFilterer{contract: contract}, nil
}

// bindIEconomic binds a generic wrapper to an already deployed contract.
func bindIEconomic(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IEconomicABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IEconomic *IEconomicRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IEconomic.Contract.IEconomicCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IEconomic *IEconomicRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IEconomic.Contract.IEconomicTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IEconomic *IEconomicRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IEconomic.Contract.IEconomicTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IEconomic *IEconomicCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IEconomic.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IEconomic *IEconomicTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IEconomic.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IEconomic *IEconomicTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IEconomic.Contract.contract.Transact(opts, method, params...)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_IEconomic *IEconomicCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IEconomic.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_IEconomic *IEconomicSession) Name() (string, error) {
	return _IEconomic.Contract.Name(&_IEconomic.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_IEconomic *IEconomicCallerSession) Name() (string, error) {
	return _IEconomic.Contract.Name(&_IEconomic.CallOpts)
}

// Reward is a free data retrieval call binding the contract method 0x228cb733.
//
// Solidity: function reward() view returns(bytes)
func (_IEconomic *IEconomicCaller) Reward(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _IEconomic.contract.Call(opts, &out, "reward")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Reward is a free data retrieval call binding the contract method 0x228cb733.
//
// Solidity: function reward() view returns(bytes)
func (_IEconomic *IEconomicSession) Reward() ([]byte, error) {
	return _IEconomic.Contract.Reward(&_IEconomic.CallOpts)
}

// Reward is a free data retrieval call binding the contract method 0x228cb733.
//
// Solidity: function reward() view returns(bytes)
func (_IEconomic *IEconomicCallerSession) Reward() ([]byte, error) {
	return _IEconomic.Contract.Reward(&_IEconomic.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IEconomic *IEconomicCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IEconomic.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IEconomic *IEconomicSession) TotalSupply() (*big.Int, error) {
	return _IEconomic.Contract.TotalSupply(&_IEconomic.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IEconomic *IEconomicCallerSession) TotalSupply() (*big.Int, error) {
	return _IEconomic.Contract.TotalSupply(&_IEconomic.CallOpts)
}
