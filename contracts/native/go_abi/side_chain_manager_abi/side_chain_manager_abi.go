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

// SideChainManagerABI is the input ABI used to generate the binding from.
const SideChainManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"}],\"name\":\"EventApproveQuitSideChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"}],\"name\":\"EventApproveRegisterSideChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"}],\"name\":\"EventApproveUpdateSideChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"}],\"name\":\"EventQuitSideChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"rk\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ContractAddress\",\"type\":\"string\"}],\"name\":\"EventRegisterRedeem\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"Router\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"BlocksToWait\",\"type\":\"uint64\"}],\"name\":\"EventRegisterSideChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"rk\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"RedeemChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"FeeRate\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"MinChange\",\"type\":\"uint64\"}],\"name\":\"EventSetBtcTxParam\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"Router\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"BlocksToWait\",\"type\":\"uint64\"}],\"name\":\"EventUpdateSideChain\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"Chainid\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodApproveQuitSideChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"Chainid\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodApproveRegisterSideChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"Chainid\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodApproveUpdateSideChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MethodContractName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"Chainid\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"name\":\"MethodQuitSideChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"RedeemChainID\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"ContractChainID\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"Redeem\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"CVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"ContractAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"Signs\",\"type\":\"bytes[]\"}],\"name\":\"MethodRegisterRedeem\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"Router\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"BlocksToWait\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"CCMCAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"ExtraInfo\",\"type\":\"bytes\"}],\"name\":\"MethodRegisterSideChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"Redeem\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"RedeemChainId\",\"type\":\"uint64\"},{\"internalType\":\"bytes[]\",\"name\":\"Sigs\",\"type\":\"bytes[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"PVersion\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"FeeRate\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"MinChange\",\"type\":\"uint64\"}],\"internalType\":\"structside_chain_manager.BtcTxParamDetial\",\"name\":\"Detial\",\"type\":\"tuple\"}],\"name\":\"MethodSetBtcTxParam\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"ChainId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"Router\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"BlocksToWait\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"CCMCAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"ExtraInfo\",\"type\":\"bytes\"}],\"name\":\"MethodUpdateSideChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SideChainManagerFuncSigs maps the 4-byte function signature to its string representation.
var SideChainManagerFuncSigs = map[string]string{
	"2a98a777": "MethodApproveQuitSideChain(uint64,address)",
	"1bd71458": "MethodApproveRegisterSideChain(uint64,address)",
	"a854e742": "MethodApproveUpdateSideChain(uint64,address)",
	"e50f8f44": "MethodContractName()",
	"114668b5": "MethodQuitSideChain(uint64,address)",
	"7c048283": "MethodRegisterRedeem(uint64,uint64,bytes,uint64,bytes,bytes[])",
	"ab940b32": "MethodRegisterSideChain(address,uint64,uint64,string,uint64,bytes,bytes)",
	"4af71c48": "MethodSetBtcTxParam(bytes,uint64,bytes[],(uint64,uint64,uint64))",
	"a1cb0b51": "MethodUpdateSideChain(address,uint64,uint64,string,uint64,bytes,bytes)",
}

// SideChainManagerBin is the compiled bytecode used for deploying new contracts.
var SideChainManagerBin = "0x608060405234801561001057600080fd5b50610604806100206000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c80637c048283116100665780637c048283146100db578063a1cb0b51146100f5578063a854e74214610098578063ab940b32146100f5578063e50f8f441461011057600080fd5b8063114668b5146100985780631bd71458146100985780632a98a777146100985780634af71c48146100c3575b600080fd5b6100ae6100a636600461041d565b600092915050565b60405190151581526020015b60405180910390f35b6100ae6100d1366004610358565b6000949350505050565b6100ae6100e9366004610450565b60009695505050505050565b6100ae610103366004610276565b6000979650505050505050565b60606040516100ba9190610509565b600067ffffffffffffffff831115610139576101396105b8565b61014c601f8401601f1916602001610587565b905082815283838301111561016057600080fd5b828260208301376000602084830101529392505050565b80356001600160a01b038116811461018e57600080fd5b919050565b600082601f8301126101a457600080fd5b8135602067ffffffffffffffff808311156101c1576101c16105b8565b8260051b6101d0838201610587565b8481528381019087850183890186018a10156101eb57600080fd5b60009350835b8781101561022857813586811115610207578586fd5b6102158c89838e0101610237565b85525092860192908601906001016101f1565b50909998505050505050505050565b600082601f83011261024857600080fd5b6102578383356020850161011f565b9392505050565b803567ffffffffffffffff8116811461018e57600080fd5b600080600080600080600060e0888a03121561029157600080fd5b61029a88610177565b96506102a86020890161025e565b95506102b66040890161025e565b9450606088013567ffffffffffffffff808211156102d357600080fd5b818a0191508a601f8301126102e757600080fd5b6102f68b83356020850161011f565b955061030460808b0161025e565b945060a08a013591508082111561031a57600080fd5b6103268b838c01610237565b935060c08a013591508082111561033c57600080fd5b506103498a828b01610237565b91505092959891949750929550565b60008060008084860360c081121561036f57600080fd5b853567ffffffffffffffff8082111561038757600080fd5b61039389838a01610237565b96506103a16020890161025e565b955060408801359150808211156103b757600080fd5b506103c488828901610193565b9350506060605f19820112156103d957600080fd5b506103e261055e565b6103ee6060870161025e565b81526103fc6080870161025e565b602082015261040d60a0870161025e565b6040820152939692955090935050565b6000806040838503121561043057600080fd5b6104398361025e565b915061044760208401610177565b90509250929050565b60008060008060008060c0878903121561046957600080fd5b6104728761025e565b95506104806020880161025e565b9450604087013567ffffffffffffffff8082111561049d57600080fd5b6104a98a838b01610237565b95506104b760608a0161025e565b945060808901359150808211156104cd57600080fd5b6104d98a838b01610237565b935060a08901359150808211156104ef57600080fd5b506104fc89828a01610193565b9150509295509295509295565b600060208083528351808285015260005b818110156105365785810183015185820160400152820161051a565b81811115610548576000604083870101525b50601f01601f1916929092016040019392505050565b6040516060810167ffffffffffffffff81118282101715610581576105816105b8565b60405290565b604051601f8201601f1916810167ffffffffffffffff811182821017156105b0576105b06105b8565b604052919050565b634e487b7160e01b600052604160045260246000fdfea264697066735822122058680fac4d8a506b028d5141d02e8fb84c923b57e8d1bc21843ceda3d5da2f5564736f6c63430008060033"

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

