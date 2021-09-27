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

// side_chain_managerBtcTxParamDetial is an auto generated low-level Go binding around an user-defined struct.
type side_chain_managerBtcTxParamDetial struct {
	PVersion  uint64
	FeeRate   uint64
	MinChange uint64
}

var (
	MethodApproveQuitSideChain = "approveQuitSideChain"

	MethodApproveRegisterSideChain = "approveRegisterSideChain"

	MethodApproveUpdateSideChain = "approveUpdateSideChain"

	MethodName = "name"

	MethodQuitSideChain = "quitSideChain"

	MethodRegisterRedeem = "registerRedeem"

	MethodRegisterSideChain = "registerSideChain"

	MethodSetBtcTxParam = "setBtcTxParam"

	MethodUpdateSideChain = "updateSideChain"
)

// SideChainManagerABI is the input ABI used to generate the binding from.
const SideChainManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"}],\"name\":\"evtApproveQuitSideChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"}],\"name\":\"evtApproveRegisterSideChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"}],\"name\":\"evtApproveUpdateSideChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"}],\"name\":\"evtQuitSideChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"rk\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ContractAddress\",\"type\":\"string\"}],\"name\":\"evtRegisterRedeem\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"Router\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"BlocksToWait\",\"type\":\"uint64\"}],\"name\":\"evtRegisterSideChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"rk\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"RedeemChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"FeeRate\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"MinChange\",\"type\":\"uint64\"}],\"name\":\"evtSetBtcTxParam\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"Router\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"BlocksToWait\",\"type\":\"uint64\"}],\"name\":\"evtUpdateSideChain\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"Chainid\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"approveQuitSideChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"Chainid\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"approveRegisterSideChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"Chainid\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"approveUpdateSideChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"Chainid\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"quitSideChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"RedeemChainID\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"ContractChainID\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"Redeem\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"CVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"ContractAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"Signs\",\"type\":\"bytes[]\"}],\"name\":\"registerRedeem\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"Router\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"BlocksToWait\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"CCMCAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"ExtraInfo\",\"type\":\"bytes\"}],\"name\":\"registerSideChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"Redeem\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"RedeemChainId\",\"type\":\"uint64\"},{\"internalType\":\"bytes[]\",\"name\":\"Sigs\",\"type\":\"bytes[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"PVersion\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"FeeRate\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"MinChange\",\"type\":\"uint64\"}],\"internalType\":\"structside_chain_manager.BtcTxParamDetial\",\"name\":\"Detial\",\"type\":\"tuple\"}],\"name\":\"setBtcTxParam\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"Router\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"BlocksToWait\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"CCMCAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"ExtraInfo\",\"type\":\"bytes\"}],\"name\":\"updateSideChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SideChainManagerFuncSigs maps the 4-byte function signature to its string representation.
var SideChainManagerFuncSigs = map[string]string{
	"6c8ac5c1": "approveQuitSideChain(uint64,address)",
	"65764e16": "approveRegisterSideChain(uint64,address)",
	"805b508e": "approveUpdateSideChain(uint64,address)",
	"06fdde03": "name()",
	"7460736e": "quitSideChain(uint64,address)",
	"33e1d41a": "registerRedeem(uint64,uint64,bytes,uint64,bytes,bytes[])",
	"ab7a2037": "registerSideChain(address,uint64,uint64,string,uint64,bytes,bytes)",
	"ee9891e3": "setBtcTxParam(bytes,uint64,bytes[],(uint64,uint64,uint64))",
	"f7782f81": "updateSideChain(address,uint64,uint64,string,uint64,bytes,bytes)",
}

