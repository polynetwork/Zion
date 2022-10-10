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
	MethodCancelValidator = "cancelValidator"

	MethodChangeEpoch = "changeEpoch"

	MethodCreateValidator = "createValidator"

	MethodEndBlock = "endBlock"

	MethodStake = "stake"

	MethodUnStake = "unStake"

	MethodUpdateCommission = "updateCommission"

	MethodUpdateValidator = "updateValidator"

	MethodWithdraw = "withdraw"

	MethodWithdrawCommission = "withdrawCommission"

	MethodWithdrawStakeRewards = "withdrawStakeRewards"

	MethodWithdrawValidator = "withdrawValidator"

	MethodGetAccumulatedCommission = "getAccumulatedCommission"

	MethodGetAllValidators = "getAllValidators"

	MethodGetCommunityInfo = "getCommunityInfo"

	MethodGetCurrentEpochInfo = "getCurrentEpochInfo"

	MethodGetEpochInfo = "getEpochInfo"

	MethodGetGlobalConfig = "getGlobalConfig"

	MethodGetOutstandingRewards = "getOutstandingRewards"

	MethodGetStakeInfo = "getStakeInfo"

	MethodGetStakeRewards = "getStakeRewards"

	MethodGetStakeStartingInfo = "getStakeStartingInfo"

	MethodGetTotalPool = "getTotalPool"

	MethodGetUnlockingInfo = "getUnlockingInfo"

	MethodGetValidator = "getValidator"

	MethodGetValidatorAccumulatedRewards = "getValidatorAccumulatedRewards"

	MethodGetValidatorOutstandingRewards = "getValidatorOutstandingRewards"

	MethodGetValidatorSnapshotRewards = "getValidatorSnapshotRewards"

	EventCancelValidator = "CancelValidator"

	EventChangeEpoch = "ChangeEpoch"

	EventCreateValidator = "CreateValidator"

	EventStake = "Stake"

	EventUnStake = "UnStake"

	EventUpdateCommission = "UpdateCommission"

	EventUpdateValidator = "UpdateValidator"

	EventWithdraw = "Withdraw"

	EventWithdrawCommission = "WithdrawCommission"

	EventWithdrawStakeRewards = "WithdrawStakeRewards"

	EventWithdrawValidator = "WithdrawValidator"
)

// INodeManagerABI is the input ABI used to generate the binding from.
const INodeManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"consensusAddress\",\"type\":\"string\"}],\"name\":\"CancelValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"epochID\",\"type\":\"string\"}],\"name\":\"ChangeEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"consensusAddress\",\"type\":\"string\"}],\"name\":\"CreateValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"consensusAddress\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"amount\",\"type\":\"string\"}],\"name\":\"Stake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"consensusAddress\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"amount\",\"type\":\"string\"}],\"name\":\"UnStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"consensusAddress\",\"type\":\"string\"}],\"name\":\"UpdateCommission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"consensusAddress\",\"type\":\"string\"}],\"name\":\"UpdateValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"caller\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"amount\",\"type\":\"string\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"consensusAddress\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"commission\",\"type\":\"string\"}],\"name\":\"WithdrawCommission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"consensusAddress\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"caller\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"rewards\",\"type\":\"string\"}],\"name\":\"WithdrawStakeRewards\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"consensusAddress\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"selfStake\",\"type\":\"string\"}],\"name\":\"WithdrawValidator\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"}],\"name\":\"cancelValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"changeEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposalAddress\",\"type\":\"address\"},{\"internalType\":\"int256\",\"name\":\"commission\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"initStake\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"desc\",\"type\":\"string\"}],\"name\":\"createValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"endBlock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"}],\"name\":\"getAccumulatedCommission\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllValidators\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommunityInfo\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentEpochInfo\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"id\",\"type\":\"int256\"}],\"name\":\"getEpochInfo\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGlobalConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOutstandingRewards\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stakeAddress\",\"type\":\"address\"}],\"name\":\"getStakeInfo\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stakeAddress\",\"type\":\"address\"}],\"name\":\"getStakeRewards\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stakeAddress\",\"type\":\"address\"}],\"name\":\"getStakeStartingInfo\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalPool\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"stakeAddress\",\"type\":\"address\"}],\"name\":\"getUnlockingInfo\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"}],\"name\":\"getValidator\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"}],\"name\":\"getValidatorAccumulatedRewards\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"}],\"name\":\"getValidatorOutstandingRewards\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"period\",\"type\":\"uint64\"}],\"name\":\"getValidatorSnapshotRewards\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"},{\"internalType\":\"int256\",\"name\":\"amount\",\"type\":\"int256\"}],\"name\":\"stake\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"},{\"internalType\":\"int256\",\"name\":\"amount\",\"type\":\"int256\"}],\"name\":\"unStake\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"},{\"internalType\":\"int256\",\"name\":\"commission\",\"type\":\"int256\"}],\"name\":\"updateCommission\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposalAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"desc\",\"type\":\"string\"}],\"name\":\"updateValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"}],\"name\":\"withdrawCommission\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"}],\"name\":\"withdrawStakeRewards\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusAddress\",\"type\":\"address\"}],\"name\":\"withdrawValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// INodeManagerFuncSigs maps the 4-byte function signature to its string representation.
var INodeManagerFuncSigs = map[string]string{
	"1af78584": "cancelValidator(address)",
	"fe6f86f8": "changeEpoch()",
	"b5c42005": "createValidator(address,address,address,int256,int256,string)",
	"083c6323": "endBlock()",
	"21d38c78": "getAccumulatedCommission(address)",
	"f3513a37": "getAllValidators()",
	"6e10ffd0": "getCommunityInfo()",
	"babc394f": "getCurrentEpochInfo()",
	"1af10a9c": "getEpochInfo(int256)",
	"cda92be4": "getGlobalConfig()",
	"fef97e4c": "getOutstandingRewards()",
	"d77c8f14": "getStakeInfo(address,address)",
	"ea3f32ff": "getStakeRewards(address,address)",
	"17674715": "getStakeStartingInfo(address,address)",
	"75f4d677": "getTotalPool()",
	"5e45511e": "getUnlockingInfo(address)",
	"1904bb2e": "getValidator(address)",
	"9c898a3b": "getValidatorAccumulatedRewards(address)",
	"a76d00a8": "getValidatorOutstandingRewards(address)",
	"edd0efa9": "getValidatorSnapshotRewards(address,uint64)",
	"91398b13": "stake(address,int256)",
	"dfe6bad3": "unStake(address,int256)",
	"c5e7ad1d": "updateCommission(address,int256)",
	"24712218": "updateValidator(address,address,address,string)",
	"3ccfd60b": "withdraw()",
	"16c58d04": "withdrawCommission(address)",
	"5d3abd54": "withdrawStakeRewards(address)",
	"814a2c11": "withdrawValidator(address)",
}

// GetSpecMethodID for consensus use
func GetSpecMethodID() map[string]bool {
	return map[string]bool{"fe6f86f8": true, "083c6323": true}
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

// GetAccumulatedCommission is a free data retrieval call binding the contract method 0x21d38c78.
//
// Solidity: function getAccumulatedCommission(address consensusAddress) view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetAccumulatedCommission(opts *bind.CallOpts, consensusAddress common.Address) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getAccumulatedCommission", consensusAddress)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetAccumulatedCommission is a free data retrieval call binding the contract method 0x21d38c78.
//
// Solidity: function getAccumulatedCommission(address consensusAddress) view returns(bytes)
func (_INodeManager *INodeManagerSession) GetAccumulatedCommission(consensusAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetAccumulatedCommission(&_INodeManager.CallOpts, consensusAddress)
}

// GetAccumulatedCommission is a free data retrieval call binding the contract method 0x21d38c78.
//
// Solidity: function getAccumulatedCommission(address consensusAddress) view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetAccumulatedCommission(consensusAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetAccumulatedCommission(&_INodeManager.CallOpts, consensusAddress)
}