// MethodApproveQuitSideChain is a paid mutator transaction binding the contract method 0x2a98a777.
//
// Solidity: function MethodApproveQuitSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) MethodApproveQuitSideChain(opts *bind.TransactOpts, Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "MethodApproveQuitSideChain", Chainid, Address)
}

// MethodApproveQuitSideChain is a paid mutator transaction binding the contract method 0x2a98a777.
//
// Solidity: function MethodApproveQuitSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerSession) MethodApproveQuitSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodApproveQuitSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// MethodApproveQuitSideChain is a paid mutator transaction binding the contract method 0x2a98a777.
//
// Solidity: function MethodApproveQuitSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) MethodApproveQuitSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodApproveQuitSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// MethodApproveRegisterSideChain is a paid mutator transaction binding the contract method 0x1bd71458.
//
// Solidity: function MethodApproveRegisterSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) MethodApproveRegisterSideChain(opts *bind.TransactOpts, Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "MethodApproveRegisterSideChain", Chainid, Address)
}

// MethodApproveRegisterSideChain is a paid mutator transaction binding the contract method 0x1bd71458.
//
// Solidity: function MethodApproveRegisterSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerSession) MethodApproveRegisterSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodApproveRegisterSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// MethodApproveRegisterSideChain is a paid mutator transaction binding the contract method 0x1bd71458.
//
// Solidity: function MethodApproveRegisterSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) MethodApproveRegisterSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodApproveRegisterSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// MethodApproveUpdateSideChain is a paid mutator transaction binding the contract method 0xa854e742.
//
// Solidity: function MethodApproveUpdateSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) MethodApproveUpdateSideChain(opts *bind.TransactOpts, Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "MethodApproveUpdateSideChain", Chainid, Address)
}

// MethodApproveUpdateSideChain is a paid mutator transaction binding the contract method 0xa854e742.
//
// Solidity: function MethodApproveUpdateSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerSession) MethodApproveUpdateSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodApproveUpdateSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// MethodApproveUpdateSideChain is a paid mutator transaction binding the contract method 0xa854e742.
//
// Solidity: function MethodApproveUpdateSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) MethodApproveUpdateSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodApproveUpdateSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// MethodContractName is a paid mutator transaction binding the contract method 0xe50f8f44.
//
// Solidity: function MethodContractName() returns(string Name)
func (_SideChainManager *SideChainManagerTransactor) MethodContractName(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "MethodContractName")
}

// MethodContractName is a paid mutator transaction binding the contract method 0xe50f8f44.
//
// Solidity: function MethodContractName() returns(string Name)
func (_SideChainManager *SideChainManagerSession) MethodContractName() (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodContractName(&_SideChainManager.TransactOpts)
}

// MethodContractName is a paid mutator transaction binding the contract method 0xe50f8f44.
//
// Solidity: function MethodContractName() returns(string Name)
func (_SideChainManager *SideChainManagerTransactorSession) MethodContractName() (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodContractName(&_SideChainManager.TransactOpts)
}

// MethodQuitSideChain is a paid mutator transaction binding the contract method 0x114668b5.
//
// Solidity: function MethodQuitSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) MethodQuitSideChain(opts *bind.TransactOpts, Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "MethodQuitSideChain", Chainid, Address)
}