// SideChainManagerBin is the compiled bytecode used for deploying new contracts.
var SideChainManagerBin = "0x608060405234801561001057600080fd5b50610608806100206000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c80637460736e116100665780637460736e146100da578063805b508e146100da578063ab7a2037146100f0578063ee9891e31461010b578063f7782f81146100f057600080fd5b806306fdde031461009857806333e1d41a146100b057806365764e16146100da5780636c8ac5c1146100da575b600080fd5b60606040516100a7919061050d565b60405180910390f35b6100ca6100be366004610454565b60009695505050505050565b60405190151581526020016100a7565b6100ca6100e8366004610421565b600092915050565b6100ca6100fe36600461027a565b6000979650505050505050565b6100ca61011936600461035c565b6000949350505050565b600067ffffffffffffffff83111561013d5761013d6105bc565b610150601f8401601f191660200161058b565b905082815283838301111561016457600080fd5b828260208301376000602084830101529392505050565b80356001600160a01b038116811461019257600080fd5b919050565b600082601f8301126101a857600080fd5b8135602067ffffffffffffffff808311156101c5576101c56105bc565b8260051b6101d483820161058b565b8481528381019087850183890186018a10156101ef57600080fd5b60009350835b8781101561022c5781358681111561020b578586fd5b6102198c89838e010161023b565b85525092860192908601906001016101f5565b50909998505050505050505050565b600082601f83011261024c57600080fd5b61025b83833560208501610123565b9392505050565b803567ffffffffffffffff8116811461019257600080fd5b600080600080600080600060e0888a03121561029557600080fd5b61029e8861017b565b96506102ac60208901610262565b95506102ba60408901610262565b9450606088013567ffffffffffffffff808211156102d757600080fd5b818a0191508a601f8301126102eb57600080fd5b6102fa8b833560208501610123565b955061030860808b01610262565b945060a08a013591508082111561031e57600080fd5b61032a8b838c0161023b565b935060c08a013591508082111561034057600080fd5b5061034d8a828b0161023b565b91505092959891949750929550565b60008060008084860360c081121561037357600080fd5b853567ffffffffffffffff8082111561038b57600080fd5b61039789838a0161023b565b96506103a560208901610262565b955060408801359150808211156103bb57600080fd5b506103c888828901610197565b9350506060605f19820112156103dd57600080fd5b506103e6610562565b6103f260608701610262565b815261040060808701610262565b602082015261041160a08701610262565b6040820152939692955090935050565b6000806040838503121561043457600080fd5b61043d83610262565b915061044b6020840161017b565b90509250929050565b60008060008060008060c0878903121561046d57600080fd5b61047687610262565b955061048460208801610262565b9450604087013567ffffffffffffffff808211156104a157600080fd5b6104ad8a838b0161023b565b95506104bb60608a01610262565b945060808901359150808211156104d157600080fd5b6104dd8a838b0161023b565b935060a08901359150808211156104f357600080fd5b5061050089828a01610197565b9150509295509295509295565b600060208083528351808285015260005b8181101561053a5785810183015185820160400152820161051e565b8181111561054c576000604083870101525b50601f01601f1916929092016040019392505050565b6040516060810167ffffffffffffffff81118282101715610585576105856105bc565b60405290565b604051601f8201601f1916810167ffffffffffffffff811182821017156105b4576105b46105bc565b604052919050565b634e487b7160e01b600052604160045260246000fdfea26469706673582212201e30dd90c8f8c2bc2a8d03860a8d81c9f32d74209e8ba33de259f3bdaaaa078064736f6c63430008060033"

// DeploySideChainManager deploys a new Ethereum contract, binding an instance of SideChainManager to it.
func DeploySideChainManager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SideChainManager, error) {
	parsed, err := abi.JSON(strings.NewReader(SideChainManagerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SideChainManagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SideChainManager{SideChainManagerCaller: SideChainManagerCaller{contract: contract}, SideChainManagerTransactor: SideChainManagerTransactor{contract: contract}, SideChainManagerFilterer: SideChainManagerFilterer{contract: contract}}, nil
}

// SideChainManager is an auto generated Go binding around an Ethereum contract.
type SideChainManager struct {
	SideChainManagerCaller     // Read-only binding to the contract
	SideChainManagerTransactor // Write-only binding to the contract
	SideChainManagerFilterer   // Log filterer for contract events
}

// SideChainManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type SideChainManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SideChainManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SideChainManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SideChainManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SideChainManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SideChainManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SideChainManagerSession struct {
	Contract     *SideChainManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SideChainManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SideChainManagerCallerSession struct {
	Contract *SideChainManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// SideChainManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SideChainManagerTransactorSession struct {
	Contract     *SideChainManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// SideChainManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type SideChainManagerRaw struct {
	Contract *SideChainManager // Generic contract binding to access the raw methods on
}

// SideChainManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SideChainManagerCallerRaw struct {
	Contract *SideChainManagerCaller // Generic read-only contract binding to access the raw methods on
}

// SideChainManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SideChainManagerTransactorRaw struct {
	Contract *SideChainManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSideChainManager creates a new instance of SideChainManager, bound to a specific deployed contract.
func NewSideChainManager(address common.Address, backend bind.ContractBackend) (*SideChainManager, error) {
	contract, err := bindSideChainManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SideChainManager{SideChainManagerCaller: SideChainManagerCaller{contract: contract}, SideChainManagerTransactor: SideChainManagerTransactor{contract: contract}, SideChainManagerFilterer: SideChainManagerFilterer{contract: contract}}, nil
}

// NewSideChainManagerCaller creates a new read-only instance of SideChainManager, bound to a specific deployed contract.
func NewSideChainManagerCaller(address common.Address, caller bind.ContractCaller) (*SideChainManagerCaller, error) {
	contract, err := bindSideChainManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SideChainManagerCaller{contract: contract}, nil
}

// NewSideChainManagerTransactor creates a new write-only instance of SideChainManager, bound to a specific deployed contract.
func NewSideChainManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*SideChainManagerTransactor, error) {
	contract, err := bindSideChainManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SideChainManagerTransactor{contract: contract}, nil
}

// NewSideChainManagerFilterer creates a new log filterer instance of SideChainManager, bound to a specific deployed contract.
func NewSideChainManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*SideChainManagerFilterer, error) {
	contract, err := bindSideChainManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SideChainManagerFilterer{contract: contract}, nil
}

