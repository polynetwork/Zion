// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package side_chain_manager_abi

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

// ISideChainManagerSideChain is an auto generated low-level Go binding around an user-defined struct.
type ISideChainManagerSideChain struct {
	Owner       common.Address
	ChainID     uint64
	Router      uint64
	Name        string
	CCMCAddress []byte
	ExtraInfo   []byte
}

var (
	MethodApproveQuitSideChain = "approveQuitSideChain"

	MethodApproveRegisterSideChain = "approveRegisterSideChain"

	MethodApproveUpdateSideChain = "approveUpdateSideChain"

	MethodQuitSideChain = "quitSideChain"

	MethodRegisterAsset = "registerAsset"

	MethodRegisterSideChain = "registerSideChain"

	MethodUpdateFee = "updateFee"

	MethodUpdateSideChain = "updateSideChain"

	MethodGetFee = "getFee"

	MethodGetSideChain = "getSideChain"

	EventApproveQuitSideChain = "ApproveQuitSideChain"

	EventApproveRegisterSideChain = "ApproveRegisterSideChain"

	EventApproveUpdateSideChain = "ApproveUpdateSideChain"

	EventQuitSideChain = "QuitSideChain"

	EventRegisterSideChain = "RegisterSideChain"

	EventUpdateSideChain = "UpdateSideChain"
)

// ISideChainManagerABI is the input ABI used to generate the binding from.
const ISideChainManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"}],\"name\":\"ApproveQuitSideChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"}],\"name\":\"ApproveRegisterSideChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"}],\"name\":\"ApproveUpdateSideChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"}],\"name\":\"QuitSideChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"Router\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"}],\"name\":\"RegisterSideChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"Router\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"}],\"name\":\"UpdateSideChain\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"name\":\"approveQuitSideChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"name\":\"approveRegisterSideChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"name\":\"approveUpdateSideChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"name\":\"getFee\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"name\":\"getSideChain\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"router\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"CCMCAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraInfo\",\"type\":\"bytes\"}],\"internalType\":\"structISideChainManager.SideChain\",\"name\":\"sidechain\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"name\":\"quitSideChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"uint64[]\",\"name\":\"AssetMapKey\",\"type\":\"uint64[]\"},{\"internalType\":\"bytes[]\",\"name\":\"AssetMapValue\",\"type\":\"bytes[]\"},{\"internalType\":\"uint64[]\",\"name\":\"LockProxyMapKey\",\"type\":\"uint64[]\"},{\"internalType\":\"bytes[]\",\"name\":\"LockProxyMapValue\",\"type\":\"bytes[]\"}],\"name\":\"registerAsset\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"router\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"CCMCAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraInfo\",\"type\":\"bytes\"}],\"name\":\"registerSideChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"viewNum\",\"type\":\"uint64\"},{\"internalType\":\"int256\",\"name\":\"fee\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"updateFee\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"router\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"CCMCAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraInfo\",\"type\":\"bytes\"}],\"name\":\"updateSideChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ISideChainManagerFuncSigs maps the 4-byte function signature to its string representation.
var ISideChainManagerFuncSigs = map[string]string{
	"9bcb64f0": "approveQuitSideChain(uint64)",
	"c3e7746d": "approveRegisterSideChain(uint64)",
	"678f0135": "approveUpdateSideChain(uint64)",
	"1982b1d0": "getFee(uint64)",
	"84838fb8": "getSideChain(uint64)",
	"78b94ab1": "quitSideChain(uint64)",
	"e171240f": "registerAsset(uint64,uint64[],bytes[],uint64[],bytes[])",
	"3a24101f": "registerSideChain(uint64,uint64,string,bytes,bytes)",
	"db5d3488": "updateFee(uint64,uint64,int256,bytes)",
	"956f1463": "updateSideChain(uint64,uint64,string,bytes,bytes)",
}

// ISideChainManager is an auto generated Go binding around an Ethereum contract.
type ISideChainManager struct {
	ISideChainManagerCaller     // Read-only binding to the contract
	ISideChainManagerTransactor // Write-only binding to the contract
	ISideChainManagerFilterer   // Log filterer for contract events
}

// ISideChainManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ISideChainManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISideChainManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ISideChainManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISideChainManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ISideChainManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISideChainManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ISideChainManagerSession struct {
	Contract     *ISideChainManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ISideChainManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ISideChainManagerCallerSession struct {
	Contract *ISideChainManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ISideChainManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ISideChainManagerTransactorSession struct {
	Contract     *ISideChainManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ISideChainManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ISideChainManagerRaw struct {
	Contract *ISideChainManager // Generic contract binding to access the raw methods on
}

// ISideChainManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ISideChainManagerCallerRaw struct {
	Contract *ISideChainManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ISideChainManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ISideChainManagerTransactorRaw struct {
	Contract *ISideChainManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewISideChainManager creates a new instance of ISideChainManager, bound to a specific deployed contract.
func NewISideChainManager(address common.Address, backend bind.ContractBackend) (*ISideChainManager, error) {
	contract, err := bindISideChainManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ISideChainManager{ISideChainManagerCaller: ISideChainManagerCaller{contract: contract}, ISideChainManagerTransactor: ISideChainManagerTransactor{contract: contract}, ISideChainManagerFilterer: ISideChainManagerFilterer{contract: contract}}, nil
}

// NewISideChainManagerCaller creates a new read-only instance of ISideChainManager, bound to a specific deployed contract.
func NewISideChainManagerCaller(address common.Address, caller bind.ContractCaller) (*ISideChainManagerCaller, error) {
	contract, err := bindISideChainManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ISideChainManagerCaller{contract: contract}, nil
}

// NewISideChainManagerTransactor creates a new write-only instance of ISideChainManager, bound to a specific deployed contract.
func NewISideChainManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ISideChainManagerTransactor, error) {
	contract, err := bindISideChainManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ISideChainManagerTransactor{contract: contract}, nil
}

// NewISideChainManagerFilterer creates a new log filterer instance of ISideChainManager, bound to a specific deployed contract.
func NewISideChainManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ISideChainManagerFilterer, error) {
	contract, err := bindISideChainManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ISideChainManagerFilterer{contract: contract}, nil
}

// bindISideChainManager binds a generic wrapper to an already deployed contract.
func bindISideChainManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ISideChainManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISideChainManager *ISideChainManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISideChainManager.Contract.ISideChainManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISideChainManager *ISideChainManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISideChainManager.Contract.ISideChainManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISideChainManager *ISideChainManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISideChainManager.Contract.ISideChainManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISideChainManager *ISideChainManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISideChainManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISideChainManager *ISideChainManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISideChainManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISideChainManager *ISideChainManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISideChainManager.Contract.contract.Transact(opts, method, params...)
}