// MethodQuitSideChain is a paid mutator transaction binding the contract method 0x114668b5.
//
// Solidity: function MethodQuitSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerSession) MethodQuitSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodQuitSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// MethodQuitSideChain is a paid mutator transaction binding the contract method 0x114668b5.
//
// Solidity: function MethodQuitSideChain(uint64 Chainid, address Address) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) MethodQuitSideChain(Chainid uint64, Address common.Address) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodQuitSideChain(&_SideChainManager.TransactOpts, Chainid, Address)
}

// MethodRegisterRedeem is a paid mutator transaction binding the contract method 0x7c048283.
//
// Solidity: function MethodRegisterRedeem(uint64 RedeemChainID, uint64 ContractChainID, bytes Redeem, uint64 CVersion, bytes ContractAddress, bytes[] Signs) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) MethodRegisterRedeem(opts *bind.TransactOpts, RedeemChainID uint64, ContractChainID uint64, Redeem []byte, CVersion uint64, ContractAddress []byte, Signs [][]byte) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "MethodRegisterRedeem", RedeemChainID, ContractChainID, Redeem, CVersion, ContractAddress, Signs)
}

// MethodRegisterRedeem is a paid mutator transaction binding the contract method 0x7c048283.
//
// Solidity: function MethodRegisterRedeem(uint64 RedeemChainID, uint64 ContractChainID, bytes Redeem, uint64 CVersion, bytes ContractAddress, bytes[] Signs) returns(bool success)
func (_SideChainManager *SideChainManagerSession) MethodRegisterRedeem(RedeemChainID uint64, ContractChainID uint64, Redeem []byte, CVersion uint64, ContractAddress []byte, Signs [][]byte) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodRegisterRedeem(&_SideChainManager.TransactOpts, RedeemChainID, ContractChainID, Redeem, CVersion, ContractAddress, Signs)
}

// MethodRegisterRedeem is a paid mutator transaction binding the contract method 0x7c048283.
//
// Solidity: function MethodRegisterRedeem(uint64 RedeemChainID, uint64 ContractChainID, bytes Redeem, uint64 CVersion, bytes ContractAddress, bytes[] Signs) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) MethodRegisterRedeem(RedeemChainID uint64, ContractChainID uint64, Redeem []byte, CVersion uint64, ContractAddress []byte, Signs [][]byte) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodRegisterRedeem(&_SideChainManager.TransactOpts, RedeemChainID, ContractChainID, Redeem, CVersion, ContractAddress, Signs)
}

// MethodRegisterSideChain is a paid mutator transaction binding the contract method 0xab940b32.
//
// Solidity: function MethodRegisterSideChain(address Address, uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait, bytes CCMCAddress, bytes ExtraInfo) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) MethodRegisterSideChain(opts *bind.TransactOpts, Address common.Address, ChainId uint64, Router uint64, Name string, BlocksToWait uint64, CCMCAddress []byte, ExtraInfo []byte) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "MethodRegisterSideChain", Address, ChainId, Router, Name, BlocksToWait, CCMCAddress, ExtraInfo)
}

// MethodRegisterSideChain is a paid mutator transaction binding the contract method 0xab940b32.
//
// Solidity: function MethodRegisterSideChain(address Address, uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait, bytes CCMCAddress, bytes ExtraInfo) returns(bool success)
func (_SideChainManager *SideChainManagerSession) MethodRegisterSideChain(Address common.Address, ChainId uint64, Router uint64, Name string, BlocksToWait uint64, CCMCAddress []byte, ExtraInfo []byte) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodRegisterSideChain(&_SideChainManager.TransactOpts, Address, ChainId, Router, Name, BlocksToWait, CCMCAddress, ExtraInfo)
}

// MethodRegisterSideChain is a paid mutator transaction binding the contract method 0xab940b32.
//
// Solidity: function MethodRegisterSideChain(address Address, uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait, bytes CCMCAddress, bytes ExtraInfo) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) MethodRegisterSideChain(Address common.Address, ChainId uint64, Router uint64, Name string, BlocksToWait uint64, CCMCAddress []byte, ExtraInfo []byte) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodRegisterSideChain(&_SideChainManager.TransactOpts, Address, ChainId, Router, Name, BlocksToWait, CCMCAddress, ExtraInfo)
}

// MethodSetBtcTxParam is a paid mutator transaction binding the contract method 0x4af71c48.
//
// Solidity: function MethodSetBtcTxParam(bytes Redeem, uint64 RedeemChainId, bytes[] Sigs, (uint64,uint64,uint64) Detial) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) MethodSetBtcTxParam(opts *bind.TransactOpts, Redeem []byte, RedeemChainId uint64, Sigs [][]byte, Detial side_chain_managerBtcTxParamDetial) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "MethodSetBtcTxParam", Redeem, RedeemChainId, Sigs, Detial)
}