// GetAllValidators is a free data retrieval call binding the contract method 0xf3513a37.
//
// Solidity: function getAllValidators() view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetAllValidators(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getAllValidators")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetAllValidators is a free data retrieval call binding the contract method 0xf3513a37.
//
// Solidity: function getAllValidators() view returns(bytes)
func (_INodeManager *INodeManagerSession) GetAllValidators() ([]byte, error) {
	return _INodeManager.Contract.GetAllValidators(&_INodeManager.CallOpts)
}

// GetAllValidators is a free data retrieval call binding the contract method 0xf3513a37.
//
// Solidity: function getAllValidators() view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetAllValidators() ([]byte, error) {
	return _INodeManager.Contract.GetAllValidators(&_INodeManager.CallOpts)
}

// GetCommunityInfo is a free data retrieval call binding the contract method 0x6e10ffd0.
//
// Solidity: function getCommunityInfo() view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetCommunityInfo(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getCommunityInfo")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetCommunityInfo is a free data retrieval call binding the contract method 0x6e10ffd0.
//
// Solidity: function getCommunityInfo() view returns(bytes)
func (_INodeManager *INodeManagerSession) GetCommunityInfo() ([]byte, error) {
	return _INodeManager.Contract.GetCommunityInfo(&_INodeManager.CallOpts)
}

// GetCommunityInfo is a free data retrieval call binding the contract method 0x6e10ffd0.
//
// Solidity: function getCommunityInfo() view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetCommunityInfo() ([]byte, error) {
	return _INodeManager.Contract.GetCommunityInfo(&_INodeManager.CallOpts)
}

// GetCurrentEpochInfo is a free data retrieval call binding the contract method 0xbabc394f.
//
// Solidity: function getCurrentEpochInfo() view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetCurrentEpochInfo(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getCurrentEpochInfo")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetCurrentEpochInfo is a free data retrieval call binding the contract method 0xbabc394f.
//
// Solidity: function getCurrentEpochInfo() view returns(bytes)
func (_INodeManager *INodeManagerSession) GetCurrentEpochInfo() ([]byte, error) {
	return _INodeManager.Contract.GetCurrentEpochInfo(&_INodeManager.CallOpts)
}

// GetCurrentEpochInfo is a free data retrieval call binding the contract method 0xbabc394f.
//
// Solidity: function getCurrentEpochInfo() view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetCurrentEpochInfo() ([]byte, error) {
	return _INodeManager.Contract.GetCurrentEpochInfo(&_INodeManager.CallOpts)
}

// GetEpochInfo is a free data retrieval call binding the contract method 0x1af10a9c.
//
// Solidity: function getEpochInfo(int256 id) view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetEpochInfo(opts *bind.CallOpts, id *big.Int) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getEpochInfo", id)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetEpochInfo is a free data retrieval call binding the contract method 0x1af10a9c.
//
// Solidity: function getEpochInfo(int256 id) view returns(bytes)
func (_INodeManager *INodeManagerSession) GetEpochInfo(id *big.Int) ([]byte, error) {
	return _INodeManager.Contract.GetEpochInfo(&_INodeManager.CallOpts, id)
}

// GetEpochInfo is a free data retrieval call binding the contract method 0x1af10a9c.
//
// Solidity: function getEpochInfo(int256 id) view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetEpochInfo(id *big.Int) ([]byte, error) {
	return _INodeManager.Contract.GetEpochInfo(&_INodeManager.CallOpts, id)
}

// GetGlobalConfig is a free data retrieval call binding the contract method 0xcda92be4.
//
// Solidity: function getGlobalConfig() view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetGlobalConfig(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getGlobalConfig")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetGlobalConfig is a free data retrieval call binding the contract method 0xcda92be4.
//
// Solidity: function getGlobalConfig() view returns(bytes)
func (_INodeManager *INodeManagerSession) GetGlobalConfig() ([]byte, error) {
	return _INodeManager.Contract.GetGlobalConfig(&_INodeManager.CallOpts)
}

// GetGlobalConfig is a free data retrieval call binding the contract method 0xcda92be4.
//
// Solidity: function getGlobalConfig() view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetGlobalConfig() ([]byte, error) {
	return _INodeManager.Contract.GetGlobalConfig(&_INodeManager.CallOpts)
}

// GetOutstandingRewards is a free data retrieval call binding the contract method 0xfef97e4c.
//
// Solidity: function getOutstandingRewards() view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetOutstandingRewards(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getOutstandingRewards")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetOutstandingRewards is a free data retrieval call binding the contract method 0xfef97e4c.
//
// Solidity: function getOutstandingRewards() view returns(bytes)
func (_INodeManager *INodeManagerSession) GetOutstandingRewards() ([]byte, error) {
	return _INodeManager.Contract.GetOutstandingRewards(&_INodeManager.CallOpts)
}

// GetOutstandingRewards is a free data retrieval call binding the contract method 0xfef97e4c.
//
// Solidity: function getOutstandingRewards() view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetOutstandingRewards() ([]byte, error) {
	return _INodeManager.Contract.GetOutstandingRewards(&_INodeManager.CallOpts)
}

// GetStakeInfo is a free data retrieval call binding the contract method 0xd77c8f14.
//
// Solidity: function getStakeInfo(address consensusAddress, address stakeAddress) view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetStakeInfo(opts *bind.CallOpts, consensusAddress common.Address, stakeAddress common.Address) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getStakeInfo", consensusAddress, stakeAddress)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetStakeInfo is a free data retrieval call binding the contract method 0xd77c8f14.
//
// Solidity: function getStakeInfo(address consensusAddress, address stakeAddress) view returns(bytes)
func (_INodeManager *INodeManagerSession) GetStakeInfo(consensusAddress common.Address, stakeAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetStakeInfo(&_INodeManager.CallOpts, consensusAddress, stakeAddress)
}

// GetStakeInfo is a free data retrieval call binding the contract method 0xd77c8f14.
//
// Solidity: function getStakeInfo(address consensusAddress, address stakeAddress) view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetStakeInfo(consensusAddress common.Address, stakeAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetStakeInfo(&_INodeManager.CallOpts, consensusAddress, stakeAddress)
}

// GetStakeRewards is a free data retrieval call binding the contract method 0xea3f32ff.
//
// Solidity: function getStakeRewards(address consensusAddress, address stakeAddress) view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetStakeRewards(opts *bind.CallOpts, consensusAddress common.Address, stakeAddress common.Address) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getStakeRewards", consensusAddress, stakeAddress)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetStakeRewards is a free data retrieval call binding the contract method 0xea3f32ff.
//
// Solidity: function getStakeRewards(address consensusAddress, address stakeAddress) view returns(bytes)
func (_INodeManager *INodeManagerSession) GetStakeRewards(consensusAddress common.Address, stakeAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetStakeRewards(&_INodeManager.CallOpts, consensusAddress, stakeAddress)
}

// GetStakeRewards is a free data retrieval call binding the contract method 0xea3f32ff.
//
// Solidity: function getStakeRewards(address consensusAddress, address stakeAddress) view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetStakeRewards(consensusAddress common.Address, stakeAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetStakeRewards(&_INodeManager.CallOpts, consensusAddress, stakeAddress)
}

