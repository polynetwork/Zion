// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package node_manager_abi

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
	MethodBeginBlock = "beginBlock"

	MethodCancelValidator = "cancelValidator"

	MethodChangeEpoch = "changeEpoch"

	MethodCreateValidator = "createValidator"

	MethodStake = "stake"

	MethodUnStake = "unStake"

	MethodUpdateValidator = "updateValidator"

	MethodWithdraw = "withdraw"

	MethodWithdrawCommission = "withdrawCommission"

	MethodWithdrawStakeRewards = "withdrawStakeRewards"

	MethodWithdrawValidator = "withdrawValidator"
)

// INodeManagerABI is the input ABI used to generate the binding from.
const INodeManagerABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"beginBlock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"consensusPubkey\",\"type\":\"string\"}],\"name\":\"cancelValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"changeEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"consensusPubkey\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"proposalAddress\",\"type\":\"address\"},{\"internalType\":\"int256\",\"name\":\"commission\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"initStake\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"desc\",\"type\":\"string\"}],\"name\":\"createValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"consensusPubkey\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"amount\",\"type\":\"int256\"}],\"name\":\"stake\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"consensusPubkey\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"amount\",\"type\":\"int256\"}],\"name\":\"unStake\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"consensusPubkey\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"proposalAddress\",\"type\":\"address\"},{\"internalType\":\"int256\",\"name\":\"commission\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"desc\",\"type\":\"string\"}],\"name\":\"updateValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"consensusPubkey\",\"type\":\"string\"}],\"name\":\"withdrawCommission\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"consensusPubkey\",\"type\":\"string\"}],\"name\":\"withdrawStakeRewards\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"consensusPubkey\",\"type\":\"string\"}],\"name\":\"withdrawValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// INodeManagerFuncSigs maps the 4-byte function signature to its string representation.
var INodeManagerFuncSigs = map[string]string{
	"a6903278": "beginBlock()",
	"d087b01b": "cancelValidator(string)",
	"fe6f86f8": "changeEpoch()",
	"afa1d242": "createValidator(string,address,int256,int256,string)",
	"102fc25a": "stake(string,int256)",
	"6b2e7c4e": "unStake(string,int256)",
	"b0694322": "updateValidator(string,address,int256,string)",
	"3ccfd60b": "withdraw()",
	"65371818": "withdrawCommission(string)",
	"377b97fb": "withdrawStakeRewards(string)",
	"944ac99f": "withdrawValidator(string)",
}

// INodeManager is an auto generated Go binding around an Ethereum contract.
type INodeManager struct {
	INodeManagerCaller     // Read-only binding to the contract
	INodeManagerTransactor // Write-only binding to the contract
	INodeManagerFilterer   // Log filterer for contract events
}

// INodeManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type INodeManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// INodeManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type INodeManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// INodeManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type INodeManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// INodeManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type INodeManagerSession struct {
	Contract     *INodeManager     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// INodeManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type INodeManagerCallerSession struct {
	Contract *INodeManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// INodeManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type INodeManagerTransactorSession struct {
	Contract     *INodeManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// INodeManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type INodeManagerRaw struct {
	Contract *INodeManager // Generic contract binding to access the raw methods on
}

// INodeManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type INodeManagerCallerRaw struct {
	Contract *INodeManagerCaller // Generic read-only contract binding to access the raw methods on
}

// INodeManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type INodeManagerTransactorRaw struct {
	Contract *INodeManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewINodeManager creates a new instance of INodeManager, bound to a specific deployed contract.
func NewINodeManager(address common.Address, backend bind.ContractBackend) (*INodeManager, error) {
	contract, err := bindINodeManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &INodeManager{INodeManagerCaller: INodeManagerCaller{contract: contract}, INodeManagerTransactor: INodeManagerTransactor{contract: contract}, INodeManagerFilterer: INodeManagerFilterer{contract: contract}}, nil
}

// NewINodeManagerCaller creates a new read-only instance of INodeManager, bound to a specific deployed contract.
func NewINodeManagerCaller(address common.Address, caller bind.ContractCaller) (*INodeManagerCaller, error) {
	contract, err := bindINodeManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &INodeManagerCaller{contract: contract}, nil
}