// bindSideChainManager binds a generic wrapper to an already deployed contract.
func bindSideChainManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SideChainManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SideChainManager *SideChainManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SideChainManager.Contract.SideChainManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SideChainManager *SideChainManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SideChainManager.Contract.SideChainManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SideChainManager *SideChainManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SideChainManager.Contract.SideChainManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SideChainManager *SideChainManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SideChainManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SideChainManager *SideChainManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SideChainManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SideChainManager *SideChainManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SideChainManager.Contract.contract.Transact(opts, method, params...)
}

// ApproveQuitSideChain is a paid mutator transaction binding the contract method 0x6c8ac5c1.
//
// Solidity: function approveQuitSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) ApproveQuitSideChain(opts *bind.TransactOpts, Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "approveQuitSideChain", Chainid, Address)
}

// ApproveQuitSideChain is a paid mutator transaction binding the contract method 0x6c8ac5c1.
//
// Solidity: function approveQuitSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerSession) ApproveQuitSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.ApproveQuitSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// ApproveQuitSideChain is a paid mutator transaction binding the contract method 0x6c8ac5c1.
//
// Solidity: function approveQuitSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) ApproveQuitSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.ApproveQuitSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// ApproveRegisterSideChain is a paid mutator transaction binding the contract method 0x65764e16.
//
// Solidity: function approveRegisterSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) ApproveRegisterSideChain(opts *bind.TransactOpts, Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "approveRegisterSideChain", Chainid, Address)
}

// ApproveRegisterSideChain is a paid mutator transaction binding the contract method 0x65764e16.
//
// Solidity: function approveRegisterSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerSession) ApproveRegisterSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.ApproveRegisterSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// ApproveRegisterSideChain is a paid mutator transaction binding the contract method 0x65764e16.
//
// Solidity: function approveRegisterSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) ApproveRegisterSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.ApproveRegisterSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// ApproveUpdateSideChain is a paid mutator transaction binding the contract method 0x805b508e.
//
// Solidity: function approveUpdateSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) ApproveUpdateSideChain(opts *bind.TransactOpts, Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "approveUpdateSideChain", Chainid, Address)
}

// ApproveUpdateSideChain is a paid mutator transaction binding the contract method 0x805b508e.
//
// Solidity: function approveUpdateSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerSession) ApproveUpdateSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.ApproveUpdateSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// ApproveUpdateSideChain is a paid mutator transaction binding the contract method 0x805b508e.
//
// Solidity: function approveUpdateSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) ApproveUpdateSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.ApproveUpdateSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string Name)
func (_SideChainManager *SideChainManagerTransactor) Name(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "name")
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string Name)
func (_SideChainManager *SideChainManagerSession) Name() (*types.Transaction, error) {
	return _SideChainManager.Contract.Name(&_SideChainManager.TransactOpts)
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string Name)
func (_SideChainManager *SideChainManagerTransactorSession) Name() (*types.Transaction, error) {
	return _SideChainManager.Contract.Name(&_SideChainManager.TransactOpts)
}

// QuitSideChain is a paid mutator transaction binding the contract method 0x7460736e.
//
// Solidity: function quitSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) QuitSideChain(opts *bind.TransactOpts, Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "quitSideChain", Chainid, Address)
}

// QuitSideChain is a paid mutator transaction binding the contract method 0x7460736e.
//
// Solidity: function quitSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerSession) QuitSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.QuitSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// QuitSideChain is a paid mutator transaction binding the contract method 0x7460736e.
//
// Solidity: function quitSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) QuitSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.QuitSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// RegisterRedeem is a paid mutator transaction binding the contract method 0x33e1d41a.
//
// Solidity: function registerRedeem(uint64 RedeemChainID, uint64 ContractChainID, bytes Redeem, uint64 CVersion, bytes ContractAddress, bytes[] Signs) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) RegisterRedeem(opts *bind.TransactOpts, RedeemChainID uint64, ContractChainID uint64, Redeem []byte, CVersion uint64, ContractAddress []byte, Signs [][]byte) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "registerRedeem", RedeemChainID, ContractChainID, Redeem, CVersion, ContractAddress, Signs)
}

// RegisterRedeem is a paid mutator transaction binding the contract method 0x33e1d41a.
//
// Solidity: function registerRedeem(uint64 RedeemChainID, uint64 ContractChainID, bytes Redeem, uint64 CVersion, bytes ContractAddress, bytes[] Signs) returns(bool success)
func (_SideChainManager *SideChainManagerSession) RegisterRedeem(RedeemChainID uint64, ContractChainID uint64, Redeem []byte, CVersion uint64, ContractAddress []byte, Signs [][]byte) (*types.Transaction, error) {
	return _SideChainManager.Contract.RegisterRedeem(&_SideChainManager.TransactOpts, RedeemChainID, ContractChainID, Redeem, CVersion, ContractAddress, Signs)
}