// GetStakeStartingInfo is a free data retrieval call binding the contract method 0x17674715.
//
// Solidity: function getStakeStartingInfo(address consensusAddress, address stakeAddress) view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetStakeStartingInfo(opts *bind.CallOpts, consensusAddress common.Address, stakeAddress common.Address) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getStakeStartingInfo", consensusAddress, stakeAddress)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetStakeStartingInfo is a free data retrieval call binding the contract method 0x17674715.
//
// Solidity: function getStakeStartingInfo(address consensusAddress, address stakeAddress) view returns(bytes)
func (_INodeManager *INodeManagerSession) GetStakeStartingInfo(consensusAddress common.Address, stakeAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetStakeStartingInfo(&_INodeManager.CallOpts, consensusAddress, stakeAddress)
}

// GetStakeStartingInfo is a free data retrieval call binding the contract method 0x17674715.
//
// Solidity: function getStakeStartingInfo(address consensusAddress, address stakeAddress) view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetStakeStartingInfo(consensusAddress common.Address, stakeAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetStakeStartingInfo(&_INodeManager.CallOpts, consensusAddress, stakeAddress)
}

// GetTotalPool is a free data retrieval call binding the contract method 0x75f4d677.
//
// Solidity: function getTotalPool() view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetTotalPool(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getTotalPool")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetTotalPool is a free data retrieval call binding the contract method 0x75f4d677.
//
// Solidity: function getTotalPool() view returns(bytes)
func (_INodeManager *INodeManagerSession) GetTotalPool() ([]byte, error) {
	return _INodeManager.Contract.GetTotalPool(&_INodeManager.CallOpts)
}

// GetTotalPool is a free data retrieval call binding the contract method 0x75f4d677.
//
// Solidity: function getTotalPool() view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetTotalPool() ([]byte, error) {
	return _INodeManager.Contract.GetTotalPool(&_INodeManager.CallOpts)
}

// GetUnlockingInfo is a free data retrieval call binding the contract method 0x5e45511e.
//
// Solidity: function getUnlockingInfo(address stakeAddress) view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetUnlockingInfo(opts *bind.CallOpts, stakeAddress common.Address) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getUnlockingInfo", stakeAddress)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetUnlockingInfo is a free data retrieval call binding the contract method 0x5e45511e.
//
// Solidity: function getUnlockingInfo(address stakeAddress) view returns(bytes)
func (_INodeManager *INodeManagerSession) GetUnlockingInfo(stakeAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetUnlockingInfo(&_INodeManager.CallOpts, stakeAddress)
}

// GetUnlockingInfo is a free data retrieval call binding the contract method 0x5e45511e.
//
// Solidity: function getUnlockingInfo(address stakeAddress) view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetUnlockingInfo(stakeAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetUnlockingInfo(&_INodeManager.CallOpts, stakeAddress)
}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address consensusAddress) view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetValidator(opts *bind.CallOpts, consensusAddress common.Address) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getValidator", consensusAddress)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address consensusAddress) view returns(bytes)
func (_INodeManager *INodeManagerSession) GetValidator(consensusAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetValidator(&_INodeManager.CallOpts, consensusAddress)
}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address consensusAddress) view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetValidator(consensusAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetValidator(&_INodeManager.CallOpts, consensusAddress)
}

// GetValidatorAccumulatedRewards is a free data retrieval call binding the contract method 0x9c898a3b.
//
// Solidity: function getValidatorAccumulatedRewards(address consensusAddress) view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetValidatorAccumulatedRewards(opts *bind.CallOpts, consensusAddress common.Address) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getValidatorAccumulatedRewards", consensusAddress)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetValidatorAccumulatedRewards is a free data retrieval call binding the contract method 0x9c898a3b.
//
// Solidity: function getValidatorAccumulatedRewards(address consensusAddress) view returns(bytes)
func (_INodeManager *INodeManagerSession) GetValidatorAccumulatedRewards(consensusAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetValidatorAccumulatedRewards(&_INodeManager.CallOpts, consensusAddress)
}

// GetValidatorAccumulatedRewards is a free data retrieval call binding the contract method 0x9c898a3b.
//
// Solidity: function getValidatorAccumulatedRewards(address consensusAddress) view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetValidatorAccumulatedRewards(consensusAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetValidatorAccumulatedRewards(&_INodeManager.CallOpts, consensusAddress)
}

// GetValidatorOutstandingRewards is a free data retrieval call binding the contract method 0xa76d00a8.
//
// Solidity: function getValidatorOutstandingRewards(address consensusAddress) view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetValidatorOutstandingRewards(opts *bind.CallOpts, consensusAddress common.Address) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getValidatorOutstandingRewards", consensusAddress)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetValidatorOutstandingRewards is a free data retrieval call binding the contract method 0xa76d00a8.
//
// Solidity: function getValidatorOutstandingRewards(address consensusAddress) view returns(bytes)
func (_INodeManager *INodeManagerSession) GetValidatorOutstandingRewards(consensusAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetValidatorOutstandingRewards(&_INodeManager.CallOpts, consensusAddress)
}

// GetValidatorOutstandingRewards is a free data retrieval call binding the contract method 0xa76d00a8.
//
// Solidity: function getValidatorOutstandingRewards(address consensusAddress) view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetValidatorOutstandingRewards(consensusAddress common.Address) ([]byte, error) {
	return _INodeManager.Contract.GetValidatorOutstandingRewards(&_INodeManager.CallOpts, consensusAddress)
}

// GetValidatorSnapshotRewards is a free data retrieval call binding the contract method 0xedd0efa9.
//
// Solidity: function getValidatorSnapshotRewards(address consensusAddress, uint64 period) view returns(bytes)
func (_INodeManager *INodeManagerCaller) GetValidatorSnapshotRewards(opts *bind.CallOpts, consensusAddress common.Address, period uint64) ([]byte, error) {
	var out []interface{}
	err := _INodeManager.contract.Call(opts, &out, "getValidatorSnapshotRewards", consensusAddress, period)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetValidatorSnapshotRewards is a free data retrieval call binding the contract method 0xedd0efa9.
//
// Solidity: function getValidatorSnapshotRewards(address consensusAddress, uint64 period) view returns(bytes)
func (_INodeManager *INodeManagerSession) GetValidatorSnapshotRewards(consensusAddress common.Address, period uint64) ([]byte, error) {
	return _INodeManager.Contract.GetValidatorSnapshotRewards(&_INodeManager.CallOpts, consensusAddress, period)
}

// GetValidatorSnapshotRewards is a free data retrieval call binding the contract method 0xedd0efa9.
//
// Solidity: function getValidatorSnapshotRewards(address consensusAddress, uint64 period) view returns(bytes)
func (_INodeManager *INodeManagerCallerSession) GetValidatorSnapshotRewards(consensusAddress common.Address, period uint64) ([]byte, error) {
	return _INodeManager.Contract.GetValidatorSnapshotRewards(&_INodeManager.CallOpts, consensusAddress, period)
}

// CancelValidator is a paid mutator transaction binding the contract method 0x1af78584.
//
// Solidity: function cancelValidator(address consensusAddress) returns(bool success)
func (_INodeManager *INodeManagerTransactor) CancelValidator(opts *bind.TransactOpts, consensusAddress common.Address) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "cancelValidator", consensusAddress)
}

// CancelValidator is a paid mutator transaction binding the contract method 0x1af78584.
//
// Solidity: function cancelValidator(address consensusAddress) returns(bool success)
func (_INodeManager *INodeManagerSession) CancelValidator(consensusAddress common.Address) (*types.Transaction, error) {
	return _INodeManager.Contract.CancelValidator(&_INodeManager.TransactOpts, consensusAddress)
}

// CancelValidator is a paid mutator transaction binding the contract method 0x1af78584.
//
// Solidity: function cancelValidator(address consensusAddress) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) CancelValidator(consensusAddress common.Address) (*types.Transaction, error) {
	return _INodeManager.Contract.CancelValidator(&_INodeManager.TransactOpts, consensusAddress)
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