// NewINodeManagerTransactor creates a new write-only instance of INodeManager, bound to a specific deployed contract.
func NewINodeManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*INodeManagerTransactor, error) {
	contract, err := bindINodeManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &INodeManagerTransactor{contract: contract}, nil
}

// NewINodeManagerFilterer creates a new log filterer instance of INodeManager, bound to a specific deployed contract.
func NewINodeManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*INodeManagerFilterer, error) {
	contract, err := bindINodeManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &INodeManagerFilterer{contract: contract}, nil
}

// bindINodeManager binds a generic wrapper to an already deployed contract.
func bindINodeManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(INodeManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_INodeManager *INodeManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _INodeManager.Contract.INodeManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_INodeManager *INodeManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _INodeManager.Contract.INodeManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_INodeManager *INodeManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _INodeManager.Contract.INodeManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_INodeManager *INodeManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _INodeManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_INodeManager *INodeManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _INodeManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_INodeManager *INodeManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _INodeManager.Contract.contract.Transact(opts, method, params...)
}

// BeginBlock is a paid mutator transaction binding the contract method 0xa6903278.
//
// Solidity: function beginBlock() returns(bool success)
func (_INodeManager *INodeManagerTransactor) BeginBlock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "beginBlock")
}

// BeginBlock is a paid mutator transaction binding the contract method 0xa6903278.
//
// Solidity: function beginBlock() returns(bool success)
func (_INodeManager *INodeManagerSession) BeginBlock() (*types.Transaction, error) {
	return _INodeManager.Contract.BeginBlock(&_INodeManager.TransactOpts)
}

// BeginBlock is a paid mutator transaction binding the contract method 0xa6903278.
//
// Solidity: function beginBlock() returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) BeginBlock() (*types.Transaction, error) {
	return _INodeManager.Contract.BeginBlock(&_INodeManager.TransactOpts)
}

// CancelValidator is a paid mutator transaction binding the contract method 0xd087b01b.
//
// Solidity: function cancelValidator(string consensusPubkey) returns(bool success)
func (_INodeManager *INodeManagerTransactor) CancelValidator(opts *bind.TransactOpts, consensusPubkey string) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "cancelValidator", consensusPubkey)
}

// CancelValidator is a paid mutator transaction binding the contract method 0xd087b01b.
//
// Solidity: function cancelValidator(string consensusPubkey) returns(bool success)
func (_INodeManager *INodeManagerSession) CancelValidator(consensusPubkey string) (*types.Transaction, error) {
	return _INodeManager.Contract.CancelValidator(&_INodeManager.TransactOpts, consensusPubkey)
}

// CancelValidator is a paid mutator transaction binding the contract method 0xd087b01b.
//
// Solidity: function cancelValidator(string consensusPubkey) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) CancelValidator(consensusPubkey string) (*types.Transaction, error) {
	return _INodeManager.Contract.CancelValidator(&_INodeManager.TransactOpts, consensusPubkey)
}

// ChangeEpoch is a paid mutator transaction binding the contract method 0xfe6f86f8.
//
// Solidity: function changeEpoch() returns(bool success)
func (_INodeManager *INodeManagerTransactor) ChangeEpoch(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "changeEpoch")
}

// ChangeEpoch is a paid mutator transaction binding the contract method 0xfe6f86f8.
//
// Solidity: function changeEpoch() returns(bool success)
func (_INodeManager *INodeManagerSession) ChangeEpoch() (*types.Transaction, error) {
	return _INodeManager.Contract.ChangeEpoch(&_INodeManager.TransactOpts)
}

// ChangeEpoch is a paid mutator transaction binding the contract method 0xfe6f86f8.
//
// Solidity: function changeEpoch() returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) ChangeEpoch() (*types.Transaction, error) {
	return _INodeManager.Contract.ChangeEpoch(&_INodeManager.TransactOpts)
}

// CreateValidator is a paid mutator transaction binding the contract method 0xafa1d242.
//
// Solidity: function createValidator(string consensusPubkey, address proposalAddress, int256 commission, int256 initStake, string desc) returns(bool success)
func (_INodeManager *INodeManagerTransactor) CreateValidator(opts *bind.TransactOpts, consensusPubkey string, proposalAddress common.Address, commission *big.Int, initStake *big.Int, desc string) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "createValidator", consensusPubkey, proposalAddress, commission, initStake, desc)
}