// RegisterRedeem is a paid mutator transaction binding the contract method 0x33e1d41a.
//
// Solidity: function registerRedeem(uint64 RedeemChainID, uint64 ContractChainID, bytes Redeem, uint64 CVersion, bytes ContractAddress, bytes[] Signs) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) RegisterRedeem(RedeemChainID uint64, ContractChainID uint64, Redeem []byte, CVersion uint64, ContractAddress []byte, Signs [][]byte) (*types.Transaction, error) {
	return _SideChainManager.Contract.RegisterRedeem(&_SideChainManager.TransactOpts, RedeemChainID, ContractChainID, Redeem, CVersion, ContractAddress, Signs)
}

// RegisterSideChain is a paid mutator transaction binding the contract method 0xab7a2037.
//
// Solidity: function registerSideChain(address Address, uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait, bytes CCMCAddress, bytes ExtraInfo) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) RegisterSideChain(opts *bind.TransactOpts, Address common.Address, ChainId uint64, Router uint64, Name string, BlocksToWait uint64, CCMCAddress []byte, ExtraInfo []byte) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "registerSideChain", Address, ChainId, Router, Name, BlocksToWait, CCMCAddress, ExtraInfo)
}

// RegisterSideChain is a paid mutator transaction binding the contract method 0xab7a2037.
//
// Solidity: function registerSideChain(address Address, uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait, bytes CCMCAddress, bytes ExtraInfo) returns(bool success)
func (_SideChainManager *SideChainManagerSession) RegisterSideChain(Address common.Address, ChainId uint64, Router uint64, Name string, BlocksToWait uint64, CCMCAddress []byte, ExtraInfo []byte) (*types.Transaction, error) {
	return _SideChainManager.Contract.RegisterSideChain(&_SideChainManager.TransactOpts, Address, ChainId, Router, Name, BlocksToWait, CCMCAddress, ExtraInfo)
}

// RegisterSideChain is a paid mutator transaction binding the contract method 0xab7a2037.
//
// Solidity: function registerSideChain(address Address, uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait, bytes CCMCAddress, bytes ExtraInfo) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) RegisterSideChain(Address common.Address, ChainId uint64, Router uint64, Name string, BlocksToWait uint64, CCMCAddress []byte, ExtraInfo []byte) (*types.Transaction, error) {
	return _SideChainManager.Contract.RegisterSideChain(&_SideChainManager.TransactOpts, Address, ChainId, Router, Name, BlocksToWait, CCMCAddress, ExtraInfo)
}

// SetBtcTxParam is a paid mutator transaction binding the contract method 0xee9891e3.
//
// Solidity: function setBtcTxParam(bytes Redeem, uint64 RedeemChainId, bytes[] Sigs, (uint64,uint64,uint64) Detial) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) SetBtcTxParam(opts *bind.TransactOpts, Redeem []byte, RedeemChainId uint64, Sigs [][]byte, Detial side_chain_managerBtcTxParamDetial) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "setBtcTxParam", Redeem, RedeemChainId, Sigs, Detial)
}

// SetBtcTxParam is a paid mutator transaction binding the contract method 0xee9891e3.
//
// Solidity: function setBtcTxParam(bytes Redeem, uint64 RedeemChainId, bytes[] Sigs, (uint64,uint64,uint64) Detial) returns(bool success)
func (_SideChainManager *SideChainManagerSession) SetBtcTxParam(Redeem []byte, RedeemChainId uint64, Sigs [][]byte, Detial side_chain_managerBtcTxParamDetial) (*types.Transaction, error) {
	return _SideChainManager.Contract.SetBtcTxParam(&_SideChainManager.TransactOpts, Redeem, RedeemChainId, Sigs, Detial)
}

// SetBtcTxParam is a paid mutator transaction binding the contract method 0xee9891e3.
//
// Solidity: function setBtcTxParam(bytes Redeem, uint64 RedeemChainId, bytes[] Sigs, (uint64,uint64,uint64) Detial) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) SetBtcTxParam(Redeem []byte, RedeemChainId uint64, Sigs [][]byte, Detial side_chain_managerBtcTxParamDetial) (*types.Transaction, error) {
	return _SideChainManager.Contract.SetBtcTxParam(&_SideChainManager.TransactOpts, Redeem, RedeemChainId, Sigs, Detial)
}

// UpdateSideChain is a paid mutator transaction binding the contract method 0xf7782f81.
//
// Solidity: function updateSideChain(address Address, uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait, bytes CCMCAddress, bytes ExtraInfo) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) UpdateSideChain(opts *bind.TransactOpts, Address common.Address, ChainId uint64, Router uint64, Name string, BlocksToWait uint64, CCMCAddress []byte, ExtraInfo []byte) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "updateSideChain", Address, ChainId, Router, Name, BlocksToWait, CCMCAddress, ExtraInfo)
}

// UpdateSideChain is a paid mutator transaction binding the contract method 0xf7782f81.
//
// Solidity: function updateSideChain(address Address, uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait, bytes CCMCAddress, bytes ExtraInfo) returns(bool success)
func (_SideChainManager *SideChainManagerSession) UpdateSideChain(Address common.Address, ChainId uint64, Router uint64, Name string, BlocksToWait uint64, CCMCAddress []byte, ExtraInfo []byte) (*types.Transaction, error) {
	return _SideChainManager.Contract.UpdateSideChain(&_SideChainManager.TransactOpts, Address, ChainId, Router, Name, BlocksToWait, CCMCAddress, ExtraInfo)
}