// CreateValidator is a paid mutator transaction binding the contract method 0xb5c42005.
//
// Solidity: function createValidator(address consensusAddress, address signerAddress, address proposalAddress, int256 commission, int256 initStake, string desc) returns(bool success)
func (_INodeManager *INodeManagerTransactor) CreateValidator(opts *bind.TransactOpts, consensusAddress common.Address, signerAddress common.Address, proposalAddress common.Address, commission *big.Int, initStake *big.Int, desc string) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "createValidator", consensusAddress, signerAddress, proposalAddress, commission, initStake, desc)
}

// CreateValidator is a paid mutator transaction binding the contract method 0xb5c42005.
//
// Solidity: function createValidator(address consensusAddress, address signerAddress, address proposalAddress, int256 commission, int256 initStake, string desc) returns(bool success)
func (_INodeManager *INodeManagerSession) CreateValidator(consensusAddress common.Address, signerAddress common.Address, proposalAddress common.Address, commission *big.Int, initStake *big.Int, desc string) (*types.Transaction, error) {
	return _INodeManager.Contract.CreateValidator(&_INodeManager.TransactOpts, consensusAddress, signerAddress, proposalAddress, commission, initStake, desc)
}

// CreateValidator is a paid mutator transaction binding the contract method 0xb5c42005.
//
// Solidity: function createValidator(address consensusAddress, address signerAddress, address proposalAddress, int256 commission, int256 initStake, string desc) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) CreateValidator(consensusAddress common.Address, signerAddress common.Address, proposalAddress common.Address, commission *big.Int, initStake *big.Int, desc string) (*types.Transaction, error) {
	return _INodeManager.Contract.CreateValidator(&_INodeManager.TransactOpts, consensusAddress, signerAddress, proposalAddress, commission, initStake, desc)
}

// EndBlock is a paid mutator transaction binding the contract method 0x083c6323.
//
// Solidity: function endBlock() returns(bool success)
func (_INodeManager *INodeManagerTransactor) EndBlock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "endBlock")
}

// EndBlock is a paid mutator transaction binding the contract method 0x083c6323.
//
// Solidity: function endBlock() returns(bool success)
func (_INodeManager *INodeManagerSession) EndBlock() (*types.Transaction, error) {
	return _INodeManager.Contract.EndBlock(&_INodeManager.TransactOpts)
}

// EndBlock is a paid mutator transaction binding the contract method 0x083c6323.
//
// Solidity: function endBlock() returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) EndBlock() (*types.Transaction, error) {
	return _INodeManager.Contract.EndBlock(&_INodeManager.TransactOpts)
}

// Stake is a paid mutator transaction binding the contract method 0x91398b13.
//
// Solidity: function stake(address consensusAddress, int256 amount) returns(bool success)
func (_INodeManager *INodeManagerTransactor) Stake(opts *bind.TransactOpts, consensusAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "stake", consensusAddress, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x91398b13.
//
// Solidity: function stake(address consensusAddress, int256 amount) returns(bool success)
func (_INodeManager *INodeManagerSession) Stake(consensusAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _INodeManager.Contract.Stake(&_INodeManager.TransactOpts, consensusAddress, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x91398b13.
//
// Solidity: function stake(address consensusAddress, int256 amount) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) Stake(consensusAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _INodeManager.Contract.Stake(&_INodeManager.TransactOpts, consensusAddress, amount)
}

// UnStake is a paid mutator transaction binding the contract method 0xdfe6bad3.
//
// Solidity: function unStake(address consensusAddress, int256 amount) returns(bool success)
func (_INodeManager *INodeManagerTransactor) UnStake(opts *bind.TransactOpts, consensusAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "unStake", consensusAddress, amount)
}

// UnStake is a paid mutator transaction binding the contract method 0xdfe6bad3.
//
// Solidity: function unStake(address consensusAddress, int256 amount) returns(bool success)
func (_INodeManager *INodeManagerSession) UnStake(consensusAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _INodeManager.Contract.UnStake(&_INodeManager.TransactOpts, consensusAddress, amount)
}

// UnStake is a paid mutator transaction binding the contract method 0xdfe6bad3.
//
// Solidity: function unStake(address consensusAddress, int256 amount) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) UnStake(consensusAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _INodeManager.Contract.UnStake(&_INodeManager.TransactOpts, consensusAddress, amount)
}

// UpdateCommission is a paid mutator transaction binding the contract method 0xc5e7ad1d.
//
// Solidity: function updateCommission(address consensusAddress, int256 commission) returns(bool success)
func (_INodeManager *INodeManagerTransactor) UpdateCommission(opts *bind.TransactOpts, consensusAddress common.Address, commission *big.Int) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "updateCommission", consensusAddress, commission)
}

// UpdateCommission is a paid mutator transaction binding the contract method 0xc5e7ad1d.
//
// Solidity: function updateCommission(address consensusAddress, int256 commission) returns(bool success)
func (_INodeManager *INodeManagerSession) UpdateCommission(consensusAddress common.Address, commission *big.Int) (*types.Transaction, error) {
	return _INodeManager.Contract.UpdateCommission(&_INodeManager.TransactOpts, consensusAddress, commission)
}

// UpdateCommission is a paid mutator transaction binding the contract method 0xc5e7ad1d.
//
// Solidity: function updateCommission(address consensusAddress, int256 commission) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) UpdateCommission(consensusAddress common.Address, commission *big.Int) (*types.Transaction, error) {
	return _INodeManager.Contract.UpdateCommission(&_INodeManager.TransactOpts, consensusAddress, commission)
}

// UpdateValidator is a paid mutator transaction binding the contract method 0x24712218.
//
// Solidity: function updateValidator(address consensusAddress, address signerAddress, address proposalAddress, string desc) returns(bool success)
func (_INodeManager *INodeManagerTransactor) UpdateValidator(opts *bind.TransactOpts, consensusAddress common.Address, signerAddress common.Address, proposalAddress common.Address, desc string) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "updateValidator", consensusAddress, signerAddress, proposalAddress, desc)
}

// UpdateValidator is a paid mutator transaction binding the contract method 0x24712218.
//
// Solidity: function updateValidator(address consensusAddress, address signerAddress, address proposalAddress, string desc) returns(bool success)
func (_INodeManager *INodeManagerSession) UpdateValidator(consensusAddress common.Address, signerAddress common.Address, proposalAddress common.Address, desc string) (*types.Transaction, error) {
	return _INodeManager.Contract.UpdateValidator(&_INodeManager.TransactOpts, consensusAddress, signerAddress, proposalAddress, desc)
}

// UpdateValidator is a paid mutator transaction binding the contract method 0x24712218.
//
// Solidity: function updateValidator(address consensusAddress, address signerAddress, address proposalAddress, string desc) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) UpdateValidator(consensusAddress common.Address, signerAddress common.Address, proposalAddress common.Address, desc string) (*types.Transaction, error) {
	return _INodeManager.Contract.UpdateValidator(&_INodeManager.TransactOpts, consensusAddress, signerAddress, proposalAddress, desc)
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

// WithdrawCommission is a paid mutator transaction binding the contract method 0x16c58d04.
//
// Solidity: function withdrawCommission(address consensusAddress) returns(bool success)
func (_INodeManager *INodeManagerTransactor) WithdrawCommission(opts *bind.TransactOpts, consensusAddress common.Address) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "withdrawCommission", consensusAddress)
}

// WithdrawCommission is a paid mutator transaction binding the contract method 0x16c58d04.
//
// Solidity: function withdrawCommission(address consensusAddress) returns(bool success)
func (_INodeManager *INodeManagerSession) WithdrawCommission(consensusAddress common.Address) (*types.Transaction, error) {
	return _INodeManager.Contract.WithdrawCommission(&_INodeManager.TransactOpts, consensusAddress)
}