// CreateValidator is a paid mutator transaction binding the contract method 0xafa1d242.
//
// Solidity: function createValidator(string consensusPubkey, address proposalAddress, int256 commission, int256 initStake, string desc) returns(bool success)
func (_INodeManager *INodeManagerSession) CreateValidator(consensusPubkey string, proposalAddress common.Address, commission *big.Int, initStake *big.Int, desc string) (*types.Transaction, error) {
	return _INodeManager.Contract.CreateValidator(&_INodeManager.TransactOpts, consensusPubkey, proposalAddress, commission, initStake, desc)
}

// CreateValidator is a paid mutator transaction binding the contract method 0xafa1d242.
//
// Solidity: function createValidator(string consensusPubkey, address proposalAddress, int256 commission, int256 initStake, string desc) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) CreateValidator(consensusPubkey string, proposalAddress common.Address, commission *big.Int, initStake *big.Int, desc string) (*types.Transaction, error) {
	return _INodeManager.Contract.CreateValidator(&_INodeManager.TransactOpts, consensusPubkey, proposalAddress, commission, initStake, desc)
}

// Stake is a paid mutator transaction binding the contract method 0x102fc25a.
//
// Solidity: function stake(string consensusPubkey, int256 amount) returns(bool success)
func (_INodeManager *INodeManagerTransactor) Stake(opts *bind.TransactOpts, consensusPubkey string, amount *big.Int) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "stake", consensusPubkey, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x102fc25a.
//
// Solidity: function stake(string consensusPubkey, int256 amount) returns(bool success)
func (_INodeManager *INodeManagerSession) Stake(consensusPubkey string, amount *big.Int) (*types.Transaction, error) {
	return _INodeManager.Contract.Stake(&_INodeManager.TransactOpts, consensusPubkey, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x102fc25a.
//
// Solidity: function stake(string consensusPubkey, int256 amount) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) Stake(consensusPubkey string, amount *big.Int) (*types.Transaction, error) {
	return _INodeManager.Contract.Stake(&_INodeManager.TransactOpts, consensusPubkey, amount)
}

// UnStake is a paid mutator transaction binding the contract method 0x6b2e7c4e.
//
// Solidity: function unStake(string consensusPubkey, int256 amount) returns(bool success)
func (_INodeManager *INodeManagerTransactor) UnStake(opts *bind.TransactOpts, consensusPubkey string, amount *big.Int) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "unStake", consensusPubkey, amount)
}

// UnStake is a paid mutator transaction binding the contract method 0x6b2e7c4e.
//
// Solidity: function unStake(string consensusPubkey, int256 amount) returns(bool success)
func (_INodeManager *INodeManagerSession) UnStake(consensusPubkey string, amount *big.Int) (*types.Transaction, error) {
	return _INodeManager.Contract.UnStake(&_INodeManager.TransactOpts, consensusPubkey, amount)
}

// UnStake is a paid mutator transaction binding the contract method 0x6b2e7c4e.
//
// Solidity: function unStake(string consensusPubkey, int256 amount) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) UnStake(consensusPubkey string, amount *big.Int) (*types.Transaction, error) {
	return _INodeManager.Contract.UnStake(&_INodeManager.TransactOpts, consensusPubkey, amount)
}

// UpdateValidator is a paid mutator transaction binding the contract method 0xb0694322.
//
// Solidity: function updateValidator(string consensusPubkey, address proposalAddress, int256 commission, string desc) returns(bool success)
func (_INodeManager *INodeManagerTransactor) UpdateValidator(opts *bind.TransactOpts, consensusPubkey string, proposalAddress common.Address, commission *big.Int, desc string) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "updateValidator", consensusPubkey, proposalAddress, commission, desc)
}

// UpdateValidator is a paid mutator transaction binding the contract method 0xb0694322.
//
// Solidity: function updateValidator(string consensusPubkey, address proposalAddress, int256 commission, string desc) returns(bool success)
func (_INodeManager *INodeManagerSession) UpdateValidator(consensusPubkey string, proposalAddress common.Address, commission *big.Int, desc string) (*types.Transaction, error) {
	return _INodeManager.Contract.UpdateValidator(&_INodeManager.TransactOpts, consensusPubkey, proposalAddress, commission, desc)
}