// MethodSetBtcTxParam is a paid mutator transaction binding the contract method 0x4af71c48.
//
// Solidity: function MethodSetBtcTxParam(bytes Redeem, uint64 RedeemChainId, bytes[] Sigs, (uint64,uint64,uint64) Detial) returns(bool success)
func (_SideChainManager *SideChainManagerSession) MethodSetBtcTxParam(Redeem []byte, RedeemChainId uint64, Sigs [][]byte, Detial side_chain_managerBtcTxParamDetial) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodSetBtcTxParam(&_SideChainManager.TransactOpts, Redeem, RedeemChainId, Sigs, Detial)
}

// MethodSetBtcTxParam is a paid mutator transaction binding the contract method 0x4af71c48.
//
// Solidity: function MethodSetBtcTxParam(bytes Redeem, uint64 RedeemChainId, bytes[] Sigs, (uint64,uint64,uint64) Detial) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) MethodSetBtcTxParam(Redeem []byte, RedeemChainId uint64, Sigs [][]byte, Detial side_chain_managerBtcTxParamDetial) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodSetBtcTxParam(&_SideChainManager.TransactOpts, Redeem, RedeemChainId, Sigs, Detial)
}

// MethodUpdateSideChain is a paid mutator transaction binding the contract method 0xa1cb0b51.
//
// Solidity: function MethodUpdateSideChain(address Address, uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait, bytes CCMCAddress, bytes ExtraInfo) returns(bool success)
func (_SideChainManager *SideChainManagerTransactor) MethodUpdateSideChain(opts *bind.TransactOpts, Address common.Address, ChainId uint64, Router uint64, Name string, BlocksToWait uint64, CCMCAddress []byte, ExtraInfo []byte) (*types.Transaction, error) {
	return _SideChainManager.contract.Transact(opts, "MethodUpdateSideChain", Address, ChainId, Router, Name, BlocksToWait, CCMCAddress, ExtraInfo)
}

// MethodUpdateSideChain is a paid mutator transaction binding the contract method 0xa1cb0b51.
//
// Solidity: function MethodUpdateSideChain(address Address, uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait, bytes CCMCAddress, bytes ExtraInfo) returns(bool success)
func (_SideChainManager *SideChainManagerSession) MethodUpdateSideChain(Address common.Address, ChainId uint64, Router uint64, Name string, BlocksToWait uint64, CCMCAddress []byte, ExtraInfo []byte) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodUpdateSideChain(&_SideChainManager.TransactOpts, Address, ChainId, Router, Name, BlocksToWait, CCMCAddress, ExtraInfo)
}

// MethodUpdateSideChain is a paid mutator transaction binding the contract method 0xa1cb0b51.
//
// Solidity: function MethodUpdateSideChain(address Address, uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait, bytes CCMCAddress, bytes ExtraInfo) returns(bool success)
func (_SideChainManager *SideChainManagerTransactorSession) MethodUpdateSideChain(Address common.Address, ChainId uint64, Router uint64, Name string, BlocksToWait uint64, CCMCAddress []byte, ExtraInfo []byte) (*types.Transaction, error) {
	return _SideChainManager.Contract.MethodUpdateSideChain(&_SideChainManager.TransactOpts, Address, ChainId, Router, Name, BlocksToWait, CCMCAddress, ExtraInfo)
}