// WithdrawCommission is a paid mutator transaction binding the contract method 0x16c58d04.
//
// Solidity: function withdrawCommission(address consensusAddress) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) WithdrawCommission(consensusAddress common.Address) (*types.Transaction, error) {
	return _INodeManager.Contract.WithdrawCommission(&_INodeManager.TransactOpts, consensusAddress)
}

// WithdrawStakeRewards is a paid mutator transaction binding the contract method 0x5d3abd54.
//
// Solidity: function withdrawStakeRewards(address consensusAddress) returns(bool success)
func (_INodeManager *INodeManagerTransactor) WithdrawStakeRewards(opts *bind.TransactOpts, consensusAddress common.Address) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "withdrawStakeRewards", consensusAddress)
}

// WithdrawStakeRewards is a paid mutator transaction binding the contract method 0x5d3abd54.
//
// Solidity: function withdrawStakeRewards(address consensusAddress) returns(bool success)
func (_INodeManager *INodeManagerSession) WithdrawStakeRewards(consensusAddress common.Address) (*types.Transaction, error) {
	return _INodeManager.Contract.WithdrawStakeRewards(&_INodeManager.TransactOpts, consensusAddress)
}

// WithdrawStakeRewards is a paid mutator transaction binding the contract method 0x5d3abd54.
//
// Solidity: function withdrawStakeRewards(address consensusAddress) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) WithdrawStakeRewards(consensusAddress common.Address) (*types.Transaction, error) {
	return _INodeManager.Contract.WithdrawStakeRewards(&_INodeManager.TransactOpts, consensusAddress)
}

// WithdrawValidator is a paid mutator transaction binding the contract method 0x814a2c11.
//
// Solidity: function withdrawValidator(address consensusAddress) returns(bool success)
func (_INodeManager *INodeManagerTransactor) WithdrawValidator(opts *bind.TransactOpts, consensusAddress common.Address) (*types.Transaction, error) {
	return _INodeManager.contract.Transact(opts, "withdrawValidator", consensusAddress)
}

// WithdrawValidator is a paid mutator transaction binding the contract method 0x814a2c11.
//
// Solidity: function withdrawValidator(address consensusAddress) returns(bool success)
func (_INodeManager *INodeManagerSession) WithdrawValidator(consensusAddress common.Address) (*types.Transaction, error) {
	return _INodeManager.Contract.WithdrawValidator(&_INodeManager.TransactOpts, consensusAddress)
}

// WithdrawValidator is a paid mutator transaction binding the contract method 0x814a2c11.
//
// Solidity: function withdrawValidator(address consensusAddress) returns(bool success)
func (_INodeManager *INodeManagerTransactorSession) WithdrawValidator(consensusAddress common.Address) (*types.Transaction, error) {
	return _INodeManager.Contract.WithdrawValidator(&_INodeManager.TransactOpts, consensusAddress)
}