// UpdateSideChain is a paid mutator transaction binding the contract method 0xf7782f81.
//
// Solidity: function updateSideChain(address Address, uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait, bytes CCMCAddress, bytes ExtraInfo) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) UpdateSideChain(Address common.Address, ChainId uint64, Router uint64, Name string, BlocksToWait uint64, CCMCAddress []byte, ExtraInfo []byte) (*types.Transaction, error) {
	return _SideChainManager.Contract.UpdateSideChain(&_SideChainManager.TransactOpts, Address, ChainId, Router, Name, BlocksToWait, CCMCAddress, ExtraInfo)
}

// SideChainManagerApproveQuitSideChainIterator is returned from FilterApproveQuitSideChain and is used to iterate over the raw logs and unpacked data for ApproveQuitSideChain events raised by the SideChainManager contract.
type SideChainManagerApproveQuitSideChainIterator struct {
	Event *SideChainManagerApproveQuitSideChain // Event containing the contract specifics and raw log

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
func (it *SideChainManagerApproveQuitSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerApproveQuitSideChain)
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
		it.Event = new(SideChainManagerApproveQuitSideChain)
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
func (it *SideChainManagerApproveQuitSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerApproveQuitSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerApproveQuitSideChain represents a ApproveQuitSideChain event raised by the SideChainManager contract.
type SideChainManagerApproveQuitSideChain struct {
	ChainId uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproveQuitSideChain is a free log retrieval operation binding the contract event 0x2d8d546abdecb5e9d71f9df93e8bb1e939c865274bf2e59be647e053909634a1.
//
// Solidity: event evtApproveQuitSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) FilterApproveQuitSideChain(opts *bind.FilterOpts) (*SideChainManagerApproveQuitSideChainIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "evtApproveQuitSideChain")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerApproveQuitSideChainIterator{contract: _SideChainManager.contract, event: "evtApproveQuitSideChain", logs: logs, sub: sub}, nil
}

// WatchApproveQuitSideChain is a free log subscription operation binding the contract event 0x2d8d546abdecb5e9d71f9df93e8bb1e939c865274bf2e59be647e053909634a1.
//
// Solidity: event evtApproveQuitSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) WatchApproveQuitSideChain(opts *bind.WatchOpts, sink chan<- *SideChainManagerApproveQuitSideChain) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "evtApproveQuitSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerApproveQuitSideChain)
				if err := _SideChainManager.contract.UnpackLog(event, "evtApproveQuitSideChain", log); err != nil {
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

// ParseApproveQuitSideChain is a log parse operation binding the contract event 0x2d8d546abdecb5e9d71f9df93e8bb1e939c865274bf2e59be647e053909634a1.
//
// Solidity: event evtApproveQuitSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) ParseApproveQuitSideChain(log types.Log) (*SideChainManagerApproveQuitSideChain, error) {
	event := new(SideChainManagerApproveQuitSideChain)
	if err := _SideChainManager.contract.UnpackLog(event, "evtApproveQuitSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SideChainManagerApproveRegisterSideChainIterator is returned from FilterApproveRegisterSideChain and is used to iterate over the raw logs and unpacked data for ApproveRegisterSideChain events raised by the SideChainManager contract.
type SideChainManagerApproveRegisterSideChainIterator struct {
	Event *SideChainManagerApproveRegisterSideChain // Event containing the contract specifics and raw log

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
func (it *SideChainManagerApproveRegisterSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerApproveRegisterSideChain)
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
		it.Event = new(SideChainManagerApproveRegisterSideChain)
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
func (it *SideChainManagerApproveRegisterSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerApproveRegisterSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerApproveRegisterSideChain represents a ApproveRegisterSideChain event raised by the SideChainManager contract.
type SideChainManagerApproveRegisterSideChain struct {
	ChainId uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproveRegisterSideChain is a free log retrieval operation binding the contract event 0x6517ffadca69f75ee51efa5c1e977750e009b25f0bd235ad2afc381ab9704e3e.
//
// Solidity: event evtApproveRegisterSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) FilterApproveRegisterSideChain(opts *bind.FilterOpts) (*SideChainManagerApproveRegisterSideChainIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "evtApproveRegisterSideChain")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerApproveRegisterSideChainIterator{contract: _SideChainManager.contract, event: "evtApproveRegisterSideChain", logs: logs, sub: sub}, nil
}

// WatchApproveRegisterSideChain is a free log subscription operation binding the contract event 0x6517ffadca69f75ee51efa5c1e977750e009b25f0bd235ad2afc381ab9704e3e.
//
// Solidity: event evtApproveRegisterSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) WatchApproveRegisterSideChain(opts *bind.WatchOpts, sink chan<- *SideChainManagerApproveRegisterSideChain) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "evtApproveRegisterSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerApproveRegisterSideChain)
				if err := _SideChainManager.contract.UnpackLog(event, "evtApproveRegisterSideChain", log); err != nil {
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

// ParseApproveRegisterSideChain is a log parse operation binding the contract event 0x6517ffadca69f75ee51efa5c1e977750e009b25f0bd235ad2afc381ab9704e3e.
//
// Solidity: event evtApproveRegisterSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) ParseApproveRegisterSideChain(log types.Log) (*SideChainManagerApproveRegisterSideChain, error) {
	event := new(SideChainManagerApproveRegisterSideChain)
	if err := _SideChainManager.contract.UnpackLog(event, "evtApproveRegisterSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SideChainManagerApproveUpdateSideChainIterator is returned from FilterApproveUpdateSideChain and is used to iterate over the raw logs and unpacked data for ApproveUpdateSideChain events raised by the SideChainManager contract.
type SideChainManagerApproveUpdateSideChainIterator struct {
	Event *SideChainManagerApproveUpdateSideChain // Event containing the contract specifics and raw log

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
func (it *SideChainManagerApproveUpdateSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerApproveUpdateSideChain)
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
		it.Event = new(SideChainManagerApproveUpdateSideChain)
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
func (it *SideChainManagerApproveUpdateSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerApproveUpdateSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerApproveUpdateSideChain represents a ApproveUpdateSideChain event raised by the SideChainManager contract.
type SideChainManagerApproveUpdateSideChain struct {
	ChainId uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproveUpdateSideChain is a free log retrieval operation binding the contract event 0x4c14575dec13d2259a0de2653a8dbbce76780e169975f171b841b9764259252d.
//
// Solidity: event evtApproveUpdateSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) FilterApproveUpdateSideChain(opts *bind.FilterOpts) (*SideChainManagerApproveUpdateSideChainIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "evtApproveUpdateSideChain")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerApproveUpdateSideChainIterator{contract: _SideChainManager.contract, event: "evtApproveUpdateSideChain", logs: logs, sub: sub}, nil
}

// WatchApproveUpdateSideChain is a free log subscription operation binding the contract event 0x4c14575dec13d2259a0de2653a8dbbce76780e169975f171b841b9764259252d.
//
// Solidity: event evtApproveUpdateSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) WatchApproveUpdateSideChain(opts *bind.WatchOpts, sink chan<- *SideChainManagerApproveUpdateSideChain) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "evtApproveUpdateSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerApproveUpdateSideChain)
				if err := _SideChainManager.contract.UnpackLog(event, "evtApproveUpdateSideChain", log); err != nil {
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

// ParseApproveUpdateSideChain is a log parse operation binding the contract event 0x4c14575dec13d2259a0de2653a8dbbce76780e169975f171b841b9764259252d.
//
// Solidity: event evtApproveUpdateSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) ParseApproveUpdateSideChain(log types.Log) (*SideChainManagerApproveUpdateSideChain, error) {
	event := new(SideChainManagerApproveUpdateSideChain)
	if err := _SideChainManager.contract.UnpackLog(event, "evtApproveUpdateSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SideChainManagerQuitSideChainIterator is returned from FilterQuitSideChain and is used to iterate over the raw logs and unpacked data for QuitSideChain events raised by the SideChainManager contract.
type SideChainManagerQuitSideChainIterator struct {
	Event *SideChainManagerQuitSideChain // Event containing the contract specifics and raw log

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
func (it *SideChainManagerQuitSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerQuitSideChain)
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
		it.Event = new(SideChainManagerQuitSideChain)
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
func (it *SideChainManagerQuitSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerQuitSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerQuitSideChain represents a QuitSideChain event raised by the SideChainManager contract.
type SideChainManagerQuitSideChain struct {
	ChainId uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterQuitSideChain is a free log retrieval operation binding the contract event 0x7bf563d05c6f1e91715904f8c9e53945100b40362d910f37e54bdc4989a685c9.
//
// Solidity: event evtQuitSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) FilterQuitSideChain(opts *bind.FilterOpts) (*SideChainManagerQuitSideChainIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "evtQuitSideChain")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerQuitSideChainIterator{contract: _SideChainManager.contract, event: "evtQuitSideChain", logs: logs, sub: sub}, nil
}

// WatchQuitSideChain is a free log subscription operation binding the contract event 0x7bf563d05c6f1e91715904f8c9e53945100b40362d910f37e54bdc4989a685c9.
//
// Solidity: event evtQuitSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) WatchQuitSideChain(opts *bind.WatchOpts, sink chan<- *SideChainManagerQuitSideChain) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "evtQuitSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerQuitSideChain)
				if err := _SideChainManager.contract.UnpackLog(event, "evtQuitSideChain", log); err != nil {
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

// ParseQuitSideChain is a log parse operation binding the contract event 0x7bf563d05c6f1e91715904f8c9e53945100b40362d910f37e54bdc4989a685c9.
//
// Solidity: event evtQuitSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) ParseQuitSideChain(log types.Log) (*SideChainManagerQuitSideChain, error) {
	event := new(SideChainManagerQuitSideChain)
	if err := _SideChainManager.contract.UnpackLog(event, "evtQuitSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SideChainManagerRegisterRedeemIterator is returned from FilterRegisterRedeem and is used to iterate over the raw logs and unpacked data for RegisterRedeem events raised by the SideChainManager contract.
type SideChainManagerRegisterRedeemIterator struct {
	Event *SideChainManagerRegisterRedeem // Event containing the contract specifics and raw log

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
func (it *SideChainManagerRegisterRedeemIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerRegisterRedeem)
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
		it.Event = new(SideChainManagerRegisterRedeem)
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
func (it *SideChainManagerRegisterRedeemIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerRegisterRedeemIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerRegisterRedeem represents a RegisterRedeem event raised by the SideChainManager contract.
type SideChainManagerRegisterRedeem struct {
	Rk              string
	ContractAddress string
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterRegisterRedeem is a free log retrieval operation binding the contract event 0x59fabb884927c5d2e368d355ea11674fd5974066f70ac78f9adab4f1e6a06a36.
//
// Solidity: event evtRegisterRedeem(string rk, string ContractAddress)
func (_SideChainManager *SideChainManagerFilterer) FilterRegisterRedeem(opts *bind.FilterOpts) (*SideChainManagerRegisterRedeemIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "evtRegisterRedeem")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerRegisterRedeemIterator{contract: _SideChainManager.contract, event: "evtRegisterRedeem", logs: logs, sub: sub}, nil
}

// WatchRegisterRedeem is a free log subscription operation binding the contract event 0x59fabb884927c5d2e368d355ea11674fd5974066f70ac78f9adab4f1e6a06a36.
//
// Solidity: event evtRegisterRedeem(string rk, string ContractAddress)
func (_SideChainManager *SideChainManagerFilterer) WatchRegisterRedeem(opts *bind.WatchOpts, sink chan<- *SideChainManagerRegisterRedeem) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "evtRegisterRedeem")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerRegisterRedeem)
				if err := _SideChainManager.contract.UnpackLog(event, "evtRegisterRedeem", log); err != nil {
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

// ParseRegisterRedeem is a log parse operation binding the contract event 0x59fabb884927c5d2e368d355ea11674fd5974066f70ac78f9adab4f1e6a06a36.
//
// Solidity: event evtRegisterRedeem(string rk, string ContractAddress)
func (_SideChainManager *SideChainManagerFilterer) ParseRegisterRedeem(log types.Log) (*SideChainManagerRegisterRedeem, error) {
	event := new(SideChainManagerRegisterRedeem)
	if err := _SideChainManager.contract.UnpackLog(event, "evtRegisterRedeem", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SideChainManagerRegisterSideChainIterator is returned from FilterRegisterSideChain and is used to iterate over the raw logs and unpacked data for RegisterSideChain events raised by the SideChainManager contract.
type SideChainManagerRegisterSideChainIterator struct {
	Event *SideChainManagerRegisterSideChain // Event containing the contract specifics and raw log

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
func (it *SideChainManagerRegisterSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerRegisterSideChain)
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
		it.Event = new(SideChainManagerRegisterSideChain)
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
func (it *SideChainManagerRegisterSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerRegisterSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerRegisterSideChain represents a RegisterSideChain event raised by the SideChainManager contract.
type SideChainManagerRegisterSideChain struct {
	ChainId      uint64
	Router       uint64
	Name         string
	BlocksToWait uint64
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRegisterSideChain is a free log retrieval operation binding the contract event 0x844fd8662935c80c93372c0cb9eb5fd6db52835a3078934d65f762a9bfd13ea2.
//
// Solidity: event evtRegisterSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait)
func (_SideChainManager *SideChainManagerFilterer) FilterRegisterSideChain(opts *bind.FilterOpts) (*SideChainManagerRegisterSideChainIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "evtRegisterSideChain")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerRegisterSideChainIterator{contract: _SideChainManager.contract, event: "evtRegisterSideChain", logs: logs, sub: sub}, nil
}

// WatchRegisterSideChain is a free log subscription operation binding the contract event 0x844fd8662935c80c93372c0cb9eb5fd6db52835a3078934d65f762a9bfd13ea2.
//
// Solidity: event evtRegisterSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait)
func (_SideChainManager *SideChainManagerFilterer) WatchRegisterSideChain(opts *bind.WatchOpts, sink chan<- *SideChainManagerRegisterSideChain) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "evtRegisterSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerRegisterSideChain)
				if err := _SideChainManager.contract.UnpackLog(event, "evtRegisterSideChain", log); err != nil {
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

// ParseRegisterSideChain is a log parse operation binding the contract event 0x844fd8662935c80c93372c0cb9eb5fd6db52835a3078934d65f762a9bfd13ea2.
//
// Solidity: event evtRegisterSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait)
func (_SideChainManager *SideChainManagerFilterer) ParseRegisterSideChain(log types.Log) (*SideChainManagerRegisterSideChain, error) {
	event := new(SideChainManagerRegisterSideChain)
	if err := _SideChainManager.contract.UnpackLog(event, "evtRegisterSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SideChainManagerSetBtcTxParamIterator is returned from FilterSetBtcTxParam and is used to iterate over the raw logs and unpacked data for SetBtcTxParam events raised by the SideChainManager contract.
type SideChainManagerSetBtcTxParamIterator struct {
	Event *SideChainManagerSetBtcTxParam // Event containing the contract specifics and raw log

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
func (it *SideChainManagerSetBtcTxParamIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerSetBtcTxParam)
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
		it.Event = new(SideChainManagerSetBtcTxParam)
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
func (it *SideChainManagerSetBtcTxParamIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerSetBtcTxParamIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerSetBtcTxParam represents a SetBtcTxParam event raised by the SideChainManager contract.
type SideChainManagerSetBtcTxParam struct {
	Rk            string
	RedeemChainId uint64
	FeeRate       uint64
	MinChange     uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSetBtcTxParam is a free log retrieval operation binding the contract event 0x9d437e9e204401edf24ac2844b216c7276e0350ec67ec35417a8ca070ddbde7b.
//
// Solidity: event evtSetBtcTxParam(string rk, uint64 RedeemChainId, uint64 FeeRate, uint64 MinChange)
func (_SideChainManager *SideChainManagerFilterer) FilterSetBtcTxParam(opts *bind.FilterOpts) (*SideChainManagerSetBtcTxParamIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "evtSetBtcTxParam")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerSetBtcTxParamIterator{contract: _SideChainManager.contract, event: "evtSetBtcTxParam", logs: logs, sub: sub}, nil
}

// WatchSetBtcTxParam is a free log subscription operation binding the contract event 0x9d437e9e204401edf24ac2844b216c7276e0350ec67ec35417a8ca070ddbde7b.
//
// Solidity: event evtSetBtcTxParam(string rk, uint64 RedeemChainId, uint64 FeeRate, uint64 MinChange)
func (_SideChainManager *SideChainManagerFilterer) WatchSetBtcTxParam(opts *bind.WatchOpts, sink chan<- *SideChainManagerSetBtcTxParam) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "evtSetBtcTxParam")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerSetBtcTxParam)
				if err := _SideChainManager.contract.UnpackLog(event, "evtSetBtcTxParam", log); err != nil {
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

// ParseSetBtcTxParam is a log parse operation binding the contract event 0x9d437e9e204401edf24ac2844b216c7276e0350ec67ec35417a8ca070ddbde7b.
//
// Solidity: event evtSetBtcTxParam(string rk, uint64 RedeemChainId, uint64 FeeRate, uint64 MinChange)
func (_SideChainManager *SideChainManagerFilterer) ParseSetBtcTxParam(log types.Log) (*SideChainManagerSetBtcTxParam, error) {
	event := new(SideChainManagerSetBtcTxParam)
	if err := _SideChainManager.contract.UnpackLog(event, "evtSetBtcTxParam", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SideChainManagerUpdateSideChainIterator is returned from FilterUpdateSideChain and is used to iterate over the raw logs and unpacked data for UpdateSideChain events raised by the SideChainManager contract.
type SideChainManagerUpdateSideChainIterator struct {
	Event *SideChainManagerUpdateSideChain // Event containing the contract specifics and raw log

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
func (it *SideChainManagerUpdateSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerUpdateSideChain)
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
		it.Event = new(SideChainManagerUpdateSideChain)
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
func (it *SideChainManagerUpdateSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerUpdateSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerUpdateSideChain represents a UpdateSideChain event raised by the SideChainManager contract.
type SideChainManagerUpdateSideChain struct {
	ChainId      uint64
	Router       uint64
	Name         string
	BlocksToWait uint64
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUpdateSideChain is a free log retrieval operation binding the contract event 0xefc07386ca56a4a4f14c5dfb934a955331872da7cc24748a6ca78be8c1741bbe.
//
// Solidity: event evtUpdateSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait)
func (_SideChainManager *SideChainManagerFilterer) FilterUpdateSideChain(opts *bind.FilterOpts) (*SideChainManagerUpdateSideChainIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "evtUpdateSideChain")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerUpdateSideChainIterator{contract: _SideChainManager.contract, event: "evtUpdateSideChain", logs: logs, sub: sub}, nil
}

// WatchUpdateSideChain is a free log subscription operation binding the contract event 0xefc07386ca56a4a4f14c5dfb934a955331872da7cc24748a6ca78be8c1741bbe.
//
// Solidity: event evtUpdateSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait)
func (_SideChainManager *SideChainManagerFilterer) WatchUpdateSideChain(opts *bind.WatchOpts, sink chan<- *SideChainManagerUpdateSideChain) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "evtUpdateSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerUpdateSideChain)
				if err := _SideChainManager.contract.UnpackLog(event, "evtUpdateSideChain", log); err != nil {
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

// ParseUpdateSideChain is a log parse operation binding the contract event 0xefc07386ca56a4a4f14c5dfb934a955331872da7cc24748a6ca78be8c1741bbe.
//
// Solidity: event evtUpdateSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait)
func (_SideChainManager *SideChainManagerFilterer) ParseUpdateSideChain(log types.Log) (*SideChainManagerUpdateSideChain, error) {
	event := new(SideChainManagerUpdateSideChain)
	if err := _SideChainManager.contract.UnpackLog(event, "evtUpdateSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