// SideChainManagerEventApproveQuitSideChainIterator is returned from FilterEventApproveQuitSideChain and is used to iterate over the raw logs and unpacked data for EventApproveQuitSideChain events raised by the SideChainManager contract.
type SideChainManagerEventApproveQuitSideChainIterator struct {
	Event *SideChainManagerEventApproveQuitSideChain // Event containing the contract specifics and raw log

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
func (it *SideChainManagerEventApproveQuitSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerEventApproveQuitSideChain)
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
		it.Event = new(SideChainManagerEventApproveQuitSideChain)
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
func (it *SideChainManagerEventApproveQuitSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerEventApproveQuitSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerEventApproveQuitSideChain represents a EventApproveQuitSideChain event raised by the SideChainManager contract.
type SideChainManagerEventApproveQuitSideChain struct {
	ChainId uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterEventApproveQuitSideChain is a free log retrieval operation binding the contract event 0x5e034eb5779566da890ead43c51b13e0b17d02fec5566eecdcf8a20b5ff76a7a.
//
// Solidity: event EventApproveQuitSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) FilterEventApproveQuitSideChain(opts *bind.FilterOpts) (*SideChainManagerEventApproveQuitSideChainIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "EventApproveQuitSideChain")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerEventApproveQuitSideChainIterator{contract: _SideChainManager.contract, event: "EventApproveQuitSideChain", logs: logs, sub: sub}, nil
}

// WatchEventApproveQuitSideChain is a free log subscription operation binding the contract event 0x5e034eb5779566da890ead43c51b13e0b17d02fec5566eecdcf8a20b5ff76a7a.
//
// Solidity: event EventApproveQuitSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) WatchEventApproveQuitSideChain(opts *bind.WatchOpts, sink chan<- *SideChainManagerEventApproveQuitSideChain) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "EventApproveQuitSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerEventApproveQuitSideChain)
				if err := _SideChainManager.contract.UnpackLog(event, "EventApproveQuitSideChain", log); err != nil {
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

// ParseEventApproveQuitSideChain is a log parse operation binding the contract event 0x5e034eb5779566da890ead43c51b13e0b17d02fec5566eecdcf8a20b5ff76a7a.
//
// Solidity: event EventApproveQuitSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) ParseEventApproveQuitSideChain(log types.Log) (*SideChainManagerEventApproveQuitSideChain, error) {
	event := new(SideChainManagerEventApproveQuitSideChain)
	if err := _SideChainManager.contract.UnpackLog(event, "EventApproveQuitSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SideChainManagerEventApproveRegisterSideChainIterator is returned from FilterEventApproveRegisterSideChain and is used to iterate over the raw logs and unpacked data for EventApproveRegisterSideChain events raised by the SideChainManager contract.
type SideChainManagerEventApproveRegisterSideChainIterator struct {
	Event *SideChainManagerEventApproveRegisterSideChain // Event containing the contract specifics and raw log

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
func (it *SideChainManagerEventApproveRegisterSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerEventApproveRegisterSideChain)
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
		it.Event = new(SideChainManagerEventApproveRegisterSideChain)
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
func (it *SideChainManagerEventApproveRegisterSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerEventApproveRegisterSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerEventApproveRegisterSideChain represents a EventApproveRegisterSideChain event raised by the SideChainManager contract.
type SideChainManagerEventApproveRegisterSideChain struct {
	ChainId uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterEventApproveRegisterSideChain is a free log retrieval operation binding the contract event 0x89989b5103147110383d568e9119509f41ea48a41b74d08165786093794a014a.
//
// Solidity: event EventApproveRegisterSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) FilterEventApproveRegisterSideChain(opts *bind.FilterOpts) (*SideChainManagerEventApproveRegisterSideChainIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "EventApproveRegisterSideChain")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerEventApproveRegisterSideChainIterator{contract: _SideChainManager.contract, event: "EventApproveRegisterSideChain", logs: logs, sub: sub}, nil
}

// WatchEventApproveRegisterSideChain is a free log subscription operation binding the contract event 0x89989b5103147110383d568e9119509f41ea48a41b74d08165786093794a014a.
//
// Solidity: event EventApproveRegisterSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) WatchEventApproveRegisterSideChain(opts *bind.WatchOpts, sink chan<- *SideChainManagerEventApproveRegisterSideChain) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "EventApproveRegisterSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerEventApproveRegisterSideChain)
				if err := _SideChainManager.contract.UnpackLog(event, "EventApproveRegisterSideChain", log); err != nil {
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

// ParseEventApproveRegisterSideChain is a log parse operation binding the contract event 0x89989b5103147110383d568e9119509f41ea48a41b74d08165786093794a014a.
//
// Solidity: event EventApproveRegisterSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) ParseEventApproveRegisterSideChain(log types.Log) (*SideChainManagerEventApproveRegisterSideChain, error) {
	event := new(SideChainManagerEventApproveRegisterSideChain)
	if err := _SideChainManager.contract.UnpackLog(event, "EventApproveRegisterSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SideChainManagerEventApproveUpdateSideChainIterator is returned from FilterEventApproveUpdateSideChain and is used to iterate over the raw logs and unpacked data for EventApproveUpdateSideChain events raised by the SideChainManager contract.
type SideChainManagerEventApproveUpdateSideChainIterator struct {
	Event *SideChainManagerEventApproveUpdateSideChain // Event containing the contract specifics and raw log

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
func (it *SideChainManagerEventApproveUpdateSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerEventApproveUpdateSideChain)
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
		it.Event = new(SideChainManagerEventApproveUpdateSideChain)
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
func (it *SideChainManagerEventApproveUpdateSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerEventApproveUpdateSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerEventApproveUpdateSideChain represents a EventApproveUpdateSideChain event raised by the SideChainManager contract.
type SideChainManagerEventApproveUpdateSideChain struct {
	ChainId uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterEventApproveUpdateSideChain is a free log retrieval operation binding the contract event 0x502f871ba66d02e976863b3a83616be336ad58f13668bd0517c4fd2834d6e562.
//
// Solidity: event EventApproveUpdateSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) FilterEventApproveUpdateSideChain(opts *bind.FilterOpts) (*SideChainManagerEventApproveUpdateSideChainIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "EventApproveUpdateSideChain")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerEventApproveUpdateSideChainIterator{contract: _SideChainManager.contract, event: "EventApproveUpdateSideChain", logs: logs, sub: sub}, nil
}

// WatchEventApproveUpdateSideChain is a free log subscription operation binding the contract event 0x502f871ba66d02e976863b3a83616be336ad58f13668bd0517c4fd2834d6e562.
//
// Solidity: event EventApproveUpdateSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) WatchEventApproveUpdateSideChain(opts *bind.WatchOpts, sink chan<- *SideChainManagerEventApproveUpdateSideChain) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "EventApproveUpdateSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerEventApproveUpdateSideChain)
				if err := _SideChainManager.contract.UnpackLog(event, "EventApproveUpdateSideChain", log); err != nil {
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

// ParseEventApproveUpdateSideChain is a log parse operation binding the contract event 0x502f871ba66d02e976863b3a83616be336ad58f13668bd0517c4fd2834d6e562.
//
// Solidity: event EventApproveUpdateSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) ParseEventApproveUpdateSideChain(log types.Log) (*SideChainManagerEventApproveUpdateSideChain, error) {
	event := new(SideChainManagerEventApproveUpdateSideChain)
	if err := _SideChainManager.contract.UnpackLog(event, "EventApproveUpdateSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SideChainManagerEventQuitSideChainIterator is returned from FilterEventQuitSideChain and is used to iterate over the raw logs and unpacked data for EventQuitSideChain events raised by the SideChainManager contract.
type SideChainManagerEventQuitSideChainIterator struct {
	Event *SideChainManagerEventQuitSideChain // Event containing the contract specifics and raw log

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
func (it *SideChainManagerEventQuitSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerEventQuitSideChain)
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
		it.Event = new(SideChainManagerEventQuitSideChain)
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
func (it *SideChainManagerEventQuitSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerEventQuitSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerEventQuitSideChain represents a EventQuitSideChain event raised by the SideChainManager contract.
type SideChainManagerEventQuitSideChain struct {
	ChainId uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterEventQuitSideChain is a free log retrieval operation binding the contract event 0xd4d9b94d0ae94466eaf372660bfb91e759657517e4ee621a0a8c7555068429e1.
//
// Solidity: event EventQuitSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) FilterEventQuitSideChain(opts *bind.FilterOpts) (*SideChainManagerEventQuitSideChainIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "EventQuitSideChain")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerEventQuitSideChainIterator{contract: _SideChainManager.contract, event: "EventQuitSideChain", logs: logs, sub: sub}, nil
}

// WatchEventQuitSideChain is a free log subscription operation binding the contract event 0xd4d9b94d0ae94466eaf372660bfb91e759657517e4ee621a0a8c7555068429e1.
//
// Solidity: event EventQuitSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) WatchEventQuitSideChain(opts *bind.WatchOpts, sink chan<- *SideChainManagerEventQuitSideChain) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "EventQuitSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerEventQuitSideChain)
				if err := _SideChainManager.contract.UnpackLog(event, "EventQuitSideChain", log); err != nil {
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

// ParseEventQuitSideChain is a log parse operation binding the contract event 0xd4d9b94d0ae94466eaf372660bfb91e759657517e4ee621a0a8c7555068429e1.
//
// Solidity: event EventQuitSideChain(uint64 ChainId)
func (_SideChainManager *SideChainManagerFilterer) ParseEventQuitSideChain(log types.Log) (*SideChainManagerEventQuitSideChain, error) {
	event := new(SideChainManagerEventQuitSideChain)
	if err := _SideChainManager.contract.UnpackLog(event, "EventQuitSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SideChainManagerEventRegisterRedeemIterator is returned from FilterEventRegisterRedeem and is used to iterate over the raw logs and unpacked data for EventRegisterRedeem events raised by the SideChainManager contract.
type SideChainManagerEventRegisterRedeemIterator struct {
	Event *SideChainManagerEventRegisterRedeem // Event containing the contract specifics and raw log

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
func (it *SideChainManagerEventRegisterRedeemIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerEventRegisterRedeem)
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
		it.Event = new(SideChainManagerEventRegisterRedeem)
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
func (it *SideChainManagerEventRegisterRedeemIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerEventRegisterRedeemIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerEventRegisterRedeem represents a EventRegisterRedeem event raised by the SideChainManager contract.
type SideChainManagerEventRegisterRedeem struct {
	Rk              string
	ContractAddress string
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterEventRegisterRedeem is a free log retrieval operation binding the contract event 0x3c747cbdd6d26e0b3f91408a29fb04cc9fdba2ac0664269b5dacc50e102cfb29.
//
// Solidity: event EventRegisterRedeem(string rk, string ContractAddress)
func (_SideChainManager *SideChainManagerFilterer) FilterEventRegisterRedeem(opts *bind.FilterOpts) (*SideChainManagerEventRegisterRedeemIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "EventRegisterRedeem")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerEventRegisterRedeemIterator{contract: _SideChainManager.contract, event: "EventRegisterRedeem", logs: logs, sub: sub}, nil
}

// WatchEventRegisterRedeem is a free log subscription operation binding the contract event 0x3c747cbdd6d26e0b3f91408a29fb04cc9fdba2ac0664269b5dacc50e102cfb29.
//
// Solidity: event EventRegisterRedeem(string rk, string ContractAddress)
func (_SideChainManager *SideChainManagerFilterer) WatchEventRegisterRedeem(opts *bind.WatchOpts, sink chan<- *SideChainManagerEventRegisterRedeem) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "EventRegisterRedeem")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerEventRegisterRedeem)
				if err := _SideChainManager.contract.UnpackLog(event, "EventRegisterRedeem", log); err != nil {
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

// ParseEventRegisterRedeem is a log parse operation binding the contract event 0x3c747cbdd6d26e0b3f91408a29fb04cc9fdba2ac0664269b5dacc50e102cfb29.
//
// Solidity: event EventRegisterRedeem(string rk, string ContractAddress)
func (_SideChainManager *SideChainManagerFilterer) ParseEventRegisterRedeem(log types.Log) (*SideChainManagerEventRegisterRedeem, error) {
	event := new(SideChainManagerEventRegisterRedeem)
	if err := _SideChainManager.contract.UnpackLog(event, "EventRegisterRedeem", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SideChainManagerEventRegisterSideChainIterator is returned from FilterEventRegisterSideChain and is used to iterate over the raw logs and unpacked data for EventRegisterSideChain events raised by the SideChainManager contract.
type SideChainManagerEventRegisterSideChainIterator struct {
	Event *SideChainManagerEventRegisterSideChain // Event containing the contract specifics and raw log

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
func (it *SideChainManagerEventRegisterSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerEventRegisterSideChain)
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
		it.Event = new(SideChainManagerEventRegisterSideChain)
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
func (it *SideChainManagerEventRegisterSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerEventRegisterSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerEventRegisterSideChain represents a EventRegisterSideChain event raised by the SideChainManager contract.
type SideChainManagerEventRegisterSideChain struct {
	ChainId      uint64
	Router       uint64
	Name         string
	BlocksToWait uint64
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterEventRegisterSideChain is a free log retrieval operation binding the contract event 0xc9b771b2fe54416c7d44e132629c68bca376cf90f6896d8ddeb2660bf8658f8d.
//
// Solidity: event EventRegisterSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait)
func (_SideChainManager *SideChainManagerFilterer) FilterEventRegisterSideChain(opts *bind.FilterOpts) (*SideChainManagerEventRegisterSideChainIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "EventRegisterSideChain")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerEventRegisterSideChainIterator{contract: _SideChainManager.contract, event: "EventRegisterSideChain", logs: logs, sub: sub}, nil
}

// WatchEventRegisterSideChain is a free log subscription operation binding the contract event 0xc9b771b2fe54416c7d44e132629c68bca376cf90f6896d8ddeb2660bf8658f8d.
//
// Solidity: event EventRegisterSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait)
func (_SideChainManager *SideChainManagerFilterer) WatchEventRegisterSideChain(opts *bind.WatchOpts, sink chan<- *SideChainManagerEventRegisterSideChain) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "EventRegisterSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerEventRegisterSideChain)
				if err := _SideChainManager.contract.UnpackLog(event, "EventRegisterSideChain", log); err != nil {
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

// ParseEventRegisterSideChain is a log parse operation binding the contract event 0xc9b771b2fe54416c7d44e132629c68bca376cf90f6896d8ddeb2660bf8658f8d.
//
// Solidity: event EventRegisterSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait)
func (_SideChainManager *SideChainManagerFilterer) ParseEventRegisterSideChain(log types.Log) (*SideChainManagerEventRegisterSideChain, error) {
	event := new(SideChainManagerEventRegisterSideChain)
	if err := _SideChainManager.contract.UnpackLog(event, "EventRegisterSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SideChainManagerEventSetBtcTxParamIterator is returned from FilterEventSetBtcTxParam and is used to iterate over the raw logs and unpacked data for EventSetBtcTxParam events raised by the SideChainManager contract.
type SideChainManagerEventSetBtcTxParamIterator struct {
	Event *SideChainManagerEventSetBtcTxParam // Event containing the contract specifics and raw log

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
func (it *SideChainManagerEventSetBtcTxParamIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerEventSetBtcTxParam)
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
		it.Event = new(SideChainManagerEventSetBtcTxParam)
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
func (it *SideChainManagerEventSetBtcTxParamIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerEventSetBtcTxParamIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerEventSetBtcTxParam represents a EventSetBtcTxParam event raised by the SideChainManager contract.
type SideChainManagerEventSetBtcTxParam struct {
	Rk            string
	RedeemChainId uint64
	FeeRate       uint64
	MinChange     uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterEventSetBtcTxParam is a free log retrieval operation binding the contract event 0x87ebb80202803f960c6ab08a7825dece987a38bb0894c8e0b49a013d8d0eb44c.
//
// Solidity: event EventSetBtcTxParam(string rk, uint64 RedeemChainId, uint64 FeeRate, uint64 MinChange)
func (_SideChainManager *SideChainManagerFilterer) FilterEventSetBtcTxParam(opts *bind.FilterOpts) (*SideChainManagerEventSetBtcTxParamIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "EventSetBtcTxParam")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerEventSetBtcTxParamIterator{contract: _SideChainManager.contract, event: "EventSetBtcTxParam", logs: logs, sub: sub}, nil
}

// WatchEventSetBtcTxParam is a free log subscription operation binding the contract event 0x87ebb80202803f960c6ab08a7825dece987a38bb0894c8e0b49a013d8d0eb44c.
//
// Solidity: event EventSetBtcTxParam(string rk, uint64 RedeemChainId, uint64 FeeRate, uint64 MinChange)
func (_SideChainManager *SideChainManagerFilterer) WatchEventSetBtcTxParam(opts *bind.WatchOpts, sink chan<- *SideChainManagerEventSetBtcTxParam) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "EventSetBtcTxParam")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerEventSetBtcTxParam)
				if err := _SideChainManager.contract.UnpackLog(event, "EventSetBtcTxParam", log); err != nil {
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

// ParseEventSetBtcTxParam is a log parse operation binding the contract event 0x87ebb80202803f960c6ab08a7825dece987a38bb0894c8e0b49a013d8d0eb44c.
//
// Solidity: event EventSetBtcTxParam(string rk, uint64 RedeemChainId, uint64 FeeRate, uint64 MinChange)
func (_SideChainManager *SideChainManagerFilterer) ParseEventSetBtcTxParam(log types.Log) (*SideChainManagerEventSetBtcTxParam, error) {
	event := new(SideChainManagerEventSetBtcTxParam)
	if err := _SideChainManager.contract.UnpackLog(event, "EventSetBtcTxParam", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SideChainManagerEventUpdateSideChainIterator is returned from FilterEventUpdateSideChain and is used to iterate over the raw logs and unpacked data for EventUpdateSideChain events raised by the SideChainManager contract.
type SideChainManagerEventUpdateSideChainIterator struct {
	Event *SideChainManagerEventUpdateSideChain // Event containing the contract specifics and raw log

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
func (it *SideChainManagerEventUpdateSideChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SideChainManagerEventUpdateSideChain)
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
		it.Event = new(SideChainManagerEventUpdateSideChain)
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
func (it *SideChainManagerEventUpdateSideChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SideChainManagerEventUpdateSideChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SideChainManagerEventUpdateSideChain represents a EventUpdateSideChain event raised by the SideChainManager contract.
type SideChainManagerEventUpdateSideChain struct {
	ChainId      uint64
	Router       uint64
	Name         string
	BlocksToWait uint64
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterEventUpdateSideChain is a free log retrieval operation binding the contract event 0x95139edda7c5321bac4a5c0743e6080c8b40fd0d666c16c0657cf6051ecb645d.
//
// Solidity: event EventUpdateSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait)
func (_SideChainManager *SideChainManagerFilterer) FilterEventUpdateSideChain(opts *bind.FilterOpts) (*SideChainManagerEventUpdateSideChainIterator, error) {

	logs, sub, err := _SideChainManager.contract.FilterLogs(opts, "EventUpdateSideChain")
	if err != nil {
		return nil, err
	}
	return &SideChainManagerEventUpdateSideChainIterator{contract: _SideChainManager.contract, event: "EventUpdateSideChain", logs: logs, sub: sub}, nil
}

// WatchEventUpdateSideChain is a free log subscription operation binding the contract event 0x95139edda7c5321bac4a5c0743e6080c8b40fd0d666c16c0657cf6051ecb645d.
//
// Solidity: event EventUpdateSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait)
func (_SideChainManager *SideChainManagerFilterer) WatchEventUpdateSideChain(opts *bind.WatchOpts, sink chan<- *SideChainManagerEventUpdateSideChain) (event.Subscription, error) {

	logs, sub, err := _SideChainManager.contract.WatchLogs(opts, "EventUpdateSideChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SideChainManagerEventUpdateSideChain)
				if err := _SideChainManager.contract.UnpackLog(event, "EventUpdateSideChain", log); err != nil {
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

// ParseEventUpdateSideChain is a log parse operation binding the contract event 0x95139edda7c5321bac4a5c0743e6080c8b40fd0d666c16c0657cf6051ecb645d.
//
// Solidity: event EventUpdateSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait)
func (_SideChainManager *SideChainManagerFilterer) ParseEventUpdateSideChain(log types.Log) (*SideChainManagerEventUpdateSideChain, error) {
	event := new(SideChainManagerEventUpdateSideChain)
	if err := _SideChainManager.contract.UnpackLog(event, "EventUpdateSideChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