// INodeManagerCancelValidatorIterator is returned from FilterCancelValidator and is used to iterate over the raw logs and unpacked data for CancelValidator events raised by the INodeManager contract.
type INodeManagerCancelValidatorIterator struct {
	Event *INodeManagerCancelValidator // Event containing the contract specifics and raw log

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
func (it *INodeManagerCancelValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(INodeManagerCancelValidator)
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
		it.Event = new(INodeManagerCancelValidator)
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
func (it *INodeManagerCancelValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *INodeManagerCancelValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// INodeManagerCancelValidator represents a CancelValidator event raised by the INodeManager contract.
type INodeManagerCancelValidator struct {
	ConsensusAddress string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterCancelValidator is a free log retrieval operation binding the contract event 0x958925709318fd39dab7c4c9812b315994b68e8d17a3408c1ca4bf0dc594473d.
//
// Solidity: event CancelValidator(string consensusAddress)
func (_INodeManager *INodeManagerFilterer) FilterCancelValidator(opts *bind.FilterOpts) (*INodeManagerCancelValidatorIterator, error) {

	logs, sub, err := _INodeManager.contract.FilterLogs(opts, "CancelValidator")
	if err != nil {
		return nil, err
	}
	return &INodeManagerCancelValidatorIterator{contract: _INodeManager.contract, event: "CancelValidator", logs: logs, sub: sub}, nil
}

// WatchCancelValidator is a free log subscription operation binding the contract event 0x958925709318fd39dab7c4c9812b315994b68e8d17a3408c1ca4bf0dc594473d.
//
// Solidity: event CancelValidator(string consensusAddress)
func (_INodeManager *INodeManagerFilterer) WatchCancelValidator(opts *bind.WatchOpts, sink chan<- *INodeManagerCancelValidator) (event.Subscription, error) {

	logs, sub, err := _INodeManager.contract.WatchLogs(opts, "CancelValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(INodeManagerCancelValidator)
				if err := _INodeManager.contract.UnpackLog(event, "CancelValidator", log); err != nil {
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

// ParseCancelValidator is a log parse operation binding the contract event 0x958925709318fd39dab7c4c9812b315994b68e8d17a3408c1ca4bf0dc594473d.
//
// Solidity: event CancelValidator(string consensusAddress)
func (_INodeManager *INodeManagerFilterer) ParseCancelValidator(log types.Log) (*INodeManagerCancelValidator, error) {
	event := new(INodeManagerCancelValidator)
	if err := _INodeManager.contract.UnpackLog(event, "CancelValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// INodeManagerChangeEpochIterator is returned from FilterChangeEpoch and is used to iterate over the raw logs and unpacked data for ChangeEpoch events raised by the INodeManager contract.
type INodeManagerChangeEpochIterator struct {
	Event *INodeManagerChangeEpoch // Event containing the contract specifics and raw log

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
func (it *INodeManagerChangeEpochIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(INodeManagerChangeEpoch)
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
		it.Event = new(INodeManagerChangeEpoch)
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
func (it *INodeManagerChangeEpochIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *INodeManagerChangeEpochIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// INodeManagerChangeEpoch represents a ChangeEpoch event raised by the INodeManager contract.
type INodeManagerChangeEpoch struct {
	EpochID string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterChangeEpoch is a free log retrieval operation binding the contract event 0xf92b44dcfd5f229b4a15bfb78fe69d9bc16bfeae87450b6fe3b33d11074d4330.
//
// Solidity: event ChangeEpoch(string epochID)
func (_INodeManager *INodeManagerFilterer) FilterChangeEpoch(opts *bind.FilterOpts) (*INodeManagerChangeEpochIterator, error) {

	logs, sub, err := _INodeManager.contract.FilterLogs(opts, "ChangeEpoch")
	if err != nil {
		return nil, err
	}
	return &INodeManagerChangeEpochIterator{contract: _INodeManager.contract, event: "ChangeEpoch", logs: logs, sub: sub}, nil
}

// WatchChangeEpoch is a free log subscription operation binding the contract event 0xf92b44dcfd5f229b4a15bfb78fe69d9bc16bfeae87450b6fe3b33d11074d4330.
//
// Solidity: event ChangeEpoch(string epochID)
func (_INodeManager *INodeManagerFilterer) WatchChangeEpoch(opts *bind.WatchOpts, sink chan<- *INodeManagerChangeEpoch) (event.Subscription, error) {

	logs, sub, err := _INodeManager.contract.WatchLogs(opts, "ChangeEpoch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(INodeManagerChangeEpoch)
				if err := _INodeManager.contract.UnpackLog(event, "ChangeEpoch", log); err != nil {
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

// ParseChangeEpoch is a log parse operation binding the contract event 0xf92b44dcfd5f229b4a15bfb78fe69d9bc16bfeae87450b6fe3b33d11074d4330.
//
// Solidity: event ChangeEpoch(string epochID)
func (_INodeManager *INodeManagerFilterer) ParseChangeEpoch(log types.Log) (*INodeManagerChangeEpoch, error) {
	event := new(INodeManagerChangeEpoch)
	if err := _INodeManager.contract.UnpackLog(event, "ChangeEpoch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// INodeManagerCreateValidatorIterator is returned from FilterCreateValidator and is used to iterate over the raw logs and unpacked data for CreateValidator events raised by the INodeManager contract.
type INodeManagerCreateValidatorIterator struct {
	Event *INodeManagerCreateValidator // Event containing the contract specifics and raw log

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
func (it *INodeManagerCreateValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(INodeManagerCreateValidator)
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
		it.Event = new(INodeManagerCreateValidator)
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
func (it *INodeManagerCreateValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *INodeManagerCreateValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// INodeManagerCreateValidator represents a CreateValidator event raised by the INodeManager contract.
type INodeManagerCreateValidator struct {
	ConsensusAddress string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterCreateValidator is a free log retrieval operation binding the contract event 0xb7f4cfc364000338326cb6f95799a39e25447cd02f70d1c7568f2d8d0a6fa2df.
//
// Solidity: event CreateValidator(string consensusAddress)
func (_INodeManager *INodeManagerFilterer) FilterCreateValidator(opts *bind.FilterOpts) (*INodeManagerCreateValidatorIterator, error) {

	logs, sub, err := _INodeManager.contract.FilterLogs(opts, "CreateValidator")
	if err != nil {
		return nil, err
	}
	return &INodeManagerCreateValidatorIterator{contract: _INodeManager.contract, event: "CreateValidator", logs: logs, sub: sub}, nil
}

// WatchCreateValidator is a free log subscription operation binding the contract event 0xb7f4cfc364000338326cb6f95799a39e25447cd02f70d1c7568f2d8d0a6fa2df.
//
// Solidity: event CreateValidator(string consensusAddress)
func (_INodeManager *INodeManagerFilterer) WatchCreateValidator(opts *bind.WatchOpts, sink chan<- *INodeManagerCreateValidator) (event.Subscription, error) {

	logs, sub, err := _INodeManager.contract.WatchLogs(opts, "CreateValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(INodeManagerCreateValidator)
				if err := _INodeManager.contract.UnpackLog(event, "CreateValidator", log); err != nil {
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

// ParseCreateValidator is a log parse operation binding the contract event 0xb7f4cfc364000338326cb6f95799a39e25447cd02f70d1c7568f2d8d0a6fa2df.
//
// Solidity: event CreateValidator(string consensusAddress)
func (_INodeManager *INodeManagerFilterer) ParseCreateValidator(log types.Log) (*INodeManagerCreateValidator, error) {
	event := new(INodeManagerCreateValidator)
	if err := _INodeManager.contract.UnpackLog(event, "CreateValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// INodeManagerStakeIterator is returned from FilterStake and is used to iterate over the raw logs and unpacked data for Stake events raised by the INodeManager contract.
type INodeManagerStakeIterator struct {
	Event *INodeManagerStake // Event containing the contract specifics and raw log

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
func (it *INodeManagerStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(INodeManagerStake)
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
		it.Event = new(INodeManagerStake)
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
func (it *INodeManagerStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *INodeManagerStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// INodeManagerStake represents a Stake event raised by the INodeManager contract.
type INodeManagerStake struct {
	ConsensusAddress string
	Amount           string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterStake is a free log retrieval operation binding the contract event 0x28c6fab16944691610b9d650c8296cb50ed85a965e81865eb3f1f3ccc274a272.
//
// Solidity: event Stake(string consensusAddress, string amount)
func (_INodeManager *INodeManagerFilterer) FilterStake(opts *bind.FilterOpts) (*INodeManagerStakeIterator, error) {

	logs, sub, err := _INodeManager.contract.FilterLogs(opts, "Stake")
	if err != nil {
		return nil, err
	}
	return &INodeManagerStakeIterator{contract: _INodeManager.contract, event: "Stake", logs: logs, sub: sub}, nil
}

// WatchStake is a free log subscription operation binding the contract event 0x28c6fab16944691610b9d650c8296cb50ed85a965e81865eb3f1f3ccc274a272.
//
// Solidity: event Stake(string consensusAddress, string amount)
func (_INodeManager *INodeManagerFilterer) WatchStake(opts *bind.WatchOpts, sink chan<- *INodeManagerStake) (event.Subscription, error) {

	logs, sub, err := _INodeManager.contract.WatchLogs(opts, "Stake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(INodeManagerStake)
				if err := _INodeManager.contract.UnpackLog(event, "Stake", log); err != nil {
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

// ParseStake is a log parse operation binding the contract event 0x28c6fab16944691610b9d650c8296cb50ed85a965e81865eb3f1f3ccc274a272.
//
// Solidity: event Stake(string consensusAddress, string amount)
func (_INodeManager *INodeManagerFilterer) ParseStake(log types.Log) (*INodeManagerStake, error) {
	event := new(INodeManagerStake)
	if err := _INodeManager.contract.UnpackLog(event, "Stake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// INodeManagerUnStakeIterator is returned from FilterUnStake and is used to iterate over the raw logs and unpacked data for UnStake events raised by the INodeManager contract.
type INodeManagerUnStakeIterator struct {
	Event *INodeManagerUnStake // Event containing the contract specifics and raw log

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
func (it *INodeManagerUnStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(INodeManagerUnStake)
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
		it.Event = new(INodeManagerUnStake)
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
func (it *INodeManagerUnStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *INodeManagerUnStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// INodeManagerUnStake represents a UnStake event raised by the INodeManager contract.
type INodeManagerUnStake struct {
	ConsensusAddress string
	Amount           string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUnStake is a free log retrieval operation binding the contract event 0x09c079860913a3cf3561020297d9982bdeb613ecdb83920f88063bf6e3e19088.
//
// Solidity: event UnStake(string consensusAddress, string amount)
func (_INodeManager *INodeManagerFilterer) FilterUnStake(opts *bind.FilterOpts) (*INodeManagerUnStakeIterator, error) {

	logs, sub, err := _INodeManager.contract.FilterLogs(opts, "UnStake")
	if err != nil {
		return nil, err
	}
	return &INodeManagerUnStakeIterator{contract: _INodeManager.contract, event: "UnStake", logs: logs, sub: sub}, nil
}

// WatchUnStake is a free log subscription operation binding the contract event 0x09c079860913a3cf3561020297d9982bdeb613ecdb83920f88063bf6e3e19088.
//
// Solidity: event UnStake(string consensusAddress, string amount)
func (_INodeManager *INodeManagerFilterer) WatchUnStake(opts *bind.WatchOpts, sink chan<- *INodeManagerUnStake) (event.Subscription, error) {

	logs, sub, err := _INodeManager.contract.WatchLogs(opts, "UnStake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(INodeManagerUnStake)
				if err := _INodeManager.contract.UnpackLog(event, "UnStake", log); err != nil {
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

// ParseUnStake is a log parse operation binding the contract event 0x09c079860913a3cf3561020297d9982bdeb613ecdb83920f88063bf6e3e19088.
//
// Solidity: event UnStake(string consensusAddress, string amount)
func (_INodeManager *INodeManagerFilterer) ParseUnStake(log types.Log) (*INodeManagerUnStake, error) {
	event := new(INodeManagerUnStake)
	if err := _INodeManager.contract.UnpackLog(event, "UnStake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// INodeManagerUpdateCommissionIterator is returned from FilterUpdateCommission and is used to iterate over the raw logs and unpacked data for UpdateCommission events raised by the INodeManager contract.
type INodeManagerUpdateCommissionIterator struct {
	Event *INodeManagerUpdateCommission // Event containing the contract specifics and raw log

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
func (it *INodeManagerUpdateCommissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(INodeManagerUpdateCommission)
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
		it.Event = new(INodeManagerUpdateCommission)
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
func (it *INodeManagerUpdateCommissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *INodeManagerUpdateCommissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// INodeManagerUpdateCommission represents a UpdateCommission event raised by the INodeManager contract.
type INodeManagerUpdateCommission struct {
	ConsensusAddress string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpdateCommission is a free log retrieval operation binding the contract event 0x4122b09ab8922a4e7eb98b7aade17dd076d3eb78158a17a0ab826b857b17e2b2.
//
// Solidity: event UpdateCommission(string consensusAddress)
func (_INodeManager *INodeManagerFilterer) FilterUpdateCommission(opts *bind.FilterOpts) (*INodeManagerUpdateCommissionIterator, error) {

	logs, sub, err := _INodeManager.contract.FilterLogs(opts, "UpdateCommission")
	if err != nil {
		return nil, err
	}
	return &INodeManagerUpdateCommissionIterator{contract: _INodeManager.contract, event: "UpdateCommission", logs: logs, sub: sub}, nil
}

// WatchUpdateCommission is a free log subscription operation binding the contract event 0x4122b09ab8922a4e7eb98b7aade17dd076d3eb78158a17a0ab826b857b17e2b2.
//
// Solidity: event UpdateCommission(string consensusAddress)
func (_INodeManager *INodeManagerFilterer) WatchUpdateCommission(opts *bind.WatchOpts, sink chan<- *INodeManagerUpdateCommission) (event.Subscription, error) {

	logs, sub, err := _INodeManager.contract.WatchLogs(opts, "UpdateCommission")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(INodeManagerUpdateCommission)
				if err := _INodeManager.contract.UnpackLog(event, "UpdateCommission", log); err != nil {
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

// ParseUpdateCommission is a log parse operation binding the contract event 0x4122b09ab8922a4e7eb98b7aade17dd076d3eb78158a17a0ab826b857b17e2b2.
//
// Solidity: event UpdateCommission(string consensusAddress)
func (_INodeManager *INodeManagerFilterer) ParseUpdateCommission(log types.Log) (*INodeManagerUpdateCommission, error) {
	event := new(INodeManagerUpdateCommission)
	if err := _INodeManager.contract.UnpackLog(event, "UpdateCommission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// INodeManagerUpdateValidatorIterator is returned from FilterUpdateValidator and is used to iterate over the raw logs and unpacked data for UpdateValidator events raised by the INodeManager contract.
type INodeManagerUpdateValidatorIterator struct {
	Event *INodeManagerUpdateValidator // Event containing the contract specifics and raw log

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
func (it *INodeManagerUpdateValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(INodeManagerUpdateValidator)
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
		it.Event = new(INodeManagerUpdateValidator)
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
func (it *INodeManagerUpdateValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *INodeManagerUpdateValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// INodeManagerUpdateValidator represents a UpdateValidator event raised by the INodeManager contract.
type INodeManagerUpdateValidator struct {
	ConsensusAddress string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpdateValidator is a free log retrieval operation binding the contract event 0xf6d6118bb8579adc14af14122184c92bc5fa2d973b612298019af0fab85640c1.
//
// Solidity: event UpdateValidator(string consensusAddress)
func (_INodeManager *INodeManagerFilterer) FilterUpdateValidator(opts *bind.FilterOpts) (*INodeManagerUpdateValidatorIterator, error) {

	logs, sub, err := _INodeManager.contract.FilterLogs(opts, "UpdateValidator")
	if err != nil {
		return nil, err
	}
	return &INodeManagerUpdateValidatorIterator{contract: _INodeManager.contract, event: "UpdateValidator", logs: logs, sub: sub}, nil
}

// WatchUpdateValidator is a free log subscription operation binding the contract event 0xf6d6118bb8579adc14af14122184c92bc5fa2d973b612298019af0fab85640c1.
//
// Solidity: event UpdateValidator(string consensusAddress)
func (_INodeManager *INodeManagerFilterer) WatchUpdateValidator(opts *bind.WatchOpts, sink chan<- *INodeManagerUpdateValidator) (event.Subscription, error) {

	logs, sub, err := _INodeManager.contract.WatchLogs(opts, "UpdateValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(INodeManagerUpdateValidator)
				if err := _INodeManager.contract.UnpackLog(event, "UpdateValidator", log); err != nil {
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

// ParseUpdateValidator is a log parse operation binding the contract event 0xf6d6118bb8579adc14af14122184c92bc5fa2d973b612298019af0fab85640c1.
//
// Solidity: event UpdateValidator(string consensusAddress)
func (_INodeManager *INodeManagerFilterer) ParseUpdateValidator(log types.Log) (*INodeManagerUpdateValidator, error) {
	event := new(INodeManagerUpdateValidator)
	if err := _INodeManager.contract.UnpackLog(event, "UpdateValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// INodeManagerWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the INodeManager contract.
type INodeManagerWithdrawIterator struct {
	Event *INodeManagerWithdraw // Event containing the contract specifics and raw log

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
func (it *INodeManagerWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(INodeManagerWithdraw)
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
		it.Event = new(INodeManagerWithdraw)
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
func (it *INodeManagerWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *INodeManagerWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// INodeManagerWithdraw represents a Withdraw event raised by the INodeManager contract.
type INodeManagerWithdraw struct {
	Caller string
	Amount string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x8611643a3aa3bfaaff3142871c24f5b4013939f2c2b7a5a36aaf32d5e1c68994.
//
// Solidity: event Withdraw(string caller, string amount)
func (_INodeManager *INodeManagerFilterer) FilterWithdraw(opts *bind.FilterOpts) (*INodeManagerWithdrawIterator, error) {

	logs, sub, err := _INodeManager.contract.FilterLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return &INodeManagerWithdrawIterator{contract: _INodeManager.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x8611643a3aa3bfaaff3142871c24f5b4013939f2c2b7a5a36aaf32d5e1c68994.
//
// Solidity: event Withdraw(string caller, string amount)
func (_INodeManager *INodeManagerFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *INodeManagerWithdraw) (event.Subscription, error) {

	logs, sub, err := _INodeManager.contract.WatchLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(INodeManagerWithdraw)
				if err := _INodeManager.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0x8611643a3aa3bfaaff3142871c24f5b4013939f2c2b7a5a36aaf32d5e1c68994.
//
// Solidity: event Withdraw(string caller, string amount)
func (_INodeManager *INodeManagerFilterer) ParseWithdraw(log types.Log) (*INodeManagerWithdraw, error) {
	event := new(INodeManagerWithdraw)
	if err := _INodeManager.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// INodeManagerWithdrawCommissionIterator is returned from FilterWithdrawCommission and is used to iterate over the raw logs and unpacked data for WithdrawCommission events raised by the INodeManager contract.
type INodeManagerWithdrawCommissionIterator struct {
	Event *INodeManagerWithdrawCommission // Event containing the contract specifics and raw log

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
func (it *INodeManagerWithdrawCommissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(INodeManagerWithdrawCommission)
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
		it.Event = new(INodeManagerWithdrawCommission)
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
func (it *INodeManagerWithdrawCommissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *INodeManagerWithdrawCommissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// INodeManagerWithdrawCommission represents a WithdrawCommission event raised by the INodeManager contract.
type INodeManagerWithdrawCommission struct {
	ConsensusAddress string
	Commission       string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterWithdrawCommission is a free log retrieval operation binding the contract event 0xb732e4a208911fcc74b6416fb6480712d8e9f03baecb7acc08372cb8a4ea64a7.
//
// Solidity: event WithdrawCommission(string consensusAddress, string commission)
func (_INodeManager *INodeManagerFilterer) FilterWithdrawCommission(opts *bind.FilterOpts) (*INodeManagerWithdrawCommissionIterator, error) {

	logs, sub, err := _INodeManager.contract.FilterLogs(opts, "WithdrawCommission")
	if err != nil {
		return nil, err
	}
	return &INodeManagerWithdrawCommissionIterator{contract: _INodeManager.contract, event: "WithdrawCommission", logs: logs, sub: sub}, nil
}

// WatchWithdrawCommission is a free log subscription operation binding the contract event 0xb732e4a208911fcc74b6416fb6480712d8e9f03baecb7acc08372cb8a4ea64a7.
//
// Solidity: event WithdrawCommission(string consensusAddress, string commission)
func (_INodeManager *INodeManagerFilterer) WatchWithdrawCommission(opts *bind.WatchOpts, sink chan<- *INodeManagerWithdrawCommission) (event.Subscription, error) {

	logs, sub, err := _INodeManager.contract.WatchLogs(opts, "WithdrawCommission")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(INodeManagerWithdrawCommission)
				if err := _INodeManager.contract.UnpackLog(event, "WithdrawCommission", log); err != nil {
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

// ParseWithdrawCommission is a log parse operation binding the contract event 0xb732e4a208911fcc74b6416fb6480712d8e9f03baecb7acc08372cb8a4ea64a7.
//
// Solidity: event WithdrawCommission(string consensusAddress, string commission)
func (_INodeManager *INodeManagerFilterer) ParseWithdrawCommission(log types.Log) (*INodeManagerWithdrawCommission, error) {
	event := new(INodeManagerWithdrawCommission)
	if err := _INodeManager.contract.UnpackLog(event, "WithdrawCommission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// INodeManagerWithdrawStakeRewardsIterator is returned from FilterWithdrawStakeRewards and is used to iterate over the raw logs and unpacked data for WithdrawStakeRewards events raised by the INodeManager contract.
type INodeManagerWithdrawStakeRewardsIterator struct {
	Event *INodeManagerWithdrawStakeRewards // Event containing the contract specifics and raw log

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
func (it *INodeManagerWithdrawStakeRewardsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(INodeManagerWithdrawStakeRewards)
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
		it.Event = new(INodeManagerWithdrawStakeRewards)
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
func (it *INodeManagerWithdrawStakeRewardsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *INodeManagerWithdrawStakeRewardsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// INodeManagerWithdrawStakeRewards represents a WithdrawStakeRewards event raised by the INodeManager contract.
type INodeManagerWithdrawStakeRewards struct {
	ConsensusAddress string
	Caller           string
	Rewards          string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterWithdrawStakeRewards is a free log retrieval operation binding the contract event 0xa1d663f7088e276bf2b90b1586311ae4ff6584f7d3cfaa901c4f369c999d152e.
//
// Solidity: event WithdrawStakeRewards(string consensusAddress, string caller, string rewards)
func (_INodeManager *INodeManagerFilterer) FilterWithdrawStakeRewards(opts *bind.FilterOpts) (*INodeManagerWithdrawStakeRewardsIterator, error) {

	logs, sub, err := _INodeManager.contract.FilterLogs(opts, "WithdrawStakeRewards")
	if err != nil {
		return nil, err
	}
	return &INodeManagerWithdrawStakeRewardsIterator{contract: _INodeManager.contract, event: "WithdrawStakeRewards", logs: logs, sub: sub}, nil
}

// WatchWithdrawStakeRewards is a free log subscription operation binding the contract event 0xa1d663f7088e276bf2b90b1586311ae4ff6584f7d3cfaa901c4f369c999d152e.
//
// Solidity: event WithdrawStakeRewards(string consensusAddress, string caller, string rewards)
func (_INodeManager *INodeManagerFilterer) WatchWithdrawStakeRewards(opts *bind.WatchOpts, sink chan<- *INodeManagerWithdrawStakeRewards) (event.Subscription, error) {

	logs, sub, err := _INodeManager.contract.WatchLogs(opts, "WithdrawStakeRewards")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(INodeManagerWithdrawStakeRewards)
				if err := _INodeManager.contract.UnpackLog(event, "WithdrawStakeRewards", log); err != nil {
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

// ParseWithdrawStakeRewards is a log parse operation binding the contract event 0xa1d663f7088e276bf2b90b1586311ae4ff6584f7d3cfaa901c4f369c999d152e.
//
// Solidity: event WithdrawStakeRewards(string consensusAddress, string caller, string rewards)
func (_INodeManager *INodeManagerFilterer) ParseWithdrawStakeRewards(log types.Log) (*INodeManagerWithdrawStakeRewards, error) {
	event := new(INodeManagerWithdrawStakeRewards)
	if err := _INodeManager.contract.UnpackLog(event, "WithdrawStakeRewards", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// INodeManagerWithdrawValidatorIterator is returned from FilterWithdrawValidator and is used to iterate over the raw logs and unpacked data for WithdrawValidator events raised by the INodeManager contract.
type INodeManagerWithdrawValidatorIterator struct {
	Event *INodeManagerWithdrawValidator // Event containing the contract specifics and raw log

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
func (it *INodeManagerWithdrawValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(INodeManagerWithdrawValidator)
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
		it.Event = new(INodeManagerWithdrawValidator)
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
func (it *INodeManagerWithdrawValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *INodeManagerWithdrawValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// INodeManagerWithdrawValidator represents a WithdrawValidator event raised by the INodeManager contract.
type INodeManagerWithdrawValidator struct {
	ConsensusAddress string
	SelfStake        string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterWithdrawValidator is a free log retrieval operation binding the contract event 0xa5bee9e1a2697e6ed6ddf963ed42431f3c8c594f49c697cf1391bbf74e4ea877.
//
// Solidity: event WithdrawValidator(string consensusAddress, string selfStake)
func (_INodeManager *INodeManagerFilterer) FilterWithdrawValidator(opts *bind.FilterOpts) (*INodeManagerWithdrawValidatorIterator, error) {

	logs, sub, err := _INodeManager.contract.FilterLogs(opts, "WithdrawValidator")
	if err != nil {
		return nil, err
	}
	return &INodeManagerWithdrawValidatorIterator{contract: _INodeManager.contract, event: "WithdrawValidator", logs: logs, sub: sub}, nil
}

// WatchWithdrawValidator is a free log subscription operation binding the contract event 0xa5bee9e1a2697e6ed6ddf963ed42431f3c8c594f49c697cf1391bbf74e4ea877.
//
// Solidity: event WithdrawValidator(string consensusAddress, string selfStake)
func (_INodeManager *INodeManagerFilterer) WatchWithdrawValidator(opts *bind.WatchOpts, sink chan<- *INodeManagerWithdrawValidator) (event.Subscription, error) {

	logs, sub, err := _INodeManager.contract.WatchLogs(opts, "WithdrawValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(INodeManagerWithdrawValidator)
				if err := _INodeManager.contract.UnpackLog(event, "WithdrawValidator", log); err != nil {
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

// ParseWithdrawValidator is a log parse operation binding the contract event 0xa5bee9e1a2697e6ed6ddf963ed42431f3c8c594f49c697cf1391bbf74e4ea877.
//
// Solidity: event WithdrawValidator(string consensusAddress, string selfStake)
func (_INodeManager *INodeManagerFilterer) ParseWithdrawValidator(log types.Log) (*INodeManagerWithdrawValidator, error) {
	event := new(INodeManagerWithdrawValidator)
	if err := _INodeManager.contract.UnpackLog(event, "WithdrawValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