// UpdateValidator is a paid mutator transaction binding the contract method 0xb0694322.
//
// Solidity: function updateValidator(string consensusPubkey, address proposalAddress, int256 commission, string desc) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) UpdateValidator(consensusPubkey string, proposalAddress common.Address, commission *big.Int, desc string) (*types.Transaction, error) {
	return _INodeManager.Contract.UpdateValidator(&_INodeManager.TransactOpts, consensusPubkey, proposalAddress, commission, desc)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns(bool success)
func (_INodeManager *INodeManagerTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns(bool success)
func (_INodeManager *INodeManagerSession) Withdraw() (*types.Transaction, error) {
	return _INodeManager.Contract.Withdraw(&_INodeManager.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) Withdraw() (*types.Transaction, error) {
	return _INodeManager.Contract.Withdraw(&_INodeManager.TransactOpts)
}

// WithdrawCommission is a paid mutator transaction binding the contract method 0x65371818.
//
// Solidity: function withdrawCommission(string consensusPubkey) returns(bool success)
func (_INodeManager *INodeManagerTransactor) WithdrawCommission(opts *bind.TransactOpts, consensusPubkey string) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "withdrawCommission", consensusPubkey)
}

// WithdrawCommission is a paid mutator transaction binding the contract method 0x65371818.
//
// Solidity: function withdrawCommission(string consensusPubkey) returns(bool success)
func (_INodeManager *INodeManagerSession) WithdrawCommission(consensusPubkey string) (*types.Transaction, error) {
	return _INodeManager.Contract.WithdrawCommission(&_INodeManager.TransactOpts, consensusPubkey)
}

// WithdrawCommission is a paid mutator transaction binding the contract method 0x65371818.
//
// Solidity: function withdrawCommission(string consensusPubkey) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) WithdrawCommission(consensusPubkey string) (*types.Transaction, error) {
	return _INodeManager.Contract.WithdrawCommission(&_INodeManager.TransactOpts, consensusPubkey)
}

// WithdrawStakeRewards is a paid mutator transaction binding the contract method 0x377b97fb.
//
// Solidity: function withdrawStakeRewards(string consensusPubkey) returns(bool success)
func (_INodeManager *INodeManagerTransactor) WithdrawStakeRewards(opts *bind.TransactOpts, consensusPubkey string) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "withdrawStakeRewards", consensusPubkey)
}

// WithdrawStakeRewards is a paid mutator transaction binding the contract method 0x377b97fb.
//
// Solidity: function withdrawStakeRewards(string consensusPubkey) returns(bool success)
func (_INodeManager *INodeManagerSession) WithdrawStakeRewards(consensusPubkey string) (*types.Transaction, error) {
	return _INodeManager.Contract.WithdrawStakeRewards(&_INodeManager.TransactOpts, consensusPubkey)
}

// WithdrawStakeRewards is a paid mutator transaction binding the contract method 0x377b97fb.
//
// Solidity: function withdrawStakeRewards(string consensusPubkey) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) WithdrawStakeRewards(consensusPubkey string) (*types.Transaction, error) {
	return _INodeManager.Contract.WithdrawStakeRewards(&_INodeManager.TransactOpts, consensusPubkey)
}

// WithdrawValidator is a paid mutator transaction binding the contract method 0x944ac99f.
//
// Solidity: function withdrawValidator(string consensusPubkey) returns(bool success)
func (_INodeManager *INodeManagerTransactor) WithdrawValidator(opts *bind.TransactOpts, consensusPubkey string) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "withdrawValidator", consensusPubkey)
}

// WithdrawValidator is a paid mutator transaction binding the contract method 0x944ac99f.
//
// Solidity: function withdrawValidator(string consensusPubkey) returns(bool success)
func (_INodeManager *INodeManagerSession) WithdrawValidator(consensusPubkey string) (*types.Transaction, error) {
	return _INodeManager.Contract.WithdrawValidator(&_INodeManager.TransactOpts, consensusPubkey)
}

// WithdrawValidator is a paid mutator transaction binding the contract method 0x944ac99f.
//
// Solidity: function withdrawValidator(string consensusPubkey) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) WithdrawValidator(consensusPubkey string) (*types.Transaction, error) {
	return _INodeManager.Contract.WithdrawValidator(&_INodeManager.TransactOpts, consensusPubkey)
}