// GetFee is a free data retrieval call binding the contract method 0x1982b1d0.
//
// Solidity: function getFee(uint64 chainID) view returns(bytes)
func (_ISideChainManager *ISideChainManagerCaller) GetFee(opts *bind.CallOpts, chainID uint64) ([]byte, error) {
	var out []interface{}
	err := _ISideChainManager.contract.Call(opts, &out, "getFee", chainID)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetFee is a free data retrieval call binding the contract method 0x1982b1d0.
//
// Solidity: function getFee(uint64 chainID) view returns(bytes)
func (_ISideChainManager *ISideChainManagerSession) GetFee(chainID uint64) ([]byte, error) {
	return _ISideChainManager.Contract.GetFee(&_ISideChainManager.CallOpts, chainID)
}

// GetFee is a free data retrieval call binding the contract method 0x1982b1d0.
//
// Solidity: function getFee(uint64 chainID) view returns(bytes)
func (_ISideChainManager *ISideChainManagerCallerSession) GetFee(chainID uint64) ([]byte, error) {
	return _ISideChainManager.Contract.GetFee(&_ISideChainManager.CallOpts, chainID)
}

// GetSideChain is a free data retrieval call binding the contract method 0x84838fb8.
//
// Solidity: function getSideChain(uint64 chainID) view returns((address,uint64,uint64,string,bytes,bytes) sidechain)
func (_ISideChainManager *ISideChainManagerCaller) GetSideChain(opts *bind.CallOpts, chainID uint64) (ISideChainManagerSideChain, error) {
	var out []interface{}
	err := _ISideChainManager.contract.Call(opts, &out, "getSideChain", chainID)

	if err != nil {
		return *new(ISideChainManagerSideChain), err
	}

	out0 := *abi.ConvertType(out[0], new(ISideChainManagerSideChain)).(*ISideChainManagerSideChain)

	return out0, err

}

// GetSideChain is a free data retrieval call binding the contract method 0x84838fb8.
//
// Solidity: function getSideChain(uint64 chainID) view returns((address,uint64,uint64,string,bytes,bytes) sidechain)
func (_ISideChainManager *ISideChainManagerSession) GetSideChain(chainID uint64) (ISideChainManagerSideChain, error) {
	return _ISideChainManager.Contract.GetSideChain(&_ISideChainManager.CallOpts, chainID)
}

// GetSideChain is a free data retrieval call binding the contract method 0x84838fb8.
//
// Solidity: function getSideChain(uint64 chainID) view returns((address,uint64,uint64,string,bytes,bytes) sidechain)
func (_ISideChainManager *ISideChainManagerCallerSession) GetSideChain(chainID uint64) (ISideChainManagerSideChain, error) {
	return _ISideChainManager.Contract.GetSideChain(&_ISideChainManager.CallOpts, chainID)
}

// ApproveQuitSideChain is a paid mutator transaction binding the contract method 0x9bcb64f0.
//
// Solidity: function approveQuitSideChain(uint64 chainID) returns(bool success)
func (_ISideChainManager *ISideChainManagerTransactor) ApproveQuitSideChain(opts *bind.TransactOpts, chainID uint64) (*types.Transaction, error) {
	return _ISideChainManager.contract.Transact(opts, "approveQuitSideChain", chainID)
}

// ApproveQuitSideChain is a paid mutator transaction binding the contract method 0x9bcb64f0.
//
// Solidity: function approveQuitSideChain(uint64 chainID) returns(bool success)
func (_ISideChainManager *ISideChainManagerSession) ApproveQuitSideChain(chainID uint64) (*types.Transaction, error) {
	return _ISideChainManager.Contract.ApproveQuitSideChain(&_ISideChainManager.TransactOpts, chainID)
}

// ApproveQuitSideChain is a paid mutator transaction binding the contract method 0x9bcb64f0.
//
// Solidity: function approveQuitSideChain(uint64 chainID) returns(bool success)
func (_ISideChainManager *ISideChainManagerTransactorSession) ApproveQuitSideChain(chainID uint64) (*types.Transaction, error) {
	return _ISideChainManager.Contract.ApproveQuitSideChain(&_ISideChainManager.TransactOpts, chainID)
}

// ApproveRegisterSideChain is a paid mutator transaction binding the contract method 0xc3e7746d.
//
// Solidity: function approveRegisterSideChain(uint64 chainID) returns(bool success)
func (_ISideChainManager *ISideChainManagerTransactor) ApproveRegisterSideChain(opts *bind.TransactOpts, chainID uint64) (*types.Transaction, error) {
	return _ISideChainManager.contract.Transact(opts, "approveRegisterSideChain", chainID)
}

// ApproveRegisterSideChain is a paid mutator transaction binding the contract method 0xc3e7746d.
//
// Solidity: function approveRegisterSideChain(uint64 chainID) returns(bool success)
func (_ISideChainManager *ISideChainManagerSession) ApproveRegisterSideChain(chainID uint64) (*types.Transaction, error) {
	return _ISideChainManager.Contract.ApproveRegisterSideChain(&_ISideChainManager.TransactOpts, chainID)
}

// ApproveRegisterSideChain is a paid mutator transaction binding the contract method 0xc3e7746d.
//
// Solidity: function approveRegisterSideChain(uint64 chainID) returns(bool success)
func (_ISideChainManager *ISideChainManagerTransactorSession) ApproveRegisterSideChain(chainID uint64) (*types.Transaction, error) {
	return _ISideChainManager.Contract.ApproveRegisterSideChain(&_ISideChainManager.TransactOpts, chainID)
}

// ApproveUpdateSideChain is a paid mutator transaction binding the contract method 0x678f0135.
//
// Solidity: function approveUpdateSideChain(uint64 chainID) returns(bool success)
func (_ISideChainManager *ISideChainManagerTransactor) ApproveUpdateSideChain(opts *bind.TransactOpts, chainID uint64) (*types.Transaction, error) {
	return _ISideChainManager.contract.Transact(opts, "approveUpdateSideChain", chainID)
}

// ApproveUpdateSideChain is a paid mutator transaction binding the contract method 0x678f0135.
//
// Solidity: function approveUpdateSideChain(uint64 chainID) returns(bool success)
func (_ISideChainManager *ISideChainManagerSession) ApproveUpdateSideChain(chainID uint64) (*types.Transaction, error) {
	return _ISideChainManager.Contract.ApproveUpdateSideChain(&_ISideChainManager.TransactOpts, chainID)
}

// ApproveUpdateSideChain is a paid mutator transaction binding the contract method 0x678f0135.
//
// Solidity: function approveUpdateSideChain(uint64 chainID) returns(bool success)
func (_ISideChainManager *ISideChainManagerTransactorSession) ApproveUpdateSideChain(chainID uint64) (*types.Transaction, error) {
	return _ISideChainManager.Contract.ApproveUpdateSideChain(&_ISideChainManager.TransactOpts, chainID)
}

// QuitSideChain is a paid mutator transaction binding the contract method 0x78b94ab1.
//
// Solidity: function quitSideChain(uint64 chainID) returns()
func (_ISideChainManager *ISideChainManagerTransactor) QuitSideChain(opts *bind.TransactOpts, chainID uint64) (*types.Transaction, error) {
	return _ISideChainManager.contract.Transact(opts, "quitSideChain", chainID)
}

// QuitSideChain is a paid mutator transaction binding the contract method 0x78b94ab1.
//
// Solidity: function quitSideChain(uint64 chainID) returns()
func (_ISideChainManager *ISideChainManagerSession) QuitSideChain(chainID uint64) (*types.Transaction, error) {
	return _ISideChainManager.Contract.QuitSideChain(&_ISideChainManager.TransactOpts, chainID)
}

// QuitSideChain is a paid mutator transaction binding the contract method 0x78b94ab1.
//
// Solidity: function quitSideChain(uint64 chainID) returns()
func (_ISideChainManager *ISideChainManagerTransactorSession) QuitSideChain(chainID uint64) (*types.Transaction, error) {
	return _ISideChainManager.Contract.QuitSideChain(&_ISideChainManager.TransactOpts, chainID)
}

// RegisterAsset is a paid mutator transaction binding the contract method 0xe171240f.
//
// Solidity: function registerAsset(uint64 chainID, uint64[] AssetMapKey, bytes[] AssetMapValue, uint64[] LockProxyMapKey, bytes[] LockProxyMapValue) returns(bool success)
func (_ISideChainManager *ISideChainManagerTransactor) RegisterAsset(opts *bind.TransactOpts, chainID uint64, AssetMapKey []uint64, AssetMapValue [][]byte, LockProxyMapKey []uint64, LockProxyMapValue [][]byte) (*types.Transaction, error) {
	return _ISideChainManager.contract.Transact(opts, "registerAsset", chainID, AssetMapKey, AssetMapValue, LockProxyMapKey, LockProxyMapValue)
}

// RegisterAsset is a paid mutator transaction binding the contract method 0xe171240f.
//
// Solidity: function registerAsset(uint64 chainID, uint64[] AssetMapKey, bytes[] AssetMapValue, uint64[] LockProxyMapKey, bytes[] LockProxyMapValue) returns(bool success)
func (_ISideChainManager *ISideChainManagerSession) RegisterAsset(chainID uint64, AssetMapKey []uint64, AssetMapValue [][]byte, LockProxyMapKey []uint64, LockProxyMapValue [][]byte) (*types.Transaction, error) {
	return _ISideChainManager.Contract.RegisterAsset(&_ISideChainManager.TransactOpts, chainID, AssetMapKey, AssetMapValue, LockProxyMapKey, LockProxyMapValue)
}

// RegisterAsset is a paid mutator transaction binding the contract method 0xe171240f.
//
// Solidity: function registerAsset(uint64 chainID, uint64[] AssetMapKey, bytes[] AssetMapValue, uint64[] LockProxyMapKey, bytes[] LockProxyMapValue) returns(bool success)
func (_ISideChainManager *ISideChainManagerTransactorSession) RegisterAsset(chainID uint64, AssetMapKey []uint64, AssetMapValue [][]byte, LockProxyMapKey []uint64, LockProxyMapValue [][]byte) (*types.Transaction, error) {
	return _ISideChainManager.Contract.RegisterAsset(&_ISideChainManager.TransactOpts, chainID, AssetMapKey, AssetMapValue, LockProxyMapKey, LockProxyMapValue)
}

// RegisterSideChain is a paid mutator transaction binding the contract method 0x3a24101f.
//
// Solidity: function registerSideChain(uint64 chainID, uint64 router, string name, bytes CCMCAddress, bytes extraInfo) returns()
func (_ISideChainManager *ISideChainManagerTransactor) RegisterSideChain(opts *bind.TransactOpts, chainID uint64, router uint64, name string, CCMCAddress []byte, extraInfo []byte) (*types.Transaction, error) {
	return _ISideChainManager.contract.Transact(opts, "registerSideChain", chainID, router, name, CCMCAddress, extraInfo)
}

// RegisterSideChain is a paid mutator transaction binding the contract method 0x3a24101f.
//
// Solidity: function registerSideChain(uint64 chainID, uint64 router, string name, bytes CCMCAddress, bytes extraInfo) returns()
func (_ISideChainManager *ISideChainManagerSession) RegisterSideChain(chainID uint64, router uint64, name string, CCMCAddress []byte, extraInfo []byte) (*types.Transaction, error) {
	return _ISideChainManager.Contract.RegisterSideChain(&_ISideChainManager.TransactOpts, chainID, router, name, CCMCAddress, extraInfo)
}

// RegisterSideChain is a paid mutator transaction binding the contract method 0x3a24101f.
//
// Solidity: function registerSideChain(uint64 chainID, uint64 router, string name, bytes CCMCAddress, bytes extraInfo) returns()
func (_ISideChainManager *ISideChainManagerTransactorSession) RegisterSideChain(chainID uint64, router uint64, name string, CCMCAddress []byte, extraInfo []byte) (*types.Transaction, error) {
	return _ISideChainManager.Contract.RegisterSideChain(&_ISideChainManager.TransactOpts, chainID, router, name, CCMCAddress, extraInfo)
}

// UpdateFee is a paid mutator transaction binding the contract method 0xdb5d3488.
//
// Solidity: function updateFee(uint64 chainID, uint64 viewNum, int256 fee, bytes signature) returns(bool success)
func (_ISideChainManager *ISideChainManagerTransactor) UpdateFee(opts *bind.TransactOpts, chainID uint64, viewNum uint64, fee *big.Int, signature []byte) (*types.Transaction, error) {
	return _ISideChainManager.contract.Transact(opts, "updateFee", chainID, viewNum, fee, signature)
}

// UpdateFee is a paid mutator transaction binding the contract method 0xdb5d3488.
//
// Solidity: function updateFee(uint64 chainID, uint64 viewNum, int256 fee, bytes signature) returns(bool success)
func (_ISideChainManager *ISideChainManagerSession) UpdateFee(chainID uint64, viewNum uint64, fee *big.Int, signature []byte) (*types.Transaction, error) {
	return _ISideChainManager.Contract.UpdateFee(&_ISideChainManager.TransactOpts, chainID, viewNum, fee, signature)
}

// UpdateFee is a paid mutator transaction binding the contract method 0xdb5d3488.
//
// Solidity: function updateFee(uint64 chainID, uint64 viewNum, int256 fee, bytes signature) returns(bool success)
func (_ISideChainManager *ISideChainManagerTransactorSession) UpdateFee(chainID uint64, viewNum uint64, fee *big.Int, signature []byte) (*types.Transaction, error) {
	return _ISideChainManager.Contract.UpdateFee(&_ISideChainManager.TransactOpts, chainID, viewNum, fee, signature)
}

// UpdateSideChain is a paid mutator transaction binding the contract method 0x956f1463.
//
// Solidity: function updateSideChain(uint64 chainID, uint64 router, string name, bytes CCMCAddress, bytes extraInfo) returns()
func (_ISideChainManager *ISideChainManagerTransactor) UpdateSideChain(opts *bind.TransactOpts, chainID uint64, router uint64, name string, CCMCAddress []byte, extraInfo []byte) (*types.Transaction, error) {
	return _ISideChainManager.contract.Transact(opts, "updateSideChain", chainID, router, name, CCMCAddress, extraInfo)
}

// UpdateSideChain is a paid mutator transaction binding the contract method 0x956f1463.
//
// Solidity: function updateSideChain(uint64 chainID, uint64 router, string name, bytes CCMCAddress, bytes extraInfo) returns()
func (_ISideChainManager *ISideChainManagerSession) UpdateSideChain(chainID uint64, router uint64, name string, CCMCAddress []byte, extraInfo []byte) (*types.Transaction, error) {
	return _ISideChainManager.Contract.UpdateSideChain(&_ISideChainManager.TransactOpts, chainID, router, name, CCMCAddress, extraInfo)
}

// UpdateSideChain is a paid mutator transaction binding the contract method 0x956f1463.
//
// Solidity: function updateSideChain(uint64 chainID, uint64 router, string name, bytes CCMCAddress, bytes extraInfo) returns()
func (_ISideChainManager *ISideChainManagerTransactorSession) UpdateSideChain(chainID uint64, router uint64, name string, CCMCAddress []byte, extraInfo []byte) (*types.Transaction, error) {
	return _ISideChainManager.Contract.UpdateSideChain(&_ISideChainManager.TransactOpts, chainID, router, name, CCMCAddress, extraInfo)
}

// ISideChainManagerApproveQuitSideChainIterator is returned from FilterApproveQuitSideChain and is used to iterate over the raw logs and unpacked data for ApproveQuitSideChain events raised by the ISideChainManager contract.
type ISideChainManagerApproveQuitSideChainIterator struct {
	Event *ISideChainManagerApproveQuitSideChain // Event containing the contract specifics and raw log

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
func (it *ISideChainManagerApproveQuitSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ISideChainManagerApproveQuitSideChain)
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
		it.Event = new(ISideChainManagerApproveQuitSideChain)
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
func (it *ISideChainManagerApproveQuitSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ISideChainManagerApproveQuitSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ISideChainManagerApproveQuitSideChain represents a ApproveQuitSideChain event raised by the ISideChainManager contract.
type ISideChainManagerApproveQuitSideChain struct {
	ChainId uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproveQuitSideChain is a free log retrieval operation binding the contract event 0x12d05899d2cf3e2ea805d35769f340779fdfa004e8b2d9342a47eb158d276e73.
//
// Solidity: event ApproveQuitSideChain(uint64 ChainId)
func (_ISideChainManager *ISideChainManagerFilterer) FilterApproveQuitSideChain(opts *bind.FilterOpts) (*ISideChainManagerApproveQuitSideChainIterator, error) {

	logs, sub, err := _ISideChainManager.contract.FilterLogs(opts, "ApproveQuitSideChain")
	if err != nil {
		return nil, err
	}
	return &ISideChainManagerApproveQuitSideChainIterator{contract: _ISideChainManager.contract, event: "ApproveQuitSideChain", logs: logs, sub: sub}, nil
}

// WatchApproveQuitSideChain is a free log subscription operation binding the contract event 0x12d05899d2cf3e2ea805d35769f340779fdfa004e8b2d9342a47eb158d276e73.
//
// Solidity: event ApproveQuitSideChain(uint64 ChainId)
func (_ISideChainManager *ISideChainManagerFilterer) WatchApproveQuitSideChain(opts *bind.WatchOpts, sink chan<- *ISideChainManagerApproveQuitSideChain) (event.Subscription, error) {

	logs, sub, err := _ISideChainManager.contract.WatchLogs(opts, "ApproveQuitSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ISideChainManagerApproveQuitSideChain)
				if err := _ISideChainManager.contract.UnpackLog(event, "ApproveQuitSideChain", log); err != nil {
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

// ParseApproveQuitSideChain is a log parse operation binding the contract event 0x12d05899d2cf3e2ea805d35769f340779fdfa004e8b2d9342a47eb158d276e73.
//
// Solidity: event ApproveQuitSideChain(uint64 ChainId)
func (_ISideChainManager *ISideChainManagerFilterer) ParseApproveQuitSideChain(log types.Log) (*ISideChainManagerApproveQuitSideChain, error) {
	event := new(ISideChainManagerApproveQuitSideChain)
	if err := _ISideChainManager.contract.UnpackLog(event, "ApproveQuitSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ISideChainManagerApproveRegisterSideChainIterator is returned from FilterApproveRegisterSideChain and is used to iterate over the raw logs and unpacked data for ApproveRegisterSideChain events raised by the ISideChainManager contract.
type ISideChainManagerApproveRegisterSideChainIterator struct {
	Event *ISideChainManagerApproveRegisterSideChain // Event containing the contract specifics and raw log

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
func (it *ISideChainManagerApproveRegisterSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ISideChainManagerApproveRegisterSideChain)
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
		it.Event = new(ISideChainManagerApproveRegisterSideChain)
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
func (it *ISideChainManagerApproveRegisterSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ISideChainManagerApproveRegisterSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ISideChainManagerApproveRegisterSideChain represents a ApproveRegisterSideChain event raised by the ISideChainManager contract.
type ISideChainManagerApproveRegisterSideChain struct {
	ChainId uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproveRegisterSideChain is a free log retrieval operation binding the contract event 0x7f80ce991f1aef1de58b0a4d15734c702a491b07b594a2655503a5e433fd5749.
//
// Solidity: event ApproveRegisterSideChain(uint64 ChainId)
func (_ISideChainManager *ISideChainManagerFilterer) FilterApproveRegisterSideChain(opts *bind.FilterOpts) (*ISideChainManagerApproveRegisterSideChainIterator, error) {

	logs, sub, err := _ISideChainManager.contract.FilterLogs(opts, "ApproveRegisterSideChain")
	if err != nil {
		return nil, err
	}
	return &ISideChainManagerApproveRegisterSideChainIterator{contract: _ISideChainManager.contract, event: "ApproveRegisterSideChain", logs: logs, sub: sub}, nil
}

// WatchApproveRegisterSideChain is a free log subscription operation binding the contract event 0x7f80ce991f1aef1de58b0a4d15734c702a491b07b594a2655503a5e433fd5749.
//
// Solidity: event ApproveRegisterSideChain(uint64 ChainId)
func (_ISideChainManager *ISideChainManagerFilterer) WatchApproveRegisterSideChain(opts *bind.WatchOpts, sink chan<- *ISideChainManagerApproveRegisterSideChain) (event.Subscription, error) {

	logs, sub, err := _ISideChainManager.contract.WatchLogs(opts, "ApproveRegisterSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ISideChainManagerApproveRegisterSideChain)
				if err := _ISideChainManager.contract.UnpackLog(event, "ApproveRegisterSideChain", log); err != nil {
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

// ParseApproveRegisterSideChain is a log parse operation binding the contract event 0x7f80ce991f1aef1de58b0a4d15734c702a491b07b594a2655503a5e433fd5749.
//
// Solidity: event ApproveRegisterSideChain(uint64 ChainId)
func (_ISideChainManager *ISideChainManagerFilterer) ParseApproveRegisterSideChain(log types.Log) (*ISideChainManagerApproveRegisterSideChain, error) {
	event := new(ISideChainManagerApproveRegisterSideChain)
	if err := _ISideChainManager.contract.UnpackLog(event, "ApproveRegisterSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ISideChainManagerApproveUpdateSideChainIterator is returned from FilterApproveUpdateSideChain and is used to iterate over the raw logs and unpacked data for ApproveUpdateSideChain events raised by the ISideChainManager contract.
type ISideChainManagerApproveUpdateSideChainIterator struct {
	Event *ISideChainManagerApproveUpdateSideChain // Event containing the contract specifics and raw log

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
func (it *ISideChainManagerApproveUpdateSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ISideChainManagerApproveUpdateSideChain)
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
		it.Event = new(ISideChainManagerApproveUpdateSideChain)
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
func (it *ISideChainManagerApproveUpdateSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ISideChainManagerApproveUpdateSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ISideChainManagerApproveUpdateSideChain represents a ApproveUpdateSideChain event raised by the ISideChainManager contract.
type ISideChainManagerApproveUpdateSideChain struct {
	ChainId uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproveUpdateSideChain is a free log retrieval operation binding the contract event 0x24eae46246c3dc63bc634070751e269a353b630665bfc8fbc057b614629e8136.
//
// Solidity: event ApproveUpdateSideChain(uint64 ChainId)
func (_ISideChainManager *ISideChainManagerFilterer) FilterApproveUpdateSideChain(opts *bind.FilterOpts) (*ISideChainManagerApproveUpdateSideChainIterator, error) {

	logs, sub, err := _ISideChainManager.contract.FilterLogs(opts, "ApproveUpdateSideChain")
	if err != nil {
		return nil, err
	}
	return &ISideChainManagerApproveUpdateSideChainIterator{contract: _ISideChainManager.contract, event: "ApproveUpdateSideChain", logs: logs, sub: sub}, nil
}

// WatchApproveUpdateSideChain is a free log subscription operation binding the contract event 0x24eae46246c3dc63bc634070751e269a353b630665bfc8fbc057b614629e8136.
//
// Solidity: event ApproveUpdateSideChain(uint64 ChainId)
func (_ISideChainManager *ISideChainManagerFilterer) WatchApproveUpdateSideChain(opts *bind.WatchOpts, sink chan<- *ISideChainManagerApproveUpdateSideChain) (event.Subscription, error) {

	logs, sub, err := _ISideChainManager.contract.WatchLogs(opts, "ApproveUpdateSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ISideChainManagerApproveUpdateSideChain)
				if err := _ISideChainManager.contract.UnpackLog(event, "ApproveUpdateSideChain", log); err != nil {
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

// ParseApproveUpdateSideChain is a log parse operation binding the contract event 0x24eae46246c3dc63bc634070751e269a353b630665bfc8fbc057b614629e8136.
//
// Solidity: event ApproveUpdateSideChain(uint64 ChainId)
func (_ISideChainManager *ISideChainManagerFilterer) ParseApproveUpdateSideChain(log types.Log) (*ISideChainManagerApproveUpdateSideChain, error) {
	event := new(ISideChainManagerApproveUpdateSideChain)
	if err := _ISideChainManager.contract.UnpackLog(event, "ApproveUpdateSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ISideChainManagerQuitSideChainIterator is returned from FilterQuitSideChain and is used to iterate over the raw logs and unpacked data for QuitSideChain events raised by the ISideChainManager contract.
type ISideChainManagerQuitSideChainIterator struct {
	Event *ISideChainManagerQuitSideChain // Event containing the contract specifics and raw log

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
func (it *ISideChainManagerQuitSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ISideChainManagerQuitSideChain)
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
		it.Event = new(ISideChainManagerQuitSideChain)
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
func (it *ISideChainManagerQuitSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ISideChainManagerQuitSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ISideChainManagerQuitSideChain represents a QuitSideChain event raised by the ISideChainManager contract.
type ISideChainManagerQuitSideChain struct {
	ChainId uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterQuitSideChain is a free log retrieval operation binding the contract event 0xd5e9edc8ae17c077aca5871ac2653f2acb7fb85796fb7c5b43f5ea907c198e90.
//
// Solidity: event QuitSideChain(uint64 ChainId)
func (_ISideChainManager *ISideChainManagerFilterer) FilterQuitSideChain(opts *bind.FilterOpts) (*ISideChainManagerQuitSideChainIterator, error) {

	logs, sub, err := _ISideChainManager.contract.FilterLogs(opts, "QuitSideChain")
	if err != nil {
		return nil, err
	}
	return &ISideChainManagerQuitSideChainIterator{contract: _ISideChainManager.contract, event: "QuitSideChain", logs: logs, sub: sub}, nil
}

// WatchQuitSideChain is a free log subscription operation binding the contract event 0xd5e9edc8ae17c077aca5871ac2653f2acb7fb85796fb7c5b43f5ea907c198e90.
//
// Solidity: event QuitSideChain(uint64 ChainId)
func (_ISideChainManager *ISideChainManagerFilterer) WatchQuitSideChain(opts *bind.WatchOpts, sink chan<- *ISideChainManagerQuitSideChain) (event.Subscription, error) {

	logs, sub, err := _ISideChainManager.contract.WatchLogs(opts, "QuitSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ISideChainManagerQuitSideChain)
				if err := _ISideChainManager.contract.UnpackLog(event, "QuitSideChain", log); err != nil {
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

// ParseQuitSideChain is a log parse operation binding the contract event 0xd5e9edc8ae17c077aca5871ac2653f2acb7fb85796fb7c5b43f5ea907c198e90.
//
// Solidity: event QuitSideChain(uint64 ChainId)
func (_ISideChainManager *ISideChainManagerFilterer) ParseQuitSideChain(log types.Log) (*ISideChainManagerQuitSideChain, error) {
	event := new(ISideChainManagerQuitSideChain)
	if err := _ISideChainManager.contract.UnpackLog(event, "QuitSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ISideChainManagerRegisterSideChainIterator is returned from FilterRegisterSideChain and is used to iterate over the raw logs and unpacked data for RegisterSideChain events raised by the ISideChainManager contract.
type ISideChainManagerRegisterSideChainIterator struct {
	Event *ISideChainManagerRegisterSideChain // Event containing the contract specifics and raw log

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
func (it *ISideChainManagerRegisterSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ISideChainManagerRegisterSideChain)
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
		it.Event = new(ISideChainManagerRegisterSideChain)
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
func (it *ISideChainManagerRegisterSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ISideChainManagerRegisterSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ISideChainManagerRegisterSideChain represents a RegisterSideChain event raised by the ISideChainManager contract.
type ISideChainManagerRegisterSideChain struct {
	ChainId uint64
	Router  uint64
	Name    string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRegisterSideChain is a free log retrieval operation binding the contract event 0x7688735bc7c43b9ea9af0a45be0fc37526cf866360acaf0c29c44185f8eee791.
//
// Solidity: event RegisterSideChain(uint64 ChainId, uint64 Router, string Name)
func (_ISideChainManager *ISideChainManagerFilterer) FilterRegisterSideChain(opts *bind.FilterOpts) (*ISideChainManagerRegisterSideChainIterator, error) {

	logs, sub, err := _ISideChainManager.contract.FilterLogs(opts, "RegisterSideChain")
	if err != nil {
		return nil, err
	}
	return &ISideChainManagerRegisterSideChainIterator{contract: _ISideChainManager.contract, event: "RegisterSideChain", logs: logs, sub: sub}, nil
}

// WatchRegisterSideChain is a free log subscription operation binding the contract event 0x7688735bc7c43b9ea9af0a45be0fc37526cf866360acaf0c29c44185f8eee791.
//
// Solidity: event RegisterSideChain(uint64 ChainId, uint64 Router, string Name)
func (_ISideChainManager *ISideChainManagerFilterer) WatchRegisterSideChain(opts *bind.WatchOpts, sink chan<- *ISideChainManagerRegisterSideChain) (event.Subscription, error) {

	logs, sub, err := _ISideChainManager.contract.WatchLogs(opts, "RegisterSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ISideChainManagerRegisterSideChain)
				if err := _ISideChainManager.contract.UnpackLog(event, "RegisterSideChain", log); err != nil {
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

// ParseRegisterSideChain is a log parse operation binding the contract event 0x7688735bc7c43b9ea9af0a45be0fc37526cf866360acaf0c29c44185f8eee791.
//
// Solidity: event RegisterSideChain(uint64 ChainId, uint64 Router, string Name)
func (_ISideChainManager *ISideChainManagerFilterer) ParseRegisterSideChain(log types.Log) (*ISideChainManagerRegisterSideChain, error) {
	event := new(ISideChainManagerRegisterSideChain)
	if err := _ISideChainManager.contract.UnpackLog(event, "RegisterSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ISideChainManagerUpdateSideChainIterator is returned from FilterUpdateSideChain and is used to iterate over the raw logs and unpacked data for UpdateSideChain events raised by the ISideChainManager contract.
type ISideChainManagerUpdateSideChainIterator struct {
	Event *ISideChainManagerUpdateSideChain // Event containing the contract specifics and raw log

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
func (it *ISideChainManagerUpdateSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ISideChainManagerUpdateSideChain)
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
		it.Event = new(ISideChainManagerUpdateSideChain)
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
func (it *ISideChainManagerUpdateSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ISideChainManagerUpdateSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ISideChainManagerUpdateSideChain represents a UpdateSideChain event raised by the ISideChainManager contract.
type ISideChainManagerUpdateSideChain struct {
	ChainId uint64
	Router  uint64
	Name    string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUpdateSideChain is a free log retrieval operation binding the contract event 0xe5e51843a978f17f0cb363b30b9592c76e998342e674c98164d970f8c7fe1201.
//
// Solidity: event UpdateSideChain(uint64 ChainId, uint64 Router, string Name)
func (_ISideChainManager *ISideChainManagerFilterer) FilterUpdateSideChain(opts *bind.FilterOpts) (*ISideChainManagerUpdateSideChainIterator, error) {

	logs, sub, err := _ISideChainManager.contract.FilterLogs(opts, "UpdateSideChain")
	if err != nil {
		return nil, err
	}
	return &ISideChainManagerUpdateSideChainIterator{contract: _ISideChainManager.contract, event: "UpdateSideChain", logs: logs, sub: sub}, nil
}

// WatchUpdateSideChain is a free log subscription operation binding the contract event 0xe5e51843a978f17f0cb363b30b9592c76e998342e674c98164d970f8c7fe1201.
//
// Solidity: event UpdateSideChain(uint64 ChainId, uint64 Router, string Name)
func (_ISideChainManager *ISideChainManagerFilterer) WatchUpdateSideChain(opts *bind.WatchOpts, sink chan<- *ISideChainManagerUpdateSideChain) (event.Subscription, error) {

	logs, sub, err := _ISideChainManager.contract.WatchLogs(opts, "UpdateSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ISideChainManagerUpdateSideChain)
				if err := _ISideChainManager.contract.UnpackLog(event, "UpdateSideChain", log); err != nil {
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

// ParseUpdateSideChain is a log parse operation binding the contract event 0xe5e51843a978f17f0cb363b30b9592c76e998342e674c98164d970f8c7fe1201.
//
// Solidity: event UpdateSideChain(uint64 ChainId, uint64 Router, string Name)
func (_ISideChainManager *ISideChainManagerFilterer) ParseUpdateSideChain(log types.Log) (*ISideChainManagerUpdateSideChain, error) {
	event := new(ISideChainManagerUpdateSideChain)
	if err := _ISideChainManager.contract.UnpackLog(event, "UpdateSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

